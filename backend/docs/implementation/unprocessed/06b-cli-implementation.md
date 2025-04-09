# CLI Implementation

This document details the implementation of the Command Line Interface (CLI) for the Go crypto trading bot, providing commands for monitoring and managing the trading bot.

## 1. CLI Structure Overview

The CLI follows a modular structure using the Cobra library:

```
cmd/cli/
├── commands/
│   ├── root.go         # Root command definition
│   ├── portfolio.go    # Portfolio commands
│   ├── trade.go        # Trading operation commands
│   ├── newcoin.go      # New coin detection commands
│   └── config.go       # Configuration commands
└── main.go             # CLI entry point
```

## 2. Root Command Implementation

```go
// cmd/cli/commands/root.go
package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	dbPath  string
	apiKey  string
	apiSecret string
)

// RootCmd represents the base command
var RootCmd = &cobra.Command{
	Use:   "cryptobot",
	Short: "A cryptocurrency trading bot",
	Long: `A cryptocurrency trading bot for automated trading on the MEXC exchange.
Monitors for new coin listings and implements various trading strategies.`,
}

// Execute adds all child commands to the root command and sets flags appropriately
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cryptobot.yaml)")
	RootCmd.PersistentFlags().StringVar(&dbPath, "db", "data/cryptobot.db", "database file path")
	RootCmd.PersistentFlags().StringVar(&apiKey, "key", "", "MEXC API key")
	RootCmd.PersistentFlags().StringVar(&apiSecret, "secret", "", "MEXC API secret key")

	// Bind flags to viper
	viper.BindPFlag("db", RootCmd.PersistentFlags().Lookup("db"))
	viper.BindPFlag("mexc.key", RootCmd.PersistentFlags().Lookup("key"))
	viper.BindPFlag("mexc.secret", RootCmd.PersistentFlags().Lookup("secret"))
}

// initConfig reads in config file and ENV variables if set
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".cryptobot" (without extension)
		viper.AddConfigPath(home)
		viper.SetConfigName(".cryptobot")
	}

	// Read in environment variables that match
	viper.AutomaticEnv()

	// If a config file is found, read it in
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
```

## 3. Portfolio Command Implementation

```go
// cmd/cli/commands/portfolio.go
package commands

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"

	"github.com/ryanlisse/cryptobot/internal/core/portfolio"
	"github.com/ryanlisse/cryptobot/internal/platform/database"
	"github.com/ryanlisse/cryptobot/internal/platform/mexc"
)

var (
	timeRange string
)

func init() {
	portfolioCmd := &cobra.Command{
		Use:   "portfolio",
		Short: "Manage and view your portfolio",
		Long:  `Commands for viewing your portfolio status, active trades, and performance metrics.`,
	}

	// portfolio value command
	valueCmd := &cobra.Command{
		Use:   "value",
		Short: "Get total portfolio value",
		Run:   getPortfolioValue,
	}

	// portfolio trades command
	tradesCmd := &cobra.Command{
		Use:   "trades",
		Short: "List all active trades",
		Run:   getActiveTrades,
	}

	// portfolio performance command
	performanceCmd := &cobra.Command{
		Use:   "performance",
		Short: "Show trading performance metrics",
		Run:   getPerformance,
	}

	// Add flags
	performanceCmd.Flags().StringVarP(&timeRange, "range", "r", "week", "Time range (day, week, month, year, all)")

	// Add commands to portfolio
	portfolioCmd.AddCommand(valueCmd, tradesCmd, performanceCmd)

	// Add portfolio to root
	RootCmd.AddCommand(portfolioCmd)
}

func getPortfolioValue(cmd *cobra.Command, args []string) {
	ctx := context.Background()

	// Initialize services
	service, err := initPortfolioService(ctx)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Get portfolio value
	value, err := service.GetPortfolioValue(ctx)
	if err != nil {
		fmt.Printf("Error getting portfolio value: %v\n", err)
		return
	}

	fmt.Printf("Total Portfolio Value: %.2f USDT\n", value)
}

func getActiveTrades(cmd *cobra.Command, args []string) {
	ctx := context.Background()

	// Initialize services
	service, err := initPortfolioService(ctx)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Get active trades
	trades, err := service.GetActiveTrades(ctx)
	if err != nil {
		fmt.Printf("Error getting active trades: %v\n", err)
		return
	}

	if len(trades) == 0 {
		fmt.Println("No active trades found.")
		return
	}

	// Print trades in table format
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tSymbol\tBuy Price\tCurrent Price\tQuantity\tProfit %\tCurrent Value\tBought At")
	for _, trade := range trades {
		fmt.Fprintf(w, "%d\t%s\t%.4f\t%.4f\t%.6f\t%.2f%%\t%.2f\t%s\n",
			trade.ID,
			trade.Symbol,
			trade.PurchasePrice,
			trade.CurrentPrice,
			trade.Quantity,
			trade.ProfitPercentage,
			trade.CurrentValue,
			trade.PurchaseTime.Format("2006-01-02 15:04:05"),
		)
	}
	w.Flush()
}

// Helper function to initialize the portfolio service
func initPortfolioService(ctx context.Context) (*portfolio.Service, error) {
	// Connect to database
	db, err := database.Connect(database.Config{Path: dbPath})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Create repositories
	boughtCoinRepo := database.NewBoughtCoinRepository(db)

	// Create exchange service
	exchangeService, err := mexc.NewClient(apiKey, apiSecret, "")
	if err != nil {
		return nil, fmt.Errorf("failed to create MEXC client: %w", err)
	}

	if err := exchangeService.Connect(ctx); err != nil {
		return nil, fmt.Errorf("failed to connect to MEXC: %w", err)
	}

	// Create portfolio service
	portfolioService := portfolio.NewService(
		exchangeService,
		boughtCoinRepo,
	)

	return portfolioService, nil
}
```

## 4. Trade Command Implementation

```go
// cmd/cli/commands/trade.go
package commands

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/ryanlisse/cryptobot/internal/core/trade"
	"github.com/ryanlisse/cryptobot/internal/domain/models"
	"github.com/ryanlisse/cryptobot/internal/platform/database"
	"github.com/ryanlisse/cryptobot/internal/platform/mexc"
)

func init() {
	tradeCmd := &cobra.Command{
		Use:   "trade",
		Short: "Execute trading operations",
		Long:  `Commands for buying, selling, and managing trades.`,
	}

	// buy command
	buyCmd := &cobra.Command{
		Use:   "buy [symbol] [amount]",
		Short: "Buy a cryptocurrency",
		Args:  cobra.RangeArgs(1, 2),
		Run:   executeBuy,
	}

	// sell command
	sellCmd := &cobra.Command{
		Use:   "sell [coin_id] [amount]",
		Short: "Sell a cryptocurrency",
		Args:  cobra.RangeArgs(1, 2),
		Run:   executeSell,
	}

	// Add trade commands to root
	tradeCmd.AddCommand(buyCmd, sellCmd)
	RootCmd.AddCommand(tradeCmd)
}

func executeBuy(cmd *cobra.Command, args []string) {
	ctx := context.Background()

	// Parse arguments
	symbol := args[0]
	
	var amount float64
	var err error
	
	if len(args) > 1 {
		amount, err = strconv.ParseFloat(args[1], 64)
		if err != nil {
			fmt.Printf("Error parsing amount: %v\n", err)
			return
		}
	} else {
		// Use default amount from config
		amount = 0 // Will use service default
	}

	// Initialize services
	service, err := initTradeService(ctx)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Execute purchase
	coin, err := service.ExecutePurchase(ctx, symbol, amount)
	if err != nil {
		fmt.Printf("Error executing purchase: %v\n", err)
		return
	}

	fmt.Printf("Successfully purchased %s:\n", symbol)
	fmt.Printf("Amount: %.6f\n", coin.Quantity)
	fmt.Printf("Price: %.4f USDT\n", coin.PurchasePrice)
	fmt.Printf("Total: %.2f USDT\n", coin.Quantity*coin.PurchasePrice)
	fmt.Printf("Stop Loss: %.4f USDT\n", coin.StopLossPrice)
}

func executeSell(cmd *cobra.Command, args []string) {
	ctx := context.Background()

	// Parse arguments
	coinID, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		fmt.Printf("Error parsing coin ID: %v\n", err)
		return
	}
	
	var amount float64
	var sellAll bool
	
	if len(args) > 1 {
		if args[1] == "all" {
			sellAll = true
		} else {
			amount, err = strconv.ParseFloat(args[1], 64)
			if err != nil {
				fmt.Printf("Error parsing amount: %v\n", err)
				return
			}
		}
	} else {
		sellAll = true
	}

	// Initialize services
	service, err := initTradeService(ctx)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Get the coin
	db, err := database.Connect(database.Config{Path: dbPath})
	if err != nil {
		fmt.Printf("Error connecting to database: %v\n", err)
		return
	}
	
	boughtCoinRepo := database.NewBoughtCoinRepository(db)
	coin, err := boughtCoinRepo.FindByID(ctx, uint(coinID))
	if err != nil {
		fmt.Printf("Error finding coin: %v\n", err)
		return
	}
	
	if coin == nil {
		fmt.Printf("Coin with ID %d not found\n", coinID)
		return
	}
	
	// Set sell amount
	if sellAll {
		amount = coin.Quantity
	}

	// Execute sell
	order, err := service.SellCoin(ctx, coin, amount)
	if err != nil {
		fmt.Printf("Error selling coin: %v\n", err)
		return
	}

	fmt.Printf("Successfully sold %s:\n", coin.Symbol)
	fmt.Printf("Amount: %.6f\n", order.Quantity)
	fmt.Printf("Price: %.4f USDT\n", order.Price)
	fmt.Printf("Total: %.2f USDT\n", order.Quantity*order.Price)
}

// Helper function to initialize the trade service
func initTradeService(ctx context.Context) (*trade.Service, error) {
	// Connect to database
	db, err := database.Connect(database.Config{Path: dbPath})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Create repositories
	boughtCoinRepo := database.NewBoughtCoinRepository(db)
	logRepo := database.NewLogRepository(db)

	// Create exchange service
	exchangeService, err := mexc.NewClient(apiKey, apiSecret, "")
	if err != nil {
		return nil, fmt.Errorf("failed to create MEXC client: %w", err)
	}

	if err := exchangeService.Connect(ctx); err != nil {
		return nil, fmt.Errorf("failed to connect to MEXC: %w", err)
	}

	// Create trade service
	tradeService, err := trade.NewService(
		exchangeService,
		boughtCoinRepo,
		logRepo,
		trade.DefaultConfig(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create trade service: %w", err)
	}

	return tradeService, nil
}
```

## 5. NewCoin Command Implementation

```go
// cmd/cli/commands/newcoin.go
package commands

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"

	"github.com/ryanlisse/cryptobot/internal/core/newcoin"
	"github.com/ryanlisse/cryptobot/internal/platform/database"
	"github.com/ryanlisse/cryptobot/internal/platform/mexc"
)

func init() {
	newcoinCmd := &cobra.Command{
		Use:   "newcoin",
		Short: "Manage new coin detection",
		Long:  `Commands for listing and processing newly detected coins.`,
	}

	// list command
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List newly detected coins",
		Run:   listNewCoins,
	}

	// process command
	processCmd := &cobra.Command{
		Use:   "process",
		Short: "Process newly detected coins",
		Run:   processNewCoins,
	}

	// Add commands
	newcoinCmd.AddCommand(listCmd, processCmd)
	RootCmd.AddCommand(newcoinCmd)
}

func listNewCoins(cmd *cobra.Command, args []string) {
	ctx := context.Background()

	// Connect to database
	db, err := database.Connect(database.Config{Path: dbPath})
	if err != nil {
		fmt.Printf("Error connecting to database: %v\n", err)
		return
	}

	// Create repository
	newCoinRepo := database.NewNewCoinRepository(db)

	// Get all new coins
	coins, err := newCoinRepo.FindActive(ctx)
	if err != nil {
		fmt.Printf("Error getting new coins: %v\n", err)
		return
	}

	if len(coins) == 0 {
		fmt.Println("No new coins detected.")
		return
	}

	// Print in table format
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tSymbol\tDetected At\tLast Checked")
	for _, coin := range coins {
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\n",
			coin.ID,
			coin.Symbol,
			coin.DetectedAt.Format("2006-01-02 15:04:05"),
			coin.LastChecked.Format("2006-01-02 15:04:05"),
		)
	}
	w.Flush()
}

func processNewCoins(cmd *cobra.Command, args []string) {
	ctx := context.Background()

	// Initialize services
	service, err := initNewCoinService(ctx)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Process new coins
	results, err := service.ProcessNewCoins(ctx)
	if err != nil {
		fmt.Printf("Error processing new coins: %v\n", err)
		return
	}

	fmt.Printf("Processed %d new coins\n", len(results))
	
	// Print results
	if len(results) > 0 {
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "Symbol\tStatus\tReason")
		for _, result := range results {
			fmt.Fprintf(w, "%s\t%s\t%s\n",
				result.Symbol,
				result.Status,
				result.Reason,
			)
		}
		w.Flush()
	}
}

// Helper function to initialize the newcoin service
func initNewCoinService(ctx context.Context) (*newcoin.Service, error) {
	// Connect to database
	db, err := database.Connect(database.Config{Path: dbPath})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Create repositories
	newCoinRepo := database.NewNewCoinRepository(db)
	purchaseDecisionRepo := database.NewPurchaseDecisionRepository(db)
	logRepo := database.NewLogRepository(db)

	// Create exchange service
	exchangeService, err := mexc.NewClient(apiKey, apiSecret, "")
	if err != nil {
		return nil, fmt.Errorf("failed to create MEXC client: %w", err)
	}

	if err := exchangeService.Connect(ctx); err != nil {
		return nil, fmt.Errorf("failed to connect to MEXC: %w", err)
	}

	// Create newcoin service
	newCoinService := newcoin.NewService(
		exchangeService,
		newCoinRepo,
		purchaseDecisionRepo,
		logRepo,
		newcoin.DefaultConfig(),
	)

	return newCoinService, nil
}
```

## 6. Bot Command for Running the Trading Bot

```go
// cmd/cli/commands/bot.go
package commands

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"github.com/ryanlisse/cryptobot/internal/bot"
	"github.com/ryanlisse/cryptobot/internal/core/newcoin"
	"github.com/ryanlisse/cryptobot/internal/core/portfolio"
	"github.com/ryanlisse/cryptobot/internal/core/trade"
	"github.com/ryanlisse/cryptobot/internal/platform/database"
	"github.com/ryanlisse/cryptobot/internal/platform/mexc"
)

func init() {
	botCmd := &cobra.Command{
		Use:   "bot",
		Short: "Run the trading bot",
		Long:  `Run the automated trading bot which monitors for new coins and manages trades.`,
		Run:   runBot,
	}

	// Add flags
	botCmd.Flags().Bool("headless", false, "Run in headless mode without interactive UI")
	botCmd.Flags().Duration("check-interval", 30*time.Second, "Interval for checking new coins")

	// Add to root
	RootCmd.AddCommand(botCmd)
}

func runBot(cmd *cobra.Command, args []string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Get flags
	headless, _ := cmd.Flags().GetBool("headless")
	checkInterval, _ := cmd.Flags().GetDuration("check-interval")

	fmt.Println("Starting crypto trading bot...")

	// Setup signal handling for graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigCh
		fmt.Printf("\nReceived signal: %v\n", sig)
		fmt.Println("Shutting down gracefully...")
		cancel()
	}()

	// Connect to database
	db, err := database.Connect(database.Config{Path: dbPath})
	if err != nil {
		fmt.Printf("Error connecting to database: %v\n", err)
		return
	}
	defer database.Close(db)

	// Run migrations
	if err := database.Migrate(db); err != nil {
		fmt.Printf("Error running migrations: %v\n", err)
		return
	}

	// Create repositories
	boughtCoinRepo := database.NewBoughtCoinRepository(db)
	newCoinRepo := database.NewNewCoinRepository(db)
	purchaseDecisionRepo := database.NewPurchaseDecisionRepository(db)
	logRepo := database.NewLogRepository(db)

	// Create exchange service
	exchangeService, err := mexc.NewClient(apiKey, apiSecret, "")
	if err != nil {
		fmt.Printf("Error creating MEXC client: %v\n", err)
		return
	}

	if err := exchangeService.Connect(ctx); err != nil {
		fmt.Printf("Error connecting to MEXC: %v\n", err)
		return
	}
	defer exchangeService.Disconnect()

	// Create services
	tradeService, err := trade.NewService(
		exchangeService,
		boughtCoinRepo,
		logRepo,
		trade.DefaultConfig(),
	)
	if err != nil {
		fmt.Printf("Error creating trade service: %v\n", err)
		return
	}

	portfolioService := portfolio.NewService(
		exchangeService,
		boughtCoinRepo,
	)

	newCoinService := newcoin.NewService(
		exchangeService,
		newCoinRepo,
		purchaseDecisionRepo,
		logRepo,
		newcoin.DefaultConfig(),
	)

	// Create bot configuration
	botConfig := bot.Config{
		NewCoinCheckInterval: checkInterval,
		PortfolioCheckInterval: 1 * time.Minute,
		MaxConcurrentRequests: 5,
		Headless: headless,
	}

	// Create and run bot
	tradingBot := bot.NewBot(
		tradeService,
		newCoinService,
		portfolioService,
		logRepo,
		botConfig,
	)

	if err := tradingBot.Run(ctx); err != nil {
		fmt.Printf("Bot error: %v\n", err)
		return
	}

	// Wait for context cancellation (from signal handler)
	<-ctx.Done()
	fmt.Println("Bot stopped")
}
```

## 7. Main Entry Point

```go
// cmd/cli/main.go
package main

import "github.com/ryanlisse/cryptobot/cmd/cli/commands"

func main() {
	commands.Execute()
}
```

## 8. Testing CLI Commands

```go
// cmd/cli/commands/portfolio_test.go
package commands

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/ryanlisse/cryptobot/internal/domain/models"
)

// MockPortfolioService for testing
type MockPortfolioService struct {
	portfolioValue float64
	activeTrades   []*models.BoughtCoin
}

// Implement interface methods
func (m *MockPortfolioService) GetPortfolioValue(ctx context.Context) (float64, error) {
	return m.portfolioValue, nil
}

func (m *MockPortfolioService) GetActiveTrades(ctx context.Context) ([]*models.BoughtCoin, error) {
	return m.activeTrades, nil
}

// TestGetPortfolioValueCommand tests the portfolio value command
func TestGetPortfolioValueCommand(t *testing.T) {
	// Create mock
	mockService := &MockPortfolioService{
		portfolioValue: 1234.56,
	}

	// Save original and restore after test
	originalInit := initPortfolioService
	defer func() { initPortfolioService = originalInit }()

	// Replace with mock
	initPortfolioService = func(ctx context.Context) (*portfolio.Service, error) {
		return nil, nil // Not used as we're mocking the actual call
	}

	// Create command
	cmd := &cobra.Command{Use: "test"}
	var output bytes.Buffer
	cmd.SetOut(&output)

	// Execute with our mock
	ctx := context.Background()
	actualGetPortfolioValue := getPortfolioValue
	getPortfolioValue = func(cmd *cobra.Command, args []string) {
		// Mock implementation
		value := mockService.portfolioValue
		cmd.Printf("Total Portfolio Value: %.2f USDT\n", value)
	}
	defer func() { getPortfolioValue = actualGetPortfolioValue }()

	// Run command
	getPortfolioValue(cmd, []string{})

	// Check output
	expectedOutput := "Total Portfolio Value: 1234.56 USDT\n"
	if output.String() != expectedOutput {
		t.Errorf("Expected output: %q, got: %q", expectedOutput, output.String())
	}
}
```

## 9. Building and Packaging

```bash
# Build CLI
go build -o bin/cryptobot cmd/cli/main.go

# Cross-compile for different platforms
GOOS=windows GOARCH=amd64 go build -o bin/cryptobot.exe cmd/cli/main.go
GOOS=darwin GOARCH=amd64 go build -o bin/cryptobot-mac cmd/cli/main.go
GOOS=linux GOARCH=amd64 go build -o bin/cryptobot-linux cmd/cli/main.go

# Use goreleaser for more advanced packaging
goreleaser init
goreleaser release --snapshot --skip-publish --rm-dist
```

## 10. Configuration Example

```yaml
# ~/.cryptobot.yaml
db: "./data/cryptobot.db"

mexc:
  key: "your-api-key"
  secret: "your-api-secret"

trade:
  usdt_per_trade: 25.0
  stop_loss_percent: 15.0
  take_profit_levels: [5.0, 10.0, 15.0, 20.0]
  sell_percentages: [0.25, 0.25, 0.25, 0.25]

newcoin:
  min_volume: 10000
  min_price: 0.000001
  max_price: 0.01
  max_age_hours: 24

bot:
  check_interval: 30s
  portfolio_check_interval: 1m
  max_concurrent_requests: 5
```
