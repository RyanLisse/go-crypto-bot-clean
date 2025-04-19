package apperror

import (
	"errors"
)

// Common application errors
var (
	// ErrNotFound is returned when a resource is not found
	ErrNotFound = errors.New("resource not found")

	// ErrUnauthorized is returned when a user is not authenticated
	ErrUnauthorized = errors.New("unauthorized")

	// ErrForbidden is returned when a user is not authorized to access a resource
	ErrForbidden = errors.New("forbidden")

	// ErrInvalidInput is returned when input validation fails
	ErrInvalidInput = errors.New("invalid input")

	// ErrInternal is returned when an internal server error occurs
	ErrInternal = errors.New("internal server error")

	// ErrConflict is returned when a resource already exists
	ErrConflict = errors.New("resource conflict")

	// ErrServiceUnavailable is returned when a service is unavailable
	ErrServiceUnavailable = errors.New("service unavailable")

	// ErrTimeout is returned when a request times out
	ErrTimeout = errors.New("request timeout")

	// ErrRateLimited is returned when a request is rate limited
	ErrRateLimited = errors.New("rate limit exceeded")

	// ErrDatabaseError is returned when a database error occurs
	ErrDatabaseError = errors.New("database error")

	// ErrExternalServiceError is returned when an external service error occurs
	ErrExternalServiceError = errors.New("external service error")

	// ErrValidation is returned when validation fails
	ErrValidation = errors.New("validation error")

	// ErrMissingParameter is returned when a required parameter is missing
	ErrMissingParameter = errors.New("missing required parameter")

	// ErrInvalidCredentials is returned when credentials are invalid
	ErrInvalidCredentials = errors.New("invalid credentials")

	// ErrSessionExpired is returned when a session has expired
	ErrSessionExpired = errors.New("session expired")

	// ErrTokenExpired is returned when a token has expired
	ErrTokenExpired = errors.New("token expired")

	// ErrInvalidToken is returned when a token is invalid
	ErrInvalidToken = errors.New("invalid token")

	// ErrUnsupportedMediaType is returned when the media type is not supported
	ErrUnsupportedMediaType = errors.New("unsupported media type")

	// ErrTooManyRequests is returned when too many requests are made
	ErrTooManyRequests = errors.New("too many requests")

	// ErrBadGateway is returned when a bad gateway error occurs
	ErrBadGateway = errors.New("bad gateway")

	// ErrGatewayTimeout is returned when a gateway timeout occurs
	ErrGatewayTimeout = errors.New("gateway timeout")
)
