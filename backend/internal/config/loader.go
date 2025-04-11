package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// ConfigType represents the type of configuration to load
type ConfigType string

const (
	// ConfigTypeMinimal loads the minimal configuration
	ConfigTypeMinimal ConfigType = "minimal"
	// ConfigTypeFull loads the full configuration
	ConfigTypeFull ConfigType = "full"
	// ConfigTypeNotification loads the notification configuration
	ConfigTypeNotification ConfigType = "notification"
)

// Environment represents the application environment
type Environment string

const (
	// EnvironmentDevelopment is the development environment
	EnvironmentDevelopment Environment = "development"
	// EnvironmentStaging is the staging environment
	EnvironmentStaging Environment = "staging"
	// EnvironmentProduction is the production environment
	EnvironmentProduction Environment = "production"
)

// ConfigLoader is responsible for loading and managing configuration
type ConfigLoader struct {
	configType    ConfigType
	environment   Environment
	configPath    string
	configFile    string
	viper         *viper.Viper
	logger        *zap.Logger
	validate      *validator.Validate
	config        interface{}
	reloadEnabled bool
	reloadMutex   sync.RWMutex
	reloadChan    chan struct{}
	watcher       *fsnotify.Watcher
}

// NewConfigLoader creates a new configuration loader
func NewConfigLoader(configType ConfigType, environment Environment, logger *zap.Logger) *ConfigLoader {
	// Default config path and file
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./configs"
	}

	// Determine config file based on type and environment
	var configFile string
	switch configType {
	case ConfigTypeMinimal:
		configFile = "config.minimal.yaml"
	case ConfigTypeFull:
		if environment == EnvironmentDevelopment {
			configFile = "config.dev.yaml"
		} else if environment == EnvironmentStaging {
			configFile = "config.staging.yaml"
		} else {
			configFile = "config.yaml"
		}
	case ConfigTypeNotification:
		configFile = "notification.yaml"
	default:
		configFile = "config.yaml"
	}

	// Override with environment variable if set
	if envConfigFile := os.Getenv("CONFIG_FILE"); envConfigFile != "" {
		configFile = envConfigFile
	}

	return &ConfigLoader{
		configType:    configType,
		environment:   environment,
		configPath:    configPath,
		configFile:    configFile,
		viper:         viper.New(),
		logger:        logger,
		validate:      validator.New(),
		reloadEnabled: false,
		reloadChan:    make(chan struct{}),
	}
}

// Load loads the configuration
func (cl *ConfigLoader) Load() (interface{}, error) {
	cl.reloadMutex.Lock()
	defer cl.reloadMutex.Unlock()

	// Set up viper
	cl.viper.SetConfigType("yaml")

	// Try to load config file if it exists
	configFilePath := filepath.Join(cl.configPath, cl.configFile)
	cl.viper.SetConfigFile(configFilePath)

	if err := cl.viper.ReadInConfig(); err != nil {
		cl.logger.Warn("Error reading config file", zap.Error(err), zap.String("path", configFilePath))
		// Continue even if config file is not found, we'll use environment variables
	} else {
		cl.logger.Info("Loaded configuration file", zap.String("path", configFilePath))
	}

	// Set up environment variable bindings
	cl.viper.AutomaticEnv()
	cl.viper.SetEnvPrefix("")

	// Set up environment variable replacements
	cl.viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Load configuration based on type
	var config interface{}
	var err error

	switch cl.configType {
	case ConfigTypeMinimal:
		config, err = cl.loadMinimalConfig()
	case ConfigTypeFull:
		config, err = cl.loadFullConfig()
	case ConfigTypeNotification:
		config, err = cl.loadNotificationConfig()
	default:
		return nil, fmt.Errorf("unknown config type: %s", cl.configType)
	}

	if err != nil {
		return nil, err
	}

	// Validate configuration
	if err := cl.validate.Struct(config); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	cl.config = config
	return config, nil
}

// loadMinimalConfig loads the minimal configuration
func (cl *ConfigLoader) loadMinimalConfig() (*MinimalConfig, error) {
	// Bind specific environment variables
	cl.viper.BindEnv("app.name", "APP_NAME")
	cl.viper.BindEnv("app.environment", "ENVIRONMENT")
	cl.viper.BindEnv("app.log_level", "LOG_LEVEL")
	cl.viper.BindEnv("app.debug", "DEBUG")
	cl.viper.BindEnv("app.port", "PORT")
	cl.viper.BindEnv("logging.file_path", "LOG_PATH")

	// Set defaults
	cl.viper.SetDefault("app.name", "Go Crypto Bot")
	cl.viper.SetDefault("app.environment", string(cl.environment))
	cl.viper.SetDefault("app.log_level", "info")
	cl.viper.SetDefault("app.debug", false)
	cl.viper.SetDefault("app.port", "8080")
	cl.viper.SetDefault("logging.file_path", "/app/data/logs")
	cl.viper.SetDefault("logging.max_size", 10)
	cl.viper.SetDefault("logging.max_backups", 3)
	cl.viper.SetDefault("logging.max_age", 30)

	// Unmarshal configuration
	var config MinimalConfig
	if err := cl.viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	return &config, nil
}

// loadFullConfig loads the full configuration
func (cl *ConfigLoader) loadFullConfig() (*Config, error) {
	// Bind specific environment variables
	cl.viper.BindEnv("mexc.api_key", "MEXC_API_KEY")
	cl.viper.BindEnv("mexc.secret_key", "MEXC_SECRET_KEY")
	cl.viper.BindEnv("mexc.base_url", "MEXC_BASE_URL")
	cl.viper.BindEnv("mexc.websocket_url", "MEXC_WEBSOCKET_URL")

	cl.viper.BindEnv("database.turso.enabled", "TURSO_ENABLED")
	cl.viper.BindEnv("database.turso.url", "TURSO_URL")
	cl.viper.BindEnv("database.turso.authToken", "TURSO_AUTH_TOKEN")
	cl.viper.BindEnv("database.turso.syncEnabled", "TURSO_SYNC_ENABLED")
	cl.viper.BindEnv("database.turso.syncIntervalSeconds", "TURSO_SYNC_INTERVAL_SECONDS")

	cl.viper.BindEnv("app.log_level", "LOG_LEVEL")
	cl.viper.BindEnv("app.environment", "ENVIRONMENT")
	cl.viper.BindEnv("database.path", "DB_PATH")
	cl.viper.BindEnv("logging.file_path", "LOG_PATH")

	// Add Clerk environment variable binding
	cl.viper.BindEnv("auth.clerk_secret_key", "CLERK_SECRET_KEY")

	// Add Gemini environment variable binding
	cl.viper.BindEnv("gemini.api_key", "GEMINI_API_KEY")

	// Set defaults based on environment
	cl.viper.SetDefault("app.environment", string(cl.environment))
	
	if cl.environment == EnvironmentDevelopment {
		cl.viper.SetDefault("app.debug", true)
		cl.viper.SetDefault("logging.file_path", "./logs")
		cl.viper.SetDefault("database.path", "./data/trading.db")
	} else {
		cl.viper.SetDefault("app.debug", false)
		cl.viper.SetDefault("logging.file_path", "/app/data/logs")
		cl.viper.SetDefault("database.path", "/app/data/trading.db")
	}

	// Unmarshal configuration
	var config Config
	if err := cl.viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	return &config, nil
}

// loadNotificationConfig loads the notification configuration
func (cl *ConfigLoader) loadNotificationConfig() (*NotificationConfig, error) {
	// Bind specific environment variables
	cl.viper.BindEnv("notification.providers.telegram.token", "TELEGRAM_BOT_TOKEN")
	cl.viper.BindEnv("notification.providers.slack.token", "SLACK_BOT_TOKEN")

	// Unmarshal configuration
	var config NotificationConfig
	if err := cl.viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	return &config, nil
}

// EnableReload enables configuration reloading
func (cl *ConfigLoader) EnableReload() error {
	cl.reloadMutex.Lock()
	defer cl.reloadMutex.Unlock()

	if cl.reloadEnabled {
		return nil
	}

	// Create file watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("failed to create file watcher: %w", err)
	}

	// Watch config file
	configFilePath := filepath.Join(cl.configPath, cl.configFile)
	if err := watcher.Add(configFilePath); err != nil {
		watcher.Close()
		return fmt.Errorf("failed to watch config file: %w", err)
	}

	cl.watcher = watcher
	cl.reloadEnabled = true

	// Start watching for changes
	go cl.watchConfigChanges()

	cl.logger.Info("Configuration reloading enabled", zap.String("path", configFilePath))
	return nil
}

// DisableReload disables configuration reloading
func (cl *ConfigLoader) DisableReload() {
	cl.reloadMutex.Lock()
	defer cl.reloadMutex.Unlock()

	if !cl.reloadEnabled {
		return
	}

	if cl.watcher != nil {
		cl.watcher.Close()
		cl.watcher = nil
	}

	cl.reloadEnabled = false
	cl.logger.Info("Configuration reloading disabled")
}

// watchConfigChanges watches for configuration file changes
func (cl *ConfigLoader) watchConfigChanges() {
	for {
		select {
		case event, ok := <-cl.watcher.Events:
			if !ok {
				return
			}

			if event.Op&fsnotify.Write == fsnotify.Write {
				cl.logger.Info("Configuration file changed, reloading", zap.String("path", event.Name))
				
				// Debounce reloads (wait a bit to ensure file is fully written)
				time.Sleep(100 * time.Millisecond)
				
				if _, err := cl.Load(); err != nil {
					cl.logger.Error("Failed to reload configuration", zap.Error(err))
				} else {
					// Notify subscribers
					select {
					case cl.reloadChan <- struct{}{}:
					default:
						// Non-blocking send
					}
				}
			}
		case err, ok := <-cl.watcher.Errors:
			if !ok {
				return
			}
			cl.logger.Error("Error watching config file", zap.Error(err))
		}
	}
}

// ReloadChan returns a channel that is notified when configuration is reloaded
func (cl *ConfigLoader) ReloadChan() <-chan struct{} {
	return cl.reloadChan
}

// GetConfig returns the current configuration
func (cl *ConfigLoader) GetConfig() interface{} {
	cl.reloadMutex.RLock()
	defer cl.reloadMutex.RUnlock()
	return cl.config
}

// GetMinimalConfig returns the current minimal configuration
func (cl *ConfigLoader) GetMinimalConfig() *MinimalConfig {
	cl.reloadMutex.RLock()
	defer cl.reloadMutex.RUnlock()
	
	if config, ok := cl.config.(*MinimalConfig); ok {
		return config
	}
	return nil
}

// GetFullConfig returns the current full configuration
func (cl *ConfigLoader) GetFullConfig() *Config {
	cl.reloadMutex.RLock()
	defer cl.reloadMutex.RUnlock()
	
	if config, ok := cl.config.(*Config); ok {
		return config
	}
	return nil
}

// GetNotificationConfig returns the current notification configuration
func (cl *ConfigLoader) GetNotificationConfig() *NotificationConfig {
	cl.reloadMutex.RLock()
	defer cl.reloadMutex.RUnlock()
	
	if config, ok := cl.config.(*NotificationConfig); ok {
		return config
	}
	return nil
}
