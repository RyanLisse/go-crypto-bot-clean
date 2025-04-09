package risk

import (
	"database/sql"

	"go-crypto-bot-clean/backend/internal/domain/risk/controls"
	"go-crypto-bot-clean/backend/internal/platform/database"
)

// Factory creates and configures a RiskService
type Factory struct {
	db            *sql.DB
	priceService  controls.PriceService
	accountService controls.AccountService
	positionRepo  controls.PositionRepository
	tradeRepo     controls.TradeRepository
	logger        Logger
}

// NewFactory creates a new Factory
func NewFactory(
	db *sql.DB,
	priceService controls.PriceService,
	accountService controls.AccountService,
	positionRepo controls.PositionRepository,
	tradeRepo controls.TradeRepository,
	logger Logger,
) *Factory {
	return &Factory{
		db:            db,
		priceService:  priceService,
		accountService: accountService,
		positionRepo:  positionRepo,
		tradeRepo:     tradeRepo,
		logger:        logger,
	}
}

// Create creates a new RiskService
func (f *Factory) Create() RiskService {
	// Create repositories
	balanceRepo := database.NewSQLiteBalanceHistoryRepository(f.db)

	// Create control components
	positionSizer := controls.NewPositionSizer(f.priceService, f.logger)
	drawdownMonitor := controls.NewDrawdownMonitor(balanceRepo, f.logger)
	exposureMonitor := controls.NewExposureMonitor(f.positionRepo, f.accountService, f.logger)
	dailyLimitMonitor := controls.NewDailyLimitMonitor(f.tradeRepo, f.accountService, f.logger)

	// Create risk manager
	riskManager := NewRiskManager(
		balanceRepo,
		positionSizer,
		drawdownMonitor,
		exposureMonitor,
		dailyLimitMonitor,
		f.logger,
	)

	return riskManager
}
