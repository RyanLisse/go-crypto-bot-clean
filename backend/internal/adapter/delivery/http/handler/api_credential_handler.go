package handler

import (
	"encoding/json"
	"net/http"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/delivery/http/validation"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/apperror"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/usecase"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
)

// APICredentialHandler handles API credential-related endpoints
type APICredentialHandler struct {
	useCase usecase.APICredentialUseCase
	logger  *zerolog.Logger
}

// NewAPICredentialHandler creates a new APICredentialHandler
func NewAPICredentialHandler(useCase usecase.APICredentialUseCase, logger *zerolog.Logger) *APICredentialHandler {
	return &APICredentialHandler{
		useCase: useCase,
		logger:  logger,
	}
}

// RegisterRoutes registers the API credential routes
func (h *APICredentialHandler) RegisterRoutes(r chi.Router) {
	r.Post("/credentials", h.CreateCredential)
	r.Get("/credentials", h.ListCredentials)
	r.Get("/credentials/{id}", h.GetCredential)
	r.Put("/credentials/{id}", h.UpdateCredential)
	r.Delete("/credentials/{id}", h.DeleteCredential)
}

// CreateCredentialRequest represents the request body for creating an API credential
type CreateCredentialRequest struct {
	Exchange  string `json:"exchange"`
	APIKey    string `json:"apiKey"`
	APISecret string `json:"apiSecret"`
	Label     string `json:"label"`
}

// CreateCredential creates a new API credential
func (h *APICredentialHandler) CreateCredential(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		h.logger.Error().Msg("User ID not found in context")
		apperror.WriteError(w, apperror.NewUnauthorized("User ID not found in context", nil))
		return
	}

	// Parse request body
	var req CreateCredentialRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error().Err(err).Msg("Failed to decode request body")
		apperror.WriteError(w, apperror.NewInvalid("Invalid request body", nil, err))
		return
	}

	// Validate request using the credential validator
	validator := validation.NewCredentialValidator()
	validator.ValidateExchange(req.Exchange).
		ValidateAPIKey(req.APIKey, req.Exchange).
		ValidateAPISecret(req.APISecret, req.Exchange).
		ValidateLabel(req.Label)

	if validator.HasErrors() {
		apperror.WriteError(w, validator.ToAppError())
		return
	}

	// Create credential
	credential := model.NewAPICredential(userID, req.Exchange, req.APIKey, req.APISecret, req.Label)

	// Save credential
	if err := h.useCase.CreateCredential(r.Context(), credential); err != nil {
		h.logger.Error().Err(err).Str("userID", userID).Msg("Failed to create API credential")
		apperror.WriteError(w, apperror.NewInternal(err))
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"id":        credential.ID,
			"exchange":  credential.Exchange,
			"apiKey":    credential.APIKey,
			"label":     credential.Label,
			"createdAt": credential.CreatedAt,
		},
	}); err != nil {
		h.logger.Error().Err(err).Msg("Failed to encode response")
	}
}

// ListCredentials lists API credentials for the current user
func (h *APICredentialHandler) ListCredentials(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		h.logger.Error().Msg("User ID not found in context")
		apperror.WriteError(w, apperror.NewUnauthorized("User ID not found in context", nil))
		return
	}

	// Get credentials
	credentials, err := h.useCase.ListCredentials(r.Context(), userID)
	if err != nil {
		h.logger.Error().Err(err).Str("userID", userID).Msg("Failed to list API credentials")
		apperror.WriteError(w, apperror.NewInternal(err))
		return
	}

	// Map credentials to response
	response := make([]map[string]interface{}, 0, len(credentials))
	for _, credential := range credentials {
		response = append(response, map[string]interface{}{
			"id":        credential.ID,
			"exchange":  credential.Exchange,
			"apiKey":    credential.APIKey,
			"label":     credential.Label,
			"createdAt": credential.CreatedAt,
		})
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    response,
	}); err != nil {
		h.logger.Error().Err(err).Msg("Failed to encode response")
	}
}

// GetCredential gets an API credential by ID
func (h *APICredentialHandler) GetCredential(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		h.logger.Error().Msg("User ID not found in context")
		apperror.WriteError(w, apperror.NewUnauthorized("User ID not found in context", nil))
		return
	}

	// Get credential ID from URL
	id := chi.URLParam(r, "id")
	if id == "" {
		apperror.WriteError(w, apperror.NewInvalid("Credential ID is required", nil, nil))
		return
	}

	// Get credential
	credential, err := h.useCase.GetCredential(r.Context(), id)
	if err != nil {
		h.logger.Error().Err(err).Str("id", id).Msg("Failed to get API credential")
		apperror.WriteError(w, apperror.NewInternal(err))
		return
	}

	// Check if credential exists
	if credential == nil {
		apperror.WriteError(w, apperror.NewNotFound("api_credential", id, nil))
		return
	}

	// Check if credential belongs to user
	if credential.UserID != userID {
		apperror.WriteError(w, apperror.NewForbidden("API credential does not belong to user", nil))
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"id":        credential.ID,
			"exchange":  credential.Exchange,
			"apiKey":    credential.APIKey,
			"label":     credential.Label,
			"createdAt": credential.CreatedAt,
			"updatedAt": credential.UpdatedAt,
		},
	}); err != nil {
		h.logger.Error().Err(err).Msg("Failed to encode response")
	}
}

// UpdateCredentialRequest represents the request body for updating an API credential
type UpdateCredentialRequest struct {
	APIKey    string `json:"apiKey"`
	APISecret string `json:"apiSecret"`
	Label     string `json:"label"`
}

// UpdateCredential updates an API credential
func (h *APICredentialHandler) UpdateCredential(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		h.logger.Error().Msg("User ID not found in context")
		apperror.WriteError(w, apperror.NewUnauthorized("User ID not found in context", nil))
		return
	}

	// Get credential ID from URL
	id := chi.URLParam(r, "id")
	if id == "" {
		apperror.WriteError(w, apperror.NewInvalid("Credential ID is required", nil, nil))
		return
	}

	// Parse request body
	var req UpdateCredentialRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error().Err(err).Msg("Failed to decode request body")
		apperror.WriteError(w, apperror.NewInvalid("Invalid request body", nil, err))
		return
	}

	// Get credential
	credential, err := h.useCase.GetCredential(r.Context(), id)
	if err != nil {
		h.logger.Error().Err(err).Str("id", id).Msg("Failed to get API credential")
		apperror.WriteError(w, apperror.NewInternal(err))
		return
	}

	// Check if credential exists
	if credential == nil {
		apperror.WriteError(w, apperror.NewNotFound("api_credential", id, nil))
		return
	}

	// Check if credential belongs to user
	if credential.UserID != userID {
		apperror.WriteError(w, apperror.NewForbidden("API credential does not belong to user", nil))
		return
	}

	// Validate request using the credential validator
	validator := validation.NewCredentialValidator()

	// Only validate fields that are being updated
	if req.APIKey != "" {
		validator.ValidateAPIKey(req.APIKey, credential.Exchange)
	}

	if req.APISecret != "" {
		validator.ValidateAPISecret(req.APISecret, credential.Exchange)
	}

	if req.Label != "" {
		validator.ValidateLabel(req.Label)
	}

	if validator.HasErrors() {
		apperror.WriteError(w, validator.ToAppError())
		return
	}

	// Update credential
	credential.Update(req.APIKey, req.APISecret, req.Label)

	// Save credential
	if err := h.useCase.UpdateCredential(r.Context(), credential); err != nil {
		h.logger.Error().Err(err).Str("id", id).Msg("Failed to update API credential")
		apperror.WriteError(w, apperror.NewInternal(err))
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"id":        credential.ID,
			"exchange":  credential.Exchange,
			"apiKey":    credential.APIKey,
			"label":     credential.Label,
			"updatedAt": credential.UpdatedAt,
		},
	}); err != nil {
		h.logger.Error().Err(err).Msg("Failed to encode response")
	}
}

// DeleteCredential deletes an API credential
func (h *APICredentialHandler) DeleteCredential(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		h.logger.Error().Msg("User ID not found in context")
		apperror.WriteError(w, apperror.NewUnauthorized("User ID not found in context", nil))
		return
	}

	// Get credential ID from URL
	id := chi.URLParam(r, "id")
	if id == "" {
		apperror.WriteError(w, apperror.NewInvalid("Credential ID is required", nil, nil))
		return
	}

	// Get credential
	credential, err := h.useCase.GetCredential(r.Context(), id)
	if err != nil {
		h.logger.Error().Err(err).Str("id", id).Msg("Failed to get API credential")
		apperror.WriteError(w, apperror.NewInternal(err))
		return
	}

	// Check if credential exists
	if credential == nil {
		apperror.WriteError(w, apperror.NewNotFound("api_credential", id, nil))
		return
	}

	// Check if credential belongs to user
	if credential.UserID != userID {
		apperror.WriteError(w, apperror.NewForbidden("API credential does not belong to user", nil))
		return
	}

	// Delete credential
	if err := h.useCase.DeleteCredential(r.Context(), id); err != nil {
		h.logger.Error().Err(err).Str("id", id).Msg("Failed to delete API credential")
		apperror.WriteError(w, apperror.NewInternal(err))
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
	}); err != nil {
		h.logger.Error().Err(err).Msg("Failed to encode response")
	}
}
