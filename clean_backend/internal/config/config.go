package config

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
)

// Config holds the root application configuration struct.
type Config struct {
	Server        ServerConfig        `mapstructure:"server"`
	Database      DatabaseConfig      `mapstructure:"database"`
	Log           LogConfig           `mapstructure:"log"`
	MEXC          MEXCConfig          `mapstructure:"mexc"`
	Encryption    EncryptionConfig    `mapstructure:"encryption"`
	Auth          AuthConfig          `mapstructure:"auth"`
	AI            AIConfig            `mapstructure:"ai"`
	RateLimit     RateLimitConfig     `mapstructure:"rate_limit"`
	CSRF          CSRFConfig          `mapstructure:"csrf"`
	SecureHeaders SecureHeadersConfig `mapstructure:"secure_headers"`

	// Environment indicator (development, testing, staging, production)
	Environment string `mapstructure:"environment" env:"APP_ENV" default:"development"`

	// Mock configuration for controlling mock implementations
	Mock MockConfig `mapstructure:"mock"`
}

// ServerConfig holds HTTP server specific configuration.
type ServerConfig struct {
	Port               int           `mapstructure:"port"`                 // Port the HTTP server listens on.
	Host               string        `mapstructure:"host"`                 // Host address the server binds to.
	ReadTimeout        time.Duration `mapstructure:"read_timeout"`         // Max duration for reading the entire request.
	WriteTimeout       time.Duration `mapstructure:"write_timeout"`        // Max duration before timing out writes of the response.
	IdleTimeout        time.Duration `mapstructure:"idle_timeout"`         // Max amount of time to wait for the next request when keep-alives are enabled.
	FrontendURL        string        `mapstructure:"frontend_url"`         // Base URL of the frontend application (for CORS, redirects etc.).
	CORSAllowedOrigins []string      `mapstructure:"cors_allowed_origins"` // List of allowed origins for CORS.
}

// DatabaseConfig holds database connection details.
type DatabaseConfig struct {
	Type                   string `mapstructure:"type"`                      // Type of database ("sqlite" or "turso").
	DSN                    string `mapstructure:"dsn"`                       // Data Source Name (connection string or file path for SQLite).
	MaxIdleConns           int    `mapstructure:"max_idle_conns"`            // Maximum number of connections in the idle connection pool.
	MaxOpenConns           int    `mapstructure:"max_open_conns"`            // Maximum number of open connections to the database.
	ConnMaxLifetimeMinutes int    `mapstructure:"conn_max_lifetime_minutes"` // Maximum amount of time a connection may be reused (in minutes).
	TursoURL               string `mapstructure:"turso_url"`                 // URL for Turso database (if type is "turso").
	AuthToken              string `mapstructure:"auth_token"`                // Auth token for Turso database (if type is "turso").
	EnableLogging          bool   `mapstructure:"enable_logging"`            // Whether to enable GORM query logging.
	AutoMigrate            bool   `mapstructure:"auto_migrate"`              // Whether to automatically run database migrations on startup.
}

// LogConfig holds logging level and format settings.
type LogConfig struct {
	Level  string `mapstructure:"level"`  // Log level (e.g., "debug", "info", "warn", "error").
	Pretty bool   `mapstructure:"pretty"` // Whether to use human-readable console output.
}

// MEXCConfig holds API credentials for the MEXC exchange.
type MEXCConfig struct {
	APIKey    string `mapstructure:"api_key" env:"MEXC_API_KEY"`       // MEXC API Key (loaded from env).
	SecretKey string `mapstructure:"secret_key" env:"MEXC_SECRET_KEY"` // MEXC Secret Key (loaded from env).
	// BaseURL   string `mapstructure:"base_url" env:"MEXC_BASE_URL" default:"https://api.mexc.com"` // Optional: Base URL for MEXC API
}

// EncryptionConfig holds settings for data encryption.
type EncryptionConfig struct {
	CurrentKeyID string `mapstructure:"current_key_id" env:"ENCRYPTION_CURRENT_KEY_ID"` // ID of the key currently used for encryption.
	Keys         string `mapstructure:"keys" env:"ENCRYPTION_KEYS"`                     // Comma-separated list of key pairs (e.g., "key1:base64encodedkey1,key2:base64...").
}

// MockConfig holds configuration for mock implementations.
type MockConfig struct {
	Enabled       bool `mapstructure:"enabled" env:"MOCK_ENABLED" default:"false"`               // Whether to enable mock implementations globally.
	MEXCClient    bool `mapstructure:"mexc_client" env:"MOCK_MEXC_CLIENT" default:"false"`       // Whether to use mock MEXC client.
	WalletService bool `mapstructure:"wallet_service" env:"MOCK_WALLET_SERVICE" default:"false"` // Whether to use mock wallet service.
	AuthService   bool `mapstructure:"auth_service" env:"MOCK_AUTH_SERVICE" default:"false"`     // Whether to use mock auth service.
	AIService     bool `mapstructure:"ai_service" env:"MOCK_AI_SERVICE" default:"false"`         // Whether to use mock AI service.
}

// AuthConfig holds authentication provider and settings.
type AuthConfig struct {
	Provider          string `mapstructure:"provider" env:"AUTH_PROVIDER" default:"clerk"`                   // Authentication provider ("clerk", "jwt", etc.).
	Disabled          bool   `mapstructure:"disabled" env:"AUTH_DISABLED" default:"false"`                   // If true, authentication checks are bypassed.
	ClerkSecretKey    string `mapstructure:"clerk_secret_key" env:"CLERK_SECRET_KEY"`                        // Secret key for Clerk backend API.
	ClerkJWTPublicKey string `mapstructure:"clerk_jwt_public_key" env:"CLERK_JWT_PUBLIC_KEY"`                // Public key for verifying Clerk JWTs.
	ClerkJWTTemplate  string `mapstructure:"clerk_jwt_template" env:"CLERK_JWT_TEMPLATE" default:"api_auth"` // Clerk session JWT template name (if applicable).
	JWTSecret         string `mapstructure:"jwt_secret" env:"JWT_SECRET" default:"test-secret-key"`          // Secret key for signing/verifying JWTs (for non-Clerk testing).
	UseEnhanced       bool   `mapstructure:"use_enhanced" env:"AUTH_USE_ENHANCED" default:"false"`           // Flag to enable enhanced authentication features (implementation specific).
}

// AIConfig holds configuration for AI services.
type AIConfig struct {
	Provider              string  `mapstructure:"provider" env:"AI_PROVIDER" default:"google_vertexai"`         // AI provider ("google_vertexai", "openai", "mock", etc.).
	APIKey                string  `mapstructure:"api_key" env:"AI_API_KEY"`                                     // API key for the AI service (optional, depends on provider/auth method).
	GoogleCredentialsPath string  `mapstructure:"google_credentials_path" env:"GOOGLE_APPLICATION_CREDENTIALS"` // Path to Google service account JSON file.
	ProjectID             string  `mapstructure:"project_id" env:"AI_PROJECT_ID"`                               // Google Cloud Project ID or equivalent.
	LocationID            string  `mapstructure:"location_id" env:"AI_LOCATION_ID" default:"us-central1"`       // Google Cloud Location ID or equivalent.
	EndpointID            string  `mapstructure:"endpoint_id" env:"AI_ENDPOINT_ID"`                             // Specific endpoint ID for Google Vertex AI.
	ModelName             string  `mapstructure:"model_name" env:"AI_MODEL_NAME" default:"gemini-pro"`          // Model name (e.g., "gemini-pro", "gpt-4").
	MaxTokens             int     `mapstructure:"max_tokens" env:"AI_MAX_TOKENS" default:"1024"`                // Maximum tokens to generate.
	Temperature           float64 `mapstructure:"temperature" env:"AI_TEMPERATURE" default:"0.7"`               // Temperature for generation (0.0-1.0).
}

// RateLimitConfig holds settings for rate limiting middleware.
type RateLimitConfig struct {
	Enabled         bool                     `mapstructure:"enabled"`
	DefaultLimit    float64                  `mapstructure:"default_limit"`    // Requests per second allowed by default.
	DefaultBurst    int                      `mapstructure:"default_burst"`    // Burst allowed by default.
	IPLimit         float64                  `mapstructure:"ip_limit"`         // Limit per IP.
	IPBurst         int                      `mapstructure:"ip_burst"`         // Burst per IP.
	UserLimit       float64                  `mapstructure:"user_limit"`       // Limit per user ID (if available).
	UserBurst       int                      `mapstructure:"user_burst"`       // Burst per user ID.
	AuthUserLimit   float64                  `mapstructure:"auth_user_limit"`  // Specific limit for authenticated users.
	AuthUserBurst   int                      `mapstructure:"auth_user_burst"`  // Burst for authenticated users.
	CleanupInterval time.Duration            `mapstructure:"cleanup_interval"` // How often to clean up old visitor data.
	BlockDuration   time.Duration            `mapstructure:"block_duration"`   // How long to block after exceeding limit.
	TrustedProxies  []string                 `mapstructure:"trusted_proxies"`  // List of trusted proxy IPs/CIDRs.
	ExcludedPaths   []string                 `mapstructure:"excluded_paths"`   // Paths to exclude from rate limiting.
	RedisEnabled    bool                     `mapstructure:"redis_enabled"`    // Use Redis for distributed rate limiting.
	RedisURL        string                   `mapstructure:"redis_url"`        // Redis URL.
	RedisKeyPrefix  string                   `mapstructure:"redis_key_prefix"` // Prefix for Redis keys.
	EndpointLimits  map[string]EndpointLimit `mapstructure:"endpoint_limits"`  // Endpoint-specific limits.
}

// EndpointLimit contains rate limiting configuration for a specific endpoint
type EndpointLimit struct {
	Path      string  `mapstructure:"path"`       // Path pattern to match
	Method    string  `mapstructure:"method"`     // HTTP method to match
	Limit     float64 `mapstructure:"limit"`      // Requests per second
	Burst     int     `mapstructure:"burst"`      // Burst size
	UserLimit float64 `mapstructure:"user_limit"` // Requests per second per user
	UserBurst int     `mapstructure:"user_burst"` // Burst size per user
}

// CSRFConfig holds settings for CSRF protection middleware.
type CSRFConfig struct {
	Enabled           bool          `mapstructure:"enabled"`
	Secret            string        `mapstructure:"secret" env:"CSRF_SECRET"` // Secret key for signing CSRF tokens (REQUIRED if enabled).
	TokenLength       int           `mapstructure:"token_length"`             // Length of the generated CSRF token.
	CookieName        string        `mapstructure:"cookie_name"`              // Name of the CSRF cookie.
	CookiePath        string        `mapstructure:"cookie_path"`              // Path for the CSRF cookie.
	CookieDomain      string        `mapstructure:"cookie_domain"`            // Domain for the CSRF cookie.
	CookieMaxAge      time.Duration `mapstructure:"cookie_max_age"`           // Max age of the cookie (changed type to time.Duration).
	CookieSecure      bool          `mapstructure:"cookie_secure"`            // Sets the Secure flag on the cookie (requires HTTPS).
	CookieHTTPOnly    bool          `mapstructure:"cookie_http_only"`         // Sets the HttpOnly flag.
	CookieSameSite    string        `mapstructure:"cookie_same_site"`         // Sets SameSite attribute (e.g., "Strict", "Lax", "None").
	HeaderName        string        `mapstructure:"header_name"`              // HTTP header name where token is expected.
	FormFieldName     string        `mapstructure:"form_field_name"`          // Form field name where token is expected.
	ExcludedPaths     []string      `mapstructure:"excluded_paths"`           // Paths excluded from CSRF protection.
	ExcludedMethods   []string      `mapstructure:"excluded_methods"`         // HTTP methods excluded from CSRF protection.
	FailureStatusCode int           `mapstructure:"failure_status_code"`      // HTTP status code to return on failure.
}

// SecureHeadersConfig holds settings for security-related HTTP headers.
type SecureHeadersConfig struct {
	Enabled                         bool              `mapstructure:"enabled"`
	ContentSecurityPolicy           string            `mapstructure:"content_security_policy"`             // Content-Security-Policy header value.
	ContentSecurityPolicyReportOnly bool              `mapstructure:"content_security_policy_report_only"` // Enable CSP report-only mode.
	ContentSecurityPolicyReportURI  string            `mapstructure:"content_security_policy_report_uri"`  // URI for CSP violation reports.
	XContentTypeOptions             string            `mapstructure:"x_content_type_options"`              // X-Content-Type-Options header value.
	XFrameOptions                   string            `mapstructure:"x_frame_options"`                     // X-Frame-Options header value.
	XXSSProtection                  string            `mapstructure:"x_xss_protection"`                    // X-XSS-Protection header value.
	ReferrerPolicy                  string            `mapstructure:"referrer_policy"`                     // Referrer-Policy header value.
	StrictTransportSecurity         string            `mapstructure:"strict_transport_security"`           // Strict-Transport-Security header value.
	PermissionsPolicy               string            `mapstructure:"permissions_policy"`                  // Permissions-Policy header value.
	CrossOriginEmbedderPolicy       string            `mapstructure:"cross_origin_embedder_policy"`        // Cross-Origin-Embedder-Policy header value.
	CrossOriginOpenerPolicy         string            `mapstructure:"cross_origin_opener_policy"`          // Cross-Origin-Opener-Policy header value.
	CrossOriginResourcePolicy       string            `mapstructure:"cross_origin_resource_policy"`        // Cross-Origin-Resource-Policy header value.
	CacheControl                    string            `mapstructure:"cache_control"`                       // Cache-Control header value.
	ExcludedPaths                   []string          `mapstructure:"excluded_paths"`                      // Paths excluded from secure headers.
	CustomHeaders                   map[string]string `mapstructure:"custom_headers"`                      // Custom headers to add.
	RemoveServerHeader              bool              `mapstructure:"remove_server_header"`                // Remove the default 'Server' header.
	RemovePoweredByHeader           bool              `mapstructure:"remove_powered_by_header"`            // Remove the default 'X-Powered-By' header.
}

// LoadConfig loads configuration from file/env.
func LoadConfig(path string) (*Config, error) {
	v := viper.New()

	// Set default values
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.read_timeout", 10*time.Second)
	v.SetDefault("server.write_timeout", 10*time.Second)
	v.SetDefault("server.idle_timeout", 60*time.Second)
	v.SetDefault("server.cors_allowed_origins", []string{"*"})

	v.SetDefault("database.type", "sqlite")
	v.SetDefault("database.dsn", "./data/crypto_bot.db")
	v.SetDefault("database.max_idle_conns", 5)
	v.SetDefault("database.max_open_conns", 10)
	v.SetDefault("database.conn_max_lifetime_minutes", 60)
	v.SetDefault("database.enable_logging", false)
	v.SetDefault("database.auto_migrate", true)

	v.SetDefault("log.level", "info")
	v.SetDefault("log.pretty", true)

	// Rate limit defaults
	rateLimitConfig := GetDefaultRateLimitConfig()
	v.SetDefault("rate_limit.enabled", rateLimitConfig.Enabled)
	v.SetDefault("rate_limit.default_limit", rateLimitConfig.DefaultLimit)
	v.SetDefault("rate_limit.default_burst", rateLimitConfig.DefaultBurst)
	v.SetDefault("rate_limit.ip_limit", rateLimitConfig.IPLimit)
	v.SetDefault("rate_limit.ip_burst", rateLimitConfig.IPBurst)
	v.SetDefault("rate_limit.user_limit", rateLimitConfig.UserLimit)
	v.SetDefault("rate_limit.user_burst", rateLimitConfig.UserBurst)
	v.SetDefault("rate_limit.auth_user_limit", rateLimitConfig.AuthUserLimit)
	v.SetDefault("rate_limit.auth_user_burst", rateLimitConfig.AuthUserBurst)
	v.SetDefault("rate_limit.cleanup_interval", rateLimitConfig.CleanupInterval)
	v.SetDefault("rate_limit.block_duration", rateLimitConfig.BlockDuration)
	v.SetDefault("rate_limit.trusted_proxies", rateLimitConfig.TrustedProxies)
	v.SetDefault("rate_limit.excluded_paths", rateLimitConfig.ExcludedPaths)
	v.SetDefault("rate_limit.redis_enabled", rateLimitConfig.RedisEnabled)
	v.SetDefault("rate_limit.redis_url", rateLimitConfig.RedisURL)
	v.SetDefault("rate_limit.redis_key_prefix", rateLimitConfig.RedisKeyPrefix)
	v.SetDefault("rate_limit.endpoint_limits", rateLimitConfig.EndpointLimits)

	// CSRF defaults
	csrfConfig := GetDefaultCSRFConfig()
	v.SetDefault("csrf.enabled", csrfConfig.Enabled)
	v.SetDefault("csrf.secret", csrfConfig.Secret)
	v.SetDefault("csrf.token_length", csrfConfig.TokenLength)
	v.SetDefault("csrf.cookie_name", csrfConfig.CookieName)
	v.SetDefault("csrf.cookie_path", csrfConfig.CookiePath)
	v.SetDefault("csrf.cookie_max_age", csrfConfig.CookieMaxAge)
	v.SetDefault("csrf.cookie_secure", csrfConfig.CookieSecure)
	v.SetDefault("csrf.cookie_http_only", csrfConfig.CookieHTTPOnly)
	v.SetDefault("csrf.cookie_same_site", csrfConfig.CookieSameSite)
	v.SetDefault("csrf.header_name", csrfConfig.HeaderName)
	v.SetDefault("csrf.form_field_name", csrfConfig.FormFieldName)
	v.SetDefault("csrf.excluded_paths", csrfConfig.ExcludedPaths)
	v.SetDefault("csrf.excluded_methods", csrfConfig.ExcludedMethods)
	v.SetDefault("csrf.failure_status_code", csrfConfig.FailureStatusCode)

	// Secure headers defaults
	secureHeadersConfig := GetDefaultSecureHeadersConfig()
	v.SetDefault("secure_headers.enabled", secureHeadersConfig.Enabled)
	v.SetDefault("secure_headers.content_security_policy", secureHeadersConfig.ContentSecurityPolicy)
	v.SetDefault("secure_headers.content_security_policy_report_only", secureHeadersConfig.ContentSecurityPolicyReportOnly)
	v.SetDefault("secure_headers.content_security_policy_report_uri", secureHeadersConfig.ContentSecurityPolicyReportURI)
	v.SetDefault("secure_headers.x_content_type_options", secureHeadersConfig.XContentTypeOptions)
	v.SetDefault("secure_headers.x_frame_options", secureHeadersConfig.XFrameOptions)
	v.SetDefault("secure_headers.x_xss_protection", secureHeadersConfig.XXSSProtection)
	v.SetDefault("secure_headers.referrer_policy", secureHeadersConfig.ReferrerPolicy)
	v.SetDefault("secure_headers.strict_transport_security", secureHeadersConfig.StrictTransportSecurity)
	v.SetDefault("secure_headers.permissions_policy", secureHeadersConfig.PermissionsPolicy)
	v.SetDefault("secure_headers.cross_origin_embedder_policy", secureHeadersConfig.CrossOriginEmbedderPolicy)
	v.SetDefault("secure_headers.cross_origin_opener_policy", secureHeadersConfig.CrossOriginOpenerPolicy)
	v.SetDefault("secure_headers.cross_origin_resource_policy", secureHeadersConfig.CrossOriginResourcePolicy)
	v.SetDefault("secure_headers.cache_control", secureHeadersConfig.CacheControl)
	v.SetDefault("secure_headers.excluded_paths", secureHeadersConfig.ExcludedPaths)
	v.SetDefault("secure_headers.custom_headers", secureHeadersConfig.CustomHeaders)
	v.SetDefault("secure_headers.remove_server_header", secureHeadersConfig.RemoveServerHeader)
	v.SetDefault("secure_headers.remove_powered_by_header", secureHeadersConfig.RemovePoweredByHeader)

	// Encryption defaults
	encryptionConfig := GetDefaultEncryptionConfig()
	v.SetDefault("encryption.current_key_id", encryptionConfig.CurrentKeyID)
	v.SetDefault("encryption.keys", encryptionConfig.Keys)

	// Mock defaults
	mockConfig := GetDefaultMockConfig()
	v.SetDefault("mock.enabled", mockConfig.Enabled)
	v.SetDefault("mock.mexc_client", mockConfig.MEXCClient)
	v.SetDefault("mock.wallet_service", mockConfig.WalletService)
	v.SetDefault("mock.auth_service", mockConfig.AuthService)
	v.SetDefault("mock.ai_service", mockConfig.AIService)

	// AI defaults
	aiConfig := GetDefaultAIConfig()
	v.SetDefault("ai.provider", aiConfig.Provider)
	v.SetDefault("ai.api_key", aiConfig.APIKey)
	v.SetDefault("ai.model", aiConfig.ModelName)
	v.SetDefault("ai.max_tokens", aiConfig.MaxTokens)
	v.SetDefault("ai.temperature", aiConfig.Temperature)

	// Environment defaults
	v.SetDefault("environment", "development")

	// Environment variables
	v.AutomaticEnv()

	// Map environment variables
	v.BindEnv("server.port", "SERVER_PORT")
	v.BindEnv("server.host", "SERVER_HOST")
	v.BindEnv("database.type", "DATABASE_TYPE")
	v.BindEnv("database.dsn", "DATABASE_DSN")
	v.BindEnv("database.max_idle_conns", "DATABASE_MAX_IDLE_CONNS")
	v.BindEnv("database.max_open_conns", "DATABASE_MAX_OPEN_CONNS")
	v.BindEnv("database.conn_max_lifetime_minutes", "DATABASE_CONN_MAX_LIFETIME_MINUTES")
	v.BindEnv("database.enable_logging", "DATABASE_ENABLE_LOGGING")
	v.BindEnv("database.auto_migrate", "DATABASE_AUTO_MIGRATE")
	v.BindEnv("database.turso_url", "TURSO_DB_URL")
	v.BindEnv("database.auth_token", "TURSO_AUTH_TOKEN")
	v.BindEnv("log.level", "LOG_LEVEL")
	v.BindEnv("log.pretty", "LOG_PRETTY")
	v.BindEnv("mexc.api_key", "MEXC_API_KEY")
	v.BindEnv("mexc.secret_key", "MEXC_SECRET_KEY")
	v.BindEnv("encryption.current_key_id", "CURRENT_KEY_ID")
	v.BindEnv("encryption.keys", "ENCRYPTION_KEYS")
	v.BindEnv("auth.provider", "AUTH_PROVIDER")
	v.BindEnv("auth.disabled", "DISABLE_AUTH")
	v.BindEnv("auth.clerk_secret_key", "CLERK_SECRET_KEY")
	v.BindEnv("auth.clerk_jwt_public_key", "CLERK_JWT_PUBLIC_KEY")
	v.BindEnv("auth.clerk_jwt_template", "CLERK_JWT_TEMPLATE")
	v.BindEnv("auth.jwt_secret", "JWT_SECRET")
	v.BindEnv("auth.use_enhanced", "USE_ENHANCED_AUTH")

	// Environment and mock settings
	v.BindEnv("environment", "APP_ENV")
	v.BindEnv("mock.enabled", "MOCK_ENABLED")
	v.BindEnv("mock.mexc_client", "MOCK_MEXC_CLIENT")
	v.BindEnv("mock.wallet_service", "MOCK_WALLET_SERVICE")
	v.BindEnv("mock.auth_service", "MOCK_AUTH_SERVICE")
	v.BindEnv("mock.ai_service", "MOCK_AI_SERVICE")

	// AI settings
	v.BindEnv("ai.provider", "AI_PROVIDER")
	v.BindEnv("ai.api_key", "AI_API_KEY")
	v.BindEnv("ai.model", "AI_MODEL")
	v.BindEnv("ai.max_tokens", "AI_MAX_TOKENS")
	v.BindEnv("ai.temperature", "AI_TEMPERATURE")

	// Try to read config file if it exists
	v.AddConfigPath(path)
	v.SetConfigName("config")
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		// It's okay if config file doesn't exist
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	// Unmarshal config
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	// Validate the loaded configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	// Override with environment variables
	if port := os.Getenv("SERVER_PORT"); port != "" {
		fmt.Sscanf(port, "%d", &cfg.Server.Port)
	}

	return &cfg, nil
}

// Validate checks for required configuration values.
func (c *Config) Validate() error {
	// Example Validations:
	if c.Database.Type == "turso" && (c.Database.TursoURL == "" || c.Database.AuthToken == "") {
		return fmt.Errorf("database type is turso, but TURSO_DB_URL or TURSO_AUTH_TOKEN is missing")
	}
	if c.Database.Type == "sqlite" && c.Database.DSN == "" {
		return fmt.Errorf("database type is sqlite, but DATABASE_DSN is missing")
	}

	if c.Auth.Provider == "clerk" && c.Auth.ClerkSecretKey == "" {
		return fmt.Errorf("auth provider is clerk, but CLERK_SECRET_KEY is missing")
	}

	if c.Encryption.CurrentKeyID == "" || c.Encryption.Keys == "" {
		return fmt.Errorf("encryption keys (CURRENT_KEY_ID, ENCRYPTION_KEYS) are missing")
	}

	// Add more validation checks as needed for other critical fields

	// Prevent mock implementations in production
	if c.Environment == "production" && c.Mock.Enabled {
		return fmt.Errorf("mock implementations are not allowed in production environment")
	}

	return nil
}

// GetDefaultRateLimitConfig returns the default rate limit configuration.
func GetDefaultRateLimitConfig() RateLimitConfig {
	return RateLimitConfig{
		Enabled:         true,
		DefaultLimit:    60.0, // Changed to float64 to match struct
		DefaultBurst:    10,
		IPLimit:         300.0, // Changed to float64
		IPBurst:         20,
		UserLimit:       600.0, // Changed to float64
		UserBurst:       30,
		AuthUserLimit:   1200.0, // Changed to float64
		AuthUserBurst:   60,
		CleanupInterval: 5 * time.Minute,
		BlockDuration:   15 * time.Minute,
		TrustedProxies:  []string{"127.0.0.1", "::1"}, // Use defaults from original file
		ExcludedPaths:   []string{"/health", "/metrics", "/favicon.ico"},
		RedisEnabled:    false,
		RedisKeyPrefix:  "ratelimit:",
		EndpointLimits:  map[string]EndpointLimit{}, // Default to empty map
	}
}

// GetDefaultCSRFConfig returns the default CSRF configuration.
func GetDefaultCSRFConfig() CSRFConfig {
	return CSRFConfig{
		Enabled: true,
		// Secret:         "", // No default secret, MUST be set via env/config
		TokenLength: 32,
		CookieName:  "csrf_token",
		CookiePath:  "/",
		// CookieDomain:   "", // No default domain
		CookieMaxAge:      24 * time.Hour, // Use time.Duration
		CookieSecure:      true,
		CookieHTTPOnly:    true,
		CookieSameSite:    "Lax", // Changed from Strict to Lax as per original
		HeaderName:        "X-CSRF-Token",
		FormFieldName:     "csrf_token", // Changed from _csrf as per original
		ExcludedPaths:     []string{"/health", "/metrics", "/favicon.ico"},
		ExcludedMethods:   []string{"GET", "HEAD", "OPTIONS", "TRACE"},
		FailureStatusCode: 403,
	}
}

// GetDefaultSecureHeadersConfig returns the default secure headers configuration.
func GetDefaultSecureHeadersConfig() SecureHeadersConfig {
	return SecureHeadersConfig{
		Enabled:                         true,
		ContentSecurityPolicy:           "default-src 'self'; script-src 'self'; object-src 'none'; img-src 'self' data:; style-src 'self' 'unsafe-inline'; font-src 'self'; frame-src 'none'; connect-src 'self'", // From original
		ContentSecurityPolicyReportOnly: false,
		ContentSecurityPolicyReportURI:  "",
		XContentTypeOptions:             "nosniff",
		XFrameOptions:                   "DENY",
		XXSSProtection:                  "1; mode=block",
		ReferrerPolicy:                  "strict-origin-when-cross-origin",
		StrictTransportSecurity:         "max-age=31536000; includeSubDomains",
		PermissionsPolicy:               "camera=(), microphone=(), geolocation=(), interest-cohort=()", // From original
		CrossOriginEmbedderPolicy:       "require-corp",
		CrossOriginOpenerPolicy:         "same-origin",
		CrossOriginResourcePolicy:       "same-origin",
		CacheControl:                    "no-store, max-age=0", // From original
		ExcludedPaths:                   []string{"/health", "/metrics", "/favicon.ico"},
		CustomHeaders:                   map[string]string{},
		RemoveServerHeader:              true,
		RemovePoweredByHeader:           true,
	}
}

// GetDefaultEncryptionConfig returns the default encryption configuration.
func GetDefaultEncryptionConfig() EncryptionConfig {
	return EncryptionConfig{
		CurrentKeyID: "default",
		Keys:         "default:c29tZVJhbmRvbUJhc2U2NEtleQ==", // Placeholder key
	} // Added missing closing brace
}

// GetDefaultMockConfig returns the default mock configuration.
func GetDefaultMockConfig() MockConfig {
	return MockConfig{
		Enabled:       false,
		MEXCClient:    false,
		WalletService: false,
		AuthService:   false,
		AIService:     false,
	}
}

// GetDefaultAIConfig returns the default AI configuration.
func GetDefaultAIConfig() AIConfig {
	return AIConfig{
		Provider:    "google_vertexai",
		LocationID:  "us-central1",
		ModelName:   "gemini-pro",
		MaxTokens:   1024,
		Temperature: 0.7,
	}
}
