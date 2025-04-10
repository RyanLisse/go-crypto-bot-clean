package websocket

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/websocket"

	"go-crypto-bot-clean/backend/internal/core/newcoin"
	"go-crypto-bot-clean/backend/internal/core/trade"
	"go-crypto-bot-clean/backend/internal/domain/service"
)

// Handler manages WebSocket connections and broadcasts.
type Handler struct {
	hub             *Hub
	tradeService    trade.TradeService
	newCoinService  newcoin.NewCoinService
	exchangeService service.ExchangeService
	upgrader        websocket.Upgrader
	marketService   *MarketDataService
	ctx             context.Context
	cancel          context.CancelFunc
}

// NewHandler creates a new WebSocket handler with injected services.
func NewHandler(hub *Hub, tradeSvc trade.TradeService, newCoinService newcoin.NewCoinService, exchangeService service.ExchangeService) *Handler {
	ctx, cancel := context.WithCancel(context.Background())
	h := &Handler{
		hub:             hub,
		tradeService:    tradeSvc,
		newCoinService:  newCoinService,
		exchangeService: exchangeService,
		ctx:             ctx,
		cancel:          cancel,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins; restrict in production
			},
		},
	}

	return h
}

// StartMarketDataService initializes and starts the market data service
func (h *Handler) StartMarketDataService(symbols []string) {
	// Initialize market data service with provided symbols
	h.marketService = NewMarketDataService(h, h.exchangeService, symbols, 5*time.Second)

	// Start the market data service
	go h.marketService.Start(h.ctx)
}

// ServeWS handles WebSocket requests.
//
// @Summary      WebSocket Endpoint
// @Description  Upgrade HTTP connection to WebSocket for real-time updates
// @Tags         websocket
// @Produce      json
// @Success      101 {string} string "Switching Protocols"
// @Router       /ws [get]
func (h *Handler) ServeWS(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not upgrade connection", http.StatusInternalServerError)
		return
	}
	client := NewClient(h.hub, conn)
	h.hub.register <- client

	go client.WritePump()
	go client.ReadPump()
}

// BroadcastMarketData broadcasts market data to all clients.
func (h *Handler) BroadcastMarketData(data MarketDataPayload) {
	h.hub.Broadcast(WSMessage{
		Type:      MarketDataType,
		Timestamp: time.Now().Unix(),
		Payload:   data,
	})
}

// BroadcastTradeNotification broadcasts trade notifications.
func (h *Handler) BroadcastTradeNotification(data TradeNotificationPayload) {
	h.hub.Broadcast(WSMessage{
		Type:      TradeNotificationType,
		Timestamp: time.Now().Unix(),
		Payload:   data,
	})
}

// BroadcastNewCoinAlert broadcasts new coin alerts.
func (h *Handler) BroadcastNewCoinAlert(data NewCoinAlertPayload) {
	h.hub.Broadcast(WSMessage{
		Type:      NewCoinAlertType,
		Timestamp: time.Now().Unix(),
		Payload:   data,
	})
}

// BroadcastPortfolioUpdate broadcasts portfolio updates.
func (h *Handler) BroadcastPortfolioUpdate(data PortfolioUpdatePayload) {
	h.hub.Broadcast(WSMessage{
		Type:      PortfolioUpdateType,
		Timestamp: time.Now().Unix(),
		Payload:   data,
	})
}

// BroadcastTradeUpdate broadcasts trade updates.
func (h *Handler) BroadcastTradeUpdate(data TradeUpdatePayload) {
	h.hub.Broadcast(WSMessage{
		Type:      TradeUpdateType,
		Timestamp: time.Now().Unix(),
		Payload:   data,
	})
}

// Cleanup stops all background services and releases resources.
func (h *Handler) Cleanup() {
	if h.marketService != nil {
		h.marketService.Stop()
	}
	if h.cancel != nil {
		h.cancel()
	}
}
