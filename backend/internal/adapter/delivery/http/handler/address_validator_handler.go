package handler

import (
	"encoding/json"
	"net/http"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/apperror"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/usecase"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
)

// AddressValidatorHandler handles address validation HTTP requests
type AddressValidatorHandler struct {
	addressValidatorService usecase.AddressValidatorService
	logger                  *zerolog.Logger
}

// NewAddressValidatorHandler creates a new AddressValidatorHandler
func NewAddressValidatorHandler(
	addressValidatorService usecase.AddressValidatorService,
	logger *zerolog.Logger,
) *AddressValidatorHandler {
	return &AddressValidatorHandler{
		addressValidatorService: addressValidatorService,
		logger:                  logger,
	}
}

// RegisterRoutes registers the address validator routes
func (h *AddressValidatorHandler) RegisterRoutes(r chi.Router) {
	r.Route("/address-validator", func(r chi.Router) {
		r.Get("/networks", h.GetSupportedNetworks)
		r.Post("/validate", h.ValidateAddress)
		r.Post("/info", h.GetAddressInfo)
	})
}

// ValidateAddress handles the validate address endpoint
func (h *AddressValidatorHandler) ValidateAddress(w http.ResponseWriter, r *http.Request) {
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
	valid, err := h.addressValidatorService.ValidateAddress(r.Context(), request.Network, request.Address)
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

// GetAddressInfo handles the get address info endpoint
func (h *AddressValidatorHandler) GetAddressInfo(w http.ResponseWriter, r *http.Request) {
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

	// Get address info
	info, err := h.addressValidatorService.GetAddressInfo(r.Context(), request.Network, request.Address)
	if err != nil {
		h.logger.Error().Err(err).Str("network", request.Network).Str("address", request.Address).Msg("Failed to get address info")
		apperror.WriteError(w, apperror.NewInternal(err))
		return
	}

	// Write response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(info)
}

// GetSupportedNetworks handles the get supported networks endpoint
func (h *AddressValidatorHandler) GetSupportedNetworks(w http.ResponseWriter, r *http.Request) {
	// Get supported networks
	networks, err := h.addressValidatorService.GetSupportedNetworks(r.Context())
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
