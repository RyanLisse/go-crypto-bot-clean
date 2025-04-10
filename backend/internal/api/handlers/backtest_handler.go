package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"go-crypto-bot-clean/backend/internal/backtest"
	"go-crypto-bot-clean/backend/internal/domain/models"

	"github.com/go-chi/chi/v5"
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
func (h *BacktestHandler) RunBacktest(w http.ResponseWriter, r *http.Request) {
	var req BacktestRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": err.Error()})
		return
	}

	serviceConfig := &backtest.BacktestRequestConfig{
		Strategy:       req.Strategy,
		Symbol:         req.Symbol,
		Timeframe:      req.Timeframe,
		StartTime:      req.StartDate,
		EndTime:        req.EndDate,
		InitialCapital: req.InitialCapital,
		RiskPerTrade:   req.RiskPerTrade,
	}

	result, err := h.backtestService.RunBacktest(r.Context(), serviceConfig)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": err.Error()})
		return
	}

	response := &BacktestResponse{
		EquityCurve:        result.EquityCurve,
		DrawdownCurve:      result.DrawdownCurve,
		PerformanceMetrics: result.PerformanceMetrics,
		Trades:             result.Trades,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// GetBacktestResults handles the request to get backtest results
func (h *BacktestHandler) GetBacktestResults(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	result, err := h.backtestService.GetBacktestResult(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": err.Error()})
		return
	}

	response := &BacktestResponse{
		EquityCurve:        result.EquityCurve,
		DrawdownCurve:      result.DrawdownCurve,
		PerformanceMetrics: result.PerformanceMetrics,
		Trades:             result.Trades,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// ListBacktestResults handles the request to list backtest results
func (h *BacktestHandler) ListBacktestResults(w http.ResponseWriter, r *http.Request) {
	results, err := h.backtestService.ListBacktestResults(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(results)
}
