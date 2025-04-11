package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// MinimalConfig represents a simplified configuration for the minimal API
type MinimalConfig struct {
	App struct {
		Name        string `mapstructure:"name"`
		Environment string `mapstructure:"environment"`
		LogLevel    string `mapstructure:"log_level"`
		Debug       bool   `mapstructure:"debug"`
		Port        string `mapstructure:"port"`
	} `mapstructure:"app"`

	Logging struct {
		FilePath   string `mapstructure:"file_path"`
		MaxSize    int    `mapstructure:"max_size"`
		MaxBackups int    `mapstructure:"max_backups"`
		MaxAge     int    `mapstructure:"max_age"`
	} `mapstructure:"logging"`
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

	// Unmarshal configuration
	var config MinimalConfig
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	return &config, nil
}
