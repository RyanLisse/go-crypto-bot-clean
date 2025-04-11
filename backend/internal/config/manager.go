package config

import (
	"fmt"
	"sync"

	"go.uber.org/zap"
)

// Manager is responsible for managing application configuration
type Manager struct {
	logger        *zap.Logger
	environment   Environment
	minimalLoader *ConfigLoader
	fullLoader    *ConfigLoader
	notifLoader   *ConfigLoader
	minimalConfig *MinimalConfig
	fullConfig    *Config
	notifConfig   *NotificationConfig
	mutex         sync.RWMutex
}

// NewManager creates a new configuration manager
func NewManager(logger *zap.Logger, environment Environment) *Manager {
	return &Manager{
		logger:      logger,
		environment: environment,
	}
}

// LoadMinimalConfig loads the minimal configuration
func (m *Manager) LoadMinimalConfig() (*MinimalConfig, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Create loader if it doesn't exist
	if m.minimalLoader == nil {
		m.minimalLoader = NewConfigLoader(ConfigTypeMinimal, m.environment, m.logger)
	}

	// Load configuration
	config, err := m.minimalLoader.Load()
	if err != nil {
		return nil, err
	}

	// Store configuration
	minimalConfig, ok := config.(*MinimalConfig)
	if !ok {
		return nil, fmt.Errorf("invalid configuration type")
	}

	m.minimalConfig = minimalConfig
	return minimalConfig, nil
}

// LoadFullConfig loads the full configuration
func (m *Manager) LoadFullConfig() (*Config, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Create loader if it doesn't exist
	if m.fullLoader == nil {
		m.fullLoader = NewConfigLoader(ConfigTypeFull, m.environment, m.logger)
	}

	// Load configuration
	config, err := m.fullLoader.Load()
	if err != nil {
		return nil, err
	}

	// Store configuration
	fullConfig, ok := config.(*Config)
	if !ok {
		return nil, fmt.Errorf("invalid configuration type")
	}

	m.fullConfig = fullConfig
	return fullConfig, nil
}

// LoadNotificationConfig loads the notification configuration
func (m *Manager) LoadNotificationConfig() (*NotificationConfig, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Create loader if it doesn't exist
	if m.notifLoader == nil {
		m.notifLoader = NewConfigLoader(ConfigTypeNotification, m.environment, m.logger)
	}

	// Load configuration
	config, err := m.notifLoader.Load()
	if err != nil {
		return nil, err
	}

	// Store configuration
	notifConfig, ok := config.(*NotificationConfig)
	if !ok {
		return nil, fmt.Errorf("invalid configuration type")
	}

	m.notifConfig = notifConfig
	return notifConfig, nil
}

// GetMinimalConfig returns the current minimal configuration
func (m *Manager) GetMinimalConfig() *MinimalConfig {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.minimalConfig
}

// GetFullConfig returns the current full configuration
func (m *Manager) GetFullConfig() *Config {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.fullConfig
}

// GetNotificationConfig returns the current notification configuration
func (m *Manager) GetNotificationConfig() *NotificationConfig {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.notifConfig
}

// EnableReload enables configuration reloading for all loaders
func (m *Manager) EnableReload() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Enable reload for minimal config if it exists
	if m.minimalLoader != nil {
		if err := m.minimalLoader.EnableReload(); err != nil {
			return fmt.Errorf("failed to enable reload for minimal config: %w", err)
		}

		// Listen for changes
		go m.listenForMinimalChanges()
	}

	// Enable reload for full config if it exists
	if m.fullLoader != nil {
		if err := m.fullLoader.EnableReload(); err != nil {
			return fmt.Errorf("failed to enable reload for full config: %w", err)
		}

		// Listen for changes
		go m.listenForFullChanges()
	}

	// Enable reload for notification config if it exists
	if m.notifLoader != nil {
		if err := m.notifLoader.EnableReload(); err != nil {
			return fmt.Errorf("failed to enable reload for notification config: %w", err)
		}

		// Listen for changes
		go m.listenForNotificationChanges()
	}

	return nil
}

// DisableReload disables configuration reloading for all loaders
func (m *Manager) DisableReload() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Disable reload for all loaders
	if m.minimalLoader != nil {
		m.minimalLoader.DisableReload()
	}

	if m.fullLoader != nil {
		m.fullLoader.DisableReload()
	}

	if m.notifLoader != nil {
		m.notifLoader.DisableReload()
	}
}

// listenForMinimalChanges listens for changes to the minimal configuration
func (m *Manager) listenForMinimalChanges() {
	for range m.minimalLoader.ReloadChan() {
		m.mutex.Lock()
		m.minimalConfig = m.minimalLoader.GetMinimalConfig()
		m.mutex.Unlock()
		m.logger.Info("Minimal configuration reloaded")
	}
}

// listenForFullChanges listens for changes to the full configuration
func (m *Manager) listenForFullChanges() {
	for range m.fullLoader.ReloadChan() {
		m.mutex.Lock()
		m.fullConfig = m.fullLoader.GetFullConfig()
		m.mutex.Unlock()
		m.logger.Info("Full configuration reloaded")
	}
}

// listenForNotificationChanges listens for changes to the notification configuration
func (m *Manager) listenForNotificationChanges() {
	for range m.notifLoader.ReloadChan() {
		m.mutex.Lock()
		m.notifConfig = m.notifLoader.GetNotificationConfig()
		m.mutex.Unlock()
		m.logger.Info("Notification configuration reloaded")
	}
}
