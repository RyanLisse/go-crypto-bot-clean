package factory

import (
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/cache/standard"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/config"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/rs/zerolog"
)

// CacheFactory creates cache instances
type CacheFactory struct {
	config *config.Config
	logger *zerolog.Logger
}

// NewCacheFactory creates a new CacheFactory
func NewCacheFactory(config *config.Config, logger *zerolog.Logger) *CacheFactory {
	return &CacheFactory{
		config: config,
		logger: logger,
	}
}

// CreateMarketCache creates a new MarketCache instance using the go-cache library
func (f *CacheFactory) CreateMarketCache() port.MarketCache {
	// Default cache configuration
	defaultTTL := 5 * time.Minute
	cleanupInterval := 10 * time.Minute

	// Use configuration if available
	if f.config != nil {
		// Default TTL and cleanup are hardcoded, only the specific TTLs are configurable
		f.logger.Info().
			Dur("defaultTTL", defaultTTL).
			Dur("cleanupInterval", cleanupInterval).
			Msg("Creating standard market cache")
	}

	cache := standard.NewStandardCache(defaultTTL, cleanupInterval)

	// Apply specific TTLs if configured
	if f.config != nil {
		if f.config.Market.Cache.TickerTTL > 0 {
			tickerTTL := time.Duration(f.config.Market.Cache.TickerTTL) * time.Second
			cache.SetTickerExpiry(tickerTTL)
			f.logger.Debug().Dur("ttl", tickerTTL).Msg("Set ticker cache TTL")
		}

		if f.config.Market.Cache.CandleTTL > 0 {
			candleTTL := time.Duration(f.config.Market.Cache.CandleTTL) * time.Second
			cache.SetCandleExpiry(candleTTL)
			f.logger.Debug().Dur("ttl", candleTTL).Msg("Set candle cache TTL")
		}

		if f.config.Market.Cache.OrderbookTTL > 0 {
			orderbookTTL := time.Duration(f.config.Market.Cache.OrderbookTTL) * time.Second
			cache.SetOrderbookExpiry(orderbookTTL)
			f.logger.Debug().Dur("ttl", orderbookTTL).Msg("Set orderbook cache TTL")
		}
	}

	return cache
}

// CreateExtendedMarketCache creates a new ExtendedMarketCache instance with error handling capabilities
func (f *CacheFactory) CreateExtendedMarketCache() port.ExtendedMarketCache {
	// Default cache configuration
	defaultTTL := 5 * time.Minute
	cleanupInterval := 10 * time.Minute

	// Use configuration if available
	if f.config != nil {
		// Default TTL and cleanup are hardcoded, only the specific TTLs are configurable
		f.logger.Info().
			Dur("defaultTTL", defaultTTL).
			Dur("cleanupInterval", cleanupInterval).
			Msg("Creating extended market cache with error handling")
	}

	cache := standard.NewStandardCache(defaultTTL, cleanupInterval)

	// Apply specific TTLs if configured
	if f.config != nil {
		if f.config.Market.Cache.TickerTTL > 0 {
			tickerTTL := time.Duration(f.config.Market.Cache.TickerTTL) * time.Second
			cache.SetTickerExpiry(tickerTTL)
			f.logger.Debug().Dur("ttl", tickerTTL).Msg("Set ticker cache TTL")
		}

		if f.config.Market.Cache.CandleTTL > 0 {
			candleTTL := time.Duration(f.config.Market.Cache.CandleTTL) * time.Second
			cache.SetCandleExpiry(candleTTL)
			f.logger.Debug().Dur("ttl", candleTTL).Msg("Set candle cache TTL")
		}

		if f.config.Market.Cache.OrderbookTTL > 0 {
			orderbookTTL := time.Duration(f.config.Market.Cache.OrderbookTTL) * time.Second
			cache.SetOrderbookExpiry(orderbookTTL)
			f.logger.Debug().Dur("ttl", orderbookTTL).Msg("Set orderbook cache TTL")
		}
	}

	return cache
}
