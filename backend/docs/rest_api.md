# REST API Documentation

This document describes the REST API for market data.

## Base URL

All API endpoints are relative to the base URL:

```
http://<host>/api/v1
```

## Authentication

Currently, the API does not require authentication. This will be implemented in a future update.

## Rate Limiting

The API is rate-limited to protect the server from excessive requests. The current limits are:

- 60 requests per minute per IP address
- 1000 requests per day per IP address

If you exceed these limits, you will receive a `429 Too Many Requests` response.

## Endpoints

### Get Latest Tickers

Retrieves the latest ticker data for all symbols.

**Endpoint:** `GET /market/tickers`

**Parameters:**
- None

**Response:**

```json
{
  "status": "success",
  "data": [
    {
      "symbol": "BTCUSDT",
      "exchange": "mexc",
      "price": 50000.0,
      "volume": 100.0,
      "high24h": 51000.0,
      "low24h": 49000.0,
      "priceChange": 1000.0,
      "percentChange": 2.0,
      "lastUpdated": "2023-06-01T12:00:00Z"
    },
    {
      "symbol": "ETHUSDT",
      "exchange": "mexc",
      "price": 3000.0,
      "volume": 200.0,
      "high24h": 3100.0,
      "low24h": 2900.0,
      "priceChange": 100.0,
      "percentChange": 3.0,
      "lastUpdated": "2023-06-01T12:00:00Z"
    }
  ]
}
```

### Get Ticker

Retrieves the latest ticker data for a specific symbol.

**Endpoint:** `GET /market/ticker/:symbol`

**Parameters:**
- `symbol` (path parameter): The symbol to retrieve (e.g., BTCUSDT)

**Response:**

```json
{
  "status": "success",
  "data": {
    "symbol": "BTCUSDT",
    "exchange": "mexc",
    "price": 50000.0,
    "volume": 100.0,
    "high24h": 51000.0,
    "low24h": 49000.0,
    "priceChange": 1000.0,
    "percentChange": 2.0,
    "lastUpdated": "2023-06-01T12:00:00Z"
  }
}
```

### Get Candles

Retrieves candle (k-line) data for a specific symbol and interval.

**Endpoint:** `GET /market/candles/:symbol`

**Parameters:**
- `symbol` (path parameter): The symbol to retrieve (e.g., BTCUSDT)
- `interval` (query parameter): The candle interval (e.g., 1m, 5m, 15m, 30m, 1h, 4h, 1d, 1w, 1M)
- `limit` (query parameter, optional): The number of candles to retrieve (default: 100, max: 1000)
- `start` (query parameter, optional): The start time in ISO8601 format (e.g., 2023-06-01T00:00:00Z)
- `end` (query parameter, optional): The end time in ISO8601 format (e.g., 2023-06-01T12:00:00Z)

**Response:**

```json
{
  "status": "success",
  "data": [
    {
      "symbol": "BTCUSDT",
      "exchange": "mexc",
      "interval": "1h",
      "openTime": "2023-06-01T11:00:00Z",
      "closeTime": "2023-06-01T12:00:00Z",
      "open": 49500.0,
      "high": 50500.0,
      "low": 49000.0,
      "close": 50000.0,
      "volume": 100.0,
      "quoteVolume": 5000000.0,
      "tradeCount": 1000,
      "complete": true
    },
    {
      "symbol": "BTCUSDT",
      "exchange": "mexc",
      "interval": "1h",
      "openTime": "2023-06-01T10:00:00Z",
      "closeTime": "2023-06-01T11:00:00Z",
      "open": 49000.0,
      "high": 49800.0,
      "low": 48500.0,
      "close": 49500.0,
      "volume": 90.0,
      "quoteVolume": 4500000.0,
      "tradeCount": 900,
      "complete": true
    }
  ]
}
```

### Get All Symbols

Retrieves information about all available trading symbols.

**Endpoint:** `GET /market/symbols`

**Parameters:**
- None

**Response:**

```json
{
  "status": "success",
  "data": [
    {
      "symbol": "BTCUSDT",
      "baseAsset": "BTC",
      "quoteAsset": "USDT",
      "exchange": "mexc",
      "status": "TRADING",
      "minPrice": 0.01,
      "maxPrice": 100000.0,
      "pricePrecision": 2,
      "minQty": 0.0001,
      "maxQty": 1000.0,
      "qtyPrecision": 4,
      "allowedOrderTypes": ["LIMIT", "MARKET"]
    },
    {
      "symbol": "ETHUSDT",
      "baseAsset": "ETH",
      "quoteAsset": "USDT",
      "exchange": "mexc",
      "status": "TRADING",
      "minPrice": 0.01,
      "maxPrice": 10000.0,
      "pricePrecision": 2,
      "minQty": 0.001,
      "maxQty": 1000.0,
      "qtyPrecision": 3,
      "allowedOrderTypes": ["LIMIT", "MARKET"]
    }
  ]
}
```

### Get Symbol Info

Retrieves detailed information about a specific trading symbol.

**Endpoint:** `GET /market/symbol/:symbol`

**Parameters:**
- `symbol` (path parameter): The symbol to retrieve (e.g., BTCUSDT)

**Response:**

```json
{
  "status": "success",
  "data": {
    "symbol": "BTCUSDT",
    "baseAsset": "BTC",
    "quoteAsset": "USDT",
    "exchange": "mexc",
    "status": "TRADING",
    "minPrice": 0.01,
    "maxPrice": 100000.0,
    "pricePrecision": 2,
    "minQty": 0.0001,
    "maxQty": 1000.0,
    "qtyPrecision": 4,
    "allowedOrderTypes": ["LIMIT", "MARKET"]
  }
}
```

## Error Handling

If an error occurs, the API will return an error response with an appropriate HTTP status code:

```json
{
  "status": "error",
  "error": {
    "code": "not_found",
    "message": "Symbol not found"
  }
}
```

Common error codes:

- `bad_request`: The request was invalid (400)
- `not_found`: The requested resource was not found (404)
- `rate_limit_exceeded`: Rate limit exceeded (429)
- `internal_error`: An internal server error occurred (500)

## Example Usage

### cURL

```bash
# Get latest tickers
curl -X GET "http://localhost:8080/api/v1/market/tickers"

# Get ticker for BTCUSDT
curl -X GET "http://localhost:8080/api/v1/market/ticker/BTCUSDT"

# Get candles for BTCUSDT with 1h interval
curl -X GET "http://localhost:8080/api/v1/market/candles/BTCUSDT?interval=1h&limit=10"

# Get all symbols
curl -X GET "http://localhost:8080/api/v1/market/symbols"

# Get symbol info for BTCUSDT
curl -X GET "http://localhost:8080/api/v1/market/symbol/BTCUSDT"
```

### Go

```go
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type ApiResponse struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
	Error  *ApiError   `json:"error,omitempty"`
}

type ApiError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func main() {
	// Get ticker for BTCUSDT
	resp, err := http.Get("http://localhost:8080/api/v1/market/ticker/BTCUSDT")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	var apiResp ApiResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		fmt.Println("Error parsing response:", err)
		return
	}

	if apiResp.Status == "error" {
		fmt.Printf("API Error: %s - %s\n", apiResp.Error.Code, apiResp.Error.Message)
		return
	}

	fmt.Printf("Response: %+v\n", apiResp.Data)
}
```

### JavaScript

```javascript
// Get ticker for BTCUSDT
fetch('http://localhost:8080/api/v1/market/ticker/BTCUSDT')
  .then(response => response.json())
  .then(data => {
    if (data.status === 'error') {
      console.error(`API Error: ${data.error.code} - ${data.error.message}`);
      return;
    }
    console.log('Ticker data:', data.data);
  })
  .catch(error => {
    console.error('Error fetching data:', error);
  });
```
