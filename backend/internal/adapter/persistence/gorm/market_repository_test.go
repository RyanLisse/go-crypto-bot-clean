package gorm

import (
	"context"
	"testing"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
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

	// Convert market.Ticker to model.Ticker
	modelTicker := &model.Ticker{
		ID:                 ticker.ID,
		Symbol:             ticker.Symbol,
		Exchange:           ticker.Exchange,
		LastPrice:          ticker.Price,
		Volume:             ticker.Volume,
		HighPrice:          ticker.High24h,
		LowPrice:           ticker.Low24h,
		PriceChange:        ticker.PriceChange,
		PriceChangePercent: ticker.PercentChange,
		Timestamp:          ticker.LastUpdated,
	}

	// Save the ticker
	err := repo.SaveTicker(ctx, modelTicker)
	require.NoError(t, err)

	// Retrieve the ticker
	retrievedTicker, err := repo.GetTicker(ctx, "BTCUSDT", "mexc")
	require.NoError(t, err)

	// Verify the ticker was saved correctly
	assert.Equal(t, ticker.ID, retrievedTicker.ID)
	assert.Equal(t, ticker.Symbol, retrievedTicker.Symbol)
	assert.Equal(t, ticker.Exchange, retrievedTicker.Exchange)
	assert.Equal(t, ticker.Price, retrievedTicker.LastPrice)
	assert.Equal(t, ticker.Volume, retrievedTicker.Volume)
	assert.Equal(t, ticker.High24h, retrievedTicker.HighPrice)
	assert.Equal(t, ticker.Low24h, retrievedTicker.LowPrice)
	assert.Equal(t, ticker.PriceChange, retrievedTicker.PriceChange)
	assert.Equal(t, ticker.PercentChange, retrievedTicker.PriceChangePercent)
	assert.Equal(t, ticker.LastUpdated.Unix(), retrievedTicker.Timestamp.Unix())
}

func TestGetAllTickers(t *testing.T) {
	repo, cleanup := setupTestRepository(t)
	defer cleanup()

	ctx := context.Background()

	// Create test tickers
	marketTicker1 := &market.Ticker{
		ID:          "test-ticker-1",
		Symbol:      "BTCUSDT",
		Exchange:    "mexc",
		Price:       50000.0,
		LastUpdated: time.Now().Round(time.Millisecond),
	}

	marketTicker2 := &market.Ticker{
		ID:          "test-ticker-2",
		Symbol:      "ETHUSDT",
		Exchange:    "mexc",
		Price:       3000.0,
		LastUpdated: time.Now().Round(time.Millisecond),
	}

	// Convert to model.Ticker
	ticker1 := &model.Ticker{
		ID:        marketTicker1.ID,
		Symbol:    marketTicker1.Symbol,
		Exchange:  marketTicker1.Exchange,
		LastPrice: marketTicker1.Price,
		Timestamp: marketTicker1.LastUpdated,
	}

	ticker2 := &model.Ticker{
		ID:        marketTicker2.ID,
		Symbol:    marketTicker2.Symbol,
		Exchange:  marketTicker2.Exchange,
		LastPrice: marketTicker2.Price,
		Timestamp: marketTicker2.LastUpdated,
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
	marketCandle := &market.Candle{
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

	// Convert to model.Kline
	kline := &model.Kline{
		Symbol:      marketCandle.Symbol,
		Exchange:    marketCandle.Exchange,
		Interval:    model.KlineInterval(string(marketCandle.Interval)),
		OpenTime:    marketCandle.OpenTime,
		CloseTime:   marketCandle.CloseTime,
		Open:        marketCandle.Open,
		High:        marketCandle.High,
		Low:         marketCandle.Low,
		Close:       marketCandle.Close,
		Volume:      marketCandle.Volume,
		QuoteVolume: marketCandle.QuoteVolume,
		TradeCount:  marketCandle.TradeCount,
		Complete:    marketCandle.Complete,
	}

	// Save the candle
	err := repo.SaveCandle(ctx, kline)
	require.NoError(t, err)

	// Retrieve the candle
	retrievedCandle, err := repo.GetCandle(ctx, "BTCUSDT", "mexc", model.KlineInterval("1h"), now)
	require.NoError(t, err)

	// Verify the candle was saved correctly
	assert.Equal(t, marketCandle.Symbol, retrievedCandle.Symbol)
	assert.Equal(t, marketCandle.Exchange, retrievedCandle.Exchange)
	assert.Equal(t, string(marketCandle.Interval), string(retrievedCandle.Interval))
	assert.Equal(t, marketCandle.OpenTime.Unix(), retrievedCandle.OpenTime.Unix())
	assert.Equal(t, marketCandle.CloseTime.Unix(), retrievedCandle.CloseTime.Unix())
	assert.Equal(t, marketCandle.Open, retrievedCandle.Open)
	assert.Equal(t, marketCandle.High, retrievedCandle.High)
	assert.Equal(t, marketCandle.Low, retrievedCandle.Low)
	assert.Equal(t, marketCandle.Close, retrievedCandle.Close)
	assert.Equal(t, marketCandle.Volume, retrievedCandle.Volume)
	assert.Equal(t, marketCandle.QuoteVolume, retrievedCandle.QuoteVolume)
	assert.Equal(t, marketCandle.TradeCount, retrievedCandle.TradeCount)
	assert.Equal(t, marketCandle.Complete, retrievedCandle.Complete)
}

func TestGetCandles(t *testing.T) {
	repo, cleanup := setupTestRepository(t)
	defer cleanup()

	ctx := context.Background()

	// Create test candles
	now := time.Now().Round(time.Millisecond)
	marketCandle1 := &market.Candle{
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

	marketCandle2 := &market.Candle{
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

	// Convert to model.Kline
	kline1 := &model.Kline{
		Symbol:    marketCandle1.Symbol,
		Exchange:  marketCandle1.Exchange,
		Interval:  model.KlineInterval(string(marketCandle1.Interval)),
		OpenTime:  marketCandle1.OpenTime,
		CloseTime: marketCandle1.CloseTime,
		Open:      marketCandle1.Open,
		High:      marketCandle1.High,
		Low:       marketCandle1.Low,
		Close:     marketCandle1.Close,
		Volume:    marketCandle1.Volume,
		Complete:  marketCandle1.Complete,
	}

	kline2 := &model.Kline{
		Symbol:    marketCandle2.Symbol,
		Exchange:  marketCandle2.Exchange,
		Interval:  model.KlineInterval(string(marketCandle2.Interval)),
		OpenTime:  marketCandle2.OpenTime,
		CloseTime: marketCandle2.CloseTime,
		Open:      marketCandle2.Open,
		High:      marketCandle2.High,
		Low:       marketCandle2.Low,
		Close:     marketCandle2.Close,
		Volume:    marketCandle2.Volume,
		Complete:  marketCandle2.Complete,
	}

	// Save the candles
	err := repo.SaveCandle(ctx, kline1)
	require.NoError(t, err)

	err = repo.SaveCandle(ctx, kline2)
	require.NoError(t, err)

	// Retrieve candles within a time range
	start := now.Add(-3 * time.Hour)
	end := now.Add(1 * time.Hour)
	candles, err := repo.GetCandles(ctx, "BTCUSDT", "mexc", model.KlineInterval("1h"), start, end, 10)
	require.NoError(t, err)

	// Verify the candles were retrieved correctly
	assert.Equal(t, 2, len(candles))
	assert.Equal(t, marketCandle1.OpenTime.Unix(), candles[0].OpenTime.Unix())
	assert.Equal(t, marketCandle2.OpenTime.Unix(), candles[1].OpenTime.Unix())
}

func TestGetLatestCandle(t *testing.T) {
	repo, cleanup := setupTestRepository(t)
	defer cleanup()

	ctx := context.Background()

	// Create test candles with different times
	now := time.Now().Round(time.Millisecond)
	marketCandle1 := &market.Candle{
		Symbol:    "BTCUSDT",
		Exchange:  "mexc",
		Interval:  market.Interval1h,
		OpenTime:  now.Add(-2 * time.Hour),
		CloseTime: now.Add(-1 * time.Hour),
		Open:      49000.0,
		Close:     49500.0,
		Complete:  true,
	}

	marketCandle2 := &market.Candle{
		Symbol:    "BTCUSDT",
		Exchange:  "mexc",
		Interval:  market.Interval1h,
		OpenTime:  now.Add(-1 * time.Hour),
		CloseTime: now,
		Open:      49500.0,
		Close:     50500.0,
		Complete:  true,
	}

	// Convert to model.Kline
	kline1 := &model.Kline{
		Symbol:    marketCandle1.Symbol,
		Exchange:  marketCandle1.Exchange,
		Interval:  model.KlineInterval(string(marketCandle1.Interval)),
		OpenTime:  marketCandle1.OpenTime,
		CloseTime: marketCandle1.CloseTime,
		Open:      marketCandle1.Open,
		Close:     marketCandle1.Close,
		Complete:  marketCandle1.Complete,
	}

	kline2 := &model.Kline{
		Symbol:    marketCandle2.Symbol,
		Exchange:  marketCandle2.Exchange,
		Interval:  model.KlineInterval(string(marketCandle2.Interval)),
		OpenTime:  marketCandle2.OpenTime,
		CloseTime: marketCandle2.CloseTime,
		Open:      marketCandle2.Open,
		Close:     marketCandle2.Close,
		Complete:  marketCandle2.Complete,
	}

	// Save the candles
	err := repo.SaveCandle(ctx, kline1)
	require.NoError(t, err)

	err = repo.SaveCandle(ctx, kline2)
	require.NoError(t, err)

	// Retrieve the latest candle
	latestCandle, err := repo.GetLatestCandle(ctx, "BTCUSDT", "mexc", model.KlineInterval("1h"))
	require.NoError(t, err)

	// Verify the latest candle was returned
	assert.Equal(t, marketCandle2.OpenTime.Unix(), latestCandle.OpenTime.Unix())
	assert.Equal(t, marketCandle2.Close, latestCandle.Close)
}

func TestSaveAndGetSymbol(t *testing.T) {
	repo, cleanup := setupTestRepository(t)
	defer cleanup()

	ctx := context.Background()

	// Create a test symbol
	marketSymbol := &market.Symbol{
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

	// Convert to model.Symbol
	symbol := &model.Symbol{
		Symbol:            marketSymbol.Symbol,
		BaseAsset:         marketSymbol.BaseAsset,
		QuoteAsset:        marketSymbol.QuoteAsset,
		Exchange:          marketSymbol.Exchange,
		Status:            model.SymbolStatus(marketSymbol.Status),
		MinPrice:          marketSymbol.MinPrice,
		MaxPrice:          marketSymbol.MaxPrice,
		PricePrecision:    marketSymbol.PricePrecision,
		MinQuantity:       marketSymbol.MinQty,
		MaxQuantity:       marketSymbol.MaxQty,
		QuantityPrecision: marketSymbol.QtyPrecision,
		AllowedOrderTypes: marketSymbol.AllowedOrderTypes,
	}

	// Save the symbol
	err := repo.Create(ctx, symbol)
	require.NoError(t, err)

	// Retrieve the symbol
	retrievedSymbol, err := repo.GetBySymbol(ctx, "BTCUSDT")
	require.NoError(t, err)

	// Verify the symbol was saved correctly
	assert.Equal(t, marketSymbol.Symbol, retrievedSymbol.Symbol)
	assert.Equal(t, marketSymbol.BaseAsset, retrievedSymbol.BaseAsset)
	assert.Equal(t, marketSymbol.QuoteAsset, retrievedSymbol.QuoteAsset)
	assert.Equal(t, marketSymbol.Exchange, retrievedSymbol.Exchange)
	assert.Equal(t, marketSymbol.Status, string(retrievedSymbol.Status))
	assert.Equal(t, marketSymbol.MinPrice, retrievedSymbol.MinPrice)
	assert.Equal(t, marketSymbol.MaxPrice, retrievedSymbol.MaxPrice)
	assert.Equal(t, marketSymbol.PricePrecision, retrievedSymbol.PricePrecision)
	assert.Equal(t, marketSymbol.MinQty, retrievedSymbol.MinQuantity)
	assert.Equal(t, marketSymbol.MaxQty, retrievedSymbol.MaxQuantity)
	assert.Equal(t, marketSymbol.QtyPrecision, retrievedSymbol.QuantityPrecision)
	assert.ElementsMatch(t, marketSymbol.AllowedOrderTypes, retrievedSymbol.AllowedOrderTypes)
}

func TestGetByExchange(t *testing.T) {
	repo, cleanup := setupTestRepository(t)
	defer cleanup()

	ctx := context.Background()

	// Create test symbols
	marketSymbol1 := &market.Symbol{
		Symbol:    "BTCUSDT",
		BaseAsset: "BTC",
		Exchange:  "mexc",
		Status:    "TRADING",
	}

	marketSymbol2 := &market.Symbol{
		Symbol:    "ETHUSDT",
		BaseAsset: "ETH",
		Exchange:  "mexc",
		Status:    "TRADING",
	}

	marketSymbol3 := &market.Symbol{
		Symbol:    "BTCUSDT",
		BaseAsset: "BTC",
		Exchange:  "binance",
		Status:    "TRADING",
	}

	// Convert to model.Symbol
	symbol1 := &model.Symbol{
		Symbol:    marketSymbol1.Symbol,
		BaseAsset: marketSymbol1.BaseAsset,
		Exchange:  marketSymbol1.Exchange,
		Status:    model.SymbolStatus(marketSymbol1.Status),
	}

	symbol2 := &model.Symbol{
		Symbol:    marketSymbol2.Symbol,
		BaseAsset: marketSymbol2.BaseAsset,
		Exchange:  marketSymbol2.Exchange,
		Status:    model.SymbolStatus(marketSymbol2.Status),
	}

	symbol3 := &model.Symbol{
		Symbol:    marketSymbol3.Symbol,
		BaseAsset: marketSymbol3.BaseAsset,
		Exchange:  marketSymbol3.Exchange,
		Status:    model.SymbolStatus(marketSymbol3.Status),
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
	marketSymbol1 := &market.Symbol{
		Symbol:    "BTCUSDT",
		BaseAsset: "BTC",
		Exchange:  "mexc",
		Status:    "TRADING",
	}

	marketSymbol2 := &market.Symbol{
		Symbol:    "ETHUSDT",
		BaseAsset: "ETH",
		Exchange:  "mexc",
		Status:    "TRADING",
	}

	// Convert to model.Symbol
	symbol1 := &model.Symbol{
		Symbol:    marketSymbol1.Symbol,
		BaseAsset: marketSymbol1.BaseAsset,
		Exchange:  marketSymbol1.Exchange,
		Status:    model.SymbolStatus(marketSymbol1.Status),
	}

	symbol2 := &model.Symbol{
		Symbol:    marketSymbol2.Symbol,
		BaseAsset: marketSymbol2.BaseAsset,
		Exchange:  marketSymbol2.Exchange,
		Status:    model.SymbolStatus(marketSymbol2.Status),
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
	marketSymbol := &market.Symbol{
		Symbol:    "BTCUSDT",
		BaseAsset: "BTC",
		Exchange:  "mexc",
		Status:    "TRADING",
	}

	// Convert to model.Symbol
	symbol := &model.Symbol{
		Symbol:    marketSymbol.Symbol,
		BaseAsset: marketSymbol.BaseAsset,
		Exchange:  marketSymbol.Exchange,
		Status:    model.SymbolStatus(marketSymbol.Status),
	}

	// Save the symbol
	err := repo.Create(ctx, symbol)
	require.NoError(t, err)

	// Update the symbol
	symbol.Status = model.SymbolStatus("BREAK")
	err = repo.Update(ctx, symbol)
	require.NoError(t, err)

	// Retrieve the updated symbol
	updatedSymbol, err := repo.GetBySymbol(ctx, "BTCUSDT")
	require.NoError(t, err)

	// Verify the symbol was updated correctly
	assert.Equal(t, model.SymbolStatus("BREAK"), updatedSymbol.Status)
}

func TestDeleteSymbol(t *testing.T) {
	repo, cleanup := setupTestRepository(t)
	defer cleanup()

	ctx := context.Background()

	// Create a test symbol
	marketSymbol := &market.Symbol{
		Symbol:    "BTCUSDT",
		BaseAsset: "BTC",
		Exchange:  "mexc",
		Status:    "TRADING",
	}

	// Convert to model.Symbol
	symbol := &model.Symbol{
		Symbol:    marketSymbol.Symbol,
		BaseAsset: marketSymbol.BaseAsset,
		Exchange:  marketSymbol.Exchange,
		Status:    model.SymbolStatus(marketSymbol.Status),
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
	oldMarketTicker := &market.Ticker{
		ID:          "old-ticker",
		Symbol:      "BTCUSDT",
		Exchange:    "mexc",
		Price:       48000.0,
		LastUpdated: oldTime,
	}

	// Create new ticker
	newMarketTicker := &market.Ticker{
		ID:          "new-ticker",
		Symbol:      "BTCUSDT",
		Exchange:    "mexc",
		Price:       50000.0,
		LastUpdated: newTime,
	}

	// Convert to model.Ticker
	oldTicker := &model.Ticker{
		ID:        oldMarketTicker.ID,
		Symbol:    oldMarketTicker.Symbol,
		Exchange:  oldMarketTicker.Exchange,
		LastPrice: oldMarketTicker.Price,
		Timestamp: oldMarketTicker.LastUpdated,
	}

	newTicker := &model.Ticker{
		ID:        newMarketTicker.ID,
		Symbol:    newMarketTicker.Symbol,
		Exchange:  newMarketTicker.Exchange,
		LastPrice: newMarketTicker.Price,
		Timestamp: newMarketTicker.LastUpdated,
	}

	// Save the tickers
	err := repo.SaveTicker(ctx, oldTicker)
	require.NoError(t, err)

	err = repo.SaveTicker(ctx, newTicker)
	require.NoError(t, err)

	// Create old candle
	oldMarketCandle := &market.Candle{
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
	newMarketCandle := &market.Candle{
		Symbol:    "BTCUSDT",
		Exchange:  "mexc",
		Interval:  market.Interval1h,
		OpenTime:  newTime,
		CloseTime: newTime.Add(1 * time.Hour),
		Open:      50000.0,
		Close:     50500.0,
		Complete:  true,
	}

	// Convert to model.Kline
	oldCandle := &model.Kline{
		Symbol:    oldMarketCandle.Symbol,
		Exchange:  oldMarketCandle.Exchange,
		Interval:  model.KlineInterval(string(oldMarketCandle.Interval)),
		OpenTime:  oldMarketCandle.OpenTime,
		CloseTime: oldMarketCandle.CloseTime,
		Open:      oldMarketCandle.Open,
		Close:     oldMarketCandle.Close,
		Complete:  oldMarketCandle.Complete,
	}

	newCandle := &model.Kline{
		Symbol:    newMarketCandle.Symbol,
		Exchange:  newMarketCandle.Exchange,
		Interval:  model.KlineInterval(string(newMarketCandle.Interval)),
		OpenTime:  newMarketCandle.OpenTime,
		CloseTime: newMarketCandle.CloseTime,
		Open:      newMarketCandle.Open,
		Close:     newMarketCandle.Close,
		Complete:  newMarketCandle.Complete,
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

	candles, err := repo.GetCandles(ctx, "BTCUSDT", "mexc", model.KlineInterval("1h"), oldTime, newTime.Add(1*time.Hour), 10)
	require.NoError(t, err)
	assert.Equal(t, 1, len(candles))
	assert.Equal(t, newTime.Unix(), candles[0].OpenTime.Unix())
}
