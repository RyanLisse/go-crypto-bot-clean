package rest

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/neo/crypto-bot/internal/domain/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestClient provides a test client with a mock server
func setupTestClient(handler http.Handler) (*Client, *httptest.Server, func()) {
	// Create a test server
	server := httptest.NewServer(handler)

	// Create a client that uses the test server URL
	client := NewClient("testApiKey", "testSecretKey", WithBaseURL(server.URL))

	// Return a cleanup function
	cleanup := func() {
		server.Close()
	}

	return client, server, cleanup
}

// mockHandler returns a handler that serves a predefined response
func mockHandler(statusCode int, responseBody string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set content type
		w.Header().Set("Content-Type", "application/json")

		// For the test endpoints, we'll just accept the requests without checking auth
		// since we're testing the response parsing, not the authentication

		// Set the status code
		w.WriteHeader(statusCode)

		// Write the response body
		w.Write([]byte(responseBody))
	})
}

func TestGetAccount(t *testing.T) {
	// Sample response from MEXC API documentation
	responseBody := `{
		"makerCommission": 15,
		"takerCommission": 15,
		"buyerCommission": 0,
		"sellerCommission": 0,
		"canTrade": true,
		"canWithdraw": true,
		"canDeposit": true,
		"updateTime": 1641182585000,
		"accountType": "SPOT",
		"balances": [
			{
				"asset": "BTC",
				"free": "0.1",
				"locked": "0.05"
			},
			{
				"asset": "ETH",
				"free": "2.5",
				"locked": "0.0"
			},
			{
				"asset": "USDT",
				"free": "1000.0",
				"locked": "500.0"
			}
		]
	}`

	client, _, cleanup := setupTestClient(mockHandler(http.StatusOK, responseBody))
	defer cleanup()

	// Test GetAccount
	wallet, err := client.GetAccount()

	// Verify results
	require.NoError(t, err)
	require.NotNil(t, wallet)

	// Check if wallet has the correct balances
	assert.NotNil(t, wallet.Wallet)
	assert.Len(t, wallet.Wallet.Balances, 3)
	assert.Equal(t, 0.1, wallet.Wallet.Balances[model.Asset("BTC")].Free)
	assert.Equal(t, 0.05, wallet.Wallet.Balances[model.Asset("BTC")].Locked)
	assert.InDelta(t, 0.15, wallet.Wallet.Balances[model.Asset("BTC")].Total, 0.000001)

	assert.Equal(t, 2.5, wallet.Wallet.Balances[model.Asset("ETH")].Free)
	assert.Equal(t, 0.0, wallet.Wallet.Balances[model.Asset("ETH")].Locked)
	assert.Equal(t, 2.5, wallet.Wallet.Balances[model.Asset("ETH")].Total)

	assert.Equal(t, 1000.0, wallet.Wallet.Balances[model.Asset("USDT")].Free)
	assert.Equal(t, 500.0, wallet.Wallet.Balances[model.Asset("USDT")].Locked)
	assert.Equal(t, 1500.0, wallet.Wallet.Balances[model.Asset("USDT")].Total)
}

func TestGetMarketData(t *testing.T) {
	// Sample response based on MEXC API documentation
	responseBody := `{
		"symbol": "BTCUSDT",
		"lastPrice": "42000.0",
		"priceChange": "100.0",
		"priceChangePercent": "2.5",
		"highPrice": "42200.0",
		"lowPrice": "41700.0",
		"volume": "100.0",
		"quoteVolume": "4200000.0",
		"bidPrice": "41995.0",
		"bidQty": "1.5",
		"askPrice": "42005.0",
		"askQty": "2.0",
		"count": 5000
	}`

	client, _, cleanup := setupTestClient(mockHandler(http.StatusOK, responseBody))
	defer cleanup()

	// Test GetMarketData
	marketData, err := client.GetMarketData(context.Background(), "BTCUSDT")

	// Verify results
	require.NoError(t, err)
	require.NotNil(t, marketData)

	// Check if market data has the correct values
	assert.Equal(t, "BTCUSDT", marketData.Symbol)

	// Check ticker data
	require.NotNil(t, marketData.Ticker)              // Ensure it's using the correct type from model
	assert.Equal(t, 42000.0, marketData.Ticker.Price) // Use Price from market.Ticker
	assert.Equal(t, 100.0, marketData.Ticker.PriceChange)
	assert.Equal(t, 2.5, marketData.Ticker.PercentChange) // Use PercentChange from market.Ticker
	assert.Equal(t, 42200.0, marketData.Ticker.High24h)   // Use High24h from market.Ticker
	assert.Equal(t, 41700.0, marketData.Ticker.Low24h)    // Use Low24h from market.Ticker
	assert.Equal(t, 100.0, marketData.Ticker.Volume)
	// QuoteVolume, BidPrice, BidQty, AskPrice, AskQty, Count are not in market.Ticker
	// Add assertions for fields that *are* in market.Ticker
	assert.Equal(t, "MEXC", marketData.Ticker.Exchange)
	assert.NotZero(t, marketData.Ticker.LastUpdated)
}

func TestPlaceOrder(t *testing.T) {
	// Create test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/order", r.URL.Path)

		// Send response
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{
			"symbol": "BTC_USDT",
			"orderId": 12345,
			"clientOrderId": "test123",
			"transactTime": 1499827319559,
			"price": "0.1",
			"origQty": "1.0",
			"executedQty": "0.0",
			"status": "NEW",
			"timeInForce": "GTC",
			"type": "LIMIT",
			"side": "BUY"
		}`)
	}))
	defer ts.Close()

	// Create client
	client := NewClient("test-api-key", "test-secret-key", WithBaseURL(ts.URL))

	// Test place order
	order, err := client.PlaceOrder(context.Background(),
		"BTC_USDT",
		model.OrderSideBuy,
		model.OrderTypeLimit,
		1.0,
		0.1,
		model.TimeInForceGTC)

	// Verify response
	require.NoError(t, err)
	assert.NotNil(t, order)
	assert.Equal(t, "12345", order.OrderID)
	assert.Equal(t, "test123", order.ClientOrderID)
	assert.Equal(t, model.OrderSideBuy, order.Side)
	assert.Equal(t, model.OrderTypeLimit, order.Type)
	assert.Equal(t, model.OrderStatusNew, order.Status)
}

func TestGetOrderStatus(t *testing.T) {
	// Create test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/order", r.URL.Path)

		// Send response
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{
			"symbol": "BTC_USDT",
			"orderId": 123456789,
			"clientOrderId": "test123",
			"price": "0.1",
			"origQty": "1.0",
			"executedQty": "0.0",
			"status": "NEW",
			"timeInForce": "GTC",
			"type": "LIMIT",
			"side": "BUY",
			"time": 1499827319559,
			"updateTime": 1499827319559
		}`)
	}))
	defer ts.Close()

	// Create client
	client := NewClient("test-api-key", "test-secret-key", WithBaseURL(ts.URL))

	// Test get order status
	order, err := client.GetOrderStatus(context.Background(), "BTC_USDT", "123456789")

	// Verify response
	require.NoError(t, err)
	assert.NotNil(t, order)
	assert.Equal(t, "123456789", order.OrderID)
	assert.Equal(t, "test123", order.ClientOrderID)
	assert.Equal(t, model.OrderSideBuy, order.Side)
	assert.Equal(t, model.OrderTypeLimit, order.Type)
	assert.Equal(t, model.OrderStatusNew, order.Status)
}

func TestCancelOrder(t *testing.T) {
	// Create test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "/order", r.URL.Path)

		// Send response
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{
			"symbol": "BTC_USDT",
			"orderId": 123456789,
			"clientOrderId": "test123",
			"price": "0.1",
			"origQty": "1.0",
			"executedQty": "0.0",
			"status": "CANCELED",
			"timeInForce": "GTC",
			"type": "LIMIT",
			"side": "BUY"
		}`)
	}))
	defer ts.Close()

	// Create client
	client := NewClient("test-api-key", "test-secret-key", WithBaseURL(ts.URL))

	// Test cancel order
	err := client.CancelOrder(context.Background(), "BTC_USDT", "123456789")

	// Verify response
	require.NoError(t, err)
	// CancelOrder only returns error, no order object
}

func TestGetOrderBook(t *testing.T) {
	// Sample response based on MEXC API documentation
	responseBody := `{
		"lastUpdateId": 1641183100,
		"bids": [
			["41995.0", "1.5"],
			["41990.0", "2.5"],
			["41985.0", "3.0"]
		],
		"asks": [
			["42005.0", "2.0"],
			["42010.0", "3.0"],
			["42015.0", "1.0"]
		]
	}`

	client, _, cleanup := setupTestClient(mockHandler(http.StatusOK, responseBody))
	defer cleanup()

	// Test GetOrderBook
	orderBook, err := client.GetOrderBook(context.Background(), "BTCUSDT", 10)

	// Verify results
	require.NoError(t, err)
	require.NotNil(t, orderBook)

	// Check if the order book has the correct values
	assert.Equal(t, "BTCUSDT", orderBook.Symbol)
	assert.Equal(t, int64(1641183100), orderBook.LastUpdateID)

	// Check bids
	require.Len(t, orderBook.Bids, 3)
	assert.Equal(t, 41995.0, orderBook.Bids[0].Price)
	assert.Equal(t, 1.5, orderBook.Bids[0].Quantity)
	assert.Equal(t, 41990.0, orderBook.Bids[1].Price)
	assert.Equal(t, 2.5, orderBook.Bids[1].Quantity)
	assert.Equal(t, 41985.0, orderBook.Bids[2].Price)
	assert.Equal(t, 3.0, orderBook.Bids[2].Quantity)

	// Check asks
	require.Len(t, orderBook.Asks, 3)
	assert.Equal(t, 42005.0, orderBook.Asks[0].Price)
	assert.Equal(t, 2.0, orderBook.Asks[0].Quantity)
	assert.Equal(t, 42010.0, orderBook.Asks[1].Price)
	assert.Equal(t, 3.0, orderBook.Asks[1].Quantity)
	assert.Equal(t, 42015.0, orderBook.Asks[2].Price)
	assert.Equal(t, 1.0, orderBook.Asks[2].Quantity)
}

func TestGetKlines(t *testing.T) {
	// Sample response based on MEXC API documentation
	responseBody := `[
		[1641182400000, "41800.0", "41900.0", "41750.0", "41850.0", "10.5", 1641186000000, "440000.0", 100, "5.5", "230000.0", "0"],
		[1641186000000, "41850.0", "42000.0", "41800.0", "42000.0", "15.0", 1641189600000, "630000.0", 150, "8.0", "336000.0", "0"]
	]`

	client, _, cleanup := setupTestClient(mockHandler(http.StatusOK, responseBody))
	defer cleanup()

	// Test GetKlines
	klines, err := client.GetKlines(context.Background(), "BTCUSDT", string(model.KlineInterval1h), 2)

	// Verify results
	require.NoError(t, err)
	require.NotNil(t, klines)
	require.Len(t, klines, 2)

	// Check first kline
	assert.Equal(t, "BTCUSDT", klines[0].Symbol)
	assert.Equal(t, model.KlineInterval1h, klines[0].Interval)
	assert.Equal(t, int64(1641182400000), klines[0].OpenTime.UnixNano()/int64(time.Millisecond))
	assert.Equal(t, 41800.0, klines[0].Open)
	assert.Equal(t, 41900.0, klines[0].High)
	assert.Equal(t, 41750.0, klines[0].Low)
	assert.Equal(t, 41850.0, klines[0].Close)
	assert.Equal(t, 10.5, klines[0].Volume)
	assert.Equal(t, 440000.0, klines[0].QuoteVolume)
	assert.Equal(t, int64(100), klines[0].TradeCount)
	assert.True(t, klines[0].IsClosed)

	// Check second kline
	assert.Equal(t, "BTCUSDT", klines[1].Symbol)
	assert.Equal(t, model.KlineInterval1h, klines[1].Interval)
	assert.Equal(t, int64(1641186000000), klines[1].OpenTime.UnixNano()/int64(time.Millisecond))
	assert.Equal(t, 41850.0, klines[1].Open)
	assert.Equal(t, 42000.0, klines[1].High)
	assert.Equal(t, 41800.0, klines[1].Low)
	assert.Equal(t, 42000.0, klines[1].Close)
	assert.Equal(t, 15.0, klines[1].Volume)
	assert.Equal(t, 630000.0, klines[1].QuoteVolume)
	assert.Equal(t, int64(150), klines[1].TradeCount)
	assert.True(t, klines[1].IsClosed)
}

func TestErrorHandling(t *testing.T) {
	// Test error response
	errorResponse := `{"code": 400, "msg": "Invalid parameter"}`
	client, _, cleanup := setupTestClient(mockHandler(http.StatusBadRequest, errorResponse))
	defer cleanup()

	// Test GetMarketData with error
	_, err := client.GetMarketData(context.Background(), "INVALID")

	// Verify error
	require.Error(t, err)
	assert.Contains(t, err.Error(), "API error")
}
