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
	shortPeriod int
	longPeriod  int
	logger      *zap.Logger
	shortMA     map[string][]float64
	longMA      map[string][]float64
	prices      map[string][]float64
}

// NewSimpleMAStrategy creates a new simple moving average strategy
func NewSimpleMAStrategy(shortPeriod, longPeriod int, logger *zap.Logger) *SimpleMAStrategy {
	return &SimpleMAStrategy{
		shortPeriod: shortPeriod,
		longPeriod:  longPeriod,
		logger:      logger,
		shortMA:     make(map[string][]float64),
		longMA:      make(map[string][]float64),
		prices:      make(map[string][]float64),
	}
}

// Initialize implements the BacktestStrategy interface
func (s *SimpleMAStrategy) Initialize(ctx context.Context, config interface{}) error {
	return nil // No initialization needed
}

// OnTick implements the BacktestStrategy interface
func (s *SimpleMAStrategy) OnTick(ctx context.Context, symbol string, timestamp time.Time, data interface{}) ([]*backtest.Signal, error) {
	kline, ok := data.(*models.Kline)
	if !ok {
		return nil, nil
	}

	// Initialize price arrays if needed
	if _, ok := s.prices[symbol]; !ok {
		s.prices[symbol] = make([]float64, 0)
		s.shortMA[symbol] = make([]float64, 0)
		s.longMA[symbol] = make([]float64, 0)
	}

	// Add new price
	s.prices[symbol] = append(s.prices[symbol], kline.Close)

	// Calculate moving averages
	if len(s.prices[symbol]) >= s.longPeriod {
		shortMA := calculateSMA(s.prices[symbol], s.shortPeriod)
		longMA := calculateSMA(s.prices[symbol], s.longPeriod)

		s.shortMA[symbol] = append(s.shortMA[symbol], shortMA)
		s.longMA[symbol] = append(s.longMA[symbol], longMA)

		// Generate signals based on MA crossover
		if len(s.shortMA[symbol]) >= 2 && len(s.longMA[symbol]) >= 2 {
			prevShortMA := s.shortMA[symbol][len(s.shortMA[symbol])-2]
			prevLongMA := s.longMA[symbol][len(s.longMA[symbol])-2]

			// Buy signal: short MA crosses above long MA
			if prevShortMA <= prevLongMA && shortMA > longMA {
				s.logger.Info("Buy signal generated",
					zap.String("symbol", symbol),
					zap.Float64("price", kline.Close),
					zap.Float64("shortMA", shortMA),
					zap.Float64("longMA", longMA),
				)
				return []*backtest.Signal{
					{
						Symbol:    symbol,
						Side:      string(models.OrderSideBuy),
						Price:     kline.Close,
						Quantity:  1.0, // Quantity will be adjusted by position sizing
						Timestamp: timestamp,
						Reason:    "MA Crossover - Buy Signal",
					},
				}, nil
			}

			// Sell signal: short MA crosses below long MA
			if prevShortMA >= prevLongMA && shortMA < longMA {
				s.logger.Info("Sell signal generated",
					zap.String("symbol", symbol),
					zap.Float64("price", kline.Close),
					zap.Float64("shortMA", shortMA),
					zap.Float64("longMA", longMA),
				)
				return []*backtest.Signal{
					{
						Symbol:    symbol,
						Side:      string(models.OrderSideSell),
						Price:     kline.Close,
						Quantity:  1.0, // Full position will be closed by the engine
						Timestamp: timestamp,
						Reason:    "MA Crossover - Sell Signal",
					},
				}, nil
			}
		}
	}

	return nil, nil
}

// OnOrderFilled implements the BacktestStrategy interface
func (s *SimpleMAStrategy) OnOrderFilled(ctx context.Context, order *models.Order) error {
	return nil // No special handling needed
}

// ClosePositions implements the BacktestStrategy interface
func (s *SimpleMAStrategy) ClosePositions(ctx context.Context) ([]*backtest.Signal, error) {
	return nil, nil // No special handling needed
}

// calculateSMA calculates the Simple Moving Average
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
