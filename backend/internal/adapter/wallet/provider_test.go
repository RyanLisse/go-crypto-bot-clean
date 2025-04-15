package wallet

import (
	"context"
	"testing"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockMEXCClient is a mock implementation of the MEXCClient interface
type MockMEXCClient struct {
	mock.Mock
}

// MockCoinbaseClient is a mock implementation of the CoinbaseClient interface
type MockCoinbaseClient struct {
	mock.Mock
}

func (m *MockCoinbaseClient) GetAccount(ctx context.Context) (*model.Wallet, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Wallet), args.Error(1)
}

// Ensure MockCoinbaseClient implements port.CoinbaseClient
var _ port.CoinbaseClient = (*MockCoinbaseClient)(nil);

// MockMEXCClient methods below

func (m *MockMEXCClient) GetAccount(ctx context.Context) (*model.Wallet, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Wallet), args.Error(1)
}

// Implement other required methods of the MEXCClient interface
func (m *MockMEXCClient) GetNewListings(ctx context.Context) ([]*model.NewCoin, error) {
	return nil, nil
}

func (m *MockMEXCClient) GetSymbolInfo(ctx context.Context, symbol string) (*model.SymbolInfo, error) {
	return nil, nil
}

func (m *MockMEXCClient) GetSymbolStatus(ctx context.Context, symbol string) (model.Status, error) {
	return "", nil
}

func (m *MockMEXCClient) GetTradingSchedule(ctx context.Context, symbol string) (model.TradingSchedule, error) {
	return model.TradingSchedule{}, nil
}

func (m *MockMEXCClient) GetSymbolConstraints(ctx context.Context, symbol string) (*model.SymbolConstraints, error) {
	return nil, nil
}

func (m *MockMEXCClient) GetExchangeInfo(ctx context.Context) (*model.ExchangeInfo, error) {
	return nil, nil
}

func (m *MockMEXCClient) GetMarketData(ctx context.Context, symbol string) (*model.Ticker, error) {
	return nil, nil
}

func (m *MockMEXCClient) GetKlines(ctx context.Context, symbol string, interval model.KlineInterval, limit int) ([]*model.Kline, error) {
	return nil, nil
}

func (m *MockMEXCClient) GetOrderBook(ctx context.Context, symbol string, depth int) (*model.OrderBook, error) {
	return nil, nil
}

func (m *MockMEXCClient) PlaceOrder(ctx context.Context, symbol string, side model.OrderSide, orderType model.OrderType, quantity float64, price float64, timeInForce model.TimeInForce) (*model.Order, error) {
	return nil, nil
}

func (m *MockMEXCClient) CancelOrder(ctx context.Context, symbol string, orderID string) error {
	return nil
}

func (m *MockMEXCClient) GetOrderStatus(ctx context.Context, symbol string, orderID string) (*model.Order, error) {
	return nil, nil
}

func (m *MockMEXCClient) GetOpenOrders(ctx context.Context, symbol string) ([]*model.Order, error) {
	return nil, nil
}

func (m *MockMEXCClient) GetOrderHistory(ctx context.Context, symbol string, limit, offset int) ([]*model.Order, error) {
	return nil, nil
}

// Ensure MockMEXCClient implements port.MEXCClient
var _ port.MEXCClient = (*MockMEXCClient)(nil)

func TestCoinbaseProvider(t *testing.T) {
	logger := zerolog.New(zerolog.NewTestWriter(t))
	mockClient := new(MockCoinbaseClient)
	provider := NewCoinbaseProvider(mockClient, &logger)

	assert.Equal(t, "Coinbase", provider.GetName())
	assert.Equal(t, model.WalletTypeExchange, provider.GetType())

	mockWallet := &model.Wallet{
		ID:            "cb_wallet",
		UserID:        "usercb",
		Exchange:      "Coinbase",
		Type:          model.WalletTypeExchange,
		Status:        model.WalletStatusActive,
		TotalUSDValue: 2000.0,
		Balances:      make(map[model.Asset]*model.Balance),
	}
	mockWallet.Balances[model.AssetBTC] = &model.Balance{
		Asset:    model.AssetBTC,
		Free:     0.5,
		Locked:   0.0,
		Total:    0.5,
		USDValue: 30000.0,
	}

	mockClient.On("GetAccount", mock.Anything).Return(mockWallet, nil)

	params := map[string]interface{}{
		"user_id":    "usercb",
		"api_key":    "cb_api_key",
		"api_secret": "cb_api_secret",
	}
	wallet, err := provider.Connect(context.Background(), params)
	require.NoError(t, err)
	assert.NotNil(t, wallet)
	assert.Equal(t, "usercb", wallet.UserID)
	assert.Equal(t, "Coinbase", wallet.Exchange)
	assert.Equal(t, model.WalletTypeExchange, wallet.Type)
	assert.NotEmpty(t, wallet.Balances)

	mockClient.On("GetAccount", mock.Anything).Return(mockWallet, nil)
	updatedWallet, err := provider.GetBalance(context.Background(), wallet)
	require.NoError(t, err)
	assert.NotNil(t, updatedWallet)
	assert.Equal(t, mockWallet.TotalUSDValue, updatedWallet.TotalUSDValue)
	assert.NotEmpty(t, updatedWallet.Balances)

	mockClient.On("GetAccount", mock.Anything).Return(mockWallet, nil)
	verified, err := provider.Verify(context.Background(), "Coinbase", "test_message", "test_signature")
	require.NoError(t, err)
	assert.True(t, verified)

	err = provider.Disconnect(context.Background(), "cb_wallet")
	require.NoError(t, err)

	mockClient.AssertExpectations(t)
}


func TestMEXCProvider(t *testing.T) {
	// Setup
	logger := zerolog.New(zerolog.NewTestWriter(t))
	mockClient := new(MockMEXCClient)
	provider := NewMEXCProvider(mockClient, &logger)

	// Test GetName
	assert.Equal(t, "MEXC", provider.GetName())

	// Test GetType
	assert.Equal(t, model.WalletTypeExchange, provider.GetType())

	// Setup mock for Connect test
	mockWallet := &model.Wallet{
		ID:            "test_wallet",
		UserID:        "user123",
		Exchange:      "MEXC",
		Type:          model.WalletTypeExchange,
		Status:        model.WalletStatusActive,
		TotalUSDValue: 1000.0,
		Balances:      make(map[model.Asset]*model.Balance),
	}
	mockWallet.Balances[model.AssetBTC] = &model.Balance{
		Asset:    model.AssetBTC,
		Free:     1.0,
		Locked:   0.1,
		Total:    1.1,
		USDValue: 50000.0,
	}

	mockClient.On("GetAccount", mock.Anything).Return(mockWallet, nil)

	// Test Connect
	params := map[string]interface{}{
		"user_id":    "user123",
		"api_key":    "test_api_key",
		"api_secret": "test_api_secret",
	}
	wallet, err := provider.Connect(context.Background(), params)
	require.NoError(t, err)
	assert.NotNil(t, wallet)
	assert.Equal(t, "user123", wallet.UserID)
	assert.Equal(t, "MEXC", wallet.Exchange)
	assert.Equal(t, model.WalletTypeExchange, wallet.Type)
	assert.Equal(t, model.WalletStatusActive, wallet.Status)
	assert.NotEmpty(t, wallet.Balances)

	// Test GetBalance
	mockClient.On("GetAccount", mock.Anything).Return(mockWallet, nil)
	updatedWallet, err := provider.GetBalance(context.Background(), wallet)
	require.NoError(t, err)
	assert.NotNil(t, updatedWallet)
	assert.Equal(t, mockWallet.TotalUSDValue, updatedWallet.TotalUSDValue)
	assert.NotEmpty(t, updatedWallet.Balances)

	// Test Verify
	mockClient.On("GetAccount", mock.Anything).Return(mockWallet, nil)
	verified, err := provider.Verify(context.Background(), "MEXC", "test_message", "test_signature")
	require.NoError(t, err)
	assert.True(t, verified)

	// Test Disconnect
	err = provider.Disconnect(context.Background(), "test_wallet")
	require.NoError(t, err)

	// Verify all expectations were met
	mockClient.AssertExpectations(t)
}

func TestEthereumProvider(t *testing.T) {
	// Setup
	logger := zerolog.New(zerolog.NewTestWriter(t))
	provider := NewEthereumProvider(1, "Ethereum", "https://mainnet.infura.io/v3/test_key", "https://etherscan.io", &logger)

	// Test GetName
	assert.Equal(t, "Ethereum", provider.GetName())

	// Test GetType
	assert.Equal(t, model.WalletTypeWeb3, provider.GetType())

	// Test GetChainID
	assert.Equal(t, int64(1), provider.GetChainID())

	// Test GetNetwork
	assert.Equal(t, "Ethereum", provider.GetNetwork())

	// Test IsValidAddress with valid address
	valid, err := provider.IsValidAddress(context.Background(), "0x742d35Cc6634C0532925a3b844Bc454e4438f44e")
	require.NoError(t, err)
	assert.True(t, valid)

	// Test IsValidAddress with invalid address
	valid, err = provider.IsValidAddress(context.Background(), "invalid_address")
	require.NoError(t, err)
	assert.False(t, valid)

	// Test Connect
	params := map[string]interface{}{
		"user_id": "user123",
		"address": "0x742d35Cc6634C0532925a3b844Bc454e4438f44e",
	}
	wallet, err := provider.Connect(context.Background(), params)
	require.NoError(t, err)
	assert.NotNil(t, wallet)
	assert.Equal(t, "user123", wallet.UserID)
	assert.Equal(t, model.WalletTypeWeb3, wallet.Type)
	assert.Equal(t, model.WalletStatusActive, wallet.Status)
	assert.Equal(t, "Ethereum", wallet.Metadata.Network)
	assert.Equal(t, "0x742d35Cc6634C0532925a3b844Bc454e4438f44e", wallet.Metadata.Address)
	assert.Equal(t, int64(1), wallet.Metadata.ChainID)
}

func TestProviderRegistry(t *testing.T) {
	// Setup
	logger := zerolog.New(zerolog.NewTestWriter(t))
	mockMEXC := new(MockMEXCClient)
	mockCoinbase := new(MockCoinbaseClient)
	registry := NewProviderRegistry()

	// Register providers
	mexcProvider := NewMEXCProvider(mockMEXC, &logger)
	coinbaseProvider := NewCoinbaseProvider(mockCoinbase, &logger)
	ethereumProvider := NewEthereumProvider(1, "Ethereum", "https://mainnet.infura.io/v3/test_key", "https://etherscan.io", &logger)
	registry.RegisterProvider(mexcProvider)
	registry.RegisterProvider(coinbaseProvider)
	registry.RegisterProvider(ethereumProvider)

	// Test GetProvider
	provider, err := registry.GetProvider("MEXC")
	require.NoError(t, err)
	assert.Equal(t, "MEXC", provider.GetName())

	provider, err = registry.GetProvider("Coinbase")
	require.NoError(t, err)
	assert.Equal(t, "Coinbase", provider.GetName())

	provider, err = registry.GetProvider("Ethereum")
	require.NoError(t, err)
	assert.Equal(t, "Ethereum", provider.GetName())

	// Test GetExchangeProvider
	exchangeProvider, err := registry.GetExchangeProvider("MEXC")
	require.NoError(t, err)
	assert.Equal(t, "MEXC", exchangeProvider.GetName())

	exchangeProvider, err = registry.GetExchangeProvider("Coinbase")
	require.NoError(t, err)
	assert.Equal(t, "Coinbase", exchangeProvider.GetName())

	// Test GetWeb3Provider
	web3Provider, err := registry.GetWeb3Provider("Ethereum")
	require.NoError(t, err)
	assert.Equal(t, "Ethereum", web3Provider.GetName())

	// Test GetProviderByType
	exchangeProviders, err := registry.GetProviderByType(model.WalletTypeExchange)
	require.NoError(t, err)
	assert.Len(t, exchangeProviders, 2)
	assert.ElementsMatch(t, []string{"MEXC", "Coinbase"}, []string{exchangeProviders[0].GetName(), exchangeProviders[1].GetName()})

	web3Providers, err := registry.GetProviderByType(model.WalletTypeWeb3)
	require.NoError(t, err)
	assert.Len(t, web3Providers, 1)
	assert.Equal(t, "Ethereum", web3Providers[0].GetName())

	// Test GetAllProviders
	allProviders := registry.GetAllProviders()
	assert.Len(t, allProviders, 3)

	// Test GetAllExchangeProviders
	allExchangeProviders := registry.GetAllExchangeProviders()
	assert.Len(t, allExchangeProviders, 2)
	assert.ElementsMatch(t, []string{"MEXC", "Coinbase"}, []string{allExchangeProviders[0].GetName(), allExchangeProviders[1].GetName()})

	// Test GetAllWeb3Providers
	allWeb3Providers := registry.GetAllWeb3Providers()
	assert.Len(t, allWeb3Providers, 1)
	assert.Equal(t, "Ethereum", allWeb3Providers[0].GetName())
}

