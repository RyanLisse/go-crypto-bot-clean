package factory

import (
	"github.com/neo/crypto-bot/internal/adapter/http/handler"
	"github.com/neo/crypto-bot/internal/config"
	"github.com/neo/crypto-bot/internal/domain/port"
	"github.com/neo/crypto-bot/internal/domain/service"
	"github.com/neo/crypto-bot/internal/usecase"
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
	mexcAPI port.MexcAPI,
	marketDataService *service.MarketDataService,
	symbolRepo port.SymbolRepository,
	orderRepo port.OrderRepository,
) port.TradeService {
	// Create the trade service with necessary dependencies
	return service.NewMexcTradeService(
		mexcAPI,
		marketDataService,
		symbolRepo,
		orderRepo,
		f.logger,
	)
}

// CreateTradeUseCase creates a new TradeUseCase implementation
func (f *TradeFactory) CreateTradeUseCase(
	mexcAPI port.MexcAPI,
	symbolRepo port.SymbolRepository,
	orderRepo port.OrderRepository,
	tradeService port.TradeService,
) usecase.TradeUseCase {
	// Create the trade use case with necessary dependencies
	return usecase.NewTradeUseCase(
		mexcAPI,
		orderRepo,
		symbolRepo,
		tradeService,
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
	return gorm.NewOrderRepository(f.db, f.logger)
}
