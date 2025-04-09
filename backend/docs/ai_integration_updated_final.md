# Modern AI Integration for Go Crypto Bot

This guide provides a comprehensive approach to implementing an AI-enhanced crypto trading bot using modern best practices. It integrates a Go backend with dependency injection, GORM for database operations, and advanced AI providers, alongside a React frontend using the Vercel AI SDK and Drizzle ORM with Turso for a responsive chat interface.

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

1. **React Frontend**: Handles user interactions and displays AI responses using the Vercel AI SDK for streaming capabilities and Drizzle ORM for data persistence
2. **Go Backend**: Processes requests, enriches them with trading context, and communicates with the AI provider using GORM for all database operations
3. **AI Provider**: Generates responses based on the provided context and user queries
4. **Database Layer**:
   - Backend: Uses GORM with SQLite for all database operations including conversation history
   - Frontend: Uses Drizzle ORM with Turso for client-side data persistence
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
    db          *gorm.DB
    portfolioSvc portfolio.Service
    tradeSvc     trade.Service
    riskSvc      risk.Service
}

// NewGeminiAIService creates a new GeminiAIService
func NewGeminiAIService(
    client *generativeai.Client,
    db *gorm.DB,
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
            Name:        "get_market_data",
            Description: "Get current market data for a specific cryptocurrency",
            Parameters: map[string]interface{}{
                "symbol": map[string]string{
                    "type":        "string",
                    "description": "The trading symbol (e.g., BTC, ETH)",
                },
                "timeframe": map[string]string{
                    "type":        "string",
                    "description": "Timeframe for the data (e.g., 1h, 4h, 1d)",
                    "enum":        []string{"1m", "5m", "15m", "1h", "4h", "1d", "1w"},
                },
            },
            Required: []string{"symbol"},
        },
        {
            Name:        "analyze_technical_indicators",
            Description: "Analyze technical indicators for a specific cryptocurrency",
            Parameters: map[string]interface{}{
                "symbol": map[string]string{
                    "type":        "string",
                    "description": "The trading symbol (e.g., BTC, ETH)",
                },
                "indicators": map[string]interface{}{
                    "type":        "array",
                    "description": "List of indicators to analyze",
                    "items": map[string]string{
                        "type": "string",
                        "enum": []string{"rsi", "macd", "bollinger", "ema", "sma", "fibonacci"},
                    },
                },
            },
            Required: []string{"symbol", "indicators"},
        },
        {
            Name:        "execute_trade",
            Description: "Execute a trade (requires confirmation)",
            Parameters: map[string]interface{}{
                "symbol": map[string]string{
                    "type":        "string",
                    "description": "The trading symbol (e.g., BTC, ETH)",
                },
                "action": map[string]string{
                    "type":        "string",
                    "description": "Buy or sell",
                    "enum":        []string{"buy", "sell"},
                },
                "amount": map[string]string{
                    "type":        "number",
                    "description": "Amount to trade (in USD or percentage of portfolio)",
                },
                "price_type": map[string]string{
                    "type":        "string",
                    "description": "Market or limit order",
                    "enum":        []string{"market", "limit"},
                },
                "limit_price": map[string]string{
                    "type":        "number",
                    "description": "Price for limit orders",
                },
                "stop_loss": map[string]string{
                    "type":        "number",
                    "description": "Stop loss price",
                },
                "take_profit": map[string]string{
                    "type":        "number",
                    "description": "Take profit price",
                },
            },
            Required: []string{"symbol", "action", "amount", "price_type"},
        },
    }
}
```

### 1.6 Risk Management Integration

Connect the AI assistant with the existing risk management system to ensure all AI-suggested trading actions adhere to risk parameters:

```go
// Risk-aware trade validation
func validateAITradeWithRiskSystem(
    ctx context.Context,
    trade *TradeRequest,
    userID int,
    riskSvc risk.Service,
) (*RiskAssessment, error) {
    assessment := &RiskAssessment{
        TradeAllowed: false,
        RiskFactors:  []string{},
        Explanation:  "",
    }

    // Check daily loss limit
    dailyLossCheck, err := riskSvc.CheckDailyLossLimit(ctx, userID)
    if err != nil {
        return nil, fmt.Errorf("failed to check daily loss limit: %w", err)
    }
    if !dailyLossCheck.Allowed {
        assessment.RiskFactors = append(assessment.RiskFactors, "daily_loss_limit_exceeded")
        assessment.Explanation += fmt.Sprintf("Daily loss limit of %.2f%% reached. ", dailyLossCheck.Threshold)
    }

    // Check maximum drawdown
    drawdownCheck, err := riskSvc.CheckMaximumDrawdown(ctx, userID)
    if err != nil {
        return nil, fmt.Errorf("failed to check maximum drawdown: %w", err)
    }
    if !drawdownCheck.Allowed {
        assessment.RiskFactors = append(assessment.RiskFactors, "max_drawdown_exceeded")
        assessment.Explanation += fmt.Sprintf("Maximum drawdown of %.2f%% reached. ", drawdownCheck.Threshold)
    }

    // Check exposure limit for the specific asset
    exposureCheck, err := riskSvc.CheckExposureLimit(ctx, userID, trade.Symbol)
    if err != nil {
        return nil, fmt.Errorf("failed to check exposure limit: %w", err)
    }
    if !exposureCheck.Allowed {
        assessment.RiskFactors = append(assessment.RiskFactors, "exposure_limit_exceeded")
        assessment.Explanation += fmt.Sprintf("Exposure limit of %.2f%% for %s reached. ",
            exposureCheck.Threshold, trade.Symbol)
    }

    // Trade is allowed if no risk factors were triggered
    assessment.TradeAllowed = len(assessment.RiskFactors) == 0

    return assessment, nil
}
```

### 1.7 Conversation Memory System with GORM

Implement a sophisticated memory system for AI conversations using GORM for database operations:

```go
// ConversationMemory stores conversation history
type ConversationMemory struct {
    UserID       int       `json:"user_id" gorm:"primaryKey;column:user_id"`
    SessionID    string    `json:"session_id" gorm:"primaryKey;column:session_id"`
    MessagesJSON string    `json:"-" gorm:"column:messages_json"`
    Summary      string    `json:"summary" gorm:"column:summary"`
    LastAccessed time.Time `json:"last_accessed" gorm:"column:last_accessed;index"`
    CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime"`

    // Virtual field for messages (not stored directly in database)
    Messages     []Message `json:"messages" gorm:"-"`
}

type Message struct {
    Role      string                 `json:"role"` // "user" or "assistant"
    Content   string                 `json:"content"`
    Timestamp time.Time              `json:"timestamp"`
    Metadata  map[string]interface{} `json:"metadata,omitempty"` // For tracking context, actions, etc.
}

// TableName specifies the table name for GORM
func (ConversationMemory) TableName() string {
    return "conversation_memories"
}

// BeforeSave hook to serialize messages to JSON before saving
func (c *ConversationMemory) BeforeSave(tx *gorm.DB) error {
    messagesJSON, err := json.Marshal(c.Messages)
    if err != nil {
        return err
    }
    c.MessagesJSON = string(messagesJSON)
    return nil
}

// AfterFind hook to deserialize messages from JSON after fetching
func (c *ConversationMemory) AfterFind(tx *gorm.DB) error {
    var messages []Message
    if err := json.Unmarshal([]byte(c.MessagesJSON), &messages); err != nil {
        return err
    }
    c.Messages = messages
    return nil
}
```

### 1.8 Chat API Endpoint

Implement the `/chat` endpoint that handles user messages and generates AI responses:

```go
// ChatHandler handles chat requests
func ChatHandler(
    aiSvc AIService,
) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Extract user ID from context (set by auth middleware)
        userID, ok := r.Context().Value("userID").(int)
        if !ok {
            http.Error(w, "User ID not found in context", http.StatusUnauthorized)
            return
        }

        // Parse request body
        var requestBody struct {
            Messages []struct {
                Role    string `json:"role"`
                Content string `json:"content"`
            } `json:"messages"`
            SessionID string `json:"session_id,omitempty"`
        }
        if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
            http.Error(w, "Invalid request body", http.StatusBadRequest)
            return
        }

        // Get the last user message
        var userMessage string
        for i := len(requestBody.Messages) - 1; i >= 0; i-- {
            if requestBody.Messages[i].Role == "user" {
                userMessage = requestBody.Messages[i].Content
                break
            }
        }
        if userMessage == "" {
            http.Error(w, "No user message provided", http.StatusBadRequest)
            return
        }

        // Generate session ID if not provided
        sessionID := requestBody.SessionID
        if sessionID == "" {
            sessionID = uuid.New().String()
        }

        // Convert messages to internal format
        var messages []Message
        for _, msg := range requestBody.Messages {
            messages = append(messages, Message{
                Role:      msg.Role,
                Content:   msg.Content,
                Timestamp: time.Now(),
            })
        }

        // Store conversation
        if err := aiSvc.StoreConversation(r.Context(), userID, sessionID, messages); err != nil {
            log.Printf("Error storing conversation: %v", err)
            // Continue without storing
        }

        // Generate response
        response, err := aiSvc.GenerateResponse(r.Context(), userID, userMessage)
        if err != nil {
            http.Error(w, fmt.Sprintf("Failed to generate response: %v", err), http.StatusInternalServerError)
            return
        }

        // Return response
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]interface{}{
            "output":     response,
            "session_id": sessionID,
        })
    }
}
```

---

## 2. Frontend Implementation with Drizzle ORM

### 2.1 Drizzle Schema for Conversation History

Define the schema for storing conversation history in the frontend:

```typescript
// src/db/schema/conversations.ts
import { sqliteTable, text, integer, blob } from 'drizzle-orm/sqlite-core';
import { sql } from 'drizzle-orm';
import { users } from './users';

export const conversations = sqliteTable('conversations', {
  id: text('id').primaryKey(),
  userId: text('user_id').notNull().references(() => users.id),
  title: text('title').notNull(),
  createdAt: integer('created_at', { mode: 'timestamp' }).notNull().default(sql`CURRENT_TIMESTAMP`),
  updatedAt: integer('updated_at', { mode: 'timestamp' }).notNull().default(sql`CURRENT_TIMESTAMP`),
});

export const conversationMessages = sqliteTable('conversation_messages', {
  id: text('id').primaryKey(),
  conversationId: text('conversation_id').notNull().references(() => conversations.id),
  role: text('role').notNull(), // 'user' or 'assistant'
  content: text('content').notNull(),
  timestamp: integer('timestamp', { mode: 'timestamp' }).notNull().default(sql`CURRENT_TIMESTAMP`),
  metadata: blob('metadata', { mode: 'json' }),
});
```

### 2.2 React Chat Interface with Vercel AI SDK

Create a responsive and interactive chat interface using the latest Vercel AI SDK features:

```jsx
// TradingAssistant.jsx
import { useChat } from 'ai/react';
import { useState } from 'react';
import { CryptoChart, PortfolioCard } from '../components';

export function TradingAssistant() {
  const [visualMode, setVisualMode] = useState('chat'); // 'chat', 'chart', 'portfolio'

  const { messages, input, handleInputChange, handleSubmit, isLoading } = useChat({
    api: '/api/chat',
    headers: {
      'Authorization': `Bearer ${localStorage.getItem('token')}`,
    },
    onFinish: (message) => {
      // Check if response contains chart data
      try {
        const data = JSON.parse(message.content);
        if (data.chartData) {
          setVisualMode('chart');
        } else if (data.portfolioSnapshot) {
          setVisualMode('portfolio');
        }
      } catch (e) {
        // Not JSON, regular chat message
        setVisualMode('chat');
      }
    },
  });

  return (
    <div className="trading-assistant">
      <div className="message-container">
        {messages.map((message) => (
          <div key={message.id} className={`message ${message.role}`}>
            {message.role === 'user' ? (
              <div className="user-message">{message.content}</div>
            ) : (
              <div className="ai-message">
                {visualMode === 'chat' && <div className="text-content">{message.content}</div>}
                {visualMode === 'chart' && <CryptoChart data={JSON.parse(message.content).chartData} />}
                {visualMode === 'portfolio' && <PortfolioCard data={JSON.parse(message.content).portfolioSnapshot} />}
              </div>
            )}
          </div>
        ))}
      </div>

      <form onSubmit={handleSubmit} className="input-form">
        <input
          value={input}
          onChange={handleInputChange}
          placeholder="Ask about your portfolio or trading strategies..."
          disabled={isLoading}
        />
        <button type="submit" disabled={isLoading || !input.trim()}>
          {isLoading ? 'Thinking...' : 'Send'}
        </button>
      </form>
    </div>
  );
}
```

### 2.3 Frontend API Client with Drizzle Integration

Create a client for interacting with the backend API:

```jsx
// api/aiClient.js
export async function sendChatMessage(message, sessionId = null) {
  const token = localStorage.getItem('token');
  if (!token) {
    throw new Error('Authentication required');
  }

  // Store message locally using Drizzle
  if (sessionId) {
    try {
      await db.insert(conversationMessages).values({
        id: crypto.randomUUID(),
        conversationId: sessionId,
        role: 'user',
        content: message,
        timestamp: new Date(),
      });
    } catch (error) {
      console.error('Failed to store message locally:', error);
      // Continue with API call even if local storage fails
    }
  }

  const response = await fetch('/api/chat', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`,
    },
    body: JSON.stringify({
      messages: [{ role: 'user', content: message }],
      session_id: sessionId,
    }),
  });

  if (!response.ok) {
    throw new Error(`API error: ${response.status}`);
  }

  const result = await response.json();

  // Store AI response locally using Drizzle
  if (sessionId) {
    try {
      await db.insert(conversationMessages).values({
        id: crypto.randomUUID(),
        conversationId: sessionId,
        role: 'assistant',
        content: result.output,
        timestamp: new Date(),
      });
    } catch (error) {
      console.error('Failed to store AI response locally:', error);
    }
  }

  return result;
}

export async function executeTradingFunction(functionName, parameters) {
  const token = localStorage.getItem('token');
  if (!token) {
    throw new Error('Authentication required');
  }

  const response = await fetch('/api/function', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`,
    },
    body: JSON.stringify({
      function_name: functionName,
      parameters: parameters,
    }),
  });

  if (!response.ok) {
    throw new Error(`API error: ${response.status}`);
  }

  return response.json();
}
```

### 2.4 Data Visualization Components

Create components for visualizing trading data:

```jsx
// components/CryptoChart.jsx
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts';

export function CryptoChart({ data }) {
  return (
    <div className="crypto-chart">
      <h3>{data.title}</h3>
      <ResponsiveContainer width="100%" height={300}>
        <LineChart data={data.points}>
          <CartesianGrid strokeDasharray="3 3" />
          <XAxis dataKey="time" />
          <YAxis />
          <Tooltip />
          <Legend />
          <Line type="monotone" dataKey="price" stroke="#8884d8" activeDot={{ r: 8 }} />
          {data.indicators && data.indicators.map((indicator, index) => (
            <Line
              key={indicator.name}
              type="monotone"
              dataKey={indicator.key}
              stroke={indicator.color}
              strokeDasharray={indicator.dashed ? "5 5" : "0"}
            />
          ))}
        </LineChart>
      </ResponsiveContainer>
    </div>
  );
}

// components/PortfolioCard.jsx
export function PortfolioCard({ data }) {
  return (
    <div className="portfolio-card">
      <h3>Portfolio Summary</h3>
      <div className="portfolio-value">
        <span className="label">Total Value:</span>
        <span className="value">${data.totalValue.toFixed(2)}</span>
        <span className={`change ${data.change >= 0 ? 'positive' : 'negative'}`}>
          {data.change >= 0 ? '+' : ''}{data.change.toFixed(2)}%
        </span>
      </div>
      <div className="assets">
        <h4>Assets</h4>
        <ul>
          {data.assets.map((asset) => (
            <li key={asset.symbol}>
              <span className="symbol">{asset.symbol}</span>
              <span className="amount">{asset.amount}</span>
              <span className="value">${asset.value.toFixed(2)}</span>
              <span className={`change ${asset.change >= 0 ? 'positive' : 'negative'}`}>
                {asset.change >= 0 ? '+' : ''}{asset.change.toFixed(2)}%
              </span>
            </li>
          ))}
        </ul>
      </div>
    </div>
  );
}
```
