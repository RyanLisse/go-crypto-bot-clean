package backtest

import (
	"context"
	"os"
	"testing"
	"time"

	"go-crypto-bot-clean/backend/internal/domain/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestDatabaseDataLoader(t *testing.T) {
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
	err = db.AutoMigrate(&models.Kline{}, &models.Ticker{}, &models.OrderBook{}, &models.CandleOrderBookEntry{})
	require.NoError(t, err)

	// Insert test data
	symbol := "BTCUSDT"
	interval := "1h"
	startTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

	// Insert klines
	klines := []models.Kline{
		{
			Symbol:     symbol,
			Interval:   interval,
			OpenTime:   startTime,
			CloseTime:  startTime.Add(time.Hour),
			OpenPrice:  100.0,
			HighPrice:  105.0,
			LowPrice:   95.0,
			ClosePrice: 102.0,
			Volume:     1000.0,
		},
		{
			Symbol:     symbol,
			Interval:   interval,
			OpenTime:   startTime.Add(time.Hour),
			CloseTime:  startTime.Add(2 * time.Hour),
			OpenPrice:  102.0,
			HighPrice:  107.0,
			LowPrice:   101.0,
			ClosePrice: 106.0,
			Volume:     1200.0,
		},
		{
			Symbol:     symbol,
			Interval:   interval,
			OpenTime:   startTime.Add(2 * time.Hour),
			CloseTime:  startTime.Add(3 * time.Hour),
			OpenPrice:  106.0,
			HighPrice:  110.0,
			LowPrice:   104.0,
			ClosePrice: 108.0,
			Volume:     1500.0,
		},
	}

	err = db.Create(&klines).Error
	require.NoError(t, err)

	// Create DatabaseDataLoader
	loader, err := NewDatabaseDataLoader(dbPath)
	require.NoError(t, err)

	// Test loading data
	endTime := startTime.Add(3 * time.Hour)

	dataset, err := loader.LoadData(context.Background(), symbol, interval, startTime, endTime)
	require.NoError(t, err)

	// Verify dataset
	assert.Equal(t, symbol, dataset.Symbol)
	assert.Equal(t, interval, dataset.Interval)
	assert.Len(t, dataset.Klines, 3)

	// Verify kline data
	for i, kline := range dataset.Klines {
		assert.Equal(t, klines[i].Symbol, kline.Symbol)
		assert.Equal(t, klines[i].Interval, kline.Interval)
		assert.Equal(t, klines[i].OpenTime, kline.OpenTime)
		assert.Equal(t, klines[i].CloseTime, kline.CloseTime)
		assert.Equal(t, klines[i].OpenPrice, kline.OpenPrice)
		assert.Equal(t, klines[i].HighPrice, kline.HighPrice)
		assert.Equal(t, klines[i].LowPrice, kline.LowPrice)
		assert.Equal(t, klines[i].ClosePrice, kline.ClosePrice)
		assert.Equal(t, klines[i].Volume, kline.Volume)
	}
}

func TestDatabaseDataLoaderWithOptions(t *testing.T) {
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
	err = db.AutoMigrate(&models.Kline{}, &models.Ticker{}, &models.OrderBook{}, &models.CandleOrderBookEntry{})
	require.NoError(t, err)

	// Insert test data with missing values and outliers
	symbol := "BTCUSDT"
	interval := "1h"
	startTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

	// Insert klines with a gap and an outlier
	klines := []models.Kline{
		{
			Symbol:     symbol,
			Interval:   interval,
			OpenTime:   startTime,
			CloseTime:  startTime.Add(time.Hour),
			OpenPrice:  100.0,
			HighPrice:  105.0,
			LowPrice:   95.0,
			ClosePrice: 102.0,
			Volume:     1000.0,
		},
		{
			Symbol:     symbol,
			Interval:   interval,
			OpenTime:   startTime.Add(time.Hour),
			CloseTime:  startTime.Add(2 * time.Hour),
			OpenPrice:  102.0,
			HighPrice:  107.0,
			LowPrice:   101.0,
			ClosePrice: 106.0,
			Volume:     1200.0,
		},
		// Gap at 2h
		{
			Symbol:     symbol,
			Interval:   interval,
			OpenTime:   startTime.Add(3 * time.Hour),
			CloseTime:  startTime.Add(4 * time.Hour),
			OpenPrice:  108.0,
			HighPrice:  1000.0, // Outlier
			LowPrice:   104.0,
			ClosePrice: 110.0,
			Volume:     1300.0,
		},
		{
			Symbol:     symbol,
			Interval:   interval,
			OpenTime:   startTime.Add(4 * time.Hour),
			CloseTime:  startTime.Add(5 * time.Hour),
			OpenPrice:  110.0,
			HighPrice:  115.0,
			LowPrice:   108.0,
			ClosePrice: 112.0,
			Volume:     1400.0,
		},
	}

	err = db.Create(&klines).Error
	require.NoError(t, err)

	// Create DatabaseDataLoader with preprocessing options
	options := &DataLoaderOptions{
		FillMissingValues: true,
		DetectOutliers:    true,
		OutlierThreshold:  3.0,
	}

	loader, err := NewDatabaseDataLoaderWithOptions(dbPath, options)
	require.NoError(t, err)

	// Test loading data with preprocessing
	endTime := startTime.Add(5 * time.Hour)

	dataset, err := loader.LoadData(context.Background(), symbol, interval, startTime, endTime)
	require.NoError(t, err)

	// Verify dataset
	assert.Equal(t, symbol, dataset.Symbol)
	assert.Equal(t, interval, dataset.Interval)
	assert.Len(t, dataset.Klines, 5) // With fill_missing_values=true, we should have 5 records (0h, 1h, 2h, 3h, 4h)

	// Verify that the outlier was fixed
	for _, kline := range dataset.Klines {
		assert.Less(t, kline.HighPrice, 200.0, "High price should be within reasonable bounds after outlier detection")
	}
}

func TestDatabaseDataLoader_LoadData(t *testing.T) {
	// Create a temporary database
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	require.NoError(t, err)

	// Migrate the schema
	err = db.AutoMigrate(&models.Kline{})
	require.NoError(t, err)

	// Create test data
	testData := []*models.Kline{
		{
			Symbol:    "BTCUSDT",
			Interval:  "1h",
			OpenTime:  time.Now().Add(-2 * time.Hour),
			CloseTime: time.Now().Add(-1 * time.Hour),
			Open:      50000,
			High:      51000,
			Low:       49000,
			Close:     50500,
			Volume:    100,
			IsClosed:  true,
		},
	}

	// Insert test data
	err = db.Create(&testData).Error
	require.NoError(t, err)

	// Create the data loader
	loader := NewDatabaseDataLoader(db)

	// Load the data
	dataset, err := loader.LoadData(context.Background(), "BTCUSDT", "1h", time.Now().Add(-3*time.Hour), time.Now())
	require.NoError(t, err)

	// Verify the loaded data
	require.Equal(t, 1, len(dataset))
	require.Equal(t, testData[0].Open, dataset[0].Open)
	require.Equal(t, testData[0].High, dataset[0].High)
	require.Equal(t, testData[0].Low, dataset[0].Low)
	require.Equal(t, testData[0].Close, dataset[0].Close)
	require.Equal(t, testData[0].Volume, dataset[0].Volume)
}
