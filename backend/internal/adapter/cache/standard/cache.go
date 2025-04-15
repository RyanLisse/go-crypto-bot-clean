package standard

import (
	"context"
	"fmt"
	"strings"
	"time"

	gocache "github.com/patrickmn/go-cache"
	"github.com/rs/zerolog/log"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/apperror"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model/market"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
)

// Cache error types
const (
	ErrCacheKeyNotFound = "cache_key_not_found"
	ErrCacheExpired     = "cache_expired"
	ErrCacheInvalidType = "cache_invalid_type"
	ErrCacheNilValue    = "cache_nil_value"
)

// CacheError represents an error from the cache operations
type CacheError struct {
	Code     string
	Message  string
	Resource string
	Err      error
}

// Error returns the error message
func (e *CacheError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s for %s: %v", e.Code, e.Message, e.Resource, e.Err)
	}
	return fmt.Sprintf("%s: %s for %s", e.Code, e.Message, e.Resource)
}

// Unwrap returns the wrapped error
func (e *CacheError) Unwrap() error {
	return e.Err
}

// Is checks if the target error is a CacheError with the same code
func (e *CacheError) Is(target error) bool {
	var cacheErr *CacheError
	if target == nil {
		return false
	}
	if err, ok := target.(*CacheError); ok {
		return err.Code == e.Code
	}
	return cacheErr != nil && cacheErr.Code == e.Code
}

// NewCacheKeyNotFoundError creates a new cache key not found error
func NewCacheKeyNotFoundError(resource string, err error) *CacheError {
	return &CacheError{
		Code:     ErrCacheKeyNotFound,
		Message:  "Cache key not found",
		Resource: resource,
		Err:      err,
	}
}

// NewCacheExpiredError creates a new cache expired error
func NewCacheExpiredError(resource string, err error) *CacheError {
	return &CacheError{
		Code:     ErrCacheExpired,
		Message:  "Cache entry expired",
		Resource: resource,
		Err:      err,
	}
}

// NewCacheInvalidTypeError creates a new cache invalid type error
func NewCacheInvalidTypeError(resource string, err error) *CacheError {
	return &CacheError{
		Code:     ErrCacheInvalidType,
		Message:  "Invalid type in cache",
		Resource: resource,
		Err:      err,
	}
}

// NewCacheNilValueError creates a new cache nil value error
func NewCacheNilValueError(resource string, err error) *CacheError {
	return &CacheError{
		Code:     ErrCacheNilValue,
		Message:  "Nil value provided to cache",
		Resource: resource,
		Err:      err,
	}
}

// StandardCache implements port.MarketCache using go-cache library
type StandardCache struct {
	tickerCache    *gocache.Cache
	candleCache    *gocache.Cache
	orderBookCache *gocache.Cache
	defaultTTL     time.Duration
}

// NewStandardCache creates a new instance of StandardCache
func NewStandardCache(defaultTTL time.Duration, cleanupInterval time.Duration) port.ExtendedMarketCache {
	return &StandardCache{
		tickerCache:    gocache.New(defaultTTL, cleanupInterval),
		candleCache:    gocache.New(defaultTTL, cleanupInterval),
		orderBookCache: gocache.New(defaultTTL, cleanupInterval),
		defaultTTL:     defaultTTL,
	}
}

// CacheTicker stores a ticker in the cache with the default TTL
func (c *StandardCache) CacheTicker(ticker *market.Ticker) {
	if ticker == nil {
		return
	}
	key := c.generateTickerKey(ticker.Exchange, ticker.Symbol)
	c.tickerCache.Set(key, ticker, c.defaultTTL)

	// Also store in latest tickers collection
	latestKey := c.generateLatestTickerKey(ticker.Symbol)
	c.tickerCache.Set(latestKey, ticker, c.defaultTTL)
}

// CacheTickerWithTTL stores a ticker in the cache with a custom TTL
func (c *StandardCache) CacheTickerWithTTL(symbol string, ticker *market.Ticker, ttl time.Duration) {
	key := c.generateTickerKey(symbol, ticker.Symbol)
	c.tickerCache.Set(key, ticker, ttl)
}

// GetTicker retrieves a ticker from the cache
func (c *StandardCache) GetTicker(ctx context.Context, exchange, symbol string) (*market.Ticker, bool) {
	key := c.generateTickerKey(exchange, symbol)
	ticker, found := c.tickerCache.Get(key)
	if !found {
		return nil, false
	}

	cachedTicker, ok := ticker.(*market.Ticker)
	if !ok {
		log.Warn().Str("exchange", exchange).Str("symbol", symbol).Msg("Invalid ticker type in cache")
		return nil, false
	}

	return cachedTicker, true
}

// GetTickerWithError retrieves a ticker from the cache with error handling
func (c *StandardCache) GetTickerWithError(ctx context.Context, exchange, symbol string) (*market.Ticker, error) {
	key := c.generateTickerKey(exchange, symbol)
	ticker, found := c.tickerCache.Get(key)
	if !found {
		return nil, NewCacheKeyNotFoundError(fmt.Sprintf("ticker:%s:%s", exchange, symbol), nil)
	}

	cachedTicker, ok := ticker.(*market.Ticker)
	if !ok {
		return nil, NewCacheInvalidTypeError(fmt.Sprintf("ticker:%s:%s", exchange, symbol), nil)
	}

	return cachedTicker, nil
}

// GetAllTickers retrieves all tickers for an exchange from cache
func (c *StandardCache) GetAllTickers(ctx context.Context, exchange string) ([]*market.Ticker, bool) {
	// Get all items from cache with exchange prefix
	prefix := fmt.Sprintf("ticker:%s:", exchange)
	tickers := make([]*market.Ticker, 0)

	// Iterate through all items in the cache and filter by prefix
	for key, item := range c.tickerCache.Items() {
		if strings.HasPrefix(key, prefix) {
			if ticker, ok := item.Object.(*market.Ticker); ok {
				tickers = append(tickers, ticker)
			}
		}
	}

	if len(tickers) == 0 {
		return nil, false
	}

	return tickers, true
}

// GetLatestTickers retrieves the most recent tickers across all exchanges
func (c *StandardCache) GetLatestTickers(ctx context.Context) ([]*market.Ticker, bool) {
	prefix := "latest_ticker:"
	tickers := make([]*market.Ticker, 0)

	// Get all latest tickers from cache
	for key, item := range c.tickerCache.Items() {
		if strings.HasPrefix(key, prefix) {
			if ticker, ok := item.Object.(*market.Ticker); ok {
				tickers = append(tickers, ticker)
			}
		}
	}

	if len(tickers) == 0 {
		return nil, false
	}

	return tickers, true
}

// GetAllTickersWithError retrieves all tickers for an exchange from cache with error handling
func (c *StandardCache) GetAllTickersWithError(ctx context.Context, exchange string) ([]*market.Ticker, error) {
	// Get all items from cache with exchange prefix
	prefix := fmt.Sprintf("ticker:%s:", exchange)
	tickers := make([]*market.Ticker, 0)

	// Iterate through all items in the cache and filter by prefix
	for key, item := range c.tickerCache.Items() {
		if strings.HasPrefix(key, prefix) {
			if ticker, ok := item.Object.(*market.Ticker); ok {
				// Check if item is expired
				if item.Expiration != 0 && time.Now().UnixNano() > item.Expiration {
					continue // Skip expired items
				}
				tickers = append(tickers, ticker)
			}
		}
	}

	if len(tickers) == 0 {
		return nil, NewCacheKeyNotFoundError(fmt.Sprintf("tickers for exchange:%s", exchange), nil)
	}

	return tickers, nil
}

// GetLatestTickersWithError retrieves the most recent tickers across all exchanges with error handling
func (c *StandardCache) GetLatestTickersWithError(ctx context.Context) ([]*market.Ticker, error) {
	prefix := "latest_ticker:"
	tickers := make([]*market.Ticker, 0)

	// Get all latest tickers from cache
	for key, item := range c.tickerCache.Items() {
		if strings.HasPrefix(key, prefix) {
			if ticker, ok := item.Object.(*market.Ticker); ok {
				// Check if item is expired
				if item.Expiration != 0 && time.Now().UnixNano() > item.Expiration {
					continue // Skip expired items
				}
				tickers = append(tickers, ticker)
			}
		}
	}

	if len(tickers) == 0 {
		return nil, NewCacheKeyNotFoundError("latest tickers", nil)
	}

	return tickers, nil
}

// CacheCandle stores a candle in the cache with the default TTL
func (c *StandardCache) CacheCandle(candle *market.Candle) {
	if candle == nil {
		return
	}
	key := c.generateCandleKey(candle.Exchange, candle.Symbol, string(candle.Interval), candle.OpenTime)
	c.candleCache.Set(key, candle, c.defaultTTL)

	// Also cache as latest candle for this symbol and interval
	latestKey := c.generateLatestCandleKey(candle.Exchange, candle.Symbol, string(candle.Interval))
	c.candleCache.Set(latestKey, candle, c.defaultTTL)
}

// CacheCandleWithTTL stores a candle in the cache with a custom TTL
func (c *StandardCache) CacheCandleWithTTL(symbol string, interval string, candle *market.Candle, ttl time.Duration) {
	key := c.generateCandleKey(symbol, interval, string(candle.Interval), candle.OpenTime)
	c.candleCache.Set(key, candle, ttl)
}

// GetCandle retrieves a candle from the cache
func (c *StandardCache) GetCandle(ctx context.Context, exchange, symbol string, interval market.Interval, openTime time.Time) (*market.Candle, bool) {
	key := c.generateCandleKey(exchange, symbol, string(interval), openTime)
	candle, found := c.candleCache.Get(key)
	if !found {
		return nil, false
	}

	cachedCandle, ok := candle.(*market.Candle)
	if !ok {
		log.Warn().Str("exchange", exchange).Str("symbol", symbol).Str("interval", string(interval)).Msg("Invalid candle type in cache")
		return nil, false
	}

	return cachedCandle, true
}

// GetCandleWithError retrieves a candle from the cache with error handling
func (c *StandardCache) GetCandleWithError(ctx context.Context, exchange, symbol string, interval market.Interval, openTime time.Time) (*market.Candle, error) {
	key := c.generateCandleKey(exchange, symbol, string(interval), openTime)
	candle, found := c.candleCache.Get(key)
	if !found {
		return nil, NewCacheKeyNotFoundError(fmt.Sprintf("candle:%s:%s:%s:%s",
			exchange, symbol, string(interval), openTime.Format(time.RFC3339)), nil)
	}

	cachedCandle, ok := candle.(*market.Candle)
	if !ok {
		return nil, NewCacheInvalidTypeError(fmt.Sprintf("candle:%s:%s:%s:%s",
			exchange, symbol, string(interval), openTime.Format(time.RFC3339)), nil)
	}

	return cachedCandle, nil
}

// GetLatestCandle retrieves the most recent candle for a symbol and interval
func (c *StandardCache) GetLatestCandle(ctx context.Context, exchange, symbol string, interval market.Interval) (*market.Candle, bool) {
	key := c.generateLatestCandleKey(exchange, symbol, string(interval))
	candle, found := c.candleCache.Get(key)
	if !found {
		return nil, false
	}

	cachedCandle, ok := candle.(*market.Candle)
	if !ok {
		log.Warn().Str("exchange", exchange).Str("symbol", symbol).Str("interval", string(interval)).Msg("Invalid candle type in cache")
		return nil, false
	}

	return cachedCandle, true
}

// GetLatestCandleWithError retrieves the most recent candle for a symbol and interval with error handling
func (c *StandardCache) GetLatestCandleWithError(ctx context.Context, exchange, symbol string, interval market.Interval) (*market.Candle, error) {
	key := c.generateLatestCandleKey(exchange, symbol, string(interval))
	candle, found := c.candleCache.Get(key)
	if !found {
		return nil, NewCacheKeyNotFoundError(fmt.Sprintf("latest_candle:%s:%s:%s",
			exchange, symbol, string(interval)), nil)
	}

	cachedCandle, ok := candle.(*market.Candle)
	if !ok {
		return nil, NewCacheInvalidTypeError(fmt.Sprintf("latest_candle:%s:%s:%s",
			exchange, symbol, string(interval)), nil)
	}

	return cachedCandle, nil
}

// CacheOrderBook stores an order book in the cache with the default TTL
func (c *StandardCache) CacheOrderBook(orderBook *market.OrderBook) {
	if orderBook == nil {
		return
	}
	key := c.generateOrderBookKey(orderBook.Exchange, orderBook.Symbol)
	c.orderBookCache.Set(key, orderBook, c.defaultTTL)
}

// CacheOrderBookWithTTL stores an order book in the cache with a custom TTL
func (c *StandardCache) CacheOrderBookWithTTL(symbol string, orderBook *market.OrderBook, ttl time.Duration) {
	key := c.generateOrderBookKey(symbol, orderBook.Symbol)
	c.orderBookCache.Set(key, orderBook, ttl)
}

// GetOrderBook retrieves an order book from the cache
func (c *StandardCache) GetOrderBook(ctx context.Context, exchange, symbol string) (*market.OrderBook, bool) {
	key := c.generateOrderBookKey(exchange, symbol)
	orderBook, found := c.orderBookCache.Get(key)
	if !found {
		return nil, false
	}

	cachedOrderBook, ok := orderBook.(*market.OrderBook)
	if !ok {
		log.Warn().Str("exchange", exchange).Str("symbol", symbol).Msg("Invalid order book type in cache")
		return nil, false
	}

	return cachedOrderBook, true
}

// GetOrderBookWithError retrieves an order book from the cache with error handling
func (c *StandardCache) GetOrderBookWithError(ctx context.Context, exchange, symbol string) (*market.OrderBook, error) {
	key := c.generateOrderBookKey(exchange, symbol)
	orderBook, found := c.orderBookCache.Get(key)
	if !found {
		return nil, NewCacheKeyNotFoundError(fmt.Sprintf("orderbook:%s:%s", exchange, symbol), nil)
	}

	cachedOrderBook, ok := orderBook.(*market.OrderBook)
	if !ok {
		return nil, NewCacheInvalidTypeError(fmt.Sprintf("orderbook:%s:%s", exchange, symbol), nil)
	}

	return cachedOrderBook, nil
}

// Clear removes all cached data
func (c *StandardCache) Clear() {
	c.tickerCache.Flush()
	c.candleCache.Flush()
	c.orderBookCache.Flush()
}

// SetTickerExpiry sets the ticker cache expiration duration
func (c *StandardCache) SetTickerExpiry(d time.Duration) {
	c.tickerCache.Flush()
	c.tickerCache = gocache.New(d, d/2)
}

// SetCandleExpiry sets the candle cache expiration duration
func (c *StandardCache) SetCandleExpiry(d time.Duration) {
	c.candleCache.Flush()
	c.candleCache = gocache.New(d, d/2)
}

// SetOrderbookExpiry sets the orderbook cache expiration duration
func (c *StandardCache) SetOrderbookExpiry(d time.Duration) {
	c.orderBookCache.Flush()
	c.orderBookCache = gocache.New(d, d/2)
}

// StartCleanupTask is not needed as go-cache handles cleanup internally
func (c *StandardCache) StartCleanupTask(ctx context.Context, interval time.Duration) {
	// No-op: go-cache handles cleanup with the interval specified at creation time
}

// Helper methods for key generation
func (c *StandardCache) generateTickerKey(exchange, symbol string) string {
	return fmt.Sprintf("ticker:%s:%s", exchange, symbol)
}

func (c *StandardCache) generateLatestTickerKey(symbol string) string {
	return fmt.Sprintf("latest_ticker:%s", symbol)
}

func (c *StandardCache) generateCandleKey(exchange, symbol, interval string, openTime time.Time) string {
	return fmt.Sprintf("candle:%s:%s:%s:%s", exchange, symbol, interval, openTime.Format(time.RFC3339))
}

func (c *StandardCache) generateLatestCandleKey(exchange, symbol, interval string) string {
	return fmt.Sprintf("latest_candle:%s:%s:%s", exchange, symbol, interval)
}

func (c *StandardCache) generateOrderBookKey(exchange, symbol string) string {
	return fmt.Sprintf("orderbook:%s:%s", exchange, symbol)
}

// Helper function to convert cache errors to app errors
func ConvertCacheError(err error) error {
	var cacheErr *CacheError
	if err == nil {
		return nil
	}

	if !apperror.As(err, &cacheErr) {
		return apperror.NewInternal(err)
	}

	switch cacheErr.Code {
	case ErrCacheKeyNotFound:
		return apperror.NewNotFound(cacheErr.Resource, nil, cacheErr)
	case ErrCacheExpired:
		return apperror.NewNotFound(cacheErr.Resource, nil, cacheErr)
	case ErrCacheInvalidType:
		return apperror.NewInternal(cacheErr)
	case ErrCacheNilValue:
		return apperror.NewInvalid("Nil value provided to cache", nil, cacheErr)
	default:
		return apperror.NewInternal(cacheErr)
	}
}

// IsExpired checks if an item is expired in the cache
func (c *StandardCache) IsExpired(cache interface{}, key string) bool {
	if gcache, ok := cache.(*gocache.Cache); ok {
		item, found := gcache.Items()[key]
		if !found {
			return true
		}
		// Expiration is an int64 representing nanoseconds
		return item.Expiration != 0 && time.Now().UnixNano() > item.Expiration
	}

	// If we can't cast to the right type, assume it's expired
	return true
}

// Cache operations with custom TTL handling
// CacheTickerWithCustomTTL stores a ticker with a specific TTL
func (c *StandardCache) CacheTickerWithCustomTTL(ticker *market.Ticker, ttl time.Duration) {
	if ticker == nil {
		return
	}
	key := c.generateTickerKey(ticker.Exchange, ticker.Symbol)
	c.tickerCache.Set(key, ticker, ttl)

	// Also store in latest tickers collection
	latestKey := c.generateLatestTickerKey(ticker.Symbol)
	c.tickerCache.Set(latestKey, ticker, ttl)
}
