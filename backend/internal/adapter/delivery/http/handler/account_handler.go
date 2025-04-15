package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/http/middleware"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/apperror"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/usecase"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
)

// AccountHandler handles account-related endpoints
type AccountHandler struct {
	useCase usecase.AccountUsecase
	logger  *zerolog.Logger
}

// NewAccountHandler creates a new AccountHandler
func NewAccountHandler(useCase usecase.AccountUsecase, logger *zerolog.Logger) *AccountHandler {
	return &AccountHandler{
		useCase: useCase,
		logger:  logger,
	}
}

// RegisterRoutes registers the account routes
func (h *AccountHandler) RegisterRoutes(r chi.Router) {
	h.logger.Info().Msg("Registering account routes")

	r.Route("/account", func(r chi.Router) {
		h.logger.Info().Msg("Setting up /account routes")
		r.Get("/wallet", h.GetWallet)
		r.Get("/balance/{asset}", h.GetBalanceHistory)
		r.Post("/refresh", h.RefreshWallet)
	})
	h.logger.Info().Msg("Account routes registered")
}

// GetWallet returns the user's wallet
func (h *AccountHandler) GetWallet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get user ID from context or use a default for testing
	userID := "MEXC_USER" // Default user ID for direct API access

	// Get MEXC API credentials from context
	credentials := middleware.GetMEXCAPICredentials(ctx)
	if credentials != nil {
		h.logger.Info().Str("apiKey", credentials.APIKey).Msg("Using MEXC API credentials from context")
	} else {
		h.logger.Warn().Msg("No MEXC API credentials found in context")
	}

	// Get the wallet from the use case
	wallet, err := h.useCase.GetWallet(ctx, userID)
	if err != nil {
		h.logger.Error().Err(err).Str("userID", userID).Msg("Failed to get wallet")
		apperror.WriteError(w, apperror.NewInternal(err))
		return
	}

	// Return the wallet
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    wallet,
	}); err != nil {
		h.logger.Error().Err(err).Msg("Failed to encode wallet")
	}
}

// GetBalanceHistory returns the balance history for a specific asset
func (h *AccountHandler) GetBalanceHistory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get the asset from the URL
	assetStr := chi.URLParam(r, "asset")
	if assetStr == "" {
		apperror.WriteError(w, apperror.NewInvalid("Asset is required", nil, nil))
		return
	}

	// Parse days parameter
	daysStr := r.URL.Query().Get("days")
	days := 30 // Default to 30 days
	if daysStr != "" {
		parsedDays, err := strconv.Atoi(daysStr)
		if err != nil || parsedDays <= 0 {
			apperror.WriteError(w, apperror.NewInvalid("Days must be a positive integer", nil, err))
			return
		}
		days = parsedDays
	}

	// Get user ID from context or use a default for testing
	userID := "MEXC_USER" // Default user ID for direct API access

	// Get MEXC API credentials from context
	credentials := middleware.GetMEXCAPICredentials(ctx)
	if credentials != nil {
		h.logger.Info().Str("apiKey", credentials.APIKey).Msg("Using MEXC API credentials from context")
	} else {
		h.logger.Warn().Msg("No MEXC API credentials found in context")
	}

	// Calculate from and to dates
	to := time.Now()
	from := to.AddDate(0, 0, -days)

	// Get the balance history from the use case
	history, err := h.useCase.GetBalanceHistory(ctx, userID, model.Asset(assetStr), from, to)
	if err != nil {
		h.logger.Error().Err(err).Str("userID", userID).Str("asset", assetStr).Int("days", days).Msg("Failed to get balance history")
		apperror.WriteError(w, apperror.NewInternal(err))
		return
	}

	// Return the balance history
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    history,
	}); err != nil {
		h.logger.Error().Err(err).Msg("Failed to encode balance history")
	}
}

// RefreshWallet refreshes the user's wallet
func (h *AccountHandler) RefreshWallet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get user ID from context or use a default for testing
	userID := "MEXC_USER" // Default user ID for direct API access

	// Get MEXC API credentials from context
	credentials := middleware.GetMEXCAPICredentials(ctx)
	if credentials != nil {
		h.logger.Info().Str("apiKey", credentials.APIKey).Msg("Using MEXC API credentials from context")
	} else {
		h.logger.Warn().Msg("No MEXC API credentials found in context")
	}

	// Refresh the wallet
	err := h.useCase.RefreshWallet(ctx, userID)
	if err != nil {
		h.logger.Error().Err(err).Str("userID", userID).Msg("Failed to refresh wallet")
		apperror.WriteError(w, apperror.NewInternal(err))
		return
	}

	// Return success
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"success":   true,
		"message":   "Wallet refreshed successfully",
		"timestamp": time.Now(),
	}); err != nil {
		h.logger.Error().Err(err).Msg("Failed to encode response")
	}
}
