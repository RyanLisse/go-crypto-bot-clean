package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// MinimalConfig represents a simplified configuration for the minimal API
type MinimalConfig struct {
	Auth struct {
		Enabled        bool     `mapstructure:"enabled"`
		JWTSecret      string   `mapstructure:"jwt_secret" validate:"required_if=Enabled true"`
		JWTExpiry      int      `mapstructure:"jwt_expiry" validate:"omitempty,min=1"` // in hours
		CookieName     string   `mapstructure:"cookie_name" validate:"required_if=Enabled true"`
		APIKeys        []string `mapstructure:"api_keys"`
		ClerkSecretKey string   `mapstructure:"clerk_secret_key" validate:"required_if=Enabled true"`
		ClerkDomain    string   `mapstructure:"clerk_domain" validate:"required_if=Enabled true"`
	} `mapstructure:"auth" validate:"required"`
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
		Turso                  struct {
			Enabled             bool   `mapstructure:"enabled"`
			URL                 string `mapstructure:"url" validate:"required_if=Enabled true"`
			AuthToken           string `mapstructure:"auth_token" validate:"required_if=Enabled true"`
			SyncEnabled         bool   `mapstructure:"sync_enabled"`
			SyncIntervalSeconds int    `mapstructure:"sync_interval_seconds" validate:"omitempty,min=1"`
			BatchSize           int    `mapstructure:"batch_size" validate:"omitempty,min=1"`
			MaxRetries          int    `mapstructure:"max_retries" validate:"omitempty,min=0"`
			RetryDelaySeconds   int    `mapstructure:"retry_delay_seconds" validate:"omitempty,min=1"`
		} `mapstructure:"turso"`
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
	viper.BindEnv("database.turso.enabled", "TURSO_ENABLED")
	viper.BindEnv("database.turso.url", "TURSO_URL")
	viper.BindEnv("database.turso.auth_token", "TURSO_AUTH_TOKEN")
	viper.BindEnv("database.turso.sync_enabled", "TURSO_SYNC_ENABLED")
	viper.BindEnv("database.turso.sync_interval_seconds", "TURSO_SYNC_INTERVAL_SECONDS")
	viper.BindEnv("database.turso.batch_size", "TURSO_BATCH_SIZE")
	viper.BindEnv("database.turso.max_retries", "TURSO_MAX_RETRIES")
	viper.BindEnv("database.turso.retry_delay_seconds", "TURSO_RETRY_DELAY_SECONDS")
	viper.BindEnv("auth.enabled", "AUTH_ENABLED")
	viper.BindEnv("auth.jwt_secret", "JWT_SECRET")
	viper.BindEnv("auth.jwt_expiry", "JWT_EXPIRY")
	viper.BindEnv("auth.cookie_name", "AUTH_COOKIE_NAME")
	viper.BindEnv("auth.clerk_secret_key", "CLERK_SECRET_KEY")
	viper.BindEnv("auth.clerk_domain", "CLERK_DOMAIN")

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
	viper.SetDefault("database.turso.enabled", false)
	viper.SetDefault("database.turso.sync_enabled", false)
	viper.SetDefault("database.turso.sync_interval_seconds", 60)
	viper.SetDefault("database.turso.batch_size", 100)
	viper.SetDefault("database.turso.max_retries", 3)
	viper.SetDefault("database.turso.retry_delay_seconds", 5)
	viper.SetDefault("auth.enabled", false)
	viper.SetDefault("auth.jwt_expiry", 24)
	viper.SetDefault("auth.cookie_name", "auth_token")

	// Unmarshal configuration
	var config MinimalConfig
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	return &config, nil
}
