package response

import (
	"time"

	"github.com/ryanlisse/go-crypto-bot/internal/domain/models"
)

// TradeAnalyticsResponse represents the response for trade analytics
type TradeAnalyticsResponse struct {
	TimeFrame           string                           `json:"time_frame"`
	StartTime           time.Time                        `json:"start_time"`
	EndTime             time.Time                        `json:"end_time"`
	TotalTrades         int                              `json:"total_trades"`
	WinningTrades       int                              `json:"winning_trades"`
	LosingTrades        int                              `json:"losing_trades"`
	WinRate             float64                          `json:"win_rate"`
	TotalProfit         float64                          `json:"total_profit"`
	TotalLoss           float64                          `json:"total_loss"`
	NetProfit           float64                          `json:"net_profit"`
	ProfitFactor        float64                          `json:"profit_factor"`
	AverageProfit       float64                          `json:"average_profit"`
	AverageLoss         float64                          `json:"average_loss"`
	LargestProfit       float64                          `json:"largest_profit"`
	LargestLoss         float64                          `json:"largest_loss"`
	MaxDrawdown         float64                          `json:"max_drawdown"`
	MaxDrawdownPercent  float64                          `json:"max_drawdown_percent"`
	SharpeRatio         float64                          `json:"sharpe_ratio"`
	SortinoRatio        float64                          `json:"sortino_ratio"`
	RiskRewardRatio     float64                          `json:"risk_reward_ratio"`
	AverageHoldingTime  string                           `json:"average_holding_time"`
	TradesPerDay        float64                          `json:"trades_per_day"`
	TradesPerWeek       float64                          `json:"trades_per_week"`
	TradesPerMonth      float64                          `json:"trades_per_month"`
	PerformanceByReason map[string]ReasonPerformanceResponse  `json:"performance_by_reason"`
	PerformanceBySymbol map[string]SymbolPerformanceResponse  `json:"performance_by_symbol"`
	PerformanceByStrategy map[string]StrategyPerformanceResponse `json:"performance_by_strategy"`
	BalanceHistory      []BalancePointResponse           `json:"balance_history"`
	EquityCurve         []EquityPointResponse            `json:"equity_curve"`
}

// TradePerformanceResponse represents the response for trade performance
type TradePerformanceResponse struct {
	TradeID         string    `json:"trade_id"`
	Symbol          string    `json:"symbol"`
	EntryTime       time.Time `json:"entry_time"`
	ExitTime        time.Time `json:"exit_time"`
	EntryPrice      float64   `json:"entry_price"`
	ExitPrice       float64   `json:"exit_price"`
	Quantity        float64   `json:"quantity"`
	ProfitLoss      float64   `json:"profit_loss"`
	ProfitLossPercent float64 `json:"profit_loss_percent"`
	HoldingTime     string    `json:"holding_time"`
	HoldingTimeMs   int64     `json:"holding_time_ms"`
	EntryReason     string    `json:"entry_reason"`
	ExitReason      string    `json:"exit_reason"`
	Strategy        string    `json:"strategy"`
	StopLoss        float64   `json:"stop_loss"`
	TakeProfit      float64   `json:"take_profit"`
	RiskRewardRatio float64   `json:"risk_reward_ratio"`
	ExpectedValue   float64   `json:"expected_value"`
	ActualRR        float64   `json:"actual_rr"`
}

// ReasonPerformanceResponse represents the response for reason performance
type ReasonPerformanceResponse struct {
	Reason          string  `json:"reason"`
	TotalTrades     int     `json:"total_trades"`
	WinningTrades   int     `json:"winning_trades"`
	LosingTrades    int     `json:"losing_trades"`
	WinRate         float64 `json:"win_rate"`
	TotalProfit     float64 `json:"total_profit"`
	AverageProfit   float64 `json:"average_profit"`
	ProfitFactor    float64 `json:"profit_factor"`
}

// SymbolPerformanceResponse represents the response for symbol performance
type SymbolPerformanceResponse struct {
	Symbol          string  `json:"symbol"`
	TotalTrades     int     `json:"total_trades"`
	WinningTrades   int     `json:"winning_trades"`
	LosingTrades    int     `json:"losing_trades"`
	WinRate         float64 `json:"win_rate"`
	TotalProfit     float64 `json:"total_profit"`
	AverageProfit   float64 `json:"average_profit"`
	ProfitFactor    float64 `json:"profit_factor"`
}

// StrategyPerformanceResponse represents the response for strategy performance
type StrategyPerformanceResponse struct {
	Strategy        string  `json:"strategy"`
	TotalTrades     int     `json:"total_trades"`
	WinningTrades   int     `json:"winning_trades"`
	LosingTrades    int     `json:"losing_trades"`
	WinRate         float64 `json:"win_rate"`
	TotalProfit     float64 `json:"total_profit"`
	AverageProfit   float64 `json:"average_profit"`
	ProfitFactor    float64 `json:"profit_factor"`
}

// BalancePointResponse represents the response for a balance point
type BalancePointResponse struct {
	Timestamp time.Time `json:"timestamp"`
	Balance   float64   `json:"balance"`
}

// EquityPointResponse represents the response for an equity point
type EquityPointResponse struct {
	Timestamp time.Time `json:"timestamp"`
	Equity    float64   `json:"equity"`
}

// TradeAnalyticsFromModel converts a model to a response DTO
func TradeAnalyticsFromModel(model *models.TradeAnalytics) TradeAnalyticsResponse {
	resp := TradeAnalyticsResponse{
		TimeFrame:          string(model.TimeFrame),
		StartTime:          model.StartTime,
		EndTime:            model.EndTime,
		TotalTrades:        model.TotalTrades,
		WinningTrades:      model.WinningTrades,
		LosingTrades:       model.LosingTrades,
		WinRate:            model.WinRate,
		TotalProfit:        model.TotalProfit,
		TotalLoss:          model.TotalLoss,
		NetProfit:          model.NetProfit,
		ProfitFactor:       model.ProfitFactor,
		AverageProfit:      model.AverageProfit,
		AverageLoss:        model.AverageLoss,
		LargestProfit:      model.LargestProfit,
		LargestLoss:        model.LargestLoss,
		MaxDrawdown:        model.MaxDrawdown,
		MaxDrawdownPercent: model.MaxDrawdownPercent,
		SharpeRatio:        model.SharpeRatio,
		SortinoRatio:       model.SortinoRatio,
		RiskRewardRatio:    model.RiskRewardRatio,
		AverageHoldingTime: model.AverageHoldingTime,
		TradesPerDay:       model.TradesPerDay,
		TradesPerWeek:      model.TradesPerWeek,
		TradesPerMonth:     model.TradesPerMonth,
		PerformanceByReason: make(map[string]ReasonPerformanceResponse),
		PerformanceBySymbol: make(map[string]SymbolPerformanceResponse),
		PerformanceByStrategy: make(map[string]StrategyPerformanceResponse),
		BalanceHistory:     make([]BalancePointResponse, len(model.BalanceHistory)),
		EquityCurve:        make([]EquityPointResponse, len(model.EquityCurve)),
	}

	// Convert performance by reason
	for reason, perf := range model.PerformanceByReason {
		resp.PerformanceByReason[reason] = ReasonPerformanceFromModel(perf)
	}

	// Convert performance by symbol
	for symbol, perf := range model.PerformanceBySymbol {
		resp.PerformanceBySymbol[symbol] = SymbolPerformanceFromModel(perf)
	}

	// Convert performance by strategy
	for strategy, perf := range model.PerformanceByStrategy {
		resp.PerformanceByStrategy[strategy] = StrategyPerformanceFromModel(perf)
	}

	// Convert balance history
	for i, point := range model.BalanceHistory {
		resp.BalanceHistory[i] = BalancePointResponse{
			Timestamp: point.Timestamp,
			Balance:   point.Balance,
		}
	}

	// Convert equity curve
	for i, point := range model.EquityCurve {
		resp.EquityCurve[i] = EquityPointResponse{
			Timestamp: point.Timestamp,
			Equity:    point.Equity,
		}
	}

	return resp
}

// TradePerformanceFromModel converts a model to a response DTO
func TradePerformanceFromModel(model *models.TradePerformance) TradePerformanceResponse {
	return TradePerformanceResponse{
		TradeID:          model.TradeID,
		Symbol:           model.Symbol,
		EntryTime:        model.EntryTime,
		ExitTime:         model.ExitTime,
		EntryPrice:       model.EntryPrice,
		ExitPrice:        model.ExitPrice,
		Quantity:         model.Quantity,
		ProfitLoss:       model.ProfitLoss,
		ProfitLossPercent: model.ProfitLossPercent,
		HoldingTime:      model.HoldingTime,
		HoldingTimeMs:    model.HoldingTimeMs,
		EntryReason:      model.EntryReason,
		ExitReason:       model.ExitReason,
		Strategy:         model.Strategy,
		StopLoss:         model.StopLoss,
		TakeProfit:       model.TakeProfit,
		RiskRewardRatio:  model.RiskRewardRatio,
		ExpectedValue:    model.ExpectedValue,
		ActualRR:         model.ActualRR,
	}
}

// ReasonPerformanceFromModel converts a model to a response DTO
func ReasonPerformanceFromModel(model models.ReasonPerformance) ReasonPerformanceResponse {
	return ReasonPerformanceResponse{
		Reason:        model.Reason,
		TotalTrades:   model.TotalTrades,
		WinningTrades: model.WinningTrades,
		LosingTrades:  model.LosingTrades,
		WinRate:       model.WinRate,
		TotalProfit:   model.TotalProfit,
		AverageProfit: model.AverageProfit,
		ProfitFactor:  model.ProfitFactor,
	}
}

// SymbolPerformanceFromModel converts a model to a response DTO
func SymbolPerformanceFromModel(model models.SymbolPerformance) SymbolPerformanceResponse {
	return SymbolPerformanceResponse{
		Symbol:        model.Symbol,
		TotalTrades:   model.TotalTrades,
		WinningTrades: model.WinningTrades,
		LosingTrades:  model.LosingTrades,
		WinRate:       model.WinRate,
		TotalProfit:   model.TotalProfit,
		AverageProfit: model.AverageProfit,
		ProfitFactor:  model.ProfitFactor,
	}
}

// StrategyPerformanceFromModel converts a model to a response DTO
func StrategyPerformanceFromModel(model models.StrategyPerformance) StrategyPerformanceResponse {
	return StrategyPerformanceResponse{
		Strategy:      model.Strategy,
		TotalTrades:   model.TotalTrades,
		WinningTrades: model.WinningTrades,
		LosingTrades:  model.LosingTrades,
		WinRate:       model.WinRate,
		TotalProfit:   model.TotalProfit,
		AverageProfit: model.AverageProfit,
		ProfitFactor:  model.ProfitFactor,
	}
}
