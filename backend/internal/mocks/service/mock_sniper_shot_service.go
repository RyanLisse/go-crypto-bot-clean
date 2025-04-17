package mocks

import (
	"context"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/service"
	"github.com/stretchr/testify/mock"
)

// MockSniperShotService is a mock implementation of the SniperShotService
type MockSniperShotService struct {
	mock.Mock
}

// ExecuteSniper mocks the ExecuteSniper method
func (m *MockSniperShotService) ExecuteSniper(ctx context.Context, req *service.SniperShotRequest) (*service.SniperShotResult, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.SniperShotResult), args.Error(1)
}

// CancelSniper mocks the CancelSniper method
func (m *MockSniperShotService) CancelSniper(ctx context.Context, symbol, orderID string) error {
	args := m.Called(ctx, symbol, orderID)
	return args.Error(0)
}

// GetSniperOrderStatus mocks the GetSniperOrderStatus method
func (m *MockSniperShotService) GetSniperOrderStatus(ctx context.Context, symbol, orderID string) (*model.Order, error) {
	args := m.Called(ctx, symbol, orderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Order), args.Error(1)
}
