package rest

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/neo/crypto-bot/internal/domain/model"
)

// Response types for MEXC API responses
type (
	// AccountResponse represents the response for account information
	AccountResponse struct {
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
	}

	// TickerResponse represents the ticker information
	TickerResponse struct {
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

	// KlineResponse represents a single kline/candlestick
	KlineResponse []interface{}

	// OrderBookResponse represents the order book
	OrderBookResponse struct {
		LastUpdateID int64      `json:"lastUpdateId"`
		Bids         [][]string `json:"bids"` // [price, quantity]
		Asks         [][]string `json:"asks"` // [price, quantity]
	}

	// OrderResponse represents the response for placing an order
	OrderResponse struct {
		Symbol              string      `json:"symbol"`
		OrderID             interface{} `json:"orderId"`     // Can be either string or int64
		OrderListID         interface{} `json:"orderListId"` // Can be either string or int64
		ClientOrderID       string      `json:"clientOrderId"`
		TransactTime        int64       `json:"transactTime"`
		Price               string      `json:"price"`
		OrigQty             string      `json:"origQty"`
		ExecutedQty         string      `json:"executedQty"`
		CummulativeQuoteQty string      `json:"cummulativeQuoteQty"`
		Status              string      `json:"status"`
		TimeInForce         string      `json:"timeInForce"`
		Type                string      `json:"type"`
		Side                string      `json:"side"`
		Fills               []struct {
			Price           string `json:"price"`
			Qty             string `json:"qty"`
			Commission      string `json:"commission"`
			CommissionAsset string `json:"commissionAsset"`
		} `json:"fills"`
	}
)

// GetAccount fetches account information
func (c *Client) GetAccount(ctx context.Context) (*model.Wallet, error) {
	resp, err := c.callPrivateAPI(ctx, http.MethodGet, "/account", nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	var accountResp AccountResponse
	if err := json.Unmarshal(resp, &accountResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal account response: %w", err)
	}

	// Convert response to domain model
	wallet := &model.Wallet{
		Balances:    make(map[model.Asset]*model.Balance),
		LastUpdated: time.Now(),
	}

	for _, balance := range accountResp.Balances {
		free, err := strconv.ParseFloat(balance.Free, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse free balance: %w", err)
		}

		locked, err := strconv.ParseFloat(balance.Locked, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse locked balance: %w", err)
		}

		// Only include assets with non-zero balance
		if free > 0 || locked > 0 {
			asset := model.Asset(balance.Asset)
			wallet.Balances[asset] = &model.Balance{
				Asset:    asset,
				Free:     free,
				Locked:   locked,
				Total:    free + locked,
				USDValue: 0, // We don't have USD value from MEXC API directly
			}
		}
	}

	return wallet, nil
}

// GetMarketData fetches market data for a symbol
func (c *Client) GetMarketData(ctx context.Context, symbol string) (*model.Ticker, error) {
	params := map[string]string{
		"symbol": symbol,
	}

	resp, err := c.callPublicAPI(ctx, http.MethodGet, "/ticker/24hr", params)
	if err != nil {
		return nil, fmt.Errorf("failed to get market data: %w", err)
	}

	var tickerResp TickerResponse
	if err := json.Unmarshal(resp, &tickerResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal ticker response: %w", err)
	}

	// Convert response to domain model
	lastPrice, err := strconv.ParseFloat(tickerResp.LastPrice, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse last price: %w", err)
	}

	bidPrice, err := strconv.ParseFloat(tickerResp.BidPrice, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse bid price: %w", err)
	}

	askPrice, err := strconv.ParseFloat(tickerResp.AskPrice, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ask price: %w", err)
	}

	highPrice, err := strconv.ParseFloat(tickerResp.HighPrice, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse high price: %w", err)
	}

	lowPrice, err := strconv.ParseFloat(tickerResp.LowPrice, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse low price: %w", err)
	}

	volume, err := strconv.ParseFloat(tickerResp.Volume, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse volume: %w", err)
	}

	quoteVolume, err := strconv.ParseFloat(tickerResp.QuoteVolume, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse quote volume: %w", err)
	}

	priceChange, err := strconv.ParseFloat(tickerResp.PriceChange, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse price change: %w", err)
	}

	priceChangePercent, err := strconv.ParseFloat(tickerResp.PriceChangePercent, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse price change percent: %w", err)
	}

	openPrice, err := strconv.ParseFloat(tickerResp.OpenPrice, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse open price: %w", err)
	}

	prevClosePrice, err := strconv.ParseFloat(tickerResp.PrevClosePrice, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse previous close price: %w", err)
	}

	bidQty, err := strconv.ParseFloat(tickerResp.BidQty, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse bid quantity: %w", err)
	}

	askQty, err := strconv.ParseFloat(tickerResp.AskQty, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ask quantity: %w", err)
	}

	ticker := &model.Ticker{
		Symbol:             symbol,
		LastPrice:          lastPrice,
		PriceChange:        priceChange,
		PriceChangePercent: priceChangePercent,
		HighPrice:          highPrice,
		LowPrice:           lowPrice,
		Volume:             volume,
		QuoteVolume:        quoteVolume,
		OpenPrice:          openPrice,
		PrevClosePrice:     prevClosePrice,
		BidPrice:           bidPrice,
		BidQty:             bidQty,
		AskPrice:           askPrice,
		AskQty:             askQty,
		Count:              int64(tickerResp.Count),
		Timestamp:          time.Now(),
	}

	return ticker, nil
}

// GetKlines fetches kline/candlestick data for a symbol
func (c *Client) GetKlines(ctx context.Context, symbol string, interval model.KlineInterval, limit int) ([]*model.Kline, error) {
	if limit <= 0 || limit > 1000 {
		limit = 100 // Default to 100 klines if invalid limit
	}

	// Convert domain interval to MEXC interval format
	var intervalStr string
	switch interval {
	case model.KlineInterval1m:
		intervalStr = "1m"
	case model.KlineInterval5m:
		intervalStr = "5m"
	case model.KlineInterval15m:
		intervalStr = "15m"
	case model.KlineInterval30m:
		intervalStr = "30m"
	case model.KlineInterval1h:
		intervalStr = "1h"
	case model.KlineInterval4h:
		intervalStr = "4h"
	case model.KlineInterval1d:
		intervalStr = "1d"
	case model.KlineInterval1w:
		intervalStr = "1w"
	default:
		return nil, fmt.Errorf("invalid kline interval: %s", interval)
	}

	params := map[string]string{
		"symbol":   symbol,
		"interval": intervalStr,
		"limit":    strconv.Itoa(limit),
	}

	resp, err := c.callPublicAPI(ctx, http.MethodGet, "/klines", params)
	if err != nil {
		return nil, fmt.Errorf("failed to get klines: %w", err)
	}

	var klinesResp [][]interface{}
	if err := json.Unmarshal(resp, &klinesResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal klines response: %w", err)
	}

	// Convert response to domain model
	klines := make([]*model.Kline, 0, len(klinesResp))
	for _, k := range klinesResp {
		// MEXC kline format: [OpenTime, Open, High, Low, Close, Volume, CloseTime, QuoteVolume, NumTrades, TakerBuyBaseVol, TakerBuyQuoteVol, Ignored]
		if len(k) < 12 {
			return nil, fmt.Errorf("invalid kline data format")
		}

		// Extract and parse values
		openTimeMs := int64(k[0].(float64))
		closeTimeMs := int64(k[6].(float64))
		open, err := strconv.ParseFloat(k[1].(string), 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse open price: %w", err)
		}
		high, err := strconv.ParseFloat(k[2].(string), 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse high price: %w", err)
		}
		low, err := strconv.ParseFloat(k[3].(string), 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse low price: %w", err)
		}
		closePrice, err := strconv.ParseFloat(k[4].(string), 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse close price: %w", err)
		}
		volume, err := strconv.ParseFloat(k[5].(string), 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse volume: %w", err)
		}
		quoteVolume, err := strconv.ParseFloat(k[7].(string), 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse quote volume: %w", err)
		}
		numTrades := int64(k[8].(float64))

		kline := &model.Kline{
			Symbol:      symbol,
			Interval:    interval,
			OpenTime:    time.Unix(0, openTimeMs*int64(time.Millisecond)),
			CloseTime:   time.Unix(0, closeTimeMs*int64(time.Millisecond)),
			Open:        open,
			High:        high,
			Low:         low,
			Close:       closePrice,
			Volume:      volume,
			QuoteVolume: quoteVolume,
			TradeCount:  numTrades,
			IsClosed:    true, // Assuming historical klines are closed
		}
		klines = append(klines, kline)
	}

	return klines, nil
}

// GetOrderBook fetches the order book for a symbol
func (c *Client) GetOrderBook(ctx context.Context, symbol string, limit int) (*model.OrderBook, error) {
	if limit <= 0 || limit > 5000 {
		limit = 100 // Default to 100 depth levels if invalid limit
	}

	params := map[string]string{
		"symbol": symbol,
		"limit":  strconv.Itoa(limit),
	}

	resp, err := c.callPublicAPI(ctx, http.MethodGet, "/depth", params)
	if err != nil {
		return nil, fmt.Errorf("failed to get order book: %w", err)
	}

	var orderBookResp OrderBookResponse
	if err := json.Unmarshal(resp, &orderBookResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal order book response: %w", err)
	}

	// Convert response to domain model
	orderBook := &model.OrderBook{
		Symbol:       symbol,
		LastUpdateID: orderBookResp.LastUpdateID,
		Bids:         make([]model.OrderBookEntry, 0, len(orderBookResp.Bids)),
		Asks:         make([]model.OrderBookEntry, 0, len(orderBookResp.Asks)),
		Timestamp:    time.Now(),
	}

	// Process bids
	for _, bid := range orderBookResp.Bids {
		if len(bid) < 2 {
			continue
		}
		price, err := strconv.ParseFloat(bid[0], 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse bid price: %w", err)
		}
		quantity, err := strconv.ParseFloat(bid[1], 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse bid quantity: %w", err)
		}
		orderBook.Bids = append(orderBook.Bids, model.OrderBookEntry{
			Price:    price,
			Quantity: quantity,
		})
	}

	// Process asks
	for _, ask := range orderBookResp.Asks {
		if len(ask) < 2 {
			continue
		}
		price, err := strconv.ParseFloat(ask[0], 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse ask price: %w", err)
		}
		quantity, err := strconv.ParseFloat(ask[1], 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse ask quantity: %w", err)
		}
		orderBook.Asks = append(orderBook.Asks, model.OrderBookEntry{
			Price:    price,
			Quantity: quantity,
		})
	}

	return orderBook, nil
}

// sendSignedRequest sends a signed request to the MEXC API
func (c *Client) sendSignedRequest(ctx context.Context, method, endpoint string, params url.Values) (*http.Response, error) {
	signature := c.generateSignature(params)
	params.Add("signature", signature)

	reqURL := c.baseURL + endpoint

	if method == "GET" {
		reqURL = reqURL + "?" + params.Encode()
	}

	var req *http.Request
	var err error

	if method == "POST" {
		req, err = http.NewRequestWithContext(ctx, method, reqURL, strings.NewReader(params.Encode()))
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	} else {
		req, err = http.NewRequestWithContext(ctx, method, reqURL, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}
	}

	req.Header.Set("X-MBX-APIKEY", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	if resp.StatusCode >= 400 {
		var errResp struct {
			Code    int    `json:"code"`
			Message string `json:"msg"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			resp.Body.Close()
			return nil, fmt.Errorf("HTTP error %d", resp.StatusCode)
		}

		resp.Body.Close()
		return nil, fmt.Errorf("API error: %d - %s", errResp.Code, errResp.Message)
	}

	return resp, nil
}

// generateSignature generates a HMAC SHA256 signature for a request
func (c *Client) generateSignature(params url.Values) string {
	queryString := params.Encode()
	mac := hmac.New(sha256.New, []byte(c.secretKey))
	mac.Write([]byte(queryString))
	return hex.EncodeToString(mac.Sum(nil))
}

// PlaceOrder places a new order on the exchange
func (c *Client) PlaceOrder(ctx context.Context, symbol string, side model.OrderSide, orderType model.OrderType, quantity float64, price float64, timeInForce model.TimeInForce) (*model.Order, error) {
	// Create parameters map
	params := map[string]string{
		"symbol": symbol,
		"side":   string(side),
		"type":   string(orderType),
	}

	// Convert quantity to string with appropriate precision
	quantityStr := strconv.FormatFloat(quantity, 'f', -1, 64)
	params["quantity"] = quantityStr

	// Add price for limit orders
	if orderType == model.OrderTypeLimit {
		priceStr := strconv.FormatFloat(price, 'f', -1, 64)
		params["price"] = priceStr
		params["timeInForce"] = string(timeInForce)
	}

	// Generate client order ID
	clientOrderID := fmt.Sprintf("go_bot_%d", time.Now().UnixNano())
	params["newClientOrderId"] = clientOrderID

	// Send request
	resp, err := c.callPrivateAPI(ctx, http.MethodPost, "/order", params, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to place order: %w", err)
	}

	// Parse response
	var response struct {
		Symbol        string `json:"symbol"`
		OrderID       string `json:"orderId"`
		ClientOrderID string `json:"clientOrderId"`
		TransactTime  int64  `json:"transactTime"`
		Price         string `json:"price"`
		OrigQty       string `json:"origQty"`
		ExecutedQty   string `json:"executedQty"`
		Status        string `json:"status"`
		TimeInForce   string `json:"timeInForce"`
		Type          string `json:"type"`
		Side          string `json:"side"`
	}

	if err := json.Unmarshal(resp, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal order response: %w", err)
	}

	// Parse price and quantities
	parsedPrice, err := strconv.ParseFloat(response.Price, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid price in response: %w", err)
	}

	parsedQty, err := strconv.ParseFloat(response.OrigQty, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid quantity in response: %w", err)
	}

	executedQty, err := strconv.ParseFloat(response.ExecutedQty, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid executed quantity in response: %w", err)
	}

	// Create and return order
	order := &model.Order{
		OrderID:       response.OrderID,
		ClientOrderID: response.ClientOrderID,
		Symbol:        response.Symbol,
		Side:          model.OrderSide(response.Side),
		Type:          model.OrderType(response.Type),
		Status:        model.OrderStatus(response.Status),
		Price:         parsedPrice,
		Quantity:      parsedQty,
		ExecutedQty:   executedQty,
		CreatedAt:     time.Unix(0, response.TransactTime*int64(time.Millisecond)),
		UpdatedAt:     time.Now(),
	}

	return order, nil
}

// CancelOrder cancels an existing order
func (c *Client) CancelOrder(ctx context.Context, symbol string, orderID string) error {
	endpoint := "/order"

	params := make(map[string]string)
	params["symbol"] = symbol
	params["orderId"] = orderID

	resp, err := c.callPrivateAPI(ctx, http.MethodDelete, endpoint, params, nil)
	if err != nil {
		return fmt.Errorf("failed to cancel order: %w", err)
	}

	// We don't need to parse the response, just check for errors
	_ = resp

	return nil
}

// GetOrderStatus checks the status of an order
func (c *Client) GetOrderStatus(ctx context.Context, symbol string, orderID string) (*model.Order, error) {
	// Create parameters
	params := map[string]string{
		"symbol": symbol,
	}

	// Determine if orderID is numeric or a client order ID
	if _, err := strconv.ParseInt(orderID, 10, 64); err == nil {
		// OrderID is numeric
		params["orderId"] = orderID
	} else {
		// OrderID is likely a client order ID
		params["origClientOrderId"] = orderID
	}

	// Send request
	resp, err := c.callPrivateAPI(ctx, http.MethodGet, "/order", params, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get order status: %w", err)
	}

	// Parse response
	var response struct {
		Symbol        string `json:"symbol"`
		OrderID       string `json:"orderId"`
		ClientOrderID string `json:"clientOrderId"`
		Price         string `json:"price"`
		OrigQty       string `json:"origQty"`
		ExecutedQty   string `json:"executedQty"`
		Status        string `json:"status"`
		TimeInForce   string `json:"timeInForce"`
		Type          string `json:"type"`
		Side          string `json:"side"`
		Time          int64  `json:"time"`
		UpdateTime    int64  `json:"updateTime"`
	}

	if err := json.Unmarshal(resp, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal order response: %w", err)
	}

	// Parse price and quantities
	parsedPrice, err := strconv.ParseFloat(response.Price, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid price in response: %w", err)
	}

	parsedQty, err := strconv.ParseFloat(response.OrigQty, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid quantity in response: %w", err)
	}

	executedQty, err := strconv.ParseFloat(response.ExecutedQty, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid executed quantity in response: %w", err)
	}

	// Create and return order
	order := &model.Order{
		OrderID:       response.OrderID,
		ClientOrderID: response.ClientOrderID,
		Symbol:        response.Symbol,
		Side:          model.OrderSide(response.Side),
		Type:          model.OrderType(response.Type),
		Status:        model.OrderStatus(response.Status),
		Price:         parsedPrice,
		Quantity:      parsedQty,
		ExecutedQty:   executedQty,
		CreatedAt:     time.Unix(0, response.Time*int64(time.Millisecond)),
		UpdatedAt:     time.Unix(0, response.UpdateTime*int64(time.Millisecond)),
	}

	return order, nil
}
