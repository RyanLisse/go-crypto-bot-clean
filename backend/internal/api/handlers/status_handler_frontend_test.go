package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go-crypto-bot-clean/backend/internal/core/status"
)

type dummyStatusService struct {
	status *status.SystemStatus
	err    error
}

func (d *dummyStatusService) GetStatus() (*status.SystemStatus, error) {
	return d.status, d.err
}

// Add dummy methods to satisfy status.Service interface
func (d *dummyStatusService) StartProcesses(ctx context.Context) error {
	return nil // No-op for this test
}

func (d *dummyStatusService) StopProcesses() error {
	return nil // No-op for this test
}

// TestStatusHandler_GetStatusForFrontend tests that the status handler returns data in the format expected by the frontend
func TestStatusHandler_GetStatusForFrontend(t *testing.T) {

	// Create mock system status
	mockSystemStatus := &status.SystemStatus{
		SystemInfo: status.SystemInfo{
			Version:   "1.0.0",
			GoVersion: "go1.16",
			StartTime: time.Now().Add(-1 * time.Hour),
			Uptime:    "1h 0m 0s",
		},
		Components: []status.ComponentStatus{
			{
				Name:      "NewCoinWatcher",
				IsRunning: true,
				Status:    "running",
			},
			{
				Name:      "PositionMonitor",
				IsRunning: true,
				Status:    "running",
			},
		},
		OverallStatus: "healthy",
	}

	// Create mock service
	mockSvc := &dummyStatusService{status: mockSystemStatus, err: nil}

	// Create handler
	handler := NewStatusHandler(mockSvc)

	// Create router

	// Create request
	req := httptest.NewRequest(http.MethodGet, "/api/v1/status", nil)
	w := httptest.NewRecorder()

	// Serve request
	handler.GetStatus(w, req)

	// Check response
	assert.Equal(t, http.StatusOK, w.Code)

	// Parse response
	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Check that the response has the fields expected by the frontend
	assert.Contains(t, response, "status")
	assert.Contains(t, response, "version")
	assert.Contains(t, response, "uptime")
	assert.Contains(t, response, "memory_usage")

	// Verify mock expectations
}
