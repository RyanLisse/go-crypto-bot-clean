package usecase

import (
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/stretchr/testify/mock"
)

// MockNewCoinUseCase is a mock implementation of the NewCoinUseCase interface for testing
type MockNewCoinUseCase struct {
	mock.Mock
}

// DetectNewCoins mocks the method to check for newly listed coins on MEXC
func (m *MockNewCoinUseCase) DetectNewCoins() error {
	args := m.Called()
	return args.Error(0)
}

// UpdateCoinStatus mocks the method to update a coin's status and create an event
func (m *MockNewCoinUseCase) UpdateCoinStatus(coinID string, newStatus model.CoinStatus) error {
	args := m.Called(coinID, newStatus)
	return args.Error(0)
}

// GetCoinDetails mocks the method to retrieve detailed information about a coin
func (m *MockNewCoinUseCase) GetCoinDetails(coinID string) (*model.Coin, error) {
	args := m.Called(coinID)
	if coin := args.Get(0); coin != nil {
		return coin.(*model.Coin), args.Error(1)
	}
	return nil, args.Error(1)
}

// SubscribeToEvents mocks the method to allow subscribing to new coin events
func (m *MockNewCoinUseCase) SubscribeToEvents(handler func(*model.CoinEvent)) error {
	args := m.Called(handler)
	return args.Error(0)
}
