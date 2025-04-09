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
