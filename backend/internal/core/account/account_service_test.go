package account

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"go-crypto-bot-clean/backend/internal/domain/models"
)

// Mocks
type MockMexcClient struct {
	mock.Mock
}

func (m *MockMexcClient) FetchBalances(ctx context.Context) (models.Balance, error) {
	args := m.Called(ctx)
	return args.Get(0).(models.Balance), args.Error(1)
}

func (m *MockMexcClient) ValidateKeys(ctx context.Context) (bool, error) {
	args := m.Called(ctx)
	return args.Bool(0), args.Error(1)
}

func (m *MockMexcClient) GetAccountBalance(ctx context.Context) (float64, error) {
	args := m.Called(ctx)
	return args.Get(0).(float64), args.Error(1)
}

type MockBoughtCoinRepo struct {
	mock.Mock
}

func (m *MockBoughtCoinRepo) GetAll(ctx context.Context) ([]string, error) {
	args := m.Called(ctx)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockBoughtCoinRepo) GetPosition(ctx context.Context, symbol string) (float64, error) {
	args := m.Called(ctx, symbol)
	return args.Get(0).(float64), args.Error(1)
}

type MockWalletRepository struct {
	mock.Mock
}

func (m *MockWalletRepository) GetWallet(ctx context.Context) (*models.Wallet, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Wallet), args.Error(1)
}

func (m *MockWalletRepository) SaveWallet(ctx context.Context, wallet *models.Wallet) (*models.Wallet, error) {
	args := m.Called(ctx, wallet)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Wallet), args.Error(1)
}

type MockTransactionRepository struct {
	mock.Mock
}

func (m *MockTransactionRepository) Create(ctx context.Context, transaction *models.Transaction) (*models.Transaction, error) {
	args := m.Called(ctx, transaction)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) FindByID(ctx context.Context, id string) (*models.Transaction, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) FindByTimeRange(ctx context.Context, startTime, endTime time.Time) ([]*models.Transaction, error) {
	args := m.Called(ctx, startTime, endTime)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) FindAll(ctx context.Context) ([]*models.Transaction, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Transaction), args.Error(1)
}

// MockConfig implements the Config interface for testing
type MockConfig struct {
	RiskThreshold float64
	CacheTTLValue time.Duration
}

func (c *MockConfig) GetRiskThreshold() float64 {
	return c.RiskThreshold
}

func (c *MockConfig) GetCacheTTL() time.Duration {
	return c.CacheTTLValue
}

// --- TESTS ---

func TestGetAccountBalance_Success(t *testing.T) {
	mexc := new(MockMexcClient)
	coinRepo := new(MockBoughtCoinRepo)
	walletRepo := new(MockWalletRepository)
	txRepo := new(MockTransactionRepository)
	cfg := &MockConfig{CacheTTLValue: time.Minute}

	svc := NewAccountService(mexc, coinRepo, walletRepo, txRepo, cfg)

	expected := models.Balance{Fiat: 1000, Available: map[string]float64{"BTC": 0.5}}
	mexc.On("FetchBalances", mock.Anything).Return(expected, nil)

	bal, err := svc.GetAccountBalance(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, expected, bal)
}

func TestGetAccountBalance_Error(t *testing.T) {
	mexc := new(MockMexcClient)
	coinRepo := new(MockBoughtCoinRepo)
	walletRepo := new(MockWalletRepository)
	txRepo := new(MockTransactionRepository)
	cfg := &MockConfig{CacheTTLValue: time.Minute}

	svc := NewAccountService(mexc, coinRepo, walletRepo, txRepo, cfg)

	mexc.On("FetchBalances", mock.Anything).Return(models.Balance{}, errors.New("network error"))

	_, err := svc.GetAccountBalance(context.Background())
	assert.Error(t, err)
}

func TestGetWallet_Success(t *testing.T) {
	mexc := new(MockMexcClient)
	coinRepo := new(MockBoughtCoinRepo)
	walletRepo := new(MockWalletRepository)
	txRepo := new(MockTransactionRepository)
	cfg := &MockConfig{CacheTTLValue: time.Minute}

	svc := NewAccountService(mexc, coinRepo, walletRepo, txRepo, cfg)

	expected := &models.Wallet{
		Balances: map[string]*models.AssetBalance{
			"USDT": {
				Asset:  "USDT",
				Free:   1000.0,
				Locked: 0.0,
				Total:  1000.0,
			},
		},
		UpdatedAt: time.Now(),
	}

	walletRepo.On("GetWallet", mock.Anything).Return(expected, nil)

	wallet, err := svc.GetWallet(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, expected, wallet)
}

func TestGetPortfolioValue_Success(t *testing.T) {
	mexc := new(MockMexcClient)
	coinRepo := new(MockBoughtCoinRepo)
	walletRepo := new(MockWalletRepository)
	txRepo := new(MockTransactionRepository)
	cfg := &MockConfig{CacheTTLValue: time.Minute}

	svc := NewAccountService(mexc, coinRepo, walletRepo, txRepo, cfg)

	balance := models.Balance{Fiat: 500, Available: map[string]float64{"BTC": 1.0}}
	mexc.On("FetchBalances", mock.Anything).Return(balance, nil)

	// Assume BTC price is fetched via repo (simplified)
	coinRepo.On("GetAll", mock.Anything).Return([]string{"BTC"}, nil)
	coinRepo.On("GetPosition", mock.Anything, "BTC").Return(60000.0, nil)

	val, err := svc.GetPortfolioValue(context.Background())
	assert.NoError(t, err)
	assert.Greater(t, val, 0.0)
}

func TestGetPositionRisk_Success(t *testing.T) {
	mexc := new(MockMexcClient)
	coinRepo := new(MockBoughtCoinRepo)
	walletRepo := new(MockWalletRepository)
	txRepo := new(MockTransactionRepository)
	cfg := &MockConfig{RiskThreshold: 0.2, CacheTTLValue: time.Minute}

	svc := NewAccountService(mexc, coinRepo, walletRepo, txRepo, cfg)

	coinRepo.On("GetPosition", mock.Anything, "BTC").Return(60000.0, nil)

	risk, err := svc.GetPositionRisk(context.Background(), "BTC")
	assert.NoError(t, err)
	assert.Equal(t, "BTC", risk.Symbol)
}

func TestGetAllPositionRisks_Success(t *testing.T) {
	mexc := new(MockMexcClient)
	coinRepo := new(MockBoughtCoinRepo)
	walletRepo := new(MockWalletRepository)
	txRepo := new(MockTransactionRepository)
	cfg := &MockConfig{RiskThreshold: 0.2, CacheTTLValue: time.Minute}

	svc := NewAccountService(mexc, coinRepo, walletRepo, txRepo, cfg)

	coinRepo.On("GetAll", mock.Anything).Return([]string{"BTC", "ETH"}, nil)
	coinRepo.On("GetPosition", mock.Anything, "BTC").Return(60000.0, nil)
	coinRepo.On("GetPosition", mock.Anything, "ETH").Return(3000.0, nil)

	risks, err := svc.GetAllPositionRisks(context.Background())
	assert.NoError(t, err)
	assert.Len(t, risks, 2)
}

func TestValidateAPIKeys_Success(t *testing.T) {
	mexc := new(MockMexcClient)
	coinRepo := new(MockBoughtCoinRepo)
	walletRepo := new(MockWalletRepository)
	txRepo := new(MockTransactionRepository)
	cfg := &MockConfig{CacheTTLValue: time.Minute}

	svc := NewAccountService(mexc, coinRepo, walletRepo, txRepo, cfg)

	mexc.On("ValidateKeys", mock.Anything).Return(true, nil)

	valid, err := svc.ValidateAPIKeys(context.Background())
	assert.NoError(t, err)
	assert.True(t, valid)
}

func TestValidateAPIKeys_Error(t *testing.T) {
	mexc := new(MockMexcClient)
	coinRepo := new(MockBoughtCoinRepo)
	walletRepo := new(MockWalletRepository)
	txRepo := new(MockTransactionRepository)
	cfg := &MockConfig{CacheTTLValue: time.Minute}

	svc := NewAccountService(mexc, coinRepo, walletRepo, txRepo, cfg)

	mexc.On("ValidateKeys", mock.Anything).Return(false, errors.New("invalid keys"))

	valid, err := svc.ValidateAPIKeys(context.Background())
	assert.Error(t, err)
	assert.False(t, valid)
}

func TestUpdateBalance_Success(t *testing.T) {
	mexc := new(MockMexcClient)
	coinRepo := new(MockBoughtCoinRepo)
	walletRepo := new(MockWalletRepository)
	txRepo := new(MockTransactionRepository)
	cfg := &MockConfig{CacheTTLValue: time.Minute}

	svc := NewAccountService(mexc, coinRepo, walletRepo, txRepo, cfg)

	// Setup initial wallet
	wallet := &models.Wallet{
		Balances: map[string]*models.AssetBalance{
			"USDT": {
				Asset:  "USDT",
				Free:   1000.0,
				Locked: 0.0,
				Total:  1000.0,
			},
		},
		UpdatedAt: time.Now(),
	}

	// Setup balance from exchange
	balance := models.Balance{Fiat: 1000.0, Available: map[string]float64{"USDT": 1000.0}}

	// Setup expectations
	walletRepo.On("GetWallet", mock.Anything).Return(wallet, nil)
	mexc.On("FetchBalances", mock.Anything).Return(balance, nil)
	walletRepo.On("SaveWallet", mock.Anything, mock.Anything).Return(wallet, nil)
	txRepo.On("Create", mock.Anything, mock.Anything).Return(&models.Transaction{ID: "00000000-0000-0000-0000-000000000001"}, nil)

	// Test updating balance
	err := svc.UpdateBalance(context.Background(), 500.0, "Test deposit")
	assert.NoError(t, err)

	// Verify mocks were called
	walletRepo.AssertExpectations(t)
	txRepo.AssertExpectations(t)
}

func TestSyncWithExchange_Success(t *testing.T) {
	mexc := new(MockMexcClient)
	coinRepo := new(MockBoughtCoinRepo)
	walletRepo := new(MockWalletRepository)
	txRepo := new(MockTransactionRepository)
	cfg := &MockConfig{CacheTTLValue: time.Minute}

	svc := NewAccountService(mexc, coinRepo, walletRepo, txRepo, cfg)

	// Setup initial wallet
	wallet := &models.Wallet{
		Balances: map[string]*models.AssetBalance{
			"USDT": {
				Asset:  "USDT",
				Free:   1000.0,
				Locked: 0.0,
				Total:  1000.0,
			},
		},
		UpdatedAt: time.Now(),
	}

	// Setup expectations
	walletRepo.On("GetWallet", mock.Anything).Return(wallet, nil)
	mexc.On("GetAccountBalance", mock.Anything).Return(1500.0, nil)
	mexc.On("FetchBalances", mock.Anything).Return(models.Balance{Fiat: 1500.0}, nil)
	walletRepo.On("SaveWallet", mock.Anything, mock.Anything).Return(wallet, nil)
	txRepo.On("Create", mock.Anything, mock.Anything).Return(&models.Transaction{ID: "00000000-0000-0000-0000-000000000001"}, nil)

	// Test syncing with exchange
	err := svc.SyncWithExchange(context.Background())
	assert.NoError(t, err)

	// Verify mocks were called
	mexc.AssertExpectations(t)
	walletRepo.AssertExpectations(t)
}

func TestGetBalanceSummary_Success(t *testing.T) {
	mexc := new(MockMexcClient)
	coinRepo := new(MockBoughtCoinRepo)
	walletRepo := new(MockWalletRepository)
	txRepo := new(MockTransactionRepository)
	cfg := &MockConfig{CacheTTLValue: time.Minute}

	svc := NewAccountService(mexc, coinRepo, walletRepo, txRepo, cfg)

	// Setup balance from exchange
	balance := models.Balance{Fiat: 1000.0, Available: map[string]float64{"USDT": 1000.0}}

	// Setup transactions
	now := time.Now()
	transactions := []*models.Transaction{
		{ID: "00000000-0000-0000-0000-000000000001", Amount: 500.0, Balance: 500.0, Reason: "Initial deposit", Timestamp: now.Add(-48 * time.Hour)},
		{ID: "00000000-0000-0000-0000-000000000002", Amount: 300.0, Balance: 800.0, Reason: "Deposit", Timestamp: now.Add(-24 * time.Hour)},
		{ID: "00000000-0000-0000-0000-000000000003", Amount: -100.0, Balance: 700.0, Reason: "Withdrawal", Timestamp: now.Add(-12 * time.Hour)},
		{ID: "00000000-0000-0000-0000-000000000004", Amount: 300.0, Balance: 1000.0, Reason: "Deposit", Timestamp: now.Add(-6 * time.Hour)},
	}

	// Setup expectations
	mexc.On("FetchBalances", mock.Anything).Return(balance, nil)
	txRepo.On("FindByTimeRange", mock.Anything, mock.Anything, mock.Anything).Return(transactions, nil)

	// Test getting balance summary
	summary, err := svc.GetBalanceSummary(context.Background(), 3)
	assert.NoError(t, err)
	assert.Equal(t, 1000.0, summary.CurrentBalance)
	assert.Equal(t, 1100.0, summary.Deposits)
	assert.Equal(t, 100.0, summary.Withdrawals)
	assert.Equal(t, 1000.0, summary.NetChange)
	assert.Equal(t, 4, summary.TransactionCount)

	// Verify mocks were called
	mexc.AssertExpectations(t)
	txRepo.AssertExpectations(t)
}

func TestAnalyzeTransactions_Success(t *testing.T) {
	mexc := new(MockMexcClient)
	coinRepo := new(MockBoughtCoinRepo)
	walletRepo := new(MockWalletRepository)
	txRepo := new(MockTransactionRepository)
	cfg := &MockConfig{CacheTTLValue: time.Minute}

	svc := NewAccountService(mexc, coinRepo, walletRepo, txRepo, cfg)

	// Setup transactions
	now := time.Now()
	transactions := []*models.Transaction{
		{ID: "00000000-0000-0000-0000-000000000001", Amount: 500.0, Balance: 500.0, Reason: "Initial deposit", Timestamp: now.Add(-48 * time.Hour)},
		{ID: "00000000-0000-0000-0000-000000000002", Amount: 300.0, Balance: 800.0, Reason: "Buy BTC", Timestamp: now.Add(-24 * time.Hour)},
		{ID: "00000000-0000-0000-0000-000000000003", Amount: -100.0, Balance: 700.0, Reason: "Sell ETH", Timestamp: now.Add(-12 * time.Hour)},
		{ID: "00000000-0000-0000-0000-000000000004", Amount: 300.0, Balance: 1000.0, Reason: "Deposit", Timestamp: now.Add(-6 * time.Hour)},
	}

	// Setup expectations
	txRepo.On("FindByTimeRange", mock.Anything, mock.Anything, mock.Anything).Return(transactions, nil)

	// Test analyzing transactions
	analysis, err := svc.AnalyzeTransactions(context.Background(), now.Add(-72*time.Hour), now)
	assert.NoError(t, err)
	assert.Equal(t, 4, analysis.TotalCount)
	// The isBuyTransaction function matches "buy", "purchase", and "deposit"
	assert.Equal(t, 3, analysis.BuyCount)
	// The isSellTransaction function matches "sell", "sale", and "withdrawal"
	assert.Equal(t, 1, analysis.SellCount)

	// Verify mocks were called
	txRepo.AssertExpectations(t)
}
