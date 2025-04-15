package handler

import (
	"encoding/json"
	"net/http"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/apperror"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/usecase"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
)

type StatusHandler struct {
	useCase usecase.StatusUseCase
	logger  *zerolog.Logger
}

func NewStatusHandler(useCase usecase.StatusUseCase, logger *zerolog.Logger) *StatusHandler {
	return &StatusHandler{
		useCase: useCase,
		logger:  logger,
	}
}

func (h *StatusHandler) RegisterRoutes(r chi.Router) {
	r.Route("/status", func(r chi.Router) {
		r.Get("/services", h.GetServicesStatus)
		r.Get("/exchange", h.GetExchangeStatus)
		r.Get("/exchanges", h.GetExchangesStatus)
	})
}

// GetServicesStatus returns the status of all services
func (h *StatusHandler) GetServicesStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.logger.Debug().Msg("Getting services status")

	// Get system status from use case
	systemStatus, err := h.useCase.GetSystemStatus(ctx)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get system status")
		apperror.WriteError(w, apperror.NewInternal(err))
		return
	}

	// Return the status
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(systemStatus); err != nil {
		h.logger.Error().Err(err).Msg("Failed to encode system status")
	}
}

// GetExchangeStatus returns the status of the exchange
func (h *StatusHandler) GetExchangeStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.logger.Debug().Msg("Getting exchange status")

	// Get component status for the exchange
	componentStatus, err := h.useCase.GetComponentStatus(ctx, "mexc_api")
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get exchange status")

		// Use more specific error types based on the error
		switch {
		case err.Error() == "component not found":
			apperror.WriteError(w, apperror.NewNotFound("exchange_status", nil, err))
		case err.Error() == "context canceled" || err.Error() == "context deadline exceeded":
			apperror.WriteError(w, apperror.NewExternalService("exchange_api", "Exchange status check timed out", err))
		default:
			apperror.WriteError(w, apperror.NewInternal(err))
		}
		return
	}

	// Check if the component is in error state and provide more detailed response
	if componentStatus.Status == "error" {
		h.logger.Warn().Str("component", "mexc_api").Str("error", componentStatus.LastError).Msg("Exchange is in error state")
	}

	// Return the status
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(componentStatus); err != nil {
		h.logger.Error().Err(err).Msg("Failed to encode exchange status")
	}
}

// GetExchangesStatus returns the status of all exchanges
func (h *StatusHandler) GetExchangesStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.logger.Debug().Msg("Getting all exchanges status")

	// For now, we just return the MEXC exchange status
	componentStatus, err := h.useCase.GetComponentStatus(ctx, "mexc_api")
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get exchange status")
		apperror.WriteError(w, apperror.NewInternal(err))
		return
	}

	// Return the status as a map of exchanges
	exchanges := map[string]interface{}{
		"mexc": componentStatus,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    exchanges,
	}); err != nil {
		h.logger.Error().Err(err).Msg("Failed to encode exchanges status")
	}
}
