package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go-crypto-bot-clean/backend/internal/api/dto/request"
	"go-crypto-bot-clean/backend/internal/api/dto/response"
	"go-crypto-bot-clean/backend/internal/config"
)

// ConfigHandler handles configuration-related endpoints
type ConfigHandler struct {
	Config *config.Config
}

// NewConfigHandler creates a new ConfigHandler
func NewConfigHandler(cfg *config.Config) *ConfigHandler {
	return &ConfigHandler{
		Config: cfg,
	}
}

// GetCurrentConfig godoc
// @Summary Get current configuration
// @Description Returns the current bot configuration
// @Tags config
// @Accept json
// @Produce json
// @Success 200 {object} response.ConfigResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/config [get]
func (h *ConfigHandler) GetCurrentConfig(c *gin.Context) {
	// Build response
	resp := response.ConfigResponse{
		USDTPerTrade:     h.Config.Trading.DefaultQuantity,
		StopLossPercent:  h.Config.Trading.StopLossPercent,
		TakeProfitLevels: h.Config.Trading.TakeProfitLevels,
		SellPercentages:  h.Config.Trading.SellPercentages,
		UpdatedAt:        time.Now(),
	}

	c.JSON(http.StatusOK, resp)
}

// UpdateConfig godoc
// @Summary Update configuration
// @Description Updates the bot configuration
// @Tags config
// @Accept json
// @Produce json
// @Param config body request.ConfigUpdateRequest true "Configuration to update"
// @Success 200 {object} response.ConfigResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/config [put]
func (h *ConfigHandler) UpdateConfig(c *gin.Context) {
	// Parse request
	var req request.ConfigUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Code:    "invalid_request",
			Message: "Invalid request format",
			Details: err.Error(),
		})
		return
	}

	// Update config
	if req.USDTPerTrade != nil {
		h.Config.Trading.DefaultQuantity = *req.USDTPerTrade
	}
	if req.StopLossPercent != nil {
		h.Config.Trading.StopLossPercent = *req.StopLossPercent
	}
	if req.TakeProfitLevels != nil {
		h.Config.Trading.TakeProfitLevels = req.TakeProfitLevels
	}
	if req.SellPercentages != nil {
		h.Config.Trading.SellPercentages = req.SellPercentages
	}

	// Build response
	resp := response.ConfigResponse{
		USDTPerTrade:     h.Config.Trading.DefaultQuantity,
		StopLossPercent:  h.Config.Trading.StopLossPercent,
		TakeProfitLevels: h.Config.Trading.TakeProfitLevels,
		SellPercentages:  h.Config.Trading.SellPercentages,
		UpdatedAt:        time.Now(),
	}

	c.JSON(http.StatusOK, resp)
}

// GetDefaultConfig godoc
// @Summary Get default configuration
// @Description Returns the default bot configuration
// @Tags config
// @Accept json
// @Produce json
// @Success 200 {object} response.ConfigResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/config/defaults [get]
func (h *ConfigHandler) GetDefaultConfig(c *gin.Context) {
	// Build response with default values
	resp := response.ConfigResponse{
		USDTPerTrade:     20.0,
		StopLossPercent:  10.0,
		TakeProfitLevels: []float64{5.0, 10.0, 15.0, 20.0},
		SellPercentages:  []float64{0.25, 0.25, 0.25, 0.25},
		UpdatedAt:        time.Now(),
	}

	c.JSON(http.StatusOK, resp)
}
