package gateway

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/apperror"
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/config"
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/port/gateway"
	"github.com/rs/zerolog"
	"github.com/shopspring/decimal"
)

// Ensure mexcAdapter implements the MEXCGateway interface.
var _ gateway.MEXCGateway = (*mexcAdapter)(nil)

const (
	mexcBaseURL  = "https://api.mexc.com"
	apiKeyHeader = "X-MEXC-APIKEY"
)

type mexcAdapter struct {
	config     config.MEXCConfig
	httpClient *http.Client
	logger     *zerolog.Logger
}

// NewMEXCGateway creates a new adapter for interacting with MEXC.
func NewMEXCGateway(cfg config.MEXCConfig, logger *zerolog.Logger) (gateway.MEXCGateway, error) {
	if cfg.APIKey == "" || cfg.SecretKey == "" {
		logger.Warn().Msg("MEXC APIKey or SecretKey not set. Authenticated endpoints will fail.")
		// Allow creation but authenticated calls will fail
	}
	return &mexcAdapter{
		config: cfg,
		httpClient: &http.Client{
			Timeout: 10 * time.Second, // Consider making timeout configurable
		},
		logger: logger,
	}, nil
}

// --- Helper Methods ---

// signRequest generates the signature for MEXC API requests.
func (a *mexcAdapter) signRequest(method, path string, params url.Values) string {
	timestamp := strconv.FormatInt(time.Now().Unix()*1000, 10)
	params.Set("timestamp", timestamp)

	// Create string to sign
	var dataToSign string
	if method == http.MethodGet || method == http.MethodDelete {
		// Sort parameters alphabetically for GET/DELETE
		var keys []string
		for k := range params {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		var sortedParams []string
		for _, k := range keys {
			sortedParams = append(sortedParams, k+"="+params.Get(k))
		}
		queryString := strings.Join(sortedParams, "&")
		dataToSign = path + "?" + queryString
		// MEXC v3 GET/DELETE doesn't seem to use timestamp/signature according to some docs, verify!
		// For simplicity, let's assume timestamp is needed but signature might depend on endpoint
		// Let's try signing with timestamp only for now based on some interpretations
		dataToSign = timestamp // Revisit this based on official V3 signing docs
	} else { // POST/PUT
		// For POST/PUT, sign the timestamp + request body
		// Assuming request body is handled separately and passed here if needed
		// For simplicity, let's assume params are in body for POST
		// The exact V3 signing can be complex (query params vs body)
		// Let's use timestamp + query params for now as a starting point
		bodyString := params.Encode()
		dataToSign = timestamp + bodyString
	}

	h := hmac.New(sha256.New, []byte(a.config.SecretKey))
	h.Write([]byte(dataToSign))
	signature := hex.EncodeToString(h.Sum(nil))
	return signature
}

// createRequest creates a new HTTP request, adding necessary headers and signature.
func (a *mexcAdapter) createRequest(ctx context.Context, method, endpointPath string, params url.Values, requiresAuth bool) (*http.Request, error) {
	fullURL := mexcBaseURL + endpointPath
	var reqBody io.Reader

	if requiresAuth {
		if a.config.APIKey == "" || a.config.SecretKey == "" {
			return nil, apperror.NewUnauthorized("MEXC API credentials not configured", nil)
		}
		timestamp := strconv.FormatInt(time.Now().Unix()*1000, 10)
		params.Set("timestamp", timestamp)
		// Re-calculate signature here based on final params
		// signature := a.signRequest(method, endpointPath, params) // Signing logic needs refinement
		// params.Set("signature", signature)
	}

	if method == http.MethodGet || method == http.MethodDelete {
		if len(params) > 0 {
			fullURL += "?" + params.Encode()
		}
	} else {
		// Assume params go in the body for POST/PUT
		reqBody = strings.NewReader(params.Encode())
	}

	req, err := http.NewRequestWithContext(ctx, method, fullURL, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if requiresAuth {
		req.Header.Set(apiKeyHeader, a.config.APIKey)
	}

	if method == http.MethodPost || method == http.MethodPut {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded") // Or application/json if API requires
	}
	req.Header.Set("Accept", "application/json")

	return req, nil
}

// sendRequest sends the request and handles the response.
func (a *mexcAdapter) sendRequest(req *http.Request, result interface{}) error {
	a.logger.Debug().Str("method", req.Method).Str("url", req.URL.String()).Msg("Sending MEXC request")

	resp, err := a.httpClient.Do(req)
	if err != nil {
		a.logger.Error().Err(err).Str("url", req.URL.String()).Msg("MEXC request failed")
		// Consider wrapping specific network errors (e.g., timeout) with apperror
		return apperror.NewExternalServiceError("MEXC request failed: "+err.Error(), err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		a.logger.Error().Err(err).Str("url", req.URL.String()).Msg("Failed to read MEXC response body")
		return apperror.NewInternal(fmt.Errorf("failed to read MEXC response body: %w", err))
	}

	a.logger.Debug().Int("status_code", resp.StatusCode).Str("url", req.URL.String()).RawJSON("body", bodyBytes).Msg("Received MEXC response")

	// TODO: Implement robust MEXC error handling based on status code and response body
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// Try to parse MEXC error response format
		// var mexcErr MexcErrorResponse
		// if json.Unmarshal(bodyBytes, &mexcErr) == nil && mexcErr.Code != 0 {
		//    return apperror.NewExternalServiceError(fmt.Sprintf("MEXC API error %d: %s", mexcErr.Code, mexcErr.Msg), errors.New(string(bodyBytes)))
		// }
		return apperror.NewExternalServiceError(fmt.Sprintf("MEXC API request failed with status %d: %s", resp.StatusCode, string(bodyBytes)), errors.New("mexc api error"))
	}

	if result != nil {
		if err := json.Unmarshal(bodyBytes, result); err != nil {
			a.logger.Error().Err(err).Str("url", req.URL.String()).Msg("Failed to decode MEXC response JSON")
			return apperror.NewInternal(fmt.Errorf("failed to decode MEXC response: %w", err))
		}
	}

	return nil
}

// --- Interface Implementations ---

func (a *mexcAdapter) GetSymbols(ctx context.Context) ([]model.SymbolInfo, error) {
	params := url.Values{}
	req, err := a.createRequest(ctx, http.MethodGet, "/api/v3/exchangeInfo", params, false)
	if err != nil {
		return nil, err
	}
	type MexcExchangeInfoResponse struct {
		Symbols []struct {
			Symbol              string `json:"symbol"`
			Status              string `json:"status"`
			BaseAsset           string `json:"baseAsset"`
			BaseAssetPrecision  int    `json:"baseAssetPrecision"`
			QuoteAsset          string `json:"quoteAsset"`
			QuoteAssetPrecision int    `json:"quotePrecision"`
		} `json:"symbols"`
	}
	var response MexcExchangeInfoResponse
	if err := a.sendRequest(req, &response); err != nil {
		return nil, err
	}
	symbols := make([]model.SymbolInfo, 0, len(response.Symbols))
	for _, mexcSymbol := range response.Symbols {
		symbols = append(symbols, model.SymbolInfo{
			Symbol:              mexcSymbol.Symbol,
			Status:              mexcSymbol.Status,
			BaseAsset:           mexcSymbol.BaseAsset,
			BaseAssetPrecision:  mexcSymbol.BaseAssetPrecision,
			QuoteAsset:          mexcSymbol.QuoteAsset,
			QuoteAssetPrecision: mexcSymbol.QuoteAssetPrecision,
		})
	}
	return symbols, nil
}

func (a *mexcAdapter) GetTicker(ctx context.Context, symbol string) (model.Ticker, error) {
	params := url.Values{}
	params.Set("symbol", symbol)
	req, err := a.createRequest(ctx, http.MethodGet, "/api/v3/ticker/24hr", params, false)
	if err != nil {
		return model.Ticker{}, err
	}
	type MexcTickerResponse struct {
		Symbol             string `json:"symbol"`
		PriceChange        string `json:"priceChange"`
		PriceChangePercent string `json:"priceChangePercent"`
		LastPrice          string `json:"lastPrice"`
		BidPrice           string `json:"bidPrice"`
		AskPrice           string `json:"askPrice"`
		OpenPrice          string `json:"openPrice"`
		HighPrice          string `json:"highPrice"`
		LowPrice           string `json:"lowPrice"`
		Volume             string `json:"volume"`
		QuoteVolume        string `json:"quoteVolume"`
		OpenTime           int64  `json:"openTime"`
		CloseTime          int64  `json:"closeTime"`
		// Add WeightedAvgPrice, PrevClosePrice, LastQty, BidQty, AskQty, Count if needed
	}
	var response MexcTickerResponse
	if err := a.sendRequest(req, &response); err != nil {
		return model.Ticker{}, err
	}

	parseStrToFloat := func(s string, fieldName string) float64 {
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			a.logger.Warn().Err(err).Str("field", fieldName).Str("value", s).Msg("Failed to parse float from MEXC ticker string")
			return 0.0
		}
		return f
	}

	ticker := model.Ticker{
		Symbol:             response.Symbol,
		Exchange:           "MEXC", // Hardcode exchange name for this adapter
		PriceChange:        parseStrToFloat(response.PriceChange, "PriceChange"),
		PriceChangePercent: parseStrToFloat(response.PriceChangePercent, "PriceChangePercent"),
		LastPrice:          parseStrToFloat(response.LastPrice, "LastPrice"),
		BidPrice:           parseStrToFloat(response.BidPrice, "BidPrice"),
		AskPrice:           parseStrToFloat(response.AskPrice, "AskPrice"),
		OpenPrice:          parseStrToFloat(response.OpenPrice, "OpenPrice"),
		HighPrice:          parseStrToFloat(response.HighPrice, "HighPrice"),
		LowPrice:           parseStrToFloat(response.LowPrice, "LowPrice"),
		Volume:             parseStrToFloat(response.Volume, "Volume"),
		QuoteVolume:        parseStrToFloat(response.QuoteVolume, "QuoteVolume"),
		Timestamp:          time.UnixMilli(response.CloseTime), // Use CloseTime as the ticker timestamp
		// Map other fields like PrevClosePrice, BidQty, AskQty, Count if available and needed
	}

	return ticker, nil
}

func (a *mexcAdapter) GetOrderBook(ctx context.Context, symbol string, depth int) (model.OrderBook, error) {
	params := url.Values{}
	params.Set("symbol", symbol)
	if depth > 0 {
		if depth > 5000 {
			depth = 5000 // MEXC max depth is 5000
		}
		params.Set("limit", strconv.Itoa(depth))
	}
	req, err := a.createRequest(ctx, http.MethodGet, "/api/v3/depth", params, false)
	if err != nil {
		return model.OrderBook{}, err
	}
	type MexcOrderBookResponse struct {
		LastUpdateID int64      `json:"lastUpdateId"`
		Bids         [][]string `json:"bids"`
		Asks         [][]string `json:"asks"`
	}
	var response MexcOrderBookResponse
	if err := a.sendRequest(req, &response); err != nil {
		return model.OrderBook{}, err
	}

	parseLevel := func(level []string, entryType string) (model.OrderBookEntry, error) {
		if len(level) != 2 {
			return model.OrderBookEntry{}, fmt.Errorf("invalid order book level format: %v", level)
		}
		price, err := strconv.ParseFloat(level[0], 64)
		if err != nil {
			a.logger.Warn().Err(err).Str("price", level[0]).Msg("Failed to parse price in order book")
			return model.OrderBookEntry{}, fmt.Errorf("failed to parse price '%s': %w", level[0], err)
		}
		qty, err := strconv.ParseFloat(level[1], 64)
		if err != nil {
			a.logger.Warn().Err(err).Str("quantity", level[1]).Msg("Failed to parse quantity in order book")
			return model.OrderBookEntry{}, fmt.Errorf("failed to parse quantity '%s': %w", level[1], err)
		}
		return model.OrderBookEntry{Price: price, Quantity: qty}, nil
	}

	ob := model.OrderBook{
		Symbol:       symbol,
		Exchange:     "MEXC",
		LastUpdateID: response.LastUpdateID,
		Bids:         make([]model.OrderBookEntry, 0, len(response.Bids)),
		Asks:         make([]model.OrderBookEntry, 0, len(response.Asks)),
		Timestamp:    time.Now().UTC(), // Use current time as MEXC doesn't provide timestamp here
	}

	for _, bidLevel := range response.Bids {
		entry, err := parseLevel(bidLevel, "bid")
		if err != nil {
			// Log or skip invalid level?
			a.logger.Error().Err(err).Msg("Skipping invalid bid level in order book")
			continue
		}
		ob.Bids = append(ob.Bids, entry)
	}
	for _, askLevel := range response.Asks {
		entry, err := parseLevel(askLevel, "ask")
		if err != nil {
			a.logger.Error().Err(err).Msg("Skipping invalid ask level in order book")
			continue
		}
		ob.Asks = append(ob.Asks, entry)
	}

	return ob, nil
}

func (a *mexcAdapter) GetAccountInfo(ctx context.Context) (model.AccountInfo, error) {
	params := url.Values{}
	req, err := a.createRequest(ctx, http.MethodGet, "/api/v3/account", params, true)
	if err != nil {
		return model.AccountInfo{}, err
	}
	type MexcAccountInfoResponse struct {
		CanTrade    bool  `json:"canTrade"`
		CanWithdraw bool  `json:"canWithdraw"`
		CanDeposit  bool  `json:"canDeposit"`
		UpdateTime  int64 `json:"updateTime"`
		Balances    []struct {
			Asset  string `json:"asset"`
			Free   string `json:"free"`
			Locked string `json:"locked"`
		} `json:"balances"`
		// Add other fields if needed (commissions, accountType, permissions)
	}
	var response MexcAccountInfoResponse
	if err := a.sendRequest(req, &response); err != nil {
		return model.AccountInfo{}, err
	}

	accountInfo := model.AccountInfo{
		// UserID needs to be inferred from context or API key association if possible
		CanTrade:    response.CanTrade,
		CanWithdraw: response.CanWithdraw,
		CanDeposit:  response.CanDeposit,
		Balances:    make([]model.Balance, len(response.Balances)),
		LastUpdated: time.UnixMilli(response.UpdateTime),
	}
	for i, b := range response.Balances {
		freeDec, errF := decimal.NewFromString(b.Free)
		if errF != nil {
			a.logger.Error().Err(errF).Str("asset", b.Asset).Str("value", b.Free).Msg("Failed to parse free balance")
			// Decide how to handle - skip asset, return error, set to zero?
			continue
		}
		lockedDec, errL := decimal.NewFromString(b.Locked)
		if errL != nil {
			a.logger.Error().Err(errL).Str("asset", b.Asset).Str("value", b.Locked).Msg("Failed to parse locked balance")
			continue
		}

		freeF64, _ := freeDec.Float64()
		lockedF64, _ := lockedDec.Float64()

		accountInfo.Balances[i] = model.Balance{
			Asset:  model.Asset(b.Asset), // Cast string to model.Asset
			Free:   freeF64,              // Convert decimal to float64
			Locked: lockedF64,            // Convert decimal to float64
			Total:  freeF64 + lockedF64,
			// USDValue would require fetching price data separately
		}
	}

	return accountInfo, nil
}

func (a *mexcAdapter) GetAssetBalance(ctx context.Context, asset string) (decimal.Decimal, error) {
	accountInfo, err := a.GetAccountInfo(ctx)
	if err != nil {
		return decimal.Zero, err
	}

	for _, balance := range accountInfo.Balances {
		// Cast model.Asset to string for comparison
		if strings.EqualFold(string(balance.Asset), asset) {
			// Need to return decimal.Decimal, but model.Balance stores float64.
			// For consistency, let's re-parse from the original strings or fetch account info again.
			// Re-parsing for simplicity here:
			for _, b := range accountInfo.Balances { // Assuming original response structure if needed again
				if strings.EqualFold(string(b.Asset), asset) { // Check asset again
					freeDec, err := decimal.NewFromString(fmt.Sprintf("%f", b.Free)) // Hacky: convert float back to string
					if err != nil {
						return decimal.Zero, apperror.NewInternal(fmt.Errorf("failed to re-parse free balance for %s: %w", asset, err))
					}
					return freeDec, nil
				}
			}
			// Fallback if re-parsing fails (shouldn't happen)
			return decimal.NewFromFloat(balance.Free), nil
		}
	}

	return decimal.Zero, apperror.NewNotFound(fmt.Sprintf("asset %s not found in account balances", asset), nil, nil)
}

// PlaceOrder places an order on the exchange, aligning with portgateway.MEXCGateway interface.
func (a *mexcAdapter) PlaceOrder(ctx context.Context, order model.OrderRequest) (model.OrderResponse, error) {
	params := url.Values{}
	params.Set("symbol", order.Symbol)
	params.Set("side", string(order.Side))
	params.Set("type", string(order.Type))
	params.Set("quantity", strconv.FormatFloat(order.Quantity, 'f', -1, 64))

	if order.Type == model.OrderTypeLimit {
		if order.Price <= 0 {
			// Provide message, nil details, and nil original error for NewBadRequest
			return model.OrderResponse{}, apperror.NewBadRequest("Price must be positive for limit orders", nil, nil)
		}
		params.Set("price", strconv.FormatFloat(order.Price, 'f', -1, 64))
	}
	// TODO: Handle timeInForce if required by MEXC API for specific order types

	req, err := a.createRequest(ctx, http.MethodPost, "/api/v3/order", params, true)
	if err != nil {
		return model.OrderResponse{}, err
	}

	// Define the expected successful response structure from MEXC for placing an order
	type MexcPlaceOrderSuccessResponse struct {
		Symbol        string `json:"symbol"`
		OrderID       string `json:"orderId"`
		OrderListId   int64  `json:"orderListId"`
		ClientOrderID string `json:"clientOrderId"`
		TransactTime  int64  `json:"transactTime"`
	}

	var successResp MexcPlaceOrderSuccessResponse
	if err := a.sendRequest(req, &successResp); err != nil {
		return model.OrderResponse{IsSuccess: false}, err
	}

	// Map the successful MEXC response to the domain model.OrderResponse
	orderResponse := model.OrderResponse{
		IsSuccess: true,
		Order: model.Order{
			Exchange:      "MEXC",
			ID:            successResp.OrderID,
			OrderID:       successResp.OrderID,
			ClientOrderID: successResp.ClientOrderID,
			Symbol:        successResp.Symbol,
			Side:          order.Side,
			Type:          order.Type,
			Status:        model.OrderStatusNew,
			Price:         order.Price,
			Quantity:      order.Quantity,
			TimeInForce:   order.TimeInForce,
			CreatedAt:     time.UnixMilli(successResp.TransactTime),
			UpdatedAt:     time.UnixMilli(successResp.TransactTime),
		},
	}

	return orderResponse, nil
}

func (a *mexcAdapter) CancelOrder(ctx context.Context, symbol, orderID string) error {
	params := url.Values{}
	params.Set("symbol", symbol)
	params.Set("orderId", orderID)
	req, err := a.createRequest(ctx, http.MethodDelete, "/api/v3/order", params, true)
	if err != nil {
		return err
	}
	if err := a.sendRequest(req, nil); err != nil {
		return err
	}
	return nil
}

// GetExchangeInfo retrieves general exchange information
func (a *mexcAdapter) GetExchangeInfo(ctx context.Context) (*model.ExchangeInfo, error) {
	params := url.Values{}
	req, err := a.createRequest(ctx, http.MethodGet, "/api/v3/exchangeInfo", params, false)
	if err != nil {
		return nil, err
	}
	// Using model.ExchangeInfo directly if it matches MEXC's structure
	var response model.ExchangeInfo
	if err := a.sendRequest(req, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// GetNewListings retrieves recently listed or updated coins from the source.
// This method makes mexcAdapter satisfy the port.ListingDetector interface.
func (a *mexcAdapter) GetNewListings(ctx context.Context) ([]*model.NewCoin, error) {
	a.logger.Debug().
		Str("component", "mexcAdapter").
		Str("method", "GetNewListings").
		Str("data_source", "REAL_API").
		Msg("Fetching new listings from MEXC API (placeholder)")

	// TODO: Implement the actual API call to fetch new listings from MEXC.
	// The current implementation returns an empty list as a placeholder.
	// The actual API endpoint and response structure need to be determined
	// from the MEXC API documentation.
	return []*model.NewCoin{}, nil // Placeholder implementation
}

// Ensure mexcAdapter implements ListingDetector
// var _ port.ListingDetector = (*mexcAdapter)(nil) // Need to import port package
