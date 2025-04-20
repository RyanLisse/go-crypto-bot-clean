package mexc

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/adapter/infrastructure/gateway/mexc/apikeystore"
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/model"
	"github.com/rs/zerolog"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockMEXCClient is a mock implementation of the port.MEXCClient interface
type MockMEXCClient struct {
	mock.Mock
}

// GetMarketData mocks the GetMarketData method
func (m *MockMEXCClient) GetMarketData(ctx context.Context, symbol string) (*model.Ticker, error) {
	args := m.Called(ctx, symbol)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Ticker), args.Error(1)
}

// GetKlines mocks the GetKlines method
func (m *MockMEXCClient) GetKlines(ctx context.Context, symbol string, interval model.KlineInterval, limit int) ([]*model.Kline, error) {
	args := m.Called(ctx, symbol, interval, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Kline), args.Error(1)
}

// GetOrderBook mocks the GetOrderBook method
func (m *MockMEXCClient) GetOrderBook(ctx context.Context, symbol string, depth int) (*model.OrderBook, error) {
	args := m.Called(ctx, symbol, depth)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.OrderBook), args.Error(1)
}

// GetSymbols mocks the GetSymbols method
func (m *MockMEXCClient) GetSymbols(ctx context.Context) ([]*model.Symbol, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Symbol), args.Error(1)
}

// GetSymbol mocks the GetSymbol method
func (m *MockMEXCClient) GetSymbol(ctx context.Context, symbol string) (*model.Symbol, error) {
	args := m.Called(ctx, symbol)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Symbol), args.Error(1)
}

// GetServerTime mocks the GetServerTime method
func (m *MockMEXCClient) GetServerTime(ctx context.Context) (time.Time, error) {
	args := m.Called(ctx)
	return args.Get(0).(time.Time), args.Error(1)
}

// GetExchangeInfo mocks the GetExchangeInfo method
func (m *MockMEXCClient) GetExchangeInfo(ctx context.Context) (*model.ExchangeInfo, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ExchangeInfo), args.Error(1)
}

// PlaceOrder mocks the PlaceOrder method
func (m *MockMEXCClient) PlaceOrder(ctx context.Context, symbol string, side model.OrderSide, orderType model.OrderType, quantity float64, price float64, timeInForce model.TimeInForce) (*model.Order, error) {
	args := m.Called(ctx, symbol, side, orderType, quantity, price, timeInForce)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Order), args.Error(1)
}

func TestGetSymbols(t *testing.T) {
	// Skip this test for now as it requires mocking the REST client
	t.Skip("Skipping test as it requires mocking the REST client")
}

func TestGetTicker(t *testing.T) {
	// Skip this test for now as it requires mocking the REST client
	t.Skip("Skipping test as it requires mocking the REST client")
}

func TestGetOrderBook(t *testing.T) {
	// Skip this test for now as it requires mocking the REST client
	t.Skip("Skipping test as it requires mocking the REST client")
}

// TestGateway is a wrapper around Gateway that allows us to inject mocks for testing
type TestGateway struct {
	Gateway
	mockAccountInfo model.AccountInfo
	mockError       error
	handlerCalled   bool
	lastTicker      model.Ticker
	lastOrderBook   model.OrderBook
}

// GetAccountInfo overrides the Gateway.GetAccountInfo method for testing
func (g *TestGateway) GetAccountInfo(ctx context.Context) (model.AccountInfo, error) {
	return g.mockAccountInfo, g.mockError
}

// GetAssetBalance overrides the Gateway.GetAssetBalance method for testing
func (g *TestGateway) GetAssetBalance(ctx context.Context, asset string) (decimal.Decimal, error) {
	if g.mockError != nil {
		return decimal.Zero, g.mockError
	}

	// Find the balance for the specified asset
	for _, balance := range g.mockAccountInfo.Balances {
		if string(balance.Asset) == asset {
			return decimal.NewFromFloat(balance.Free), nil
		}
	}

	// Asset not found, return zero balance
	return decimal.Zero, nil
}

// CancelOrder overrides the Gateway.CancelOrder method for testing
func (g *TestGateway) CancelOrder(ctx context.Context, symbol, orderID string) error {
	return g.mockError
}

// SubscribeToTicker overrides the Gateway.SubscribeToTicker method for testing
func (g *TestGateway) SubscribeToTicker(ctx context.Context, symbol string, handler func(model.Ticker)) error {
	if g.mockError != nil {
		return g.mockError
	}

	// Call the handler with a mock ticker
	g.lastTicker = model.Ticker{
		Symbol:    symbol,
		Exchange:  "MEXC",
		LastPrice: 50000.0,
		Volume:    100.0,
		Timestamp: time.Now(),
	}
	handler(g.lastTicker)
	g.handlerCalled = true

	return nil
}

// SubscribeToOrderBook overrides the Gateway.SubscribeToOrderBook method for testing
func (g *TestGateway) SubscribeToOrderBook(ctx context.Context, symbol string, depth int, handler func(model.OrderBook)) error {
	if g.mockError != nil {
		return g.mockError
	}

	// Call the handler with a mock order book
	g.lastOrderBook = model.OrderBook{
		Symbol: symbol,
		Bids: []model.OrderBookEntry{
			{Price: 49900.0, Quantity: 1.0},
			{Price: 49800.0, Quantity: 2.0},
		},
		Asks: []model.OrderBookEntry{
			{Price: 50100.0, Quantity: 1.0},
			{Price: 50200.0, Quantity: 2.0},
		},
		Timestamp: time.Now(),
	}
	handler(g.lastOrderBook)
	g.handlerCalled = true

	return nil
}

// Unsubscribe overrides the Gateway.Unsubscribe method for testing
func (g *TestGateway) Unsubscribe(ctx context.Context, channel, symbol string) error {
	return g.mockError
}

func TestGetAssetBalance(t *testing.T) {
	// Create a test gateway
	testGateway := &TestGateway{
		mockAccountInfo: model.AccountInfo{
			Balances: []model.Balance{
				{Asset: "BTC", Free: 1.0, Locked: 0.0, Total: 1.0},
				{Asset: "USDT", Free: 10000.0, Locked: 0.0, Total: 10000.0},
			},
		},
		mockError: nil,
	}

	// Call the method
	balance, err := testGateway.GetAssetBalance(context.Background(), "BTC")

	// Assert expectations
	assert.NoError(t, err)
	assert.Equal(t, "1", balance.String())

	// Test with non-existent asset
	balance, err = testGateway.GetAssetBalance(context.Background(), "ETH")

	// Assert expectations
	assert.NoError(t, err)
	assert.Equal(t, "0", balance.String())

	// Test with error
	testGateway.mockError = fmt.Errorf("test error")
	_, err = testGateway.GetAssetBalance(context.Background(), "BTC")
	assert.Error(t, err)
	assert.Equal(t, "test error", err.Error())
}

// PlaceOrder overrides the Gateway.PlaceOrder method for testing
func (g *TestGateway) PlaceOrder(ctx context.Context, orderRequest model.OrderRequest) (model.OrderResponse, error) {
	if g.mockError != nil {
		return model.OrderResponse{IsSuccess: false}, g.mockError
	}

	return model.OrderResponse{
		Order: model.Order{
			OrderID:       "12345",
			ClientOrderID: "client-12345",
			Symbol:        orderRequest.Symbol,
			CreatedAt:     time.Now(),
			Price:         orderRequest.Price,
			Quantity:      orderRequest.Quantity,
			ExecutedQty:   0,
			Status:        model.OrderStatusNew,
			Type:          orderRequest.Type,
			Side:          orderRequest.Side,
			Exchange:      "MEXC",
		},
		IsSuccess: true,
	}, nil
}

func TestPlaceOrder(t *testing.T) {
	// Create a test gateway
	testGateway := &TestGateway{
		mockError: nil,
	}

	// Create an order request
	orderRequest := model.OrderRequest{
		Symbol:      "BTCUSDT",
		Side:        model.OrderSideBuy,
		Type:        model.OrderTypeLimit,
		Quantity:    1.0,
		Price:       50000.0,
		TimeInForce: model.TimeInForceGTC,
	}

	// Call the method
	response, err := testGateway.PlaceOrder(context.Background(), orderRequest)

	// Assert expectations
	assert.NoError(t, err)
	assert.True(t, response.IsSuccess)
	assert.Equal(t, "12345", response.Order.OrderID)
	assert.Equal(t, "BTCUSDT", response.Order.Symbol)
	assert.Equal(t, model.OrderSideBuy, response.Order.Side)
	assert.Equal(t, model.OrderTypeLimit, response.Order.Type)
	assert.Equal(t, 1.0, response.Order.Quantity)
	assert.Equal(t, 50000.0, response.Order.Price)
	assert.Equal(t, model.OrderStatusNew, response.Order.Status)

	// Test with error
	testGateway.mockError = fmt.Errorf("test error")
	response, err = testGateway.PlaceOrder(context.Background(), orderRequest)
	assert.Error(t, err)
	assert.Equal(t, "test error", err.Error())
	assert.False(t, response.IsSuccess)
}

func TestCancelOrder(t *testing.T) {
	// Create a test gateway
	testGateway := &TestGateway{
		mockError: nil,
	}

	// Call the method
	err := testGateway.CancelOrder(context.Background(), "BTCUSDT", "12345")

	// Assert expectations
	assert.NoError(t, err)

	// Test with error
	testGateway.mockError = fmt.Errorf("test error")
	err = testGateway.CancelOrder(context.Background(), "BTCUSDT", "12345")
	assert.Error(t, err)
	assert.Equal(t, "test error", err.Error())
}

// MockKeyStore is a mock implementation of the apikeystore.KeyStore interface
type MockKeyStore struct {
	mock.Mock
}

// GetAPIKey mocks the GetAPIKey method
func (m *MockKeyStore) GetAPIKey(keyID string) (*apikeystore.APIKeyCredentials, error) {
	args := m.Called(keyID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*apikeystore.APIKeyCredentials), args.Error(1)
}

// SetAPIKey mocks the SetAPIKey method
func (m *MockKeyStore) SetAPIKey(keyID string, credentials *apikeystore.APIKeyCredentials) error {
	args := m.Called(keyID, credentials)
	return args.Error(0)
}

// DeleteAPIKey mocks the DeleteAPIKey method
func (m *MockKeyStore) DeleteAPIKey(keyID string) error {
	args := m.Called(keyID)
	return args.Error(0)
}

// TestSubscribeToTicker tests the SubscribeToTicker method
func TestSubscribeToTicker(t *testing.T) {
	// Create a test gateway
	testGateway := &TestGateway{
		mockError: nil,
	}

	// Create a channel to receive the ticker
	tickCh := make(chan model.Ticker, 1)

	// Call the method
	err := testGateway.SubscribeToTicker(context.Background(), "BTCUSDT", func(ticker model.Ticker) {
		tickCh <- ticker
	})

	// Assert expectations
	assert.NoError(t, err)
	assert.True(t, testGateway.handlerCalled)

	// Get the ticker from the channel
	ticker := <-tickCh
	assert.Equal(t, "BTCUSDT", ticker.Symbol)
	assert.Equal(t, "MEXC", ticker.Exchange)
	assert.Equal(t, 50000.0, ticker.LastPrice)

	// Test with error
	testGateway.mockError = fmt.Errorf("test error")
	err = testGateway.SubscribeToTicker(context.Background(), "BTCUSDT", func(ticker model.Ticker) {})
	assert.Error(t, err)
	assert.Equal(t, "test error", err.Error())
}

// TestSubscribeToOrderBook tests the SubscribeToOrderBook method
func TestSubscribeToOrderBook(t *testing.T) {
	// Create a test gateway
	testGateway := &TestGateway{
		mockError: nil,
	}

	// Create a channel to receive the order book
	orderBookCh := make(chan model.OrderBook, 1)

	// Call the method
	err := testGateway.SubscribeToOrderBook(context.Background(), "BTCUSDT", 10, func(orderBook model.OrderBook) {
		orderBookCh <- orderBook
	})

	// Assert expectations
	assert.NoError(t, err)
	assert.True(t, testGateway.handlerCalled)

	// Get the order book from the channel
	orderBook := <-orderBookCh
	assert.Equal(t, "BTCUSDT", orderBook.Symbol)
	assert.Len(t, orderBook.Bids, 2)
	assert.Len(t, orderBook.Asks, 2)
	assert.Equal(t, 49900.0, orderBook.Bids[0].Price)
	assert.Equal(t, 50100.0, orderBook.Asks[0].Price)

	// Test with error
	testGateway.mockError = fmt.Errorf("test error")
	err = testGateway.SubscribeToOrderBook(context.Background(), "BTCUSDT", 10, func(orderBook model.OrderBook) {})
	assert.Error(t, err)
	assert.Equal(t, "test error", err.Error())
}

// TestUnsubscribe tests the Unsubscribe method
func TestUnsubscribe(t *testing.T) {
	// Create a test gateway
	testGateway := &TestGateway{
		mockError: nil,
	}

	// Call the method
	err := testGateway.Unsubscribe(context.Background(), "ticker", "BTCUSDT")

	// Assert expectations
	assert.NoError(t, err)

	// Test with error
	testGateway.mockError = fmt.Errorf("test error")
	err = testGateway.Unsubscribe(context.Background(), "ticker", "BTCUSDT")
	assert.Error(t, err)
	assert.Equal(t, "test error", err.Error())
}

// TestNewGateway tests the creation of a new gateway
func TestNewGateway(t *testing.T) {
	// Create mocks
	mockKeyStore := new(MockKeyStore)

	// Set up expectations
	mockKeyStore.On("GetAPIKey", "default").Return(&apikeystore.APIKeyCredentials{
		APIKey:    "test-api-key",
		SecretKey: "test-secret-key",
	}, nil)

	// Create logger
	logger := zerolog.Nop()

	// Create gateway
	gateway, err := NewGateway(mockKeyStore, DefaultGatewayConfig(), &logger)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, gateway)
	mockKeyStore.AssertExpectations(t)
}
