package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/neo/crypto-bot/internal/adapter/http/response"
	"github.com/neo/crypto-bot/internal/domain/model"
	"github.com/neo/crypto-bot/internal/usecase"
	"github.com/rs/zerolog"
)

// PositionHandler handles HTTP requests for position management
type PositionHandler struct {
	useCase usecase.PositionUseCase
	logger  *zerolog.Logger
}

// NewPositionHandler creates a new PositionHandler
func NewPositionHandler(uc usecase.PositionUseCase, logger *zerolog.Logger) *PositionHandler {
	return &PositionHandler{
		useCase: uc,
		logger:  logger,
	}
}

// RegisterRoutes registers position-related routes with the Gin engine
func (h *PositionHandler) RegisterRoutes(router *gin.RouterGroup) {
	positionGroup := router.Group("/positions")
	{
		positionGroup.POST("", h.CreatePosition)
		positionGroup.GET("", h.GetPositions)
		positionGroup.GET("/open", h.GetOpenPositions)
		positionGroup.GET("/closed", h.GetClosedPositions)
		positionGroup.GET("/symbol/:symbol", h.GetPositionsBySymbol)
		positionGroup.GET("/type/:type", h.GetPositionsByType)
		positionGroup.GET("/:id", h.GetPositionByID)
		positionGroup.PUT("/:id", h.UpdatePosition)
		positionGroup.PUT("/:id/price", h.UpdatePositionPrice)
		positionGroup.PUT("/:id/close", h.ClosePosition)
		positionGroup.PUT("/:id/stop-loss", h.SetStopLoss)
		positionGroup.PUT("/:id/take-profit", h.SetTakeProfit)
		positionGroup.DELETE("/:id", h.DeletePosition)
	}
}

// CreatePosition godoc
// @Summary Create a new position
// @Description Create a new trading position
// @Tags positions
// @Accept json
// @Produce json
// @Param position body model.PositionCreateRequest true "Position data"
// @Success 201 {object} response.APIResponse{data=model.Position}
// @Failure 400 {object} response.APIResponse{error=response.APIError}
// @Failure 500 {object} response.APIResponse{error=response.APIError}
// @Router /positions [post]
func (h *PositionHandler) CreatePosition(c *gin.Context) {
	var req model.PositionCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error().Err(err).Msg("Invalid position create request")
		c.JSON(http.StatusBadRequest, response.Error(response.ErrorCodeBadRequest, "Invalid position data: "+err.Error()))
		return
	}

	position, err := h.useCase.CreatePosition(c.Request.Context(), req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		errCode := response.ErrorCodeInternalError
		errMsg := "Failed to create position"

		if err == usecase.ErrSymbolNotFound {
			statusCode = http.StatusBadRequest
			errCode = response.ErrorCodeBadRequest
			errMsg = "Symbol not found"
		} else if err == usecase.ErrInvalidPositionData {
			statusCode = http.StatusBadRequest
			errCode = response.ErrorCodeBadRequest
			errMsg = "Invalid position data"
		}

		h.logger.Error().Err(err).Msg(errMsg)
		c.JSON(statusCode, response.Error(errCode, errMsg))
		return
	}

	c.JSON(http.StatusCreated, response.Success(position))
}

// GetPositions godoc
// @Summary Get all positions
// @Description Returns a list of all positions with optional pagination
// @Tags positions
// @Accept json
// @Produce json
// @Param limit query int false "Limit number of results (default: 50)"
// @Param offset query int false "Offset for pagination (default: 0)"
// @Success 200 {object} response.APIResponse{data=[]model.Position}
// @Failure 500 {object} response.APIResponse{error=response.APIError}
// @Router /positions [get]
func (h *PositionHandler) GetPositions(c *gin.Context) {
	// This is a placeholder for getting all positions
	// In a real implementation, this would need pagination and filtering
	openPositions, err := h.useCase.GetOpenPositions(c.Request.Context())
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get positions")
		c.JSON(http.StatusInternalServerError, response.Error(response.ErrorCodeInternalError, "Failed to get positions"))
		return
	}

	c.JSON(http.StatusOK, response.Success(openPositions))
}

// GetOpenPositions godoc
// @Summary Get open positions
// @Description Get all currently open positions
// @Tags positions
// @Accept json
// @Produce json
// @Success 200 {object} response.APIResponse{data=[]model.Position}
// @Failure 500 {object} response.APIResponse{error=response.APIError}
// @Router /positions/open [get]
func (h *PositionHandler) GetOpenPositions(c *gin.Context) {
	positions, err := h.useCase.GetOpenPositions(c.Request.Context())
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get open positions")
		c.JSON(http.StatusInternalServerError, response.Error(response.ErrorCodeInternalError, "Failed to get open positions"))
		return
	}

	c.JSON(http.StatusOK, response.Success(positions))
}

// GetClosedPositions godoc
// @Summary Get closed positions
// @Description Get closed positions with optional time range and pagination
// @Tags positions
// @Accept json
// @Produce json
// @Param from query string false "Start date (RFC3339 format)"
// @Param to query string false "End date (RFC3339 format)"
// @Param limit query int false "Limit number of results (default: 50)"
// @Param offset query int false "Offset for pagination (default: 0)"
// @Success 200 {object} response.APIResponse{data=[]model.Position}
// @Failure 400 {object} response.APIResponse{error=response.APIError}
// @Failure 500 {object} response.APIResponse{error=response.APIError}
// @Router /positions/closed [get]
func (h *PositionHandler) GetClosedPositions(c *gin.Context) {
	// Default time range: last 30 days
	toTime := time.Now()
	fromTime := toTime.AddDate(0, -1, 0) // 1 month ago

	// Parse from date if provided
	if fromStr := c.Query("from"); fromStr != "" {
		parsedFrom, err := time.Parse(time.RFC3339, fromStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, response.Error(response.ErrorCodeBadRequest, "Invalid 'from' date format. Use RFC3339 format (e.g., 2023-01-01T00:00:00Z)"))
			return
		}
		fromTime = parsedFrom
	}

	// Parse to date if provided
	if toStr := c.Query("to"); toStr != "" {
		parsedTo, err := time.Parse(time.RFC3339, toStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, response.Error(response.ErrorCodeBadRequest, "Invalid 'to' date format. Use RFC3339 format (e.g., 2023-01-01T00:00:00Z)"))
			return
		}
		toTime = parsedTo
	}

	// Parse pagination parameters
	limit := 50 // Default limit
	if limitStr := c.Query("limit"); limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err != nil || parsedLimit <= 0 {
			c.JSON(http.StatusBadRequest, response.Error(response.ErrorCodeBadRequest, "Invalid limit parameter"))
			return
		}
		limit = parsedLimit
	}

	offset := 0 // Default offset
	if offsetStr := c.Query("offset"); offsetStr != "" {
		parsedOffset, err := strconv.Atoi(offsetStr)
		if err != nil || parsedOffset < 0 {
			c.JSON(http.StatusBadRequest, response.Error(response.ErrorCodeBadRequest, "Invalid offset parameter"))
			return
		}
		offset = parsedOffset
	}

	positions, err := h.useCase.GetClosedPositions(c.Request.Context(), fromTime, toTime, limit, offset)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get closed positions")
		c.JSON(http.StatusInternalServerError, response.Error(response.ErrorCodeInternalError, "Failed to get closed positions"))
		return
	}

	c.JSON(http.StatusOK, response.Success(positions))
}

// GetPositionsBySymbol godoc
// @Summary Get positions by symbol
// @Description Get positions for a specific trading symbol
// @Tags positions
// @Accept json
// @Produce json
// @Param symbol path string true "Trading symbol (e.g., BTCUSDT)"
// @Param limit query int false "Limit number of results (default: 50)"
// @Param offset query int false "Offset for pagination (default: 0)"
// @Success 200 {object} response.APIResponse{data=[]model.Position}
// @Failure 400 {object} response.APIResponse{error=response.APIError}
// @Failure 500 {object} response.APIResponse{error=response.APIError}
// @Router /positions/symbol/{symbol} [get]
func (h *PositionHandler) GetPositionsBySymbol(c *gin.Context) {
	symbol := c.Param("symbol")
	if symbol == "" {
		c.JSON(http.StatusBadRequest, response.Error(response.ErrorCodeBadRequest, "Symbol parameter is required"))
		return
	}

	// Parse pagination parameters
	limit := 50 // Default limit
	if limitStr := c.Query("limit"); limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err != nil || parsedLimit <= 0 {
			c.JSON(http.StatusBadRequest, response.Error(response.ErrorCodeBadRequest, "Invalid limit parameter"))
			return
		}
		limit = parsedLimit
	}

	offset := 0 // Default offset
	if offsetStr := c.Query("offset"); offsetStr != "" {
		parsedOffset, err := strconv.Atoi(offsetStr)
		if err != nil || parsedOffset < 0 {
			c.JSON(http.StatusBadRequest, response.Error(response.ErrorCodeBadRequest, "Invalid offset parameter"))
			return
		}
		offset = parsedOffset
	}

	positions, err := h.useCase.GetPositionsBySymbol(c.Request.Context(), symbol, limit, offset)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get positions by symbol")
		c.JSON(http.StatusInternalServerError, response.Error(response.ErrorCodeInternalError, "Failed to get positions by symbol"))
		return
	}

	c.JSON(http.StatusOK, response.Success(positions))
}

// GetPositionsByType godoc
// @Summary Get positions by type
// @Description Get open positions for a specific type
// @Tags positions
// @Accept json
// @Produce json
// @Param type path string true "Position type (MANUAL, AUTOMATIC, NEWCOIN)"
// @Success 200 {object} response.APIResponse{data=[]model.Position}
// @Failure 400 {object} response.APIResponse{error=response.APIError}
// @Failure 500 {object} response.APIResponse{error=response.APIError}
// @Router /positions/type/{type} [get]
func (h *PositionHandler) GetPositionsByType(c *gin.Context) {
	positionTypeStr := c.Param("type")
	if positionTypeStr == "" {
		c.JSON(http.StatusBadRequest, response.Error(response.ErrorCodeBadRequest, "Type parameter is required"))
		return
	}

	// Validate position type
	positionType := model.PositionType(positionTypeStr)
	if positionType != model.PositionTypeManual &&
		positionType != model.PositionTypeAutomatic &&
		positionType != model.PositionTypeNewCoin {
		c.JSON(http.StatusBadRequest, response.Error(response.ErrorCodeBadRequest, "Invalid position type. Must be one of: MANUAL, AUTOMATIC, NEWCOIN"))
		return
	}

	positions, err := h.useCase.GetOpenPositionsByType(c.Request.Context(), positionType)
	if err != nil {
		h.logger.Error().Err(err).Str("type", positionTypeStr).Msg("Failed to get positions by type")
		c.JSON(http.StatusInternalServerError, response.Error(response.ErrorCodeInternalError, "Failed to get positions by type"))
		return
	}

	c.JSON(http.StatusOK, response.Success(positions))
}

// GetPositionByID godoc
// @Summary Get position by ID
// @Description Get a specific position by its ID
// @Tags positions
// @Accept json
// @Produce json
// @Param id path string true "Position ID"
// @Success 200 {object} response.APIResponse{data=model.Position}
// @Failure 404 {object} response.APIResponse{error=response.APIError}
// @Failure 500 {object} response.APIResponse{error=response.APIError}
// @Router /positions/{id} [get]
func (h *PositionHandler) GetPositionByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, response.Error(response.ErrorCodeBadRequest, "ID parameter is required"))
		return
	}

	position, err := h.useCase.GetPositionByID(c.Request.Context(), id)
	if err != nil {
		h.logger.Error().Err(err).Str("id", id).Msg("Failed to get position by ID")
		if err == usecase.ErrPositionNotFound {
			c.JSON(http.StatusNotFound, response.Error(response.ErrorCodeNotFound, "Position not found"))
		} else {
			c.JSON(http.StatusInternalServerError, response.Error(response.ErrorCodeInternalError, "Failed to get position"))
		}
		return
	}

	if position == nil {
		c.JSON(http.StatusNotFound, response.Error(response.ErrorCodeNotFound, "Position not found"))
		return
	}

	c.JSON(http.StatusOK, response.Success(position))
}

// UpdatePosition godoc
// @Summary Update position
// @Description Update an existing position
// @Tags positions
// @Accept json
// @Produce json
// @Param id path string true "Position ID"
// @Param position body model.PositionUpdateRequest true "Position update data"
// @Success 200 {object} response.APIResponse{data=model.Position}
// @Failure 400 {object} response.APIResponse{error=response.APIError}
// @Failure 404 {object} response.APIResponse{error=response.APIError}
// @Failure 500 {object} response.APIResponse{error=response.APIError}
// @Router /positions/{id} [put]
func (h *PositionHandler) UpdatePosition(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, response.Error(response.ErrorCodeBadRequest, "ID parameter is required"))
		return
	}

	var req model.PositionUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error().Err(err).Msg("Invalid position update request")
		c.JSON(http.StatusBadRequest, response.Error(response.ErrorCodeBadRequest, "Invalid position update data: "+err.Error()))
		return
	}

	position, err := h.useCase.UpdatePosition(c.Request.Context(), id, req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		errCode := response.ErrorCodeInternalError
		errMsg := "Failed to update position"

		if err == usecase.ErrPositionNotFound {
			statusCode = http.StatusNotFound
			errCode = response.ErrorCodeNotFound
			errMsg = "Position not found"
		} else if err == usecase.ErrInvalidPositionData {
			statusCode = http.StatusBadRequest
			errCode = response.ErrorCodeBadRequest
			errMsg = "Invalid position data"
		}

		h.logger.Error().Err(err).Str("id", id).Msg(errMsg)
		c.JSON(statusCode, response.Error(errCode, errMsg))
		return
	}

	c.JSON(http.StatusOK, response.Success(position))
}

// UpdatePositionPrice godoc
// @Summary Update position price
// @Description Update the current price of a position
// @Tags positions
// @Accept json
// @Produce json
// @Param id path string true "Position ID"
// @Param price body struct{Price float64 `json:"price"`} true "Current price"
// @Success 200 {object} response.APIResponse{data=model.Position}
// @Failure 400 {object} response.APIResponse{error=response.APIError}
// @Failure 404 {object} response.APIResponse{error=response.APIError}
// @Failure 500 {object} response.APIResponse{error=response.APIError}
// @Router /positions/{id}/price [put]
func (h *PositionHandler) UpdatePositionPrice(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, response.Error(response.ErrorCodeBadRequest, "ID parameter is required"))
		return
	}

	var req struct {
		Price float64 `json:"price" binding:"required,gt=0"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error().Err(err).Msg("Invalid price update request")
		c.JSON(http.StatusBadRequest, response.Error(response.ErrorCodeBadRequest, "Invalid price data: "+err.Error()))
		return
	}

	position, err := h.useCase.UpdatePositionPrice(c.Request.Context(), id, req.Price)
	if err != nil {
		statusCode := http.StatusInternalServerError
		errCode := response.ErrorCodeInternalError
		errMsg := "Failed to update position price"

		if err == usecase.ErrPositionNotFound {
			statusCode = http.StatusNotFound
			errCode = response.ErrorCodeNotFound
			errMsg = "Position not found"
		}

		h.logger.Error().Err(err).Str("id", id).Msg(errMsg)
		c.JSON(statusCode, response.Error(errCode, errMsg))
		return
	}

	c.JSON(http.StatusOK, response.Success(position))
}

// ClosePosition godoc
// @Summary Close position
// @Description Close an open position
// @Tags positions
// @Accept json
// @Produce json
// @Param id path string true "Position ID"
// @Param data body struct{ExitPrice float64 `json:"exitPrice"`;ExitOrderIDs []string `json:"exitOrderIds"`} true "Closing data"
// @Success 200 {object} response.APIResponse{data=model.Position}
// @Failure 400 {object} response.APIResponse{error=response.APIError}
// @Failure 404 {object} response.APIResponse{error=response.APIError}
// @Failure 500 {object} response.APIResponse{error=response.APIError}
// @Router /positions/{id}/close [put]
func (h *PositionHandler) ClosePosition(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, response.Error(response.ErrorCodeBadRequest, "ID parameter is required"))
		return
	}

	var req struct {
		ExitPrice    float64  `json:"exitPrice" binding:"required,gt=0"`
		ExitOrderIDs []string `json:"exitOrderIds" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error().Err(err).Msg("Invalid close position request")
		c.JSON(http.StatusBadRequest, response.Error(response.ErrorCodeBadRequest, "Invalid close position data: "+err.Error()))
		return
	}

	position, err := h.useCase.ClosePosition(c.Request.Context(), id, req.ExitPrice, req.ExitOrderIDs)
	if err != nil {
		statusCode := http.StatusInternalServerError
		errCode := response.ErrorCodeInternalError
		errMsg := "Failed to close position"

		if err == usecase.ErrPositionNotFound {
			statusCode = http.StatusNotFound
			errCode = response.ErrorCodeNotFound
			errMsg = "Position not found"
		}

		h.logger.Error().Err(err).Str("id", id).Msg(errMsg)
		c.JSON(statusCode, response.Error(errCode, errMsg))
		return
	}

	c.JSON(http.StatusOK, response.Success(position))
}

// SetStopLoss godoc
// @Summary Set stop loss
// @Description Set or update the stop loss for a position
// @Tags positions
// @Accept json
// @Produce json
// @Param id path string true "Position ID"
// @Param data body struct{StopLoss float64 `json:"stopLoss"`} true "Stop loss price"
// @Success 200 {object} response.APIResponse{data=model.Position}
// @Failure 400 {object} response.APIResponse{error=response.APIError}
// @Failure 404 {object} response.APIResponse{error=response.APIError}
// @Failure 500 {object} response.APIResponse{error=response.APIError}
// @Router /positions/{id}/stop-loss [put]
func (h *PositionHandler) SetStopLoss(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, response.Error(response.ErrorCodeBadRequest, "ID parameter is required"))
		return
	}

	var req struct {
		StopLoss float64 `json:"stopLoss" binding:"required,gt=0"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error().Err(err).Msg("Invalid stop loss request")
		c.JSON(http.StatusBadRequest, response.Error(response.ErrorCodeBadRequest, "Invalid stop loss data: "+err.Error()))
		return
	}

	position, err := h.useCase.SetStopLoss(c.Request.Context(), id, req.StopLoss)
	if err != nil {
		statusCode := http.StatusInternalServerError
		errCode := response.ErrorCodeInternalError
		errMsg := "Failed to set stop loss"

		if err == usecase.ErrPositionNotFound {
			statusCode = http.StatusNotFound
			errCode = response.ErrorCodeNotFound
			errMsg = "Position not found"
		}

		h.logger.Error().Err(err).Str("id", id).Msg(errMsg)
		c.JSON(statusCode, response.Error(errCode, errMsg))
		return
	}

	c.JSON(http.StatusOK, response.Success(position))
}

// SetTakeProfit godoc
// @Summary Set take profit
// @Description Set or update the take profit for a position
// @Tags positions
// @Accept json
// @Produce json
// @Param id path string true "Position ID"
// @Param data body struct{TakeProfit float64 `json:"takeProfit"`} true "Take profit price"
// @Success 200 {object} response.APIResponse{data=model.Position}
// @Failure 400 {object} response.APIResponse{error=response.APIError}
// @Failure 404 {object} response.APIResponse{error=response.APIError}
// @Failure 500 {object} response.APIResponse{error=response.APIError}
// @Router /positions/{id}/take-profit [put]
func (h *PositionHandler) SetTakeProfit(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, response.Error(response.ErrorCodeBadRequest, "ID parameter is required"))
		return
	}

	var req struct {
		TakeProfit float64 `json:"takeProfit" binding:"required,gt=0"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error().Err(err).Msg("Invalid take profit request")
		c.JSON(http.StatusBadRequest, response.Error(response.ErrorCodeBadRequest, "Invalid take profit data: "+err.Error()))
		return
	}

	position, err := h.useCase.SetTakeProfit(c.Request.Context(), id, req.TakeProfit)
	if err != nil {
		statusCode := http.StatusInternalServerError
		errCode := response.ErrorCodeInternalError
		errMsg := "Failed to set take profit"

		if err == usecase.ErrPositionNotFound {
			statusCode = http.StatusNotFound
			errCode = response.ErrorCodeNotFound
			errMsg = "Position not found"
		}

		h.logger.Error().Err(err).Str("id", id).Msg(errMsg)
		c.JSON(statusCode, response.Error(errCode, errMsg))
		return
	}

	c.JSON(http.StatusOK, response.Success(position))
}

// DeletePosition godoc
// @Summary Delete position
// @Description Delete a position
// @Tags positions
// @Accept json
// @Produce json
// @Param id path string true "Position ID"
// @Success 204 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse{error=response.APIError}
// @Failure 500 {object} response.APIResponse{error=response.APIError}
// @Router /positions/{id} [delete]
func (h *PositionHandler) DeletePosition(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, response.Error(response.ErrorCodeBadRequest, "ID parameter is required"))
		return
	}

	err := h.useCase.DeletePosition(c.Request.Context(), id)
	if err != nil {
		statusCode := http.StatusInternalServerError
		errCode := response.ErrorCodeInternalError
		errMsg := "Failed to delete position"

		if err == usecase.ErrPositionNotFound {
			statusCode = http.StatusNotFound
			errCode = response.ErrorCodeNotFound
			errMsg = "Position not found"
		}

		h.logger.Error().Err(err).Str("id", id).Msg(errMsg)
		c.JSON(statusCode, response.Error(errCode, errMsg))
		return
	}

	c.JSON(http.StatusNoContent, response.Success(nil))
}
