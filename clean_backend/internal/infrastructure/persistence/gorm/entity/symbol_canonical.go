package entity

import (
	"strings"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/model"
)

// SymbolEntityCanonical represents symbol information
type SymbolEntityCanonical struct {
	Symbol              string    `gorm:"primaryKey;index:idx_symbols_canonical_symbol_exchange,priority:1"`
	Exchange            string    `gorm:"primaryKey;index:idx_symbols_canonical_symbol_exchange,priority:2"`
	BaseAsset           string    `gorm:"not null;index:idx_symbols_canonical_base_asset"`
	QuoteAsset          string    `gorm:"not null;index:idx_symbols_canonical_quote_asset"`
	Status              string    `gorm:"not null"` // e.g., "TRADING", "BREAK", etc.
	MinPrice            float64   `gorm:"not null"`
	MaxPrice            float64   `gorm:"not null"`
	PricePrecision      int       `gorm:"not null"`
	MinQuantity         float64   `gorm:"not null"`
	MaxQuantity         float64   `gorm:"not null"`
	QuantityPrecision   int       `gorm:"not null"`
	BaseAssetPrecision  int       `gorm:"not null"`
	QuoteAssetPrecision int       `gorm:"not null"`
	MinNotional         float64   `gorm:"not null"`
	StepSize            float64   `gorm:"not null;default:0"`
	TickSize            float64   `gorm:"not null;default:0"`
	AllowedOrderTypes   string    `gorm:"type:text"` // Comma-separated list of allowed order types
	ListingDate         time.Time `gorm:"default:null"`
	CreatedAt           time.Time `gorm:"not null;autoCreateTime"`
	UpdatedAt           time.Time `gorm:"not null;autoUpdateTime"`
	BridgeCM            string    `gorm:"not null"`
}

// TableName specifies the table name for SymbolEntityCanonical
func (SymbolEntityCanonical) TableName() string {
	return "symbols_canonical"
}

// FromModel converts a domain model to an entity
func (e *SymbolEntityCanonical) FromModel(symbol *model.Symbol) {
	e.Symbol = symbol.Symbol
	e.Exchange = symbol.Exchange
	e.BaseAsset = symbol.BaseAsset
	e.QuoteAsset = symbol.QuoteAsset
	e.Status = string(symbol.Status)
	e.MinPrice = symbol.MinPrice
	e.MaxPrice = symbol.MaxPrice
	e.PricePrecision = symbol.PricePrecision
	e.MinQuantity = symbol.MinQuantity
	e.MaxQuantity = symbol.MaxQuantity
	e.QuantityPrecision = symbol.QuantityPrecision
	e.BaseAssetPrecision = symbol.BaseAssetPrecision
	e.QuoteAssetPrecision = symbol.QuoteAssetPrecision
	e.MinNotional = symbol.MinNotional
	e.StepSize = symbol.StepSize
	e.TickSize = symbol.TickSize

	// Convert allowed order types to comma-separated string
	if len(symbol.AllowedOrderTypes) > 0 {
		e.AllowedOrderTypes = strings.Join(symbol.AllowedOrderTypes, ",")
	}
}

// ToModel converts an entity to a domain model
func (e *SymbolEntityCanonical) ToModel() *model.Symbol {
	symbol := &model.Symbol{
		Symbol:              e.Symbol,
		Exchange:            e.Exchange,
		BaseAsset:           e.BaseAsset,
		QuoteAsset:          e.QuoteAsset,
		Status:              model.SymbolStatus(e.Status),
		MinPrice:            e.MinPrice,
		MaxPrice:            e.MaxPrice,
		PricePrecision:      e.PricePrecision,
		MinQuantity:         e.MinQuantity,
		MaxQuantity:         e.MaxQuantity,
		QuantityPrecision:   e.QuantityPrecision,
		BaseAssetPrecision:  e.BaseAssetPrecision,
		QuoteAssetPrecision: e.QuoteAssetPrecision,
		MinNotional:         e.MinNotional,
		StepSize:            e.StepSize,
		TickSize:            e.TickSize,
		CreatedAt:           e.CreatedAt,
		UpdatedAt:           e.UpdatedAt,
	}

	// Convert comma-separated string to allowed order types
	if e.AllowedOrderTypes != "" {
		symbol.AllowedOrderTypes = strings.Split(e.AllowedOrderTypes, ",")
	} else {
		symbol.AllowedOrderTypes = []string{}
	}

	return symbol
}
