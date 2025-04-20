package usecase_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"    // Assuming mocks
	usecase "github.com/RyanLisse/go-crypto-bot-clean/backend/internal/usecase" // Import package under test
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockPriceChecker mocks the PriceChecker interface for testing
type MockPriceChecker struct {
	mock.Mock
}

func (m *MockPriceChecker) GetCurrentPrice(ctx context.Context, symbol string) (float64, error) {
	args := m.Called(ctx, symbol)
	return args.Get(0).(float64), args.Error(1)
}

// MockTradeExecutor mocks the TradeExecutor interface for testing
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

// MockOrderRepository mocks the OrderRepository interface for testing
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

func TestSniperShotServiceWithTriggerCondition(t *testing.T) {
	// Create a logger for testing
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	// Create mock dependencies
	mockTradeExecutor := new(MockTradeExecutor)
	mockOrderRepo := new(MockOrderRepository)
	mockPriceChecker := new(MockPriceChecker)

	// Create the service
	sniperService := usecase.NewSniperShotService(&logger, mockTradeExecutor, mockOrderRepo, mockPriceChecker)

	// Test data
	testUserID := "test-user-1"
	testSymbol := "BTCUSDT"
	testOrderID := "12345"

	t.Run("Execute order when price goes above target", func(t *testing.T) {
		// Mock price checker to return a price above target after the first call
		prices := []float64{49000.0, 51000.0}
		priceIndex := 0
		mockPriceChecker.On("GetCurrentPrice", mock.Anything, testSymbol).Return(
			func(ctx context.Context, symbol string) float64 {
				price := prices[priceIndex]
				if priceIndex < len(prices)-1 {
					priceIndex++
				}
				return price
			},
			func(ctx context.Context, symbol string) error {
				return nil
			},
		).Times(2)

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
				AvgFillPrice:    51000.0,
				Commission:      0.0001,
				CommissionAsset: "BTC",
				TimeInForce:     model.TimeInForceGTC,
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
				Exchange:        "binance",
			},
			IsSuccess: true,
		}

		// Set up the mock expectations for trade execution
		mockTradeExecutor.On("ExecuteOrder", mock.Anything, mock.Anything).Return(mockOrderResponse, nil).Once()

		// Mock repository save
		mockOrderRepo.On("Create", mock.Anything, mock.Anything).Return(nil).Once()

		// Create a trigger condition that activates when price goes above 50000
		triggerCondition := usecase.NewTriggerCondition(50000.0, ">").
			WithTimeout(5).
			WithCheckInterval(10) // Fast interval for testing

		// Create the request with trigger condition
		request := &usecase.SniperShotRequest{
			UserID:    testUserID,
			Symbol:    testSymbol,
			Side:      model.OrderSideBuy,
			Type:      model.OrderTypeMarket,
			Quantity:  0.1,
			Price:     0,
			TimeLimit: 5 * time.Second,
			Condition: &triggerCondition,
		}

		// Execute the function being tested
		result, err := sniperService.ExecuteSniper(context.Background(), request)

		// Check expectations
		assert.NoError(t, err)
		assert.True(t, result.Success)
		assert.Nil(t, result.Error)
		assert.NotNil(t, result.Order)
		assert.Equal(t, testOrderID, result.Order.OrderID)
		assert.Greater(t, result.Latency, time.Duration(0))

		// Verify all mocks were called as expected
		mockPriceChecker.AssertExpectations(t)
		mockTradeExecutor.AssertExpectations(t)
		mockOrderRepo.AssertExpectations(t)
	})

	t.Run("Timeout when price doesn't reach target", func(t *testing.T) {
		// Reset mock
		mockPriceChecker = new(MockPriceChecker)
		sniperService = usecase.NewSniperShotService(&logger, mockTradeExecutor, mockOrderRepo, mockPriceChecker)

		// Mock price checker to consistently return a price below target
		mockPriceChecker.On("GetCurrentPrice", mock.Anything, testSymbol).Return(49000.0, nil).Times(10)

		// Create a trigger condition that activates when price goes above 50000
		// with a very short timeout for test purposes
		triggerCondition := usecase.NewTriggerCondition(50000.0, ">").
			WithTimeout(1).       // 1 second timeout
			WithCheckInterval(10) // Fast interval for testing

		// Create the request with trigger condition
		request := &usecase.SniperShotRequest{
			UserID:    testUserID,
			Symbol:    testSymbol,
			Side:      model.OrderSideBuy,
			Type:      model.OrderTypeMarket,
			Quantity:  0.1,
			Price:     0,
			TimeLimit: 5 * time.Second,
			Condition: &triggerCondition,
		}

		// Execute the function being tested
		result, err := sniperService.ExecuteSniper(context.Background(), request)

		// Check expectations
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "timed out waiting for trigger condition")
		assert.False(t, result.Success)
		assert.NotNil(t, result.Error)
		assert.Nil(t, result.Order)

		// Verify all mocks were called as expected
		mockPriceChecker.AssertExpectations(t)
	})

	t.Run("Condition with price buffer", func(t *testing.T) {
		// Reset mock
		mockPriceChecker = new(MockPriceChecker)
		sniperService = usecase.NewSniperShotService(&logger, mockTradeExecutor, mockOrderRepo, mockPriceChecker)

		// Mock price checker to return a price close to but not exactly at target
		// The buffer should allow this to pass
		mockPriceChecker.On("GetCurrentPrice", mock.Anything, testSymbol).Return(49800.0, nil).Once()

		// Mock successful order response
		mockOrderResponse := &model.OrderResponse{
			Order: model.Order{
				ID:            "internal-id-2",
				OrderID:       testOrderID,
				ClientOrderID: "client-54321",
				UserID:        testUserID,
				Symbol:        testSymbol,
				Side:          model.OrderSideBuy,
				Type:          model.OrderTypeMarket,
				Status:        model.OrderStatusFilled,
				Quantity:      0.1,
				ExecutedQty:   0.1,
				AvgFillPrice:  49800.0,
				TimeInForce:   model.TimeInForceGTC,
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
				Exchange:      "binance",
			},
			IsSuccess: true,
		}

		// Set up the mock expectations for trade execution
		mockTradeExecutor.On("ExecuteOrder", mock.Anything, mock.Anything).Return(mockOrderResponse, nil).Once()

		// Mock repository save
		mockOrderRepo.On("Create", mock.Anything, mock.Anything).Return(nil).Once()

		// Create a trigger condition that activates when price is within 0.5% of 50000
		triggerCondition := usecase.NewTriggerCondition(50000.0, ">=").
			WithTimeout(5).
			WithPriceBuffer(0.005). // 0.5% buffer
			WithCheckInterval(10)   // Fast interval for testing

		// Create the request with trigger condition
		request := &usecase.SniperShotRequest{
			UserID:    testUserID,
			Symbol:    testSymbol,
			Side:      model.OrderSideBuy,
			Type:      model.OrderTypeMarket,
			Quantity:  0.1,
			Price:     0,
			TimeLimit: 5 * time.Second,
			Condition: &triggerCondition,
		}

		// Execute the function being tested
		result, err := sniperService.ExecuteSniper(context.Background(), request)

		// Check expectations
		assert.NoError(t, err)
		assert.True(t, result.Success)
		assert.Nil(t, result.Error)
		assert.NotNil(t, result.Order)
		assert.Equal(t, testOrderID, result.Order.OrderID)

		// Verify all mocks were called as expected
		mockPriceChecker.AssertExpectations(t)
		mockTradeExecutor.AssertExpectations(t)
		mockOrderRepo.AssertExpectations(t)
	})

	t.Run("Trigger condition with callback", func(t *testing.T) {
		// Reset mock
		mockPriceChecker = new(MockPriceChecker)
		sniperService = usecase.NewSniperShotService(&logger, mockTradeExecutor, mockOrderRepo, mockPriceChecker)

		// Mock price checker to return a price above the target
		mockPriceChecker.On("GetCurrentPrice", mock.Anything, testSymbol).Return(51000.0, nil).Once()

		// Mock successful order response
		mockOrderResponse := &model.OrderResponse{
			Order: model.Order{
				ID:            "internal-id-3",
				OrderID:       testOrderID,
				ClientOrderID: "client-callback",
				UserID:        testUserID,
				Symbol:        testSymbol,
				Side:          model.OrderSideBuy,
				Type:          model.OrderTypeMarket,
				Status:        model.OrderStatusFilled,
				Quantity:      0.1,
				ExecutedQty:   0.1,
				AvgFillPrice:  51000.0,
				TimeInForce:   model.TimeInForceGTC,
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
				Exchange:      "binance",
			},
			IsSuccess: true,
		}

		// Set up the mock expectations for trade execution
		mockTradeExecutor.On("ExecuteOrder", mock.Anything, mock.Anything).Return(mockOrderResponse, nil).Once()

		// Mock repository save
		mockOrderRepo.On("Create", mock.Anything, mock.Anything).Return(nil).Once()

		// Track callback execution
		callbackExecuted := false
		callbackPrice := 0.0

		// Create a trigger condition that activates when price is above 50000
		// and includes a callback function
		triggerCondition := usecase.NewTriggerCondition(50000.0, ">").
			WithTimeout(5).
			WithCheckInterval(10).
			WithCallback(func(price float64) {
				callbackExecuted = true
				callbackPrice = price
			})

		// Create the request with trigger condition
		request := &usecase.SniperShotRequest{
			UserID:    testUserID,
			Symbol:    testSymbol,
			Side:      model.OrderSideBuy,
			Type:      model.OrderTypeMarket,
			Quantity:  0.1,
			Price:     0,
			TimeLimit: 5 * time.Second,
			Condition: &triggerCondition,
		}

		// Execute the function being tested
		result, err := sniperService.ExecuteSniper(context.Background(), request)

		// Check expectations
		assert.NoError(t, err)
		assert.True(t, result.Success)
		assert.Nil(t, result.Error)
		assert.NotNil(t, result.Order)
		assert.Equal(t, testOrderID, result.Order.OrderID)

		// Verify callback was executed with the correct price
		assert.True(t, callbackExecuted, "Callback function should have been executed")
		assert.Equal(t, 51000.0, callbackPrice, "Callback should receive the correct price")

		// Verify all mocks were called as expected
		mockPriceChecker.AssertExpectations(t)
		mockTradeExecutor.AssertExpectations(t)
		mockOrderRepo.AssertExpectations(t)
	})

	t.Run("Complex price movement with multiple checks", func(t *testing.T) {
		// Reset mock
		mockPriceChecker = new(MockPriceChecker)
		sniperService = service.NewSniperShotService(&logger, mockTradeExecutor, mockOrderRepo, mockPriceChecker)

		// Mock a price sequence that moves up and down before crossing threshold
		// 48000 -> 49000 -> 49500 -> 49200 -> 50100
		prices := []float64{48000.0, 49000.0, 49500.0, 49200.0, 50100.0}
		priceIndex := 0

		mockPriceChecker.On("GetCurrentPrice", mock.Anything, testSymbol).Return(
			func(ctx context.Context, symbol string) float64 {
				price := prices[priceIndex]
				if priceIndex < len(prices)-1 {
					priceIndex++
				}
				return price
			},
			func(ctx context.Context, symbol string) error {
				return nil
			},
		).Times(5)

		// Track price check sequence
		priceSequence := []float64{}

		// Create a trigger condition that activates when price goes above 50000
		triggerCondition := service.NewTriggerCondition(50000.0, ">").
			WithTimeout(5).
			WithCheckInterval(10).
			WithCallback(func(price float64) {
				priceSequence = append(priceSequence, price)
			})

		// Mock successful order response
		mockOrderResponse := &model.OrderResponse{
			Order: model.Order{
				ID:            "internal-id-complex",
				OrderID:       testOrderID,
				ClientOrderID: "client-complex",
				UserID:        testUserID,
				Symbol:        testSymbol,
				Side:          model.OrderSideBuy,
				Type:          model.OrderTypeLimit,
				Status:        model.OrderStatusNew,
				Quantity:      0.1,
				Price:         50100.0, // Use the price at trigger point
				ExecutedQty:   0,
				TimeInForce:   model.TimeInForceGTC,
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
				Exchange:      "binance",
			},
			IsSuccess: true,
		}

		// Set up the mock expectations for trade execution
		mockTradeExecutor.On("ExecuteOrder", mock.Anything, mock.Anything).Return(mockOrderResponse, nil).Once()

		// Mock repository save
		mockOrderRepo.On("Create", mock.Anything, mock.Anything).Return(nil).Once()

		// Create the request with trigger condition
		request := &service.SniperShotRequest{
			UserID:    testUserID,
			Symbol:    testSymbol,
			Side:      model.OrderSideBuy,
			Type:      model.OrderTypeLimit,
			Quantity:  0.1,
			Price:     50100.0, // Use current price as limit price
			TimeLimit: 5 * time.Second,
			Condition: &triggerCondition,
		}

		// Execute the function being tested
		result, err := sniperService.ExecuteSniper(context.Background(), request)

		// Check expectations
		assert.NoError(t, err)
		assert.True(t, result.Success)
		assert.Nil(t, result.Error)
		assert.NotNil(t, result.Order)
		assert.Equal(t, testOrderID, result.Order.OrderID)

		// Verify that the last price check was the one that triggered (50100)
		assert.Len(t, priceSequence, 1, "Only one callback should be executed")
		assert.Equal(t, 50100.0, priceSequence[0], "The triggering price should be the callback value")

		// Verify the total number of price checks
		mockPriceChecker.AssertNumberOfCalls(t, "GetCurrentPrice", 5)

		// Verify all mocks were called as expected
		mockPriceChecker.AssertExpectations(t)
		mockTradeExecutor.AssertExpectations(t)
		mockOrderRepo.AssertExpectations(t)
	})
}
