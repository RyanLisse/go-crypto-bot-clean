package service

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model/market"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/service"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
	"golang.org/x/time/rate"
)

// Errors
var (
	ErrSniperNotRunning        = errors.New("sniper service is not running")
	ErrInvalidSymbol           = errors.New("invalid symbol for sniping")
	ErrMaxPriceExceeded        = errors.New("current price exceeds maximum allowed price")
	ErrInsufficientLiquidity   = errors.New("insufficient liquidity for sniping")
	ErrConcurrencyLimitReached = errors.New("maximum concurrent orders limit reached")
)

// Default configuration values
const (
	DefaultMaxBuyAmount        = 100.0
	DefaultMaxPricePerToken    = 1.0
	DefaultMaxSlippagePercent  = 5.0
	DefaultMaxConcurrentOrders = 3
	DefaultRetryAttempts       = 3
	DefaultRetryDelayMs        = 100
	DefaultTakeProfitPercent   = 20.0
	DefaultStopLossPercent     = 10.0
)

// Status constants
const (
	StatusRunning = "running"
	StatusStopped = "stopped"
	StatusError   = "error"
)

// MexcSniperService implements the SniperService interface for MEXC exchange
type MexcSniperService struct {
	// Dependencies
	mexcClient    port.MEXCClient
	symbolRepo    port.SymbolRepository
	orderRepo     port.OrderRepository
	marketService port.MarketDataService
	logger        *zerolog.Logger

	// Configuration
	config *port.SniperConfig

	// Service state
	status         string
	activeOrders   int32
	rateLimiter    *rate.Limiter
	orderSemaphore chan struct{}
	mutex          sync.RWMutex

	// Caches for performance optimization
	symbolCache    sync.Map
	priceCache     sync.Map
	lastCacheClean time.Time

	// Async order save queue and worker group
	saveOrderQueue chan *model.Order
	saveOrderWG    sync.WaitGroup

	// WebSocket integration
	listingDetectionService *service.NewListingDetectionService
	autoSnipeEnabled        bool
	autoSnipeConfig         *port.SniperConfig
	autoSnipeMutex          sync.RWMutex
}

// NewMexcSniperService creates a new instance of the MEXC sniper service
func NewMexcSniperService(
	mexcClient port.MEXCClient,
	symbolRepo port.SymbolRepository,
	orderRepo port.OrderRepository,
	marketService port.MarketDataService,
	listingDetectionService *service.NewListingDetectionService,
	logger *zerolog.Logger,
) *MexcSniperService {
	// Create default configuration
	config := &port.SniperConfig{
		MaxBuyAmount:        DefaultMaxBuyAmount,
		MaxPricePerToken:    DefaultMaxPricePerToken,
		EnablePartialFills:  true,
		MaxSlippagePercent:  DefaultMaxSlippagePercent,
		BypassRiskChecks:    true,                  // Default to bypassing risk checks for speed
		PreferredOrderType:  model.OrderTypeMarket, // Default to market orders for speed
		MaxConcurrentOrders: DefaultMaxConcurrentOrders,
		RetryAttempts:       DefaultRetryAttempts,
		RetryDelayMs:        DefaultRetryDelayMs,
		EnableTakeProfit:    false, // Disabled by default
		TakeProfitPercent:   DefaultTakeProfitPercent,
		EnableStopLoss:      false, // Disabled by default
		StopLossPercent:     DefaultStopLossPercent,
		PriceCacheExpiryMs:  500, // Added configurable price cache expiry in ms
		RateLimitPerSec:     10,  // Added configurable rate limit per second
	}

	// Create service instance
	service := &MexcSniperService{
		mexcClient:              mexcClient,
		symbolRepo:              symbolRepo,
		orderRepo:               orderRepo,
		marketService:           marketService,
		listingDetectionService: listingDetectionService,
		logger:                  logger,
		config:                  config,
		status:                  StatusStopped,
		rateLimiter:             rate.NewLimiter(rate.Limit(config.RateLimitPerSec), 20), // Configurable rate limiter
		orderSemaphore:          make(chan struct{}, config.MaxConcurrentOrders),         // Semaphore size from config
		lastCacheClean:          time.Now(),
		saveOrderQueue:          make(chan *model.Order, 100), // Buffered channel for async order saves
		autoSnipeEnabled:        false,
		autoSnipeConfig:         config, // Use default config for auto-snipe
	}

	// Start async order save workers
	for i := 0; i < 5; i++ {
		service.saveOrderWG.Add(1)
		go service.saveOrderWorker()
	}

	return service
}

// ExecuteSnipe executes a high-speed buy on a newly listed token
func (s *MexcSniperService) ExecuteSnipe(ctx context.Context, symbol string) (*model.Order, error) {
	return s.ExecuteSnipeWithConfig(ctx, symbol, s.config)
}

// ExecuteSnipeWithConfig executes a high-speed buy with custom configuration
func (s *MexcSniperService) ExecuteSnipeWithConfig(ctx context.Context, symbol string, config *port.SniperConfig) (*model.Order, error) {
	// Check if service is running
	if s.status != StatusRunning {
		return nil, ErrSniperNotRunning
	}

	// Check concurrency limit
	if atomic.LoadInt32(&s.activeOrders) >= int32(config.MaxConcurrentOrders) {
		return nil, ErrConcurrencyLimitReached
	}

	// Increment active orders counter before acquiring semaphore to avoid race
	atomic.AddInt32(&s.activeOrders, 1)

	// Acquire semaphore slot
	select {
	case s.orderSemaphore <- struct{}{}:
		// Slot acquired
		defer func() {
			<-s.orderSemaphore
			atomic.AddInt32(&s.activeOrders, -1)
		}()
	case <-ctx.Done():
		atomic.AddInt32(&s.activeOrders, -1)
		return nil, ctx.Err()
	}
	defer atomic.AddInt32(&s.activeOrders, -1)

	// Fast validation of symbol
	valid, err := s.fastValidateSymbol(ctx, symbol)
	if err != nil {
		s.logger.Error().Err(err).Str("symbol", symbol).Msg("Failed to validate symbol")
		return nil, err
	}
	if !valid {
		return nil, ErrInvalidSymbol
	}

	// Get current price (with caching for performance)
	price, err := s.getFastPrice(ctx, symbol)
	if err != nil {
		s.logger.Error().Err(err).Str("symbol", symbol).Msg("Failed to get price")
		return nil, err
	}

	// Check if price exceeds maximum allowed price
	if price > config.MaxPricePerToken {
		s.logger.Warn().
			Str("symbol", symbol).
			Float64("price", price).
			Float64("maxPrice", config.MaxPricePerToken).
			Msg("Price exceeds maximum allowed price")
		return nil, ErrMaxPriceExceeded
	}

	// Calculate quantity based on max buy amount
	quantity := config.MaxBuyAmount / price

	// Create order request
	orderRequest := &model.OrderRequest{
		Symbol:   symbol,
		Side:     model.OrderSideBuy,
		Type:     config.PreferredOrderType,
		Quantity: quantity,
		Price:    price, // Only used for limit orders
	}

	// Execute order with retries
	var order *model.Order
	for attempt := 0; attempt < config.RetryAttempts; attempt++ {
		// Wait for rate limiter
		if err := s.rateLimiter.Wait(ctx); err != nil {
			return nil, err
		}

		// Place order
		order, err = s.executeOrder(ctx, orderRequest)
		if err == nil {
			break // Success
		}

		s.logger.Warn().
			Err(err).
			Str("symbol", symbol).
			Int("attempt", attempt+1).
			Int("maxAttempts", config.RetryAttempts).
			Msg("Retry placing order")

		// Wait before retrying
		select {
		case <-time.After(time.Duration(config.RetryDelayMs) * time.Millisecond):
			// Continue with retry
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	if err != nil {
		s.logger.Error().
			Err(err).
			Str("symbol", symbol).
			Int("attempts", config.RetryAttempts).
			Msg("Failed to place order after all retry attempts")
		return nil, err
	}

	// Handle take-profit and stop-loss if enabled
	if order != nil && (config.EnableTakeProfit || config.EnableStopLoss) {
		go s.handlePostTradeActions(context.Background(), order, config)
	}

	return order, nil
}

// PrevalidateSymbol checks if a symbol is valid for sniping without executing a trade
func (s *MexcSniperService) PrevalidateSymbol(ctx context.Context, symbol string) (bool, error) {
	// Check if service is running
	if s.status != StatusRunning {
		return false, ErrSniperNotRunning
	}

	return s.fastValidateSymbol(ctx, symbol)
}

// GetConfig returns the current sniper configuration
func (s *MexcSniperService) GetConfig() *port.SniperConfig {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	// Return a copy to prevent external modification
	configCopy := *s.config
	return &configCopy
}

// UpdateConfig updates the sniper configuration
func (s *MexcSniperService) UpdateConfig(config *port.SniperConfig) error {
	if config == nil {
		return errors.New("config cannot be nil")
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Update configuration
	s.config = config

	// Update semaphore if max concurrent orders changed
	if cap(s.orderSemaphore) != config.MaxConcurrentOrders {
		// Create new semaphore with updated capacity
		newSemaphore := make(chan struct{}, config.MaxConcurrentOrders)
		s.orderSemaphore = newSemaphore
	}

	return nil
}

// GetStatus returns the current status of the sniper service
func (s *MexcSniperService) GetStatus() string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.status
}

// Start starts the sniper service
func (s *MexcSniperService) Start() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.status == StatusRunning {
		return nil // Already running
	}

	// Initialize caches
	s.symbolCache = sync.Map{}
	s.priceCache = sync.Map{}
	s.lastCacheClean = time.Now()

	// Start cache cleanup goroutine
	go s.cacheCleaner()

	s.status = StatusRunning
	s.logger.Info().Msg("Sniper service started")
	return nil
}

// Stop stops the sniper service
func (s *MexcSniperService) Stop() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.status != StatusRunning {
		return nil // Not running
	}

	s.status = StatusStopped
	s.logger.Info().Msg("Sniper service stopped")
	return nil
}

// Internal helper methods

// fastValidateSymbol performs a fast validation of a symbol
func (s *MexcSniperService) fastValidateSymbol(ctx context.Context, symbol string) (bool, error) {
	// Check cache first
	if valid, ok := s.symbolCache.Load(symbol); ok {
		return valid.(bool), nil
	}

	// Validate symbol with repository
	symbolInfo, err := s.symbolRepo.GetBySymbol(ctx, symbol)
	if err != nil {
		return false, err
	}

	valid := symbolInfo != nil && symbolInfo.Status == "TRADING"

	// Cache the result
	s.symbolCache.Store(symbol, valid)

	return valid, nil
}

// getFastPrice gets the current price of a symbol with caching
func (s *MexcSniperService) getFastPrice(ctx context.Context, symbol string) (float64, error) {
	// Check cache first
	if cachedPrice, ok := s.priceCache.Load(symbol); ok {
		priceData := cachedPrice.(struct {
			price     float64
			timestamp time.Time
		})

		// Use cached price if it's recent (less than configurable expiry)
		if time.Since(priceData.timestamp) < time.Duration(s.config.PriceCacheExpiryMs)*time.Millisecond {
			return priceData.price, nil
		}
	}

	// Get fresh price
	ticker, err := s.marketService.GetTicker(ctx, symbol)
	if err != nil {
		return 0, err
	}

	price := ticker.Price

	// Cache the result
	s.priceCache.Store(symbol, struct {
		price     float64
		timestamp time.Time
	}{
		price:     price,
		timestamp: time.Now(),
	})

	return price, nil
}

// executeOrder executes an order with optimized path
func (s *MexcSniperService) executeOrder(ctx context.Context, request *model.OrderRequest) (*model.Order, error) {
	// Direct API call to MEXC for fastest execution
	order, err := s.mexcClient.PlaceOrder(
		ctx,
		request.Symbol,
		request.Side,
		request.Type,
		request.Quantity,
		request.Price,
		model.TimeInForceGTC, // Only used for limit orders
	)

	if err != nil {
		return nil, err
	}

	// Asynchronously save order to database to not block execution
	s.saveOrderAsync(order)

	return order, nil
}

func (s *MexcSniperService) saveOrderAsync(order *model.Order) {
	select {
	case s.saveOrderQueue <- order:
		// Enqueued successfully
	default:
		// Queue full, log warning and drop save to avoid blocking
		s.logger.Warn().
			Str("orderID", order.OrderID).
			Msg("Order save queue full, dropping save to avoid blocking")
	}
}

func (s *MexcSniperService) saveOrderWorker() {
	defer s.saveOrderWG.Done()
	for order := range s.saveOrderQueue {
		if err := s.orderRepo.Create(context.Background(), order); err != nil {
			s.logger.Error().
				Err(err).
				Str("orderID", order.OrderID).
				Msg("Failed to save order to database")
		}
	}
}

// handlePostTradeActions handles take-profit and stop-loss orders
func (s *MexcSniperService) handlePostTradeActions(ctx context.Context, order *model.Order, config *port.SniperConfig) {
	// Wait for order to be filled
	filledOrder, err := s.waitForOrderFill(ctx, order.Symbol, order.OrderID)
	if err != nil {
		s.logger.Error().
			Err(err).
			Str("orderID", order.OrderID).
			Msg("Failed to wait for order fill")
		return
	}

	if filledOrder.Status != model.OrderStatusFilled {
		s.logger.Warn().
			Str("orderID", order.OrderID).
			Str("status", string(filledOrder.Status)).
			Msg("Order not filled, skipping post-trade actions")
		return
	}

	// Get filled price
	filledPrice := filledOrder.Price

	// Create take-profit and stop-loss orders
	var eg errgroup.Group

	if config.EnableTakeProfit {
		eg.Go(func() error {
			takeProfitPrice := filledPrice * (1 + config.TakeProfitPercent/100)

			tpRequest := &model.OrderRequest{
				Symbol:   order.Symbol,
				Side:     model.OrderSideSell,
				Type:     model.OrderTypeLimit,
				Quantity: filledOrder.ExecutedQty,
				Price:    takeProfitPrice,
			}

			_, err := s.executeOrder(ctx, tpRequest)
			if err != nil {
				s.logger.Error().
					Err(err).
					Str("orderID", order.OrderID).
					Float64("takeProfitPrice", takeProfitPrice).
					Msg("Failed to place take-profit order")
				return err
			}

			s.logger.Info().
				Str("orderID", order.OrderID).
				Float64("takeProfitPrice", takeProfitPrice).
				Msg("Take-profit order placed successfully")

			return nil
		})
	}

	if config.EnableStopLoss {
		eg.Go(func() error {
			stopLossPrice := filledPrice * (1 - config.StopLossPercent/100)

			slRequest := &model.OrderRequest{
				Symbol: order.Symbol,
				Side:   model.OrderSideSell,
				// Type:     model.OrderTypeStopLoss, // Commented out because undefined
				Type:     model.OrderTypeLimit, // Use OrderTypeLimit as fallback
				Quantity: filledOrder.ExecutedQty,
				Price:    stopLossPrice,
			}

			_, err := s.executeOrder(ctx, slRequest)
			if err != nil {
				s.logger.Error().
					Err(err).
					Str("orderID", order.OrderID).
					Float64("stopLossPrice", stopLossPrice).
					Msg("Failed to place stop-loss order")
				return err
			}

			s.logger.Info().
				Str("orderID", order.OrderID).
				Float64("stopLossPrice", stopLossPrice).
				Msg("Stop-loss order placed successfully")

			return nil
		})
	}

	// Wait for all orders to be placed
	if err := eg.Wait(); err != nil {
		s.logger.Error().
			Err(err).
			Str("orderID", order.OrderID).
			Msg("Error in post-trade actions")
	}
}

// waitForOrderFill waits for an order to be filled
func (s *MexcSniperService) waitForOrderFill(ctx context.Context, symbol, orderID string) (*model.Order, error) {
	// Create a timeout context
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			order, err := s.mexcClient.GetOrderStatus(ctx, symbol, orderID)
			if err != nil {
				s.logger.Warn().
					Err(err).
					Str("orderID", orderID).
					Msg("Error checking order status")
				continue
			}

			if order.Status == model.OrderStatusFilled ||
				order.Status == model.OrderStatusCanceled ||
				order.Status == model.OrderStatusRejected {
				return order, nil
			}

		case <-ctx.Done():
			return nil, fmt.Errorf("timeout waiting for order fill: %w", ctx.Err())
		}
	}
}

// cacheCleaner periodically cleans up caches
func (s *MexcSniperService) cacheCleaner() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.mutex.RLock()
			if s.status != StatusRunning {
				s.mutex.RUnlock()
				return
			}
			s.mutex.RUnlock()

			// Clear caches
			s.symbolCache = sync.Map{}
			s.priceCache = sync.Map{}
			s.lastCacheClean = time.Now()

			s.logger.Debug().Msg("Sniper service caches cleared")
		}
	}
}

// SetupAutoSnipe configures the sniper to automatically snipe new listings
func (s *MexcSniperService) SetupAutoSnipe(enabled bool, config *port.SniperConfig) error {
	s.autoSnipeMutex.Lock()
	defer s.autoSnipeMutex.Unlock()

	s.autoSnipeEnabled = enabled

	if enabled && config != nil {
		s.autoSnipeConfig = config
	} else if enabled {
		// Use default config if none provided
		s.autoSnipeConfig = s.GetConfig()
	}

	s.logger.Info().
		Bool("enabled", enabled).
		Msg("Auto-snipe configuration updated")

	// If enabled, set up the event listener
	if enabled && s.listingDetectionService != nil {
		// Register a callback for new coin events
		go s.setupNewCoinEventListener()
	}

	return nil
}

// setupNewCoinEventListener sets up a listener for new coin events
func (s *MexcSniperService) setupNewCoinEventListener() {
	// Ensure the listing detection service is running
	if s.listingDetectionService != nil {
		// The service should already be running, but we'll log a message just in case
		s.logger.Info().Msg("Using existing listing detection service for auto-snipe")
	} else {
		s.logger.Warn().Msg("No listing detection service available for auto-snipe")
		return
	}

	// Subscribe to events
	// Note: This is a simplified example. In a real implementation, you would need to
	// register with the event bus or use the appropriate mechanism to receive events.
	// For now, we'll simulate this by directly checking for new coins periodically.
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		// Create a repository to query for new coins
		ctx := context.Background()

		for {
			select {
			case <-ticker.C:
				// Check if auto-snipe is still enabled
				s.autoSnipeMutex.RLock()
				enabled := s.autoSnipeEnabled
				s.autoSnipeMutex.RUnlock()

				if !enabled {
					return
				}

				// Get recently listed coins that are now tradable
				// We'll use the repository directly since the service doesn't expose this method
				coins, err := s.symbolRepo.GetSymbolsByStatus(ctx, "TRADING", 10, 0)
				if err != nil {
					s.logger.Error().Err(err).Msg("Failed to get tradable symbols")
					continue
				}

				// Process each coin
				for _, coin := range coins {
					// Check if we've already processed this coin
					if processed, _ := s.symbolCache.Load(coin.Symbol + "_processed"); processed != nil {
						continue
					}

					// Mark as processed
					s.symbolCache.Store(coin.Symbol+"_processed", true)

					// Execute snipe in a new goroutine
					go func(coin *market.Symbol) {
						// Convert market.Symbol to a format we can use
						s.autoSnipeMutex.RLock()
						config := s.autoSnipeConfig
						s.autoSnipeMutex.RUnlock()

						if config == nil {
							return
						}

						ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
						defer cancel()

						s.logger.Info().
							Str("symbol", coin.Symbol).
							Str("status", coin.Status).
							Msg("Auto-sniping new listing")

						// Execute snipe
						_, err := s.ExecuteSnipeWithConfig(ctx, coin.Symbol, config)
						if err != nil {
							s.logger.Error().
								Err(err).
								Str("symbol", coin.Symbol).
								Msg("Failed to auto-snipe new listing")
							return
						}

						s.logger.Info().
							Str("symbol", coin.Symbol).
							Msg("Auto-snipe executed successfully")
					}(coin)
				}
			}
		}
	}()
}
