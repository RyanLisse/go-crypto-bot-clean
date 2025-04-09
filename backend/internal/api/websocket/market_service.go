package websocket

import (
	"context"
	"log"
	"time"

	"github.com/ryanlisse/go-crypto-bot/internal/domain/models"
	"github.com/ryanlisse/go-crypto-bot/internal/domain/service"
)

// MarketDataService handles fetching and broadcasting market data
type MarketDataService struct {
	handler        *Handler
	exchangeClient service.ExchangeService
	symbols        []string
	interval       time.Duration
	stopChan       chan struct{}
}

// NewMarketDataService creates a new market data service
func NewMarketDataService(handler *Handler, exchangeClient service.ExchangeService, symbols []string, interval time.Duration) *MarketDataService {
	return &MarketDataService{
		handler:        handler,
		exchangeClient: exchangeClient,
		symbols:        symbols,
		interval:       interval,
		stopChan:       make(chan struct{}),
	}
}

// Start begins the market data broadcasting service
func (s *MarketDataService) Start(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	// Broadcast initial data
	s.broadcastMarketData(ctx)

	for {
		select {
		case <-ticker.C:
			s.broadcastMarketData(ctx)
		case <-s.stopChan:
			log.Println("Market data service stopped")
			return
		case <-ctx.Done():
			log.Println("Market data service stopped due to context cancellation")
			return
		}
	}
}

// Stop stops the market data broadcasting service
func (s *MarketDataService) Stop() {
	close(s.stopChan)
}

// broadcastMarketData fetches and broadcasts market data for all tracked symbols
func (s *MarketDataService) broadcastMarketData(ctx context.Context) {
	for _, symbol := range s.symbols {
		ticker, err := s.exchangeClient.GetTicker(ctx, symbol)
		if err != nil {
			log.Printf("Error fetching ticker for %s: %v", symbol, err)
			continue
		}

		s.broadcastTicker(ticker)
	}
}

// broadcastTicker broadcasts a single ticker update
func (s *MarketDataService) broadcastTicker(ticker *models.Ticker) {
	payload := MarketDataPayload{
		Symbol:    ticker.Symbol,
		Price:     ticker.Price,
		Volume:    ticker.Volume,
		Timestamp: time.Now().Unix(),
	}

	s.handler.BroadcastMarketData(payload)
}

// AddSymbol adds a symbol to the tracked list
func (s *MarketDataService) AddSymbol(symbol string) {
	// Check if symbol already exists
	for _, existingSymbol := range s.symbols {
		if existingSymbol == symbol {
			return
		}
	}
	s.symbols = append(s.symbols, symbol)
}

// RemoveSymbol removes a symbol from the tracked list
func (s *MarketDataService) RemoveSymbol(symbol string) {
	for i, existingSymbol := range s.symbols {
		if existingSymbol == symbol {
			// Remove the symbol by replacing it with the last element and truncating
			s.symbols[i] = s.symbols[len(s.symbols)-1]
			s.symbols = s.symbols[:len(s.symbols)-1]
			return
		}
	}
}
