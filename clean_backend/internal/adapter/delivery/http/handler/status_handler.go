package handler

import (
	"net/http"

	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/adapter/delivery/http/response"
	"github.com/rs/zerolog"
)

// StatusHandler handles status check requests.
type StatusHandler struct {
	logger *zerolog.Logger
}

// NewStatusHandler creates a new StatusHandler.
func NewStatusHandler(logger *zerolog.Logger) *StatusHandler {
	return &StatusHandler{logger: logger}
}

// GetStatus handles the GET /api/v1/status request.
func (h *StatusHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
	h.logger.Info().Msg("API status endpoint called")
	response.WriteSuccessJSON(w, http.StatusOK, map[string]string{
		"status": "operational",
	}, h.logger)
}
