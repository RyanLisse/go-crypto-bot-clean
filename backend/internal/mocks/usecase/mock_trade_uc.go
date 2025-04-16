package mocks

import (
	"context"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
)

// MockTradeUseCase is a mock implementation of the TradeUseCase interface
type MockTradeUseCase struct{}

// PlaceOrder places a new order
func (m *MockTradeUseCase) PlaceOrder(ctx context.Context, req model.OrderRequest) (*model.Order, error) {
	return &model.Order{
		ID:     "mock-order-id",
		Symbol: req.Symbol,
		Side:   req.Side,
		Type:   req.Type,
		Status: model.OrderStatusNew,
	}, nil
}

// CancelOrder cancels an existing order
func (m *MockTradeUseCase) CancelOrder(ctx context.Context, symbol, orderID string) error {
	return nil
}

// GetOrderStatus gets the current status of an order
func (m *MockTradeUseCase) GetOrderStatus(ctx context.Context, symbol, orderID string) (*model.Order, error) {
	return &model.Order{
		ID:     orderID,
		Symbol: symbol,
		Status: model.OrderStatusFilled,
	}, nil
}

// GetOpenOrders gets all open orders for a symbol
func (m *MockTradeUseCase) GetOpenOrders(ctx context.Context, symbol string) ([]*model.Order, error) {
	return []*model.Order{
		{
			ID:     "mock-open-order-1",
			Symbol: symbol,
			Status: model.OrderStatusNew,
		},
	}, nil
}

// GetOrderHistory gets order history for a symbol with pagination
func (m *MockTradeUseCase) GetOrderHistory(ctx context.Context, symbol string, limit, offset int) ([]*model.Order, error) {
	return []*model.Order{
		{
			ID:     "mock-history-order-1",
			Symbol: symbol,
			Status: model.OrderStatusFilled,
		},
	}, nil
}

// CalculateRequiredQuantity calculates the required quantity for an order based on amount in quote currency
func (m *MockTradeUseCase) CalculateRequiredQuantity(ctx context.Context, symbol string, side model.OrderSide, amount float64) (float64, error) {
	return amount, nil
}

// Ensure MockTradeUseCase implements TradeUseCase
var _ TradeUseCase = (*MockTradeUseCase)(nil)
