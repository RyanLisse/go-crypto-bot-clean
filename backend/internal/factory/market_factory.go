package factory

import (
	"fmt"

	mexcGateway "github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/gateway/mexc"
	gormAdapter "github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/persistence/gorm"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/config"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/service"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/usecase"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/pkg/platform/mexc"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// MarketFactory creates market data related components
type MarketFactory struct {
	cfg          *config.Config
	logger       *zerolog.Logger
	db           *gorm.DB
	cacheFactory *CacheFactory
	baseService  *service.MarketDataService
}

// NewMarketFactory creates a new MarketFactory
func NewMarketFactory(cfg *config.Config, logger *zerolog.Logger, db *gorm.DB) *MarketFactory {
	return &MarketFactory{
		cfg:          cfg,
		logger:       logger,
		db:           db,
		cacheFactory: NewCacheFactory(cfg, logger),
	}
}

// CreateMarketRepository creates a market data repository
func (f *MarketFactory) CreateMarketRepository() (port.MarketRepository, port.SymbolRepository) {
	// Use the direct repository that implements the new model interfaces directly
	repo := gormAdapter.NewMarketRepositoryDirect(f.db, f.logger)
	// MarketRepositoryDirect implements both interfaces
	return repo, repo
}

// CreateMarketCache creates a market data cache
func (f *MarketFactory) CreateMarketCache() port.MarketCache {
	// Use the new CacheFactory to create a StandardCache
	return f.cacheFactory.CreateMarketCache()
}

// CreateExtendedMarketCache creates a market data cache with error handling capabilities
func (f *MarketFactory) CreateExtendedMarketCache() port.ExtendedMarketCache {
	// Use the CacheFactory to create an extended cache with error handling
	return f.cacheFactory.CreateExtendedMarketCache()
}

// CreateMarketDataUseCase creates the market data use case
func (f *MarketFactory) CreateMarketDataUseCase() (*usecase.MarketDataUseCase, error) {
	marketRepo, symbolRepo := f.CreateMarketRepository()
	cache := f.CreateMarketCache()

	uc := usecase.NewMarketDataUseCase(marketRepo, symbolRepo, cache, f.logger)
	return uc, nil
}

// CreateMEXCClient creates a MEXC API client
func (f *MarketFactory) CreateMEXCClient() port.MEXCClient {
	// Get API credentials from config
	apiKey := f.cfg.MEXC.APIKey
	apiSecret := f.cfg.MEXC.APISecret

	// Create the MEXC client
	return mexc.NewClient(apiKey, apiSecret, f.logger)
}

// CreateMEXCGateway creates a MEXC gateway
func (f *MarketFactory) CreateMEXCGateway() *mexcGateway.MEXCGateway {
	// Create the MEXC client
	mexcClient := f.CreateMEXCClient()

	// Create the MEXC gateway
	return mexcGateway.NewMEXCGateway(mexcClient, f.logger)
}

// CreateMEXCStatusProvider creates a MEXC status provider
func (f *MarketFactory) CreateMEXCStatusProvider() port.StatusProvider {
	// Create the MEXC client
	mexcClient := f.CreateMEXCClient()

	// Create the MEXC status provider
	return mexcGateway.NewMEXCStatusProvider(mexcClient, f.logger)
}

// CreateMarketDataService creates the market data service
func (f *MarketFactory) CreateMarketDataService() *service.MarketDataService {
	marketRepo, symbolRepo := f.CreateMarketRepository()
	cache := f.CreateMarketCache()
	mexcClient := f.CreateMEXCClient()

	return service.NewMarketDataService(
		marketRepo,
		symbolRepo,
		cache,
		mexcClient,
		f.logger,
	)
}

// CreateMarketDataServiceWithErrorHandling creates a MarketDataServiceWithErrorHandling
func (f *MarketFactory) CreateMarketDataServiceWithErrorHandling() (port.MarketDataService, error) {
	// Get dependencies
	marketRepo, symbolRepo := f.CreateMarketRepository()
	f.logger.Debug().Msg("Created market repositories")

	// Check if base service already exists
	if f.baseService != nil {
		return nil, fmt.Errorf("base market data service already exists")
	}

	// Get extended cache with error handling
	cacheService := f.CreateExtendedMarketCache()
	if cacheService == nil {
		return nil, fmt.Errorf("failed to create extended market cache")
	}

	// Get MEXC client
	mexcClient := f.CreateMEXCClient()
	if mexcClient == nil {
		f.logger.Warn().Msg("MEXC client not available, some fallback functionality may not work")
	}

	f.logger.Debug().Msg("Creating market data service with error handling")

	// Create service with error handling using the base service and passing the mexcClient
	baseService := f.CreateMarketDataService()
	if baseService == nil {
		return nil, fmt.Errorf("failed to create base market data service")
	}

	f.baseService = baseService

	return service.NewMarketDataServiceWithErrorHandlingWithService(
		marketRepo,
		symbolRepo,
		cacheService,
		baseService,
		mexcClient,
		f.logger,
	), nil
}
