package handlers

import (
	"context"
	"net/http"
	"time"

	ws "go-crypto-bot-clean/backend/internal/api/websocket"
	"go-crypto-bot-clean/backend/internal/core/account"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

// WebSocketHandler handles WebSocket connections
type WebSocketHandler struct {
	hub            *ws.Hub
	upgrader       websocket.Upgrader
	logger         *zap.Logger
	accountService account.AccountService
}

// NewWebSocketHandler creates a new WebSocketHandler
func NewWebSocketHandler(hub *ws.Hub, logger *zap.Logger) *WebSocketHandler {
	return NewWebSocketHandlerWithAccountService(hub, nil, logger)
}

// NewWebSocketHandlerWithAccountService creates a new WebSocketHandler with an account service
func NewWebSocketHandlerWithAccountService(hub *ws.Hub, accountService account.AccountService, logger *zap.Logger) *WebSocketHandler {
	return &WebSocketHandler{
		hub: hub,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				// Allow all origins in development
				return true
			},
		},
		logger:         logger,
		accountService: accountService,
	}
}

// HandleWebSocket handles WebSocket connections
func (h *WebSocketHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.logger.Error("Failed to upgrade connection", zap.Error(err))
		return
	}

	// Create a new client
	client := ws.NewClient(h.hub, conn)

	// The client will be registered in the ReadPump method

	// Start the client's read and write pumps
	go client.WritePump()
	go client.ReadPump()

	h.logger.Info("WebSocket connection established", zap.String("addr", conn.RemoteAddr().String()))
}

// RegisterRoutes registers the WebSocket handler routes
func (h *WebSocketHandler) RegisterRoutes(router chi.Router) {
	router.Get("/ws", h.HandleWebSocket)

	// Start sending account updates if account service is available
	if h.accountService != nil {
		go h.startAccountUpdates()
	}
}

// startAccountUpdates starts sending account updates to WebSocket clients
func (h *WebSocketHandler) startAccountUpdates() {
	// Send initial account update
	h.sendAccountUpdate()

	// Send account updates every 5 seconds
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		h.sendAccountUpdate()
	}
}

// sendAccountUpdate sends an account update to all WebSocket clients
func (h *WebSocketHandler) sendAccountUpdate() {
	if h.accountService == nil {
		return
	}

	// Get wallet from account service
	wallet, err := h.accountService.GetWallet(context.Background())
	if err != nil {
		h.logger.Error("Failed to get wallet", zap.Error(err))
		return
	}

	// Convert wallet to account update payload
	balances := make(map[string]ws.AssetBalancePayload)
	for symbol, balance := range wallet.Balances {
		balances[symbol] = ws.AssetBalancePayload{
			Asset:  balance.Asset,
			Free:   balance.Free,
			Locked: balance.Locked,
			Total:  balance.Total,
		}
	}

	// Create account update message
	msg := ws.WSMessage{
		Type:      ws.AccountUpdateType,
		Timestamp: time.Now().Unix(),
		Payload: ws.AccountUpdatePayload{
			Balances:  balances,
			UpdatedAt: wallet.UpdatedAt.Unix(),
		},
	}

	// Broadcast message to all clients
	h.hub.Broadcast(msg)
}
