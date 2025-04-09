package response

import "time"

// StatusResponse represents the system status
type StatusResponse struct {
	Status       string                  `json:"status"`
	Version      string                  `json:"version"`
	Uptime       string                  `json:"uptime"`
	StartTime    time.Time               `json:"start_time"`
	MemoryUsage  MemoryUsageResponse     `json:"memory_usage"`
	Goroutines   int                     `json:"goroutines"`
	ProcessCount int                     `json:"process_count"`
	Processes    []ProcessStatusResponse `json:"processes"`
}

// MemoryUsageResponse represents memory usage information
type MemoryUsageResponse struct {
	Allocated string `json:"allocated"`
	Total     string `json:"total"`
	System    string `json:"system"`
}

// ProcessStatusResponse represents the status of a system process
type ProcessStatusResponse struct {
	Name      string `json:"name"`
	Status    string `json:"status"`
	IsRunning bool   `json:"is_running"`
	Uptime    string `json:"uptime,omitempty"`
}

// ErrorResponse is defined in error.go

// ProcessControlResponse represents the result of a process control operation
type ProcessControlResponse struct {
	Process   string    `json:"process"`
	Action    string    `json:"action"`
	Success   bool      `json:"success"`
	Message   string    `json:"message,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}
