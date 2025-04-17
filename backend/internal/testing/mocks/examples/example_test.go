package examples_test

import (
	"errors"
	"testing"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/testing/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Example component that uses the NewCoinUseCase
type CoinAlertService struct {
	coinUC model.NewCoinService // Use the domain interface here
}

// Example method that uses the NewCoinUseCase
func (s *CoinAlertService) AlertOnNewTrading(limit int) ([]string, error) {
	coins, err := s.coinUC.GetRecentTradableCoins(limit)
	if err != nil {
		return nil, err
	}

	result := make([]string, len(coins))
	for i, coin := range coins {
		result[i] = coin.Symbol
	}
	return result, nil
}

// Example test showing how to use the adapter with our mock
func TestCoinAlertService_WithMock(t *testing.T) {
	// Create our mock adapter
	mockUC := mocks.NewMockNewCoinUseCase()

	// Create the service using our mock adapter
	service := &CoinAlertService{
		coinUC: mockUC, // The adapter implements the domain interface
	}

	// Setup expectations on the underlying mock
	mockUC.Mock.On("GetRecentTradableCoins", 5).Return([]*mocks.NewCoin{
		{ID: "1", Symbol: "BTC", Status: mocks.CoinStatusTrading},
		{ID: "2", Symbol: "ETH", Status: mocks.CoinStatusTrading},
		{ID: "3", Symbol: "XRP", Status: mocks.CoinStatusTrading},
	}, nil)

	// Call the service method
	symbols, err := service.AlertOnNewTrading(5)

	// Assert results
	assert.NoError(t, err)
	assert.Equal(t, []string{"BTC", "ETH", "XRP"}, symbols)
	mockUC.Mock.AssertExpectations(t)
}

// Example test showing error handling
func TestCoinAlertService_ErrorHandling(t *testing.T) {
	// Create our mock adapter
	mockUC := mocks.NewMockNewCoinUseCase()

	// Create the service using our mock adapter
	service := &CoinAlertService{
		coinUC: mockUC,
	}

	// Setup expectations with an error
	expectedErr := errors.New("database connection error")
	mockUC.Mock.On("GetRecentTradableCoins", 5).Return(nil, expectedErr)

	// Call the service method
	symbols, err := service.AlertOnNewTrading(5)

	// Assert error is passed through
	assert.Nil(t, symbols)
	assert.Equal(t, expectedErr, err)
	mockUC.Mock.AssertExpectations(t)
}

// Example test demonstrating full event handling
func TestCoinEventSubscription(t *testing.T) {
	// Create our mock adapter
	mockUC := mocks.NewMockNewCoinUseCase()

	// Variables to track callback execution
	var (
		callbackCalled = false
		receivedEvent  *model.NewCoinEvent
	)

	// Setup the subscribe expectation with callback capture
	var capturedCallback func(*mocks.NewCoinEvent)
	mockUC.Mock.On("SubscribeToEvents", mock.AnythingOfType("func(*mocks.NewCoinEvent)")).
		Run(func(args mock.Arguments) {
			capturedCallback = args.Get(0).(func(*mocks.NewCoinEvent))
		}).
		Return(nil)

	// Subscribe to events with our callback
	err := mockUC.SubscribeToEvents(func(event *model.NewCoinEvent) {
		callbackCalled = true
		receivedEvent = event
	})
	assert.NoError(t, err)

	// Simulate an event by calling the captured callback
	mockEvent := &mocks.NewCoinEvent{
		ID:        "event1",
		CoinID:    "BTC-USDT",
		EventType: "status_change",
		OldStatus: mocks.CoinStatus("listed"),
		NewStatus: mocks.CoinStatus("trading"),
		CreatedAt: time.Now(),
	}
	capturedCallback(mockEvent)

	// Verify our callback was called with the correct data
	assert.True(t, callbackCalled, "Event callback should have been called")
	assert.NotNil(t, receivedEvent)
	assert.Equal(t, mockEvent.CoinID, receivedEvent.CoinID)
	assert.Equal(t, mockEvent.EventType, receivedEvent.EventType)

	mockUC.Mock.AssertExpectations(t)
}
