package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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
func (h *AccountHandler) GetAccount(c *gin.Context) {
	wallet, err := h.AccountService.GetWallet(c.Request.Context())
	if err != nil {
		h.logger.Error("Failed to get account info", zap.Error(err))
		c.JSON(http.StatusInternalServerError, responseDto.ErrorResponse{
			Code:    "GET_ACCOUNT_FAILED",
			Message: "Failed to get account info",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, wallet)
}

// GetBalances godoc
// @Summary Get account balances
// @Description Get account balances
// @Tags account
// @Produce json
// @Success 200 {object} interface{}
// @Failure 500 {object} responseDto.ErrorResponse
// @Router /api/v1/account/balance [get]
func (h *AccountHandler) GetBalances(c *gin.Context) {
	balance, err := h.AccountService.GetAccountBalance(c.Request.Context())
	if err != nil {
		h.logger.Error("Failed to get account balances", zap.Error(err))
		c.JSON(http.StatusInternalServerError, responseDto.ErrorResponse{
			Code:    "GET_BALANCES_FAILED",
			Message: "Failed to get account balances",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, balance)
}

// GetWallet godoc
// @Summary Get wallet information
// @Description Get wallet information
// @Tags account
// @Produce json
// @Success 200 {object} interface{}
// @Failure 500 {object} responseDto.ErrorResponse
// @Router /api/v1/account/wallet [get]
func (h *AccountHandler) GetWallet(c *gin.Context) {
	wallet, err := h.AccountService.GetWallet(c.Request.Context())
	if err != nil {
		h.logger.Error("Failed to get wallet", zap.Error(err))
		c.JSON(http.StatusInternalServerError, responseDto.ErrorResponse{
			Code:    "GET_WALLET_FAILED",
			Message: "Failed to get wallet",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, wallet)
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
func (h *AccountHandler) GetBalanceSummary(c *gin.Context) {
	// Parse days parameter
	daysStr := c.DefaultQuery("days", "30")
	days, err := strconv.Atoi(daysStr)
	if err != nil {
		h.logger.Error("Invalid days parameter", zap.Error(err))
		c.JSON(http.StatusBadRequest, responseDto.ErrorResponse{
			Code:    "INVALID_PARAMETER",
			Message: "Invalid days parameter",
			Details: err.Error(),
		})
		return
	}

	// Get balance summary
	summary, err := h.AccountService.GetBalanceSummary(c.Request.Context(), days)
	if err != nil {
		h.logger.Error("Failed to get balance summary", zap.Error(err))
		c.JSON(http.StatusInternalServerError, responseDto.ErrorResponse{
			Code:    "GET_BALANCE_SUMMARY_FAILED",
			Message: "Failed to get balance summary",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, summary)
}

// ValidateAPIKeys godoc
// @Summary Validate API keys
// @Description Validate the configured API keys
// @Tags account
// @Produce json
// @Success 200 {object} interface{}
// @Failure 500 {object} responseDto.ErrorResponse
// @Router /api/v1/account/validate-keys [get]
func (h *AccountHandler) ValidateAPIKeys(c *gin.Context) {
	valid, err := h.AccountService.ValidateAPIKeys(c.Request.Context())
	if err != nil {
		h.logger.Error("Failed to validate API keys", zap.Error(err))
		c.JSON(http.StatusInternalServerError, responseDto.ErrorResponse{
			Code:    "VALIDATE_KEYS_FAILED",
			Message: "Failed to validate API keys",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"valid": valid})
}

// SyncWithExchange godoc
// @Summary Sync with exchange
// @Description Sync account data with the exchange
// @Tags account
// @Produce json
// @Success 200 {object} interface{}
// @Failure 500 {object} responseDto.ErrorResponse
// @Router /api/v1/account/sync [post]
func (h *AccountHandler) SyncWithExchange(c *gin.Context) {
	err := h.AccountService.SyncWithExchange(c.Request.Context())
	if err != nil {
		h.logger.Error("Failed to sync with exchange", zap.Error(err))
		c.JSON(http.StatusInternalServerError, responseDto.ErrorResponse{
			Code:    "SYNC_FAILED",
			Message: "Failed to sync with exchange",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
