package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/apperror"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

// UnifiedErrorMiddleware combines error handling, recovery, logging, and tracing
// into a single middleware component for consistent error handling.
type UnifiedErrorMiddleware struct {
	logger *zerolog.Logger
}

// NewUnifiedErrorMiddleware creates a new UnifiedErrorMiddleware
func NewUnifiedErrorMiddleware(logger *zerolog.Logger) *UnifiedErrorMiddleware {
	return &UnifiedErrorMiddleware{
		logger: logger,
	}
}

// Middleware returns middleware that combines error handling, recovery, tracing, and logging
func (m *UnifiedErrorMiddleware) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Generate or get request ID for tracing
			requestID := r.Header.Get("X-Request-ID")
			if requestID == "" {
				requestID = uuid.New().String()
				r.Header.Set("X-Request-ID", requestID)
			}

			// Set up wrapped response writer that captures errors
			crw := &unifiedResponseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
				requestID:      requestID,
				logger:         m.logger,
				request:        r,
			}

			// Create error handler context that uses the trace ID
			ctx := apperror.WithErrorHandler(r.Context(), func(w http.ResponseWriter, err error, _ string) {
				apperror.DefaultErrorHandler(w, err, requestID)
			})

			// Add the request ID to response headers
			w.Header().Set("X-Request-ID", requestID)

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
					var appErr *apperror.AppError
					switch err := rec.(type) {
					case *apperror.AppError:
						appErr = err
					case error:
						appErr = apperror.NewInternal(err)
					case string:
						appErr = apperror.NewInternal(apperror.WrapError(nil, err))
					default:
						appErr = apperror.NewInternal(apperror.WrapError(nil, "unknown panic"))
					}

					// Write error response with trace ID
					apperror.WriteErrorWithTraceID(w, appErr, requestID)
				}
			}()

			// Log the request
			m.logger.Info().
				Str("request_id", requestID).
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Str("remote_addr", r.RemoteAddr).
				Str("user_agent", r.UserAgent()).
				Msg("Request received")

			// Call the next handler with the context containing the error handler
			next.ServeHTTP(crw, r.WithContext(ctx))

			// Log the response
			logEvent := m.logger.Info()
			if crw.statusCode >= 400 {
				logEvent = m.logger.Error()
			}

			logEvent.
				Str("request_id", requestID).
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Int("status_code", crw.statusCode).
				Str("remote_addr", r.RemoteAddr).
				Msg("Response completed")
		})
	}
}

// unifiedResponseWriter wraps http.ResponseWriter to capture status codes
// and standardize error responses
type unifiedResponseWriter struct {
	http.ResponseWriter
	statusCode int
	requestID  string
	logger     *zerolog.Logger
	request    *http.Request
}

// WriteHeader captures status code and standardizes error responses
func (rw *unifiedResponseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Write captures the write and ensures status code is set
func (rw *unifiedResponseWriter) Write(b []byte) (int, error) {
	// If status code hasn't been explicitly set, default to 200 OK
	if rw.statusCode == 0 {
		rw.statusCode = http.StatusOK
	}
	return rw.ResponseWriter.Write(b)
}

// Status returns the status code
func (rw *unifiedResponseWriter) Status() int {
	return rw.statusCode
}

// Flush implements http.Flusher if the underlying ResponseWriter supports it
func (rw *unifiedResponseWriter) Flush() {
	if f, ok := rw.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

// Hijack implements http.Hijacker if the underlying ResponseWriter supports it
func (rw *unifiedResponseWriter) Hijack() (conn interface{}, buf interface{}, err error) {
	if hj, ok := rw.ResponseWriter.(http.Hijacker); ok {
		return hj.Hijack()
	}
	return nil, nil, http.ErrNotSupported
}
