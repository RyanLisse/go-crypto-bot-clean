package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// Define the minimum structures needed for the MEXC API
type TickerResponse struct {
	Symbol    string  `json:"symbol"`
	LastPrice string  `json:"lastPrice"`
	Volume    string  `json:"volume"`
	BidPrice  string  `json:"bidPrice"`
	AskPrice  string  `json:"askPrice"`
	PriceChange string `json:"priceChange"`
	PriceChangePercent string `json:"priceChangePercent"`
}

type OrderBookResponse struct {
	LastUpdateID int64      `json:"lastUpdateId"`
	Bids         [][]string `json:"bids"`
	Asks         [][]string `json:"asks"`
}

type SymbolInfo struct {
	Symbol            string `json:"symbol"`
	Status            string `json:"status"`
	BaseAsset         string `json:"baseAsset"`
	QuoteAsset        string `json:"quoteAsset"`
	PricePrecision    int    `json:"pricePrecision"`
	QuantityPrecision int    `json:"quantityPrecision"`
}

type ExchangeInfoResponse struct {
	Timezone   string       `json:"timezone"`
	ServerTime int64        `json:"serverTime"`
	Symbols    []SymbolInfo `json:"symbols"`
}

// Simple HTTP client with timeout
func newHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 10 * time.Second,
	}
}

// Get ticker for a symbol
func getTicker(symbol string) (*TickerResponse, error) {
	client := newHTTPClient()
	url := fmt.Sprintf("https://api.mexc.com/api/v3/ticker/24hr?symbol=%s", symbol)
	
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	
	var ticker TickerResponse
	err = json.Unmarshal(body, &ticker)
	if err != nil {
		return nil, err
	}
	
	return &ticker, nil
}

// Get order book for a symbol
func getOrderBook(symbol string, limit int) (*OrderBookResponse, error) {
	client := newHTTPClient()
	url := fmt.Sprintf("https://api.mexc.com/api/v3/depth?symbol=%s&limit=%d", symbol, limit)
	
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	
	var orderBook OrderBookResponse
	err = json.Unmarshal(body, &orderBook)
	if err != nil {
		return nil, err
	}
	
	return &orderBook, nil
}

// Get exchange info
func getExchangeInfo() (*ExchangeInfoResponse, error) {
	client := newHTTPClient()
	url := "https://api.mexc.com/api/v3/exchangeInfo"
	
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	
	var exchangeInfo ExchangeInfoResponse
	err = json.Unmarshal(body, &exchangeInfo)
	if err != nil {
		return nil, err
	}
	
	return &exchangeInfo, nil
}

func main() {
	ctx := context.Background()
	_ = ctx // Not used in this simple example
	
	// Test getting ticker for BTCUSDT
	ticker, err := getTicker("BTCUSDT")
	if err != nil {
		fmt.Printf("Error getting ticker: %v\n", err)
		os.Exit(1)
	}
	
	// Print ticker as JSON
	tickerJSON, err := json.MarshalIndent(ticker, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling ticker to JSON: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Ticker JSON:")
	fmt.Println(string(tickerJSON))
	
	// Test getting order book for BTCUSDT
	orderBook, err := getOrderBook("BTCUSDT", 5)
	if err != nil {
		fmt.Printf("Error getting order book: %v\n", err)
		os.Exit(1)
	}
	
	// Print order book as JSON
	orderBookJSON, err := json.MarshalIndent(orderBook, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling order book to JSON: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("\nOrder Book JSON:")
	fmt.Println(string(orderBookJSON))
	
	// Test getting exchange info
	exchangeInfo, err := getExchangeInfo()
	if err != nil {
		fmt.Printf("Error getting exchange info: %v\n", err)
		os.Exit(1)
	}
	
	// Print first 5 symbols from exchange info
	fmt.Println("\nFirst 5 symbols from exchange info:")
	for i, symbol := range exchangeInfo.Symbols {
		if i >= 5 {
			break
		}
		fmt.Printf("%d. %s (%s/%s)\n", i+1, symbol.Symbol, symbol.BaseAsset, symbol.QuoteAsset)
	}
	
	fmt.Println("\nTest completed successfully")
}
