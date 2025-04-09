package account

import (
	"context"
	"errors"
	"testing"
	"time"

	"go-crypto-bot-clean/backend/internal/domain/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// MockMexcRESTClient is a mock implementation of the MexcRESTClient interface
type MockMexcRESTClient struct {
	mock.Mock
}

func (m *MockMexcRESTClient) FetchBalances(ctx context.Context) (models.Balance, error) {
	args := m.Called(ctx)
	return args.Get(0).(models.Balance), args.Error(1)
}

func (m *MockMexcRESTClient) ValidateKeys(ctx context.Context) (bool, error) {
	args := m.Called(ctx)
	return args.Bool(0), args.Error(1)
}

func (m *MockMexcRESTClient) GetAccountBalance(ctx context.Context) (float64, error) {
	args := m.Called(ctx)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockMexcRESTClient) GetWallet(ctx context.Context) (*models.Wallet, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Wallet), args.Error(1)
}

// MockMexcWSClient is a mock implementation of the MexcWebSocketClient interface
type MockMexcWSClient struct {
	mock.Mock
}

func (m *MockMexcWSClient) Connect(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockMexcWSClient) Disconnect() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockMexcWSClient) IsConnected() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockMexcWSClient) SubscribeToAccountUpdates(ctx context.Context, callback func(*models.Wallet)) error {
	args := m.Called(ctx, callback)
	return args.Error(0)
}

func (m *MockMexcWSClient) UnsubscribeFromAccountUpdates(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockMexcWSClient) Authenticate(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockMexcWSClient) SetReconnectHandler(handler func() error) {
	m.Called(handler)
}

// TestRealAccountService_GetAccountBalance tests the GetAccountBalance method
func TestRealAccountService_GetAccountBalance(t *testing.T) {
	// Create mocks
	restClient := new(MockMexcRESTClient)
	wsClient := new(MockMexcWSClient)
	coinRepo := new(MockBoughtCoinRepo)
	walletRepo := new(MockWalletRepository)
	txRepo := new(MockTransactionRepository)
	cfg := &MockConfig{CacheTTLValue: time.Minute}

	// Setup mock expectations for the reconnect handler
	wsClient.On("SetReconnectHandler", mock.AnythingOfType("func() error")).Return()

	// Create service
	svc := NewRealAccountService(restClient, wsClient, coinRepo, walletRepo, txRepo, cfg)

	// Setup expected data
	expectedBalance := models.Balance{
		Fiat: 1000.0,
		Available: map[string]float64{
			"USDT": 1000.0,
			"BTC":  0.01,
		},
	}

	// Setup mock expectations
	restClient.On("FetchBalances", mock.Anything).Return(expectedBalance, nil)

	// Call the method
	balance, err := svc.GetAccountBalance(context.Background())

	// Assert expectations
	require.NoError(t, err)
	assert.Equal(t, expectedBalance, balance)
	restClient.AssertExpectations(t)
	wsClient.AssertExpectations(t)
}

// TestRealAccountService_GetWallet tests the GetWallet method
func TestRealAccountService_GetWallet(t *testing.T) {
	// Create mocks
	restClient := new(MockMexcRESTClient)
	wsClient := new(MockMexcWSClient)
	coinRepo := new(MockBoughtCoinRepo)
	walletRepo := new(MockWalletRepository)
	txRepo := new(MockTransactionRepository)
	cfg := &MockConfig{CacheTTLValue: time.Minute}

	// Setup mock expectations for the reconnect handler
	wsClient.On("SetReconnectHandler", mock.AnythingOfType("func() error")).Return()

	// Create service
	svc := NewRealAccountService(restClient, wsClient, coinRepo, walletRepo, txRepo, cfg)

	// Setup expected data
	expectedWallet := &models.Wallet{
		Balances: map[string]*models.AssetBalance{
			"USDT": {
				Asset:  "USDT",
				Free:   1000.0,
				Locked: 0.0,
				Total:  1000.0,
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
	walletRepo.On("GetWallet", mock.Anything).Return(expectedWallet, nil)

	// Call the method
	wallet, err := svc.GetWallet(context.Background())

	// Assert expectations
	require.NoError(t, err)
	assert.Equal(t, expectedWallet, wallet)
	walletRepo.AssertExpectations(t)
}

// TestRealAccountService_SubscribeToBalanceUpdates tests the SubscribeToBalanceUpdates method
func TestRealAccountService_SubscribeToBalanceUpdates(t *testing.T) {
	// Create mocks
	restClient := new(MockMexcRESTClient)
	wsClient := new(MockMexcWSClient)
	coinRepo := new(MockBoughtCoinRepo)
	walletRepo := new(MockWalletRepository)
	txRepo := new(MockTransactionRepository)
	cfg := &MockConfig{CacheTTLValue: time.Minute}

	// Setup mock expectations for the reconnect handler
	wsClient.On("SetReconnectHandler", mock.AnythingOfType("func() error")).Return()

	// Create service
	svc := NewRealAccountService(restClient, wsClient, coinRepo, walletRepo, txRepo, cfg)

	// Setup mock expectations
	wsClient.On("IsConnected").Return(false)
	wsClient.On("Connect", mock.Anything).Return(nil)
	wsClient.On("Authenticate", mock.Anything).Return(nil)
	wsClient.On("SubscribeToAccountUpdates", mock.Anything, mock.AnythingOfType("func(*models.Wallet)")).Return(nil)

	// Create a callback function with a channel to signal when it's called
	callbackCh := make(chan struct{}, 1)
	callback := func(wallet *models.Wallet) {
		callbackCh <- struct{}{}
	}

	// Call the method
	err := svc.SubscribeToBalanceUpdates(context.Background(), callback)

	// Assert expectations
	require.NoError(t, err)
	wsClient.AssertExpectations(t)

	// Simulate a wallet update
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

	// Setup expectation for SaveWallet
	walletRepo.On("SaveWallet", mock.Anything, mock.Anything).Return(wallet, nil)

	// Print the calls to debug
	t.Logf("Number of calls: %d", len(wsClient.Calls))
	for i, call := range wsClient.Calls {
		t.Logf("Call %d: %s with %d arguments", i, call.Method, len(call.Arguments))
	}

	// Find the SubscribeToAccountUpdates call
	var wsCallback func(*models.Wallet)
	for _, call := range wsClient.Calls {
		if call.Method == "SubscribeToAccountUpdates" && len(call.Arguments) > 1 {
			wsCallback = call.Arguments[1].(func(*models.Wallet))
			break
		}
	}

	if wsCallback == nil {
		t.Fatal("Could not find SubscribeToAccountUpdates call with callback")
	}

	// Call the callback function
	wsCallback(wallet)

	// Wait for the callback to be called or timeout
	select {
	case <-callbackCh:
		// Callback was called successfully
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Callback was not called within timeout")
	}
}

// TestRealAccountService_SyncWithExchange tests the SyncWithExchange method
func TestRealAccountService_SyncWithExchange(t *testing.T) {
	// Create mocks
	restClient := new(MockMexcRESTClient)
	wsClient := new(MockMexcWSClient)
	coinRepo := new(MockBoughtCoinRepo)
	walletRepo := new(MockWalletRepository)
	txRepo := new(MockTransactionRepository)
	cfg := &MockConfig{CacheTTLValue: time.Minute}

	// Setup mock expectations for the reconnect handler
	wsClient.On("SetReconnectHandler", mock.AnythingOfType("func() error")).Return()

	// Create service
	svc := NewRealAccountService(restClient, wsClient, coinRepo, walletRepo, txRepo, cfg)

	// Setup expected data
	expectedBalance := models.Balance{
		Fiat: 1000.0,
		Available: map[string]float64{
			"USDT": 1000.0,
			"BTC":  0.01,
		},
	}

	expectedWallet := &models.Wallet{
		Balances: map[string]*models.AssetBalance{
			"USDT": {
				Asset:  "USDT",
				Free:   1000.0,
				Locked: 0.0,
				Total:  1000.0,
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
	restClient.On("FetchBalances", mock.Anything).Return(expectedBalance, nil)
	walletRepo.On("SaveWallet", mock.Anything, mock.Anything).Return(expectedWallet, nil)

	// Call the method
	err := svc.SyncWithExchange(context.Background())

	// Assert expectations
	require.NoError(t, err)
	restClient.AssertExpectations(t)
	walletRepo.AssertExpectations(t)
}

// TestRealAccountService_HandleAPIError tests error handling for API calls
func TestRealAccountService_HandleAPIError(t *testing.T) {
	// Create mocks
	restClient := new(MockMexcRESTClient)
	wsClient := new(MockMexcWSClient)
	coinRepo := new(MockBoughtCoinRepo)
	walletRepo := new(MockWalletRepository)
	txRepo := new(MockTransactionRepository)
	cfg := &MockConfig{CacheTTLValue: time.Minute}

	// Setup mock expectations for the reconnect handler
	wsClient.On("SetReconnectHandler", mock.AnythingOfType("func() error")).Return()

	// Create service with logger
	logger, _ := zap.NewDevelopment()
	svc := NewRealAccountServiceWithLogger(restClient, wsClient, coinRepo, walletRepo, txRepo, cfg, logger)

	// Setup mock expectations with error
	apiError := errors.New("API connection error")
	restClient.On("FetchBalances", mock.Anything).Return(models.Balance{}, apiError)

	// Call the method
	_, err := svc.GetAccountBalance(context.Background())

	// Assert expectations
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to fetch balances")
	restClient.AssertExpectations(t)
}

// TestRealAccountService_WebSocketReconnection tests the WebSocket reconnection mechanism
func TestRealAccountService_WebSocketReconnection(t *testing.T) {
	// Create mocks
	restClient := new(MockMexcRESTClient)
	wsClient := new(MockMexcWSClient)
	coinRepo := new(MockBoughtCoinRepo)
	walletRepo := new(MockWalletRepository)
	txRepo := new(MockTransactionRepository)
	cfg := &MockConfig{CacheTTLValue: time.Minute}

	// Setup mock expectations for the reconnect handler
	wsClient.On("SetReconnectHandler", mock.AnythingOfType("func() error")).Return()

	// Create service
	logger, _ := zap.NewDevelopment()
	svc := NewRealAccountServiceWithLogger(restClient, wsClient, coinRepo, walletRepo, txRepo, cfg, logger)

	// Setup mock expectations
	wsClient.On("IsConnected").Return(false)
	wsClient.On("Connect", mock.Anything).Return(nil)
	wsClient.On("Authenticate", mock.Anything).Return(nil)
	wsClient.On("SetReconnectHandler", mock.AnythingOfType("func() error")).Return()
	wsClient.On("SubscribeToAccountUpdates", mock.Anything, mock.AnythingOfType("func(*models.Wallet)")).Return(nil)

	// Call the method to subscribe (which should trigger connection)
	callback := func(wallet *models.Wallet) {}
	err := svc.SubscribeToBalanceUpdates(context.Background(), callback)

	// Assert expectations
	require.NoError(t, err)
	wsClient.AssertExpectations(t)

	// Print the calls to debug
	t.Logf("Number of calls: %d", len(wsClient.Calls))
	for i, call := range wsClient.Calls {
		t.Logf("Call %d: %s with %d arguments", i, call.Method, len(call.Arguments))
	}

	// Find the SetReconnectHandler call
	var reconnectHandler func() error
	for _, call := range wsClient.Calls {
		if call.Method == "SetReconnectHandler" && len(call.Arguments) > 0 {
			reconnectHandler = call.Arguments[0].(func() error)
			break
		}
	}

	if reconnectHandler == nil {
		t.Fatal("Could not find SetReconnectHandler call with handler function")
	}

	// Reset mock expectations for reconnection
	wsClient = new(MockMexcWSClient)
	wsClient.On("Connect", mock.Anything).Return(nil)
	wsClient.On("Authenticate", mock.Anything).Return(nil)
	wsClient.On("SubscribeToAccountUpdates", mock.Anything, mock.AnythingOfType("func(*models.Wallet)")).Return(nil)

	// Replace the WebSocket client in the service
	svc.(*realAccountService).wsClient = wsClient

	// Call the reconnect handler
	err = reconnectHandler()

	// Assert expectations
	require.NoError(t, err)
	wsClient.AssertExpectations(t)
}

// TestRealAccountService_GetBalanceSummary tests the GetBalanceSummary method
func TestRealAccountService_GetBalanceSummary(t *testing.T) {
	// Create mocks
	restClient := new(MockMexcRESTClient)
	wsClient := new(MockMexcWSClient)
	coinRepo := new(MockBoughtCoinRepo)
	walletRepo := new(MockWalletRepository)
	txRepo := new(MockTransactionRepository)
	cfg := &MockConfig{CacheTTLValue: time.Minute}

	// Setup mock expectations for the reconnect handler
	wsClient.On("SetReconnectHandler", mock.AnythingOfType("func() error")).Return()

	// Create service
	svc := NewRealAccountService(restClient, wsClient, coinRepo, walletRepo, txRepo, cfg)

	// Setup expected data
	expectedBalance := models.Balance{
		Fiat: 1000.0,
		Available: map[string]float64{
			"USDT": 1000.0,
		},
	}

	// Create some test transactions
	now := time.Now()
	transactions := []*models.Transaction{
		{
			Amount:    100.0,
			Balance:   1000.0,
			Reason:    "deposit",
			Timestamp: now.Add(-24 * time.Hour),
		},
		{
			Amount:    -50.0,
			Balance:   950.0,
			Reason:    "withdrawal",
			Timestamp: now.Add(-12 * time.Hour),
		},
		{
			Amount:    100.0,
			Balance:   1050.0,
			Reason:    "deposit",
			Timestamp: now.Add(-6 * time.Hour),
		},
		{
			Amount:    -50.0,
			Balance:   1000.0,
			Reason:    "withdrawal",
			Timestamp: now.Add(-1 * time.Hour),
		},
	}

	// Setup mock expectations
	restClient.On("FetchBalances", mock.Anything).Return(expectedBalance, nil)
	txRepo.On("FindByTimeRange", mock.Anything, mock.Anything, mock.Anything).Return(transactions, nil)

	// Call the method
	summary, err := svc.GetBalanceSummary(context.Background(), 1)

	// Assert expectations
	require.NoError(t, err)
	assert.Equal(t, 1000.0, summary.CurrentBalance)
	assert.Equal(t, 200.0, summary.Deposits)
	assert.Equal(t, 100.0, summary.Withdrawals)
	assert.Equal(t, 100.0, summary.NetChange)
	assert.Equal(t, 4, summary.TransactionCount)
	assert.Equal(t, 1, summary.Period)
	restClient.AssertExpectations(t)
	txRepo.AssertExpectations(t)
}
