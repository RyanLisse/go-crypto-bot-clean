## 3. Security and Best Practices

### 3.1 API Key Management

Proper management of AI provider API keys is critical for security:

```go
// Environment variable configuration
type Config struct {
    AIProviderAPIKey string
    AIProviderModel  string
    DatabaseURL      string
    JWTSecret        string
    // Other configuration options...
}

// Load configuration from environment variables
func LoadConfig() (*Config, error) {
    config := &Config{
        AIProviderAPIKey: os.Getenv("AI_PROVIDER_API_KEY"),
        AIProviderModel:  os.Getenv("AI_PROVIDER_MODEL"),
        DatabaseURL:      os.Getenv("DATABASE_URL"),
        JWTSecret:        os.Getenv("JWT_SECRET"),
    }
    
    // Validate required configuration
    if config.AIProviderAPIKey == "" {
        return nil, errors.New("AI_PROVIDER_API_KEY environment variable is required")
    }
    if config.DatabaseURL == "" {
        return nil, errors.New("DATABASE_URL environment variable is required")
    }
    
    return config, nil
}
```

### 3.2 Rate Limiting and Circuit Breaking

Implement rate limiting to prevent excessive API usage and costs:

```go
// RateLimiter limits the rate of AI API calls
type RateLimiter struct {
    mu           sync.Mutex
    requestCount int
    lastReset    time.Time
    maxRequests  int
    resetPeriod  time.Duration
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(maxRequests int, resetPeriod time.Duration) *RateLimiter {
    return &RateLimiter{
        maxRequests: maxRequests,
        resetPeriod: resetPeriod,
        lastReset:   time.Now(),
    }
}

// Allow checks if a request is allowed
func (r *RateLimiter) Allow() bool {
    r.mu.Lock()
    defer r.mu.Unlock()
    
    now := time.Now()
    if now.Sub(r.lastReset) > r.resetPeriod {
        r.requestCount = 0
        r.lastReset = now
    }
    
    if r.requestCount >= r.maxRequests {
        return false
    }
    
    r.requestCount++
    return true
}
```

### 3.3 Input Validation and Sanitization

Validate and sanitize all user inputs to prevent prompt injection and other security issues:

```go
// SanitizeUserInput sanitizes user input to prevent prompt injection
func SanitizeUserInput(input string) string {
    // Remove potentially harmful characters and sequences
    input = strings.ReplaceAll(input, "\n", " ")
    input = strings.ReplaceAll(input, "\r", " ")
    input = strings.ReplaceAll(input, "\t", " ")
    
    // Limit input length
    const maxInputLength = 1000
    if len(input) > maxInputLength {
        input = input[:maxInputLength]
    }
    
    return input
}

// ValidateTradeRequest validates a trade request
func ValidateTradeRequest(req *TradeRequest) error {
    if req.Symbol == "" {
        return errors.New("symbol is required")
    }
    
    if req.Amount <= 0 {
        return errors.New("amount must be positive")
    }
    
    if req.PriceType != "market" && req.PriceType != "limit" {
        return errors.New("price_type must be 'market' or 'limit'")
    }
    
    if req.PriceType == "limit" && req.LimitPrice <= 0 {
        return errors.New("limit_price must be positive for limit orders")
    }
    
    return nil
}
```
