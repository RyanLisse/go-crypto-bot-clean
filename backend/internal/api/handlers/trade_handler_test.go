package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ryanlisse/go-crypto-bot/internal/domain/models"
	"github.com/stretchr/testify/assert"
)

type mockTradeService struct {
	EvaluatePurchaseDecisionFunc func(ctx context.Context, symbol string) (*models.PurchaseDecision, error)
	ExecutePurchaseFunc          func(ctx context.Context, symbol string, amount float64) (*models.BoughtCoin, error)
	CheckStopLossFunc            func(ctx context.Context, coin *models.BoughtCoin) (bool, error)
	CheckTakeProfitFunc          func(ctx context.Context, coin *models.BoughtCoin) (bool, error)
	SellCoinFunc                 func(ctx context.Context, coin *models.BoughtCoin, amount float64) (*models.Order, error)
}

func (m *mockTradeService) EvaluatePurchaseDecision(ctx context.Context, symbol string) (*models.PurchaseDecision, error) {
	if m.EvaluatePurchaseDecisionFunc != nil {
		return m.EvaluatePurchaseDecisionFunc(ctx, symbol)
	}
	return &models.PurchaseDecision{Decision: true, Reason: ""}, nil
}

func (m *mockTradeService) CheckStopLoss(ctx context.Context, coin *models.BoughtCoin) (bool, error) {
	if m.CheckStopLossFunc != nil {
		return m.CheckStopLossFunc(ctx, coin)
	}
	return false, nil
}

func (m *mockTradeService) CheckTakeProfit(ctx context.Context, coin *models.BoughtCoin) (bool, error) {
	if m.CheckTakeProfitFunc != nil {
		return m.CheckTakeProfitFunc(ctx, coin)
	}
	return false, nil
}

func (m *mockTradeService) ExecutePurchase(ctx context.Context, symbol string, amount float64, options *models.PurchaseOptions) (*models.BoughtCoin, error) {
	if m.ExecutePurchaseFunc != nil {
		return m.ExecutePurchaseFunc(ctx, symbol, amount)
	}
	return nil, errors.New("not implemented")
}

func (m *mockTradeService) SellCoin(ctx context.Context, coin *models.BoughtCoin, amount float64) (*models.Order, error) {
	if m.SellCoinFunc != nil {
		return m.SellCoinFunc(ctx, coin, amount)
	}
	return nil, errors.New("not implemented")
}

// Mock BoughtCoinRepository for testing
type mockBoughtCoinRepository struct {
	FindAllFunc       func(ctx context.Context) ([]*models.BoughtCoin, error)
	FindByIDFunc      func(ctx context.Context, id int64) (*models.BoughtCoin, error)
	FindBySymbolFunc  func(ctx context.Context, symbol string) (*models.BoughtCoin, error)
	SaveFunc          func(ctx context.Context, coin *models.BoughtCoin) error
	DeleteFunc        func(ctx context.Context, symbol string) error
	DeleteByIDFunc    func(ctx context.Context, id int64) error
	UpdatePriceFunc   func(ctx context.Context, symbol string, price float64) error
	FindAllActiveFunc func(ctx context.Context) ([]*models.BoughtCoin, error)
	HardDeleteFunc    func(ctx context.Context, symbol string) error
	CountFunc         func(ctx context.Context) (int64, error)
}

func (m *mockBoughtCoinRepository) FindAll(ctx context.Context) ([]*models.BoughtCoin, error) {
	if m.FindAllFunc != nil {
		return m.FindAllFunc(ctx)
	}
	return []*models.BoughtCoin{}, nil
}

func (m *mockBoughtCoinRepository) FindByID(ctx context.Context, id int64) (*models.BoughtCoin, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *mockBoughtCoinRepository) FindBySymbol(ctx context.Context, symbol string) (*models.BoughtCoin, error) {
	if m.FindBySymbolFunc != nil {
		return m.FindBySymbolFunc(ctx, symbol)
	}
	return nil, nil
}

func (m *mockBoughtCoinRepository) Save(ctx context.Context, coin *models.BoughtCoin) error {
	if m.SaveFunc != nil {
		return m.SaveFunc(ctx, coin)
	}
	return nil
}

func (m *mockBoughtCoinRepository) Delete(ctx context.Context, symbol string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, symbol)
	}
	return nil
}

func (m *mockBoughtCoinRepository) DeleteByID(ctx context.Context, id int64) error {
	if m.DeleteByIDFunc != nil {
		return m.DeleteByIDFunc(ctx, id)
	}
	return nil
}

func (m *mockBoughtCoinRepository) UpdatePrice(ctx context.Context, symbol string, price float64) error {
	if m.UpdatePriceFunc != nil {
		return m.UpdatePriceFunc(ctx, symbol, price)
	}
	return nil
}

func (m *mockBoughtCoinRepository) FindAllActive(ctx context.Context) ([]*models.BoughtCoin, error) {
	if m.FindAllActiveFunc != nil {
		return m.FindAllActiveFunc(ctx)
	}
	return []*models.BoughtCoin{}, nil
}

func (m *mockBoughtCoinRepository) HardDelete(ctx context.Context, symbol string) error {
	if m.HardDeleteFunc != nil {
		return m.HardDeleteFunc(ctx, symbol)
	}
	return nil
}

func (m *mockBoughtCoinRepository) Count(ctx context.Context) (int64, error) {
	if m.CountFunc != nil {
		return m.CountFunc(ctx)
	}
	return 0, nil
}

func (m *mockTradeService) CancelOrder(ctx context.Context, orderID string) error {
	return nil
}

func (m *mockTradeService) GetPendingOrders(ctx context.Context) ([]*models.Order, error) {
	return nil, nil
}

func (m *mockTradeService) GetOrderStatus(ctx context.Context, orderID string) (*models.Order, error) {
	return nil, nil
}

func (m *mockTradeService) ExecuteTrade(ctx context.Context, order *models.Order) (*models.Order, error) {
	return &models.Order{}, nil
}

func (m *mockTradeService) GetTradeHistory(ctx context.Context, startTime time.Time, limit int) ([]*models.Order, error) {
	return []*models.Order{}, nil
}

func (m *mockTradeService) GetActiveTrades(ctx context.Context) ([]*models.BoughtCoin, error) {
	return []*models.BoughtCoin{}, nil
}

func TestTradeHandler_ExecutePurchaseHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	validPayload := map[string]interface{}{
		"symbol":   "BTCUSDT",
		"price":    50000,
		"quantity": 0.1,
	}
	payloadBytes, _ := json.Marshal(validPayload)

	tests := []struct {
		name             string
		body             []byte
		mockExecPurchase func(ctx context.Context, symbol string, amount float64) (*models.BoughtCoin, error)
		expectedStatus   int
		expectError      string
	}{
		{
			name: "success",
			body: payloadBytes,
			mockExecPurchase: func(ctx context.Context, symbol string, amount float64) (*models.BoughtCoin, error) {
				return &models.BoughtCoin{
					Symbol:   symbol,
					BuyPrice: 50000,
					Quantity: amount,
				}, nil
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "invalid json",
			body: []byte(`{"symbol":""}`),
			mockExecPurchase: func(ctx context.Context, symbol string, amount float64) (*models.BoughtCoin, error) {
				return nil, nil
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    "INVALID_REQUEST",
		},
		{
			name: "service error",
			body: payloadBytes,
			mockExecPurchase: func(ctx context.Context, symbol string, amount float64) (*models.BoughtCoin, error) {
				return nil, errors.New("fail")
			},
			expectedStatus: http.StatusInternalServerError,
			expectError:    "TRADE_EXECUTION_FAILED",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockTradeService{
				ExecutePurchaseFunc: tt.mockExecPurchase,
			}
			mockRepo := &mockBoughtCoinRepository{}
			handler := &TradeHandler{
				TradeService:   mockSvc,
				BoughtCoinRepo: mockRepo,
			}

			router := gin.New()
			router.POST("/api/v1/trades", handler.ExecutePurchaseHandler)

			req := httptest.NewRequest(http.MethodPost, "/api/v1/trades", bytes.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectError != "" {
				assert.Contains(t, w.Body.String(), tt.expectError)
			}
		})
	}
}

/*
func TestTradeHandler_ListTrades(t *testing.T) {
	// Commented out because ListTrades handler is not implemented currently
}
*/

/*
func TestTradeHandler_CancelTrade(t *testing.T) {
	// Commented out because CancelTrade handler is not implemented currently
}
*/
