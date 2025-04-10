package handlers

import (
	"net/http"

	"go-crypto-bot-clean/backend/internal/core/status"
)

// StatusHandler handles status-related HTTP requests.
type StatusHandler struct {
	statusService status.Service // Use the new interface type
}

// NewStatusHandler creates a new status handler.
func NewStatusHandler(statusService status.Service) *StatusHandler { // Use the new interface type
	return &StatusHandler{
		statusService: statusService,
	}
}

// GetStatus returns the current system status.
func (h *StatusHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
	status, err := h.statusService.GetStatus()
	if err != nil {
		SendError(w, http.StatusInternalServerError, "Failed to get status: "+err.Error())
		return
	}

	SendSuccess(w, status)
}

// StartProcesses starts all system processes.
func (h *StatusHandler) StartProcesses(w http.ResponseWriter, r *http.Request) {
	// Pass context and handle both return values
	err := h.statusService.StartProcesses(r.Context()) // Now returns only error
	if err != nil {
		SendError(w, http.StatusInternalServerError, "Failed to start processes: "+err.Error())
		return
	}
	// We ignore the returned status here, just report success/failure

	SendSuccess(w, map[string]string{"message": "Processes started successfully"})
}

// StopProcesses stops all system processes.
func (h *StatusHandler) StopProcesses(w http.ResponseWriter, r *http.Request) {
	// Handle both return values
	err := h.statusService.StopProcesses() // Now returns only error
	if err != nil {
		SendError(w, http.StatusInternalServerError, "Failed to stop processes: "+err.Error())
		return
	}
	// We ignore the returned status here, just report success/failure

	SendSuccess(w, map[string]string{"message": "Processes stopped successfully"})
}
