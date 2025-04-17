package usecase

import (
	"context"
	"testing"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	portmocks "github.com/RyanLisse/go-crypto-bot-clean/backend/internal/mocks/domain/port"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockNewCoinUseCase is a test-specific mock implementation
type MockNewCoinUseCase struct {
	mock.Mock
}

// DetectNewCoins mocks the method to check for newly listed coins
func (m *MockNewCoinUseCase) DetectNewCoins() error {
	args := m.Called()
	return args.Error(0)
}

// UpdateCoinStatus mocks the method to update a coin's status
func (m *MockNewCoinUseCase) UpdateCoinStatus(coinID string, newStatus model.CoinStatus) error {
	args := m.Called(coinID, newStatus)
	return args.Error(0)
}

// GetCoinDetails mocks the method to retrieve a coin's details
func (m *MockNewCoinUseCase) GetCoinDetails(coinID string) (*model.Coin, error) {
	args := m.Called(coinID)
	if coin := args.Get(0); coin != nil {
		return coin.(*model.Coin), args.Error(1)
	}
	return nil, args.Error(1)
}

// SubscribeToEvents mocks the method to subscribe to new coin events
func (m *MockNewCoinUseCase) SubscribeToEvents(handler func(*model.CoinEvent)) error {
	args := m.Called(handler)
	return args.Error(0)
}

func TestSniperUseCase_ExecuteSnipe(t *testing.T) {
	// Create mocks
	mockSniperService := &portmocks.SniperService{}
	mockNewCoinUC := &MockNewCoinUseCase{}

	// Create logger
	logger := zerolog.New(zerolog.NewTestWriter(t))

	// Create test symbol
	symbol := "BTC_USDT"

	// Create test order
	order := &model.Order{
		OrderID:  "test-order-id",
		Symbol:   symbol,
		Side:     model.OrderSideBuy,
		Type:     model.OrderTypeMarket,
		Quantity: 0.002,
		Price:    50000.0,
		Status:   model.OrderStatusFilled,
	}

	// Set up mock expectations
	mockSniperService.On("Start").Return(nil)
	mockSniperService.On("ExecuteSnipe", mock.Anything, symbol).Return(order, nil)
	mockNewCoinUC.On("SubscribeToEvents", mock.AnythingOfType("func(*model.CoinEvent)")).Return(nil)

	// Create sniper use case
	sniperUC := NewSniperUseCase(mockSniperService, mockNewCoinUC, &logger)

	// Execute snipe
	ctx := context.Background()
	result, err := sniperUC.ExecuteSnipe(ctx, symbol)

	// Verify results
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, symbol, result.Symbol)
	assert.Equal(t, model.OrderSideBuy, result.Side)
	assert.Equal(t, model.OrderTypeMarket, result.Type)

	// Verify mock expectations
	mockSniperService.AssertExpectations(t)
	mockNewCoinUC.AssertExpectations(t)
}

func TestSniperUseCase_ExecuteSnipeWithConfig(t *testing.T) {
	// Create mocks
	mockSniperService := &portmocks.SniperService{}
	mockNewCoinUC := &MockNewCoinUseCase{}

	// Create logger
	logger := zerolog.New(zerolog.NewTestWriter(t))

	// Create test symbol
	symbol := "BTC_USDT"

	// Create test config
	config := &port.SniperConfig{
		MaxBuyAmount:     200.0,
		MaxPricePerToken: 2.0,
	}

	// Create test order
	order := &model.Order{
		OrderID:  "test-order-id",
		Symbol:   symbol,
		Side:     model.OrderSideBuy,
		Type:     model.OrderTypeMarket,
		Quantity: 0.002,
		Price:    50000.0,
		Status:   model.OrderStatusFilled,
	}

	// Set up mock expectations
	mockSniperService.On("Start").Return(nil)
	mockSniperService.On("ExecuteSnipeWithConfig", mock.Anything, symbol, config).Return(order, nil)
	mockNewCoinUC.On("SubscribeToEvents", mock.AnythingOfType("func(*model.CoinEvent)")).Return(nil)

	// Create sniper use case
	sniperUC := NewSniperUseCase(mockSniperService, mockNewCoinUC, &logger)

	// Execute snipe with config
	ctx := context.Background()
	result, err := sniperUC.ExecuteSnipeWithConfig(ctx, symbol, config)

	// Verify results
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, symbol, result.Symbol)
	assert.Equal(t, model.OrderSideBuy, result.Side)
	assert.Equal(t, model.OrderTypeMarket, result.Type)

	// Verify mock expectations
	mockSniperService.AssertExpectations(t)
	mockNewCoinUC.AssertExpectations(t)
}

func TestSniperUseCase_SetupAutoSnipe(t *testing.T) {
	// Create mocks
	mockSniperService := &portmocks.SniperService{}
	mockNewCoinUC := &MockNewCoinUseCase{}

	// Create logger
	logger := zerolog.New(zerolog.NewTestWriter(t))

	// Create test config
	config := &port.SniperConfig{
		MaxBuyAmount:     200.0,
		MaxPricePerToken: 2.0,
	}

	// Set up mock expectations
	mockSniperService.On("Start").Return(nil)
	mockSniperService.On("GetConfig").Return(config)
	mockNewCoinUC.On("SubscribeToEvents", mock.AnythingOfType("func(*model.CoinEvent)")).Return(nil)

	// Create sniper use case
	sniperUC := NewSniperUseCase(mockSniperService, mockNewCoinUC, &logger)

	// Test enabling auto-snipe
	err := sniperUC.SetupAutoSnipe(true, config)
	assert.NoError(t, err)

	// We need to use reflection or type assertion to access private fields
	// For testing purposes, we'll check the behavior instead
	// Call GetSniperConfig to verify the config was set
	actualConfig, err := sniperUC.GetSniperConfig()
	assert.NoError(t, err)
	assert.Equal(t, config, actualConfig)

	// Test disabling auto-snipe
	err = sniperUC.SetupAutoSnipe(false, nil)
	assert.NoError(t, err)

	// Verify mock expectations
	mockSniperService.AssertExpectations(t)
	mockNewCoinUC.AssertExpectations(t)
}
