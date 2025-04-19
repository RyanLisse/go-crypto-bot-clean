package middleware

import (
	"net/http"

	"github.com/rs/zerolog"
)

// DisabledMiddleware is a middleware that bypasses authentication for all requests
type DisabledMiddleware struct {
	logger *zerolog.Logger
}

// NewDisabledMiddleware creates a new DisabledMiddleware
func NewDisabledMiddleware(logger *zerolog.Logger) AuthMiddleware {
	return &DisabledMiddleware{
		logger: logger,
	}
}

// Middleware returns a middleware function that does nothing
func (m *DisabledMiddleware) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Authentication is disabled, just pass through
			m.logger.Debug().Msg("Authentication disabled, bypassing checks")
			next.ServeHTTP(w, r)
		})
	}
}

// RequireAuthentication is a middleware that requires authentication (disabled)
func (m *DisabledMiddleware) RequireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Authentication is disabled, just pass through
		m.logger.Debug().Msg("Authentication requirement bypassed (disabled)")
		next.ServeHTTP(w, r)
	})
}

// RequireRole is a middleware that requires a specific role (disabled)
func (m *DisabledMiddleware) RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Authentication is disabled, just pass through
			m.logger.Debug().Str("role", role).Msg("Role requirement bypassed (disabled)")
			next.ServeHTTP(w, r)
		})
	}
}
