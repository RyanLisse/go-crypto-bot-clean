package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"go-crypto-bot-clean/backend/internal/api/dto/response"
	"go-crypto-bot-clean/backend/internal/domain/models"
)

// PortfolioServiceInterface defines the interface for portfolio service
type PortfolioServiceInterface interface {
	GetPortfolioValue(ctx context.Context) (float64, error)
	GetActiveTrades(ctx context.Context) ([]*models.BoughtCoin, error)
	GetTradePerformance(ctx context.Context, timeRange string) (*models.PerformanceMetrics, error)
}

// PortfolioHandler handles portfolio-related API endpoints
type PortfolioHandler struct {
	portfolioService PortfolioServiceInterface
}

// NewPortfolioHandler creates a new portfolio handler
func NewPortfolioHandler(portfolioService PortfolioServiceInterface) *PortfolioHandler {
	return &PortfolioHandler{
		portfolioService: portfolioService,
	}
}

// GetPortfolioSummary godoc
// @Summary Get portfolio summary
// @Description Returns a summary of the current portfolio including total value and active trades
// @Tags portfolio
// @Accept json
// @Produce json
// @Success 200 {object} response.PortfolioSummaryResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/portfolio [get]
func (h *PortfolioHandler) GetPortfolioSummary(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get active trades
	activeTrades, err := h.portfolioService.GetActiveTrades(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response.ErrorResponse{
			Code:    "internal_error",
			Message: "Failed to get active trades",
			Details: err.Error(),
		})
		return
	}

	// Get portfolio value
	totalValue, err := h.portfolioService.GetPortfolioValue(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response.ErrorResponse{
			Code:    "internal_error",
			Message: "Failed to get portfolio value",
			Details: err.Error(),
		})
		return
	}

	// Get performance metrics
	metrics, err := h.portfolioService.GetTradePerformance(ctx, "all")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response.ErrorResponse{
			Code:    "internal_error",
			Message: "Failed to get performance metrics",
			Details: err.Error(),
		})
		return
	}

	// Build response
	resp := response.PortfolioSummaryResponse{
		TotalValue:       totalValue,
		ActiveTradeCount: len(activeTrades),
		ActiveTrades:     mapToTradeResponses(activeTrades),
		Performance:      mapToPerformanceResponse(metrics),
		Timestamp:        time.Now(),
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// GetActiveTrades godoc
// @Summary Get active trades
// @Description Returns a list of all active trading positions
// @Tags portfolio
// @Accept json
// @Produce json
// @Success 200 {object} response.ActiveTradesResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/portfolio/active [get]
func (h *PortfolioHandler) GetActiveTrades(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get active trades
	activeTrades, err := h.portfolioService.GetActiveTrades(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response.ErrorResponse{
			Code:    "internal_error",
			Message: "Failed to get active trades",
			Details: err.Error(),
		})
		return
	}

	// Build response
	resp := response.ActiveTradesResponse{
		Trades:    mapToTradeResponses(activeTrades),
		Count:     len(activeTrades),
		Timestamp: time.Now(),
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// GetPerformanceMetrics godoc
// @Summary Get performance metrics
// @Description Returns trading performance metrics for a specified time range
// @Tags portfolio
// @Accept json
// @Produce json
// @Param timeRange query string false "Time range (day, week, month, year, all)" default(all)
// @Success 200 {object} response.PerformanceResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/portfolio/performance [get]
func (h *PortfolioHandler) GetPerformanceMetrics(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	timeRange := r.URL.Query().Get("timeRange")
	if timeRange == "" {
		timeRange = "all"
	}

	// Get performance metrics
	metrics, err := h.portfolioService.GetTradePerformance(ctx, timeRange)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response.ErrorResponse{
			Code:    "internal_error",
			Message: "Failed to get performance metrics",
			Details: err.Error(),
		})
		return
	}

	// Build response
	resp := mapToPerformanceResponse(metrics)
	resp.TimeRange = timeRange

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// GetTotalValue godoc
// @Summary Get total portfolio value
// @Description Returns the total value of the portfolio in USDT
// @Tags portfolio
// @Accept json
// @Produce json
// @Success 200 {object} response.TotalValueResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/portfolio/value [get]
func (h *PortfolioHandler) GetTotalValue(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get portfolio value
	totalValue, err := h.portfolioService.GetPortfolioValue(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response.ErrorResponse{
			Code:    "internal_error",
			Message: "Failed to get portfolio value",
			Details: err.Error(),
		})
		return
	}

	// Build response
	resp := response.TotalValueResponse{
		Value:     totalValue,
		Timestamp: time.Now(),
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// Helper functions to map domain models to DTOs
func mapToTradeResponses(trades []*models.BoughtCoin) []response.TradeResponse {
	resp := make([]response.TradeResponse, len(trades))
	for i, trade := range trades {
		resp[i] = response.TradeResponse{
			ID:               uint(trade.ID),
			Symbol:           trade.Symbol,
			PurchasePrice:    trade.PurchasePrice,
			CurrentPrice:     trade.CurrentPrice,
			Quantity:         trade.Quantity,
			PurchaseTime:     trade.BoughtAt,
			ProfitPercent:    calculateProfitPercent(trade.PurchasePrice, trade.CurrentPrice),
			CurrentValue:     trade.CurrentPrice * trade.Quantity,
			StopLossPrice:    trade.StopLoss,
			TakeProfitLevels: mapToTakeProfitLevels(trade.TakeProfit),
		}
	}
	return resp
}

func mapToPerformanceResponse(metrics *models.PerformanceMetrics) response.PerformanceResponse {
	return response.PerformanceResponse{
		TotalTrades:           metrics.TotalTrades,
		WinningTrades:         metrics.WinningTrades,
		LosingTrades:          metrics.LosingTrades,
		WinRate:               metrics.WinRate,
		TotalProfitLoss:       metrics.TotalProfitLoss,
		AverageProfitPerTrade: metrics.AverageProfitPerTrade,
		LargestProfit:         metrics.LargestProfit,
		LargestLoss:           metrics.LargestLoss,
	}
}

func mapToTakeProfitLevels(takeProfit float64) []response.TakeProfitLevelResponse {
	// For now, we'll just create a single take profit level
	// In the future, this could be expanded to support multiple levels
	if takeProfit <= 0 {
		return []response.TakeProfitLevelResponse{}
	}

	return []response.TakeProfitLevelResponse{
		{
			Price:    takeProfit,
			Percent:  0, // We don't have this information yet
			Executed: false,
		},
	}
}

func calculateProfitPercent(buyPrice, currentPrice float64) float64 {
	if buyPrice == 0 {
		return 0
	}
	return ((currentPrice - buyPrice) / buyPrice) * 100
}
