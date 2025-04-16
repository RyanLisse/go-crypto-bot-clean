package usecase

import (
	"context"
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

// MockWalletProvider is a mock implementation of the WalletProvider interface
type MockWalletProvider struct {
	mock.Mock
}

func (m *MockWalletProvider) GetName() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockWalletProvider) GetType() model.WalletType {
	args := m.Called()
	return model.WalletType(args.String(0))
}

func (m *MockWalletProvider) Connect(ctx context.Context, params map[string]interface{}) (*model.Wallet, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Wallet), args.Error(1)
}

func (m *MockWalletProvider) Disconnect(ctx context.Context, walletID string) error {
	args := m.Called(ctx, walletID)
	return args.Error(0)
}

func (m *MockWalletProvider) Verify(ctx context.Context, address string, message string, signature string) (bool, error) {
	args := m.Called(ctx, address, message, signature)
	return args.Bool(0), args.Error(1)
}

func (m *MockWalletProvider) GetBalance(ctx context.Context, wallet *model.Wallet) (*model.Wallet, error) {
	args := m.Called(ctx, wallet)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Wallet), args.Error(1)
}

func (m *MockWalletProvider) IsValidAddress(ctx context.Context, address string) (bool, error) {
	args := m.Called(ctx, address)
	return args.Bool(0), args.Error(1)
}

// SVMockWalletRepository is a mock implementation of the WalletRepository interface
type SVMockWalletRepository struct {
	mock.Mock
}

func (m *SVMockWalletRepository) Save(ctx context.Context, wallet *model.Wallet) error {
	args := m.Called(ctx, wallet)
	return args.Error(0)
}

func (m *SVMockWalletRepository) GetByID(ctx context.Context, id string) (*model.Wallet, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Wallet), args.Error(1)
}

func (m *SVMockWalletRepository) GetByUserID(ctx context.Context, userID string) (*model.Wallet, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Wallet), args.Error(1)
}

func (m *SVMockWalletRepository) GetWalletsByUserID(ctx context.Context, userID string) ([]*model.Wallet, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Wallet), args.Error(1)
}

func (m *SVMockWalletRepository) DeleteWallet(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *SVMockWalletRepository) SaveBalanceHistory(ctx context.Context, history *model.BalanceHistory) error {
	args := m.Called(ctx, history)
	return args.Error(0)
}

func (m *SVMockWalletRepository) GetBalanceHistory(ctx context.Context, userID string, asset model.Asset, from, to time.Time) ([]*model.BalanceHistory, error) {
	args := m.Called(ctx, userID, asset, from, to)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.BalanceHistory), args.Error(1)
}

// MockProviderRegistry is a mock implementation of the wallet.ProviderRegistry
type MockProviderRegistry struct {
	wallet.ProviderRegistry
	mock.Mock
}

func NewMockProviderRegistry() *MockProviderRegistry {
	return &MockProviderRegistry{
		ProviderRegistry: *wallet.NewProviderRegistry(),
	}
}

func (m *MockProviderRegistry) GetProvider(name string) (port.WalletProvider, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(port.WalletProvider), args.Error(1)
}

func TestSignatureVerificationService_GenerateChallenge(t *testing.T) {
	// Setup
	ctx := context.Background()
	logger := zerolog.New(zerolog.NewTestWriter(t))
	mockRepo := new(SVMockWalletRepository)
	mockRegistry := wallet.NewProviderRegistry()
	service := NewSignatureVerificationService(mockRegistry, mockRepo, &logger)

	// Test data
	walletID := "wallet123"
	userID := "user123"
	wallet := &model.Wallet{
		ID:       walletID,
		UserID:   userID,
		Exchange: "MEXC",
		Type:     model.WalletTypeExchange,
		Status:   model.WalletStatusActive,
		Metadata: &model.WalletMetadata{},
	}

	// Setup expectations
	mockRepo.On("GetByID", ctx, walletID).Return(wallet, nil)

	// Call the method
	challenge, err := service.GenerateChallenge(ctx, walletID)

	// Assertions
	require.NoError(t, err)
	assert.NotEmpty(t, challenge)
	assert.Contains(t, challenge, walletID)
	assert.Contains(t, challenge, "Sign this message")
	mockRepo.AssertExpectations(t)
}

func TestSignatureVerificationService_VerifySignature(t *testing.T) {
	// Setup
	ctx := context.Background()
	logger := zerolog.New(zerolog.NewTestWriter(t))
	mockRepo := new(SVMockWalletRepository)
	mockRegistry := wallet.NewProviderRegistry()
	mockProvider := new(MockWalletProvider)
	service := NewSignatureVerificationService(mockRegistry, mockRepo, &logger)

	// Test data
	walletID := "wallet123"
	userID := "user123"
	wallet := &model.Wallet{
		ID:       walletID,
		UserID:   userID,
		Exchange: "MEXC",
		Type:     model.WalletTypeExchange,
		Status:   model.WalletStatusActive,
		Metadata: &model.WalletMetadata{},
	}
	challenge := "Sign this message to verify your wallet ownership: abc123\nWallet ID: wallet123\nTimestamp: 1234567890"
	signature := "valid_signature"

	// Setup expectations
	mockRepo.On("GetByID", ctx, walletID).Return(wallet, nil)
	// Register the mock provider in the registry
	mockProvider.On("GetName").Return("MEXC")
	mockRegistry.RegisterProvider(mockProvider)
	mockProvider.On("Verify", ctx, "MEXC", challenge, signature).Return(true, nil)
	mockRepo.On("Save", ctx, mock.AnythingOfType("*model.Wallet")).Return(nil)

	// Generate a challenge first
	svc := service.(*signatureVerificationService)
	svc.challenges[walletID] = &Challenge{
		WalletID:  walletID,
		Message:   challenge,
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}

	// Call the method
	verified, err := service.VerifySignature(ctx, walletID, challenge, signature)

	// Assertions
	require.NoError(t, err)
	assert.True(t, verified)
	mockRepo.AssertExpectations(t)
	mockProvider.AssertExpectations(t)

	// Verify that the wallet status was updated
	// There should be at least one call to Save
	saveCallCount := 0
	for _, call := range mockRepo.Calls {
		if call.Method == "Save" {
			saveCallCount++
		}
	}
	assert.GreaterOrEqual(t, saveCallCount, 1)
	for _, call := range mockRepo.Calls {
		if call.Method == "Save" {
			wallet := call.Arguments.Get(1).(*model.Wallet)
			assert.Equal(t, model.WalletStatusVerified, wallet.Status)
		}
	}
}

func TestSignatureVerificationService_GetWalletStatus(t *testing.T) {
	// Setup
	ctx := context.Background()
	logger := zerolog.New(zerolog.NewTestWriter(t))
	mockRepo := new(SVMockWalletRepository)
	mockRegistry := wallet.NewProviderRegistry()
	service := NewSignatureVerificationService(mockRegistry, mockRepo, &logger)

	// Test data
	walletID := "wallet123"
	userID := "user123"
	wallet := &model.Wallet{
		ID:       walletID,
		UserID:   userID,
		Exchange: "MEXC",
		Type:     model.WalletTypeExchange,
		Status:   model.WalletStatusVerified,
		Metadata: &model.WalletMetadata{},
	}

	// Setup expectations
	mockRepo.On("GetByID", ctx, walletID).Return(wallet, nil)

	// Call the method
	status, err := service.GetWalletStatus(ctx, walletID)

	// Assertions
	require.NoError(t, err)
	assert.Equal(t, model.WalletStatusVerified, status)
	mockRepo.AssertExpectations(t)
}

func TestSignatureVerificationService_SetWalletStatus(t *testing.T) {
	// Setup
	ctx := context.Background()
	logger := zerolog.New(zerolog.NewTestWriter(t))
	mockRepo := new(SVMockWalletRepository)
	mockRegistry := wallet.NewProviderRegistry()
	service := NewSignatureVerificationService(mockRegistry, mockRepo, &logger)

	// Test data
	walletID := "wallet123"
	userID := "user123"
	wallet := &model.Wallet{
		ID:       walletID,
		UserID:   userID,
		Exchange: "MEXC",
		Type:     model.WalletTypeExchange,
		Status:   model.WalletStatusActive,
		Metadata: &model.WalletMetadata{},
	}

	// Setup expectations
	mockRepo.On("GetByID", ctx, walletID).Return(wallet, nil)
	mockRepo.On("Save", ctx, mock.AnythingOfType("*model.Wallet")).Return(nil)

	// Call the method
	err := service.SetWalletStatus(ctx, walletID, model.WalletStatusVerified)

	// Assertions
	require.NoError(t, err)
	mockRepo.AssertExpectations(t)

	// Verify that the wallet status was updated
	// There should be at least one call to Save
	saveCallCount := 0
	for _, call := range mockRepo.Calls {
		if call.Method == "Save" {
			saveCallCount++
		}
	}
	assert.GreaterOrEqual(t, saveCallCount, 1)
	for _, call := range mockRepo.Calls {
		if call.Method == "Save" {
			wallet := call.Arguments.Get(1).(*model.Wallet)
			assert.Equal(t, model.WalletStatusVerified, wallet.Status)
		}
	}
}
