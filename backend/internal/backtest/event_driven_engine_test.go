package backtest

import (
	"context"
	"os"
	"testing"
	"time"

	"go-crypto-bot-clean/backend/internal/domain/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TestEventDrivenEngine tests the event-driven backtesting engine
func TestEventDrivenEngine(t *testing.T) {
	// Create a temporary database file
	dbFile, err := os.CreateTemp("", "backtest_test_*.db")
	require.NoError(t, err)
	dbPath := dbFile.Name()
	dbFile.Close()
	defer os.Remove(dbPath)

	// Connect to the database
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	require.NoError(t, err)

	// Create tables
	err = db.AutoMigrate(
		&models.Candle{},
		&models.Ticker{},
		&models.OrderBook{},
		&models.CandleOrderBookEntry{},
		&models.BacktestResult{},
		&models.Position{},
		&models.ClosedPosition{},
		&models.Order{},
	)
	require.NoError(t, err)

	// Create test data
	symbol := "BTCUSDT"
	interval := "1h"
	startTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	endTime := time.Date(2023, 1, 10, 0, 0, 0, 0, time.UTC)

	// Create test candles
	candles := createTestCandles(db, symbol, startTime, endTime, interval)
	require.NotEmpty(t, candles)

	// Create a test strategy
	strategy := &TestEventDrivenStrategy{}

	// Create a database data loader
	dataLoader, err := NewDatabaseDataLoader(dbPath)
	require.NoError(t, err)

	// Create event-driven engine config
	config := &EventDrivenEngineConfig{
		StartTime:          startTime,
		EndTime:            endTime,
		InitialCapital:     10000,
		Symbols:            []string{symbol},
		Interval:           interval,
		FeeModel:           NewFixedFeeModel(0.001), // 0.1% fee
		SlippageModel:      &NoSlippage{},
		EnableShortSelling: false,
		DataLoader:         dataLoader,
		Strategy:           strategy,
		DB:                 db,
		Logger:             zap.NewNop(),
	}

	// Create event-driven engine
	engine := NewEventDrivenEngine(config)

	// Run backtest
	result, err := engine.Run(context.Background())
	require.NoError(t, err)
	require.NotNil(t, result)

	// Check results
	assert.Equal(t, startTime, result.StartTime)
	assert.Equal(t, endTime, result.EndTime)
	assert.Equal(t, 10000.0, result.InitialCapital)
	assert.Greater(t, result.FinalCapital, 0.0)
	assert.NotEmpty(t, result.EquityCurve)
	assert.NotEmpty(t, result.DrawdownCurve)
	assert.NotEmpty(t, result.Events)

	// Check that the strategy was called
	assert.True(t, strategy.initializeCalled)
	assert.True(t, strategy.onTickCalled)
	assert.True(t, strategy.onMarketEventCalled)

	// Check that the results were saved to the database
	var savedResults []models.BacktestResult
	err = db.Find(&savedResults).Error
	require.NoError(t, err)
	assert.NotEmpty(t, savedResults)
}

// TestEventDrivenStrategy is a simple strategy for testing the event-driven engine
type TestEventDrivenStrategy struct {
	initializeCalled    bool
	onTickCalled        bool
	onMarketEventCalled bool
	onOrderFilledCalled bool
}

// Initialize initializes the strategy
func (s *TestEventDrivenStrategy) Initialize(ctx context.Context, config interface{}) error {
	s.initializeCalled = true
	// Can add config parsing here if needed
	return nil
}

// OnTick is called for each new data point
func (s *TestEventDrivenStrategy) OnTick(ctx context.Context, symbol string, timestamp time.Time, data interface{}) ([]*Signal, error) {
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

// OnMarketEvent is called for each market event
func (s *TestEventDrivenStrategy) OnMarketEvent(ctx context.Context, event *MarketEvent) ([]*Signal, error) {
	s.onMarketEventCalled = true
	return nil, nil
}

// OnOrderFilled is called when an order is filled
func (s *TestEventDrivenStrategy) OnOrderFilled(ctx context.Context, order *models.Order) error {
	s.onOrderFilledCalled = true
	return nil
}

// ClosePositions implements the BacktestStrategy interface
func (s *TestEventDrivenStrategy) ClosePositions(ctx context.Context) ([]*Signal, error) {
	// Test strategy doesn't need to do anything specific on close
	return nil, nil
}

// OnPositionClosed is called when a position is closed
// Note: This method is NOT part of the BacktestStrategy or EventDrivenStrategy interface
// It was likely left over from a previous design.
// func (s *TestEventDrivenStrategy) OnPositionClosed(ctx context.Context, position *models.ClosedPosition) error {
// 	return nil
// }

// createTestCandles creates test candles in the database
func createTestCandles(db *gorm.DB, symbol string, startTime, endTime time.Time, interval string) []models.Candle {
	var candles []models.Candle
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

	// Create candles with a simple price pattern
	price := 20000.0
	for t := startTime; t.Before(endTime); t = t.Add(intervalDuration) {
		// Simple price movement: oscillate between 19000 and 21000
		price = price + (200 * (0.5 - float64(t.Unix()%100)/100))
		if price < 19000 {
			price = 19000
		}
		if price > 21000 {
			price = 21000
		}

		candle := models.Candle{
			Symbol:     symbol,
			Interval:   interval,
			OpenTime:   t,
			CloseTime:  t.Add(intervalDuration),
			OpenPrice:  price - 50,
			HighPrice:  price + 100,
			LowPrice:   price - 100,
			ClosePrice: price,
			Volume:     1000 + float64(t.Unix()%1000),
		}

		candles = append(candles, candle)
	}

	// Insert candles into the database
	db.Create(&candles)

	return candles
}
