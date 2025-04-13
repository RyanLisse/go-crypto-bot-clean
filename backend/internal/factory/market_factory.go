package factory

import (
	"time"

	cacheAdapter "github.com/neo/crypto-bot/internal/adapter/cache/memory"
	gormAdapter "github.com/neo/crypto-bot/internal/adapter/persistence/gorm"
	"github.com/neo/crypto-bot/internal/config"
	"github.com/neo/crypto-bot/internal/domain/port"
	"github.com/neo/crypto-bot/internal/domain/service"
	"github.com/neo/crypto-bot/internal/usecase"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// MarketFactory creates market data related components
type MarketFactory struct {
	cfg    *config.Config
	logger *zerolog.Logger
	db     *gorm.DB
}

// NewMarketFactory creates a new MarketFactory
func NewMarketFactory(cfg *config.Config, logger *zerolog.Logger, db *gorm.DB) *MarketFactory {
	return &MarketFactory{
		cfg:    cfg,
		logger: logger,
		db:     db,
	}
}

// CreateMarketRepository creates a market data repository
func (f *MarketFactory) CreateMarketRepository() (port.MarketRepository, port.SymbolRepository) {
	repo := gormAdapter.NewMarketRepository(f.db, f.logger)
	// GORM MarketRepository implements both interfaces
	return repo, repo
}

// CreateMarketCache creates a market data cache
func (f *MarketFactory) CreateMarketCache() port.MarketCache {
	// Create the cache instance
	cache := cacheAdapter.NewMarketCache(f.logger)

	// Define default cache TTLs
	tickerTTL := 5 * time.Minute
	candleTTL := 15 * time.Minute
	orderbookTTL := 30 * time.Second

	// Use market config TTLs if configured
	if f.cfg.Market.Cache.TickerTTL > 0 {
		tickerTTL = time.Duration(f.cfg.Market.Cache.TickerTTL) * time.Second
	}
	if f.cfg.Market.Cache.CandleTTL > 0 {
		candleTTL = time.Duration(f.cfg.Market.Cache.CandleTTL) * time.Second
	}
	if f.cfg.Market.Cache.OrderbookTTL > 0 {
		orderbookTTL = time.Duration(f.cfg.Market.Cache.OrderbookTTL) * time.Second
	}

	// Configure the cache TTLs
	cache.SetTickerExpiry(tickerTTL)
	cache.SetCandleExpiry(candleTTL)
	cache.SetOrderbookExpiry(orderbookTTL)

	return cache
}

// CreateMarketDataUseCase creates the market data use case
func (f *MarketFactory) CreateMarketDataUseCase() (*usecase.MarketDataUseCase, error) {
	marketRepo, symbolRepo := f.CreateMarketRepository()
	cache := f.CreateMarketCache()

	uc := usecase.NewMarketDataUseCase(marketRepo, symbolRepo, cache, f.logger)
	return uc, nil
}

// CreateMarketDataService creates the market data service
func (f *MarketFactory) CreateMarketDataService(
	marketRepo port.MarketRepository,
	symbolRepo port.SymbolRepository,
	cache port.MarketCache,
	mexcAPI port.MexcAPI,
) *service.MarketDataService {
	return service.NewMarketDataService(
		marketRepo,
		symbolRepo,
		cache,
		mexcAPI,
		f.logger,
	)
}
