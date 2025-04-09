package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func main() {
	// Use the API keys from the config file
	apiKey := "mx0vglSSj7Lg3y2Y27"
	secretKey := "bf25f45c440d4550977e1c65ca664fd0"
	fmt.Printf("Using API key: %s\n", apiKey)

	// Test 1: Get account information
	fmt.Println("=== Testing GetAccount ===")
	account, err := getAccountInfo(apiKey, secretKey)
	if err != nil {
		fmt.Printf("Error getting account: %v\n", err)
	} else {
		fmt.Printf("Account Information:\n")
		fmt.Printf("  Maker Commission: %d\n", account.MakerCommission)
		fmt.Printf("  Taker Commission: %d\n", account.TakerCommission)
		fmt.Printf("  Balances:\n")

		// Print non-zero balances
		for _, balance := range account.Balances {
			free, _ := strconv.ParseFloat(balance.Free, 64)
			locked, _ := strconv.ParseFloat(balance.Locked, 64)
			if free > 0 || locked > 0 {
				fmt.Printf("    %s: Free=%s, Locked=%s\n", balance.Asset, balance.Free, balance.Locked)
			}
		}
	}

	// Test 2: Get ticker for BTC
	fmt.Println("\n=== Testing GetTicker ===")
	ticker, err := getTickerInfo("BTCUSDT")
	if err != nil {
		fmt.Printf("Error getting ticker: %v\n", err)
	} else {
		fmt.Printf("Ticker for BTCUSDT:\n")
		fmt.Printf("  Price: %s\n", ticker.LastPrice)
		fmt.Printf("  Volume: %s\n", ticker.Volume)
		fmt.Printf("  High: %s\n", ticker.HighPrice)
		fmt.Printf("  Low: %s\n", ticker.LowPrice)
	}

	// Test 3: Get klines for BTC
	fmt.Println("\n=== Testing GetKlines ===")
	// Try with a different interval format
	klines, err := getKlinesData("BTCUSDT", "1hour", 5)
	// If that fails, try with another format
	if err != nil {
		fmt.Printf("Error with interval '1hour': %v\n", err)
		klines, err = getKlinesData("BTCUSDT", "1d", 5)
	}
	// If that fails too, try with another format
	if err != nil {
		fmt.Printf("Error with interval '1d': %v\n", err)
		klines, err = getKlinesData("BTCUSDT", "1m", 5)
	}
	if err != nil {
		fmt.Printf("Error getting klines: %v\n", err)
	} else {
		fmt.Printf("Klines for BTCUSDT (1h):\n")
		for i, kline := range klines {
			openTime := time.Unix(int64(kline[0].(float64))/1000, 0)
			open := kline[1].(string)
			high := kline[2].(string)
			low := kline[3].(string)
			close := kline[4].(string)
			volume := kline[5].(string)

			fmt.Printf("  %d: Time=%s, Open=%s, High=%s, Low=%s, Close=%s, Volume=%s\n",
				i+1, openTime.Format(time.RFC3339), open, high, low, close, volume)
		}
	}
}

// AccountInfo represents the account information returned by the MEXC API
type AccountInfo struct {
	MakerCommission  int           `json:"makerCommission"`
	TakerCommission  int           `json:"takerCommission"`
	BuyerCommission  int           `json:"buyerCommission"`
	SellerCommission int           `json:"sellerCommission"`
	CanTrade         bool          `json:"canTrade"`
	CanWithdraw      bool          `json:"canWithdraw"`
	CanDeposit       bool          `json:"canDeposit"`
	UpdateTime       int64         `json:"updateTime"`
	AccountType      string        `json:"accountType"`
	Balances         []BalanceInfo `json:"balances"`
}

// BalanceInfo represents a balance for a specific asset
type BalanceInfo struct {
	Asset  string `json:"asset"`
	Free   string `json:"free"`
	Locked string `json:"locked"`
}

// TickerInfo represents the ticker information returned by the MEXC API
type TickerInfo struct {
	Symbol             string `json:"symbol"`
	PriceChange        string `json:"priceChange"`
	PriceChangePercent string `json:"priceChangePercent"`
	PrevClosePrice     string `json:"prevClosePrice"`
	LastPrice          string `json:"lastPrice"`
	BidPrice           string `json:"bidPrice"`
	BidQty             string `json:"bidQty"`
	AskPrice           string `json:"askPrice"`
	AskQty             string `json:"askQty"`
	OpenPrice          string `json:"openPrice"`
	HighPrice          string `json:"highPrice"`
	LowPrice           string `json:"lowPrice"`
	Volume             string `json:"volume"`
	QuoteVolume        string `json:"quoteVolume"`
	OpenTime           int64  `json:"openTime"`
	CloseTime          int64  `json:"closeTime"`
	Count              int    `json:"count"`
}

// getAccountInfo gets the account information from the MEXC API
func getAccountInfo(apiKey, secretKey string) (*AccountInfo, error) {
	baseURL := "https://api.mexc.com"
	endpoint := "/api/v3/account"

	// Create timestamp
	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)

	// Create query parameters
	params := url.Values{}
	params.Set("timestamp", timestamp)

	// Create signature
	signature := createHmacSignature(params.Encode(), secretKey)
	params.Set("signature", signature)

	// Create request URL
	reqURL := fmt.Sprintf("%s%s?%s", baseURL, endpoint, params.Encode())
	fmt.Printf("Request URL: %s\n", reqURL)

	// Create request
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Add headers
	req.Header.Set("X-MEXC-APIKEY", apiKey)
	fmt.Printf("API Key Header: %s\n", apiKey)

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	// Check for error
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error response: %s", string(body))
	}

	// Parse response
	var account AccountInfo
	if err := json.Unmarshal(body, &account); err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}

	return &account, nil
}

// getTickerInfo gets the ticker information for a symbol from the MEXC API
func getTickerInfo(symbol string) (*TickerInfo, error) {
	baseURL := "https://api.mexc.com"
	endpoint := "/api/v3/ticker/24hr"

	// Create query parameters
	params := url.Values{}
	params.Set("symbol", symbol)

	// Create request URL
	reqURL := fmt.Sprintf("%s%s?%s", baseURL, endpoint, params.Encode())

	// Create request
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	// Check for error
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error response: %s", string(body))
	}

	// Parse response
	var ticker TickerInfo
	if err := json.Unmarshal(body, &ticker); err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}

	return &ticker, nil
}

// getKlinesData gets the klines for a symbol from the MEXC API
func getKlinesData(symbol, interval string, limit int) ([][]interface{}, error) {
	baseURL := "https://api.mexc.com"
	endpoint := "/api/v3/klines"

	// Create query parameters
	params := url.Values{}
	params.Set("symbol", symbol)
	params.Set("interval", interval)
	params.Set("limit", strconv.Itoa(limit))

	// Create request URL
	reqURL := fmt.Sprintf("%s%s?%s", baseURL, endpoint, params.Encode())

	// Create request
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	// Check for error
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error response: %s", string(body))
	}

	// Parse response
	var klines [][]interface{}
	if err := json.Unmarshal(body, &klines); err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}

	return klines, nil
}

// createHmacSignature creates a signature for the MEXC API
func createHmacSignature(payload, secretKey string) string {
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(payload))
	return hex.EncodeToString(h.Sum(nil))
}
