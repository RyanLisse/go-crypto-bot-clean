package strategies

import (
	"context"
	"fmt"
	"time"

	"go-crypto-bot-clean/backend/internal/domain/models"
	"go-crypto-bot-clean/backend/internal/trading"

	"go.uber.org/zap"
)

// VolumeSpikeConfig holds configuration for the VolumeSpikeStrategy
type VolumeSpikeConfig struct {
	VolumeThresholdPercent float64 // Volume increase threshold percentage
	MinVolume              float64 // Minimum volume threshold
	MinPrice               float64 // Minimum price threshold
	MaxPrice               float64 // Maximum price threshold
	LookbackPeriod         int     // Number of periods to look back for volume comparison
	StopLossPercent        float64 // Stop loss percentage
	TakeProfitPercent      float64 // Take profit percentage
	PositionSize           float64 // Position size as percentage of portfolio
	ExitAfterBars          int     // Exit position after this many bars if neither TP nor SL hit
}

// VolumeSpikeStrategy implements a strategy for trading volume spikes
type VolumeSpikeStrategy struct {
	BaseStrategy
	config      VolumeSpikeConfig
	entryPrices map[string]float64 // Entry prices for each symbol
	entryBars   map[string]int     // Number of bars since entry for each symbol
}

// NewVolumeSpikeStrategy creates a new VolumeSpikeStrategy
func NewVolumeSpikeStrategy(logger *zap.Logger) trading.Strategy {
	// Default configuration
	config := map[string]interface{}{
		"volume_threshold_percent": 200.0,  // 200% increase in volume
		"min_volume":               5000.0, // Minimum $5000 volume
		"min_price":                0.1,    // Minimum price
		"max_price":                1000.0, // Maximum price
		"lookback_period":          5,      // Look back 5 periods
		"stop_loss_percent":        5.0,    // 5% stop loss
		"take_profit_percent":      15.0,   // 15% take profit
		"position_size":            2.0,    // 2% of portfolio
		"exit_after_bars":          20,     // Exit after 20 bars if neither TP nor SL hit
	}

	strategy := &VolumeSpikeStrategy{
		BaseStrategy: *NewBaseStrategy("VolumeSpikeStrategy", config, logger),
		config: VolumeSpikeConfig{
			VolumeThresholdPercent: 200.0,
			MinVolume:              5000.0,
			MinPrice:               0.1,
			MaxPrice:               1000.0,
			LookbackPeriod:         5,
			StopLossPercent:        5.0,
			TakeProfitPercent:      15.0,
			PositionSize:           2.0,
			ExitAfterBars:          20,
		},
		entryPrices: make(map[string]float64),
		entryBars:   make(map[string]int),
	}

	return strategy
}

// Initialize sets up the strategy with any required configuration
func (s *VolumeSpikeStrategy) Initialize(ctx context.Context, config interface{}) error {
	// Call base implementation first
	if err := s.BaseStrategy.Initialize(ctx, config); err != nil {
		return err
	}

	// Parse strategy-specific configuration
	s.config.VolumeThresholdPercent = s.GetFloatConfigValue("volume_threshold_percent", 200.0)
	s.config.MinVolume = s.GetFloatConfigValue("min_volume", 5000.0)
	s.config.MinPrice = s.GetFloatConfigValue("min_price", 0.1)
	s.config.MaxPrice = s.GetFloatConfigValue("max_price", 1000.0)
	s.config.LookbackPeriod = s.GetIntConfigValue("lookback_period", 5)
	s.config.StopLossPercent = s.GetFloatConfigValue("stop_loss_percent", 5.0)
	s.config.TakeProfitPercent = s.GetFloatConfigValue("take_profit_percent", 15.0)
	s.config.PositionSize = s.GetFloatConfigValue("position_size", 2.0)
	s.config.ExitAfterBars = s.GetIntConfigValue("exit_after_bars", 20)

	// Initialize maps if they're nil
	if s.entryPrices == nil {
		s.entryPrices = make(map[string]float64)
	}
	if s.entryBars == nil {
		s.entryBars = make(map[string]int)
	}

	return nil
}

// OnTick processes a new tick of market data and returns any trading signals
func (s *VolumeSpikeStrategy) OnTick(ctx context.Context, symbol string, timestamp time.Time, data interface{}) ([]*trading.Signal, error) {
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

	// Check if we have enough volume history
	if len(volumes) <= s.config.LookbackPeriod {
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

	// Calculate average volume over lookback period
	var totalVolume float64
	for i := 1; i <= s.config.LookbackPeriod; i++ {
		totalVolume += volumes[len(volumes)-1-i]
	}
	avgVolume := totalVolume / float64(s.config.LookbackPeriod)

	// Check for volume spike
	volumeIncrease := (currentVolume - avgVolume) / avgVolume * 100.0
	if volumeIncrease < s.config.VolumeThresholdPercent {
		return nil, nil // No volume spike
	}

	// Generate buy signal
	s.Logger.Info("Volume spike detected",
		zap.String("symbol", symbol),
		zap.Time("timestamp", timestamp),
		zap.Float64("current_volume", currentVolume),
		zap.Float64("avg_volume", avgVolume),
		zap.Float64("volume_increase", volumeIncrease),
		zap.Float64("price", currentPrice),
	)

	signal := &trading.Signal{
		Symbol:    symbol,
		Side:      models.OrderSideBuy,
		Quantity:  s.config.PositionSize, // Will be adjusted by position sizing
		Price:     currentPrice,
		Timestamp: timestamp,
		Reason:    fmt.Sprintf("Volume spike detected (%.2f%% increase)", volumeIncrease),
	}

	return []*trading.Signal{signal}, nil
}

// OnOrderFilled is called when an order has been filled
func (s *VolumeSpikeStrategy) OnOrderFilled(ctx context.Context, order *models.Order) error {
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
func (s *VolumeSpikeStrategy) ClosePositions(ctx context.Context) ([]*trading.Signal, error) {
	// Call base implementation to close all positions
	return s.BaseStrategy.ClosePositions(ctx)
}
