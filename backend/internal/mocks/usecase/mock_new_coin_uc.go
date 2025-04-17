package mocks

import (
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/stretchr/testify/mock"
)

// MockNewCoinUseCase is a testable mock for the NewCoinUseCase interface
type MockNewCoinUseCase struct {
	mock.Mock
}

// DetectNewCoins mocks the DetectNewCoins method
func (m *MockNewCoinUseCase) DetectNewCoins() error {
	args := m.Called()
	return args.Error(0)
}

// UpdateCoinStatus mocks the UpdateCoinStatus method
func (m *MockNewCoinUseCase) UpdateCoinStatus(coinID string, newStatus model.CoinStatus) error {
	args := m.Called(coinID, newStatus)
	return args.Error(0)
}

// GetCoinDetails mocks the GetCoinDetails method
func (m *MockNewCoinUseCase) GetCoinDetails(symbol string) (*model.NewCoin, error) {
	args := m.Called(symbol)

	var coin *model.NewCoin
	if args.Get(0) != nil {
		coin = args.Get(0).(*model.NewCoin)
	}

	return coin, args.Error(1)
}

// ListNewCoins mocks the ListNewCoins method
func (m *MockNewCoinUseCase) ListNewCoins(status model.CoinStatus, limit, offset int) ([]*model.NewCoin, error) {
	args := m.Called(status, limit, offset)

	var coins []*model.NewCoin
	if args.Get(0) != nil {
		coins = args.Get(0).([]*model.NewCoin)
	}

	return coins, args.Error(1)
}

// GetRecentTradableCoins mocks the GetRecentTradableCoins method
func (m *MockNewCoinUseCase) GetRecentTradableCoins(limit int) ([]*model.NewCoin, error) {
	args := m.Called(limit)

	var coins []*model.NewCoin
	if args.Get(0) != nil {
		coins = args.Get(0).([]*model.NewCoin)
	}

	return coins, args.Error(1)
}

// SubscribeToEvents mocks the SubscribeToEvents method
func (m *MockNewCoinUseCase) SubscribeToEvents(callback func(*model.NewCoinEvent)) error {
	args := m.Called(callback)
	return args.Error(0)
}

// UnsubscribeFromEvents mocks the UnsubscribeFromEvents method
func (m *MockNewCoinUseCase) UnsubscribeFromEvents(callback func(*model.NewCoinEvent)) error {
	args := m.Called(callback)
	return args.Error(0)
}
