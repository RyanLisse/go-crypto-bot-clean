package handler

import (
	"encoding/json"
	"net/http"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/http/response"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/usecase"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
)

// SniperHandler handles HTTP requests for sniper operations
type SniperHandler struct {
	sniperUC usecase.SniperUseCase
	logger   *zerolog.Logger
}

// NewSniperHandler creates a new sniper handler
func NewSniperHandler(sniperUC usecase.SniperUseCase, logger *zerolog.Logger) *SniperHandler {
	return &SniperHandler{
		sniperUC: sniperUC,
		logger:   logger,
	}
}

// RegisterRoutes registers the sniper routes
func (h *SniperHandler) RegisterRoutes(router chi.Router) {
	router.Post("/api/v1/sniper/execute/{symbol}", h.ExecuteSnipe)
	router.Get("/api/v1/sniper/config", h.GetSniperConfig)
	router.Put("/api/v1/sniper/config", h.UpdateSniperConfig)
	router.Get("/api/v1/sniper/status", h.GetSniperStatus)
	router.Post("/api/v1/sniper/start", h.StartSniper)
	router.Post("/api/v1/sniper/stop", h.StopSniper)
	router.Post("/api/v1/sniper/auto", h.ConfigureAutoSnipe)
}

// ExecuteSnipe handles the request to execute a snipe
func (h *SniperHandler) ExecuteSnipe(w http.ResponseWriter, r *http.Request) {
	symbol := chi.URLParam(r, "symbol")

	if symbol == "" {
		response.WriteJSON(w, http.StatusBadRequest, response.Error("missing_symbol", "Symbol is required"))
		return
	}

	// Check if custom config is provided
	var config *port.SniperConfig
	if r.ContentLength > 0 {
		if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
			h.logger.Error().Err(err).Msg("Failed to decode sniper config")
			response.WriteJSON(w, http.StatusBadRequest, response.Error("invalid_format", "Invalid configuration format"))
			return
		}
	}

	var order interface{}
	var err error

	if config != nil {
		// Execute with custom config
		order, err = h.sniperUC.ExecuteSnipeWithConfig(r.Context(), symbol, config)
	} else {
		// Execute with default config
		order, err = h.sniperUC.ExecuteSnipe(r.Context(), symbol)
	}

	if err != nil {
		h.logger.Error().Err(err).Str("symbol", symbol).Msg("Failed to execute snipe")
		response.WriteJSON(w, http.StatusInternalServerError, response.Error("execution_failed", err.Error()))
		return
	}

	response.WriteJSON(w, http.StatusOK, response.Success(order))
}

// GetSniperConfig handles the request to get the sniper configuration
func (h *SniperHandler) GetSniperConfig(w http.ResponseWriter, r *http.Request) {
	config, err := h.sniperUC.GetSniperConfig()
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get sniper config")
		response.WriteJSON(w, http.StatusInternalServerError, response.Error("config_retrieval_failed", err.Error()))
		return
	}

	response.WriteJSON(w, http.StatusOK, response.Success(config))
}

// UpdateSniperConfig handles the request to update the sniper configuration
func (h *SniperHandler) UpdateSniperConfig(w http.ResponseWriter, r *http.Request) {
	var config port.SniperConfig
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		h.logger.Error().Err(err).Msg("Failed to decode sniper config")
		response.WriteJSON(w, http.StatusBadRequest, response.Error("invalid_format", "Invalid configuration format"))
		return
	}

	if err := h.sniperUC.UpdateSniperConfig(&config); err != nil {
		h.logger.Error().Err(err).Msg("Failed to update sniper config")
		response.WriteJSON(w, http.StatusInternalServerError, response.Error("config_update_failed", err.Error()))
		return
	}

	response.WriteJSON(w, http.StatusOK, response.Success(map[string]string{"status": "config updated"}))
}

// GetSniperStatus handles the request to get the sniper status
func (h *SniperHandler) GetSniperStatus(w http.ResponseWriter, r *http.Request) {
	status, err := h.sniperUC.GetSniperStatus()
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get sniper status")
		response.WriteJSON(w, http.StatusInternalServerError, response.Error("status_retrieval_failed", err.Error()))
		return
	}

	response.WriteJSON(w, http.StatusOK, response.Success(map[string]string{"status": status}))
}

// StartSniper handles the request to start the sniper service
func (h *SniperHandler) StartSniper(w http.ResponseWriter, r *http.Request) {
	if err := h.sniperUC.StartSniper(); err != nil {
		h.logger.Error().Err(err).Msg("Failed to start sniper")
		response.WriteJSON(w, http.StatusInternalServerError, response.Error("start_failed", err.Error()))
		return
	}

	response.WriteJSON(w, http.StatusOK, response.Success(map[string]string{"status": "sniper started"}))
}

// StopSniper handles the request to stop the sniper service
func (h *SniperHandler) StopSniper(w http.ResponseWriter, r *http.Request) {
	if err := h.sniperUC.StopSniper(); err != nil {
		h.logger.Error().Err(err).Msg("Failed to stop sniper")
		response.WriteJSON(w, http.StatusInternalServerError, response.Error("stop_failed", err.Error()))
		return
	}

	response.WriteJSON(w, http.StatusOK, response.Success(map[string]string{"status": "sniper stopped"}))
}

// ConfigureAutoSnipe handles the request to configure auto-sniping
func (h *SniperHandler) ConfigureAutoSnipe(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Enabled bool               `json:"enabled"`
		Config  *port.SniperConfig `json:"config,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.logger.Error().Err(err).Msg("Failed to decode auto-snipe request")
		response.WriteJSON(w, http.StatusBadRequest, response.Error("invalid_format", "Invalid request format"))
		return
	}

	if err := h.sniperUC.SetupAutoSnipe(request.Enabled, request.Config); err != nil {
		h.logger.Error().Err(err).Msg("Failed to configure auto-snipe")
		response.WriteJSON(w, http.StatusInternalServerError, response.Error("config_failed", err.Error()))
		return
	}

	response.WriteJSON(w, http.StatusOK, response.Success(map[string]interface{}{
		"status":  "auto-snipe configured",
		"enabled": request.Enabled,
	}))
}
