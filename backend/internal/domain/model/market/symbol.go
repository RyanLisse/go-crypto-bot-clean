package market

import "time"

// PrecisionConfig defines the precision settings for a trading pair
type PrecisionConfig struct {
	PricePrecision    int // Number of decimal places for price
	QuantityPrecision int // Number of decimal places for quantity
	QuotePrecision    int // Number of decimal places for quote asset
}

// FilterConfig defines trading limits for a symbol
type FilterConfig struct {
	MinPrice    float64 // Minimum valid price for orders
	MaxPrice    float64 // Maximum valid price for orders
	MinQuantity float64 // Minimum quantity for orders
	MaxQuantity float64 // Maximum quantity for orders
	StepSize    float64 // Minimum quantity increment
	TickSize    float64 // Minimum price increment
}

// Symbol represents a trading pair on an exchange
type Symbol struct {
	// Symbol is the trading pair identifier (e.g., "BTCUSDT")
	Symbol string 

	// BaseAsset is the first part of the pair (e.g., "BTC")
	BaseAsset string 

	// QuoteAsset is the second part of the pair (e.g., "USDT")
	QuoteAsset string 

	// Exchange indicates which exchange this symbol is from
	Exchange string 

	// Status indicates if trading is enabled for this symbol
	Status string 

	// MinPrice is the minimum valid price for orders
	MinPrice float64 

	// MaxPrice is the maximum valid price for orders
	MaxPrice float64 

	// PricePrecision is the number of decimal places allowed for price
	PricePrecision int 

	// MinQty is the minimum quantity for orders
	MinQty float64 

	// MaxQty is the maximum quantity for orders
	MaxQty float64 

	// QtyPrecision is the number of decimal places allowed for quantity
	QtyPrecision int 

	// BaseAssetPrecision is the precision for the base asset
	BaseAssetPrecision int 

	// QuoteAssetPrecision is the precision for the quote asset
	QuoteAssetPrecision int 

	// MinNotional is the minimum order value (price * quantity)
	MinNotional float64 

	// MinLotSize is the minimum order quantity
	MinLotSize float64 

	// MaxLotSize is the maximum order quantity
	MaxLotSize float64 `json:"maxLotSize,omitempty"` // Added

	// StepSize defines allowed quantity increments
	StepSize float64 `json:"stepSize,omitempty"` // Added

	// TickSize defines allowed price increments
	TickSize float64 `json:"tickSize,omitempty"` // Added

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
