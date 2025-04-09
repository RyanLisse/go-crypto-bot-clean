package handlers

import (
	"net/http"
	"time"

	"go-crypto-bot-clean/backend/internal/backtest"
	"go-crypto-bot-clean/backend/internal/domain/models"

	"github.com/gin-gonic/gin"
)

// BacktestHandler handles backtest-related API requests
type BacktestHandler struct {
	backtestService *backtest.Service
}

// NewBacktestHandler creates a new backtest handler
func NewBacktestHandler(backtestService *backtest.Service) *BacktestHandler {
	return &BacktestHandler{
		backtestService: backtestService,
	}
}

// BacktestRequest represents a request to run a backtest
type BacktestRequest struct {
	Strategy       string    `json:"strategy" binding:"required"`
	Symbol         string    `json:"symbol" binding:"required"`
	Timeframe      string    `json:"timeframe" binding:"required"`
	StartDate      time.Time `json:"startDate" binding:"required"`
	EndDate        time.Time `json:"endDate" binding:"required"`
	InitialCapital float64   `json:"initialCapital" binding:"required"`
	RiskPerTrade   float64   `json:"riskPerTrade" binding:"required"`
}

// BacktestResponse represents the response from a backtest
type BacktestResponse struct {
	EquityCurve        []*backtest.EquityPoint      `json:"equityCurve"`
	DrawdownCurve      []*backtest.DrawdownPoint    `json:"drawdownCurve"`
	PerformanceMetrics *backtest.PerformanceMetrics `json:"performanceMetrics"`
	Trades             []*models.Order              `json:"trades"`
}

// RunBacktest handles the request to run a backtest
func (h *BacktestHandler) RunBacktest(c *gin.Context) {
	var req BacktestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create the configuration expected by the backtest service
	serviceConfig := &backtest.BacktestRequestConfig{
		Strategy:       req.Strategy,
		Symbol:         req.Symbol,    // Service uses single symbol
		Timeframe:      req.Timeframe, // Service uses Timeframe
		StartTime:      req.StartDate,
		EndTime:        req.EndDate,
		InitialCapital: req.InitialCapital,
		RiskPerTrade:   req.RiskPerTrade, // Pass through strategy-specific params
	}

	// Run backtest using the service configuration
	result, err := h.backtestService.RunBacktest(c.Request.Context(), serviceConfig)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create response (assuming BacktestResult structure is compatible)
	// Note: The structure of BacktestResponse might need adjustment if
	// the service's BacktestResult differs significantly.
	response := &BacktestResponse{
		EquityCurve:        result.EquityCurve,
		DrawdownCurve:      result.DrawdownCurve,
		PerformanceMetrics: result.PerformanceMetrics,
		Trades:             result.Trades,
	}

	c.JSON(http.StatusOK, response)
}

// GetBacktestResults handles the request to get backtest results
func (h *BacktestHandler) GetBacktestResults(c *gin.Context) {
	id := c.Param("id")

	// Get backtest results
	result, err := h.backtestService.GetBacktestResult(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create response
	response := &BacktestResponse{
		EquityCurve:        result.EquityCurve,
		DrawdownCurve:      result.DrawdownCurve,
		PerformanceMetrics: result.PerformanceMetrics,
		Trades:             result.Trades,
	}

	c.JSON(http.StatusOK, response)
}

// ListBacktestResults handles the request to list backtest results
func (h *BacktestHandler) ListBacktestResults(c *gin.Context) {
	// Get backtest results
	results, err := h.backtestService.ListBacktestResults(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}
