package crypto

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
)

// ConfigManager manages secure configuration values
type ConfigManager struct {
	encryptionSvc EncryptionService
	configPath    string
	config        map[string]string
	mu            sync.RWMutex
}

// NewConfigManager creates a new ConfigManager
func NewConfigManager(encryptionSvc EncryptionService, configPath string) (*ConfigManager, error) {
	manager := &ConfigManager{
		encryptionSvc: encryptionSvc,
		configPath:    configPath,
		config:        make(map[string]string),
	}

	// Load config if file exists
	if _, err := os.Stat(configPath); err == nil {
		if err := manager.loadConfig(); err != nil {
			return nil, err
		}
	}

	return manager, nil
}

// loadConfig loads the configuration from the config file
func (m *ConfigManager) loadConfig() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Read encrypted config file
	data, err := os.ReadFile(m.configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	// Decrypt config
	decrypted, err := m.encryptionSvc.Decrypt(data)
	if err != nil {
		return fmt.Errorf("failed to decrypt config: %w", err)
	}

	// Parse JSON
	if err := json.Unmarshal([]byte(decrypted), &m.config); err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	return nil
}

// saveConfig saves the configuration to the config file
func (m *ConfigManager) saveConfig() error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Serialize config to JSON
	data, err := json.Marshal(m.config)
	if err != nil {
		return fmt.Errorf("failed to serialize config: %w", err)
	}

	// Encrypt config
	encrypted, err := m.encryptionSvc.Encrypt(string(data))
	if err != nil {
		return fmt.Errorf("failed to encrypt config: %w", err)
	}

	// Write encrypted config to file
	if err := os.WriteFile(m.configPath, encrypted, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// GetValue gets a configuration value
func (m *ConfigManager) GetValue(key string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	value, ok := m.config[key]
	if !ok {
		return "", errors.New("configuration value not found")
	}

	return value, nil
}

// SetValue sets a configuration value
func (m *ConfigManager) SetValue(key, value string) error {
	m.mu.Lock()
	m.config[key] = value
	m.mu.Unlock()

	return m.saveConfig()
}

// DeleteValue deletes a configuration value
func (m *ConfigManager) DeleteValue(key string) error {
	m.mu.Lock()
	delete(m.config, key)
	m.mu.Unlock()

	return m.saveConfig()
}

// GetAllValues gets all configuration values
func (m *ConfigManager) GetAllValues() map[string]string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Create a copy of the config map
	config := make(map[string]string)
	for k, v := range m.config {
		config[k] = v
	}

	return config
}

// SetMultipleValues sets multiple configuration values
func (m *ConfigManager) SetMultipleValues(values map[string]string) error {
	m.mu.Lock()
	for k, v := range values {
		m.config[k] = v
	}
	m.mu.Unlock()

	return m.saveConfig()
}

// Clear clears all configuration values
func (m *ConfigManager) Clear() error {
	m.mu.Lock()
	m.config = make(map[string]string)
	m.mu.Unlock()

	return m.saveConfig()
}
