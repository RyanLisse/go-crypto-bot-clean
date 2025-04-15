package factory

import (
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/delivery/http/handler"
	gormrepo "github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/persistence/gorm"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/config"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/usecase"
	"github.com/rs/zerolog"
	gormdb "gorm.io/gorm"
)

// AutoBuyFactory creates and provides auto-buy related components
type AutoBuyFactory struct {
	config        *config.Config
	logger        *zerolog.Logger
	db            *gormdb.DB
	marketFactory *MarketFactory
	tradeFactory  *TradeFactory
}

// NewAutoBuyFactory creates a new AutoBuyFactory
func NewAutoBuyFactory(
	config *config.Config,
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
	return gormrepo.NewAutoBuyRuleRepository(f.db, f.logger.With().Str("repository", "auto_buy_rule").Logger())
}

// CreateAutoBuyExecutionRepository creates a repository for auto-buy execution records
func (f *AutoBuyFactory) CreateAutoBuyExecutionRepository() port.AutoBuyExecutionRepository {
	return gormrepo.NewAutoBuyExecutionRepository(f.db, f.logger.With().Str("repository", "auto_buy_execution").Logger())
}

// CreateAutoBuyUseCase creates the auto-buy use case
func (f *AutoBuyFactory) CreateAutoBuyUseCase() usecase.AutoBuyUseCase {
	ruleRepo := f.CreateAutoBuyRuleRepository()
	executionRepo := f.CreateAutoBuyExecutionRepository()
	marketDataService, _ := f.marketFactory.CreateMarketDataUseCase()
	symbolRepo := f.CreateSymbolRepository()
	walletRepo := f.CreateWalletRepository()
	tradeService := f.CreateTradeService()
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

// CreateSymbolRepository creates a symbol repository
func (f *AutoBuyFactory) CreateSymbolRepository() port.SymbolRepository {
	// TODO: implement actual repository when needed
	return nil
}

// CreateWalletRepository creates a wallet repository
func (f *AutoBuyFactory) CreateWalletRepository() port.WalletRepository {
	// TODO: implement actual repository when needed
	return nil
}

// CreateTradeService creates a trade service
func (f *AutoBuyFactory) CreateTradeService() port.TradeService {
	// TODO: implement actual trade service when needed
	return nil
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
