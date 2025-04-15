package port

import (
	"context"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model/status"
)

// StatusProvider defines the interface for components that can provide status information
type StatusProvider interface {
	// GetStatus returns the current status of the component
	GetStatus(ctx context.Context) (*status.ComponentStatus, error)
	// GetName returns the name of the component
	GetName() string
	// IsRunning returns true if the component is running
	IsRunning() bool
}

// ControllableStatusProvider extends StatusProvider with control capabilities
type ControllableStatusProvider interface {
	StatusProvider
	// Start starts the component
	Start(ctx context.Context) error
	// Stop stops the component
	Stop(ctx context.Context) error
	// Restart restarts the component
	Restart(ctx context.Context) error
}

// SystemStatusRepository defines the interface for storing and retrieving system status
type SystemStatusRepository interface {
	// SaveSystemStatus saves the current system status
	SaveSystemStatus(ctx context.Context, status *status.SystemStatus) error
	// GetSystemStatus retrieves the current system status
	GetSystemStatus(ctx context.Context) (*status.SystemStatus, error)
	// SaveComponentStatus saves a component status
	SaveComponentStatus(ctx context.Context, componentStatus *status.ComponentStatus) error
	// GetComponentStatus retrieves a component status by name
	GetComponentStatus(ctx context.Context, name string) (*status.ComponentStatus, error)
	// GetComponentHistory retrieves historical status for a component
	GetComponentHistory(ctx context.Context, name string, limit int) ([]*status.ComponentStatus, error)
}

// SystemInfoProvider defines the interface for providing system resource information
type SystemInfoProvider interface {
	// GetSystemInfo returns the current system resource information
	GetSystemInfo(ctx context.Context) (*status.SystemInfo, error)
}

// StatusNotifier defines the interface for sending status notifications
type StatusNotifier interface {
	// NotifyStatusChange sends a notification about a status change
	NotifyStatusChange(ctx context.Context, component string, oldStatus, newStatus status.Status, message string) error
	// NotifySystemStatusChange sends a notification about a system status change
	NotifySystemStatusChange(ctx context.Context, oldStatus, newStatus status.Status, message string) error
}
