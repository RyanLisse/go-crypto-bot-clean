package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go-crypto-bot-clean/backend/internal/api/dto/request"
	responseDto "go-crypto-bot-clean/backend/internal/api/dto/response"
	"go-crypto-bot-clean/backend/internal/core/trade"
	"go-crypto-bot-clean/backend/internal/domain/repositories"
)

// TradeHandler handles trade-related endpoints
type TradeHandler struct {
	TradeService   trade.TradeService
	BoughtCoinRepo repositories.BoughtCoinRepository
}

// NewTradeHandler creates a new TradeHandler
func NewTradeHandler(tradeService trade.TradeService, boughtCoinRepo repositories.BoughtCoinRepository) *TradeHandler {
	return &TradeHandler{
		TradeService:   tradeService,
		BoughtCoinRepo: boughtCoinRepo,
	}
}

// GetTradeHistory godoc
// @Summary Get trade history
// @Description Returns a list of completed trades
// @Tags trade
// @Accept json
// @Produce json
// @Success 200 {object} responseDto.TradeHistoryResponse
// @Failure 500 {object} responseDto.ErrorResponse
// @Router /api/v1/trade/history [get]
func (h *TradeHandler) GetTradeHistory(c *gin.Context) {
	ctx := c.Request.Context()

	// Get all trades
	allTrades, err := h.BoughtCoinRepo.FindAll(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responseDto.ErrorResponse{
			Code:    "internal_error",
			Message: "Failed to get trade history",
			Details: err.Error(),
		})
		return
	}

	// Filter completed trades (for now, we'll consider all trades as completed for demo purposes)
	completedTrades := allTrades

	// Map to response DTOs
	tradeHistory := make([]responseDto.TradeHistoryItem, len(completedTrades))
	for i, trade := range completedTrades {
		// For demo purposes, we'll use current price as sell price
		sellPrice := trade.CurrentPrice
		if sellPrice == 0 {
			sellPrice = trade.PurchasePrice * 1.05 // Assume 5% profit for demo
		}

		tradeHistory[i] = responseDto.TradeHistoryItem{
			ID:            uint(trade.ID),
			Symbol:        trade.Symbol,
			BuyPrice:      trade.PurchasePrice,
			SellPrice:     sellPrice,
			Quantity:      trade.Quantity,
			BuyTime:       trade.BoughtAt,
			SellTime:      time.Now(), // For demo purposes
			ProfitLoss:    (sellPrice - trade.PurchasePrice) * trade.Quantity,
			ProfitPercent: ((sellPrice - trade.PurchasePrice) / trade.PurchasePrice) * 100,
		}
	}

	// Build response
	resp := responseDto.TradeHistoryResponse{
		Trades:    tradeHistory,
		Count:     len(tradeHistory),
		Timestamp: time.Now(),
	}

	c.JSON(http.StatusOK, resp)
}

// ExecuteTrade godoc
// @Summary Execute a trade
// @Description Executes a buy trade for a specified symbol
// @Tags trade
// @Accept json
// @Produce json
// @Param trade body request.TradeRequest true "Trade details"
// @Success 200 {object} responseDto.TradeExecutionResponse
// @Failure 400 {object} responseDto.ErrorResponse
// @Failure 500 {object} responseDto.ErrorResponse
// @Router /api/v1/trade/buy [post]
func (h *TradeHandler) ExecuteTrade(c *gin.Context) {
	ctx := c.Request.Context()

	// Parse request
	var req request.TradeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, responseDto.ErrorResponse{
			Code:    "invalid_request",
			Message: "Invalid request format",
			Details: err.Error(),
		})
		return
	}

	// Execute purchase
	coin, err := h.TradeService.ExecutePurchase(ctx, req.Symbol, req.Amount, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responseDto.ErrorResponse{
			Code:    "execution_failed",
			Message: "Failed to execute trade",
			Details: err.Error(),
		})
		return
	}

	// Build response
	resp := responseDto.TradeExecutionResponse{
		ID:            uint(coin.ID),
		Symbol:        coin.Symbol,
		Price:         coin.BuyPrice,
		Quantity:      coin.Quantity,
		Total:         coin.BuyPrice * coin.Quantity,
		ExecutionTime: coin.BoughtAt,
		Status:        "completed",
	}

	c.JSON(http.StatusOK, resp)
}

// SellCoin godoc
// @Summary Sell a coin
// @Description Sells a previously bought coin
// @Tags trade
// @Accept json
// @Produce json
// @Param sell body request.SellRequest true "Sell details"
// @Success 200 {object} responseDto.TradeExecutionResponse
// @Failure 400 {object} responseDto.ErrorResponse
// @Failure 404 {object} responseDto.ErrorResponse
// @Failure 500 {object} responseDto.ErrorResponse
// @Router /api/v1/trade/sell [post]
func (h *TradeHandler) SellCoin(c *gin.Context) {
	ctx := c.Request.Context()

	// Parse request
	var req request.SellRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, responseDto.ErrorResponse{
			Code:    "invalid_request",
			Message: "Invalid request format",
			Details: err.Error(),
		})
		return
	}

	// Get the coin
	coin, err := h.BoughtCoinRepo.FindByID(ctx, int64(req.CoinID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, responseDto.ErrorResponse{
			Code:    "internal_error",
			Message: "Failed to find coin",
			Details: err.Error(),
		})
		return
	}

	if coin == nil {
		c.JSON(http.StatusNotFound, responseDto.ErrorResponse{
			Code:    "not_found",
			Message: "Coin not found",
		})
		return
	}

	// Determine amount to sell
	amount := req.Amount
	if req.All || amount <= 0 {
		amount = coin.Quantity
	}

	// Execute sell
	order, err := h.TradeService.SellCoin(ctx, coin, amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responseDto.ErrorResponse{
			Code:    "execution_failed",
			Message: "Failed to sell coin",
			Details: err.Error(),
		})
		return
	}

	// Build response
	resp := responseDto.TradeExecutionResponse{
		ID:            uint(coin.ID),
		Symbol:        order.Symbol,
		Price:         order.Price,
		Quantity:      order.Quantity,
		Total:         order.Price * order.Quantity,
		ExecutionTime: time.Now(),
		Status:        "completed",
	}

	c.JSON(http.StatusOK, resp)
}

// GetTradeStatus godoc
// @Summary Get trade status
// @Description Returns the status of a specific trade
// @Tags trade
// @Accept json
// @Produce json
// @Param id path string true "Trade ID"
// @Success 200 {object} responseDto.TradeStatusResponse
// @Failure 400 {object} responseDto.ErrorResponse
// @Failure 404 {object} responseDto.ErrorResponse
// @Failure 500 {object} responseDto.ErrorResponse
// @Router /api/v1/trade/status/{id} [get]
func (h *TradeHandler) GetTradeStatus(c *gin.Context) {
	ctx := c.Request.Context()

	// Parse trade ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, responseDto.ErrorResponse{
			Code:    "invalid_id",
			Message: "Invalid trade ID format",
			Details: err.Error(),
		})
		return
	}

	// Get the coin
	coin, err := h.BoughtCoinRepo.FindByID(ctx, int64(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, responseDto.ErrorResponse{
			Code:    "internal_error",
			Message: "Failed to find trade",
			Details: err.Error(),
		})
		return
	}

	if coin == nil {
		c.JSON(http.StatusNotFound, responseDto.ErrorResponse{
			Code:    "not_found",
			Message: "Trade not found",
		})
		return
	}

	// Determine status
	status := "active"
	if coin.IsDeleted {
		status = "completed"
	}

	// Build response
	resp := responseDto.TradeStatusResponse{
		ID:            uint(coin.ID),
		Symbol:        coin.Symbol,
		Status:        status,
		PurchasePrice: coin.BuyPrice,
		CurrentPrice:  coin.CurrentPrice,
		Quantity:      coin.Quantity,
		PurchaseTime:  coin.BoughtAt,
		ProfitPercent: ((coin.CurrentPrice - coin.BuyPrice) / coin.BuyPrice) * 100,
		CurrentValue:  coin.CurrentPrice * coin.Quantity,
		StopLossPrice: coin.StopLoss,
		TakeProfitLevels: []responseDto.TakeProfitLevelResponse{
			{
				Price:    coin.TakeProfit,
				Percent:  ((coin.TakeProfit - coin.BuyPrice) / coin.BuyPrice) * 100,
				Executed: false,
			},
		},
	}

	c.JSON(http.StatusOK, resp)
}

// ExecutePurchaseHandler is the legacy handler for executing purchases
// This is kept for backward compatibility
func (h *TradeHandler) ExecutePurchaseHandler(c *gin.Context) {
	var req CreateTradeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, responseDto.ErrorResponse{
			Code:    "INVALID_REQUEST",
			Message: "Invalid request payload",
			Details: err.Error(),
		})
		return
	}

	boughtCoin, err := h.TradeService.ExecutePurchase(c.Request.Context(), req.Symbol, req.Quantity, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responseDto.ErrorResponse{
			Code:    "TRADE_EXECUTION_FAILED",
			Message: "Failed to execute purchase",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, TradeResponse{
		Symbol:   boughtCoin.Symbol,
		Price:    boughtCoin.BuyPrice,
		Quantity: boughtCoin.Quantity,
	})
}

// CreateTradeRequest represents the request body for creating a trade
type CreateTradeRequest struct {
	Symbol   string  `json:"symbol" binding:"required"`
	Price    float64 `json:"price" binding:"required,gt=0"`
	Quantity float64 `json:"quantity" binding:"required,gt=0"`
}

// TradeResponse represents a trade response
type TradeResponse struct {
	ID       string  `json:"id"`
	Symbol   string  `json:"symbol"`
	Price    float64 `json:"price"`
	Quantity float64 `json:"quantity"`
}
