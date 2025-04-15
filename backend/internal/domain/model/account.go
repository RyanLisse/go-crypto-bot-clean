package model

import (
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model/market"
)

// Account represents a user's exchange account information
type Account struct {
	UserID      string    
	Exchange    string    
	Wallet      *Wallet   
	Permissions []string  
	LastUpdated time.Time 
}

// MarketData represents aggregated market data for a trading pair
type MarketData struct {
	Symbol      string           `json:"symbol"`
	Ticker      *market.Ticker   `json:"ticker"`
	OrderBook   market.OrderBook `json:"orderBook"`
	LastTrade   MarketTrade      `json:"lastTrade"`
	LastUpdated time.Time        
}

// NewAccount creates a new exchange account for a user
func NewAccount(userID string, exchange string) *Account {
	return &Account{
		UserID:      userID,
		Exchange:    exchange,
		Wallet:      NewWallet(userID),
		Permissions: make([]string, 0),
		LastUpdated: time.Now(),
	}
}

// NewMarketData creates a new market data instance for a symbol
func NewMarketData(symbol string) *MarketData {
	return &MarketData{
		Symbol:      symbol,
		LastUpdated: time.Now(),
	}
}

// UpdateMarketData updates all market data components
func (m *MarketData) UpdateMarketData(ticker *market.Ticker, orderBook market.OrderBook, lastTrade MarketTrade) {
	m.Ticker = ticker
	m.OrderBook = orderBook
	m.LastTrade = lastTrade
	m.LastUpdated = time.Now()
}
