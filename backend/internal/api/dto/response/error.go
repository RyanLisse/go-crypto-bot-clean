// Package response contains API response DTOs.
package response

import "time"

// ErrorResponse represents a standardized error response format
type ErrorResponse struct {
	Code      string      `json:"code"`                // Error code for client handling
	Message   string      `json:"message"`             // User-friendly error message
	Details   interface{} `json:"details,omitempty"`   // Optional detailed error information
	Help      string      `json:"help,omitempty"`      // Optional help text or troubleshooting guidance
	RequestID string      `json:"requestId,omitempty"` // Request ID for tracing
	Path      string      `json:"path"`                // Request path
	Method    string      `json:"method"`              // HTTP method
	Timestamp time.Time   `json:"timestamp"`           // Error timestamp
	Latency   string      `json:"latency,omitempty"`   // Request processing time
}

// SuccessResponse represents a standardized API success response.
// swagger:model SuccessResponse
type SuccessResponse struct {
	Message   string    `json:"message"`   // Human-readable success message
	Timestamp time.Time `json:"timestamp"` // Timestamp of the response
}
