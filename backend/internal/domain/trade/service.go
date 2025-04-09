package trade

import (
	"context"
	"time"

	"go-crypto-bot-clean/backend/internal/domain/models"
)

// Service defines the interface for trading operations
type Service interface {
	// Core trading operations
	ExecutePurchase(ctx context.Context, userID int, symbol string, amount float64, options *models.PurchaseOptions) (*models.BoughtCoin, error)
	SellCoin(ctx context.Context, userID int, coin *models.BoughtCoin, amount float64) (*models.Order, error)
	GetPendingOrders(ctx context.Context, userID int) ([]*models.Order, error)
	GetOrderStatus(ctx context.Context, userID int, orderID string) (*models.Order, error)
	CancelOrder(ctx context.Context, userID int, orderID string) error

	// CLI command support
	GetActiveTrades(ctx context.Context, userID int) ([]*models.BoughtCoin, error)
	ExecuteTrade(ctx context.Context, userID int, order *models.Order) (*models.Order, error)
	GetTradeHistory(ctx context.Context, userID int, startTime time.Time, limit int) ([]*models.Order, error)
}
