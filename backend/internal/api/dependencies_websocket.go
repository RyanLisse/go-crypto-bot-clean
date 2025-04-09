package api

import (
	"log"

	"go-crypto-bot-clean/backend/internal/api/websocket"
	"go-crypto-bot-clean/backend/internal/core/newcoin"
	"go-crypto-bot-clean/backend/internal/core/trade"
	"go-crypto-bot-clean/backend/internal/domain/service"
	"go-crypto-bot-clean/backend/internal/platform/mexc"
	"go.uber.org/zap"
)

// InitializeWebSocketDependencies initializes the WebSocket dependencies
func (d *Dependencies) InitializeWebSocketDependencies() {
	// Create logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Printf("Failed to create logger: %v", err)
		return
	}

	// Create WebSocket hub
	hub := websocket.NewHub()
	go hub.Run()

	// Create MEXC client for exchange service
	var exchangeService service.ExchangeService
	var tradeService trade.TradeService
	var newCoinService newcoin.NewCoinService

	// Try to create real exchange service
	exchangeService, err = mexc.NewClient(d.Config)
	if err != nil {
		logger.Error("Failed to create MEXC client for WebSocket", zap.Error(err))
		// Use mock services
		exchangeService = &mockExchangeService{}
		tradeService = &mockTradeService{}
		newCoinService = &mockNewCoinService{}
	}

	// Create WebSocket handler
	d.WebSocketHandler = websocket.NewHandler(hub, tradeService, newCoinService, exchangeService)

	// Start market data service
	d.WebSocketHandler.StartMarketDataService([]string{"BTCUSDT", "ETHUSDT"})
}
