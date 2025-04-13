package rest

import (
	"context"
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
	wallet, err := client.GetAccount(context.Background())

	// Verify results
	require.NoError(t, err)
	require.NotNil(t, wallet)

	// Check if wallet has the correct balances
	assert.Len(t, wallet.Balances, 3)
	assert.Equal(t, 0.1, wallet.Balances[model.Asset("BTC")].Free)
	assert.Equal(t, 0.05, wallet.Balances[model.Asset("BTC")].Locked)
	assert.InDelta(t, 0.15, wallet.Balances[model.Asset("BTC")].Total, 0.000001)

	assert.Equal(t, 2.5, wallet.Balances[model.Asset("ETH")].Free)
	assert.Equal(t, 0.0, wallet.Balances[model.Asset("ETH")].Locked)
	assert.Equal(t, 2.5, wallet.Balances[model.Asset("ETH")].Total)

	assert.Equal(t, 1000.0, wallet.Balances[model.Asset("USDT")].Free)
	assert.Equal(t, 500.0, wallet.Balances[model.Asset("USDT")].Locked)
	assert.Equal(t, 1500.0, wallet.Balances[model.Asset("USDT")].Total)
}

func TestGetMarketData(t *testing.T) {
	// Sample response based on MEXC API documentation
	responseBody := `{
		"symbol": "BTCUSDT",
		"priceChange": "100.0",
		"priceChangePercent": "2.5",
		"weightedAvgPrice": "42000.0",
		"prevClosePrice": "41900.0",
		"lastPrice": "42000.0",
		"lastQty": "0.01",
		"bidPrice": "41995.0",
		"bidQty": "1.5",
		"askPrice": "42005.0",
		"askQty": "2.0",
		"openPrice": "41800.0",
		"highPrice": "42200.0",
		"lowPrice": "41700.0",
		"volume": "100.0",
		"quoteVolume": "4200000.0",
		"openTime": 1641182500000,
		"closeTime": 1641182800000,
		"count": 5000
	}`

	client, _, cleanup := setupTestClient(mockHandler(http.StatusOK, responseBody))
	defer cleanup()

	// Test GetMarketData
	ticker, err := client.GetMarketData(context.Background(), "BTCUSDT")

	// Verify results
	require.NoError(t, err)
	require.NotNil(t, ticker)

	// Check if ticker has the correct values
	assert.Equal(t, "BTCUSDT", ticker.Symbol)
	assert.Equal(t, 42000.0, ticker.LastPrice)
	assert.Equal(t, 100.0, ticker.PriceChange)
	assert.Equal(t, 2.5, ticker.PriceChangePercent)
	assert.Equal(t, 42200.0, ticker.HighPrice)
	assert.Equal(t, 41700.0, ticker.LowPrice)
	assert.Equal(t, 100.0, ticker.Volume)
	assert.Equal(t, 4200000.0, ticker.QuoteVolume)
	assert.Equal(t, 41995.0, ticker.BidPrice)
	assert.Equal(t, 1.5, ticker.BidQty)
	assert.Equal(t, 42005.0, ticker.AskPrice)
	assert.Equal(t, 2.0, ticker.AskQty)
	assert.Equal(t, int64(5000), ticker.Count)
}

func TestPlaceOrder(t *testing.T) {
	// Sample response based on MEXC API documentation
	responseBody := `{
		"symbol": "BTCUSDT",
		"orderId": "123456789",
		"clientOrderId": "client123",
		"transactTime": 1641182900000,
		"price": "42000.0",
		"origQty": "0.01",
		"executedQty": "0.0",
		"status": "NEW",
		"timeInForce": "GTC",
		"type": "LIMIT",
		"side": "BUY"
	}`

	client, _, cleanup := setupTestClient(mockHandler(http.StatusOK, responseBody))
	defer cleanup()

	// Test PlaceOrder with individual parameters
	placedOrder, err := client.PlaceOrder(
		context.Background(),
		"BTCUSDT",
		model.OrderSideBuy,
		model.OrderTypeLimit,
		0.01,
		42000.0,
		model.TimeInForceGTC,
	)

	// Verify results
	require.NoError(t, err)
	require.NotNil(t, placedOrder)

	// Check if the order has the correct values
	assert.Equal(t, "123456789", placedOrder.OrderID)
	assert.Equal(t, "client123", placedOrder.ClientOrderID)
	assert.Equal(t, "BTCUSDT", placedOrder.Symbol)
	assert.Equal(t, model.OrderSideBuy, placedOrder.Side)
	assert.Equal(t, model.OrderTypeLimit, placedOrder.Type)
	assert.Equal(t, model.OrderStatusNew, placedOrder.Status)
	assert.Equal(t, model.TimeInForceGTC, placedOrder.TimeInForce)
	assert.Equal(t, 42000.0, placedOrder.Price)
	assert.Equal(t, 0.01, placedOrder.Quantity)
	assert.Equal(t, 0.0, placedOrder.ExecutedQty)
	assert.True(t, placedOrder.CreatedAt.Unix() > 0)
	assert.True(t, placedOrder.UpdatedAt.Unix() > 0)

	// Verify the order is not complete
	assert.False(t, placedOrder.IsComplete())

	// Verify the remaining quantity
	assert.Equal(t, 0.01, placedOrder.RemainingQuantity())
}

func TestGetOrderStatus(t *testing.T) {
	// Sample response based on MEXC API documentation
	responseBody := `{
		"symbol": "BTCUSDT",
		"orderId": "123456789",
		"clientOrderId": "client123",
		"price": "42000.0",
		"origQty": "0.01",
		"executedQty": "0.005",
		"status": "PARTIALLY_FILLED",
		"timeInForce": "GTC",
		"type": "LIMIT",
		"side": "BUY",
		"time": 1641182900000,
		"updateTime": 1641183000000
	}`

	client, _, cleanup := setupTestClient(mockHandler(http.StatusOK, responseBody))
	defer cleanup()

	// Test GetOrderStatus with individual parameters
	order, err := client.GetOrderStatus(context.Background(), "BTCUSDT", "123456789")

	// Verify results
	require.NoError(t, err)
	require.NotNil(t, order)

	// Check if the order has the correct values
	assert.Equal(t, "123456789", order.OrderID)
	assert.Equal(t, "client123", order.ClientOrderID)
	assert.Equal(t, "BTCUSDT", order.Symbol)
	assert.Equal(t, model.OrderSideBuy, order.Side)
	assert.Equal(t, model.OrderTypeLimit, order.Type)
	assert.Equal(t, model.OrderStatusPartiallyFilled, order.Status)
	assert.Equal(t, model.TimeInForceGTC, order.TimeInForce)
	assert.Equal(t, 42000.0, order.Price)
	assert.Equal(t, 0.01, order.Quantity)
	assert.Equal(t, 0.005, order.ExecutedQty)

	// Verify the order is not complete
	assert.False(t, order.IsComplete())

	// Verify the remaining quantity
	assert.Equal(t, 0.005, order.RemainingQuantity())
}

func TestCancelOrder(t *testing.T) {
	// Sample response for a successful cancellation
	responseBody := `{
		"symbol": "BTCUSDT",
		"orderId": "123456789",
		"clientOrderId": "client123",
		"status": "CANCELED"
	}`

	client, _, cleanup := setupTestClient(mockHandler(http.StatusOK, responseBody))
	defer cleanup()

	// Test CancelOrder
	err := client.CancelOrder(context.Background(), "BTCUSDT", "123456789")

	// Verify results
	require.NoError(t, err)
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
	klines, err := client.GetKlines(context.Background(), "BTCUSDT", model.KlineInterval1h, 2)

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
