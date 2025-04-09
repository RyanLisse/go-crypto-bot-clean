// Package response contains API response DTOs.
package response

import "time"

// ErrorResponse represents a standardized API error response.
// swagger:model ErrorResponse
type ErrorResponse struct {
	Code    string `json:"code"`              // Error code identifier
	Message string `json:"message"`           // Human-readable error message
	Details string `json:"details,omitempty"` // Optional additional details
}

// SuccessResponse represents a standardized API success response.
// swagger:model SuccessResponse
type SuccessResponse struct {
	Message   string    `json:"message"`   // Human-readable success message
	Timestamp time.Time `json:"timestamp"` // Timestamp of the response
}
