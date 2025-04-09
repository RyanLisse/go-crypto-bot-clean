package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// TradingConfig represents trading-specific configuration
type TradingConfig struct {
	DefaultSymbol    string    `mapstructure:"default_symbol"`
	DefaultOrderType string    `mapstructure:"default_order_type"`
	DefaultQuantity  float64   `mapstructure:"default_quantity"`
	StopLossPercent  float64   `mapstructure:"stop_loss_percent"`
	TakeProfitLevels []float64 `mapstructure:"take_profit_levels"`
	SellPercentages  []float64 `mapstructure:"sell_percentages"`
}

// Config represents the application configuration
type Config struct {
	App struct {
		Name        string `mapstructure:"name"`
		Environment string `mapstructure:"environment"`
		LogLevel    string `mapstructure:"log_level"`
		Debug       bool   `mapstructure:"debug"`
	} `mapstructure:"app"`

	Auth struct {
		Enabled    bool   `mapstructure:"enabled"`
		JWTSecret  string `mapstructure:"jwt_secret"`
		JWTExpiry  int    `mapstructure:"jwt_expiry"` // in hours
		CookieName string `mapstructure:"cookie_name"`
		// For API key auth
		APIKeys []string `mapstructure:"api_keys"`
		// Clerk configuration
		ClerkSecretKey string `mapstructure:"clerk_secret_key"`
	} `mapstructure:"auth"`

	Mexc struct {
		APIKey       string `mapstructure:"api_key"`
		SecretKey    string `mapstructure:"secret_key"`
		BaseURL      string `mapstructure:"base_url"`
		WebsocketURL string `mapstructure:"websocket_url"`
	} `mapstructure:"mexc"`

	Gemini struct {
		APIKey string `mapstructure:"api_key"`
	} `mapstructure:"gemini"`

	Reporting struct {
		Interval int `mapstructure:"interval"` // in minutes
	} `mapstructure:"reporting"`

	WebSocket struct {
		ReconnectDelay       time.Duration `mapstructure:"reconnect_delay"`
		MaxReconnectAttempts int           `mapstructure:"max_reconnect_attempts"`
		PingInterval         time.Duration `mapstructure:"ping_interval"`
		AutoReconnect        bool          `mapstructure:"auto_reconnect"`
	} `mapstructure:"websocket"`

	ConnectionRateLimiter struct {
		RequestsPerSecond float64 `mapstructure:"requests_per_second"`
		BurstCapacity     int     `mapstructure:"burst_capacity"`
	} `mapstructure:"connection_rate_limiter"`

	SubscriptionRateLimiter struct {
		RequestsPerSecond float64 `mapstructure:"requests_per_second"`
		BurstCapacity     int     `mapstructure:"burst_capacity"`
	} `mapstructure:"subscription_rate_limiter"`

	Trading TradingConfig `mapstructure:"trading"`

	Logging struct {
		FilePath   string `mapstructure:"file_path"`
		MaxSize    int    `mapstructure:"max_size"`
		MaxBackups int    `mapstructure:"max_backups"`
		MaxAge     int    `mapstructure:"max_age"`
	} `mapstructure:"logging"`

	Database struct {
		Type                   string `mapstructure:"type"`
		Path                   string `mapstructure:"path"`
		MaxOpenConns           int    `mapstructure:"maxOpenConns"`
		MaxIdleConns           int    `mapstructure:"maxIdleConns"`
		ConnMaxLifetimeSeconds int    `mapstructure:"connMaxLifetimeSeconds"`
		Turso                  struct {
			Enabled             bool   `mapstructure:"enabled"`
			URL                 string `mapstructure:"url"`
			AuthToken           string `mapstructure:"authToken"`
			SyncEnabled         bool   `mapstructure:"syncEnabled"`
			SyncIntervalSeconds int    `mapstructure:"syncIntervalSeconds"`
		} `mapstructure:"turso"`
		ShadowMode bool `mapstructure:"shadowMode"`
	} `mapstructure:"database"`
}

// LoadConfig loads the configuration from a YAML file and environment variables
func LoadConfig(path string) (*Config, error) {
	// Set configuration file type
	viper.SetConfigType("yaml")

	// Set configuration file path
	viper.SetConfigFile(path)

	// Read configuration
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Warning: error reading config file: %v\n", err)
		// Continue even if config file is not found, we'll use environment variables
	}

	// Set up environment variable bindings
	viper.AutomaticEnv()
	viper.SetEnvPrefix("")

	// Bind specific environment variables
	viper.BindEnv("mexc.api_key", "MEXC_API_KEY")
	viper.BindEnv("mexc.secret_key", "MEXC_SECRET_KEY")
	viper.BindEnv("mexc.base_url", "MEXC_BASE_URL")
	viper.BindEnv("mexc.websocket_url", "MEXC_WEBSOCKET_URL")

	viper.BindEnv("database.turso.enabled", "TURSO_ENABLED")
	viper.BindEnv("database.turso.url", "TURSO_URL")
	viper.BindEnv("database.turso.authToken", "TURSO_AUTH_TOKEN")
	viper.BindEnv("database.turso.syncEnabled", "TURSO_SYNC_ENABLED")
	viper.BindEnv("database.turso.syncIntervalSeconds", "TURSO_SYNC_INTERVAL_SECONDS")

	viper.BindEnv("app.log_level", "LOG_LEVEL")
	viper.BindEnv("app.environment", "ENVIRONMENT")
	viper.BindEnv("database.path", "DB_PATH")
	viper.BindEnv("logging.file_path", "LOG_PATH")

	// Add Clerk environment variable binding
	viper.BindEnv("auth.clerk_secret_key", "CLERK_SECRET_KEY")

	// Unmarshal configuration
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	return &config, nil
}
