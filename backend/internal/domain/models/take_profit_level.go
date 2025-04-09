package models

import "time"

// TakeProfitLevel represents a level at which to take profit
type TakeProfitLevel struct {
	ID          string    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	PositionID  string    `gorm:"index;not null" json:"position_id"`
	Level       int       `gorm:"not null" json:"level"`
	Price       float64   `gorm:"not null" json:"price"`
	Percentage  float64   `gorm:"not null" json:"percentage"`
	Quantity    float64   `gorm:"not null" json:"quantity"`
	QuantityPct float64   `gorm:"not null" json:"quantity_pct"`
	Triggered   bool      `gorm:"not null;default:false" json:"triggered"`
	Executed    bool      `gorm:"not null;default:false" json:"executed"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
