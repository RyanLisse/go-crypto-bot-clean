package portfolio

import (
	"context"
	"testing"
	"time"

	"go-crypto-bot-clean/backend/internal/domain/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockWalletProvider is a mock implementation of the WalletProvider interface
type MockWalletProvider struct {
	mock.Mock
}

// GetWallet mocks the GetWallet method
func (m *MockWalletProvider) GetWallet(ctx context.Context) (*models.Wallet, error) {
	args := m.Called(ctx)
	return args.Get(0).(*models.Wallet), args.Error(1)
}

// MockBoughtCoinRepo is a mock implementation of the BoughtCoinRepository
type MockBoughtCoinRepo struct {
	mock.Mock
}

// FindAll mocks the FindAll method
func (m *MockBoughtCoinRepo) FindAll(ctx context.Context) ([]models.BoughtCoin, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.BoughtCoin), args.Error(1)
}

// FindByID mocks the FindByID method
func (m *MockBoughtCoinRepo) FindByID(ctx context.Context, id int64) (*models.BoughtCoin, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.BoughtCoin), args.Error(1)
}

// FindBySymbol mocks the FindBySymbol method
func (m *MockBoughtCoinRepo) FindBySymbol(ctx context.Context, symbol string) (*models.BoughtCoin, error) {
	args := m.Called(ctx, symbol)
	return args.Get(0).(*models.BoughtCoin), args.Error(1)
}

// Create mocks the Create method
func (m *MockBoughtCoinRepo) Create(ctx context.Context, coin *models.BoughtCoin) (int64, error) {
	args := m.Called(ctx, coin)
	return args.Get(0).(int64), args.Error(1)
}

// Update mocks the Update method
func (m *MockBoughtCoinRepo) Update(ctx context.Context, coin *models.BoughtCoin) error {
	args := m.Called(ctx, coin)
	return args.Error(0)
}

// Delete mocks the Delete method
func (m *MockBoughtCoinRepo) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestGetWallet(t *testing.T) {
	// Create mocks
	mockWallet := new(MockWalletProvider)
	mockRepo := new(MockBoughtCoinRepo)

	// Setup expected wallet data
	expectedWallet := &models.Wallet{
		Balances: map[string]*models.AssetBalance{
			"USDT": {
				Asset:  "USDT",
				Free:   500.0,
				Locked: 0.0,
				Total:  500.0,
			},
			"BTC": {
				Asset:  "BTC",
				Free:   0.01,
				Locked: 0.0,
				Total:  0.01,
			},
		},
		UpdatedAt: time.Now(),
	}

	// Setup mock expectations
	mockWallet.On("GetWallet", mock.Anything).Return(expectedWallet, nil)

	// Create a test service
	service := &portfolioService{
		walletProvider: mockWallet,
		boughtCoinRepo: mockRepo,
	}

	// Call the method
	wallet, err := service.GetWallet(context.Background())

	// Assert expectations
	assert.NoError(t, err)
	assert.Equal(t, expectedWallet, wallet)
	mockWallet.AssertExpectations(t)
}

func TestGetPositions(t *testing.T) {
	// Create mocks
	mockWallet := new(MockWalletProvider)
	mockRepo := new(MockBoughtCoinRepo)

	// Setup test data
	boughtCoins := []models.BoughtCoin{
		{
			ID:            1,
			Symbol:        "BTC/USDT",
			Quantity:      0.01,
			PurchasePrice: 20000.0,
			BoughtAt:      time.Now(),
			StopLoss:      19000.0,
			TakeProfit:    21000.0,
		},
		{
			ID:            2,
			Symbol:        "ETH/USDT",
			Quantity:      0.5,
			PurchasePrice: 2000.0,
			BoughtAt:      time.Now(),
			StopLoss:      1900.0,
			TakeProfit:    2100.0,
		},
	}

	// Expected positions after conversion
	expectedPositions := []models.Position{
		{
			ID:         "BTC/USDT",
			Symbol:     "BTC/USDT",
			Quantity:   0.01,
			EntryPrice: 20000.0,
			OpenedAt:   boughtCoins[0].BoughtAt,
			StopLoss:   19000.0,
			TakeProfit: 21000.0,
			Status:     "open",
		},
		{
			ID:         "ETH/USDT",
			Symbol:     "ETH/USDT",
			Quantity:   0.5,
			EntryPrice: 2000.0,
			OpenedAt:   boughtCoins[1].BoughtAt,
			StopLoss:   1900.0,
			TakeProfit: 2100.0,
			Status:     "open",
		},
	}

	// Setup mock expectations
	mockRepo.On("FindAll", mock.Anything).Return(boughtCoins, nil)

	// Create a test service
	service := &portfolioService{
		walletProvider: mockWallet,
		boughtCoinRepo: mockRepo,
	}

	// Call the method
	positions, err := service.GetPositions(context.Background())

	// Assert expectations
	assert.NoError(t, err)
	assert.Equal(t, expectedPositions, positions)
	mockRepo.AssertExpectations(t)
}
