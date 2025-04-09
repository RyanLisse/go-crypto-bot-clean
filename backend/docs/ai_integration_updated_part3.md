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

---

## 4. Testing and Monitoring

### 4.1 Unit Testing

Create comprehensive unit tests for all AI-related components:

```go
// Test for trade recommendation prompt
func TestTradeRecommendationPrompt(t *testing.T) {
    // Test case setup
    prompt := &TradeRecommendationPrompt{
        PortfolioContext: "Portfolio value: $10,000, BTC: 0.5, ETH: 2.0",
        RiskProfile:      "Medium risk tolerance, max 5% per position",
        MarketConditions: "BTC up 2% in last 24h, ETH down 1%",
        TargetAsset:      "BTC",
    }
    
    // Generate the prompt
    promptText := prompt.GeneratePrompt()
    
    // Verify prompt contains all required elements
    requiredElements := []string{
        "Portfolio", "Risk Profile", "Market Conditions",
        "BTC", "recommendation", "JSON",
    }
    
    for _, element := range requiredElements {
        if !strings.Contains(promptText, element) {
            t.Errorf("Prompt missing required element: %s", element)
        }
    }
}

// Mock AI service for testing
type MockAIService struct {
    responses map[string]string
}

func NewMockAIService() *MockAIService {
    return &MockAIService{
        responses: map[string]string{
            "BTC": `{"recommendation":"BUY","confidence":0.85,"reasoning":"Strong uptrend","risk_level":"MEDIUM"}`,
            "ETH": `{"recommendation":"HOLD","confidence":0.65,"reasoning":"Consolidating","risk_level":"MEDIUM"}`,
        },
    }
}

func (m *MockAIService) GenerateResponse(ctx context.Context, userID int, message string) (string, error) {
    // Simple mock implementation that returns predefined responses based on message content
    for key, response := range m.responses {
        if strings.Contains(message, key) {
            return response, nil
        }
    }
    return `{"recommendation":"HOLD","confidence":0.5,"reasoning":"Insufficient data","risk_level":"MEDIUM"}`, nil
}
```

### 4.2 Integration Testing

Create integration tests for the complete AI workflow:

```go
// Integration test for chat endpoint
func TestChatEndpoint(t *testing.T) {
    // Create mock services
    mockAI := NewMockAIService()
    mockDB, err := sql.Open("sqlite3", ":memory:")
    if err != nil {
        t.Fatalf("Failed to open in-memory database: %v", err)
    }
    defer mockDB.Close()
    
    // Create handler with mock services
    handler := ChatHandler(mockAI)
    
    // Create test server
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Mock auth middleware
        ctx := context.WithValue(r.Context(), "userID", 123)
        handler.ServeHTTP(w, r.WithContext(ctx))
    }))
    defer server.Close()
    
    // Create test request
    reqBody := map[string]interface{}{
        "messages": []map[string]string{
            {"role": "user", "content": "Should I buy BTC?"},
        },
    }
    reqJSON, err := json.Marshal(reqBody)
    if err != nil {
        t.Fatalf("Failed to marshal request body: %v", err)
    }
    
    // Send request
    resp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(reqJSON))
    if err != nil {
        t.Fatalf("Failed to send request: %v", err)
    }
    defer resp.Body.Close()
    
    // Check response
    if resp.StatusCode != http.StatusOK {
        t.Errorf("Expected status OK, got %v", resp.Status)
    }
    
    // Parse response
    var respBody map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
        t.Fatalf("Failed to decode response body: %v", err)
    }
    
    // Check response contains expected fields
    if _, ok := respBody["output"]; !ok {
        t.Errorf("Response missing 'output' field")
    }
    if _, ok := respBody["session_id"]; !ok {
        t.Errorf("Response missing 'session_id' field")
    }
}
```

### 4.3 Monitoring and Alerting

Implement monitoring for AI usage, costs, and performance metrics:

```go
// AIMetrics tracks AI usage metrics
type AIMetrics struct {
    RequestCount     int64
    TokenCount       int64
    TotalLatency     time.Duration
    ErrorCount       int64
    LastRequestTime  time.Time
    CostEstimate     float64
    mu               sync.Mutex
}

// NewAIMetrics creates a new AIMetrics
func NewAIMetrics() *AIMetrics {
    return &AIMetrics{}
}

// RecordRequest records a request
func (m *AIMetrics) RecordRequest(tokens int, latency time.Duration, err error, costPerToken float64) {
    m.mu.Lock()
    defer m.mu.Unlock()
    
    m.RequestCount++
    m.TokenCount += int64(tokens)
    m.TotalLatency += latency
    m.LastRequestTime = time.Now()
    m.CostEstimate += float64(tokens) * costPerToken
    
    if err != nil {
        m.ErrorCount++
    }
}

// GetMetrics gets the current metrics
func (m *AIMetrics) GetMetrics() map[string]interface{} {
    m.mu.Lock()
    defer m.mu.Unlock()
    
    var avgLatency time.Duration
    if m.RequestCount > 0 {
        avgLatency = m.TotalLatency / time.Duration(m.RequestCount)
    }
    
    return map[string]interface{}{
        "request_count":      m.RequestCount,
        "token_count":        m.TokenCount,
        "avg_latency_ms":     avgLatency.Milliseconds(),
        "error_count":        m.ErrorCount,
        "last_request_time":  m.LastRequestTime,
        "cost_estimate_usd":  m.CostEstimate,
        "error_rate":         float64(m.ErrorCount) / float64(m.RequestCount),
    }
}
```

---

## 5. Deployment and Configuration

### 5.1 Environment Variables

Configure the application using environment variables:

```bash
# AI Provider
AI_PROVIDER=gemini  # or openai, anthropic
AI_PROVIDER_API_KEY=your_api_key_here
AI_PROVIDER_MODEL=gemini-flash  # or gpt-4, claude-3-opus

# Database
DATABASE_URL=libsql://your-turso-db-url
DATABASE_AUTH_TOKEN=your_turso_auth_token

# Security
JWT_SECRET=your_jwt_secret_here
CORS_ALLOWED_ORIGINS=http://localhost:3000,https://your-production-domain.com

# Rate Limiting
MAX_REQUESTS_PER_MINUTE=60
MAX_TOKENS_PER_DAY=100000

# Logging
LOG_LEVEL=info  # debug, info, warn, error
```

### 5.2 Docker Deployment

Deploy the application using Docker:

```dockerfile
# Dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/app ./cmd/server

# Create a minimal image
FROM alpine:3.18

COPY --from=builder /go/bin/app /app

# Set environment variables
ENV AI_PROVIDER=gemini
ENV AI_PROVIDER_MODEL=gemini-flash
ENV LOG_LEVEL=info

# Expose the port
EXPOSE 8080

# Run the application
CMD ["/app"]
```

### 5.3 Kubernetes Deployment

Deploy the application to Kubernetes:

```yaml
# kubernetes/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: crypto-bot-ai
  labels:
    app: crypto-bot-ai
spec:
  replicas: 3
  selector:
    matchLabels:
      app: crypto-bot-ai
  template:
    metadata:
      labels:
        app: crypto-bot-ai
    spec:
      containers:
      - name: crypto-bot-ai
        image: your-registry/crypto-bot-ai:latest
        ports:
        - containerPort: 8080
        env:
        - name: AI_PROVIDER
          value: "gemini"
        - name: AI_PROVIDER_MODEL
          value: "gemini-flash"
        - name: AI_PROVIDER_API_KEY
          valueFrom:
            secretKeyRef:
              name: ai-secrets
              key: api-key
        - name: DATABASE_URL
          valueFrom:
            configMapKeyRef:
              name: db-config
              key: url
        - name: DATABASE_AUTH_TOKEN
          valueFrom:
            secretKeyRef:
              name: db-secrets
              key: auth-token
        - name: JWT_SECRET
          valueFrom:
            secretKeyRef:
              name: auth-secrets
              key: jwt-secret
        resources:
          limits:
            cpu: "500m"
            memory: "512Mi"
          requests:
            cpu: "100m"
            memory: "128Mi"
```

## 6. Conclusion

This implementation guide provides a comprehensive approach to integrating AI capabilities into the go-crypto-bot-migration project. By following these patterns and best practices, you can create a robust, secure, and performant AI assistant that enhances the trading experience for your users.

Key benefits of this implementation include:

1. **Clean Architecture**: Following the project's dependency injection pattern for maintainable code
2. **Security**: Proper API key management, input validation, and rate limiting
3. **Performance**: Optimized prompts, caching, and streaming responses
4. **User Experience**: Rich, interactive interface with data visualization
5. **Reliability**: Comprehensive testing and monitoring

Remember to regularly update your AI models and prompts as new capabilities become available, and continuously monitor usage to optimize costs and performance.
