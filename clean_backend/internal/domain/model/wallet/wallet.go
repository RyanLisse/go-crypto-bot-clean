package wallet

import (
	"time"

	"github.com/google/uuid"
)

// WalletType represents the type of wallet (Web3 or Exchange)
type WalletType string

const (
	// Web3Wallet represents a blockchain wallet
	Web3Wallet WalletType = "web3"
	// ExchangeWallet represents an exchange-connected wallet
	ExchangeWallet WalletType = "exchange"
)

// VerificationStatus represents the status of wallet verification
type VerificationStatus string

const (
	// Pending indicates the wallet is awaiting verification
	Pending VerificationStatus = "pending"
	// Verified indicates the wallet has been successfully verified
	Verified VerificationStatus = "verified"
	// Failed indicates the wallet verification has failed
	Failed VerificationStatus = "failed"
	// Verifying indicates the wallet verification is in progress (optional status)
	// Verifying VerificationStatus = "verifying"
)

// Type represents the type of wallet (e.g., exchange, web3)
type Type string

const (
	TypeExchange Type = "exchange"
	TypeWeb3     Type = "web3"
	TypeCustom   Type = "custom"
)

// Status represents the status of a wallet
type Status string

const (
	StatusActive   Status = "active"
	StatusInactive Status = "inactive"
	StatusPending  Status = "pending"
	StatusDisabled Status = "disabled"
	StatusVerified Status = "verified"
)

// SyncStatus represents the synchronization status of a wallet
type SyncStatus string

const (
	SyncStatusNone       SyncStatus = "none"
	SyncStatusScheduled  SyncStatus = "scheduled"
	SyncStatusInProgress SyncStatus = "in_progress"
	SyncStatusSuccess    SyncStatus = "success"
	SyncStatusFailed     SyncStatus = "failed"
)

// Metadata contains additional metadata for a wallet
type Metadata struct {
	Name        string            `json:"name,omitempty"`
	Description string            `json:"description,omitempty"`
	Tags        []string          `json:"tags,omitempty"`
	IsPrimary   bool              `json:"is_primary,omitempty"`
	Network     string            `json:"network,omitempty"`
	Address     string            `json:"address,omitempty"`
	ChainID     int64             `json:"chain_id,omitempty"`
	Explorer    string            `json:"explorer,omitempty"`
	Custom      map[string]string `json:"custom,omitempty"`
}

// Balance represents a token balance in a wallet
type Balance struct {
	Token     string    `json:"token"`
	Amount    string    `json:"amount"` // Using string for precision with large/small crypto amounts
	Value     float64   `json:"value"`  // Optional: Estimated fiat value (e.g., USD)
	UpdatedAt time.Time `json:"updated_at"`
}

// Wallet represents a user's cryptocurrency wallet
type Wallet struct {
	ID        uuid.UUID          `json:"id"`
	UserID    uuid.UUID          `json:"user_id"`
	Type      WalletType         `json:"type"`
	Status    VerificationStatus `json:"status"`
	IsPrimary bool               `json:"is_primary"`

	// Web3-specific fields
	Address string `json:"address,omitempty"`
	Network string `json:"network,omitempty"`

	// Exchange-specific fields
	ExchangeID   string `json:"exchange_id,omitempty"`
	ExchangeName string `json:"exchange_name,omitempty"`

	// Common fields
	Balances   []*Balance `json:"balances"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	VerifiedAt *time.Time `json:"verified_at,omitempty"`
}

// NewWeb3Wallet creates a new Web3 wallet instance
func NewWeb3Wallet(userID uuid.UUID, address, network string) *Wallet {
	now := time.Now()
	return &Wallet{
		ID:        uuid.New(),
		UserID:    userID,
		Type:      Web3Wallet,
		Status:    Pending,
		Address:   address,
		Network:   network,
		Balances:  make([]*Balance, 0),
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// NewExchangeWallet creates a new exchange wallet instance
func NewExchangeWallet(userID uuid.UUID, exchangeID, exchangeName string) *Wallet {
	now := time.Now()
	return &Wallet{
		ID:           uuid.New(),
		UserID:       userID,
		Type:         ExchangeWallet,
		Status:       Pending,
		ExchangeID:   exchangeID,
		ExchangeName: exchangeName,
		Balances:     make([]*Balance, 0),
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// SetVerified marks the wallet as verified
func (w *Wallet) SetVerified() {
	now := time.Now()
	w.Status = Verified
	w.VerifiedAt = &now
	w.UpdatedAt = now
}

// SetFailed marks the wallet as failed verification
func (w *Wallet) SetFailed() {
	w.Status = Failed
	w.UpdatedAt = time.Now()
}

// UpdateBalance updates or adds a new balance for a token
func (w *Wallet) UpdateBalance(token string, amount string, value float64) {
	now := time.Now()

	// Look for existing balance
	for _, bal := range w.Balances {
		if bal.Token == token {
			bal.Amount = amount
			bal.Value = value
			bal.UpdatedAt = now
			return
		}
	}

	// Add new balance
	w.Balances = append(w.Balances, &Balance{
		Token:     token,
		Amount:    amount,
		Value:     value,
		UpdatedAt: now,
	})
}

// GetBalance retrieves the balance for a specific token
func (w *Wallet) GetBalance(token string) *Balance {
	for _, bal := range w.Balances {
		if bal.Token == token {
			return bal
		}
	}
	return nil
}

// SetPrimary marks the wallet as primary
func (w *Wallet) SetPrimary(isPrimary bool) {
	w.IsPrimary = isPrimary
	w.UpdatedAt = time.Now()
}

// IsWeb3 checks if the wallet is a Web3 wallet
func (w *Wallet) IsWeb3() bool {
	return w.Type == Web3Wallet
}

// IsExchange checks if the wallet is an exchange wallet
func (w *Wallet) IsExchange() bool {
	return w.Type == ExchangeWallet
}

// IsVerified checks if the wallet is verified
func (w *Wallet) IsVerified() bool {
	return w.Status == Verified
}

// Validate performs basic validation of the wallet data
func (w *Wallet) Validate() error {
	if w.UserID == uuid.Nil {
		return ErrMissingUserID
	}

	switch w.Type {
	case Web3Wallet:
		if w.Address == "" {
			return ErrMissingAddress
		}
		if w.Network == "" {
			return ErrMissingNetwork
		}
	case ExchangeWallet:
		if w.ExchangeID == "" {
			return ErrMissingExchangeID
		}
	default:
		return ErrInvalidWalletType
	}

	return nil
}
