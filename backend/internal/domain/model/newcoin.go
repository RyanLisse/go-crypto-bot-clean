package model

import (
	"errors"
	"time"
)

// NewCoinStatus represents the status of a new coin listing.
type NewCoinStatus string

const (
	// StatusExpected means the coin is announced but not yet tradable.
	StatusExpected NewCoinStatus = "EXPECTED"
	// StatusTrading means the coin is now actively tradable.
	StatusTrading NewCoinStatus = "TRADING"
	// StatusProcessed means the coin has been detected as tradable and the event published.
	StatusProcessed NewCoinStatus = "PROCESSED"
)

// ErrNotFound is returned when a requested entity is not found.
// Consider moving this to a common errors package if used more broadly.
var ErrNotFound = errors.New("entity not found")

// NewCoin represents information about a newly listed coin on an exchange.
type NewCoin struct {
	ID                    string        `json:"id" gorm:"primaryKey"` // Use UUID or similar
	Symbol                string        `json:"symbol" gorm:"uniqueIndex"`
	ExpectedListingTime   time.Time     `json:"expected_listing_time"`
	Status                NewCoinStatus `json:"status" gorm:"index"`
	BecameTradableAt      *time.Time    `json:"became_tradable_at,omitempty"` // Pointer to allow null
	IsProcessedForAutobuy bool          `json:"is_processed_for_autobuy" gorm:"default:false"`
	CreatedAt             time.Time     `json:"created_at"`
	UpdatedAt             time.Time     `json:"updated_at"`
}

// MarkAsTradable updates the status and timestamp when a coin becomes tradable.
func (nc *NewCoin) MarkAsTradable(tradableTime time.Time) {
	if nc.Status == StatusExpected {
		nc.Status = StatusTrading
		nc.BecameTradableAt = &tradableTime
		nc.UpdatedAt = tradableTime // Also update UpdatedAt
	}
}

// MarkAsProcessed updates the status after the tradable event has been handled.
func (nc *NewCoin) MarkAsProcessed(processedTime time.Time) {
	if nc.Status == StatusTrading {
		nc.Status = StatusProcessed
		nc.IsProcessedForAutobuy = true // Assuming processing implies autobuy check
		nc.UpdatedAt = processedTime
	}
}
