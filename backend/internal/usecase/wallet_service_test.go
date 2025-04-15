package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockWalletRepository is a mock implementation of the WalletRepository interface
type MockWalletRepository struct {
	mock.Mock
}

func (m *MockWalletRepository) Save(ctx context.Context, wallet *model.Wallet) error {
	args := m.Called(ctx, wallet)
	return args.Error(0)
}

func (m *MockWalletRepository) GetByID(ctx context.Context, id string) (*model.Wallet, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Wallet), args.Error(1)
}

func (m *MockWalletRepository) GetByUserID(ctx context.Context, userID string) (*model.Wallet, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Wallet), args.Error(1)
}

func (m *MockWalletRepository) GetWalletsByUserID(ctx context.Context, userID string) ([]*model.Wallet, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Wallet), args.Error(1)
}

func (m *MockWalletRepository) DeleteWallet(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockWalletRepository) SaveBalanceHistory(ctx context.Context, history *model.BalanceHistory) error {
	args := m.Called(ctx, history)
	return args.Error(0)
}

func (m *MockWalletRepository) GetBalanceHistory(ctx context.Context, userID string, asset model.Asset, from, to time.Time) ([]*model.BalanceHistory, error) {
	args := m.Called(ctx, userID, asset, from, to)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.BalanceHistory), args.Error(1)
}

// WalletMockMEXCClient is a mock implementation of the MEXCClient interface for wallet tests
type WalletMockMEXCClient struct {
	mock.Mock
}

// GetAccount implements the MEXCClient interface
func (m *WalletMockMEXCClient) GetAccount(ctx context.Context) (*model.Wallet, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Wallet), args.Error(1)
}

// GetNewListings implements the MEXCClient interface
func (m *WalletMockMEXCClient) GetNewListings(ctx context.Context) ([]*model.NewCoin, error) {
	return nil, nil
}

// GetSymbolInfo implements the MEXCClient interface
func (m *WalletMockMEXCClient) GetSymbolInfo(ctx context.Context, symbol string) (*model.SymbolInfo, error) {
	return nil, nil
}

// GetSymbolStatus implements the MEXCClient interface
func (m *WalletMockMEXCClient) GetSymbolStatus(ctx context.Context, symbol string) (model.Status, error) {
	return "", nil
}

// GetTradingSchedule implements the MEXCClient interface
func (m *WalletMockMEXCClient) GetTradingSchedule(ctx context.Context, symbol string) (model.TradingSchedule, error) {
	return model.TradingSchedule{}, nil
}

// GetSymbolConstraints implements the MEXCClient interface
func (m *WalletMockMEXCClient) GetSymbolConstraints(ctx context.Context, symbol string) (*model.SymbolConstraints, error) {
	return nil, nil
}

// GetExchangeInfo implements the MEXCClient interface
func (m *WalletMockMEXCClient) GetExchangeInfo(ctx context.Context) (*model.ExchangeInfo, error) {
	return nil, nil
}

// GetMarketData implements the MEXCClient interface
func (m *WalletMockMEXCClient) GetMarketData(ctx context.Context, symbol string) (*model.Ticker, error) {
	return nil, nil
}

// GetKlines implements the MEXCClient interface
func (m *WalletMockMEXCClient) GetKlines(ctx context.Context, symbol string, interval model.KlineInterval, limit int) ([]*model.Kline, error) {
	return nil, nil
}

// GetOrderBook implements the MEXCClient interface
func (m *WalletMockMEXCClient) GetOrderBook(ctx context.Context, symbol string, depth int) (*model.OrderBook, error) {
	return nil, nil
}

// PlaceOrder implements the MEXCClient interface
func (m *WalletMockMEXCClient) PlaceOrder(ctx context.Context, symbol string, side model.OrderSide, orderType model.OrderType, quantity float64, price float64, timeInForce model.TimeInForce) (*model.Order, error) {
	return nil, nil
}

// CancelOrder implements the MEXCClient interface
func (m *WalletMockMEXCClient) CancelOrder(ctx context.Context, symbol string, orderID string) error {
	return nil
}

// GetOrderStatus implements the MEXCClient interface
func (m *WalletMockMEXCClient) GetOrderStatus(ctx context.Context, symbol string, orderID string) (*model.Order, error) {
	return nil, nil
}

// GetOpenOrders implements the MEXCClient interface
func (m *WalletMockMEXCClient) GetOpenOrders(ctx context.Context, symbol string) ([]*model.Order, error) {
	return nil, nil
}

// GetOrderHistory implements the MEXCClient interface
func (m *WalletMockMEXCClient) GetOrderHistory(ctx context.Context, symbol string, limit, offset int) ([]*model.Order, error) {
	return nil, nil
}

func TestWalletService_SetPrimaryWallet(t *testing.T) {
	// Setup
	ctx := context.Background()
	logger := zerolog.New(zerolog.NewTestWriter(t))
	mockRepo := new(MockWalletRepository)
	mockClient := new(WalletMockMEXCClient)
	service := NewWalletService(mockRepo, mockClient, &logger)

	// Test data
	userID := "user123"
	walletID1 := "wallet1"
	walletID2 := "wallet2"
	walletID3 := "wallet3"

	// Create test wallets
	wallet1 := &model.Wallet{
		ID:       walletID1,
		UserID:   userID,
		Exchange: "MEXC",
		Type:     model.WalletTypeExchange,
		Status:   model.WalletStatusActive,
		Metadata: &model.WalletMetadata{
			IsPrimary: true,
		},
	}

	wallet2 := &model.Wallet{
		ID:       walletID2,
		UserID:   userID,
		Exchange: "Binance",
		Type:     model.WalletTypeExchange,
		Status:   model.WalletStatusActive,
		Metadata: &model.WalletMetadata{
			IsPrimary: false,
		},
	}

	wallet3 := &model.Wallet{
		ID:     walletID3,
		UserID: userID,
		Type:   model.WalletTypeWeb3,
		Status: model.WalletStatusActive,
		Metadata: &model.WalletMetadata{
			IsPrimary: false,
			Network:   "Ethereum",
			Address:   "0x1234567890abcdef",
		},
	}

	wallets := []*model.Wallet{wallet1, wallet2, wallet3}

	// Setup expectations
	mockRepo.On("GetWalletsByUserID", ctx, userID).Return(wallets, nil)
	mockRepo.On("Save", ctx, mock.AnythingOfType("*model.Wallet")).Return(nil).Times(3)

	// Call the method
	err := service.SetPrimaryWallet(ctx, userID, walletID2)

	// Assertions
	require.NoError(t, err)
	mockRepo.AssertExpectations(t)

	// Verify that the primary wallet was updated
	for _, call := range mockRepo.Calls {
		if call.Method == "Save" {
			wallet := call.Arguments.Get(1).(*model.Wallet)
			if wallet.ID == walletID1 {
				assert.False(t, wallet.Metadata.IsPrimary, "Wallet 1 should not be primary")
			} else if wallet.ID == walletID2 {
				assert.True(t, wallet.Metadata.IsPrimary, "Wallet 2 should be primary")
			} else if wallet.ID == walletID3 {
				assert.False(t, wallet.Metadata.IsPrimary, "Wallet 3 should not be primary")
			}
		}
	}
}

func TestWalletService_SetPrimaryWallet_WalletNotFound(t *testing.T) {
	// Setup
	ctx := context.Background()
	logger := zerolog.New(zerolog.NewTestWriter(t))
	mockRepo := new(MockWalletRepository)
	mockClient := new(WalletMockMEXCClient)
	service := NewWalletService(mockRepo, mockClient, &logger)

	// Test data
	userID := "user123"
	walletID1 := "wallet1"
	walletID2 := "wallet2"
	nonExistentWalletID := "wallet999"

	// Create test wallets
	wallet1 := &model.Wallet{
		ID:       walletID1,
		UserID:   userID,
		Exchange: "MEXC",
		Type:     model.WalletTypeExchange,
		Status:   model.WalletStatusActive,
		Metadata: &model.WalletMetadata{
			IsPrimary: true,
		},
	}

	wallet2 := &model.Wallet{
		ID:       walletID2,
		UserID:   userID,
		Exchange: "Binance",
		Type:     model.WalletTypeExchange,
		Status:   model.WalletStatusActive,
		Metadata: &model.WalletMetadata{
			IsPrimary: false,
		},
	}

	wallets := []*model.Wallet{wallet1, wallet2}

	// Setup expectations
	mockRepo.On("GetWalletsByUserID", ctx, userID).Return(wallets, nil)

	// Call the method
	err := service.SetPrimaryWallet(ctx, userID, nonExistentWalletID)

	// Assertions
	require.Error(t, err)
	assert.Contains(t, err.Error(), "wallet not found")
	mockRepo.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "Save")
}

func TestWalletService_SetWalletMetadata(t *testing.T) {
	// Setup
	ctx := context.Background()
	logger := zerolog.New(zerolog.NewTestWriter(t))
	mockRepo := new(MockWalletRepository)
	mockClient := new(WalletMockMEXCClient)
	service := NewWalletService(mockRepo, mockClient, &logger)

	// Test data
	walletID := "wallet123"
	userID := "user123"
	name := "My Primary Wallet"
	description := "This is my main trading wallet"
	tags := []string{"trading", "main", "exchange"}

	// Create test wallet
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
	// The wallet should be saved with updated metadata
	mockRepo.On("Save", ctx, mock.MatchedBy(func(w *model.Wallet) bool {
		return w.ID == walletID &&
			w.Metadata.Name == name &&
			w.Metadata.Description == description &&
			len(w.Metadata.Tags) == len(tags)
	})).Return(nil)

	// Call the method
	err := service.SetWalletMetadata(ctx, walletID, name, description, tags)

	// Assertions
	require.NoError(t, err)
	mockRepo.AssertExpectations(t)

	// Verify that the metadata was updated
	assert.Equal(t, name, wallet.Metadata.Name)
	assert.Equal(t, description, wallet.Metadata.Description)
	assert.Equal(t, tags, wallet.Metadata.Tags)
}

func TestWalletService_SetWalletMetadata_WalletNotFound(t *testing.T) {
	// Setup
	ctx := context.Background()
	logger := zerolog.New(zerolog.NewTestWriter(t))
	mockRepo := new(MockWalletRepository)
	mockClient := new(WalletMockMEXCClient)
	service := NewWalletService(mockRepo, mockClient, &logger)

	// Test data
	walletID := "wallet123"
	name := "My Primary Wallet"
	description := "This is my main trading wallet"
	tags := []string{"trading", "main", "exchange"}

	// Setup expectations
	mockRepo.On("GetByID", ctx, walletID).Return(nil, nil) // Wallet not found

	// Call the method
	err := service.SetWalletMetadata(ctx, walletID, name, description, tags)

	// Assertions
	require.Error(t, err)
	assert.Contains(t, err.Error(), "wallet not found")
	mockRepo.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "Save")
}

func TestWalletService_AddCustomMetadata(t *testing.T) {
	// Setup
	ctx := context.Background()
	logger := zerolog.New(zerolog.NewTestWriter(t))
	mockRepo := new(MockWalletRepository)
	mockClient := new(WalletMockMEXCClient)
	service := NewWalletService(mockRepo, mockClient, &logger)

	// Test data
	walletID := "wallet123"
	userID := "user123"
	key := "risk_level"
	value := "high"

	// Create test wallet
	wallet := &model.Wallet{
		ID:       walletID,
		UserID:   userID,
		Exchange: "MEXC",
		Type:     model.WalletTypeExchange,
		Status:   model.WalletStatusActive,
		Metadata: &model.WalletMetadata{
			Custom: make(map[string]string),
		},
	}

	// Setup expectations
	mockRepo.On("GetByID", ctx, walletID).Return(wallet, nil)
	// The wallet should be saved with updated custom metadata
	mockRepo.On("Save", ctx, mock.MatchedBy(func(w *model.Wallet) bool {
		return w.ID == walletID && w.Metadata.Custom[key] == value
	})).Return(nil)

	// Call the method
	err := service.AddCustomMetadata(ctx, walletID, key, value)

	// Assertions
	require.NoError(t, err)
	mockRepo.AssertExpectations(t)

	// Verify that the custom metadata was added
	assert.Equal(t, value, wallet.Metadata.Custom[key])
}
