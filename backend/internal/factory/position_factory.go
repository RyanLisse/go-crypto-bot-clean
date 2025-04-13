package factory

import (
	"github.com/neo/crypto-bot/internal/adapter/persistence/gorm"
	"github.com/neo/crypto-bot/internal/config"
	"github.com/neo/crypto-bot/internal/domain/port"
	"github.com/neo/crypto-bot/internal/domain/service"
	"github.com/neo/crypto-bot/internal/usecase"
	"github.com/rs/zerolog"
	gormdb "gorm.io/gorm"
)

// PositionFactory creates position management related components
type PositionFactory struct {
	cfg    *config.Config
	logger *zerolog.Logger
	db     *gormdb.DB
}

// NewPositionFactory creates a new PositionFactory
func NewPositionFactory(cfg *config.Config, logger *zerolog.Logger, db *gormdb.DB) *PositionFactory {
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

// CreatePositionUseCase creates the position use case
func (f *PositionFactory) CreatePositionUseCase(
	marketRepo port.MarketRepository,
	symbolRepo port.SymbolRepository,
) (usecase.PositionUseCase, error) {
	positionRepo := f.CreatePositionRepository()

	uc := usecase.NewPositionUseCase(
		positionRepo,
		marketRepo,
		symbolRepo,
		*f.logger,
	)

	return uc, nil
}

// CreatePositionMonitor creates the position monitor service
func (f *PositionFactory) CreatePositionMonitor(
	positionUC usecase.PositionUseCase,
	marketService *service.MarketDataService,
	tradeUC usecase.TradeUseCase,
) *service.PositionMonitor {
	return service.NewPositionMonitor(
		positionUC,
		marketService,
		tradeUC,
		f.logger,
	)
}
