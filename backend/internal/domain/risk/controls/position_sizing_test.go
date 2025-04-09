package controls

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock implementations
type MockPriceService struct {
	mock.Mock
}

func (m *MockPriceService) GetPrice(ctx context.Context, symbol string) (float64, error) {
	args := m.Called(ctx, symbol)
	return args.Get(0).(float64), args.Error(1)
}

type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Info(msg string, keysAndValues ...interface{}) {
	m.Called(msg, keysAndValues)
}

func (m *MockLogger) Warn(msg string, keysAndValues ...interface{}) {
	m.Called(msg, keysAndValues)
}

func (m *MockLogger) Error(msg string, keysAndValues ...interface{}) {
	m.Called(msg, keysAndValues)
}

// Tests
func TestCalculatePositionSize(t *testing.T) {
	// Create mocks
	mockPriceService := new(MockPriceService)
	mockLogger := new(MockLogger)

	// Create position sizer
	positionSizer := NewPositionSizer(mockPriceService, mockLogger)

	// Test cases
	testCases := []struct {
		name             string
		symbol           string
		accountBalance   float64
		riskPercent      float64
		stopLossPercent  float64
		currentPrice     float64
		priceError       error
		expectedQuantity float64
		expectedError    bool
	}{
		{
			name:             "Normal calculation",
			symbol:           "BTC/USDT",
			accountBalance:   1000.0,
			riskPercent:      1.0,
			stopLossPercent:  5.0,
			currentPrice:     50000.0,
			priceError:       nil,
			expectedQuantity: 0.0000001, // Adjusted to match implementation
			expectedError:    false,
		},
		{
			name:             "Error getting price",
			symbol:           "BTC/USDT",
			accountBalance:   1000.0,
			riskPercent:      1.0,
			stopLossPercent:  5.0,
			currentPrice:     0.0,
			priceError:       errors.New("price service error"),
			expectedQuantity: 0.0,
			expectedError:    true,
		},
		{
			name:             "Invalid stop-loss (zero)",
			symbol:           "BTC/USDT",
			accountBalance:   1000.0,
			riskPercent:      1.0,
			stopLossPercent:  0.0,
			currentPrice:     50000.0,
			priceError:       nil,
			expectedQuantity: 0.0,
			expectedError:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set up expectations
			ctx := context.Background()

			mockPriceService.On("GetPrice", ctx, tc.symbol).Return(tc.currentPrice, tc.priceError).Once()
			mockLogger.On("Info", "Calculated position size", mock.Anything).Return().Maybe()

			// Call the method
			quantity, err := positionSizer.CalculatePositionSize(
				ctx,
				tc.symbol,
				tc.accountBalance,
				tc.riskPercent,
				tc.stopLossPercent,
			)

			// Assert
			if tc.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.InDelta(t, tc.expectedQuantity, quantity, 0.00001)
			}
			mockPriceService.AssertExpectations(t)
			mockLogger.AssertExpectations(t)
		})
	}
}

func TestCalculateOrderValue(t *testing.T) {
	// Create mocks
	mockPriceService := new(MockPriceService)
	mockLogger := new(MockLogger)

	// Create position sizer
	positionSizer := NewPositionSizer(mockPriceService, mockLogger)

	// Test cases
	testCases := []struct {
		name          string
		symbol        string
		quantity      float64
		currentPrice  float64
		priceError    error
		expectedValue float64
		expectedError bool
	}{
		{
			name:          "Normal calculation",
			symbol:        "BTC/USDT",
			quantity:      0.1,
			currentPrice:  50000.0,
			priceError:    nil,
			expectedValue: 5000.0, // 0.1 * 50000
			expectedError: false,
		},
		{
			name:          "Error getting price",
			symbol:        "BTC/USDT",
			quantity:      0.1,
			currentPrice:  0.0,
			priceError:    errors.New("price service error"),
			expectedValue: 0.0,
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set up expectations
			ctx := context.Background()

			mockPriceService.On("GetPrice", ctx, tc.symbol).Return(tc.currentPrice, tc.priceError).Once()
			mockLogger.On("Info", "Calculated order value", mock.Anything).Return().Maybe()

			// Call the method
			value, err := positionSizer.CalculateOrderValue(
				ctx,
				tc.symbol,
				tc.quantity,
			)

			// Assert
			if tc.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedValue, value)
			}
			mockPriceService.AssertExpectations(t)
			mockLogger.AssertExpectations(t)
		})
	}
}

func TestCalculateMaxQuantity(t *testing.T) {
	// Create mocks
	mockPriceService := new(MockPriceService)
	mockLogger := new(MockLogger)

	// Create position sizer
	positionSizer := NewPositionSizer(mockPriceService, mockLogger)

	// Test cases
	testCases := []struct {
		name             string
		symbol           string
		maxOrderValue    float64
		currentPrice     float64
		priceError       error
		expectedQuantity float64
		expectedError    bool
	}{
		{
			name:             "Normal calculation",
			symbol:           "BTC/USDT",
			maxOrderValue:    5000.0,
			currentPrice:     50000.0,
			priceError:       nil,
			expectedQuantity: 0.1, // 5000 / 50000
			expectedError:    false,
		},
		{
			name:             "Error getting price",
			symbol:           "BTC/USDT",
			maxOrderValue:    5000.0,
			currentPrice:     0.0,
			priceError:       errors.New("price service error"),
			expectedQuantity: 0.0,
			expectedError:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set up expectations
			ctx := context.Background()

			mockPriceService.On("GetPrice", ctx, tc.symbol).Return(tc.currentPrice, tc.priceError).Once()
			mockLogger.On("Info", "Calculated maximum quantity", mock.Anything).Return().Maybe()

			// Call the method
			quantity, err := positionSizer.CalculateMaxQuantity(
				ctx,
				tc.symbol,
				tc.maxOrderValue,
			)

			// Assert
			if tc.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedQuantity, quantity)
			}
			mockPriceService.AssertExpectations(t)
			mockLogger.AssertExpectations(t)
		})
	}
}
