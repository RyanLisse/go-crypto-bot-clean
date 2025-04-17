package service

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	mockrepo "github.com/RyanLisse/go-crypto-bot-clean/backend/internal/mocks/repository"
	mocksvc "github.com/RyanLisse/go-crypto-bot-clean/backend/internal/mocks/service"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTradeExecutor is a mock for the TradeExecutor interface
type MockTradeExecutor struct {
	mock.Mock
}

// ExecuteOrder mocks the ExecuteOrder method
func (m *MockTradeExecutor) ExecuteOrder(ctx context.Context, order *model.OrderRequest) (*model.OrderResponse, error) {
	args := m.Called(ctx, order)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.OrderResponse), args.Error(1)
}

// CancelOrderWithRetry mocks the CancelOrderWithRetry method
func (m *MockTradeExecutor) CancelOrderWithRetry(ctx context.Context, symbol, orderID string) error {
	args := m.Called(ctx, symbol, orderID)
	return args.Error(0)
}

// GetOrderStatusWithRetry mocks the GetOrderStatusWithRetry method
func (m *MockTradeExecutor) GetOrderStatusWithRetry(ctx context.Context, symbol, orderID string) (*model.Order, error) {
	args := m.Called(ctx, symbol, orderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Order), args.Error(1)
}

// MockOrderRepository is a mock for the OrderRepository interface
type MockOrderRepository struct {
	mock.Mock
}

// Save mocks the Save method
func (m *MockOrderRepository) Save(ctx context.Context, order *model.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

// Get mocks the Get method
func (m *MockOrderRepository) Get(ctx context.Context, id string) (*model.Order, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Order), args.Error(1)
}

// GetByOrderID mocks the GetByOrderID method
func (m *MockOrderRepository) GetByOrderID(ctx context.Context, orderID string) (*model.Order, error) {
	args := m.Called(ctx, orderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Order), args.Error(1)
}

// List mocks the List method
func (m *MockOrderRepository) List(ctx context.Context, userID string, limit, offset int) ([]*model.Order, error) {
	args := m.Called(ctx, userID, limit, offset)
	return args.Get(0).([]*model.Order), args.Error(1)
}

// Create mocks the Create method
func (m *MockOrderRepository) Create(ctx context.Context, order *model.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

// MockPriceChecker is a mock implementation of the PriceChecker interface
type MockPriceChecker struct {
	mock.Mock
}

// GetCurrentPrice mocks the GetCurrentPrice method
func (m *MockPriceChecker) GetCurrentPrice(ctx context.Context, symbol string) (float64, error) {
	args := m.Called(ctx, symbol)
	return args.Get(0).(float64), args.Error(1)
}

func TestSniperShotService_ExecuteSniper(t *testing.T) {
	// Create a logger for testing
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	// Create mock dependencies
	mockTradeExecutor := new(mocksvc.MockTradeExecutor)
	mockOrderRepo := new(mockrepo.MockOrderRepository)
	mockPriceChecker := new(MockPriceChecker)

	// Create the service
	service := NewSniperShotService(&logger, mockTradeExecutor, mockOrderRepo, mockPriceChecker)

	// Test data
	testUserID := "test-user-1"
	testSymbol := "BTCUSDT"
	testOrderID := "12345"

	// Create test cases
	tests := []struct {
		name          string
		request       *SniperShotRequest
		mockSetup     func()
		expectedError bool
		checkResult   func(*SniperShotResult)
	}{
		{
			name: "Successful Market Order Execution",
			request: &SniperShotRequest{
				UserID:    testUserID,
				Symbol:    testSymbol,
				Side:      model.OrderSideBuy,
				Type:      model.OrderTypeMarket,
				Quantity:  0.1,
				Price:     0,
				TimeLimit: 5 * time.Second,
			},
			mockSetup: func() {
				// Expected order request
				expectedOrderReq := &model.OrderRequest{
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
					return req.UserID == expectedOrderReq.UserID &&
						req.Symbol == expectedOrderReq.Symbol &&
						req.Side == expectedOrderReq.Side &&
						req.Type == expectedOrderReq.Type &&
						req.Quantity == expectedOrderReq.Quantity &&
						req.TimeInForce == expectedOrderReq.TimeInForce
				})).Return(mockOrderResponse, nil).Once()

				// Mock repository save
				mockOrderRepo.On("Create", mock.Anything, mock.MatchedBy(func(order *model.Order) bool {
					return order.OrderID == testOrderID &&
						order.UserID == testUserID &&
						order.Symbol == testSymbol
				})).Return(nil).Once()
			},
			expectedError: false,
			checkResult: func(result *SniperShotResult) {
				assert.True(t, result.Success)
				assert.Nil(t, result.Error)
				assert.NotNil(t, result.Order)
				assert.Equal(t, testOrderID, result.Order.OrderID)
				assert.Equal(t, testUserID, result.Order.UserID)
				assert.Equal(t, testSymbol, result.Order.Symbol)
				assert.Equal(t, model.OrderStatusFilled, result.Order.Status)
				assert.Equal(t, 0.1, result.Order.ExecutedQty)
				assert.Greater(t, result.Latency, time.Duration(0))
			},
		},
		{
			name: "Failed Order Execution",
			request: &SniperShotRequest{
				UserID:    testUserID,
				Symbol:    testSymbol,
				Side:      model.OrderSideBuy,
				Type:      model.OrderTypeLimit,
				Quantity:  0.2,
				Price:     45000.0,
				TimeLimit: 5 * time.Second,
			},
			mockSetup: func() {
				// Set up the mock expectations for a failed execution
				mockTradeExecutor.On("ExecuteOrder", mock.Anything, mock.Anything).
					Return(nil, errors.New("insufficient balance")).Once()
			},
			expectedError: true,
			checkResult: func(result *SniperShotResult) {
				assert.False(t, result.Success)
				assert.NotNil(t, result.Error)
				assert.Nil(t, result.Order)
				assert.Contains(t, result.Error.Error(), "insufficient balance")
				assert.Greater(t, result.Latency, time.Duration(0))
			},
		},
		{
			name: "Repository Save Error (but trade successful)",
			request: &SniperShotRequest{
				UserID:    testUserID,
				Symbol:    testSymbol,
				Side:      model.OrderSideSell,
				Type:      model.OrderTypeLimit,
				Quantity:  0.05,
				Price:     52000.0,
				TimeLimit: 5 * time.Second,
			},
			mockSetup: func() {
				// Mock successful order response
				mockOrderResponse := &model.OrderResponse{
					Order: model.Order{
						ID:            "internal-id-2",
						OrderID:       "sell-order-1",
						ClientOrderID: "client-67890",
						UserID:        testUserID,
						Symbol:        testSymbol,
						Side:          model.OrderSideSell,
						Type:          model.OrderTypeLimit,
						Status:        model.OrderStatusNew,
						Quantity:      0.05,
						Price:         52000.0,
						ExecutedQty:   0,
						TimeInForce:   model.TimeInForceGTC,
						CreatedAt:     time.Now(),
						UpdatedAt:     time.Now(),
						Exchange:      "binance",
					},
					IsSuccess: true,
				}

				// Set up the mock expectations
				mockTradeExecutor.On("ExecuteOrder", mock.Anything, mock.Anything).
					Return(mockOrderResponse, nil).Once()

				// Mock repository error
				mockOrderRepo.On("Create", mock.Anything, mock.Anything).
					Return(errors.New("database connection error")).Once()
			},
			expectedError: false, // Trade is still successful even if DB save fails
			checkResult: func(result *SniperShotResult) {
				assert.True(t, result.Success)
				assert.Nil(t, result.Error)
				assert.NotNil(t, result.Order)
				assert.Equal(t, "sell-order-1", result.Order.OrderID)
				assert.Equal(t, model.OrderSideSell, result.Order.Side)
				assert.Equal(t, model.OrderStatusNew, result.Order.Status)
				assert.Greater(t, result.Latency, time.Duration(0))
			},
		},
	}

	// Run test cases
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Set up the mocks for this test case
			tc.mockSetup()

			// Execute the function being tested
			result, err := service.ExecuteSniper(context.Background(), tc.request)

			// Check error expectations
			if tc.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			// Check the result
			tc.checkResult(result)

			// Verify that all expectations were met
			mockTradeExecutor.AssertExpectations(t)
			mockOrderRepo.AssertExpectations(t)
		})
	}
}

func TestSniperShotService_CancelSniper(t *testing.T) {
	// Create a logger for testing
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	// Create mock dependencies
	mockTradeExecutor := new(mocksvc.MockTradeExecutor)
	mockOrderRepo := new(mockrepo.MockOrderRepository)

	// Create the service
	service := NewSniperShotService(&logger, mockTradeExecutor, mockOrderRepo)

	// Test data
	testSymbol := "BTCUSDT"
	testOrderID := "order-to-cancel-123"

	t.Run("Successful order cancellation", func(t *testing.T) {
		// Set up mock expectations
		mockTradeExecutor.On("CancelOrderWithRetry", mock.Anything, testSymbol, testOrderID).
			Return(nil).Once()

		// Execute the function
		err := service.CancelSniper(context.Background(), testSymbol, testOrderID)

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
		err := service.CancelSniper(context.Background(), testSymbol, testOrderID)

		// Check expectations
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to cancel order")
		mockTradeExecutor.AssertExpectations(t)
	})
}

func TestSniperShotService_GetSniperOrderStatus(t *testing.T) {
	// Create a logger for testing
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	// Create mock dependencies
	mockTradeExecutor := new(mocksvc.MockTradeExecutor)
	mockOrderRepo := new(mockrepo.MockOrderRepository)

	// Create the service
	service := NewSniperShotService(&logger, mockTradeExecutor, mockOrderRepo)

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
		order, err := service.GetSniperOrderStatus(context.Background(), testSymbol, testOrderID)

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
		order, err := service.GetSniperOrderStatus(context.Background(), testSymbol, testOrderID)

		// Check expectations
		assert.Error(t, err)
		assert.Nil(t, order)
		assert.Contains(t, err.Error(), "failed to get order status")
		mockTradeExecutor.AssertExpectations(t)
	})
}

func TestExecuteSniperWithCallbacks(t *testing.T) {
	// Create mocks
	mockTradeExecutor := new(MockTradeExecutor)
	mockOrderRepo := new(MockOrderRepository)

	// Override getCurrentPrice for testing
	origGetCurrentPrice := getCurrentPrice
	defer func() { getCurrentPrice = origGetCurrentPrice }()

	currentPrice := 49600.0
	getCurrentPrice = func(ctx context.Context, symbol string) (float64, error) {
		return currentPrice, nil
	}

	// Logger setup
	logger := zerolog.New(os.Stdout)

	// Create the service
	service := NewSniperShotService(&logger, mockTradeExecutor, mockOrderRepo)

	// Test data
	userID := "test-user"
	symbol := "BTCUSDT"
	side := model.OrderSideBuy
	quantity := 1.0
	price := 50000.0
	targetPrice := 49500.0 // Less than currentPrice to trigger condition

	// Setup the order response
	orderResponse := &model.OrderResponse{
		Order: model.Order{
			ID:           "test-id",
			OrderID:      "test-order-id",
			UserID:       userID,
			Symbol:       symbol,
			Side:         side,
			Type:         model.OrderTypeLimit,
			Status:       model.OrderStatusFilled,
			Price:        price,
			Quantity:     quantity,
			ExecutedQty:  quantity,
			AvgFillPrice: price,
			Exchange:     "test-exchange",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		IsSuccess: true,
	}

	// Create a callback channel to verify it was called
	callbackCalled := make(chan float64, 1)
	callback := func(price float64) {
		callbackCalled <- price
	}

	// Create a request with a trigger condition and callback
	req := &SniperShotRequest{
		UserID:    userID,
		Symbol:    symbol,
		Side:      side,
		Quantity:  quantity,
		Price:     price,
		Type:      model.OrderTypeLimit,
		TimeLimit: 5 * time.Second,
		Condition: &TriggerCondition{
			TargetPrice:     targetPrice,
			Operator:        ">=", // currentPrice >= targetPrice will be true
			MaxTimeoutSecs:  2,
			PriceBufferPct:  0.01,
			CheckIntervalMs: 100,
			Callbacks:       []func(price float64){callback},
		},
	}

	// Mock the expected order execution
	mockTradeExecutor.On("ExecuteOrder", mock.Anything, mock.MatchedBy(func(order *model.OrderRequest) bool {
		// For buy orders with a buffer
		bufferAmount := targetPrice * 0.01
		expectedPrice := targetPrice + bufferAmount

		return order.UserID == userID &&
			order.Symbol == symbol &&
			order.Side == side &&
			order.Type == model.OrderTypeLimit &&
			order.Quantity == quantity &&
			order.Price == expectedPrice
	})).Return(orderResponse, nil)

	// Mock the order repository
	mockOrderRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

	// Execute the sniper shot
	result, err := service.ExecuteSniper(context.Background(), req)

	// Verify results
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.Success)
	assert.Equal(t, "test-order-id", result.Order.OrderID)

	// Verify callback was called with the current price
	select {
	case receivedPrice := <-callbackCalled:
		assert.Equal(t, currentPrice, receivedPrice)
	case <-time.After(1 * time.Second):
		t.Fatal("Callback was not called within timeout")
	}

	// Verify expectations
	mockTradeExecutor.AssertExpectations(t)
	mockOrderRepo.AssertExpectations(t)
}
