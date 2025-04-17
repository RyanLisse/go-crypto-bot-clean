package mocks

import (
	"errors"
	"testing"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	mocks "github.com/RyanLisse/go-crypto-bot-clean/backend/internal/mocks/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Test to verify that our MockNewCoinUseCase works correctly
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
		status := model.CoinStatus("trading")
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
		expectedCoin := &model.NewCoin{
			ID:     "1",
			Symbol: symbol,
			Name:   "Bitcoin",
			Status: model.CoinStatus("trading"),
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
		status := model.CoinStatus("listed")
		limit := 10
		offset := 0
		expectedCoins := []*model.NewCoin{
			{ID: "1", Symbol: "BTC", Name: "Bitcoin", Status: model.CoinStatus("trading")},
			{ID: "2", Symbol: "ETH", Name: "Ethereum", Status: model.CoinStatus("listed")},
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
		expectedCoins := []*model.NewCoin{
			{ID: "1", Symbol: "BTC", Name: "Bitcoin", Status: model.CoinStatus("trading")},
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
		callback := func(*model.NewCoinEvent) {}
		mockUC.On("SubscribeToEvents", mock.AnythingOfType("func(*model.NewCoinEvent)")).Return(nil)

		// Call method
		err := mockUC.SubscribeToEvents(callback)

		// Assert results
		assert.NoError(t, err)
		mockUC.AssertExpectations(t)
	})

	t.Run("UnsubscribeFromEvents", func(t *testing.T) {
		// Setup mock
		callback := func(*model.NewCoinEvent) {}
		mockUC.On("UnsubscribeFromEvents", mock.AnythingOfType("func(*model.NewCoinEvent)")).Return(nil)

		// Call method
		err := mockUC.UnsubscribeFromEvents(callback)

		// Assert results
		assert.NoError(t, err)
		mockUC.AssertExpectations(t)
	})
}
