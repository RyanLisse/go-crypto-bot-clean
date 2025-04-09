package backtest

import (
	"math"
	"math/rand"
	"sort"
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
	// Advanced metrics
	CalmarRatio      float64 // Annualized Return / Max Drawdown
	OmegaRatio       float64 // Probability weighted ratio of gains vs losses
	InformationRatio float64 // Excess return per unit of risk
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
	Metrics               *PerformanceMetrics
	MonthlyReturns        map[string]float64
	TradeStats            *TradeStats
	EquityCurve           []*EquityPoint
	DrawdownCurve         []*DrawdownPoint
	MonteCarloSimulations [][]float64
}

// TradeStats contains statistics about trades
type TradeStats struct {
	AverageWin          float64
	AverageLoss         float64
	LargestWin          float64
	LargestLoss         float64
	AverageHoldingTime  time.Duration
	MedianHoldingTime   time.Duration
	WinningHoldingTime  time.Duration
	LosingHoldingTime   time.Duration
	ConsecutiveWins     int
	ConsecutiveLosses   int
	ProfitableMonths    int
	UnprofitableMonths  int
	BestMonth           float64
	WorstMonth          float64
	StandardDeviation   float64
	DownsideDeviation   float64
	SharpeRatio         float64
	SortinoRatio        float64
	CalmarRatio         float64
	MAR                 float64
	OmegaRatio          float64
	TailRatio           float64
	ValueAtRisk         float64
	ConditionalVaR      float64
	ExpectedShortfall   float64
	UlcerIndex          float64
	PainIndex           float64
	PainRatio           float64
	BurkesRatio         float64
	SterlingRatio       float64
	KellyRatio          float64
	InformationRatio    float64
	TreynorRatio        float64
	JensenAlpha         float64
	Beta                float64
	R2                  float64
	TrackingError       float64
	TrackingErrorAnnual float64
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

	// Calculate Calmar ratio (annualized return / maximum drawdown)
	if metrics.MaxDrawdownPercent > 0 {
		metrics.CalmarRatio = metrics.AnnualizedReturn / metrics.MaxDrawdownPercent
	}

	// Calculate Omega ratio (probability weighted ratio of gains versus losses)
	positiveProfits := 0.0
	negativeProfits := 0.0
	for _, r := range dailyReturns {
		if r >= 0 {
			positiveProfits += r
		} else {
			negativeProfits += math.Abs(r)
		}
	}
	if negativeProfits > 0 {
		metrics.OmegaRatio = positiveProfits / negativeProfits
	} else if positiveProfits > 0 {
		metrics.OmegaRatio = positiveProfits
	}

	// Calculate Information ratio (excess return per unit of risk)
	// For simplicity, we'll use a benchmark return of 0%
	benchmarkReturn := 0.0
	excessReturn := meanReturn - benchmarkReturn
	trackingError := stdDev
	if trackingError > 0 {
		metrics.InformationRatio = excessReturn / trackingError * math.Sqrt(252)
	}

	return metrics, nil
}

// GenerateReport generates a detailed performance report
func (a *DefaultPerformanceAnalyzer) GenerateReport(result *BacktestResult, metrics *PerformanceMetrics) (*BacktestReport, error) {
	// Create the report structure
	report := &BacktestReport{
		Metrics:        metrics,
		MonthlyReturns: make(map[string]float64),
		EquityCurve:    result.EquityCurve,
		DrawdownCurve:  result.DrawdownCurve,
	}

	// Calculate monthly returns using the dedicated method
	monthlyReturns, err := a.CalculateMonthlyReturns(result)
	if err != nil {
		return nil, err
	}
	report.MonthlyReturns = monthlyReturns

	// Calculate detailed trade statistics
	report.TradeStats = a.calculateTradeStats(result, metrics)

	// Run Monte Carlo simulation (100 simulations)
	simulations, err := a.RunMonteCarloSimulation(result, 100)
	if err == nil {
		report.MonteCarloSimulations = simulations
	}

	return report, nil
}

// calculateTradeStats calculates detailed statistics about trades
func (a *DefaultPerformanceAnalyzer) calculateTradeStats(result *BacktestResult, metrics *PerformanceMetrics) *TradeStats {
	stats := &TradeStats{
		AverageWin:         metrics.AverageProfitTrade,
		AverageLoss:        metrics.AverageLossTrade,
		LargestWin:         metrics.LargestProfitTrade,
		LargestLoss:        metrics.LargestLossTrade,
		AverageHoldingTime: metrics.AverageHoldingTime,
		SharpeRatio:        metrics.SharpeRatio,
		SortinoRatio:       metrics.SortinoRatio,
		CalmarRatio:        metrics.CalmarRatio,
		OmegaRatio:         metrics.OmegaRatio,
		InformationRatio:   metrics.InformationRatio,
	}

	// Calculate consecutive wins/losses
	var currentConsecutiveWins, currentConsecutiveLosses int
	var maxConsecutiveWins, maxConsecutiveLosses int
	var prevTradeWasWin bool
	var firstTrade = true

	// Calculate holding times for winning and losing trades separately
	var totalWinningHoldingTime, totalLosingHoldingTime time.Duration
	var winningTradeCount, losingTradeCount int

	// Collect all holding times for median calculation
	holdingTimes := make([]time.Duration, 0, len(result.ClosedPositions))

	for _, position := range result.ClosedPositions {
		holdingTime := position.CloseTime.Sub(position.OpenTime)
		holdingTimes = append(holdingTimes, holdingTime)

		isWin := position.Profit > 0

		// Track consecutive wins/losses
		if firstTrade {
			firstTrade = false
			prevTradeWasWin = isWin
			if isWin {
				currentConsecutiveWins = 1
			} else {
				currentConsecutiveLosses = 1
			}
		} else {
			if isWin {
				if prevTradeWasWin {
					currentConsecutiveWins++
				} else {
					currentConsecutiveWins = 1
					prevTradeWasWin = true
				}
				currentConsecutiveLosses = 0
			} else {
				if !prevTradeWasWin {
					currentConsecutiveLosses++
				} else {
					currentConsecutiveLosses = 1
					prevTradeWasWin = false
				}
				currentConsecutiveWins = 0
			}
		}

		// Update max consecutive wins/losses
		if currentConsecutiveWins > maxConsecutiveWins {
			maxConsecutiveWins = currentConsecutiveWins
		}
		if currentConsecutiveLosses > maxConsecutiveLosses {
			maxConsecutiveLosses = currentConsecutiveLosses
		}

		// Track holding times by trade outcome
		if isWin {
			totalWinningHoldingTime += holdingTime
			winningTradeCount++
		} else {
			totalLosingHoldingTime += holdingTime
			losingTradeCount++
		}
	}

	// Calculate average holding times by trade outcome
	if winningTradeCount > 0 {
		stats.WinningHoldingTime = totalWinningHoldingTime / time.Duration(winningTradeCount)
	}
	if losingTradeCount > 0 {
		stats.LosingHoldingTime = totalLosingHoldingTime / time.Duration(losingTradeCount)
	}

	// Calculate median holding time
	if len(holdingTimes) > 0 {
		sort.Slice(holdingTimes, func(i, j int) bool {
			return holdingTimes[i] < holdingTimes[j]
		})
		stats.MedianHoldingTime = holdingTimes[len(holdingTimes)/2]
	}

	// Set consecutive wins/losses
	stats.ConsecutiveWins = maxConsecutiveWins
	stats.ConsecutiveLosses = maxConsecutiveLosses

	// Calculate monthly statistics
	monthlyReturns, _ := a.CalculateMonthlyReturns(result)
	var profitableMonths, unprofitableMonths int
	var bestMonth, worstMonth float64
	var firstMonth = true

	for _, monthReturn := range monthlyReturns {
		if firstMonth {
			bestMonth = monthReturn
			worstMonth = monthReturn
			firstMonth = false
		} else {
			if monthReturn > bestMonth {
				bestMonth = monthReturn
			}
			if monthReturn < worstMonth {
				worstMonth = monthReturn
			}
		}

		if monthReturn > 0 {
			profitableMonths++
		} else {
			unprofitableMonths++
		}
	}

	stats.ProfitableMonths = profitableMonths
	stats.UnprofitableMonths = unprofitableMonths
	stats.BestMonth = bestMonth
	stats.WorstMonth = worstMonth

	// Calculate Value at Risk (VaR) - 95% confidence level
	dailyReturns := make([]float64, 0)
	prevEquity := result.InitialCapital
	for _, point := range result.EquityCurve {
		if prevEquity > 0 {
			dailyReturn := (point.Equity - prevEquity) / prevEquity
			dailyReturns = append(dailyReturns, dailyReturn)
		}
		prevEquity = point.Equity
	}

	if len(dailyReturns) > 0 {
		sort.Float64s(dailyReturns)
		varIndex := int(float64(len(dailyReturns)) * 0.05)
		if varIndex < len(dailyReturns) {
			stats.ValueAtRisk = -dailyReturns[varIndex] * 100 // Convert to percentage
		}

		// Calculate Conditional VaR (Expected Shortfall)
		var totalCVaR float64
		for i := 0; i < varIndex; i++ {
			totalCVaR += dailyReturns[i]
		}
		if varIndex > 0 {
			stats.ConditionalVaR = -totalCVaR / float64(varIndex) * 100 // Convert to percentage
			stats.ExpectedShortfall = stats.ConditionalVaR
		}
	}

	return stats
}

// GenerateEquityCurve generates an equity curve from backtest results
func (a *DefaultPerformanceAnalyzer) GenerateEquityCurve(result *BacktestResult) ([]*EquityPoint, error) {
	return result.EquityCurve, nil
}

// GenerateDrawdownCurve generates a drawdown curve from backtest results
func (a *DefaultPerformanceAnalyzer) GenerateDrawdownCurve(result *BacktestResult) ([]*DrawdownPoint, error) {
	return result.DrawdownCurve, nil
}

// CalculateMonthlyReturns calculates monthly returns from equity curve
func (a *DefaultPerformanceAnalyzer) CalculateMonthlyReturns(result *BacktestResult) (map[string]float64, error) {
	monthlyReturns := make(map[string]float64)
	monthlyEquity := make(map[string]float64)

	// Get the last equity value for each month
	for _, point := range result.EquityCurve {
		month := point.Timestamp.Format("2006-01")
		monthlyEquity[month] = point.Equity
	}

	// Sort months chronologically
	months := make([]string, 0, len(monthlyEquity))
	for month := range monthlyEquity {
		months = append(months, month)
	}
	sort.Strings(months)

	// Calculate returns for each month
	prevEquity := result.InitialCapital
	for _, month := range months {
		equity := monthlyEquity[month]
		monthlyReturn := (equity - prevEquity) / prevEquity * 100
		monthlyReturns[month] = monthlyReturn
		prevEquity = equity
	}

	return monthlyReturns, nil
}

// RunMonteCarloSimulation performs Monte Carlo simulation on the backtest results
func (a *DefaultPerformanceAnalyzer) RunMonteCarloSimulation(result *BacktestResult, numSimulations int) ([][]float64, error) {
	// Extract daily returns from equity curve
	dailyReturns := make([]float64, 0)
	prevEquity := result.InitialCapital

	for _, point := range result.EquityCurve {
		if prevEquity > 0 {
			dailyReturn := (point.Equity - prevEquity) / prevEquity
			dailyReturns = append(dailyReturns, dailyReturn)
		}
		prevEquity = point.Equity
	}

	// Run simulations
	simulations := make([][]float64, numSimulations)
	for i := 0; i < numSimulations; i++ {
		// Initialize simulation with initial capital
		simulation := make([]float64, len(dailyReturns)+1)
		simulation[0] = result.InitialCapital

		// Shuffle returns for randomization
		shuffledReturns := make([]float64, len(dailyReturns))
		copy(shuffledReturns, dailyReturns)
		for j := len(shuffledReturns) - 1; j > 0; j-- {
			k := int(math.Floor(float64(j+1) * rand.Float64()))
			shuffledReturns[j], shuffledReturns[k] = shuffledReturns[k], shuffledReturns[j]
		}

		// Apply returns to generate equity curve
		for j := 0; j < len(shuffledReturns); j++ {
			simulation[j+1] = simulation[j] * (1 + shuffledReturns[j])
		}

		simulations[i] = simulation
	}

	return simulations, nil
}

// AnalyzeRegimes analyzes performance in different market regimes
func (a *DefaultPerformanceAnalyzer) AnalyzeRegimes(result *BacktestResult, benchmarkReturns map[string]float64) (map[string]map[string]float64, error) {
	// Define regimes based on benchmark returns
	// For example: Bull market (>1%), Sideways (-1% to 1%), Bear market (<-1%)
	regimes := map[string]map[string]float64{
		"bull":     make(map[string]float64),
		"sideways": make(map[string]float64),
		"bear":     make(map[string]float64),
	}

	// Calculate monthly returns for the strategy
	monthlyReturns, err := a.CalculateMonthlyReturns(result)
	if err != nil {
		return nil, err
	}

	// Categorize returns by regime
	for month, benchmarkReturn := range benchmarkReturns {
		strategyReturn, exists := monthlyReturns[month]
		if !exists {
			continue
		}

		if benchmarkReturn > 1.0 {
			regimes["bull"][month] = strategyReturn
		} else if benchmarkReturn < -1.0 {
			regimes["bear"][month] = strategyReturn
		} else {
			regimes["sideways"][month] = strategyReturn
		}
	}

	return regimes, nil
}

// CalculateCorrelation calculates correlation between strategy returns and benchmark returns
func (a *DefaultPerformanceAnalyzer) CalculateCorrelation(result *BacktestResult, benchmarkReturns map[string]float64) (float64, error) {
	// Calculate monthly returns for the strategy
	monthlyReturns, err := a.CalculateMonthlyReturns(result)
	if err != nil {
		return 0, err
	}

	// Extract matching months
	var strategyReturnValues []float64
	var benchmarkReturnValues []float64

	for month, benchmarkReturn := range benchmarkReturns {
		strategyReturn, exists := monthlyReturns[month]
		if exists {
			strategyReturnValues = append(strategyReturnValues, strategyReturn)
			benchmarkReturnValues = append(benchmarkReturnValues, benchmarkReturn)
		}
	}

	// Calculate correlation
	if len(strategyReturnValues) < 2 {
		return 0, nil
	}

	return calculateCorrelation(strategyReturnValues, benchmarkReturnValues), nil
}

// Helper function to calculate correlation between two series
func calculateCorrelation(x, y []float64) float64 {
	if len(x) != len(y) || len(x) == 0 {
		return 0
	}

	n := float64(len(x))
	var sumX, sumY, sumXY, sumX2, sumY2 float64

	for i := 0; i < len(x); i++ {
		sumX += x[i]
		sumY += y[i]
		sumXY += x[i] * y[i]
		sumX2 += x[i] * x[i]
		sumY2 += y[i] * y[i]
	}

	numerator := sumXY - (sumX * sumY / n)
	denominator := math.Sqrt((sumX2 - (sumX * sumX / n)) * (sumY2 - (sumY * sumY / n)))

	if denominator == 0 {
		return 0
	}

	return numerator / denominator
}
