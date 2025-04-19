package health

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"

	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/adapter/delivery/http/response"
)

// Handler handles health check requests
type Handler struct {
	logger *zerolog.Logger
}

// NewHandler creates a new health check handler
func NewHandler(logger *zerolog.Logger) *Handler {
	return &Handler{
		logger: logger,
	}
}

// RegisterRoutes registers the health check routes
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Get("/health", h.HealthCheck)
	r.Get("/readiness", h.ReadinessCheck)
	r.Get("/liveness", h.LivenessCheck)
}

// HealthCheck handles the health check endpoint
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug().Msg("Health check requested")
	
	response.WriteSuccessJSON(w, http.StatusOK, map[string]interface{}{
		"status": "ok",
		"time":   time.Now().Format(time.RFC3339),
	}, h.logger)
}

// ReadinessCheck handles the readiness check endpoint
func (h *Handler) ReadinessCheck(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug().Msg("Readiness check requested")
	
	response.WriteSuccessJSON(w, http.StatusOK, map[string]interface{}{
		"status": "ready",
		"time":   time.Now().Format(time.RFC3339),
	}, h.logger)
}

// LivenessCheck handles the liveness check endpoint
func (h *Handler) LivenessCheck(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug().Msg("Liveness check requested")
	
	response.WriteSuccessJSON(w, http.StatusOK, map[string]interface{}{
		"status": "alive",
		"time":   time.Now().Format(time.RFC3339),
	}, h.logger)
}
