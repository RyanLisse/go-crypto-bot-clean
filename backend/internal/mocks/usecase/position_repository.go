package mocks

import (
	"context"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/stretchr/testify/mock"
)

// MockPositionRepository is a mock implementation of port.PositionRepository
type MockPositionRepository struct {
	mock.Mock
}

// Create mocks the Create method
func (m *MockPositionRepository) Create(ctx context.Context, position *model.Position) error {
	args := m.Called(ctx, position)
	return args.Error(0)
}

// Update mocks the Update method
func (m *MockPositionRepository) Update(ctx context.Context, position *model.Position) error {
	args := m.Called(ctx, position)
	return args.Error(0)
}

// GetByID mocks the GetByID method
func (m *MockPositionRepository) GetByID(ctx context.Context, id string) (*model.Position, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Position), args.Error(1)
}

// GetByUserID mocks the GetByUserID method
func (m *MockPositionRepository) GetByUserID(ctx context.Context, userID string, page, limit int) ([]*model.Position, error) {
	args := m.Called(ctx, userID, page, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Position), args.Error(1)
}

// GetOpenPositionsByUserID mocks the GetOpenPositionsByUserID method
func (m *MockPositionRepository) GetOpenPositionsByUserID(ctx context.Context, userID string) ([]*model.Position, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Position), args.Error(1)
}

// GetBySymbol mocks the GetBySymbol method
func (m *MockPositionRepository) GetBySymbol(ctx context.Context, symbol string, page, limit int) ([]*model.Position, error) {
	args := m.Called(ctx, symbol, page, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Position), args.Error(1)
}

// GetBySymbolAndUser mocks the GetBySymbolAndUser method
func (m *MockPositionRepository) GetBySymbolAndUser(ctx context.Context, symbol, userID string, page, limit int) ([]*model.Position, error) {
	args := m.Called(ctx, symbol, userID, page, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Position), args.Error(1)
}

// GetActiveByUser mocks the GetActiveByUser method
func (m *MockPositionRepository) GetActiveByUser(ctx context.Context, userID string) ([]*model.Position, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Position), args.Error(1)
}

// Delete mocks the Delete method
func (m *MockPositionRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// Count mocks the Count method
func (m *MockPositionRepository) Count(ctx context.Context, filters map[string]interface{}) (int64, error) {
	args := m.Called(ctx, filters)
	return args.Get(0).(int64), args.Error(1)
}
