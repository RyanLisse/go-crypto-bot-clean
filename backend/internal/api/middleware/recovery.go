package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"runtime/debug"
	"strings"

	"go-crypto-bot-clean/backend/internal/auth"
)

// RequestIDKey is the key type for request ID in context
type requestIDKey int

// RequestIDContextKey is the context key for request ID
const RequestIDContextKey requestIDKey = iota

// GetRequestID retrieves the request ID from the context
func GetRequestID(ctx context.Context) string {
	if id, ok := ctx.Value(RequestIDContextKey).(string); ok {
		return id
	}
	return "unknown"
}

// RecoveryMiddleware recovers from panics and returns an appropriate error response
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Get stack trace
				stack := debug.Stack()
				stackTrace := processStackTrace(stack, 10) // Limit to 10 frames

				// Determine error type and status
				errorType, _, statusCode := determineErrorType(err)

				// Get sanitized request info
				reqInfo := getSanitizedRequestInfo(r, []string{"Authorization", "Cookie"}, []string{"token", "password"})

				// Create error response with metadata in details
				metadata := map[string]interface{}{
					"stack_trace": stackTrace,
					"request":     reqInfo,
				}

				// Create error response
				errResp := auth.NewAuthError(
					auth.ErrorType(errorType),
					"An unexpected error occurred",
					statusCode,
				).WithDetails(metadata).
					WithHelp("Please try again later or contact support if the issue persists").
					WithRequestID(GetRequestID(r.Context()))

				// Set response headers
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(statusCode)

				// Write error response
				errResp.WriteJSON(w)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// WriteJSON writes a JSON response with proper error handling
func WriteJSON(w http.ResponseWriter, v interface{}) error {
	return json.NewEncoder(w).Encode(v)
}

// processStackTrace formats and limits the stack trace
func processStackTrace(stack []byte, maxLines int) string {
	lines := strings.Split(string(stack), "\n")
	if len(lines) > maxLines*2 { // Each frame is 2 lines
		lines = lines[:maxLines*2]
		lines = append(lines, "... (truncated)")
	}
	return strings.Join(lines, "\n")
}

// getSanitizedRequestInfo extracts and sanitizes request information
func getSanitizedRequestInfo(r *http.Request, redactHeaders, redactParams []string) map[string]interface{} {
	// Create a map to hold request info
	info := map[string]interface{}{
		"url":     r.URL.String(),
		"method":  r.Method,
		"host":    r.Host,
		"path":    r.URL.Path,
		"remote":  r.RemoteAddr,
		"headers": make(map[string]string),
	}

	// Add sanitized headers
	headers := make(map[string]string)
	for name, values := range r.Header {
		// Check if this header should be redacted
		should_redact := false
		for _, h := range redactHeaders {
			if strings.EqualFold(name, h) {
				should_redact = true
				break
			}
		}

		if should_redact {
			headers[name] = "[REDACTED]"
		} else if len(values) > 0 {
			headers[name] = values[0]
		}
	}
	info["headers"] = headers

	// Add sanitized query parameters
	query := make(map[string]string)
	for name, values := range r.URL.Query() {
		// Check if this parameter should be redacted
		should_redact := false
		for _, p := range redactParams {
			if strings.EqualFold(name, p) {
				should_redact = true
				break
			}
		}

		if should_redact {
			query[name] = "[REDACTED]"
		} else if len(values) > 0 {
			query[name] = values[0]
		}
	}
	info["query"] = query

	return info
}

// determineErrorType analyzes the panic error and returns appropriate error type, code and status
func determineErrorType(err interface{}) (string, string, int) {
	// Default values
	errorType := "internal_error"
	errorCode := "INTERNAL_SERVER_ERROR"
	statusCode := http.StatusInternalServerError

	// Try to determine more specific error type based on the panic value
	switch e := err.(type) {
	case *auth.AuthError:
		// Already an auth error, use its values
		return string(e.Type), string(e.Type), e.Code
	case error:
		// Check for common error patterns
		errStr := e.Error()
		switch {
		case strings.Contains(errStr, "token") && strings.Contains(errStr, "invalid"):
			return "invalid_token", "INVALID_TOKEN", http.StatusUnauthorized
		case strings.Contains(errStr, "token") && strings.Contains(errStr, "expired"):
			return "expired_token", "EXPIRED_TOKEN", http.StatusUnauthorized
		case strings.Contains(errStr, "permission") || strings.Contains(errStr, "forbidden"):
			return "permission_denied", "PERMISSION_DENIED", http.StatusForbidden
		case strings.Contains(errStr, "not found") || strings.Contains(errStr, "404"):
			return "not_found", "NOT_FOUND", http.StatusNotFound
		case strings.Contains(errStr, "timeout") || strings.Contains(errStr, "deadline"):
			return "timeout", "REQUEST_TIMEOUT", http.StatusGatewayTimeout
		case strings.Contains(errStr, "rate limit") || strings.Contains(errStr, "too many"):
			return "rate_limit_exceeded", "RATE_LIMIT_EXCEEDED", http.StatusTooManyRequests
		}
	}

	return errorType, errorCode, statusCode
}
