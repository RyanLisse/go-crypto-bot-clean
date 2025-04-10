package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

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
func (h *AnalyticsHandler) GetTradeAnalytics(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	// Parse request parameters
	var req request.TradeAnalyticsRequest
	req.TimeFrame = query.Get("timeFrame")

	startTimeStr := query.Get("startTime")
	endTimeStr := query.Get("endTime")

	var err error
	if startTimeStr != "" {
		req.StartTime, err = time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			http.Error(w, "Invalid start time", http.StatusBadRequest)
			return
		}
	}

	if endTimeStr != "" {
		req.EndTime, err = time.Parse(time.RFC3339, endTimeStr)
		if err != nil {
			http.Error(w, "Invalid end time", http.StatusBadRequest)
			return
		}
	}

	// Set default time range if not provided
	if req.StartTime.IsZero() {
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
	analytics, err := h.analyticsService.GetTradeAnalytics(r.Context(), timeFrame, req.StartTime, req.EndTime)
	if err != nil {
		h.logger.Error("Failed to get trade analytics", zap.Error(err))
		http.Error(w, "Failed to get trade analytics", http.StatusInternalServerError)
		return
	}

	// Convert to response DTO
	resp := response.TradeAnalyticsFromModel(analytics)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// GetTradePerformance handles requests for individual trade performance
func (h *AnalyticsHandler) GetTradePerformance(w http.ResponseWriter, r *http.Request) {
	// Get trade ID from path
	tradeID := r.URL.Query().Get("id")
	if tradeID == "" {
		http.Error(w, "Trade ID is required", http.StatusBadRequest)
		return
	}

	// Get trade performance from service
	performance, err := h.analyticsService.GetTradePerformance(r.Context(), tradeID)
	if err != nil {
		h.logger.Error("Failed to get trade performance", zap.Error(err), zap.String("tradeID", tradeID))
		http.Error(w, "Failed to get trade performance", http.StatusInternalServerError)
		return
	}

	// Convert to response DTO
	resp := response.TradePerformanceFromModel(performance)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// GetAllTradePerformance handles requests for all trade performances
func (h *AnalyticsHandler) GetAllTradePerformance(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	// Parse request parameters
	startTimeStr := query.Get("start_time")
	endTimeStr := query.Get("end_time")
	limitStr := query.Get("limit")

	var startTime, endTime time.Time
	var err error

	// Parse start time
	if startTimeStr != "" {
		startTime, err = time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			http.Error(w, "Invalid start_time format", http.StatusBadRequest)
			return
		}
	} else {
		startTime = time.Now().AddDate(0, 0, -30)
	}

	// Parse end time
	if endTimeStr != "" {
		endTime, err = time.Parse(time.RFC3339, endTimeStr)
		if err != nil {
			http.Error(w, "Invalid end_time format", http.StatusBadRequest)
			return
		}
	} else {
		endTime = time.Now()
	}

	// Get trade performances from service
	performances, err := h.analyticsService.GetAllTradePerformance(r.Context(), startTime, endTime)
	if err != nil {
		h.logger.Error("Failed to get trade performances", zap.Error(err))
		http.Error(w, "Failed to get trade performances", http.StatusInternalServerError)
		return
	}

	// Apply limit if provided
	if limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			http.Error(w, "Invalid limit", http.StatusBadRequest)
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// GetWinRate handles requests for win rate
func (h *AnalyticsHandler) GetWinRate(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	// Parse request parameters
	startTimeStr := query.Get("start_time")
	endTimeStr := query.Get("end_time")

	var startTime, endTime time.Time
	var err error

	// Parse start time
	if startTimeStr != "" {
		startTime, err = time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			http.Error(w, "Invalid start_time format", http.StatusBadRequest)
			return
		}
	} else {
		startTime = time.Now().AddDate(0, 0, -30)
	}

	// Parse end time
	if endTimeStr != "" {
		endTime, err = time.Parse(time.RFC3339, endTimeStr)
		if err != nil {
			http.Error(w, "Invalid end_time format", http.StatusBadRequest)
			return
		}
	} else {
		endTime = time.Now()
	}

	// Get win rate from service
	winRate, err := h.analyticsService.GetWinRate(r.Context(), startTime, endTime)
	if err != nil {
		h.logger.Error("Failed to get win rate", zap.Error(err))
		http.Error(w, "Failed to get win rate", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]float64{"win_rate": winRate})
}

// GetBalanceHistory handles requests for balance history
func (h *AnalyticsHandler) GetBalanceHistory(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	// Parse request parameters
	startTimeStr := query.Get("start_time")
	endTimeStr := query.Get("end_time")
	intervalStr := query.Get("interval")

	var startTime, endTime time.Time
	var interval time.Duration
	var err error

	// Parse start time
	if startTimeStr != "" {
		startTime, err = time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			http.Error(w, "Invalid start_time format", http.StatusBadRequest)
			return
		}
	} else {
		startTime = time.Now().AddDate(0, 0, -30)
	}

	// Parse end time
	if endTimeStr != "" {
		endTime, err = time.Parse(time.RFC3339, endTimeStr)
		if err != nil {
			http.Error(w, "Invalid end_time format", http.StatusBadRequest)
			return
		}
	} else {
		endTime = time.Now()
	}

	// Parse interval
	if intervalStr != "" {
		interval, err = time.ParseDuration(intervalStr)
		if err != nil {
			http.Error(w, "Invalid interval format", http.StatusBadRequest)
			return
		}
	} else {
		interval = 24 * time.Hour
	}

	// Get balance history from service
	balanceHistory, err := h.analyticsService.GetBalanceHistory(r.Context(), startTime, endTime, interval)
	if err != nil {
		h.logger.Error("Failed to get balance history", zap.Error(err))
		http.Error(w, "Failed to get balance history", http.StatusInternalServerError)
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// GetPerformanceBySymbol handles requests for performance by symbol
func (h *AnalyticsHandler) GetPerformanceBySymbol(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	// Parse request parameters
	startTimeStr := query.Get("start_time")
	endTimeStr := query.Get("end_time")

	var startTime, endTime time.Time
	var err error

	// Parse start time
	if startTimeStr != "" {
		startTime, err = time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			http.Error(w, "Invalid start_time format", http.StatusBadRequest)
			return
		}
	} else {
		startTime = time.Now().AddDate(0, 0, -30)
	}

	// Parse end time
	if endTimeStr != "" {
		endTime, err = time.Parse(time.RFC3339, endTimeStr)
		if err != nil {
			http.Error(w, "Invalid end_time format", http.StatusBadRequest)
			return
		}
	} else {
		endTime = time.Now()
	}

	// Get performance by symbol from service
	performance, err := h.analyticsService.GetPerformanceBySymbol(r.Context(), startTime, endTime)
	if err != nil {
		h.logger.Error("Failed to get performance by symbol", zap.Error(err))
		http.Error(w, "Failed to get performance by symbol", http.StatusInternalServerError)
		return
	}

	// Convert to response DTOs
	resp := make(map[string]response.SymbolPerformanceResponse)
	for symbol, perf := range performance {
		resp[symbol] = response.SymbolPerformanceFromModel(perf)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// GetPerformanceByReason handles requests for performance by reason
func (h *AnalyticsHandler) GetPerformanceByReason(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	// Parse request parameters
	startTimeStr := query.Get("start_time")
	endTimeStr := query.Get("end_time")

	var startTime, endTime time.Time
	var err error

	// Parse start time
	if startTimeStr != "" {
		startTime, err = time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			http.Error(w, "Invalid start_time format", http.StatusBadRequest)
			return
		}
	} else {
		startTime = time.Now().AddDate(0, 0, -30)
	}

	// Parse end time
	if endTimeStr != "" {
		endTime, err = time.Parse(time.RFC3339, endTimeStr)
		if err != nil {
			http.Error(w, "Invalid end_time format", http.StatusBadRequest)
			return
		}
	} else {
		endTime = time.Now()
	}

	// Get performance by reason from service
	performance, err := h.analyticsService.GetPerformanceByReason(r.Context(), startTime, endTime)
	if err != nil {
		h.logger.Error("Failed to get performance by reason", zap.Error(err))
		http.Error(w, "Failed to get performance by reason", http.StatusInternalServerError)
		return
	}

	// Convert to response DTOs
	resp := make(map[string]response.ReasonPerformanceResponse)
	for reason, perf := range performance {
		resp[reason] = response.ReasonPerformanceFromModel(perf)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// GetPerformanceByStrategy handles requests for performance by strategy
func (h *AnalyticsHandler) GetPerformanceByStrategy(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	// Parse request parameters
	startTimeStr := query.Get("start_time")
	endTimeStr := query.Get("end_time")

	var startTime, endTime time.Time
	var err error

	// Parse start time
	if startTimeStr != "" {
		startTime, err = time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			http.Error(w, "Invalid start_time format", http.StatusBadRequest)
			return
		}
	} else {
		startTime = time.Now().AddDate(0, 0, -30)
	}

	// Parse end time
	if endTimeStr != "" {
		endTime, err = time.Parse(time.RFC3339, endTimeStr)
		if err != nil {
			http.Error(w, "Invalid end_time format", http.StatusBadRequest)
			return
		}
	} else {
		endTime = time.Now()
	}

	// Get performance by strategy from service
	performance, err := h.analyticsService.GetPerformanceByStrategy(r.Context(), startTime, endTime)
	if err != nil {
		h.logger.Error("Failed to get performance by strategy", zap.Error(err))
		http.Error(w, "Failed to get performance by strategy", http.StatusInternalServerError)
		return
	}

	// Convert to response DTOs
	resp := make(map[string]response.StrategyPerformanceResponse)
	for strategy, perf := range performance {
		resp[strategy] = response.StrategyPerformanceFromModel(perf)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
