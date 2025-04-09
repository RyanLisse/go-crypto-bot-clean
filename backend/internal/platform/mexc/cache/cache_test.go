package cache

import (
	"testing"
	"time"

	"github.com/ryanlisse/go-crypto-bot/internal/domain/models"
	"github.com/stretchr/testify/assert"
)

func TestTickerCache(t *testing.T) {
	cache := NewTickerCache()

	// Test setting and getting a ticker
	symbol := "BTC/USDT"
	ticker := &models.Ticker{
		Symbol: symbol,
		Price:  50000.0,
	}

	// Set the ticker with a 1 second TTL
	cache.SetTicker(symbol, ticker, 1*time.Second)

	// Get the ticker
	cachedTicker, exists := cache.GetTicker(symbol)
	assert.True(t, exists)
	assert.Equal(t, ticker, cachedTicker)

	// Wait for the ticker to expire
	time.Sleep(1100 * time.Millisecond)

	// Try to get the expired ticker
	cachedTicker, exists = cache.GetTicker(symbol)
	assert.False(t, exists)
	assert.Nil(t, cachedTicker)

	// Test deleting a ticker
	cache.SetTicker(symbol, ticker, 10*time.Second)
	cache.Delete(symbol)
	cachedTicker, exists = cache.GetTicker(symbol)
	assert.False(t, exists)
	assert.Nil(t, cachedTicker)

	// Test clearing the cache
	cache.SetTicker(symbol, ticker, 10*time.Second)
	cache.Clear()
	cachedTicker, exists = cache.GetTicker(symbol)
	assert.False(t, exists)
	assert.Nil(t, cachedTicker)
}

func TestAllTickersCache(t *testing.T) {
	cache := NewTickerCache()

	// Test setting and getting all tickers
	tickers := map[string]*models.Ticker{
		"BTC/USDT": {
			Symbol: "BTC/USDT",
			Price:  50000.0,
		},
		"ETH/USDT": {
			Symbol: "ETH/USDT",
			Price:  3000.0,
		},
	}

	// Set all tickers with a 1 second TTL
	cache.SetAllTickers(tickers, 1*time.Second)

	// Get all tickers
	cachedTickers, exists := cache.GetAllTickers()
	assert.True(t, exists)
	assert.Equal(t, tickers, cachedTickers)

	// Wait for the tickers to expire
	time.Sleep(1100 * time.Millisecond)

	// Try to get the expired tickers
	cachedTickers, exists = cache.GetAllTickers()
	assert.False(t, exists)
	assert.Nil(t, cachedTickers)
}

func TestOrderBookCache(t *testing.T) {
	cache := NewOrderBookCache()

	// Test setting and getting an order book
	symbol := "BTC/USDT"
	orderBook := &models.OrderBookUpdate{
		Symbol: symbol,
		Bids: []models.OrderBookEntry{
			{Price: 50000.0, Quantity: 1.0},
			{Price: 49900.0, Quantity: 2.0},
		},
		Asks: []models.OrderBookEntry{
			{Price: 50100.0, Quantity: 1.0},
			{Price: 50200.0, Quantity: 2.0},
		},
	}

	// Set the order book with a 1 second TTL
	cache.SetOrderBook(symbol, orderBook, 1*time.Second)

	// Get the order book
	cachedOrderBook, exists := cache.GetOrderBook(symbol)
	assert.True(t, exists)
	assert.Equal(t, orderBook, cachedOrderBook)

	// Wait for the order book to expire
	time.Sleep(1100 * time.Millisecond)

	// Try to get the expired order book
	cachedOrderBook, exists = cache.GetOrderBook(symbol)
	assert.False(t, exists)
	assert.Nil(t, cachedOrderBook)
}

func TestKlineCache(t *testing.T) {
	cache := NewKlineCache()

	// Test setting and getting klines
	symbol := "BTC/USDT"
	interval := "1h"
	limit := 10
	klines := []*models.Kline{
		{
			Symbol:    symbol,
			Interval:  interval,
			OpenTime:  time.Now().Add(-2 * time.Hour),
			CloseTime: time.Now().Add(-1 * time.Hour),
			Open:      50000.0,
			High:      51000.0,
			Low:       49000.0,
			Close:     50500.0,
			Volume:    100.0,
		},
		{
			Symbol:    symbol,
			Interval:  interval,
			OpenTime:  time.Now().Add(-1 * time.Hour),
			CloseTime: time.Now(),
			Open:      50500.0,
			High:      52000.0,
			Low:       50000.0,
			Close:     51500.0,
			Volume:    200.0,
		},
	}

	// Set the klines with a 1 second TTL
	cache.SetKlines(symbol, interval, limit, klines, 1*time.Second)

	// Get the klines
	cachedKlines, exists := cache.GetKlines(symbol, interval, limit)
	assert.True(t, exists)
	assert.Equal(t, klines, cachedKlines)

	// Wait for the klines to expire
	time.Sleep(1100 * time.Millisecond)

	// Try to get the expired klines
	cachedKlines, exists = cache.GetKlines(symbol, interval, limit)
	assert.False(t, exists)
	assert.Nil(t, cachedKlines)
}

func TestNewCoinCache(t *testing.T) {
	cache := NewNewCoinCache()

	// Test setting and getting new coins
	newCoins := []*models.NewCoin{
		{
			Symbol:      "NEW/USDT",
			FoundAt:     time.Now().Add(-1 * time.Hour),
			BaseVolume:  1000.0,
			QuoteVolume: 1000000.0,
			IsProcessed: false,
		},
		{
			Symbol:      "NEWER/USDT",
			FoundAt:     time.Now().Add(-30 * time.Minute),
			BaseVolume:  2000.0,
			QuoteVolume: 2000000.0,
			IsProcessed: false,
		},
	}

	// Set the new coins with a 1 second TTL
	cache.SetNewCoins(newCoins, 1*time.Second)

	// Get the new coins
	cachedNewCoins, exists := cache.GetNewCoins()
	assert.True(t, exists)
	assert.Equal(t, newCoins, cachedNewCoins)

	// Wait for the new coins to expire
	time.Sleep(1100 * time.Millisecond)

	// Try to get the expired new coins
	cachedNewCoins, exists = cache.GetNewCoins()
	assert.False(t, exists)
	assert.Nil(t, cachedNewCoins)
}
