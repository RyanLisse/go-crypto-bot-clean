package backtest

import (
	"math"
	"time"
)

// PerformanceMetrics contains performance metrics for a backtest
type PerformanceMetrics struct {
	TotalReturn        float64
	AnnualizedReturn   float64
	SharpeRatio        float64
	SortinoRatio       float64
	MaxDrawdown        float64
	MaxDrawdownPercent float64
	WinRate            float64
	ProfitFactor       float64
	ExpectedPayoff     float64
	TotalTrades        int
	WinningTrades      int
	LosingTrades       int
	BreakEvenTrades    int
	AverageProfitTrade float64
	AverageLossTrade   float64
	LargestProfitTrade float64
	LargestLossTrade   float64
	AverageHoldingTime time.Duration
}

// PerformanceAnalyzer defines the interface for analyzing backtest performance
type PerformanceAnalyzer interface {
	// CalculateMetrics calculates performance metrics from backtest results
	CalculateMetrics(result *BacktestResult) (*PerformanceMetrics, error)

	// GenerateReport generates a detailed performance report
	GenerateReport(result *BacktestResult, metrics *PerformanceMetrics) (*BacktestReport, error)

	// GenerateEquityCurve generates an equity curve from backtest results
	GenerateEquityCurve(result *BacktestResult) ([]*EquityPoint, error)

	// GenerateDrawdownCurve generates a drawdown curve from backtest results
	GenerateDrawdownCurve(result *BacktestResult) ([]*DrawdownPoint, error)
}

// BacktestReport contains a detailed report of backtest performance
type BacktestReport struct {
	Metrics       *PerformanceMetrics
	MonthlyReturns map[string]float64
	TradeStats    *TradeStats
	EquityCurve   []*EquityPoint
	DrawdownCurve []*DrawdownPoint
}

// TradeStats contains statistics about trades
type TradeStats struct {
	AverageWin           float64
	AverageLoss          float64
	LargestWin           float64
	LargestLoss          float64
	AverageHoldingTime   time.Duration
	MedianHoldingTime    time.Duration
	WinningHoldingTime   time.Duration
	LosingHoldingTime    time.Duration
	ConsecutiveWins      int
	ConsecutiveLosses    int
	ProfitableMonths     int
	UnprofitableMonths   int
	BestMonth            float64
	WorstMonth           float64
	StandardDeviation    float64
	DownsideDeviation    float64
	SharpeRatio          float64
	SortinoRatio         float64
	CalmarRatio          float64
	MAR                  float64
	OmegaRatio           float64
	TailRatio            float64
	ValueAtRisk          float64
	ConditionalVaR       float64
	ExpectedShortfall    float64
	UlcerIndex           float64
	PainIndex            float64
	PainRatio            float64
	BurkesRatio          float64
	SterlingRatio        float64
	KellyRatio           float64
	InformationRatio     float64
	TreynorRatio         float64
	JensenAlpha          float64
	Beta                 float64
	R2                   float64
	TrackingError        float64
	TrackingErrorAnnual  float64
}

// DefaultPerformanceAnalyzer implements the PerformanceAnalyzer interface
type DefaultPerformanceAnalyzer struct{}

// NewPerformanceAnalyzer creates a new DefaultPerformanceAnalyzer
func NewPerformanceAnalyzer() *DefaultPerformanceAnalyzer {
	return &DefaultPerformanceAnalyzer{}
}

// CalculateMetrics calculates performance metrics from backtest results
func (a *DefaultPerformanceAnalyzer) CalculateMetrics(result *BacktestResult) (*PerformanceMetrics, error) {
	metrics := &PerformanceMetrics{}

	// Calculate total return
	metrics.TotalReturn = (result.FinalCapital - result.InitialCapital) / result.InitialCapital * 100

	// Calculate annualized return
	years := result.EndTime.Sub(result.StartTime).Hours() / 24 / 365
	metrics.AnnualizedReturn = math.Pow(1+metrics.TotalReturn/100, 1/years) - 1
	metrics.AnnualizedReturn *= 100

	// Calculate trade statistics
	metrics.TotalTrades = len(result.ClosedPositions)
	metrics.WinningTrades = 0
	metrics.LosingTrades = 0
	metrics.BreakEvenTrades = 0

	totalProfit := 0.0
	totalLoss := 0.0
	metrics.LargestProfitTrade = 0.0
	metrics.LargestLossTrade = 0.0
	totalHoldingTime := time.Duration(0)

	for _, position := range result.ClosedPositions {
		// Calculate holding time
		holdingTime := position.CloseTime.Sub(position.OpenTime)
		totalHoldingTime += holdingTime

		// Categorize trade
		if position.Profit > 0 {
			metrics.WinningTrades++
			totalProfit += position.Profit
			if position.Profit > metrics.LargestProfitTrade {
				metrics.LargestProfitTrade = position.Profit
			}
		} else if position.Profit < 0 {
			metrics.LosingTrades++
			totalLoss += math.Abs(position.Profit)
			if position.Profit < metrics.LargestLossTrade {
				metrics.LargestLossTrade = position.Profit
			}
		} else {
			metrics.BreakEvenTrades++
		}
	}

	// Calculate win rate
	if metrics.TotalTrades > 0 {
		metrics.WinRate = float64(metrics.WinningTrades) / float64(metrics.TotalTrades) * 100
	}

	// Calculate profit factor
	if totalLoss > 0 {
		metrics.ProfitFactor = totalProfit / totalLoss
	} else {
		metrics.ProfitFactor = totalProfit
	}

	// Calculate average profit/loss
	if metrics.WinningTrades > 0 {
		metrics.AverageProfitTrade = totalProfit / float64(metrics.WinningTrades)
	}
	if metrics.LosingTrades > 0 {
		metrics.AverageLossTrade = totalLoss / float64(metrics.LosingTrades)
	}

	// Calculate expected payoff
	if metrics.TotalTrades > 0 {
		metrics.ExpectedPayoff = (totalProfit - totalLoss) / float64(metrics.TotalTrades)
	}

	// Calculate average holding time
	if metrics.TotalTrades > 0 {
		metrics.AverageHoldingTime = totalHoldingTime / time.Duration(metrics.TotalTrades)
	}

	// Calculate maximum drawdown
	maxDrawdown := 0.0
	maxDrawdownPercent := 0.0
	highWaterMark := result.InitialCapital
	
	for _, point := range result.EquityCurve {
		if point.Equity > highWaterMark {
			highWaterMark = point.Equity
		}
		
		drawdown := highWaterMark - point.Equity
		drawdownPercent := drawdown / highWaterMark * 100
		
		if drawdown > maxDrawdown {
			maxDrawdown = drawdown
		}
		
		if drawdownPercent > maxDrawdownPercent {
			maxDrawdownPercent = drawdownPercent
		}
	}
	
	metrics.MaxDrawdown = maxDrawdown
	metrics.MaxDrawdownPercent = maxDrawdownPercent

	// Calculate Sharpe ratio (assuming risk-free rate of 0%)
	// First, calculate daily returns
	dailyReturns := make([]float64, 0)
	prevEquity := result.InitialCapital
	
	for _, point := range result.EquityCurve {
		if prevEquity > 0 {
			dailyReturn := (point.Equity - prevEquity) / prevEquity
			dailyReturns = append(dailyReturns, dailyReturn)
		}
		prevEquity = point.Equity
	}
	
	// Calculate mean and standard deviation of daily returns
	meanReturn := 0.0
	for _, r := range dailyReturns {
		meanReturn += r
	}
	if len(dailyReturns) > 0 {
		meanReturn /= float64(len(dailyReturns))
	}
	
	variance := 0.0
	for _, r := range dailyReturns {
		variance += math.Pow(r-meanReturn, 2)
	}
	if len(dailyReturns) > 1 {
		variance /= float64(len(dailyReturns) - 1)
	}
	
	stdDev := math.Sqrt(variance)
	
	// Calculate Sharpe ratio (annualized)
	if stdDev > 0 {
		metrics.SharpeRatio = meanReturn / stdDev * math.Sqrt(252) // Assuming 252 trading days in a year
	}
	
	// Calculate Sortino ratio (using downside deviation)
	downsideVariance := 0.0
	downsideCount := 0
	
	for _, r := range dailyReturns {
		if r < 0 {
			downsideVariance += math.Pow(r, 2)
			downsideCount++
		}
	}
	
	if downsideCount > 0 {
		downsideVariance /= float64(downsideCount)
	}
	
	downsideDeviation := math.Sqrt(downsideVariance)
	
	if downsideDeviation > 0 {
		metrics.SortinoRatio = meanReturn / downsideDeviation * math.Sqrt(252) // Assuming 252 trading days in a year
	}

	return metrics, nil
}

// GenerateReport generates a detailed performance report
func (a *DefaultPerformanceAnalyzer) GenerateReport(result *BacktestResult, metrics *PerformanceMetrics) (*BacktestReport, error) {
	// This is a simplified implementation
	report := &BacktestReport{
		Metrics:       metrics,
		MonthlyReturns: make(map[string]float64),
		EquityCurve:   result.EquityCurve,
		DrawdownCurve: result.DrawdownCurve,
	}
	
	// Calculate monthly returns
	monthlyEquity := make(map[string]float64)
	
	for _, point := range result.EquityCurve {
		month := point.Timestamp.Format("2006-01")
		monthlyEquity[month] = point.Equity
	}
	
	// Calculate returns for each month
	prevMonth := ""
	prevEquity := result.InitialCapital
	
	for month, equity := range monthlyEquity {
		if prevMonth != "" {
			monthlyReturn := (equity - prevEquity) / prevEquity * 100
			report.MonthlyReturns[month] = monthlyReturn
		}
		prevMonth = month
		prevEquity = equity
	}
	
	// Calculate trade statistics
	// This would be a more detailed analysis of trades
	// For now, we'll just create a basic TradeStats object
	report.TradeStats = &TradeStats{
		AverageWin:         metrics.AverageProfitTrade,
		AverageLoss:        metrics.AverageLossTrade,
		LargestWin:         metrics.LargestProfitTrade,
		LargestLoss:        metrics.LargestLossTrade,
		AverageHoldingTime: metrics.AverageHoldingTime,
		SharpeRatio:        metrics.SharpeRatio,
		SortinoRatio:       metrics.SortinoRatio,
	}
	
	return report, nil
}

// GenerateEquityCurve generates an equity curve from backtest results
func (a *DefaultPerformanceAnalyzer) GenerateEquityCurve(result *BacktestResult) ([]*EquityPoint, error) {
	return result.EquityCurve, nil
}

// GenerateDrawdownCurve generates a drawdown curve from backtest results
func (a *DefaultPerformanceAnalyzer) GenerateDrawdownCurve(result *BacktestResult) ([]*DrawdownPoint, error) {
	return result.DrawdownCurve, nil
}
