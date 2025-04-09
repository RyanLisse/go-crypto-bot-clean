package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go-crypto-bot-clean/backend/internal/backtest"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// NewVisualizeCmd creates a new visualize command
func NewVisualizeCmd() *cobra.Command {
	var (
		resultFile  string
		outputDir   string
		chartTypes  []string
		interactive bool
	)

	cmd := &cobra.Command{
		Use:   "visualize",
		Short: "Generate visualizations from backtest results",
		Long:  `Generate charts and reports from backtest results.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Set up logger
			logger, _ := zap.NewDevelopment()
			defer logger.Sync()

			// Check if result file exists
			if _, err := os.Stat(resultFile); os.IsNotExist(err) {
				return fmt.Errorf("result file not found: %s", resultFile)
			}

			// Create output directory if it doesn't exist
			if _, err := os.Stat(outputDir); os.IsNotExist(err) {
				if err := os.MkdirAll(outputDir, 0755); err != nil {
					return fmt.Errorf("failed to create output directory: %w", err)
				}
			}

			// Load backtest result from file
			result, err := loadBacktestResult(resultFile)
			if err != nil {
				return fmt.Errorf("failed to load backtest result: %w", err)
			}

			// Create performance analyzer
			analyzer := backtest.NewPerformanceAnalyzer()

			// Generate report if not already present
			if result.PerformanceMetrics == nil {
				metrics, err := analyzer.CalculateMetrics(result)
				if err != nil {
					return fmt.Errorf("failed to calculate performance metrics: %w", err)
				}
				result.PerformanceMetrics = metrics
			}

			// Generate visualizations
			if err := generateVisualizations(result, outputDir, chartTypes, interactive); err != nil {
				return fmt.Errorf("failed to generate visualizations: %w", err)
			}

			logger.Info("Visualizations generated successfully",
				zap.String("output_dir", outputDir),
				zap.Strings("chart_types", chartTypes),
			)

			return nil
		},
	}

	// Add flags
	cmd.Flags().StringVar(&resultFile, "result", "", "Path to backtest result file")
	cmd.Flags().StringVar(&outputDir, "output", "visualizations", "Output directory for visualizations")
	cmd.Flags().StringSliceVar(&chartTypes, "charts", []string{"equity", "drawdown", "monthly", "trades", "monte-carlo"}, "Chart types to generate")
	cmd.Flags().BoolVar(&interactive, "interactive", false, "Generate interactive HTML charts")

	// Mark required flags
	cmd.MarkFlagRequired("result")

	return cmd
}

// loadBacktestResult loads a backtest result from a file
func loadBacktestResult(filePath string) (*backtest.BacktestResult, error) {
	// Read file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Unmarshal JSON
	var result backtest.BacktestResult
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return &result, nil
}

// generateVisualizations generates visualizations from backtest results
func generateVisualizations(result *backtest.BacktestResult, outputDir string, chartTypes []string, interactive bool) error {
	// Create performance analyzer
	analyzer := backtest.NewPerformanceAnalyzer()

	// Generate report
	report, err := analyzer.GenerateReport(result, result.PerformanceMetrics)
	if err != nil {
		return fmt.Errorf("failed to generate report: %w", err)
	}

	// Generate charts based on chart types
	for _, chartType := range chartTypes {
		switch chartType {
		case "equity":
			if err := generateEquityCurveChart(report, outputDir, interactive); err != nil {
				return fmt.Errorf("failed to generate equity curve chart: %w", err)
			}
		case "drawdown":
			if err := generateDrawdownChart(report, outputDir, interactive); err != nil {
				return fmt.Errorf("failed to generate drawdown chart: %w", err)
			}
		case "monthly":
			if err := generateMonthlyReturnsChart(report, outputDir, interactive); err != nil {
				return fmt.Errorf("failed to generate monthly returns chart: %w", err)
			}
		case "trades":
			if err := generateTradeDistributionChart(report, outputDir, interactive); err != nil {
				return fmt.Errorf("failed to generate trade distribution chart: %w", err)
			}
		case "monte-carlo":
			if err := generateMonteCarloChart(report, outputDir, interactive); err != nil {
				return fmt.Errorf("failed to generate Monte Carlo chart: %w", err)
			}
		default:
			fmt.Printf("Unknown chart type: %s\n", chartType)
		}
	}

	// Generate summary report
	if err := generateSummaryReport(report, outputDir); err != nil {
		return fmt.Errorf("failed to generate summary report: %w", err)
	}

	return nil
}

// generateEquityCurveChart generates an equity curve chart
func generateEquityCurveChart(report *backtest.BacktestReport, outputDir string, interactive bool) error {
	// For now, just save the equity curve data to a JSON file
	// In a real implementation, this would generate a chart image or HTML file
	equityCurveData := make([]map[string]interface{}, 0, len(report.EquityCurve))
	for _, point := range report.EquityCurve {
		equityCurveData = append(equityCurveData, map[string]interface{}{
			"timestamp": point.Timestamp.Format(time.RFC3339),
			"equity":    point.Equity,
		})
	}

	// Save to file
	outputFile := filepath.Join(outputDir, "equity_curve.json")
	return saveJSONToFile(equityCurveData, outputFile)
}

// generateDrawdownChart generates a drawdown chart
func generateDrawdownChart(report *backtest.BacktestReport, outputDir string, interactive bool) error {
	// For now, just save the drawdown curve data to a JSON file
	drawdownData := make([]map[string]interface{}, 0, len(report.DrawdownCurve))
	for _, point := range report.DrawdownCurve {
		drawdownData = append(drawdownData, map[string]interface{}{
			"timestamp": point.Timestamp.Format(time.RFC3339),
			"drawdown":  point.Drawdown,
		})
	}

	// Save to file
	outputFile := filepath.Join(outputDir, "drawdown_curve.json")
	return saveJSONToFile(drawdownData, outputFile)
}

// generateMonthlyReturnsChart generates a monthly returns chart
func generateMonthlyReturnsChart(report *backtest.BacktestReport, outputDir string, interactive bool) error {
	// For now, just save the monthly returns data to a JSON file
	monthlyReturnsData := make([]map[string]interface{}, 0, len(report.MonthlyReturns))
	for month, returnValue := range report.MonthlyReturns {
		monthlyReturnsData = append(monthlyReturnsData, map[string]interface{}{
			"month":  month,
			"return": returnValue,
		})
	}

	// Save to file
	outputFile := filepath.Join(outputDir, "monthly_returns.json")
	return saveJSONToFile(monthlyReturnsData, outputFile)
}

// generateTradeDistributionChart generates a trade distribution chart
func generateTradeDistributionChart(report *backtest.BacktestReport, outputDir string, interactive bool) error {
	// For now, just save the trade statistics to a JSON file
	tradeData := map[string]interface{}{
		"winningTrades":      report.Metrics.WinningTrades,
		"losingTrades":       report.Metrics.LosingTrades,
		"breakEvenTrades":    report.Metrics.BreakEvenTrades,
		"winRate":            report.Metrics.WinRate,
		"averageProfitTrade": report.Metrics.AverageProfitTrade,
		"averageLossTrade":   report.Metrics.AverageLossTrade,
		"largestProfitTrade": report.Metrics.LargestProfitTrade,
		"largestLossTrade":   report.Metrics.LargestLossTrade,
	}

	// Save to file
	outputFile := filepath.Join(outputDir, "trade_distribution.json")
	return saveJSONToFile(tradeData, outputFile)
}

// generateMonteCarloChart generates a Monte Carlo simulation chart
func generateMonteCarloChart(report *backtest.BacktestReport, outputDir string, interactive bool) error {
	// For now, just save the Monte Carlo simulation data to a JSON file
	if len(report.MonteCarloSimulations) == 0 {
		// If no Monte Carlo simulations were run, run them now
		analyzer := backtest.NewPerformanceAnalyzer()
		result := &backtest.BacktestResult{
			InitialCapital: 10000, // Default value, should be replaced with actual value
			EquityCurve:    report.EquityCurve,
		}
		simulations, err := analyzer.RunMonteCarloSimulation(result, 100)
		if err != nil {
			return fmt.Errorf("failed to run Monte Carlo simulation: %w", err)
		}
		report.MonteCarloSimulations = simulations
	}

	// Save to file
	outputFile := filepath.Join(outputDir, "monte_carlo.json")
	return saveJSONToFile(report.MonteCarloSimulations, outputFile)
}

// generateSummaryReport generates a summary report
func generateSummaryReport(report *backtest.BacktestReport, outputDir string) error {
	// Create summary data
	summary := map[string]interface{}{
		"totalReturn":        report.Metrics.TotalReturn,
		"annualizedReturn":   report.Metrics.AnnualizedReturn,
		"sharpeRatio":        report.Metrics.SharpeRatio,
		"sortinoRatio":       report.Metrics.SortinoRatio,
		"maxDrawdown":        report.Metrics.MaxDrawdown,
		"maxDrawdownPercent": report.Metrics.MaxDrawdownPercent,
		"winRate":            report.Metrics.WinRate,
		"profitFactor":       report.Metrics.ProfitFactor,
		"totalTrades":        report.Metrics.TotalTrades,
		"winningTrades":      report.Metrics.WinningTrades,
		"losingTrades":       report.Metrics.LosingTrades,
		"calmarRatio":        report.Metrics.CalmarRatio,
		"omegaRatio":         report.Metrics.OmegaRatio,
		"informationRatio":   report.Metrics.InformationRatio,
	}

	// Add trade statistics
	if report.TradeStats != nil {
		summary["tradeStats"] = map[string]interface{}{
			"consecutiveWins":    report.TradeStats.ConsecutiveWins,
			"consecutiveLosses":  report.TradeStats.ConsecutiveLosses,
			"profitableMonths":   report.TradeStats.ProfitableMonths,
			"unprofitableMonths": report.TradeStats.UnprofitableMonths,
			"bestMonth":          report.TradeStats.BestMonth,
			"worstMonth":         report.TradeStats.WorstMonth,
			"valueAtRisk":        report.TradeStats.ValueAtRisk,
			"conditionalVaR":     report.TradeStats.ConditionalVaR,
		}
	}

	// Save to file
	outputFile := filepath.Join(outputDir, "summary.json")
	return saveJSONToFile(summary, outputFile)
}

// saveJSONToFile saves data as JSON to a file
func saveJSONToFile(data interface{}, filePath string) error {
	// Marshal to JSON
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	// Write to file
	if err := os.WriteFile(filePath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
