package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
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

// Balance represents a token balance in a wallet
type Balance struct {
	Asset    Asset   `json:"asset"`
	Free     float64 `json:"free"`     // Available balance
	Locked   float64 `json:"locked"`   // Frozen/locked balance
	Total    float64 `json:"total"`    // Total balance (free + locked)
	USDValue float64 `json:"usdValue"` // USD value of the total balance
}

// WalletType represents the type of wallet (e.g., exchange, web3)
type WalletType string

const (
	WalletTypeExchange WalletType = "exchange"
	WalletTypeWeb3     WalletType = "web3"
)

// WalletStatus represents the status of a wallet
type WalletStatus string

const (
	WalletStatusActive   WalletStatus = "active"
	WalletStatusInactive WalletStatus = "inactive"
	WalletStatusPending  WalletStatus = "pending"
	WalletStatusDisabled WalletStatus = "disabled"
	WalletStatusVerified WalletStatus = "verified"
)

// VerificationStatus represents the verification status of a wallet
type VerificationStatus string

const (
	VerificationStatusNone     VerificationStatus = "none"
	VerificationStatusPending  VerificationStatus = "pending"
	VerificationStatusVerified VerificationStatus = "verified"
	VerificationStatusFailed   VerificationStatus = "failed"
)

// SyncStatus represents the synchronization status of a wallet
type SyncStatus string

// Wallet type constants
const (
	WalletTypeCustom WalletType = "CUSTOM" // Custom wallet type
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

// Wallet represents a user's cryptocurrency wallet
type Wallet struct {
	ID        uuid.UUID    `json:"id"`
	UserID    uuid.UUID    `json:"user_id"`
	Type      WalletType   `json:"type"`
	Status    WalletStatus `json:"status"`
	IsPrimary bool         `json:"is_primary"`

	// Web3-specific fields
	Address string `json:"address,omitempty"`
	Network string `json:"network,omitempty"`

	// Exchange-specific fields
	ExchangeID   string `json:"exchange_id,omitempty"`
	ExchangeName string `json:"exchange_name,omitempty"`

	// Common fields
	Balances      []*Balance             `json:"balances"`
	TotalUSDValue float64                `json:"total_usd_value"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
	VerifiedAt    *time.Time             `json:"verified_at,omitempty"`
}

// WalletMetadataUpdate represents the request body for updating wallet metadata
type WalletMetadataUpdate struct {
	Metadata map[string]interface{} `json:"metadata"`
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

// NewWallet creates a new wallet instance
func NewWallet(userID string) *Wallet {
	now := time.Now()
	uuidValue, err := uuid.Parse(userID)
	if err != nil {
		uuidValue = uuid.New() // Generate a new UUID if parsing fails
	}
	return &Wallet{
		ID:        uuid.New(),
		UserID:    uuidValue,
		Status:    WalletStatusPending,
		Metadata:  make(map[string]interface{}),
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// NewExchangeWallet creates a new exchange wallet
func NewExchangeWallet(userID uuid.UUID, exchangeID, exchangeName string) *Wallet {
	now := time.Now()
	return &Wallet{
		UserID:       userID,
		Type:         WalletTypeExchange,
		ExchangeID:   exchangeID,
		ExchangeName: exchangeName,
		Balances:     make([]*Balance, 0),
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// NewWeb3Wallet creates a new Web3 wallet
func NewWeb3Wallet(userID uuid.UUID, address, network string) *Wallet {
	now := time.Now()
	return &Wallet{
		UserID:    userID,
		Type:      WalletTypeWeb3,
		Address:   address,
		Network:   network,
		Balances:  make([]*Balance, 0),
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// SetPrimary sets this wallet as the primary wallet
func (w *Wallet) SetPrimary(isPrimary bool) {
	w.IsPrimary = isPrimary
	w.UpdatedAt = time.Now()
}

// SetMetadata updates the wallet's metadata
func (w *Wallet) SetMetadata(name, description string, tags []string) {
	if w.Metadata == nil {
		w.Metadata = make(map[string]interface{})
	}
	w.Metadata["name"] = name
	w.Metadata["description"] = description
	w.Metadata["tags"] = tags
	w.UpdatedAt = time.Now()
}

// AddCustomMetadata adds a custom metadata key-value pair
func (w *Wallet) AddCustomMetadata(key string, value interface{}) {
	if w.Metadata == nil {
		w.Metadata = make(map[string]interface{})
	}
	w.Metadata[key] = value
	w.UpdatedAt = time.Now()
}

// Validate checks if the wallet has all required fields
func (w *Wallet) Validate() error {
	if w.UserID == uuid.Nil {
		return errors.New("user ID is required")
	}

	if w.Type == "" {
		return errors.New("wallet type is required")
	}

	if w.Type == WalletTypeExchange && w.Address == "" {
		return errors.New("address is required for exchange wallets")
	}

	if w.Type == WalletTypeWeb3 {
		if w.Address == "" {
			return errors.New("address is required for Web3 wallets")
		}
		network, ok := w.Metadata["network"].(string)
		if !ok || network == "" {
			return errors.New("network is required for Web3 wallets")
		}
	}

	return nil
}

// UpdateBalance updates or adds a new balance for a token
func (w *Wallet) UpdateBalance(asset Asset, free, locked, usdValue float64) {
	found := false
	for _, bal := range w.Balances {
		if bal.Asset == asset {
			bal.Free = free
			bal.Locked = locked
			bal.Total = free + locked
			bal.USDValue = usdValue
			found = true
			break
		}
	}
	if !found {
		w.Balances = append(w.Balances, &Balance{
			Asset:    asset,
			Free:     free,
			Locked:   locked,
			Total:    free + locked,
			USDValue: usdValue,
		})
	}
	w.recalculateTotalUSDValue()
	w.UpdatedAt = time.Now()
}

// GetBalance retrieves the balance for a specific asset
func (w *Wallet) GetBalance(asset Asset) *Balance {
	for _, balance := range w.Balances {
		if balance.Asset == asset {
			return balance
		}
	}
	return nil
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
