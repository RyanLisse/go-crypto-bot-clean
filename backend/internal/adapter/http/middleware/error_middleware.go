package middleware

import (
	"errors"
	"net/http"
	"runtime/debug"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/apperror"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

// ErrorMiddleware is a middleware that handles errors
type ErrorMiddleware struct {
	logger *zerolog.Logger
}

// NewErrorMiddleware creates a new ErrorMiddleware
func NewErrorMiddleware(logger *zerolog.Logger) *ErrorMiddleware {
	return &ErrorMiddleware{
		logger: logger,
	}
}

// Middleware returns a middleware function that handles errors
func (m *ErrorMiddleware) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Generate a request ID if not already present
			requestID := r.Header.Get("X-Request-ID")
			if requestID == "" {
				requestID = uuid.New().String()
				r.Header.Set("X-Request-ID", requestID)
			}

			// Create a response writer that can capture the status code
			crw := &captureResponseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			// Recover from panics
			defer func() {
				if err := recover(); err != nil {
					// Log the panic
					m.logger.Error().
						Str("request_id", requestID).
						Interface("panic", err).
						Str("stack", string(debug.Stack())).
						Msg("Panic recovered in error middleware")

					// Create an internal server error
					var appErr *apperror.AppError
					switch e := err.(type) {
					case error:
						appErr = apperror.NewInternal(e)
					case string:
						appErr = apperror.NewInternal(errors.New(e))
					default:
						appErr = apperror.NewInternal(errors.New("unknown panic"))
					}

					// Write the error response
					w.Header().Set("X-Request-ID", requestID)
					apperror.WriteError(w, appErr)
				}
			}()

			// Set the request ID header in the response
			crw.Header().Set("X-Request-ID", requestID)

			// Call the next handler
			next.ServeHTTP(crw, r)
		})
	}
}

// captureResponseWriter is a wrapper around http.ResponseWriter that captures the status code
type captureResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code and calls the underlying ResponseWriter's WriteHeader
func (crw *captureResponseWriter) WriteHeader(statusCode int) {
	crw.statusCode = statusCode
	crw.ResponseWriter.WriteHeader(statusCode)
}

// Status returns the captured status code
func (crw *captureResponseWriter) Status() int {
	return crw.statusCode
}
