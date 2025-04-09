package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	// Server configuration
	Server struct {
		Host string `mapstructure:"host"`
		Port string `mapstructure:"port"`
	} `mapstructure:"server"`

	// Database configuration
	Database struct {
		URL string `mapstructure:"url"`
	} `mapstructure:"database"`

	// API configuration
	API struct {
		BasePath string `mapstructure:"base_path"`
	} `mapstructure:"api"`

	// Logging configuration
	LogLevel string `mapstructure:"log_level"`
}

// Load loads the configuration from environment variables
func Load() (*Config, error) {
	// Create a new viper instance
	v := viper.New()

	// Set default values
	setDefaults(v)

	// Read environment variables
	v.AutomaticEnv()
	v.SetEnvPrefix("APP")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Override with environment variables
	if port := os.Getenv("PORT"); port != "" {
		v.Set("server.port", port)
	}

	// Database URL (Railway provides this)
	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		v.Set("database.url", dbURL)
	}

	// Unmarshal configuration
	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

// setDefaults sets default values for configuration
func setDefaults(v *viper.Viper) {
	// Server defaults
	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.port", "8080")

	// Database defaults
	v.SetDefault("database.url", "sqlite3://crypto-bot.db")

	// API defaults
	v.SetDefault("api.base_path", "/api/v1")

	// Logging defaults
	v.SetDefault("log_level", "info")
}
