package memory

import (
	"context"
	"testing"
	"time"

	"github.com/neo/crypto-bot/internal/domain/model/market"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func setupTestCache() *MarketCache {
	logger := zerolog.New(zerolog.NewConsoleWriter()).With().Timestamp().Logger()
	cache := NewMarketCache(&logger).(*MarketCache)
	cache.SetTickerExpiry(1 * time.Minute)
	cache.SetCandleExpiry(5 * time.Minute)
	cache.SetOrderbookExpiry(30 * time.Second)
	return cache
}

func TestCacheTicker(t *testing.T) {
	cache := setupTestCache()
	
	// Create a test ticker
	ticker := &market.Ticker{
		Symbol:        "BTCUSDT",
		Exchange:      "mexc",
		Price:         50000.0,
		Volume:        100.0,
		High24h:       51000.0,
		Low24h:        49000.0,
		PriceChange:   1000.0,
		PercentChange: 2.0,
		LastUpdated:   time.Now(),
	}
	
	// Cache the ticker
	cache.CacheTicker(ticker)
	
	// Retrieve the ticker
	ctx := context.Background()
	cachedTicker, exists := cache.GetTicker(ctx, "mexc", "BTCUSDT")
	
	// Verify the ticker was cached correctly
	assert.True(t, exists)
	assert.Equal(t, ticker.Symbol, cachedTicker.Symbol)
	assert.Equal(t, ticker.Price, cachedTicker.Price)
	assert.Equal(t, ticker.Exchange, cachedTicker.Exchange)
}

func TestGetAllTickers(t *testing.T) {
	cache := setupTestCache()
	
	// Create test tickers
	ticker1 := &market.Ticker{
		Symbol:        "BTCUSDT",
		Exchange:      "mexc",
		Price:         50000.0,
		LastUpdated:   time.Now(),
	}
	
	ticker2 := &market.Ticker{
		Symbol:        "ETHUSDT",
		Exchange:      "mexc",
		Price:         3000.0,
		LastUpdated:   time.Now(),
	}
	
	// Cache the tickers
	cache.CacheTicker(ticker1)
	cache.CacheTicker(ticker2)
	
	// Set the exchange update time
	cache.mu.Lock()
	cache.lastUpdated["exchange:mexc"] = time.Now()
	cache.mu.Unlock()
	
	// Retrieve all tickers for the exchange
	ctx := context.Background()
	tickers, exists := cache.GetAllTickers(ctx, "mexc")
	
	// Verify the tickers were cached correctly
	assert.True(t, exists)
	assert.Equal(t, 2, len(tickers))
	
	// Verify the ticker symbols
	symbols := []string{tickers[0].Symbol, tickers[1].Symbol}
	assert.Contains(t, symbols, "BTCUSDT")
	assert.Contains(t, symbols, "ETHUSDT")
}

func TestGetLatestTickers(t *testing.T) {
	cache := setupTestCache()
	
	// Create test tickers
	ticker1 := &market.Ticker{
		Symbol:        "BTCUSDT",
		Exchange:      "mexc",
		Price:         50000.0,
		LastUpdated:   time.Now(),
	}
	
	ticker2 := &market.Ticker{
		Symbol:        "ETHUSDT",
		Exchange:      "mexc",
		Price:         3000.0,
		LastUpdated:   time.Now(),
	}
	
	// Cache the tickers
	cache.CacheTicker(ticker1)
	cache.CacheTicker(ticker2)
	
	// Retrieve latest tickers
	ctx := context.Background()
	tickers, exists := cache.GetLatestTickers(ctx)
	
	// Verify the tickers were cached correctly
	assert.True(t, exists)
	assert.Equal(t, 2, len(tickers))
	
	// Verify the ticker symbols
	symbols := []string{tickers[0].Symbol, tickers[1].Symbol}
	assert.Contains(t, symbols, "BTCUSDT")
	assert.Contains(t, symbols, "ETHUSDT")
}

func TestCacheCandle(t *testing.T) {
	cache := setupTestCache()
	
	// Create a test candle
	now := time.Now()
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
	
	// Cache the candle
	cache.CacheCandle(candle)
	
	// Retrieve the candle
	ctx := context.Background()
	cachedCandle, exists := cache.GetCandle(ctx, "mexc", "BTCUSDT", market.Interval1h, now)
	
	// Verify the candle was cached correctly
	assert.True(t, exists)
	assert.Equal(t, candle.Symbol, cachedCandle.Symbol)
	assert.Equal(t, candle.Exchange, cachedCandle.Exchange)
	assert.Equal(t, candle.Interval, cachedCandle.Interval)
	assert.Equal(t, candle.Open, cachedCandle.Open)
	assert.Equal(t, candle.Close, cachedCandle.Close)
}

func TestGetLatestCandle(t *testing.T) {
	cache := setupTestCache()
	
	// Create test candles with different times
	now := time.Now()
	candle1 := &market.Candle{
		Symbol:      "BTCUSDT",
		Exchange:    "mexc",
		Interval:    market.Interval1h,
		OpenTime:    now.Add(-2 * time.Hour),
		CloseTime:   now.Add(-1 * time.Hour),
		Open:        49000.0,
		High:        50000.0,
		Low:         48000.0,
		Close:       49500.0,
		Volume:      90.0,
		Complete:    true,
	}
	
	candle2 := &market.Candle{
		Symbol:      "BTCUSDT",
		Exchange:    "mexc",
		Interval:    market.Interval1h,
		OpenTime:    now.Add(-1 * time.Hour),
		CloseTime:   now,
		Open:        49500.0,
		High:        51000.0,
		Low:         49000.0,
		Close:       50500.0,
		Volume:      100.0,
		Complete:    true,
	}
	
	// Cache the candles
	cache.CacheCandle(candle1)
	cache.CacheCandle(candle2)
	
	// Retrieve the latest candle
	ctx := context.Background()
	latestCandle, exists := cache.GetLatestCandle(ctx, "mexc", "BTCUSDT", market.Interval1h)
	
	// Verify the latest candle was returned
	assert.True(t, exists)
	assert.Equal(t, candle2.OpenTime, latestCandle.OpenTime)
	assert.Equal(t, candle2.Close, latestCandle.Close)
}

func TestCacheOrderBook(t *testing.T) {
	cache := setupTestCache()
	
	// Create a test orderbook
	orderbook := &market.OrderBook{
		Symbol:      "BTCUSDT",
		Exchange:    "mexc",
		LastUpdated: time.Now(),
		Bids: []market.OrderBookEntry{
			{Price: 49900.0, Quantity: 1.5},
			{Price: 49800.0, Quantity: 2.0},
		},
		Asks: []market.OrderBookEntry{
			{Price: 50000.0, Quantity: 1.0},
			{Price: 50100.0, Quantity: 2.5},
		},
	}
	
	// Cache the orderbook
	cache.CacheOrderBook(orderbook)
	
	// Retrieve the orderbook
	ctx := context.Background()
	cachedOrderbook, exists := cache.GetOrderBook(ctx, "mexc", "BTCUSDT")
	
	// Verify the orderbook was cached correctly
	assert.True(t, exists)
	assert.Equal(t, orderbook.Symbol, cachedOrderbook.Symbol)
	assert.Equal(t, orderbook.Exchange, cachedOrderbook.Exchange)
	assert.Equal(t, 2, len(cachedOrderbook.Bids))
	assert.Equal(t, 2, len(cachedOrderbook.Asks))
	assert.Equal(t, 49900.0, cachedOrderbook.Bids[0].Price)
	assert.Equal(t, 50000.0, cachedOrderbook.Asks[0].Price)
}

func TestCacheExpiry(t *testing.T) {
	cache := setupTestCache()
	
	// Set very short expiry for testing
	cache.SetTickerExpiry(50 * time.Millisecond)
	
	// Create a test ticker
	ticker := &market.Ticker{
		Symbol:        "BTCUSDT",
		Exchange:      "mexc",
		Price:         50000.0,
		LastUpdated:   time.Now(),
	}
	
	// Cache the ticker
	cache.CacheTicker(ticker)
	
	// Verify it exists immediately
	ctx := context.Background()
	_, exists := cache.GetTicker(ctx, "mexc", "BTCUSDT")
	assert.True(t, exists)
	
	// Wait for expiry
	time.Sleep(100 * time.Millisecond)
	
	// Verify it's expired
	_, exists = cache.GetTicker(ctx, "mexc", "BTCUSDT")
	assert.False(t, exists)
}

func TestClear(t *testing.T) {
	cache := setupTestCache()
	
	// Create and cache test data
	ticker := &market.Ticker{
		Symbol:        "BTCUSDT",
		Exchange:      "mexc",
		Price:         50000.0,
		LastUpdated:   time.Now(),
	}
	cache.CacheTicker(ticker)
	
	candle := &market.Candle{
		Symbol:      "BTCUSDT",
		Exchange:    "mexc",
		Interval:    market.Interval1h,
		OpenTime:    time.Now(),
		Close:       50500.0,
	}
	cache.CacheCandle(candle)
	
	// Clear the cache
	cache.Clear()
	
	// Verify everything is cleared
	ctx := context.Background()
	_, tickerExists := cache.GetTicker(ctx, "mexc", "BTCUSDT")
	_, candleExists := cache.GetCandle(ctx, "mexc", "BTCUSDT", market.Interval1h, candle.OpenTime)
	
	assert.False(t, tickerExists)
	assert.False(t, candleExists)
}
