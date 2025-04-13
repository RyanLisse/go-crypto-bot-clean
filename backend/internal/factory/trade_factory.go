package factory

import (
	"context"

	"github.com/neo/crypto-bot/internal/adapter/http/handler"
	"github.com/neo/crypto-bot/internal/config"
	"github.com/neo/crypto-bot/internal/domain/model"
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

// DummyOrderRepository is a stub implementation of port.OrderRepository for testing purposes.
type DummyOrderRepository struct{}

func (d *DummyOrderRepository) Create(ctx context.Context, order *model.Order) error { return nil }
func (d *DummyOrderRepository) GetByID(ctx context.Context, id string) (*model.Order, error) {
	return nil, nil
}
func (d *DummyOrderRepository) GetByClientOrderID(ctx context.Context, clientOrderID string) (*model.Order, error) {
	return nil, nil
}
func (d *DummyOrderRepository) Update(ctx context.Context, order *model.Order) error { return nil }
func (d *DummyOrderRepository) GetBySymbol(ctx context.Context, symbol string, limit, offset int) ([]*model.Order, error) {
	return []*model.Order{}, nil
}
func (d *DummyOrderRepository) GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*model.Order, error) {
	return []*model.Order{}, nil
}
func (d *DummyOrderRepository) GetByStatus(ctx context.Context, status model.OrderStatus, limit, offset int) ([]*model.Order, error) {
	return []*model.Order{}, nil
}
func (d *DummyOrderRepository) Count(ctx context.Context, filters map[string]interface{}) (int64, error) {
	return 0, nil
}
func (d *DummyOrderRepository) Delete(ctx context.Context, id string) error { return nil }

// NewOrderRepository is a dummy stub for order repository to allow build. It returns an instance of DummyOrderRepository.
func NewOrderRepository(db interface{}, logger interface{}) port.OrderRepository {
	return &DummyOrderRepository{}
}

// Dummy stub for NewOrderRepository to allow build. This should be replaced with the actual implementation.

func (f *TradeFactory) CreateOrderRepository() port.OrderRepository {
	return NewOrderRepository(f.db, f.logger)
}
