package gorm

import (
	"context"
	"testing"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model/market"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupTestDB(t *testing.T) (*gorm.DB, func()) {
	// Use in-memory SQLite database for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)

	// Run migrations
	err = db.AutoMigrate(
		&TickerEntity{},
		&CandleEntity{},
		&OrderBookEntity{},
		&OrderBookEntryEntity{},
		&SymbolEntity{},
	)
	require.NoError(t, err)

	// No file to clean up, just close the DB connection
	cleanup := func() {
		sqlDB, err := db.DB()
		if err == nil {
			sqlDB.Close()
		}
	}

	return db, cleanup
}

func setupTestRepository(t *testing.T) (*MarketRepository, func()) {
	db, cleanup := setupTestDB(t)

	// Create a logger
	logger := zerolog.New(zerolog.NewConsoleWriter()).With().Timestamp().Logger()

	// Create the repository
	repo := NewMarketRepository(db, &logger)

	return repo, cleanup
}

func TestSaveAndGetTicker(t *testing.T) {
	repo, cleanup := setupTestRepository(t)
	defer cleanup()

	ctx := context.Background()

	// Create a test ticker
	ticker := &market.Ticker{
		ID:            "test-ticker-1",
		Symbol:        "BTCUSDT",
		Exchange:      "mexc",
		Price:         50000.0,
		Volume:        100.0,
		High24h:       51000.0,
		Low24h:        49000.0,
		PriceChange:   1000.0,
		PercentChange: 2.0,
		LastUpdated:   time.Now().Round(time.Millisecond), // Round to avoid precision issues
	}

	// Save the ticker
	err := repo.SaveTicker(ctx, ticker)
	require.NoError(t, err)

	// Retrieve the ticker
	retrievedTicker, err := repo.GetTicker(ctx, "BTCUSDT", "mexc")
	require.NoError(t, err)

	// Verify the ticker was saved correctly
	assert.Equal(t, ticker.ID, retrievedTicker.ID)
	assert.Equal(t, ticker.Symbol, retrievedTicker.Symbol)
	assert.Equal(t, ticker.Exchange, retrievedTicker.Exchange)
	assert.Equal(t, ticker.Price, retrievedTicker.Price)
	assert.Equal(t, ticker.Volume, retrievedTicker.Volume)
	assert.Equal(t, ticker.High24h, retrievedTicker.High24h)
	assert.Equal(t, ticker.Low24h, retrievedTicker.Low24h)
	assert.Equal(t, ticker.PriceChange, retrievedTicker.PriceChange)
	assert.Equal(t, ticker.PercentChange, retrievedTicker.PercentChange)
	assert.Equal(t, ticker.LastUpdated.Unix(), retrievedTicker.LastUpdated.Unix())
}

func TestGetAllTickers(t *testing.T) {
	repo, cleanup := setupTestRepository(t)
	defer cleanup()

	ctx := context.Background()

	// Create test tickers
	ticker1 := &market.Ticker{
		ID:          "test-ticker-1",
		Symbol:      "BTCUSDT",
		Exchange:    "mexc",
		Price:       50000.0,
		LastUpdated: time.Now().Round(time.Millisecond),
	}

	ticker2 := &market.Ticker{
		ID:          "test-ticker-2",
		Symbol:      "ETHUSDT",
		Exchange:    "mexc",
		Price:       3000.0,
		LastUpdated: time.Now().Round(time.Millisecond),
	}

	// Save the tickers
	err := repo.SaveTicker(ctx, ticker1)
	require.NoError(t, err)

	err = repo.SaveTicker(ctx, ticker2)
	require.NoError(t, err)

	// Retrieve all tickers for the exchange
	tickers, err := repo.GetAllTickers(ctx, "mexc")
	require.NoError(t, err)

	// Verify the tickers were retrieved correctly
	assert.Equal(t, 2, len(tickers))

	// Verify the ticker symbols
	symbols := []string{tickers[0].Symbol, tickers[1].Symbol}
	assert.Contains(t, symbols, "BTCUSDT")
	assert.Contains(t, symbols, "ETHUSDT")
}

func TestSaveAndGetCandle(t *testing.T) {
	repo, cleanup := setupTestRepository(t)
	defer cleanup()

	ctx := context.Background()

	// Create a test candle
	now := time.Now().Round(time.Millisecond)
	candle := &market.Candle{
		Symbol:      "BTCUSDT",
		Exchange:    "mexc",
		Interval:    market.Interval1h,
		OpenTime:    now,
		CloseTime:   now.Add(1 * time.Hour),
		Open:        50000.0,
		High:        51000.0,
		Low:         49000.0,
		Close:       50500.0,
		Volume:      100.0,
		QuoteVolume: 5000000.0,
		TradeCount:  1000,
		Complete:    true,
	}

	// Save the candle
	err := repo.SaveCandle(ctx, candle)
	require.NoError(t, err)

	// Retrieve the candle
	retrievedCandle, err := repo.GetCandle(ctx, "BTCUSDT", "mexc", market.Interval1h, now)
	require.NoError(t, err)

	// Verify the candle was saved correctly
	assert.Equal(t, candle.Symbol, retrievedCandle.Symbol)
	assert.Equal(t, candle.Exchange, retrievedCandle.Exchange)
	assert.Equal(t, candle.Interval, retrievedCandle.Interval)
	assert.Equal(t, candle.OpenTime.Unix(), retrievedCandle.OpenTime.Unix())
	assert.Equal(t, candle.CloseTime.Unix(), retrievedCandle.CloseTime.Unix())
	assert.Equal(t, candle.Open, retrievedCandle.Open)
	assert.Equal(t, candle.High, retrievedCandle.High)
	assert.Equal(t, candle.Low, retrievedCandle.Low)
	assert.Equal(t, candle.Close, retrievedCandle.Close)
	assert.Equal(t, candle.Volume, retrievedCandle.Volume)
	assert.Equal(t, candle.QuoteVolume, retrievedCandle.QuoteVolume)
	assert.Equal(t, candle.TradeCount, retrievedCandle.TradeCount)
	assert.Equal(t, candle.Complete, retrievedCandle.Complete)
}

func TestGetCandles(t *testing.T) {
	repo, cleanup := setupTestRepository(t)
	defer cleanup()

	ctx := context.Background()

	// Create test candles
	now := time.Now().Round(time.Millisecond)
	candle1 := &market.Candle{
		Symbol:    "BTCUSDT",
		Exchange:  "mexc",
		Interval:  market.Interval1h,
		OpenTime:  now.Add(-2 * time.Hour),
		CloseTime: now.Add(-1 * time.Hour),
		Open:      49000.0,
		High:      50000.0,
		Low:       48000.0,
		Close:     49500.0,
		Volume:    90.0,
		Complete:  true,
	}

	candle2 := &market.Candle{
		Symbol:    "BTCUSDT",
		Exchange:  "mexc",
		Interval:  market.Interval1h,
		OpenTime:  now.Add(-1 * time.Hour),
		CloseTime: now,
		Open:      49500.0,
		High:      51000.0,
		Low:       49000.0,
		Close:     50500.0,
		Volume:    100.0,
		Complete:  true,
	}

	// Save the candles
	err := repo.SaveCandle(ctx, candle1)
	require.NoError(t, err)

	err = repo.SaveCandle(ctx, candle2)
	require.NoError(t, err)

	// Retrieve candles within a time range
	start := now.Add(-3 * time.Hour)
	end := now.Add(1 * time.Hour)
	candles, err := repo.GetCandles(ctx, "BTCUSDT", "mexc", market.Interval1h, start, end, 10)
	require.NoError(t, err)

	// Verify the candles were retrieved correctly
	assert.Equal(t, 2, len(candles))
	assert.Equal(t, candle1.OpenTime.Unix(), candles[0].OpenTime.Unix())
	assert.Equal(t, candle2.OpenTime.Unix(), candles[1].OpenTime.Unix())
}

func TestGetLatestCandle(t *testing.T) {
	repo, cleanup := setupTestRepository(t)
	defer cleanup()

	ctx := context.Background()

	// Create test candles with different times
	now := time.Now().Round(time.Millisecond)
	candle1 := &market.Candle{
		Symbol:    "BTCUSDT",
		Exchange:  "mexc",
		Interval:  market.Interval1h,
		OpenTime:  now.Add(-2 * time.Hour),
		CloseTime: now.Add(-1 * time.Hour),
		Open:      49000.0,
		Close:     49500.0,
		Complete:  true,
	}

	candle2 := &market.Candle{
		Symbol:    "BTCUSDT",
		Exchange:  "mexc",
		Interval:  market.Interval1h,
		OpenTime:  now.Add(-1 * time.Hour),
		CloseTime: now,
		Open:      49500.0,
		Close:     50500.0,
		Complete:  true,
	}

	// Save the candles
	err := repo.SaveCandle(ctx, candle1)
	require.NoError(t, err)

	err = repo.SaveCandle(ctx, candle2)
	require.NoError(t, err)

	// Retrieve the latest candle
	latestCandle, err := repo.GetLatestCandle(ctx, "BTCUSDT", "mexc", market.Interval1h)
	require.NoError(t, err)

	// Verify the latest candle was returned
	assert.Equal(t, candle2.OpenTime.Unix(), latestCandle.OpenTime.Unix())
	assert.Equal(t, candle2.Close, latestCandle.Close)
}

func TestSaveAndGetSymbol(t *testing.T) {
	repo, cleanup := setupTestRepository(t)
	defer cleanup()

	ctx := context.Background()

	// Create a test symbol
	symbol := &market.Symbol{
		Symbol:            "BTCUSDT",
		BaseAsset:         "BTC",
		QuoteAsset:        "USDT",
		Exchange:          "mexc",
		Status:            "TRADING",
		MinPrice:          0.01,
		MaxPrice:          100000.0,
		PricePrecision:    2,
		MinQty:            0.0001,
		MaxQty:            1000.0,
		QtyPrecision:      4,
		AllowedOrderTypes: []string{"LIMIT", "MARKET"},
	}

	// Save the symbol
	err := repo.Create(ctx, symbol)
	require.NoError(t, err)

	// Retrieve the symbol
	retrievedSymbol, err := repo.GetBySymbol(ctx, "BTCUSDT")
	require.NoError(t, err)

	// Verify the symbol was saved correctly
	assert.Equal(t, symbol.Symbol, retrievedSymbol.Symbol)
	assert.Equal(t, symbol.BaseAsset, retrievedSymbol.BaseAsset)
	assert.Equal(t, symbol.QuoteAsset, retrievedSymbol.QuoteAsset)
	assert.Equal(t, symbol.Exchange, retrievedSymbol.Exchange)
	assert.Equal(t, symbol.Status, retrievedSymbol.Status)
	assert.Equal(t, symbol.MinPrice, retrievedSymbol.MinPrice)
	assert.Equal(t, symbol.MaxPrice, retrievedSymbol.MaxPrice)
	assert.Equal(t, symbol.PricePrecision, retrievedSymbol.PricePrecision)
	assert.Equal(t, symbol.MinQty, retrievedSymbol.MinQty)
	assert.Equal(t, symbol.MaxQty, retrievedSymbol.MaxQty)
	assert.Equal(t, symbol.QtyPrecision, retrievedSymbol.QtyPrecision)
	assert.ElementsMatch(t, symbol.AllowedOrderTypes, retrievedSymbol.AllowedOrderTypes)
}

func TestGetByExchange(t *testing.T) {
	repo, cleanup := setupTestRepository(t)
	defer cleanup()

	ctx := context.Background()

	// Create test symbols
	symbol1 := &market.Symbol{
		Symbol:    "BTCUSDT",
		BaseAsset: "BTC",
		Exchange:  "mexc",
		Status:    "TRADING",
	}

	symbol2 := &market.Symbol{
		Symbol:    "ETHUSDT",
		BaseAsset: "ETH",
		Exchange:  "mexc",
		Status:    "TRADING",
	}

	symbol3 := &market.Symbol{
		Symbol:    "BTCUSDT",
		BaseAsset: "BTC",
		Exchange:  "binance",
		Status:    "TRADING",
	}

	// Save the symbols
	err := repo.Create(ctx, symbol1)
	require.NoError(t, err)

	err = repo.Create(ctx, symbol2)
	require.NoError(t, err)

	err = repo.Create(ctx, symbol3)
	require.NoError(t, err)

	// Retrieve symbols by exchange
	mexcSymbols, err := repo.GetByExchange(ctx, "mexc")
	require.NoError(t, err)

	// Verify the correct symbols were returned
	assert.Equal(t, 2, len(mexcSymbols))

	// Verify the symbol details
	symbols := []string{mexcSymbols[0].Symbol, mexcSymbols[1].Symbol}
	assert.Contains(t, symbols, "BTCUSDT")
	assert.Contains(t, symbols, "ETHUSDT")

	// Verify all symbols have the correct exchange
	for _, s := range mexcSymbols {
		assert.Equal(t, "mexc", s.Exchange)
	}
}

func TestGetAll(t *testing.T) {
	repo, cleanup := setupTestRepository(t)
	defer cleanup()

	ctx := context.Background()

	// Create test symbols
	symbol1 := &market.Symbol{
		Symbol:    "BTCUSDT",
		BaseAsset: "BTC",
		Exchange:  "mexc",
		Status:    "TRADING",
	}

	symbol2 := &market.Symbol{
		Symbol:    "ETHUSDT",
		BaseAsset: "ETH",
		Exchange:  "mexc",
		Status:    "TRADING",
	}

	// Save the symbols
	err := repo.Create(ctx, symbol1)
	require.NoError(t, err)

	err = repo.Create(ctx, symbol2)
	require.NoError(t, err)

	// Retrieve all symbols
	symbols, err := repo.GetAll(ctx)
	require.NoError(t, err)

	// Verify all symbols were returned
	assert.Equal(t, 2, len(symbols))

	// Verify the symbol details
	symbolNames := []string{symbols[0].Symbol, symbols[1].Symbol}
	assert.Contains(t, symbolNames, "BTCUSDT")
	assert.Contains(t, symbolNames, "ETHUSDT")
}

func TestUpdateSymbol(t *testing.T) {
	repo, cleanup := setupTestRepository(t)
	defer cleanup()

	ctx := context.Background()

	// Create a test symbol
	symbol := &market.Symbol{
		Symbol:    "BTCUSDT",
		BaseAsset: "BTC",
		Exchange:  "mexc",
		Status:    "TRADING",
	}

	// Save the symbol
	err := repo.Create(ctx, symbol)
	require.NoError(t, err)

	// Update the symbol
	symbol.Status = "BREAK"
	err = repo.Update(ctx, symbol)
	require.NoError(t, err)

	// Retrieve the updated symbol
	updatedSymbol, err := repo.GetBySymbol(ctx, "BTCUSDT")
	require.NoError(t, err)

	// Verify the symbol was updated correctly
	assert.Equal(t, "BREAK", updatedSymbol.Status)
}

func TestDeleteSymbol(t *testing.T) {
	repo, cleanup := setupTestRepository(t)
	defer cleanup()

	ctx := context.Background()

	// Create a test symbol
	symbol := &market.Symbol{
		Symbol:    "BTCUSDT",
		BaseAsset: "BTC",
		Exchange:  "mexc",
		Status:    "TRADING",
	}

	// Save the symbol
	err := repo.Create(ctx, symbol)
	require.NoError(t, err)

	// Delete the symbol
	err = repo.Delete(ctx, "BTCUSDT")
	require.NoError(t, err)

	// Try to retrieve the deleted symbol
	_, err = repo.GetBySymbol(ctx, "BTCUSDT")
	assert.Error(t, err)
}

func TestPurgeOldData(t *testing.T) {
	repo, cleanup := setupTestRepository(t)
	defer cleanup()

	ctx := context.Background()

	// Create test data with different timestamps
	oldTime := time.Now().Add(-24 * time.Hour).Round(time.Millisecond)
	newTime := time.Now().Round(time.Millisecond)

	// Create old ticker
	oldTicker := &market.Ticker{
		ID:          "old-ticker",
		Symbol:      "BTCUSDT",
		Exchange:    "mexc",
		Price:       48000.0,
		LastUpdated: oldTime,
	}

	// Create new ticker
	newTicker := &market.Ticker{
		ID:          "new-ticker",
		Symbol:      "BTCUSDT",
		Exchange:    "mexc",
		Price:       50000.0,
		LastUpdated: newTime,
	}

	// Save the tickers
	err := repo.SaveTicker(ctx, oldTicker)
	require.NoError(t, err)

	err = repo.SaveTicker(ctx, newTicker)
	require.NoError(t, err)

	// Create old candle
	oldCandle := &market.Candle{
		Symbol:    "BTCUSDT",
		Exchange:  "mexc",
		Interval:  market.Interval1h,
		OpenTime:  oldTime,
		CloseTime: oldTime.Add(1 * time.Hour),
		Open:      48000.0,
		Close:     48500.0,
		Complete:  true,
	}

	// Create new candle
	newCandle := &market.Candle{
		Symbol:    "BTCUSDT",
		Exchange:  "mexc",
		Interval:  market.Interval1h,
		OpenTime:  newTime,
		CloseTime: newTime.Add(1 * time.Hour),
		Open:      50000.0,
		Close:     50500.0,
		Complete:  true,
	}

	// Save the candles
	err = repo.SaveCandle(ctx, oldCandle)
	require.NoError(t, err)

	err = repo.SaveCandle(ctx, newCandle)
	require.NoError(t, err)

	// Purge data older than 12 hours
	purgeTime := time.Now().Add(-12 * time.Hour)
	err = repo.PurgeOldData(ctx, purgeTime)
	require.NoError(t, err)

	// Verify old data was purged
	tickers, err := repo.GetTickerHistory(ctx, "BTCUSDT", "mexc", oldTime, newTime.Add(1*time.Hour))
	require.NoError(t, err)
	assert.Equal(t, 1, len(tickers))
	assert.Equal(t, "new-ticker", tickers[0].ID)

	candles, err := repo.GetCandles(ctx, "BTCUSDT", "mexc", market.Interval1h, oldTime, newTime.Add(1*time.Hour), 10)
	require.NoError(t, err)
	assert.Equal(t, 1, len(candles))
	assert.Equal(t, newTime.Unix(), candles[0].OpenTime.Unix())
}
