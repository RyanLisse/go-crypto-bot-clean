// Package middleware contains API middleware components.
package middleware

import (
	"net/http"
	"time"
)

// Logger defines the minimal logger interface for middleware.
type Logger interface {
	Info(args ...interface{})
	Error(args ...interface{})
}

// LoggingMiddleware logs incoming requests.
//
//	@summary	Logging middleware
//	@description	Logs method, path, status, latency, request ID, and errors.
func LoggingMiddleware(logger Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Wrap the ResponseWriter to capture status code
			rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			next.ServeHTTP(rw, r)

			latency := time.Since(start)
			status := rw.statusCode
			method := r.Method
			path := r.URL.Path
			requestID := r.Header.Get("X-Request-ID")

			logFields := []interface{}{
				"method", method,
				"path", path,
				"status", status,
				"latency", latency,
				"request_id", requestID,
			}

			logger.Info(logFields...)
		})
	}
}

// responseWriter wraps http.ResponseWriter to capture status code
