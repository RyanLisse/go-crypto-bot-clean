package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Config holds all configuration settings
type Config struct {
	LogLevel string `mapstructure:"log_level"`
	ENV      string `mapstructure:"env"`
	Server   struct {
		Port         int           `mapstructure:"port"`
		Host         string        `mapstructure:"host"`
		ReadTimeout  time.Duration `mapstructure:"read_timeout"`
		WriteTimeout time.Duration `mapstructure:"write_timeout"`
		IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
	} `mapstructure:"server"`
	Database struct {
		Driver   string `mapstructure:"driver"`
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		Name     string `mapstructure:"name"`
		SSLMode  string `mapstructure:"ssl_mode"`
	} `mapstructure:"database"`
	MEXC struct {
		APIKey     string `mapstructure:"api_key"`
		APISecret  string `mapstructure:"api_secret"`
		BaseURL    string `mapstructure:"base_url"`
		WSBaseURL  string `mapstructure:"ws_base_url"`
		UseTestnet bool   `mapstructure:"use_testnet"`
		RateLimit  struct {
			RequestsPerMinute int `mapstructure:"requests_per_minute"`
			BurstSize         int `mapstructure:"burst_size"`
		} `mapstructure:"rate_limit"`
	} `mapstructure:"mexc"`
	AI struct {
		Provider string `mapstructure:"provider"`
		APIKey   string `mapstructure:"api_key"`
		Model    string `mapstructure:"model"`
	} `mapstructure:"ai"`
}

// Load loads configuration from file and environment variables
func Load() (*Config, error) {
	// First load .env file if it exists
	_ = godotenv.Load() // ignore error if .env file doesn't exist

	// Create a new viper instance
	v := viper.New()

	// Set default values
	setDefaults(v)

	// Load config from file
	configFile := getConfigFilePath()
	if configFile != "" {
		v.SetConfigFile(configFile)
		if err := v.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				return nil, fmt.Errorf("error reading config file: %w", err)
			}
			// Config file not found, will use defaults and environment variables
		}
	}

	// Override with environment variables
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Unmarshal config
	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode config: %w", err)
	}

	// Validate config
	if err := validateConfig(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

// setDefaults sets the default values for configuration
func setDefaults(v *viper.Viper) {
	// Server defaults
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.read_timeout", 30*time.Second)
	v.SetDefault("server.write_timeout", 30*time.Second)
	v.SetDefault("server.idle_timeout", 60*time.Second)

	// Environment defaults
	v.SetDefault("env", "development")
	v.SetDefault("log_level", "info")

	// Database defaults
	v.SetDefault("database.driver", "postgres")
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 5432)
	v.SetDefault("database.name", "crypto_bot")
	v.SetDefault("database.ssl_mode", "disable")

	// MEXC defaults
	v.SetDefault("mexc.base_url", "https://api.mexc.com")
	v.SetDefault("mexc.ws_base_url", "wss://wbs.mexc.com/ws")
	v.SetDefault("mexc.use_testnet", false)
	v.SetDefault("mexc.rate_limit.requests_per_minute", 1200)
	v.SetDefault("mexc.rate_limit.burst_size", 10)

	// AI defaults
	v.SetDefault("ai.provider", "gemini")
	v.SetDefault("ai.model", "gemini-pro")
}

// validateConfig validates the configuration
func validateConfig(cfg *Config) error {
	// Validate server port
	if cfg.Server.Port < 1 || cfg.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", cfg.Server.Port)
	}

	// Validate required API keys in production
	if cfg.ENV == "production" {
		if cfg.MEXC.APIKey == "" || cfg.MEXC.APISecret == "" {
			return fmt.Errorf("MEXC API credentials are required in production")
		}
	}

	return nil
}

// getConfigFilePath determines the config file path
func getConfigFilePath() string {
	// Check if CONFIG_FILE environment variable is set
	if configFile := os.Getenv("CONFIG_FILE"); configFile != "" {
		return configFile
	}

	// Check for config files in standard locations
	configName := "config"
	if env := os.Getenv("ENV"); env != "" {
		configName = fmt.Sprintf("config.%s", strings.ToLower(env))
	}

	// Check current directory
	if fileExists(configName + ".yaml") {
		return configName + ".yaml"
	}
	if fileExists(configName + ".yml") {
		return configName + ".yml"
	}

	// Check ./configs directory
	configsDir := "./configs"
	if fileExists(filepath.Join(configsDir, configName+".yaml")) {
		return filepath.Join(configsDir, configName+".yaml")
	}
	if fileExists(filepath.Join(configsDir, configName+".yml")) {
		return filepath.Join(configsDir, configName+".yml")
	}

	return ""
}

// fileExists checks if a file exists
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
