package strategies

import (
	"context"
	"time"

	"go-crypto-bot-clean/backend/internal/backtest"
	"go-crypto-bot-clean/backend/internal/domain/models"

	"go.uber.org/zap"
)

// SimpleMAStrategy implements a simple moving average crossover strategy
type SimpleMAStrategy struct {
	backtest.BaseStrategy
	ShortPeriod int
	LongPeriod  int
	shortMA     map[string][]float64
	longMA      map[string][]float64
	prices      map[string][]float64
	position    map[string]bool // true if long position is open
	logger      *zap.Logger
}

// NewSimpleMAStrategy creates a new SimpleMAStrategy
func NewSimpleMAStrategy(shortPeriod, longPeriod int, logger *zap.Logger) *SimpleMAStrategy {
	if logger == nil {
		logger, _ = zap.NewDevelopment()
	}

	return &SimpleMAStrategy{
		BaseStrategy: backtest.BaseStrategy{ /* No Name field */ },
		ShortPeriod:  shortPeriod,
		LongPeriod:   longPeriod,
		shortMA:      make(map[string][]float64),
		longMA:       make(map[string][]float64),
		prices:       make(map[string][]float64),
		position:     make(map[string]bool),
		logger:       logger,
	}
}

// Initialize initializes the strategy with backtest-specific parameters
// It now accepts interface{} to match the BacktestStrategy interface
func (s *SimpleMAStrategy) Initialize(ctx context.Context, config interface{}) error {
	// Optionally call base initialize if BaseStrategy requires it
	err := s.BaseStrategy.Initialize(ctx, config)
	if err != nil {
		return err
	}

	// Attempt to cast config to the expected map type
	if configMap, ok := config.(map[string]interface{}); ok && configMap != nil {
		// Override periods if provided in config
		if shortPeriod, ok := configMap["short_period"].(float64); ok {
			s.ShortPeriod = int(shortPeriod) // Cast float64 from potential JSON parsing
		}
		if longPeriod, ok := configMap["long_period"].(float64); ok {
			s.LongPeriod = int(longPeriod) // Cast float64 from potential JSON parsing
		}
	} // else: config is nil or not the expected type, use defaults

	s.logger.Info("Initialized SimpleMAStrategy",
		zap.Int("short_period", s.ShortPeriod),
		zap.Int("long_period", s.LongPeriod),
	)

	return nil
}

// OnTick is called for each new data point (candle, ticker, etc.)
func (s *SimpleMAStrategy) OnTick(ctx context.Context, symbol string, timestamp time.Time, data interface{}) ([]*backtest.Signal, error) {
	// Check if data is a Kline
	kline, ok := data.(*models.Kline)
	if !ok {
		return nil, nil
	}

	// Initialize arrays for this symbol if they don't exist
	if _, ok := s.prices[symbol]; !ok {
		s.prices[symbol] = make([]float64, 0)
		s.shortMA[symbol] = make([]float64, 0)
		s.longMA[symbol] = make([]float64, 0)
		s.position[symbol] = false
	}

	// Add price to the price array
	s.prices[symbol] = append(s.prices[symbol], kline.Close)

	// Calculate moving averages
	shortMA := calculateSMA(s.prices[symbol], s.ShortPeriod)
	longMA := calculateSMA(s.prices[symbol], s.LongPeriod)

	// Add moving averages to the arrays
	s.shortMA[symbol] = append(s.shortMA[symbol], shortMA)
	s.longMA[symbol] = append(s.longMA[symbol], longMA)

	// Generate signals
	var signals []*backtest.Signal

	// We need enough data to calculate both moving averages
	if len(s.prices[symbol]) >= s.LongPeriod {
		// Check for crossover
		if len(s.shortMA[symbol]) >= 2 && len(s.longMA[symbol]) >= 2 {
			prevShortMA := s.shortMA[symbol][len(s.shortMA[symbol])-2]
			prevLongMA := s.longMA[symbol][len(s.longMA[symbol])-2]
			currentShortMA := s.shortMA[symbol][len(s.shortMA[symbol])-1]
			currentLongMA := s.longMA[symbol][len(s.longMA[symbol])-1]

			// Buy signal: short MA crosses above long MA
			if prevShortMA <= prevLongMA && currentShortMA > currentLongMA && !s.position[symbol] {
				s.logger.Info("Buy signal generated",
					zap.String("symbol", symbol),
					zap.Time("timestamp", timestamp),
					zap.Float64("price", kline.Close),
					zap.Float64("short_ma", currentShortMA),
					zap.Float64("long_ma", currentLongMA),
				)

				signals = append(signals, &backtest.Signal{
					Symbol:    symbol,
					Side:      "BUY",
					Quantity:  1.0, // Fixed quantity for simplicity
					Price:     kline.Close,
					Timestamp: timestamp,
					Reason:    "MA crossover (short > long)",
				})

				s.position[symbol] = true
			}

			// Sell signal: short MA crosses below long MA
			if prevShortMA >= prevLongMA && currentShortMA < currentLongMA && s.position[symbol] {
				s.logger.Info("Sell signal generated",
					zap.String("symbol", symbol),
					zap.Time("timestamp", timestamp),
					zap.Float64("price", kline.Close),
					zap.Float64("short_ma", currentShortMA),
					zap.Float64("long_ma", currentLongMA),
				)

				signals = append(signals, &backtest.Signal{
					Symbol:    symbol,
					Side:      "SELL",
					Quantity:  1.0, // Fixed quantity for simplicity
					Price:     kline.Close,
					Timestamp: timestamp,
					Reason:    "MA crossover (short < long)",
				})

				s.position[symbol] = false
			}
		}
	}

	return signals, nil
}

// OnOrderFilled is called when an order is filled during the backtest
func (s *SimpleMAStrategy) OnOrderFilled(ctx context.Context, order *models.Order) error {
	s.logger.Info("Order filled",
		zap.String("symbol", order.Symbol),
		zap.String("side", string(order.Side)),
		zap.Float64("quantity", order.Quantity),
		zap.Float64("price", order.Price),
		zap.Time("time", order.Time),
	)
	return nil
}

// OnPositionClosed is called when a position is closed during the backtest
func (s *SimpleMAStrategy) OnPositionClosed(ctx context.Context, position *models.ClosedPosition) error {
	s.logger.Info("Position closed",
		zap.String("symbol", position.Symbol),
		zap.String("side", string(position.Side)),
		zap.Float64("entry_price", position.EntryPrice),
		zap.Float64("exit_price", position.ExitPrice),
		zap.Float64("quantity", position.Quantity),
		zap.Float64("profit", position.Profit),
		zap.Time("open_time", position.OpenTime),
		zap.Time("close_time", position.CloseTime),
	)
	return nil
}

// calculateSMA calculates the simple moving average for a given period
func calculateSMA(prices []float64, period int) float64 {
	if len(prices) < period {
		return 0
	}

	sum := 0.0
	for i := len(prices) - period; i < len(prices); i++ {
		sum += prices[i]
	}

	return sum / float64(period)
}
