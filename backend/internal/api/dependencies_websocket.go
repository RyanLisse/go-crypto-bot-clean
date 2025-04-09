package api

import (
	"log"

	"github.com/ryanlisse/go-crypto-bot/internal/api/websocket"
	"github.com/ryanlisse/go-crypto-bot/internal/core/newcoin"
	"github.com/ryanlisse/go-crypto-bot/internal/core/trade"
	"github.com/ryanlisse/go-crypto-bot/internal/domain/service"
	"github.com/ryanlisse/go-crypto-bot/internal/platform/mexc"
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
