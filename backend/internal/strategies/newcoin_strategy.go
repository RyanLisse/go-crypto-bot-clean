package strategies

import (
	"context"
	"fmt"
	"time"

	"go-crypto-bot-clean/backend/internal/domain/models"
	"go-crypto-bot-clean/backend/internal/trading"

	"go.uber.org/zap"
)

// NewCoinConfig holds configuration for the NewCoinStrategy
type NewCoinConfig struct {
	MinVolume         float64 // Minimum volume threshold
	MinPrice          float64 // Minimum price threshold
	MaxPrice          float64 // Maximum price threshold
	EntryDelay        int     // Delay in minutes before entering after detection
	StopLossPercent   float64 // Stop loss percentage
	TakeProfitPercent float64 // Take profit percentage
	PositionSize      float64 // Position size as percentage of portfolio
	MaxAge            int     // Maximum age in hours to consider a coin "new"
}

// NewCoinStrategy implements a strategy for trading newly listed coins
type NewCoinStrategy struct {
	BaseStrategy
	config        NewCoinConfig
	detectedCoins map[string]time.Time // Map of detected coins and their first seen timestamp
	entryPrices   map[string]float64   // Entry prices for each symbol
}

// NewNewCoinStrategy creates a new NewCoinStrategy
func NewNewCoinStrategy(logger *zap.Logger) trading.Strategy {
	// Default configuration
	config := map[string]interface{}{
		"min_volume":          1000.0,  // Default $1000 minimum volume
		"min_price":           0.00001, // Default minimum price
		"max_price":           10.0,    // Default maximum price
		"entry_delay":         5,       // Default 5 minutes delay
		"stop_loss_percent":   5.0,     // Default 5% stop loss
		"take_profit_percent": 10.0,    // Default 10% take profit
		"position_size":       1.0,     // Default 1% of portfolio
		"max_age":             24,      // Default 24 hours
	}

	strategy := &NewCoinStrategy{
		BaseStrategy: *NewBaseStrategy("NewCoinStrategy", config, logger),
		config: NewCoinConfig{
			MinVolume:         1000.0,
			MinPrice:          0.00001,
			MaxPrice:          10.0,
			EntryDelay:        5,
			StopLossPercent:   5.0,
			TakeProfitPercent: 10.0,
			PositionSize:      1.0,
			MaxAge:            24,
		},
		detectedCoins: make(map[string]time.Time),
		entryPrices:   make(map[string]float64),
	}

	return strategy
}

// Initialize sets up the strategy with any required configuration
func (s *NewCoinStrategy) Initialize(ctx context.Context, config interface{}) error {
	// Call base implementation first
	if err := s.BaseStrategy.Initialize(ctx, config); err != nil {
		return err
	}

	// Parse strategy-specific configuration
	s.config.MinVolume = s.GetFloatConfigValue("min_volume", 1000.0)
	s.config.MinPrice = s.GetFloatConfigValue("min_price", 0.00001)
	s.config.MaxPrice = s.GetFloatConfigValue("max_price", 10.0)
	s.config.EntryDelay = s.GetIntConfigValue("entry_delay", 5)
	s.config.StopLossPercent = s.GetFloatConfigValue("stop_loss_percent", 5.0)
	s.config.TakeProfitPercent = s.GetFloatConfigValue("take_profit_percent", 10.0)
	s.config.PositionSize = s.GetFloatConfigValue("position_size", 1.0)
	s.config.MaxAge = s.GetIntConfigValue("max_age", 24)

	// Initialize maps if they're nil
	if s.detectedCoins == nil {
		s.detectedCoins = make(map[string]time.Time)
	}
	if s.entryPrices == nil {
		s.entryPrices = make(map[string]float64)
	}

	return nil
}

// OnTick processes a new tick of market data and returns any trading signals
func (s *NewCoinStrategy) OnTick(ctx context.Context, symbol string, timestamp time.Time, data interface{}) ([]*trading.Signal, error) {
	// Call base implementation to update price history
	_, err := s.BaseStrategy.OnTick(ctx, symbol, timestamp, data)
	if err != nil {
		return nil, err
	}

	// Get the current price
	prices, ok := s.PriceHistory[symbol]
	if !ok || len(prices) == 0 {
		return nil, nil // No price data yet
	}
	currentPrice := prices[len(prices)-1]

	// Get the current volume
	volumes, ok := s.VolumeHistory[symbol]
	if !ok || len(volumes) == 0 {
		return nil, nil // No volume data yet
	}
	currentVolume := volumes[len(volumes)-1]

	// Check if we have an open position
	if s.Positions[symbol] {
		// Check for take profit or stop loss
		entryPrice, ok := s.entryPrices[symbol]
		if !ok {
			// This shouldn't happen, but just in case
			s.Logger.Warn("Position exists but no entry price found",
				zap.String("symbol", symbol),
				zap.Float64("current_price", currentPrice),
			)
			return nil, nil
		}

		// Calculate profit/loss percentage
		pnlPercent := (currentPrice - entryPrice) / entryPrice * 100.0

		// Check for take profit
		if pnlPercent >= s.config.TakeProfitPercent {
			s.Logger.Info("Take profit triggered",
				zap.String("symbol", symbol),
				zap.Float64("entry_price", entryPrice),
				zap.Float64("current_price", currentPrice),
				zap.Float64("pnl_percent", pnlPercent),
			)

			signal := &trading.Signal{
				Symbol:    symbol,
				Side:      models.OrderSideSell,
				Quantity:  1.0, // Will be adjusted based on position size
				Price:     currentPrice,
				Timestamp: timestamp,
				Reason:    fmt.Sprintf("Take profit triggered (%.2f%%)", pnlPercent),
			}
			return []*trading.Signal{signal}, nil
		}

		// Check for stop loss
		if pnlPercent <= -s.config.StopLossPercent {
			s.Logger.Info("Stop loss triggered",
				zap.String("symbol", symbol),
				zap.Float64("entry_price", entryPrice),
				zap.Float64("current_price", currentPrice),
				zap.Float64("pnl_percent", pnlPercent),
			)

			signal := &trading.Signal{
				Symbol:    symbol,
				Side:      models.OrderSideSell,
				Quantity:  1.0, // Will be adjusted based on position size
				Price:     currentPrice,
				Timestamp: timestamp,
				Reason:    fmt.Sprintf("Stop loss triggered (%.2f%%)", pnlPercent),
			}
			return []*trading.Signal{signal}, nil
		}

		// No signals for existing positions
		return nil, nil
	}

	// Check if this is a new coin we haven't seen before
	firstSeen, exists := s.detectedCoins[symbol]
	if !exists {
		// Only consider coins that meet our price criteria
		if currentPrice >= s.config.MinPrice && currentPrice <= s.config.MaxPrice {
			s.detectedCoins[symbol] = timestamp
			s.Logger.Info("New coin detected",
				zap.String("symbol", symbol),
				zap.Time("timestamp", timestamp),
				zap.Float64("price", currentPrice),
				zap.Float64("volume", currentVolume),
			)
		}
		return nil, nil // Just record it for now
	}

	// Check if the coin is still within our "new" window
	if timestamp.Sub(firstSeen).Hours() > float64(s.config.MaxAge) {
		delete(s.detectedCoins, symbol) // Remove old coins
		s.Logger.Info("Coin no longer considered new",
			zap.String("symbol", symbol),
			zap.Time("first_seen", firstSeen),
			zap.Time("current_time", timestamp),
		)
		return nil, nil
	}

	// Check if it meets our criteria
	if currentPrice < s.config.MinPrice || currentPrice > s.config.MaxPrice {
		return nil, nil
	}

	// Check volume criteria
	if currentVolume < s.config.MinVolume {
		return nil, nil
	}

	// Check if enough time has passed since detection
	if timestamp.Sub(firstSeen).Minutes() < float64(s.config.EntryDelay) {
		return nil, nil
	}

	// Check if we already have a position
	if s.Positions[symbol] {
		return nil, nil // Already have a position
	}

	// Generate buy signal
	s.Logger.Info("Buy signal generated for new coin",
		zap.String("symbol", symbol),
		zap.Time("timestamp", timestamp),
		zap.Time("first_seen", firstSeen),
		zap.Float64("price", currentPrice),
		zap.Float64("volume", currentVolume),
	)

	signal := &trading.Signal{
		Symbol:    symbol,
		Side:      models.OrderSideBuy,
		Quantity:  s.config.PositionSize, // Will be adjusted by position sizing
		Price:     currentPrice,
		Timestamp: timestamp,
		Reason:    fmt.Sprintf("New coin detected %s ago", timestamp.Sub(firstSeen).Round(time.Minute).String()),
	}

	return []*trading.Signal{signal}, nil
}

// OnOrderFilled is called when an order has been filled
func (s *NewCoinStrategy) OnOrderFilled(ctx context.Context, order *models.Order) error {
	// Call base implementation to update position tracking
	if err := s.BaseStrategy.OnOrderFilled(ctx, order); err != nil {
		return err
	}

	// Store entry price for buy orders
	if order.Side == models.OrderSideBuy {
		s.entryPrices[order.Symbol] = order.Price
	} else if order.Side == models.OrderSideSell {
		// Remove entry price for sell orders
		delete(s.entryPrices, order.Symbol)
	}

	return nil
}

// ClosePositions is called at the end of backtesting to close any open positions
func (s *NewCoinStrategy) ClosePositions(ctx context.Context) ([]*trading.Signal, error) {
	// Call base implementation to close all positions
	return s.BaseStrategy.ClosePositions(ctx)
}
