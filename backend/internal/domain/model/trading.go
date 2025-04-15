package model

import "time"

// TradingSchedule represents the listing and trading schedule for a symbol
type TradingSchedule struct {
	ListingTime time.Time `json:"listing_time"` // When the symbol was listed
	TradingTime time.Time `json:"trading_time"` // When trading begins
}

// SymbolConstraints represents trading constraints for a symbol
type SymbolConstraints struct {
	MinPrice   float64 `json:"min_price"`    // Minimum price
	MaxPrice   float64 `json:"max_price"`    // Maximum price
	MinQty     float64 `json:"min_qty"`      // Minimum quantity
	MaxQty     float64 `json:"max_qty"`      // Maximum quantity
	PriceScale int     `json:"price_scale"`  // Price precision
	QtyScale   int     `json:"qty_scale"`    // Quantity precision
}

// NewSymbolConstraints creates a new SymbolConstraints
func NewSymbolConstraints(minPrice, maxPrice, minQty, maxQty float64, priceScale, qtyScale int) *SymbolConstraints {
	return &SymbolConstraints{
		MinPrice:   minPrice,
		MaxPrice:   maxPrice,
		MinQty:     minQty,
		MaxQty:     maxQty,
		PriceScale: priceScale,
		QtyScale:   qtyScale,
	}
}
