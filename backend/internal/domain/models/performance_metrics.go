package models

type PerformanceMetrics struct {
	TotalTrades           int
	WinningTrades         int
	LosingTrades          int
	WinRate               float64
	TotalProfitLoss       float64
	AverageProfitPerTrade float64
	LargestProfit         float64
	LargestLoss           float64
	ROI                   float64
	NumTrades             int // Alias for TotalTrades for backward compatibility
}
