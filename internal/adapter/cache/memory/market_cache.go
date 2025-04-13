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
func NewMarketCache(
	logger *zerolog.Logger,
	tickerTTL time.Duration,
	candleTTL time.Duration,
	orderbookTTL time.Duration,
) port.MarketCache {
	return &MarketCache{
		tickers:           make(map[string]*market.Ticker),
		latestTickers:     make(map[string]*market.Ticker),
		tickersByExchange: make(map[string][]*market.Ticker),
		candles:           make(map[string]*market.Candle),
		latestCandles:     make(map[string]*market.Candle),
		orderbooks:        make(map[string]*market.OrderBook),
		lastUpdated:       make(map[string]time.Time),
		tickerExpiry:      tickerTTL,
		candleExpiry:      candleTTL,
		orderbookExpiry:   orderbookTTL,
		logger:            logger,
	}
}

// CacheTicker stores a ticker in the cache
func (mc *MarketCache) CacheTicker(ticker *market.Ticker) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	if ticker == nil {
		return
	}

	key := mc.tickerKey(ticker.Exchange, ticker.Symbol)
	mc.tickers[key] = ticker
	mc.latestTickers[ticker.Symbol] = ticker
	mc.lastUpdated[key] = time.Now()

	// Update the by-exchange index
	exchangeTickers, exists := mc.tickersByExchange[ticker.Exchange]
	if !exists {
		exchangeTickers = make([]*market.Ticker, 0)
	}

	// Check if we already have this symbol for this exchange
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

	mc.tickersByExchange[ticker.Exchange] = exchangeTickers

	mc.logger.Debug().
		Str("exchange", ticker.Exchange).
		Str("symbol", ticker.Symbol).
		Float64("price", ticker.Price).
		Msg("Ticker cached")
}

// GetTicker retrieves a ticker from the cache if it exists and is not expired
func (mc *MarketCache) GetTicker(ctx context.Context, exchange, symbol string) (*market.Ticker, bool) {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	key := mc.tickerKey(exchange, symbol)
	ticker, exists := mc.tickers[key]
	if !exists {
		return nil, false
	}

	// Check if the ticker has expired
	lastUpdate, exists := mc.lastUpdated[key]
	if !exists || time.Since(lastUpdate) > mc.tickerExpiry {
		return nil, false
	}

	return ticker, true
}

// GetAllTickers retrieves all tickers for a specific exchange
func (mc *MarketCache) GetAllTickers(ctx context.Context, exchange string) ([]*market.Ticker, bool) {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	tickers, exists := mc.tickersByExchange[exchange]
	if !exists || len(tickers) == 0 {
		return nil, false
	}

	// Make a copy to avoid external modifications
	result := make([]*market.Ticker, len(tickers))
	for i, ticker := range tickers {
		// Check if expired
		key := mc.tickerKey(ticker.Exchange, ticker.Symbol)
		lastUpdate, exists := mc.lastUpdated[key]
		if !exists || time.Since(lastUpdate) > mc.tickerExpiry {
			return nil, false
		}
		result[i] = ticker
	}

	return result, true
}

// GetLatestTickers retrieves the latest tickers for all symbols
func (mc *MarketCache) GetLatestTickers(ctx context.Context) ([]*market.Ticker, bool) {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	if len(mc.latestTickers) == 0 {
		return nil, false
	}

	// Make a copy of the latest tickers
	result := make([]*market.Ticker, 0, len(mc.latestTickers))
	for _, ticker := range mc.latestTickers {
		// Check if expired
		key := mc.tickerKey(ticker.Exchange, ticker.Symbol)
		lastUpdate, exists := mc.lastUpdated[key]
		if exists && time.Since(lastUpdate) <= mc.tickerExpiry {
			result = append(result, ticker)
		}
	}

	if len(result) == 0 {
		return nil, false
	}

	return result, true
}

// CacheCandle stores a candle in the cache
func (mc *MarketCache) CacheCandle(candle *market.Candle) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	if candle == nil {
		return
	}

	// Cache the candle with its specific timestamp
	key := mc.candleKey(candle.Exchange, candle.Symbol, candle.Interval, candle.OpenTime)
	mc.candles[key] = candle
	mc.lastUpdated[key] = time.Now()

	// Update the latest candle for this symbol+interval
	latestKey := mc.latestCandleKey(candle.Exchange, candle.Symbol, candle.Interval)

	// Only update if this is a newer candle or we don't have one yet
	latestCandle, exists := mc.latestCandles[latestKey]
	if !exists || latestCandle.OpenTime.Before(candle.OpenTime) {
		mc.latestCandles[latestKey] = candle
	}

	mc.logger.Debug().
		Str("exchange", candle.Exchange).
		Str("symbol", candle.Symbol).
		Str("interval", string(candle.Interval)).
		Time("openTime", candle.OpenTime).
		Msg("Candle cached")
}

// GetCandle retrieves a specific candle from the cache
func (mc *MarketCache) GetCandle(ctx context.Context, exchange, symbol string, interval market.Interval, openTime time.Time) (*market.Candle, bool) {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	key := mc.candleKey(exchange, symbol, interval, openTime)
	candle, exists := mc.candles[key]
	if !exists {
		return nil, false
	}

	// Check if the candle has expired
	lastUpdate, exists := mc.lastUpdated[key]
	if !exists || time.Since(lastUpdate) > mc.candleExpiry {
		return nil, false
	}

	return candle, true
}

// GetLatestCandle retrieves the most recent candle for a symbol and interval
func (mc *MarketCache) GetLatestCandle(ctx context.Context, exchange, symbol string, interval market.Interval) (*market.Candle, bool) {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	key := mc.latestCandleKey(exchange, symbol, interval)
	candle, exists := mc.latestCandles[key]
	if !exists {
		return nil, false
	}

	// Check if the candle has expired
	candleKey := mc.candleKey(candle.Exchange, candle.Symbol, candle.Interval, candle.OpenTime)
	lastUpdate, exists := mc.lastUpdated[candleKey]
	if !exists || time.Since(lastUpdate) > mc.candleExpiry {
		return nil, false
	}

	return candle, true
}

// CacheOrderBook stores an order book in the cache
func (mc *MarketCache) CacheOrderBook(orderbook *market.OrderBook) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	if orderbook == nil {
		return
	}

	key := mc.orderbookKey(orderbook.Exchange, orderbook.Symbol)
	mc.orderbooks[key] = orderbook
	mc.lastUpdated[key] = time.Now()

	mc.logger.Debug().
		Str("exchange", orderbook.Exchange).
		Str("symbol", orderbook.Symbol).
		Int("bids", len(orderbook.Bids)).
		Int("asks", len(orderbook.Asks)).
		Msg("OrderBook cached")
}

// GetOrderBook retrieves an order book from the cache
func (mc *MarketCache) GetOrderBook(ctx context.Context, exchange, symbol string) (*market.OrderBook, bool) {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	key := mc.orderbookKey(exchange, symbol)
	orderbook, exists := mc.orderbooks[key]
	if !exists {
		return nil, false
	}

	// Check if the order book has expired
	lastUpdate, exists := mc.lastUpdated[key]
	if !exists || time.Since(lastUpdate) > mc.orderbookExpiry {
		return nil, false
	}

	return orderbook, true
}

// Clear removes all items from all caches
func (mc *MarketCache) Clear() {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	mc.tickers = make(map[string]*market.Ticker)
	mc.latestTickers = make(map[string]*market.Ticker)
	mc.tickersByExchange = make(map[string][]*market.Ticker)
	mc.candles = make(map[string]*market.Candle)
	mc.latestCandles = make(map[string]*market.Candle)
	mc.orderbooks = make(map[string]*market.OrderBook)
	mc.lastUpdated = make(map[string]time.Time)

	mc.logger.Info().Msg("All caches cleared")
}

// SetTickerExpiry updates the expiration time for tickers
func (mc *MarketCache) SetTickerExpiry(d time.Duration) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	mc.tickerExpiry = d
	mc.logger.Info().Dur("duration", d).Msg("Ticker expiry updated")
}

// SetCandleExpiry updates the expiration time for candles
func (mc *MarketCache) SetCandleExpiry(d time.Duration) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	mc.candleExpiry = d
	mc.logger.Info().Dur("duration", d).Msg("Candle expiry updated")
}

// SetOrderbookExpiry updates the expiration time for orderbooks
func (mc *MarketCache) SetOrderbookExpiry(d time.Duration) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	mc.orderbookExpiry = d
	mc.logger.Info().Dur("duration", d).Msg("Orderbook expiry updated")
}

// StartCleanupTask starts a background task that periodically removes expired items
func (mc *MarketCache) StartCleanupTask(ctx context.Context, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				mc.logger.Info().Msg("Cleanup task canceled")
				return
			case <-ticker.C:
				mc.cleanup()
			}
		}
	}()

	mc.logger.Info().Dur("interval", interval).Msg("Started cache cleanup task")
}

// cleanup removes expired items from all caches
func (mc *MarketCache) cleanup() {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	now := time.Now()
	expiredCount := 0

	// Clean up tickers
	for key, lastUpdate := range mc.lastUpdated {
		if strings.HasPrefix(key, "ticker:") {
			if now.Sub(lastUpdate) > mc.tickerExpiry {
				parts := strings.Split(key, ":")
				if len(parts) >= 3 {
					exchange, symbol := parts[1], parts[2]
					delete(mc.tickers, key)

					// Also update tickersByExchange
					if tickers, exists := mc.tickersByExchange[exchange]; exists {
						updatedTickers := make([]*market.Ticker, 0, len(tickers))
						for _, t := range tickers {
							if t.Symbol != symbol {
								updatedTickers = append(updatedTickers, t)
							}
						}
						if len(updatedTickers) == 0 {
							delete(mc.tickersByExchange, exchange)
						} else {
							mc.tickersByExchange[exchange] = updatedTickers
						}
					}

					// If this is the latest ticker for the symbol, remove it
					if lt, exists := mc.latestTickers[symbol]; exists && lt.Exchange == exchange {
						delete(mc.latestTickers, symbol)
					}

					delete(mc.lastUpdated, key)
					expiredCount++
				}
			}
		} else if strings.HasPrefix(key, "candle:") {
			if now.Sub(lastUpdate) > mc.candleExpiry {
				parts := strings.Split(key, ":")
				if len(parts) >= 5 {
					exchange, symbol, interval := parts[1], parts[2], parts[3]
					delete(mc.candles, key)

					// Check if this is the latest candle and remove it if necessary
					latestKey := mc.latestCandleKey(exchange, symbol, market.Interval(interval))
					if lc, exists := mc.latestCandles[latestKey]; exists {
						candleKey := mc.candleKey(lc.Exchange, lc.Symbol, lc.Interval, lc.OpenTime)
						if candleKey == key {
							delete(mc.latestCandles, latestKey)
						}
					}

					delete(mc.lastUpdated, key)
					expiredCount++
				}
			}
		} else if strings.HasPrefix(key, "orderbook:") {
			if now.Sub(lastUpdate) > mc.orderbookExpiry {
				delete(mc.orderbooks, key)
				delete(mc.lastUpdated, key)
				expiredCount++
			}
		}
	}

	if expiredCount > 0 {
		mc.logger.Debug().Int("count", expiredCount).Msg("Expired cache items removed")
	}
}

// Helper methods for cache key generation

func (mc *MarketCache) tickerKey(exchange, symbol string) string {
	return "ticker:" + exchange + ":" + symbol
}

func (mc *MarketCache) candleKey(exchange, symbol string, interval market.Interval, openTime time.Time) string {
	return "candle:" + exchange + ":" + symbol + ":" + string(interval) + ":" + openTime.Format(time.RFC3339)
}

func (mc *MarketCache) latestCandleKey(exchange, symbol string, interval market.Interval) string {
	return "latest_candle:" + exchange + ":" + symbol + ":" + string(interval)
}

func (mc *MarketCache) orderbookKey(exchange, symbol string) string {
	return "orderbook:" + exchange + ":" + symbol
}
