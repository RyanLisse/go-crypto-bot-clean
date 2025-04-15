package middleware

import (
	"context"
	"net/http"
	"os"

	"github.com/rs/zerolog"
)

// MEXCAPICredentials holds the API credentials for MEXC
type MEXCAPICredentials struct {
	APIKey    string
	APISecret string
}

// MEXCAPICredentialsKey is the context key for MEXC API credentials
type MEXCAPICredentialsKey struct{}

// MEXCAPIMiddleware adds MEXC API credentials to the request context
type MEXCAPIMiddleware struct {
	logger *zerolog.Logger
}

// NewMEXCAPIMiddleware creates a new MEXCAPIMiddleware
func NewMEXCAPIMiddleware(logger *zerolog.Logger) *MEXCAPIMiddleware {
	return &MEXCAPIMiddleware{
		logger: logger,
	}
}

// Middleware returns a middleware function that adds MEXC API credentials to the request context
func (m *MEXCAPIMiddleware) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get API credentials from environment variables
			apiKey := os.Getenv("MEXC_API_KEY")
			apiSecret := os.Getenv("MEXC_SECRET_KEY")

			if apiKey == "" || apiSecret == "" {
				m.logger.Warn().Msg("MEXC_API_KEY or MEXC_SECRET_KEY not set")
				// Continue without credentials
				next.ServeHTTP(w, r)
				return
			}

			// Add credentials to context
			credentials := &MEXCAPICredentials{
				APIKey:    apiKey,
				APISecret: apiSecret,
			}
			ctx := context.WithValue(r.Context(), MEXCAPICredentialsKey{}, credentials)

			// Call the next handler with updated context
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetMEXCAPICredentials gets the MEXC API credentials from the context
func GetMEXCAPICredentials(ctx context.Context) *MEXCAPICredentials {
	credentials, ok := ctx.Value(MEXCAPICredentialsKey{}).(*MEXCAPICredentials)
	if !ok {
		return nil
	}
	return credentials
}
