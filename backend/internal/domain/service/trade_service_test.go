package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/neo/crypto-bot/internal/domain/model"
	"github.com/neo/crypto-bot/internal/domain/model/market"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockMexcClient is a mock implementation of the MEXC API client
type MockMexcClient struct {
	mock.Mock
}

func (m *MockMexcClient) GetAccount(ctx context.Context) (*model.Wallet, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Wallet), args.Error(1)
}

func (m *MockMexcClient) PlaceOrder(ctx context.Context, symbol string, side model.OrderSide, orderType model.OrderType, quantity float64, price float64, timeInForce model.TimeInForce) (*model.Order, error) {
	args := m.Called(ctx, symbol, side, orderType, quantity, price, timeInForce)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Order), args.Error(1)
}

func (m *MockMexcClient) CancelOrder(ctx context.Context, symbol, orderID string) error {
	args := m.Called(ctx, symbol, orderID)
	return args.Error(0)
}

func (m *MockMexcClient) GetOrderStatus(ctx context.Context, symbol, orderID string) (*model.Order, error) {
	args := m.Called(ctx, symbol, orderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Order), args.Error(1)
}

func (m *MockMexcClient) GetTicker(ctx context.Context, symbol string) (*market.Ticker, error) {
	args := m.Called(ctx, symbol)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*market.Ticker), args.Error(1)
}

func (m *MockMexcClient) GetMarketData(ctx context.Context, symbol string) (*model.Ticker, error) {
	args := m.Called(ctx, symbol)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Ticker), args.Error(1)
}

func (m *MockMexcClient) GetKlines(ctx context.Context, symbol string, interval model.KlineInterval, limit int) ([]*model.Kline, error) {
	args := m.Called(ctx, symbol, interval, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Kline), args.Error(1)
}

func (m *MockMexcClient) GetOrderBook(ctx context.Context, symbol string, limit int) (*model.OrderBook, error) {
	args := m.Called(ctx, symbol, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.OrderBook), args.Error(1)
}

// MockOrderRepository is a mock implementation of the OrderRepository interface
type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) Create(ctx context.Context, order *model.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func (m *MockOrderRepository) Update(ctx context.Context, order *model.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func (m *MockOrderRepository) GetByID(ctx context.Context, id string) (*model.Order, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Order), args.Error(1)
}

func (m *MockOrderRepository) GetByExternalID(ctx context.Context, externalID string) (*model.Order, error) {
	args := m.Called(ctx, externalID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Order), args.Error(1)
}

func (m *MockOrderRepository) GetAll(ctx context.Context) ([]*model.Order, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Order), args.Error(1)
}

func (m *MockOrderRepository) GetAllByStatus(ctx context.Context, status model.OrderStatus) ([]*model.Order, error) {
	args := m.Called(ctx, status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Order), args.Error(1)
}

func (m *MockOrderRepository) GetBySymbol(ctx context.Context, symbol string, limit, offset int) ([]*model.Order, error) {
	args := m.Called(ctx, symbol, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Order), args.Error(1)
}

func (m *MockOrderRepository) GetByClientOrderID(ctx context.Context, clientOrderID string) (*model.Order, error) {
	args := m.Called(ctx, clientOrderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Order), args.Error(1)
}

func (m *MockOrderRepository) GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*model.Order, error) {
	args := m.Called(ctx, userID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Order), args.Error(1)
}

func (m *MockOrderRepository) GetByStatus(ctx context.Context, status model.OrderStatus, limit, offset int) ([]*model.Order, error) {
	args := m.Called(ctx, status, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Order), args.Error(1)
}

func (m *MockOrderRepository) Count(ctx context.Context, filters map[string]interface{}) (int64, error) {
	args := m.Called(ctx, filters)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockOrderRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockSymbolRepository is a mock implementation of the SymbolRepository interface
type MockSymbolRepository struct {
	mock.Mock
}

func (m *MockSymbolRepository) Create(ctx context.Context, symbol *market.Symbol) error {
	args := m.Called(ctx, symbol)
	return args.Error(0)
}

func (m *MockSymbolRepository) Update(ctx context.Context, symbol *market.Symbol) error {
	args := m.Called(ctx, symbol)
	return args.Error(0)
}

func (m *MockSymbolRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockSymbolRepository) GetByID(ctx context.Context, id string) (*market.Symbol, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*market.Symbol), args.Error(1)
}

func (m *MockSymbolRepository) GetBySymbol(ctx context.Context, symbol string) (*market.Symbol, error) {
	args := m.Called(ctx, symbol)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*market.Symbol), args.Error(1)
}

func (m *MockSymbolRepository) GetAll(ctx context.Context) ([]*market.Symbol, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*market.Symbol), args.Error(1)
}

func (m *MockSymbolRepository) GetByExchange(ctx context.Context, exchange string) ([]*market.Symbol, error) {
	args := m.Called(ctx, exchange)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*market.Symbol), args.Error(1)
}

// MockTradeMarketDataService is a mock specific for trade service tests
type MockTradeMarketDataService struct {
	mock.Mock
}

func (m *MockTradeMarketDataService) RefreshTicker(ctx context.Context, symbol string) (*market.Ticker, error) {
	args := m.Called(ctx, symbol)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*market.Ticker), args.Error(1)
}

// MockMarketCache is a mock implementation of the market cache interface
type MockMarketCache struct {
	mock.Mock
}

func (m *MockMarketCache) CacheTicker(ticker *market.Ticker) {
	m.Called(ticker)
}

func (m *MockMarketCache) GetTicker(ctx context.Context, exchange, symbol string) (*market.Ticker, bool) {
	args := m.Called(ctx, exchange, symbol)
	if args.Get(0) == nil {
		return nil, args.Bool(1)
	}
	return args.Get(0).(*market.Ticker), args.Bool(1)
}

func (m *MockMarketCache) GetAllTickers(ctx context.Context, exchange string) ([]*market.Ticker, bool) {
	args := m.Called(ctx, exchange)
	if args.Get(0) == nil {
		return nil, args.Bool(1)
	}
	return args.Get(0).([]*market.Ticker), args.Bool(1)
}

func (m *MockMarketCache) GetLatestTickers(ctx context.Context) ([]*market.Ticker, bool) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Bool(1)
	}
	return args.Get(0).([]*market.Ticker), args.Bool(1)
}

func (m *MockMarketCache) CacheCandle(candle *market.Candle) {
	m.Called(candle)
}

func (m *MockMarketCache) GetCandle(ctx context.Context, exchange, symbol string, interval market.Interval, openTime time.Time) (*market.Candle, bool) {
	args := m.Called(ctx, exchange, symbol, interval, openTime)
	if args.Get(0) == nil {
		return nil, args.Bool(1)
	}
	return args.Get(0).(*market.Candle), args.Bool(1)
}

func (m *MockMarketCache) GetLatestCandle(ctx context.Context, exchange, symbol string, interval market.Interval) (*market.Candle, bool) {
	args := m.Called(ctx, exchange, symbol, interval)
	if args.Get(0) == nil {
		return nil, args.Bool(1)
	}
	return args.Get(0).(*market.Candle), args.Bool(1)
}

func (m *MockMarketCache) CacheOrderBook(orderbook *market.OrderBook) {
	m.Called(orderbook)
}

func (m *MockMarketCache) GetOrderBook(ctx context.Context, exchange, symbol string) (*market.OrderBook, bool) {
	args := m.Called(ctx, exchange, symbol)
	if args.Get(0) == nil {
		return nil, args.Bool(1)
	}
	return args.Get(0).(*market.OrderBook), args.Bool(1)
}

func (m *MockMarketCache) Clear() {
	m.Called()
}

func (m *MockMarketCache) SetTickerExpiry(d time.Duration) {
	m.Called(d)
}

func (m *MockMarketCache) SetCandleExpiry(d time.Duration) {
	m.Called(d)
}

func (m *MockMarketCache) SetOrderbookExpiry(d time.Duration) {
	m.Called(d)
}

func (m *MockMarketCache) StartCleanupTask(ctx context.Context, interval time.Duration) {
	m.Called(ctx, interval)
}

// TestPlaceOrder tests the PlaceOrder method of MexcTradeService
func TestPlaceOrder(t *testing.T) {
	// Create mocks
	mockClient := new(MockMexcClient)
	mockOrderRepo := new(MockOrderRepository)
	mockSymbolRepo := new(MockSymbolRepository)
	mockMarketData := new(MockTradeMarketDataService)
	logger := zerolog.New(zerolog.NewConsoleWriter()).With().Timestamp().Logger()

	// Create the service
	marketService := &MarketDataService{
		marketRepo: &mockMarketRepoWrapper{mockMarketData},
		symbolRepo: mockSymbolRepo,
		logger:     &logger,
	}

	service := NewMexcTradeService(
		mockClient,
		marketService,
		mockSymbolRepo,
		mockOrderRepo,
		&logger,
	)

	// Setup test data
	ctx := context.Background()
	symbol := "BTC-USDT"
	side := model.OrderSideBuy
	orderType := model.OrderTypeLimit
	quantity := 0.001
	price := 50000.0
	timeInForce := model.TimeInForceGTC

	symbolInfo := &market.Symbol{
		Symbol:     symbol,
		BaseAsset:  "BTC",
		QuoteAsset: "USDT",
		MinQty:     0.0001,
	}

	orderRequest := &model.OrderRequest{
		Symbol:   symbol,
		Side:     side,
		Type:     orderType,
		Quantity: quantity,
		Price:    price,
	}

	order := &model.Order{
		Symbol:      symbol,
		Side:        side,
		Type:        orderType,
		Quantity:    quantity,
		Price:       price,
		TimeInForce: timeInForce,
		Status:      model.OrderStatusNew,
		OrderID:     "123456",
	}

	// Setup expectations
	mockSymbolRepo.On("GetBySymbol", ctx, symbol).Return(symbolInfo, nil)
	mockClient.On("PlaceOrder", ctx, symbol, side, orderType, quantity, price, timeInForce).Return(order, nil)
	mockOrderRepo.On("Create", ctx, mock.Anything).Return(nil)

	// Call the service method
	result, err := service.PlaceOrder(ctx, orderRequest)

	// Assertions
	require.NoError(t, err)
	assert.Equal(t, order.OrderID, result.Order.OrderID)
	assert.Equal(t, model.OrderStatusNew, result.Order.Status)

	// Verify mocks
	mockClient.AssertExpectations(t)
	mockOrderRepo.AssertExpectations(t)
	mockSymbolRepo.AssertExpectations(t)
}

// TestPlaceOrderWithError tests error handling in PlaceOrder
func TestPlaceOrderWithError(t *testing.T) {
	// Create mocks
	mockClient := new(MockMexcClient)
	mockOrderRepo := new(MockOrderRepository)
	mockSymbolRepo := new(MockSymbolRepository)
	mockMarketData := new(MockTradeMarketDataService)
	logger := zerolog.New(zerolog.NewConsoleWriter()).With().Timestamp().Logger()

	// Create the service
	marketService := &MarketDataService{
		marketRepo: &mockMarketRepoWrapper{mockMarketData},
		symbolRepo: mockSymbolRepo,
		logger:     &logger,
	}

	service := NewMexcTradeService(
		mockClient,
		marketService,
		mockSymbolRepo,
		mockOrderRepo,
		&logger,
	)

	// Setup test data
	ctx := context.Background()
	symbol := "BTC-USDT"
	side := model.OrderSideBuy
	orderType := model.OrderTypeLimit
	quantity := 0.001
	price := 50000.0

	orderRequest := &model.OrderRequest{
		Symbol:   symbol,
		Side:     side,
		Type:     orderType,
		Quantity: quantity,
		Price:    price,
	}

	symbolInfo := &market.Symbol{
		Symbol:     symbol,
		BaseAsset:  "BTC",
		QuoteAsset: "USDT",
		MinQty:     0.0001,
	}

	expectedErr := errors.New("API error")

	// Setup expectations
	mockSymbolRepo.On("GetBySymbol", ctx, symbol).Return(symbolInfo, nil)
	mockClient.On("PlaceOrder", ctx, symbol, side, orderType, quantity, price, model.TimeInForceGTC).Return(nil, expectedErr)

	// Call the service method
	result, err := service.PlaceOrder(ctx, orderRequest)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, result)

	// Verify mocks
	mockClient.AssertExpectations(t)
	mockSymbolRepo.AssertExpectations(t)
}

// TestCancelOrder tests the CancelOrder method
func TestCancelOrder(t *testing.T) {
	// Create mocks
	mockClient := new(MockMexcClient)
	mockOrderRepo := new(MockOrderRepository)
	mockSymbolRepo := new(MockSymbolRepository)
	mockMarketData := new(MockTradeMarketDataService)
	logger := zerolog.New(zerolog.NewConsoleWriter()).With().Timestamp().Logger()

	// Create the service
	marketService := &MarketDataService{
		marketRepo: &mockMarketRepoWrapper{mockMarketData},
		symbolRepo: mockSymbolRepo,
		logger:     &logger,
	}

	service := NewMexcTradeService(
		mockClient,
		marketService,
		mockSymbolRepo,
		mockOrderRepo,
		&logger,
	)

	// Setup test data
	ctx := context.Background()
	symbol := "BTC-USDT"
	orderID := "order123"
	order := &model.Order{
		ID:      "internal123",
		Symbol:  symbol,
		OrderID: orderID,
		Status:  model.OrderStatusNew,
	}

	// Setup expectations
	mockOrderRepo.On("GetByID", ctx, orderID).Return(order, nil)
	mockClient.On("CancelOrder", ctx, symbol, orderID).Return(nil)
	mockOrderRepo.On("Update", ctx, mock.Anything).Return(nil)

	// Call the service method
	err := service.CancelOrder(ctx, symbol, orderID)

	// Assertions
	require.NoError(t, err)

	// Verify mocks
	mockClient.AssertExpectations(t)
	mockOrderRepo.AssertExpectations(t)
}

// TestGetOrderStatus tests the GetOrderStatus method
func TestGetOrderStatus(t *testing.T) {
	// Create mocks
	mockClient := new(MockMexcClient)
	mockOrderRepo := new(MockOrderRepository)
	mockSymbolRepo := new(MockSymbolRepository)
	mockMarketData := new(MockTradeMarketDataService)
	logger := zerolog.New(zerolog.NewConsoleWriter()).With().Timestamp().Logger()

	// Create the service
	marketService := &MarketDataService{
		marketRepo: &mockMarketRepoWrapper{mockMarketData},
		symbolRepo: mockSymbolRepo,
		logger:     &logger,
	}

	service := NewMexcTradeService(
		mockClient,
		marketService,
		mockSymbolRepo,
		mockOrderRepo,
		&logger,
	)

	// Setup test data
	ctx := context.Background()
	symbol := "BTC-USDT"
	orderID := "order123"
	order := &model.Order{
		ID:      "internal123",
		Symbol:  symbol,
		OrderID: orderID,
		Status:  model.OrderStatusNew,
	}

	updatedOrder := &model.Order{
		ID:      "internal123",
		Symbol:  symbol,
		OrderID: orderID,
		Status:  model.OrderStatusFilled,
	}

	// Setup expectations
	mockOrderRepo.On("GetByID", ctx, orderID).Return(order, nil)
	mockClient.On("GetOrderStatus", ctx, symbol, orderID).Return(updatedOrder, nil)
	mockOrderRepo.On("Update", ctx, mock.Anything).Return(nil)

	// Call the service method
	result, err := service.GetOrderStatus(ctx, symbol, orderID)

	// Assertions
	require.NoError(t, err)
	assert.Equal(t, model.OrderStatusFilled, result.Status)

	// Verify mocks
	mockClient.AssertExpectations(t)
	mockOrderRepo.AssertExpectations(t)
}

// TestCalculateRequiredQuantity tests the CalculateRequiredQuantity method
func TestCalculateRequiredQuantity(t *testing.T) {
	// Create mocks
	mockClient := new(MockMexcClient)
	mockOrderRepo := new(MockOrderRepository)
	mockSymbolRepo := new(MockSymbolRepository)
	mockMarketData := new(MockTradeMarketDataService)
	mockCache := new(MockMarketCache)
	logger := zerolog.New(zerolog.NewConsoleWriter()).With().Timestamp().Logger()

	// Create the market service with proper logger and cache
	marketService := &MarketDataService{
		marketRepo: &mockMarketRepoWrapper{mockMarketData},
		symbolRepo: mockSymbolRepo,
		logger:     &logger,
		mexcAPI:    mockClient,
		cache:      mockCache,
	}

	// Setup the trade service
	service := &MexcTradeService{
		mexcAPI:       mockClient,
		marketService: marketService,
		symbolRepo:    mockSymbolRepo,
		orderRepo:     mockOrderRepo,
		logger:        &logger,
	}

	// Setup test data
	ctx := context.Background()
	symbol := "BTC-USDT"
	amount := 10000.0 // $10,000

	ticker := &market.Ticker{
		Exchange: "MEXC",
		Symbol:   symbol,
		Price:    50000.0, // BTC at $50,000
	}

	symbolInfo := &market.Symbol{
		Symbol: symbol,
		Status: "TRADING",
		MinQty: 0.0001,
	}

	modelTicker := &model.Ticker{
		Symbol:    symbol,
		LastPrice: 50000.0,
	}

	// Setup expectations
	mockCache.On("GetTicker", ctx, "mexc", symbol).Return(nil, false)
	mockClient.On("GetMarketData", ctx, symbol).Return(modelTicker, nil)
	mockCache.On("CacheTicker", mock.AnythingOfType("*market.Ticker")).Return()
	mockMarketData.On("RefreshTicker", ctx, symbol).Return(ticker, nil)
	mockSymbolRepo.On("GetBySymbol", ctx, symbol).Return(symbolInfo, nil)

	// Call the service method
	quantity, err := service.CalculateRequiredQuantity(ctx, symbol, model.OrderSideBuy, amount)

	// Assertions
	require.NoError(t, err)
	assert.Equal(t, 0.2, quantity) // $10,000 / $50,000 = 0.2 BTC

	// Verify mocks
	mockClient.AssertExpectations(t)
	mockSymbolRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

// mockMarketRepoWrapper is a helper to forward calls to our mock
type mockMarketRepoWrapper struct {
	mock *MockTradeMarketDataService
}

func (w *mockMarketRepoWrapper) SaveTicker(ctx context.Context, ticker *market.Ticker) error {
	return nil
}

func (w *mockMarketRepoWrapper) GetTicker(ctx context.Context, symbol, exchange string) (*market.Ticker, error) {
	return w.mock.RefreshTicker(ctx, symbol)
}

func (w *mockMarketRepoWrapper) GetAllTickers(ctx context.Context, exchange string) ([]*market.Ticker, error) {
	return nil, nil
}

func (w *mockMarketRepoWrapper) GetTickerHistory(ctx context.Context, symbol, exchange string, start, end time.Time) ([]*market.Ticker, error) {
	return nil, nil
}

func (w *mockMarketRepoWrapper) SaveCandle(ctx context.Context, candle *market.Candle) error {
	return nil
}

func (w *mockMarketRepoWrapper) SaveCandles(ctx context.Context, candles []*market.Candle) error {
	return nil
}

func (w *mockMarketRepoWrapper) GetCandle(ctx context.Context, symbol, exchange string, interval market.Interval, openTime time.Time) (*market.Candle, error) {
	return nil, nil
}

func (w *mockMarketRepoWrapper) GetCandles(ctx context.Context, symbol, exchange string, interval market.Interval, start, end time.Time, limit int) ([]*market.Candle, error) {
	return nil, nil
}

func (w *mockMarketRepoWrapper) GetLatestCandle(ctx context.Context, symbol, exchange string, interval market.Interval) (*market.Candle, error) {
	return nil, nil
}

func (w *mockMarketRepoWrapper) PurgeOldData(ctx context.Context, olderThan time.Time) error {
	return nil
}

func (w *mockMarketRepoWrapper) GetLatestTickers(ctx context.Context, limit int) ([]*market.Ticker, error) {
	return nil, nil
}

func (w *mockMarketRepoWrapper) GetTickersBySymbol(ctx context.Context, symbol string, limit int) ([]*market.Ticker, error) {
	return nil, nil
}

func (w *mockMarketRepoWrapper) RefreshTicker(ctx context.Context, symbol string) (*market.Ticker, error) {
	return w.mock.RefreshTicker(ctx, symbol)
}

// TestGetOpenOrders tests the GetOpenOrders method
func TestGetOpenOrders(t *testing.T) {
	// Create mocks
	mockClient := new(MockMexcClient)
	mockOrderRepo := new(MockOrderRepository)
	mockSymbolRepo := new(MockSymbolRepository)
	mockMarketData := new(MockTradeMarketDataService)
	logger := zerolog.New(zerolog.NewConsoleWriter()).With().Timestamp().Logger()

	// Create the service
	marketService := &MarketDataService{
		marketRepo: &mockMarketRepoWrapper{mockMarketData},
		symbolRepo: mockSymbolRepo,
		logger:     &logger,
	}

	service := NewMexcTradeService(
		mockClient,
		marketService,
		mockSymbolRepo,
		mockOrderRepo,
		&logger,
	)

	// Setup test data
	ctx := context.Background()
	symbol := "BTC-USDT"
	limit := 100
	offset := 0

	// Create a slice of orders with different statuses
	allOrders := []*model.Order{
		{
			OrderID:   "order1",
			Symbol:    symbol,
			Status:    model.OrderStatusNew,
			CreatedAt: time.Now(),
		},
		{
			OrderID:   "order2",
			Symbol:    symbol,
			Status:    model.OrderStatusPartiallyFilled,
			CreatedAt: time.Now(),
		},
		{
			OrderID:   "order3",
			Symbol:    symbol,
			Status:    model.OrderStatusFilled,
			CreatedAt: time.Now(),
		},
		{
			OrderID:   "order4",
			Symbol:    symbol,
			Status:    model.OrderStatusCanceled,
			CreatedAt: time.Now(),
		},
	}

	// Setup expectations
	mockOrderRepo.On("GetBySymbol", ctx, symbol, limit, offset).Return(allOrders, nil)

	// Call the method
	orders, err := service.GetOpenOrders(ctx, symbol)

	// Assert results
	require.NoError(t, err)
	require.NotNil(t, orders)
	require.Equal(t, 2, len(orders), "Should return only open orders (NEW or PARTIALLY_FILLED)")

	// Verify that returned orders have the expected statuses
	for _, order := range orders {
		assert.True(t, order.Status == model.OrderStatusNew || order.Status == model.OrderStatusPartiallyFilled)
	}

	// Verify expectations were met
	mockOrderRepo.AssertExpectations(t)
}

// TestGetOpenOrdersError tests error handling in GetOpenOrders
func TestGetOpenOrdersError(t *testing.T) {
	// Create mocks
	mockClient := new(MockMexcClient)
	mockOrderRepo := new(MockOrderRepository)
	mockSymbolRepo := new(MockSymbolRepository)
	mockMarketData := new(MockTradeMarketDataService)
	logger := zerolog.New(zerolog.NewConsoleWriter()).With().Timestamp().Logger()

	// Create the service
	marketService := &MarketDataService{
		marketRepo: &mockMarketRepoWrapper{mockMarketData},
		symbolRepo: mockSymbolRepo,
		logger:     &logger,
	}

	service := NewMexcTradeService(
		mockClient,
		marketService,
		mockSymbolRepo,
		mockOrderRepo,
		&logger,
	)

	// Setup test data
	ctx := context.Background()
	symbol := "BTC-USDT"
	limit := 100
	offset := 0
	expectedError := errors.New("database error")

	// Setup expectations
	mockOrderRepo.On("GetBySymbol", ctx, symbol, limit, offset).Return(nil, expectedError)

	// Call the method
	orders, err := service.GetOpenOrders(ctx, symbol)

	// Assert results
	require.Error(t, err)
	require.Nil(t, orders)
	assert.Contains(t, err.Error(), "failed to get open orders")

	// Verify expectations were met
	mockOrderRepo.AssertExpectations(t)
}

// TestGetOrderHistory tests the GetOrderHistory method
func TestGetOrderHistory(t *testing.T) {
	// Create mocks
	mockClient := new(MockMexcClient)
	mockOrderRepo := new(MockOrderRepository)
	mockSymbolRepo := new(MockSymbolRepository)
	mockMarketData := new(MockTradeMarketDataService)
	logger := zerolog.New(zerolog.NewConsoleWriter()).With().Timestamp().Logger()

	// Create the service
	marketService := &MarketDataService{
		marketRepo: &mockMarketRepoWrapper{mockMarketData},
		symbolRepo: mockSymbolRepo,
		logger:     &logger,
	}

	service := NewMexcTradeService(
		mockClient,
		marketService,
		mockSymbolRepo,
		mockOrderRepo,
		&logger,
	)

	// Setup test data
	ctx := context.Background()
	symbol := "BTC-USDT"
	limit := 50
	offset := 10

	// Create a slice of historical orders
	historicalOrders := []*model.Order{
		{
			OrderID:   "order1",
			Symbol:    symbol,
			Status:    model.OrderStatusFilled,
			CreatedAt: time.Now().Add(-24 * time.Hour),
		},
		{
			OrderID:   "order2",
			Symbol:    symbol,
			Status:    model.OrderStatusCanceled,
			CreatedAt: time.Now().Add(-48 * time.Hour),
		},
	}

	// Setup expectations
	mockOrderRepo.On("GetBySymbol", ctx, symbol, limit, offset).Return(historicalOrders, nil)

	// Call the method
	orders, err := service.GetOrderHistory(ctx, symbol, limit, offset)

	// Assert results
	require.NoError(t, err)
	require.NotNil(t, orders)
	require.Equal(t, 2, len(orders))
	assert.Equal(t, historicalOrders, orders)

	// Verify expectations were met
	mockOrderRepo.AssertExpectations(t)
}

// TestGetOrderHistoryError tests error handling in GetOrderHistory
func TestGetOrderHistoryError(t *testing.T) {
	// Create mocks
	mockClient := new(MockMexcClient)
	mockOrderRepo := new(MockOrderRepository)
	mockSymbolRepo := new(MockSymbolRepository)
	mockMarketData := new(MockTradeMarketDataService)
	logger := zerolog.New(zerolog.NewConsoleWriter()).With().Timestamp().Logger()

	// Create the service
	marketService := &MarketDataService{
		marketRepo: &mockMarketRepoWrapper{mockMarketData},
		symbolRepo: mockSymbolRepo,
		logger:     &logger,
	}

	service := NewMexcTradeService(
		mockClient,
		marketService,
		mockSymbolRepo,
		mockOrderRepo,
		&logger,
	)

	// Setup test data
	ctx := context.Background()
	symbol := "BTC-USDT"
	limit := 50
	offset := 10
	expectedError := errors.New("database error")

	// Setup expectations
	mockOrderRepo.On("GetBySymbol", ctx, symbol, limit, offset).Return(nil, expectedError)

	// Call the method
	orders, err := service.GetOrderHistory(ctx, symbol, limit, offset)

	// Assert results
	require.Error(t, err)
	require.Nil(t, orders)
	assert.Contains(t, err.Error(), "failed to get order history")

	// Verify expectations were met
	mockOrderRepo.AssertExpectations(t)
}
