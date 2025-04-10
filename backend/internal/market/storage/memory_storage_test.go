package storage

import (
	"context"
	"testing"
	"time"

	"go-crypto-bot-clean/backend/internal/domain/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMemoryStorage(t *testing.T) {
	storage := NewMemoryStorage(1000)
	assert.NotNil(t, storage)
	assert.Equal(t, 1000, storage.maxSize)
}

func TestMemoryStorage_StoreAndGetCandles(t *testing.T) {
	ctx := context.Background()
	storage := NewMemoryStorage(1000)
	symbol := "BTC-USD"
	now := time.Now()

	candles := []*models.Candle{
		{
			Symbol:     symbol,
			OpenTime:   now,
			OpenPrice:  100.0,
			HighPrice:  110.0,
			LowPrice:   90.0,
			ClosePrice: 105.0,
			Volume:     10.0,
			CloseTime:  now,
		},
		{
			Symbol:     symbol,
			OpenTime:   now.Add(time.Minute),
			OpenPrice:  105.0,
			HighPrice:  115.0,
			LowPrice:   95.0,
			ClosePrice: 110.0,
			Volume:     15.0,
			CloseTime:  now.Add(time.Minute),
		},
	}

	err := storage.StoreCandles(ctx, symbol, candles)
	require.NoError(t, err)

	// Test retrieval
	retrieved, err := storage.GetCandles(ctx, MarketDataFilter{
		Symbol:    symbol,
		StartTime: now.Add(-time.Hour),
		EndTime:   now.Add(time.Hour),
	})
	require.NoError(t, err)
	assert.Len(t, retrieved, 2)
	assert.Equal(t, candles[0].ClosePrice, retrieved[0].ClosePrice)
	assert.Equal(t, candles[1].ClosePrice, retrieved[1].ClosePrice)
}

func TestMemoryStorage_StoreAndGetTrades(t *testing.T) {
	ctx := context.Background()
	storage := NewMemoryStorage(1000)
	symbol := "BTC-USD"
	now := time.Now()

	trades := []*models.MarketTrade{
		{
			Symbol:    symbol,
			Price:     100.0,
			Quantity:  1.0,
			IsBuyer:   true,
			Timestamp: now,
		},
		{
			Symbol:    symbol,
			Price:     101.0,
			Quantity:  2.0,
			IsBuyer:   false,
			Timestamp: now.Add(time.Minute),
		},
	}

	err := storage.StoreTrades(ctx, symbol, trades)
	require.NoError(t, err)

	// Test retrieval
	filter := MarketDataFilter{
		Symbol:    symbol,
		StartTime: now.Add(-time.Hour),
		EndTime:   now.Add(time.Hour),
	}

	retrieved, err := storage.GetTrades(ctx, filter)
	require.NoError(t, err)
	assert.Len(t, retrieved, 2)
	assert.Equal(t, trades[0].Price, retrieved[0].Price)
	assert.Equal(t, trades[1].Price, retrieved[1].Price)
}

func TestMemoryStorage_StoreAndGetOrderBook(t *testing.T) {
	ctx := context.Background()
	storage := NewMemoryStorage(1000)
	symbol := "BTC-USD"
	now := time.Now()

	book := &models.OrderBook{
		Symbol: symbol,
		Bids: []models.CandleOrderBookEntry{
			{Price: 100.0, Quantity: 1.0, Type: "bid"},
			{Price: 99.0, Quantity: 2.0, Type: "bid"},
		},
		Asks: []models.CandleOrderBookEntry{
			{Price: 101.0, Quantity: 1.0, Type: "ask"},
			{Price: 102.0, Quantity: 2.0, Type: "ask"},
		},
		Timestamp: now,
	}

	err := storage.StoreOrderBook(ctx, symbol, book)
	require.NoError(t, err)

	// Test retrieval
	retrieved, err := storage.GetOrderBook(ctx, symbol, now)
	require.NoError(t, err)
	assert.Equal(t, book.Bids[0].Price, retrieved.Bids[0].Price)
	assert.Equal(t, book.Asks[0].Price, retrieved.Asks[0].Price)
}

func TestMemoryStorage_StoreAndGetTicker(t *testing.T) {
	ctx := context.Background()
	storage := NewMemoryStorage(1000)
	symbol := "BTC-USD"
	now := time.Now()

	ticker := &models.Ticker{
		Symbol:      symbol,
		Price:       100.5,
		PriceChange: 0.5,
		Volume:      10.0,
		High24h:     101.0,
		Low24h:      100.0,
		Timestamp:   now,
	}

	err := storage.StoreTicker(ctx, symbol, ticker)
	require.NoError(t, err)

	// Test retrieval
	retrieved, err := storage.GetTicker(ctx, symbol)
	require.NoError(t, err)
	assert.Equal(t, ticker.Price, retrieved.Price)
	assert.Equal(t, ticker.Volume, retrieved.Volume)
	assert.Equal(t, ticker.Timestamp.Unix(), retrieved.Timestamp.Unix())
}

func TestMemoryStorage_GetVWAP(t *testing.T) {
	ctx := context.Background()
	storage := NewMemoryStorage(1000)
	symbol := "BTC-USD"
	now := time.Now()

	trades := []*models.MarketTrade{
		{
			Symbol:    symbol,
			Price:     100.0,
			Quantity:  1.0,
			Timestamp: now,
		},
		{
			Symbol:    symbol,
			Price:     200.0,
			Quantity:  2.0,
			Timestamp: now.Add(time.Minute),
		},
	}

	err := storage.StoreTrades(ctx, symbol, trades)
	require.NoError(t, err)

	// Expected VWAP = (100*1 + 200*2) / (1 + 2) = 166.67
	vwap, err := storage.GetVWAP(ctx, symbol, now.Add(-time.Hour), now.Add(time.Hour))
	require.NoError(t, err)
	assert.InDelta(t, 166.67, vwap, 0.01)
}

func TestMemoryStorage_GetVolume(t *testing.T) {
	ctx := context.Background()
	storage := NewMemoryStorage(1000)
	symbol := "BTC-USD"
	now := time.Now()

	trades := []*models.MarketTrade{
		{
			Symbol:    symbol,
			Price:     100.0,
			Quantity:  1.0,
			Timestamp: now,
		},
		{
			Symbol:    symbol,
			Price:     200.0,
			Quantity:  2.0,
			Timestamp: now.Add(time.Minute),
		},
	}

	err := storage.StoreTrades(ctx, symbol, trades)
	require.NoError(t, err)

	volume, err := storage.GetVolume(ctx, symbol, now.Add(-time.Hour), now.Add(time.Hour))
	require.NoError(t, err)
	assert.Equal(t, 3.0, volume)
}

func TestMemoryStorage_Cleanup(t *testing.T) {
	ctx := context.Background()
	storage := NewMemoryStorage(1000)
	symbol := "BTC-USD"
	now := time.Now()

	trades := []*models.MarketTrade{
		{
			Symbol:    symbol,
			Price:     100.0,
			Quantity:  1.0,
			Timestamp: now.Add(-2 * time.Hour),
		},
		{
			Symbol:    symbol,
			Price:     200.0,
			Quantity:  2.0,
			Timestamp: now,
		},
	}

	err := storage.StoreTrades(ctx, symbol, trades)
	require.NoError(t, err)

	err = storage.Cleanup(ctx, now.Add(-time.Hour), MarketDataTypeTrade)
	require.NoError(t, err)

	retrieved, err := storage.GetTrades(ctx, MarketDataFilter{
		Symbol:    symbol,
		StartTime: now.Add(-3 * time.Hour),
		EndTime:   now.Add(time.Hour),
	})
	require.NoError(t, err)
	assert.Len(t, retrieved, 1)
	assert.Equal(t, 200.0, retrieved[0].Price)
}

func TestMemoryStorage_GetDataTypes(t *testing.T) {
	ctx := context.Background()
	storage := NewMemoryStorage(1000)
	symbol := "BTC-USD"
	now := time.Now()

	// Store different types of data
	candles := []*models.Candle{{Symbol: symbol, OpenTime: now, CloseTime: now}}
	trades := []*models.MarketTrade{{Symbol: symbol, Timestamp: now}}
	book := &models.OrderBook{Symbol: symbol, Timestamp: now}
	ticker := &models.Ticker{Symbol: symbol, Timestamp: now}

	require.NoError(t, storage.StoreCandles(ctx, symbol, candles))
	require.NoError(t, storage.StoreTrades(ctx, symbol, trades))
	require.NoError(t, storage.StoreOrderBook(ctx, symbol, book))
	require.NoError(t, storage.StoreTicker(ctx, symbol, ticker))

	types, err := storage.GetDataTypes(ctx, symbol)
	require.NoError(t, err)
	assert.Len(t, types, 4)
	assert.Contains(t, types, MarketDataTypeCandle)
	assert.Contains(t, types, MarketDataTypeTrade)
	assert.Contains(t, types, MarketDataTypeOrderBook)
	assert.Contains(t, types, MarketDataTypeTicker)
}

func TestMemoryStorage_GetTimeRange(t *testing.T) {
	ctx := context.Background()
	storage := NewMemoryStorage(1000)
	symbol := "BTC-USD"
	now := time.Now()

	candles := []*models.Candle{
		{Symbol: symbol, OpenTime: now.Add(-time.Hour), CloseTime: now.Add(-time.Hour)},
		{Symbol: symbol, OpenTime: now, CloseTime: now},
	}

	require.NoError(t, storage.StoreCandles(ctx, symbol, candles))

	start, end, err := storage.GetTimeRange(ctx, symbol, MarketDataTypeCandle)
	require.NoError(t, err)
	assert.Equal(t, now.Add(-time.Hour).Unix(), start.Unix())
	assert.Equal(t, now.Unix(), end.Unix())
}

func TestMemoryStorage_GetSymbols(t *testing.T) {
	ctx := context.Background()
	storage := NewMemoryStorage(1000)
	now := time.Now()

	symbols := []string{"BTC-USD", "ETH-USD"}
	for _, symbol := range symbols {
		candles := []*models.Candle{{Symbol: symbol, OpenTime: now, CloseTime: now}}
		require.NoError(t, storage.StoreCandles(ctx, symbol, candles))
	}

	retrieved, err := storage.GetSymbols(ctx)
	require.NoError(t, err)
	assert.Len(t, retrieved, 2)
	assert.Contains(t, retrieved, "BTC-USD")
	assert.Contains(t, retrieved, "ETH-USD")
}

func TestMemoryStorage_MaxSize(t *testing.T) {
	ctx := context.Background()
	storage := NewMemoryStorage(2)
	symbol := "BTC-USD"
	now := time.Now()

	candles := []*models.Candle{
		{Symbol: symbol, OpenPrice: 100.0, OpenTime: now, CloseTime: now},
		{Symbol: symbol, OpenPrice: 101.0, OpenTime: now.Add(time.Minute), CloseTime: now.Add(time.Minute)},
		{Symbol: symbol, OpenPrice: 102.0, OpenTime: now.Add(2 * time.Minute), CloseTime: now.Add(2 * time.Minute)},
	}

	require.NoError(t, storage.StoreCandles(ctx, symbol, candles))

	retrieved, err := storage.GetCandles(ctx, MarketDataFilter{
		Symbol:    symbol,
		StartTime: now.Add(-time.Hour),
		EndTime:   now.Add(time.Hour),
	})
	require.NoError(t, err)
	assert.Len(t, retrieved, 2)
	assert.Equal(t, 101.0, retrieved[0].OpenPrice)
	assert.Equal(t, 102.0, retrieved[1].OpenPrice)
}
