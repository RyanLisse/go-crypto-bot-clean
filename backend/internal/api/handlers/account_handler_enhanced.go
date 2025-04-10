package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	responseDto "go-crypto-bot-clean/backend/internal/api/dto/response"
	"go-crypto-bot-clean/backend/internal/domain/models"
)

// AccountServiceInterface defines the interface for account service
type AccountServiceInterface interface {
	GetAccountBalance(ctx context.Context) (models.Balance, error)
	GetWallet(ctx context.Context) (*models.Wallet, error)
	ValidateAPIKeys(ctx context.Context) (bool, error)
	GetCurrentExposure(ctx context.Context) (float64, error)
	GetListenKey(ctx context.Context) (string, error)
	RenewListenKey(ctx context.Context, listenKey string) error
	CloseListenKey(ctx context.Context, listenKey string) error
}

// EnhancedAccountHandler handles enhanced account-related API endpoints
type EnhancedAccountHandler struct {
	accountService AccountServiceInterface
}

// NewEnhancedAccountHandler creates a new enhanced account handler
func NewEnhancedAccountHandler(accountService AccountServiceInterface) *EnhancedAccountHandler {
	return &EnhancedAccountHandler{
		accountService: accountService,
	}
}

// GetAccountDetails godoc
// @Summary Get detailed account information
// @Description Returns detailed account information including balances and exposure
// @Tags account
// @Accept json
// @Produce json
// @Success 200 {object} responseDto.AccountResponse
// @Failure 500 {object} responseDto.ErrorResponse
// @Router /api/v1/account/details [get]
func (h *EnhancedAccountHandler) GetAccountDetails(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get account balance
	balance, err := h.accountService.GetAccountBalance(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(responseDto.ErrorResponse{
			Code:    "internal_error",
			Message: "Failed to get account balance",
			Details: err.Error(),
		})
		return
	}

	// Get wallet
	wallet, err := h.accountService.GetWallet(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(responseDto.ErrorResponse{
			Code:    "internal_error",
			Message: "Failed to get wallet",
			Details: err.Error(),
		})
		return
	}

	// Get current exposure
	exposure, err := h.accountService.GetCurrentExposure(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(responseDto.ErrorResponse{
			Code:    "internal_error",
			Message: "Failed to get current exposure",
			Details: err.Error(),
		})
		return
	}

	// Build response
	resp := responseDto.AccountResponse{
		TotalBalance:    balance.Fiat,
		AvailableFunds:  balance.Fiat - exposure,
		CurrentExposure: exposure,
		Assets:          mapToAssetResponses(wallet),
		Timestamp:       time.Now(),
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// ValidateAPIKeys godoc
// @Summary Validate API keys
// @Description Validates the API keys configured in the system
// @Tags account
// @Accept json
// @Produce json
// @Success 200 {object} responseDto.APIKeyValidationResponse
// @Failure 500 {object} responseDto.ErrorResponse
// @Router /api/v1/account/validate-keys [get]
func (h *EnhancedAccountHandler) ValidateAPIKeys(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Validate API keys
	valid, err := h.accountService.ValidateAPIKeys(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(responseDto.ErrorResponse{
			Code:    "internal_error",
			Message: "Failed to validate API keys",
			Details: err.Error(),
		})
		return
	}

	// Build response
	resp := responseDto.APIKeyValidationResponse{
		Valid:     valid,
		Timestamp: time.Now(),
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// GetListenKey godoc
// @Summary Get a listen key for WebSocket authentication
// @Description Returns a listen key that can be used to authenticate WebSocket connections
// @Tags account
// @Accept json
// @Produce json
// @Success 200 {object} responseDto.ListenKeyResponse
// @Failure 500 {object} responseDto.ErrorResponse
// @Router /api/v1/account/listen-key [get]
func (h *EnhancedAccountHandler) GetListenKey(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get listen key
	listenKey, err := h.accountService.GetListenKey(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(responseDto.ErrorResponse{
			Code:    "internal_error",
			Message: "Failed to get listen key",
			Details: err.Error(),
		})
		return
	}

	// Build response
	resp := responseDto.ListenKeyResponse{
		ListenKey: listenKey,
		Expires:   time.Now().Add(60 * time.Minute), // Listen keys typically expire after 60 minutes
		Timestamp: time.Now(),
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// RenewListenKey godoc
// @Summary Renew a listen key
// @Description Renews a listen key to extend its validity
// @Tags account
// @Accept json
// @Produce json
// @Param listen_key query string true "Listen key to renew"
// @Success 200 {object} responseDto.ListenKeyResponse
// @Failure 400 {object} responseDto.ErrorResponse
// @Failure 500 {object} responseDto.ErrorResponse
// @Router /api/v1/account/listen-key/renew [put]
func (h *EnhancedAccountHandler) RenewListenKey(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get listen key from query parameter
	listenKey := r.URL.Query().Get("listen_key")
	if listenKey == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(responseDto.ErrorResponse{
			Code:    "missing_parameter",
			Message: "Missing listen_key parameter",
		})
		return
	}

	// Renew listen key
	err := h.accountService.RenewListenKey(ctx, listenKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(responseDto.ErrorResponse{
			Code:    "internal_error",
			Message: "Failed to renew listen key",
			Details: err.Error(),
		})
		return
	}

	// Build response
	resp := responseDto.ListenKeyResponse{
		ListenKey: listenKey,
		Expires:   time.Now().Add(60 * time.Minute), // Listen keys typically expire after 60 minutes
		Timestamp: time.Now(),
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// CloseListenKey godoc
// @Summary Close a listen key
// @Description Closes a listen key to invalidate it
// @Tags account
// @Accept json
// @Produce json
// @Param listen_key query string true "Listen key to close"
// @Success 200 {object} responseDto.SuccessResponse
// @Failure 400 {object} responseDto.ErrorResponse
// @Failure 500 {object} responseDto.ErrorResponse
// @Router /api/v1/account/listen-key/close [delete]
func (h *EnhancedAccountHandler) CloseListenKey(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get listen key from query parameter
	listenKey := r.URL.Query().Get("listen_key")
	if listenKey == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(responseDto.ErrorResponse{
			Code:    "missing_parameter",
			Message: "Missing listen_key parameter",
		})
		return
	}

	// Close listen key
	err := h.accountService.CloseListenKey(ctx, listenKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(responseDto.ErrorResponse{
			Code:    "internal_error",
			Message: "Failed to close listen key",
			Details: err.Error(),
		})
		return
	}

	// Build response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responseDto.SuccessResponse{
		Message:   "Listen key closed successfully",
		Timestamp: time.Now(),
	})
}

// Helper functions to map domain models to DTOs
func mapToAssetResponses(wallet *models.Wallet) []responseDto.AssetResponse {
	assets := make([]responseDto.AssetResponse, 0, len(wallet.Balances))
	for symbol, balance := range wallet.Balances {
		if balance.Total > 0 {
			assets = append(assets, responseDto.AssetResponse{
				Symbol:    symbol,
				Free:      balance.Free,
				Locked:    balance.Locked,
				Total:     balance.Total,
				Price:     balance.Price,
				ValueUSDT: balance.Total * balance.Price,
			})
		}
	}
	return assets
}
