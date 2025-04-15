package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model/market"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/usecase"
	"github.com/rs/zerolog"
)

// MockNewCoinRepository is a mock type for the NewCoinRepository type
type MockNewCoinRepository struct {
	mock.Mock
}

func (m *MockNewCoinRepository) Save(ctx context.Context, coin *model.NewCoin) error {
	args := m.Called(ctx, coin)
	return args.Error(0)
}

func (m *MockNewCoinRepository) GetBySymbol(ctx context.Context, symbol string) (*model.NewCoin, error) {
	args := m.Called(ctx, symbol)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.NewCoin), args.Error(1)
}

func (m *MockNewCoinRepository) Update(ctx context.Context, coin *model.NewCoin) error {
	args := m.Called(ctx, coin)
	return args.Error(0)
}

func (m *MockNewCoinRepository) GetByStatus(ctx context.Context, status model.Status) ([]*model.NewCoin, error) {
	args := m.Called(ctx, status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.NewCoin), args.Error(1)
}

func (m *MockNewCoinRepository) FindRecentlyListed(ctx context.Context, thresholdTime time.Time) ([]*model.NewCoin, error) {
	args := m.Called(ctx, thresholdTime)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.NewCoin), args.Error(1)
}

func (m *MockNewCoinRepository) GetRecent(ctx context.Context, limit int) ([]*model.NewCoin, error) {
	args := m.Called(ctx, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.NewCoin), args.Error(1)
}

// FindByStatus is an alias to GetByStatus to satisfy interface requirements.
func (m *MockNewCoinRepository) FindByStatus(ctx context.Context, status model.Status) ([]*model.NewCoin, error) {
	return m.GetByStatus(ctx, status)
}

// MockEventBus is a mock type for the EventBus type
type MockEventBus struct {
	mock.Mock
}

func (m *MockEventBus) Publish(event *model.NewCoinEvent) {
	m.Called(event)
}

func (m *MockEventBus) Subscribe(listener func(*model.NewCoinEvent)) {
	m.Called(listener)
}

func (m *MockEventBus) Unsubscribe(listener func(*model.NewCoinEvent)) {
	m.Called(listener)
}

// MockEventRepository is a mock type for the EventRepository type
type MockEventRepository struct {
	mock.Mock
}

func (m *MockEventRepository) SaveEvent(ctx context.Context, event *model.NewCoinEvent) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockEventRepository) GetEvents(ctx context.Context, coinID string, limit, offset int) ([]*model.NewCoinEvent, error) {
	args := m.Called(ctx, coinID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.NewCoinEvent), args.Error(1)
}

// MockMarketDataService is a mock type for the MEXCClient type
type MockMarketDataService struct {
	mock.Mock
}

func (m *MockMarketDataService) GetSymbolInfo(ctx context.Context, symbol string) (*model.SymbolInfo, error) {
	args := m.Called(ctx, symbol)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.SymbolInfo), args.Error(1)
}

func (m *MockMarketDataService) CancelOrder(ctx context.Context, symbol string, orderId string) error {
	args := m.Called(ctx, symbol, orderId)
	return args.Error(0)
}

// Additional methods to satisfy the MEXCClient interface
func (m *MockMarketDataService) GetTicker(ctx context.Context, symbol string) (*model.Ticker, error) {
	args := m.Called(ctx, symbol)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Ticker), args.Error(1)
}

func (m *MockMarketDataService) GetCandles(ctx context.Context, symbol string, interval string, limit int) ([]*market.Candle, error) {
	args := m.Called(ctx, symbol, interval, limit)
	return args.Get(0).([]*market.Candle), args.Error(1)
}

// Update GetOrderBook to match the MEXCClient interface
func (m *MockMarketDataService) GetOrderBook(ctx context.Context, symbol string, depth int) (*model.OrderBook, error) {
	args := m.Called(ctx, symbol, depth)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.OrderBook), args.Error(1)
}

func (m *MockMarketDataService) GetSymbols(ctx context.Context) ([]*model.Symbol, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Symbol), args.Error(1)
}

func (m *MockMarketDataService) GetNewCoins(ctx context.Context) ([]*model.NewCoin, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.NewCoin), args.Error(1)
}

func (m *MockMarketDataService) GetAccount(ctx context.Context) (*model.Wallet, error) {
	args := m.Called(ctx)
	return args.Get(0).(*model.Wallet), args.Error(1)
}

// Update PlaceOrder to match MEXCClient interface
func (m *MockMarketDataService) PlaceOrder(ctx context.Context, symbol string, side model.OrderSide, orderType model.OrderType, quantity float64, price float64, timeInForce model.TimeInForce) (*model.Order, error) {
	args := m.Called(ctx, symbol, side, orderType, quantity, price, timeInForce)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Order), args.Error(1)
}

// GetExchangeInfo retrieves information about all symbols on the exchange
func (m *MockMarketDataService) GetExchangeInfo(ctx context.Context) (*model.ExchangeInfo, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ExchangeInfo), args.Error(1)
}

// Additional missing methods to implement MEXCClient interface
func (m *MockMarketDataService) GetTradingSchedule(ctx context.Context, symbol string) (model.TradingSchedule, error) {
	args := m.Called(ctx, symbol)
	var schedule model.TradingSchedule
	if args.Get(0) != nil {
		schedule = args.Get(0).(model.TradingSchedule)
	}
	return schedule, args.Error(1)
}

func (m *MockMarketDataService) GetSymbolConstraints(ctx context.Context, symbol string) (*model.SymbolConstraints, error) {
	args := m.Called(ctx, symbol)
	var constraints *model.SymbolConstraints
	if args.Get(0) != nil {
		constraints = args.Get(0).(*model.SymbolConstraints)
	}
	return constraints, args.Error(1)
}

func (m *MockMarketDataService) GetKlines(ctx context.Context, symbol string, interval model.KlineInterval, limit int) ([]*model.Kline, error) {
	args := m.Called(ctx, symbol, interval, limit)
	return args.Get(0).([]*model.Kline), args.Error(1)
}

func (m *MockMarketDataService) GetOrderStatus(ctx context.Context, symbol string, orderID string) (*model.Order, error) {
	args := m.Called(ctx, symbol, orderID)
	return args.Get(0).(*model.Order), args.Error(1)
}

func (m *MockMarketDataService) GetOpenOrders(ctx context.Context, symbol string) ([]*model.Order, error) {
	args := m.Called(ctx, symbol)
	return args.Get(0).([]*model.Order), args.Error(1)
}

func (m *MockMarketDataService) GetOrderHistory(ctx context.Context, symbol string, limit, offset int) ([]*model.Order, error) {
	args := m.Called(ctx, symbol, limit, offset)
	return args.Get(0).([]*model.Order), args.Error(1)
}

// Updated MockMarketDataService with proper GetMarketData method
func (m *MockMarketDataService) GetMarketData(ctx context.Context, symbol string) (*model.Ticker, error) {
	args := m.Called(ctx, symbol)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Ticker), args.Error(1)
}

// GetNewListings retrieves information about newly listed coins
func (m *MockMarketDataService) GetNewListings(ctx context.Context) ([]*model.NewCoin, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*model.NewCoin), args.Error(1)
}

// GetSymbolStatus checks if a symbol is currently tradeable
func (m *MockMarketDataService) GetSymbolStatus(ctx context.Context, symbol string) (model.Status, error) {
	args := m.Called(ctx, symbol)
	return args.Get(0).(model.Status), args.Error(1)
}

// --- Test Cases ---

func TestNewCoinUsecase_DetectNewCoins_CoinBecomesTradable(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockNewCoinRepository)
	mockEventRepo := new(MockEventRepository)
	mockBus := new(MockEventBus)
	mockMarketSvc := new(MockMarketDataService)
	logger := zerolog.New(zerolog.NewTestWriter(t))

	// --- Arrange ---
	symbol := "NEWCOIN/USDT"
	expectedListingTime := time.Now().Add(-1 * time.Hour)
	now := time.Now()

	existingCoin := &model.NewCoin{
		ID:                  "uuid-1",
		Symbol:              symbol,
		ExpectedListingTime: expectedListingTime,
		Status:              model.StatusExpected,
		CreatedAt:           now.Add(-2 * time.Hour),
		UpdatedAt:           now.Add(-2 * time.Hour),
	}
	// Expectation for GetNewListings as called by DetectNewCoins
	mockMarketSvc.On("GetNewListings", mock.Anything).Return([]*model.NewCoin{existingCoin}, nil).Once()

	mockRepo.On("GetBySymbol", mock.Anything, symbol).Return(existingCoin, nil).Once()

	symbolInfo := &model.SymbolInfo{
		Symbol: symbol,
		Status: "TRADING",
	}
	mockMarketSvc.On("GetSymbolInfo", ctx, symbol).Return(symbolInfo, nil).Maybe()

	mockRepo.On("Update", ctx, mock.MatchedBy(func(coin *model.NewCoin) bool {
		return coin.Symbol == symbol &&
			coin.Status == model.StatusTrading &&
			coin.BecameTradableAt != nil &&
			!coin.IsProcessedForAutobuy
	})).Return(nil).Maybe()

	mockBus.On("Publish", mock.Anything).Return().Maybe()

	// --- Act ---
	uc := usecase.NewNewCoinUseCase(mockRepo, mockEventRepo, mockBus, mockMarketSvc, &logger)
	err := uc.DetectNewCoins()

	// --- Assert ---
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockMarketSvc.AssertExpectations(t)
	mockBus.AssertExpectations(t)
}
