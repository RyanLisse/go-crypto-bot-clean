# MEXC WebSocket and New Coin Detection Guide

This guide provides a comprehensive overview of implementing and using the MEXC WebSocket client and New Coin Detection service in the Go Crypto Bot project.

## Table of Contents

1. [Overview](#overview)
2. [MEXC WebSocket Client](#mexc-websocket-client)
   - [Features](#features)
   - [Implementation Details](#implementation-details)
   - [Usage Examples](#usage-examples)
   - [Best Practices](#best-practices)
3. [New Coin Detection](#new-coin-detection)
   - [Detection Strategies](#detection-strategies)
   - [Implementation Details](#implementation-details-1)
   - [Usage Examples](#usage-examples-1)
   - [Integration with WebSocket](#integration-with-websocket)
4. [Hybrid Approach](#hybrid-approach)
5. [Troubleshooting](#troubleshooting)
6. [Future Improvements](#future-improvements)

## Overview

The Go Crypto Bot project uses two primary mechanisms to interact with the MEXC exchange:

1. **REST API Client**: For fetching market data, account information, and placing orders
2. **WebSocket Client**: For real-time market data streaming and event notifications

These components work together to provide a robust foundation for trading strategies, particularly for detecting and trading newly listed coins.

## MEXC WebSocket Client

### Features

The WebSocket client (`internal/platform/mexc/websocket/client.go`) provides:

- Thread-safe WebSocket connection management
- Automatic reconnection with exponential backoff
- Ping/pong handling for connection health monitoring
- Subscription management for market data streams
- Rate limiting to prevent API abuse
- Context propagation for proper cancellation

### Implementation Details

The WebSocket client is implemented with the following components:

1. **Connection Management**:
   ```go
   // Connect establishes a WebSocket connection
   func (c *Client) Connect(ctx context.Context) error {
       // Connection logic with rate limiting and error handling
   }
   
   // Disconnect closes the WebSocket connection
   func (c *Client) Disconnect() error {
       // Clean disconnection logic
   }
   
   // reconnect attempts to reestablish the WebSocket connection
   func (c *Client) reconnect() error {
       // Exponential backoff reconnection logic
   }
   ```

2. **Subscription Management**:
   ```go
   // SubscribeToTickers subscribes to ticker updates for given symbols
   func (c *Client) SubscribeToTickers(ctx context.Context, symbols []string) error {
       // Subscription logic with rate limiting
   }
   
   // UnsubscribeFromTickers unsubscribes from ticker updates
   func (c *Client) UnsubscribeFromTickers(symbols []string) error {
       // Unsubscription logic
   }
   ```

3. **Message Processing**:
   ```go
   // handleMessages processes incoming WebSocket messages
   func (c *Client) handleMessages() {
       // Message processing loop
   }
   
   // processMessage handles different types of messages
   func (c *Client) processMessage(message []byte) {
       // Message type detection and routing
   }
   ```

4. **Health Monitoring**:
   ```go
   // startPingTicker starts a ticker to send periodic ping messages
   func (c *Client) startPingTicker() {
       // Ping scheduling and monitoring
   }
   
   // handlePing responds to ping messages from the server
   func (c *Client) handlePing(msg map[string]interface{}) {
       // Ping response handling
   }
   ```

### Usage Examples

**Basic Usage**:

```go
// Create a new WebSocket client
client, err := websocket.NewClient(cfg)
if err != nil {
    log.Fatalf("Failed to create WebSocket client: %v", err)
}

// Connect to the WebSocket server
err = client.Connect(context.Background())
if err != nil {
    log.Fatalf("Failed to connect: %v", err)
}

// Subscribe to ticker updates
err = client.SubscribeToTickers(context.Background(), []string{"BTCUSDT", "ETHUSDT"})
if err != nil {
    log.Printf("Failed to subscribe: %v", err)
}

// Process ticker updates
go func() {
    for ticker := range client.TickerChannel() {
        fmt.Printf("Ticker update: %s - Price: %f\n", ticker.Symbol, ticker.Price)
    }
}()

// Graceful shutdown
defer client.Disconnect()
```

**Advanced Usage with Context Handling**:

```go
// Create a cancellable context
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

// Create a new WebSocket client
client, err := websocket.NewClient(cfg)
if err != nil {
    log.Fatalf("Failed to create WebSocket client: %v", err)
}

// Connect with context
err = client.Connect(ctx)
if err != nil {
    log.Fatalf("Failed to connect: %v", err)
}

// Subscribe with context
err = client.SubscribeToTickers(ctx, []string{"BTCUSDT", "ETHUSDT"})
if err != nil {
    log.Printf("Failed to subscribe: %v", err)
}

// Process ticker updates with graceful shutdown
go func() {
    for {
        select {
        case ticker, ok := <-client.TickerChannel():
            if !ok {
                return // Channel closed
            }
            fmt.Printf("Ticker update: %s - Price: %f\n", ticker.Symbol, ticker.Price)
        case <-ctx.Done():
            return // Context cancelled
        }
    }
}()

// Disconnect on signal
signalCh := make(chan os.Signal, 1)
signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
<-signalCh
cancel() // Cancel context
client.Disconnect()
```

### Best Practices

1. **Always use context for cancellation**:
   ```go
   ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
   defer cancel()
   client.Connect(ctx)
   ```

2. **Handle reconnection events**:
   ```go
   client.SetReconnectHandler(func() error {
       log.Println("Reconnected to WebSocket")
       // Perform any necessary re-initialization
       return nil
   })
   ```

3. **Implement proper error handling**:
   ```go
   err = client.SubscribeToTickers(ctx, symbols)
   if err != nil {
       if errors.Is(err, context.DeadlineExceeded) {
           // Handle timeout
       } else if errors.Is(err, websocket.ErrNotConnected) {
           // Handle connection issue
       } else {
           // Handle other errors
       }
   }
   ```

4. **Use rate limiting appropriately**:
   ```go
   // Configure rate limiters based on exchange limits
   connLimiter := ratelimiter.NewTokenBucketRateLimiter(1.0, 1.0) // 1 conn/sec
   subLimiter := ratelimiter.NewTokenBucketRateLimiter(5.0, 10.0) // 5 subs/sec, burst of 10
   
   client, err := websocket.NewClient(
       cfg,
       websocket.WithConnRateLimiter(connLimiter),
       websocket.WithSubRateLimiter(subLimiter),
   )
   ```

## New Coin Detection

### Detection Strategies

The Go Crypto Bot implements two complementary strategies for detecting new coin listings:

1. **REST API Polling**: Periodically fetching exchange information to identify new symbols
2. **WebSocket Market Activity**: Monitoring real-time market data for sudden activity in new symbols

### Implementation Details

The New Coin Service (`internal/core/newcoin/newcoin_service.go`) provides:

1. **Detection Logic**:
   ```go
   // DetectNewCoins identifies new coin listings on MEXC
   func (s *newCoinService) DetectNewCoins(ctx context.Context) ([]*models.NewCoin, error) {
       // Fetch new coins from the exchange
       coins, err := s.mexcClient.GetNewCoins(ctx)
       if err != nil {
           return nil, err
       }
       return coins, nil
   }
   ```

2. **Persistence**:
   ```go
   // SaveNewCoins saves the detected coins to the repository
   func (s *newCoinService) SaveNewCoins(ctx context.Context, coins []*models.NewCoin) error {
       for _, coin := range coins {
           _, err := s.newCoinRepo.Create(ctx, coin)
           if err != nil {
               return err
           }
       }
       return nil
   }
   ```

3. **Continuous Monitoring**:
   ```go
   // StartWatching begins watching for new coins at a specified interval
   func (s *newCoinService) StartWatching(ctx context.Context, interval time.Duration) error {
       ticker := time.NewTicker(interval)
       defer ticker.Stop()

       for {
           select {
           case <-ticker.C:
               // Detect and save new coins
               coins, err := s.DetectNewCoins(ctx)
               if err != nil {
                   continue
               }
               if len(coins) > 0 {
                   _ = s.SaveNewCoins(ctx, coins)
               }
           case <-s.stopChan:
               return nil
           case <-ctx.Done():
               return ctx.Err()
           }
       }
   }
   ```

The REST client (`internal/platform/mexc/rest/client.go`) implements the `GetNewCoins` method:

```go
// GetNewCoins gets newly listed coins
func (c *Client) GetNewCoins(ctx context.Context) ([]*models.NewCoin, error) {
    // Check cache first
    if newCoins, found := c.newCoinCache.GetNewCoins(); found {
        return newCoins, nil
    }

    // Fetch exchange information
    var exchangeInfo struct {
        Symbols []struct {
            Symbol     string   `json:"symbol"`
            Status     string   `json:"status"`
            BaseAsset  string   `json:"baseAsset"`
            QuoteAsset string   `json:"quoteAsset"`
            // Other fields...
        } `json:"symbols"`
    }

    err := c.makeRequest(ctx, http.MethodGet, "/api/v3/exchangeInfo", nil, false, &exchangeInfo)
    if err != nil {
        return nil, err
    }

    // Get all tickers for volume information
    tickers, err := c.GetAllTickers(ctx)
    if err != nil {
        return nil, err
    }

    // Filter for USDT pairs and create new coins
    newCoins := make([]*models.NewCoin, 0)
    for _, symbol := range exchangeInfo.Symbols {
        if symbol.Status == "ENABLED" && symbol.QuoteAsset == "USDT" {
            ticker, exists := tickers[symbol.Symbol]
            if !exists {
                continue
            }

            newCoins = append(newCoins, &models.NewCoin{
                Symbol:      symbol.Symbol,
                BaseVolume:  0,
                QuoteVolume: ticker.Volume,
                FoundAt:     time.Now(),
            })
        }
    }

    // Sort by volume and limit results
    sort.Slice(newCoins, func(i, j int) bool {
        return newCoins[i].QuoteVolume > newCoins[j].QuoteVolume
    })

    if len(newCoins) > 20 {
        newCoins = newCoins[:20]
    }

    // Cache results
    c.newCoinCache.SetNewCoins(newCoins, 5*time.Minute)

    return newCoins, nil
}
```

### Usage Examples

**Basic Usage**:

```go
// Create dependencies
mexcClient, _ := rest.NewClient(apiKey, secretKey)
newCoinRepo := repository.NewNewCoinRepository(db)

// Create new coin service
newCoinService := newcoin.NewNewCoinService(mexcClient, newCoinRepo)

// Detect new coins
ctx := context.Background()
coins, err := newCoinService.DetectNewCoins(ctx)
if err != nil {
    log.Fatalf("Failed to detect new coins: %v", err)
}

// Process detected coins
for _, coin := range coins {
    fmt.Printf("New coin: %s, Volume: %.2f USDT\n", coin.Symbol, coin.QuoteVolume)
}

// Start continuous monitoring
go func() {
    err := newCoinService.StartWatching(ctx, 1*time.Minute)
    if err != nil {
        log.Printf("Watching stopped: %v", err)
    }
}()

// Stop monitoring when done
defer newCoinService.StopWatching()
```

**Advanced Usage with Processing Logic**:

```go
// Create dependencies
mexcClient, _ := rest.NewClient(apiKey, secretKey)
newCoinRepo := repository.NewNewCoinRepository(db)
tradeService := trade.NewTradeService(mexcClient, tradeRepo)

// Create new coin service
newCoinService := newcoin.NewNewCoinService(mexcClient, newCoinRepo)

// Start monitoring with custom processing
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

go func() {
    for {
        // Detect new coins
        coins, err := newCoinService.DetectNewCoins(ctx)
        if err != nil {
            log.Printf("Error detecting new coins: %v", err)
            time.Sleep(10 * time.Second)
            continue
        }

        // Process each new coin
        for _, coin := range coins {
            // Check if it's already in our database
            existing, err := newCoinRepo.FindBySymbol(ctx, coin.Symbol)
            if err == nil && existing != nil {
                // Already tracked, skip
                continue
            }

            // Save the new coin
            _, err = newCoinRepo.Create(ctx, coin)
            if err != nil {
                log.Printf("Error saving new coin %s: %v", coin.Symbol, err)
                continue
            }

            log.Printf("New coin detected: %s with volume %.2f USDT", 
                coin.Symbol, coin.QuoteVolume)

            // Implement your trading strategy here
            // For example, place a market buy order
            if coin.QuoteVolume > 10000 { // Only trade high volume coins
                order := &models.Order{
                    Symbol:   coin.Symbol,
                    Side:     models.OrderSideBuy,
                    Type:     models.OrderTypeMarket,
                    Quantity: 10.0, // Fixed USDT amount
                }
                
                _, err = tradeService.PlaceOrder(ctx, order)
                if err != nil {
                    log.Printf("Failed to place order for %s: %v", coin.Symbol, err)
                } else {
                    log.Printf("Successfully placed order for new coin %s", coin.Symbol)
                }
            }
        }

        // Wait before next check
        time.Sleep(30 * time.Second)
    }
}()
```

### Integration with WebSocket

While REST API polling is effective, integrating WebSocket data can provide faster detection of new coins through market activity:

```go
// Create WebSocket client
wsClient, _ := websocket.NewClient(cfg)
wsClient.Connect(ctx)

// Subscribe to all tickers
allSymbols, _ := mexcClient.GetAllSymbols(ctx)
wsClient.SubscribeToTickers(ctx, allSymbols)

// Track unknown symbols from WebSocket
unknownSymbols := make(map[string]time.Time)
symbolMutex := sync.RWMutex{}

// Process ticker updates
go func() {
    for ticker := range wsClient.TickerChannel() {
        // Check if this is a known symbol
        existing, err := newCoinRepo.FindBySymbol(ctx, ticker.Symbol)
        if err == nil && existing != nil {
            // Already tracked, skip
            continue
        }

        // Check if we've already seen this unknown symbol
        symbolMutex.RLock()
        firstSeen, exists := unknownSymbols[ticker.Symbol]
        symbolMutex.RUnlock()

        if !exists {
            // First time seeing this symbol, record it
            symbolMutex.Lock()
            unknownSymbols[ticker.Symbol] = time.Now()
            symbolMutex.Unlock()
            
            log.Printf("Potential new coin detected via WebSocket: %s", ticker.Symbol)
            
            // Verify with REST API
            go func(symbol string) {
                // Get detailed information about this symbol
                symbolInfo, err := mexcClient.GetSymbolInfo(ctx, symbol)
                if err != nil {
                    log.Printf("Failed to get info for symbol %s: %v", symbol, err)
                    return
                }
                
                // Create and save new coin
                newCoin := &models.NewCoin{
                    Symbol:      symbol,
                    BaseVolume:  ticker.Volume,
                    QuoteVolume: ticker.Volume * ticker.Price,
                    FoundAt:     time.Now(),
                }
                
                _, err = newCoinRepo.Create(ctx, newCoin)
                if err != nil {
                    log.Printf("Failed to save new coin %s: %v", symbol, err)
                }
                
                log.Printf("New coin confirmed and saved: %s", symbol)
            }(ticker.Symbol)
        }
    }
}()
```

## Hybrid Approach

For optimal new coin detection, implement a hybrid approach combining both REST API polling and WebSocket monitoring:

```go
func StartNewCoinDetection(
    ctx context.Context,
    mexcClient *rest.Client,
    wsClient *websocket.Client,
    newCoinRepo repository.NewCoinRepository,
) {
    // 1. Start REST API polling
    go func() {
        ticker := time.NewTicker(1 * time.Minute)
        defer ticker.Stop()
        
        for {
            select {
            case <-ticker.C:
                coins, err := mexcClient.GetNewCoins(ctx)
                if err != nil {
                    log.Printf("Error in REST polling: %v", err)
                    continue
                }
                
                for _, coin := range coins {
                    // Check if already exists
                    existing, _ := newCoinRepo.FindBySymbol(ctx, coin.Symbol)
                    if existing != nil {
                        continue
                    }
                    
                    // Save new coin
                    _, err := newCoinRepo.Create(ctx, coin)
                    if err != nil {
                        log.Printf("Failed to save coin %s: %v", coin.Symbol, err)
                    } else {
                        log.Printf("New coin detected via REST: %s", coin.Symbol)
                    }
                }
                
            case <-ctx.Done():
                return
            }
        }
    }()
    
    // 2. Start WebSocket monitoring
    unknownSymbols := make(map[string]time.Time)
    symbolMutex := sync.RWMutex{}
    
    // Subscribe to all tickers
    allSymbols, _ := mexcClient.GetAllSymbols(ctx)
    wsClient.SubscribeToTickers(ctx, allSymbols)
    
    go func() {
        for {
            select {
            case ticker, ok := <-wsClient.TickerChannel():
                if !ok {
                    return
                }
                
                // Process ticker for new coin detection
                existing, _ := newCoinRepo.FindBySymbol(ctx, ticker.Symbol)
                if existing != nil {
                    continue
                }
                
                symbolMutex.RLock()
                firstSeen, exists := unknownSymbols[ticker.Symbol]
                symbolMutex.RUnlock()
                
                if !exists {
                    symbolMutex.Lock()
                    unknownSymbols[ticker.Symbol] = time.Now()
                    symbolMutex.Unlock()
                    
                    log.Printf("Potential new coin via WebSocket: %s", ticker.Symbol)
                    
                    // Verify with REST API
                    symbolInfo, err := mexcClient.GetSymbolInfo(ctx, ticker.Symbol)
                    if err != nil {
                        continue
                    }
                    
                    newCoin := &models.NewCoin{
                        Symbol:      ticker.Symbol,
                        BaseVolume:  ticker.Volume,
                        QuoteVolume: ticker.Volume * ticker.Price,
                        FoundAt:     time.Now(),
                    }
                    
                    _, err = newCoinRepo.Create(ctx, newCoin)
                    if err == nil {
                        log.Printf("New coin confirmed via WebSocket: %s", ticker.Symbol)
                    }
                }
                
            case <-ctx.Done():
                return
            }
        }
    }()
    
    // 3. Periodically clean up old unknown symbols
    go func() {
        ticker := time.NewTicker(1 * time.Hour)
        defer ticker.Stop()
        
        for {
            select {
            case <-ticker.C:
                now := time.Now()
                symbolMutex.Lock()
                for symbol, firstSeen := range unknownSymbols {
                    if now.Sub(firstSeen) > 24*time.Hour {
                        delete(unknownSymbols, symbol)
                    }
                }
                symbolMutex.Unlock()
                
            case <-ctx.Done():
                return
            }
        }
    }()
}
```

## Troubleshooting

### Common Issues and Solutions

1. **WebSocket Connection Failures**
   - **Symptom**: Frequent disconnections or failure to connect
   - **Solution**: 
     - Check network connectivity
     - Verify API endpoint is correct
     - Ensure rate limits are respected
     - Implement exponential backoff for reconnection attempts

2. **Missing New Coin Notifications**
   - **Symptom**: New coins appear on the exchange but aren't detected
   - **Solution**:
     - Decrease polling interval
     - Ensure WebSocket subscriptions are active
     - Check filtering logic for false negatives
     - Verify that the exchange information endpoint is returning all symbols

3. **High CPU/Memory Usage**
   - **Symptom**: Application consumes excessive resources
   - **Solution**:
     - Increase polling intervals
     - Implement more efficient filtering
     - Use buffered channels with appropriate sizes
     - Add caching with reasonable TTLs

4. **Rate Limiting Issues**
   - **Symptom**: Frequent 429 (Too Many Requests) errors
   - **Solution**:
     - Adjust rate limiters to respect exchange limits
     - Implement exponential backoff for retries
     - Use caching to reduce API calls

### Debugging Tips

1. **Enable Detailed Logging**:
   ```go
   logger, _ := zap.NewDevelopment()
   client, _ := websocket.NewClient(cfg, websocket.WithLogger(logger))
   ```

2. **Monitor WebSocket Traffic**:
   ```go
   // Add a message interceptor for debugging
   client.SetMessageInterceptor(func(direction string, message []byte) {
       fmt.Printf("[%s] %s\n", direction, string(message))
   })
   ```

3. **Test Connection Health**:
   ```go
   // Periodically check connection status
   go func() {
       ticker := time.NewTicker(30 * time.Second)
       defer ticker.Stop()
       
       for range ticker.C {
           if client.IsConnected() {
               log.Println("WebSocket connection is healthy")
           } else {
               log.Println("WebSocket connection is down")
           }
       }
   }()
   ```

## Future Improvements

1. **Enhanced Detection Algorithms**
   - Implement machine learning for volume spike detection
   - Add sentiment analysis from social media for new coin prediction
   - Develop historical pattern recognition for listing events

2. **Performance Optimizations**
   - Implement more efficient filtering mechanisms
   - Use goroutine pools for parallel processing
   - Optimize database queries for faster lookups

3. **Reliability Enhancements**
   - Add circuit breakers for API calls
   - Implement fallback mechanisms for critical operations
   - Add comprehensive metrics and alerting

4. **Feature Additions**
   - Support for additional exchanges
   - Configurable trading strategies for new coins
   - Integration with notification systems (Telegram, Discord, etc.)
   - Dashboard for monitoring new coin performance
