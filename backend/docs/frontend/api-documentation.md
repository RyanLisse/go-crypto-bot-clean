# Crypto Trading Bot API Documentation

This document provides comprehensive documentation for the Crypto Trading Bot API, including authentication, endpoints, request/response formats, and WebSocket integration.

## Table of Contents

1. [Authentication](#authentication)
2. [API Endpoints](#api-endpoints)
   - [Health Check](#health-check)
   - [Authentication](#authentication-endpoints)
   - [Portfolio](#portfolio-endpoints)
   - [Trading](#trading-endpoints)
   - [New Coins](#new-coins-endpoints)
   - [Configuration](#configuration-endpoints)
   - [Status](#status-endpoints)
   - [AI Assistant](#ai-assistant-endpoints)
3. [WebSocket API](#websocket-api)
4. [Error Handling](#error-handling)
5. [Rate Limiting](#rate-limiting)
6. [CORS Support](#cors-support)

## Authentication

The API uses JWT (JSON Web Token) for authentication. To access protected endpoints, you need to:

1. Obtain a JWT token by authenticating with the `/auth/login` endpoint
2. Include the token in the `Authorization` header of subsequent requests

### Token Format

```
Authorization: Bearer <token>
```

### Token Expiration

Tokens expire after the time specified in the configuration (default: 24 hours). When a token expires, you need to obtain a new one by authenticating again.

## API Endpoints

All API endpoints are prefixed with `/api/v1` unless otherwise specified.

### Health Check

#### GET /health

Check if the API is running.

**Response:**

```json
{
  "status": "ok",
  "timestamp": "2023-06-01T12:00:00Z"
}
```

### Authentication Endpoints

#### POST /auth/login

Authenticate a user and get a JWT token.

**Request:**

```json
{
  "username": "admin",
  "password": "admin123"
}
```

**Response:**

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_at": "2023-06-02T12:00:00Z",
  "user_id": "admin",
  "role": "admin"
}
```

#### POST /auth/logout

Log out the current user (clears the authentication cookie if used).

**Response:**

```json
{
  "message": "Successfully logged out"
}
```

#### GET /auth/me

Get information about the currently authenticated user.

**Response:**

```json
{
  "user_id": "admin",
  "role": "admin"
}
```

### Portfolio Endpoints

#### GET /api/v1/portfolio

Get a summary of the user's portfolio.

**Response:**

```json
{
  "total_value": 10000.50,
  "active_trades": 5,
  "profit_loss": 500.25,
  "profit_loss_percent": 5.25,
  "holdings": [
    {
      "symbol": "BTCUSDT",
      "quantity": 0.1,
      "value": 5000.0,
      "profit_loss": 250.0,
      "profit_loss_percent": 5.0
    },
    {
      "symbol": "ETHUSDT",
      "quantity": 1.0,
      "value": 3000.0,
      "profit_loss": 150.0,
      "profit_loss_percent": 5.0
    }
  ]
}
```

#### GET /api/v1/portfolio/active

Get a list of active trades.

**Response:**

```json
{
  "active_trades": [
    {
      "id": 1,
      "symbol": "BTCUSDT",
      "buy_price": 50000.0,
      "current_price": 55000.0,
      "quantity": 0.1,
      "bought_at": "2023-05-31T12:00:00Z",
      "profit_loss": 500.0,
      "profit_loss_percent": 10.0
    },
    {
      "id": 2,
      "symbol": "ETHUSDT",
      "buy_price": 3000.0,
      "current_price": 3200.0,
      "quantity": 1.0,
      "bought_at": "2023-05-31T18:00:00Z",
      "profit_loss": 200.0,
      "profit_loss_percent": 6.67
    }
  ]
}
```

#### GET /api/v1/portfolio/performance

Get performance metrics for the portfolio.

**Query Parameters:**

- `time_range`: Time range for metrics (7d, 30d, 90d, 1y, all). Default: 30d

**Response:**

```json
{
  "total_trades": 10,
  "winning_trades": 7,
  "losing_trades": 3,
  "win_rate": 0.7,
  "total_profit_loss": 1000.0,
  "average_profit_per_trade": 100.0,
  "largest_profit": 500.0,
  "largest_loss": -200.0,
  "time_range": "30d"
}
```

#### GET /api/v1/portfolio/value

Get the total value of the portfolio.

**Response:**

```json
{
  "value": 10000.50,
  "timestamp": "2023-06-01T12:00:00Z"
}
```

### Trading Endpoints

#### GET /api/v1/trade/history

Get trade history.

**Query Parameters:**

- `limit`: Maximum number of trades to return. Default: 50
- `start_time`: Start time for filtering trades (ISO 8601 format). Optional

**Response:**

```json
{
  "trades": [
    {
      "id": "1",
      "order_id": "MEXC123456",
      "symbol": "BTCUSDT",
      "side": "BUY",
      "type": "MARKET",
      "quantity": 0.1,
      "price": 50000.0,
      "status": "FILLED",
      "created_at": "2023-05-30T12:00:00Z",
      "time": "2023-05-30T12:00:00Z"
    },
    {
      "id": "2",
      "order_id": "MEXC123457",
      "symbol": "ETHUSDT",
      "side": "BUY",
      "type": "MARKET",
      "quantity": 1.0,
      "price": 3000.0,
      "status": "FILLED",
      "created_at": "2023-05-30T18:00:00Z",
      "time": "2023-05-30T18:00:00Z"
    }
  ],
  "count": 2
}
```

#### POST /api/v1/trade/buy

Execute a buy trade.

**Request:**

```json
{
  "symbol": "BTCUSDT",
  "quantity": 0.1,
  "order_type": "MARKET",
  "price": 50000.0  // Only required for LIMIT orders
}
```

**Response:**

```json
{
  "id": "3",
  "order_id": "MEXC123458",
  "symbol": "BTCUSDT",
  "side": "BUY",
  "type": "MARKET",
  "quantity": 0.1,
  "price": 50000.0,
  "status": "FILLED",
  "created_at": "2023-06-01T12:00:00Z",
  "time": "2023-06-01T12:00:00Z"
}
```

#### POST /api/v1/trade/sell

Sell a coin.

**Request:**

```json
{
  "symbol": "BTCUSDT",
  "quantity": 0.1,
  "order_type": "MARKET",
  "price": 55000.0  // Only required for LIMIT orders
}
```

**Response:**

```json
{
  "id": "4",
  "order_id": "MEXC123459",
  "symbol": "BTCUSDT",
  "side": "SELL",
  "type": "MARKET",
  "quantity": 0.1,
  "price": 55000.0,
  "status": "FILLED",
  "created_at": "2023-06-01T18:00:00Z",
  "time": "2023-06-01T18:00:00Z"
}
```

#### GET /api/v1/trade/status/:id

Get the status of a trade.

**Path Parameters:**

- `id`: Trade ID

**Response:**

```json
{
  "id": "3",
  "order_id": "MEXC123458",
  "symbol": "BTCUSDT",
  "side": "BUY",
  "type": "MARKET",
  "quantity": 0.1,
  "price": 50000.0,
  "status": "FILLED",
  "created_at": "2023-06-01T12:00:00Z",
  "time": "2023-06-01T12:00:00Z"
}
```

### New Coins Endpoints

#### GET /api/v1/newcoins

Get a list of newly detected coins.

**Query Parameters:**

- `processed`: Filter by processed status (true/false). Optional

**Response:**

```json
{
  "coins": [
    {
      "id": 1,
      "symbol": "NEWUSDT",
      "found_at": "2023-06-01T12:00:00Z",
      "base_volume": 1000.0,
      "quote_volume": 1000000.0,
      "is_processed": false
    },
    {
      "id": 2,
      "symbol": "NEWERUSDT",
      "found_at": "2023-06-01T18:00:00Z",
      "base_volume": 2000.0,
      "quote_volume": 2000000.0,
      "is_processed": false
    }
  ],
  "count": 2
}
```

#### POST /api/v1/newcoins/process

Process newly detected coins.

**Response:**

```json
{
  "processed_coins": [
    {
      "id": 1,
      "symbol": "NEWUSDT",
      "found_at": "2023-06-01T12:00:00Z",
      "base_volume": 1000.0,
      "quote_volume": 1000000.0,
      "is_processed": true
    },
    {
      "id": 2,
      "symbol": "NEWERUSDT",
      "found_at": "2023-06-01T18:00:00Z",
      "base_volume": 2000.0,
      "quote_volume": 2000000.0,
      "is_processed": true
    }
  ],
  "count": 2,
  "timestamp": "2023-06-02T12:00:00Z"
}
```

#### POST /api/v1/newcoins/detect

Trigger detection of new coins.

**Response:**

```json
{
  "coins": [
    {
      "id": 3,
      "symbol": "NEWESTUSDT",
      "found_at": "2023-06-02T12:00:00Z",
      "base_volume": 3000.0,
      "quote_volume": 3000000.0,
      "is_processed": false
    }
  ],
  "count": 1,
  "timestamp": "2023-06-02T12:00:00Z"
}
```

### Configuration Endpoints

#### GET /api/v1/config

Get the current configuration.

**Response:**

```json
{
  "trading": {
    "default_symbol": "BTCUSDT",
    "default_order_type": "MARKET",
    "default_quantity": 0.001,
    "stop_loss_percent": 5.0,
    "take_profit_levels": [5.0, 10.0, 15.0, 20.0],
    "sell_percentages": [0.25, 0.25, 0.25, 0.25]
  },
  "websocket": {
    "reconnect_delay": "5s",
    "max_reconnect_attempts": 5,
    "ping_interval": "30s",
    "auto_reconnect": true
  }
}
```

#### PUT /api/v1/config

Update the configuration.

**Request:**

```json
{
  "trading": {
    "default_symbol": "ETHUSDT",
    "default_order_type": "LIMIT",
    "default_quantity": 0.01,
    "stop_loss_percent": 10.0,
    "take_profit_levels": [5.0, 15.0, 30.0],
    "sell_percentages": [0.3, 0.3, 0.4]
  }
}
```

**Response:**

```json
{
  "message": "Configuration updated successfully",
  "config": {
    "trading": {
      "default_symbol": "ETHUSDT",
      "default_order_type": "LIMIT",
      "default_quantity": 0.01,
      "stop_loss_percent": 10.0,
      "take_profit_levels": [5.0, 15.0, 30.0],
      "sell_percentages": [0.3, 0.3, 0.4]
    },
    "websocket": {
      "reconnect_delay": "5s",
      "max_reconnect_attempts": 5,
      "ping_interval": "30s",
      "auto_reconnect": true
    }
  }
}
```

#### GET /api/v1/config/defaults

Get the default configuration.

**Response:**

```json
{
  "trading": {
    "default_symbol": "BTCUSDT",
    "default_order_type": "MARKET",
    "default_quantity": 0.001,
    "stop_loss_percent": 5.0,
    "take_profit_levels": [5.0, 10.0, 15.0, 20.0],
    "sell_percentages": [0.25, 0.25, 0.25, 0.25]
  },
  "websocket": {
    "reconnect_delay": "5s",
    "max_reconnect_attempts": 5,
    "ping_interval": "30s",
    "auto_reconnect": true
  }
}
```

### Status Endpoints

#### GET /api/v1/status

Get the system status.

**Response:**

```json
{
  "status": "running",
  "version": "1.0.0",
  "uptime": "1d 2h 34m",
  "processes": {
    "new_coin_watcher": {
      "status": "running",
      "last_run": "2023-06-01T12:00:00Z"
    },
    "position_monitor": {
      "status": "running",
      "last_run": "2023-06-01T12:05:00Z"
    }
  },
  "system_info": {
    "cpu_usage": 25.5,
    "memory_usage": 512.0,
    "disk_usage": 10.5
  }
}
```

#### POST /api/v1/status/start

Start system processes.

**Request:**

```json
{
  "processes": ["new_coin_watcher", "position_monitor"]
}
```

**Response:**

```json
{
  "message": "Processes started successfully",
  "processes": {
    "new_coin_watcher": {
      "status": "running",
      "started_at": "2023-06-02T12:00:00Z"
    },
    "position_monitor": {
      "status": "running",
      "started_at": "2023-06-02T12:00:00Z"
    }
  }
}
```

#### POST /api/v1/status/stop

Stop system processes.

**Request:**

```json
{
  "processes": ["new_coin_watcher", "position_monitor"]
}
```

**Response:**

```json
{
  "message": "Processes stopped successfully",
  "processes": {
    "new_coin_watcher": {
      "status": "stopped",
      "stopped_at": "2023-06-02T18:00:00Z"
    },
    "position_monitor": {
      "status": "stopped",
      "stopped_at": "2023-06-02T18:00:00Z"
    }
  }
}
```

### AI Assistant Endpoints

#### POST /api/v1/ai/chat

Send a message to the AI assistant and get a response.

**Request:**

```json
{
  "message": "What is my current portfolio value?",
  "conversation_id": "optional-conversation-id"
}
```

**Response:**

```json
{
  "response": "Your current portfolio value is $12,345.67, which is up 5.2% in the last 24 hours. Your largest holdings are BTC (45%), ETH (30%), and SOL (15%).",
  "conversation_id": "abc123",
  "timestamp": "2023-06-01T12:00:00Z"
}
```

#### GET /api/v1/ai/conversations

Get a list of recent AI assistant conversations.

**Response:**

```json
{
  "conversations": [
    {
      "id": "abc123",
      "title": "Portfolio Analysis",
      "created_at": "2023-06-01T10:00:00Z",
      "last_message": "What is my current portfolio value?"
    },
    {
      "id": "def456",
      "title": "Trading Strategy",
      "created_at": "2023-05-30T15:30:00Z",
      "last_message": "Explain the current trading strategy"
    }
  ],
  "count": 2
}
```

#### GET /api/v1/ai/conversations/{id}

Get the messages in a specific AI assistant conversation.

**Response:**

```json
{
  "id": "abc123",
  "title": "Portfolio Analysis",
  "messages": [
    {
      "role": "user",
      "content": "What is my current portfolio value?",
      "timestamp": "2023-06-01T10:00:00Z"
    },
    {
      "role": "assistant",
      "content": "Your current portfolio value is $12,345.67, which is up 5.2% in the last 24 hours. Your largest holdings are BTC (45%), ETH (30%), and SOL (15%).",
      "timestamp": "2023-06-01T10:00:05Z"
    }
  ],
  "created_at": "2023-06-01T10:00:00Z",
  "updated_at": "2023-06-01T10:00:05Z"
}
```

#### POST /api/v1/ai/analyze/performance

Get an AI analysis of your trading performance.

**Request:**

```json
{
  "time_range": "30d"
}
```

**Response:**

```json
{
  "analysis": "Your portfolio has performed well over the last 30 days with a 12.5% increase in value. Your most successful trades were in ETH (+18%) and SOL (+22%). Your strategy of buying during market dips has been effective, with an average return of 8.3% on these opportunistic purchases. Consider increasing your position in mid-cap altcoins which have shown strong momentum in recent weeks.",
  "key_metrics": [
    "12.5% overall portfolio growth",
    "18% return on ETH trades",
    "8.3% average return on dip-buying strategy"
  ],
  "recommendations": [
    "Increase position in mid-cap altcoins",
    "Continue dip-buying strategy",
    "Consider taking profits on SOL position"
  ],
  "timestamp": "2023-06-01T12:00:00Z"
}
```

## WebSocket API

The WebSocket API provides real-time updates for market data, trade notifications, and new coin alerts.

### Connection

Connect to the WebSocket endpoint:

```
ws://localhost:8080/ws
```

### Message Format

All WebSocket messages follow this format:

```json
{
  "type": "message_type",
  "timestamp": 1685635200,
  "payload": {
    // Message-specific data
  }
}
```

### Message Types

#### Market Data

```json
{
  "type": "market_data",
  "timestamp": 1685635200,
  "payload": {
    "symbol": "BTCUSDT",
    "price": 50000.0,
    "volume": 1000000.0,
    "timestamp": 1685635200
  }
}
```

#### Trade Notification

```json
{
  "type": "trade_notification",
  "timestamp": 1685635200,
  "payload": {
    "id": "5",
    "symbol": "BTCUSDT",
    "side": "BUY",
    "quantity": 0.1,
    "price": 50000.0,
    "timestamp": 1685635200
  }
}
```

#### New Coin Alert

```json
{
  "type": "new_coin_alert",
  "timestamp": 1685635200,
  "payload": {
    "id": 3,
    "symbol": "NEWESTUSDT",
    "found_at": 1685635200,
    "base_volume": 3000.0,
    "quote_volume": 3000000.0
  }
}
```

### Subscribing to Updates

To subscribe to ticker updates, send a message:

```json
{
  "type": "subscribe_ticker",
  "payload": {
    "symbols": ["BTCUSDT", "ETHUSDT"]
  }
}
```

### Subscription Confirmation

After subscribing, you'll receive a confirmation:

```json
{
  "type": "subscription_success",
  "timestamp": 1685635200,
  "payload": {
    "message": "Subscribed successfully"
  }
}
```

## Error Handling

The API uses standard HTTP status codes and returns error responses in a consistent format:

```json
{
  "code": "error_code",
  "message": "Human-readable error message",
  "details": "Additional error details (optional)"
}
```

### Common Error Codes

- `unauthorized`: Authentication failed or token expired
- `invalid_request`: Invalid request parameters
- `not_found`: Resource not found
- `internal_error`: Server error

## Rate Limiting

The API implements rate limiting to prevent abuse. The default limits are:

- 10 requests per second
- Burst capacity of 20 requests

When rate limits are exceeded, the API returns a 429 Too Many Requests response:

```json
{
  "code": "rate_limit_exceeded",
  "message": "Rate limit exceeded. Try again later.",
  "details": "10 requests per second allowed"
}
```

## CORS Support

The API supports Cross-Origin Resource Sharing (CORS) with the following headers:

```
Access-Control-Allow-Origin: *
Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS
Access-Control-Allow-Headers: Content-Type, Authorization
```

This allows frontend applications from any origin to interact with the API.
