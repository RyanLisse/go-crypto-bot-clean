package interfaces

import (
	"context"
	"time"

	"github.com/ryanlisse/go-crypto-bot/internal/domain/models"
)

// PriceService defines the interface for price data
type PriceService interface {
	// GetPrice returns the current price for a symbol
	GetPrice(ctx context.Context, symbol string) (float64, error)

	// GetTicker returns the current ticker for a symbol
	GetTicker(ctx context.Context, symbol string) (*models.Ticker, error)

	// GetKlines returns historical kline data
	GetKlines(ctx context.Context, symbol, interval string, limit int) ([]*models.Kline, error)

	// GetPriceHistory returns historical price data
	GetPriceHistory(ctx context.Context, symbol string, startTime, endTime time.Time) ([]float64, error)
}

// OrderService defines the interface for order execution
type OrderService interface {
	// ExecuteOrder executes a trade order
	ExecuteOrder(ctx context.Context, order *models.Order) (*models.Order, error)

	// CancelOrder cancels an existing order
	CancelOrder(ctx context.Context, orderID string) error

	// GetOrderStatus retrieves the status of an order
	GetOrderStatus(ctx context.Context, orderID string) (*models.Order, error)

	// GetOpenOrders retrieves all open orders
	GetOpenOrders(ctx context.Context, symbol string) ([]*models.Order, error)
}

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
