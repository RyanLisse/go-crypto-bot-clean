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

### 1.7 Conversation Memory System

Implement a sophisticated memory system for AI conversations to improve response quality:

```go
// ConversationMemory stores conversation history
type ConversationMemory struct {
    UserID       int       `json:"user_id"`
    SessionID    string    `json:"session_id"`
    Messages     []Message `json:"messages"`
    Summary      string    `json:"summary"`
    LastAccessed time.Time `json:"last_accessed"`
}

type Message struct {
    Role      string    `json:"role"` // "user" or "assistant"
    Content   string    `json:"content"`
    Timestamp time.Time `json:"timestamp"`
    Metadata  map[string]interface{} `json:"metadata,omitempty"` // For tracking context, actions, etc.
}

// Database schema for conversation memory
const createConversationMemoryTableSQL = `
CREATE TABLE IF NOT EXISTS conversation_memories (
    user_id INTEGER NOT NULL,
    session_id TEXT NOT NULL,
    messages_json TEXT NOT NULL,
    summary TEXT,
    last_accessed TIMESTAMP NOT NULL,
    PRIMARY KEY (user_id, session_id)
);
CREATE INDEX IF NOT EXISTS idx_conversation_memories_user_id ON conversation_memories(user_id);
CREATE INDEX IF NOT EXISTS idx_conversation_memories_last_accessed ON conversation_memories(last_accessed);
`
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

## 2. Frontend Implementation

### 2.1 React Chat Interface with Vercel AI SDK

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

### 2.2 Frontend API Client

Create a client for interacting with the backend API:

```jsx
// api/aiClient.js
export async function sendChatMessage(message, sessionId = null) {
  const token = localStorage.getItem('token');
  if (!token) {
    throw new Error('Authentication required');
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

  return response.json();
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

### 2.3 Data Visualization Components

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
