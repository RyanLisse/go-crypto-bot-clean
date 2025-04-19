package middleware

import (
	"context"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	// Use clean_backend apperror
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/apperror"
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/model" // Import model for GetUserFromContext
	"github.com/google/uuid"                                                       // Keep for request ID generation
	"github.com/rs/zerolog"

	// Add missing import for chi middleware
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

// contextKey is an unexported type for context keys defined in this package.
type contextKey string

// UserIDKey is the key used to store/retrieve the user ID in the request context.
const UserIDKey contextKey = "userID"

// RolesKey is the context key for user roles
const RolesKey contextKey = "roles"

// UserKey is the context key for user model
const UserKey contextKey = "user"

// AuthMiddleware is the interface for authentication middleware
type AuthMiddleware interface {
	Middleware() func(http.Handler) http.Handler
	RequireAuthentication(http.Handler) http.Handler
	RequireRole(role string) func(http.Handler) http.Handler
}

// GetUserFromContext retrieves the *model.User from the context
func GetUserFromContext(ctx context.Context) (*model.User, bool) {
	user, ok := ctx.Value(UserKey).(*model.User)
	return user, ok
}

// GetUserIDFromContext retrieves the user ID string from the context
func GetUserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(UserIDKey).(string)
	return userID, ok
}

// GetRolesFromContext retrieves the user roles slice from the context
func GetRolesFromContext(ctx context.Context) ([]string, bool) {
	roles, ok := ctx.Value(RolesKey).([]string)
	return roles, ok
}

// --- Start: Migrated UnifiedErrorMiddleware logic from old backend ---

// ErrorContextKey is a context key for passing error handling functions through the request context
const ErrorContextKey contextKey = "errorHandler"

// ErrorHandler is a function that handles an error and writes a response
type ErrorHandler func(w http.ResponseWriter, err error, traceID string)

// DefaultErrorHandler is the default error handler used when none is provided
// Uses apperror from clean_backend
func DefaultErrorHandler(w http.ResponseWriter, err error, traceID string) {
	// Use existing AppError creators instead of FromError for now
	var appErr *apperror.AppError
	if ae, ok := err.(*apperror.AppError); ok {
		appErr = ae
	} else {
		appErr = apperror.NewInternal(err)
	}

	// Add trace ID (temporarily in Details)
	detailsMap := map[string]interface{}{"trace_id": traceID}
	if appErr.Details != nil {
		if existingMap, ok := appErr.Details.(map[string]interface{}); ok {
			for k, v := range existingMap {
				detailsMap[k] = v // Merge existing details
			}
		} else {
			detailsMap["original_details"] = appErr.Details
		}
	}
	appErr.Details = detailsMap

	apperror.WriteError(w, appErr)
}

// WithErrorHandler attaches an error handler to the context
func WithErrorHandler(ctx context.Context, handler ErrorHandler) context.Context {
	return context.WithValue(ctx, ErrorContextKey, handler)
}

// GetErrorHandler retrieves the error handler from the context
func GetErrorHandler(ctx context.Context) ErrorHandler {
	if handler, ok := ctx.Value(ErrorContextKey).(ErrorHandler); ok {
		return handler
	}
	return DefaultErrorHandler
}

// RespondWithError writes an error response using the handler from context
func RespondWithError(w http.ResponseWriter, r *http.Request, err error) {
	traceID := GetTraceID(r)
	handler := GetErrorHandler(r.Context())
	handler(w, err, traceID)
}

// GetTraceID extracts the trace ID from the request (using Chi's RequestID middleware)
func GetTraceID(r *http.Request) string {
	if reqID := r.Context().Value(chimiddleware.RequestIDKey); reqID != nil {
		if id, ok := reqID.(string); ok {
			return id
		}
	}
	// Fallback if Chi's RequestID isn't available (shouldn't happen if middleware is used)
	return r.Header.Get("X-Request-ID")
}

// UnifiedErrorMiddleware combines error handling, recovery, logging, and tracing.
type UnifiedErrorMiddleware struct {
	logger *zerolog.Logger
}

// NewUnifiedErrorMiddleware creates a new UnifiedErrorMiddleware.
func NewUnifiedErrorMiddleware(logger *zerolog.Logger) *UnifiedErrorMiddleware {
	return &UnifiedErrorMiddleware{
		logger: logger,
	}
}

// Middleware returns middleware that combines error handling, recovery, tracing, and logging.
func (m *UnifiedErrorMiddleware) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get request ID (should be set by Chi's middleware upstream)
			requestID := GetTraceID(r)
			if requestID == "" {
				requestID = uuid.New().String() // Generate if missing
			}

			// Ensure request ID is on the response header
			w.Header().Set("X-Request-ID", requestID)

			// Set up wrapped response writer
			crw := &unifiedResponseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK, // Default before WriteHeader/Write
				requestID:      requestID,
				logger:         m.logger,
				request:        r,
			}

			// Create error handler context that uses the trace ID
			ctx := WithErrorHandler(r.Context(), func(w http.ResponseWriter, err error, _ string) {
				DefaultErrorHandler(w, err, requestID)
			})

			// Recover from panics
			defer func() {
				if rec := recover(); rec != nil {
					stack := debug.Stack()

					// Log the panic
					m.logger.Error().
						Str("request_id", requestID).
						Interface("panic", rec).
						Str("stack", string(stack)).
						Str("method", r.Method).
						Str("path", r.URL.Path).
						Str("remote_addr", r.RemoteAddr).
						Msg("Panic recovered")

					// Create an app error from the panic
					var appErr *apperror.AppError // Use clean_backend apperror
					switch err := rec.(type) {
					case *apperror.AppError:
						appErr = err
					case error:
						appErr = apperror.NewInternal(err) // Use clean_backend NewInternal
					// case string: // Avoid WrapError for now
					// 	appErr = apperror.NewInternal(apperror.WrapError(nil, err))
					default:
						// Handle non-error panic types more gracefully
						errMsg := fmt.Sprintf("unknown panic type: %T, value: %v", err, err)
						appErr = apperror.NewInternal(fmt.Errorf(errMsg))
					}

					// Add trace ID (temporarily in Details) and write error response
					detailsMap := map[string]interface{}{"trace_id": requestID}
					if appErr.Details != nil {
						if existingMap, ok := appErr.Details.(map[string]interface{}); ok {
							for k, v := range existingMap {
								detailsMap[k] = v // Merge existing details
							}
						} else {
							detailsMap["original_details"] = appErr.Details
						}
					}
					appErr.Details = detailsMap
					apperror.WriteError(w, appErr) // Use clean_backend apperror writer
				}
			}()

			// Log the request
			start := time.Now()
			m.logger.Info().
				Str("request_id", requestID).
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Str("remote_addr", r.RemoteAddr).
				Msg("Request received")

			// Serve the next handler
			next.ServeHTTP(crw, r.WithContext(ctx))

			// Log the response
			duration := time.Since(start)
			m.logger.Info().
				Str("request_id", requestID).
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Int("status", crw.Status()).
				Dur("duration", duration). // Log the duration
				Msg("Request completed")
		})
	}
}

// RequestTimingMiddleware logs request duration. (Example of another middleware)
func RequestTimingMiddleware(logger *zerolog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			duration := time.Since(start)
			// Get request ID if available
			requestID := "N/A"
			if reqID := r.Context().Value(chimiddleware.RequestIDKey); reqID != nil {
				if id, ok := reqID.(string); ok {
					requestID = id
				}
			}
			logger.Info().Str("request_id", requestID).Str("method", r.Method).Str("path", r.URL.Path).Dur("duration", duration).Msg("Request timed")
		})
	}
}

// unifiedResponseWriter is a wrapper to capture status code (from old backend, needs review)
type unifiedResponseWriter struct {
	http.ResponseWriter
	statusCode    int
	requestID     string
	logger        *zerolog.Logger
	request       *http.Request
	headerWritten bool // Track if header has been written
}

func (rw *unifiedResponseWriter) WriteHeader(code int) {
	if rw.headerWritten {
		return // Ensure header is written only once
	}
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
	rw.headerWritten = true
}

func (rw *unifiedResponseWriter) Write(b []byte) (int, error) {
	if !rw.headerWritten {
		// If WriteHeader hasn't been called explicitly, call it with 200 OK
		rw.WriteHeader(http.StatusOK)
	}
	return rw.ResponseWriter.Write(b)
}

func (rw *unifiedResponseWriter) Status() int {
	return rw.statusCode
}

// min returns the minimum of two integers
// Currently unused but kept for future use
// func min(a, b int) int {
// 	if a < b {
// 		return a
// 	}
// 	return b
// }
