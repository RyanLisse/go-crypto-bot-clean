package service

import (
	"context"

	"go-crypto-bot-clean/backend/internal/domain/models"
)

// ExchangeService defines the interface for interacting with cryptocurrency exchanges
type ExchangeService interface {
	// Connection management
	Connect(ctx context.Context) error
	Disconnect() error

	// Market data operations
	GetTicker(ctx context.Context, symbol string) (*models.Ticker, error)
	GetAllTickers(ctx context.Context) (map[string]*models.Ticker, error)
	GetKlines(ctx context.Context, symbol, interval string, limit int) ([]*models.Kline, error)
	GetNewCoins(ctx context.Context) ([]*models.NewCoin, error)

	// Account operations
	GetWallet(ctx context.Context) (*models.Wallet, error)

	// Order operations
	PlaceOrder(ctx context.Context, order *models.Order) (*models.Order, error)
	CancelOrder(ctx context.Context, orderID, symbol string) error
	GetOrder(ctx context.Context, orderID, symbol string) (*models.Order, error)
	GetOpenOrders(ctx context.Context, symbol string) ([]*models.Order, error)

	// WebSocket operations
	SubscribeToTickers(ctx context.Context, symbols []string, updates chan<- *models.Ticker) error
	UnsubscribeFromTickers(ctx context.Context, symbols []string) error
}

// ExchangeServiceFactory defines a factory for creating exchange services
type ExchangeServiceFactory interface {
	Create() (ExchangeService, error)
}
