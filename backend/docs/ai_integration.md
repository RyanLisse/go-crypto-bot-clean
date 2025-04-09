### Comprehensive Implementation Guide for React Frontend and Go Backend

This guide provides a step-by-step approach to implementing an AI-enhanced crypto trading bot. It integrates a Go backend with a Turso SQLite database and the Gemini Flash model, alongside a React frontend using the Vercel AI SDK for a seamless chat interface. The system allows users to interact with an AI agent to get information about their portfolio and trades, including daily summaries.

---

#### Step 1: Set Up the Go Backend

The backend will handle requests from the frontend, interact with the Turso database, and communicate with the Gemini Flash model via the Google Generative AI API.

##### 1.1 Install Necessary Packages
You'll need the following Go packages:
- **Database Connection:** `"github.com/mattn/go-sqlite3"` for connecting to the Turso SQLite database.
- **HTTP Handling:** Standard `"net/http"` package for handling HTTP requests.
- **Google Generative AI:** `"cloud.google.com/go/generative-ai"` for interacting with the Gemini Flash model.
- **CORS Handling:** `"github.com/rs/cors"` to enable Cross-Origin Resource Sharing (CORS) for frontend requests.

Install them using:
```bash
go get github.com/mattn/go-sqlite3
go get cloud.google.com/go/generative-ai
go get github.com/rs/cors
```

##### 1.2 Database Connection and Schema
- **Connection:** Use the SQLite driver to connect to the Turso database. Replace `"path/to/turso.db"` with your actual Turso connection string (refer to [Turso Documentation](https://docs.turso.tech) for details).
- **Schema:** Assume the database includes:
  - `portfolios`: Stores portfolio data (e.g., `user_id`, `portfolio_value`, `timestamp`).
  - `trades`: Stores trade history (e.g., `user_id`, `asset`, `quantity`, `price`, `timestamp`).

##### 1.3 Authentication Middleware
Implement middleware to verify JWT tokens and extract the user ID:
```go
package main

import (
    "context"
    "net/http"
    "strings"
    // Assume a package or function for JWT verification
    "yourpackage/auth"
)

func authMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "Missing authorization header", http.StatusUnauthorized)
            return
        }
        token := strings.TrimPrefix(authHeader, "Bearer ")
        userID, err := auth.VerifyToken(token) // Replace with your JWT verification logic
        if err != nil {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }
        ctx := context.WithValue(r.Context(), "userID", userID)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```
- Replace `auth.VerifyToken` with your actual JWT verification function, which should return the user ID (e.g., an integer).

##### 1.4 Implement the `/chat` Endpoint
The `/chat` endpoint (POST) will:
- Extract the user ID from the request context.
- Fetch the user's portfolio summary and recent trades.
- Construct a prompt with context and the user's message.
- Send it to the Gemini Flash model.
- Return the AI's response as JSON.

```go
import (
    "database/sql"
    "encoding/json"
    "fmt"
    "net/http"
    "os"
    "cloud.google.com/go/generative-ai"
    "google.golang.org/api/option"
)

func chatHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Extract user ID from context
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
        }
        if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
            http.Error(w, "Invalid request body", http.StatusBadRequest)
            return
        }
        // Get the last user message
        var userMessage string
        for _, msg := range requestBody.Messages {
            if msg.Role == "user" {
                userMessage = msg.Content
            }
        }
        if userMessage == "" {
            http.Error(w, "No user message provided", http.StatusBadRequest)
            return
        }

        // Fetch data from database
        portfolio, err := getPortfolioSummary(db, userID)
        if err != nil {
            http.Error(w, "Failed to fetch portfolio", http.StatusInternalServerError)
            return
        }
        trades, err := getRecentTrades(db, userID)
        if err != nil {
            http.Error(w, "Failed to fetch trades", http.StatusInternalServerError)
            return
        }

        // Construct prompt
        prompt := fmt.Sprintf(
            "You are a helpful assistant for a crypto trading bot. Here is the user's portfolio summary: %s. Here are their recent trades: %s. The user asked: %s",
            portfolio, trades, userMessage,
        )

        // Set up AI client
        ctx := r.Context()
        client, err := generativeai.NewClient(ctx, option.WithAPIKey(os.Getenv("GOOGLE_API_KEY")))
        if err != nil {
            http.Error(w, "Failed to create AI client", http.StatusInternalServerError)
            return
        }
        defer client.Close()

        model := client.GenerativeModel("gemini-flash")
        resp, err := model.GenerateContent(ctx, generativeai.Text(prompt))
        if err != nil {
            http.Error(w, "Failed to generate response", http.StatusInternalServerError)
            return
        }

        // Extract AI response
        aiResponse := string(resp.Candidates[0].Content.Parts[0].(generativeai.Text))

        // Return response in expected format
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]string{"output": aiResponse})
    }
}
```

##### 1.5 Helper Functions for Data Retrieval
- **Portfolio Summary:**
```go
func getPortfolioSummary(db *sql.DB, userID int) (string, error) {
    var portfolioValue float64
    err := db.QueryRow(
        "SELECT portfolio_value FROM portfolios WHERE user_id = ? ORDER BY timestamp DESC LIMIT 1",
        userID,
    ).Scan(&portfolioValue)
    if err != nil {
        return "", err
    }
    return fmt.Sprintf("Current portfolio value: $%.2f", portfolioValue), nil
}
```
- **Recent Trades:**
```go
func getRecentTrades(db *sql.DB, userID int) (string, error) {
    rows, err := db.Query(
        "SELECT asset, quantity, price, timestamp FROM trades WHERE user_id = ? ORDER BY timestamp DESC LIMIT 5",
        userID,
    )
    if err != nil {
        return "", err
    }
    defer rows.Close()

    var trades []string
    for rows.Next() {
        var asset, timestamp string
        var quantity, price float64
        if err := rows.Scan(&asset, &quantity, &price, &timestamp); err != nil {
            return "", err
        }
        trades = append(trades, fmt.Sprintf("Traded %.2f %s at $%.2f on %s", quantity, asset, price, timestamp))
    }
    if len(trades) == 0 {
        return "No recent trades.", nil
    }
    return "Recent trades: " + strings.Join(trades, ", "), nil
}
```

##### 1.6 Set Up the Server
Start the server with CORS and authentication middleware:
```go
package main

import (
    "database/sql"
    "log"
    "net/http"
    _ "github.com/mattn/go-sqlite3"
    "github.com/rs/cors"
)

func main() {
    // Connect to Turso database
    db, err := sql.Open("sqlite3", "path/to/turso.db")
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    defer db.Close()

    // Set up router
    router := http.NewServeMux()
    router.Handle("/chat", authMiddleware(http.HandlerFunc(chatHandler(db))))

    // Configure CORS
    c := cors.New(cors.Options{
        AllowedOrigins:   []string{"http://localhost:3000"},
        AllowedMethods:   []string{"POST"},
        AllowedHeaders:   []string{"Authorization", "Content-Type"},
        AllowCredentials: true,
    })
    handler := c.Handler(router)

    // Start server
    log.Println("Server starting on :8080")
    log.Fatal(http.ListenAndServe(":8080", handler))
}
```
- Set the `GOOGLE_API_KEY` environment variable with your Google Generative AI API key.

---

#### Step 2: Set Up the React Frontend

The frontend will provide a chat interface using the Vercel AI SDK, communicating with the Go backend.

##### 2.1 Install Vercel AI SDK
In your React project, install the SDK:
```bash
npm install @vercel/ai
```

##### 2.2 Implement the Chat Interface
Use the `useChat` hook to manage chat state and send requests to the backend:
```jsx
import React from 'react';
import { useChat } from '@vercel/ai/react';

function ChatComponent() {
  // Retrieve JWT token (e.g., from local storage)
  const token = localStorage.getItem('token');

  const { messages, input, handleInputChange, handleSubmit } = useChat({
    api: 'http://localhost:8080/chat',
    headers: {
      'Authorization': `Bearer ${token}`,
    },
  });

  return (
    <div style={{ maxWidth: '600px', margin: '20px auto' }}>
      <h2>Crypto Trading Bot Assistant</h2>
      <div style={{ border: '1px solid #ccc', padding: '10px', maxHeight: '400px', overflowY: 'auto' }}>
        {messages.map((m, i) => (
          <div key={i} style={{ margin: '10px 0' }}>
            <strong>{m.role === 'user' ? 'You' : 'Assistant'}:</strong> {m.content}
          </div>
        ))}
      </div>
      <form onSubmit={handleSubmit} style={{ marginTop: '20px' }}>
        <input
          value={input}
          onChange={handleInputChange}
          placeholder="Ask about your portfolio or trades"
          style={{ width: '80%', padding: '5px' }}
        />
        <button type="submit" style={{ padding: '5px 10px', marginLeft: '10px' }}>
          Send
        </button>
      </form>
    </div>
  );
}

export default ChatComponent;
```
- **Assumption:** The JWT token is stored in `localStorage` after user login. Adjust the token retrieval logic based on your authentication setup.

##### 2.3 Handle User Authentication
- Ensure users log in to obtain a JWT token, stored in `localStorage` or cookies.
- The token is sent with each request via the `headers` option in `useChat`.

---

#### Step 3: Enhancements and Considerations

##### 3.1 Daily Summaries
- **Approach:** Include daily summaries in the portfolio summary. Modify `getPortfolioSummary` to compute or fetch a daily summary (e.g., value change over the last 24 hours).
- **Example Modification:**
```go
func getPortfolioSummary(db *sql.DB, userID int) (string, error) {
    var currentValue, previousValue float64
    err := db.QueryRow(
        "SELECT portfolio_value FROM portfolios WHERE user_id = ? ORDER BY timestamp DESC LIMIT 1",
        userID,
    ).Scan(&currentValue)
    if err != nil {
        return "", err
    }
    // Fetch value from 24 hours ago (simplified; adjust query as needed)
    err = db.QueryRow(
        "SELECT portfolio_value FROM portfolios WHERE user_id = ? AND timestamp < datetime('now', '-24 hours') ORDER BY timestamp DESC LIMIT 1",
        userID,
    ).Scan(&previousValue)
    if err != nil && err != sql.ErrNoRows {
        return "", err
    }
    change := currentValue - previousValue
    return fmt.Sprintf("Current portfolio value: $%.2f. Daily summary: Value changed by $%.2f today.", currentValue, change), nil
}
```

##### 3.2 Conversation History
- **Stateless Approach:** Each request is independent (current implementation).
- **Stateful Enhancement:** Store messages in the database (e.g., a `chat_history` table) and include recent messages in the prompt for context. Manage token limits by truncating older messages.

##### 3.3 Security and Scaling
- **Security:** Validate all inputs, secure the API key, and ensure JWT tokens are properly signed and verified.
- **CORS:** In production, replace `"http://localhost:3000"` with your frontend domain.
- **Scaling:** Use connection pooling for the database and consider caching frequent queries.

---

#### Running the Application

1. **Backend:**
   - Set the `GOOGLE_API_KEY` environment variable:
     ```bash
     export GOOGLE_API_KEY="your-api-key"
     ```
   - Run the Go server:
     ```bash
     go run main.go
     ```

2. **Frontend:**
   - Start your React app (e.g., with `npm start` if using Create React App).
   - Ensure the backend is running at `http://localhost:8080`.

3. **Test:**
   - Open the frontend (e.g., `http://localhost:3000`), log in to set the token, and use the chat interface to ask questions like "What's my portfolio value?" or "Show my recent trades."

---

#### Conclusion
This guide provides a complete implementation for a React frontend and Go backend, integrating an AI agent with the Gemini Flash model and a Turso database. Users can interact with the bot to get real-time portfolio and trade information, enhanced with daily summaries. For further improvements, consider adding conversation history, optimizing data fetching, or deploying with proper security measures.

**Key References:**
- [Vercel AI SDK Documentation](https://sdk.vercel.ai/docs)
- [Google Generative AI Go Client](https://cloud.google.com/generative-ai/docs)
- [Turso Documentation](https://docs.turso.tech)