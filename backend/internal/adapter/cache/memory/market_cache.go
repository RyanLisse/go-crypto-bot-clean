package memory

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/neo/crypto-bot/internal/domain/model/market"
	"github.com/neo/crypto-bot/internal/domain/port"
	"github.com/rs/zerolog"
)

// MarketCache provides in-memory caching for market data
type MarketCache struct {
	// Ticker caches
	tickers           map[string]*market.Ticker // key: exchange:symbol
	latestTickers     map[string]*market.Ticker // key: symbol
	tickersByExchange map[string][]*market.Ticker

	// Candle caches
	candles       map[string]*market.Candle // key: exchange:symbol:interval:timestamp
	latestCandles map[string]*market.Candle // key: exchange:symbol:interval

	// OrderBook caches
	orderbooks map[string]*market.OrderBook // key: exchange:symbol

	// Expiration settings
	tickerExpiry    time.Duration
	candleExpiry    time.Duration
	orderbookExpiry time.Duration

	// Cache metadata
	lastUpdated map[string]time.Time

	// Logger
	logger *zerolog.Logger

	// Mutex for thread safety
	mu sync.RWMutex
}

// NewMarketCache creates a new in-memory market data cache
func NewMarketCache(logger *zerolog.Logger) port.MarketCache {
	return &MarketCache{
		tickers:           make(map[string]*market.Ticker),
		latestTickers:     make(map[string]*market.Ticker),
		tickersByExchange: make(map[string][]*market.Ticker),
		candles:           make(map[string]*market.Candle),
		latestCandles:     make(map[string]*market.Candle),
		orderbooks:        make(map[string]*market.OrderBook),
		lastUpdated:       make(map[string]time.Time),
		tickerExpiry:      5 * time.Minute,  // Default ticker expiry
		candleExpiry:      60 * time.Minute, // Default candle expiry
		orderbookExpiry:   30 * time.Second, // Default orderbook expiry
		logger:            logger,
	}
}

// SetTickerExpiry sets the ticker cache expiration duration
func (c *MarketCache) SetTickerExpiry(d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.tickerExpiry = d
}

// SetCandleExpiry sets the candle cache expiration duration
func (c *MarketCache) SetCandleExpiry(d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.candleExpiry = d
}

// SetOrderbookExpiry sets the orderbook cache expiration duration
func (c *MarketCache) SetOrderbookExpiry(d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.orderbookExpiry = d
}

// CacheTicker stores a ticker in cache
func (c *MarketCache) CacheTicker(ticker *market.Ticker) {
	if ticker == nil {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	key := c.getTickerKey(ticker.Exchange, ticker.Symbol)
	c.tickers[key] = ticker
	c.latestTickers[ticker.Symbol] = ticker
	c.lastUpdated[key] = time.Now()

	// Update the tickers by exchange cache
	exchangeTickers, exists := c.tickersByExchange[ticker.Exchange]
	if !exists {
		exchangeTickers = make([]*market.Ticker, 0)
	}

	// Check if we already have this symbol in the exchange tickers
	found := false
	for i, t := range exchangeTickers {
		if t.Symbol == ticker.Symbol {
			exchangeTickers[i] = ticker
			found = true
			break
		}
	}

	if !found {
		exchangeTickers = append(exchangeTickers, ticker)
	}

	c.tickersByExchange[ticker.Exchange] = exchangeTickers
	c.logger.Debug().Str("exchange", ticker.Exchange).Str("symbol", ticker.Symbol).Msg("Ticker cached")
}

// GetTicker retrieves a ticker from cache
func (c *MarketCache) GetTicker(ctx context.Context, exchange, symbol string) (*market.Ticker, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	key := c.getTickerKey(exchange, symbol)
	ticker, exists := c.tickers[key]
	if !exists {
		return nil, false
	}

	// Check if ticker has expired
	lastUpdate, updated := c.lastUpdated[key]
	if !updated || time.Since(lastUpdate) > c.tickerExpiry {
		c.logger.Debug().Str("exchange", exchange).Str("symbol", symbol).Msg("Ticker cache expired")
		return nil, false
	}

	return ticker, true
}

// GetAllTickers retrieves all tickers for an exchange from cache
func (c *MarketCache) GetAllTickers(ctx context.Context, exchange string) ([]*market.Ticker, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	tickers, exists := c.tickersByExchange[exchange]
	if !exists || len(tickers) == 0 {
		return nil, false
	}

	// Check if any exchange tickers have expired
	key := "exchange:" + exchange
	lastUpdate, updated := c.lastUpdated[key]
	if !updated || time.Since(lastUpdate) > c.tickerExpiry {
		c.logger.Debug().Str("exchange", exchange).Msg("Exchange tickers cache expired")
		return nil, false
	}

	return tickers, true
}

// GetLatestTickers retrieves the most recent tickers across all exchanges
func (c *MarketCache) GetLatestTickers(ctx context.Context) ([]*market.Ticker, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if len(c.latestTickers) == 0 {
		return nil, false
	}

	tickers := make([]*market.Ticker, 0, len(c.latestTickers))
	for _, ticker := range c.latestTickers {
		tickers = append(tickers, ticker)
	}

	return tickers, true
}

// CacheCandle stores a candle in cache
func (c *MarketCache) CacheCandle(candle *market.Candle) {
	if candle == nil {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	key := c.getCandleKey(candle.Exchange, candle.Symbol, string(candle.Interval), candle.OpenTime)
	c.candles[key] = candle
	c.lastUpdated[key] = time.Now()

	// Update latest candle for this symbol and interval
	latestKey := c.getLatestCandleKey(candle.Exchange, candle.Symbol, string(candle.Interval))

	// Only update if this is a newer candle or there's no existing candle
	existingLatest, exists := c.latestCandles[latestKey]
	if !exists || candle.OpenTime.After(existingLatest.OpenTime) {
		c.latestCandles[latestKey] = candle
		c.logger.Debug().
			Str("exchange", candle.Symbol).
			Str("symbol", candle.Symbol).
			Str("interval", string(candle.Interval)).
			Time("openTime", candle.OpenTime).
			Msg("Latest candle cached")
	}
}

// GetCandle retrieves a candle from cache
func (c *MarketCache) GetCandle(ctx context.Context, exchange, symbol string, interval market.Interval, openTime time.Time) (*market.Candle, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	key := c.getCandleKey(exchange, symbol, string(interval), openTime)
	candle, exists := c.candles[key]
	if !exists {
		return nil, false
	}

	// Check if candle has expired
	lastUpdate, updated := c.lastUpdated[key]
	if !updated || time.Since(lastUpdate) > c.candleExpiry {
		c.logger.Debug().
			Str("exchange", exchange).
			Str("symbol", symbol).
			Str("interval", string(interval)).
			Time("openTime", openTime).
			Msg("Candle cache expired")
		return nil, false
	}

	return candle, true
}

// GetLatestCandle retrieves the most recent candle for a symbol and interval
func (c *MarketCache) GetLatestCandle(ctx context.Context, exchange, symbol string, interval market.Interval) (*market.Candle, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	key := c.getLatestCandleKey(exchange, symbol, string(interval))
	candle, exists := c.latestCandles[key]
	if !exists {
		return nil, false
	}

	// Check if candle has expired
	lastUpdate, updated := c.lastUpdated[key]
	if !updated || time.Since(lastUpdate) > c.candleExpiry {
		c.logger.Debug().
			Str("exchange", exchange).
			Str("symbol", symbol).
			Str("interval", string(interval)).
			Msg("Latest candle cache expired")
		return nil, false
	}

	return candle, true
}

// CacheOrderBook stores an orderbook in cache
func (c *MarketCache) CacheOrderBook(orderbook *market.OrderBook) {
	if orderbook == nil {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	key := c.getOrderBookKey(orderbook.Exchange, orderbook.Symbol)
	c.orderbooks[key] = orderbook
	c.lastUpdated[key] = time.Now()
	c.logger.Debug().
		Str("exchange", orderbook.Exchange).
		Str("symbol", orderbook.Symbol).
		Msg("OrderBook cached")
}

// GetOrderBook retrieves an orderbook from cache
func (c *MarketCache) GetOrderBook(ctx context.Context, exchange, symbol string) (*market.OrderBook, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	key := c.getOrderBookKey(exchange, symbol)
	orderbook, exists := c.orderbooks[key]
	if !exists {
		return nil, false
	}

	// Check if orderbook has expired
	lastUpdate, updated := c.lastUpdated[key]
	if !updated || time.Since(lastUpdate) > c.orderbookExpiry {
		c.logger.Debug().
			Str("exchange", exchange).
			Str("symbol", symbol).
			Msg("OrderBook cache expired")
		return nil, false
	}

	return orderbook, true
}

// Clear removes all cached data
func (c *MarketCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.tickers = make(map[string]*market.Ticker)
	c.latestTickers = make(map[string]*market.Ticker)
	c.tickersByExchange = make(map[string][]*market.Ticker)
	c.candles = make(map[string]*market.Candle)
	c.latestCandles = make(map[string]*market.Candle)
	c.orderbooks = make(map[string]*market.OrderBook)
	c.lastUpdated = make(map[string]time.Time)
	c.logger.Debug().Msg("Cache cleared")
}

// Helper method to generate ticker cache keys
func (c *MarketCache) getTickerKey(exchange, symbol string) string {
	return exchange + ":" + symbol
}

// Helper method to generate candle cache keys
func (c *MarketCache) getCandleKey(exchange, symbol, interval string, openTime time.Time) string {
	return exchange + ":" + symbol + ":" + interval + ":" + openTime.Format(time.RFC3339)
}

// Helper method to generate latest candle cache keys
func (c *MarketCache) getLatestCandleKey(exchange, symbol, interval string) string {
	return "latest:" + exchange + ":" + symbol + ":" + interval
}

// Helper method to generate orderbook cache keys
func (c *MarketCache) getOrderBookKey(exchange, symbol string) string {
	return exchange + ":" + symbol
}

// StartCleanupTask starts a periodic task to clean up expired cache entries
func (c *MarketCache) StartCleanupTask(ctx context.Context, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				c.cleanupExpiredEntries()
			}
		}
	}()
}

// cleanupExpiredEntries removes expired entries from all caches
func (c *MarketCache) cleanupExpiredEntries() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()

	// Clean up tickers
	for key, lastUpdate := range c.lastUpdated {
		if strings.HasPrefix(key, "exchange:") {
			// Skip exchange keys
			continue
		}

		// Determine which cache this key belongs to
		if strings.Count(key, ":") == 1 {
			// Ticker key: exchange:symbol
			if now.Sub(lastUpdate) > c.tickerExpiry {
				parts := strings.Split(key, ":")
				if len(parts) >= 2 {
					exchange, symbol := parts[0], parts[1]
					delete(c.tickers, key)

					// Also update the exchange tickers
					if exchangeTickers, ok := c.tickersByExchange[exchange]; ok {
						for i, ticker := range exchangeTickers {
							if ticker.Symbol == symbol {
								// Remove this ticker from the slice
								exchangeTickers = append(exchangeTickers[:i], exchangeTickers[i+1:]...)
								break
							}
						}
						if len(exchangeTickers) > 0 {
							c.tickersByExchange[exchange] = exchangeTickers
						} else {
							delete(c.tickersByExchange, exchange)
						}
					}

					// Also clean up from latestTickers if it's there
					if ticker, ok := c.latestTickers[symbol]; ok {
						if ticker.Exchange == exchange {
							delete(c.latestTickers, symbol)
						}
					}

					delete(c.lastUpdated, key)
					c.logger.Debug().Str("key", key).Msg("Expired ticker removed from cache")
				}
			}
		} else if strings.Count(key, ":") == 3 && !strings.HasPrefix(key, "latest:") {
			// Candle key: exchange:symbol:interval:timestamp
			if now.Sub(lastUpdate) > c.candleExpiry {
				delete(c.candles, key)
				delete(c.lastUpdated, key)
				c.logger.Debug().Str("key", key).Msg("Expired candle removed from cache")
			}
		}
	}

	// Clean up latest candles separately
	for key, lastUpdate := range c.lastUpdated {
		if strings.HasPrefix(key, "latest:") && now.Sub(lastUpdate) > c.candleExpiry {
			delete(c.latestCandles, key)
			delete(c.lastUpdated, key)
			c.logger.Debug().Str("key", key).Msg("Expired latest candle removed from cache")
		}
	}

	// Clean up orderbooks
	for key, lastUpdate := range c.lastUpdated {
		if strings.Count(key, ":") == 1 {
			// Could be either a ticker or orderbook, but we've already handled tickers
			if _, exists := c.tickers[key]; !exists {
				if now.Sub(lastUpdate) > c.orderbookExpiry {
					delete(c.orderbooks, key)
					delete(c.lastUpdated, key)
					c.logger.Debug().Str("key", key).Msg("Expired orderbook removed from cache")
				}
			}
		}
	}
}
