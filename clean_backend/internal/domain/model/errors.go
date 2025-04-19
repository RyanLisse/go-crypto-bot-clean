package model

import "errors"

// Domain model errors
var (
	ErrInvalidUserID      = errors.New("invalid user ID")
	ErrInvalidExchange    = errors.New("invalid exchange")
	ErrInvalidAPIKey      = errors.New("invalid API key")
	ErrInvalidAPISecret   = errors.New("invalid API secret")
	ErrInvalidWalletID    = errors.New("invalid wallet ID")
	ErrInvalidAsset       = errors.New("invalid asset")
	ErrInvalidAmount      = errors.New("invalid amount")
	ErrInsufficientFunds  = errors.New("insufficient funds")
	ErrWalletNotFound     = errors.New("wallet not found")
	ErrCredentialNotFound = errors.New("credential not found")

	// Network and API errors
	ErrNetworkFailure     = errors.New("network failure")
	ErrTimeout            = errors.New("request timeout")
	ErrConnectionReset    = errors.New("connection reset")
	ErrRateLimitExceeded  = errors.New("rate limit exceeded")
	ErrServerError        = errors.New("server error")
	ErrBadGateway         = errors.New("bad gateway")
	ErrServiceUnavailable = errors.New("service unavailable")
	ErrGatewayTimeout     = errors.New("gateway timeout")
)

// HTTPError represents an HTTP error with status code and headers
type HTTPError struct {
	StatusCode int
	Message    string
	Headers    map[string]string
}

// Error implements the error interface
func (e *HTTPError) Error() string {
	return e.Message
}
