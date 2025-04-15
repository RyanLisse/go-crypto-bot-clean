package usecase

import (
	"context"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model/status"
)

// StatusUseCase defines the interface for status operations
type StatusUseCase interface {
	// Start starts the status use case
	Start(ctx context.Context) error
	
	// Stop stops the status use case
	Stop()
	
	// GetSystemStatus returns the current system status
	GetSystemStatus(ctx context.Context) (*status.SystemStatus, error)
	
	// GetComponentStatus returns the status of a specific component
	GetComponentStatus(ctx context.Context, name string) (*status.ComponentStatus, error)
	
	// ControlComponent controls a component (start, stop, restart)
	ControlComponent(ctx context.Context, control status.ProcessControl) (*status.ProcessControlResponse, error)
	
	// RegisterProvider registers a status provider
	RegisterProvider(provider interface{})
}
