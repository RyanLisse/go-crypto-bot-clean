package model

// HealthResponse represents the response from the health check endpoint
type HealthResponse struct {
	// Status of the API
	Status string `json:"status"`

	// Version of the API
	Version string `json:"version"`

	// Timestamp of the response
	Timestamp string `json:"timestamp"`
}
