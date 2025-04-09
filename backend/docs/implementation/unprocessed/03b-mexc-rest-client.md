# MEXC REST API Client Implementation

This document provides implementation details for the REST API client part of the MEXC interface.

## Client Setup and Authentication

```go
// internal/platform/mexc/rest/client.go
package rest

import (
    "bytes"
    "crypto/hmac"
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"
    "io/ioutil"
    "net/http"
    "net/url"
    "sort"
    "strconv"
    "strings"
    "time"
)

// Client is the REST API client for MEXC
type Client struct {
    apiKey     string
    secretKey  string
    baseURL    string
    httpClient *http.Client
}

// NewClient creates a new REST client for MEXC API
func NewClient(apiKey, secretKey, baseURL string) (*Client, error) {
    if baseURL == "" {
        baseURL = "https://api.mexc.com" // Default API URL
    }
    
    return &Client{
        apiKey:     apiKey,
        secretKey:  secretKey,
        baseURL:    baseURL,
        httpClient: &http.Client{
            Timeout: 10 * time.Second,
        },
    }, nil
}

// sign creates a HMAC SHA256 signature required for authenticated API calls
func (c *Client) sign(payload string) string {
    h := hmac.New(sha256.New, []byte(c.secretKey))
    h.Write([]byte(payload))
    return hex.EncodeToString(h.Sum(nil))
}

// buildQueryString creates a sorted query string from parameters
func buildQueryString(params map[string]string) string {
    if len(params) == 0 {
        return ""
    }
    
    keys := make([]string, 0, len(params))
    for k := range params {
        keys = append(keys, k)
    }
    sort.Strings(keys)
    
    var b strings.Builder
    for i, k := range keys {
        if i > 0 {
            b.WriteString("&")
        }
        b.WriteString(url.QueryEscape(k))
        b.WriteString("=")
        b.WriteString(url.QueryEscape(params[k]))
    }
    
    return b.String()
}

// createAuthenticatedRequest creates a request with necessary authentication headers
func (c *Client) createAuthenticatedRequest(method, endpoint string, params map[string]string) (*http.Request, error) {
    timestamp := time.Now().UnixMilli()
    
    // Build query string
    query := buildQueryString(params)
    
    // Create signature string
    signString := fmt.Sprintf("timestamp=%d&%s", timestamp, query)
    signature := c.sign(signString)
    
    // Create full URL
    var fullURL string
    if method == http.MethodGet {
        fullURL = fmt.Sprintf("%s%s?%s&timestamp=%d&signature=%s", 
            c.baseURL, endpoint, query, timestamp, signature)
    } else {
        fullURL = fmt.Sprintf("%s%s?timestamp=%d&signature=%s", 
            c.baseURL, endpoint, timestamp, signature)
    }
    
    // Create request
    req, err := http.NewRequest(method, fullURL, nil)
    if err != nil {
        return nil, err
    }
    
    // Add headers
    req.Header.Add("Content-Type", "application/json")
    req.Header.Add("X-MEXC-APIKEY", c.apiKey)
    
    // Add body for POST/PUT/DELETE requests
    if method != http.MethodGet && len(params) > 0 {
        jsonBody, err := json.Marshal(params)
        if err != nil {
            return nil, err
        }
        req.Body = ioutil.NopCloser(bytes.NewBuffer(jsonBody))
    }
    
    return req, nil
}
```

## Market Data Implementation

```go
// internal/platform/mexc/rest/market.go
package rest

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "strconv"
    "time"
    
    "github.com/ryanlisse/cryptobot/internal/domain/models"
)

// MexcTicker represents the MEXC API response for a ticker
type MexcTicker struct {
    Symbol             string `json:"symbol"`
    LastPrice          string `json:"lastPrice"`
    BidPrice           string `json:"bidPrice"`
    AskPrice           string `json:"askPrice"`
    Volume             string `json:"volume"`
    High               string `json:"highPrice"`
    Low                string `json:"lowPrice"`
    PriceChange        string `json:"priceChange"`
    PriceChangePercent string `json:"priceChangePercent"`
}

// GetTicker fetches the current ticker data for a specific symbol
func (c *Client) GetTicker(ctx context.Context, symbol string) (*models.Ticker, error) {
    endpoint := "/api/v3/ticker/24hr"
    url := fmt.Sprintf("%s%s?symbol=%s", c.baseURL, endpoint, symbol)
    
    req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
    if err != nil {
        return nil, fmt.Errorf("creating ticker request: %w", err)
    }
    
    resp, err := c.httpClient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("executing ticker request: %w", err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("API error: %s", resp.Status)
    }
    
    var mexcTicker MexcTicker
    if err := json.NewDecoder(resp.Body).Decode(&mexcTicker); err != nil {
        return nil, fmt.Errorf("decoding ticker response: %w", err)
    }
    
    return convertToTicker(&mexcTicker)
}

// GetAllTickers fetches tickers for all symbols
func (c *Client) GetAllTickers(ctx context.Context) (map[string]*models.Ticker, error) {
    endpoint := "/api/v3/ticker/24hr"
    url := fmt.Sprintf("%s%s", c.baseURL, endpoint)
    
    req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
    if err != nil {
        return nil, fmt.Errorf("creating all tickers request: %w", err)
    }
    
    resp, err := c.httpClient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("executing all tickers request: %w", err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("API error: %s", resp.Status)
    }
    
    var mexcTickers []MexcTicker
    if err := json.NewDecoder(resp.Body).Decode(&mexcTickers); err != nil {
        return nil, fmt.Errorf("decoding tickers response: %w", err)
    }
    
    result := make(map[string]*models.Ticker)
    for _, mexcTicker := range mexcTickers {
        ticker, err := convertToTicker(&mexcTicker)
        if err != nil {
            continue // Skip invalid tickers
        }
        result[mexcTicker.Symbol] = ticker
    }
    
    return result, nil
}

// GetKlines fetches candlestick data for a symbol
func (c *Client) GetKlines(ctx context.Context, symbol string, interval string, limit int) ([]*models.Kline, error) {
    endpoint := "/api/v3/klines"
    url := fmt.Sprintf("%s%s?symbol=%s&interval=%s&limit=%d", 
        c.baseURL, endpoint, symbol, interval, limit)
    
    req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
    if err != nil {
        return nil, fmt.Errorf("creating klines request: %w", err)
    }
    
    resp, err := c.httpClient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("executing klines request: %w", err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("API error: %s", resp.Status)
    }
    
    var rawKlines [][]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&rawKlines); err != nil {
        return nil, fmt.Errorf("decoding klines response: %w", err)
    }
    
    klines := make([]*models.Kline, 0, len(rawKlines))
    for _, raw := range rawKlines {
        if len(raw) < 6 {
            continue
        }
        
        openTime := int64(raw[0].(float64))
        closeTime := int64(raw[6].(float64))
        
        kline := &models.Kline{
            Symbol:    symbol,
            Interval:  interval,
            OpenTime:  time.Unix(openTime/1000, 0),
            CloseTime: time.Unix(closeTime/1000, 0),
            Open:      parseFloat64(raw[1]),
            High:      parseFloat64(raw[2]),
            Low:       parseFloat64(raw[3]),
            Close:     parseFloat64(raw[4]),
            Volume:    parseFloat64(raw[5]),
        }
        
        klines = append(klines, kline)
    }
    
    return klines, nil
}

// Helper function to parse float64 from interface{}
func parseFloat64(value interface{}) float64 {
    switch v := value.(type) {
    case float64:
        return v
    case string:
        f, _ := strconv.ParseFloat(v, 64)
        return f
    default:
        return 0
    }
}

// Helper function to convert MEXC ticker to domain model
func convertToTicker(mexcTicker *MexcTicker) (*models.Ticker, error) {
    price, err := strconv.ParseFloat(mexcTicker.LastPrice, 64)
    if err != nil {
        return nil, fmt.Errorf("parsing price: %w", err)
    }
    
    priceChange, err := strconv.ParseFloat(mexcTicker.PriceChange, 64)
    if err != nil {
        priceChange = 0
    }
    
    priceChangePct, err := strconv.ParseFloat(mexcTicker.PriceChangePercent, 64)
    if err != nil {
        priceChangePct = 0
    }
    
    volume, err := strconv.ParseFloat(mexcTicker.Volume, 64)
    if err != nil {
        volume = 0
    }
    
    high, err := strconv.ParseFloat(mexcTicker.High, 64)
    if err != nil {
        high = 0
    }
    
    low, err := strconv.ParseFloat(mexcTicker.Low, 64)
    if err != nil {
        low = 0
    }
    
    return &models.Ticker{
        Symbol:         mexcTicker.Symbol,
        Price:          price,
        PriceChange:    priceChange,
        PriceChangePct: priceChangePct,
        Volume:         volume,
        High24h:        high,
        Low24h:         low,
        Timestamp:      time.Now(),
    }, nil
}
```

## Error Handling and Custom Error Types

```go
// internal/platform/mexc/rest/errors.go
package rest

import (
    "fmt"
    "net/http"
)

// Common API error codes
const (
    ErrCodeUnknown      = -1000
    ErrCodeUnauthorized = -2015
    ErrCodeRateLimited  = -1003
    ErrCodeIPBanned     = -1004
)

// APIError represents an error returned by the MEXC API
type APIError struct {
    Code    int    `json:"code"`
    Message string `json:"msg"`
    Status  int    // HTTP status code
}

// Error implements the error interface
func (e *APIError) Error() string {
    return fmt.Sprintf("MEXC API error %d: %s", e.Code, e.Message)
}

// IsRateLimited returns true if the error indicates we've been rate limited
func (e *APIError) IsRateLimited() bool {
    return e.Code == ErrCodeRateLimited
}

// IsUnauthorized returns true if the error indicates invalid API credentials
func (e *APIError) IsUnauthorized() bool {
    return e.Code == ErrCodeUnauthorized
}

// IsIPBanned returns true if our IP has been temporarily banned
func (e *APIError) IsIPBanned() bool {
    return e.Code == ErrCodeIPBanned
}

// parseAPIError attempts to parse an API error from the response
func parseAPIError(resp *http.Response) error {
    var apiErr APIError
    apiErr.Status = resp.StatusCode
    
    if err := json.NewDecoder(resp.Body).Decode(&apiErr); err != nil {
        // If we can't decode the JSON error, create a generic one
        return &APIError{
            Code:    ErrCodeUnknown,
            Message: fmt.Sprintf("HTTP error: %s", resp.Status),
            Status:  resp.StatusCode,
        }
    }
    
    return &apiErr
}
```

## Retry Mechanism with Exponential Backoff

```go
// internal/platform/mexc/rest/retry.go
package rest

import (
    "context"
    "math"
    "math/rand"
    "time"
)

// RetryConfig holds parameters for the retry mechanism
type RetryConfig struct {
    MaxRetries  int
    InitialWait time.Duration
    MaxWait     time.Duration
    Factor      float64  // Multiplier for each retry
    Jitter      float64  // Random jitter factor between 0-1
}

// DefaultRetryConfig provides sensible defaults for API retries
func DefaultRetryConfig() RetryConfig {
    return RetryConfig{
        MaxRetries:  3,
        InitialWait: 500 * time.Millisecond,
        MaxWait:     5 * time.Second,
        Factor:      2.0,
        Jitter:      0.2,
    }
}

// RetryableFunc is a function that can be retried
type RetryableFunc func() error

// Retry executes the given function with exponential backoff
func (c *Client) Retry(ctx context.Context, config RetryConfig, fn RetryableFunc) error {
    var err error
    wait := config.InitialWait
    
    for attempt := 0; attempt <= config.MaxRetries; attempt++ {
        // Execute the function
        err = fn()
        
        // Success, return immediately
        if err == nil {
            return nil
        }
        
        // Check if it's an API error
        if apiErr, ok := err.(*APIError); ok {
            // Don't retry certain errors
            if !apiErr.IsRateLimited() && !isTransientError(apiErr) {
                return err
            }
        }
        
        // Last attempt, return the error
        if attempt == config.MaxRetries {
            return err
        }
        
        // Add some randomness to the wait time
        jitter := 1.0
        if config.Jitter > 0 {
            jitter = 1.0 + rand.Float64()*config.Jitter
        }
        
        // Calculate next wait time with exponential backoff
        wait = time.Duration(float64(wait) * config.Factor * jitter)
        if wait > config.MaxWait {
            wait = config.MaxWait
        }
        
        // Wait before next attempt, or return if context is done
        select {
        case <-ctx.Done():
            return ctx.Err()
        case <-time.After(wait):
            // Continue to next attempt
        }
    }
    
    return err
}

// isTransientError determines if an error is temporary and should be retried
func isTransientError(err *APIError) bool {
    // Network or server errors (5xx status codes)
    if err.Status >= 500 && err.Status < 600 {
        return true
    }
    
    // Specific API error codes that are worth retrying
    transientCodes := map[int]bool{
        -1003: true, // Rate limit
        -1004: false, // IP ban (don't retry immediately)
        -1008: true, // Server busy
        -1021: true, // Timestamp outside of recvWindow
    }
    
    retry, exists := transientCodes[err.Code]
    return exists && retry
}
```

This separation into multiple files makes the codebase more manageable and easier to understand, with each file focusing on a specific aspect of the REST client implementation.
