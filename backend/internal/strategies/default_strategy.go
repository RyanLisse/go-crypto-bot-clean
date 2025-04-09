package strategies

import (
	"context"
	"fmt"
	"time"

	"go-crypto-bot-clean/backend/internal/domain/models"
	"go-crypto-bot-clean/backend/internal/trading"

	"go.uber.org/zap"
)

// DefaultStrategy is a simple moving average crossover strategy
type DefaultStrategy struct {
	BaseStrategy // Embed BaseStrategy
	// Strategy-specific fields
	ShortPeriod int
	LongPeriod  int
}

// NewDefaultStrategy creates a new DefaultStrategy
func NewDefaultStrategy(logger *zap.Logger) trading.Strategy {
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
	s.ShortPeriod = s.GetIntConfigValue("short_period", 10)
	s.LongPeriod = s.GetIntConfigValue("long_period", 30)

	return nil
}

// OnTick implements the trading logic for DefaultStrategy
func (s *DefaultStrategy) OnTick(ctx context.Context, symbol string, timestamp time.Time, data interface{}) ([]*trading.Signal, error) {
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
	var signals []*trading.Signal

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

			signal := &trading.Signal{
				Symbol:    symbol,
				Side:      models.OrderSideBuy,
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

			signal := &trading.Signal{
				Symbol:    symbol,
				Side:      models.OrderSideSell,
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
