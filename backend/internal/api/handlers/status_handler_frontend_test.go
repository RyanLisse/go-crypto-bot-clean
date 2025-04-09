package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"go-crypto-bot-clean/backend/internal/core/status"
)

// TestStatusHandler_GetStatusForFrontend tests that the status handler returns data in the format expected by the frontend
func TestStatusHandler_GetStatusForFrontend(t *testing.T) {
	gin.SetMode(gin.TestMode)

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
	mockSvc := new(MockStatusService)
	mockSvc.On("GetStatus").Return(mockSystemStatus, nil)

	// Create handler
	handler := NewStatusHandler(mockSvc)

	// Create router
	router := gin.New()
	router.GET("/api/v1/status", handler.GetStatus)

	// Create request
	req := httptest.NewRequest(http.MethodGet, "/api/v1/status", nil)
	w := httptest.NewRecorder()

	// Serve request
	router.ServeHTTP(w, req)

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
	mockSvc.AssertExpectations(t)
}
