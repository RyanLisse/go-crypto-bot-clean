package memory

import (
	"context"
	"sync"
	"time"

	"github.com/neo/crypto-bot/internal/domain/model/market"
	"github.com/neo/crypto-bot/internal/domain/port"
	"github.com/rs/zerolog"
)

// TickerCache provides specialized caching for ticker data
type TickerCache struct {
	// Cache for tickers by symbol and exchange
	tickersCache map[string]port.Cache[market.Ticker]

	// Cache for latest tickers by symbol
	latestCache map[string]port.Cache[market.Ticker]

	// Cache for tickers by exchange
	exchangeCache map[string]port.Cache[[]*market.Ticker]

	// Mutex for thread safety
	mu sync.RWMutex

	// TTL for different caches
	tickerTTL   time.Duration
	latestTTL   time.Duration
	exchangeTTL time.Duration

	// Logger
	logger *zerolog.Logger
}

// NewTickerCache creates a new ticker cache with the specified TTLs
func NewTickerCache(logger *zerolog.Logger) *TickerCache {
	return &TickerCache{
		tickersCache:  make(map[string]port.Cache[market.Ticker]),
		latestCache:   make(map[string]port.Cache[market.Ticker]),
		exchangeCache: make(map[string]port.Cache[[]*market.Ticker]),
		tickerTTL:     5 * time.Minute,
		latestTTL:     5 * time.Minute,
		exchangeTTL:   5 * time.Minute,
		logger:        logger,
	}
}

// CacheTicker stores a ticker in cache
func (c *TickerCache) CacheTicker(ticker *market.Ticker) {
	if ticker == nil {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	// Create key for the ticker
	key := getTickerKey(ticker.Exchange, ticker.Symbol)

	// Get or create cache for this ticker
	cache, exists := c.tickersCache[key]
	if !exists {
		cache = NewGenericCache[market.Ticker](c.tickerTTL)
		c.tickersCache[key] = cache
	}

	// Cache the ticker
	cache.Set(ticker)
	c.logger.Debug().Str("exchange", ticker.Exchange).Str("symbol", ticker.Symbol).Msg("Ticker cached")

	// Also update the latest ticker for this symbol
	latestCache, exists := c.latestCache[ticker.Symbol]
	if !exists {
		latestCache = NewGenericCache[market.Ticker](c.latestTTL)
		c.latestCache[ticker.Symbol] = latestCache
	}
	latestCache.Set(ticker)

	// Update the exchange cache
	c.updateExchangeCache(ticker)
}

// GetTicker retrieves a ticker from cache
func (c *TickerCache) GetTicker(ctx context.Context, exchange, symbol string) (*market.Ticker, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	key := getTickerKey(exchange, symbol)
	cache, exists := c.tickersCache[key]
	if !exists {
		return nil, false
	}

	ticker, found := cache.Get()
	if !found {
		c.logger.Debug().Str("exchange", exchange).Str("symbol", symbol).Msg("Ticker not found in cache or expired")
		return nil, false
	}

	return ticker, true
}

// GetAllTickers retrieves all tickers for an exchange from cache
func (c *TickerCache) GetAllTickers(ctx context.Context, exchange string) ([]*market.Ticker, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	cache, exists := c.exchangeCache[exchange]
	if !exists {
		return nil, false
	}

	tickers, found := cache.Get()
	if !found || tickers == nil || len(*tickers) == 0 {
		c.logger.Debug().Str("exchange", exchange).Msg("Exchange tickers not found in cache or expired")
		return nil, false
	}

	return *tickers, true
}

// GetLatestTickers retrieves the most recent tickers across all exchanges
func (c *TickerCache) GetLatestTickers(ctx context.Context) ([]*market.Ticker, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if len(c.latestCache) == 0 {
		return nil, false
	}

	tickers := make([]*market.Ticker, 0, len(c.latestCache))
	for symbol, cache := range c.latestCache {
		ticker, found := cache.Get()
		if found && ticker != nil {
			tickers = append(tickers, ticker)
		} else {
			c.logger.Debug().Str("symbol", symbol).Msg("Latest ticker not found or expired")
		}
	}

	if len(tickers) == 0 {
		return nil, false
	}

	return tickers, true
}

// updateExchangeCache updates the cache of tickers by exchange
func (c *TickerCache) updateExchangeCache(newTicker *market.Ticker) {
	// Get the existing cached tickers for this exchange
	cache, exists := c.exchangeCache[newTicker.Exchange]
	if !exists {
		cache = NewGenericCache[[]*market.Ticker](c.exchangeTTL)
		c.exchangeCache[newTicker.Exchange] = cache
	}

	// Get existing tickers or create a new slice
	var exchangeTickers []*market.Ticker
	existing, found := cache.Get()
	if found && existing != nil {
		exchangeTickers = *existing
	} else {
		exchangeTickers = make([]*market.Ticker, 0)
	}

	// Update or add the ticker in the slice
	updated := false
	for i, ticker := range exchangeTickers {
		if ticker.Symbol == newTicker.Symbol {
			exchangeTickers[i] = newTicker
			updated = true
			break
		}
	}

	if !updated {
		exchangeTickers = append(exchangeTickers, newTicker)
	}

	// Cache the updated slice
	cache.Set(&exchangeTickers)
}

// SetTickerTTL sets the TTL for individual ticker caches
func (c *TickerCache) SetTickerTTL(ttl time.Duration) {
	if ttl <= 0 {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	c.tickerTTL = ttl

	// Update TTL for existing caches
	for _, cache := range c.tickersCache {
		cache.UpdateTTL(ttl)
	}
}

// SetLatestTTL sets the TTL for latest ticker caches
func (c *TickerCache) SetLatestTTL(ttl time.Duration) {
	if ttl <= 0 {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	c.latestTTL = ttl

	// Update TTL for existing caches
	for _, cache := range c.latestCache {
		cache.UpdateTTL(ttl)
	}
}

// SetExchangeTTL sets the TTL for exchange ticker caches
func (c *TickerCache) SetExchangeTTL(ttl time.Duration) {
	if ttl <= 0 {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	c.exchangeTTL = ttl

	// Update TTL for existing caches
	for _, cache := range c.exchangeCache {
		cache.UpdateTTL(ttl)
	}
}

// Clear removes all cached ticker data
func (c *TickerCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Invalidate all caches
	for _, cache := range c.tickersCache {
		cache.Invalidate()
	}
	for _, cache := range c.latestCache {
		cache.Invalidate()
	}
	for _, cache := range c.exchangeCache {
		cache.Invalidate()
	}

	// Reset maps
	c.tickersCache = make(map[string]port.Cache[market.Ticker])
	c.latestCache = make(map[string]port.Cache[market.Ticker])
	c.exchangeCache = make(map[string]port.Cache[[]*market.Ticker])

	c.logger.Debug().Msg("Ticker cache cleared")
}

// Helper function to generate ticker cache keys
func getTickerKey(exchange, symbol string) string {
	return exchange + ":" + symbol
}
