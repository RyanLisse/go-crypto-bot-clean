package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/wallet"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockWeb3WalletProvider is a mock implementation of the Web3WalletProvider interface
type MockWeb3WalletProvider struct {
	mock.Mock
}

func (m *MockWeb3WalletProvider) GetName() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockWeb3WalletProvider) GetType() model.WalletType {
	args := m.Called()
	return model.WalletType(args.String(0))
}

func (m *MockWeb3WalletProvider) GetChainID() int64 {
	args := m.Called()
	return args.Get(0).(int64)
}

func (m *MockWeb3WalletProvider) GetNetwork() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockWeb3WalletProvider) Connect(ctx context.Context, params map[string]interface{}) (*model.Wallet, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Wallet), args.Error(1)
}

func (m *MockWeb3WalletProvider) Disconnect(ctx context.Context, walletID string) error {
	args := m.Called(ctx, walletID)
	return args.Error(0)
}

func (m *MockWeb3WalletProvider) Verify(ctx context.Context, address string, message string, signature string) (bool, error) {
	args := m.Called(ctx, address, message, signature)
	return args.Bool(0), args.Error(1)
}

func (m *MockWeb3WalletProvider) GetBalance(ctx context.Context, wallet *model.Wallet) (*model.Wallet, error) {
	args := m.Called(ctx, wallet)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Wallet), args.Error(1)
}

func (m *MockWeb3WalletProvider) IsValidAddress(ctx context.Context, address string) (bool, error) {
	args := m.Called(ctx, address)
	return args.Bool(0), args.Error(1)
}

// SignMessage implements port.Web3WalletProvider
func (m *MockWeb3WalletProvider) SignMessage(ctx context.Context, message string) (string, error) {
	args := m.Called(ctx, message)
	return args.String(0), args.Error(1)
}

// MockWeb3ProviderRegistry is a mock implementation of the ProviderRegistry
type MockWeb3ProviderRegistry struct {
	mock.Mock
}

func (m *MockWeb3ProviderRegistry) RegisterProvider(provider port.WalletProvider) {
	m.Called(provider)
}

func (m *MockWeb3ProviderRegistry) GetProvider(name string) (port.WalletProvider, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(port.WalletProvider), args.Error(1)
}

func (m *MockWeb3ProviderRegistry) GetExchangeProvider(name string) (port.ExchangeWalletProvider, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(port.ExchangeWalletProvider), args.Error(1)
}

func (m *MockWeb3ProviderRegistry) GetWeb3Provider(name string) (port.Web3WalletProvider, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(port.Web3WalletProvider), args.Error(1)
}

func (m *MockWeb3ProviderRegistry) GetProviderByType(typ model.WalletType) ([]port.WalletProvider, error) {
	args := m.Called(typ)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]port.WalletProvider), args.Error(1)
}

func (m *MockWeb3ProviderRegistry) GetAllProviders() []port.WalletProvider {
	args := m.Called()
	return args.Get(0).([]port.WalletProvider)
}

func (m *MockWeb3ProviderRegistry) GetAllExchangeProviders() []port.ExchangeWalletProvider {
	args := m.Called()
	return args.Get(0).([]port.ExchangeWalletProvider)
}

func (m *MockWeb3ProviderRegistry) GetAllWeb3Providers() []port.Web3WalletProvider {
	args := m.Called()
	return args.Get(0).([]port.Web3WalletProvider)
}

// MockWeb3WalletRepository is a mock implementation of the WalletRepository interface
type MockWeb3WalletRepository struct {
	mock.Mock
}

func (m *MockWeb3WalletRepository) Save(ctx context.Context, wallet *model.Wallet) error {
	args := m.Called(ctx, wallet)
	return args.Error(0)
}

func (m *MockWeb3WalletRepository) GetByID(ctx context.Context, id string) (*model.Wallet, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Wallet), args.Error(1)
}

func (m *MockWeb3WalletRepository) GetByUserID(ctx context.Context, userID string) (*model.Wallet, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Wallet), args.Error(1)
}

func (m *MockWeb3WalletRepository) GetWalletsByUserID(ctx context.Context, userID string) ([]*model.Wallet, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Wallet), args.Error(1)
}

func (m *MockWeb3WalletRepository) DeleteWallet(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockWeb3WalletRepository) SaveBalanceHistory(ctx context.Context, history *model.BalanceHistory) error {
	args := m.Called(ctx, history)
	return args.Error(0)
}

func (m *MockWeb3WalletRepository) GetBalanceHistory(ctx context.Context, userID string, asset model.Asset, from, to time.Time) ([]*model.BalanceHistory, error) {
	args := m.Called(ctx, userID, asset, from, to)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.BalanceHistory), args.Error(1)
}

func TestWeb3WalletService_ConnectWallet(t *testing.T) {
	// Setup
	ctx := context.Background()
	logger := zerolog.New(zerolog.NewTestWriter(t))
	mockRepo := new(MockWeb3WalletRepository)
	realRegistry := wallet.NewProviderRegistry()
	mockProvider := new(MockWeb3WalletProvider)
	realRegistry.RegisterProvider(mockProvider)
	service := NewWeb3WalletService(mockRepo, realRegistry, &logger)

	// Test data
	userID := "user123"
	network := "Ethereum"
	address := "0x742d35Cc6634C0532925a3b844Bc454e4438f44e"
	walletID := "wallet123"
	wallet := &model.Wallet{
		ID:       walletID,
		UserID:   userID,
		Type:     model.WalletTypeWeb3,
		Status:   model.WalletStatusActive,
		Metadata: &model.WalletMetadata{
			Network: network,
			Address: address,
			ChainID: 1,
		},
	}

	// Setup expectations
	mockRegistry := new(MockWeb3ProviderRegistry)
	mockRegistry.On("GetWeb3Provider", network).Return(mockProvider, nil)
	mockProvider.On("IsValidAddress", ctx, address).Return(true, nil)
	mockRepo.On("GetWalletsByUserID", ctx, "").Return([]*model.Wallet{}, nil)
	mockProvider.On("Connect", ctx, mock.Anything).Return(wallet, nil)
	mockRepo.On("Save", ctx, wallet).Return(nil)

	// Call the method
	result, err := service.ConnectWallet(ctx, userID, network, address)

	// Assertions
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, walletID, result.ID)
	assert.Equal(t, userID, result.UserID)
	assert.Equal(t, model.WalletTypeWeb3, result.Type)
	assert.Equal(t, network, result.Metadata.Network)
	assert.Equal(t, address, result.Metadata.Address)
	
	mockProvider.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestWeb3WalletService_ConnectWallet_InvalidAddress(t *testing.T) {
	// Setup
	ctx := context.Background()
	logger := zerolog.New(zerolog.NewTestWriter(t))
	mockRepo := new(MockWeb3WalletRepository)
	realRegistry := wallet.NewProviderRegistry()
	mockProvider := new(MockWeb3WalletProvider)
	realRegistry.RegisterProvider(mockProvider)
	service := NewWeb3WalletService(mockRepo, realRegistry, &logger)

	// Test data
	userID := "user123"
	network := "Ethereum"
	address := "invalid_address"

	// Setup expectations
	mockRegistry := new(MockWeb3ProviderRegistry)
	mockRegistry.On("GetWeb3Provider", network).Return(mockProvider, nil)
	mockProvider.On("IsValidAddress", ctx, address).Return(false, nil)

	// Call the method
	result, err := service.ConnectWallet(ctx, userID, network, address)

	// Assertions
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "invalid address", err.Error())
	
	mockProvider.AssertExpectations(t)
}

func TestWeb3WalletService_ConnectWallet_UnsupportedNetwork(t *testing.T) {
	// Setup
	ctx := context.Background()
	logger := zerolog.New(zerolog.NewTestWriter(t))
	mockRepo := new(MockWeb3WalletRepository)
	realRegistry := wallet.NewProviderRegistry()
	mockProvider := new(MockWeb3WalletProvider)
	realRegistry.RegisterProvider(mockProvider)
	service := NewWeb3WalletService(mockRepo, realRegistry, &logger)

	// Test data
	userID := "user123"
	network := "UnsupportedNetwork"
	address := "0x742d35Cc6634C0532925a3b844Bc454e4438f44e"

	// Setup expectations
	mockRegistry := new(MockWeb3ProviderRegistry)
	mockRegistry.On("GetWeb3Provider", network).Return(nil, errors.New("unsupported network"))

	// Call the method
	result, err := service.ConnectWallet(ctx, userID, network, address)

	// Assertions
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "unsupported network: UnsupportedNetwork", err.Error())
	
}

func TestWeb3WalletService_DisconnectWallet(t *testing.T) {
	// Setup
	ctx := context.Background()
	logger := zerolog.New(zerolog.NewTestWriter(t))
	mockRepo := new(MockWeb3WalletRepository)
	realRegistry := wallet.NewProviderRegistry()
	mockProvider := new(MockWeb3WalletProvider)
	realRegistry.RegisterProvider(mockProvider)
	service := NewWeb3WalletService(mockRepo, realRegistry, &logger)

	// Test data
	walletID := "wallet123"
	userID := "user123"
	network := "Ethereum"
	address := "0x742d35Cc6634C0532925a3b844Bc454e4438f44e"
	wallet := &model.Wallet{
		ID:       walletID,
		UserID:   userID,
		Type:     model.WalletTypeWeb3,
		Status:   model.WalletStatusActive,
		Metadata: &model.WalletMetadata{
			Network: network,
			Address: address,
			ChainID: 1,
		},
	}

	// Setup expectations
	mockRepo.On("GetByID", ctx, walletID).Return(wallet, nil)
	mockRegistry := new(MockWeb3ProviderRegistry)
	mockRegistry.On("GetWeb3Provider", network).Return(mockProvider, nil)
	mockProvider.On("Disconnect", ctx, walletID).Return(nil)
	mockRepo.On("DeleteWallet", ctx, walletID).Return(nil)

	// Call the method
	err := service.DisconnectWallet(ctx, walletID)

	// Assertions
	require.NoError(t, err)
	mockRepo.AssertExpectations(t)
	
	mockProvider.AssertExpectations(t)
}

func TestWeb3WalletService_GetWalletBalance(t *testing.T) {
	// Setup
	ctx := context.Background()
	logger := zerolog.New(zerolog.NewTestWriter(t))
	mockRepo := new(MockWeb3WalletRepository)
	realRegistry := wallet.NewProviderRegistry()
	mockProvider := new(MockWeb3WalletProvider)
	realRegistry.RegisterProvider(mockProvider)
	service := NewWeb3WalletService(mockRepo, realRegistry, &logger)

	// Test data
	walletID := "wallet123"
	userID := "user123"
	network := "Ethereum"
	address := "0x742d35Cc6634C0532925a3b844Bc454e4438f44e"
	wallet := &model.Wallet{
		ID:       walletID,
		UserID:   userID,
		Type:     model.WalletTypeWeb3,
		Status:   model.WalletStatusActive,
		Metadata: &model.WalletMetadata{
			Network: network,
			Address: address,
			ChainID: 1,
		},
		Balances: make(map[model.Asset]*model.Balance),
	}
	updatedWallet := &model.Wallet{
		ID:       walletID,
		UserID:   userID,
		Type:     model.WalletTypeWeb3,
		Status:   model.WalletStatusActive,
		Metadata: &model.WalletMetadata{
			Network: network,
			Address: address,
			ChainID: 1,
		},
		Balances: map[model.Asset]*model.Balance{
			model.AssetETH: {
				Asset:    model.AssetETH,
				Free:     1.0,
				Locked:   0.0,
				Total:    1.0,
				USDValue: 2000.0,
			},
		},
		TotalUSDValue: 2000.0,
	}

	// Setup expectations
	mockRepo.On("GetByID", ctx, walletID).Return(wallet, nil)
	mockRegistry := new(MockWeb3ProviderRegistry)
	mockRegistry.On("GetWeb3Provider", network).Return(mockProvider, nil)
	mockProvider.On("GetBalance", ctx, wallet).Return(updatedWallet, nil)
	mockRepo.On("Save", ctx, mock.AnythingOfType("*model.Wallet")).Return(nil)

	// Call the method
	result, err := service.GetWalletBalance(ctx, walletID)

	// Assertions
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, walletID, result.ID)
	assert.Equal(t, userID, result.UserID)
	assert.Equal(t, model.WalletTypeWeb3, result.Type)
	assert.Equal(t, network, result.Metadata.Network)
	assert.Equal(t, address, result.Metadata.Address)
	assert.Equal(t, 2000.0, result.TotalUSDValue)
	assert.NotNil(t, result.Balances[model.AssetETH])
	assert.Equal(t, 1.0, result.Balances[model.AssetETH].Total)
	mockRepo.AssertExpectations(t)
	
	mockProvider.AssertExpectations(t)
}

func TestWeb3WalletService_IsValidAddress(t *testing.T) {
	// Setup
	ctx := context.Background()
	logger := zerolog.New(zerolog.NewTestWriter(t))
	mockRepo := new(MockWeb3WalletRepository)
	realRegistry := wallet.NewProviderRegistry()
	mockProvider := new(MockWeb3WalletProvider)
	realRegistry.RegisterProvider(mockProvider)
	service := NewWeb3WalletService(mockRepo, realRegistry, &logger)

	// Test data
	network := "Ethereum"
	validAddress := "0x742d35Cc6634C0532925a3b844Bc454e4438f44e"
	invalidAddress := "invalid_address"

	// Setup expectations for valid address
	mockProvider.On("IsValidAddress", ctx, validAddress).Return(true, nil)
	mockProvider.On("IsValidAddress", ctx, invalidAddress).Return(false, nil)

	// Test valid address
	valid, err := service.IsValidAddress(ctx, network, validAddress)
	require.NoError(t, err)
	assert.True(t, valid)

	// Test invalid address
	valid, err = service.IsValidAddress(ctx, network, invalidAddress)
	require.NoError(t, err)
	assert.False(t, valid)

	
	mockProvider.AssertExpectations(t)
}

func TestWeb3WalletService_GetSupportedNetworks(t *testing.T) {
	// Setup
	ctx := context.Background()
	logger := zerolog.New(zerolog.NewTestWriter(t))
	mockRepo := new(MockWeb3WalletRepository)
	realRegistry := wallet.NewProviderRegistry()
	mockProvider1 := new(MockWeb3WalletProvider)
	mockProvider2 := new(MockWeb3WalletProvider)
	realRegistry.RegisterProvider(mockProvider1)
	realRegistry.RegisterProvider(mockProvider2)
	service := NewWeb3WalletService(mockRepo, realRegistry, &logger)

	// Test data
	mockProvider1.On("GetName").Return("Ethereum")
	mockProvider2.On("GetName").Return("Polygon")

	// Call the method
	networks, err := service.GetSupportedNetworks(ctx)

	// Assertions
	require.NoError(t, err)
	assert.Len(t, networks, 2)
	assert.Contains(t, networks, "Ethereum")
	assert.Contains(t, networks, "Polygon")
	
	mockProvider1.AssertExpectations(t)
	mockProvider2.AssertExpectations(t)
}
