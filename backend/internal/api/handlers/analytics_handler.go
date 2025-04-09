package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"go-crypto-bot-clean/backend/internal/api/dto/request"
	"go-crypto-bot-clean/backend/internal/api/dto/response"
	"go-crypto-bot-clean/backend/internal/core/analytics"
	"go-crypto-bot-clean/backend/internal/domain/models"
)

// AnalyticsHandler handles API requests for trade analytics
type AnalyticsHandler struct {
	analyticsService analytics.TradeAnalyticsService
	logger           *zap.Logger
}

// NewAnalyticsHandler creates a new analytics handler
func NewAnalyticsHandler(analyticsService analytics.TradeAnalyticsService, logger *zap.Logger) *AnalyticsHandler {
	return &AnalyticsHandler{
		analyticsService: analyticsService,
		logger:           logger,
	}
}

// GetTradeAnalytics handles requests for trade analytics
func (h *AnalyticsHandler) GetTradeAnalytics(c *gin.Context) {
	// Parse request parameters
	var req request.TradeAnalyticsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set default time range if not provided
	if req.StartTime.IsZero() {
		// Default to 30 days ago
		req.StartTime = time.Now().AddDate(0, 0, -30)
	}
	if req.EndTime.IsZero() {
		req.EndTime = time.Now()
	}

	// Convert timeframe string to enum
	timeFrame := models.TimeFrameAll
	switch req.TimeFrame {
	case "day":
		timeFrame = models.TimeFrameDay
	case "week":
		timeFrame = models.TimeFrameWeek
	case "month":
		timeFrame = models.TimeFrameMonth
	case "quarter":
		timeFrame = models.TimeFrameQuarter
	case "year":
		timeFrame = models.TimeFrameYear
	}

	// Get analytics from service
	analytics, err := h.analyticsService.GetTradeAnalytics(c.Request.Context(), timeFrame, req.StartTime, req.EndTime)
	if err != nil {
		h.logger.Error("Failed to get trade analytics", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get trade analytics"})
		return
	}

	// Convert to response DTO
	resp := response.TradeAnalyticsFromModel(analytics)
	c.JSON(http.StatusOK, resp)
}

// GetTradePerformance handles requests for individual trade performance
func (h *AnalyticsHandler) GetTradePerformance(c *gin.Context) {
	// Get trade ID from path
	tradeID := c.Param("id")
	if tradeID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Trade ID is required"})
		return
	}

	// Get trade performance from service
	performance, err := h.analyticsService.GetTradePerformance(c.Request.Context(), tradeID)
	if err != nil {
		h.logger.Error("Failed to get trade performance", zap.Error(err), zap.String("tradeID", tradeID))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get trade performance"})
		return
	}

	// Convert to response DTO
	resp := response.TradePerformanceFromModel(performance)
	c.JSON(http.StatusOK, resp)
}

// GetAllTradePerformance handles requests for all trade performances
func (h *AnalyticsHandler) GetAllTradePerformance(c *gin.Context) {
	// Parse request parameters
	startTimeStr := c.Query("start_time")
	endTimeStr := c.Query("end_time")
	limitStr := c.Query("limit")

	var startTime, endTime time.Time
	var err error

	// Parse start time
	if startTimeStr != "" {
		startTime, err = time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_time format"})
			return
		}
	} else {
		// Default to 30 days ago
		startTime = time.Now().AddDate(0, 0, -30)
	}

	// Parse end time
	if endTimeStr != "" {
		endTime, err = time.Parse(time.RFC3339, endTimeStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_time format"})
			return
		}
	} else {
		endTime = time.Now()
	}

	// Get trade performances from service
	performances, err := h.analyticsService.GetAllTradePerformance(c.Request.Context(), startTime, endTime)
	if err != nil {
		h.logger.Error("Failed to get trade performances", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get trade performances"})
		return
	}

	// Apply limit if provided
	if limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit"})
			return
		}

		if limit < len(performances) {
			performances = performances[:limit]
		}
	}

	// Convert to response DTOs
	resp := make([]response.TradePerformanceResponse, len(performances))
	for i, perf := range performances {
		resp[i] = response.TradePerformanceFromModel(perf)
	}

	c.JSON(http.StatusOK, resp)
}

// GetWinRate handles requests for win rate
func (h *AnalyticsHandler) GetWinRate(c *gin.Context) {
	// Parse request parameters
	startTimeStr := c.Query("start_time")
	endTimeStr := c.Query("end_time")

	var startTime, endTime time.Time
	var err error

	// Parse start time
	if startTimeStr != "" {
		startTime, err = time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_time format"})
			return
		}
	} else {
		// Default to 30 days ago
		startTime = time.Now().AddDate(0, 0, -30)
	}

	// Parse end time
	if endTimeStr != "" {
		endTime, err = time.Parse(time.RFC3339, endTimeStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_time format"})
			return
		}
	} else {
		endTime = time.Now()
	}

	// Get win rate from service
	winRate, err := h.analyticsService.GetWinRate(c.Request.Context(), startTime, endTime)
	if err != nil {
		h.logger.Error("Failed to get win rate", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get win rate"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"win_rate": winRate})
}

// GetBalanceHistory handles requests for balance history
func (h *AnalyticsHandler) GetBalanceHistory(c *gin.Context) {
	// Parse request parameters
	startTimeStr := c.Query("start_time")
	endTimeStr := c.Query("end_time")
	intervalStr := c.Query("interval")

	var startTime, endTime time.Time
	var interval time.Duration
	var err error

	// Parse start time
	if startTimeStr != "" {
		startTime, err = time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_time format"})
			return
		}
	} else {
		// Default to 30 days ago
		startTime = time.Now().AddDate(0, 0, -30)
	}

	// Parse end time
	if endTimeStr != "" {
		endTime, err = time.Parse(time.RFC3339, endTimeStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_time format"})
			return
		}
	} else {
		endTime = time.Now()
	}

	// Parse interval
	if intervalStr != "" {
		interval, err = time.ParseDuration(intervalStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid interval format"})
			return
		}
	} else {
		// Default to 1 day
		interval = 24 * time.Hour
	}

	// Get balance history from service
	balanceHistory, err := h.analyticsService.GetBalanceHistory(c.Request.Context(), startTime, endTime, interval)
	if err != nil {
		h.logger.Error("Failed to get balance history", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get balance history"})
		return
	}

	// Convert to response DTOs
	resp := make([]response.BalancePointResponse, len(balanceHistory))
	for i, point := range balanceHistory {
		resp[i] = response.BalancePointResponse{
			Timestamp: point.Timestamp,
			Balance:   point.Balance,
		}
	}

	c.JSON(http.StatusOK, resp)
}

// GetPerformanceBySymbol handles requests for performance by symbol
func (h *AnalyticsHandler) GetPerformanceBySymbol(c *gin.Context) {
	// Parse request parameters
	startTimeStr := c.Query("start_time")
	endTimeStr := c.Query("end_time")

	var startTime, endTime time.Time
	var err error

	// Parse start time
	if startTimeStr != "" {
		startTime, err = time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_time format"})
			return
		}
	} else {
		// Default to 30 days ago
		startTime = time.Now().AddDate(0, 0, -30)
	}

	// Parse end time
	if endTimeStr != "" {
		endTime, err = time.Parse(time.RFC3339, endTimeStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_time format"})
			return
		}
	} else {
		endTime = time.Now()
	}

	// Get performance by symbol from service
	performance, err := h.analyticsService.GetPerformanceBySymbol(c.Request.Context(), startTime, endTime)
	if err != nil {
		h.logger.Error("Failed to get performance by symbol", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get performance by symbol"})
		return
	}

	// Convert to response DTOs
	resp := make(map[string]response.SymbolPerformanceResponse)
	for symbol, perf := range performance {
		resp[symbol] = response.SymbolPerformanceFromModel(perf)
	}

	c.JSON(http.StatusOK, resp)
}

// GetPerformanceByReason handles requests for performance by reason
func (h *AnalyticsHandler) GetPerformanceByReason(c *gin.Context) {
	// Parse request parameters
	startTimeStr := c.Query("start_time")
	endTimeStr := c.Query("end_time")

	var startTime, endTime time.Time
	var err error

	// Parse start time
	if startTimeStr != "" {
		startTime, err = time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_time format"})
			return
		}
	} else {
		// Default to 30 days ago
		startTime = time.Now().AddDate(0, 0, -30)
	}

	// Parse end time
	if endTimeStr != "" {
		endTime, err = time.Parse(time.RFC3339, endTimeStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_time format"})
			return
		}
	} else {
		endTime = time.Now()
	}

	// Get performance by reason from service
	performance, err := h.analyticsService.GetPerformanceByReason(c.Request.Context(), startTime, endTime)
	if err != nil {
		h.logger.Error("Failed to get performance by reason", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get performance by reason"})
		return
	}

	// Convert to response DTOs
	resp := make(map[string]response.ReasonPerformanceResponse)
	for reason, perf := range performance {
		resp[reason] = response.ReasonPerformanceFromModel(perf)
	}

	c.JSON(http.StatusOK, resp)
}

// GetPerformanceByStrategy handles requests for performance by strategy
func (h *AnalyticsHandler) GetPerformanceByStrategy(c *gin.Context) {
	// Parse request parameters
	startTimeStr := c.Query("start_time")
	endTimeStr := c.Query("end_time")

	var startTime, endTime time.Time
	var err error

	// Parse start time
	if startTimeStr != "" {
		startTime, err = time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_time format"})
			return
		}
	} else {
		// Default to 30 days ago
		startTime = time.Now().AddDate(0, 0, -30)
	}

	// Parse end time
	if endTimeStr != "" {
		endTime, err = time.Parse(time.RFC3339, endTimeStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_time format"})
			return
		}
	} else {
		endTime = time.Now()
	}

	// Get performance by strategy from service
	performance, err := h.analyticsService.GetPerformanceByStrategy(c.Request.Context(), startTime, endTime)
	if err != nil {
		h.logger.Error("Failed to get performance by strategy", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get performance by strategy"})
		return
	}

	// Convert to response DTOs
	resp := make(map[string]response.StrategyPerformanceResponse)
	for strategy, perf := range performance {
		resp[strategy] = response.StrategyPerformanceFromModel(perf)
	}

	c.JSON(http.StatusOK, resp)
}
