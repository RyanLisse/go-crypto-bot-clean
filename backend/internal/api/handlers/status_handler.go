package handlers

import (
	"context"
	"fmt"
	"net/http"
	"runtime"

	"github.com/gin-gonic/gin"
	responseDto "github.com/ryanlisse/go-crypto-bot/internal/api/dto/response"
	"github.com/ryanlisse/go-crypto-bot/internal/core/status"
)

// StatusServiceInterface defines the interface for status service
type StatusServiceInterface interface {
	GetStatus() (*status.SystemStatus, error)
	StartProcesses(ctx context.Context) (*status.SystemStatus, error)
	StopProcesses() (*status.SystemStatus, error)
}

// StatusHandler handles status endpoints
type StatusHandler struct {
	StatusService StatusServiceInterface
}

// NewStatusHandler creates a new StatusHandler
func NewStatusHandler(statusService StatusServiceInterface) *StatusHandler {
	return &StatusHandler{StatusService: statusService}
}

// GetStatus godoc
// @Summary Get system status
// @Description Get system status
// @Tags status
// @Produce json
// @Success 200 {object} responseDto.StatusResponse
// @Failure 500 {object} responseDto.ErrorResponse
// @Router /api/v1/status [get]
func (h *StatusHandler) GetStatus(c *gin.Context) {
	sysStatus, err := h.StatusService.GetStatus()
	if err != nil {
		c.JSON(http.StatusInternalServerError, responseDto.ErrorResponse{
			Code:    "GET_STATUS_FAILED",
			Message: "Failed to get system status",
			Details: err.Error(),
		})
		return
	}

	// Convert to frontend-friendly format
	processes := make([]responseDto.ProcessStatusResponse, 0, len(sysStatus.Components))
	for _, comp := range sysStatus.Components {
		processes = append(processes, responseDto.ProcessStatusResponse{
			Name:      comp.Name,
			Status:    comp.Status,
			IsRunning: comp.IsRunning,
		})
	}

	// Get memory stats
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	response := responseDto.StatusResponse{
		Status:       sysStatus.OverallStatus,
		Version:      sysStatus.SystemInfo.Version,
		Uptime:       sysStatus.SystemInfo.Uptime,
		StartTime:    sysStatus.SystemInfo.StartTime,
		Goroutines:   runtime.NumGoroutine(),
		ProcessCount: len(processes),
		Processes:    processes,
		MemoryUsage: responseDto.MemoryUsageResponse{
			Allocated: fmt.Sprintf("%dMB", m.Alloc/1024/1024),
			Total:     fmt.Sprintf("%dMB", m.TotalAlloc/1024/1024),
			System:    fmt.Sprintf("%dMB", m.Sys/1024/1024),
		},
	}

	c.JSON(http.StatusOK, response)
}

// StartProcesses godoc
// @Summary Start system processes
// @Description Start all system processes
// @Tags status
// @Produce json
// @Success 200 {object} status.SystemStatus
// @Failure 500 {object} responseDto.ErrorResponse
// @Router /api/v1/status/start [post]
func (h *StatusHandler) StartProcesses(c *gin.Context) {
	status, err := h.StatusService.StartProcesses(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, responseDto.ErrorResponse{
			Code:    "START_PROCESSES_FAILED",
			Message: "Failed to start system processes",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, status)
}

// StopProcesses godoc
// @Summary Stop system processes
// @Description Stop all system processes
// @Tags status
// @Produce json
// @Success 200 {object} status.SystemStatus
// @Failure 500 {object} responseDto.ErrorResponse
// @Router /api/v1/status/stop [post]
func (h *StatusHandler) StopProcesses(c *gin.Context) {
	status, err := h.StatusService.StopProcesses()
	if err != nil {
		c.JSON(http.StatusInternalServerError, responseDto.ErrorResponse{
			Code:    "STOP_PROCESSES_FAILED",
			Message: "Failed to stop system processes",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, status)
}
