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
	LogLevel      string              `mapstructure:"log_level"`
	ENV           string              `mapstructure:"env"`
	Version       string              `mapstructure:"version"`
	Notifications Notifications       `mapstructure:"notifications"`
	Auth          Auth                `mapstructure:"auth"`
	RateLimit     RateLimitConfig     `mapstructure:"rate_limit"`
	CSRF          CSRFConfig          `mapstructure:"csrf"`
	SecureHeaders SecureHeadersConfig `mapstructure:"secure_headers"`
	InfuraAPIKey  string              `mapstructure:"infura_api_key"`
	Server        struct {
		Port               int           `mapstructure:"port"`
		Host               string        `mapstructure:"host"`
		ReadTimeout        time.Duration `mapstructure:"read_timeout"`
		WriteTimeout       time.Duration `mapstructure:"write_timeout"`
		IdleTimeout        time.Duration `mapstructure:"idle_timeout"`
		FrontendURL        string        `mapstructure:"frontend_url"`
		CORSAllowedOrigins []string      `mapstructure:"cors_allowed_origins"`
	} `mapstructure:"server"`
	Database struct {
		Driver   string `mapstructure:"driver"`
		Path     string `mapstructure:"path"`
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		Name     string `mapstructure:"name"`
		SSLMode  string `mapstructure:"ssl_mode"`
		Turso    struct {
			Enabled   bool   `mapstructure:"enabled"`
			URL       string `mapstructure:"url"`
			AuthToken string `mapstructure:"auth_token"`
		} `mapstructure:"turso"`
	} `mapstructure:"database"`
	Market struct {
		Cache struct {
			TickerTTL    int `mapstructure:"ticker_ttl"`
			CandleTTL    int `mapstructure:"candle_ttl"`
			OrderbookTTL int `mapstructure:"orderbook_ttl"`
		} `mapstructure:"cache"`
	} `mapstructure:"market"`
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
		Provider     string  `mapstructure:"provider"`
		APIKey       string  `mapstructure:"api_key"`
		Model        string  `mapstructure:"model"`
		GeminiAPIKey string  `mapstructure:"gemini_api_key"`
		GeminiModel  string  `mapstructure:"gemini_model"`
		SystemPrompt string  `mapstructure:"system_prompt"`
		Temperature  float32 `mapstructure:"temperature"`
		TopP         float32 `mapstructure:"top_p"`
		TopK         int32   `mapstructure:"top_k"`
		MaxTokens    int32   `mapstructure:"max_tokens"`
	} `mapstructure:"ai"`
}

// Auth holds authentication configuration
type Auth struct {
	Enabled           bool          `mapstructure:"enabled"`
	Provider          string        `mapstructure:"provider"` // "clerk", "jwt", etc.
	ClerkAPIKey       string        `mapstructure:"clerk_api_key"`
	ClerkSecretKey    string        `mapstructure:"clerk_secret_key"`
	ClerkJWTPublicKey string        `mapstructure:"clerk_jwt_public_key"`
	ClerkJWTTemplate  string        `mapstructure:"clerk_jwt_template"`
	JWTSecret         string        `mapstructure:"jwt_secret"`
	TokenDuration     time.Duration `mapstructure:"token_duration"`
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

// Notifications holds notification configuration
type Notifications struct {
	Email   EmailNotification   `mapstructure:"email"`
	Webhook WebhookNotification `mapstructure:"webhook"`
}

// EmailNotification holds email notification configuration
type EmailNotification struct {
	Enabled       bool     `mapstructure:"enabled"`
	SMTPServer    string   `mapstructure:"smtp_server"`
	SMTPPort      int      `mapstructure:"smtp_port"`
	Username      string   `mapstructure:"username"`
	Password      string   `mapstructure:"password"`
	FromAddress   string   `mapstructure:"from_address"`
	ToAddresses   []string `mapstructure:"to_addresses"`
	MinLevel      string   `mapstructure:"min_level"`
	SubjectPrefix string   `mapstructure:"subject_prefix"`
}

// WebhookNotification holds webhook notification configuration
type WebhookNotification struct {
	Enabled   bool              `mapstructure:"enabled"`
	URL       string            `mapstructure:"url"`
	Method    string            `mapstructure:"method"`
	Headers   map[string]string `mapstructure:"headers"`
	MinLevel  string            `mapstructure:"min_level"`
	Timeout   time.Duration     `mapstructure:"timeout"`
	BatchSize int               `mapstructure:"batch_size"`
}

// setDefaults sets the default values for configuration
func setDefaults(v *viper.Viper) {
	// Server defaults
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.read_timeout", 30*time.Second)
	v.SetDefault("server.write_timeout", 30*time.Second)
	v.SetDefault("server.idle_timeout", 60*time.Second)
	v.SetDefault("server.frontend_url", "http://localhost:3000")
	v.SetDefault("server.cors_allowed_origins", []string{"http://localhost:3000"})

	// Environment defaults
	v.SetDefault("env", "development")
	v.SetDefault("log_level", "info")
	v.SetDefault("version", "1.0.0")

	// Auth defaults
	v.SetDefault("auth.enabled", true)
	v.SetDefault("auth.provider", "clerk")
	v.SetDefault("auth.token_duration", 24*time.Hour)
	v.SetDefault("auth.clerk_jwt_template", "api_auth")

	// Notification defaults
	v.SetDefault("notifications.email.enabled", false)
	v.SetDefault("notifications.email.smtp_port", 587)
	v.SetDefault("notifications.email.min_level", "error")
	v.SetDefault("notifications.email.subject_prefix", "[CryptoBot]")

	v.SetDefault("notifications.webhook.enabled", false)
	v.SetDefault("notifications.webhook.method", "POST")
	v.SetDefault("notifications.webhook.min_level", "warning")
	v.SetDefault("notifications.webhook.timeout", 10*time.Second)
	v.SetDefault("notifications.webhook.batch_size", 1)

	// Database defaults
	v.SetDefault("database.driver", "sqlite")
	v.SetDefault("database.path", "./data/crypto_bot.db")
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 5432)
	v.SetDefault("database.name", "crypto_bot")
	v.SetDefault("database.ssl_mode", "disable")
	v.SetDefault("database.turso.enabled", false)

	// Market defaults
	v.SetDefault("market.cache.ticker_ttl", 300)   // 5 minutes in seconds
	v.SetDefault("market.cache.candle_ttl", 900)   // 15 minutes in seconds
	v.SetDefault("market.cache.orderbook_ttl", 30) // 30 seconds

	// MEXC defaults
	v.SetDefault("mexc.base_url", "https://api.mexc.com")
	v.SetDefault("mexc.ws_base_url", "wss://wbs.mexc.com/ws")
	v.SetDefault("mexc.use_testnet", false)
	v.SetDefault("mexc.rate_limit.requests_per_minute", 1200)
	v.SetDefault("mexc.rate_limit.burst_size", 10)

	// Rate limiting defaults
	defaultRateLimit := GetDefaultRateLimitConfig()
	v.SetDefault("rate_limit.enabled", defaultRateLimit.Enabled)
	v.SetDefault("rate_limit.default_limit", defaultRateLimit.DefaultLimit)
	v.SetDefault("rate_limit.default_burst", defaultRateLimit.DefaultBurst)
	v.SetDefault("rate_limit.ip_limit", defaultRateLimit.IPLimit)
	v.SetDefault("rate_limit.ip_burst", defaultRateLimit.IPBurst)
	v.SetDefault("rate_limit.user_limit", defaultRateLimit.UserLimit)
	v.SetDefault("rate_limit.user_burst", defaultRateLimit.UserBurst)
	v.SetDefault("rate_limit.auth_user_limit", defaultRateLimit.AuthUserLimit)
	v.SetDefault("rate_limit.auth_user_burst", defaultRateLimit.AuthUserBurst)
	v.SetDefault("rate_limit.cleanup_interval", defaultRateLimit.CleanupInterval)
	v.SetDefault("rate_limit.block_duration", defaultRateLimit.BlockDuration)
	v.SetDefault("rate_limit.trusted_proxies", defaultRateLimit.TrustedProxies)
	v.SetDefault("rate_limit.excluded_paths", defaultRateLimit.ExcludedPaths)
	v.SetDefault("rate_limit.redis_enabled", defaultRateLimit.RedisEnabled)
	v.SetDefault("rate_limit.redis_key_prefix", defaultRateLimit.RedisKeyPrefix)

	// CSRF defaults
	defaultCSRF := GetDefaultCSRFConfig()
	v.SetDefault("csrf.enabled", defaultCSRF.Enabled)
	v.SetDefault("csrf.token_length", defaultCSRF.TokenLength)
	v.SetDefault("csrf.cookie_name", defaultCSRF.CookieName)
	v.SetDefault("csrf.cookie_path", defaultCSRF.CookiePath)
	v.SetDefault("csrf.cookie_max_age", defaultCSRF.CookieMaxAge)
	v.SetDefault("csrf.cookie_secure", defaultCSRF.CookieSecure)
	v.SetDefault("csrf.cookie_http_only", defaultCSRF.CookieHTTPOnly)
	v.SetDefault("csrf.cookie_same_site", defaultCSRF.CookieSameSite)
	v.SetDefault("csrf.header_name", defaultCSRF.HeaderName)
	v.SetDefault("csrf.form_field_name", defaultCSRF.FormFieldName)
	v.SetDefault("csrf.excluded_paths", defaultCSRF.ExcludedPaths)
	v.SetDefault("csrf.excluded_methods", defaultCSRF.ExcludedMethods)
	v.SetDefault("csrf.failure_status_code", defaultCSRF.FailureStatusCode)

	// Secure headers defaults
	defaultSecureHeaders := GetDefaultSecureHeadersConfig()
	v.SetDefault("secure_headers.enabled", defaultSecureHeaders.Enabled)
	v.SetDefault("secure_headers.content_security_policy", defaultSecureHeaders.ContentSecurityPolicy)
	v.SetDefault("secure_headers.x_content_type_options", defaultSecureHeaders.XContentTypeOptions)
	v.SetDefault("secure_headers.x_frame_options", defaultSecureHeaders.XFrameOptions)
	v.SetDefault("secure_headers.x_xss_protection", defaultSecureHeaders.XXSSProtection)
	v.SetDefault("secure_headers.referrer_policy", defaultSecureHeaders.ReferrerPolicy)
	v.SetDefault("secure_headers.strict_transport_security", defaultSecureHeaders.StrictTransportSecurity)
	v.SetDefault("secure_headers.permissions_policy", defaultSecureHeaders.PermissionsPolicy)
	v.SetDefault("secure_headers.cross_origin_embedder_policy", defaultSecureHeaders.CrossOriginEmbedderPolicy)
	v.SetDefault("secure_headers.cross_origin_opener_policy", defaultSecureHeaders.CrossOriginOpenerPolicy)
	v.SetDefault("secure_headers.cross_origin_resource_policy", defaultSecureHeaders.CrossOriginResourcePolicy)
	v.SetDefault("secure_headers.cache_control", defaultSecureHeaders.CacheControl)
	v.SetDefault("secure_headers.excluded_paths", defaultSecureHeaders.ExcludedPaths)
	v.SetDefault("secure_headers.custom_headers", defaultSecureHeaders.CustomHeaders)
	v.SetDefault("secure_headers.remove_server_header", defaultSecureHeaders.RemoveServerHeader)
	v.SetDefault("secure_headers.remove_powered_by_header", defaultSecureHeaders.RemovePoweredByHeader)
	v.SetDefault("secure_headers.content_security_policy_report_only", defaultSecureHeaders.ContentSecurityPolicyReportOnly)
	v.SetDefault("secure_headers.content_security_policy_report_uri", defaultSecureHeaders.ContentSecurityPolicyReportURI)

	// AI defaults
	v.SetDefault("ai.provider", "gemini")
	v.SetDefault("ai.model", "gemini-pro")
	v.SetDefault("ai.gemini_model", "gemini-1.5-flash")
	v.SetDefault("ai.system_prompt", "You are a crypto trading assistant. You help users understand their portfolio, market trends, and provide trading advice. Keep responses concise and focused on crypto trading.")
	v.SetDefault("ai.temperature", 0.7)
	v.SetDefault("ai.top_p", 0.95)
	v.SetDefault("ai.top_k", 40)
	v.SetDefault("ai.max_tokens", 1024)

	// Web3 defaults
	v.SetDefault("infura_api_key", "")
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
