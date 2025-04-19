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

	// BaseAssetPrecision is the precision for the base asset
	BaseAssetPrecision int `json:"baseAssetPrecision"`

	// QuoteAssetPrecision is the precision for the quote asset
	QuoteAssetPrecision int `json:"quoteAssetPrecision"`

	// MinNotional is the minimum order value (price * quantity)
	MinNotional float64 `json:"minNotional"`

	// StepSize defines allowed quantity increments
	StepSize float64 `json:"stepSize"`

	// TickSize defines allowed price increments
	TickSize float64 `json:"tickSize"`

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
	return price >= s.MinPrice && (s.MaxPrice == 0 || price <= s.MaxPrice)
}

// ValidateQuantity checks if a quantity is within the allowed range and precision
func (s *Symbol) ValidateQuantity(quantity float64) bool {
	return quantity >= s.MinQuantity && (s.MaxQuantity == 0 || quantity <= s.MaxQuantity)
}

// ToMarketSymbol converts a Symbol to MarketSymbol format for backward compatibility
func (s *Symbol) ToMarketSymbol() *MarketSymbol {
	return &MarketSymbol{
		Symbol:              s.Symbol,
		BaseAsset:           s.BaseAsset,
		QuoteAsset:          s.QuoteAsset,
		Exchange:            s.Exchange,
		Status:              string(s.Status),
		MinPrice:            s.MinPrice,
		MaxPrice:            s.MaxPrice,
		PricePrecision:      s.PricePrecision,
		MinQty:              s.MinQuantity,
		MaxQty:              s.MaxQuantity,
		QtyPrecision:        s.QuantityPrecision,
		BaseAssetPrecision:  s.BaseAssetPrecision,
		QuoteAssetPrecision: s.QuoteAssetPrecision,
		MinNotional:         s.MinNotional,
		StepSize:            s.StepSize,
		TickSize:            s.TickSize,
		AllowedOrderTypes:   s.AllowedOrderTypes,
		CreatedAt:           s.CreatedAt,
		UpdatedAt:           s.UpdatedAt,
	}
}

// MarketSymbol represents the legacy market/symbol.go model
// This is provided for backward compatibility during transition
type MarketSymbol struct {
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

	// StepSize defines allowed quantity increments
	StepSize float64

	// TickSize defines allowed price increments
	TickSize float64

	// AllowedOrderTypes contains the order types supported for this symbol
	AllowedOrderTypes []string

	// CreatedAt is when this symbol was added to our system
	CreatedAt time.Time

	// UpdatedAt is when this symbol was last updated
	UpdatedAt time.Time
}

// ToSymbol converts a MarketSymbol to the canonical Symbol model
func (ms *MarketSymbol) ToSymbol() *Symbol {
	var status SymbolStatus = SymbolStatusHalt
	if ms.Status == string(SymbolStatusTrading) {
		status = SymbolStatusTrading
	} else if ms.Status == string(SymbolStatusBreak) {
		status = SymbolStatusBreak
	}

	return &Symbol{
		Symbol:              ms.Symbol,
		BaseAsset:           ms.BaseAsset,
		QuoteAsset:          ms.QuoteAsset,
		Exchange:            ms.Exchange,
		Status:              status,
		MinPrice:            ms.MinPrice,
		MaxPrice:            ms.MaxPrice,
		PricePrecision:      ms.PricePrecision,
		MinQuantity:         ms.MinQty,
		MaxQuantity:         ms.MaxQty,
		QuantityPrecision:   ms.QtyPrecision,
		BaseAssetPrecision:  ms.BaseAssetPrecision,
		QuoteAssetPrecision: ms.QuoteAssetPrecision,
		MinNotional:         ms.MinNotional,
		StepSize:            ms.StepSize,
		TickSize:            ms.TickSize,
		AllowedOrderTypes:   ms.AllowedOrderTypes,
		CreatedAt:           ms.CreatedAt,
		UpdatedAt:           ms.UpdatedAt,
	}
}
