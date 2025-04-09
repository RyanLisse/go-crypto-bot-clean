package interfaces

import (
	"context"
	"time"

	"go-crypto-bot-clean/backend/internal/domain/models"
)

// TradeService defines the interface for trading operations
type TradeService interface {
	// Core trading operations
	ExecutePurchase(ctx context.Context, symbol string, amount float64, options *models.PurchaseOptions) (*models.BoughtCoin, error)
	SellCoin(ctx context.Context, coin *models.BoughtCoin, amount float64) (*models.Order, error)
	GetPendingOrders(ctx context.Context) ([]*models.Order, error)
	GetOrderStatus(ctx context.Context, orderID string) (*models.Order, error)
	CancelOrder(ctx context.Context, orderID string) error

	// CLI command support
	GetActiveTrades(ctx context.Context) ([]*models.BoughtCoin, error)
	ExecuteTrade(ctx context.Context, order *models.Order) (*models.Order, error)
	GetTradeHistory(ctx context.Context, startTime time.Time, limit int) ([]*models.Order, error)
}
