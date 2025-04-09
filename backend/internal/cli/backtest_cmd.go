package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/ryanlisse/go-crypto-bot/internal/backtest"
	"github.com/ryanlisse/go-crypto-bot/internal/backtest/strategies"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// NewBacktestCmd creates a new backtest command
func NewBacktestCmd() *cobra.Command {
	var (
		strategyName   string
		symbols        []string
		startDate      string
		endDate        string
		initialCapital float64
		interval       string
		dataDir        string
		shortPeriod    int
		longPeriod     int
		outputFile     string
	)

	cmd := &cobra.Command{
		Use:   "backtest",
		Short: "Run a backtest for a trading strategy",
		Long:  `Run a backtest for a trading strategy against historical market data.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Parse dates
			start, err := time.Parse("2006-01-02", startDate)
			if err != nil {
				return fmt.Errorf("invalid start date: %w", err)
			}

			end, err := time.Parse("2006-01-02", endDate)
			if err != nil {
				return fmt.Errorf("invalid end date: %w", err)
			}

			// Set up logger
			logger, _ := zap.NewDevelopment()
			defer logger.Sync()

			// Create data provider
			var dataProvider backtest.DataProvider
			if dataDir != "" {
				dataProvider = backtest.NewCSVDataProvider(dataDir)
			} else {
				// Use in-memory data provider for testing
				dataProvider = backtest.NewInMemoryDataProvider()
				// TODO: Load test data
			}

			// Create strategy
			var strategy backtest.BacktestStrategy
			switch strategyName {
			case "simple_ma":
				strategy = strategies.NewSimpleMAStrategy(shortPeriod, longPeriod, logger)
			default:
				return fmt.Errorf("unknown strategy: %s", strategyName)
			}

			// Create slippage model
			slippageModel := backtest.NewFixedSlippage(0.1) // 0.1% slippage

			// Create backtest config
			config := &backtest.BacktestConfig{
				StartTime:          start,
				EndTime:            end,
				InitialCapital:     initialCapital,
				Symbols:            symbols,
				Interval:           interval,
				CommissionRate:     0.001, // 0.1% commission
				SlippageModel:      slippageModel,
				EnableShortSelling: false,
				DataProvider:       dataProvider,
				Strategy:           strategy,
				Logger:             logger,
			}

			// Create backtest engine
			engine := backtest.NewEngine(config)

			// Run backtest
			logger.Info("Running backtest...",
				zap.String("strategy", strategyName),
				zap.Strings("symbols", symbols),
				zap.String("start_date", startDate),
				zap.String("end_date", endDate),
				zap.Float64("initial_capital", initialCapital),
				zap.String("interval", interval),
			)

			result, err := engine.Run(context.Background())
			if err != nil {
				return fmt.Errorf("backtest failed: %w", err)
			}

			// Print results
			printBacktestResults(result)

			// Save results to file if specified
			if outputFile != "" {
				err = saveBacktestResults(result, outputFile)
				if err != nil {
					return fmt.Errorf("failed to save results: %w", err)
				}
				logger.Info("Results saved to file", zap.String("file", outputFile))
			}

			return nil
		},
	}

	// Add flags
	cmd.Flags().StringVar(&strategyName, "strategy", "simple_ma", "Strategy to backtest (simple_ma)")
	cmd.Flags().StringSliceVar(&symbols, "symbols", []string{"BTCUSDT"}, "Symbols to backtest")
	cmd.Flags().StringVar(&startDate, "start", "2023-01-01", "Start date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&endDate, "end", "2023-12-31", "End date (YYYY-MM-DD)")
	cmd.Flags().Float64Var(&initialCapital, "capital", 10000, "Initial capital")
	cmd.Flags().StringVar(&interval, "interval", "1h", "Candle interval (1m, 5m, 15m, 1h, 4h, 1d)")
	cmd.Flags().StringVar(&dataDir, "data-dir", "", "Directory containing historical data CSV files")
	cmd.Flags().IntVar(&shortPeriod, "short-period", 10, "Short period for MA strategy")
	cmd.Flags().IntVar(&longPeriod, "long-period", 50, "Long period for MA strategy")
	cmd.Flags().StringVar(&outputFile, "output", "", "Output file for backtest results")

	return cmd
}

// printBacktestResults prints the results of a backtest
func printBacktestResults(result *backtest.BacktestResult) {
	metrics := result.PerformanceMetrics
	if metrics == nil {
		fmt.Println("No performance metrics available")
		return
	}

	fmt.Println("=== Backtest Results ===")
	fmt.Printf("Initial Capital: $%.2f\n", result.InitialCapital)
	fmt.Printf("Final Capital: $%.2f\n", result.FinalCapital)
	fmt.Printf("Total Return: %.2f%%\n", metrics.TotalReturn)
	fmt.Printf("Annualized Return: %.2f%%\n", metrics.AnnualizedReturn)
	fmt.Printf("Sharpe Ratio: %.2f\n", metrics.SharpeRatio)
	fmt.Printf("Sortino Ratio: %.2f\n", metrics.SortinoRatio)
	fmt.Printf("Max Drawdown: $%.2f (%.2f%%)\n", metrics.MaxDrawdown, metrics.MaxDrawdownPercent)
	fmt.Printf("Win Rate: %.2f%%\n", metrics.WinRate)
	fmt.Printf("Profit Factor: %.2f\n", metrics.ProfitFactor)
	fmt.Printf("Expected Payoff: $%.2f\n", metrics.ExpectedPayoff)
	fmt.Printf("Total Trades: %d\n", metrics.TotalTrades)
	fmt.Printf("Winning Trades: %d\n", metrics.WinningTrades)
	fmt.Printf("Losing Trades: %d\n", metrics.LosingTrades)
	fmt.Printf("Break-Even Trades: %d\n", metrics.BreakEvenTrades)
	fmt.Printf("Average Profit Trade: $%.2f\n", metrics.AverageProfitTrade)
	fmt.Printf("Average Loss Trade: $%.2f\n", metrics.AverageLossTrade)
	fmt.Printf("Largest Profit Trade: $%.2f\n", metrics.LargestProfitTrade)
	fmt.Printf("Largest Loss Trade: $%.2f\n", metrics.LargestLossTrade)
	fmt.Printf("Average Holding Time: %s\n", metrics.AverageHoldingTime)
}

// saveBacktestResults saves the results of a backtest to a file
func saveBacktestResults(result *backtest.BacktestResult, outputFile string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(outputFile)
	if dir != "" {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}

	// Create file
	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Write header
	fmt.Fprintf(file, "Backtest Results\n")
	fmt.Fprintf(file, "===============\n\n")

	// Write configuration
	fmt.Fprintf(file, "Configuration:\n")
	fmt.Fprintf(file, "  Strategy: %s\n", result.Config.Strategy.(*backtest.BaseStrategy).Name)
	fmt.Fprintf(file, "  Symbols: %v\n", result.Config.Symbols)
	fmt.Fprintf(file, "  Start Date: %s\n", result.StartTime.Format("2006-01-02"))
	fmt.Fprintf(file, "  End Date: %s\n", result.EndTime.Format("2006-01-02"))
	fmt.Fprintf(file, "  Initial Capital: $%.2f\n", result.InitialCapital)
	fmt.Fprintf(file, "  Interval: %s\n", result.Config.Interval)
	fmt.Fprintf(file, "  Commission Rate: %.2f%%\n", result.Config.CommissionRate*100)
	fmt.Fprintf(file, "\n")

	// Write performance metrics
	metrics := result.PerformanceMetrics
	if metrics != nil {
		fmt.Fprintf(file, "Performance Metrics:\n")
		fmt.Fprintf(file, "  Final Capital: $%.2f\n", result.FinalCapital)
		fmt.Fprintf(file, "  Total Return: %.2f%%\n", metrics.TotalReturn)
		fmt.Fprintf(file, "  Annualized Return: %.2f%%\n", metrics.AnnualizedReturn)
		fmt.Fprintf(file, "  Sharpe Ratio: %.2f\n", metrics.SharpeRatio)
		fmt.Fprintf(file, "  Sortino Ratio: %.2f\n", metrics.SortinoRatio)
		fmt.Fprintf(file, "  Max Drawdown: $%.2f (%.2f%%)\n", metrics.MaxDrawdown, metrics.MaxDrawdownPercent)
		fmt.Fprintf(file, "  Win Rate: %.2f%%\n", metrics.WinRate)
		fmt.Fprintf(file, "  Profit Factor: %.2f\n", metrics.ProfitFactor)
		fmt.Fprintf(file, "  Expected Payoff: $%.2f\n", metrics.ExpectedPayoff)
		fmt.Fprintf(file, "  Total Trades: %d\n", metrics.TotalTrades)
		fmt.Fprintf(file, "  Winning Trades: %d\n", metrics.WinningTrades)
		fmt.Fprintf(file, "  Losing Trades: %d\n", metrics.LosingTrades)
		fmt.Fprintf(file, "  Break-Even Trades: %d\n", metrics.BreakEvenTrades)
		fmt.Fprintf(file, "  Average Profit Trade: $%.2f\n", metrics.AverageProfitTrade)
		fmt.Fprintf(file, "  Average Loss Trade: $%.2f\n", metrics.AverageLossTrade)
		fmt.Fprintf(file, "  Largest Profit Trade: $%.2f\n", metrics.LargestProfitTrade)
		fmt.Fprintf(file, "  Largest Loss Trade: $%.2f\n", metrics.LargestLossTrade)
		fmt.Fprintf(file, "  Average Holding Time: %s\n", metrics.AverageHoldingTime)
		fmt.Fprintf(file, "\n")
	}

	// Write trade summary
	fmt.Fprintf(file, "Trade Summary:\n")
	fmt.Fprintf(file, "  Symbol | Side | Entry Price | Exit Price | Quantity | Profit | Open Time | Close Time\n")
	fmt.Fprintf(file, "  ------ | ---- | ----------- | ---------- | -------- | ------ | --------- | ----------\n")
	for _, position := range result.ClosedPositions {
		fmt.Fprintf(file, "  %s | %s | $%.2f | $%.2f | %.4f | $%.2f | %s | %s\n",
			position.Symbol,
			string(position.Side),
			position.EntryPrice,
			position.ExitPrice,
			position.Quantity,
			position.Profit,
			position.OpenTime.Format("2006-01-02 15:04:05"),
			position.CloseTime.Format("2006-01-02 15:04:05"),
		)
	}

	return nil
}
