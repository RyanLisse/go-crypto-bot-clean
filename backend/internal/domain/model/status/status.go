package status

import (
	"time"
)

// Status represents the operational status of a component
type Status string

const (
	// StatusRunning indicates the component is running normally
	StatusRunning Status = "running"
	// StatusStopped indicates the component is intentionally stopped
	StatusStopped Status = "stopped"
	// StatusError indicates the component is in an error state
	StatusError Status = "error"
	// StatusWarning indicates the component is running but with warnings
	StatusWarning Status = "warning"
	// StatusStarting indicates the component is in the process of starting
	StatusStarting Status = "starting"
	// StatusStopping indicates the component is in the process of stopping
	StatusStopping Status = "stopping"
	// StatusUnknown indicates the component's status cannot be determined
	StatusUnknown Status = "unknown"
)

// ComponentStatus represents the status of a system component
type ComponentStatus struct {
	// Name is the name of the component
	Name string `json:"name"`
	// Status is the current operational status
	Status Status `json:"status"`
	// Message provides additional details about the status
	Message string `json:"message,omitempty"`
	// StartedAt is when the component was started
	StartedAt *time.Time `json:"started_at,omitempty"`
	// StoppedAt is when the component was stopped
	StoppedAt *time.Time `json:"stopped_at,omitempty"`
	// LastError contains the last error message if status is error
	LastError string `json:"last_error,omitempty"`
	// LastCheckedAt is when the status was last checked
	LastCheckedAt time.Time `json:"last_checked_at"`
	// Metrics contains component-specific metrics
	Metrics map[string]interface{} `json:"metrics,omitempty"`
}

// SystemStatus represents the overall status of the system
type SystemStatus struct {
	// Status is the overall system status
	Status Status `json:"status"`
	// Version is the application version
	Version string `json:"version"`
	// Uptime is how long the system has been running
	Uptime string `json:"uptime"`
	// StartedAt is when the system was started
	StartedAt time.Time `json:"started_at"`
	// Components contains the status of all system components
	Components map[string]*ComponentStatus `json:"components"`
	// SystemInfo contains system resource information
	SystemInfo *SystemInfo `json:"system_info"`
	// LastUpdated is when the status was last updated
	LastUpdated time.Time `json:"last_updated"`
}

// SystemInfo represents system resource information
type SystemInfo struct {
	// CPUUsage is the current CPU usage percentage
	CPUUsage float64 `json:"cpu_usage"`
	// MemoryUsage is the current memory usage percentage
	MemoryUsage float64 `json:"memory_usage"`
	// DiskUsage is the current disk usage percentage
	DiskUsage float64 `json:"disk_usage"`
	// NumGoroutines is the current number of goroutines
	NumGoroutines int `json:"num_goroutines"`
	// AllocatedMemory is the current allocated memory in bytes
	AllocatedMemory uint64 `json:"allocated_memory"`
	// TotalAllocatedMemory is the total allocated memory since startup in bytes
	TotalAllocatedMemory uint64 `json:"total_allocated_memory"`
	// GCPauseTotal is the total time spent in GC pauses in nanoseconds
	GCPauseTotal uint64 `json:"gc_pause_total"`
	// LastGCPause is the last GC pause time in nanoseconds
	LastGCPause uint64 `json:"last_gc_pause"`
}

// ProcessControl represents a command to control a system process
type ProcessControl struct {
	// Action is the control action to perform (start, stop, restart)
	Action string `json:"action"`
	// Component is the name of the component to control
	Component string `json:"component"`
	// Timeout is the maximum time to wait for the action to complete
	Timeout time.Duration `json:"timeout,omitempty"`
}

// ProcessControlResponse represents the response to a process control command
type ProcessControlResponse struct {
	// Success indicates whether the control action was successful
	Success bool `json:"success"`
	// Message provides additional details about the result
	Message string `json:"message,omitempty"`
	// Component is the name of the component that was controlled
	Component string `json:"component"`
	// Action is the control action that was performed
	Action string `json:"action"`
	// NewStatus is the new status of the component after the action
	NewStatus Status `json:"new_status"`
	// CompletedAt is when the action was completed
	CompletedAt time.Time `json:"completed_at"`
}

// NewComponentStatus creates a new component status
func NewComponentStatus(name string, status Status) *ComponentStatus {
	now := time.Now()
	return &ComponentStatus{
		Name:          name,
		Status:        status,
		LastCheckedAt: now,
		Metrics:       make(map[string]interface{}),
	}
}

// NewSystemStatus creates a new system status
func NewSystemStatus(version string, startedAt time.Time) *SystemStatus {
	return &SystemStatus{
		Status:      StatusRunning,
		Version:     version,
		StartedAt:   startedAt,
		Uptime:      time.Since(startedAt).String(),
		Components:  make(map[string]*ComponentStatus),
		SystemInfo:  &SystemInfo{},
		LastUpdated: time.Now(),
	}
}

// UpdateStatus updates the status of a component
func (c *ComponentStatus) UpdateStatus(status Status, message string) {
	now := time.Now()
	c.Status = status
	c.Message = message
	c.LastCheckedAt = now

	if status == StatusRunning && c.StartedAt == nil {
		c.StartedAt = &now
		c.StoppedAt = nil
	} else if status == StatusStopped && c.StoppedAt == nil {
		c.StoppedAt = &now
	}
}

// AddMetric adds or updates a metric for the component
func (c *ComponentStatus) AddMetric(name string, value interface{}) {
	c.Metrics[name] = value
}

// SetError sets the component status to error with the given error message
func (c *ComponentStatus) SetError(err error) {
	if err == nil {
		return
	}
	c.Status = StatusError
	c.LastError = err.Error()
	c.LastCheckedAt = time.Now()
}

// UpdateSystemStatus updates the overall system status based on component statuses
func (s *SystemStatus) UpdateSystemStatus() {
	s.LastUpdated = time.Now()
	s.Uptime = time.Since(s.StartedAt).String()

	// Determine overall status based on component statuses
	overallStatus := StatusRunning
	for _, component := range s.Components {
		switch component.Status {
		case StatusError:
			overallStatus = StatusError
			break
		case StatusWarning:
			if overallStatus != StatusError {
				overallStatus = StatusWarning
			}
		case StatusStopped:
			if overallStatus != StatusError && overallStatus != StatusWarning {
				overallStatus = StatusStopped
			}
		}
	}
	s.Status = overallStatus
}

// AddComponent adds a component to the system status
func (s *SystemStatus) AddComponent(component *ComponentStatus) {
	s.Components[component.Name] = component
	s.UpdateSystemStatus()
}

// GetComponent gets a component by name
func (s *SystemStatus) GetComponent(name string) *ComponentStatus {
	return s.Components[name]
}

// RemoveComponent removes a component from the system status
func (s *SystemStatus) RemoveComponent(name string) {
	delete(s.Components, name)
	s.UpdateSystemStatus()
}
