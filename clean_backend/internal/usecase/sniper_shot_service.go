package usecase

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/port"
	"github.com/rs/zerolog"
)

// PriceChecker is an interface for getting the current price of a symbol
type PriceChecker interface {
	// GetCurrentPrice returns the current price for a given symbol
	GetCurrentPrice(ctx context.Context, symbol string) (float64, error)
}

// marketDataPriceChecker adapts a MarketDataService to the PriceChecker interface
type marketDataPriceChecker struct {
	marketDataService port.MarketDataService
}

// NewMarketDataPriceChecker creates a new PriceChecker from a MarketDataService
func NewMarketDataPriceChecker(marketDataService port.MarketDataService) PriceChecker {
	return &marketDataPriceChecker{
		marketDataService: marketDataService,
	}
}

// GetCurrentPrice implements the PriceChecker interface
func (m *marketDataPriceChecker) GetCurrentPrice(ctx context.Context, symbol string) (float64, error) {
	ticker, err := m.marketDataService.GetTicker(ctx, symbol)
	if err != nil {
		return 0, fmt.Errorf("failed to get ticker for %s: %w", symbol, err)
	}

	// Return the last price as the current price
	return ticker.LastPrice, nil
}

// SimplePriceChecker is a function type that gets the current price
type SimplePriceChecker func(ctx context.Context, symbol string) (float64, error)

// GetCurrentPrice calls the function
func (f SimplePriceChecker) GetCurrentPrice(ctx context.Context, symbol string) (float64, error) {
	return f(ctx, symbol)
}

// WithFixedPrice returns a PriceChecker that always returns the specified price
func WithFixedPrice(price float64) PriceChecker {
	return SimplePriceChecker(func(ctx context.Context, symbol string) (float64, error) {
		return price, nil
	})
}

// SniperShotService is responsible for executing rapid trades when market conditions meet criteria
type SniperShotService struct {
	logger         *zerolog.Logger
	tradeExecutor  port.TradeExecutor
	orderRepo      port.OrderRepository
	priceChecker   PriceChecker
	triggerService port.TriggerService
	mutex          sync.Mutex
}

// SniperShotConfig contains configuration for the sniper shot service
type SniperShotConfig struct {
	MaxRetries      int
	RetryDelay      time.Duration
	ExecutionExpiry time.Duration
}

// NewSniperShotService creates a new instance of the SniperShot service
func NewSniperShotService(
	logger *zerolog.Logger,
	tradeExecutor port.TradeExecutor,
	orderRepo port.OrderRepository,
	priceChecker port.PriceChecker,
	triggerService port.TriggerService,
) *SniperShotService {
	return &SniperShotService{
		logger:         logger,
		tradeExecutor:  tradeExecutor,
		orderRepo:      orderRepo,
		priceChecker:   priceChecker,
		triggerService: triggerService,
	}
}

// NewSniperShotServiceWithMarketData creates a new SniperShotService using a MarketDataService
func NewSniperShotServiceWithMarketData(
	logger *zerolog.Logger,
	tradeExecutor port.TradeExecutor,
	orderRepo port.OrderRepository,
	marketDataService port.MarketDataService,
) *SniperShotService {
	priceChecker := NewMarketDataPriceChecker(marketDataService)
	return NewSniperShotService(logger, tradeExecutor, orderRepo, priceChecker, nil)
}

// ExecuteSniper executes a sniper shot trade according to the provided request.
// It implements the method from port.SniperShotServicePort.
func (s *SniperShotService) ExecuteSniper(ctx context.Context, req *model.SniperShotRequest) (*model.SniperShotResult, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	startTime := time.Now()
	logger := s.logger.With().
		Str("user_id", req.UserID).
		Str("symbol", req.Symbol).
		Str("side", string(req.Side)).
		Str("type", string(req.Type)).
		Float64("quantity", req.Quantity).
		Logger()

	logger.Info().Msg("SniperShot execution initiated")

	// Create an execution context with timeout
	execCtx, cancel := context.WithTimeout(ctx, req.TimeLimit)
	defer cancel()

	// Create order details
	orderReq := &model.OrderRequest{
		UserID:      req.UserID,
		Symbol:      req.Symbol,
		Side:        req.Side,
		Type:        req.Type,
		Quantity:    req.Quantity,
		Price:       req.Price,
		TimeInForce: model.TimeInForceGTC, // Good-Till-Cancel by default
	}

	// If there's a trigger condition, wait for it to be met before executing
	if req.Condition != nil {
		// Use model.TriggerCondition
		condition := req.Condition
		logger.Info().
			Float64("target_price", condition.TargetPrice).
			Str("operator", condition.Operator).             // Keep string operator for logging flexibility if needed
			Str("comparison", string(condition.Comparison)). // Log enum comparison type
			Int("timeout_secs", condition.MaxTimeoutSecs).
			Float64("buffer_pct", condition.PriceBufferPct).
			Msg("Waiting for trigger condition")

		// Set up timeout channel if needed
		var timeoutChan <-chan time.Time
		if condition.MaxTimeoutSecs > 0 {
			timeoutChan = time.After(time.Duration(condition.MaxTimeoutSecs) * time.Second)
		}

		// Default check interval if not specified
		checkInterval := 500 * time.Millisecond
		if condition.CheckIntervalMs > 0 {
			checkInterval = time.Duration(condition.CheckIntervalMs) * time.Millisecond
		}

		ticker := time.NewTicker(checkInterval)
		defer ticker.Stop()

		// Wait for the condition to be met or timeout
		conditionMet := false
		var currentPrice float64
		for !conditionMet {
			select {
			case <-ctx.Done():
				logger.Warn().Msg("Context cancelled while waiting for trigger condition")
				return &model.SniperShotResult{ // Use model type
					Success:   false,
					Error:     ctx.Err(),
					Timestamp: time.Now(),
					Latency:   time.Since(startTime),
				}, ctx.Err()
			case <-timeoutChan:
				err := fmt.Errorf("timeout waiting for trigger condition after %d seconds", condition.MaxTimeoutSecs)
				logger.Warn().Msg(err.Error())
				return &model.SniperShotResult{ // Use model type
					Success:   false,
					Error:     err,
					Timestamp: time.Now(),
					Latency:   time.Since(startTime),
				}, err
			case <-ticker.C:
				// Get current price using the injected priceChecker
				var err error
				currentPrice, err = s.priceChecker.GetCurrentPrice(ctx, req.Symbol) // Use injected checker
				if err != nil {
					logger.Error().Err(err).Msg("Failed to get current price")
					continue
				}

				// Check if condition is met using the injected TriggerService
				conditionMet = s.triggerService.CheckCondition(condition, currentPrice) // Use injected service
				if conditionMet {
					logger.Info().
						Float64("current_price", currentPrice).
						Float64("target_price", condition.TargetPrice).
						Msg("Trigger condition met")

					// Execute any callbacks if configured
					if condition.Callbacks != nil {
						for _, callback := range condition.Callbacks {
							if callback != nil {
								// Execute callbacks in a goroutine to prevent blocking
								go func(cb func(price float64), p float64) {
									cb(p)
								}(callback, currentPrice)
							}
						}
					}
				}
			}
		}

		// Apply price buffer if specified
		if condition.PriceBufferPct > 0 {
			if req.Type == model.OrderTypeLimit {
				// Apply buffer differently based on order side
				if req.Side == model.OrderSideBuy {
					// For buy orders, we can pay a bit more (add buffer)
					bufferAmount := condition.TargetPrice * condition.PriceBufferPct
					orderReq.Price = condition.TargetPrice + bufferAmount
					logger.Info().
						Float64("original_price", condition.TargetPrice).
						Float64("buffered_price", orderReq.Price).
						Msg("Applied price buffer for buy order")
				} else {
					// For sell orders, we can accept a bit less (subtract buffer)
					bufferAmount := condition.TargetPrice * condition.PriceBufferPct
					orderReq.Price = condition.TargetPrice - bufferAmount
					logger.Info().
						Float64("original_price", condition.TargetPrice).
						Float64("buffered_price", orderReq.Price).
						Msg("Applied price buffer for sell order")
				}
			}
		}
	}

	// Execute the order
	logger.Info().Msg("Executing order now")
	orderResponse, err := s.tradeExecutor.ExecuteOrder(execCtx, orderReq)
	if err != nil {
		logger.Error().
			Err(err).
			Str("symbol", req.Symbol).
			Msg("Failed to execute sniper shot")

		return &model.SniperShotResult{ // Use model type
			Success:   false,
			Error:     err,
			Timestamp: time.Now(),
			Latency:   time.Since(startTime),
		}, err
	}

	// Create a new Order struct pointer from the embedded order in the response
	order := &model.Order{
		ID:              orderResponse.ID,
		OrderID:         orderResponse.OrderID,
		ClientOrderID:   orderResponse.ClientOrderID,
		UserID:          orderResponse.UserID,
		Symbol:          orderResponse.Symbol,
		Side:            orderResponse.Side,
		Type:            orderResponse.Type,
		Status:          orderResponse.Status,
		Price:           orderResponse.Price,
		Quantity:        orderResponse.Quantity,
		ExecutedQty:     orderResponse.ExecutedQty,
		AvgFillPrice:    orderResponse.AvgFillPrice,
		Commission:      orderResponse.Commission,
		CommissionAsset: orderResponse.CommissionAsset,
		TimeInForce:     orderResponse.TimeInForce,
		CreatedAt:       orderResponse.CreatedAt,
		UpdatedAt:       orderResponse.UpdatedAt,
		Exchange:        orderResponse.Exchange,
	}

	// Log successful execution
	logger.Info().
		Str("order_id", order.OrderID).
		Str("status", string(order.Status)).
		Float64("executed_qty", order.ExecutedQty).
		Float64("avg_price", order.AvgFillPrice).
		Msg("SniperShot execution completed successfully")

	// Save order to repository
	if err := s.orderRepo.Create(ctx, order); err != nil {
		logger.Warn().
			Err(err).
			Str("order_id", order.OrderID).
			Msg("Failed to save order to repository, but trade was executed")
	}

	return &model.SniperShotResult{ // Use model type
		Success:   true,
		Order:     order,
		Timestamp: time.Now(),
		Latency:   time.Since(startTime),
	}, nil
}

// CancelSniper attempts to cancel a previously executed sniper shot order
func (s *SniperShotService) CancelSniper(ctx context.Context, symbol, orderID string) error {
	s.logger.Info().
		Str("orderID", orderID).
		Str("symbol", symbol).
		Msg("Attempting to cancel sniper shot order")

	err := s.tradeExecutor.CancelOrderWithRetry(ctx, symbol, orderID)
	if err != nil {
		s.logger.Error().
			Err(err).
			Str("orderID", orderID).
			Str("symbol", symbol).
			Msg("Failed to cancel sniper shot order")
		return fmt.Errorf("failed to cancel order: %w", err)
	}

	s.logger.Info().
		Str("orderID", orderID).
		Str("symbol", symbol).
		Msg("Successfully cancelled sniper shot order")

	return nil
}

// GetSniperOrderStatus retrieves the current status of a sniper shot order
func (s *SniperShotService) GetSniperOrderStatus(ctx context.Context, symbol, orderID string) (*model.Order, error) {
	s.logger.Debug().
		Str("orderID", orderID).
		Str("symbol", symbol).
		Msg("Retrieving sniper shot order status")

	order, err := s.tradeExecutor.GetOrderStatusWithRetry(ctx, symbol, orderID)
	if err != nil {
		s.logger.Error().
			Err(err).
			Str("orderID", orderID).
			Str("symbol", symbol).
			Msg("Failed to retrieve sniper shot order status")
		return nil, fmt.Errorf("failed to get order status: %w", err)
	}

	return order, nil
}
