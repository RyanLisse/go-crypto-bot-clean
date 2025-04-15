package config

import (
	"time"
)

// RateLimitConfig contains rate limiting configuration
type RateLimitConfig struct {
	Enabled            bool          `mapstructure:"enabled"`
	DefaultLimit       int           `mapstructure:"default_limit"`       // Requests per minute
	DefaultBurst       int           `mapstructure:"default_burst"`       // Burst size
	IPLimit            int           `mapstructure:"ip_limit"`            // Requests per minute per IP
	IPBurst            int           `mapstructure:"ip_burst"`            // Burst size per IP
	UserLimit          int           `mapstructure:"user_limit"`          // Requests per minute per user
	UserBurst          int           `mapstructure:"user_burst"`          // Burst size per user
	AuthUserLimit      int           `mapstructure:"auth_user_limit"`     // Requests per minute for authenticated users
	AuthUserBurst      int           `mapstructure:"auth_user_burst"`     // Burst size for authenticated users
	CleanupInterval    time.Duration `mapstructure:"cleanup_interval"`    // Interval to clean up expired limiters
	BlockDuration      time.Duration `mapstructure:"block_duration"`      // Duration to block after exceeding limit
	TrustedProxies     []string      `mapstructure:"trusted_proxies"`     // List of trusted proxies
	ExcludedPaths      []string      `mapstructure:"excluded_paths"`      // Paths to exclude from rate limiting
	RedisEnabled       bool          `mapstructure:"redis_enabled"`       // Use Redis for distributed rate limiting
	RedisURL           string        `mapstructure:"redis_url"`           // Redis URL
	RedisKeyPrefix     string        `mapstructure:"redis_key_prefix"`    // Prefix for Redis keys
	EndpointLimits     map[string]EndpointLimit `mapstructure:"endpoint_limits"` // Endpoint-specific limits
}

// EndpointLimit contains rate limiting configuration for a specific endpoint
type EndpointLimit struct {
	Path      string `mapstructure:"path"`      // Path pattern to match
	Method    string `mapstructure:"method"`    // HTTP method to match
	Limit     int    `mapstructure:"limit"`     // Requests per minute
	Burst     int    `mapstructure:"burst"`     // Burst size
	UserLimit int    `mapstructure:"user_limit"` // Requests per minute per user
	UserBurst int    `mapstructure:"user_burst"` // Burst size per user
}

// GetDefaultRateLimitConfig returns the default rate limiting configuration
func GetDefaultRateLimitConfig() RateLimitConfig {
	return RateLimitConfig{
		Enabled:         true,
		DefaultLimit:    60,  // 1 request per second
		DefaultBurst:    10,  // Allow bursts of 10 requests
		IPLimit:         300, // 5 requests per second per IP
		IPBurst:         20,  // Allow bursts of 20 requests per IP
		UserLimit:       600, // 10 requests per second per user
		UserBurst:       30,  // Allow bursts of 30 requests per user
		AuthUserLimit:   1200, // 20 requests per second for authenticated users
		AuthUserBurst:   60,   // Allow bursts of 60 requests for authenticated users
		CleanupInterval: 5 * time.Minute,
		BlockDuration:   15 * time.Minute,
		TrustedProxies:  []string{"127.0.0.1", "::1"},
		ExcludedPaths:   []string{"/health", "/metrics", "/favicon.ico"},
		RedisEnabled:    false,
		RedisKeyPrefix:  "ratelimit:",
		EndpointLimits: map[string]EndpointLimit{
			"api_create": {
				Path:      "/api/v1/.*",
				Method:    "POST",
				Limit:     30,  // 0.5 requests per second
				Burst:     5,   // Allow bursts of 5 requests
				UserLimit: 60,  // 1 request per second per user
				UserBurst: 10,  // Allow bursts of 10 requests per user
			},
		},
	}
}
