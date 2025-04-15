package handler

import (
	"encoding/json"
	"net/http"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/http/middleware"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/apperror"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/usecase"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
)

// WalletConnectionHandler handles wallet connection HTTP requests
type WalletConnectionHandler struct {
	connectionService usecase.WalletConnectionService
	logger            *zerolog.Logger
}

// NewWalletConnectionHandler creates a new WalletConnectionHandler
func NewWalletConnectionHandler(
	connectionService usecase.WalletConnectionService,
	logger *zerolog.Logger,
) *WalletConnectionHandler {
	return &WalletConnectionHandler{
		connectionService: connectionService,
		logger:            logger,
	}
}

// RegisterRoutes registers the wallet connection routes
func (h *WalletConnectionHandler) RegisterRoutes(r chi.Router, authMiddleware *middleware.EnhancedClerkMiddleware) {
	r.Route("/wallet-connection", func(r chi.Router) {
		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.RequireAuthentication)

			r.Get("/providers", h.GetProviders)
			r.Get("/providers/{type}", h.GetProvidersByType)
			r.Post("/connect", h.Connect)
			r.Post("/disconnect/{id}", h.Disconnect)
			r.Post("/verify/{id}", h.Verify)
			r.Post("/refresh/{id}", h.RefreshWallet)
			r.Get("/validate-address/{provider}/{address}", h.ValidateAddress)
		})
	})
}

// GetProviders handles the get providers endpoint
func (h *WalletConnectionHandler) GetProviders(w http.ResponseWriter, r *http.Request) {
	// Get providers
	providers, err := h.connectionService.GetProviders(r.Context())
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get providers")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Write response
	json.NewEncoder(w).Encode(map[string]interface{}{
		"providers": providers,
	})
}

// GetProvidersByType handles the get providers by type endpoint
func (h *WalletConnectionHandler) GetProvidersByType(w http.ResponseWriter, r *http.Request) {
	// Get type from URL
	typStr := chi.URLParam(r, "type")
	if typStr == "" {
		apperror.WriteError(w, apperror.NewInvalid("Type is required", nil, nil))
		return
	}

	// Convert to wallet type
	var typ model.WalletType
	switch typStr {
	case "exchange":
		typ = model.WalletTypeExchange
	case "web3":
		typ = model.WalletTypeWeb3
	default:
		apperror.WriteError(w, apperror.NewInvalid("Invalid wallet type", nil, nil))
		return
	}

	// Get providers
	providers, err := h.connectionService.GetProvidersByType(r.Context(), typ)
	if err != nil {
		h.logger.Error().Err(err).Str("type", string(typ)).Msg("Failed to get providers by type")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Write response
	json.NewEncoder(w).Encode(map[string]interface{}{
		"providers": providers,
	})
}

// Connect handles the connect endpoint
func (h *WalletConnectionHandler) Connect(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var request struct {
		Provider string                 `json:"provider"`
		Params   map[string]interface{} `json:"params"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		apperror.WriteError(w, apperror.NewInvalid("Invalid request body", nil, nil))
		return
	}

	// Get user ID from context
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		apperror.WriteError(w, apperror.NewUnauthorized("User ID not found in context", nil))
		return
	}

	// Connect to provider
	wallet, err := h.connectionService.Connect(r.Context(), userID, request.Provider, request.Params)
	if err != nil {
		h.logger.Error().Err(err).Str("provider", request.Provider).Msg("Failed to connect to provider")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Write response
	json.NewEncoder(w).Encode(wallet)
}

// Disconnect handles the disconnect endpoint
func (h *WalletConnectionHandler) Disconnect(w http.ResponseWriter, r *http.Request) {
	// Get wallet ID from URL
	walletID := chi.URLParam(r, "id")
	if walletID == "" {
		http.Error(w, "Wallet ID is required", http.StatusBadRequest)
		return
	}

	// Disconnect from provider
	if err := h.connectionService.Disconnect(r.Context(), walletID); err != nil {
		h.logger.Error().Err(err).Str("id", walletID).Msg("Failed to disconnect from provider")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Write response
	w.WriteHeader(http.StatusNoContent)
}

// Verify handles the verify endpoint
func (h *WalletConnectionHandler) Verify(w http.ResponseWriter, r *http.Request) {
	// Get wallet ID from URL
	walletID := chi.URLParam(r, "id")
	if walletID == "" {
		http.Error(w, "Wallet ID is required", http.StatusBadRequest)
		return
	}

	// Parse request body
	var request struct {
		Message   string `json:"message"`
		Signature string `json:"signature"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		apperror.WriteError(w, apperror.NewInvalid("Invalid request body", nil, nil))
		return
	}

	// Verify signature
	verified, err := h.connectionService.Verify(r.Context(), walletID, request.Message, request.Signature)
	if err != nil {
		h.logger.Error().Err(err).Str("id", walletID).Msg("Failed to verify signature")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Write response
	json.NewEncoder(w).Encode(map[string]interface{}{
		"verified": verified,
	})
}

// RefreshWallet handles the refresh wallet endpoint
func (h *WalletConnectionHandler) RefreshWallet(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	_, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		apperror.WriteError(w, apperror.NewUnauthorized("User ID not found in context", nil))
		return
	}

	// Get wallet ID from URL
	walletID := chi.URLParam(r, "id")
	if walletID == "" {
		http.Error(w, "Wallet ID is required", http.StatusBadRequest)
		return
	}

	// Refresh wallet
	wallet, err := h.connectionService.RefreshWallet(r.Context(), walletID)
	if err != nil {
		h.logger.Error().Err(err).Str("id", walletID).Msg("Failed to refresh wallet")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Write response
	json.NewEncoder(w).Encode(wallet)
}

// ValidateAddress handles the validate address endpoint
func (h *WalletConnectionHandler) ValidateAddress(w http.ResponseWriter, r *http.Request) {
	// Get provider from URL
	provider := chi.URLParam(r, "provider")
	if provider == "" {
		http.Error(w, "Provider is required", http.StatusBadRequest)
		return
	}

	// Get address from URL
	address := chi.URLParam(r, "address")
	if address == "" {
		http.Error(w, "Address is required", http.StatusBadRequest)
		return
	}

	// Validate address
	valid, err := h.connectionService.IsValidAddress(r.Context(), provider, address)
	if err != nil {
		h.logger.Error().Err(err).Str("provider", provider).Str("address", address).Msg("Failed to validate address")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Write response
	json.NewEncoder(w).Encode(map[string]interface{}{
		"valid": valid,
	})
}
