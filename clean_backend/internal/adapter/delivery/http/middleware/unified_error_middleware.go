package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/apperror"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

// EnhancedErrorMiddleware handles errors and panics in a consistent way
type EnhancedErrorMiddleware struct {
	logger *zerolog.Logger
}

// NewEnhancedErrorMiddleware creates a new EnhancedErrorMiddleware
func NewEnhancedErrorMiddleware(logger *zerolog.Logger) *EnhancedErrorMiddleware {
	return &EnhancedErrorMiddleware{
		logger: logger,
	}
}

// Middleware returns a middleware function that handles errors and panics
func (m *EnhancedErrorMiddleware) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Add request ID if not present
			requestID := r.Header.Get("X-Request-ID")
			if requestID == "" {
				requestID = uuid.New().String()
				r.Header.Set("X-Request-ID", requestID)
			}
			w.Header().Set("X-Request-ID", requestID)

			// Create a response writer wrapper to capture status code
			rw := &enhancedResponseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			// Recover from panics
			defer func() {
				if err := recover(); err != nil {
					// Log the panic with stack trace
					stackTrace := debug.Stack()
					m.logger.Error().
						Str("request_id", requestID).
						Interface("error", err).
						Str("stack_trace", string(stackTrace)).
						Msg("Panic recovered in HTTP handler")

					// Convert panic to appropriate error response
					switch e := err.(type) {
					case *apperror.AppError:
						// If it's already an AppError, use it directly
						apperror.RespondWithError(rw, r, e)
					case error:
						// If it's an error, wrap it in an internal server error
						appErr := apperror.NewInternal(e)
						apperror.RespondWithError(rw, r, appErr)
					default:
						// For any other type, convert to string and wrap in internal server error
						appErr := apperror.NewInternal(fmt.Errorf("%v", e))
						apperror.RespondWithError(rw, r, appErr)
					}
				}
			}()

			// Call the next handler
			next.ServeHTTP(rw, r)
		})
	}
}

// enhancedResponseWriter is a wrapper around http.ResponseWriter that captures the status code
type enhancedResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code and calls the underlying WriteHeader
func (rw *enhancedResponseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

// Status returns the captured status code
func (rw *enhancedResponseWriter) Status() int {
	return rw.statusCode
}
