package usecase

import (
	"context"
	"testing"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model/market"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRiskUseCase is a mock implementation of RiskUseCase for testing
type MockRiskUseCase struct {
	mock.Mock
}

func (m *MockRiskUseCase) EvaluateOrderRisk(ctx context.Context, userID string, orderRequest model.OrderRequest) (bool, []*model.RiskAssessment, error) {
	args := m.Called(ctx, userID, orderRequest)
	return args.Bool(0), args.Get(1).([]*model.RiskAssessment), args.Error(2)
}

func (m *MockRiskUseCase) EvaluatePositionRisk(ctx context.Context, userID string, positionID string) ([]*model.RiskAssessment, error) {
	args := m.Called(ctx, userID, positionID)
	return args.Get(0).([]*model.RiskAssessment), args.Error(1)
}

func (m *MockRiskUseCase) EvaluatePortfolioRisk(ctx context.Context, userID string) ([]*model.RiskAssessment, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*model.RiskAssessment), args.Error(1)
}

func (m *MockRiskUseCase) GetRiskMetrics(ctx context.Context, userID string) (*model.RiskMetrics, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(*model.RiskMetrics), args.Error(1)
}

func (m *MockRiskUseCase) GetHistoricalRiskMetrics(ctx context.Context, userID string, days int) ([]*model.RiskMetrics, error) {
	args := m.Called(ctx, userID, days)
	return args.Get(0).([]*model.RiskMetrics), args.Error(1)
}

func (m *MockRiskUseCase) GetActiveRisks(ctx context.Context, userID string) ([]*model.RiskAssessment, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*model.RiskAssessment), args.Error(1)
}

func (m *MockRiskUseCase) GetRiskAssessments(ctx context.Context, userID string, riskType *model.RiskType, level *model.RiskLevel, limit, offset int) ([]*model.RiskAssessment, error) {
	args := m.Called(ctx, userID, riskType, level, limit, offset)
	return args.Get(0).([]*model.RiskAssessment), args.Error(1)
}

func (m *MockRiskUseCase) ResolveRisk(ctx context.Context, riskID string) error {
	args := m.Called(ctx, riskID)
	return args.Error(0)
}

func (m *MockRiskUseCase) IgnoreRisk(ctx context.Context, riskID string) error {
	args := m.Called(ctx, riskID)
	return args.Error(0)
}

func (m *MockRiskUseCase) GetRiskProfile(ctx context.Context, userID string) (*model.RiskProfile, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(*model.RiskProfile), args.Error(1)
}

func (m *MockRiskUseCase) UpdateRiskProfile(ctx context.Context, profile *model.RiskProfile) error {
	args := m.Called(ctx, profile)
	return args.Error(0)
}

func (m *MockRiskUseCase) SaveRiskConstraint(ctx context.Context, constraint *model.RiskConstraint) error {
	args := m.Called(ctx, constraint)
	return args.Error(0)
}

func (m *MockRiskUseCase) DeleteRiskConstraint(ctx context.Context, constraintID string) error {
	args := m.Called(ctx, constraintID)
	return args.Error(0)
}

func (m *MockRiskUseCase) GetActiveConstraints(ctx context.Context, userID string) ([]*model.RiskConstraint, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*model.RiskConstraint), args.Error(1)
}

// MockTradeService is a mock implementation of TradeService for testing
type MockTradeService struct {
	mock.Mock
}

func (m *MockTradeService) PlaceOrder(ctx context.Context, req *model.OrderRequest) (*model.PlaceOrderResponse, error) {
	args := m.Called(ctx, req)
	var resp *model.PlaceOrderResponse
	if arg0 := args.Get(0); arg0 != nil {
		resp = arg0.(*model.PlaceOrderResponse)
	}
	return resp, args.Error(1)
}

func (m *MockTradeService) CancelOrder(ctx context.Context, symbol, orderID string) error {
	args := m.Called(ctx, symbol, orderID)
	return args.Error(0)
}

func (m *MockTradeService) GetOrderStatus(ctx context.Context, symbol, orderID string) (*model.Order, error) {
	args := m.Called(ctx, symbol, orderID)
	return args.Get(0).(*model.Order), args.Error(1)
}

func (m *MockTradeService) GetOpenOrders(ctx context.Context, symbol string) ([]*model.Order, error) {
	args := m.Called(ctx, symbol)
	return args.Get(0).([]*model.Order), args.Error(1)
}

func (m *MockTradeService) GetOrderHistory(ctx context.Context, symbol string, limit, offset int) ([]*model.Order, error) {
	args := m.Called(ctx, symbol, limit, offset)
	return args.Get(0).([]*model.Order), args.Error(1)
}

func (m *MockTradeService) CalculateRequiredQuantity(ctx context.Context, symbol string, side model.OrderSide, amount float64) (float64, error) {
	args := m.Called(ctx, symbol, side, amount)
	val := args.Get(0)
	if val == nil {
		return 0.0, args.Error(1)
	}
	return val.(float64), args.Error(1)
}

// Minimal mocks for required dependencies

type MockMEXCClient struct{ mock.Mock }

func (m *MockMEXCClient) CancelOrder(ctx context.Context, symbol, orderID string) error { return nil }
func (m *MockMEXCClient) GetAccount(ctx context.Context) (*model.Wallet, error)         { return nil, nil }
func (m *MockMEXCClient) GetNewListings(ctx context.Context) ([]*model.NewCoin, error) {
	return nil, nil
}
func (m *MockMEXCClient) GetSymbolInfo(ctx context.Context, symbol string) (*model.SymbolInfo, error) {
	return nil, nil
}
func (m *MockMEXCClient) GetSymbolStatus(ctx context.Context, symbol string) (model.Status, error) {
	return "", nil
}
func (m *MockMEXCClient) GetTradingSchedule(ctx context.Context, symbol string) (model.TradingSchedule, error) {
	return model.TradingSchedule{}, nil
}
func (m *MockMEXCClient) GetSymbolConstraints(ctx context.Context, symbol string) (*model.SymbolConstraints, error) {
	return &model.SymbolConstraints{
		MinPrice:   0.00000001,
		MaxPrice:   100000.0,
		MinQty:     0.00000001,
		MaxQty:     10000.0,
		PriceScale: 8,
		QtyScale:   8,
	}, nil
}
func (m *MockMEXCClient) GetExchangeInfo(ctx context.Context) (*model.ExchangeInfo, error) {
	return nil, nil
}
func (m *MockMEXCClient) GetMarketData(ctx context.Context, symbol string) (*model.Ticker, error) {
	return nil, nil
}
func (m *MockMEXCClient) GetKlines(ctx context.Context, symbol string, interval model.KlineInterval, limit int) ([]*model.Kline, error) {
	return nil, nil
}
func (m *MockMEXCClient) GetOrderBook(ctx context.Context, symbol string, depth int) (*model.OrderBook, error) {
	return nil, nil
}
func (m *MockMEXCClient) PlaceOrder(ctx context.Context, symbol string, side model.OrderSide, orderType model.OrderType, quantity float64, price float64, timeInForce model.TimeInForce) (*model.Order, error) {
	return nil, nil
}
func (m *MockMEXCClient) GetOrderStatus(ctx context.Context, symbol string, orderID string) (*model.Order, error) {
	return nil, nil
}
func (m *MockMEXCClient) GetOpenOrders(ctx context.Context, symbol string) ([]*model.Order, error) {
	return nil, nil
}
func (m *MockMEXCClient) GetOrderHistory(ctx context.Context, symbol string, limit, offset int) ([]*model.Order, error) {
	return nil, nil
}

// OrderRepository

type MockOrderRepository struct{ mock.Mock }

func (m *MockOrderRepository) Create(ctx context.Context, order *model.Order) error { return nil }
func (m *MockOrderRepository) GetByID(ctx context.Context, id string) (*model.Order, error) {
	return nil, nil
}
func (m *MockOrderRepository) GetByClientOrderID(ctx context.Context, clientOrderID string) (*model.Order, error) {
	return nil, nil
}
func (m *MockOrderRepository) Update(ctx context.Context, order *model.Order) error { return nil }
func (m *MockOrderRepository) GetBySymbol(ctx context.Context, symbol string, limit, offset int) ([]*model.Order, error) {
	return nil, nil
}
func (m *MockOrderRepository) GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*model.Order, error) {
	return nil, nil
}
func (m *MockOrderRepository) GetByStatus(ctx context.Context, status model.OrderStatus, limit, offset int) ([]*model.Order, error) {
	return nil, nil
}
func (m *MockOrderRepository) Count(ctx context.Context, filters map[string]interface{}) (int64, error) {
	return 0, nil
}
func (m *MockOrderRepository) Delete(ctx context.Context, id string) error { return nil }

// SymbolRepository (use *market.Symbol)
type MockSymbolRepository struct{ mock.Mock }

func (m *MockSymbolRepository) Create(ctx context.Context, symbol *market.Symbol) error { return nil }
func (m *MockSymbolRepository) GetBySymbol(ctx context.Context, symbol string) (*market.Symbol, error) {
	return &market.Symbol{Symbol: "BTCUSDT"}, nil
}
func (m *MockSymbolRepository) GetByExchange(ctx context.Context, exchange string) ([]*market.Symbol, error) {
	return nil, nil
}
func (m *MockSymbolRepository) GetAll(ctx context.Context) ([]*market.Symbol, error)    { return nil, nil }
func (m *MockSymbolRepository) Update(ctx context.Context, symbol *market.Symbol) error { return nil }
func (m *MockSymbolRepository) Delete(ctx context.Context, symbol string) error         { return nil }

type MockTransactionManager struct {
	mock.Mock
	ShouldCallFn bool
}

func (m *MockTransactionManager) WithTransaction(ctx context.Context, fn func(context.Context) error) error {
	if m.ShouldCallFn && fn != nil {
		_ = fn(ctx)
	}
	return nil
}

func TestTradeRiskIntegration(t *testing.T) {
	// Create testify mocks for all dependencies

	mockSymbolRepo := new(MockSymbolRepository)
	mockTxManager := &MockTransactionManager{ShouldCallFn: true}

	// Always return a valid symbol for any GetBySymbol call
	mockSymbolRepo.On("GetBySymbol", mock.Anything, mock.Anything).Return(&market.Symbol{Symbol: "BTCUSDT"}, nil)

	// Test case 1: Order allowed by risk assessment
	t.Run("Order allowed by risk assessment", func(t *testing.T) {
		mockTradeService := new(MockTradeService)
		mockRiskUC := new(MockRiskUseCase)
		tradeUC := NewTradeUseCase(
			new(MockMEXCClient),
			new(MockOrderRepository),
			new(MockSymbolRepository),
			mockTradeService,
			mockRiskUC,
			mockTxManager,
			zerolog.Logger{},
		)
		ctx := context.Background()
		orderReq := model.OrderRequest{
			UserID:   "user123",
			Symbol:   "BTCUSDT",
			Side:     model.OrderSideBuy,
			Type:     model.OrderTypeLimit,
			Quantity: 0.1,
			Price:    50000,
		}

		mockRiskUC.On("EvaluateOrderRisk", ctx, "user123", orderReq).
			Return(true, []*model.RiskAssessment{}, nil).Once()

		orderResponse := &model.PlaceOrderResponse{
			Order: model.Order{
				OrderID: "order123",
				Symbol:  "BTCUSDT",
				Side:    model.OrderSideBuy,
				Type:    model.OrderTypeLimit,
				Status:  model.OrderStatusNew,
			},
			IsSuccess: true,
		}
		mockTradeService.On("PlaceOrder", ctx, &orderReq).Return(orderResponse, nil).Once()

		order, err := tradeUC.PlaceOrder(ctx, orderReq)

		assert.NoError(t, err)
		assert.NotNil(t, order)
		assert.Equal(t, "order123", order.OrderID)
		mockRiskUC.AssertExpectations(t)
		mockTradeService.AssertExpectations(t)
	})

	// Test case 2: Order rejected by risk assessment
	t.Run("Order rejected by risk assessment", func(t *testing.T) {
		// Prevent transaction function from being called if risk fails
		mockTxManager.ShouldCallFn = false
		mockTradeService := new(MockTradeService)
		mockRiskUC := new(MockRiskUseCase)
		tradeUC := NewTradeUseCase(
			new(MockMEXCClient),
			new(MockOrderRepository),
			new(MockSymbolRepository),
			mockTradeService,
			mockRiskUC,
			mockTxManager,
			zerolog.Logger{},
		)
		ctx := context.Background()
		orderReq := model.OrderRequest{
			UserID:   "user123",
			Symbol:   "BTCUSDT",
			Side:     model.OrderSideBuy,
			Type:     model.OrderTypeLimit,
			Quantity: 10.0, // Large quantity that exceeds risk limits
			Price:    50000,
		}

		riskAssessment := model.NewRiskAssessment(
			"user123",
			model.RiskTypePosition,
			model.RiskLevelHigh,
			"Order value exceeds maximum position size",
		)
		riskAssessment.Recommendation = "Reduce order size"

		mockRiskUC.On("EvaluateOrderRisk", mock.Anything, mock.Anything, mock.Anything).
			Return(false, []*model.RiskAssessment{riskAssessment}, nil).Once()

		order, err := tradeUC.PlaceOrder(ctx, orderReq)

		assert.Error(t, err)
		assert.Nil(t, order)
		assert.Contains(t, err.Error(), "Order value exceeds maximum position size")
		mockRiskUC.AssertExpectations(t)
		mockTradeService.AssertNotCalled(t, "PlaceOrder", mock.Anything, mock.Anything)
	})
}
