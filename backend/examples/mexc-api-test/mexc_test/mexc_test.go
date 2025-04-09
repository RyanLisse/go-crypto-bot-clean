package mexc_test

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestMexcAPI(t *testing.T) {
	// Get API keys from environment variables
	apiKey := os.Getenv("MEXC_API_KEY")
	secretKey := os.Getenv("MEXC_SECRET_KEY")

	if apiKey == "" || secretKey == "" {
		t.Skip("MEXC_API_KEY and MEXC_SECRET_KEY environment variables must be set")
	}

	// Test 1: Get account information
	t.Run("GetAccount", func(t *testing.T) {
		account, err := getAccount(apiKey, secretKey)
		if err != nil {
			t.Fatalf("Error getting account: %v", err)
		}

		t.Logf("Account Information:")
		t.Logf("  Maker Commission: %d", account.MakerCommission)
		t.Logf("  Taker Commission: %d", account.TakerCommission)
		t.Logf("  Balances:")

		// Print non-zero balances
		for _, balance := range account.Balances {
			free, _ := strconv.ParseFloat(balance.Free, 64)
			locked, _ := strconv.ParseFloat(balance.Locked, 64)
			if free > 0 || locked > 0 {
				t.Logf("    %s: Free=%s, Locked=%s", balance.Asset, balance.Free, balance.Locked)
			}
		}
	})

	// Test 2: Get ticker for BTC
	t.Run("GetTicker", func(t *testing.T) {
		ticker, err := getTicker("BTCUSDT")
		if err != nil {
			t.Fatalf("Error getting ticker: %v", err)
		}

		t.Logf("Ticker for BTCUSDT:")
		t.Logf("  Price: %s", ticker.LastPrice)
		t.Logf("  Volume: %s", ticker.Volume)
		t.Logf("  High: %s", ticker.HighPrice)
		t.Logf("  Low: %s", ticker.LowPrice)
	})

	// Test 3: Get klines for BTC
	t.Run("GetKlines", func(t *testing.T) {
		klines, err := getKlines("BTCUSDT", "1h", 5)
		if err != nil {
			t.Fatalf("Error getting klines: %v", err)
		}

		t.Logf("Klines for BTCUSDT (1h):")
		for i, kline := range klines {
			openTime := time.Unix(int64(kline[0].(float64))/1000, 0)
			open := kline[1].(string)
			high := kline[2].(string)
			low := kline[3].(string)
			close := kline[4].(string)
			volume := kline[5].(string)

			t.Logf("  %d: Time=%s, Open=%s, High=%s, Low=%s, Close=%s, Volume=%s",
				i+1, openTime.Format(time.RFC3339), open, high, low, close, volume)
		}
	})
}

// Account represents the account information returned by the MEXC API
type Account struct {
	MakerCommission  int       `json:"makerCommission"`
	TakerCommission  int       `json:"takerCommission"`
	BuyerCommission  int       `json:"buyerCommission"`
	SellerCommission int       `json:"sellerCommission"`
	CanTrade         bool      `json:"canTrade"`
	CanWithdraw      bool      `json:"canWithdraw"`
	CanDeposit       bool      `json:"canDeposit"`
	UpdateTime       int64     `json:"updateTime"`
	AccountType      string    `json:"accountType"`
	Balances         []Balance `json:"balances"`
}

// Balance represents a balance for a specific asset
type Balance struct {
	Asset  string `json:"asset"`
	Free   string `json:"free"`
	Locked string `json:"locked"`
}

// Ticker represents the ticker information returned by the MEXC API
type Ticker struct {
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

// getAccount gets the account information from the MEXC API
func getAccount(apiKey, secretKey string) (*Account, error) {
	baseURL := "https://api.mexc.com"
	endpoint := "/api/v3/account"

	// Create timestamp
	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)

	// Create query parameters
	params := url.Values{}
	params.Set("timestamp", timestamp)

	// Create signature
	signature := createSignature(params.Encode(), secretKey)
	params.Set("signature", signature)

	// Create request URL
	reqURL := fmt.Sprintf("%s%s?%s", baseURL, endpoint, params.Encode())

	// Create request
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Add headers
	req.Header.Set("X-MBX-APIKEY", apiKey)

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
	var account Account
	if err := json.Unmarshal(body, &account); err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}

	return &account, nil
}

// getTicker gets the ticker information for a symbol from the MEXC API
func getTicker(symbol string) (*Ticker, error) {
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
	var ticker Ticker
	if err := json.Unmarshal(body, &ticker); err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}

	return &ticker, nil
}

// getKlines gets the klines for a symbol from the MEXC API
func getKlines(symbol, interval string, limit int) ([][]interface{}, error) {
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

// createSignature creates a signature for the MEXC API
func createSignature(payload, secretKey string) string {
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(payload))
	return hex.EncodeToString(h.Sum(nil))
}
