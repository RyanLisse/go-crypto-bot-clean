package scripts

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/pkg/platform/mexc"
)

// TestMexcEndpoints tests all MEXC endpoints
func TestMexcEndpoints() {
	// Setup logger
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	logger := log.With().Str("component", "mexc-endpoints-test").Logger()

	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		// Try parent directory
		err = godotenv.Load("../.env")
		if err != nil {
			logger.Fatal().Err(err).Msg("Error loading .env file")
		}
	}

	// Get API credentials from environment variables
	apiKey := os.Getenv("MEXC_API_KEY")
	apiSecret := os.Getenv("MEXC_SECRET_KEY")

	if apiKey == "" || apiSecret == "" {
		logger.Fatal().Msg("MEXC_API_KEY and MEXC_SECRET_KEY environment variables must be set")
	}

	// Log the API key and secret (truncated for security)
	logger.Info().
		Str("API Key (truncated)", apiKey[:5]+"..."+apiKey[len(apiKey)-4:]).
		Str("API Secret (truncated)", apiSecret[:5]+"..."+apiSecret[len(apiSecret)-4:]).
		Msg("Using MEXC credentials")

	// Create MEXC client
	client := mexc.NewClient(apiKey, apiSecret, &logger)
	logger.Info().Msg("MEXC client created")

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test all endpoints
	testAccount(ctx, client, &logger)
	testExchangeInfo(ctx, client, &logger)
	testTicker(ctx, client, &logger)
	testOrderBook(ctx, client, &logger)
	testKlines(ctx, client, &logger)
	testMarketData(ctx, client, &logger)

	logger.Info().Msg("All tests completed successfully")
}

// testAccount tests the account endpoint
func testAccount(ctx context.Context, client *mexc.Client, logger *zerolog.Logger) {
	logger.Info().Msg("Testing account endpoint...")

	wallet, err := client.GetAccount(ctx)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to get account information")
	}

	logger.Info().
		Int("Number of balances", len(wallet.Balances)).
		Msg("Account information retrieved successfully")

	// Print balances
	fmt.Println("\n=== MEXC Account Balances ===")
	fmt.Printf("%-10s %-15s %-15s\n", "Asset", "Free", "Locked")
	fmt.Println("----------------------------------------")

	// Filter out zero balances
	nonZeroBalances := 0
	for asset, balance := range wallet.Balances {
		if balance.Free > 0 || balance.Locked > 0 {
			fmt.Printf("%-10s %-15f %-15f\n", asset, balance.Free, balance.Locked)
			nonZeroBalances++
		}
	}

	if nonZeroBalances == 0 {
		fmt.Println("No non-zero balances found.")
	}

	// Save wallet data to file
	saveToFile(wallet, "mexc_wallet.json", logger)
}

// testExchangeInfo tests the exchange info endpoint
func testExchangeInfo(ctx context.Context, client *mexc.Client, logger *zerolog.Logger) {
	logger.Info().Msg("Testing exchange info endpoint...")

	exchangeInfo, err := client.GetExchangeInfo(ctx)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to get exchange information")
	}

	logger.Info().
		Int("Number of symbols", len(exchangeInfo.Symbols)).
		Msg("Exchange information retrieved successfully")

	// Print a few symbols
	fmt.Println("\n=== MEXC Exchange Symbols ===")
	fmt.Printf("%-10s %-10s %-10s %-10s\n", "Symbol", "Base", "Quote", "Status")
	fmt.Println("----------------------------------------")

	// Print first 10 symbols
	for i := 0; i < 10 && i < len(exchangeInfo.Symbols); i++ {
		symbol := exchangeInfo.Symbols[i]
		fmt.Printf("%-10s %-10s %-10s %-10s\n", symbol.Symbol, symbol.BaseAsset, symbol.QuoteAsset, symbol.Status)
	}

	// Save exchange info to file
	saveToFile(exchangeInfo, "mexc_exchange_info.json", logger)
}

// testTicker tests the ticker endpoint
func testTicker(ctx context.Context, client *mexc.Client, logger *zerolog.Logger) {
	logger.Info().Msg("Testing ticker endpoint...")

	// Test with a few popular symbols
	symbols := []string{"BTCUSDT", "ETHUSDT", "SOLUSDT"}

	for _, symbol := range symbols {
		ticker, err := client.GetMarketData(ctx, symbol)
		if err != nil {
			logger.Error().Err(err).Str("symbol", symbol).Msg("Failed to get ticker")
			continue
		}

		logger.Info().
			Str("Symbol", ticker.Symbol).
			Float64("Price", ticker.LastPrice).
			Float64("24h Volume", ticker.Volume).
			Float64("24h High", ticker.HighPrice).
			Float64("24h Low", ticker.LowPrice).
			Float64("24h Change %", ticker.PriceChangePercent).
			Msg("Ticker retrieved successfully")

		// Save ticker to file
		saveToFile(ticker, fmt.Sprintf("mexc_ticker_%s.json", symbol), logger)
	}
}

// testOrderBook tests the order book endpoint
func testOrderBook(ctx context.Context, client *mexc.Client, logger *zerolog.Logger) {
	logger.Info().Msg("Testing order book endpoint...")

	// Test with a few popular symbols
	symbols := []string{"BTCUSDT", "ETHUSDT", "SOLUSDT"}

	for _, symbol := range symbols {
		orderBook, err := client.GetOrderBook(ctx, symbol, 10)
		if err != nil {
			logger.Error().Err(err).Str("symbol", symbol).Msg("Failed to get order book")
			continue
		}

		logger.Info().
			Str("Symbol", orderBook.Symbol).
			Int64("Last Update ID", orderBook.LastUpdateID).
			Int("Bids", len(orderBook.Bids)).
			Int("Asks", len(orderBook.Asks)).
			Msg("Order book retrieved successfully")

		// Print top 3 bids and asks
		fmt.Printf("\n=== MEXC Order Book for %s ===\n", symbol)
		fmt.Println("Top 3 Bids:")
		fmt.Printf("%-15s %-15s\n", "Price", "Quantity")
		fmt.Println("------------------------------")
		for i := 0; i < 3 && i < len(orderBook.Bids); i++ {
			fmt.Printf("%-15f %-15f\n", orderBook.Bids[i].Price, orderBook.Bids[i].Quantity)
		}

		fmt.Println("\nTop 3 Asks:")
		fmt.Printf("%-15s %-15s\n", "Price", "Quantity")
		fmt.Println("------------------------------")
		for i := 0; i < 3 && i < len(orderBook.Asks); i++ {
			fmt.Printf("%-15f %-15f\n", orderBook.Asks[i].Price, orderBook.Asks[i].Quantity)
		}

		// Save order book to file
		saveToFile(orderBook, fmt.Sprintf("mexc_orderbook_%s.json", symbol), logger)
	}
}

// testKlines tests the klines endpoint
func testKlines(ctx context.Context, client *mexc.Client, logger *zerolog.Logger) {
	logger.Info().Msg("Testing klines endpoint...")

	// Test with a few popular symbols and intervals
	symbols := []string{"BTCUSDT", "ETHUSDT", "SOLUSDT"}
	// MEXC uses different interval format, so we'll use string literals
	intervals := []model.KlineInterval{"60m", "4h", "1d"}

	for _, symbol := range symbols {
		for _, interval := range intervals {
			klines, err := client.GetKlines(ctx, symbol, interval, 10)
			if err != nil {
				logger.Error().Err(err).Str("symbol", symbol).Str("interval", string(interval)).Msg("Failed to get klines")
				continue
			}

			logger.Info().
				Str("Symbol", symbol).
				Str("Interval", string(interval)).
				Int("Count", len(klines)).
				Msg("Klines retrieved successfully")

			// Print the most recent kline
			if len(klines) > 0 {
				kline := klines[0]
				fmt.Printf("\n=== Most Recent Kline for %s (%s) ===\n", symbol, interval)
				fmt.Printf("Open Time: %s\n", kline.OpenTime.Format(time.RFC3339))
				fmt.Printf("Open: %.8f\n", kline.Open)
				fmt.Printf("High: %.8f\n", kline.High)
				fmt.Printf("Low: %.8f\n", kline.Low)
				fmt.Printf("Close: %.8f\n", kline.Close)
				fmt.Printf("Volume: %.8f\n", kline.Volume)
				fmt.Printf("Close Time: %s\n", kline.CloseTime.Format(time.RFC3339))
			}

			// Save klines to file
			saveToFile(klines, fmt.Sprintf("mexc_klines_%s_%s.json", symbol, interval), logger)
		}
	}
}

// testMarketData tests the market data endpoint
func testMarketData(ctx context.Context, client *mexc.Client, logger *zerolog.Logger) {
	logger.Info().Msg("Testing market data endpoint...")

	// Test with a few popular symbols
	symbols := []string{"BTCUSDT", "ETHUSDT", "SOLUSDT"}

	for _, symbol := range symbols {
		// Get ticker
		ticker, err := client.GetMarketData(ctx, symbol)
		if err != nil {
			logger.Error().Err(err).Str("symbol", symbol).Msg("Failed to get ticker")
			continue
		}

		// Get order book
		orderBook, err := client.GetOrderBook(ctx, symbol, 10)
		if err != nil {
			logger.Error().Err(err).Str("symbol", symbol).Msg("Failed to get order book")
			continue
		}

		// Create a combined market data object
		marketData := struct {
			Symbol    string           `json:"symbol"`
			Ticker    *model.Ticker    `json:"ticker"`
			OrderBook *model.OrderBook `json:"orderBook"`
		}{
			Symbol:    symbol,
			Ticker:    ticker,
			OrderBook: orderBook,
		}

		logger.Info().
			Str("Symbol", symbol).
			Float64("Price", ticker.LastPrice).
			Int("Order Book Bids", len(orderBook.Bids)).
			Int("Order Book Asks", len(orderBook.Asks)).
			Msg("Market data retrieved successfully")

		// Save market data to file
		saveToFile(marketData, fmt.Sprintf("mexc_market_data_%s.json", symbol), logger)
	}
}

// saveToFile saves data to a JSON file
func saveToFile(data interface{}, filename string, logger *zerolog.Logger) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		logger.Error().Err(err).Str("filename", filename).Msg("Failed to marshal data to JSON")
		return
	}

	err = os.WriteFile(filename, jsonData, 0644)
	if err != nil {
		logger.Error().Err(err).Str("filename", filename).Msg("Failed to write data to file")
		return
	}

	logger.Info().Str("filename", filename).Msg("Data saved to file")
}
