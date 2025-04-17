package mocks

import (
	"context"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/stretchr/testify/mock"
)

// SymbolRepository is a mock implementation of port.SymbolRepository
type SymbolRepository struct {
	mock.Mock
}

// Create mocks the Create method
func (m *SymbolRepository) Create(ctx context.Context, symbol *model.Symbol) error {
	args := m.Called(ctx, symbol)
	return args.Error(0)
}

// GetBySymbol mocks the GetBySymbol method
func (m *SymbolRepository) GetBySymbol(ctx context.Context, symbol string) (*model.Symbol, error) {
	args := m.Called(ctx, symbol)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Symbol), args.Error(1)
}

// GetByExchange mocks the GetByExchange method
func (m *SymbolRepository) GetByExchange(ctx context.Context, exchange string) ([]*model.Symbol, error) {
	args := m.Called(ctx, exchange)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Symbol), args.Error(1)
}

// GetAll mocks the GetAll method
func (m *SymbolRepository) GetAll(ctx context.Context) ([]*model.Symbol, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Symbol), args.Error(1)
}

// Update mocks the Update method
func (m *SymbolRepository) Update(ctx context.Context, symbol *model.Symbol) error {
	args := m.Called(ctx, symbol)
	return args.Error(0)
}

// Delete mocks the Delete method
func (m *SymbolRepository) Delete(ctx context.Context, symbol string) error {
	args := m.Called(ctx, symbol)
	return args.Error(0)
}

// GetSymbolsByStatus mocks the GetSymbolsByStatus method
func (m *SymbolRepository) GetSymbolsByStatus(ctx context.Context, status string, limit int, offset int) ([]*model.Symbol, error) {
	args := m.Called(ctx, status, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Symbol), args.Error(1)
}
