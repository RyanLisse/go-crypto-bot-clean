package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/wallet"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// (Removed duplicate MockWalletRepository definition; use the one from wallet_service_test.go)
// If needed, import or reference the shared mock from a test helper.



// GetAllWallets is not part of the current WalletRepository interface and can be removed
// DeleteByID is not part of the current WalletRepository interface and can be removed
// Old SaveBalanceHistory and GetBalanceHistory signatures removed (see above for correct ones)

// MockAPICredentialManagerService is a mock implementation of APICredentialManagerService
// Updated to match the current APICredentialManagerService interface
type MockAPICredentialManagerService struct {
	mock.Mock
}

func (m *MockAPICredentialManagerService) CreateCredential(ctx context.Context, userID, exchange, apiKey, apiSecret, label string) (*model.APICredential, error) {
	args := m.Called(ctx, userID, exchange, apiKey, apiSecret, label)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.APICredential), args.Error(1)
}

func (m *MockAPICredentialManagerService) GetCredential(ctx context.Context, id string) (*model.APICredential, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.APICredential), args.Error(1)
}

func (m *MockAPICredentialManagerService) GetCredentialByUserIDAndExchange(ctx context.Context, userID, exchange string) (*model.APICredential, error) {
	args := m.Called(ctx, userID, exchange)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.APICredential), args.Error(1)
}

func (m *MockAPICredentialManagerService) GetCredentialByUserIDAndLabel(ctx context.Context, userID, exchange, label string) (*model.APICredential, error) {
	args := m.Called(ctx, userID, exchange, label)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.APICredential), args.Error(1)
}

func (m *MockAPICredentialManagerService) UpdateCredential(ctx context.Context, id, apiKey, apiSecret, label string) (*model.APICredential, error) {
	args := m.Called(ctx, id, apiKey, apiSecret, label)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.APICredential), args.Error(1)
}

func (m *MockAPICredentialManagerService) DeleteCredential(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockAPICredentialManagerService) ListCredentialsByUserID(ctx context.Context, userID string) ([]*model.APICredential, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.APICredential), args.Error(1)
}

func (m *MockAPICredentialManagerService) VerifyCredential(ctx context.Context, id string) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockAPICredentialManagerService) RotateCredential(ctx context.Context, id, newAPIKey, newAPISecret string) (*model.APICredential, error) {
	args := m.Called(ctx, id, newAPIKey, newAPISecret)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.APICredential), args.Error(1)
}

func (m *MockAPICredentialManagerService) MarkCredentialAsUsed(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockAPICredentialManagerService) GetCredentialForExchange(ctx context.Context, userID, exchange string) (*model.APICredential, error) {
	args := m.Called(ctx, userID, exchange)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.APICredential), args.Error(1)
}

// WDSMockProviderRegistry is a mock implementation of port.ProviderRegistry
// WDSMockProviderRegistry removed (no longer needed)

func TestSyncWallet(t *testing.T) {
	// Setup
	ctx := context.Background()
	logger := zerolog.New(zerolog.NewTestWriter(t))
	mockRepo := new(MockWalletRepository)
	mockCredentialManager := new(MockAPICredentialManagerService)
	providerRegistry := wallet.NewProviderRegistry()
	mockProvider := new(MockExchangeWalletProvider)

	// Register the mock provider in the registry
	mockProvider.On("GetName").Return("MEXC")
	providerRegistry.RegisterProvider(mockProvider)

	// Create service
	service := NewWalletDataSyncService(mockRepo, mockCredentialManager, providerRegistry, &logger)

	// Setup mock wallet
	wallet := &model.Wallet{
		ID:       "wallet123",
		UserID:   "user123",
		Type:     model.WalletTypeExchange,
		Exchange: "MEXC",
		Status:   model.WalletStatusActive,
	}

	// Setup mock credential
	credential := &model.APICredential{
		ID:        "cred123",
		UserID:    "user123",
		Exchange:  "MEXC",
		APIKey:    "api_key",
		APISecret: "api_secret",
		Status:    model.APICredentialStatusActive,
	}

	// Setup mock repository
	mockRepo.On("GetByID", ctx, "wallet123").Return(wallet, nil)
	mockRepo.On("Save", ctx, mock.AnythingOfType("*model.Wallet")).Return(nil)
	mockRepo.On("SaveBalanceHistory", ctx, mock.AnythingOfType("*model.BalanceHistory")).Return(nil)

	// Setup mock credential manager
	mockCredentialManager.On("GetCredentialForExchange", ctx, "user123", "MEXC").Return(credential, nil)
	mockCredentialManager.On("MarkCredentialAsUsed", ctx, "cred123").Return(nil)

	// Setup mock provider
	mockProvider.On("SetAPICredentials", ctx, "api_key", "api_secret").Return(nil)
	mockProvider.On("GetBalance", ctx, wallet).Return(&model.Wallet{
		ID:            "wallet123",
		UserID:        "user123",
		Type:          model.WalletTypeExchange,
		Exchange:      "MEXC",
		Status:        model.WalletStatusActive,
		Balances:      map[model.Asset]*model.Balance{
			model.AssetBTC: {Asset: model.AssetBTC, Free: 1.0, Locked: 0, Total: 1.0, USDValue: 30000.0},
			model.AssetETH: {Asset: model.AssetETH, Free: 10.0, Locked: 0, Total: 10.0, USDValue: 20000.0},
		},
		TotalUSDValue: 50000.0,
	}, nil)

	// Test
	syncedWallet, err := service.SyncWallet(ctx, "wallet123")

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, syncedWallet)
	assert.Equal(t, "wallet123", syncedWallet.ID)
	assert.Equal(t, "user123", syncedWallet.UserID)
	assert.Equal(t, model.WalletTypeExchange, syncedWallet.Type)
	assert.Equal(t, "MEXC", syncedWallet.Exchange)
	assert.Equal(t, model.WalletStatusActive, syncedWallet.Status)
	assert.Len(t, syncedWallet.Balances, 2)
	assert.Contains(t, syncedWallet.Balances, model.AssetBTC)
	assert.Contains(t, syncedWallet.Balances, model.AssetETH)
	assert.Equal(t, float64(1.0), syncedWallet.Balances[model.AssetBTC].Free)
	assert.Equal(t, float64(10.0), syncedWallet.Balances[model.AssetETH].Free)
	assert.Equal(t, 50000.0, syncedWallet.TotalUSDValue)
	assert.NotNil(t, syncedWallet.LastSynced)

	// Verify mocks
	mockRepo.AssertExpectations(t)
	mockCredentialManager.AssertExpectations(t)
	mockProvider.AssertExpectations(t)
}

func TestSyncWalletError(t *testing.T) {
	// Setup
	ctx := context.Background()
	logger := zerolog.New(zerolog.NewTestWriter(t))
	mockRepo := new(MockWalletRepository)
	mockCredentialManager := new(MockAPICredentialManagerService)
	providerRegistry := wallet.NewProviderRegistry()
	mockProvider := new(MockExchangeWalletProvider)

	// Register the mock provider in the registry
	mockProvider.On("GetName").Return("MEXC")
	providerRegistry.RegisterProvider(mockProvider)

	// Create service
	service := NewWalletDataSyncService(mockRepo, mockCredentialManager, providerRegistry, &logger)

	// Setup mock wallet
	wallet := &model.Wallet{
		ID:       "wallet123",
		UserID:   "user123",
		Type:     model.WalletTypeExchange,
		Exchange: "MEXC",
		Status:   model.WalletStatusActive,
	}

	// Setup mock credential
	credential := &model.APICredential{
		ID:        "cred123",
		UserID:    "user123",
		Exchange:  "MEXC",
		APIKey:    "api_key",
		APISecret: "api_secret",
		Status:    model.APICredentialStatusActive,
	}

	// Setup mock repository
	mockRepo.On("GetByID", ctx, "wallet123").Return(wallet, nil)

	// Setup mock credential manager
	mockCredentialManager.On("GetCredentialForExchange", ctx, "user123", "MEXC").Return(credential, nil)
	mockCredentialManager.On("MarkCredentialAsUsed", ctx, "cred123").Return(nil)

	// Setup mock provider
	mockProvider.On("SetAPICredentials", ctx, "api_key", "api_secret").Return(nil)
	mockProvider.On("GetBalance", ctx, wallet).Return(nil, errors.New("API error"))

	// Test
	_, err := service.SyncWallet(ctx, "wallet123")

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get balance")

	// Verify mocks
	mockRepo.AssertExpectations(t)
	mockCredentialManager.AssertExpectations(t)
	mockProvider.AssertExpectations(t)
}

func TestSyncWalletsByUserID(t *testing.T) {
	// Setup
	ctx := context.Background()
	logger := zerolog.New(zerolog.NewTestWriter(t))
	mockRepo := new(MockWalletRepository)
	mockCredentialManager := new(MockAPICredentialManagerService)
	providerRegistry := wallet.NewProviderRegistry()
	mockProvider := new(MockExchangeWalletProvider)

	// Register the mock provider in the registry
	mockProvider.On("GetName").Return("MEXC")
	providerRegistry.RegisterProvider(mockProvider)

	// Create service
	service := NewWalletDataSyncService(mockRepo, mockCredentialManager, providerRegistry, &logger)

	// Setup mock wallets
	wallet1 := &model.Wallet{
		ID:       "wallet1",
		UserID:   "user123",
		Type:     model.WalletTypeExchange,
		Exchange: "MEXC",
		Status:   model.WalletStatusActive,
	}
	wallet2 := &model.Wallet{
		ID:       "wallet2",
		UserID:   "user123",
		Type:     model.WalletTypeExchange,
		Exchange: "MEXC",
		Status:   model.WalletStatusActive,
	}

	// Setup mock credential
	credential := &model.APICredential{
		ID:        "cred123",
		UserID:    "user123",
		Exchange:  "MEXC",
		APIKey:    "api_key",
		APISecret: "api_secret",
		Status:    model.APICredentialStatusActive,
	}

	// Setup mock repository
	mockRepo.On("GetWalletsByUserID", ctx, "user123").Return([]*model.Wallet{wallet1, wallet2}, nil)
	mockRepo.On("GetByID", ctx, "wallet1").Return(wallet1, nil)
	mockRepo.On("GetByID", ctx, "wallet2").Return(wallet2, nil)
	mockRepo.On("Save", ctx, mock.AnythingOfType("*model.Wallet")).Return(nil).Times(2)
	mockRepo.On("SaveBalanceHistory", ctx, mock.AnythingOfType("*model.BalanceHistory")).Return(nil).Times(2)

	// Setup mock credential manager
	mockCredentialManager.On("GetCredentialForExchange", ctx, "user123", "MEXC").Return(credential, nil).Times(2)
	mockCredentialManager.On("MarkCredentialAsUsed", ctx, "cred123").Return(nil).Times(2)

	// Setup mock provider
	mockProvider.On("SetAPICredentials", ctx, "api_key", "api_secret").Return(nil).Times(2)
	mockProvider.On("GetBalance", ctx, wallet1).Return(&model.Wallet{
		ID:            "wallet1",
		UserID:        "user123",
		Type:          model.WalletTypeExchange,
		Exchange:      "MEXC",
		Status:        model.WalletStatusActive,
		Balances:      map[model.Asset]*model.Balance{
	model.AssetBTC: {Asset: model.AssetBTC, Free: 1.0, Locked: 0, Total: 1.0, USDValue: 30000.0},
},
		TotalUSDValue: 30000.0,
	}, nil)
	mockProvider.On("GetBalance", ctx, wallet2).Return(&model.Wallet{
		ID:            "wallet2",
		UserID:        "user123",
		Type:          model.WalletTypeExchange,
		Exchange:      "MEXC",
		Status:        model.WalletStatusActive,
		Balances:      map[model.Asset]*model.Balance{
	model.AssetETH: {Asset: model.AssetETH, Free: 10.0, Locked: 0, Total: 10.0, USDValue: 20000.0},
},
		TotalUSDValue: 20000.0,
	}, nil)

	// Test
	syncedWallets, err := service.SyncWalletsByUserID(ctx, "user123")

	// Assert
	require.NoError(t, err)
	assert.Len(t, syncedWallets, 2)

	// Verify mocks
	mockRepo.AssertExpectations(t)
	mockCredentialManager.AssertExpectations(t)
	mockProvider.AssertExpectations(t)
}

func TestScheduleWalletSync(t *testing.T) {
	// Setup
	ctx := context.Background()
	logger := zerolog.New(zerolog.NewTestWriter(t))
	mockRepo := new(MockWalletRepository)
	mockCredentialManager := new(MockAPICredentialManagerService)
	providerRegistry := wallet.NewProviderRegistry()

	// Create service
	service := NewWalletDataSyncService(mockRepo, mockCredentialManager, providerRegistry, &logger)

	// Setup mock wallet
	wallet := &model.Wallet{
		ID:       "wallet123",
		UserID:   "user123",
		Type:     model.WalletTypeExchange,
		Exchange: "MEXC",
		Status:   model.WalletStatusActive,
	}

	// Setup mock repository
	mockRepo.On("GetByID", ctx, "wallet123").Return(wallet, nil)

	// Test
	err := service.ScheduleWalletSync(ctx, "wallet123", 5*time.Minute)

	// Assert
	require.NoError(t, err)

	// Verify mocks
	mockRepo.AssertExpectations(t)

	// Cleanup
	err = service.CancelWalletSync(ctx, "wallet123")
	require.NoError(t, err)
}

func TestGetLastSyncTime(t *testing.T) {
	// Setup
	ctx := context.Background()
	logger := zerolog.New(zerolog.NewTestWriter(t))
	mockRepo := new(MockWalletRepository)
	mockCredentialManager := new(MockAPICredentialManagerService)
	providerRegistry := wallet.NewProviderRegistry()

	// Create service
	service := NewWalletDataSyncService(mockRepo, mockCredentialManager, providerRegistry, &logger)

	// Setup mock wallet
	now := time.Now()
	wallet := &model.Wallet{
		ID:         "wallet123",
		UserID:     "user123",
		Type:       model.WalletTypeExchange,
		Exchange:   "MEXC",
		Status:     model.WalletStatusActive,
		LastSynced: &now,
	}

	// Setup mock repository
	mockRepo.On("GetByID", ctx, "wallet123").Return(wallet, nil)

	// Test
	lastSync, err := service.GetLastSyncTime(ctx, "wallet123")

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, lastSync)
	assert.Equal(t, now.Unix(), lastSync.Unix())

	// Verify mocks
	mockRepo.AssertExpectations(t)
}

func TestGetSyncStatus(t *testing.T) {
	// Setup
	ctx := context.Background()
	logger := zerolog.New(zerolog.NewTestWriter(t))
	mockRepo := new(MockWalletRepository)
	mockCredentialManager := new(MockAPICredentialManagerService)
	providerRegistry := wallet.NewProviderRegistry()

	// Create service
	service := NewWalletDataSyncService(mockRepo, mockCredentialManager, providerRegistry, &logger)

	// Setup mock wallet
	wallet := &model.Wallet{
		ID:         "wallet123",
		UserID:     "user123",
		Type:       model.WalletTypeExchange,
		Exchange:   "MEXC",
		Status:     model.WalletStatusActive,
		SyncStatus: model.SyncStatusSuccess,
	}

	// Setup mock repository
	mockRepo.On("GetByID", ctx, "wallet123").Return(wallet, nil)

	// Test
	status, err := service.GetSyncStatus(ctx, "wallet123")

	// Assert
	require.NoError(t, err)
	assert.Equal(t, model.SyncStatusSuccess, status)

	// Verify mocks
	mockRepo.AssertExpectations(t)
}

func TestSaveBalanceHistory(t *testing.T) {
	// Setup
	ctx := context.Background()
	logger := zerolog.New(zerolog.NewTestWriter(t))
	mockRepo := new(MockWalletRepository)
	mockCredentialManager := new(MockAPICredentialManagerService)
	providerRegistry := wallet.NewProviderRegistry()

	// Create service
	service := NewWalletDataSyncService(mockRepo, mockCredentialManager, providerRegistry, &logger)

	// Setup mock wallet
	wallet := &model.Wallet{
		ID:            "wallet123",
		UserID:        "user123",
		Type:          model.WalletTypeExchange,
		Exchange:      "MEXC",
		Status:        model.WalletStatusActive,
		Balances:      map[model.Asset]*model.Balance{
			model.AssetBTC: {Asset: model.AssetBTC, Free: 1.0, Locked: 0, Total: 1.0, USDValue: 30000.0},
			model.AssetETH: {Asset: model.AssetETH, Free: 10.0, Locked: 0, Total: 10.0, USDValue: 20000.0},
		},
		TotalUSDValue: 50000.0,
	}

	// Setup mock repository
	mockRepo.On("GetByID", ctx, "wallet123").Return(wallet, nil)
	mockRepo.On("SaveBalanceHistory", ctx, mock.AnythingOfType("*model.BalanceHistory")).Return(nil)

	// Test
	err := service.SaveBalanceHistory(ctx, "wallet123")

	// Assert
	require.NoError(t, err)

	// Verify mocks
	mockRepo.AssertExpectations(t)
}
