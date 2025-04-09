package handlers

import (
	"context"

	"go-crypto-bot-clean/backend/internal/domain/models"
)

type mockExchangeService struct {
	Tickers map[string]*models.Ticker
	Ticker  *models.Ticker
	Err     error
}

func (m *mockExchangeService) GetAllTickers(ctx context.Context) (map[string]*models.Ticker, error) {
	return m.Tickers, m.Err
}

func (m *mockExchangeService) GetTicker(ctx context.Context, symbol string) (*models.Ticker, error) {
	return m.Ticker, m.Err
}

func (m *mockExchangeService) Connect(ctx context.Context) error { return nil }
func (m *mockExchangeService) Disconnect() error                 { return nil }
func (m *mockExchangeService) GetKlines(ctx context.Context, symbol, interval string, limit int) ([]*models.Kline, error) {
	return nil, nil
}
func (m *mockExchangeService) GetNewCoins(ctx context.Context) ([]*models.NewCoin, error) {
	return nil, nil
}
func (m *mockExchangeService) GetWallet(ctx context.Context) (*models.Wallet, error) { return nil, nil }
func (m *mockExchangeService) PlaceOrder(ctx context.Context, order *models.Order) (*models.Order, error) {
	return nil, nil
}
func (m *mockExchangeService) CancelOrder(ctx context.Context, orderID, symbol string) error {
	return nil
}
func (m *mockExchangeService) GetOrder(ctx context.Context, orderID, symbol string) (*models.Order, error) {
	return nil, nil
}
func (m *mockExchangeService) GetOpenOrders(ctx context.Context, symbol string) ([]*models.Order, error) {
	return nil, nil
}
func (m *mockExchangeService) SubscribeToTickers(ctx context.Context, symbols []string, updates chan<- *models.Ticker) error {
	return nil
}
func (m *mockExchangeService) UnsubscribeFromTickers(ctx context.Context, symbols []string) error {
	return nil
}

type mockNewCoinWatcher struct {
	Coins []*models.NewCoin
	Err   error
}

func (m *mockNewCoinWatcher) DetectNewCoins(ctx context.Context) ([]*models.NewCoin, error) {
	return m.Coins, m.Err
}

func (m *mockNewCoinWatcher) StartWatching(ctx context.Context) error {
	return nil
}
