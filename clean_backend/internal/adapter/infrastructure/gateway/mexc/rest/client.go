package rest

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"golang.org/x/time/rate"
)

const (
	// API endpoints
	apiBaseURL     = "https://api.mexc.com"
	apiV3Path      = "/api/v3"
	
	// Rate limits
	defaultRPS     = 10
	defaultBurst   = 20
	
	// HTTP methods
	methodGET      = "GET"
	methodPOST     = "POST"
	methodDELETE   = "DELETE"
	
	// Headers
	headerAPIKey   = "X-MEXC-APIKEY"
	headerSignature = "signature"
	
	// Query parameters
	paramTimestamp = "timestamp"
	paramSignature = "signature"
)

// Client represents a REST client for MEXC API
type Client struct {
	baseURL    string
	apiKey     string
	apiSecret  string
	httpClient *http.Client
	logger     *zerolog.Logger
	limiter    *rate.Limiter
}

// NewClient creates a new MEXC API client
func NewClient(apiKey, apiSecret string, logger *zerolog.Logger) *Client {
	return &Client{
		baseURL:    apiBaseURL,
		apiKey:     apiKey,
		apiSecret:  apiSecret,
		httpClient: &http.Client{Timeout: 10 * time.Second},
		logger:     logger,
		limiter:    rate.NewLimiter(rate.Limit(defaultRPS), defaultBurst),
	}
}

// SetRateLimit sets the rate limit for API requests
func (c *Client) SetRateLimit(rps, burst int) {
	c.limiter = rate.NewLimiter(rate.Limit(rps), burst)
}

// Get sends a GET request to the MEXC API
func (c *Client) Get(ctx context.Context, endpoint string, params url.Values, signed bool) ([]byte, error) {
	return c.Request(ctx, methodGET, endpoint, params, nil, signed)
}

// Post sends a POST request to the MEXC API
func (c *Client) Post(ctx context.Context, endpoint string, params url.Values, body interface{}, signed bool) ([]byte, error) {
	return c.Request(ctx, methodPOST, endpoint, params, body, signed)
}

// Delete sends a DELETE request to the MEXC API
func (c *Client) Delete(ctx context.Context, endpoint string, params url.Values, signed bool) ([]byte, error) {
	return c.Request(ctx, methodDELETE, endpoint, params, nil, signed)
}

// Request sends a request to the MEXC API
func (c *Client) Request(ctx context.Context, method, endpoint string, params url.Values, body interface{}, signed bool) ([]byte, error) {
	// Apply rate limiting
	if err := c.limiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit exceeded: %w", err)
	}
	
	// Build URL
	fullURL := c.baseURL + apiV3Path + endpoint
	
	// Add timestamp for signed requests
	if signed {
		timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
		if params == nil {
			params = url.Values{}
		}
		params.Set(paramTimestamp, timestamp)
	}
	
	// Create request
	var req *http.Request
	var err error
	
	if method == methodGET || method == methodDELETE {
		// For GET and DELETE, add params to URL
		if len(params) > 0 {
			fullURL += "?" + params.Encode()
		}
		
		// Add signature for signed requests
		if signed {
			signature := c.sign(params.Encode())
			fullURL += "&" + paramSignature + "=" + signature
		}
		
		req, err = http.NewRequestWithContext(ctx, method, fullURL, nil)
	} else {
		// For POST, add params to URL and body to request
		if len(params) > 0 {
			fullURL += "?" + params.Encode()
		}
		
		var bodyData []byte
		if body != nil {
			bodyData, err = json.Marshal(body)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal request body: %w", err)
			}
		}
		
		// Add signature for signed requests
		if signed {
			var queryString string
			if len(params) > 0 {
				queryString = params.Encode()
			}
			
			if body != nil {
				if len(queryString) > 0 {
					queryString += "&"
				}
				queryString += string(bodyData)
			}
			
			signature := c.sign(queryString)
			if len(params) > 0 {
				fullURL += "&"
			} else {
				fullURL += "?"
			}
			fullURL += paramSignature + "=" + signature
		}
		
		req, err = http.NewRequestWithContext(ctx, method, fullURL, bytes.NewReader(bodyData))
		if err == nil && body != nil {
			req.Header.Set("Content-Type", "application/json")
		}
	}
	
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	// Add API key header for signed requests
	if signed {
		req.Header.Set(headerAPIKey, c.apiKey)
	}
	
	// Send request
	c.logger.Debug().
		Str("method", method).
		Str("url", fullURL).
		Msg("Sending request to MEXC API")
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	
	// Read response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	
	// Check response status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %s, status code: %d, body: %s", resp.Status, resp.StatusCode, string(respBody))
	}
	
	return respBody, nil
}

// sign creates a signature for a request
func (c *Client) sign(queryString string) string {
	h := hmac.New(sha256.New, []byte(c.apiSecret))
	h.Write([]byte(queryString))
	return hex.EncodeToString(h.Sum(nil))
}

// GetServerTime gets the server time
func (c *Client) GetServerTime(ctx context.Context) (int64, error) {
	resp, err := c.Get(ctx, "/time", nil, false)
	if err != nil {
		return 0, err
	}
	
	var result struct {
		ServerTime int64 `json:"serverTime"`
	}
	
	if err := json.Unmarshal(resp, &result); err != nil {
		return 0, fmt.Errorf("failed to parse server time: %w", err)
	}
	
	return result.ServerTime, nil
}

// GetExchangeInfo gets the exchange information
func (c *Client) GetExchangeInfo(ctx context.Context) ([]byte, error) {
	return c.Get(ctx, "/exchangeInfo", nil, false)
}

// GetAccountInfo gets the account information
func (c *Client) GetAccountInfo(ctx context.Context) ([]byte, error) {
	return c.Get(ctx, "/account", nil, true)
}

// GetSymbolPrice gets the price of a symbol
func (c *Client) GetSymbolPrice(ctx context.Context, symbol string) (float64, error) {
	params := url.Values{}
	params.Set("symbol", strings.ToUpper(symbol))
	
	resp, err := c.Get(ctx, "/ticker/price", params, false)
	if err != nil {
		return 0, err
	}
	
	var result struct {
		Symbol string `json:"symbol"`
		Price  string `json:"price"`
	}
	
	if err := json.Unmarshal(resp, &result); err != nil {
		return 0, fmt.Errorf("failed to parse price: %w", err)
	}
	
	price, err := strconv.ParseFloat(result.Price, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse price as float: %w", err)
	}
	
	return price, nil
}

// GetOrderBook gets the order book for a symbol
func (c *Client) GetOrderBook(ctx context.Context, symbol string, limit int) ([]byte, error) {
	params := url.Values{}
	params.Set("symbol", strings.ToUpper(symbol))
	
	if limit > 0 {
		params.Set("limit", strconv.Itoa(limit))
	}
	
	return c.Get(ctx, "/depth", params, false)
}

// GetRecentTrades gets recent trades for a symbol
func (c *Client) GetRecentTrades(ctx context.Context, symbol string, limit int) ([]byte, error) {
	params := url.Values{}
	params.Set("symbol", strings.ToUpper(symbol))
	
	if limit > 0 {
		params.Set("limit", strconv.Itoa(limit))
	}
	
	return c.Get(ctx, "/trades", params, false)
}

// CreateOrder creates a new order
func (c *Client) CreateOrder(ctx context.Context, symbol, side, orderType string, quantity, price float64) ([]byte, error) {
	params := url.Values{}
	params.Set("symbol", strings.ToUpper(symbol))
	params.Set("side", strings.ToUpper(side))
	params.Set("type", strings.ToUpper(orderType))
	params.Set("quantity", strconv.FormatFloat(quantity, 'f', -1, 64))
	
	if price > 0 {
		params.Set("price", strconv.FormatFloat(price, 'f', -1, 64))
	}
	
	return c.Post(ctx, "/order", params, nil, true)
}

// CancelOrder cancels an order
func (c *Client) CancelOrder(ctx context.Context, symbol string, orderId int64) ([]byte, error) {
	params := url.Values{}
	params.Set("symbol", strings.ToUpper(symbol))
	params.Set("orderId", strconv.FormatInt(orderId, 10))
	
	return c.Delete(ctx, "/order", params, true)
}

// GetOrder gets an order
func (c *Client) GetOrder(ctx context.Context, symbol string, orderId int64) ([]byte, error) {
	params := url.Values{}
	params.Set("symbol", strings.ToUpper(symbol))
	params.Set("orderId", strconv.FormatInt(orderId, 10))
	
	return c.Get(ctx, "/order", params, true)
}

// GetOpenOrders gets all open orders
func (c *Client) GetOpenOrders(ctx context.Context, symbol string) ([]byte, error) {
	params := url.Values{}
	
	if symbol != "" {
		params.Set("symbol", strings.ToUpper(symbol))
	}
	
	return c.Get(ctx, "/openOrders", params, true)
}

// GetAllOrders gets all orders
func (c *Client) GetAllOrders(ctx context.Context, symbol string, limit int) ([]byte, error) {
	params := url.Values{}
	params.Set("symbol", strings.ToUpper(symbol))
	
	if limit > 0 {
		params.Set("limit", strconv.Itoa(limit))
	}
	
	return c.Get(ctx, "/allOrders", params, true)
}

// GetAccountTrades gets account trades
func (c *Client) GetAccountTrades(ctx context.Context, symbol string, limit int) ([]byte, error) {
	params := url.Values{}
	params.Set("symbol", strings.ToUpper(symbol))
	
	if limit > 0 {
		params.Set("limit", strconv.Itoa(limit))
	}
	
	return c.Get(ctx, "/myTrades", params, true)
}
