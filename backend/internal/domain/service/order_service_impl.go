package service

import (
	"context"
	"fmt"
	"strconv"

	"go.uber.org/zap"

	"github.com/ryanlisse/go-crypto-bot/internal/domain/interfaces"
	"github.com/ryanlisse/go-crypto-bot/internal/domain/models"
)

// OrderServiceImpl implements the interfaces.OrderService interface using the TradeService
type OrderServiceImpl struct {
	tradeService interfaces.TradeService
	logger       *zap.Logger
}

// NewOrderService creates a new interfaces.OrderService
func NewOrderService(tradeService interfaces.TradeService, logger *zap.Logger) interfaces.OrderService {
	return &OrderServiceImpl{
		tradeService: tradeService,
		logger:       logger,
	}
}

// ExecuteOrder executes a trade order
func (s *OrderServiceImpl) ExecuteOrder(ctx context.Context, order *models.Order) (*models.Order, error) {
	// Create purchase options
	options := &models.PurchaseOptions{
		OrderType: string(order.Type),
		Price:     order.Price,
	}

	// Execute the order based on side
	if order.Side == "BUY" {
		// For buy orders, use ExecutePurchase
		boughtCoin, err := s.tradeService.ExecutePurchase(ctx, order.Symbol, order.Quantity, options)
		if err != nil {
			return nil, fmt.Errorf("failed to execute buy order: %w", err)
		}

		// Convert BoughtCoin to Order
		return &models.Order{
			ID:        strconv.FormatInt(boughtCoin.ID, 10),
			OrderID:   strconv.FormatInt(boughtCoin.ID, 10),
			Symbol:    boughtCoin.Symbol,
			Side:      "BUY",
			Type:      order.Type,
			Quantity:  boughtCoin.Quantity,
			Price:     boughtCoin.BuyPrice,
			Status:    "FILLED",
			CreatedAt: boughtCoin.BoughtAt,
			Time:      boughtCoin.BoughtAt,
		}, nil
	} else if order.Side == "SELL" {
		// For sell orders, we need to find the BoughtCoin first
		// This is a simplified implementation - in a real system, you'd need to look up the BoughtCoin
		// based on the order details
		boughtCoin := &models.BoughtCoin{
			ID:       1, // Dummy ID
			Symbol:   order.Symbol,
			Quantity: order.Quantity,
		}

		// Execute the sell
		sellOrder, err := s.tradeService.SellCoin(ctx, boughtCoin, order.Quantity)
		if err != nil {
			return nil, fmt.Errorf("failed to execute sell order: %w", err)
		}

		return sellOrder, nil
	}

	return nil, fmt.Errorf("unsupported order side: %s", order.Side)
}

// CancelOrder cancels an existing order
func (s *OrderServiceImpl) CancelOrder(ctx context.Context, orderID string) error {
	return s.tradeService.CancelOrder(ctx, orderID)
}

// GetOrderStatus retrieves the status of an order
func (s *OrderServiceImpl) GetOrderStatus(ctx context.Context, orderID string) (*models.Order, error) {
	return s.tradeService.GetOrderStatus(ctx, orderID)
}

// GetOpenOrders retrieves all open orders
func (s *OrderServiceImpl) GetOpenOrders(ctx context.Context, symbol string) ([]*models.Order, error) {
	return s.tradeService.GetPendingOrders(ctx)
}
