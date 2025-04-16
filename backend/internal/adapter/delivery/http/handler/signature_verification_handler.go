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

// SignatureVerificationHandler handles signature verification HTTP requests
type SignatureVerificationHandler struct {
	verificationService usecase.SignatureVerificationService
	logger              *zerolog.Logger
}

// NewSignatureVerificationHandler creates a new SignatureVerificationHandler
func NewSignatureVerificationHandler(
	verificationService usecase.SignatureVerificationService,
	logger *zerolog.Logger,
) *SignatureVerificationHandler {
	return &SignatureVerificationHandler{
		verificationService: verificationService,
		logger:              logger,
	}
}

// RegisterRoutes registers the signature verification routes
func (h *SignatureVerificationHandler) RegisterRoutes(r chi.Router, authMiddleware middleware.AuthMiddleware) {
	r.Route("/wallet-verification", func(r chi.Router) {
		// Protected routes
		r.Group(func(r chi.Router) {
			// Use auth middleware from the router
			r.Use(authMiddleware.RequireAuthentication)
			r.Post("/challenge/{id}", h.GenerateChallenge)
			r.Post("/verify/{id}", h.VerifySignature)
			r.Get("/status/{id}", h.GetWalletStatus)
			r.Put("/status/{id}", h.SetWalletStatus)
		})
	})
}

// GenerateChallenge handles the generate challenge endpoint
func (h *SignatureVerificationHandler) GenerateChallenge(w http.ResponseWriter, r *http.Request) {
	// Get wallet ID from URL
	walletID := chi.URLParam(r, "id")
	if walletID == "" {
		http.Error(w, "Wallet ID is required", http.StatusBadRequest)
		return
	}

	// Generate challenge
	challenge, err := h.verificationService.GenerateChallenge(r.Context(), walletID)
	if err != nil {
		h.logger.Error().Err(err).Str("id", walletID).Msg("Failed to generate challenge")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Write response
	json.NewEncoder(w).Encode(map[string]interface{}{
		"challenge": challenge,
	})
}

// VerifySignature handles the verify signature endpoint
func (h *SignatureVerificationHandler) VerifySignature(w http.ResponseWriter, r *http.Request) {
	// Get wallet ID from URL
	walletID := chi.URLParam(r, "id")
	if walletID == "" {
		http.Error(w, "Wallet ID is required", http.StatusBadRequest)
		return
	}

	// Parse request body
	var request struct {
		Challenge string `json:"challenge"`
		Signature string `json:"signature"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		apperror.WriteError(w, apperror.NewInvalid("Invalid request body", nil, nil))
		return
	}

	// Verify signature
	verified, err := h.verificationService.VerifySignature(r.Context(), walletID, request.Challenge, request.Signature)
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

// GetWalletStatus handles the get wallet status endpoint
func (h *SignatureVerificationHandler) GetWalletStatus(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	_, ok := r.Context().Value("userID").(string)
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

	// Get wallet status
	status, err := h.verificationService.GetWalletStatus(r.Context(), walletID)
	if err != nil {
		h.logger.Error().Err(err).Str("id", walletID).Msg("Failed to get wallet status")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Write response
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": status,
	})
}

// SetWalletStatus handles the set wallet status endpoint
func (h *SignatureVerificationHandler) SetWalletStatus(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	_, ok := r.Context().Value("userID").(string)
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

	// Parse request body
	var request struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		apperror.WriteError(w, apperror.NewInvalid("Invalid request body", nil, nil))
		return
	}

	// Validate status
	var status model.WalletStatus
	switch model.WalletStatus(request.Status) {
	case model.WalletStatusActive, model.WalletStatusInactive, model.WalletStatusPending, model.WalletStatusVerified, model.WalletStatusFailed:
		status = model.WalletStatus(request.Status)
	default:
		apperror.WriteError(w, apperror.NewInvalid("Invalid wallet status", nil, nil))
		return
	}

	// Set wallet status
	if err := h.verificationService.SetWalletStatus(r.Context(), walletID, status); err != nil {
		h.logger.Error().Err(err).Str("id", walletID).Msg("Failed to set wallet status")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Write response
	w.WriteHeader(http.StatusNoContent)
}
