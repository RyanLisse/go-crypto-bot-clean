# Cache Refactoring - Implementation Plan

## Current State Analysis

1. **Multiple Cache Implementations**:
   - `GenericCache` (`internal/adapter/cache/memory/generic_cache.go`): Simple TTL-based cache for single values
   - `MarketCache` (`internal/adapter/cache/memory/market_cache.go`): Complex custom implementation for market data
   - `TickerCache` (`internal/adapter/cache/memory/ticker_cache.go`): Specific implementation for ticker data
   - `StandardCache` (`internal/adapter/cache/standard/cache.go`): In-progress implementation using go-cache

2. **Interfaces**:
   - `port.Cache[T]`: Generic cache interface for any type
   - `port.MarketCache`: Specialized interface for market data

3. **Issues with Current Implementation**:
   - Duplicate functionality across cache implementations
   - Complex manual cleanup and expiry checking in `MarketCache`
   - Potential concurrency issues

## Implementation Steps

### 1. Complete the StandardCache Implementation

We've already started this with `internal/adapter/cache/standard/cache.go`, but it needs:
- Fixed import issues (entity vs model)
- Implementation of all required methods from `port.MarketCache`
- Consistent key generation strategy

### 2. Update Factory to Use StandardCache

1. Create or update the cache factory to use StandardCache instead of the custom implementations.

Example factory update:
```go
// NewMarketCache creates a new cache instance
func NewMarketCache(logger *zerolog.Logger) port.MarketCache {
    // Old custom implementation
    // return memory.NewMarketCache(logger)
    
    // New implementation with go-cache
    return standard.NewStandardCache(
        5*time.Minute, // default TTL
        10*time.Minute, // cleanup interval
    )
}
```

### 3. Update Dependency Injection Chain

1. Update all code that creates a cache instance to use the new factory method
2. Make sure all dependencies (usecase, services) receive the interface, not the concrete type

### 4. Implement Missing Features

1. Complete the ticker tracking functionality in StandardCache
2. Add support for batch operations if needed
3. Implement proper logging

### 5. Test the New Implementation

1. Create unit tests for StandardCache
2. Test all operations (get, set, expire, etc.)
3. Test with high concurrency to ensure thread safety

### 6. Remove Old Implementations

Once the new implementation is in place and tested:
1. Remove the old cache implementation files
2. Remove related tests for old implementations

## Implementation Details

### StandardCache Improvement

The `internal/adapter/cache/standard/cache.go` implementation should be completed to include:

1. Proper key generation:
```go
func (c *StandardCache) generateTickerKey(exchange, symbol string) string {
    return fmt.Sprintf("ticker:%s:%s", exchange, symbol)
}
```

2. Track collections of items using prefix-based keys:
```go
// GetAllTickers retrieves all tickers for an exchange from cache
func (c *StandardCache) GetAllTickers(ctx context.Context, exchange string) ([]*market.Ticker, bool) {
    // Use Items() to get all items with the prefix "ticker:{exchange}:"
    prefix := fmt.Sprintf("ticker:%s:", exchange)
    tickers := make([]*market.Ticker, 0)
    
    for k, v := range c.tickerCache.Items() {
        if strings.HasPrefix(k, prefix) {
            if ticker, ok := v.Object.(*market.Ticker); ok {
                tickers = append(tickers, ticker)
            }
        }
    }
    
    if len(tickers) == 0 {
        return nil, false
    }
    
    return tickers, true
}
```

3. Manage expiration properly:
```go
func (c *StandardCache) SetTickerExpiry(d time.Duration) {
    c.tickerCache.Flush()
    c.tickerCache = gocache.New(d, d/2)
}
```

## Migration Checklist

- [x] Fix imports in StandardCache
- [ ] Complete all required methods in StandardCache
- [ ] Add appropriate tests for StandardCache
- [ ] Update factory to use StandardCache
- [ ] Update dependency injection
- [ ] Verify all functionality with integration tests
- [ ] Remove old cache implementations 