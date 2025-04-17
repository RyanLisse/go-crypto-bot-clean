package mocks

import (
	"context"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/stretchr/testify/mock"
)

// MockSniperUseCase is a mock implementation of the SniperUseCase interface
type MockSniperUseCase struct {
	mock.Mock
}

// ExecuteSnipe mocks the ExecuteSnipe method
func (m *MockSniperUseCase) ExecuteSnipe(ctx context.Context, symbol string) (*model.Order, error) {
	args := m.Called(ctx, symbol)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Order), args.Error(1)
}

// ExecuteSnipeWithConfig mocks the ExecuteSnipeWithConfig method
func (m *MockSniperUseCase) ExecuteSnipeWithConfig(ctx context.Context, symbol string, config *port.SniperConfig) (*model.Order, error) {
	args := m.Called(ctx, symbol, config)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Order), args.Error(1)
}

// GetSniperConfig mocks the GetSniperConfig method
func (m *MockSniperUseCase) GetSniperConfig() (*port.SniperConfig, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*port.SniperConfig), args.Error(1)
}

// UpdateSniperConfig mocks the UpdateSniperConfig method
func (m *MockSniperUseCase) UpdateSniperConfig(config *port.SniperConfig) error {
	args := m.Called(config)
	return args.Error(0)
}

// StartSniper mocks the StartSniper method
func (m *MockSniperUseCase) StartSniper() error {
	args := m.Called()
	return args.Error(0)
}

// StopSniper mocks the StopSniper method
func (m *MockSniperUseCase) StopSniper() error {
	args := m.Called()
	return args.Error(0)
}

// GetSniperStatus mocks the GetSniperStatus method
func (m *MockSniperUseCase) GetSniperStatus() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

// SetupAutoSnipe mocks the SetupAutoSnipe method
func (m *MockSniperUseCase) SetupAutoSnipe(enabled bool, config *port.SniperConfig) error {
	args := m.Called(enabled, config)
	return args.Error(0)
}
