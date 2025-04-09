package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ryanlisse/go-crypto-bot/internal/core/newcoin"
	"github.com/ryanlisse/go-crypto-bot/internal/domain/models"
	"github.com/stretchr/testify/assert"
)

// MockCoinService implements the ExchangeService interface for testing
type MockCoinService struct {
	ShouldErr  bool
	MockTicker *models.Ticker
}

// Connect implements ExchangeService.Connect
func (m *MockCoinService) Connect(ctx context.Context) error {
	return nil
}

// Disconnect implements ExchangeService.Disconnect
func (m *MockCoinService) Disconnect() error {
	return nil
}

// GetTicker implements ExchangeService.GetTicker
func (m *MockCoinService) GetTicker(ctx context.Context, symbol string) (*models.Ticker, error) {
	if m.ShouldErr {
		return nil, errors.New("mock service error")
	}

	if m.MockTicker != nil {
		return m.MockTicker, nil
	}

	if symbol == "BTC/USDT" {
		return &models.Ticker{
			Symbol:      "BTC/USDT",
			Price:       40000.0,
			Volume:      100.0,
			QuoteVolume: 4000000.0,
			Timestamp:   time.Now(),
		}, nil
	} else if symbol == "ETH/USDT" {
		return &models.Ticker{
			Symbol:      "ETH/USDT",
			Price:       2000.0,
			Volume:      1000.0,
			QuoteVolume: 2000000.0,
			Timestamp:   time.Now(),
		}, nil
	}

	// Return nil for non-existent symbol
	return nil, nil
}

// GetKlines implements ExchangeService.GetKlines
func (m *MockCoinService) GetKlines(ctx context.Context, symbol, interval string, limit int) ([]*models.Kline, error) {
	return nil, nil
}

// GetNewCoins implements ExchangeService.GetNewCoins
func (m *MockCoinService) GetNewCoins(ctx context.Context) ([]*models.NewCoin, error) {
	return nil, nil
}

// GetWallet implements ExchangeService.GetWallet
func (m *MockCoinService) GetWallet(ctx context.Context) (*models.Wallet, error) {
	return nil, nil
}

// PlaceOrder implements ExchangeService.PlaceOrder
func (m *MockCoinService) PlaceOrder(ctx context.Context, order *models.Order) (*models.Order, error) {
	return nil, nil
}

// CancelOrder implements ExchangeService.CancelOrder
func (m *MockCoinService) CancelOrder(ctx context.Context, orderID, symbol string) error {
	return nil
}

// GetOrder implements ExchangeService.GetOrder
func (m *MockCoinService) GetOrder(ctx context.Context, orderID, symbol string) (*models.Order, error) {
	return nil, nil
}

// GetOpenOrders implements ExchangeService.GetOpenOrders
func (m *MockCoinService) GetOpenOrders(ctx context.Context, symbol string) ([]*models.Order, error) {
	return nil, nil
}

// SubscribeToTickers implements ExchangeService.SubscribeToTickers
func (m *MockCoinService) SubscribeToTickers(ctx context.Context, symbols []string, updates chan<- *models.Ticker) error {
	return nil
}

// UnsubscribeFromTickers implements ExchangeService.UnsubscribeFromTickers
func (m *MockCoinService) UnsubscribeFromTickers(ctx context.Context, symbols []string) error {
	return nil
}

// Override GetAllTickers to return test-specific values
func (m *MockCoinService) GetAllTickers(ctx context.Context) (map[string]*models.Ticker, error) {
	if m.ShouldErr {
		return nil, errors.New("mock service error")
	}

	// Return mock tickers
	tickers := map[string]*models.Ticker{
		"BTC/USDT": {
			Symbol:      "BTC/USDT",
			Price:       40000.0,
			Volume:      100.0,
			QuoteVolume: 4000000.0,
			Timestamp:   time.Now(),
		},
		"ETH/USDT": {
			Symbol:      "ETH/USDT",
			Price:       2000.0,
			Volume:      1000.0,
			QuoteVolume: 2000000.0,
			Timestamp:   time.Now(),
		},
	}

	return tickers, nil
}

// This is a duplicate method, removing it

// MockNewCoinWatcher implements the NewCoinWatcher interface for testing
type MockNewCoinWatcher struct {
	ShouldErr bool
	MockCoins []*models.NewCoin
}

func (m *MockNewCoinWatcher) DetectNewCoins(ctx context.Context) ([]*models.NewCoin, error) {
	if m.ShouldErr {
		return nil, errors.New("mock watcher error")
	}
	return m.MockCoins, nil
}

func (m *MockNewCoinWatcher) StartWatching(ctx context.Context, interval time.Duration) error {
	return nil
}

func (m *MockNewCoinWatcher) StopWatching() {
	// No implementation needed for tests
}
func (m *MockNewCoinWatcher) SaveNewCoins(ctx context.Context, coins []*models.NewCoin) error {
	return nil
}

func (m *MockNewCoinWatcher) GetCoinByID(ctx context.Context, id int64) (*models.NewCoin, error) {
	if m.ShouldErr {
		return nil, errors.New("mock watcher error")
	}

	// Return a mock coin
	return &models.NewCoin{
		ID:      id,
		Symbol:  "TEST/USDT",
		FoundAt: time.Now(),
	}, nil
}

func (m *MockNewCoinWatcher) GetAllNewCoins(ctx context.Context) ([]models.NewCoin, error) {
	return nil, nil
}

func (m *MockNewCoinWatcher) MarkAsProcessed(ctx context.Context, id int64) error {
	if m.ShouldErr {
		return errors.New("mock watcher error")
	}
	return nil
}

func (m *MockNewCoinWatcher) GetCoinsByDate(ctx context.Context, date time.Time) ([]models.NewCoin, error) {
	if m.ShouldErr {
		return nil, errors.New("mock watcher error")
	}

	// Convert pointer slice to value slice
	result := make([]models.NewCoin, 0, len(m.MockCoins))
	for _, coin := range m.MockCoins {
		result = append(result, *coin)
	}

	return result, nil
}

func (m *MockNewCoinWatcher) GetCoinsByDateRange(ctx context.Context, startDate, endDate time.Time) ([]models.NewCoin, error) {
	if m.ShouldErr {
		return nil, errors.New("mock watcher error")
	}

	// Convert pointer slice to value slice
	result := make([]models.NewCoin, 0, len(m.MockCoins))
	for _, coin := range m.MockCoins {
		result = append(result, *coin)
	}

	return result, nil
}

func (m *MockNewCoinWatcher) GetUpcomingCoins(ctx context.Context) ([]models.NewCoin, error) {
	if m.ShouldErr {
		return nil, errors.New("mock error")
	}

	// Convert pointer slice to value slice
	result := make([]models.NewCoin, 0, len(m.MockCoins))
	for _, coin := range m.MockCoins {
		result = append(result, *coin)
	}

	return result, nil
}

func (m *MockNewCoinWatcher) GetUpcomingCoinsByDate(ctx context.Context, date time.Time) ([]models.NewCoin, error) {
	if m.ShouldErr {
		return nil, errors.New("mock error")
	}

	// Convert pointer slice to value slice
	result := make([]models.NewCoin, 0, len(m.MockCoins))
	for _, coin := range m.MockCoins {
		result = append(result, *coin)
	}

	return result, nil
}

func (m *MockNewCoinWatcher) GetUpcomingCoinsForTodayAndTomorrow(ctx context.Context) ([]models.NewCoin, error) {
	if m.ShouldErr {
		return nil, errors.New("mock error")
	}

	// Convert pointer slice to value slice
	result := make([]models.NewCoin, 0, len(m.MockCoins))
	for _, coin := range m.MockCoins {
		result = append(result, *coin)
	}

	return result, nil
}

func (m *MockNewCoinWatcher) UpdateCoinStatus(ctx context.Context, symbol string, status string) error {
	if m.ShouldErr {
		return errors.New("mock error")
	}
	return nil
}

func (m *MockNewCoinWatcher) GetTradableCoins(ctx context.Context) ([]models.NewCoin, error) {
	if m.ShouldErr {
		return nil, errors.New("mock error")
	}

	// Convert pointer slice to value slice
	result := make([]models.NewCoin, 0, len(m.MockCoins))
	for _, coin := range m.MockCoins {
		result = append(result, *coin)
	}

	return result, nil
}

func (m *MockNewCoinWatcher) GetTradableCoinsByDate(ctx context.Context, date time.Time) ([]models.NewCoin, error) {
	if m.ShouldErr {
		return nil, errors.New("mock error")
	}

	// Convert pointer slice to value slice
	result := make([]models.NewCoin, 0, len(m.MockCoins))
	for _, coin := range m.MockCoins {
		result = append(result, *coin)
	}

	return result, nil
}

func (m *MockNewCoinWatcher) GetTradableCoinsToday(ctx context.Context) ([]models.NewCoin, error) {
	if m.ShouldErr {
		return nil, errors.New("mock error")
	}

	// Convert pointer slice to value slice
	result := make([]models.NewCoin, 0, len(m.MockCoins))
	for _, coin := range m.MockCoins {
		result = append(result, *coin)
	}

	return result, nil
}

func TestListTradableCoinsToday(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		mockCoins  []*models.NewCoin
		mockError  bool
		wantStatus int
		wantEmpty  bool
	}{
		{
			name: "Success with coins",
			mockCoins: []*models.NewCoin{
				{
					ID:               1,
					Symbol:           "TODAYTRADABLECOIN1USDT",
					FoundAt:          time.Now().Add(-12 * time.Hour),
					FirstOpenTime:    func() *time.Time { t := time.Now().Add(-6 * time.Hour); return &t }(),
					QuoteVolume:      7000.0,
					Status:           "1",
					BecameTradableAt: func() *time.Time { t := time.Now().Add(-3 * time.Hour); return &t }(),
				},
			},
			mockError:  false,
			wantStatus: http.StatusOK,
			wantEmpty:  false,
		},
		{
			name:       "Success with empty list",
			mockCoins:  []*models.NewCoin{},
			mockError:  false,
			wantStatus: http.StatusOK,
			wantEmpty:  true,
		},
		{
			name:       "Error from service",
			mockCoins:  nil,
			mockError:  true,
			wantStatus: http.StatusInternalServerError,
			wantEmpty:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock service
			mockService := &MockCoinService{}

			// Setup mock watcher with test data
			mockWatcher := &MockNewCoinWatcher{
				ShouldErr: tt.mockError,
				MockCoins: tt.mockCoins,
			}

			// Create handler
			handler := NewCoinHandler(mockService, mockWatcher)

			// Setup router
			router := gin.New()
			router.GET("/api/v1/newcoins/tradable/today", handler.ListTradableCoinsToday)

			// Create request
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/newcoins/tradable/today", nil)
			resp := httptest.NewRecorder()

			// Perform request
			router.ServeHTTP(resp, req)

			// Assert status code
			assert.Equal(t, tt.wantStatus, resp.Code)

			// Check response based on test case
			if tt.mockError {
				assert.Contains(t, resp.Body.String(), "error")
			} else if tt.wantEmpty {
				assert.Equal(t, "[]", resp.Body.String())
			} else {
				assert.NotEqual(t, "[]", resp.Body.String())
				assert.Contains(t, resp.Body.String(), "TODAYTRADABLECOIN1USDT")
			}
		})
	}
}

func TestListTradableCoins(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		mockCoins  []*models.NewCoin
		mockError  bool
		wantStatus int
		wantEmpty  bool
	}{
		{
			name: "Success with coins",
			mockCoins: []*models.NewCoin{
				{
					ID:               1,
					Symbol:           "TRADABLECOIN1USDT",
					FoundAt:          time.Now().Add(-24 * time.Hour),
					FirstOpenTime:    func() *time.Time { t := time.Now().Add(-12 * time.Hour); return &t }(),
					QuoteVolume:      5000.0,
					Status:           "1",
					BecameTradableAt: func() *time.Time { t := time.Now().Add(-6 * time.Hour); return &t }(),
				},
			},
			mockError:  false,
			wantStatus: http.StatusOK,
			wantEmpty:  false,
		},
		{
			name:       "Success with empty list",
			mockCoins:  []*models.NewCoin{},
			mockError:  false,
			wantStatus: http.StatusOK,
			wantEmpty:  true,
		},
		{
			name:       "Error from service",
			mockCoins:  nil,
			mockError:  true,
			wantStatus: http.StatusInternalServerError,
			wantEmpty:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock service
			mockService := &MockCoinService{}

			// Setup mock watcher with test data
			mockWatcher := &MockNewCoinWatcher{
				ShouldErr: tt.mockError,
				MockCoins: tt.mockCoins,
			}

			// Create handler
			handler := NewCoinHandler(mockService, mockWatcher)

			// Setup router
			router := gin.New()
			router.GET("/api/v1/newcoins/tradable", handler.ListTradableCoins)

			// Create request
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/newcoins/tradable", nil)
			resp := httptest.NewRecorder()

			// Perform request
			router.ServeHTTP(resp, req)

			// Assert status code
			assert.Equal(t, tt.wantStatus, resp.Code)

			// Check response based on test case
			if tt.mockError {
				assert.Contains(t, resp.Body.String(), "error")
			} else if tt.wantEmpty {
				assert.Equal(t, "[]", resp.Body.String())
			} else {
				assert.NotEqual(t, "[]", resp.Body.String())
				assert.Contains(t, resp.Body.String(), "TRADABLECOIN1USDT")
			}
		})
	}
}

func TestListMarkets(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		mockError      bool
		expectedStatus int
		expectedCount  int
	}{
		{
			name:           "returns markets successfully",
			mockError:      false,
			expectedStatus: http.StatusOK,
			expectedCount:  2,
		},
		{
			name:           "handles service error",
			mockError:      true,
			expectedStatus: http.StatusInternalServerError,
			expectedCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up mock services
			mockService := &MockCoinService{
				ShouldErr: tt.mockError,
			}
			// Explicitly declaring type to use the newcoin package
			var mockWatcher newcoin.NewCoinService = &MockNewCoinWatcher{}

			// Create handler
			handler := NewCoinHandler(mockService, mockWatcher)

			// Set up router
			router := gin.New()
			router.GET("/test", handler.ListMarkets)

			// Create request
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			w := httptest.NewRecorder()

			// Perform request
			router.ServeHTTP(w, req)

			// Assert status code
			assert.Equal(t, tt.expectedStatus, w.Code)

			if !tt.mockError {
				// Check for specific data in the response
				assert.Contains(t, w.Body.String(), "BTC/USDT")
				assert.Contains(t, w.Body.String(), "ETH/USDT")
				assert.Contains(t, w.Body.String(), "40000")
				assert.Contains(t, w.Body.String(), "2000")
			} else {
				// Check error message
				assert.Contains(t, w.Body.String(), "failed to fetch markets")
			}
		})
	}
}

func TestGetMarket(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		symbol         string
		mockError      bool
		expectedStatus int
		symbolExists   bool
	}{
		{
			name:           "returns existing market successfully",
			symbol:         "BTC/USDT",
			mockError:      false,
			expectedStatus: http.StatusOK,
			symbolExists:   true,
		},
		{
			name:           "returns 404 for non-existent market",
			symbol:         "NONEXIST/USDT",
			mockError:      false,
			expectedStatus: http.StatusNotFound,
			symbolExists:   false,
		},
		{
			name:           "handles service error",
			symbol:         "BTC/USDT",
			mockError:      true,
			expectedStatus: http.StatusInternalServerError,
			symbolExists:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up mock services
			mockService := &MockCoinService{
				ShouldErr: tt.mockError,
			}
			// Explicitly declaring type to use the newcoin package
			var mockWatcher newcoin.NewCoinService = &MockNewCoinWatcher{}

			// Create handler
			handler := NewCoinHandler(mockService, mockWatcher)

			// For testing Gin handlers with path parameters, it's better to create a context directly
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Create request
			req := httptest.NewRequest(http.MethodGet, "/market/"+tt.symbol, nil)
			c.Request = req

			// Set the path parameter directly
			c.Params = gin.Params{{Key: "symbol", Value: tt.symbol}}

			// Call the handler directly
			handler.GetMarket(c)

			// Assert status code
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				assert.Contains(t, w.Body.String(), tt.symbol)
				if tt.symbol == "BTC/USDT" {
					assert.Contains(t, w.Body.String(), "40000")
				} else if tt.symbol == "ETH/USDT" {
					assert.Contains(t, w.Body.String(), "2000")
				}
			} else if tt.expectedStatus == http.StatusNotFound {
				assert.Contains(t, w.Body.String(), "market not found")
			} else {
				assert.Contains(t, w.Body.String(), "failed to fetch market")
			}
		})
	}
}

func TestListNewCoins(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create test new coins
	testTime := time.Now()
	testCoins := []*models.NewCoin{
		{
			Symbol:      "NEW1/USDT",
			FoundAt:     testTime,
			BaseVolume:  1000.0,
			QuoteVolume: 50000.0,
		},
		{
			Symbol:      "NEW2/USDT",
			FoundAt:     testTime,
			BaseVolume:  500.0,
			QuoteVolume: 25000.0,
		},
	}

	tests := []struct {
		name           string
		mockCoins      []*models.NewCoin
		mockError      bool
		expectedStatus int
		expectedCount  int
	}{
		{
			name:           "returns new coins successfully",
			mockCoins:      testCoins,
			mockError:      false,
			expectedStatus: http.StatusOK,
			expectedCount:  2,
		},
		{
			name:           "returns empty list when no new coins",
			mockCoins:      []*models.NewCoin{},
			mockError:      false,
			expectedStatus: http.StatusOK,
			expectedCount:  0,
		},
		{
			name:           "handles service error",
			mockCoins:      nil,
			mockError:      true,
			expectedStatus: http.StatusInternalServerError,
			expectedCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up mock services
			mockService := &MockCoinService{}

			// Explicitly declaring type to use the newcoin package
			var mockWatcher newcoin.NewCoinService = &MockNewCoinWatcher{
				ShouldErr: tt.mockError,
				MockCoins: tt.mockCoins,
			}

			// Create handler
			handler := NewCoinHandler(mockService, mockWatcher)

			// For testing Gin handlers, create a context directly
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Create request
			req := httptest.NewRequest(http.MethodGet, "/newcoins", nil)
			c.Request = req

			// Call the handler directly
			handler.ListNewCoins(c)

			// Assert status code
			assert.Equal(t, tt.expectedStatus, w.Code)

			if !tt.mockError {
				if tt.expectedCount > 0 {
					assert.Contains(t, w.Body.String(), "NEW1/USDT")
					assert.Contains(t, w.Body.String(), "NEW2/USDT")
					assert.Contains(t, w.Body.String(), "50000")
					assert.Contains(t, w.Body.String(), "25000")
				}

				// Empty array should be returned for no coins case
				if tt.expectedCount == 0 && !tt.mockError {
					assert.Equal(t, "[]", w.Body.String())
				}
			} else {
				// Check error message
				assert.Contains(t, w.Body.String(), "failed to fetch new coins")
			}
		})
	}
}
