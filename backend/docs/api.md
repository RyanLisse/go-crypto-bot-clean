# Crypto Bot API Documentation

## Authentication

The API uses JWT tokens for authentication, with support for Clerk authentication service.

### Authentication Methods

#### Clerk Authentication

When using Clerk (default), you'll need to:

1. Register your application with Clerk and obtain API keys
2. Configure the backend with the Clerk JWT public key and API key
3. Include the Bearer token in API requests:

```
Authorization: Bearer <token>
```

#### JWT Authentication

When using standard JWT authentication:

1. Obtain a JWT token by authenticating with the `/api/v1/auth/login` endpoint
2. Include the Bearer token in API requests:

```
Authorization: Bearer <token>
```

### Protected Endpoints

Endpoints that require authentication will return a `401 Unauthorized` response if the token is missing or invalid.

## API Endpoints

### Status Endpoints

#### Get Exchange Status

```
GET /api/v1/status/exchange
```

Returns the current status of the connected exchange (MEXC).

**Response:**

```json
{
  "name": "mexc_api",
  "status": "running",
  "message": "API connection is healthy",
  "metrics": {
    "symbols_count": 1242,
    "last_check_time": "2023-06-15T12:34:56Z",
    "response_time_ms": 345,
    "rate_limits_count": 1242
  },
  "last_update": "2023-06-15T12:34:56Z"
}
```

#### Get Services Status

```
GET /api/v1/status/services
```

Returns the status of all system services.

**Response:**

```json
{
  "version": "1.0.0",
  "status": "healthy",
  "components": [
    {
      "name": "market_data",
      "status": "running",
      "message": "Market data service is operational",
      "metrics": {
        "cache_hits": 1245,
        "cache_misses": 123,
        "response_time_ms": 45
      },
      "last_update": "2023-06-15T12:34:56Z"
    },
    {
      "name": "mexc_api",
      "status": "running",
      "message": "API connection is healthy",
      "metrics": {
        "symbols_count": 1242,
        "last_check_time": "2023-06-15T12:34:56Z",
        "response_time_ms": 345
      },
      "last_update": "2023-06-15T12:34:56Z"
    }
  ],
  "system_metrics": {
    "cpu_usage": 34.5,
    "memory_usage": 67.8,
    "disk_usage": 45.6,
    "uptime_seconds": 86400
  }
}
```

### Market Data Endpoints

#### Get Ticker

```
GET /api/v1/market/ticker?symbol=BTCUSDT&exchange=mexc
```

Returns the current ticker for the specified symbol.

**Parameters:**

- `symbol` (required): Trading pair symbol (e.g., BTCUSDT)
- `exchange` (optional): Exchange name, defaults to MEXC

**Response:**

```json
{
  "symbol": "BTCUSDT",
  "price": 42123.45,
  "high": 43000.00,
  "low": 41500.00,
  "volume": 1234.56,
  "change_percent": 2.34,
  "exchange": "MEXC",
  "timestamp": "2023-06-15T12:34:56Z"
}
```

#### Get All Tickers

```
GET /api/v1/market/tickers?exchange=mexc
```

Returns all available tickers.

**Parameters:**

- `exchange` (optional): Exchange name, defaults to MEXC

**Response:**

```json
[
  {
    "symbol": "BTCUSDT",
    "price": 42123.45,
    "high": 43000.00,
    "low": 41500.00,
    "volume": 1234.56,
    "change_percent": 2.34,
    "exchange": "MEXC",
    "timestamp": "2023-06-15T12:34:56Z"
  },
  {
    "symbol": "ETHUSDT",
    "price": 2345.67,
    "high": 2400.00,
    "low": 2300.00,
    "volume": 5678.90,
    "change_percent": 1.23,
    "exchange": "MEXC",
    "timestamp": "2023-06-15T12:34:56Z"
  }
]
```

#### Get Order Book

```
GET /api/v1/market/orderbook?symbol=BTCUSDT&exchange=mexc&depth=10
```

Returns the order book for the specified symbol.

**Parameters:**

- `symbol` (required): Trading pair symbol (e.g., BTCUSDT)
- `exchange` (optional): Exchange name, defaults to MEXC
- `depth` (optional): Order book depth, defaults to 10

**Response:**

```json
{
  "symbol": "BTCUSDT",
  "exchange": "MEXC",
  "bids": [
    { "price": 42100.00, "quantity": 1.5 },
    { "price": 42050.00, "quantity": 2.3 }
  ],
  "asks": [
    { "price": 42150.00, "quantity": 1.2 },
    { "price": 42200.00, "quantity": 3.4 }
  ],
  "last_updated": "2023-06-15T12:34:56Z"
}
```

### Alert Endpoints (Protected)

These endpoints require authentication.

#### Get Alerts

```
GET /api/v1/alerts
```

Returns all alerts for the authenticated user.

**Headers:**

```
Authorization: Bearer <token>
```

**Response:**

```json
[
  {
    "id": "d290f1ee-6c54-4b01-90e6-d701748f0851",
    "type": "price_alert",
    "symbol": "BTCUSDT",
    "condition": "above",
    "threshold": 45000.00,
    "status": "active",
    "created_at": "2023-06-10T12:34:56Z",
    "triggered_at": null
  },
  {
    "id": "e290f1ee-6c54-4b01-90e6-d701748f0852",
    "type": "price_alert",
    "symbol": "ETHUSDT",
    "condition": "below",
    "threshold": 2000.00,
    "status": "triggered",
    "created_at": "2023-06-11T12:34:56Z",
    "triggered_at": "2023-06-12T10:30:45Z"
  }
]
```

#### Create Alert

```
POST /api/v1/alerts
```

Creates a new alert for the authenticated user.

**Headers:**

```
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body:**

```json
{
  "type": "price_alert",
  "symbol": "BTCUSDT",
  "condition": "above",
  "threshold": 45000.00
}
```

**Response:**

```json
{
  "id": "d290f1ee-6c54-4b01-90e6-d701748f0851",
  "type": "price_alert",
  "symbol": "BTCUSDT",
  "condition": "above",
  "threshold": 45000.00,
  "status": "active",
  "created_at": "2023-06-15T12:34:56Z",
  "triggered_at": null
}
```

#### Delete Alert

```
DELETE /api/v1/alerts/{id}
```

Deletes an alert for the authenticated user.

**Headers:**

```
Authorization: Bearer <token>
```

**Parameters:**

- `id` (required): Alert ID

**Response:**

```
204 No Content
```

## Error Responses

All endpoints return standard error responses with the following format:

```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "Error description",
    "details": {} // Optional additional details
  }
}
```

### Common Error Codes

- `INVALID_INPUT`: The request contains invalid parameters
- `NOT_FOUND`: The requested resource was not found
- `UNAUTHORIZED`: Authentication is required or failed
- `FORBIDDEN`: The authenticated user does not have permission
- `INTERNAL_ERROR`: An internal server error occurred
- `EXTERNAL_SERVICE_ERROR`: Error communicating with external service
- `VALIDATION_ERROR`: Request validation failed
- `RATE_LIMIT`: Rate limit has been exceeded

## Pagination

Endpoints that return lists of items support pagination using the following query parameters:

- `limit`: Maximum number of items to return (default: 50, max: 100)
- `offset`: Number of items to skip (default: 0)

**Example:**

```
GET /api/v1/market/tickers?exchange=mexc&limit=10&offset=20
```

## Usage Examples

### Fetching Market Data

```bash
# Get the current Bitcoin ticker
curl -X GET "https://api.example.com/api/v1/market/ticker?symbol=BTCUSDT&exchange=mexc"

# Get order book for Ethereum
curl -X GET "https://api.example.com/api/v1/market/orderbook?symbol=ETHUSDT&exchange=mexc&depth=20"
```

### Using Authentication

```bash
# Set your auth token
TOKEN="your-auth-token"

# Create a new price alert
curl -X POST "https://api.example.com/api/v1/alerts" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "price_alert",
    "symbol": "BTCUSDT",
    "condition": "above",
    "threshold": 45000.00
  }'

# Get all your alerts
curl -X GET "https://api.example.com/api/v1/alerts" \
  -H "Authorization: Bearer $TOKEN"
``` 