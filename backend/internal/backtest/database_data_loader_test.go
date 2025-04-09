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

// createTempDBFile creates a temporary SQLite database file and registers its cleanup.
func createTempDBFile(t *testing.T) string {
	tmpFile, err := os.CreateTemp("", "backtest_test_*.db")
	require.NoError(t, err)
	path := tmpFile.Name()
	tmpFile.Close()
	t.Cleanup(func() {
		_ = os.Remove(path)
	})
	return path
}

// setupFullSchema creates the full set of tables (klines, tickers, order_books)
// using SQLite-compatible CREATE TABLE statements.
func setupFullSchema(t *testing.T, db *gorm.DB) {
	queries := []string{
		`CREATE TABLE klines (
			symbol TEXT NOT NULL,
			interval TEXT NOT NULL,
			open_time DATETIME NOT NULL,
			close_time DATETIME NOT NULL,
			open REAL NOT NULL,
			high REAL NOT NULL,
			low REAL NOT NULL,
			close REAL NOT NULL,
			volume REAL NOT NULL,
			is_closed BOOLEAN NOT NULL DEFAULT 0
		)`,
		`CREATE TABLE tickers (
			symbol TEXT NOT NULL,
			price REAL NOT NULL,
			volume REAL NOT NULL,
			price_change REAL NOT NULL,
			price_change_pct REAL NOT NULL,
			high24h REAL NOT NULL,
			low24h REAL NOT NULL,
			timestamp DATETIME NOT NULL
		)`,
		`CREATE TABLE order_books (
			symbol TEXT NOT NULL,
			timestamp DATETIME NOT NULL
		)`,
	}

	for _, query := range queries {
		err := db.Exec(query).Error
		require.NoError(t, err)
	}
}

// setupKlinesSchema creates only the klines table.
func setupKlinesSchema(t *testing.T, db *gorm.DB) {
	query := `CREATE TABLE klines (
		symbol TEXT NOT NULL,
		interval TEXT NOT NULL,
		open_time DATETIME NOT NULL,
		close_time DATETIME NOT NULL,
		open REAL NOT NULL,
		high REAL NOT NULL,
		low REAL NOT NULL,
		close REAL NOT NULL,
		volume REAL NOT NULL,
		is_closed BOOLEAN NOT NULL DEFAULT 0
	)`
	err := db.Exec(query).Error
	require.NoError(t, err)
}

func TestDatabaseDataLoader(t *testing.T) {
	t.Log("Creating temporary database file...")
	dbPath := createTempDBFile(t)

	t.Log("Opening database connection...")
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	require.NoError(t, err)

	t.Log("Setting up schema...")
	setupFullSchema(t, db)

	symbol := "BTCUSDT"
	interval := "1h"
	startTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

	t.Log("Inserting test data...")
	klines := []models.Kline{
		{
			Symbol:    symbol,
			Interval:  interval,
			OpenTime:  startTime,
			CloseTime: startTime.Add(time.Hour),
			Open:      100.0,
			High:      105.0,
			Low:       95.0,
			Close:     102.0,
			Volume:    1000.0,
		},
		{
			Symbol:    symbol,
			Interval:  interval,
			OpenTime:  startTime.Add(time.Hour),
			CloseTime: startTime.Add(2 * time.Hour),
			Open:      102.0,
			High:      107.0,
			Low:       101.0,
			Close:     106.0,
			Volume:    1200.0,
		},
		{
			Symbol:    symbol,
			Interval:  interval,
			OpenTime:  startTime.Add(2 * time.Hour),
			CloseTime: startTime.Add(3 * time.Hour),
			Open:      106.0,
			High:      110.0,
			Low:       104.0,
			Close:     108.0,
			Volume:    1500.0,
		},
	}

	err = db.Create(&klines).Error
	require.NoError(t, err)

	t.Log("Creating DatabaseDataLoader...")
	loader, err := NewDatabaseDataLoader(dbPath)
	require.NoError(t, err)

	t.Log("Loading data...")
	endTime := startTime.Add(3 * time.Hour)
	dataset, err := loader.LoadData(context.Background(), symbol, interval, startTime, endTime)
	require.NoError(t, err)

	t.Log("Verifying dataset...")
	assert.Equal(t, symbol, dataset.Symbol)
	assert.Equal(t, interval, dataset.Interval)
	assert.Len(t, dataset.Klines, 3)

	t.Log("Validating kline records...")
	for i, kline := range dataset.Klines {
		assert.Equal(t, klines[i].Symbol, kline.Symbol)
		assert.Equal(t, klines[i].Interval, kline.Interval)
		assert.Equal(t, klines[i].OpenTime, kline.OpenTime)
		assert.Equal(t, klines[i].CloseTime, kline.CloseTime)
		assert.Equal(t, klines[i].Open, kline.Open)
		assert.Equal(t, klines[i].High, kline.High)
		assert.Equal(t, klines[i].Low, kline.Low)
		assert.Equal(t, klines[i].Close, kline.Close)
		assert.Equal(t, klines[i].Volume, kline.Volume)
	}
	t.Log("Test completed successfully")
}

func TestDatabaseDataLoaderWithOptions(t *testing.T) {
	t.Log("Creating temporary database file...")
	dbPath := createTempDBFile(t)

	t.Log("Opening database connection...")
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	require.NoError(t, err)

	t.Log("Setting up schema...")
	setupFullSchema(t, db)

	symbol := "BTCUSDT"
	interval := "1h"
	startTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

	t.Log("Inserting test data...")
	klines := []models.Kline{
		{
			Symbol:    symbol,
			Interval:  interval,
			OpenTime:  startTime,
			CloseTime: startTime.Add(time.Hour),
			Open:      100.0,
			High:      105.0,
			Low:       95.0,
			Close:     102.0,
			Volume:    1000.0,
		},
		{
			Symbol:    symbol,
			Interval:  interval,
			OpenTime:  startTime.Add(time.Hour),
			CloseTime: startTime.Add(2 * time.Hour),
			Open:      102.0,
			High:      107.0,
			Low:       101.0,
			Close:     106.0,
			Volume:    1200.0,
		},
		// Gap at 2h followed by an outlier at 3h.
		{
			Symbol:    symbol,
			Interval:  interval,
			OpenTime:  startTime.Add(3 * time.Hour),
			CloseTime: startTime.Add(4 * time.Hour),
			Open:      108.0,
			High:      1000.0, // Outlier value.
			Low:       104.0,
			Close:     110.0,
			Volume:    1300.0,
		},
		{
			Symbol:    symbol,
			Interval:  interval,
			OpenTime:  startTime.Add(4 * time.Hour),
			CloseTime: startTime.Add(5 * time.Hour),
			Open:      110.0,
			High:      115.0,
			Low:       108.0,
			Close:     112.0,
			Volume:    1400.0,
		},
	}

	err = db.Create(&klines).Error
	require.NoError(t, err)

	t.Log("Creating DatabaseDataLoader with options...")
	options := &DataLoaderOptions{
		FillMissingValues: true,
		DetectOutliers:    true,
		OutlierThreshold:  3.0,
	}
	loader, err := NewDatabaseDataLoaderWithOptions(dbPath, options)
	require.NoError(t, err)

	t.Log("Loading data...")
	endTime := startTime.Add(5 * time.Hour)
	dataset, err := loader.LoadData(context.Background(), symbol, interval, startTime, endTime)
	require.NoError(t, err)

	t.Log("Verifying dataset...")
	assert.Equal(t, symbol, dataset.Symbol)
	assert.Equal(t, interval, dataset.Interval)
	assert.Len(t, dataset.Klines, 5)

	t.Log("Checking outlier adjustments...")
	for _, kline := range dataset.Klines {
		assert.Less(t, kline.High, 200.0, "High price should be within reasonable bounds after outlier detection")
	}
	t.Log("Test completed successfully")
}

func TestDatabaseDataLoader_LoadData(t *testing.T) {
	t.Log("Opening in-memory database...")
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	require.NoError(t, err)

	t.Log("Setting up schema...")
	setupKlinesSchema(t, db)

	t.Log("Creating test data...")
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

	err = db.Create(&testData).Error
	require.NoError(t, err)

	t.Log("Creating DatabaseDataLoader...")
	loader, err := NewDatabaseDataLoader("file::memory:?cache=shared")
	require.NoError(t, err)

	t.Log("Loading data...")
	dataset, err := loader.LoadData(context.Background(), "BTCUSDT", "1h", time.Now().Add(-3*time.Hour), time.Now())
	require.NoError(t, err)

	t.Log("Verifying dataset...")
	require.NotNil(t, dataset)
	require.Equal(t, 1, len(dataset.Klines))
	require.Equal(t, testData[0].Open, dataset.Klines[0].Open)
	require.Equal(t, testData[0].High, dataset.Klines[0].High)
	require.Equal(t, testData[0].Low, dataset.Klines[0].Low)
	require.Equal(t, testData[0].Close, dataset.Klines[0].Close)
	require.Equal(t, testData[0].Volume, dataset.Klines[0].Volume)
	t.Log("Test completed successfully")
}
