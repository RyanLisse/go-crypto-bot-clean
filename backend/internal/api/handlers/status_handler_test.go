package handlers

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"go-crypto-bot-clean/backend/internal/core/status"
)

// MockStatusService is a mock implementation of status.StatusService
type MockStatusService struct {
	mock.Mock
}

func (m *MockStatusService) GetStatus() (*status.SystemStatus, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*status.SystemStatus), args.Error(1)
}

// Match signature expected by handler call
func (m *MockStatusService) StartProcesses(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0) // Only return error
}

// Match signature expected by handler call
func (m *MockStatusService) StopProcesses() error {
	args := m.Called()
	return args.Error(0) // Only return error
}

func TestStatusHandler_GetStatus(t *testing.T) {

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

	tests := []struct {
		name           string
		mockStatus     *status.SystemStatus
		mockErr        error
		expectedStatus int
		expectError    bool
	}{
		{
			name:           "success",
			mockStatus:     mockSystemStatus,
			mockErr:        nil,
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "service error",
			mockStatus:     nil,
			mockErr:        errors.New("service error"),
			expectedStatus: http.StatusInternalServerError,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockStatusService)
			mockSvc.On("GetStatus").Return(tt.mockStatus, tt.mockErr)

			handler := NewStatusHandler(mockSvc)

			req := httptest.NewRequest(http.MethodGet, "/api/v1/status", nil)
			w := httptest.NewRecorder()

			handler.GetStatus(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectError {
				assert.Contains(t, w.Body.String(), "GET_STATUS_FAILED")
			} else {
				assert.Contains(t, w.Body.String(), "healthy")
				assert.Contains(t, w.Body.String(), "1.0.0")
			}

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestStatusHandler_StartProcesses(t *testing.T) {

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

	tests := []struct {
		name           string
		mockStatus     *status.SystemStatus
		mockErr        error
		expectedStatus int
		expectError    bool
	}{
		{
			name:           "success",
			mockStatus:     mockSystemStatus,
			mockErr:        nil,
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "service error",
			mockStatus:     nil,
			mockErr:        errors.New("service error"),
			expectedStatus: http.StatusInternalServerError,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockStatusService)
			// Adjust mock setup to only return error
			mockSvc.On("StartProcesses", mock.Anything).Return(tt.mockErr)

			handler := NewStatusHandler(mockSvc)

			req := httptest.NewRequest(http.MethodPost, "/api/v1/status/start", bytes.NewBuffer([]byte{}))
			w := httptest.NewRecorder()

			handler.StartProcesses(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectError {
				assert.Contains(t, w.Body.String(), "START_PROCESSES_FAILED")
			} else {
				assert.Contains(t, w.Body.String(), "healthy")
			}

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestStatusHandler_StopProcesses(t *testing.T) {

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
				IsRunning: false,
				Status:    "stopped",
			},
			{
				Name:      "PositionMonitor",
				IsRunning: false,
				Status:    "stopped",
			},
		},
		OverallStatus: "degraded",
	}

	tests := []struct {
		name           string
		mockStatus     *status.SystemStatus
		mockErr        error
		expectedStatus int
		expectError    bool
	}{
		{
			name:           "success",
			mockStatus:     mockSystemStatus,
			mockErr:        nil,
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "service error",
			mockStatus:     nil,
			mockErr:        errors.New("service error"),
			expectedStatus: http.StatusInternalServerError,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockStatusService)
			// Adjust mock setup to only return error
			mockSvc.On("StopProcesses").Return(tt.mockErr)

			handler := NewStatusHandler(mockSvc)

			req := httptest.NewRequest(http.MethodPost, "/api/v1/status/stop", bytes.NewBuffer([]byte{}))
			w := httptest.NewRecorder()

			handler.StopProcesses(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectError {
				assert.Contains(t, w.Body.String(), "STOP_PROCESSES_FAILED")
			} else {
				assert.Contains(t, w.Body.String(), "degraded")
			}

			mockSvc.AssertExpectations(t)
		})
	}
}
