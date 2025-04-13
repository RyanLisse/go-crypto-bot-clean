package model

import (
	"time"
)

// Account represents a user's exchange account information
type Account struct {
	UserID      string    `json:"userId"`
	Exchange    string    `json:"exchange"`
	Wallet      *Wallet   `json:"wallet"`
	Permissions []string  `json:"permissions"`
	LastUpdated time.Time `json:"lastUpdated"`
}

// MarketData represents aggregated market data for a trading pair
type MarketData struct {
	Symbol      string      `json:"symbol"`
	Ticker      *Ticker     `json:"ticker"`
	OrderBook   OrderBook   `json:"orderBook"`
	LastTrade   MarketTrade `json:"lastTrade"`
	LastUpdated time.Time   `json:"lastUpdated"`
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
func (m *MarketData) UpdateMarketData(ticker *Ticker, orderBook OrderBook, lastTrade MarketTrade) {
	m.Ticker = ticker
	m.OrderBook = orderBook
	m.LastTrade = lastTrade
	m.LastUpdated = time.Now()
}
