package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// ErrorType represents the type of authentication error
type ErrorType string

const (
	// Error types
	ErrorTypeUnauthorized     ErrorType = "unauthorized"
	ErrorTypeForbidden        ErrorType = "forbidden"
	ErrorTypeInvalidToken     ErrorType = "invalid_token"
	ErrorTypeExpiredToken     ErrorType = "expired_token"
	ErrorTypeInvalidRole      ErrorType = "invalid_role"
	ErrorTypeInvalidSession   ErrorType = "invalid_session"
	ErrorTypeUserNotFound     ErrorType = "user_not_found"
	ErrorTypeInternalError    ErrorType = "internal_error"
	ErrorTypeInvalidMetadata  ErrorType = "invalid_metadata"
	ErrorTypePermissionDenied ErrorType = "permission_denied"

	// Additional error types
	ErrorTypeInvalidRequest     ErrorType = "invalid_request"
	ErrorTypeRateLimitExceeded  ErrorType = "rate_limit_exceeded"
	ErrorTypeServiceUnavailable ErrorType = "service_unavailable"
	ErrorTypeTokenRevoked       ErrorType = "token_revoked"
	ErrorTypeInvalidScope       ErrorType = "invalid_scope"
	ErrorTypeAccountLocked      ErrorType = "account_locked"
	ErrorTypeAccountDisabled    ErrorType = "account_disabled"
)

// ErrorResponse represents the standardized error response structure
type ErrorResponse struct {
	Error     *AuthError `json:"error"`
	RequestID string     `json:"request_id,omitempty"`
	Timestamp time.Time  `json:"timestamp"`
	Path      string     `json:"path,omitempty"`
	Method    string     `json:"method,omitempty"`
}

// AuthError represents a structured authentication error
type AuthError struct {
	Type        ErrorType `json:"type"`
	Message     string    `json:"message"`
	Code        int       `json:"code"`
	RequestID   string    `json:"request_id,omitempty"`
	Details     any       `json:"details,omitempty"`
	InternalErr error     `json:"-"` // Not exposed in JSON
	Timestamp   time.Time `json:"timestamp"`
	Retryable   bool      `json:"retryable,omitempty"`
	Help        string    `json:"help,omitempty"`
}

// Error implements the error interface
func (e *AuthError) Error() string {
	if e.InternalErr != nil {
		return fmt.Sprintf("%s: %s (internal: %v)", e.Type, e.Message, e.InternalErr)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// WriteJSON writes the error as a JSON response with proper headers
func (e *AuthError) WriteJSON(w http.ResponseWriter) {
	response := ErrorResponse{
		Error:     e,
		Timestamp: time.Now().UTC(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(e.Code)
	json.NewEncoder(w).Encode(response)
}

// NewAuthError creates a new AuthError with the given type and message
func NewAuthError(errType ErrorType, message string, code int) *AuthError {
	return &AuthError{
		Type:      errType,
		Message:   message,
		Code:      code,
		Timestamp: time.Now().UTC(),
		Retryable: isRetryable(errType),
	}
}

// WithRequestID adds a request ID to the error
func (e *AuthError) WithRequestID(requestID string) *AuthError {
	e.RequestID = requestID
	return e
}

// WithDetails adds additional details to the error
func (e *AuthError) WithDetails(details any) *AuthError {
	e.Details = details
	return e
}

// WithInternalError adds an internal error
func (e *AuthError) WithInternalError(err error) *AuthError {
	e.InternalErr = err
	return e
}

// WithHelp adds a help message to the error
func (e *AuthError) WithHelp(help string) *AuthError {
	e.Help = help
	return e
}

// isRetryable determines if an error type is retryable
func isRetryable(errType ErrorType) bool {
	switch errType {
	case ErrorTypeServiceUnavailable, ErrorTypeRateLimitExceeded:
		return true
	default:
		return false
	}
}

// Common authentication errors with helpful messages
var (
	ErrUnauthorized = NewAuthError(
		ErrorTypeUnauthorized,
		"Authentication required",
		http.StatusUnauthorized,
	).WithHelp("Please provide valid authentication credentials")

	ErrInvalidToken = NewAuthError(
		ErrorTypeInvalidToken,
		"Invalid or malformed token",
		http.StatusUnauthorized,
	).WithHelp("Please check your authentication token format")

	ErrExpiredToken = NewAuthError(
		ErrorTypeExpiredToken,
		"Token has expired",
		http.StatusUnauthorized,
	).WithHelp("Please obtain a new authentication token")

	ErrInvalidSession = NewAuthError(
		ErrorTypeInvalidSession,
		"Invalid or expired session",
		http.StatusUnauthorized,
	).WithHelp("Please log in again to create a new session")

	ErrUserNotFound = NewAuthError(
		ErrorTypeUserNotFound,
		"User not found",
		http.StatusNotFound,
	).WithHelp("Please verify the user exists and try again")

	ErrInvalidRole = NewAuthError(
		ErrorTypeInvalidRole,
		"Invalid or insufficient role",
		http.StatusForbidden,
	).WithHelp("Please ensure you have the required role for this action")

	ErrPermissionDenied = NewAuthError(
		ErrorTypePermissionDenied,
		"Permission denied",
		http.StatusForbidden,
	).WithHelp("You do not have permission to perform this action")

	ErrInternalServer = NewAuthError(
		ErrorTypeInternalError,
		"Internal server error",
		http.StatusInternalServerError,
	).WithHelp("An unexpected error occurred. Please try again later")

	ErrInvalidMetadata = NewAuthError(
		ErrorTypeInvalidMetadata,
		"Invalid user metadata",
		http.StatusBadRequest,
	).WithHelp("Please check the format of your user metadata")

	ErrRateLimitExceeded = NewAuthError(
		ErrorTypeRateLimitExceeded,
		"Rate limit exceeded",
		http.StatusTooManyRequests,
	).WithHelp("Please wait before making more requests")

	ErrServiceUnavailable = NewAuthError(
		ErrorTypeServiceUnavailable,
		"Authentication service unavailable",
		http.StatusServiceUnavailable,
	).WithHelp("The service is temporarily unavailable. Please try again later")
)
