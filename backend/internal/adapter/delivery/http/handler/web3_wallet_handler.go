package handler

import (
	"encoding/json"
	"net/http"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/http/middleware"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/apperror"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/usecase"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
)

// Web3WalletHandler handles Web3 wallet HTTP requests
type Web3WalletHandler struct {
	web3WalletService usecase.Web3WalletService
	logger            *zerolog.Logger
}

// NewWeb3WalletHandler creates a new Web3WalletHandler
func NewWeb3WalletHandler(
	web3WalletService usecase.Web3WalletService,
	logger *zerolog.Logger,
) *Web3WalletHandler {
	return &Web3WalletHandler{
		web3WalletService: web3WalletService,
		logger:            logger,
	}
}

// RegisterRoutes registers the Web3 wallet routes
func (h *Web3WalletHandler) RegisterRoutes(r chi.Router, authMiddleware middleware.AuthMiddleware) {
	r.Route("/web3-wallets", func(r chi.Router) {
		// Public routes
		r.Get("/networks", h.GetSupportedNetworks)
		r.Post("/validate-address", h.ValidateAddress)

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.RequireAuthentication)
			r.Post("/", h.ConnectWallet)
			r.Delete("/{id}", h.DisconnectWallet)
			r.Get("/{id}/balance", h.GetWalletBalance)
		})
	})
}

// ConnectWallet handles the connect wallet endpoint
func (h *Web3WalletHandler) ConnectWallet(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		apperror.WriteError(w, apperror.NewUnauthorized("User ID not found in context", nil))
		return
	}

	// Parse request body
	var request struct {
		Network string `json:"network"`
		Address string `json:"address"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		apperror.WriteError(w, apperror.NewInvalid("Invalid request body", nil, nil))
		return
	}

	// Validate request
	if request.Network == "" {
		apperror.WriteError(w, apperror.NewInvalid("Network is required", nil, nil))
		return
	}
	if request.Address == "" {
		apperror.WriteError(w, apperror.NewInvalid("Address is required", nil, nil))
		return
	}

	// Connect wallet
	wallet, err := h.web3WalletService.ConnectWallet(r.Context(), userID, request.Network, request.Address)
	if err != nil {
		h.logger.Error().Err(err).Str("network", request.Network).Str("address", request.Address).Msg("Failed to connect wallet")
		apperror.WriteError(w, apperror.NewInternal(err))
		return
	}

	// Write response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(wallet)
}

// DisconnectWallet handles the disconnect wallet endpoint
func (h *Web3WalletHandler) DisconnectWallet(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	_, ok := middleware.GetUserIDFromContext(r.Context())
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

	// Disconnect wallet
	err := h.web3WalletService.DisconnectWallet(r.Context(), walletID)
	if err != nil {
		h.logger.Error().Err(err).Str("id", walletID).Msg("Failed to disconnect wallet")
		apperror.WriteError(w, apperror.NewInternal(err))
		return
	}

	// Write response
	w.WriteHeader(http.StatusNoContent)
}

// GetWalletBalance handles the get wallet balance endpoint
func (h *Web3WalletHandler) GetWalletBalance(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	_, ok := middleware.GetUserIDFromContext(r.Context())
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

	// Get wallet balance
	wallet, err := h.web3WalletService.GetWalletBalance(r.Context(), walletID)
	if err != nil {
		h.logger.Error().Err(err).Str("id", walletID).Msg("Failed to get wallet balance")
		apperror.WriteError(w, apperror.NewInternal(err))
		return
	}

	// Write response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(wallet)
}

// ValidateAddress handles the validate address endpoint
func (h *Web3WalletHandler) ValidateAddress(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var request struct {
		Network string `json:"network"`
		Address string `json:"address"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		apperror.WriteError(w, apperror.NewInvalid("Invalid request body", nil, nil))
		return
	}

	// Validate request
	if request.Network == "" {
		apperror.WriteError(w, apperror.NewInvalid("Network is required", nil, nil))
		return
	}
	if request.Address == "" {
		apperror.WriteError(w, apperror.NewInvalid("Address is required", nil, nil))
		return
	}

	// Validate address
	valid, err := h.web3WalletService.IsValidAddress(r.Context(), request.Network, request.Address)
	if err != nil {
		h.logger.Error().Err(err).Str("network", request.Network).Str("address", request.Address).Msg("Failed to validate address")
		apperror.WriteError(w, apperror.NewInternal(err))
		return
	}

	// Write response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"valid": valid,
	})
}

// GetSupportedNetworks handles the get supported networks endpoint
func (h *Web3WalletHandler) GetSupportedNetworks(w http.ResponseWriter, r *http.Request) {
	// Get supported networks
	networks, err := h.web3WalletService.GetSupportedNetworks(r.Context())
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get supported networks")
		apperror.WriteError(w, apperror.NewInternal(err))
		return
	}

	// Write response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"networks": networks,
	})
}
