// Package health provides health check functionality for the application
package health

import (
	"encoding/json"
	"net/http"
	"runtime"
	"sync"
	"time"

	"go.uber.org/zap"
)

// Status represents the status of a component
type Status string

const (
	// StatusUp indicates the component is functioning properly
	StatusUp Status = "UP"
	// StatusDown indicates the component is not functioning properly
	StatusDown Status = "DOWN"
	// StatusDegraded indicates the component is functioning but with issues
	StatusDegraded Status = "DEGRADED"
)

// Component represents a component of the application that can be health checked
type Component struct {
	Name   string `json:"name"`
	Status Status `json:"status"`
	Detail string `json:"detail,omitempty"`
}

// HealthCheck represents the overall health of the application
type HealthCheck struct {
	Status       Status               `json:"status"`
	Components   map[string]Component `json:"components"`
	Timestamp    time.Time            `json:"timestamp"`
	Version      string               `json:"version"`
	Uptime       string               `json:"uptime"`
	GoVersion    string               `json:"goVersion"`
	GOOS         string               `json:"os"`
	GOARCH       string               `json:"arch"`
	NumGoroutine int                  `json:"numGoroutine"`
	mutex        sync.RWMutex
	startTime    time.Time
	logger       *zap.Logger
}

// NewHealthCheck creates a new HealthCheck
func NewHealthCheck(version string, logger *zap.Logger) *HealthCheck {
	return &HealthCheck{
		Status:       StatusUp,
		Components:   make(map[string]Component),
		Timestamp:    time.Now(),
		Version:      version,
		startTime:    time.Now(),
		GoVersion:    runtime.Version(),
		GOOS:         runtime.GOOS,
		GOARCH:       runtime.GOARCH,
		NumGoroutine: runtime.NumGoroutine(),
		logger:       logger,
	}
}

// AddComponent adds a component to the health check
func (h *HealthCheck) AddComponent(name string, status Status, detail string) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	h.Components[name] = Component{
		Name:   name,
		Status: status,
		Detail: detail,
	}

	h.updateOverallStatus()
}

// RemoveComponent removes a component from the health check
func (h *HealthCheck) RemoveComponent(name string) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	delete(h.Components, name)
	h.updateOverallStatus()
}

// updateOverallStatus updates the overall status based on component statuses
func (h *HealthCheck) updateOverallStatus() {
	h.Status = StatusUp

	for _, component := range h.Components {
		if component.Status == StatusDown {
			h.Status = StatusDown
			return
		}
		if component.Status == StatusDegraded {
			h.Status = StatusDegraded
		}
	}
}

// Check performs a health check and returns the result
func (h *HealthCheck) Check() *HealthCheck {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	result := &HealthCheck{
		Status:       h.Status,
		Components:   make(map[string]Component),
		Timestamp:    time.Now(),
		Version:      h.Version,
		Uptime:       time.Since(h.startTime).String(),
		GoVersion:    h.GoVersion,
		GOOS:         h.GOOS,
		GOARCH:       h.GOARCH,
		NumGoroutine: runtime.NumGoroutine(),
	}

	for name, component := range h.Components {
		result.Components[name] = component
	}

	return result
}

// Handler returns an HTTP handler for the health check
func (h *HealthCheck) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result := h.Check()

		// Set status code based on health status
		if result.Status == StatusDown {
			w.WriteHeader(http.StatusServiceUnavailable)
		} else if result.Status == StatusDegraded {
			w.WriteHeader(http.StatusOK) // Still return 200 but with degraded status
		} else {
			w.WriteHeader(http.StatusOK)
		}

		// Check if client explicitly wants plain text
		accept := r.Header.Get("Accept")
		if accept == "text/plain" {
			w.Header().Set("Content-Type", "text/plain")
			w.Write([]byte(string(result.Status)))
			return
		}

		// Default to JSON for all other cases
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(result); err != nil {
			h.logger.Error("Failed to encode health check response", zap.Error(err))
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
}

// SimpleHandler returns a simple HTTP handler that just returns "OK"
func (h *HealthCheck) SimpleHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
}
