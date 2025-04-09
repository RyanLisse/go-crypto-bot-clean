package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/ryanlisse/go-crypto-bot/internal/backtest"
	"github.com/ryanlisse/go-crypto-bot/internal/backtest/strategies"
	"github.com/ryanlisse/go-crypto-bot/internal/domain/models"
	"go.uber.org/zap"
)

func main() {
	// Create logger
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	// Create data provider with test data
	dataProvider := createTestDataProvider()

	// Create strategy
	strategy := strategies.NewSimpleMAStrategy(10, 50, logger)

	// Create slippage model
	slippageModel := backtest.NewFixedSlippage(0.1) // 0.1% slippage

	// Create backtest config
	config := &backtest.BacktestConfig{
		StartTime:          time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		EndTime:            time.Date(2023, 1, 31, 0, 0, 0, 0, time.UTC),
		InitialCapital:     10000,
		Symbols:            []string{"BTCUSDT"},
		Interval:           "1h",
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
	logger.Info("Running backtest...")
	result, err := engine.Run(context.Background())
	if err != nil {
		logger.Fatal("Backtest failed", zap.Error(err))
	}

	// Print results
	fmt.Println("=== Backtest Results ===")
	fmt.Printf("Initial Capital: $%.2f\n", result.InitialCapital)
	fmt.Printf("Final Capital: $%.2f\n", result.FinalCapital)
	fmt.Printf("Total Return: %.2f%%\n", result.PerformanceMetrics.TotalReturn)
	fmt.Printf("Sharpe Ratio: %.2f\n", result.PerformanceMetrics.SharpeRatio)
	fmt.Printf("Max Drawdown: %.2f%%\n", result.PerformanceMetrics.MaxDrawdownPercent)
	fmt.Printf("Win Rate: %.2f%%\n", result.PerformanceMetrics.WinRate)
	fmt.Printf("Total Trades: %d\n", result.PerformanceMetrics.TotalTrades)
	fmt.Printf("Winning Trades: %d\n", result.PerformanceMetrics.WinningTrades)
	fmt.Printf("Losing Trades: %d\n", result.PerformanceMetrics.LosingTrades)
	fmt.Printf("Average Profit Trade: $%.2f\n", result.PerformanceMetrics.AverageProfitTrade)
	fmt.Printf("Average Loss Trade: $%.2f\n", result.PerformanceMetrics.AverageLossTrade)
	fmt.Printf("Average Holding Time: %s\n", result.PerformanceMetrics.AverageHoldingTime)

	// Print trades
	fmt.Println("\n=== Trades ===")
	fmt.Println("Symbol | Side | Price | Quantity | Time")
	fmt.Println("------ | ---- | ----- | -------- | ----")
	for _, trade := range result.Trades {
		fmt.Printf("%s | %s | $%.2f | %.4f | %s\n",
			trade.Symbol,
			trade.Side,
			trade.Price,
			trade.Quantity,
			trade.Time.Format("2006-01-02 15:04:05"),
		)
	}
}

// createTestDataProvider creates a test data provider with simulated price data
func createTestDataProvider() *backtest.InMemoryDataProvider {
	dataProvider := backtest.NewInMemoryDataProvider()

	// Create test klines for BTCUSDT
	symbol := "BTCUSDT"
	interval := "1h"
	startTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	endTime := time.Date(2023, 1, 31, 0, 0, 0, 0, time.UTC)

	// Create klines with a simple price pattern
	klines := createTestKlines(symbol, startTime, endTime, interval)
	dataProvider.AddKlines(symbol, interval, klines)

	return dataProvider
}

// createTestKlines creates test klines for backtesting
func createTestKlines(symbol string, startTime, endTime time.Time, interval string) []*models.Kline {
	var klines []*models.Kline
	var intervalDuration time.Duration

	switch interval {
	case "1m":
		intervalDuration = time.Minute
	case "5m":
		intervalDuration = 5 * time.Minute
	case "15m":
		intervalDuration = 15 * time.Minute
	case "1h":
		intervalDuration = time.Hour
	case "4h":
		intervalDuration = 4 * time.Hour
	case "1d":
		intervalDuration = 24 * time.Hour
	default:
		intervalDuration = time.Hour
	}

	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Create klines with a simple price pattern
	price := 20000.0
	for t := startTime; t.Before(endTime); t = t.Add(intervalDuration) {
		// Simple price movement: oscillate between 19000 and 21000
		price = price + (200 * (0.5 - rand.Float64()))
		if price < 19000 {
			price = 19000
		}
		if price > 21000 {
			price = 21000
		}

		kline := &models.Kline{
			Symbol:    symbol,
			Interval:  interval,
			OpenTime:  t,
			CloseTime: t.Add(intervalDuration),
			Open:      price - 50,
			High:      price + 100,
			Low:       price - 100,
			Close:     price,
			Volume:    1000 + rand.Float64()*1000,
			IsClosed:  true,
		}

		klines = append(klines, kline)
	}

	return klines
}
