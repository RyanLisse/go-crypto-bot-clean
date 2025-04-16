package middleware

import (
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

// LoggingMiddleware is a middleware that logs HTTP requests
type LoggingMiddleware struct {
	logger *zerolog.Logger
}

// NewLoggingMiddleware creates a new LoggingMiddleware
func NewLoggingMiddleware(logger *zerolog.Logger) *LoggingMiddleware {
	return &LoggingMiddleware{
		logger: logger,
	}
}

// Middleware returns a middleware function that logs HTTP requests
func (m *LoggingMiddleware) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Start timer
			start := time.Now()

			// Create a response writer that can capture the status code
			crw := &captureResponseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			// Get request ID
			requestID := r.Header.Get("X-Request-ID")

			// Prepare the logger
			requestLogger := m.logger.With().
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Str("remote_addr", r.RemoteAddr).
				Str("user_agent", r.UserAgent()).
				Str("request_id", requestID).
				Logger()

			// Log the request
			requestLogger.Info().Msg("Request started")

			// Call the next handler
			next.ServeHTTP(crw, r)

			// Calculate duration
			duration := time.Since(start)

			// Log the response
			event := requestLogger.Info()
			if crw.Status() >= 400 {
				event = requestLogger.Error()
			}

			event.
				Int("status", crw.Status()).
				Dur("duration", duration).
				Str("duration_human", duration.String()).
				Msg("Request completed")
		})
	}
}
