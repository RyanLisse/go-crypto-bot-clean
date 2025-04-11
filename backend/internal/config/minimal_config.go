package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// MinimalConfig represents a simplified configuration for the minimal API
type MinimalConfig struct {
	App struct {
		Name        string `mapstructure:"name" validate:"required"`
		Environment string `mapstructure:"environment" validate:"required,oneof=development staging production"`
		LogLevel    string `mapstructure:"log_level" validate:"required,oneof=debug info warn error"`
		Debug       bool   `mapstructure:"debug"`
		Port        string `mapstructure:"port" validate:"required,numeric"`
	} `mapstructure:"app" validate:"required"`

	Logging struct {
		FilePath   string `mapstructure:"file_path" validate:"required"`
		MaxSize    int    `mapstructure:"max_size" validate:"required,min=1"`
		MaxBackups int    `mapstructure:"max_backups" validate:"required,min=0"`
		MaxAge     int    `mapstructure:"max_age" validate:"required,min=1"`
	} `mapstructure:"logging" validate:"required"`

	Database struct {
		Enabled                bool   `mapstructure:"enabled"`
		Path                   string `mapstructure:"path" validate:"required_if=Enabled true"`
		MaxOpenConns           int    `mapstructure:"max_open_conns" validate:"omitempty,min=1"`
		MaxIdleConns           int    `mapstructure:"max_idle_conns" validate:"omitempty,min=1"`
		ConnMaxLifetimeSeconds int    `mapstructure:"conn_max_lifetime_seconds" validate:"omitempty,min=1"`
	} `mapstructure:"database" validate:"required"`
}

// LoadMinimalConfig loads a simplified configuration from environment variables
func LoadMinimalConfig() (*MinimalConfig, error) {
	// Set up viper
	viper.SetConfigType("yaml")

	// Try to load config file if it exists
	configPath := os.Getenv("CONFIG_PATH")
	configFile := os.Getenv("CONFIG_FILE")

	if configPath != "" && configFile != "" {
		viper.SetConfigFile(fmt.Sprintf("%s/%s", configPath, configFile))
		if err := viper.ReadInConfig(); err != nil {
			fmt.Printf("Warning: error reading config file: %v\n", err)
			// Continue even if config file is not found, we'll use environment variables
		}
	}

	// Set up environment variable bindings
	viper.AutomaticEnv()
	viper.SetEnvPrefix("")

	// Bind specific environment variables
	viper.BindEnv("app.name", "APP_NAME")
	viper.BindEnv("app.environment", "ENVIRONMENT")
	viper.BindEnv("app.log_level", "LOG_LEVEL")
	viper.BindEnv("app.debug", "DEBUG")
	viper.BindEnv("app.port", "PORT")
	viper.BindEnv("logging.file_path", "LOG_PATH")
	viper.BindEnv("database.enabled", "DATABASE_ENABLED")
	viper.BindEnv("database.path", "DB_PATH")
	viper.BindEnv("database.max_open_conns", "DB_MAX_OPEN_CONNS")
	viper.BindEnv("database.max_idle_conns", "DB_MAX_IDLE_CONNS")
	viper.BindEnv("database.conn_max_lifetime_seconds", "DB_CONN_MAX_LIFETIME_SECONDS")

	// Set defaults
	viper.SetDefault("app.name", "Go Crypto Bot")
	viper.SetDefault("app.environment", "production")
	viper.SetDefault("app.log_level", "info")
	viper.SetDefault("app.debug", false)
	viper.SetDefault("app.port", "8080")
	viper.SetDefault("logging.file_path", "/app/data/logs")
	viper.SetDefault("logging.max_size", 10)
	viper.SetDefault("logging.max_backups", 3)
	viper.SetDefault("logging.max_age", 30)
	viper.SetDefault("database.enabled", true)
	viper.SetDefault("database.path", "/app/data/minimal.db")
	viper.SetDefault("database.max_open_conns", 10)
	viper.SetDefault("database.max_idle_conns", 5)
	viper.SetDefault("database.conn_max_lifetime_seconds", 300)

	// Unmarshal configuration
	var config MinimalConfig
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	return &config, nil
}
