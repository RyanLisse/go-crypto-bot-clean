package rest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/neo/crypto-bot/internal/domain/model"
	"github.com/neo/crypto-bot/internal/domain/model/market"
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

// GetAccount retrieves account information including balances
func (c *Client) GetAccount() (*model.Account, error) {
	endpoint := "/api/v3/account"

	var response struct {
		MakerCommission  int  `json:"makerCommission"`
		TakerCommission  int  `json:"takerCommission"`
		BuyerCommission  int  `json:"buyerCommission"`
		SellerCommission int  `json:"sellerCommission"`
		CanTrade         bool `json:"canTrade"`
		CanWithdraw      bool `json:"canWithdraw"`
		CanDeposit       bool `json:"canDeposit"`
		Balances         []struct {
			Asset  string `json:"asset"`
			Free   string `json:"free"`
			Locked string `json:"locked"`
		} `json:"balances"`
	}

	// Assuming callPrivateAPI is the correct method for signed requests
	data, err := c.callPrivateAPI(context.Background(), "GET", endpoint, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get account info: %w", err)
	}

	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal account response: %w", err)
	}

	account := model.NewAccount("", "MEXC") // UserID will be set by the application layer

	// Set permissions based on account capabilities
	permissions := make([]string, 0)
	if response.CanTrade {
		permissions = append(permissions, "trade")
	}
	if response.CanWithdraw {
		permissions = append(permissions, "withdraw")
	}
	if response.CanDeposit {
		permissions = append(permissions, "deposit")
	}
	account.Permissions = permissions

	// Update wallet balances
	for _, balance := range response.Balances {
		free, err := strconv.ParseFloat(balance.Free, 64)
		if err != nil {
			continue
		}
		locked, err := strconv.ParseFloat(balance.Locked, 64)
		if err != nil {
			continue
		}

		// Assuming Wallet has an UpdateBalance method with this signature
		account.Wallet.UpdateBalance(model.Asset(balance.Asset), free, locked, 0.0)
	}

	return account, nil
}

// GetMarketData retrieves market data for a symbol
func (c *Client) GetMarketData(ctx context.Context, symbol string) (*model.MarketData, error) {
	// Create new market data instance
	marketData := model.NewMarketData(symbol)

	// Get ticker data
	ticker, err := c.GetTicker(ctx, symbol)
	if err != nil {
		return nil, fmt.Errorf("failed to get ticker: %w", err)
	}
	marketData.Ticker = ticker

	// Get order book data
	orderBook, err := c.GetOrderBook(ctx, symbol, 20)
	if err != nil {
		return nil, fmt.Errorf("failed to get order book: %w", err)
	}

	// Convert model.OrderBookEntry to market.OrderBookEntry
	bids := make([]market.OrderBookEntry, len(orderBook.Bids))
	for i, bid := range orderBook.Bids {
		bids[i] = market.OrderBookEntry{
			Price:    bid.Price,
			Quantity: bid.Quantity,
		}
	}

	asks := make([]market.OrderBookEntry, len(orderBook.Asks))
	for i, ask := range orderBook.Asks {
		asks[i] = market.OrderBookEntry{
			Price:    ask.Price,
			Quantity: ask.Quantity,
		}
	}

	marketData.OrderBook = market.OrderBook{
		Exchange:     "MEXC",
		Symbol:       symbol,
		Bids:         bids,
		Asks:         asks,
		LastUpdated:  time.Now(),
		LastUpdateID: orderBook.LastUpdateID,
	}

	// Get recent trades
	trades, err := c.GetRecentTrades(symbol, 1)
	if err == nil && len(trades) > 0 {
		marketData.LastTrade = trades[0]
	}

	return marketData, nil
}

// GetKlines retrieves kline/candlestick data for a symbol
func (c *Client) GetKlines(ctx context.Context, symbol string, interval string, limit int) ([]model.Kline, error) {
	params := map[string]string{
		"symbol":   symbol,
		"interval": interval,
		"limit":    strconv.Itoa(limit),
	}

	data, err := c.callPublicAPI(ctx, "GET", "/api/v3/klines", params)
	if err != nil {
		return nil, fmt.Errorf("failed to get klines: %w", err)
	}

	var klineData [][]interface{}
	if err := json.Unmarshal(data, &klineData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal klines response: %w", err)
	}

	klines := make([]model.Kline, len(klineData))
	for i, k := range klineData {
		openTime := time.UnixMilli(int64(k[0].(float64)))
		closeTime := time.UnixMilli(int64(k[6].(float64)))

		open, _ := strconv.ParseFloat(k[1].(string), 64)
		high, _ := strconv.ParseFloat(k[2].(string), 64)
		low, _ := strconv.ParseFloat(k[3].(string), 64)
		close, _ := strconv.ParseFloat(k[4].(string), 64)
		volume, _ := strconv.ParseFloat(k[5].(string), 64)

		klines[i] = model.Kline{
			Symbol:    symbol,
			Interval:  model.KlineInterval(interval),
			OpenTime:  openTime,
			Open:      open,
			High:      high,
			Low:       low,
			Close:     close,
			Volume:    volume,
			CloseTime: closeTime,
			IsClosed:  true, // MEXC returns only closed klines
		}
	}

	return klines, nil
}

// GetOrderBook retrieves order book data for a symbol
func (c *Client) GetOrderBook(ctx context.Context, symbol string, limit int) (*model.OrderBook, error) {
	params := map[string]string{
		"symbol": symbol,
		"limit":  strconv.Itoa(limit),
	}

	data, err := c.callPublicAPI(ctx, "GET", "/api/v3/depth", params)
	if err != nil {
		return nil, fmt.Errorf("failed to get order book: %w", err)
	}

	var orderBookData struct {
		LastUpdateID int64      `json:"lastUpdateId"`
		Bids         [][]string `json:"bids"`
		Asks         [][]string `json:"asks"`
	}
	if err := json.Unmarshal(data, &orderBookData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal order book response: %w", err)
	}

	bids := make([]model.OrderBookEntry, len(orderBookData.Bids))
	for i, bid := range orderBookData.Bids {
		price, _ := strconv.ParseFloat(bid[0], 64)
		quantity, _ := strconv.ParseFloat(bid[1], 64)
		bids[i] = model.OrderBookEntry{
			Price:    price,
			Quantity: quantity,
		}
	}

	asks := make([]model.OrderBookEntry, len(orderBookData.Asks))
	for i, ask := range orderBookData.Asks {
		price, _ := strconv.ParseFloat(ask[0], 64)
		quantity, _ := strconv.ParseFloat(ask[1], 64)
		asks[i] = model.OrderBookEntry{
			Price:    price,
			Quantity: quantity,
		}
	}

	return &model.OrderBook{
		Symbol:       symbol,
		LastUpdateID: orderBookData.LastUpdateID,
		Bids:         bids,
		Asks:         asks,
	}, nil
}

// PlaceOrder places a new order on the exchange
func (c *Client) PlaceOrder(ctx context.Context, symbol string, side model.OrderSide, orderType model.OrderType, quantity float64, price float64, timeInForce model.TimeInForce) (*model.Order, error) {
	params := map[string]string{
		"symbol":   symbol,
		"side":     string(side),
		"type":     string(orderType),
		"quantity": strconv.FormatFloat(quantity, 'f', -1, 64),
	}

	if orderType == model.OrderTypeLimit {
		params["timeInForce"] = string(timeInForce)
		params["price"] = strconv.FormatFloat(price, 'f', -1, 64)
	}

	data, err := c.callPrivateAPI(ctx, "POST", "/api/v3/order", params, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to place order: %w", err)
	}

	var orderResp struct {
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
	if err := json.Unmarshal(data, &orderResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal order response: %w", err)
	}

	price, _ = strconv.ParseFloat(orderResp.Price, 64)
	origQty, _ := strconv.ParseFloat(orderResp.OrigQty, 64)
	executedQty, _ := strconv.ParseFloat(orderResp.ExecutedQty, 64)

	return &model.Order{
		OrderID:       orderResp.OrderID,
		ClientOrderID: orderResp.ClientOrderID,
		Symbol:        orderResp.Symbol,
		Side:          model.OrderSide(strings.ToUpper(orderResp.Side)),
		Type:          model.OrderType(strings.ToUpper(orderResp.Type)),
		Status:        model.OrderStatus(strings.ToUpper(orderResp.Status)),
		Price:         price,
		Quantity:      origQty,
		ExecutedQty:   executedQty,
		CreatedAt:     time.UnixMilli(orderResp.Time),
		UpdatedAt:     time.UnixMilli(orderResp.UpdateTime),
	}, nil
}

// CancelOrder cancels an existing order
func (c *Client) CancelOrder(ctx context.Context, symbol string, orderID string) error {
	params := map[string]string{
		"symbol":  symbol,
		"orderId": orderID,
	}

	_, err := c.callPrivateAPI(ctx, "DELETE", "/api/v3/order", params, nil)
	if err != nil {
		return fmt.Errorf("failed to cancel order: %w", err)
	}

	return nil
}

// GetOrderStatus retrieves the status of an order from the exchange
func (c *Client) GetOrderStatus(ctx context.Context, symbol string, orderID string) (*model.Order, error) {
	params := map[string]string{
		"symbol":  symbol,
		"orderId": orderID,
	}

	data, err := c.callPrivateAPI(ctx, "GET", "/api/v3/order", params, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get order status: %w", err)
	}

	var orderResp struct {
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
	if err := json.Unmarshal(data, &orderResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal order response: %w", err)
	}

	price, _ := strconv.ParseFloat(orderResp.Price, 64)
	quantity, _ := strconv.ParseFloat(orderResp.OrigQty, 64)
	executedQty, _ := strconv.ParseFloat(orderResp.ExecutedQty, 64)

	return &model.Order{
		OrderID:       orderResp.OrderID,
		ClientOrderID: orderResp.ClientOrderID,
		Symbol:        orderResp.Symbol,
		Side:          model.OrderSide(strings.ToUpper(orderResp.Side)),
		Type:          model.OrderType(strings.ToUpper(orderResp.Type)),
		Status:        model.OrderStatus(strings.ToUpper(orderResp.Status)),
		Price:         price,
		Quantity:      quantity,
		ExecutedQty:   executedQty,
		CreatedAt:     time.UnixMilli(orderResp.Time),
		UpdatedAt:     time.UnixMilli(orderResp.UpdateTime),
	}, nil
}

// GetTicker retrieves current ticker data for a symbol
func (c *Client) GetTicker(ctx context.Context, symbol string) (*market.Ticker, error) {
	endpoint := "/api/v3/ticker/24hr"
	params := map[string]string{
		"symbol": symbol,
	}

	resp, err := c.callPublicAPI(ctx, "GET", endpoint, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get ticker: %w", err)
	}

	var data struct {
		Symbol             string `json:"symbol"`
		PriceChange        string `json:"priceChange"`
		PriceChangePercent string `json:"priceChangePercent"`
		LastPrice          string `json:"lastPrice"`
		Volume             string `json:"volume"`
		QuoteVolume        string `json:"quoteVolume"`
		High               string `json:"highPrice"`
		Low                string `json:"lowPrice"`
		Bid                string `json:"bidPrice"`
		Ask                string `json:"askPrice"`
	}

	if err := json.Unmarshal(resp, &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal ticker response: %w", err)
	}

	// Convert string values to float64
	price, _ := strconv.ParseFloat(data.LastPrice, 64)
	volume, _ := strconv.ParseFloat(data.Volume, 64)
	high24h, _ := strconv.ParseFloat(data.High, 64)
	low24h, _ := strconv.ParseFloat(data.Low, 64)
	priceChange, _ := strconv.ParseFloat(data.PriceChange, 64)
	percentChange, _ := strconv.ParseFloat(data.PriceChangePercent, 64)

	return &market.Ticker{
		Exchange:      "MEXC",
		Symbol:        data.Symbol,
		Price:         price,
		Volume:        volume,
		High24h:       high24h,
		Low24h:        low24h,
		PriceChange:   priceChange,
		PercentChange: percentChange,
		LastUpdated:   time.Now(),
	}, nil
}

// Placeholder implementations for missing methods - These need to be fully implemented!
func (c *Client) GetRecentTrades(symbol string, limit int) ([]model.MarketTrade, error) {
	// TODO: Implement GetRecentTrades API call
	return nil, errors.New("GetRecentTrades not implemented")
}
