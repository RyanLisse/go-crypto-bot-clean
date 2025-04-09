package backtest

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/ryanlisse/go-crypto-bot/internal/domain/models"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// TestSimpleBacktest tests a simple backtest scenario
func TestSimpleBacktest(t *testing.T) {
	// Create a test strategy
	strategy := &TestStrategy{}

	// Create an in-memory data provider
	dataProvider := NewInMemoryDataProvider()

	// Add test data
	symbol := "BTCUSDT"
	interval := "1h"
	startTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	endTime := time.Date(2023, 1, 10, 0, 0, 0, 0, time.UTC)

	// Create test klines
	klines := createTestKlines(symbol, startTime, endTime, interval)
	dataProvider.AddKlines(symbol, interval, klines)

	// Create backtest config
	config := &BacktestConfig{
		StartTime:          startTime,
		EndTime:            endTime,
		InitialCapital:     10000,
		Symbols:            []string{symbol},
		Interval:           interval,
		CommissionRate:     0.001, // 0.1% commission
		SlippageModel:      &NoSlippage{},
		EnableShortSelling: false,
		DataProvider:       dataProvider,
		Strategy:           strategy,
		Logger:             zap.NewNop(),
	}

	// Create backtest engine
	engine := NewEngine(config)

	// Run backtest
	result, err := engine.Run(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, result)

	// Check results
	assert.Equal(t, config, result.Config)
	assert.Equal(t, startTime, result.StartTime)
	assert.Equal(t, endTime, result.EndTime)
	assert.Equal(t, 10000.0, result.InitialCapital)
	assert.Greater(t, result.FinalCapital, 0.0)
	assert.NotEmpty(t, result.EquityCurve)
	assert.NotEmpty(t, result.DrawdownCurve)
	assert.NotEmpty(t, result.Events)
	assert.NotNil(t, result.PerformanceMetrics)

	// Check that the strategy was called
	assert.True(t, strategy.initializeCalled)
	assert.True(t, strategy.onTickCalled)
}

// TestStrategy is a simple strategy for testing
type TestStrategy struct {
	initializeCalled bool
	onTickCalled     bool
}

// Initialize initializes the strategy
func (s *TestStrategy) Initialize(ctx context.Context, config map[string]interface{}) error {
	s.initializeCalled = true
	return nil
}

// OnTick is called for each new data point
func (s *TestStrategy) OnTick(ctx context.Context, symbol string, timestamp time.Time, data interface{}) ([]*Signal, error) {
	s.onTickCalled = true

	// Generate a buy signal on the first tick
	kline, ok := data.(*models.Kline)
	if !ok {
		return nil, nil
	}

	// Buy at the beginning
	if timestamp.Equal(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)) {
		return []*Signal{
			{
				Symbol:    symbol,
				Side:      "BUY",
				Quantity:  1.0,
				Price:     kline.Close,
				Timestamp: timestamp,
				Reason:    "Test buy signal",
			},
		}, nil
	}

	// Sell at the end
	if timestamp.Equal(time.Date(2023, 1, 9, 0, 0, 0, 0, time.UTC)) {
		return []*Signal{
			{
				Symbol:    symbol,
				Side:      "SELL",
				Quantity:  1.0,
				Price:     kline.Close,
				Timestamp: timestamp,
				Reason:    "Test sell signal",
			},
		}, nil
	}

	return nil, nil
}

// OnOrderFilled is called when an order is filled
func (s *TestStrategy) OnOrderFilled(ctx context.Context, order *models.Order) error {
	return nil
}

// OnPositionClosed is called when a position is closed
func (s *TestStrategy) OnPositionClosed(ctx context.Context, position *models.ClosedPosition) error {
	return nil
}

// createTestKlines creates test klines for backtesting
func createTestKlines(symbol string, startTime, endTime time.Time, interval string) []*models.Kline {
	var klines []*models.Kline
	var intervalDuration time.Duration

	switch interval {
	case "1m":
		intervalDuration = time.Minute
	case "5m":
		intervalDuration = 5 * time.Minute
	case "15m":
		intervalDuration = 15 * time.Minute
	case "1h":
		intervalDuration = time.Hour
	case "4h":
		intervalDuration = 4 * time.Hour
	case "1d":
		intervalDuration = 24 * time.Hour
	default:
		intervalDuration = time.Hour
	}

	// Create klines with a simple price pattern
	price := 20000.0
	for t := startTime; t.Before(endTime); t = t.Add(intervalDuration) {
		// Simple price movement: oscillate between 19000 and 21000
		price = price + (200 * (0.5 - rand.Float64()))
		if price < 19000 {
			price = 19000
		}
		if price > 21000 {
			price = 21000
		}

		kline := &models.Kline{
			Symbol:    symbol,
			Interval:  interval,
			OpenTime:  t,
			CloseTime: t.Add(intervalDuration),
			Open:      price - 50,
			High:      price + 100,
			Low:       price - 100,
			Close:     price,
			Volume:    1000 + rand.Float64()*1000,
			IsClosed:  true,
		}

		klines = append(klines, kline)
	}

	return klines
}
