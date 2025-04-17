package mexc

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/rs/zerolog"
)

const (
	baseURL = "https://api.mexc.com"
)

// Client implements port.MEXCClient interface
// Note: MEXC API requires the APIKEY header (not X-MBX-APIKEY) for authentication
type Client struct {
	httpClient *http.Client
	apiKey     string
	apiSecret  string
	logger     *zerolog.Logger
}

// NewClient creates a new MEXC API client
func NewClient(apiKey, apiSecret string, logger *zerolog.Logger) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		apiKey:    apiKey,
		apiSecret: apiSecret,
		logger:    logger,
	}
}

// GetNewListings retrieves information about newly listed coins
func (c *Client) GetNewListings(ctx context.Context) ([]*model.NewCoin, error) {
	endpoint := "/api/v3/ticker/new"

	resp, err := c.sendRequest(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get new listings: %w", err)
	}
	defer resp.Body.Close()

	var response struct {
		Symbols []struct {
			Symbol      string `json:"symbol"`
			BaseAsset   string `json:"baseAsset"`
			QuoteAsset  string `json:"quoteAsset"`
			ListingTime int64  `json:"listingTime"`
			TradingTime int64  `json:"tradingTime"`
			Status      string `json:"status"`
			MinPrice    string `json:"minPrice"`
			MaxPrice    string `json:"maxPrice"`
			MinQty      string `json:"minQty"`
			MaxQty      string `json:"maxQty"`
			PriceScale  int    `json:"priceScale"`
			QtyScale    int    `json:"qtyScale"`
		} `json:"symbols"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	var coins []*model.NewCoin
	for _, s := range response.Symbols {
		minPrice, _ := strconv.ParseFloat(s.MinPrice, 64)
		maxPrice, _ := strconv.ParseFloat(s.MaxPrice, 64)
		minQty, _ := strconv.ParseFloat(s.MinQty, 64)
		maxQty, _ := strconv.ParseFloat(s.MaxQty, 64)

   var status model.CoinStatus
   switch s.Status {
   case "PRE_TRADING":
       status = model.CoinStatusExpected
   case "TRADING":
       status = model.CoinStatusTrading
   case "BREAK":
       status = model.CoinStatusFailed
   default:
       status = model.CoinStatusListed
		}

		coin := &model.NewCoin{
			Symbol:              s.Symbol,
			BaseAsset:           s.BaseAsset,
			QuoteAsset:          s.QuoteAsset,
			Status:              status,
			ExpectedListingTime: time.Unix(s.ListingTime/1000, 0),
			MinPrice:            minPrice,
			MaxPrice:            maxPrice,
			MinQty:              minQty,
			MaxQty:              maxQty,
			PriceScale:          s.PriceScale,
			QtyScale:            s.QtyScale,
			CreatedAt:           time.Now(),
			UpdatedAt:           time.Now(),
		}

		if s.TradingTime > 0 {
			tradingTime := time.Unix(s.TradingTime/1000, 0)
			coin.BecameTradableAt = &tradingTime
		}

		coins = append(coins, coin)
	}

	return coins, nil
}

// GetSymbolInfo retrieves detailed information about a trading symbol
func (c *Client) GetSymbolInfo(ctx context.Context, symbol string) (*model.SymbolInfo, error) { // Changed return type
	endpoint := fmt.Sprintf("/api/v3/exchangeInfo?symbol=%s", symbol)

	resp, err := c.sendRequest(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get symbol info: %w", err)
	}
	defer resp.Body.Close()

	// Define a struct that matches the expected API response structure for exchangeInfo
	var response struct {
		Symbols []struct {
			Symbol               string                   `json:"symbol"`
			Status               string                   `json:"status"`
			BaseAsset            string                   `json:"baseAsset"`
			BaseAssetPrecision   int                      `json:"baseAssetPrecision"`
			QuoteAsset           string                   `json:"quoteAsset"`
			QuotePrecision       int                      `json:"quotePrecision"` // API uses quotePrecision
			OrderTypes           []string                 `json:"orderTypes"`
			IsSpotTradingAllowed bool                     `json:"isSpotTradingAllowed"`
			Permissions          []string                 `json:"permissions"`
			Filters              []map[string]interface{} `json:"filters"` // Filters are complex, parse later
		} `json:"symbols"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode exchangeInfo response: %w", err)
	}

	if len(response.Symbols) == 0 {
		return nil, fmt.Errorf("symbol %s not found in exchangeInfo", symbol)
	}

	s := response.Symbols[0]

	// Extract filter values (needs robust parsing based on filter type)
	var minNotional, minLotSize, maxLotSize, stepSize, tickSize string
	for _, filter := range s.Filters {
		filterType := filter["filterType"].(string)
		switch filterType {
		case "PRICE_FILTER":
			tickSize = filter["tickSize"].(string)
		case "LOT_SIZE":
			minLotSize = filter["minQty"].(string)
			maxLotSize = filter["maxQty"].(string)
			stepSize = filter["stepSize"].(string)
		case "MIN_NOTIONAL":
			minNotional = filter["minNotional"].(string)
			// Add other filter types as needed
		}
	}

	return &model.SymbolInfo{
		Symbol:               s.Symbol,
		Status:               s.Status,
		BaseAsset:            s.BaseAsset,
		BaseAssetPrecision:   s.BaseAssetPrecision,
		QuoteAsset:           s.QuoteAsset,
		QuoteAssetPrecision:  s.QuotePrecision,
		OrderTypes:           s.OrderTypes,
		IsSpotTradingAllowed: s.IsSpotTradingAllowed,
		Permissions:          s.Permissions,
		MinNotional:          minNotional,
		MinLotSize:           minLotSize,
		MaxLotSize:           maxLotSize,
		StepSize:             stepSize,
		TickSize:             tickSize,
	}, nil
}

// GetSymbolStatus checks if a symbol is currently tradeable
func (c *Client) GetSymbolStatus(ctx context.Context, symbol string) (model.CoinStatus, error) {
	symbolInfo, err := c.GetSymbolInfo(ctx, symbol)
	if err != nil {
		return model.CoinStatusFailed, err
	}

	// Convert string status to model.Status
	switch symbolInfo.Status {
   case "TRADING":
       return model.CoinStatusTrading, nil
   case "PRE_TRADING":
       return model.CoinStatusExpected, nil
   case "BREAK":
       return model.CoinStatusFailed, nil
   case "HALT":
       return model.CoinStatusListed, nil // Using CoinStatusListed for suspended
   default:
       return model.CoinStatusListed, nil
	}
}

// GetTradingSchedule retrieves the listing and trading schedule for a symbol
func (c *Client) GetTradingSchedule(ctx context.Context, symbol string) (model.TradingSchedule, error) {
	endpoint := fmt.Sprintf("/api/v3/ticker/new?symbol=%s", symbol)

	resp, err := c.sendRequest(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return model.TradingSchedule{}, fmt.Errorf("failed to get trading schedule: %w", err)
	}
	defer resp.Body.Close()

	var response struct {
		ListingTime int64 `json:"listingTime"`
		TradingTime int64 `json:"tradingTime"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return model.TradingSchedule{}, fmt.Errorf("failed to decode response: %w", err)
	}

	schedule := model.TradingSchedule{
		ListingTime: time.Unix(response.ListingTime/1000, 0),
	}

	if response.TradingTime > 0 {
		schedule.TradingTime = time.Unix(response.TradingTime/1000, 0)
	}

	return schedule, nil
}

// GetSymbolConstraints retrieves trading constraints for a symbol
func (c *Client) GetSymbolConstraints(ctx context.Context, symbol string) (*model.SymbolConstraints, error) {
	symbolInfo, err := c.GetSymbolInfo(ctx, symbol)
	if err != nil {
		return nil, err
	}

	// Parse string values from SymbolInfo filters to the required types
	// Default values if parsing fails or fields are empty
	minPrice := 0.0 // Min price is often not directly available, TickSize is more relevant
	maxPrice := 0.0 // Max price is often not directly available
	minQty, _ := strconv.ParseFloat(symbolInfo.MinLotSize, 64)
	maxQty, _ := strconv.ParseFloat(symbolInfo.MaxLotSize, 64)
	priceScale := symbolInfo.QuoteAssetPrecision // Use QuoteAssetPrecision for price scale
	qtyScale := symbolInfo.BaseAssetPrecision    // Use BaseAssetPrecision for quantity scale

	// Note: TickSize and StepSize might be more useful than min/max price/qty for validation
	// tickSize, _ := strconv.ParseFloat(symbolInfo.TickSize, 64)
	// stepSize, _ := strconv.ParseFloat(symbolInfo.StepSize, 64)

	return &model.SymbolConstraints{
		MinPrice:   minPrice,
		MaxPrice:   maxPrice,
		MinQty:     minQty,
		MaxQty:     maxQty,
		PriceScale: priceScale,
		QtyScale:   qtyScale,
	}, nil
}

// sendRequest sends an HTTP request to the MEXC API
func (c *Client) sendRequest(ctx context.Context, method, endpoint string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, baseURL+endpoint, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		// MEXC API requires the APIKEY header, not X-MBX-APIKEY
		req.Header.Set("APIKEY", c.apiKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		var errResp struct {
			Code    int    `json:"code"`
			Message string `json:"msg"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return nil, fmt.Errorf("request failed with status %d", resp.StatusCode)
		}
		return nil, fmt.Errorf("API error %d: %s", errResp.Code, errResp.Message)
	}

	return resp, nil
}

// GetAccount retrieves account information from MEXC
func (c *Client) GetAccount(ctx context.Context) (*model.Wallet, error) {
	c.logger.Debug().Msg("Fetching account information from MEXC")

	// Create timestamp for the request
	timestamp := time.Now().UnixMilli()

	// Create query parameters
	params := fmt.Sprintf("timestamp=%d", timestamp)

	// Generate signature
	signature := c.generateSignature(params)

	// Add signature to parameters
	endpoint := fmt.Sprintf("/api/v3/account?%s&signature=%s", params, signature)

	// Create request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add API key header
	c.logger.Debug().Str("APIKEY", c.apiKey).Msg("Setting API key header")
	req.Header.Set("APIKEY", c.apiKey)

	// Send request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Error().Err(err).Msg("Failed to send request to MEXC API")
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		var errResp struct {
			Code    int    `json:"code"`
			Message string `json:"msg"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			c.logger.Error().Err(err).Int("status", resp.StatusCode).Msg("Failed to decode error response")
			return nil, fmt.Errorf("request failed with status %d", resp.StatusCode)
		}
		c.logger.Error().Int("code", errResp.Code).Str("message", errResp.Message).Msg("MEXC API error")
		return nil, fmt.Errorf("API error %d: %s", errResp.Code, errResp.Message)
	}

	// Parse response
	var accountInfo struct {
		MakerCommission  int    `json:"makerCommission"`
		TakerCommission  int    `json:"takerCommission"`
		BuyerCommission  int    `json:"buyerCommission"`
		SellerCommission int    `json:"sellerCommission"`
		CanTrade         bool   `json:"canTrade"`
		CanWithdraw      bool   `json:"canWithdraw"`
		CanDeposit       bool   `json:"canDeposit"`
		UpdateTime       int64  `json:"updateTime"`
		AccountType      string `json:"accountType"`
		Balances         []struct {
			Asset  string `json:"asset"`
			Free   string `json:"free"`
			Locked string `json:"locked"`
		} `json:"balances"`
		Permissions []string `json:"permissions"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&accountInfo); err != nil {
		c.logger.Error().Err(err).Msg("Failed to decode account response")
		return nil, fmt.Errorf("failed to decode account response: %w", err)
	}

	// Convert to model.Wallet
	wallet := &model.Wallet{
		UserID:      "MEXC_USER", // Default user ID
		Exchange:    "MEXC",
		Balances:    make(map[model.Asset]*model.Balance),
		LastUpdated: time.Now(),
		LastSyncAt:  time.Now(),
	}

	// Add balances
	for _, b := range accountInfo.Balances {
		// Skip zero balances
		if b.Free == "0" && b.Locked == "0" {
			continue
		}

		free, _ := strconv.ParseFloat(b.Free, 64)
		locked, _ := strconv.ParseFloat(b.Locked, 64)
		total := free + locked

		wallet.Balances[model.Asset(b.Asset)] = &model.Balance{
			Asset:  model.Asset(b.Asset),
			Free:   free,
			Locked: locked,
			Total:  total,
			// Note: USDValue will be calculated later when we have price data
		}
	}

	c.logger.Info().Int("balances_count", len(wallet.Balances)).Msg("Successfully fetched account information")
	return wallet, nil
}

// generateSignature generates the HMAC SHA256 signature for authenticated requests
func (c *Client) generateSignature(data string) string {
	h := hmac.New(sha256.New, []byte(c.apiSecret))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

// CancelOrder cancels an existing order
func (c *Client) CancelOrder(ctx context.Context, symbol string, orderID string) error {
	// For now, this is a stub implementation
	c.logger.Warn().Str("symbol", symbol).Str("orderID", orderID).Msg("CancelOrder not fully implemented")
	return fmt.Errorf("CancelOrder method not fully implemented")
}

// PlaceOrder places a new order
func (c *Client) PlaceOrder(ctx context.Context, symbol string, side model.OrderSide, orderType model.OrderType, quantity float64, price float64, timeInForce model.TimeInForce) (*model.Order, error) {
	// For now, this is a stub implementation
	c.logger.Warn().Str("symbol", symbol).Str("side", string(side)).Msg("PlaceOrder not fully implemented")
	return nil, fmt.Errorf("PlaceOrder method not fully implemented")
}

// GetOrderStatus retrieves the status of an order
func (c *Client) GetOrderStatus(ctx context.Context, symbol string, orderID string) (*model.Order, error) {
	// For now, this is a stub implementation
	c.logger.Warn().Str("symbol", symbol).Str("orderID", orderID).Msg("GetOrderStatus not fully implemented")
	return nil, fmt.Errorf("GetOrderStatus method not fully implemented")
}

// GetOpenOrders retrieves all open orders for a symbol
func (c *Client) GetOpenOrders(ctx context.Context, symbol string) ([]*model.Order, error) {
	// For now, this is a stub implementation
	c.logger.Warn().Str("symbol", symbol).Msg("GetOpenOrders not fully implemented")
	return nil, fmt.Errorf("GetOpenOrders method not fully implemented")
}

// GetOrderHistory retrieves historical orders for a symbol
func (c *Client) GetOrderHistory(ctx context.Context, symbol string, limit, offset int) ([]*model.Order, error) {
	// For now, this is a stub implementation
	c.logger.Warn().Str("symbol", symbol).Int("limit", limit).Int("offset", offset).Msg("GetOrderHistory not fully implemented")
	return nil, fmt.Errorf("GetOrderHistory method not fully implemented")
}

// GetExchangeInfo retrieves information about all symbols on the exchange
func (c *Client) GetExchangeInfo(ctx context.Context) (*model.ExchangeInfo, error) {
	endpoint := "/api/v3/exchangeInfo"

	resp, err := c.sendRequest(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get exchange info: %w", err)
	}
	defer resp.Body.Close()

	var response struct {
		Timezone   string `json:"timezone"`
		ServerTime int64  `json:"serverTime"`
		Symbols    []struct {
			Symbol               string                   `json:"symbol"`
			Status               string                   `json:"status"`
			BaseAsset            string                   `json:"baseAsset"`
			BaseAssetPrecision   int                      `json:"baseAssetPrecision"`
			QuoteAsset           string                   `json:"quoteAsset"`
			QuotePrecision       int                      `json:"quotePrecision"`
			OrderTypes           []string                 `json:"orderTypes"`
			IsSpotTradingAllowed bool                     `json:"isSpotTradingAllowed"`
			Permissions          []string                 `json:"permissions"`
			Filters              []map[string]interface{} `json:"filters"`
		} `json:"symbols"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Convert to model.ExchangeInfo
	exchangeInfo := &model.ExchangeInfo{
		Symbols: make([]model.SymbolInfo, len(response.Symbols)),
	}

	for i, s := range response.Symbols {
		// Extract filter values
		var minNotional, minLotSize, maxLotSize, stepSize, tickSize string
		for _, filter := range s.Filters {
			filterType := filter["filterType"].(string)
			switch filterType {
			case "PRICE_FILTER":
				tickSize = filter["tickSize"].(string)
			case "LOT_SIZE":
				minLotSize = filter["minQty"].(string)
				maxLotSize = filter["maxQty"].(string)
				stepSize = filter["stepSize"].(string)
			case "MIN_NOTIONAL":
				minNotional = filter["minNotional"].(string)
			}
		}

		exchangeInfo.Symbols[i] = model.SymbolInfo{
			Symbol:               s.Symbol,
			Status:               s.Status,
			BaseAsset:            s.BaseAsset,
			BaseAssetPrecision:   s.BaseAssetPrecision,
			QuoteAsset:           s.QuoteAsset,
			QuoteAssetPrecision:  s.QuotePrecision,
			OrderTypes:           s.OrderTypes,
			IsSpotTradingAllowed: s.IsSpotTradingAllowed,
			Permissions:          s.Permissions,
			MinNotional:          minNotional,
			MinLotSize:           minLotSize,
			MaxLotSize:           maxLotSize,
			StepSize:             stepSize,
			TickSize:             tickSize,
		}
	}

	return exchangeInfo, nil
}

// GetKlines retrieves candle data for a symbol, interval, and limit
func (c *Client) GetKlines(ctx context.Context, symbol string, interval model.KlineInterval, limit int) ([]*model.Kline, error) {
	endpoint := fmt.Sprintf("/api/v3/klines?symbol=%s&interval=%s&limit=%d", symbol, interval, limit)

	resp, err := c.sendRequest(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get klines: %w", err)
	}
	defer resp.Body.Close()

	var rawKlines [][]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&rawKlines); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	klines := make([]*model.Kline, 0, len(rawKlines))
	for _, raw := range rawKlines {
		if len(raw) < 11 {
			c.logger.Warn().
				Str("symbol", symbol).
				Str("interval", string(interval)).
				Interface("data", raw).
				Msg("Incomplete kline data received from API")
			continue // Skip incomplete data
		}

		// Parse the timestamp
		openTime, err := strconv.ParseInt(fmt.Sprintf("%v", raw[0]), 10, 64)
		if err != nil {
			c.logger.Warn().Err(err).Msg("Failed to parse kline open time")
			continue
		}

		closeTime, err := strconv.ParseInt(fmt.Sprintf("%v", raw[6]), 10, 64)
		if err != nil {
			c.logger.Warn().Err(err).Msg("Failed to parse kline close time")
			continue
		}

		// Parse price and volume values
		open, _ := strconv.ParseFloat(fmt.Sprintf("%v", raw[1]), 64)
		high, _ := strconv.ParseFloat(fmt.Sprintf("%v", raw[2]), 64)
		low, _ := strconv.ParseFloat(fmt.Sprintf("%v", raw[3]), 64)
		close, _ := strconv.ParseFloat(fmt.Sprintf("%v", raw[4]), 64)
		volume, _ := strconv.ParseFloat(fmt.Sprintf("%v", raw[5]), 64)

		kline := &model.Kline{
			Symbol:    symbol,
			Interval:  interval,
			OpenTime:  time.Unix(openTime/1000, 0),
			CloseTime: time.Unix(closeTime/1000, 0),
			Open:      open,
			High:      high,
			Low:       low,
			Close:     close,
			Volume:    volume,
		}

		klines = append(klines, kline)
	}

	return klines, nil
}

// GetMarketData retrieves current market data for a symbol
func (c *Client) GetMarketData(ctx context.Context, symbol string) (*model.Ticker, error) {
	endpoint := fmt.Sprintf("/api/v3/ticker/24hr?symbol=%s", symbol)

	resp, err := c.sendRequest(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get market data: %w", err)
	}
	defer resp.Body.Close()

	var response struct {
		Symbol             string `json:"symbol"`
		PriceChange        string `json:"priceChange"`
		PriceChangePercent string `json:"priceChangePercent"`
		WeightedAvgPrice   string `json:"weightedAvgPrice"`
		PrevClosePrice     string `json:"prevClosePrice"`
		LastPrice          string `json:"lastPrice"`
		LastQty            string `json:"lastQty"`
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

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Parse string values to float64
	lastPrice, _ := strconv.ParseFloat(response.LastPrice, 64)
	volume, _ := strconv.ParseFloat(response.Volume, 64)
	highPrice, _ := strconv.ParseFloat(response.HighPrice, 64)
	lowPrice, _ := strconv.ParseFloat(response.LowPrice, 64)
	priceChange, _ := strconv.ParseFloat(response.PriceChange, 64)
	priceChangePercent, _ := strconv.ParseFloat(response.PriceChangePercent, 64)

	ticker := &model.Ticker{
		Symbol:             response.Symbol,
		LastPrice:          lastPrice,
		Volume:             volume,
		HighPrice:          highPrice,
		LowPrice:           lowPrice,
		PriceChange:        priceChange,
		PriceChangePercent: priceChangePercent,
	}

	return ticker, nil
}

// GetOrderBook retrieves the order book for a symbol
func (c *Client) GetOrderBook(ctx context.Context, symbol string, depth int) (*model.OrderBook, error) {
	// Validate depth parameter (MEXC supports 5, 10, 20, 50, 100, 500, 1000)
	validDepths := []int{5, 10, 20, 50, 100, 500, 1000}
	isValidDepth := false
	for _, d := range validDepths {
		if depth == d {
			isValidDepth = true
			break
		}
	}
	if !isValidDepth {
		depth = 10 // Default to 10 if invalid
		c.logger.Warn().Int("depth", depth).Msg("Using default depth for order book")
	}

	endpoint := fmt.Sprintf("/api/v3/depth?symbol=%s&limit=%d", symbol, depth)

	resp, err := c.sendRequest(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get order book: %w", err)
	}
	defer resp.Body.Close()

	var response struct {
		LastUpdateID int64      `json:"lastUpdateId"`
		Bids         [][]string `json:"bids"` // [price, quantity]
		Asks         [][]string `json:"asks"` // [price, quantity]
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Convert to model.OrderBook
	orderBook := &model.OrderBook{
		Symbol:       symbol,
		LastUpdateID: response.LastUpdateID,
		Bids:         make([]model.OrderBookEntry, len(response.Bids)),
		Asks:         make([]model.OrderBookEntry, len(response.Asks)),
	}

	// Parse bids
	for i, bid := range response.Bids {
		if len(bid) < 2 {
			continue // Skip invalid entries
		}
		price, _ := strconv.ParseFloat(bid[0], 64)
		quantity, _ := strconv.ParseFloat(bid[1], 64)
		orderBook.Bids[i] = model.OrderBookEntry{
			Price:    price,
			Quantity: quantity,
		}
	}

	// Parse asks
	for i, ask := range response.Asks {
		if len(ask) < 2 {
			continue // Skip invalid entries
		}
		price, _ := strconv.ParseFloat(ask[0], 64)
		quantity, _ := strconv.ParseFloat(ask[1], 64)
		orderBook.Asks[i] = model.OrderBookEntry{
			Price:    price,
			Quantity: quantity,
		}
	}

	return orderBook, nil
}
