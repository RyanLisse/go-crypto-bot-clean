package port

import (
	"context"
	"time"

	"github.com/neo/crypto-bot/internal/domain/model/market"
)

// MarketCache defines the interface for market data caching
type MarketCache interface {
	// Ticker operations
	CacheTicker(ticker *market.Ticker)
	GetTicker(ctx context.Context, exchange, symbol string) (*market.Ticker, bool)
	GetAllTickers(ctx context.Context, exchange string) ([]*market.Ticker, bool)
	GetLatestTickers(ctx context.Context) ([]*market.Ticker, bool)

	// Candle operations
	CacheCandle(candle *market.Candle)
	GetCandle(ctx context.Context, exchange, symbol string, interval market.Interval, openTime time.Time) (*market.Candle, bool)
	GetLatestCandle(ctx context.Context, exchange, symbol string, interval market.Interval) (*market.Candle, bool)

	// OrderBook operations
	CacheOrderBook(orderbook *market.OrderBook)
	GetOrderBook(ctx context.Context, exchange, symbol string) (*market.OrderBook, bool)

	// Cache management
	Clear()
	SetTickerExpiry(d time.Duration)
	SetCandleExpiry(d time.Duration)
	SetOrderbookExpiry(d time.Duration)
	StartCleanupTask(ctx context.Context, interval time.Duration)
}

// Cache provides a generic caching interface for any type
type Cache[T any] interface {
	// Get retrieves the cached value if it exists and is not expired
	Get() (*T, bool)

	// Set stores a value in the cache with the configured TTL
	Set(value *T)

	// GetOrSet retrieves the cached value if valid, or sets it using the provided function
	GetOrSet(fetchFn func() (*T, error)) (*T, error)

	// Invalidate clears the cached value
	Invalidate()

	// UpdateTTL changes the TTL for the cache
	UpdateTTL(ttl time.Duration)
}
