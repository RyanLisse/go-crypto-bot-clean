package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/http/middleware"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/apperror"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/usecase"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
)

// WalletHandler handles wallet-related HTTP requests
type WalletHandler struct {
	walletService usecase.WalletService
	logger        *zerolog.Logger
}

// NewWalletHandler creates a new WalletHandler
func NewWalletHandler(walletService usecase.WalletService, logger *zerolog.Logger) *WalletHandler {
	return &WalletHandler{
		walletService: walletService,
		logger:        logger,
	}
}

// RegisterRoutes registers the wallet routes
func (h *WalletHandler) RegisterRoutes(r chi.Router, authMiddleware middleware.AuthMiddleware) {
	r.Route("/wallets", func(r chi.Router) {
		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.RequireAuthentication)

			r.Get("/", h.GetWallets)
			r.Post("/", h.CreateWallet)
			r.Get("/{id}", h.GetWallet)
			r.Delete("/{id}", h.DeleteWallet)
			r.Put("/{id}/metadata", h.UpdateWalletMetadata)
			r.Put("/{id}/primary", h.SetPrimaryWallet)
			r.Get("/{id}/balances/{asset}", h.GetBalance)
			r.Post("/{id}/refresh", h.RefreshWallet)
			r.Get("/history/{asset}", h.GetBalanceHistory)
		})
	})
}

// GetWallets handles the get wallets endpoint
func (h *WalletHandler) GetWallets(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		apperror.WriteError(w, apperror.NewUnauthorized("User ID not found in context", nil))
		return
	}

	// Get wallets
	wallets, err := h.walletService.GetWalletsByUserID(r.Context(), userID)
	if err != nil {
		h.logger.Error().Err(err).Str("userID", userID).Msg("Failed to get wallets")
		apperror.WriteError(w, apperror.NewInternal(err))
		return
	}

	// Write response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(wallets)
}

// CreateWallet handles the create wallet endpoint
func (h *WalletHandler) CreateWallet(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		apperror.WriteError(w, apperror.NewUnauthorized("User ID not found in context", nil))
		return
	}

	// Parse request body
	var request struct {
		Exchange string          `json:"exchange"`
		Type     model.WalletType `json:"type"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		apperror.WriteError(w, apperror.NewInvalid("Invalid request body", nil, nil))
		return
	}

	// Create wallet
	wallet, err := h.walletService.CreateWallet(r.Context(), userID, request.Exchange, request.Type)
	if err != nil {
		h.logger.Error().Err(err).Str("userID", userID).Msg("Failed to create wallet")
		apperror.WriteError(w, apperror.NewInternal(err))
		return
	}

	// Write response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(wallet)
}

// GetWallet handles the get wallet endpoint
func (h *WalletHandler) GetWallet(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		apperror.WriteError(w, apperror.NewUnauthorized("User ID not found in context", nil))
		return
	}

	// Get wallet ID from URL
	walletID := chi.URLParam(r, "id")
	if walletID == "" {
		apperror.WriteError(w, apperror.NewInvalid("Wallet ID is required", nil, nil))
		return
	}

	// Get wallet
	wallet, err := h.walletService.GetWallet(r.Context(), walletID)
	if err != nil {
		h.logger.Error().Err(err).Str("id", walletID).Msg("Failed to get wallet")
		apperror.WriteError(w, apperror.NewInternal(err))
		return
	}
	if wallet == nil {
		apperror.WriteError(w, apperror.NewNotFound("Wallet", walletID, nil))
		return
	}

	// Check if wallet belongs to user
	if wallet.UserID != userID {
		apperror.WriteError(w, apperror.NewForbidden("Wallet does not belong to user", nil))
		return
	}

	// Write response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(wallet)
}

// DeleteWallet handles the delete wallet endpoint
func (h *WalletHandler) DeleteWallet(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		apperror.WriteError(w, apperror.NewUnauthorized("User ID not found in context", nil))
		return
	}

	// Get wallet ID from URL
	walletID := chi.URLParam(r, "id")
	if walletID == "" {
		apperror.WriteError(w, apperror.NewInvalid("Wallet ID is required", nil, nil))
		return
	}

	// Get wallet
	wallet, err := h.walletService.GetWallet(r.Context(), walletID)
	if err != nil {
		h.logger.Error().Err(err).Str("id", walletID).Msg("Failed to get wallet")
		apperror.WriteError(w, apperror.NewInternal(err))
		return
	}
	if wallet == nil {
		apperror.WriteError(w, apperror.NewNotFound("Wallet", walletID, nil))
		return
	}

	// Check if wallet belongs to user
	if wallet.UserID != userID {
		apperror.WriteError(w, apperror.NewForbidden("Wallet does not belong to user", nil))
		return
	}

	// Delete wallet
	if err := h.walletService.DeleteWallet(r.Context(), walletID); err != nil {
		h.logger.Error().Err(err).Str("id", walletID).Msg("Failed to delete wallet")
		apperror.WriteError(w, apperror.NewInternal(err))
		return
	}

	// Write response
	w.WriteHeader(http.StatusNoContent)
}

// UpdateWalletMetadata handles the update wallet metadata endpoint
func (h *WalletHandler) UpdateWalletMetadata(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		apperror.WriteError(w, apperror.NewUnauthorized("User ID not found in context", nil))
		return
	}

	// Get wallet ID from URL
	walletID := chi.URLParam(r, "id")
	if walletID == "" {
		apperror.WriteError(w, apperror.NewInvalid("Wallet ID is required", nil, nil))
		return
	}

	// Get wallet
	wallet, err := h.walletService.GetWallet(r.Context(), walletID)
	if err != nil {
		h.logger.Error().Err(err).Str("id", walletID).Msg("Failed to get wallet")
		apperror.WriteError(w, apperror.NewInternal(err))
		return
	}
	if wallet == nil {
		apperror.WriteError(w, apperror.NewNotFound("Wallet", walletID, nil))
		return
	}

	// Check if wallet belongs to user
	if wallet.UserID != userID {
		apperror.WriteError(w, apperror.NewForbidden("Wallet does not belong to user", nil))
		return
	}

	// Parse request body
	var request struct {
		Name        string   `json:"name"`
		Description string   `json:"description"`
		Tags        []string `json:"tags"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		apperror.WriteError(w, apperror.NewInvalid("Invalid request body", nil, nil))
		return
	}

	// Update metadata
	if err := h.walletService.SetWalletMetadata(r.Context(), walletID, request.Name, request.Description, request.Tags); err != nil {
		h.logger.Error().Err(err).Str("id", walletID).Msg("Failed to update wallet metadata")
		apperror.WriteError(w, apperror.NewInternal(err))
		return
	}

	// Get updated wallet
	wallet, err = h.walletService.GetWallet(r.Context(), walletID)
	if err != nil {
		h.logger.Error().Err(err).Str("id", walletID).Msg("Failed to get updated wallet")
		apperror.WriteError(w, apperror.NewInternal(err))
		return
	}

	// Write response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(wallet)
}

// SetPrimaryWallet handles the set primary wallet endpoint
func (h *WalletHandler) SetPrimaryWallet(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		apperror.WriteError(w, apperror.NewUnauthorized("User ID not found in context", nil))
		return
	}

	// Get wallet ID from URL
	walletID := chi.URLParam(r, "id")
	if walletID == "" {
		apperror.WriteError(w, apperror.NewInvalid("Wallet ID is required", nil, nil))
		return
	}

	// Get wallet
	wallet, err := h.walletService.GetWallet(r.Context(), walletID)
	if err != nil {
		h.logger.Error().Err(err).Str("id", walletID).Msg("Failed to get wallet")
		apperror.WriteError(w, apperror.NewInternal(err))
		return
	}
	if wallet == nil {
		apperror.WriteError(w, apperror.NewNotFound("Wallet", walletID, nil))
		return
	}

	// Check if wallet belongs to user
	if wallet.UserID != userID {
		apperror.WriteError(w, apperror.NewForbidden("Wallet does not belong to user", nil))
		return
	}

	// Set primary wallet
	if err := h.walletService.SetPrimaryWallet(r.Context(), userID, walletID); err != nil {
		h.logger.Error().Err(err).Str("id", walletID).Msg("Failed to set primary wallet")
		apperror.WriteError(w, apperror.NewInternal(err))
		return
	}

	// Get updated wallet
	wallet, err = h.walletService.GetWallet(r.Context(), walletID)
	if err != nil {
		h.logger.Error().Err(err).Str("id", walletID).Msg("Failed to get updated wallet")
		apperror.WriteError(w, apperror.NewInternal(err))
		return
	}

	// Write response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(wallet)
}

// GetBalance handles the get balance endpoint
func (h *WalletHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		apperror.WriteError(w, apperror.NewUnauthorized("User ID not found in context", nil))
		return
	}

	// Get wallet ID from URL
	walletID := chi.URLParam(r, "id")
	if walletID == "" {
		apperror.WriteError(w, apperror.NewInvalid("Wallet ID is required", nil, nil))
		return
	}

	// Get asset from URL
	assetStr := chi.URLParam(r, "asset")
	if assetStr == "" {
		apperror.WriteError(w, apperror.NewInvalid("Asset is required", nil, nil))
		return
	}
	asset := model.Asset(assetStr)

	// Get wallet
	wallet, err := h.walletService.GetWallet(r.Context(), walletID)
	if err != nil {
		h.logger.Error().Err(err).Str("id", walletID).Msg("Failed to get wallet")
		apperror.WriteError(w, apperror.NewInternal(err))
		return
	}
	if wallet == nil {
		apperror.WriteError(w, apperror.NewNotFound("Wallet", walletID, nil))
		return
	}

	// Check if wallet belongs to user
	if wallet.UserID != userID {
		apperror.WriteError(w, apperror.NewForbidden("Wallet does not belong to user", nil))
		return
	}

	// Get balance
	balance, err := h.walletService.GetBalance(r.Context(), walletID, asset)
	if err != nil {
		h.logger.Error().Err(err).Str("id", walletID).Str("asset", string(asset)).Msg("Failed to get balance")
		apperror.WriteError(w, apperror.NewInternal(err))
		return
	}

	// Write response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(balance)
}

// RefreshWallet handles the refresh wallet endpoint
func (h *WalletHandler) RefreshWallet(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		apperror.WriteError(w, apperror.NewUnauthorized("User ID not found in context", nil))
		return
	}

	// Get wallet ID from URL
	walletID := chi.URLParam(r, "id")
	if walletID == "" {
		apperror.WriteError(w, apperror.NewInvalid("Wallet ID is required", nil, nil))
		return
	}

	// Get wallet
	wallet, err := h.walletService.GetWallet(r.Context(), walletID)
	if err != nil {
		h.logger.Error().Err(err).Str("id", walletID).Msg("Failed to get wallet")
		apperror.WriteError(w, apperror.NewInternal(err))
		return
	}
	if wallet == nil {
		apperror.WriteError(w, apperror.NewNotFound("Wallet", walletID, nil))
		return
	}

	// Check if wallet belongs to user
	if wallet.UserID != userID {
		apperror.WriteError(w, apperror.NewForbidden("Wallet does not belong to user", nil))
		return
	}

	// Refresh wallet
	if err := h.walletService.RefreshWallet(r.Context(), walletID); err != nil {
		h.logger.Error().Err(err).Str("id", walletID).Msg("Failed to refresh wallet")
		apperror.WriteError(w, apperror.NewInternal(err))
		return
	}

	// Get updated wallet
	wallet, err = h.walletService.GetWallet(r.Context(), walletID)
	if err != nil {
		h.logger.Error().Err(err).Str("id", walletID).Msg("Failed to get updated wallet")
		apperror.WriteError(w, apperror.NewInternal(err))
		return
	}

	// Write response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(wallet)
}

// GetBalanceHistory handles the get balance history endpoint
func (h *WalletHandler) GetBalanceHistory(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		apperror.WriteError(w, apperror.NewUnauthorized("User ID not found in context", nil))
		return
	}

	// Get asset from URL
	assetStr := chi.URLParam(r, "asset")
	if assetStr == "" {
		apperror.WriteError(w, apperror.NewInvalid("Asset is required", nil, nil))
		return
	}
	asset := model.Asset(assetStr)

	// Parse query parameters
	fromStr := r.URL.Query().Get("from")
	toStr := r.URL.Query().Get("to")

	var from, to time.Time
	var err error

	if fromStr != "" {
		from, err = time.Parse(time.RFC3339, fromStr)
		if err != nil {
			apperror.WriteError(w, apperror.NewInvalid("Invalid from date format", nil, nil))
			return
		}
	} else {
		// Default to 30 days ago
		from = time.Now().AddDate(0, 0, -30)
	}

	if toStr != "" {
		to, err = time.Parse(time.RFC3339, toStr)
		if err != nil {
			apperror.WriteError(w, apperror.NewInvalid("Invalid to date format", nil, nil))
			return
		}
	} else {
		// Default to now
		to = time.Now()
	}

	// Get balance history
	history, err := h.walletService.GetBalanceHistory(r.Context(), userID, asset, from, to)
	if err != nil {
		h.logger.Error().Err(err).Str("userID", userID).Str("asset", string(asset)).Msg("Failed to get balance history")
		apperror.WriteError(w, apperror.NewInternal(err))
		return
	}

	// Write response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(history)
}
