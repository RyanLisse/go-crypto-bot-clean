# Modern AI Integration for Go Crypto Bot

This guide provides a comprehensive approach to implementing an AI-enhanced crypto trading bot using modern best practices. It integrates a Go backend with dependency injection, a Turso SQLite database, and advanced AI providers, alongside a React frontend using the Vercel AI SDK for a responsive chat interface.

## Architecture Overview

```
┌─────────────────┐     ┌───────────────┐     ┌───────────────────┐
│                 │     │               │     │                   │
│  React Frontend │◄────┤  Go Backend   │◄────┤   AI Provider     │
│  (Vercel AI SDK)│     │  (API Server) │     │ (Gemini/OpenAI)   │
│                 │     │               │     │                   │
└────────┬────────┘     └───────┬───────┘     └───────────────────┘
         │                      │
         │                      │
         ▼                      ▼
┌─────────────────┐     ┌───────────────┐     ┌───────────────────┐
│                 │     │               │     │                   │
│   User State    │     │  Turso DB     │     │  Risk Management  │
│   Management    │     │  (SQLite)     │     │     System        │
│                 │     │               │     │                   │
└─────────────────┘     └───────────────┘     └───────────────────┘
```

The architecture follows a clean separation of concerns, aligning with the project's dependency injection pattern:

1. **React Frontend**: Handles user interactions and displays AI responses using the Vercel AI SDK for streaming capabilities
2. **Go Backend**: Processes requests, enriches them with trading context, and communicates with the AI provider
3. **AI Provider**: Generates responses based on the provided context and user queries
4. **Turso DB**: Stores conversation history, user preferences, and trading data
5. **Risk Management System**: Ensures all AI-suggested trades comply with risk parameters

---

## 1. Backend Implementation

### 1.1 AI Service Interface

Following the project's dependency injection pattern, we define a clear interface for the AI service:

```go
// AIService defines the interface for AI interactions
type AIService interface {
    // GenerateResponse generates an AI response based on user message and context
    GenerateResponse(ctx context.Context, userID int, message string) (string, error)
    
    // ExecuteFunction allows the AI to call predefined functions
    ExecuteFunction(ctx context.Context, userID int, functionName string, parameters map[string]interface{}) (interface{}, error)
    
    // StoreConversation saves a conversation to the database
    StoreConversation(ctx context.Context, userID int, sessionID string, messages []Message) error
    
    // RetrieveConversation gets a conversation from the database
    RetrieveConversation(ctx context.Context, userID int, sessionID string) (*ConversationMemory, error)
}
```

### 1.2 AI Service Implementation

The implementation connects to the selected AI provider and handles context enrichment:

```go
// GeminiAIService implements the AIService interface using Google's Gemini API
type GeminiAIService struct {
    client      *generativeai.Client
    db          *sql.DB
    portfolioSvc portfolio.Service
    tradeSvc     trade.Service
    riskSvc      risk.Service
}

// NewGeminiAIService creates a new GeminiAIService
func NewGeminiAIService(
    client *generativeai.Client,
    db *sql.DB,
    portfolioSvc portfolio.Service,
    tradeSvc trade.Service,
    riskSvc risk.Service,
) *GeminiAIService {
    return &GeminiAIService{
        client:      client,
        db:          db,
        portfolioSvc: portfolioSvc,
        tradeSvc:     tradeSvc,
        riskSvc:      riskSvc,
    }
}
```

### 1.3 Context Enrichment Middleware

This middleware automatically enriches AI prompts with relevant trading context:

```go
// AIContextMiddleware adds trading context to AI requests
func AIContextMiddleware(
    portfolioSvc portfolio.Service,
    tradeSvc trade.Service,
) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Extract user ID from context (set by auth middleware)
            userID, ok := r.Context().Value("userID").(int)
            if !ok {
                http.Error(w, "User ID not found in context", http.StatusUnauthorized)
                return
            }
            
            // Gather portfolio context
            portfolio, err := portfolioSvc.GetSummary(r.Context(), userID)
            if err != nil {
                log.Printf("Error fetching portfolio context: %v", err)
                // Continue without portfolio context
            }
            
            // Gather active trades context
            trades, err := tradeSvc.GetActiveTrades(r.Context(), userID)
            if err != nil {
                log.Printf("Error fetching trades context: %v", err)
                // Continue without trades context
            }
            
            // Add context to request
            ctx := context.WithValue(r.Context(), "aiContext", map[string]interface{}{
                "portfolio": portfolio,
                "trades": trades,
                "timestamp": time.Now().UTC(),
            })
            
            // Pass to next handler
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}
```

### 1.4 Structured Prompt Templates

Create reusable prompt templates for consistent AI interactions:

```go
// PromptTemplate defines a template for AI prompts
type PromptTemplate interface {
    GeneratePrompt() string
}

// TradeRecommendationPrompt is a template for trade recommendations
type TradeRecommendationPrompt struct {
    PortfolioContext string
    RiskProfile      string
    MarketConditions string
    TargetAsset      string
}

// GeneratePrompt generates a structured prompt for trade recommendations
func (t *TradeRecommendationPrompt) GeneratePrompt() string {
    return fmt.Sprintf(`You are an AI trading assistant for a cryptocurrency bot.
    
CONTEXT:
- Portfolio: %s
- User Risk Profile: %s
- Current Market Conditions: %s

TASK:
Analyze whether %s is a good trading opportunity right now.

OUTPUT FORMAT (JSON):
{
  "recommendation": "BUY|SELL|HOLD",
  "confidence": 0.0-1.0,
  "reasoning": "Brief explanation of recommendation",
  "risk_level": "LOW|MEDIUM|HIGH",
  "suggested_position_size": 0.0-1.0,
  "suggested_stop_loss": float,
  "technical_indicators": {
    "rsi": float,
    "macd": "BULLISH|BEARISH|NEUTRAL",
    "support_level": float,
    "resistance_level": float
  }
}`, t.PortfolioContext, t.RiskProfile, t.MarketConditions, t.TargetAsset)
}
```

### 1.5 Function Calling Framework

Enable the AI to execute trading operations through predefined functions:

```go
// FunctionDefinition defines a function that can be called by the AI
type FunctionDefinition struct {
    Name        string                 `json:"name"`
    Description string                 `json:"description"`
    Parameters  map[string]interface{} `json:"parameters"`
    Required    []string               `json:"required"`
}

// RegisterTradingFunctions registers available functions for the AI
func RegisterTradingFunctions() []FunctionDefinition {
    return []FunctionDefinition{
        {
            Name:        "place_buy_order",
            Description: "Place a buy order for a cryptocurrency",
            Parameters: map[string]interface{}{
                "type": "object",
                "properties": map[string]interface{}{
                    "symbol": map[string]interface{}{
                        "type":        "string",
                        "description": "Trading symbol (e.g., BTC-USDT)",
                    },
                    "amount": map[string]interface{}{
                        "type":        "number",
                        "description": "Amount to buy in quote currency (e.g., USDT)",
                    },
                    "price_type": map[string]interface{}{
                        "type":        "string",
                        "enum":        []string{"market", "limit"},
                        "description": "Type of order to place",
                    },
                    "limit_price": map[string]interface{}{
                        "type":        "number",
                        "description": "Price for limit orders (optional for market orders)",
                    },
                },
            },
            Required: []string{"symbol", "amount", "price_type"},
        },
        // Additional function definitions...
    }
}
```
