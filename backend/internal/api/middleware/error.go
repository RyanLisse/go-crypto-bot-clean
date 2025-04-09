package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go-crypto-bot-clean/backend/internal/auth"
	"go-crypto-bot-clean/backend/internal/validation"

	"go.uber.org/zap"
)

// ErrorResponse represents a standardized error response
type ErrorResponse struct {
	Code      string      `json:"code"`                // Error code for client handling
	Message   string      `json:"message"`             // User-friendly error message
	Details   interface{} `json:"details,omitempty"`   // Optional detailed error information
	Help      string      `json:"help,omitempty"`      // Optional help text
	RequestID string      `json:"requestId,omitempty"` // Request ID for tracing
	Path      string      `json:"path"`                // Request path
	Method    string      `json:"method"`              // HTTP method
	Timestamp time.Time   `json:"timestamp"`           // Error timestamp
	Latency   int64       `json:"latency,omitempty"`   // Request processing time
}

// ErrorCtxKey is the context key for storing errors
const ErrorCtxKey = "error"

type responseWriter struct {
	http.ResponseWriter
	status int
	size   int64
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(b)
	rw.size += int64(size)
	return size, err
}

// ErrorHandlingMiddleware wraps an http.Handler and handles errors
func ErrorHandlingMiddleware(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}

			// Generate request ID if not present
			requestID := r.Header.Get("X-Request-ID")
			if requestID == "" {
				requestID = fmt.Sprintf("%d", time.Now().UnixNano())
				r.Header.Set("X-Request-ID", requestID)
			}

			// Add request logger with context
			reqLogger := logger.With(
				zap.String("request_id", requestID),
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.String("remote_addr", r.RemoteAddr),
				zap.String("user_agent", r.UserAgent()),
			)

			next.ServeHTTP(rw, r)

			latency := time.Since(start).Milliseconds()

			if err := r.Context().Value(ErrorCtxKey); err != nil {
				var response ErrorResponse
				response.Path = r.URL.Path
				response.Method = r.Method
				response.Timestamp = time.Now()
				response.Latency = latency
				response.RequestID = requestID

				switch e := err.(type) {
				case *auth.AuthError:
					response.Code = string(e.Type)
					response.Message = e.Message
					response.Details = e.Details
					response.Help = e.Help

					reqLogger.Warn("Authentication error",
						zap.String("error_type", string(e.Type)),
						zap.String("error_message", e.Message),
						zap.Int("status_code", e.Code),
						zap.Int64("latency_ms", latency),
						zap.Any("details", e.Details),
					)

					rw.WriteHeader(e.Code)

				case *validation.ValidationError:
					response.Code = "validation_error"
					response.Message = e.Error()
					response.Details = e.Details
					response.Help = "Please check the provided data and try again"

					reqLogger.Warn("Validation error",
						zap.String("error_message", e.Error()),
						zap.Int64("latency_ms", latency),
						zap.Any("details", e.Details),
					)

					rw.WriteHeader(http.StatusBadRequest)

				default:
					response.Code = "internal_error"
					response.Message = "An unexpected error occurred"
					response.Help = "Please try again later or contact support if the problem persists"

					reqLogger.Error("Unexpected error",
						zap.Error(e.(error)),
						zap.Int64("latency_ms", latency),
						zap.Int("response_size", int(rw.size)),
					)

					rw.WriteHeader(http.StatusInternalServerError)
				}

				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("X-Content-Type-Options", "nosniff")
				w.Header().Set("X-Request-ID", requestID)
				json.NewEncoder(w).Encode(response)
			} else {
				// Log successful requests
				reqLogger.Info("Request completed",
					zap.Int("status", rw.status),
					zap.Int64("latency_ms", latency),
					zap.Int64("response_size", rw.size),
				)
			}
		})
	}
}
