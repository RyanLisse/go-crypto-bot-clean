package strategies

import (
	"context"
	"fmt"
	"time"

	"go-crypto-bot-clean/backend/internal/backtest"
	"go-crypto-bot-clean/backend/internal/domain/models"
	"go-crypto-bot-clean/backend/internal/indicators"
)

// MovingAverageCrossoverStrategy implements a simple moving average crossover strategy
type MovingAverageCrossoverStrategy struct {
	backtest.BaseStrategy
	FastPeriod        int
	SlowPeriod        int
	RiskPercentage    float64
	StopLossPercent   float64
	TakeProfitPercent float64
	Positions         map[string]bool      // Track open positions by symbol
	PriceHistory      map[string][]float64 // Store price history for each symbol
}

// NewMovingAverageCrossoverStrategy creates a new MovingAverageCrossoverStrategy
func NewMovingAverageCrossoverStrategy(fastPeriod, slowPeriod int, riskPercentage, stopLossPercent, takeProfitPercent float64) *MovingAverageCrossoverStrategy {
	return &MovingAverageCrossoverStrategy{
		BaseStrategy:      backtest.BaseStrategy{ /* No Name field */ },
		FastPeriod:        fastPeriod,
		SlowPeriod:        slowPeriod,
		RiskPercentage:    riskPercentage,
		StopLossPercent:   stopLossPercent,
		TakeProfitPercent: takeProfitPercent,
		Positions:         make(map[string]bool),
		PriceHistory:      make(map[string][]float64),
	}
}

// Initialize initializes the strategy with backtest-specific parameters
func (s *MovingAverageCrossoverStrategy) Initialize(ctx context.Context, config map[string]interface{}) error {
	err := s.BaseStrategy.Initialize(ctx, config)
	if err != nil {
		return err
	}

	// Override parameters from config if provided
	if fastPeriod, ok := config["fast_period"].(int); ok {
		s.FastPeriod = fastPeriod
	}
	if slowPeriod, ok := config["slow_period"].(int); ok {
		s.SlowPeriod = slowPeriod
	}
	if riskPercentage, ok := config["risk_percentage"].(float64); ok {
		s.RiskPercentage = riskPercentage
	}
	if stopLossPercent, ok := config["stop_loss_percent"].(float64); ok {
		s.StopLossPercent = stopLossPercent
	}
	if takeProfitPercent, ok := config["take_profit_percent"].(float64); ok {
		s.TakeProfitPercent = takeProfitPercent
	}

	return nil
}

// OnTick is called for each new data point (candle, ticker, etc.)
func (s *MovingAverageCrossoverStrategy) OnTick(ctx context.Context, symbol string, timestamp time.Time, data interface{}) ([]*backtest.Signal, error) {
	// Check if data is a Kline
	kline, ok := data.(*models.Kline)
	if !ok {
		return nil, fmt.Errorf("expected *models.Kline, got %T", data)
	}

	// Update price history
	if _, ok := s.PriceHistory[symbol]; !ok {
		s.PriceHistory[symbol] = make([]float64, 0)
	}
	s.PriceHistory[symbol] = append(s.PriceHistory[symbol], kline.Close)

	// Check if we have enough data points
	if len(s.PriceHistory[symbol]) < s.SlowPeriod {
		return nil, nil
	}

	// Calculate moving averages
	fastMA, err := indicators.SMA(s.PriceHistory[symbol], s.FastPeriod)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate fast MA: %w", err)
	}

	slowMA, err := indicators.SMA(s.PriceHistory[symbol], s.SlowPeriod)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate slow MA: %w", err)
	}

	// Check for crossover
	var signals []*backtest.Signal

	// Get the current and previous values
	currentFastMA := fastMA[len(fastMA)-1]
	currentSlowMA := slowMA[len(slowMA)-1]

	// Need at least 2 data points to check for crossover
	if len(fastMA) < 2 || len(slowMA) < 2 {
		return nil, nil
	}

	previousFastMA := fastMA[len(fastMA)-2]
	previousSlowMA := slowMA[len(slowMA)-2]

	// Check for bullish crossover (fast MA crosses above slow MA)
	if previousFastMA <= previousSlowMA && currentFastMA > currentSlowMA {
		// Only generate a buy signal if we don't already have a position
		if !s.Positions[symbol] {
			signal := &backtest.Signal{
				Symbol:    symbol,
				Side:      "BUY",
				Quantity:  1.0, // Will be adjusted by position sizing
				Price:     kline.Close,
				Timestamp: timestamp,
				Reason:    "MA Crossover (Bullish)",
			}
			signals = append(signals, signal)
		}
	}

	// Check for bearish crossover (fast MA crosses below slow MA)
	if previousFastMA >= previousSlowMA && currentFastMA < currentSlowMA {
		// Only generate a sell signal if we have a position
		if s.Positions[symbol] {
			signal := &backtest.Signal{
				Symbol:    symbol,
				Side:      "SELL",
				Quantity:  1.0, // Will be adjusted to match the position size
				Price:     kline.Close,
				Timestamp: timestamp,
				Reason:    "MA Crossover (Bearish)",
			}
			signals = append(signals, signal)
		}
	}

	return signals, nil
}

// OnOrderFilled is called when an order is filled during the backtest
func (s *MovingAverageCrossoverStrategy) OnOrderFilled(ctx context.Context, order *models.Order) error {
	// Update position tracking
	if order.Side == models.OrderSideBuy {
		s.Positions[order.Symbol] = true
	} else if order.Side == models.OrderSideSell {
		s.Positions[order.Symbol] = false
	}
	return nil
}

// OnPositionClosed is called when a position is closed during the backtest
func (s *MovingAverageCrossoverStrategy) OnPositionClosed(ctx context.Context, position *models.ClosedPosition) error {
	// Update position tracking
	s.Positions[position.Symbol] = false
	return nil
}
