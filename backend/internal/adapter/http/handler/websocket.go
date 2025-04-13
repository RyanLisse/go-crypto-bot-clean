package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/neo/crypto-bot/internal/domain/model/market"
	"github.com/neo/crypto-bot/internal/domain/port"
	"github.com/rs/zerolog"
	"golang.org/x/time/rate"
)

// WebSocketHandler handles WebSocket connections for real-time market data
type WebSocketHandler struct {
	useCase       port.MarketDataUseCaseInterface
	logger        *zerolog.Logger
	clients       map[*websocket.Conn]*ClientInfo
	clientsMutex  sync.Mutex
	upgrader      websocket.Upgrader
	tickerChannel chan *market.Ticker
	stopChan      chan struct{}
}

// ClientInfo stores information about a connected WebSocket client
type ClientInfo struct {
	conn          *websocket.Conn
	limiter       *rate.Limiter
	subscriptions map[string][]string // map[channel][]symbols
	connectedAt   time.Time
	lastMessageAt time.Time
}

// WebSocketMessage represents a message sent over the WebSocket connection
type WebSocketMessage struct {
	Type    string      `json:"type"`
	Channel string      `json:"channel"`
	Symbol  string      `json:"symbol,omitempty"`
	Data    interface{} `json:"data"`
}

// SubscriptionRequest represents a subscription request from a client
type SubscriptionRequest struct {
	Action   string   `json:"action"`  // "subscribe" or "unsubscribe"
	Channel  string   `json:"channel"` // "tickers", "candles", etc.
	Symbols  []string `json:"symbols,omitempty"`
	Interval string   `json:"interval,omitempty"` // For candles
}

// NewWebSocketHandler creates a new WebSocketHandler
func NewWebSocketHandler(uc port.MarketDataUseCaseInterface, logger *zerolog.Logger) *WebSocketHandler {
	return &WebSocketHandler{
		useCase: uc,
		logger:  logger,
		clients: make(map[*websocket.Conn]*ClientInfo),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				// Allow all origins for now
				// In production, this should be restricted
				return true
			},
		},
		tickerChannel: make(chan *market.Ticker, 100),
		stopChan:      make(chan struct{}),
	}
}

// RegisterRoutes registers WebSocket routes with the Gin engine
func (h *WebSocketHandler) RegisterRoutes(router *gin.RouterGroup) {
	wsGroup := router.Group("/ws")
	{
		wsGroup.GET("/market", h.HandleWebSocket)
	}
}

// HandleWebSocket handles WebSocket connections
func (h *WebSocketHandler) HandleWebSocket(c *gin.Context) {
	// Upgrade HTTP connection to WebSocket
	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to upgrade connection to WebSocket")
		return
	}

	// Register client with rate limiter (10 messages per second, burst of 30)
	clientInfo := &ClientInfo{
		conn:          conn,
		limiter:       rate.NewLimiter(10, 30),
		subscriptions: make(map[string][]string),
		connectedAt:   time.Now(),
		lastMessageAt: time.Now(),
	}

	h.clientsMutex.Lock()
	h.clients[conn] = clientInfo
	h.clientsMutex.Unlock()

	// Send welcome message
	welcomeMsg := WebSocketMessage{
		Type:    "info",
		Channel: "system",
		Data:    "Connected to market data WebSocket",
	}
	if err := conn.WriteJSON(welcomeMsg); err != nil {
		h.logger.Error().Err(err).Msg("Failed to send welcome message")
		conn.Close()
		h.removeClient(conn)
		return
	}

	// Handle client messages in a goroutine
	go h.handleClientMessages(conn)
}

// handleClientMessages processes messages from a WebSocket client
func (h *WebSocketHandler) handleClientMessages(conn *websocket.Conn) {
	defer func() {
		conn.Close()
		h.removeClient(conn)
	}()

	// Set read deadline
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	// Process messages
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				h.logger.Error().Err(err).Msg("WebSocket connection closed unexpectedly")
			}
			break
		}

		// Apply rate limiting
		h.clientsMutex.Lock()
		clientInfo := h.clients[conn]
		h.clientsMutex.Unlock()

		if !clientInfo.limiter.Allow() {
			// Rate limit exceeded
			errorMsg := WebSocketMessage{
				Type:    "error",
				Channel: "system",
				Data:    "Rate limit exceeded. Please slow down your requests.",
			}
			conn.WriteJSON(errorMsg)
			continue
		}

		// Update last message time
		clientInfo.lastMessageAt = time.Now()

		// Parse subscription request
		var request SubscriptionRequest
		if err := json.Unmarshal(message, &request); err != nil {
			h.logger.Error().Err(err).Str("message", string(message)).Msg("Failed to parse subscription request")
			errorMsg := WebSocketMessage{
				Type:    "error",
				Channel: "system",
				Data:    "Invalid subscription request format",
			}
			conn.WriteJSON(errorMsg)
			continue
		}

		// Handle subscription request
		h.handleSubscription(conn, &request)
	}
}

// handleSubscription processes a subscription request
func (h *WebSocketHandler) handleSubscription(conn *websocket.Conn, request *SubscriptionRequest) {
	// Validate request
	if request.Action != "subscribe" && request.Action != "unsubscribe" {
		errorMsg := WebSocketMessage{
			Type:    "error",
			Channel: "system",
			Data:    "Invalid action. Must be 'subscribe' or 'unsubscribe'",
		}
		conn.WriteJSON(errorMsg)
		return
	}

	// Handle based on channel
	switch request.Channel {
	case "tickers":
		h.handleTickerSubscription(conn, request)
	case "candles":
		h.handleCandleSubscription(conn, request)
	default:
		errorMsg := WebSocketMessage{
			Type:    "error",
			Channel: "system",
			Data:    "Unsupported channel. Supported channels: 'tickers', 'candles'",
		}
		conn.WriteJSON(errorMsg)
	}
}

// handleTickerSubscription handles ticker subscription requests
func (h *WebSocketHandler) handleTickerSubscription(conn *websocket.Conn, request *SubscriptionRequest) {
	// Get client info
	h.clientsMutex.Lock()
	clientInfo := h.clients[conn]
	h.clientsMutex.Unlock()

	// Handle subscription/unsubscription
	if request.Action == "subscribe" {
		// Add subscription
		clientInfo.subscriptions["tickers"] = request.Symbols
	} else if request.Action == "unsubscribe" {
		// Remove subscription
		delete(clientInfo.subscriptions, "tickers")
	}

	// Acknowledge the subscription
	response := WebSocketMessage{
		Type:    "subscription",
		Channel: "tickers",
		Data: map[string]interface{}{
			"status":  "success",
			"action":  request.Action,
			"symbols": request.Symbols,
		},
	}
	conn.WriteJSON(response)

	// If subscribing, send initial ticker data
	if request.Action == "subscribe" && len(request.Symbols) > 0 {
		ctx := context.Background()
		for _, symbol := range request.Symbols {
			// Default to MEXC exchange for now
			ticker, err := h.useCase.GetTicker(ctx, "mexc", symbol)
			if err != nil {
				h.logger.Error().Err(err).Str("symbol", symbol).Msg("Failed to get ticker for initial data")
				continue
			}

			if ticker != nil {
				tickerMsg := WebSocketMessage{
					Type:    "data",
					Channel: "tickers",
					Symbol:  symbol,
					Data:    ticker,
				}
				conn.WriteJSON(tickerMsg)
			}
		}
	}
}

// handleCandleSubscription handles candle subscription requests
func (h *WebSocketHandler) handleCandleSubscription(conn *websocket.Conn, request *SubscriptionRequest) {
	// Validate interval
	if request.Interval == "" {
		errorMsg := WebSocketMessage{
			Type:    "error",
			Channel: "candles",
			Data:    "Interval is required for candle subscriptions",
		}
		conn.WriteJSON(errorMsg)
		return
	}

	// Get client info
	h.clientsMutex.Lock()
	clientInfo := h.clients[conn]
	h.clientsMutex.Unlock()

	// Create subscription key with interval
	subKey := "candles:" + request.Interval

	// Handle subscription/unsubscription
	if request.Action == "subscribe" {
		// Add subscription
		clientInfo.subscriptions[subKey] = request.Symbols
	} else if request.Action == "unsubscribe" {
		// Remove subscription
		delete(clientInfo.subscriptions, subKey)
	}

	// Acknowledge subscription
	response := WebSocketMessage{
		Type:    "subscription",
		Channel: "candles",
		Data: map[string]interface{}{
			"status":   "success",
			"action":   request.Action,
			"symbols":  request.Symbols,
			"interval": request.Interval,
		},
	}
	conn.WriteJSON(response)

	// If subscribing, send initial candle data
	if request.Action == "subscribe" && len(request.Symbols) > 0 {
		ctx := context.Background()
		interval := market.Interval(request.Interval)

		// Default time range: last 24 hours
		endTime := time.Now()
		startTime := endTime.Add(-24 * time.Hour)

		for _, symbol := range request.Symbols {
			// Default to MEXC exchange for now
			candles, err := h.useCase.GetCandles(ctx, "mexc", symbol, interval, startTime, endTime, 100)
			if err != nil {
				h.logger.Error().Err(err).Str("symbol", symbol).Msg("Failed to get candles for initial data")
				continue
			}

			candleMsg := WebSocketMessage{
				Type:    "data",
				Channel: "candles",
				Symbol:  symbol,
				Data:    candles,
			}
			conn.WriteJSON(candleMsg)
		}
	}
}

// Start starts the WebSocket handler
func (h *WebSocketHandler) Start() {
	// Start ticker broadcaster
	go h.broadcastTickers()
}

// Stop stops the WebSocket handler
func (h *WebSocketHandler) Stop() {
	close(h.stopChan)
}

// broadcastTickers periodically broadcasts ticker updates to all clients
func (h *WebSocketHandler) broadcastTickers() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-h.stopChan:
			return
		case <-ticker.C:
			// Get latest tickers
			ctx := context.Background()
			tickers, err := h.useCase.GetLatestTickers(ctx)
			if err != nil {
				h.logger.Error().Err(err).Msg("Failed to get latest tickers for broadcast")
				continue
			}

			// Broadcast to all clients
			h.clientsMutex.Lock()
			for conn, clientInfo := range h.clients {
				// Check if client is subscribed to tickers
				subscribedSymbols, hasTickerSubscription := clientInfo.subscriptions["tickers"]
				if !hasTickerSubscription {
					continue
				}

				// Send only subscribed tickers or all if subscribed to all
				for _, t := range tickers {
					// Check if client is subscribed to this symbol or to all symbols
					isSubscribed := len(subscribedSymbols) == 0 // Empty means all symbols
					if !isSubscribed {
						for _, symbol := range subscribedSymbols {
							if symbol == t.Symbol {
								isSubscribed = true
								break
							}
						}
					}

					if isSubscribed {
						tickerMsg := WebSocketMessage{
							Type:    "data",
							Channel: "tickers",
							Symbol:  t.Symbol,
							Data:    t,
						}
						if err := conn.WriteJSON(tickerMsg); err != nil {
							h.logger.Error().Err(err).Msg("Failed to send ticker update")
							conn.Close()
							delete(h.clients, conn)
						}
					}
				}
			}
			h.clientsMutex.Unlock()
		}
	}
}

// removeClient removes a client from the clients map
func (h *WebSocketHandler) removeClient(conn *websocket.Conn) {
	h.clientsMutex.Lock()
	clientInfo, exists := h.clients[conn]
	if exists {
		// Log client statistics
		duration := time.Since(clientInfo.connectedAt)
		h.logger.Info().
			Dur("connection_duration", duration).
			Time("connected_at", clientInfo.connectedAt).
			Time("last_message_at", clientInfo.lastMessageAt).
			Int("subscription_count", len(clientInfo.subscriptions)).
			Msg("Client disconnected")

		delete(h.clients, conn)
	}
	h.clientsMutex.Unlock()
}
