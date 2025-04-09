package strategies

import (
	"context"
	"fmt"
	"math"
	"time"

	"go-crypto-bot-clean/backend/internal/domain/models"
	"go-crypto-bot-clean/backend/internal/trading"

	"go.uber.org/zap"
)

// BreakoutConfig holds configuration for the BreakoutStrategy
type BreakoutConfig struct {
	RangePeriod         int     // Number of periods to look back for range calculation
	MinRangePercent     float64 // Minimum range as percentage of price
	BreakoutPercent     float64 // Percentage beyond range to confirm breakout
	MinVolume           float64 // Minimum volume threshold
	MinPrice            float64 // Minimum price threshold
	MaxPrice            float64 // Maximum price threshold
	StopLossPercent     float64 // Stop loss percentage
	TakeProfitPercent   float64 // Take profit percentage
	PositionSize        float64 // Position size as percentage of portfolio
	ExitAfterBars       int     // Exit position after this many bars if neither TP nor SL hit
	ConfirmationPeriods int     // Number of periods to confirm breakout
}

// BreakoutStrategy implements a strategy for trading price breakouts
type BreakoutStrategy struct {
	BaseStrategy
	config           BreakoutConfig
	entryPrices      map[string]float64 // Entry prices for each symbol
	entryBars        map[string]int     // Number of bars since entry for each symbol
	resistanceLevels map[string]float64 // Resistance levels for each symbol
	supportLevels    map[string]float64 // Support levels for each symbol
	confirmations    map[string]int     // Confirmation counter for breakouts
}

// NewBreakoutStrategy creates a new BreakoutStrategy
func NewBreakoutStrategy(logger *zap.Logger) trading.Strategy {
	// Default configuration
	config := map[string]interface{}{
		"range_period":         20,     // Look back 20 periods for range
		"min_range_percent":    5.0,    // Minimum 5% range
		"breakout_percent":     1.0,    // 1% beyond range to confirm breakout
		"min_volume":           5000.0, // Minimum $5000 volume
		"min_price":            0.1,    // Minimum price
		"max_price":            1000.0, // Maximum price
		"stop_loss_percent":    5.0,    // 5% stop loss
		"take_profit_percent":  15.0,   // 15% take profit
		"position_size":        2.0,    // 2% of portfolio
		"exit_after_bars":      20,     // Exit after 20 bars if neither TP nor SL hit
		"confirmation_periods": 2,      // Require 2 periods to confirm breakout
	}

	strategy := &BreakoutStrategy{
		BaseStrategy: *NewBaseStrategy("BreakoutStrategy", config, logger),
		config: BreakoutConfig{
			RangePeriod:         20,
			MinRangePercent:     5.0,
			BreakoutPercent:     1.0,
			MinVolume:           5000.0,
			MinPrice:            0.1,
			MaxPrice:            1000.0,
			StopLossPercent:     5.0,
			TakeProfitPercent:   15.0,
			PositionSize:        2.0,
			ExitAfterBars:       20,
			ConfirmationPeriods: 2,
		},
		entryPrices:      make(map[string]float64),
		entryBars:        make(map[string]int),
		resistanceLevels: make(map[string]float64),
		supportLevels:    make(map[string]float64),
		confirmations:    make(map[string]int),
	}

	return strategy
}

// Initialize sets up the strategy with any required configuration
func (s *BreakoutStrategy) Initialize(ctx context.Context, config interface{}) error {
	// Call base implementation first
	if err := s.BaseStrategy.Initialize(ctx, config); err != nil {
		return err
	}

	// Parse strategy-specific configuration
	s.config.RangePeriod = s.GetIntConfigValue("range_period", 20)
	s.config.MinRangePercent = s.GetFloatConfigValue("min_range_percent", 5.0)
	s.config.BreakoutPercent = s.GetFloatConfigValue("breakout_percent", 1.0)
	s.config.MinVolume = s.GetFloatConfigValue("min_volume", 5000.0)
	s.config.MinPrice = s.GetFloatConfigValue("min_price", 0.1)
	s.config.MaxPrice = s.GetFloatConfigValue("max_price", 1000.0)
	s.config.StopLossPercent = s.GetFloatConfigValue("stop_loss_percent", 5.0)
	s.config.TakeProfitPercent = s.GetFloatConfigValue("take_profit_percent", 15.0)
	s.config.PositionSize = s.GetFloatConfigValue("position_size", 2.0)
	s.config.ExitAfterBars = s.GetIntConfigValue("exit_after_bars", 20)
	s.config.ConfirmationPeriods = s.GetIntConfigValue("confirmation_periods", 2)

	// Initialize maps if they're nil
	if s.entryPrices == nil {
		s.entryPrices = make(map[string]float64)
	}
	if s.entryBars == nil {
		s.entryBars = make(map[string]int)
	}
	if s.resistanceLevels == nil {
		s.resistanceLevels = make(map[string]float64)
	}
	if s.supportLevels == nil {
		s.supportLevels = make(map[string]float64)
	}
	if s.confirmations == nil {
		s.confirmations = make(map[string]int)
	}

	return nil
}

// OnTick processes a new tick of market data and returns any trading signals
func (s *BreakoutStrategy) OnTick(ctx context.Context, symbol string, timestamp time.Time, data interface{}) ([]*trading.Signal, error) {
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
		// Increment the number of bars since entry
		s.entryBars[symbol]++

		// Check for take profit, stop loss, or time-based exit
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

		// Check for time-based exit
		if s.entryBars[symbol] >= s.config.ExitAfterBars {
			s.Logger.Info("Time-based exit triggered",
				zap.String("symbol", symbol),
				zap.Float64("entry_price", entryPrice),
				zap.Float64("current_price", currentPrice),
				zap.Float64("pnl_percent", pnlPercent),
				zap.Int("bars_held", s.entryBars[symbol]),
			)

			signal := &trading.Signal{
				Symbol:    symbol,
				Side:      models.OrderSideSell,
				Quantity:  1.0, // Will be adjusted based on position size
				Price:     currentPrice,
				Timestamp: timestamp,
				Reason:    fmt.Sprintf("Time-based exit after %d bars (%.2f%%)", s.entryBars[symbol], pnlPercent),
			}
			return []*trading.Signal{signal}, nil
		}

		// No signals for existing positions
		return nil, nil
	}

	// Check if we have enough price history
	if len(prices) <= s.config.RangePeriod {
		return nil, nil // Not enough data yet
	}

	// Check price criteria
	if currentPrice < s.config.MinPrice || currentPrice > s.config.MaxPrice {
		return nil, nil
	}

	// Check volume criteria
	if currentVolume < s.config.MinVolume {
		return nil, nil
	}

	// Calculate support and resistance levels
	s.calculateSupportResistance(symbol)

	// Get support and resistance levels
	support, ok1 := s.supportLevels[symbol]
	resistance, ok2 := s.resistanceLevels[symbol]
	if !ok1 || !ok2 {
		return nil, nil // No support/resistance levels calculated
	}

	// Calculate range as percentage of average price
	avgPrice := (support + resistance) / 2
	rangePercent := (resistance - support) / avgPrice * 100.0

	// Check if range is significant enough
	if rangePercent < s.config.MinRangePercent {
		return nil, nil // Range too small
	}

	// Check for breakout
	breakoutThreshold := resistance * (1 + s.config.BreakoutPercent/100.0)
	if currentPrice > breakoutThreshold {
		// Increment confirmation counter
		s.confirmations[symbol]++

		// Check if we have enough confirmations
		if s.confirmations[symbol] >= s.config.ConfirmationPeriods {
			// Generate buy signal
			s.Logger.Info("Resistance breakout detected",
				zap.String("symbol", symbol),
				zap.Time("timestamp", timestamp),
				zap.Float64("resistance", resistance),
				zap.Float64("breakout_threshold", breakoutThreshold),
				zap.Float64("current_price", currentPrice),
				zap.Float64("range_percent", rangePercent),
			)

			signal := &trading.Signal{
				Symbol:    symbol,
				Side:      models.OrderSideBuy,
				Quantity:  s.config.PositionSize, // Will be adjusted by position sizing
				Price:     currentPrice,
				Timestamp: timestamp,
				Reason:    fmt.Sprintf("Resistance breakout (%.2f%% above)", (currentPrice-resistance)/resistance*100.0),
			}

			// Reset confirmation counter
			s.confirmations[symbol] = 0

			return []*trading.Signal{signal}, nil
		}
	} else {
		// Reset confirmation counter if price falls below threshold
		s.confirmations[symbol] = 0
	}

	return nil, nil
}

// calculateSupportResistance calculates support and resistance levels for a symbol
func (s *BreakoutStrategy) calculateSupportResistance(symbol string) {
	prices, ok := s.PriceHistory[symbol]
	if !ok || len(prices) <= s.config.RangePeriod {
		return // Not enough data
	}

	// Get the range of prices to analyze
	rangePrices := prices[len(prices)-s.config.RangePeriod:]

	// Find highest high and lowest low
	high := math.Inf(-1)
	low := math.Inf(1)

	for _, price := range rangePrices {
		if price > high {
			high = price
		}
		if price < low {
			low = price
		}
	}

	// Update support and resistance levels
	s.supportLevels[symbol] = low
	s.resistanceLevels[symbol] = high
}

// OnOrderFilled is called when an order has been filled
func (s *BreakoutStrategy) OnOrderFilled(ctx context.Context, order *models.Order) error {
	// Call base implementation to update position tracking
	if err := s.BaseStrategy.OnOrderFilled(ctx, order); err != nil {
		return err
	}

	// Store entry price and reset bar counter for buy orders
	if order.Side == models.OrderSideBuy {
		s.entryPrices[order.Symbol] = order.Price
		s.entryBars[order.Symbol] = 0
	} else if order.Side == models.OrderSideSell {
		// Remove entry price and bar counter for sell orders
		delete(s.entryPrices, order.Symbol)
		delete(s.entryBars, order.Symbol)
	}

	return nil
}

// ClosePositions is called at the end of backtesting to close any open positions
func (s *BreakoutStrategy) ClosePositions(ctx context.Context) ([]*trading.Signal, error) {
	// Call base implementation to close all positions
	return s.BaseStrategy.ClosePositions(ctx)
}
