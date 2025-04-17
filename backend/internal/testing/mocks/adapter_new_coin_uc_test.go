package mocks_test

import (
	"testing"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/testing/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewCoinUseCaseAdapter(t *testing.T) {
	adapter := mocks.NewMockNewCoinUseCase()

	t.Run("AdapterDelegatesDetectNewCoins", func(t *testing.T) {
		// Setup mock
		adapter.Mock.On("DetectNewCoins").Return(nil).Once()

		// Call method through adapter
		err := adapter.DetectNewCoins()

		// Assert results
		assert.NoError(t, err)
		adapter.Mock.AssertExpectations(t)
	})

	t.Run("AdapterConvertsUpdateCoinStatus", func(t *testing.T) {
		// Setup mock with our internal CoinStatus type
		coinID := "BTC-USDT"
		mockStatus := mocks.CoinStatus("trading")
		adapter.Mock.On("UpdateCoinStatus", coinID, mockStatus).Return(nil).Once()

		// Call method through adapter with model.CoinStatus
		err := adapter.UpdateCoinStatus(coinID, model.CoinStatus("trading"))

		// Assert results
		assert.NoError(t, err)
		adapter.Mock.AssertExpectations(t)
	})

	t.Run("AdapterConvertsGetCoinDetails", func(t *testing.T) {
		// Setup mock with our internal NewCoin type
		symbol := "BTC"
		now := time.Now()
		mockCoin := &mocks.NewCoin{
			ID:                  "1",
			Symbol:              symbol,
			Name:                "Bitcoin",
			Status:              mocks.CoinStatus("trading"),
			ExpectedListingTime: now,
			CreatedAt:           now,
			UpdatedAt:           now,
		}
		adapter.Mock.On("GetCoinDetails", symbol).Return(mockCoin, nil).Once()

		// Call method through adapter
		coin, err := adapter.GetCoinDetails(symbol)

		// Assert results
		assert.NoError(t, err)
		assert.NotNil(t, coin)
		assert.Equal(t, mockCoin.ID, coin.ID)
		assert.Equal(t, mockCoin.Symbol, coin.Symbol)
		assert.Equal(t, mockCoin.Name, coin.Name)
		assert.Equal(t, model.CoinStatus(mockCoin.Status), coin.Status)
		assert.Equal(t, mockCoin.ExpectedListingTime, coin.ExpectedListingTime)
		adapter.Mock.AssertExpectations(t)
	})

	t.Run("AdapterConvertsListNewCoins", func(t *testing.T) {
		// Setup mock
		status := mocks.CoinStatus("listed")
		limit := 10
		offset := 0
		mockCoins := []*mocks.NewCoin{
			{ID: "1", Symbol: "BTC", Name: "Bitcoin", Status: mocks.CoinStatus("trading")},
			{ID: "2", Symbol: "ETH", Name: "Ethereum", Status: mocks.CoinStatus("listed")},
		}
		adapter.Mock.On("ListNewCoins", status, limit, offset).Return(mockCoins, nil).Once()

		// Call method through adapter
		coins, err := adapter.ListNewCoins(model.CoinStatus("listed"), limit, offset)

		// Assert results
		assert.NoError(t, err)
		assert.Len(t, coins, 2)
		assert.Equal(t, mockCoins[0].ID, coins[0].ID)
		assert.Equal(t, mockCoins[0].Symbol, coins[0].Symbol)
		assert.Equal(t, model.CoinStatus(mockCoins[0].Status), coins[0].Status)
		adapter.Mock.AssertExpectations(t)
	})

	t.Run("AdapterHandlesSubscribeToEvents", func(t *testing.T) {
		// Setup mock
		var capturedCallback func(*mocks.NewCoinEvent)
		adapter.Mock.On("SubscribeToEvents", mock.AnythingOfType("func(*mocks.NewCoinEvent)")).
			Run(func(args mock.Arguments) {
				capturedCallback = args.Get(0).(func(*mocks.NewCoinEvent))
			}).
			Return(nil).
			Once()

		// Track if our callback was called
		callbackCalled := false
		var receivedEvent *model.NewCoinEvent

		// Call method through adapter
		err := adapter.SubscribeToEvents(func(event *model.NewCoinEvent) {
			callbackCalled = true
			receivedEvent = event
		})

		assert.NoError(t, err)

		// Simulate event through captured callback
		if capturedCallback != nil {
			mockEvent := &mocks.NewCoinEvent{
				ID:        "event1",
				CoinID:    "BTC-USDT",
				EventType: "status_change",
				OldStatus: mocks.CoinStatus("listed"),
				NewStatus: mocks.CoinStatus("trading"),
				CreatedAt: time.Now(),
			}
			capturedCallback(mockEvent)

			// Verify our adapter's callback was called with converted event
			assert.True(t, callbackCalled)
			assert.NotNil(t, receivedEvent)
			assert.Equal(t, mockEvent.ID, receivedEvent.ID)
			assert.Equal(t, mockEvent.CoinID, receivedEvent.CoinID)
			assert.Equal(t, mockEvent.EventType, receivedEvent.EventType)
			assert.Equal(t, model.CoinStatus(mockEvent.OldStatus), receivedEvent.OldStatus)
			assert.Equal(t, model.CoinStatus(mockEvent.NewStatus), receivedEvent.NewStatus)
		}

		adapter.Mock.AssertExpectations(t)
	})
}
