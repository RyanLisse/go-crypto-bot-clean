package model

import (
	"errors"
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

// WalletType represents the type of wallet
type WalletType string

// WalletStatus represents the status of a wallet
type WalletStatus string

// SyncStatus represents the synchronization status of a wallet
type SyncStatus string

// Wallet type constants
const (
	WalletTypeExchange WalletType = "EXCHANGE" // Exchange wallet (e.g., MEXC, Binance)
	WalletTypeWeb3     WalletType = "WEB3"     // Web3 wallet (e.g., MetaMask, Trust Wallet)
	WalletTypeCustom   WalletType = "CUSTOM"   // Custom wallet type
)

// Wallet status constants
const (
	WalletStatusActive   WalletStatus = "ACTIVE"   // Wallet is active and can be used
	WalletStatusInactive WalletStatus = "INACTIVE" // Wallet is inactive and should not be used
	WalletStatusPending  WalletStatus = "PENDING"  // Wallet is pending activation
	WalletStatusFailed   WalletStatus = "FAILED"   // Wallet connection failed
	WalletStatusVerified WalletStatus = "VERIFIED" // Wallet has been verified
)

// Sync status constants
const (
	SyncStatusNone       SyncStatus = "NONE"        // Wallet has never been synced
	SyncStatusScheduled  SyncStatus = "SCHEDULED"   // Wallet sync is scheduled
	SyncStatusInProgress SyncStatus = "IN_PROGRESS" // Wallet sync is in progress
	SyncStatusSuccess    SyncStatus = "SUCCESS"     // Wallet sync completed successfully
	SyncStatusFailed     SyncStatus = "FAILED"      // Wallet sync failed
)

// WalletMetadata contains additional metadata for a wallet
type WalletMetadata struct {
	Name        string            `json:"name,omitempty"`        // User-defined name for the wallet
	Description string            `json:"description,omitempty"` // User-defined description
	Tags        []string          `json:"tags,omitempty"`        // Tags for categorizing wallets
	IsPrimary   bool              `json:"is_primary,omitempty"`  // Whether this is the primary wallet
	Network     string            `json:"network,omitempty"`     // Network for Web3 wallets (e.g., Ethereum, Binance Smart Chain)
	Address     string            `json:"address,omitempty"`     // Address for Web3 wallets
	ChainID     int64             `json:"chain_id,omitempty"`    // Chain ID for Web3 wallets
	Explorer    string            `json:"explorer,omitempty"`    // Block explorer URL for Web3 wallets
	Custom      map[string]string `json:"custom,omitempty"`      // Custom metadata
}

// Wallet represents a user's wallet with multiple asset balances
type Wallet struct {
	ID            string             // Unique identifier
	UserID        string             // User ID that owns this wallet
	Exchange      string             // Exchange name (for exchange wallets)
	Type          WalletType         // Type of wallet
	Status        WalletStatus       // Status of wallet
	SyncStatus    SyncStatus         // Synchronization status
	Balances      map[Asset]*Balance // Map of asset to Balance struct
	TotalUSDValue float64            // Total USD value of all balances
	Metadata      *WalletMetadata    // Additional metadata
	LastUpdated   time.Time          // When the wallet was last updated
	LastSynced    *time.Time         // When the wallet was last synced with the exchange
	LastSyncAt    time.Time          // When the wallet was last synced with the exchange (for database compatibility)
	CreatedAt     time.Time          // When the wallet was created
	UpdatedAt     time.Time          // When the wallet was last updated in the database
	Network       string             // Network for Web3 wallets (e.g., Ethereum, Binance Smart Chain)

}

// BalanceHistory represents a historical record of balance for a wallet
type BalanceHistory struct {
	ID            string             `json:"id"`
	UserID        string             `json:"user_id"`
	WalletID      string             `json:"wallet_id"`
	Balances      map[Asset]*Balance `json:"balances"`
	TotalUSDValue float64            `json:"total_usd_value"`
	Timestamp     time.Time          `json:"timestamp"`
}

// NewWallet creates a new wallet for a user
func NewWallet(userID string) *Wallet {
	now := time.Now()
	return &Wallet{
		ID:          GenerateID(),
		UserID:      userID,
		Type:        WalletTypeExchange,
		Status:      WalletStatusActive,
		SyncStatus:  SyncStatusNone,
		Balances:    make(map[Asset]*Balance),
		Metadata:    &WalletMetadata{},
		LastUpdated: now,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// NewExchangeWallet creates a new exchange wallet for a user
func NewExchangeWallet(userID, exchange string) *Wallet {
	wallet := NewWallet(userID)
	wallet.Exchange = exchange
	wallet.Type = WalletTypeExchange
	return wallet
}

// NewWeb3Wallet creates a new Web3 wallet for a user
func NewWeb3Wallet(userID, network, address string) *Wallet {
	wallet := NewWallet(userID)
	wallet.Type = WalletTypeWeb3
	wallet.Metadata.Network = network
	wallet.Metadata.Address = address
	return wallet
}

// SetPrimary sets this wallet as the primary wallet
func (w *Wallet) SetPrimary(isPrimary bool) {
	if w.Metadata == nil {
		w.Metadata = &WalletMetadata{}
	}
	w.Metadata.IsPrimary = isPrimary
	w.UpdatedAt = time.Now()
}

// SetMetadata sets the wallet metadata
func (w *Wallet) SetMetadata(name, description string, tags []string) {
	if w.Metadata == nil {
		w.Metadata = &WalletMetadata{}
	}
	w.Metadata.Name = name
	w.Metadata.Description = description
	w.Metadata.Tags = tags
	w.UpdatedAt = time.Now()
}

// AddCustomMetadata adds a custom metadata key-value pair
func (w *Wallet) AddCustomMetadata(key, value string) {
	if w.Metadata == nil {
		w.Metadata = &WalletMetadata{}
	}
	if w.Metadata.Custom == nil {
		w.Metadata.Custom = make(map[string]string)
	}
	w.Metadata.Custom[key] = value
	w.UpdatedAt = time.Now()
}

// Validate validates the wallet
func (w *Wallet) Validate() error {
	if w.UserID == "" {
		return errors.New("user ID is required")
	}

	if w.Type == "" {
		return errors.New("wallet type is required")
	}

	if w.Type == WalletTypeExchange && w.Exchange == "" {
		return errors.New("exchange is required for exchange wallets")
	}

	if w.Type == WalletTypeWeb3 {
		if w.Metadata == nil || w.Metadata.Address == "" {
			return errors.New("address is required for Web3 wallets")
		}
	}

	return nil
}

// UpdateBalance updates or adds a balance for an asset
func (w *Wallet) UpdateBalance(asset Asset, free, locked, usdValue float64) {
	w.Balances[asset] = &Balance{
		Asset:    asset,
		Free:     free,
		Locked:   locked,
		Total:    free + locked,
		USDValue: usdValue,
	}
	w.recalculateTotalUSDValue()
	w.LastUpdated = time.Now()
}

// GetBalance returns the balance for a specific asset
func (w *Wallet) GetBalance(asset Asset) *Balance {
	balance, exists := w.Balances[asset]
	if !exists {
		return nil
	}
	return balance
}

// HasSufficientBalance checks if there's sufficient balance for an asset
func (w *Wallet) HasSufficientBalance(asset Asset, requiredAmount float64) bool {
	balance := w.GetBalance(asset)
	if balance == nil {
		return false
	}
	return balance.Free >= requiredAmount
}

// recalculateTotalUSDValue recalculates the total USD value of all assets
func (w *Wallet) recalculateTotalUSDValue() {
	total := 0.0
	for _, balance := range w.Balances {
		if balance != nil {
			total += balance.USDValue
		}
	}
	w.TotalUSDValue = total
	w.UpdatedAt = time.Now()
}

// GenerateID generates a unique ID for a wallet
func GenerateID() string {
	return "wlt_" + generateUUID()
}

// generateUUID generates a UUID
func generateUUID() string {
	// This is a placeholder - in a real implementation, use a proper UUID library
	return time.Now().Format("20060102150405") + randomString(8)
}

// randomString generates a random string of the specified length
func randomString(length int) string {
	// This is a placeholder - in a real implementation, use a proper random string generator
	return "abcdefgh"[:length]
}
