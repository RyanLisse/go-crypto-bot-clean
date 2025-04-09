# MEXC API Research and Implementation Guide

## Executive Summary

This document provides comprehensive research on the MEXC cryptocurrency exchange API, focusing on implementation details for the Go crypto bot project. It covers both REST and WebSocket APIs for spot and futures trading, with specific attention to Go language implementation patterns, error handling, and best practices for automated trading systems.

## 1. API Overview

### 1.1 Key Components
- **REST API**: For account data, order management, and market information
- **WebSocket API**: For real-time data streams and subscription-based updates
- **Dual Market Support**: Both spot and futures trading with distinct endpoints
- **Authentication**: HMAC-SHA256 signature-based authentication system

### 1.2 Base URLs and Endpoints

| API Type | Base URL | Purpose |
|----------|----------|----------|
| Spot REST | `https://api.mexc.com/api/v3` | Spot trading operations |
| Futures REST | `https://contract.mexc.com/api/v1` | Futures contract operations |
| Spot WebSocket | `wss://wbs.mexc.com/ws` | Real-time spot market data |
| Futures WebSocket | `wss://contract.mexc.com/ws` | Real-time futures market data |

### 1.3 Recent API Changes (As of April 2025)
- WebSocket service upgraded to use Protobuf serialization for improved efficiency
- Legacy WebSocket URL (`wss://wbs.mexc.com/ws`) will be discontinued by August 2025
- Rate limits adjusted to 5 orders per second for spot API (effective March 25, 2025)

## 2. Authentication and Security

### 2.1 API Key Management
- Generate API keys through MEXC account settings
- Restrict keys to specific IP addresses for enhanced security
- Separate trading and withdrawal permissions
- Implement key rotation policies (recommended: 90-day rotation)

### 2.2 Request Signing
```go
func createSignature(secretKey, queryString string) string {
    h := hmac.New(sha256.New, []byte(secretKey))
    h.Write([]byte(queryString))
    return hex.EncodeToString(h.Sum(nil))
}

func buildSignedRequest(apiKey, secretKey string) {
    timestamp := fmt.Sprintf("%d", time.Now().UnixMilli())
    queryParams := fmt.Sprintf("timestamp=%s", timestamp)
    signature := createSignature(secretKey, queryParams)
    fullParams := fmt.Sprintf("%s&signature=%s", queryParams, signature)
    
    req, _ := http.NewRequest("GET", "https://api.mexc.com/api/v3/account?"+fullParams, nil)
    req.Header.Set("X-MEXC-APIKEY", apiKey)
}
```

## 3. REST API Implementation

### 3.1 Account and Balance
```go
// Get account information including balances
func (c *MexcClient) GetAccount() (*AccountInfo, error) {
    endpoint := "/api/v3/account"
    timestamp := fmt.Sprintf("%d", time.Now().UnixMilli())
    queryParams := fmt.Sprintf("timestamp=%s", timestamp)
    signature := c.createSignature(queryParams)
    
    url := fmt.Sprintf("%s%s?%s&signature=%s", c.baseURL, endpoint, queryParams, signature)
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, err
    }
    
    req.Header.Set("X-MEXC-APIKEY", c.apiKey)
    // Execute request and parse response
    // ...
}
```

### 3.2 Order Management
```go
// Place a new order
func (c *MexcClient) PlaceOrder(symbol, side, orderType string, quantity float64, price float64) (*OrderResponse, error) {
    endpoint := "/api/v3/order"
    timestamp := fmt.Sprintf("%d", time.Now().UnixMilli())
    
    params := map[string]string{
        "symbol": symbol,
        "side": side,
        "type": orderType,
        "quantity": fmt.Sprintf("%f", quantity),
        "timestamp": timestamp,
    }
    
    if price > 0 && orderType == "LIMIT" {
        params["price"] = fmt.Sprintf("%f", price)
    }
    
    // Build query string and signature
    // ...
    
    // Execute POST request
    // ...
}
```

### 3.3 Market Data
```go
// Get kline/candlestick data
func (c *MexcClient) GetKlines(symbol, interval string, limit int) ([]Kline, error) {
    endpoint := "/api/v3/klines"
    params := fmt.Sprintf("symbol=%s&interval=%s&limit=%d", symbol, interval, limit)
    
    url := fmt.Sprintf("%s%s?%s", c.baseURL, endpoint, params)
    // Execute request and parse response
    // ...
}
```

## 4. WebSocket API Implementation

### 4.1 Connection Management
```go
type MexcWebSocketClient struct {
    conn              *websocket.Conn
    url               string
    subscribedStreams map[string]bool
    handlers          map[string]func([]byte)
    mutex             sync.Mutex
    done              chan struct{}
    reconnectDelay    time.Duration
    maxReconnects     int
    connectionAttempts int
}

func NewWebSocketClient(url string) *MexcWebSocketClient {
    return &MexcWebSocketClient{
        url:               url,
        subscribedStreams: make(map[string]bool),
        handlers:          make(map[string]func([]byte)),
        done:              make(chan struct{}),
        reconnectDelay:    5 * time.Second,
        maxReconnects:     10,
    }
}

func (c *MexcWebSocketClient) Connect() error {
    c.mutex.Lock()
    c.connectionAttempts++
    c.mutex.Unlock()
    
    dialer := websocket.Dialer{}
    conn, _, err := dialer.Dial(c.url, nil)
    if err != nil {
        return err
    }
    
    c.conn = conn
    
    // Start message handling goroutine
    go c.handleMessages()
    
    // Resubscribe to streams if reconnecting
    if len(c.subscribedStreams) > 0 {
        c.resubscribe()
    }
    
    return nil
}
```

### 4.2 Subscription Management
```go
func (c *MexcWebSocketClient) Subscribe(stream string, handler func([]byte)) error {
    c.mutex.Lock()
    defer c.mutex.Unlock()
    
    if len(c.subscribedStreams) >= 30 {
        return errors.New("maximum number of subscriptions reached (30)")
    }
    
    msg := map[string]interface{}{
        "method": "SUBSCRIBE",
        "params": []string{stream},
        "id": time.Now().Unix(),
    }
    
    data, err := json.Marshal(msg)
    if err != nil {
        return err
    }
    
    if err := c.conn.WriteMessage(websocket.TextMessage, data); err != nil {
        return err
    }
    
    c.subscribedStreams[stream] = true
    c.handlers[stream] = handler
    
    return nil
}
```

### 4.3 Message Handling
```go
func (c *MexcWebSocketClient) handleMessages() {
    defer func() {
        c.conn.Close()
    }()
    
    // Setup ping handler
    c.conn.SetPingHandler(func(data string) error {
        return c.conn.WriteControl(websocket.PongMessage, []byte{}, time.Now().Add(10*time.Second))
    })
    
    for {
        select {
        case <-c.done:
            return
        default:
            _, message, err := c.conn.ReadMessage()
            if err != nil {
                // Handle reconnection logic
                c.handleReconnect()
                return
            }
            
            // Process message
            c.processMessage(message)
        }
    }
}

func (c *MexcWebSocketClient) processMessage(message []byte) {
    // Parse message to determine stream
    var msg map[string]interface{}
    if err := json.Unmarshal(message, &msg); err != nil {
        return
    }
    
    // Check if it's a ping message
    if msg["method"] == "PING" {
        c.handlePing()
        return
    }
    
    // Handle stream data
    if stream, ok := msg["stream"].(string); ok {
        if handler, exists := c.handlers[stream]; exists {
            handler(message)
        }
    }
}
```

## 5. Futures API Implementation

### 5.1 Futures-Specific Endpoints
```go
// Get futures positions
func (c *MexcFuturesClient) GetPositions(symbol string) ([]Position, error) {
    endpoint := "/api/v1/private/position/list"
    timestamp := fmt.Sprintf("%d", time.Now().UnixMilli())
    
    params := map[string]string{
        "symbol": symbol,
        "timestamp": timestamp,
    }
    
    // Build query string and signature
    // ...
    
    // Execute request and parse response
    // ...
}

// Place futures order
func (c *MexcFuturesClient) PlaceFuturesOrder(symbol string, side int, orderType int, 
                                           price float64, volume float64) (*FuturesOrderResponse, error) {
    endpoint := "/api/v1/private/order/submit"
    timestamp := fmt.Sprintf("%d", time.Now().UnixMilli())
    
    params := map[string]interface{}{
        "symbol": symbol,
        "price": price,
        "vol": volume,
        "side": side,        // 1: open long, 2: close short, 3: open short, 4: close long
        "type": orderType,   // 1: limit order, 2: post only, 3: IOC, 4: FOK, 5: market order
        "timestamp": timestamp,
    }
    
    // Build request body and signature
    // ...
    
    // Execute POST request
    // ...
}
```

## 6. Error Handling and Rate Limiting

### 6.1 Error Response Structure
```go
type MexcError struct {
    Code    int    `json:"code"`
    Message string `json:"msg"`
}

func (e *MexcError) Error() string {
    return fmt.Sprintf("MEXC API error: code=%d, message=%s", e.Code, e.Message)
}
```

### 6.2 Common Error Codes

| Code | Description | Handling Strategy |
|------|-------------|-------------------|
| 400 | Bad request | Fix request parameters |
| 401 | Unauthorized | Check API key and signature |
| 429 | Rate limit exceeded | Implement backoff and retry |
| 500 | Internal server error | Retry with exponential backoff |
| 513 | Timestamp drift | Synchronize with server time |

### 6.3 Rate Limit Implementation
```go
type RateLimiter struct {
    requests     int
    period       time.Duration
    lastRefill   time.Time
    tokens       int
    mutex        sync.Mutex
}

func NewRateLimiter(requests int, period time.Duration) *RateLimiter {
    return &RateLimiter{
        requests:   requests,
        period:     period,
        lastRefill: time.Now(),
        tokens:     requests,
    }
}

func (r *RateLimiter) Allow() bool {
    r.mutex.Lock()
    defer r.mutex.Unlock()
    
    now := time.Now()
    elapsed := now.Sub(r.lastRefill)
    
    // Refill tokens based on elapsed time
    if elapsed >= r.period {
        r.tokens = r.requests
        r.lastRefill = now
    } else if elapsed > 0 {
        newTokens := int(float64(r.requests) * (float64(elapsed) / float64(r.period)))
        if newTokens > 0 {
            r.tokens = min(r.tokens+newTokens, r.requests)
            r.lastRefill = r.lastRefill.Add(elapsed)
        }
    }
    
    if r.tokens <= 0 {
        return false
    }
    
    r.tokens--
    return true
}
```

## 7. Integration with Go Crypto Bot Architecture

### 7.1 Client Interface
```go
type ExchangeClient interface {
    // Account methods
    GetAccount() (*AccountInfo, error)
    GetBalances() (map[string]Balance, error)
    
    // Market data methods
    GetTicker(symbol string) (*Ticker, error)
    GetKlines(symbol, interval string, limit int) ([]Kline, error)
    GetOrderBook(symbol string, limit int) (*OrderBook, error)
    
    // Trading methods
    PlaceOrder(symbol, side, orderType string, quantity float64, price float64) (*OrderResponse, error)
    CancelOrder(symbol, orderId string) error
    GetOrder(symbol, orderId string) (*Order, error)
    GetOpenOrders(symbol string) ([]Order, error)
    
    // WebSocket methods
    SubscribeToTicker(symbol string, handler func(*Ticker))
    SubscribeToKlines(symbol, interval string, handler func([]Kline))
    SubscribeToTrades(symbol string, handler func([]Trade))
    SubscribeToOrderBook(symbol string, handler func(*OrderBook))
    SubscribeToUserData(handler func(interface{}))
}
```

### 7.2 Implementation for MEXC
```go
type MexcClient struct {
    baseURL     string
    apiKey      string
    secretKey   string
    httpClient  *http.Client
    rateLimiter *RateLimiter
    wsClient    *MexcWebSocketClient
}

func NewMexcClient(apiKey, secretKey string) *MexcClient {
    client := &MexcClient{
        baseURL:    "https://api.mexc.com/api/v3",
        apiKey:     apiKey,
        secretKey:  secretKey,
        httpClient: &http.Client{Timeout: 10 * time.Second},
        rateLimiter: NewRateLimiter(60, time.Minute), // 60 requests per minute
    }
    
    client.wsClient = NewWebSocketClient("wss://wbs.mexc.com/ws")
    client.wsClient.Connect()
    
    return client
}

// Implement interface methods
// ...
```

## 8. Best Practices for Production Use

### 8.1 Security Best Practices
- Store API keys in environment variables or secure vaults (HashiCorp Vault, AWS Secrets Manager)
- Implement IP whitelisting for API access
- Use read-only API keys when possible
- Rotate API keys regularly (every 90 days)
- Implement request signing verification

### 8.2 Reliability Patterns
- Implement circuit breakers for API calls
- Use exponential backoff for retries
- Maintain connection health checks for WebSockets
- Implement graceful degradation for non-critical features

```go
func (c *MexcClient) executeRequestWithRetry(req *http.Request) (*http.Response, error) {
    var resp *http.Response
    var err error
    
    backoff := 500 * time.Millisecond
    maxRetries := 3
    
    for i := 0; i < maxRetries; i++ {
        if !c.rateLimiter.Allow() {
            time.Sleep(100 * time.Millisecond)
            continue
        }
        
        resp, err = c.httpClient.Do(req)
        if err != nil {
            time.Sleep(backoff)
            backoff *= 2 // Exponential backoff
            continue
        }
        
        if resp.StatusCode == 429 {
            retryAfter := resp.Header.Get("Retry-After")
            if retryAfter != "" {
                if seconds, err := strconv.Atoi(retryAfter); err == nil {
                    time.Sleep(time.Duration(seconds) * time.Second)
                }
            } else {
                time.Sleep(backoff)
                backoff *= 2
            }
            continue
        }
        
        return resp, nil
    }
    
    return resp, err
}
```

### 8.3 Performance Optimization
- Cache frequently accessed data (tickers, order books)
- Use connection pooling for HTTP requests
- Implement request batching where supported
- Use goroutines and channels for concurrent processing

## 9. Testing and Validation

### 9.1 Unit Testing Example
```go
func TestMexcClient_GetAccount(t *testing.T) {
    // Setup mock server
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Verify request headers and parameters
        assert.Equal(t, "GET", r.Method)
        assert.Contains(t, r.URL.Path, "/api/v3/account")
        assert.Contains(t, r.Header.Get("X-MEXC-APIKEY"), "testkey")
        
        // Return mock response
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        fmt.Fprint(w, `{
            "makerCommission": 0,
            "takerCommission": 0,
            "buyerCommission": 0,
            "sellerCommission": 0,
            "canTrade": true,
            "canWithdraw": true,
            "canDeposit": true,
            "updateTime": 1617939135818,
            "accountType": "SPOT",
            "balances": [
                {
                    "asset": "BTC",
                    "free": "0.00000000",
                    "locked": "0.00000000"
                },
                {
                    "asset": "USDT",
                    "free": "1000.00000000",
                    "locked": "0.00000000"
                }
            ]
        }`)
    }))
    defer server.Close()
    
    // Create client with mock server URL
    client := NewMexcClient("testkey", "testsecret")
    client.baseURL = server.URL
    
    // Test the method
    account, err := client.GetAccount()
    assert.NoError(t, err)
    assert.NotNil(t, account)
    assert.Equal(t, true, account.CanTrade)
    assert.Equal(t, 2, len(account.Balances))
    assert.Equal(t, "USDT", account.Balances[1].Asset)
    assert.Equal(t, "1000.00000000", account.Balances[1].Free)
}
```

### 9.2 Integration Testing Strategy
- Use MEXC's testnet environment when available
- Implement sandbox mode for testing without real orders
- Create dedicated test accounts with minimal funds
- Test rate limiting behavior with concurrent requests
- Validate WebSocket reconnection logic

## 10. Conclusion and Recommendations

### 10.1 Key Findings
- MEXC API provides comprehensive endpoints for both spot and futures trading
- WebSocket implementation offers efficient real-time data access
- Rate limiting requires careful implementation for production use
- Error handling needs to account for various response codes and formats

### 10.2 Implementation Recommendations
- Use the interface-based approach for easy exchange abstraction
- Implement robust error handling and retry mechanisms
- Leverage WebSockets for real-time data to reduce API load
- Maintain connection health with regular ping/pong messages
- Implement proper rate limiting to avoid IP bans

### 10.3 Next Steps
- Develop comprehensive test suite for all API endpoints
- Implement logging and monitoring for API interactions
- Create configuration options for different environments
- Document API usage patterns for future development

## 11. References

- [MEXC API Documentation](https://mexcdevelop.github.io/apidocs/)
- [MEXC Futures API Documentation](https://mexcdevelop.github.io/apidocs/contract_v1_en/)
- [MEXC WebSocket API Documentation](https://mexcdevelop.github.io/apidocs/spot_v3_en/#websocket-market-streams)
- [MEXC Rate Limit Adjustments (2025)](https://www.mexc.com/en-GB/support/articles/17827791522801)

---

### Survey Note: Detailed Analysis of MEXC API for App Development

This note provides a comprehensive overview of the MEXC API, focusing on REST and WebSocket endpoints to support an app that offers users access to account balance, trading history, profit or loss, and real-time data. The analysis is based on the official documentation available at [MEXC API Docs](https://mexcdevelop.github.io/apidocs).

#### Overview of MEXC API Endpoints
The MEXC API provides both REST and WebSocket interfaces, catering to different needs for data retrieval and real-time updates. REST endpoints are suitable for fetching static or historical data, while WebSocket endpoints enable live streaming, essential for dynamic app features.

#### REST Endpoints for Core Features
For the app's core functionalities, several REST endpoints are critical:

- **Account Balance:** The endpoint `/api/v3/account` retrieves current account information, including balances for all assets. This is essential for displaying the user's current holdings. The response typically includes fields like `canTrade`, `canWithdraw`, and a list of balances, as noted in the documentation.

- **Trading History:** The endpoint `/api/v3/myTrades` provides a list of past trades, limited to the last month by default, with options to filter by symbol, order ID, and time range. This is crucial for users to review their trading activities. The response includes details like symbol, price, quantity, and commission, facilitating detailed analysis.

- **Orders:** For order-related data, `/api/v3/openOrders` fetches current open orders, while `/api/v3/allOrders` retrieves historical orders, useful for tracking order status and history. These endpoints support parameters like symbol and time range, enhancing flexibility.

- **Deposit and Withdrawal History:** To calculate overall profit or loss, deposit and withdrawal history are vital. The documentation indicates endpoints like `/api/v3/capital/deposit/hisrec` for deposit history and similar for withdrawals, though exact URLs suggest checking [Deposit History](https://mexcdevelop.github.io/apidocs/spot_v3_en/#deposit-history-supporting-network) and [Withdrawal History](https://mexcdevelop.github.io/apidocs/spot_v3_en/#withdraw-history-supporting-network). These allow summing up total deposits and withdrawals for portfolio valuation.

- **Current Prices for Portfolio Value:** The endpoint `/api/v3/ticker/price` provides the latest prices for symbols, necessary for calculating the current portfolio value by multiplying balances by current prices and summing in a base currency like USDT.

#### WebSocket Endpoints for Real-Time Data
For real-time updates, WebSocket is the preferred method, offering low-latency data streams. The base WebSocket endpoint is `ws://wbs-api.mexc.com/ws`, with a connection validity of 24 hours and a limit of 30 subscriptions per connection.

- **User Data Streams:** To access private data, obtain a listen key via `POST /api/v3/userDataStream`, then connect to `ws://wbs-api.mexc.com/ws?listenKey=<listenKey>`. The listen key is valid for 60 minutes and extendable via `PUT /api/v3/userDataStream`, with a maximum of 60 listen keys per user ID and 5 WebSocket connections per key, totaling 300 links.

- **Relevant Streams:** For the app, subscribe to:
  - `spot@private.account.v3.api.pb` for real-time balance updates, reflecting changes like deposits or trades.
  - `spot@private.deals.v3.api.pb` for live trade updates, essential for tracking new trades as they occur.
  - `spot@private.orders.v3.api.pb` for real-time order updates, useful for monitoring order status changes.

- **Market Data Streams:** For additional context, market streams like `spot@public.aggre.deals.v3.api.pb` for trades or `spot@public.kline.v3.api.pb` for candlestick data can enhance the app, though primarily for market analysis rather than user-specific data.

#### Calculating Profit or Loss
Profit or loss can be approached in two ways, depending on user needs:

- **Realized Profit/Loss:** Calculated from trade history using `/api/v3/myTrades`. For each sell trade, match with corresponding buy trades to compute P/L per trade, summing up for total realized P/L. This requires parsing trade details like price, quantity, and commission, which may be computationally intensive for large histories.

- **Overall Portfolio Value:** For a broader view, calculate the current portfolio value by:
  - Fetching balances via `/api/v3/account`.
  - Getting current prices via `/api/v3/ticker/price`.
  - Computing total value as sum of (balance * price) in base currency.
  - Adding total withdrawals and subtracting total deposits (from respective histories) to get overall P/L as:  
    `Overall P/L = Current Portfolio Value + Total Withdrawals - Total Deposits`.

This approach provides a holistic view, reflecting both realized and unrealized gains, suitable for a simple app interface.

#### App Feature Suggestions
Based on these endpoints, the app can include:

1. **Account Balances Section:** Display current balances, with a refresh button using `/api/v3/account` and real-time updates via WebSocket `spot@private.account.v3.api.pb`. Show total value in base currency for quick overview.

2. **Trading History Section:** List past trades from `/api/v3/myTrades`, with filters for symbol or date. Use WebSocket `spot@private.deals.v3.api.pb` for live updates, enhancing user awareness of new trades.

3. **Orders Section:** Show open orders via `/api/v3/openOrders` and historical orders via `/api/v3/allOrders`, with real-time updates via `spot@private.orders.v3.api.pb` for order status changes.

4. **Profit/Loss Section:** Offer two views:
   - Realized P/L, calculated from trade history, possibly summarized per asset.
   - Current Portfolio Value, updated with prices and balances, including deposit/withdrawal impact for overall P/L.

5. **Real-Time Market Data:** Optionally, include market streams for live prices or depth, enhancing trading decisions, though secondary to user data.

#### Implementation Considerations
- **Authentication:** All private endpoints require API keys and signatures, handled via standard HTTP headers as per [MEXC API Docs](https://mexcdevelop.github.io/apidocs).
- **Rate Limits:** REST endpoints have weights (e.g., `/api/v3/account` at weight 10), and WebSocket has a 100 times/second limit, requiring careful management to avoid disconnections or IP bans.
- **Listen Key Management:** Regularly extend listen keys (every 30 minutes recommended) to maintain WebSocket connections, given the 60-minute validity.

#### Best Way for Real-Time Data
The best approach for real-time data is using WebSocket, connecting to `ws://wbs-api.mexc.com/ws` with a listen key. Subscribe to necessary streams, handle PING/PONG messages (send `{"method": "PING"}` to receive `{"id": 0, "code": 0, "msg": "PONG"}`), and process updates to update the app state. This avoids polling, reducing latency and API load, ideal for a responsive user experience.

#### Detailed Endpoint Tables
Below are tables summarizing key REST and WebSocket endpoints for clarity:

**REST Endpoints for User Data:**

| Endpoint                     | Method | Description                                      | Permission        | Weight (IP) |
|------------------------------|--------|--------------------------------------------------|-------------------|-------------|
| /api/v3/account              | GET    | Account information, including balances          | SPOT_ACCOUNT_READ | 10          |
| /api/v3/myTrades             | GET    | Account trade list, last 1 month                | SPOT_ACCOUNT_READ | 10          |
| /api/v3/openOrders           | GET    | Current open orders                             | SPOT_DEAL_READ    | 3           |
| /api/v3/allOrders            | GET    | All orders, last 24h by default                | SPOT_DEAL_READ    | 10          |
| /api/v3/ticker/price         | GET    | Symbol price ticker                             | N/A               | 1 (symbol), 2 (all) |
| Deposit History              | GET    | Query deposit history supporting network        | SPOT_WITHDRAW_READ| 1           |
| Withdrawal History           | GET    | Query withdrawal history supporting network     | SPOT_WITHDRAW_READ| 1           |

**WebSocket Endpoints for Real-Time Updates:**

| Stream                          | Description                                      | Notes                                                                 |
|---------------------------------|--------------------------------------------------|----------------------------------------------------------------------|
| spot@private.account.v3.api.pb  | Spot account updates (balance changes)           | Requires listen key, for real-time balance updates                   |
| spot@private.deals.v3.api.pb    | Spot account deals (trades)                      | For live trade updates, essential for trading history                |
| spot@private.orders.v3.api.pb   | Spot account orders (order updates)              | For real-time order status changes, enhances order monitoring        |
| Base Endpoint                   | ws://wbs-api.mexc.com/ws                        | Connect with listen key, validity 24 hours, max 30 subscriptions     |

These tables encapsulate the endpoints necessary for building the app, ensuring comprehensive coverage of user needs.

#### Conclusion
By leveraging the MEXC API's REST and WebSocket capabilities, the app can provide a robust platform for users to monitor account balances, review trading history, calculate profit or loss, and receive real-time updates. The integration of deposit/withdrawal history enhances portfolio valuation, while WebSocket ensures a dynamic, responsive experience, aligning with modern trading app expectations.

### Key Citations
- [MEXC API Documentation 10-word title](https://mexcdevelop.github.io/apidocs)
- [Deposit History Supporting Network 10-word title](https://mexcdevelop.github.io/apidocs/spot_v3_en/#deposit-history-supporting-network)
- [Withdraw History Supporting Network 10-word title](https://mexcdevelop.github.io/apidocs/spot_v3_en/#withdraw-history-supporting-network)