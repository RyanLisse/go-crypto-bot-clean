package models

import (
	"time"
)

// Ticker represents real-time price information for a trading pair
type Ticker struct {
	ID             string    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	Symbol         string    `gorm:"index;not null;size:20" json:"symbol"`
	Price          float64   `gorm:"not null" json:"price"`
	PriceChange    float64   `gorm:"not null" json:"priceChange"`
	PriceChangePct float64   `gorm:"not null" json:"priceChangePercent"`
	Volume         float64   `gorm:"not null" json:"volume"`
	QuoteVolume    float64   `gorm:"not null" json:"quoteVolume"`
	High24h        float64   `gorm:"not null" json:"high24h"`
	Low24h         float64   `gorm:"not null" json:"low24h"`
	Timestamp      time.Time `gorm:"index;not null" json:"timestamp"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
