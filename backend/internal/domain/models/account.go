package models

import "time"

type Wallet struct {
	Balances  map[string]*AssetBalance `json:"balances"`
	UpdatedAt time.Time                `json:"updatedAt"`
}
