package model

import (
	"time"
)

// SymbolStatus represents the trading status of a symbol
type SymbolStatus string

const (
	SymbolStatusTrading SymbolStatus = "TRADING"
	SymbolStatusHalt    SymbolStatus = "HALT"
	SymbolStatusBreak   SymbolStatus = "BREAK"
)

// Symbol represents a trading pair on an exchange
type Symbol struct {
	// Symbol is the trading pair identifier (e.g., "BTCUSDT")
	Symbol string 

	// BaseAsset is the first part of the pair (e.g., "BTC")
	BaseAsset string 

	// QuoteAsset is the second part of the pair (e.g., "USDT")
	QuoteAsset string 

	// Exchange indicates which exchange this symbol is from
	Exchange string `json:"exchange"`

	// Status indicates if trading is enabled for this symbol
	Status SymbolStatus 

	// MinPrice is the minimum valid price for orders
	MinPrice float64 `json:"minPrice"`

	// MaxPrice is the maximum valid price for orders
	MaxPrice float64 `json:"maxPrice"`

	// PricePrecision is the number of decimal places allowed for price
	PricePrecision int `json:"pricePrecision"`

	// MinQuantity is the minimum quantity for orders
	MinQuantity float64 `json:"minQuantity"`

	// MaxQuantity is the maximum quantity for orders
	MaxQuantity float64 `json:"maxQuantity"`

	// QuantityPrecision is the number of decimal places allowed for quantity
	QuantityPrecision int `json:"quantityPrecision"`

	// AllowedOrderTypes contains the order types supported for this symbol
	AllowedOrderTypes []string `json:"allowedOrderTypes"`

	// CreatedAt is when this symbol was added to our system
	CreatedAt time.Time `json:"createdAt"`

	// UpdatedAt is when this symbol was last updated
	UpdatedAt time.Time `json:"updatedAt"`
}

// IsActive returns true if the symbol is available for trading
func (s *Symbol) IsActive() bool {
	return s.Status == SymbolStatusTrading
}

// ValidatePrice checks if a price is within the allowed range and precision
func (s *Symbol) ValidatePrice(price float64) bool {
	return price >= s.MinPrice && price <= s.MaxPrice
}

// ValidateQuantity checks if a quantity is within the allowed range and precision
func (s *Symbol) ValidateQuantity(quantity float64) bool {
	return quantity >= s.MinQuantity && quantity <= s.MaxQuantity
}
