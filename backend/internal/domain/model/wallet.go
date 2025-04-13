package model

import (
	"time"
)

// BalanceType represents the type of balance (available, frozen, etc.)
type BalanceType string

// Asset represents a cryptocurrency/token asset
type Asset string

// Balance type constants
const (
	BalanceTypeAvailable BalanceType = "AVAILABLE"
	BalanceTypeFrozen    BalanceType = "FROZEN"
	BalanceTypeTotal     BalanceType = "TOTAL"
)

// Common assets
const (
	AssetUSDT Asset = "USDT"
	AssetBTC  Asset = "BTC"
	AssetETH  Asset = "ETH"
)

// Balance represents a balance of a specific asset
type Balance struct {
	Asset    Asset   `json:"asset"`
	Free     float64 `json:"free"`     // Available balance
	Locked   float64 `json:"locked"`   // Frozen/locked balance
	Total    float64 `json:"total"`    // Total balance (free + locked)
	USDValue float64 `json:"usdValue"` // USD value of the total balance
}

// Wallet represents a user's wallet with multiple asset balances
type Wallet struct {
	UserID        string             `json:"userId"`
	Balances      map[Asset]*Balance `json:"balances"`
	TotalUSDValue float64            `json:"totalUsdValue"`
	LastUpdated   time.Time          `json:"lastUpdated"`
}

// BalanceHistory represents a historical record of balance for an asset
type BalanceHistory struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	Asset     Asset     `json:"asset"`
	Free      float64   `json:"free"`
	Locked    float64   `json:"locked"`
	Total     float64   `json:"total"`
	USDValue  float64   `json:"usdValue"`
	Timestamp time.Time `json:"timestamp"`
}

// NewWallet creates a new wallet for a user
func NewWallet(userID string) *Wallet {
	return &Wallet{
		UserID:      userID,
		Balances:    make(map[Asset]*Balance),
		LastUpdated: time.Now(),
	}
}

// UpdateBalance updates or adds a balance for an asset
func (w *Wallet) UpdateBalance(asset Asset, free, locked float64, usdValue float64) {
	balance, exists := w.Balances[asset]

	if !exists {
		balance = &Balance{
			Asset: asset,
		}
		w.Balances[asset] = balance
	}

	balance.Free = free
	balance.Locked = locked
	balance.Total = free + locked
	balance.USDValue = usdValue

	w.recalculateTotalUSDValue()
	w.LastUpdated = time.Now()
}

// GetBalance returns the balance for a specific asset
func (w *Wallet) GetBalance(asset Asset) *Balance {
	balance, exists := w.Balances[asset]
	if !exists {
		return &Balance{
			Asset:    asset,
			Free:     0,
			Locked:   0,
			Total:    0,
			USDValue: 0,
		}
	}

	return balance
}

// HasSufficientBalance checks if there's sufficient free balance for an asset
func (w *Wallet) HasSufficientBalance(asset Asset, requiredAmount float64) bool {
	balance := w.GetBalance(asset)
	return balance.Free >= requiredAmount
}

// recalculateTotalUSDValue recalculates the total USD value of all assets
func (w *Wallet) recalculateTotalUSDValue() {
	total := 0.0

	for _, balance := range w.Balances {
		total += balance.USDValue
	}

	w.TotalUSDValue = total
}
