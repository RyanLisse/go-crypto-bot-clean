package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	responseDto "go-crypto-bot-clean/backend/internal/api/dto/response"
	"go-crypto-bot-clean/backend/internal/core/account"
	"go-crypto-bot-clean/backend/internal/domain/service"

	"go.uber.org/zap"
)

// AccountHandler handles account-related endpoints
type AccountHandler struct {
	ExchangeService service.ExchangeService
	AccountService  account.AccountService
	logger          *zap.Logger
}

// NewAccountHandler creates a new AccountHandler
func NewAccountHandler(exchangeService service.ExchangeService, accountService account.AccountService, logger *zap.Logger) *AccountHandler {
	return &AccountHandler{
		ExchangeService: exchangeService,
		AccountService:  accountService,
		logger:          logger,
	}
}

// GetAccount godoc
// @Summary Get account information
// @Description Get account information
// @Tags account
// @Produce json
// @Success 200 {object} interface{}
// @Failure 500 {object} responseDto.ErrorResponse
// @Router /api/v1/account [get]
func (h *AccountHandler) GetAccount(w http.ResponseWriter, r *http.Request) {
	wallet, err := h.AccountService.GetWallet(r.Context())
	if err != nil {
		h.logger.Error("Failed to get account info", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(responseDto.ErrorResponse{
			Code:    "GET_ACCOUNT_FAILED",
			Message: "Failed to get account info",
			Details: err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(wallet)
}

// GetBalances godoc
// @Summary Get account balances
// @Description Get account balances
// @Tags account
// @Produce json
// @Success 200 {object} interface{}
// @Failure 500 {object} responseDto.ErrorResponse
// @Router /api/v1/account/balance [get]
func (h *AccountHandler) GetBalances(w http.ResponseWriter, r *http.Request) {
	balance, err := h.AccountService.GetAccountBalance(r.Context())
	if err != nil {
		h.logger.Error("Failed to get account balances", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(responseDto.ErrorResponse{
			Code:    "GET_BALANCES_FAILED",
			Message: "Failed to get account balances",
			Details: err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(balance)
}

// GetWallet godoc
// @Summary Get wallet information
// @Description Get wallet information
// @Tags account
// @Produce json
// @Success 200 {object} interface{}
// @Failure 500 {object} responseDto.ErrorResponse
// @Router /api/v1/account/wallet [get]
func (h *AccountHandler) GetWallet(w http.ResponseWriter, r *http.Request) {
	wallet, err := h.AccountService.GetWallet(r.Context())
	if err != nil {
		h.logger.Error("Failed to get wallet", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(responseDto.ErrorResponse{
			Code:    "GET_WALLET_FAILED",
			Message: "Failed to get wallet",
			Details: err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(wallet)
}

// GetBalanceSummary godoc
// @Summary Get balance summary
// @Description Get balance summary for a specified period
// @Tags account
// @Produce json
// @Param days query int false "Number of days to include in summary (default: 30)"
// @Success 200 {object} interface{}
// @Failure 400 {object} responseDto.ErrorResponse
// @Failure 500 {object} responseDto.ErrorResponse
// @Router /api/v1/account/balance-summary [get]
func (h *AccountHandler) GetBalanceSummary(w http.ResponseWriter, r *http.Request) {
	// Parse days parameter
	daysStr := r.URL.Query().Get("days")
	if daysStr == "" {
		daysStr = "30"
	}
	days, err := strconv.Atoi(daysStr)
	if err != nil {
		h.logger.Error("Invalid days parameter", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(responseDto.ErrorResponse{
			Code:    "INVALID_PARAMETER",
			Message: "Invalid days parameter",
			Details: err.Error(),
		})
		return
	}

	// Get balance summary
	summary, err := h.AccountService.GetBalanceSummary(r.Context(), days)
	if err != nil {
		h.logger.Error("Failed to get balance summary", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(responseDto.ErrorResponse{
			Code:    "GET_BALANCE_SUMMARY_FAILED",
			Message: "Failed to get balance summary",
			Details: err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(summary)
}

// ValidateAPIKeys godoc
// @Summary Validate API keys
// @Description Validate the configured API keys
// @Tags account
// @Produce json
// @Success 200 {object} interface{}
// @Failure 500 {object} responseDto.ErrorResponse
// @Router /api/v1/account/validate-keys [get]
func (h *AccountHandler) ValidateAPIKeys(w http.ResponseWriter, r *http.Request) {
	valid, err := h.AccountService.ValidateAPIKeys(r.Context())
	if err != nil {
		h.logger.Error("Failed to validate API keys", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(responseDto.ErrorResponse{
			Code:    "VALIDATE_KEYS_FAILED",
			Message: "Failed to validate API keys",
			Details: err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"valid": valid})
}

// SyncWithExchange godoc
// @Summary Sync with exchange
// @Description Sync account data with the exchange
// @Tags account
// @Produce json
// @Success 200 {object} interface{}
// @Failure 500 {object} responseDto.ErrorResponse
// @Router /api/v1/account/sync [post]
func (h *AccountHandler) SyncWithExchange(w http.ResponseWriter, r *http.Request) {
	err := h.AccountService.SyncWithExchange(r.Context())
	if err != nil {
		h.logger.Error("Failed to sync with exchange", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(responseDto.ErrorResponse{
			Code:    "SYNC_FAILED",
			Message: "Failed to sync with exchange",
			Details: err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"status": "success"})
}
