package middleware

import (
	"net/http"
	"strings"

	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/config"
	"github.com/rs/zerolog"
)

// SecureHeadersMiddleware is a middleware that adds secure HTTP headers
type SecureHeadersMiddleware struct {
	config *config.SecureHeadersConfig
	logger *zerolog.Logger
}

// NewSecureHeadersMiddleware creates a new SecureHeadersMiddleware
func NewSecureHeadersMiddleware(cfg *config.SecureHeadersConfig, logger *zerolog.Logger) *SecureHeadersMiddleware {
	return &SecureHeadersMiddleware{
		config: cfg,
		logger: logger,
	}
}

// Middleware returns a middleware function that adds secure HTTP headers
func (m *SecureHeadersMiddleware) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check if secure headers are enabled
			if !m.config.Enabled {
				next.ServeHTTP(w, r)
				return
			}

			// Check if the path is excluded
			path := r.URL.Path
			for _, excludedPath := range m.config.ExcludedPaths {
				if strings.HasPrefix(path, excludedPath) {
					next.ServeHTTP(w, r)
					return
				}
			}

			// Set Content-Security-Policy header
			if m.config.ContentSecurityPolicy != "" {
				if m.config.ContentSecurityPolicyReportOnly {
					w.Header().Set("Content-Security-Policy-Report-Only", m.config.ContentSecurityPolicy)
					if m.config.ContentSecurityPolicyReportURI != "" {
						w.Header().Set("Content-Security-Policy-Report-Only", w.Header().Get("Content-Security-Policy-Report-Only")+"; report-uri "+m.config.ContentSecurityPolicyReportURI)
					}
				} else {
					w.Header().Set("Content-Security-Policy", m.config.ContentSecurityPolicy)
					if m.config.ContentSecurityPolicyReportURI != "" {
						w.Header().Set("Content-Security-Policy", w.Header().Get("Content-Security-Policy")+"; report-uri "+m.config.ContentSecurityPolicyReportURI)
					}
				}
			}

			// Set X-Content-Type-Options header
			if m.config.XContentTypeOptions != "" {
				w.Header().Set("X-Content-Type-Options", m.config.XContentTypeOptions)
			}

			// Set X-Frame-Options header
			if m.config.XFrameOptions != "" {
				w.Header().Set("X-Frame-Options", m.config.XFrameOptions)
			}

			// Set X-XSS-Protection header
			if m.config.XXSSProtection != "" {
				w.Header().Set("X-XSS-Protection", m.config.XXSSProtection)
			}

			// Set Referrer-Policy header
			if m.config.ReferrerPolicy != "" {
				w.Header().Set("Referrer-Policy", m.config.ReferrerPolicy)
			}

			// Set Strict-Transport-Security header
			if m.config.StrictTransportSecurity != "" {
				w.Header().Set("Strict-Transport-Security", m.config.StrictTransportSecurity)
			}

			// Set Permissions-Policy header
			if m.config.PermissionsPolicy != "" {
				w.Header().Set("Permissions-Policy", m.config.PermissionsPolicy)
			}

			// Set Cross-Origin-Embedder-Policy header
			if m.config.CrossOriginEmbedderPolicy != "" {
				w.Header().Set("Cross-Origin-Embedder-Policy", m.config.CrossOriginEmbedderPolicy)
			}

			// Set Cross-Origin-Opener-Policy header
			if m.config.CrossOriginOpenerPolicy != "" {
				w.Header().Set("Cross-Origin-Opener-Policy", m.config.CrossOriginOpenerPolicy)
			}

			// Set Cross-Origin-Resource-Policy header
			if m.config.CrossOriginResourcePolicy != "" {
				w.Header().Set("Cross-Origin-Resource-Policy", m.config.CrossOriginResourcePolicy)
			}

			// Set Cache-Control header
			if m.config.CacheControl != "" {
				w.Header().Set("Cache-Control", m.config.CacheControl)
			}

			// Set custom headers
			for name, value := range m.config.CustomHeaders {
				w.Header().Set(name, value)
			}

			// Remove Server header
			if m.config.RemoveServerHeader {
				w.Header().Del("Server")
			}

			// Remove X-Powered-By header
			if m.config.RemovePoweredByHeader {
				w.Header().Del("X-Powered-By")
			}

			// Call the next handler
			next.ServeHTTP(w, r)
		})
	}
}

// SecureHeadersHandler is a handler that adds secure HTTP headers
func SecureHeadersHandler(cfg *config.SecureHeadersConfig, logger *zerolog.Logger) func(http.Handler) http.Handler {
	middleware := NewSecureHeadersMiddleware(cfg, logger)
	return middleware.Middleware()
}
