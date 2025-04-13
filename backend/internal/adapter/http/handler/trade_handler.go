package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/neo/crypto-bot/internal/adapter/http/response"
	"github.com/neo/crypto-bot/internal/domain/model"
	"github.com/neo/crypto-bot/internal/usecase"
	"github.com/rs/zerolog"
)

// TradeHandler handles HTTP requests for trade execution
type TradeHandler struct {
	useCase usecase.TradeUseCase
	logger  *zerolog.Logger
}

// NewTradeHandler creates a new TradeHandler
func NewTradeHandler(uc usecase.TradeUseCase, logger *zerolog.Logger) *TradeHandler {
	return &TradeHandler{
		useCase: uc,
		logger:  logger,
	}
}

// RegisterRoutes registers trade-related routes with the Gin engine
func (h *TradeHandler) RegisterRoutes(router *gin.RouterGroup) {
	tradeGroup := router.Group("/trades")
	{
		tradeGroup.POST("/orders", h.PlaceOrder)
		tradeGroup.DELETE("/orders/:symbol/:orderId", h.CancelOrder)
		tradeGroup.GET("/orders/:symbol/:orderId", h.GetOrderStatus)
		tradeGroup.GET("/orders/open/:symbol", h.GetOpenOrders)
		tradeGroup.GET("/orders/history/:symbol", h.GetOrderHistory)
		tradeGroup.GET("/calculate/:symbol", h.CalculateQuantity)
	}
}

// PlaceOrder godoc
// @Summary Place a new order
// @Description Place a market or limit order
// @Tags trades
// @Accept json
// @Produce json
// @Param order body model.OrderRequest true "Order data"
// @Success 201 {object} response.APIResponse{data=model.OrderResponse}
// @Failure 400 {object} response.APIResponse{error=response.APIError}
// @Failure 500 {object} response.APIResponse{error=response.APIError}
// @Router /trades/orders [post]
func (h *TradeHandler) PlaceOrder(c *gin.Context) {
	var req model.OrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error().Err(err).Msg("Invalid order request")
		c.JSON(http.StatusBadRequest, response.Error(response.ErrorCodeBadRequest, "Invalid order data: "+err.Error()))
		return
	}

	// Validate limit order has price
	if req.Type == model.OrderTypeLimit && req.Price <= 0 {
		h.logger.Error().Msg("Limit order must include a valid price")
		c.JSON(http.StatusBadRequest, response.Error(response.ErrorCodeBadRequest, "Limit order must include a valid price"))
		return
	}

	order, err := h.useCase.PlaceOrder(c.Request.Context(), req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		errCode := response.ErrorCodeInternalError
		errMsg := "Failed to place order"

		if err == usecase.ErrSymbolNotFound {
			statusCode = http.StatusBadRequest
			errCode = response.ErrorCodeBadRequest
			errMsg = "Symbol not found"
		} else if err == usecase.ErrInvalidOrderData {
			statusCode = http.StatusBadRequest
			errCode = response.ErrorCodeBadRequest
			errMsg = "Invalid order data"
		} else if err == usecase.ErrInsufficientBalance {
			statusCode = http.StatusBadRequest
			errCode = response.ErrorCodeBadRequest
			errMsg = "Insufficient balance for order"
		}

		h.logger.Error().Err(err).
			Str("symbol", req.Symbol).
			Str("side", string(req.Side)).
			Str("type", string(req.Type)).
			Float64("quantity", req.Quantity).
			Float64("price", req.Price).
			Msg(errMsg)

		c.JSON(statusCode, response.Error(errCode, errMsg))
		return
	}

	h.logger.Info().
		Str("orderId", order.ID).
		Str("symbol", order.Symbol).
		Str("side", string(order.Side)).
		Str("type", string(order.Type)).
		Float64("quantity", order.Quantity).
		Float64("price", order.Price).
		Msg("Order placed successfully")

	c.JSON(http.StatusCreated, response.Success(order))
}

// CancelOrder godoc
// @Summary Cancel an existing order
// @Description Cancel an open order
// @Tags trades
// @Accept json
// @Produce json
// @Param symbol path string true "Trading symbol"
// @Param orderId path string true "Order ID"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse{error=response.APIError}
// @Failure 404 {object} response.APIResponse{error=response.APIError}
// @Failure 500 {object} response.APIResponse{error=response.APIError}
// @Router /trades/orders/{symbol}/{orderId} [delete]
func (h *TradeHandler) CancelOrder(c *gin.Context) {
	symbol := c.Param("symbol")
	orderID := c.Param("orderId")

	if symbol == "" || orderID == "" {
		c.JSON(http.StatusBadRequest, response.Error(response.ErrorCodeBadRequest, "Symbol and order ID are required"))
		return
	}

	err := h.useCase.CancelOrder(c.Request.Context(), symbol, orderID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		errCode := response.ErrorCodeInternalError
		errMsg := "Failed to cancel order"

		if err == usecase.ErrOrderNotFound {
			statusCode = http.StatusNotFound
			errCode = response.ErrorCodeNotFound
			errMsg = "Order not found"
		}

		h.logger.Error().Err(err).
			Str("symbol", symbol).
			Str("orderId", orderID).
			Msg(errMsg)

		c.JSON(statusCode, response.Error(errCode, errMsg))
		return
	}

	h.logger.Info().
		Str("symbol", symbol).
		Str("orderId", orderID).
		Msg("Order canceled successfully")

	c.JSON(http.StatusOK, response.Success(nil))
}

// GetOrderStatus godoc
// @Summary Get order status
// @Description Get the current status of an order
// @Tags trades
// @Accept json
// @Produce json
// @Param symbol path string true "Trading symbol"
// @Param orderId path string true "Order ID"
// @Success 200 {object} response.APIResponse{data=model.Order}
// @Failure 400 {object} response.APIResponse{error=response.APIError}
// @Failure 404 {object} response.APIResponse{error=response.APIError}
// @Failure 500 {object} response.APIResponse{error=response.APIError}
// @Router /trades/orders/{symbol}/{orderId} [get]
func (h *TradeHandler) GetOrderStatus(c *gin.Context) {
	symbol := c.Param("symbol")
	orderID := c.Param("orderId")

	if symbol == "" || orderID == "" {
		c.JSON(http.StatusBadRequest, response.Error(response.ErrorCodeBadRequest, "Symbol and order ID are required"))
		return
	}

	order, err := h.useCase.GetOrderStatus(c.Request.Context(), symbol, orderID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		errCode := response.ErrorCodeInternalError
		errMsg := "Failed to get order status"

		if err == usecase.ErrOrderNotFound {
			statusCode = http.StatusNotFound
			errCode = response.ErrorCodeNotFound
			errMsg = "Order not found"
		}

		h.logger.Error().Err(err).
			Str("symbol", symbol).
			Str("orderId", orderID).
			Msg(errMsg)

		c.JSON(statusCode, response.Error(errCode, errMsg))
		return
	}

	c.JSON(http.StatusOK, response.Success(order))
}

// GetOpenOrders godoc
// @Summary Get all open orders
// @Description Get all currently open orders for a symbol
// @Tags trades
// @Accept json
// @Produce json
// @Param symbol path string true "Trading symbol"
// @Success 200 {object} response.APIResponse{data=[]model.Order}
// @Failure 400 {object} response.APIResponse{error=response.APIError}
// @Failure 500 {object} response.APIResponse{error=response.APIError}
// @Router /trades/orders/open/{symbol} [get]
func (h *TradeHandler) GetOpenOrders(c *gin.Context) {
	symbol := c.Param("symbol")

	if symbol == "" {
		c.JSON(http.StatusBadRequest, response.Error(response.ErrorCodeBadRequest, "Symbol is required"))
		return
	}

	orders, err := h.useCase.GetOpenOrders(c.Request.Context(), symbol)
	if err != nil {
		h.logger.Error().Err(err).
			Str("symbol", symbol).
			Msg("Failed to get open orders")

		c.JSON(http.StatusInternalServerError, response.Error(response.ErrorCodeInternalError, "Failed to get open orders"))
		return
	}

	c.JSON(http.StatusOK, response.Success(orders))
}

// GetOrderHistory godoc
// @Summary Get order history
// @Description Get historical orders for a symbol with pagination
// @Tags trades
// @Accept json
// @Produce json
// @Param symbol path string true "Trading symbol"
// @Param limit query int false "Limit number of results (default: 50)"
// @Param offset query int false "Offset for pagination (default: 0)"
// @Success 200 {object} response.APIResponse{data=[]model.Order}
// @Failure 400 {object} response.APIResponse{error=response.APIError}
// @Failure 500 {object} response.APIResponse{error=response.APIError}
// @Router /trades/orders/history/{symbol} [get]
func (h *TradeHandler) GetOrderHistory(c *gin.Context) {
	symbol := c.Param("symbol")
	if symbol == "" {
		c.JSON(http.StatusBadRequest, response.Error(response.ErrorCodeBadRequest, "Symbol is required"))
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

	orders, err := h.useCase.GetOrderHistory(c.Request.Context(), symbol, limit, offset)
	if err != nil {
		h.logger.Error().Err(err).
			Str("symbol", symbol).
			Int("limit", limit).
			Int("offset", offset).
			Msg("Failed to get order history")

		c.JSON(http.StatusInternalServerError, response.Error(response.ErrorCodeInternalError, "Failed to get order history"))
		return
	}

	c.JSON(http.StatusOK, response.Success(orders))
}

// CalculateQuantity godoc
// @Summary Calculate required quantity
// @Description Calculate the quantity of an asset needed for a trade based on amount
// @Tags trades
// @Accept json
// @Produce json
// @Param symbol path string true "Trading symbol"
// @Param side query string true "Order side (BUY or SELL)"
// @Param amount query number true "Amount in quote currency to spend/receive"
// @Success 200 {object} response.APIResponse{data=map[string]float64}
// @Failure 400 {object} response.APIResponse{error=response.APIError}
// @Failure 500 {object} response.APIResponse{error=response.APIError}
// @Router /trades/calculate/{symbol} [get]
func (h *TradeHandler) CalculateQuantity(c *gin.Context) {
	symbol := c.Param("symbol")
	if symbol == "" {
		c.JSON(http.StatusBadRequest, response.Error(response.ErrorCodeBadRequest, "Symbol is required"))
		return
	}

	sideStr := c.Query("side")
	if sideStr == "" {
		c.JSON(http.StatusBadRequest, response.Error(response.ErrorCodeBadRequest, "Order side is required"))
		return
	}
	side := model.OrderSide(sideStr)
	if side != model.OrderSideBuy && side != model.OrderSideSell {
		c.JSON(http.StatusBadRequest, response.Error(response.ErrorCodeBadRequest, "Invalid order side. Must be BUY or SELL"))
		return
	}

	amountStr := c.Query("amount")
	if amountStr == "" {
		c.JSON(http.StatusBadRequest, response.Error(response.ErrorCodeBadRequest, "Amount is required"))
		return
	}
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil || amount <= 0 {
		c.JSON(http.StatusBadRequest, response.Error(response.ErrorCodeBadRequest, "Invalid amount. Must be a positive number"))
		return
	}

	quantity, err := h.useCase.CalculateRequiredQuantity(c.Request.Context(), symbol, side, amount)
	if err != nil {
		statusCode := http.StatusInternalServerError
		errCode := response.ErrorCodeInternalError
		errMsg := "Failed to calculate quantity"

		if err == usecase.ErrSymbolNotFound {
			statusCode = http.StatusBadRequest
			errCode = response.ErrorCodeBadRequest
			errMsg = "Symbol not found"
		}

		h.logger.Error().Err(err).
			Str("symbol", symbol).
			Str("side", string(side)).
			Float64("amount", amount).
			Msg(errMsg)

		c.JSON(statusCode, response.Error(errCode, errMsg))
		return
	}

	result := map[string]float64{
		"quantity": quantity,
		"amount":   amount,
	}

	c.JSON(http.StatusOK, response.Success(result))
}
