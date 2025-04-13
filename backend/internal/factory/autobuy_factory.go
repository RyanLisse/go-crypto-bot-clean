package factory

import (
	"github.com/neo/crypto-bot/internal/adapter/handler"
	"github.com/neo/crypto-bot/internal/adapter/persistence/gorm"
	"github.com/neo/crypto-bot/internal/domain/model"
	"github.com/neo/crypto-bot/internal/domain/port"
	"github.com/neo/crypto-bot/internal/usecase"
	"github.com/rs/zerolog"
	gormdb "gorm.io/gorm"
)

// AutoBuyFactory creates and provides auto-buy related components
type AutoBuyFactory struct {
	config        *model.Config
	logger        *zerolog.Logger
	db            *gormdb.DB
	marketFactory *MarketFactory
	tradeFactory  *TradeFactory
}

// NewAutoBuyFactory creates a new AutoBuyFactory
func NewAutoBuyFactory(
	config *model.Config,
	logger *zerolog.Logger,
	db *gormdb.DB,
	marketFactory *MarketFactory,
	tradeFactory *TradeFactory,
) *AutoBuyFactory {
	return &AutoBuyFactory{
		config:        config,
		logger:        logger,
		db:            db,
		marketFactory: marketFactory,
		tradeFactory:  tradeFactory,
	}
}

// CreateAutoBuyRuleRepository creates a repository for auto-buy rules
func (f *AutoBuyFactory) CreateAutoBuyRuleRepository() port.AutoBuyRuleRepository {
	return gorm.NewAutoBuyRuleRepository(f.db, f.logger.With().Str("component", "auto_buy_rule_repository").Logger())
}

// CreateAutoBuyExecutionRepository creates a repository for auto-buy execution records
func (f *AutoBuyFactory) CreateAutoBuyExecutionRepository() port.AutoBuyExecutionRepository {
	return gorm.NewAutoBuyExecutionRepository(f.db, f.logger.With().Str("component", "auto_buy_execution_repository").Logger())
}

// CreateAutoBuyUseCase creates the auto-buy use case
func (f *AutoBuyFactory) CreateAutoBuyUseCase() usecase.AutoBuyUseCase {
	ruleRepo := f.CreateAutoBuyRuleRepository()
	executionRepo := f.CreateAutoBuyExecutionRepository()
	marketDataService := f.marketFactory.CreateMarketDataUseCase()
	symbolRepo := f.marketFactory.CreateSymbolRepository()
	walletRepo := f.tradeFactory.CreateWalletRepository()
	tradeService := f.tradeFactory.CreateTradeService()
	riskService := f.CreateRiskService()

	return usecase.NewAutoBuyUseCase(
		ruleRepo,
		executionRepo,
		marketDataService,
		symbolRepo,
		walletRepo,
		tradeService,
		riskService,
		f.logger.With().Str("component", "auto_buy_usecase").Logger(),
	)
}

// CreateRiskService creates a mock risk service
func (f *AutoBuyFactory) CreateRiskService() port.RiskService {
	// TODO: implement actual risk service when needed
	return nil
}

// CreateAutoBuyHandler creates the HTTP handler for auto-buy functionality
func (f *AutoBuyFactory) CreateAutoBuyHandler() *handler.AutoBuyHandler {
	autoBuyUseCase := f.CreateAutoBuyUseCase()
	return handler.NewAutoBuyHandler(autoBuyUseCase, f.logger)
}
