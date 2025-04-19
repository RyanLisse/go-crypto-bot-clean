package entity

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/model/wallet"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// WalletEntity represents the wallet data model for GORM
type WalletEntity struct {
	ID        uuid.UUID `gorm:"type:text;primaryKey"`
	UserID    uuid.UUID `gorm:"type:text;not null;index"`
	Type      string    `gorm:"type:text;not null"`                                       // Corresponds to wallet.WalletType
	Status    string    `gorm:"type:text;not null"`                                       // Corresponds to wallet.VerificationStatus
	IsPrimary bool      `gorm:"not null;default:false;index:idx_user_primary,priority:1"` // Combined index for user primary lookup

	// Web3-specific fields
	Address *string `gorm:"index"` // Use pointer for nullability
	Network *string

	// Exchange-specific fields
	ExchangeID   *string `gorm:"index"` // Use pointer for nullability
	ExchangeName *string

	// Balances stored as JSON text
	Balances WalletBalances `gorm:"type:text"` // Custom type for JSON marshalling

	CreatedAt  time.Time
	UpdatedAt  time.Time
	VerifiedAt *time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"` // GORM soft delete support
}

// TableName specifies the table name for GORM
func (WalletEntity) TableName() string {
	return "wallets"
}

// --- JSON Handling for Balances ---

// WalletBalances is a custom type wrapping []*wallet.Balance for JSON marshalling/unmarshalling
type WalletBalances []*wallet.Balance

// Value implements the driver.Valuer interface for database serialization.
func (b WalletBalances) Value() (driver.Value, error) {
	if len(b) == 0 {
		// Store empty array as '[]' or NULL? GORM often prefers NULL for empty.
		// Let's try NULL first. If issues arise, use "[]".
		return nil, nil
		// return "[]", nil
	}
	jsonData, err := json.Marshal(b)
	if err != nil {
		return nil, err
	}
	return string(jsonData), nil // Store as JSON string
}

// Scan implements the sql.Scanner interface for database deserialization.
func (b *WalletBalances) Scan(value interface{}) error {
	if value == nil {
		*b = make([]*wallet.Balance, 0) // Treat NULL as empty slice
		return nil
	}

	bytes, ok := value.([]byte) // SQLite TEXT might come as []byte
	if !ok {
		str, okStr := value.(string) // Or potentially as string
		if !okStr {
			return errors.New("type assertion to []byte or string failed for WalletBalances")
		}
		bytes = []byte(str)
	}

	if len(bytes) == 0 || string(bytes) == "null" || string(bytes) == "[]" {
		*b = make([]*wallet.Balance, 0)
		return nil
	}

	return json.Unmarshal(bytes, b)
}

// --- Conversion Methods (Optional but Recommended) ---

// ToDomain converts WalletEntity GORM model to domain model wallet.Wallet
func (e *WalletEntity) ToDomain() *wallet.Wallet {
	domainWallet := &wallet.Wallet{
		ID:         e.ID,
		UserID:     e.UserID,
		Type:       wallet.WalletType(e.Type),           // Cast string to custom type
		Status:     wallet.VerificationStatus(e.Status), // Cast string to custom type
		IsPrimary:  e.IsPrimary,
		Balances:   e.Balances, // Use the scanned balances
		CreatedAt:  e.CreatedAt,
		UpdatedAt:  e.UpdatedAt,
		VerifiedAt: e.VerifiedAt,
	}
	if e.Address != nil {
		domainWallet.Address = *e.Address
	}
	if e.Network != nil {
		domainWallet.Network = *e.Network
	}
	if e.ExchangeID != nil {
		domainWallet.ExchangeID = *e.ExchangeID
	}
	if e.ExchangeName != nil {
		domainWallet.ExchangeName = *e.ExchangeName
	}
	return domainWallet
}

// FromDomain converts domain model wallet.Wallet to GORM entity WalletEntity
func FromDomain(w *wallet.Wallet) *WalletEntity {
	entity := &WalletEntity{
		ID:         w.ID,
		UserID:     w.UserID,
		Type:       string(w.Type),   // Cast custom type to string
		Status:     string(w.Status), // Cast custom type to string
		IsPrimary:  w.IsPrimary,
		Balances:   w.Balances, // Assign balances directly
		CreatedAt:  w.CreatedAt,
		UpdatedAt:  w.UpdatedAt,
		VerifiedAt: w.VerifiedAt,
	}
	// Handle nullable fields correctly
	if w.Address != "" {
		entity.Address = &w.Address
	}
	if w.Network != "" {
		entity.Network = &w.Network
	}
	if w.ExchangeID != "" {
		entity.ExchangeID = &w.ExchangeID
	}
	if w.ExchangeName != "" {
		entity.ExchangeName = &w.ExchangeName
	}
	// Handle potential zero time for CreatedAt/UpdatedAt if necessary
	if entity.CreatedAt.IsZero() {
		entity.CreatedAt = time.Now()
	}
	if entity.UpdatedAt.IsZero() {
		entity.UpdatedAt = time.Now()
	}
	return entity
}
