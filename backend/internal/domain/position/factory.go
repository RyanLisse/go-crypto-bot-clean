package position

import (
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"go-crypto-bot-clean/backend/internal/domain/interfaces"
	"go-crypto-bot-clean/backend/internal/domain/position/management"
	"go-crypto-bot-clean/backend/internal/domain/service"
	"go-crypto-bot-clean/backend/internal/platform/database/repositories"
)

// Factory creates position management components
type Factory struct {
	db              *sqlx.DB
	tradeService    interfaces.TradeService
	exchangeService interfaces.ExchangeService
	logger          *zap.Logger
}

// NewFactory creates a new position factory
func NewFactory(
	db *sqlx.DB,
	tradeService interfaces.TradeService,
	exchangeService interfaces.ExchangeService,
	logger *zap.Logger,
) *Factory {
	return &Factory{
		db:              db,
		tradeService:    tradeService,
		exchangeService: exchangeService,
		logger:          logger,
	}
}

// CreatePositionService creates a new position service
func (f *Factory) CreatePositionService() PositionService {
	// Create repositories
	positionRepo := repositories.NewSQLitePositionRepository(f.db)

	// Create services
	orderService := service.NewOrderService(f.tradeService, f.logger)
	priceService := service.NewPriceService(f.exchangeService, f.logger)

	// Create position manager
	return management.NewPositionManager(positionRepo, orderService, priceService, f.logger)
}

// CreatePositionRepository creates a new position repository
func (f *Factory) CreatePositionRepository() interfaces.PositionRepository {
	return repositories.NewSQLitePositionRepository(f.db)
}
