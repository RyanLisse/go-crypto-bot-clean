package config

import (
	"time"
)

// CSRFConfig contains CSRF protection configuration
type CSRFConfig struct {
	Enabled           bool          `mapstructure:"enabled"`
	Secret            string        `mapstructure:"secret"`
	TokenLength       int           `mapstructure:"token_length"`
	CookieName        string        `mapstructure:"cookie_name"`
	CookiePath        string        `mapstructure:"cookie_path"`
	CookieDomain      string        `mapstructure:"cookie_domain"`
	CookieMaxAge      time.Duration `mapstructure:"cookie_max_age"`
	CookieSecure      bool          `mapstructure:"cookie_secure"`
	CookieHTTPOnly    bool          `mapstructure:"cookie_http_only"`
	CookieSameSite    string        `mapstructure:"cookie_same_site"`
	HeaderName        string        `mapstructure:"header_name"`
	FormFieldName     string        `mapstructure:"form_field_name"`
	ExcludedPaths     []string      `mapstructure:"excluded_paths"`
	ExcludedMethods   []string      `mapstructure:"excluded_methods"`
	FailureStatusCode int           `mapstructure:"failure_status_code"`
}

// GetDefaultCSRFConfig returns the default CSRF configuration
func GetDefaultCSRFConfig() CSRFConfig {
	return CSRFConfig{
		Enabled:           true,
		TokenLength:       32,
		CookieName:        "csrf_token",
		CookiePath:        "/",
		CookieMaxAge:      24 * time.Hour,
		CookieSecure:      true,
		CookieHTTPOnly:    true,
		CookieSameSite:    "Lax",
		HeaderName:        "X-CSRF-Token",
		FormFieldName:     "csrf_token",
		ExcludedPaths:     []string{"/health", "/metrics", "/favicon.ico"},
		ExcludedMethods:   []string{"GET", "HEAD", "OPTIONS", "TRACE"},
		FailureStatusCode: 403,
	}
}
