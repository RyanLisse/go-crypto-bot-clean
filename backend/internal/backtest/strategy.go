package backtest

import (
	"context"
	"fmt"
	"time"

	"go-crypto-bot-clean/backend/internal/domain/models"
	"go-crypto-bot-clean/backend/internal/indicators"

	"go.uber.org/zap"
)

// BaseStrategy provides a basic implementation of the BacktestStrategy interface
type BaseStrategy struct {
	// Common fields for strategies
	Name             string                 // Strategy name
	Config           map[string]interface{} // Strategy configuration
	Positions        map[string]bool        // Track open positions by symbol
	PriceHistory     map[string][]float64   // Store price history for each symbol
	MaxHistoryLength int                    // Maximum length of price history to keep
	Logger           *zap.Logger            // Logger for strategy events
	Metadata         map[string]interface{} // Additional strategy metadata
}

// NewBaseStrategy creates a new BaseStrategy with the given name and configuration
func NewBaseStrategy(name string, config map[string]interface{}, logger *zap.Logger) *BaseStrategy {
	if logger == nil {
		logger, _ = zap.NewProduction()
	}

	return &BaseStrategy{
		Name:             name,
		Config:           config,
		Positions:        make(map[string]bool),
		PriceHistory:     make(map[string][]float64),
		MaxHistoryLength: 1000, // Default to 1000 data points
		Logger:           logger,
		Metadata:         make(map[string]interface{}),
	}
}

// Initialize sets up the strategy
func (s *BaseStrategy) Initialize(ctx context.Context, config interface{}) error {
	// If config is a map, update the strategy configuration
	if cfg, ok := config.(map[string]interface{}); ok {
		s.Config = cfg
	}

	// Initialize other fields if they're nil
	if s.Positions == nil {
		s.Positions = make(map[string]bool)
	}
	if s.PriceHistory == nil {
		s.PriceHistory = make(map[string][]float64)
	}
	if s.Metadata == nil {
		s.Metadata = make(map[string]interface{})
	}

	return nil
}

// OnTick processes a market tick
func (s *BaseStrategy) OnTick(ctx context.Context, symbol string, timestamp time.Time, data interface{}) ([]*Signal, error) {
	// Update price history if data is a Kline or Candle
	price := s.extractPrice(data)
	if price > 0 {
		s.updatePriceHistory(symbol, price)
	}

	// Base implementation returns no signals, must be overridden
	return nil, nil
}

// extractPrice attempts to extract a price from various data types
func (s *BaseStrategy) extractPrice(data interface{}) float64 {
	switch v := data.(type) {
	case *models.Kline:
		return v.Close
	case *models.Candle:
		return v.ClosePrice
	case *models.Ticker:
		return v.Price
	case *models.Trade:
		return v.Price
	case float64:
		return v
	default:
		return 0
	}
}

// updatePriceHistory adds a price to the history for a symbol
func (s *BaseStrategy) updatePriceHistory(symbol string, price float64) {
	// Initialize price history for this symbol if it doesn't exist
	if _, ok := s.PriceHistory[symbol]; !ok {
		s.PriceHistory[symbol] = make([]float64, 0, s.MaxHistoryLength)
	}

	// Add the price to the history
	s.PriceHistory[symbol] = append(s.PriceHistory[symbol], price)

	// Trim the history if it exceeds the maximum length
	if len(s.PriceHistory[symbol]) > s.MaxHistoryLength {
		s.PriceHistory[symbol] = s.PriceHistory[symbol][len(s.PriceHistory[symbol])-s.MaxHistoryLength:]
	}
}

// OnOrderFilled handles filled orders
func (s *BaseStrategy) OnOrderFilled(ctx context.Context, order *models.Order) error {
	// Update position tracking
	if order.Side == models.OrderSideBuy {
		s.Positions[order.Symbol] = true
		s.Logger.Info("Long position opened",
			zap.String("symbol", order.Symbol),
			zap.Float64("price", order.Price),
			zap.Float64("quantity", order.Quantity),
			zap.String("strategy", s.Name),
		)
	} else if order.Side == models.OrderSideSell {
		s.Positions[order.Symbol] = false
		s.Logger.Info("Position closed",
			zap.String("symbol", order.Symbol),
			zap.Float64("price", order.Price),
			zap.Float64("quantity", order.Quantity),
			zap.String("strategy", s.Name),
		)
	}
	return nil
}

// ClosePositions is called at the end of the backtest
func (s *BaseStrategy) ClosePositions(ctx context.Context) ([]*Signal, error) {
	// Generate signals to close all open positions
	signals := make([]*Signal, 0)

	for symbol, isOpen := range s.Positions {
		if isOpen {
			s.Logger.Info("Closing position at end of backtest",
				zap.String("symbol", symbol),
				zap.String("strategy", s.Name),
			)

			// Create a sell signal to close the position
			signal := &Signal{
				Symbol:    symbol,
				Side:      "SELL",
				Quantity:  0, // Will be filled in by the backtester based on position size
				Price:     0, // Will be filled in by the backtester based on current price
				Timestamp: time.Now(),
				Reason:    "End of backtest - closing position",
			}
			signals = append(signals, signal)
		}
	}

	return signals, nil
}

// CalculateIndicator calculates a technical indicator on the price history
func (s *BaseStrategy) CalculateIndicator(symbol string, indicator string, period int) ([]float64, error) {
	// Check if we have price history for this symbol
	prices, ok := s.PriceHistory[symbol]
	if !ok || len(prices) < period {
		return nil, fmt.Errorf("insufficient price history for %s", symbol)
	}

	// Calculate the requested indicator
	switch indicator {
	case "sma":
		return indicators.SMA(prices, period)
	case "ema":
		return indicators.EMA(prices, period)
	case "rsi":
		return indicators.RSI(prices, period)
	case "bollinger":
		upper, middle, lower, err := indicators.BollingerBands(prices, period, 2.0)
		if err != nil {
			return nil, err
		}
		return middle, nil // Return the middle band by default
	default:
		return nil, fmt.Errorf("unsupported indicator: %s", indicator)
	}
}

// DefaultStrategy is a simple example strategy
type DefaultStrategy struct {
	BaseStrategy // Embed BaseStrategy
	// Strategy-specific fields
	ShortPeriod int
	LongPeriod  int
}

// NewDefaultStrategy creates a new DefaultStrategy
func NewDefaultStrategy(logger *zap.Logger) BacktestStrategy {
	config := map[string]interface{}{
		"short_period": 10,
		"long_period":  30,
	}

	strategy := &DefaultStrategy{
		BaseStrategy: *NewBaseStrategy("DefaultStrategy", config, logger),
		ShortPeriod:  10,
		LongPeriod:   30,
	}

	return strategy
}

// Initialize overrides the base Initialize method to match the interface
func (s *DefaultStrategy) Initialize(ctx context.Context, config interface{}) error {
	// Call base implementation first
	if err := s.BaseStrategy.Initialize(ctx, config); err != nil {
		return err
	}

	// Parse strategy-specific configuration
	if cfg, ok := config.(map[string]interface{}); ok {
		if shortPeriod, ok := cfg["short_period"].(int); ok {
			s.ShortPeriod = shortPeriod
		}
		if longPeriod, ok := cfg["long_period"].(int); ok {
			s.LongPeriod = longPeriod
		}
	}

	return nil
}

// OnTick implements the trading logic for DefaultStrategy
func (s *DefaultStrategy) OnTick(ctx context.Context, symbol string, timestamp time.Time, data interface{}) ([]*Signal, error) {
	// Call the base implementation to update price history
	_, err := s.BaseStrategy.OnTick(ctx, symbol, timestamp, data)
	if err != nil {
		return nil, err
	}

	// Check if we have enough price history
	prices, ok := s.PriceHistory[symbol]
	if !ok || len(prices) < s.LongPeriod {
		return nil, nil // Not enough data yet
	}

	// Calculate moving averages
	shortMA, err := s.CalculateIndicator(symbol, "sma", s.ShortPeriod)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate short MA: %w", err)
	}

	longMA, err := s.CalculateIndicator(symbol, "sma", s.LongPeriod)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate long MA: %w", err)
	}

	// Check for crossover signals
	var signals []*Signal

	// Need at least 2 data points to check for crossover
	if len(shortMA) < 2 || len(longMA) < 2 {
		return nil, nil
	}

	// Get the current and previous values
	currentShortMA := shortMA[len(shortMA)-1]
	currentLongMA := longMA[len(longMA)-1]
	previousShortMA := shortMA[len(shortMA)-2]
	previousLongMA := longMA[len(longMA)-2]

	// Check for bullish crossover (short MA crosses above long MA)
	if previousShortMA <= previousLongMA && currentShortMA > currentLongMA {
		// Only generate a buy signal if we don't already have a position
		if !s.Positions[symbol] {
			s.Logger.Info("Buy signal generated",
				zap.String("symbol", symbol),
				zap.Time("timestamp", timestamp),
				zap.Float64("price", prices[len(prices)-1]),
				zap.Float64("short_ma", currentShortMA),
				zap.Float64("long_ma", currentLongMA),
			)

			signal := &Signal{
				Symbol:    symbol,
				Side:      "BUY",
				Quantity:  1.0, // Will be adjusted by position sizing
				Price:     prices[len(prices)-1],
				Timestamp: timestamp,
				Reason:    "MA Crossover (Bullish)",
			}
			signals = append(signals, signal)
		}
	}

	// Check for bearish crossover (short MA crosses below long MA)
	if previousShortMA >= previousLongMA && currentShortMA < currentLongMA {
		// Only generate a sell signal if we have an open position
		if s.Positions[symbol] {
			s.Logger.Info("Sell signal generated",
				zap.String("symbol", symbol),
				zap.Time("timestamp", timestamp),
				zap.Float64("price", prices[len(prices)-1]),
				zap.Float64("short_ma", currentShortMA),
				zap.Float64("long_ma", currentLongMA),
			)

			signal := &Signal{
				Symbol:    symbol,
				Side:      "SELL",
				Quantity:  1.0, // Will be adjusted based on position size
				Price:     prices[len(prices)-1],
				Timestamp: timestamp,
				Reason:    "MA Crossover (Bearish)",
			}
			signals = append(signals, signal)
		}
	}

	return signals, nil
}

// ClosePositions implements the method required by the interface
// DefaultStrategy might not need to do anything specific on close
func (s *DefaultStrategy) ClosePositions(ctx context.Context) ([]*Signal, error) {
	// You could call the base implementation if it ever does something:
	// return s.BaseStrategy.ClosePositions(ctx)
	return nil, nil // Default implementation returns no closing signals
}

// Ensure DefaultStrategy implements BacktestStrategy
var _ BacktestStrategy = (*DefaultStrategy)(nil)
