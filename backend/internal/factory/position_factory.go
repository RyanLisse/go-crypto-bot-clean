package factory

import (
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/delivery/http/handler"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/persistence/gorm"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/service"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/usecase"
	"github.com/rs/zerolog"
	gormdb "gorm.io/gorm"
)

// PositionFactory creates position-related components
type PositionFactory struct {
	cfg    *PositionFactoryConfig
	logger *zerolog.Logger
	db     *gormdb.DB
}

// PositionFactoryConfig provides configuration for the position factory
type PositionFactoryConfig struct {
	MonitorInterval int // seconds
}

// NewPositionFactory creates a new position factory
func NewPositionFactory(cfg *PositionFactoryConfig, logger *zerolog.Logger, db *gormdb.DB) *PositionFactory {
	return &PositionFactory{
		cfg:    cfg,
		logger: logger,
		db:     db,
	}
}

// CreatePositionRepository creates a position repository
func (f *PositionFactory) CreatePositionRepository() port.PositionRepository {
	return gorm.NewPositionRepository(f.db)
}

// CreateMarketRepository creates a market repository
func (f *PositionFactory) CreateMarketRepository() port.MarketRepository {
	return gorm.NewMarketRepository(f.db, f.logger)
}

// CreateSymbolRepository creates a symbol repository
func (f *PositionFactory) CreateSymbolRepository() port.SymbolRepository {
	return gorm.NewSymbolRepository(f.db, f.logger)
}

// CreatePositionUseCase creates a position use case
func (f *PositionFactory) CreatePositionUseCase(repo port.PositionRepository) usecase.PositionUseCase {
	marketRepo := f.CreateMarketRepository()
	symbolRepo := f.CreateSymbolRepository()
	return usecase.NewPositionUseCase(repo, marketRepo, symbolRepo, *f.logger)
}

// CreatePositionMonitor creates a position monitor service
func (f *PositionFactory) CreatePositionMonitor(
	positionUC usecase.PositionUseCase,
	marketDataService port.MarketDataService,
	tradeUC usecase.TradeUseCase,
) *service.PositionMonitor {
	// Create an adapter that converts port.MarketDataService to service.MarketDataServiceInterface
	marketDataAdapter := service.NewMarketDataServiceAdapter(marketDataService, f.logger)

	monitor := service.NewPositionMonitor(
		positionUC,
		marketDataAdapter,
		tradeUC,
		f.logger,
	)

	if f.cfg.MonitorInterval > 0 {
		monitor.SetInterval(time.Duration(f.cfg.MonitorInterval) * time.Second)
	}

	return monitor
}

// CreatePositionHandler creates a position handler for HTTP API
func (f *PositionFactory) CreatePositionHandler(positionUC usecase.PositionUseCase) *handler.PositionHandler {
	return handler.NewPositionHandler(positionUC, f.logger)
}
