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

// MockAPICredentialRepository is a mock implementation of port.APICredentialRepository
type MockAPICredentialRepository struct {
	mock.Mock
	// ... existing fields
}

func (m *MockAPICredentialRepository) ListAll(ctx context.Context) ([]*model.APICredential, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.APICredential), args.Error(1)
}

func (m *MockAPICredentialRepository) Save(ctx context.Context, credential *model.APICredential) error {
	args := m.Called(ctx, credential)
	return args.Error(0)
}

func (m *MockAPICredentialRepository) GetByID(ctx context.Context, id string) (*model.APICredential, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.APICredential), args.Error(1)
}

func (m *MockAPICredentialRepository) GetByUserIDAndExchange(ctx context.Context, userID, exchange string) (*model.APICredential, error) {
	args := m.Called(ctx, userID, exchange)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.APICredential), args.Error(1)
}

func (m *MockAPICredentialRepository) GetByUserIDAndLabel(ctx context.Context, userID, exchange, label string) (*model.APICredential, error) {
	args := m.Called(ctx, userID, exchange, label)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.APICredential), args.Error(1)
}

func (m *MockAPICredentialRepository) DeleteByID(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockAPICredentialRepository) ListByUserID(ctx context.Context, userID string) ([]*model.APICredential, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.APICredential), args.Error(1)
}

func (m *MockAPICredentialRepository) UpdateStatus(ctx context.Context, id string, status model.APICredentialStatus) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockAPICredentialRepository) UpdateLastUsed(ctx context.Context, id string, lastUsed time.Time) error {
	args := m.Called(ctx, id, lastUsed)
	return args.Error(0)
}

func (m *MockAPICredentialRepository) UpdateLastVerified(ctx context.Context, id string, lastVerified time.Time) error {
	args := m.Called(ctx, id, lastVerified)
	return args.Error(0)
}

func (m *MockAPICredentialRepository) IncrementFailureCount(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockAPICredentialRepository) ResetFailureCount(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockEncryptionService is a mock implementation of crypto.EncryptionService
type MockEncryptionService struct {
	mock.Mock
}

func (m *MockEncryptionService) Encrypt(plaintext string) ([]byte, error) {
	args := m.Called(plaintext)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockEncryptionService) Decrypt(ciphertext []byte) (string, error) {
	args := m.Called(ciphertext)
	return args.String(0), args.Error(1)
}

// MockExchangeWalletProvider is a mock implementation of port.ExchangeWalletProvider
type MockExchangeWalletProvider struct {
	mock.Mock
}

func (m *MockExchangeWalletProvider) GetName() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockExchangeWalletProvider) GetType() model.WalletType {
	args := m.Called()
	return args.Get(0).(model.WalletType)
}

func (m *MockExchangeWalletProvider) Connect(ctx context.Context, params map[string]any) (*model.Wallet, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Wallet), args.Error(1)
}

func (m *MockExchangeWalletProvider) Disconnect(ctx context.Context, walletID string) error {
	args := m.Called(ctx, walletID)
	return args.Error(0)
}

func (m *MockExchangeWalletProvider) Verify(ctx context.Context, address string, message string, signature string) (bool, error) {
	args := m.Called(ctx, address, message, signature)
	return args.Bool(0), args.Error(1)
}

func (m *MockExchangeWalletProvider) GetBalance(ctx context.Context, wallet *model.Wallet) (*model.Wallet, error) {
	args := m.Called(ctx, wallet)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Wallet), args.Error(1)
}

func (m *MockExchangeWalletProvider) IsValidAddress(ctx context.Context, address string) (bool, error) {
	args := m.Called(ctx, address)
	return args.Bool(0), args.Error(1)
}

func (m *MockExchangeWalletProvider) SetAPICredentials(ctx context.Context, apiKey, apiSecret string) error {
	args := m.Called(ctx, apiKey, apiSecret)
	return args.Error(0)
}

// ACMockProviderRegistry is a mock implementation of port.ProviderRegistry
type ACMockProviderRegistry struct {
	ProviderRegistry *wallet.ProviderRegistry
	mock.Mock
}

// No need to mock RegisterProvider as we're using the real implementation

// No need to mock GetProvider as we're using the real implementation

// No need to mock GetExchangeProvider as we're using the real implementation

// No need to mock GetWeb3Provider as we're using the real implementation

// No need to mock GetProviderByType as we're using the real implementation

// No need to mock GetAllProviders as we're using the real implementation

// No need to mock GetAllExchangeProviders as we're using the real implementation

// No need to mock GetAllWeb3Providers as we're using the real implementation

func TestCreateCredential(t *testing.T) {
	// Setup
	ctx := context.Background()
	logger := zerolog.New(zerolog.NewTestWriter(t))
	mockRepo := new(MockAPICredentialRepository)
	mockEncryption := new(MockEncryptionService)
	mockProvider := new(MockExchangeWalletProvider)
	providerRegistry := wallet.NewProviderRegistry()
	// Create a mock registry (not used, just for compatibility)

	// Setup mock provider
	mockProvider.On("GetName").Return("MEXC")
	mockProvider.On("SetAPICredentials", mock.Anything, "api_key", "api_secret").Return(nil)
	mockProvider.On("GetBalance", mock.Anything, mock.AnythingOfType("*model.Wallet")).Return(&model.Wallet{}, nil)

	// Register the mock provider in the registry
	mockProvider.On("GetName").Return("MEXC")
	providerRegistry.RegisterProvider(mockProvider)

	// Setup mock repository
	mockRepo.On("GetByUserIDAndExchange", ctx, "user123", "MEXC").Return(nil, model.ErrCredentialNotFound)
	mockRepo.On("Save", ctx, mock.AnythingOfType("*model.APICredential")).Return(nil)

	// Create service
	service := NewAPICredentialManagerService(mockRepo, mockEncryption, providerRegistry, &logger)

	// Test
	credential, err := service.CreateCredential(ctx, "user123", "MEXC", "api_key", "api_secret", "")

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, credential)
	assert.Equal(t, "user123", credential.UserID)
	assert.Equal(t, "MEXC", credential.Exchange)
	assert.Equal(t, "api_key", credential.APIKey)
	assert.Equal(t, "api_secret", credential.APISecret)
	assert.Equal(t, model.APICredentialStatusActive, credential.Status)

	// Verify mocks
	mockRepo.AssertExpectations(t)
	mockProvider.AssertExpectations(t)
}

func TestCreateCredentialWithExistingLabel(t *testing.T) {
	// Setup
	ctx := context.Background()
	logger := zerolog.New(zerolog.NewTestWriter(t))
	mockRepo := new(MockAPICredentialRepository)
	mockEncryption := new(MockEncryptionService)
	mockProvider := new(MockExchangeWalletProvider)
	providerRegistry := wallet.NewProviderRegistry()

	// Setup mock provider
	mockProvider.On("GetName").Return("MEXC")

	// Register the mock provider in the registry
	mockProvider.On("GetName").Return("MEXC")
	providerRegistry.RegisterProvider(mockProvider)

	// Setup mock repository
	existingCred := &model.APICredential{
		ID:        "cred123",
		UserID:    "user123",
		Exchange:  "MEXC",
		APIKey:    "old_key",
		APISecret: "old_secret",
		Label:     "main",
		Status:    model.APICredentialStatusActive,
	}
	mockRepo.On("GetByUserIDAndExchange", ctx, "user123", "MEXC").Return(existingCred, nil)
	mockRepo.On("GetByUserIDAndLabel", ctx, "user123", "MEXC", "main").Return(existingCred, nil)

	// Create service
	service := NewAPICredentialManagerService(mockRepo, mockEncryption, providerRegistry, &logger)

	// Test
	_, err := service.CreateCredential(ctx, "user123", "MEXC", "api_key", "api_secret", "main")

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "credential with this label already exists")

	// Verify mocks
	mockRepo.AssertExpectations(t)
}

func TestVerifyCredential(t *testing.T) {
	// Setup
	ctx := context.Background()
	logger := zerolog.New(zerolog.NewTestWriter(t))
	mockRepo := new(MockAPICredentialRepository)
	mockEncryption := new(MockEncryptionService)
	mockProvider := new(MockExchangeWalletProvider)
	providerRegistry := wallet.NewProviderRegistry()

	// Setup mock provider
	mockProvider.On("GetName").Return("MEXC")
	mockProvider.On("SetAPICredentials", mock.Anything, "api_key", "api_secret").Return(nil)
	mockProvider.On("GetBalance", mock.Anything, mock.AnythingOfType("*model.Wallet")).Return(&model.Wallet{}, nil)

	// Register the mock provider in the registry
	mockProvider.On("GetName").Return("MEXC")
	providerRegistry.RegisterProvider(mockProvider)

	// Setup mock repository
	credential := &model.APICredential{
		ID:        "cred123",
		UserID:    "user123",
		Exchange:  "MEXC",
		APIKey:    "api_key",
		APISecret: "api_secret",
		Status:    model.APICredentialStatusActive,
	}
	mockRepo.On("GetByID", ctx, "cred123").Return(credential, nil)
	mockRepo.On("Save", ctx, mock.AnythingOfType("*model.APICredential")).Return(nil)

	// Create service
	service := NewAPICredentialManagerService(mockRepo, mockEncryption, providerRegistry, &logger)

	// Test
	verified, err := service.VerifyCredential(ctx, "cred123")

	// Assert
	require.NoError(t, err)
	assert.True(t, verified)

	// Verify mocks
	mockRepo.AssertExpectations(t)
	mockProvider.AssertExpectations(t)
}

func TestVerifyCredentialFailure(t *testing.T) {
	// Setup
	ctx := context.Background()
	logger := zerolog.New(zerolog.NewTestWriter(t))
	mockRepo := new(MockAPICredentialRepository)
	mockEncryption := new(MockEncryptionService)
	mockProvider := new(MockExchangeWalletProvider)
	providerRegistry := wallet.NewProviderRegistry()

	// Setup mock provider
	mockProvider.On("GetName").Return("MEXC")
	mockProvider.On("SetAPICredentials", mock.Anything, "api_key", "api_secret").Return(nil)
	mockProvider.On("GetBalance", mock.Anything, mock.AnythingOfType("*model.Wallet")).Return(nil, errors.New("invalid credentials"))

	// Register the mock provider in the registry
	mockProvider.On("GetName").Return("MEXC")
	providerRegistry.RegisterProvider(mockProvider)

	// Setup mock repository
	credential := &model.APICredential{
		ID:        "cred123",
		UserID:    "user123",
		Exchange:  "MEXC",
		APIKey:    "api_key",
		APISecret: "api_secret",
		Status:    model.APICredentialStatusActive,
	}
	mockRepo.On("GetByID", ctx, "cred123").Return(credential, nil)
	mockRepo.On("Save", ctx, mock.AnythingOfType("*model.APICredential")).Return(nil)

	// Create service
	service := NewAPICredentialManagerService(mockRepo, mockEncryption, providerRegistry, &logger)

	// Test
	verified, err := service.VerifyCredential(ctx, "cred123")

	// Assert
	require.Error(t, err)
	assert.False(t, verified)
	assert.Contains(t, err.Error(), "invalid credentials")

	// Verify mocks
	mockRepo.AssertExpectations(t)
	mockProvider.AssertExpectations(t)
}

func TestGetCredentialForExchange(t *testing.T) {
	// Setup
	ctx := context.Background()
	logger := zerolog.New(zerolog.NewTestWriter(t))
	mockRepo := new(MockAPICredentialRepository)
	mockEncryption := new(MockEncryptionService)
	mockProvider := new(MockExchangeWalletProvider)
	providerRegistry := wallet.NewProviderRegistry()

	// Setup mock provider
	mockProvider.On("GetName").Return("MEXC")

	// Register the mock provider in the registry
	mockProvider.On("GetName").Return("MEXC")
	providerRegistry.RegisterProvider(mockProvider)

	// Setup mock repository
	credential := &model.APICredential{
		ID:        "cred123",
		UserID:    "user123",
		Exchange:  "MEXC",
		APIKey:    "api_key",
		APISecret: "api_secret",
		Status:    model.APICredentialStatusActive,
	}
	mockRepo.On("GetByUserIDAndExchange", ctx, "user123", "MEXC").Return(credential, nil)
	mockRepo.On("UpdateLastUsed", ctx, "cred123", mock.AnythingOfType("time.Time")).Return(nil)

	// Create service
	service := NewAPICredentialManagerService(mockRepo, mockEncryption, providerRegistry, &logger)

	// Test
	result, err := service.GetCredentialForExchange(ctx, "user123", "MEXC")

	// Assert
	require.NoError(t, err)
	assert.Equal(t, credential, result)

	// Verify mocks
	mockRepo.AssertExpectations(t)
}

func TestGetCredentialForExchangeInactive(t *testing.T) {
	// Setup
	ctx := context.Background()
	logger := zerolog.New(zerolog.NewTestWriter(t))
	mockRepo := new(MockAPICredentialRepository)
	mockEncryption := new(MockEncryptionService)
	mockProvider := new(MockExchangeWalletProvider)
	providerRegistry := wallet.NewProviderRegistry()

	// Setup mock provider
	mockProvider.On("GetName").Return("MEXC")

	// Register the mock provider in the registry
	mockProvider.On("GetName").Return("MEXC")
	providerRegistry.RegisterProvider(mockProvider)

	// Setup mock repository
	credential := &model.APICredential{
		ID:        "cred123",
		UserID:    "user123",
		Exchange:  "MEXC",
		APIKey:    "api_key",
		APISecret: "api_secret",
		Status:    model.APICredentialStatusInactive,
	}
	mockRepo.On("GetByUserIDAndExchange", ctx, "user123", "MEXC").Return(credential, nil)

	// Create service
	service := NewAPICredentialManagerService(mockRepo, mockEncryption, providerRegistry, &logger)

	// Test
	_, err := service.GetCredentialForExchange(ctx, "user123", "MEXC")

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "credential is not active")

	// Verify mocks
	mockRepo.AssertExpectations(t)
}
