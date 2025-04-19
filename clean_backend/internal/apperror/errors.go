package apperror

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// ErrorCode represents a standardized error code
type ErrorCode string

// Standard error codes
const (
	ErrCodeBadRequest           ErrorCode = "BAD_REQUEST"
	ErrCodeUnauthorized         ErrorCode = "UNAUTHORIZED"
	ErrCodeForbidden            ErrorCode = "FORBIDDEN"
	ErrCodeNotFound             ErrorCode = "NOT_FOUND"
	ErrCodeMethodNotAllowed     ErrorCode = "METHOD_NOT_ALLOWED"
	ErrCodeConflict             ErrorCode = "CONFLICT"
	ErrCodeValidationError      ErrorCode = "VALIDATION_ERROR"
	ErrCodeRateLimitExceeded    ErrorCode = "RATE_LIMIT_EXCEEDED"
	ErrCodeInternalError        ErrorCode = "INTERNAL_ERROR"
	ErrCodeServiceUnavailable   ErrorCode = "SERVICE_UNAVAILABLE"
	ErrCodeBadGateway           ErrorCode = "BAD_GATEWAY"
	ErrCodeGatewayTimeout       ErrorCode = "GATEWAY_TIMEOUT"
	ErrCodeUnknownError         ErrorCode = "UNKNOWN_ERROR"
	ErrCodeDatabaseError        ErrorCode = "DATABASE_ERROR"
	ErrCodeExternalServiceError ErrorCode = "EXTERNAL_SERVICE_ERROR"
)

// ErrorResponse represents the standardized error response format
type ErrorResponse struct {
	Success bool        `json:"success"`
	Error   ErrorDetail `json:"error"`
}

// ErrorDetail contains the details of an error
type ErrorDetail struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	TraceID string      `json:"trace_id,omitempty"`
	Details interface{} `json:"details,omitempty"`
}

// AppError represents an application error
type AppError struct {
	Code       int         `json:"-"`
	ErrorCode  ErrorCode   `json:"-"`
	Message    string      `json:"message"`
	Details    interface{} `json:"details,omitempty"`
	Err        error       `json:"-"`
	StatusCode int         `json:"-"`
}

// Error returns the error message
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap returns the wrapped error
func (e *AppError) Unwrap() error {
	return e.Err
}

// GetStatusCode returns the HTTP status code for the error
func (e *AppError) GetStatusCode() int {
	if e.StatusCode != 0 {
		return e.StatusCode
	}
	return e.Code // Fallback to Code for backward compatibility
}

// GetErrorCode returns the error code as a string
func (e *AppError) GetErrorCode() string {
	if e.ErrorCode != "" {
		return string(e.ErrorCode)
	}
	return getErrorCode(e.Code) // Fallback to HTTP status code mapping
}

// NewNotFound creates a new not found error
func NewNotFound(resource string, identifier any, err error) *AppError {
	var details any
	if identifier != nil {
		details = map[string]any{"identifier": identifier}
	}

	return &AppError{
		Code:       http.StatusNotFound,
		ErrorCode:  ErrCodeNotFound,
		Message:    fmt.Sprintf("%s with identifier %v not found", resource, identifier),
		Details:    details,
		Err:        err,
		StatusCode: http.StatusNotFound,
	}
}

// NewBadRequest creates a new bad request error
func NewBadRequest(message string, details any, err error) *AppError {
	return &AppError{
		Code:       http.StatusBadRequest,
		ErrorCode:  ErrCodeBadRequest,
		Message:    message,
		Details:    details,
		Err:        err,
		StatusCode: http.StatusBadRequest,
	}
}

// NewUnauthorized creates a new unauthorized error
func NewUnauthorized(message string, err error) *AppError {
	return &AppError{
		Code:       http.StatusUnauthorized,
		ErrorCode:  ErrCodeUnauthorized,
		Message:    message,
		Err:        err,
		StatusCode: http.StatusUnauthorized,
	}
}

// NewForbidden creates a new forbidden error
func NewForbidden(message string, err error) *AppError {
	return &AppError{
		Code:       http.StatusForbidden,
		ErrorCode:  ErrCodeForbidden,
		Message:    message,
		Err:        err,
		StatusCode: http.StatusForbidden,
	}
}

// NewInternal creates a new internal server error
func NewInternal(err error) *AppError {
	return &AppError{
		Code:       http.StatusInternalServerError,
		ErrorCode:  ErrCodeInternalError,
		Message:    "Internal server error",
		Err:        err,
		StatusCode: http.StatusInternalServerError,
	}
}

// NewValidationError creates a new validation error
func NewValidationError(message string, details any, err error) *AppError {
	return &AppError{
		Code:       http.StatusUnprocessableEntity,
		ErrorCode:  ErrCodeValidationError,
		Message:    message,
		Details:    details,
		Err:        err,
		StatusCode: http.StatusUnprocessableEntity,
	}
}

// NewRateLimit creates a new rate limit error
func NewRateLimit(message string, err error) *AppError {
	return &AppError{
		Code:       http.StatusTooManyRequests,
		ErrorCode:  ErrCodeRateLimitExceeded,
		Message:    message,
		Err:        err,
		StatusCode: http.StatusTooManyRequests,
	}
}

// NewConflict creates a new conflict error
func NewConflict(message string, details any, err error) *AppError {
	return &AppError{
		Code:       http.StatusConflict,
		ErrorCode:  ErrCodeConflict,
		Message:    message,
		Details:    details,
		Err:        err,
		StatusCode: http.StatusConflict,
	}
}

// NewServiceUnavailable creates a new service unavailable error
func NewServiceUnavailable(message string, err error) *AppError {
	return &AppError{
		Code:       http.StatusServiceUnavailable,
		ErrorCode:  ErrCodeServiceUnavailable,
		Message:    message,
		Err:        err,
		StatusCode: http.StatusServiceUnavailable,
	}
}

// NewBadGateway creates a new bad gateway error
func NewBadGateway(message string, err error) *AppError {
	return &AppError{
		Code:       http.StatusBadGateway,
		ErrorCode:  ErrCodeBadGateway,
		Message:    message,
		Err:        err,
		StatusCode: http.StatusBadGateway,
	}
}

// NewGatewayTimeout creates a new gateway timeout error
func NewGatewayTimeout(message string, err error) *AppError {
	return &AppError{
		Code:       http.StatusGatewayTimeout,
		ErrorCode:  ErrCodeGatewayTimeout,
		Message:    message,
		Err:        err,
		StatusCode: http.StatusGatewayTimeout,
	}
}

// NewDatabaseError creates a new database error
func NewDatabaseError(message string, err error) *AppError {
	return &AppError{
		Code:       http.StatusInternalServerError,
		ErrorCode:  ErrCodeDatabaseError,
		Message:    message,
		Err:        err,
		StatusCode: http.StatusInternalServerError,
	}
}

// NewExternalServiceError creates a new external service error
func NewExternalServiceError(message string, err error) *AppError {
	return &AppError{
		Code:       http.StatusInternalServerError,
		ErrorCode:  ErrCodeExternalServiceError,
		Message:    message,
		Err:        err,
		StatusCode: http.StatusInternalServerError,
	}
}

// WriteError writes an error response
func WriteError(w http.ResponseWriter, err *AppError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.GetStatusCode())

	// Create error response
	response := ErrorResponse{
		Success: false,
		Error: ErrorDetail{
			Code:    err.GetErrorCode(),
			Message: err.Message,
		},
	}

	// Add details if present
	if err.Details != nil {
		response.Error.Details = err.Details
	}

	// Write response
	if encodeErr := json.NewEncoder(w).Encode(response); encodeErr != nil {
		// If encoding fails, write a simple error message
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"success":false,"error":{"message":"Failed to encode error response"}}`))
	}
}

// RespondWithError writes an error response with request context
func RespondWithError(w http.ResponseWriter, r *http.Request, err *AppError) {
	// Get request ID from context if available
	requestID := r.Header.Get("X-Request-ID")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.GetStatusCode())

	// Create error response
	response := ErrorResponse{
		Success: false,
		Error: ErrorDetail{
			Code:    err.GetErrorCode(),
			Message: err.Message,
			TraceID: requestID,
		},
	}

	// Add details if present
	if err.Details != nil {
		response.Error.Details = err.Details
	}

	// Write response
	if encodeErr := json.NewEncoder(w).Encode(response); encodeErr != nil {
		// If encoding fails, write a simple error message
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"success":false,"error":{"message":"Failed to encode error response"}}`))
	}
}

// FromError converts a standard error to an AppError
func FromError(err error) *AppError {
	// If it's already an AppError, return it
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}

	// Check for common error types
	switch {
	case errors.Is(err, ErrNotFound):
		return NewNotFound("Resource", nil, err)
	case errors.Is(err, ErrUnauthorized):
		return NewUnauthorized("Unauthorized", err)
	case errors.Is(err, ErrForbidden):
		return NewForbidden("Forbidden", err)
	case errors.Is(err, ErrInvalidInput):
		return NewBadRequest("Invalid input", nil, err)
	case errors.Is(err, ErrValidation):
		return NewValidationError("Validation error", nil, err)
	case errors.Is(err, ErrConflict):
		return NewConflict("Resource conflict", nil, err)
	case errors.Is(err, ErrDatabaseError):
		return NewDatabaseError("Database error", err)
	case errors.Is(err, ErrExternalServiceError):
		return NewExternalServiceError("External service error", err)
	case errors.Is(err, ErrRateLimited):
		return NewRateLimit("Rate limit exceeded", err)
	case errors.Is(err, ErrServiceUnavailable):
		return NewServiceUnavailable("Service unavailable", err)
	case errors.Is(err, ErrBadGateway):
		return NewBadGateway("Bad gateway", err)
	case errors.Is(err, ErrGatewayTimeout):
		return NewGatewayTimeout("Gateway timeout", err)
	default:
		return NewInternal(err)
	}
}

// getErrorCode returns a string error code based on HTTP status code
func getErrorCode(statusCode int) string {
	switch statusCode {
	case http.StatusBadRequest:
		return "BAD_REQUEST"
	case http.StatusUnauthorized:
		return "UNAUTHORIZED"
	case http.StatusForbidden:
		return "FORBIDDEN"
	case http.StatusNotFound:
		return "NOT_FOUND"
	case http.StatusMethodNotAllowed:
		return "METHOD_NOT_ALLOWED"
	case http.StatusConflict:
		return "CONFLICT"
	case http.StatusUnprocessableEntity:
		return "VALIDATION_ERROR"
	case http.StatusTooManyRequests:
		return "RATE_LIMIT_EXCEEDED"
	case http.StatusInternalServerError:
		return "INTERNAL_ERROR"
	case http.StatusServiceUnavailable:
		return "SERVICE_UNAVAILABLE"
	default:
		return "UNKNOWN_ERROR"
	}
}
