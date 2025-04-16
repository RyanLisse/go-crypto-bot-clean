package apperror

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// ErrorResponse is a standardized error response format
// that provides a consistent structure for all API errors
type ErrorResponse struct {
	Status      int               `json:"status"`                 // HTTP status code
	Code        string            `json:"code"`                   // Error code (e.g., INVALID_INPUT, NOT_FOUND)
	Message     string            `json:"message"`                // User-friendly error message
	Details     interface{}       `json:"details,omitempty"`      // Additional error details (may contain validation errors)
	TraceID     string            `json:"trace_id,omitempty"`     // Request trace ID for correlation
	FieldErrors map[string]string `json:"field_errors,omitempty"` // Field-specific validation errors
}

// ErrorContext is a context key for passing error handling functions through the request context
type ErrorContext struct{}

// ErrorHandler is a function that handles an error and writes a response
type ErrorHandler func(w http.ResponseWriter, err error, traceID string)

// DefaultErrorHandler is the default error handler used when none is provided
func DefaultErrorHandler(w http.ResponseWriter, err error, traceID string) {
	var appErr *AppError
	if !errors.As(err, &appErr) {
		// Convert non-AppError to AppError
		appErr = NewInternal(err)
	}

	// Write the response with trace ID
	WriteErrorWithTraceID(w, appErr, traceID)
}

// WriteErrorWithTraceID writes an error response with a trace ID
func WriteErrorWithTraceID(w http.ResponseWriter, err *AppError, traceID string) {
	resp := err.ToResponse()

	// Add trace ID to the response if provided
	if traceID != "" {
		errorMap := resp["error"].(map[string]interface{})
		errorMap["trace_id"] = traceID
	}

	// Write JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.StatusCode)
	json.NewEncoder(w).Encode(resp)
}

// WithErrorHandler attaches an error handler to the context
func WithErrorHandler(ctx context.Context, handler ErrorHandler) context.Context {
	return context.WithValue(ctx, ErrorContext{}, handler)
}

// GetErrorHandler retrieves the error handler from the context
func GetErrorHandler(ctx context.Context) ErrorHandler {
	if handler, ok := ctx.Value(ErrorContext{}).(ErrorHandler); ok {
		return handler
	}
	return DefaultErrorHandler
}

// RespondWithError writes an error response using the handler from context
func RespondWithError(w http.ResponseWriter, r *http.Request, err error) {
	// Get trace ID from request
	traceID := GetTraceID(r)

	// Get error handler from context
	handler := GetErrorHandler(r.Context())

	// Handle the error
	handler(w, err, traceID)
}

// GetTraceID extracts the trace ID from the request
func GetTraceID(r *http.Request) string {
	return r.Header.Get("X-Request-ID")
}

// WriteValidationError writes a validation error response
func WriteValidationError(w http.ResponseWriter, field string, message string, traceID string) {
	fieldErrors := map[string]string{field: message}
	err := NewInvalid("Validation error", fieldErrors, nil)

	resp := err.ToResponse()

	// Add trace ID and field errors to the response
	errorMap := resp["error"].(map[string]interface{})
	if traceID != "" {
		errorMap["trace_id"] = traceID
	}
	errorMap["field_errors"] = fieldErrors

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.StatusCode)
	json.NewEncoder(w).Encode(resp)
}

// WriteValidationErrors writes multiple validation errors
func WriteValidationErrors(w http.ResponseWriter, fieldErrors map[string]string, traceID string) {
	err := NewInvalid("Validation errors", fieldErrors, nil)

	resp := err.ToResponse()

	// Add trace ID and field errors to the response
	errorMap := resp["error"].(map[string]interface{})
	if traceID != "" {
		errorMap["trace_id"] = traceID
	}
	errorMap["field_errors"] = fieldErrors

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.StatusCode)
	json.NewEncoder(w).Encode(resp)
}

// WrapError wraps an existing error with additional context
func WrapError(err error, message string) error {
	if err == nil {
		return nil
	}

	var appErr *AppError
	if errors.As(err, &appErr) {
		// Clone the AppError but with the new message
		newErr := *appErr
		if message != "" {
			newErr.Message = fmt.Sprintf("%s: %s", message, appErr.Message)
		}
		return &newErr
	}

	// Wrap a non-AppError
	return fmt.Errorf("%s: %w", message, err)
}

// IsNotFound checks if the error is a not found error
func IsNotFound(err error) bool {
	var appErr *AppError
	return errors.As(err, &appErr) && appErr.Code == "NOT_FOUND"
}

// IsUnauthorized checks if the error is an unauthorized error
func IsUnauthorized(err error) bool {
	var appErr *AppError
	return errors.As(err, &appErr) && appErr.Code == "UNAUTHORIZED"
}

// IsForbidden checks if the error is a forbidden error
func IsForbidden(err error) bool {
	var appErr *AppError
	return errors.As(err, &appErr) && appErr.Code == "FORBIDDEN"
}

// IsInvalid checks if the error is an invalid input error
func IsInvalid(err error) bool {
	var appErr *AppError
	return errors.As(err, &appErr) && appErr.Code == "INVALID_INPUT"
}

// IsInternal checks if the error is an internal server error
func IsInternal(err error) bool {
	var appErr *AppError
	return errors.As(err, &appErr) && appErr.Code == "INTERNAL_ERROR"
}

// GetStatusCode returns the HTTP status code for an error
func GetStatusCode(err error) int {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.StatusCode
	}
	return http.StatusInternalServerError
}

// ContainsErrorMessage checks if the error contains a specific message
func ContainsErrorMessage(err error, substring string) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), substring)
}
