package services

import (
	"context"

	"go-crypto-bot-clean/backend/internal/domain/models"
)

// TradingService defines the interface for trading operations
type TradingService interface {
	// Order operations
	PlaceOrder(ctx context.Context, order *models.Order) error
	CancelOrder(ctx context.Context, orderID string) error
	GetOrderStatus(ctx context.Context, orderID string) (*models.Order, error)
	ListOrders(ctx context.Context, symbol string, status models.OrderStatus) ([]*models.Order, error)

	// Position operations
	OpenPosition(ctx context.Context, symbol string, side string, amount float64, price float64) (*models.Position, error)
	ClosePosition(ctx context.Context, positionID string, price float64) error
	GetPosition(ctx context.Context, positionID string) (*models.Position, error)
	ListPositions(ctx context.Context, status models.PositionStatus) ([]*models.Position, error)
	UpdatePositionPrice(ctx context.Context, positionID string, currentPrice float64) error

	// Trade operations
	GetTrades(ctx context.Context, symbol string, limit int) ([]*models.Trade, error)
	GetTradeByID(ctx context.Context, tradeID string) (*models.Trade, error)
}
