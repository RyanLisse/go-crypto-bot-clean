package factory

import (
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/delivery/http/handler"
	persistence "github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/persistence/gorm"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/config"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/service"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/usecase"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// TradeFactory creates trade execution related components
type TradeFactory struct {
	config *config.Config
	logger *zerolog.Logger
	db     *gorm.DB
}

// NewTradeFactory creates a new TradeFactory
func NewTradeFactory(cfg *config.Config, logger *zerolog.Logger, db *gorm.DB) *TradeFactory {
	return &TradeFactory{
		config: cfg,
		logger: logger,
		db:     db,
	}
}

// CreateTradeService creates a new implementation of the TradeService
func (f *TradeFactory) CreateTradeService(
	mexcClient port.MEXCClient,
	marketDataService *service.MarketDataService,
	symbolRepo port.SymbolRepository,
	orderRepo port.OrderRepository,
) port.TradeService {
	// Create the trade service with necessary dependencies
	return service.NewMexcTradeService(
		mexcClient,
		marketDataService,
		symbolRepo,
		orderRepo,
		f.logger,
	)
}

// CreateTradeUseCase creates a new TradeUseCase implementation
func (f *TradeFactory) CreateTradeUseCase(
	mexcClient port.MEXCClient,
	symbolRepo port.SymbolRepository,
	orderRepo port.OrderRepository,
	tradeService port.TradeService,
	riskUC usecase.RiskUseCase,
	txManager port.TransactionManager,
) usecase.TradeUseCase {
	// Create the trade use case with necessary dependencies
	return usecase.NewTradeUseCase(
		mexcClient,
		orderRepo,
		symbolRepo,
		tradeService,
		riskUC,
		txManager,
		f.logger.With().Str("component", "trade_usecase").Logger(),
	)
}

// CreateTradeHandler creates a new TradeHandler for HTTP API
func (f *TradeFactory) CreateTradeHandler(tradeUseCase usecase.TradeUseCase) *handler.TradeHandler {
	// Create the trade handler with the use case
	return handler.NewTradeHandler(tradeUseCase, f.logger)
}

// CreateOrderRepository creates a repository for order persistence
func (f *TradeFactory) CreateOrderRepository() port.OrderRepository {
	// Use the persistence/gorm implementation
	logger := f.logger.With().Str("component", "order_repository").Logger()
	return persistence.NewOrderRepository(f.db, &logger)
}
