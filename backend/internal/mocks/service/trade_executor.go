package mocks

import (
	"context"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/stretchr/testify/mock"
)

// MockTradeExecutor is a mock implementation for the port.TradeExecutor interface
type MockTradeExecutor struct {
	mock.Mock
}

func (m *MockTradeExecutor) ExecuteOrder(ctx context.Context, request *model.OrderRequest) (*model.OrderResponse, error) {
	args := m.Called(ctx, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.OrderResponse), args.Error(1)
}

func (m *MockTradeExecutor) CancelOrderWithRetry(ctx context.Context, symbol, orderID string) error {
	args := m.Called(ctx, symbol, orderID)
	return args.Error(0)
}

func (m *MockTradeExecutor) GetOrderStatusWithRetry(ctx context.Context, symbol, orderID string) (*model.Order, error) {
	args := m.Called(ctx, symbol, orderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Order), args.Error(1)
}
