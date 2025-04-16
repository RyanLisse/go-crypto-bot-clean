package mocks

import (
	"context"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model/status"
	usecase "github.com/RyanLisse/go-crypto-bot-clean/backend/internal/usecase"
)

// MockStatusUseCase is a mock implementation of the StatusUseCase interface
type MockStatusUseCase struct{}

// Start starts the status use case
func (m *MockStatusUseCase) Start(ctx context.Context) error {
	return nil
}

// Stop stops the status use case
func (m *MockStatusUseCase) Stop() {
	// No-op
}

// GetSystemStatus returns the current system status
func (m *MockStatusUseCase) GetSystemStatus(ctx context.Context) (*status.SystemStatus, error) {
	return status.NewSystemStatus("mock", time.Now()), nil
}

// GetComponentStatus returns the status of a specific component
func (m *MockStatusUseCase) GetComponentStatus(ctx context.Context, name string) (*status.ComponentStatus, error) {
	return status.NewComponentStatus(name, status.StatusRunning), nil
}

// ControlComponent controls a component (start, stop, restart)
func (m *MockStatusUseCase) ControlComponent(ctx context.Context, control status.ProcessControl) (*status.ProcessControlResponse, error) {
	return &status.ProcessControlResponse{
		Component: control.Component,
		Action:    control.Action,
		Success:   true,
		NewStatus: status.StatusRunning,
		Message:   "Mock control action executed",
	}, nil
}

// RegisterProvider registers a status provider
func (m *MockStatusUseCase) RegisterProvider(provider interface{}) {
	// No-op
}

// Ensure MockStatusUseCase implements usecase.StatusUseCase
var _ usecase.StatusUseCase = (*MockStatusUseCase)(nil)
