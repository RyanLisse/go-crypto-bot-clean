package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"go-crypto-bot-clean/backend/internal/domain/models"
	"github.com/stretchr/testify/assert"
)

// TestCoinHandler is a copy of the actual CoinHandler struct for testing
type TestCoinHandler struct {
	exchangeService interface {
		GetTicker(ctx context.Context, symbol string) (*models.Ticker, error)
		GetAllTickers(ctx context.Context) (map[string]*models.Ticker, error)
	}
	newCoinService interface {
		GetTradableCoins(ctx context.Context) ([]models.NewCoin, error)
		GetTradableCoinsToday(ctx context.Context) ([]models.NewCoin, error)
	}
}

// ListTradableCoinsToday handles GET /api/v1/newcoins/tradable/today
func (h *TestCoinHandler) ListTradableCoinsToday(c *gin.Context) {
	coins, err := h.newCoinService.GetTradableCoinsToday(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch tradable coins"})
		return
	}

	// Initialize with empty slice instead of nil to ensure [] instead of null in JSON
	resp := make([]map[string]interface{}, 0)
	for _, coin := range coins {
		var becameTradableAt int64
		if !coin.BecameTradableAt.IsZero() {
			becameTradableAt = coin.BecameTradableAt.Unix()
		}

		resp = append(resp, map[string]interface{}{
			"symbol":             coin.Symbol,
			"found_at":           coin.FoundAt.Unix(),
			"base_volume":        coin.BaseVolume,
			"quote_volume":       coin.QuoteVolume,
			"status":             coin.Status,
			"became_tradable_at": becameTradableAt,
		})
	}
	c.JSON(http.StatusOK, resp)
}

// ListTradableCoins handles GET /api/v1/newcoins/tradable
func (h *TestCoinHandler) ListTradableCoins(c *gin.Context) {
	coins, err := h.newCoinService.GetTradableCoins(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch tradable coins"})
		return
	}

	// Initialize with empty slice instead of nil to ensure [] instead of null in JSON
	resp := make([]map[string]interface{}, 0)
	for _, coin := range coins {
		var becameTradableAt int64
		if !coin.BecameTradableAt.IsZero() {
			becameTradableAt = coin.BecameTradableAt.Unix()
		}

		resp = append(resp, map[string]interface{}{
			"symbol":             coin.Symbol,
			"found_at":           coin.FoundAt.Unix(),
			"base_volume":        coin.BaseVolume,
			"quote_volume":       coin.QuoteVolume,
			"status":             coin.Status,
			"became_tradable_at": becameTradableAt,
		})
	}
	c.JSON(http.StatusOK, resp)
}

// MockExchangeService implements the ExchangeService interface for testing
type MockExchangeService struct {
	ShouldErr bool
}

// Connect implements service.ExchangeService.Connect
func (m *MockExchangeService) Connect(ctx context.Context) error {
	return nil
}

// Disconnect implements service.ExchangeService.Disconnect
func (m *MockExchangeService) Disconnect() error {
	return nil
}

// GetTicker implements service.ExchangeService.GetTicker
func (m *MockExchangeService) GetTicker(ctx context.Context, symbol string) (*models.Ticker, error) {
	return nil, nil
}

// GetAllTickers implements service.ExchangeService.GetAllTickers
func (m *MockExchangeService) GetAllTickers(ctx context.Context) (map[string]*models.Ticker, error) {
	return nil, nil
}

// GetKlines implements service.ExchangeService.GetKlines
func (m *MockExchangeService) GetKlines(ctx context.Context, symbol, interval string, limit int) ([]*models.Kline, error) {
	return nil, nil
}

// GetNewCoins implements service.ExchangeService.GetNewCoins
func (m *MockExchangeService) GetNewCoins(ctx context.Context) ([]*models.NewCoin, error) {
	return nil, nil
}

// GetWallet implements service.ExchangeService.GetWallet
func (m *MockExchangeService) GetWallet(ctx context.Context) (*models.Wallet, error) {
	return nil, nil
}

// PlaceOrder implements service.ExchangeService.PlaceOrder
func (m *MockExchangeService) PlaceOrder(ctx context.Context, order *models.Order) (*models.Order, error) {
	return nil, nil
}

// CancelOrder implements service.ExchangeService.CancelOrder
func (m *MockExchangeService) CancelOrder(ctx context.Context, orderID, symbol string) error {
	return nil
}

// GetOrder implements service.ExchangeService.GetOrder
func (m *MockExchangeService) GetOrder(ctx context.Context, orderID, symbol string) (*models.Order, error) {
	return nil, nil
}

// GetOpenOrders implements service.ExchangeService.GetOpenOrders
func (m *MockExchangeService) GetOpenOrders(ctx context.Context, symbol string) ([]*models.Order, error) {
	return nil, nil
}

// SubscribeToTickers implements service.ExchangeService.SubscribeToTickers
func (m *MockExchangeService) SubscribeToTickers(ctx context.Context, symbols []string, updates chan<- *models.Ticker) error {
	return nil
}

// UnsubscribeFromTickers implements service.ExchangeService.UnsubscribeFromTickers
func (m *MockExchangeService) UnsubscribeFromTickers(ctx context.Context, symbols []string) error {
	return nil
}

// MockNewCoinServiceForTest implements the newcoin.NewCoinService interface for testing
type MockNewCoinServiceForTest struct {
	ShouldErr bool
	MockCoins []*models.NewCoin
}

// DetectNewCoins implements newcoin.NewCoinService.DetectNewCoins
func (m *MockNewCoinServiceForTest) DetectNewCoins(ctx context.Context) ([]*models.NewCoin, error) {
	if m.ShouldErr {
		return nil, errors.New("mock error")
	}
	return m.MockCoins, nil
}

// SaveNewCoins implements newcoin.NewCoinService.SaveNewCoins
func (m *MockNewCoinServiceForTest) SaveNewCoins(ctx context.Context, coins []*models.NewCoin) error {
	return nil
}

// GetAllNewCoins implements newcoin.NewCoinService.GetAllNewCoins
func (m *MockNewCoinServiceForTest) GetAllNewCoins(ctx context.Context) ([]models.NewCoin, error) {
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

// StartWatching implements newcoin.NewCoinService.StartWatching
func (m *MockNewCoinServiceForTest) StartWatching(ctx context.Context, interval time.Duration) error {
	return nil
}

// StopWatching implements newcoin.NewCoinService.StopWatching
func (m *MockNewCoinServiceForTest) StopWatching() {
}

// MarkAsProcessed implements newcoin.NewCoinService.MarkAsProcessed
func (m *MockNewCoinServiceForTest) MarkAsProcessed(ctx context.Context, id int64) error {
	return nil
}

// GetCoinByID implements newcoin.NewCoinService.GetCoinByID
func (m *MockNewCoinServiceForTest) GetCoinByID(ctx context.Context, id int64) (*models.NewCoin, error) {
	return nil, nil
}

// GetCoinsByDate implements newcoin.NewCoinService.GetCoinsByDate
func (m *MockNewCoinServiceForTest) GetCoinsByDate(ctx context.Context, date time.Time) ([]models.NewCoin, error) {
	return nil, nil
}

// GetCoinsByDateRange implements newcoin.NewCoinService.GetCoinsByDateRange
func (m *MockNewCoinServiceForTest) GetCoinsByDateRange(ctx context.Context, startDate, endDate time.Time) ([]models.NewCoin, error) {
	return nil, nil
}

// GetUpcomingCoins implements newcoin.NewCoinService.GetUpcomingCoins
func (m *MockNewCoinServiceForTest) GetUpcomingCoins(ctx context.Context) ([]models.NewCoin, error) {
	return nil, nil
}

// GetUpcomingCoinsByDate implements newcoin.NewCoinService.GetUpcomingCoinsByDate
func (m *MockNewCoinServiceForTest) GetUpcomingCoinsByDate(ctx context.Context, date time.Time) ([]models.NewCoin, error) {
	return nil, nil
}

// GetUpcomingCoinsForTodayAndTomorrow implements newcoin.NewCoinService.GetUpcomingCoinsForTodayAndTomorrow
func (m *MockNewCoinServiceForTest) GetUpcomingCoinsForTodayAndTomorrow(ctx context.Context) ([]models.NewCoin, error) {
	return nil, nil
}

// UpdateCoinStatus implements newcoin.NewCoinService.UpdateCoinStatus
func (m *MockNewCoinServiceForTest) UpdateCoinStatus(ctx context.Context, symbol string, status string) error {
	return nil
}

// GetTradableCoins implements newcoin.NewCoinService.GetTradableCoins
func (m *MockNewCoinServiceForTest) GetTradableCoins(ctx context.Context) ([]models.NewCoin, error) {
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

// GetTradableCoinsByDate implements newcoin.NewCoinService.GetTradableCoinsByDate
func (m *MockNewCoinServiceForTest) GetTradableCoinsByDate(ctx context.Context, date time.Time) ([]models.NewCoin, error) {
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

// GetTradableCoinsToday implements newcoin.NewCoinService.GetTradableCoinsToday
func (m *MockNewCoinServiceForTest) GetTradableCoinsToday(ctx context.Context) ([]models.NewCoin, error) {
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

func TestListTradableCoinsTodayNew(t *testing.T) {
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
			// Setup mock services
			mockExchangeService := &MockExchangeService{}
			mockNewCoinService := &MockNewCoinServiceForTest{
				ShouldErr: tt.mockError,
				MockCoins: tt.mockCoins,
			}

			// Create handler
			handler := NewCoinHandler(mockExchangeService, mockNewCoinService)

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

func TestListTradableCoinsNew(t *testing.T) {
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
			// Setup mock services
			mockExchangeService := &MockExchangeService{}
			mockNewCoinService := &MockNewCoinServiceForTest{
				ShouldErr: tt.mockError,
				MockCoins: tt.mockCoins,
			}

			// Create handler
			handler := NewCoinHandler(mockExchangeService, mockNewCoinService)

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
