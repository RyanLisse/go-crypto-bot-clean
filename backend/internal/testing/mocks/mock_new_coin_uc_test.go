package mocks_test

import (
	"errors"
	"testing"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/testing/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMockNewCoinUseCase(t *testing.T) {
	mockUC := &mocks.MockNewCoinUseCase{}

	t.Run("DetectNewCoins", func(t *testing.T) {
		// Setup mock
		expectedErr := errors.New("detection failed")
		mockUC.On("DetectNewCoins").Return(expectedErr)

		// Call method
		err := mockUC.DetectNewCoins()

		// Assert results
		assert.Equal(t, expectedErr, err)
		mockUC.AssertExpectations(t)
	})

	t.Run("UpdateCoinStatus", func(t *testing.T) {
		// Setup mock
		coinID := "BTC-USDT"
		status := mocks.CoinStatusTrading
		mockUC.On("UpdateCoinStatus", coinID, status).Return(nil)

		// Call method
		err := mockUC.UpdateCoinStatus(coinID, status)

		// Assert results
		assert.NoError(t, err)
		mockUC.AssertExpectations(t)
	})

	t.Run("GetCoinDetails", func(t *testing.T) {
		// Setup mock
		symbol := "BTC"
		expectedCoin := &mocks.NewCoin{
			ID:     "1",
			Symbol: symbol,
			Name:   "Bitcoin",
			Status: mocks.CoinStatusTrading,
		}
		mockUC.On("GetCoinDetails", symbol).Return(expectedCoin, nil)

		// Call method
		coin, err := mockUC.GetCoinDetails(symbol)

		// Assert results
		assert.NoError(t, err)
		assert.Equal(t, expectedCoin, coin)
		mockUC.AssertExpectations(t)
	})

	t.Run("ListNewCoins", func(t *testing.T) {
		// Setup mock
		status := mocks.CoinStatusListed
		limit := 10
		offset := 0
		expectedCoins := []*mocks.NewCoin{
			{ID: "1", Symbol: "BTC", Name: "Bitcoin", Status: mocks.CoinStatusTrading},
			{ID: "2", Symbol: "ETH", Name: "Ethereum", Status: mocks.CoinStatusListed},
		}
		mockUC.On("ListNewCoins", status, limit, offset).Return(expectedCoins, nil)

		// Call method
		coins, err := mockUC.ListNewCoins(status, limit, offset)

		// Assert results
		assert.NoError(t, err)
		assert.Equal(t, expectedCoins, coins)
		mockUC.AssertExpectations(t)
	})

	t.Run("GetRecentTradableCoins", func(t *testing.T) {
		// Setup mock
		limit := 5
		expectedCoins := []*mocks.NewCoin{
			{ID: "1", Symbol: "BTC", Name: "Bitcoin", Status: mocks.CoinStatusTrading},
		}
		mockUC.On("GetRecentTradableCoins", limit).Return(expectedCoins, nil)

		// Call method
		coins, err := mockUC.GetRecentTradableCoins(limit)

		// Assert results
		assert.NoError(t, err)
		assert.Equal(t, expectedCoins, coins)
		mockUC.AssertExpectations(t)
	})

	t.Run("SubscribeToEvents", func(t *testing.T) {
		// Setup mock
		callback := func(*mocks.NewCoinEvent) {}
		mockUC.On("SubscribeToEvents", mock.AnythingOfType("func(*mocks.NewCoinEvent)")).Return(nil)

		// Call method
		err := mockUC.SubscribeToEvents(callback)

		// Assert results
		assert.NoError(t, err)
		mockUC.AssertExpectations(t)
	})

	t.Run("UnsubscribeFromEvents", func(t *testing.T) {
		// Setup mock
		callback := func(*mocks.NewCoinEvent) {}
		mockUC.On("UnsubscribeFromEvents", mock.AnythingOfType("func(*mocks.NewCoinEvent)")).Return(nil)

		// Call method
		err := mockUC.UnsubscribeFromEvents(callback)

		// Assert results
		assert.NoError(t, err)
		mockUC.AssertExpectations(t)
	})
}
