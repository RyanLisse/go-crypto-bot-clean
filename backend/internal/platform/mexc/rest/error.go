package rest

import (
	"fmt"
	"strings"
)

// APIError represents an error returned from the MEXC API
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
}

// Error returns a string representation of the error
func (e *APIError) Error() string {
	return fmt.Sprintf("MEXC API error (code %d): %s", e.Code, e.Message)
}

// IsRateLimited checks if the error is a rate limit error
func (e *APIError) IsRateLimited() bool {
	return e.Code == -1003 || // Rate limit code
		strings.Contains(strings.ToLower(e.Message), "rate limit")
}

// RequestError represents an error that occurred during an HTTP request
type RequestError struct {
	Err     error
	Message string
}

// Error returns a string representation of the error
func (e *RequestError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("MEXC request error: %s: %v", e.Message, e.Err)
	}
	return fmt.Sprintf("MEXC request error: %s", e.Message)
}

// UnmarshalError represents an error that occurred while unmarshaling a response
type UnmarshalError struct {
	Err     error
	Body    []byte
	Message string
}

// Error returns a string representation of the error
func (e *UnmarshalError) Error() string {
	return fmt.Sprintf("MEXC unmarshal error: %s: %v (body: %s)", e.Message, e.Err, string(e.Body))
}
