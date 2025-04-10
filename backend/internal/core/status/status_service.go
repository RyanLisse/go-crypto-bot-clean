package status

import (
	"context"
	"fmt"
	"runtime"
	"time"
)

// WatcherStatus represents components that can report their running state
type WatcherStatus interface {
	IsRunning() bool
	Start(ctx context.Context) error
	Stop()
}

// StatusProvider is an interface for retrieving component status
type StatusProvider interface {
	GetNewCoinWatcher() WatcherStatus
	GetPositionMonitor() WatcherStatus
}

// SystemInfo contains basic system information
type SystemInfo struct {
	Version   string    `json:"version"`
	GoVersion string    `json:"go_version"`
	StartTime time.Time `json:"start_time"`
	Uptime    string    `json:"uptime"`
}

// ComponentStatus represents the status of a system component
type ComponentStatus struct {
	Name      string `json:"name"`
	IsRunning bool   `json:"is_running"`
	Status    string `json:"status"`
}

// SystemStatus represents the overall system status
type SystemStatus struct {
	SystemInfo    SystemInfo        `json:"system_info"`
	Components    []ComponentStatus `json:"components"`
	OverallStatus string            `json:"overall_status"`
}

// Service defines the interface for status operations used by handlers
type Service interface {
	GetStatus() (*SystemStatus, error)
	StartProcesses(ctx context.Context) error
	StopProcesses() error
}

// StatusService provides methods for checking and controlling system status
type StatusService struct {
	provider   StatusProvider
	startTime  time.Time
	version    string
	components map[string]WatcherStatus
}

// NewStatusService creates a new status service
func NewStatusService(provider StatusProvider, version string) *StatusService {
	return &StatusService{
		provider:   provider,
		startTime:  time.Now(),
		version:    version,
		components: make(map[string]WatcherStatus),
	}
}

// GetStatus returns the current system status
func (s *StatusService) GetStatus() (*SystemStatus, error) {
	// Get component statuses
	components := []ComponentStatus{}

	// New Coin Watcher
	newCoinWatcher := s.provider.GetNewCoinWatcher()
	components = append(components, ComponentStatus{
		Name:      "NewCoinWatcher",
		IsRunning: newCoinWatcher.IsRunning(),
		Status:    s.getStatusString(newCoinWatcher.IsRunning()),
	})

	// Position Monitor
	positionMonitor := s.provider.GetPositionMonitor()
	components = append(components, ComponentStatus{
		Name:      "PositionMonitor",
		IsRunning: positionMonitor.IsRunning(),
		Status:    s.getStatusString(positionMonitor.IsRunning()),
	})

	// Calculate uptime
	uptime := time.Since(s.startTime).Round(time.Second).String()

	// Determine overall status
	overallStatus := "healthy"
	for _, component := range components {
		if !component.IsRunning {
			overallStatus = "degraded"
			break
		}
	}

	return &SystemStatus{
		SystemInfo: SystemInfo{
			Version:   s.version,
			GoVersion: runtime.Version(),
			StartTime: s.startTime,
			Uptime:    uptime,
		},
		Components:    components,
		OverallStatus: overallStatus,
	}, nil
}

// StartProcesses starts all system processes
// StartProcesses starts all system processes (implements Service interface)
func (s *StatusService) StartProcesses(ctx context.Context) error {
	// Start New Coin Watcher if not running
	newCoinWatcher := s.provider.GetNewCoinWatcher()
	if !newCoinWatcher.IsRunning() {
		if err := newCoinWatcher.Start(ctx); err != nil {
			return fmt.Errorf("failed to start NewCoinWatcher: %w", err)
		}
	}

	// Start Position Monitor if not running
	positionMonitor := s.provider.GetPositionMonitor()
	if !positionMonitor.IsRunning() {
		if err := positionMonitor.Start(ctx); err != nil {
			return fmt.Errorf("failed to start PositionMonitor: %w", err)
		}
	}

	return nil // Return nil error on success
}

// StopProcesses stops all system processes
// StopProcesses stops all system processes (implements Service interface)
func (s *StatusService) StopProcesses() error {
	// Stop New Coin Watcher if running
	newCoinWatcher := s.provider.GetNewCoinWatcher()
	if newCoinWatcher.IsRunning() {
		newCoinWatcher.Stop()
	}

	// Stop Position Monitor if running
	positionMonitor := s.provider.GetPositionMonitor()
	if positionMonitor.IsRunning() {
		positionMonitor.Stop()
	}

	// TODO: Should StopProcesses return an error if stopping fails?
	return nil // Return nil error on success
}

// getStatusString returns a human-readable status string
func (s *StatusService) getStatusString(isRunning bool) string {
	if isRunning {
		return "running"
	}
	return "stopped"
}
