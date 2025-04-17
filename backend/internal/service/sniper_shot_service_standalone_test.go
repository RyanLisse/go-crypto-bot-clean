package service_test

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/service"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock implementations for the necessary interfaces

// MockTradeExecutor is a mock implementation for the TradeExecutor interface
type MockTradeExecutor struct {
	mock.Mock
}

func (m *MockTradeExecutor) ExecuteOrder(ctx context.Context, request *model.OrderRequest) (*model.OrderResponse, error) {
	args := m.Called(ctx, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.OrderResponse), args.Error(1)
}

func (m *MockTradeExecutor) CancelOrderWithRetry(ctx context.Context, symbol, orderID string) error {
	args := m.Called(ctx, symbol, orderID)
	return args.Error(0)
}

func (m *MockTradeExecutor) GetOrderStatusWithRetry(ctx context.Context, symbol, orderID string) (*model.Order, error) {
	args := m.Called(ctx, symbol, orderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Order), args.Error(1)
}

// MockOrderRepository is a mock implementation for the OrderRepository interface
type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) Create(ctx context.Context, order *model.Order) error {
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

func (m *MockOrderRepository) GetByClientOrderID(ctx context.Context, clientOrderID string) (*model.Order, error) {
	args := m.Called(ctx, clientOrderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Order), args.Error(1)
}

func (m *MockOrderRepository) Update(ctx context.Context, order *model.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func (m *MockOrderRepository) GetBySymbol(ctx context.Context, symbol string, limit, offset int) ([]*model.Order, error) {
	args := m.Called(ctx, symbol, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Order), args.Error(1)
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

func TestSniperShotService_ExecuteSniper_Standalone(t *testing.T) {
	// Create a logger for testing
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	// Create mock dependencies
	mockTradeExecutor := new(MockTradeExecutor)
	mockOrderRepo := new(MockOrderRepository)

	// Create the service
	sniperService := service.NewSniperShotService(&logger, mockTradeExecutor, mockOrderRepo)

	// Test data
	testUserID := "test-user-1"
	testSymbol := "BTCUSDT"
	testOrderID := "12345"

	t.Run("Successful Market Order Execution", func(t *testing.T) {
		// Expected order request
		expectedOrderRequest := &model.OrderRequest{
			UserID:      testUserID,
			Symbol:      testSymbol,
			Side:        model.OrderSideBuy,
			Type:        model.OrderTypeMarket,
			Quantity:    0.1,
			Price:       0,
			TimeInForce: model.TimeInForceGTC,
		}

		// Mock successful order response
		mockOrderResponse := &model.OrderResponse{
			Order: model.Order{
				ID:              "internal-id-1",
				OrderID:         testOrderID,
				ClientOrderID:   "client-12345",
				UserID:          testUserID,
				Symbol:          testSymbol,
				Side:            model.OrderSideBuy,
				Type:            model.OrderTypeMarket,
				Status:          model.OrderStatusFilled,
				Quantity:        0.1,
				ExecutedQty:     0.1,
				AvgFillPrice:    50000.0,
				Commission:      0.0001,
				CommissionAsset: "BTC",
				TimeInForce:     model.TimeInForceGTC,
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
				Exchange:        "binance",
			},
			IsSuccess: true,
		}

		// Set up the mock expectations
		mockTradeExecutor.On("ExecuteOrder", mock.Anything, mock.MatchedBy(func(req *model.OrderRequest) bool {
			return req.UserID == expectedOrderRequest.UserID &&
				req.Symbol == expectedOrderRequest.Symbol &&
				req.Side == expectedOrderRequest.Side &&
				req.Type == expectedOrderRequest.Type &&
				req.Quantity == expectedOrderRequest.Quantity &&
				req.TimeInForce == expectedOrderRequest.TimeInForce
		})).Return(mockOrderResponse, nil).Once()

		// Mock repository save
		mockOrderRepo.On("Create", mock.Anything, mock.MatchedBy(func(order *model.Order) bool {
			return order.OrderID == testOrderID &&
				order.UserID == testUserID &&
				order.Symbol == testSymbol
		})).Return(nil).Once()

		// Create the request
		request := &service.SniperShotRequest{
			UserID:    testUserID,
			Symbol:    testSymbol,
			Side:      model.OrderSideBuy,
			Type:      model.OrderTypeMarket,
			Quantity:  0.1,
			Price:     0,
			TimeLimit: 5 * time.Second,
		}

		// Execute the function being tested
		result, err := sniperService.ExecuteSniper(context.Background(), request)

		// Check expectations
		assert.NoError(t, err)
		assert.True(t, result.Success)
		assert.Nil(t, result.Error)
		assert.NotNil(t, result.Order)
		assert.Equal(t, testOrderID, result.Order.OrderID)
		assert.Equal(t, testUserID, result.Order.UserID)
		assert.Equal(t, testSymbol, result.Order.Symbol)
		assert.Equal(t, model.OrderStatusFilled, result.Order.Status)
		assert.Equal(t, 0.1, result.Order.ExecutedQty)
		assert.Greater(t, result.Latency, time.Duration(0))

		// Verify that all expectations were met
		mockTradeExecutor.AssertExpectations(t)
		mockOrderRepo.AssertExpectations(t)
	})

	t.Run("Failed Order Execution", func(t *testing.T) {
		// Set up the mock expectations for a failed execution
		mockTradeExecutor.On("ExecuteOrder", mock.Anything, mock.Anything).
			Return(nil, errors.New("insufficient balance")).Once()

		// Create the request
		request := &service.SniperShotRequest{
			UserID:    testUserID,
			Symbol:    testSymbol,
			Side:      model.OrderSideBuy,
			Type:      model.OrderTypeLimit,
			Quantity:  0.2,
			Price:     45000.0,
			TimeLimit: 5 * time.Second,
		}

		// Execute the function being tested
		result, err := sniperService.ExecuteSniper(context.Background(), request)

		// Check expectations
		assert.Error(t, err)
		assert.False(t, result.Success)
		assert.NotNil(t, result.Error)
		assert.Nil(t, result.Order)
		assert.Contains(t, result.Error.Error(), "insufficient balance")
		assert.Greater(t, result.Latency, time.Duration(0))

		// Verify that all expectations were met
		mockTradeExecutor.AssertExpectations(t)
	})
}

func TestSniperShotService_CancelSniper_Standalone(t *testing.T) {
	// Create a logger for testing
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	// Create mock dependencies
	mockTradeExecutor := new(MockTradeExecutor)
	mockOrderRepo := new(MockOrderRepository)

	// Create the service
	sniperService := service.NewSniperShotService(&logger, mockTradeExecutor, mockOrderRepo)

	// Test data
	testSymbol := "BTCUSDT"
	testOrderID := "order-to-cancel-123"

	t.Run("Successful order cancellation", func(t *testing.T) {
		// Set up mock expectations
		mockTradeExecutor.On("CancelOrderWithRetry", mock.Anything, testSymbol, testOrderID).
			Return(nil).Once()

		// Execute the function
		err := sniperService.CancelSniper(context.Background(), testSymbol, testOrderID)

		// Check expectations
		assert.NoError(t, err)
		mockTradeExecutor.AssertExpectations(t)
	})

	t.Run("Failed order cancellation", func(t *testing.T) {
		cancelError := errors.New("order not found or already filled")

		// Set up mock expectations
		mockTradeExecutor.On("CancelOrderWithRetry", mock.Anything, testSymbol, testOrderID).
			Return(cancelError).Once()

		// Execute the function
		err := sniperService.CancelSniper(context.Background(), testSymbol, testOrderID)

		// Check expectations
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to cancel order")
		mockTradeExecutor.AssertExpectations(t)
	})
}

func TestSniperShotService_GetSniperOrderStatus_Standalone(t *testing.T) {
	// Create a logger for testing
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	// Create mock dependencies
	mockTradeExecutor := new(MockTradeExecutor)
	mockOrderRepo := new(MockOrderRepository)

	// Create the service
	sniperService := service.NewSniperShotService(&logger, mockTradeExecutor, mockOrderRepo)

	// Test data
	testSymbol := "BTCUSDT"
	testOrderID := "status-check-order-123"

	t.Run("Successfully get order status", func(t *testing.T) {
		// Create a mock order response
		mockOrder := &model.Order{
			OrderID:     testOrderID,
			Symbol:      testSymbol,
			Status:      model.OrderStatusPartiallyFilled,
			ExecutedQty: 0.05,
			Quantity:    0.1,
			CreatedAt:   time.Now().Add(-1 * time.Hour),
			UpdatedAt:   time.Now(),
		}

		// Set up mock expectations
		mockTradeExecutor.On("GetOrderStatusWithRetry", mock.Anything, testSymbol, testOrderID).
			Return(mockOrder, nil).Once()

		// Execute the function
		order, err := sniperService.GetSniperOrderStatus(context.Background(), testSymbol, testOrderID)

		// Check expectations
		assert.NoError(t, err)
		assert.NotNil(t, order)
		assert.Equal(t, testOrderID, order.OrderID)
		assert.Equal(t, model.OrderStatusPartiallyFilled, order.Status)
		assert.Equal(t, 0.05, order.ExecutedQty)
		mockTradeExecutor.AssertExpectations(t)
	})

	t.Run("Failed to get order status", func(t *testing.T) {
		statusError := errors.New("exchange API error")

		// Set up mock expectations
		mockTradeExecutor.On("GetOrderStatusWithRetry", mock.Anything, testSymbol, testOrderID).
			Return(nil, statusError).Once()

		// Execute the function
		order, err := sniperService.GetSniperOrderStatus(context.Background(), testSymbol, testOrderID)

		// Check expectations
		assert.Error(t, err)
		assert.Nil(t, order)
		assert.Contains(t, err.Error(), "failed to get order status")
		mockTradeExecutor.AssertExpectations(t)
	})
}
