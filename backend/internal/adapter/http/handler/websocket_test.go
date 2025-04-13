package handler

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/neo/crypto-bot/internal/domain/model/market"
	"github.com/neo/crypto-bot/internal/domain/port"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockMarketDataUseCase is a mock implementation of the MarketDataUseCaseInterface
type MockMarketDataUseCase struct {
	mock.Mock
}

// Ensure MockMarketDataUseCase implements the MarketDataUseCaseInterface
var _ port.MarketDataUseCaseInterface = (*MockMarketDataUseCase)(nil)

func (m *MockMarketDataUseCase) GetTicker(ctx context.Context, exchange, symbol string) (*market.Ticker, error) {
	args := m.Called(ctx, exchange, symbol)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*market.Ticker), args.Error(1)
}

func (m *MockMarketDataUseCase) GetLatestTickers(ctx context.Context) ([]market.Ticker, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	// Convert to the expected return type
	pointerSlice := args.Get(0).([]*market.Ticker)
	result := make([]market.Ticker, len(pointerSlice))
	for i, ticker := range pointerSlice {
		result[i] = *ticker
	}
	return result, args.Error(1)
}

func (m *MockMarketDataUseCase) GetAllTickers(ctx context.Context, exchange string) ([]*market.Ticker, error) {
	args := m.Called(ctx, exchange)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*market.Ticker), args.Error(1)
}

func (m *MockMarketDataUseCase) GetCandles(ctx context.Context, exchange, symbol string, interval market.Interval, start, end time.Time, limit int) ([]market.Candle, error) {
	args := m.Called(ctx, exchange, symbol, interval, start, end, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	// Convert to the expected return type
	pointerSlice := args.Get(0).([]*market.Candle)
	result := make([]market.Candle, len(pointerSlice))
	for i, candle := range pointerSlice {
		result[i] = *candle
	}
	return result, args.Error(1)
}

func (m *MockMarketDataUseCase) GetSymbols(ctx context.Context) ([]*market.Symbol, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*market.Symbol), args.Error(1)
}

func (m *MockMarketDataUseCase) GetSymbol(ctx context.Context, symbol string) (*market.Symbol, error) {
	args := m.Called(ctx, symbol)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*market.Symbol), args.Error(1)
}

func (m *MockMarketDataUseCase) GetOrderBook(ctx context.Context, exchange, symbol string) (*market.OrderBook, error) {
	args := m.Called(ctx, exchange, symbol)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*market.OrderBook), args.Error(1)
}

// setupWebSocketTest creates a test server with the WebSocket handler
func setupWebSocketTest(_ *testing.T) (*MockMarketDataUseCase, *httptest.Server, string) {
	gin.SetMode(gin.TestMode)

	// Create mock use case
	mockUseCase := new(MockMarketDataUseCase)

	// Create logger
	logger := zerolog.New(zerolog.NewConsoleWriter()).With().Timestamp().Logger()

	// Create router
	router := gin.New()

	// Create WebSocket handler
	wsHandler := NewWebSocketHandler(mockUseCase, &logger)

	// Register routes
	apiGroup := router.Group("/api/v1")
	wsHandler.RegisterRoutes(apiGroup)

	// Create test server
	server := httptest.NewServer(router)

	// Convert HTTP URL to WebSocket URL
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/api/v1/ws/market"

	return mockUseCase, server, wsURL
}

func TestWebSocketConnection(t *testing.T) {
	mockUseCase, server, wsURL := setupWebSocketTest(t)
	defer server.Close()

	// Connect to WebSocket
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err)
	defer ws.Close()

	// Read welcome message
	var welcomeMsg WebSocketMessage
	err = ws.ReadJSON(&welcomeMsg)
	require.NoError(t, err)

	// Verify welcome message
	assert.Equal(t, "info", welcomeMsg.Type)
	assert.Equal(t, "system", welcomeMsg.Channel)
	assert.Equal(t, "Connected to market data WebSocket", welcomeMsg.Data)

	// No calls to the use case should have been made yet
	mockUseCase.AssertNotCalled(t, "GetTicker")
	mockUseCase.AssertNotCalled(t, "GetLatestTickers")
}

func TestTickerSubscription(t *testing.T) {
	mockUseCase, server, wsURL := setupWebSocketTest(t)
	defer server.Close()

	// Create a test ticker
	testTicker := &market.Ticker{
		Symbol:        "BTCUSDT",
		Exchange:      "mexc",
		Price:         50000.0,
		Volume:        100.0,
		High24h:       51000.0,
		Low24h:        49000.0,
		PriceChange:   1000.0,
		PercentChange: 2.0,
		LastUpdated:   time.Now(),
	}

	// Setup mock expectations
	mockUseCase.On("GetTicker", mock.Anything, "mexc", "BTCUSDT").Return(testTicker, nil)

	// Connect to WebSocket
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err)
	defer ws.Close()

	// Read welcome message
	var welcomeMsg WebSocketMessage
	err = ws.ReadJSON(&welcomeMsg)
	require.NoError(t, err)

	// Send subscription request
	subscriptionReq := SubscriptionRequest{
		Action:  "subscribe",
		Channel: "tickers",
		Symbols: []string{"BTCUSDT"},
	}
	err = ws.WriteJSON(subscriptionReq)
	require.NoError(t, err)

	// Read subscription response
	var subResponse WebSocketMessage
	err = ws.ReadJSON(&subResponse)
	require.NoError(t, err)

	// Verify subscription response
	assert.Equal(t, "subscription", subResponse.Type)
	assert.Equal(t, "tickers", subResponse.Channel)

	// Read ticker data
	var tickerMsg WebSocketMessage
	err = ws.ReadJSON(&tickerMsg)
	require.NoError(t, err)

	// Verify ticker data
	assert.Equal(t, "data", tickerMsg.Type)
	assert.Equal(t, "tickers", tickerMsg.Channel)
	assert.Equal(t, "BTCUSDT", tickerMsg.Symbol)

	// Verify mock was called
	mockUseCase.AssertExpectations(t)
}

func TestCandleSubscription(t *testing.T) {
	mockUseCase, server, wsURL := setupWebSocketTest(t)
	defer server.Close()

	// Create test candles
	now := time.Now()
	testCandles := []*market.Candle{
		{
			Symbol:    "BTCUSDT",
			Exchange:  "mexc",
			Interval:  market.Interval1h,
			OpenTime:  now.Add(-1 * time.Hour),
			CloseTime: now,
			Open:      49000.0,
			High:      50000.0,
			Low:       48000.0,
			Close:     49500.0,
			Volume:    100.0,
			Complete:  true,
		},
	}

	// Setup mock expectations
	mockUseCase.On("GetCandles",
		mock.Anything,
		"mexc",
		"BTCUSDT",
		market.Interval1h,
		mock.AnythingOfType("time.Time"),
		mock.AnythingOfType("time.Time"),
		100,
	).Return(testCandles, nil)

	// Connect to WebSocket
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err)
	defer ws.Close()

	// Read welcome message
	var welcomeMsg WebSocketMessage
	err = ws.ReadJSON(&welcomeMsg)
	require.NoError(t, err)

	// Send subscription request
	subscriptionReq := SubscriptionRequest{
		Action:   "subscribe",
		Channel:  "candles",
		Symbols:  []string{"BTCUSDT"},
		Interval: "1h",
	}
	err = ws.WriteJSON(subscriptionReq)
	require.NoError(t, err)

	// Read subscription response
	var subResponse WebSocketMessage
	err = ws.ReadJSON(&subResponse)
	require.NoError(t, err)

	// Verify subscription response
	assert.Equal(t, "subscription", subResponse.Type)
	assert.Equal(t, "candles", subResponse.Channel)

	// Read candle data
	var candleMsg WebSocketMessage
	err = ws.ReadJSON(&candleMsg)
	require.NoError(t, err)

	// Verify candle data
	assert.Equal(t, "data", candleMsg.Type)
	assert.Equal(t, "candles", candleMsg.Channel)
	assert.Equal(t, "BTCUSDT", candleMsg.Symbol)

	// Verify mock was called
	mockUseCase.AssertExpectations(t)
}

func TestInvalidSubscriptionRequest(t *testing.T) {
	_, server, wsURL := setupWebSocketTest(t)
	defer server.Close()

	// Connect to WebSocket
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err)
	defer ws.Close()

	// Read welcome message
	var welcomeMsg WebSocketMessage
	err = ws.ReadJSON(&welcomeMsg)
	require.NoError(t, err)

	// Send invalid subscription request (missing action)
	invalidReq := map[string]interface{}{
		"channel": "tickers",
		"symbols": []string{"BTCUSDT"},
	}
	err = ws.WriteJSON(invalidReq)
	require.NoError(t, err)

	// Read error response
	var errorMsg WebSocketMessage
	err = ws.ReadJSON(&errorMsg)
	require.NoError(t, err)

	// Verify error response
	assert.Equal(t, "error", errorMsg.Type)
	assert.Equal(t, "system", errorMsg.Channel)
	assert.Contains(t, errorMsg.Data.(string), "Invalid")
}

func TestUnsupportedChannel(t *testing.T) {
	_, server, wsURL := setupWebSocketTest(t)
	defer server.Close()

	// Connect to WebSocket
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err)
	defer ws.Close()

	// Read welcome message
	var welcomeMsg WebSocketMessage
	err = ws.ReadJSON(&welcomeMsg)
	require.NoError(t, err)

	// Send subscription request with unsupported channel
	subscriptionReq := SubscriptionRequest{
		Action:  "subscribe",
		Channel: "unsupported",
		Symbols: []string{"BTCUSDT"},
	}
	err = ws.WriteJSON(subscriptionReq)
	require.NoError(t, err)

	// Read error response
	var errorMsg WebSocketMessage
	err = ws.ReadJSON(&errorMsg)
	require.NoError(t, err)

	// Verify error response
	assert.Equal(t, "error", errorMsg.Type)
	assert.Equal(t, "system", errorMsg.Channel)
	assert.Contains(t, errorMsg.Data.(string), "Unsupported channel")
}

func TestMissingIntervalForCandles(t *testing.T) {
	_, server, wsURL := setupWebSocketTest(t)
	defer server.Close()

	// Connect to WebSocket
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err)
	defer ws.Close()

	// Read welcome message
	var welcomeMsg WebSocketMessage
	err = ws.ReadJSON(&welcomeMsg)
	require.NoError(t, err)

	// Send candle subscription without interval
	subscriptionReq := SubscriptionRequest{
		Action:  "subscribe",
		Channel: "candles",
		Symbols: []string{"BTCUSDT"},
		// Missing interval
	}
	err = ws.WriteJSON(subscriptionReq)
	require.NoError(t, err)

	// Read error response
	var errorMsg WebSocketMessage
	err = ws.ReadJSON(&errorMsg)
	require.NoError(t, err)

	// Verify error response
	assert.Equal(t, "error", errorMsg.Type)
	assert.Equal(t, "candles", errorMsg.Channel)
	assert.Contains(t, errorMsg.Data.(string), "Interval is required")
}
