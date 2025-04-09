package middleware

import (
	"bytes"
	"context"
	"net/http"
	"regexp"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// SecurityConfig contains configuration for security middleware
type SecurityConfig struct {
	// EnableInputSanitization enables input sanitization
	EnableInputSanitization bool

	// EnableOutputFiltering enables output filtering
	EnableOutputFiltering bool

	// MaxRequestSize is the maximum size of a request body in bytes
	MaxRequestSize int64

	// AllowedDomains is a list of allowed domains for CORS
	AllowedDomains []string

	// Logger is the logger to use
	Logger *zap.Logger
}

// DefaultSecurityConfig returns the default security configuration
func DefaultSecurityConfig() SecurityConfig {
	return SecurityConfig{
		EnableInputSanitization: true,
		EnableOutputFiltering:   true,
		MaxRequestSize:          1024 * 1024, // 1MB
		AllowedDomains:          []string{"localhost"},
		Logger:                  zap.NewNop(),
	}
}

// SecurityMiddleware returns a middleware that adds security features
func SecurityMiddleware(config SecurityConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Add security headers
			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.Header().Set("X-Frame-Options", "DENY")
			w.Header().Set("X-XSS-Protection", "1; mode=block")
			w.Header().Set("Content-Security-Policy", "default-src 'self'")
			w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

			// Limit request size
			r.Body = http.MaxBytesReader(w, r.Body, config.MaxRequestSize)

			// Check for suspicious input patterns
			if config.EnableInputSanitization {
				path := r.URL.Path
				query := r.URL.RawQuery

				// Check for SQL injection patterns
				sqlInjectionPattern := regexp.MustCompile(`(?i)(union\s+select|select\s+.*\s+from|insert\s+into|update\s+.*\s+set|delete\s+from|drop\s+table|exec\s+xp_|exec\s+sp_|exec\s+master|declare\s+@|;--|\bor\s+1=1\b)`)
				if sqlInjectionPattern.MatchString(path) || sqlInjectionPattern.MatchString(query) {
					config.Logger.Warn("Potential SQL injection attempt",
						zap.String("path", path),
						zap.String("query", query),
						zap.String("ip", r.RemoteAddr),
					)
					http.Error(w, "Invalid request", http.StatusBadRequest)
					return
				}

				// Check for XSS patterns
				xssPattern := regexp.MustCompile(`(?i)(<script|javascript:|on\w+\s*=|alert\s*\(|eval\s*\(|document\.cookie)`)
				if xssPattern.MatchString(path) || xssPattern.MatchString(query) {
					config.Logger.Warn("Potential XSS attempt",
						zap.String("path", path),
						zap.String("query", query),
						zap.String("ip", r.RemoteAddr),
					)
					http.Error(w, "Invalid request", http.StatusBadRequest)
					return
				}
			}

			// Add request ID to context for tracing
			requestID := r.Header.Get("X-Request-ID")
			if requestID == "" {
				requestID = generateRequestID()
				r.Header.Set("X-Request-ID", requestID)
			}
			w.Header().Set("X-Request-ID", requestID)

			// Add request ID to context
			ctx := context.WithValue(r.Context(), "request_id", requestID)
			r = r.WithContext(ctx)

			// Call the next handler
			next.ServeHTTP(w, r)
		})
	}
}

// OutputFilterMiddleware filters sensitive information from responses
func OutputFilterMiddleware(config SecurityConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !config.EnableOutputFiltering {
				next.ServeHTTP(w, r)
				return
			}

			// Create a response writer wrapper that filters output
			fw := &filterWriter{
				ResponseWriter: w,
				config:         config,
				logger:         config.Logger,
				requestID:      r.Context().Value("request_id").(string),
			}

			// Call the next handler with the wrapped response writer
			next.ServeHTTP(fw, r)
		})
	}
}

// filterWriter is a response writer that filters sensitive information
type filterWriter struct {
	http.ResponseWriter
	config    SecurityConfig
	logger    *zap.Logger
	requestID string
	body      []byte
}

// Write implements http.ResponseWriter
func (fw *filterWriter) Write(b []byte) (int, error) {
	// Filter sensitive information
	filtered := filterSensitiveData(b)

	// Log if sensitive data was found
	if !bytes.Equal(b, filtered) {
		fw.logger.Warn("Sensitive data filtered from response",
			zap.String("request_id", fw.requestID),
		)
	}

	// Write filtered data
	return fw.ResponseWriter.Write(filtered)
}

// filterSensitiveData filters sensitive information from data
func filterSensitiveData(data []byte) []byte {
	// Convert to string for easier processing
	s := string(data)

	// Filter credit card numbers
	ccPattern := regexp.MustCompile(`\b(?:\d{4}[-\s]?){3}\d{4}\b`)
	s = ccPattern.ReplaceAllString(s, "[REDACTED CREDIT CARD]")

	// Filter social security numbers
	ssnPattern := regexp.MustCompile(`\b\d{3}[-\s]?\d{2}[-\s]?\d{4}\b`)
	s = ssnPattern.ReplaceAllString(s, "[REDACTED SSN]")

	// Filter API keys (common patterns)
	apiKeyPattern := regexp.MustCompile(`(?i)\b(api[-_]?key|apikey|access[-_]?key|auth[-_]?token|client[-_]?secret)[-_]?[:=]\s*["']?([a-zA-Z0-9]{16,})["']?`)
	s = apiKeyPattern.ReplaceAllString(s, "$1: [REDACTED API KEY]")

	// Filter passwords
	passwordPattern := regexp.MustCompile(`(?i)"(password|passwd|pwd)":\s*"[^"]*"`)
	s = passwordPattern.ReplaceAllString(s, `"$1":"[REDACTED]"`)

	return []byte(s)
}

// generateRequestID generates a unique request ID
func generateRequestID() string {
	// Generate a random UUID
	return uuid.New().String()
}

// RegisterSecurityMiddleware registers security middleware with a Chi router
func RegisterSecurityMiddleware(r chi.Router, logger *zap.Logger) {
	config := DefaultSecurityConfig()
	config.Logger = logger

	// Add security middleware
	r.Use(SecurityMiddleware(config))
	r.Use(OutputFilterMiddleware(config))
}
