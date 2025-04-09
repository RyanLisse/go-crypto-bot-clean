package models

import "time"

// Balance represents the account balance information
type Balance struct {
	Fiat      float64            // Available fiat currency (e.g., USDT)
	Assets    map[string]float64 // Map of asset symbol to amount
	Available map[string]float64 // Available balance for each asset
	Locked    map[string]float64 // Locked balance for each asset
	UpdatedAt time.Time          // When the balance was last updated
}

// ExtendAssetBalance adds additional fields to AssetBalance
func ExtendAssetBalance(balance *AssetBalance, price float64) *AssetBalance {
	balance.Price = price
	balance.Total = balance.Free + balance.Locked
	return balance
}
