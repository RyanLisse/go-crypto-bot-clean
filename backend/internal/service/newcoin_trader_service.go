package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/notification"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model/status"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/usecase"
	"github.com/rs/zerolog"
)

// NewCoinTraderService watches for new coins with StatusTrading status and executes market orders
type NewCoinTraderService struct {
	statusUC         usecase.StatusUseCase
	tradeExecutor    port.TradeExecutor
	tradeHistory     port.TradeHistoryRepository
	telegramNotifier *notification.TelegramNotifier
	config           *NewCoinTraderConfig
	logger           *zerolog.Logger
	ctx              context.Context
	cancel           context.CancelFunc
	isRunning        bool
}

// NewCoinTraderConfig contains configuration for the new coin trader service
type NewCoinTraderConfig struct {
	DefaultTradeAmount    float64       // Default amount to trade in quote currency (e.g., USDT)
	MaxTradeAmount        float64       // Maximum amount to trade per coin
	RetryInterval         time.Duration // Interval between retries for "not yet tradable" errors
	MaxRetries            int           // Maximum number of retries for "not yet tradable" errors
	EnableTelegram        bool          // Whether to send Telegram notifications
	EnableTradeHistory    bool          // Whether to save trade records
	BlacklistedSymbols    []string      // Symbols that should not be traded
	TradeDelayAfterStatus time.Duration // Delay between status change and trade execution
}

// DefaultNewCoinTraderConfig returns default configuration
func DefaultNewCoinTraderConfig() *NewCoinTraderConfig {
	return &NewCoinTraderConfig{
		DefaultTradeAmount:    10.0,                   // 10 USDT by default
		MaxTradeAmount:        100.0,                  // Maximum 100 USDT per coin
		RetryInterval:         5 * time.Second,        // Retry every 5 seconds
		MaxRetries:            12,                     // Retry for up to 1 minute (12 * 5 seconds)
		EnableTelegram:        true,                   // Enable Telegram notifications
		EnableTradeHistory:    true,                   // Enable trade history recording
		BlacklistedSymbols:    []string{},             // No blacklisted symbols by default
		TradeDelayAfterStatus: 500 * time.Millisecond, // Small delay after status change
	}
}

// NewNewCoinTraderService creates a new coin trader service
func NewNewCoinTraderService(
	statusUC usecase.StatusUseCase,
	tradeExecutor port.TradeExecutor,
	tradeHistory port.TradeHistoryRepository,
	telegramNotifier *notification.TelegramNotifier,
	config *NewCoinTraderConfig,
	logger *zerolog.Logger,
) *NewCoinTraderService {
	if config == nil {
		config = DefaultNewCoinTraderConfig()
	}

	serviceLogger := logger.With().Str("component", "newcoin_trader").Logger()
	ctx, cancel := context.WithCancel(context.Background())

	return &NewCoinTraderService{
		statusUC:         statusUC,
		tradeExecutor:    tradeExecutor,
		tradeHistory:     tradeHistory,
		telegramNotifier: telegramNotifier,
		config:           config,
		logger:           &serviceLogger,
		ctx:              ctx,
		cancel:           cancel,
		isRunning:        false,
	}
}

// Start begins watching for status changes and processing them
func (s *NewCoinTraderService) Start() error {
	if s.isRunning {
		return nil
	}

	// Subscribe to status changes
	statusCh := make(chan status.StatusChange, 100)
	err := s.statusUC.SubscribeToChanges(statusCh)
	if err != nil {
		return fmt.Errorf("failed to subscribe to status changes: %w", err)
	}

	s.isRunning = true
	s.logger.Info().Msg("NewCoin trader service started")

	// Start processing status changes in a goroutine
	go s.processStatusChanges(statusCh)

	return nil
}

// Stop halts the service
func (s *NewCoinTraderService) Stop() {
	if !s.isRunning {
		return
	}

	s.cancel() // Cancel the context to stop all goroutines
	s.isRunning = false
	s.logger.Info().Msg("NewCoin trader service stopped")
}

// IsRunning returns whether the service is running
func (s *NewCoinTraderService) IsRunning() bool {
	return s.isRunning
}

// processStatusChanges handles incoming status change events
func (s *NewCoinTraderService) processStatusChanges(statusCh chan status.StatusChange) {
	for {
		select {
		case change := <-statusCh:
			// Only process new coin status changes to StatusTrading
			if change.Component == "newcoin" && change.NewStatus == status.StatusTrading {
				s.logger.Info().
					Str("id", change.ID).
					Str("oldStatus", string(change.OldStatus)).
					Str("newStatus", string(change.NewStatus)).
					Interface("metadata", change.Metadata).
					Msg("Detected new coin status change to trading")

				// Process in a separate goroutine to avoid blocking
				go s.handleNewCoinTrading(change.ID, change.Metadata)
			}
		case <-s.ctx.Done():
			s.logger.Info().Msg("Status change processor stopped")
			return
		}
	}
}

// handleNewCoinTrading processes a new coin that has become tradable
func (s *NewCoinTraderService) handleNewCoinTrading(coinID string, metadata map[string]interface{}) {
	// Extract necessary information from metadata
	symbol, ok := metadata["symbol"].(string)
	if !ok || symbol == "" {
		s.logger.Error().Str("coinID", coinID).Msg("Symbol not found in metadata")
		return
	}

	// Check if symbol is blacklisted
	for _, blacklisted := range s.config.BlacklistedSymbols {
		if symbol == blacklisted {
			s.logger.Info().Str("symbol", symbol).Msg("Symbol is blacklisted, skipping trade")
			return
		}
	}

	// Add delay if configured (to allow market to stabilize)
	if s.config.TradeDelayAfterStatus > 0 {
		time.Sleep(s.config.TradeDelayAfterStatus)
	}

	// Extract trade amount from metadata or use default
	amount := s.config.DefaultTradeAmount
	if metaAmount, ok := metadata["tradeAmount"].(float64); ok && metaAmount > 0 {
		amount = metaAmount
		// Cap at maximum
		if amount > s.config.MaxTradeAmount {
			amount = s.config.MaxTradeAmount
		}
	}

	// Create order request for a market buy
	orderReq := &model.OrderRequest{
		Symbol:   symbol,
		Side:     model.OrderSideBuy,
		Type:     model.OrderTypeMarket,
		Quantity: amount,   // For market orders on MEXC, we use Quantity field
		UserID:   "system", // Use a system user ID for auto trades
	}

	// Attempt to execute the order with retries for "not yet tradable" errors
	var response *model.OrderResponse
	var err error
	var retryCount int

	for retryCount = 0; retryCount <= s.config.MaxRetries; retryCount++ {
		// Log first attempt or retries
		if retryCount == 0 {
			s.logger.Info().
				Str("symbol", symbol).
				Float64("amount", amount).
				Msg("Executing market order for new coin")
		} else {
			s.logger.Info().
				Str("symbol", symbol).
				Float64("amount", amount).
				Int("retry", retryCount).
				Msg("Retrying market order for new coin")
		}

		// Execute the order
		response, err = s.tradeExecutor.ExecuteOrder(s.ctx, orderReq)

		// If successful or a non-retryable error, break the loop
		if err == nil || !s.isNotYetTradableError(err) {
			break
		}

		// If max retries reached, break the loop
		if retryCount >= s.config.MaxRetries {
			s.logger.Warn().
				Str("symbol", symbol).
				Int("maxRetries", s.config.MaxRetries).
				Msg("Max retries reached for new coin order")
			break
		}

		// Wait before retrying
		select {
		case <-time.After(s.config.RetryInterval):
			// Continue to retry
		case <-s.ctx.Done():
			// Service is stopping, abort retries
			return
		}
	}

	// Handle the result
	if err != nil {
		s.handleOrderError(symbol, amount, err, retryCount)
		return
	}

	// Order successful, process the result
	s.handleSuccessfulOrder(symbol, response)
}

// isNotYetTradableError checks if the error is due to the coin not being tradable yet
func (s *NewCoinTraderService) isNotYetTradableError(err error) bool {
	// Check for various error messages that indicate the coin is not yet tradable
	// These strings will depend on the actual error messages from the exchange
	errorStr := err.Error()
	notTradableErrors := []string{
		"not yet tradable",
		"trading not open",
		"symbol not found",
		"invalid symbol",
		"This symbol is not available for trading",
	}

	for _, errText := range notTradableErrors {
		if s.containsIgnoreCase(errorStr, errText) {
			return true
		}
	}

	return false
}

// containsIgnoreCase checks if a string contains another string, case-insensitive
func (s *NewCoinTraderService) containsIgnoreCase(s1, s2 string) bool {
	s1, s2 = strings.ToLower(s1), strings.ToLower(s2)
	return strings.Contains(s1, s2)
}

// handleOrderError processes a failed order
func (s *NewCoinTraderService) handleOrderError(symbol string, amount float64, err error, retries int) {
	s.logger.Error().
		Err(err).
		Str("symbol", symbol).
		Float64("amount", amount).
		Int("retries", retries).
		Msg("Failed to execute market order for new coin")

	// Send notification about the failure
	if s.config.EnableTelegram && s.telegramNotifier != nil {
		errorMsg := err.Error()
		if len(errorMsg) > 200 {
			errorMsg = errorMsg[:200] + "..." // Limit length for Telegram
		}

		message := fmt.Sprintf("Failed to buy new coin %s after %d retries: %s",
			symbol, retries, errorMsg)

		_ = s.telegramNotifier.NotifyAlert(s.ctx, "error", "New Coin Order Failed", message, "newcoin_trader")
	}
}

// handleSuccessfulOrder processes a successful order
func (s *NewCoinTraderService) handleSuccessfulOrder(symbol string, response *model.OrderResponse) {
	s.logger.Info().
		Str("symbol", symbol).
		Str("orderID", response.OrderID).
		Float64("quantity", response.Quantity).
		Float64("price", response.Price).
		Float64("total", response.Quantity*response.Price).
		Msg("Successfully executed market order for new coin")

	// Save trade record if enabled
	if s.config.EnableTradeHistory && s.tradeHistory != nil {
		now := time.Now()
		tradeRecord := &model.TradeRecord{
			Symbol:        symbol,
			OrderID:       response.OrderID,
			Side:          model.OrderSideBuy,
			Type:          model.OrderTypeMarket,
			Quantity:      response.Quantity,
			Price:         response.Price,
			ExecutionTime: now,
			Strategy:      "newcoin_auto",
			Tags:          []string{"newcoin", "auto"},
			CreatedAt:     now,
			UpdatedAt:     now,
		}

		if err := s.tradeHistory.SaveTradeRecord(s.ctx, tradeRecord); err != nil {
			s.logger.Error().
				Err(err).
				Str("symbol", symbol).
				Msg("Failed to save trade record")
		}
	}

	// Send notification about successful trade
	if s.config.EnableTelegram && s.telegramNotifier != nil {
		_ = s.telegramNotifier.NotifyTrade(s.ctx, symbol, "buy", "market",
			response.Quantity, response.Price, response.OrderID)
	}
}
