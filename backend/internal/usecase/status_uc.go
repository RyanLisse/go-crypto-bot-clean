package usecase

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model/status"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/rs/zerolog"
)

// StatusUseCaseImpl implements the StatusUseCase interface
type StatusUseCaseImpl struct {
	providers       map[string]port.StatusProvider
	controllable    map[string]port.ControllableStatusProvider
	systemInfo      port.SystemInfoProvider
	statusRepo      port.SystemStatusRepository
	notifier        port.StatusNotifier
	logger          *zerolog.Logger
	systemStatus    *status.SystemStatus
	startTime       time.Time
	version         string
	updateInterval  time.Duration
	mu              sync.RWMutex
	stopChan        chan struct{}
	updateTicker    *time.Ticker
	notifyThreshold map[string]status.Status
}

// StatusUseCaseConfig contains configuration for the status use case
type StatusUseCaseConfig struct {
	Version        string
	UpdateInterval time.Duration
}

// NewStatusUseCase creates a new status use case
func NewStatusUseCase(
	systemInfo port.SystemInfoProvider,
	statusRepo port.SystemStatusRepository,
	notifier port.StatusNotifier,
	logger *zerolog.Logger,
	config StatusUseCaseConfig,
) *StatusUseCaseImpl {
	startTime := time.Now()
	updateInterval := config.UpdateInterval
	if updateInterval == 0 {
		updateInterval = 30 * time.Second
	}

	return &StatusUseCaseImpl{
		providers:       make(map[string]port.StatusProvider),
		controllable:    make(map[string]port.ControllableStatusProvider),
		systemInfo:      systemInfo,
		statusRepo:      statusRepo,
		notifier:        notifier,
		logger:          logger,
		systemStatus:    status.NewSystemStatus(config.Version, startTime),
		startTime:       startTime,
		version:         config.Version,
		updateInterval:  updateInterval,
		stopChan:        make(chan struct{}),
		notifyThreshold: make(map[string]status.Status),
	}
}

// RegisterProvider registers a status provider
func (uc *StatusUseCaseImpl) RegisterProvider(provider interface{}) {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	// Check if the provider implements the StatusProvider interface
	statusProvider, ok := provider.(port.StatusProvider)
	if !ok {
		uc.logger.Warn().Interface("provider", provider).Msg("Provider does not implement StatusProvider interface")
		return
	}

	name := statusProvider.GetName()
	uc.providers[name] = statusProvider
	uc.logger.Info().Str("component", name).Msg("Registered status provider")

	// Check if it's also a controllable provider
	if controllable, ok := provider.(port.ControllableStatusProvider); ok {
		uc.controllable[name] = controllable
		uc.logger.Info().Str("component", name).Msg("Registered controllable status provider")
	}

	// Set initial notification threshold
	uc.notifyThreshold[name] = status.StatusUnknown
}

// UnregisterProvider unregisters a status provider
func (uc *StatusUseCaseImpl) UnregisterProvider(name string) {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	delete(uc.providers, name)
	delete(uc.controllable, name)
	delete(uc.notifyThreshold, name)

	// Also remove from system status
	uc.systemStatus.RemoveComponent(name)
	uc.logger.Info().Str("component", name).Msg("Unregistered status provider")
}

// Start starts the status use case
func (uc *StatusUseCaseImpl) Start(ctx context.Context) error {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	// Initialize system status
	if err := uc.updateSystemStatus(ctx); err != nil {
		uc.logger.Warn().Err(err).Msg("Failed to initialize system status, but continuing anyway")
		// Don't return error here to allow the system to continue running
		// even if the database is not available
	}

	// Start periodic updates
	uc.updateTicker = time.NewTicker(uc.updateInterval)
	go func() {
		for {
			select {
			case <-uc.updateTicker.C:
				if err := uc.updateSystemStatus(ctx); err != nil {
					uc.logger.Error().Err(err).Msg("Failed to update system status")
				}
			case <-uc.stopChan:
				uc.updateTicker.Stop()
				return
			case <-ctx.Done():
				uc.updateTicker.Stop()
				return
			}
		}
	}()

	uc.logger.Info().Dur("interval", uc.updateInterval).Msg("Started status monitoring")
	return nil
}

// Stop stops the status use case
func (uc *StatusUseCaseImpl) Stop() {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	if uc.updateTicker != nil {
		close(uc.stopChan)
		uc.updateTicker = nil
		uc.logger.Info().Msg("Stopped status monitoring")
	}
}

// GetSystemStatus returns the current system status
func (uc *StatusUseCaseImpl) GetSystemStatus(ctx context.Context) (*status.SystemStatus, error) {
	uc.mu.RLock()
	defer uc.mu.RUnlock()

	// Update uptime
	uc.systemStatus.Uptime = time.Since(uc.startTime).String()
	uc.systemStatus.LastUpdated = time.Now()

	return uc.systemStatus, nil
}

// GetComponentStatus returns the status of a specific component
func (uc *StatusUseCaseImpl) GetComponentStatus(ctx context.Context, name string) (*status.ComponentStatus, error) {
	uc.mu.RLock()
	defer uc.mu.RUnlock()

	component := uc.systemStatus.GetComponent(name)
	if component == nil {
		return nil, fmt.Errorf("component not found: %s", name)
	}

	return component, nil
}

// ControlComponent controls a component (start, stop, restart)
func (uc *StatusUseCaseImpl) ControlComponent(ctx context.Context, control status.ProcessControl) (*status.ProcessControlResponse, error) {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	provider, ok := uc.controllable[control.Component]
	if !ok {
		return nil, fmt.Errorf("component not found or not controllable: %s", control.Component)
	}

	var err error
	response := &status.ProcessControlResponse{
		Component:   control.Component,
		Action:      control.Action,
		CompletedAt: time.Now(),
	}

	// Apply timeout if specified
	if control.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, control.Timeout)
		defer cancel()
	}

	// Execute the requested action
	switch control.Action {
	case "start":
		err = provider.Start(ctx)
		if err == nil {
			response.Success = true
			response.NewStatus = status.StatusRunning
			response.Message = "Component started successfully"
		} else {
			response.Success = false
			response.NewStatus = status.StatusError
			response.Message = fmt.Sprintf("Failed to start component: %v", err)
		}
	case "stop":
		err = provider.Stop(ctx)
		if err == nil {
			response.Success = true
			response.NewStatus = status.StatusStopped
			response.Message = "Component stopped successfully"
		} else {
			response.Success = false
			response.NewStatus = status.StatusError
			response.Message = fmt.Sprintf("Failed to stop component: %v", err)
		}
	case "restart":
		err = provider.Restart(ctx)
		if err == nil {
			response.Success = true
			response.NewStatus = status.StatusRunning
			response.Message = "Component restarted successfully"
		} else {
			response.Success = false
			response.NewStatus = status.StatusError
			response.Message = fmt.Sprintf("Failed to restart component: %v", err)
		}
	default:
		return nil, fmt.Errorf("unsupported action: %s", control.Action)
	}

	// Update component status
	if componentStatus, err := provider.GetStatus(ctx); err == nil {
		uc.systemStatus.AddComponent(componentStatus)

		// Save updated status
		if uc.statusRepo != nil {
			if err := uc.statusRepo.SaveComponentStatus(ctx, componentStatus); err != nil {
				// Check if it's a "no such table" error
				if err.Error() == "no such table: status_records" {
					// This is expected on first run before migrations are complete
					uc.logger.Debug().Str("component", control.Component).Msg("Status table not yet created, skipping component status save")
				} else {
					uc.logger.Error().Err(err).Str("component", control.Component).Msg("Failed to save component status")
				}
				// Continue anyway
			}
		}
	}

	uc.logger.Info().
		Str("component", control.Component).
		Str("action", control.Action).
		Bool("success", response.Success).
		Msg("Component control action executed")

	return response, nil
}

// updateSystemStatus updates the system status by querying all providers
func (uc *StatusUseCaseImpl) updateSystemStatus(ctx context.Context) error {
	// Update system info
	if uc.systemInfo != nil {
		systemInfo, err := uc.systemInfo.GetSystemInfo(ctx)
		if err != nil {
			uc.logger.Error().Err(err).Msg("Failed to get system info")
		} else {
			uc.systemStatus.SystemInfo = systemInfo
		}
	}

	// Update component statuses
	for name, provider := range uc.providers {
		componentStatus, err := provider.GetStatus(ctx)
		if err != nil {
			uc.logger.Error().Err(err).Str("component", name).Msg("Failed to get component status")

			// Create an error status
			componentStatus = status.NewComponentStatus(name, status.StatusError)
			componentStatus.LastError = err.Error()
		}

		// Check for status changes that need notification
		prevStatus := uc.notifyThreshold[name]
		if componentStatus.Status != prevStatus && uc.notifier != nil {
			// Only notify for significant changes
			if shouldNotify(prevStatus, componentStatus.Status) {
				err := uc.notifier.NotifyStatusChange(ctx, name, prevStatus, componentStatus.Status, componentStatus.Message)
				if err != nil {
					uc.logger.Error().Err(err).Str("component", name).Msg("Failed to send status notification")
				}
				uc.notifyThreshold[name] = componentStatus.Status
			}
		}

		// Update system status
		uc.systemStatus.AddComponent(componentStatus)

		// Save component status
		if uc.statusRepo != nil {
			if err := uc.statusRepo.SaveComponentStatus(ctx, componentStatus); err != nil {
				// Check if it's a "no such table" error
				if err.Error() == "no such table: status_records" {
					// This is expected on first run before migrations are complete
					uc.logger.Debug().Str("component", name).Msg("Status table not yet created, skipping component status save")
				} else {
					uc.logger.Error().Err(err).Str("component", name).Msg("Failed to save component status")
				}
				// Continue anyway
			}
		}
	}

	// Update overall system status
	prevStatus := uc.systemStatus.Status
	uc.systemStatus.UpdateSystemStatus()

	// Notify on system status change
	if prevStatus != uc.systemStatus.Status && uc.notifier != nil {
		if shouldNotify(prevStatus, uc.systemStatus.Status) {
			err := uc.notifier.NotifySystemStatusChange(ctx, prevStatus, uc.systemStatus.Status, "System status changed")
			if err != nil {
				uc.logger.Error().Err(err).Msg("Failed to send system status notification")
			}
		}
	}

	// Save system status
	if uc.statusRepo != nil {
		if err := uc.statusRepo.SaveSystemStatus(ctx, uc.systemStatus); err != nil {
			// Check if it's a "no such table" error
			if err.Error() == "no such table: status_records" {
				// This is expected on first run before migrations are complete
				uc.logger.Debug().Msg("Status table not yet created, skipping status save")
			} else {
				uc.logger.Error().Err(err).Msg("Failed to save system status")
			}
			// Don't return error here to allow the system to continue running
			// even if the database is not available
		}
	}

	return nil
}

// Ensure StatusUseCaseImpl implements StatusUseCase
var _ StatusUseCase = (*StatusUseCaseImpl)(nil)

// shouldNotify determines if a status change should trigger a notification
func shouldNotify(oldStatus, newStatus status.Status) bool {
	// Always notify when transitioning to or from error state
	if oldStatus == status.StatusError || newStatus == status.StatusError {
		return true
	}

	// Always notify when transitioning to warning state
	if newStatus == status.StatusWarning {
		return true
	}

	// Notify when transitioning between running and stopped
	if (oldStatus == status.StatusRunning && newStatus == status.StatusStopped) ||
		(oldStatus == status.StatusStopped && newStatus == status.StatusRunning) {
		return true
	}

	// Don't notify for transitions involving unknown state
	if oldStatus == status.StatusUnknown || newStatus == status.StatusUnknown {
		return false
	}

	// Don't notify for transitions between starting/stopping and their target states
	if (oldStatus == status.StatusStarting && newStatus == status.StatusRunning) ||
		(oldStatus == status.StatusStopping && newStatus == status.StatusStopped) {
		return false
	}

	return false
}
