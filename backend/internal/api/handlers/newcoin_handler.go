package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go-crypto-bot-clean/backend/internal/api/dto/request"
	"go-crypto-bot-clean/backend/internal/api/dto/response"
	"go-crypto-bot-clean/backend/internal/core/newcoin"
)

// NewCoinsHandler handles new coin detection endpoints
type NewCoinsHandler struct {
	NewCoinService newcoin.NewCoinService
}

// NewNewCoinsHandler creates a new NewCoinsHandler
func NewNewCoinsHandler(newCoinService newcoin.NewCoinService) *NewCoinsHandler {
	return &NewCoinsHandler{
		NewCoinService: newCoinService,
	}
}

// GetDetectedCoins godoc
// @Summary Get detected coins
// @Description Returns a list of newly detected coins
// @Tags newcoins
// @Accept json
// @Produce json
// @Success 200 {object} response.NewCoinsListResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/newcoins [get]
func (h *NewCoinsHandler) GetDetectedCoins(c *gin.Context) {
	ctx := c.Request.Context()

	// Get all new coins
	coins, err := h.NewCoinService.GetAllNewCoins(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Code:    "internal_error",
			Message: "Failed to get new coins",
			Details: err.Error(),
		})
		return
	}

	// Map to response DTOs
	coinResponses := make([]response.NewCoinResponse, len(coins))
	for i, coin := range coins {
		coinResponses[i] = response.NewCoinResponse{
			ID:            uint(coin.ID),
			Symbol:        coin.Symbol,
			FoundAt:       coin.FoundAt,
			FirstOpenTime: coin.FirstOpenTime,
			QuoteVolume:   coin.QuoteVolume,
			IsProcessed:   coin.IsProcessed,
			IsUpcoming:    coin.IsUpcoming,
		}
	}

	// Build response
	resp := response.NewCoinsListResponse{
		Coins:     coinResponses,
		Count:     len(coinResponses),
		Timestamp: time.Now(),
	}

	c.JSON(http.StatusOK, resp)
}

// ProcessNewCoins godoc
// @Summary Process new coins
// @Description Processes newly detected coins for trading
// @Tags newcoins
// @Accept json
// @Produce json
// @Param request body request.ProcessNewCoinsRequest true "Coins to process"
// @Success 200 {object} response.ProcessNewCoinsResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/newcoins/process [post]
func (h *NewCoinsHandler) ProcessNewCoins(c *gin.Context) {
	ctx := c.Request.Context()

	// Parse request
	var req request.ProcessNewCoinsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Code:    "invalid_request",
			Message: "Invalid request format",
			Details: err.Error(),
		})
		return
	}

	// Process each coin
	processedCoins := make([]response.NewCoinResponse, 0, len(req.CoinIDs))
	for _, coinID := range req.CoinIDs {
		// Mark as processed
		err := h.NewCoinService.MarkAsProcessed(ctx, int64(coinID))
		if err != nil {
			continue
		}

		// Get the updated coin
		coin, err := h.NewCoinService.GetCoinByID(ctx, int64(coinID))
		if err != nil {
			continue
		}

		// Add to processed coins
		processedCoins = append(processedCoins, response.NewCoinResponse{
			ID:          uint(coin.ID),
			Symbol:      coin.Symbol,
			FoundAt:     coin.FoundAt,
			QuoteVolume: coin.QuoteVolume,
			IsProcessed: true,
		})
	}

	// Build response
	resp := response.ProcessNewCoinsResponse{
		ProcessedCoins: processedCoins,
		Count:          len(processedCoins),
		Timestamp:      time.Now(),
	}

	c.JSON(http.StatusOK, resp)
}

// DetectNewCoins godoc
// @Summary Detect new coins
// @Description Triggers detection of new coins
// @Tags newcoins
// @Accept json
// @Produce json
// @Success 200 {object} response.NewCoinsListResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/newcoins/detect [post]
func (h *NewCoinsHandler) DetectNewCoins(c *gin.Context) {
	ctx := c.Request.Context()

	// Detect new coins
	newCoins, err := h.NewCoinService.DetectNewCoins(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Code:    "internal_error",
			Message: "Failed to detect new coins",
			Details: err.Error(),
		})
		return
	}

	// Save new coins
	if len(newCoins) > 0 {
		err = h.NewCoinService.SaveNewCoins(ctx, newCoins)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.ErrorResponse{
				Code:    "internal_error",
				Message: "Failed to save new coins",
				Details: err.Error(),
			})
			return
		}
	}

	// Map to response DTOs
	coinResponses := make([]response.NewCoinResponse, len(newCoins))
	for i, coin := range newCoins {
		coinResponses[i] = response.NewCoinResponse{
			ID:          uint(coin.ID),
			Symbol:      coin.Symbol,
			FoundAt:     coin.FoundAt,
			QuoteVolume: coin.QuoteVolume,
			IsProcessed: coin.IsProcessed,
		}
	}

	// Build response
	resp := response.NewCoinsListResponse{
		Coins:     coinResponses,
		Count:     len(coinResponses),
		Timestamp: time.Now(),
	}

	c.JSON(http.StatusOK, resp)
}

// GetCoinsByDate godoc
// @Summary Get coins by date
// @Description Returns a list of coins found on a specific date
// @Tags newcoins
// @Accept json
// @Produce json
// @Param request body request.DateFilterRequest true "Date to filter by"
// @Success 200 {object} response.NewCoinsListResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/newcoins/by-date [post]
func (h *NewCoinsHandler) GetCoinsByDate(c *gin.Context) {
	ctx := c.Request.Context()

	// Parse request
	var req request.DateFilterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Code:    "invalid_request",
			Message: "Invalid request format",
			Details: err.Error(),
		})
		return
	}

	// Parse date
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Code:    "invalid_date",
			Message: "Invalid date format. Use YYYY-MM-DD",
			Details: err.Error(),
		})
		return
	}

	// Get coins by date
	coins, err := h.NewCoinService.GetCoinsByDate(ctx, date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Code:    "internal_error",
			Message: "Failed to get coins by date",
			Details: err.Error(),
		})
		return
	}

	// Map to response DTOs
	coinResponses := make([]response.NewCoinResponse, len(coins))
	for i, coin := range coins {
		coinResponses[i] = response.NewCoinResponse{
			ID:          uint(coin.ID),
			Symbol:      coin.Symbol,
			FoundAt:     coin.FoundAt,
			QuoteVolume: coin.QuoteVolume,
			IsProcessed: coin.IsProcessed,
		}
	}

	// Build response
	resp := response.NewCoinsListResponse{
		Coins:     coinResponses,
		Count:     len(coinResponses),
		Timestamp: time.Now(),
	}

	c.JSON(http.StatusOK, resp)
}

// GetCoinsByDateRange godoc
// @Summary Get coins by date range
// @Description Returns a list of coins found within a date range
// @Tags newcoins
// @Accept json
// @Produce json
// @Param request body request.DateRangeFilterRequest true "Date range to filter by"
// @Success 200 {object} response.NewCoinsListResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/newcoins/by-date-range [post]
func (h *NewCoinsHandler) GetCoinsByDateRange(c *gin.Context) {
	ctx := c.Request.Context()

	// Parse request
	var req request.DateRangeFilterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Code:    "invalid_request",
			Message: "Invalid request format",
			Details: err.Error(),
		})
		return
	}

	// Parse dates
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Code:    "invalid_date",
			Message: "Invalid start date format. Use YYYY-MM-DD",
			Details: err.Error(),
		})
		return
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Code:    "invalid_date",
			Message: "Invalid end date format. Use YYYY-MM-DD",
			Details: err.Error(),
		})
		return
	}

	// Validate date range
	if endDate.Before(startDate) {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Code:    "invalid_date_range",
			Message: "End date must be after start date",
		})
		return
	}

	// Get coins by date range
	coins, err := h.NewCoinService.GetCoinsByDateRange(ctx, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Code:    "internal_error",
			Message: "Failed to get coins by date range",
			Details: err.Error(),
		})
		return
	}

	// Map to response DTOs
	coinResponses := make([]response.NewCoinResponse, len(coins))
	for i, coin := range coins {
		coinResponses[i] = response.NewCoinResponse{
			ID:            uint(coin.ID),
			Symbol:        coin.Symbol,
			FoundAt:       coin.FoundAt,
			FirstOpenTime: coin.FirstOpenTime,
			QuoteVolume:   coin.QuoteVolume,
			IsProcessed:   coin.IsProcessed,
			IsUpcoming:    coin.IsUpcoming,
		}
	}

	// Build response
	resp := response.NewCoinsListResponse{
		Coins:     coinResponses,
		Count:     len(coinResponses),
		Timestamp: time.Now(),
	}

	c.JSON(http.StatusOK, resp)
}

// GetUpcomingCoins godoc
// @Summary Get upcoming coins
// @Description Returns a list of coins scheduled to be listed in the future
// @Tags newcoins
// @Accept json
// @Produce json
// @Success 200 {object} response.NewCoinsListResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/newcoins/upcoming [get]
func (h *NewCoinsHandler) GetUpcomingCoins(c *gin.Context) {
	ctx := c.Request.Context()

	// Get upcoming coins
	coins, err := h.NewCoinService.GetUpcomingCoins(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Code:    "internal_error",
			Message: "Failed to get upcoming coins",
			Details: err.Error(),
		})
		return
	}

	// Map to response DTOs
	coinResponses := make([]response.NewCoinResponse, len(coins))
	for i, coin := range coins {
		coinResponses[i] = response.NewCoinResponse{
			ID:            uint(coin.ID),
			Symbol:        coin.Symbol,
			FoundAt:       coin.FoundAt,
			FirstOpenTime: coin.FirstOpenTime,
			QuoteVolume:   coin.QuoteVolume,
			IsProcessed:   coin.IsProcessed,
			IsUpcoming:    coin.IsUpcoming,
		}
	}

	// Build response
	resp := response.NewCoinsListResponse{
		Coins:     coinResponses,
		Count:     len(coinResponses),
		Timestamp: time.Now(),
	}

	c.JSON(http.StatusOK, resp)
}

// GetUpcomingCoinsForTodayAndTomorrow godoc
// @Summary Get upcoming coins for today and tomorrow
// @Description Returns a list of coins scheduled to be listed today or tomorrow
// @Tags newcoins
// @Accept json
// @Produce json
// @Success 200 {object} response.NewCoinsListResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/newcoins/upcoming/today-and-tomorrow [get]
func (h *NewCoinsHandler) GetUpcomingCoinsForTodayAndTomorrow(c *gin.Context) {
	ctx := c.Request.Context()

	// Get upcoming coins for today and tomorrow
	coins, err := h.NewCoinService.GetUpcomingCoinsForTodayAndTomorrow(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Code:    "internal_error",
			Message: "Failed to get upcoming coins for today and tomorrow",
			Details: err.Error(),
		})
		return
	}

	// Map to response DTOs
	coinResponses := make([]response.NewCoinResponse, len(coins))
	for i, coin := range coins {
		coinResponses[i] = response.NewCoinResponse{
			ID:            uint(coin.ID),
			Symbol:        coin.Symbol,
			FoundAt:       coin.FoundAt,
			FirstOpenTime: coin.FirstOpenTime,
			QuoteVolume:   coin.QuoteVolume,
			IsProcessed:   coin.IsProcessed,
			IsUpcoming:    coin.IsUpcoming,
		}
	}

	// Build response
	resp := response.NewCoinsListResponse{
		Coins:     coinResponses,
		Count:     len(coinResponses),
		Timestamp: time.Now(),
	}

	c.JSON(http.StatusOK, resp)
}

// GetUpcomingCoinsByDate godoc
// @Summary Get upcoming coins by date
// @Description Returns a list of coins scheduled to be listed on a specific date
// @Tags newcoins
// @Accept json
// @Produce json
// @Param request body request.DateRequest true "Date request"
// @Success 200 {object} response.NewCoinsListResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/newcoins/upcoming/by-date [post]
func (h *NewCoinsHandler) GetUpcomingCoinsByDate(c *gin.Context) {
	ctx := c.Request.Context()

	// Parse request
	var req request.DateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Code:    "invalid_request",
			Message: "Invalid request format",
			Details: err.Error(),
		})
		return
	}

	// Parse date
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Code:    "invalid_date",
			Message: "Invalid date format. Use YYYY-MM-DD",
			Details: err.Error(),
		})
		return
	}

	// Get upcoming coins by date
	coins, err := h.NewCoinService.GetUpcomingCoinsByDate(ctx, date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Code:    "internal_error",
			Message: "Failed to get upcoming coins by date",
			Details: err.Error(),
		})
		return
	}

	// Map to response DTOs
	coinResponses := make([]response.NewCoinResponse, len(coins))
	for i, coin := range coins {
		coinResponses[i] = response.NewCoinResponse{
			ID:            uint(coin.ID),
			Symbol:        coin.Symbol,
			FoundAt:       coin.FoundAt,
			FirstOpenTime: coin.FirstOpenTime,
			QuoteVolume:   coin.QuoteVolume,
			IsProcessed:   coin.IsProcessed,
			IsUpcoming:    coin.IsUpcoming,
		}
	}

	// Build response
	resp := response.NewCoinsListResponse{
		Coins:     coinResponses,
		Count:     len(coinResponses),
		Timestamp: time.Now(),
	}

	c.JSON(http.StatusOK, resp)
}
