package memory

import (
	"time"

	"github.com/neo/crypto-bot/internal/domain/port"
	"github.com/rs/zerolog"
)

// CacheFactory provides methods to create various cache implementations
type CacheFactory struct {
	logger *zerolog.Logger
}

// NewCacheFactory creates a new cache factory
func NewCacheFactory(logger *zerolog.Logger) *CacheFactory {
	return &CacheFactory{
		logger: logger,
	}
}

// NewMarketCache creates a market data cache implementation
func (cf *CacheFactory) NewMarketCache() port.MarketCache {
	// Configure cache TTLs with defaults
	tickerTTL := 5 * time.Minute
	candleTTL := 60 * time.Minute
	orderbookTTL := 30 * time.Second

	return NewMarketCache(cf.logger, tickerTTL, candleTTL, orderbookTTL)
}

// NewTickerCache creates a new specialized ticker cache
func (f *CacheFactory) NewTickerCache() *TickerCache {
	return NewTickerCache(f.logger)
}

// CreateGenericCache is a standalone function to create a generic cache
func CreateGenericCache[T any](logger *zerolog.Logger, cacheTTL time.Duration) port.Cache[T] {
	return NewGenericCache[T](cacheTTL)
}
