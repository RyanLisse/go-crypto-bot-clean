package mocks

import (
	"context"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/stretchr/testify/mock"
)

// SniperService is a mock implementation of port.SniperService
type SniperService struct {
	mock.Mock
}

// ExecuteSnipe mocks the ExecuteSnipe method
func (m *SniperService) ExecuteSnipe(ctx context.Context, symbol string) (*model.Order, error) {
	args := m.Called(ctx, symbol)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Order), args.Error(1)
}

// ExecuteSnipeWithConfig mocks the ExecuteSnipeWithConfig method
func (m *SniperService) ExecuteSnipeWithConfig(ctx context.Context, symbol string, config *port.SniperConfig) (*model.Order, error) {
	args := m.Called(ctx, symbol, config)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Order), args.Error(1)
}

// PrevalidateSymbol mocks the PrevalidateSymbol method
func (m *SniperService) PrevalidateSymbol(ctx context.Context, symbol string) (bool, error) {
	args := m.Called(ctx, symbol)
	return args.Bool(0), args.Error(1)
}

// GetConfig mocks the GetConfig method
func (m *SniperService) GetConfig() *port.SniperConfig {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*port.SniperConfig)
}

// UpdateConfig mocks the UpdateConfig method
func (m *SniperService) UpdateConfig(config *port.SniperConfig) error {
	args := m.Called(config)
	return args.Error(0)
}

// GetStatus mocks the GetStatus method
func (m *SniperService) GetStatus() string {
	args := m.Called()
	return args.String(0)
}

// Start mocks the Start method
func (m *SniperService) Start() error {
	args := m.Called()
	return args.Error(0)
}

// Stop mocks the Stop method
func (m *SniperService) Stop() error {
	args := m.Called()
	return args.Error(0)
}
