package mocks

import (
	"errors"
	"testing"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	uc_mocks "github.com/RyanLisse/go-crypto-bot-clean/backend/internal/mocks/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Test to verify that our MockNewCoinUseCase works correctly
func TestMockNewCoinUseCase(t *testing.T) {
	// Create the mock
	mockUC := new(uc_mocks.MockNewCoinUseCase)

	// Test DetectNewCoins
	mockUC.On("DetectNewCoins").Return(nil).Once()
	err := mockUC.DetectNewCoins()
	assert.NoError(t, err)

	// Test DetectNewCoins with error
	expectedErr := errors.New("detection error")
	mockUC.On("DetectNewCoins").Return(expectedErr).Once()
	err = mockUC.DetectNewCoins()
	assert.Equal(t, expectedErr, err)

	// Test UpdateCoinStatus
	mockUC.On("UpdateCoinStatus", "coin123", model.CoinStatus("trading")).Return(nil).Once()
	err = mockUC.UpdateCoinStatus("coin123", model.CoinStatus("trading"))
	assert.NoError(t, err)

	// Test GetCoinDetails
	mockCoin := &model.Coin{
		ID:     "coin123",
		Symbol: "BTC",
	}
	mockUC.On("GetCoinDetails", "coin123").Return(mockCoin, nil).Once()
	coin, err := mockUC.GetCoinDetails("coin123")
	assert.NoError(t, err)
	assert.Equal(t, mockCoin, coin)

	// Test SubscribeToEvents
	var capturedHandler func(*model.CoinEvent)
	mockUC.On("SubscribeToEvents", mock.AnythingOfType("func(*model.CoinEvent)")).
		Run(func(args mock.Arguments) {
			capturedHandler = args.Get(0).(func(*model.CoinEvent))
		}).
		Return(nil).
		Once()

	// Call SubscribeToEvents with a dummy handler
	dummyHandlerCalled := false
	dummyHandler := func(event *model.CoinEvent) {
		dummyHandlerCalled = true
	}

	err = mockUC.SubscribeToEvents(dummyHandler)
	assert.NoError(t, err)

	// Simulate an event by calling the captured handler
	if capturedHandler != nil {
		capturedHandler(&model.CoinEvent{CoinID: "coin123"})
		assert.True(t, dummyHandlerCalled, "Handler should have been called")
	}

	// Verify all expectations were met
	mockUC.AssertExpectations(t)
}
