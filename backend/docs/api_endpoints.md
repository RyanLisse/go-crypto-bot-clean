# API Endpoints Documentation

## Account Endpoints

### 1. GET /api/v1/account/wallet
**Purpose**: Retrieves the user's wallet information, including balances for all assets.
**Expected Behavior**: 
- Returns a wallet object with balances for each cryptocurrency
- Includes total USD value of all assets
- Requires authentication
- Example Response:
```json
{
  "success": true,
  "data": {
    "id": "wallet_123",
    "user_id": "user_123",
    "exchange": "MEXC",
    "balances": {
      "BTC": {
        "asset": "BTC",
        "free": 0.5,
        "locked": 0.1,
        "total": 0.6,
        "usd_value": 36000
      },
      "ETH": {
        "asset": "ETH",
        "free": 5.0,
        "locked": 1.0,
        "total": 6.0,
        "usd_value": 18000
      }
    },
    "total_usd_value": 54000,
    "last_updated": "2025-04-14T19:00:00Z"
  }
}
```

### 2. GET /api/v1/account/balance/{asset}
**Purpose**: Retrieves balance history for a specific asset over a time period.
**Parameters**:
- `asset`: The cryptocurrency asset code (e.g., BTC, ETH)
- `days` (query parameter): Number of days of history to retrieve (default: 30)
**Expected Behavior**: 
- Returns an array of balance snapshots for the specified asset
- Each snapshot includes free, locked, total amounts, and USD value
- Requires authentication
- Example Response:
```json
{
  "success": true,
  "data": [
    {
      "id": "BTC_20250414",
      "user_id": "user_123",
      "asset": "BTC",
      "free": 0.5,
      "locked": 0.1,
      "total": 0.6,
      "usd_value": 36000,
      "timestamp": "2025-04-14T00:00:00Z"
    },
    {
      "id": "BTC_20250413",
      "user_id": "user_123",
      "asset": "BTC",
      "free": 0.48,
      "locked": 0.1,
      "total": 0.58,
      "usd_value": 34800,
      "timestamp": "2025-04-13T00:00:00Z"
    }
  ]
}
```

### 3. POST /api/v1/account/refresh
**Purpose**: Triggers a refresh of the user's wallet data from the exchange.
**Expected Behavior**: 
- Fetches the latest wallet data from MEXC
- Updates the local database
- Returns a success message with timestamp
- Requires authentication
- Example Response:
```json
{
  "success": true,
  "message": "Wallet refreshed successfully",
  "timestamp": "2025-04-14T19:05:00Z"
}
```

## Market Endpoints

### 1. GET /api/v1/market/tickers
**Purpose**: Retrieves ticker information for all available trading pairs.
**Expected Behavior**: 
- Returns an array of ticker objects
- Each ticker includes symbol, price, volume, and other market data
- Does not require authentication
- Example Response:
```json
{
  "success": true,
  "data": [
    {
      "symbol": "BTCUSDT",
      "price": 60000,
      "volume": 1000,
      "high_24h": 61000,
      "low_24h": 59000,
      "change_24h": 2.5,
      "timestamp": "2025-04-14T19:05:00Z"
    },
    {
      "symbol": "ETHUSDT",
      "price": 3000,
      "volume": 5000,
      "high_24h": 3100,
      "low_24h": 2900,
      "change_24h": 1.8,
      "timestamp": "2025-04-14T19:05:00Z"
    }
  ]
}
```

### 2. GET /api/v1/market/ticker/{symbol}
**Purpose**: Retrieves ticker information for a specific trading pair.
**Parameters**:
- `symbol`: The trading pair symbol (e.g., BTCUSDT, ETHUSDT)
**Expected Behavior**: 
- Returns a ticker object for the specified symbol
- Includes price, volume, and other market data
- Does not require authentication
- Example Response:
```json
{
  "success": true,
  "data": {
    "symbol": "BTCUSDT",
    "price": 60000,
    "volume": 1000,
    "high_24h": 61000,
    "low_24h": 59000,
    "change_24h": 2.5,
    "timestamp": "2025-04-14T19:05:00Z"
  }
}
```

## Test Endpoints

### 1. GET /api/v1/account-test
**Purpose**: Simple test endpoint to verify account API functionality.
**Expected Behavior**: 
- Returns a success message
- Does not require authentication
- Example Response:
```json
{
  "success": true,
  "message": "Account test endpoint works!"
}
```

### 2. GET /api/v1/account-wallet-test
**Purpose**: Test endpoint that returns mock wallet data.
**Expected Behavior**: 
- Returns a mock wallet object with predefined balances
- Does not require authentication
- Example Response:
```json
{
  "success": true,
  "data": {
    "id": "wallet_test",
    "user_id": "test_user",
    "exchange": "MEXC",
    "balances": {
      "BTC": {
        "asset": "BTC",
        "free": 0.5,
        "locked": 0.1,
        "total": 0.6,
        "usd_value": 36000
      },
      "ETH": {
        "asset": "ETH",
        "free": 5.0,
        "locked": 1.0,
        "total": 6.0,
        "usd_value": 18000
      }
    },
    "total_usd_value": 54000,
    "last_updated": "2025-04-14T19:00:00Z"
  }
}
```

## Subtasks for API Endpoint Route Registration

### Subtask 1: Fix Account Endpoint Registration
- Review account handler implementation
- Implement direct account endpoints in main.go
- Test account endpoints to ensure they're working correctly
- Ensure proper authentication middleware is applied

### Subtask 2: Fix Market Endpoint Registration
- Review market handler implementation
- Ensure market endpoints are properly registered
- Test market endpoints to ensure they're working correctly
- Implement caching for market data

### Subtask 3: Implement Test Endpoints
- Create direct test endpoints for debugging
- Implement mock data endpoints for testing
- Create test scripts to verify endpoint functionality
- Document test endpoints for development use
