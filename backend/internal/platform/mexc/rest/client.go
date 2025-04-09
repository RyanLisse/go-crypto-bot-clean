package rest

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"go-crypto-bot-clean/backend/internal/domain/models"
	"go-crypto-bot-clean/backend/internal/platform/mexc/cache"
	"go-crypto-bot-clean/backend/pkg/ratelimiter"
	"go.uber.org/zap"
)

const (
	// Default API endpoints
	defaultBaseURL = "https://api.mexc.com"
	defaultWsURL   = "wss://stream.mexc.com/ws"

	// API rate limits (according to MEXC documentation)
	requestsPerSecond = 20
	requestsPerMinute = 1200

	// HTTP request timeouts
	defaultTimeout = 10 * time.Second

	// Retry configuration
	maxRetries     = 3
	retryBaseDelay = 500 * time.Millisecond
	retryMaxDelay  = 5 * time.Second
)

// Client is the MEXC REST API client
type Client struct {
	// HTTP client
	httpClient *http.Client

	// API credentials
	apiKey    string
	secretKey string

	// API endpoints
	baseURL string

	// Rate limiters
	publicRateLimiter  *ratelimiter.TokenBucketRateLimiter // For public endpoints
	privateRateLimiter *ratelimiter.TokenBucketRateLimiter // For private endpoints

	// Caches
	tickerCache    *cache.TickerCache
	orderBookCache *cache.OrderBookCache
	klineCache     *cache.KlineCache
	newCoinCache   *cache.NewCoinCache

	// Logger
	logger *zap.Logger
}

// ClientOption is a function that configures a Client
type ClientOption func(*Client)

// WithBaseURL sets the base URL for the client
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// WithPublicRateLimiter sets a custom public rate limiter
func WithPublicRateLimiter(limiter *ratelimiter.TokenBucketRateLimiter) ClientOption {
	return func(c *Client) {
		c.publicRateLimiter = limiter
	}
}

// WithPrivateRateLimiter sets a custom private rate limiter
func WithPrivateRateLimiter(limiter *ratelimiter.TokenBucketRateLimiter) ClientOption {
	return func(c *Client) {
		c.privateRateLimiter = limiter
	}
}

// WithLogger sets a custom logger
func WithLogger(logger *zap.Logger) ClientOption {
	return func(c *Client) {
		c.logger = logger
	}
}

// NewClient creates a new MEXC REST API client
func NewClient(apiKey, secretKey string, options ...ClientOption) (*Client, error) {
	// Create default logger
	logger, _ := zap.NewProduction()

	c := &Client{
		apiKey:    apiKey,
		secretKey: secretKey,
		baseURL:   defaultBaseURL,
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
		publicRateLimiter:  ratelimiter.NewTokenBucketRateLimiter(float64(requestsPerSecond), float64(requestsPerSecond)),
		privateRateLimiter: ratelimiter.NewTokenBucketRateLimiter(float64(requestsPerSecond/2), float64(requestsPerSecond/2)),
		tickerCache:        cache.NewTickerCache(),
		orderBookCache:     cache.NewOrderBookCache(),
		klineCache:         cache.NewKlineCache(),
		newCoinCache:       cache.NewNewCoinCache(),
		logger:             logger,
	}

	// Apply all options
	for _, option := range options {
		option(c)
	}

	return c, nil
}

// makeRequest makes an HTTP request to the MEXC API
func (c *Client) makeRequest(ctx context.Context, method, endpoint string, params url.Values, needAuth bool, result interface{}) error {
	// Select appropriate rate limiter based on auth requirement
	limiter := c.publicRateLimiter
	if needAuth {
		limiter = c.privateRateLimiter
	}

	// Wait for rate limit token
	if err := limiter.Wait(ctx); err != nil {
		return err
	}

	// Build full URL
	reqURL := fmt.Sprintf("%s%s", c.baseURL, endpoint)

	// Create request
	var req *http.Request
	var err error

	// Add timestamp and signature for authenticated requests
	if needAuth {
		if params == nil {
			params = url.Values{}
		}

		// Add timestamp (required for authenticated requests)
		timestamp := fmt.Sprintf("%d", time.Now().UnixMilli())
		params.Set("timestamp", timestamp)

		// Add signature
		signature := c.generateSignature(params)
		params.Set("signature", signature)

		// Add query parameters to URL for GET requests or request body for POST/DELETE
		if method == http.MethodGet || method == http.MethodDelete {
			reqURL = fmt.Sprintf("%s?%s", reqURL, params.Encode())
			req, err = http.NewRequestWithContext(ctx, method, reqURL, nil)
		} else {
			req, err = http.NewRequestWithContext(ctx, method, reqURL, strings.NewReader(params.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		}

		// Add API key header
		req.Header.Set("X-MEXC-APIKEY", c.apiKey)
	} else {
		// Public API - add query parameters to URL
		if len(params) > 0 {
			reqURL = fmt.Sprintf("%s?%s", reqURL, params.Encode())
		}

		req, err = http.NewRequestWithContext(ctx, method, reqURL, nil)
	}

	if err != nil {
		return &RequestError{Err: err, Message: "failed to create request"}
	}

	// Set common headers
	req.Header.Set("Accept", "application/json")

	// Execute request with retries
	var resp *http.Response
	var retryCount int
	var respBody []byte

	for retryCount <= maxRetries {
		// Perform the request
		resp, err = c.httpClient.Do(req)

		// Handle connection errors
		if err != nil {
			if retryCount == maxRetries {
				return &RequestError{Err: err, Message: "request failed after retries"}
			}

			// Exponential backoff with jitter
			retryDelay := retryBaseDelay * time.Duration(1<<retryCount)
			if retryDelay > retryMaxDelay {
				retryDelay = retryMaxDelay
			}

			select {
			case <-time.After(retryDelay):
				retryCount++
				continue
			case <-ctx.Done():
				return ctx.Err()
			}
		}

		defer resp.Body.Close()

		// Read response body
		respBody, err = io.ReadAll(resp.Body)
		if err != nil {
			if retryCount == maxRetries {
				return &RequestError{Err: err, Message: "failed to read response body after retries"}
			}
			retryCount++
			continue
		}

		// Check for error response
		if resp.StatusCode >= 400 {
			// Handle rate limiting
			if resp.StatusCode == http.StatusTooManyRequests {
				if retryCount == maxRetries {
					var apiErr APIError
					if err := json.Unmarshal(respBody, &apiErr); err == nil {
						return &apiErr
					}
					return &RequestError{Message: "rate limit exceeded"}
				}

				// Retry with backoff
				retryDelay := retryBaseDelay * time.Duration(1<<retryCount)
				if retryDelay > retryMaxDelay {
					retryDelay = retryMaxDelay
				}

				select {
				case <-time.After(retryDelay):
					retryCount++
					continue
				case <-ctx.Done():
					return ctx.Err()
				}
			}

			// Parse API error
			var apiErr APIError
			if err := json.Unmarshal(respBody, &apiErr); err != nil {
				return &RequestError{
					Message: fmt.Sprintf("HTTP error: %d", resp.StatusCode),
				}
			}
			return &apiErr
		}

		// Success response - unmarshal and return
		if result != nil {
			if err := json.Unmarshal(respBody, result); err != nil {
				return &UnmarshalError{
					Err:     err,
					Body:    respBody,
					Message: "failed to unmarshal response",
				}
			}
		}

		// Success!
		break
	}

	return nil
}

// generateSignature generates an HMAC SHA256 signature for authentication
func (c *Client) generateSignature(params url.Values) string {
	// Create a string to sign from the query parameters
	queryString := params.Encode()

	// Create HMAC SHA256 Hash
	h := hmac.New(sha256.New, []byte(c.secretKey))
	h.Write([]byte(queryString))

	// Return hex encoded signature
	return hex.EncodeToString(h.Sum(nil))
}

// parseOrderSide converts model order side to API parameter
func parseOrderSide(side models.OrderSide) string {
	return string(side)
}

// parseOrderType converts model order type to API parameter
func parseOrderType(orderType models.OrderType) string {
	return string(orderType)
}

// parseOrderStatusFromAPI converts API order status to model order status
func parseOrderStatusFromAPI(status string) models.OrderStatus {
	switch strings.ToUpper(status) {
	case "NEW":
		return models.OrderStatusNew
	case "PARTIALLY_FILLED":
		return models.OrderStatusPartiallyFilled
	case "FILLED":
		return models.OrderStatusFilled
	case "CANCELED", "CANCELLED":
		return models.OrderStatusCanceled
	case "REJECTED":
		return models.OrderStatusRejected
	default:
		return models.OrderStatusRejected
	}
}

// parseTime parses a UNIX timestamp in milliseconds
func parseTime(timestamp interface{}) (time.Time, error) {
	switch v := timestamp.(type) {
	case int64:
		return time.UnixMilli(v), nil
	case float64:
		return time.UnixMilli(int64(v)), nil
	case string:
		i, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return time.Time{}, err
		}
		return time.UnixMilli(i), nil
	default:
		return time.Time{}, fmt.Errorf("unsupported timestamp type: %T", timestamp)
	}
}

// GetTicker fetches real-time price information for a specific trading pair
func (c *Client) GetTicker(ctx context.Context, symbol string) (*models.Ticker, error) {
	// Check cache first
	if ticker, found := c.tickerCache.GetTicker(symbol); found {
		c.logger.Debug("Using cached ticker", zap.String("symbol", symbol))
		return ticker, nil
	}

	params := url.Values{}
	params.Set("symbol", symbol)

	var response struct {
		Symbol             string `json:"symbol"`
		LastPrice          string `json:"lastPrice"`
		Volume             string `json:"volume"`
		PriceChange        string `json:"priceChange"`
		PriceChangePercent string `json:"priceChangePercent"`
		HighPrice          string `json:"highPrice"`
		LowPrice           string `json:"lowPrice"`
	}

	err := c.makeRequest(ctx, http.MethodGet, "/api/v3/ticker/24hr", params, false, &response)
	if err != nil {
		return nil, err
	}

	// Parse string values to float64
	price, _ := strconv.ParseFloat(response.LastPrice, 64)
	volume, _ := strconv.ParseFloat(response.Volume, 64)
	priceChange, _ := strconv.ParseFloat(response.PriceChange, 64)
	priceChangePct, _ := strconv.ParseFloat(response.PriceChangePercent, 64)
	high, _ := strconv.ParseFloat(response.HighPrice, 64)
	low, _ := strconv.ParseFloat(response.LowPrice, 64)

	ticker := &models.Ticker{
		Symbol:         response.Symbol,
		Price:          price,
		Volume:         volume,
		PriceChange:    priceChange,
		PriceChangePct: priceChangePct,
		High24h:        high,
		Low24h:         low,
		Timestamp:      time.Now(),
	}

	// Cache the ticker for 5 seconds
	c.tickerCache.SetTicker(symbol, ticker, 5*time.Second)
	c.logger.Debug("Cached ticker", zap.String("symbol", symbol))

	return ticker, nil
}

// GetAllTickers fetches real-time price information for all trading pairs
func (c *Client) GetAllTickers(ctx context.Context) (map[string]*models.Ticker, error) {
	// Check cache first
	if tickers, found := c.tickerCache.GetAllTickers(); found {
		c.logger.Debug("Using cached tickers")
		return tickers, nil
	}

	var response []struct {
		Symbol             string `json:"symbol"`
		LastPrice          string `json:"lastPrice"`
		Volume             string `json:"volume"`
		PriceChange        string `json:"priceChange"`
		PriceChangePercent string `json:"priceChangePercent"`
		HighPrice          string `json:"highPrice"`
		LowPrice           string `json:"lowPrice"`
	}

	err := c.makeRequest(ctx, http.MethodGet, "/api/v3/ticker/24hr", nil, false, &response)
	if err != nil {
		return nil, err
	}

	tickers := make(map[string]*models.Ticker)
	now := time.Now()

	for _, item := range response {
		price, _ := strconv.ParseFloat(item.LastPrice, 64)
		volume, _ := strconv.ParseFloat(item.Volume, 64)
		priceChange, _ := strconv.ParseFloat(item.PriceChange, 64)
		priceChangePct, _ := strconv.ParseFloat(item.PriceChangePercent, 64)
		high, _ := strconv.ParseFloat(item.HighPrice, 64)
		low, _ := strconv.ParseFloat(item.LowPrice, 64)

		tickers[item.Symbol] = &models.Ticker{
			Symbol:         item.Symbol,
			Price:          price,
			Volume:         volume,
			PriceChange:    priceChange,
			PriceChangePct: priceChangePct,
			High24h:        high,
			Low24h:         low,
			Timestamp:      now,
		}
	}

	// Cache the tickers for 5 seconds
	c.tickerCache.SetAllTickers(tickers, 5*time.Second)
	c.logger.Debug("Cached all tickers", zap.Int("count", len(tickers)))

	return tickers, nil
}

// GetKlines fetches candlestick data for a specific trading pair and interval
func (c *Client) GetKlines(ctx context.Context, symbol, interval string, limit int) ([]*models.Kline, error) {
	// Check cache first
	if klines, found := c.klineCache.GetKlines(symbol, interval, limit); found {
		c.logger.Debug("Using cached klines",
			zap.String("symbol", symbol),
			zap.String("interval", interval),
			zap.Int("limit", limit))
		return klines, nil
	}

	params := url.Values{}
	params.Set("symbol", symbol)
	params.Set("interval", interval)
	if limit > 0 {
		params.Set("limit", strconv.Itoa(limit))
	}

	var response [][]interface{}

	err := c.makeRequest(ctx, http.MethodGet, "/api/v3/klines", params, false, &response)
	if err != nil {
		return nil, err
	}

	klines := make([]*models.Kline, 0, len(response))

	for _, k := range response {
		if len(k) < 7 {
			continue // Invalid kline data
		}

		// Parse timestamps
		openTime, err := parseTime(k[0])
		if err != nil {
			return nil, err
		}

		closeTime, err := parseTime(k[6])
		if err != nil {
			return nil, err
		}

		// Parse prices and volume
		open, _ := strconv.ParseFloat(k[1].(string), 64)
		high, _ := strconv.ParseFloat(k[2].(string), 64)
		low, _ := strconv.ParseFloat(k[3].(string), 64)
		close, _ := strconv.ParseFloat(k[4].(string), 64)
		volume, _ := strconv.ParseFloat(k[5].(string), 64)

		klines = append(klines, &models.Kline{
			Symbol:    symbol,
			Interval:  interval,
			OpenTime:  openTime,
			CloseTime: closeTime,
			Open:      open,
			High:      high,
			Low:       low,
			Close:     close,
			Volume:    volume,
		})
	}

	// Cache the klines for 30 seconds
	c.klineCache.SetKlines(symbol, interval, limit, klines, 30*time.Second)
	c.logger.Debug("Cached klines",
		zap.String("symbol", symbol),
		zap.String("interval", interval),
		zap.Int("count", len(klines)))

	return klines, nil
}

// ValidateKeys checks if the API keys are valid by making a simple authenticated request
func (c *Client) ValidateKeys(ctx context.Context) (bool, error) {
	c.logger.Debug("Validating API keys")

	// Try to get account information as a simple authenticated request
	var response struct{}
	err := c.makeRequest(ctx, http.MethodGet, "/api/v3/account", nil, true, &response)

	// If there's no error, the keys are valid
	if err == nil {
		c.logger.Debug("API keys are valid")
		return true, nil
	}

	// Check if the error is related to authentication
	if apiErr, ok := err.(*APIError); ok {
		c.logger.Error("API key validation failed",
			zap.Int("code", apiErr.Code),
			zap.String("message", apiErr.Message))

		// Check if it's an authentication error
		if apiErr.Code == -2015 || // Invalid API-key, IP, or permissions for action
			apiErr.Code == -2014 || // API-key format invalid
			apiErr.Code == -2013 { // Order does not exist (used when checking API key permissions)
			return false, nil
		}
	}

	// For other errors, return the error
	c.logger.Error("Error validating API keys", zap.Error(err))
	return false, err
}

// FetchBalances fetches all account balances and returns a structured Balance object
func (c *Client) FetchBalances(ctx context.Context) (models.Balance, error) {
	wallet, err := c.GetWallet(ctx)
	if err != nil {
		return models.Balance{}, err
	}

	balance := models.Balance{
		Available: make(map[string]float64),
		Locked:    make(map[string]float64),
	}

	// Extract USDT balance specifically for the Fiat field
	if usdtBalance, ok := wallet.Balances["USDT"]; ok {
		balance.Fiat = usdtBalance.Free
	}

	// Add all balances to the maps
	for asset, assetBalance := range wallet.Balances {
		balance.Available[asset] = assetBalance.Free
		balance.Locked[asset] = assetBalance.Locked
	}

	return balance, nil
}

// GetAccountBalance returns just the USDT balance as a float64
func (c *Client) GetAccountBalance(ctx context.Context) (float64, error) {
	balance, err := c.FetchBalances(ctx)
	if err != nil {
		return 0, err
	}
	return balance.Fiat, nil
}

// GetWallet fetches account balances
func (c *Client) GetWallet(ctx context.Context) (*models.Wallet, error) {
	c.logger.Debug("Fetching wallet from MEXC API")
	var response struct {
		Balances []struct {
			Asset  string `json:"asset"`
			Free   string `json:"free"`
			Locked string `json:"locked"`
		} `json:"balances"`
	}

	err := c.makeRequest(ctx, http.MethodGet, "/api/v3/account", nil, true, &response)
	if err != nil {
		c.logger.Error("Failed to fetch wallet from MEXC API", zap.Error(err))
		return nil, err
	}

	wallet := &models.Wallet{
		Balances:  make(map[string]*models.AssetBalance),
		UpdatedAt: time.Now(),
	}

	for _, balance := range response.Balances {
		free, _ := strconv.ParseFloat(balance.Free, 64)
		locked, _ := strconv.ParseFloat(balance.Locked, 64)

		// Only add assets with non-zero balances
		if free > 0 || locked > 0 {
			wallet.Balances[balance.Asset] = &models.AssetBalance{
				Asset:  balance.Asset,
				Free:   free,
				Locked: locked,
				Total:  free + locked,
			}
		}
	}

	c.logger.Debug("Successfully fetched wallet from MEXC API",
		zap.Int("asset_count", len(wallet.Balances)))
	return wallet, nil
}

// PlaceOrder places a new order
func (c *Client) PlaceOrder(ctx context.Context, order *models.Order) (*models.Order, error) {
	params := url.Values{}
	params.Set("symbol", order.Symbol)
	params.Set("side", parseOrderSide(order.Side))
	params.Set("type", parseOrderType(order.Type))

	// Set quantity
	params.Set("quantity", strconv.FormatFloat(order.Quantity, 'f', -1, 64))

	// Set price for limit orders
	if order.Type == models.OrderTypeLimit {
		params.Set("price", strconv.FormatFloat(order.Price, 'f', -1, 64))

		// Default timeInForce is GTC (Good Till Canceled)
		params.Set("timeInForce", "GTC")
	}

	// Set client order ID if provided
	if order.ClientID != "" {
		params.Set("newClientOrderId", order.ClientID)
	}

	var response struct {
		OrderID       string `json:"orderId"`
		ClientOrderID string `json:"clientOrderId"`
		Price         string `json:"price"`
		OrigQty       string `json:"origQty"`
		ExecutedQty   string `json:"executedQty"`
		Status        string `json:"status"`
		Type          string `json:"type"`
		Side          string `json:"side"`
		TransactTime  int64  `json:"transactTime"`
	}

	err := c.makeRequest(ctx, http.MethodPost, "/api/v3/order", params, true, &response)
	if err != nil {
		return nil, err
	}

	// Parse values
	price, _ := strconv.ParseFloat(response.Price, 64)
	origQty, _ := strconv.ParseFloat(response.OrigQty, 64)
	executedQty, _ := strconv.ParseFloat(response.ExecutedQty, 64)

	return &models.Order{
		ID:        response.OrderID,
		ClientID:  response.ClientOrderID,
		Symbol:    order.Symbol,
		Side:      models.OrderSide(response.Side),
		Type:      models.OrderType(response.Type),
		Quantity:  origQty,
		Price:     price,
		Status:    parseOrderStatusFromAPI(response.Status),
		CreatedAt: time.UnixMilli(response.TransactTime),
		UpdatedAt: time.UnixMilli(response.TransactTime),
		FilledQty: executedQty,
	}, nil
}

// CancelOrder cancels an existing order
func (c *Client) CancelOrder(ctx context.Context, orderID, symbol string) error {
	params := url.Values{}
	params.Set("symbol", symbol)
	params.Set("orderId", orderID)

	var response struct {
		OrderID string `json:"orderId"`
		Status  string `json:"status"`
	}

	err := c.makeRequest(ctx, http.MethodDelete, "/api/v3/order", params, true, &response)
	if err != nil {
		return err
	}

	if response.Status != "CANCELED" {
		return fmt.Errorf("failed to cancel order %s, status: %s", orderID, response.Status)
	}

	return nil
}

// GetOrder gets order details
func (c *Client) GetOrder(ctx context.Context, orderID, symbol string) (*models.Order, error) {
	params := url.Values{}
	params.Set("symbol", symbol)
	params.Set("orderId", orderID)

	var response struct {
		OrderID       string `json:"orderId"`
		ClientOrderID string `json:"clientOrderId"`
		Symbol        string `json:"symbol"`
		Price         string `json:"price"`
		OrigQty       string `json:"origQty"`
		ExecutedQty   string `json:"executedQty"`
		Status        string `json:"status"`
		Type          string `json:"type"`
		Side          string `json:"side"`
		Time          int64  `json:"time"`
		UpdateTime    int64  `json:"updateTime"`
	}

	err := c.makeRequest(ctx, http.MethodGet, "/api/v3/order", params, true, &response)
	if err != nil {
		return nil, err
	}

	// Parse values
	price, _ := strconv.ParseFloat(response.Price, 64)
	origQty, _ := strconv.ParseFloat(response.OrigQty, 64)
	executedQty, _ := strconv.ParseFloat(response.ExecutedQty, 64)

	return &models.Order{
		ID:        response.OrderID,
		ClientID:  response.ClientOrderID,
		Symbol:    response.Symbol,
		Side:      models.OrderSide(response.Side),
		Type:      models.OrderType(response.Type),
		Quantity:  origQty,
		Price:     price,
		Status:    parseOrderStatusFromAPI(response.Status),
		CreatedAt: time.UnixMilli(response.Time),
		UpdatedAt: time.UnixMilli(response.UpdateTime),
		FilledQty: executedQty,
	}, nil
}

// GetOpenOrders gets all open orders for a symbol
func (c *Client) GetOpenOrders(ctx context.Context, symbol string) ([]*models.Order, error) {
	params := url.Values{}
	if symbol != "" {
		params.Set("symbol", symbol)
	}

	var response []struct {
		OrderID       string `json:"orderId"`
		ClientOrderID string `json:"clientOrderId"`
		Symbol        string `json:"symbol"`
		Price         string `json:"price"`
		OrigQty       string `json:"origQty"`
		ExecutedQty   string `json:"executedQty"`
		Status        string `json:"status"`
		Type          string `json:"type"`
		Side          string `json:"side"`
		Time          int64  `json:"time"`
		UpdateTime    int64  `json:"updateTime"`
	}

	err := c.makeRequest(ctx, http.MethodGet, "/api/v3/openOrders", params, true, &response)
	if err != nil {
		return nil, err
	}

	orders := make([]*models.Order, 0, len(response))

	for _, o := range response {
		price, _ := strconv.ParseFloat(o.Price, 64)
		origQty, _ := strconv.ParseFloat(o.OrigQty, 64)
		executedQty, _ := strconv.ParseFloat(o.ExecutedQty, 64)

		orders = append(orders, &models.Order{
			ID:        o.OrderID,
			ClientID:  o.ClientOrderID,
			Symbol:    o.Symbol,
			Side:      models.OrderSide(o.Side),
			Type:      models.OrderType(o.Type),
			Quantity:  origQty,
			Price:     price,
			Status:    parseOrderStatusFromAPI(o.Status),
			CreatedAt: time.UnixMilli(o.Time),
			UpdatedAt: time.UnixMilli(o.UpdateTime),
			FilledQty: executedQty,
		})
	}

	return orders, nil
}

// GetNewCoins gets newly listed and upcoming coins
func (c *Client) GetNewCoins(ctx context.Context) ([]*models.NewCoin, error) {
	// Check cache first
	if newCoins, found := c.newCoinCache.GetNewCoins(); found {
		c.logger.Debug("Using cached new coins")
		return newCoins, nil
	}

	// Try to get real data from the MEXC calendar API first
	calendarCoins, calendarErr := c.GetNewCoinsFromCalendar(ctx)
	if calendarErr == nil && len(calendarCoins) > 0 {
		// Cache the new coins for 5 minutes
		c.newCoinCache.SetNewCoins(calendarCoins, 5*time.Minute)
		c.logger.Debug("Cached new coins from calendar API", zap.Int("count", len(calendarCoins)))
		return calendarCoins, nil
	} else if calendarErr != nil {
		c.logger.Warn("Failed to get new coins from calendar API, falling back to exchange info", zap.Error(calendarErr))
	}

	// Get exchange information with listing dates
	var exchangeInfo struct {
		Symbols []struct {
			Symbol        string `json:"symbol"`
			Status        string `json:"status"`
			BaseAsset     string `json:"baseAsset"`
			QuoteAsset    string `json:"quoteAsset"`
			FirstOpenTime int64  `json:"firstOpenTime,omitempty"`
			// Other fields...
		} `json:"symbols"`
	}

	exchangeErr := c.makeRequest(ctx, http.MethodGet, "/api/v3/exchangeInfo", nil, false, &exchangeInfo)
	if exchangeErr != nil {
		return nil, exchangeErr
	}

	// Get all tickers for volume information
	tickers, err := c.GetAllTickers(ctx)
	if err != nil {
		return nil, err
	}

	// Log some sample tickers for debugging
	tickerCount := 0
	for symbol, ticker := range tickers {
		if ticker.Volume > 0 {
			tickerCount++
			if tickerCount <= 5 {
				c.logger.Info("Sample ticker with volume",
					zap.String("symbol", symbol),
					zap.Float64("volume", ticker.Volume),
					zap.Float64("price", ticker.Price))
			}
		}
	}

	c.logger.Info("Retrieved exchange info and tickers",
		zap.Int("symbols_count", len(exchangeInfo.Symbols)),
		zap.Int("tickers_count", len(tickers)),
		zap.Int("tickers_with_volume", tickerCount))

	// Initialize empty slice for new coins
	newCoins := make([]*models.NewCoin, 0)

	// Get current time for determining upcoming listings
	currentTime := time.Now()

	// Process exchange info to find new and upcoming coins
	for _, symbol := range exchangeInfo.Symbols {
		// Only consider USDT pairs
		if symbol.QuoteAsset != "USDT" {
			continue
		}

		// Get volume information from ticker if available
		volume := 0.0
		ticker, exists := tickers[symbol.Symbol]
		if exists {
			volume = ticker.Volume
		}

		// Determine if this is an upcoming listing
		var firstOpenTime time.Time
		isUpcoming := false

		if symbol.FirstOpenTime > 0 {
			// Convert milliseconds to time.Time
			firstOpenTime = time.UnixMilli(symbol.FirstOpenTime)

			// If firstOpenTime is in the future, it's an upcoming listing
			isUpcoming = firstOpenTime.After(currentTime)

			c.logger.Info("Found coin with firstOpenTime",
				zap.String("symbol", symbol.Symbol),
				zap.Time("firstOpenTime", firstOpenTime),
				zap.Bool("isUpcoming", isUpcoming))
		} else {
			// If no firstOpenTime is provided, use current time
			firstOpenTime = currentTime
		}

		// Create new coin
		newCoins = append(newCoins, &models.NewCoin{
			Symbol:        symbol.Symbol,
			BaseVolume:    0, // Not available directly
			QuoteVolume:   volume,
			FoundAt:       currentTime,
			FirstOpenTime: &firstOpenTime,
			Status:        symbol.Status, // Store the trading status
			IsUpcoming:    isUpcoming,
		})
	}

	// Sort by firstOpenTime (ascending) to prioritize upcoming coins first
	sort.Slice(newCoins, func(i, j int) bool {
		// Handle nil pointers
		if newCoins[i].FirstOpenTime == nil {
			return false
		}
		if newCoins[j].FirstOpenTime == nil {
			return true
		}
		return newCoins[i].FirstOpenTime.Before(*newCoins[j].FirstOpenTime)
	})

	c.logger.Info("Filtered new coins",
		zap.Int("count", len(newCoins)),
		zap.Int("upcoming_count", countUpcomingCoins(newCoins)))

	// Cache the new coins for 5 minutes
	c.newCoinCache.SetNewCoins(newCoins, 5*time.Minute)
	c.logger.Debug("Cached new coins", zap.Int("count", len(newCoins)))

	return newCoins, nil
}

// Helper function to check if a symbol already exists in the coins slice
func containsSymbol(coins []*models.NewCoin, symbol string) bool {
	for _, coin := range coins {
		if coin.Symbol == symbol {
			return true
		}
	}
	return false
}

// GetNewCoinsFromCalendar gets upcoming coins from the MEXC calendar API
func (c *Client) GetNewCoinsFromCalendar(ctx context.Context) ([]*models.NewCoin, error) {
	// Initialize empty slice for new coins
	newCoins := make([]*models.NewCoin, 0)

	// Get current time for determining upcoming listings
	currentTime := time.Now()

	// Define the calendar API URL
	calendarURL := "https://www.mexc.com/api/operation/new_coin_calendar"

	// Add timestamp parameter (current time in milliseconds)
	params := url.Values{}
	params.Set("timestamp", strconv.FormatInt(currentTime.UnixMilli(), 10))

	// Create a new HTTP request
	reqURL := fmt.Sprintf("%s?%s", calendarURL, params.Encode())
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, &RequestError{Err: err, Message: "failed to create calendar API request"}
	}

	// Set common headers
	req.Header.Set("Accept", "application/json")

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, &RequestError{Err: err, Message: "failed to execute calendar API request"}
	}
	defer resp.Body.Close()

	// Check response status code
	if resp.StatusCode != http.StatusOK {
		return nil, &APIError{
			Code:    resp.StatusCode,
			Message: fmt.Sprintf("calendar API returned non-200 status code: %d", resp.StatusCode),
		}
	}

	// Parse response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &RequestError{Err: err, Message: "failed to read calendar API response body"}
	}

	// Parse JSON response
	var response struct {
		Code int `json:"code"`
		Data struct {
			List []struct {
				Symbol        string `json:"symbol"`
				Name          string `json:"name"`
				FirstOpenTime int64  `json:"firstOpenTime"`
			} `json:"list"`
		} `json:"data"`
	}

	// Log the raw response for debugging
	c.logger.Debug("Raw calendar API response", zap.String("body", string(respBody)))

	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return nil, &RequestError{Err: err, Message: "failed to parse calendar API response"}
	}

	// Check response code
	if response.Code != 200 {
		return nil, &APIError{
			Code:    response.Code,
			Message: fmt.Sprintf("calendar API returned error code: %d", response.Code),
		}
	}

	// Process response data
	for _, coin := range response.Data.List {
		// Convert firstOpenTime from milliseconds to time.Time
		firstOpenTime := time.UnixMilli(coin.FirstOpenTime)

		// Determine if this is an upcoming listing
		isUpcoming := firstOpenTime.After(currentTime)

		// Create new coin
		newCoins = append(newCoins, &models.NewCoin{
			Symbol:        coin.Symbol,
			BaseVolume:    0, // Not available directly
			QuoteVolume:   0, // Not available directly
			FoundAt:       currentTime,
			FirstOpenTime: &firstOpenTime,
			Status:        "", // Status not available from calendar API, will be updated later
			IsUpcoming:    isUpcoming,
		})
	}

	// Sort by firstOpenTime (ascending) to prioritize upcoming coins first
	sort.Slice(newCoins, func(i, j int) bool {
		// Handle nil pointers
		if newCoins[i].FirstOpenTime == nil {
			return false
		}
		if newCoins[j].FirstOpenTime == nil {
			return true
		}
		return newCoins[i].FirstOpenTime.Before(*newCoins[j].FirstOpenTime)
	})

	c.logger.Info("Got new coins from calendar API",
		zap.Int("count", len(newCoins)))

	return newCoins, nil
}

// Helper function to check if there are any upcoming coins in the slice
func containsUpcomingCoins(coins []*models.NewCoin) bool {
	for _, coin := range coins {
		if coin.IsUpcoming {
			return true
		}
	}
	return false
}

// Helper function to count upcoming coins in the slice
func countUpcomingCoins(coins []*models.NewCoin) int {
	count := 0
	for _, coin := range coins {
		if coin.IsUpcoming {
			count++
		}
	}
	return count
}

func (c *Client) GetOrderBook(ctx context.Context, symbol string, limit int) (*models.OrderBookUpdate, error) {
	// Check cache first
	if orderBook, found := c.orderBookCache.GetOrderBook(symbol); found {
		c.logger.Debug("Using cached order book", zap.String("symbol", symbol))
		return orderBook, nil
	}

	params := url.Values{}
	params.Set("symbol", symbol)
	if limit > 0 {
		params.Set("limit", strconv.Itoa(limit))
	}

	var response struct {
		LastUpdateID int64      `json:"lastUpdateId"`
		Bids         [][]string `json:"bids"`
		Asks         [][]string `json:"asks"`
	}

	err := c.makeRequest(ctx, http.MethodGet, "/api/v3/depth", params, false, &response)
	if err != nil {
		return nil, err
	}

	parseEntries := func(raw [][]string) []models.OrderBookEntry {
		entries := make([]models.OrderBookEntry, 0, len(raw))
		for _, item := range raw {
			if len(item) != 2 {
				continue
			}
			price, err1 := strconv.ParseFloat(item[0], 64)
			qty, err2 := strconv.ParseFloat(item[1], 64)
			if err1 != nil || err2 != nil {
				continue
			}
			entries = append(entries, models.OrderBookEntry{
				Price:    price,
				Quantity: qty,
			})
		}
		return entries
	}

	orderBook := &models.OrderBookUpdate{
		Symbol:        symbol,
		LastUpdateID:  response.LastUpdateID,
		FirstUpdateID: response.LastUpdateID, // REST snapshot, so first == last
		Bids:          parseEntries(response.Bids),
		Asks:          parseEntries(response.Asks),
		Timestamp:     time.Now(),
	}

	// Cache the order book for 2 seconds
	c.orderBookCache.SetOrderBook(symbol, orderBook, 2*time.Second)
	c.logger.Debug("Cached order book",
		zap.String("symbol", symbol),
		zap.Int("bids", len(orderBook.Bids)),
		zap.Int("asks", len(orderBook.Asks)))

	return orderBook, nil
}
