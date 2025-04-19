package handler

import (
	"net/http"

	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/adapter/delivery/http/response"
	"github.com/rs/zerolog"
)

// RootTestHandler handles the root test endpoint.
type RootTestHandler struct {
	logger *zerolog.Logger
}

// NewRootTestHandler creates a new RootTestHandler.
func NewRootTestHandler(logger *zerolog.Logger) *RootTestHandler {
	return &RootTestHandler{logger: logger}
}

// HandleRootTest handles the GET /root-test request.
func (h *RootTestHandler) HandleRootTest(w http.ResponseWriter, r *http.Request) {
	h.logger.Info().Msg("Root level test endpoint called")
	// Using the response helpers from clean_backend
	response.WriteSuccessJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Root level test endpoint works!",
	}, h.logger)
}
