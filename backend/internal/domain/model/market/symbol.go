package market

import "time"

// PrecisionConfig defines the precision settings for a trading pair
type PrecisionConfig struct {
	PricePrecision    int `json:"price_precision"`    // Number of decimal places for price
	QuantityPrecision int `json:"quantity_precision"` // Number of decimal places for quantity
	QuotePrecision    int `json:"quote_precision"`    // Number of decimal places for quote asset
}

// FilterConfig defines trading limits for a symbol
type FilterConfig struct {
	MinPrice    float64 `json:"min_price"`
	MaxPrice    float64 `json:"max_price"`
	MinQuantity float64 `json:"min_quantity"`
	MaxQuantity float64 `json:"max_quantity"`
	StepSize    float64 `json:"step_size"` // Minimum quantity increment
	TickSize    float64 `json:"tick_size"` // Minimum price increment
}

// Symbol represents a trading pair on an exchange
type Symbol struct {
	// Symbol is the trading pair identifier (e.g., "BTCUSDT")
	Symbol string `json:"symbol"`

	// BaseAsset is the first part of the pair (e.g., "BTC")
	BaseAsset string `json:"baseAsset"`

	// QuoteAsset is the second part of the pair (e.g., "USDT")
	QuoteAsset string `json:"quoteAsset"`

	// Exchange indicates which exchange this symbol is from
	Exchange string `json:"exchange"`

	// Status indicates if trading is enabled for this symbol
	Status string `json:"status"`

	// MinPrice is the minimum valid price for orders
	MinPrice float64 `json:"minPrice"`

	// MaxPrice is the maximum valid price for orders
	MaxPrice float64 `json:"maxPrice"`

	// PricePrecision is the number of decimal places allowed for price
	PricePrecision int `json:"pricePrecision"`

	// MinQty is the minimum quantity for orders
	MinQty float64 `json:"minQty"`

	// MaxQty is the maximum quantity for orders
	MaxQty float64 `json:"maxQty"`

	// QtyPrecision is the number of decimal places allowed for quantity
	QtyPrecision int `json:"qtyPrecision"`

	// AllowedOrderTypes contains the order types supported for this symbol
	AllowedOrderTypes []string `json:"allowedOrderTypes"`

	// CreatedAt is when this symbol was added to our system
	CreatedAt time.Time `json:"createdAt"`

	// UpdatedAt is when this symbol was last updated
	UpdatedAt time.Time `json:"updatedAt"`
}

// SymbolInfo represents exchange info for a symbol, including status.
type SymbolInfo struct {
	Symbol string
	Status string // e.g., "TRADING", "AUCTION", "BREAK", etc.
	// Add other relevant fields as needed from exchange info API
}
