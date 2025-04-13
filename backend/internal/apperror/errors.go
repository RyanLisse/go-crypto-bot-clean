package apperror

import (
	"errors"
	"fmt"
	"net/http"
)

// AppError represents an application-specific error
type AppError struct {
	StatusCode int         `json:"-"`
	Code       string      `json:"code"`
	Message    string      `json:"message"`
	Details    interface{} `json:"details,omitempty"`
	Err        error       `json:"-"`
}

// Error returns the error message
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap returns the wrapped error
func (e *AppError) Unwrap() error {
	return e.Err
}

// ToResponse returns a map suitable for JSON response
func (e *AppError) ToResponse() map[string]interface{} {
	resp := map[string]interface{}{
		"error": map[string]interface{}{
			"code":    e.Code,
			"message": e.Message,
		},
	}

	if e.Details != nil {
		resp["error"].(map[string]interface{})["details"] = e.Details
	}

	return resp
}

// Is checks if the target error is an AppError with the same code
func (e *AppError) Is(target error) bool {
	var appErr *AppError
	if !errors.As(target, &appErr) {
		return false
	}
	return appErr.Code == e.Code
}

// Common error types
var (
	ErrInvalidInput    = &AppError{StatusCode: http.StatusBadRequest, Code: "INVALID_INPUT", Message: "Invalid input provided"}
	ErrNotFound        = &AppError{StatusCode: http.StatusNotFound, Code: "NOT_FOUND", Message: "Resource not found"}
	ErrInternal        = &AppError{StatusCode: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Internal server error"}
	ErrUnauthorized    = &AppError{StatusCode: http.StatusUnauthorized, Code: "UNAUTHORIZED", Message: "Unauthorized"}
	ErrForbidden       = &AppError{StatusCode: http.StatusForbidden, Code: "FORBIDDEN", Message: "Forbidden"}
	ErrConflict        = &AppError{StatusCode: http.StatusConflict, Code: "CONFLICT", Message: "Resource conflict"}
	ErrRateLimit       = &AppError{StatusCode: http.StatusTooManyRequests, Code: "RATE_LIMIT", Message: "Rate limit exceeded"}
	ErrExternalService = &AppError{StatusCode: http.StatusServiceUnavailable, Code: "EXTERNAL_SERVICE_ERROR", Message: "External service error"}
	ErrValidation      = &AppError{StatusCode: http.StatusBadRequest, Code: "VALIDATION_ERROR", Message: "Validation error"}
)

// Common error creators

// NewInvalid creates a new invalid input error
func NewInvalid(msg string, details interface{}, err error) *AppError {
	return &AppError{
		StatusCode: http.StatusBadRequest,
		Code:       "INVALID_INPUT",
		Message:    msg,
		Details:    details,
		Err:        err,
	}
}

// NewNotFound creates a new not found error
func NewNotFound(resource string, identifier interface{}, err error) *AppError {
	var msg string
	if identifier != nil {
		msg = fmt.Sprintf("%s with identifier %v not found", resource, identifier)
	} else {
		msg = fmt.Sprintf("%s not found", resource)
	}

	return &AppError{
		StatusCode: http.StatusNotFound,
		Code:       "NOT_FOUND",
		Message:    msg,
		Err:        err,
	}
}

// NewInternal creates a new internal server error
func NewInternal(err error) *AppError {
	return &AppError{
		StatusCode: http.StatusInternalServerError,
		Code:       "INTERNAL_ERROR",
		Message:    "Internal server error",
		Err:        err,
	}
}

// NewUnauthorized creates a new unauthorized error
func NewUnauthorized(msg string, err error) *AppError {
	if msg == "" {
		msg = "Unauthorized"
	}

	return &AppError{
		StatusCode: http.StatusUnauthorized,
		Code:       "UNAUTHORIZED",
		Message:    msg,
		Err:        err,
	}
}

// NewForbidden creates a new forbidden error
func NewForbidden(msg string, err error) *AppError {
	if msg == "" {
		msg = "Forbidden"
	}

	return &AppError{
		StatusCode: http.StatusForbidden,
		Code:       "FORBIDDEN",
		Message:    msg,
		Err:        err,
	}
}

// NewValidation creates a new validation error
func NewValidation(msg string, details interface{}, err error) *AppError {
	if msg == "" {
		msg = "Validation error"
	}

	return &AppError{
		StatusCode: http.StatusBadRequest,
		Code:       "VALIDATION_ERROR",
		Message:    msg,
		Details:    details,
		Err:        err,
	}
}

// NewExternalService creates a new external service error
func NewExternalService(service string, msg string, err error) *AppError {
	if msg == "" {
		msg = fmt.Sprintf("Error communicating with %s", service)
	}

	return &AppError{
		StatusCode: http.StatusServiceUnavailable,
		Code:       "EXTERNAL_SERVICE_ERROR",
		Message:    msg,
		Err:        err,
	}
}

// As is a wrapper for errors.As
func As(err error, target interface{}) bool {
	return errors.As(err, target)
}

// Is is a wrapper for errors.Is
func Is(err, target error) bool {
	return errors.Is(err, target)
}
