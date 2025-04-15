package entity

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// WalletMetadataEntity represents the metadata for a wallet
type WalletMetadataEntity struct {
	Name        string            `json:"name,omitempty"`
	Description string            `json:"description,omitempty"`
	Tags        []string          `json:"tags,omitempty"`
	IsPrimary   bool              `json:"is_primary,omitempty"`
	Network     string            `json:"network,omitempty"`
	Address     string            `json:"address,omitempty"`
	Custom      map[string]string `json:"custom,omitempty"`
}

// Value implements the driver.Valuer interface for WalletMetadataEntity
func (m WalletMetadataEntity) Value() (driver.Value, error) {
	return json.Marshal(m)
}

// Scan implements the sql.Scanner interface for WalletMetadataEntity
func (m *WalletMetadataEntity) Scan(value interface{}) error {
	if value == nil {
		*m = WalletMetadataEntity{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, m)
}

// EnhancedWalletEntity represents the database model for enhanced wallets
type EnhancedWalletEntity struct {
	ID            string               `gorm:"primaryKey;type:varchar(50)"`
	UserID        string               `gorm:"index;type:varchar(50);not null"`
	Exchange      string               `gorm:"type:varchar(50)"`
	Type          string               `gorm:"index;type:varchar(20);not null"`
	Status        string               `gorm:"index;type:varchar(20);not null"`
	TotalUSDValue float64              `gorm:"type:decimal(18,8);not null;default:0"`
	Metadata      WalletMetadataEntity `gorm:"type:json"`
	LastUpdated   time.Time            `gorm:"not null"`
	LastSyncAt    *time.Time
	CreatedAt     time.Time `gorm:"not null;autoCreateTime"`
	UpdatedAt     time.Time `gorm:"not null;autoUpdateTime"`
}

// TableName returns the table name for EnhancedWalletEntity
func (EnhancedWalletEntity) TableName() string {
	return "enhanced_wallets"
}

// EnhancedWalletBalanceEntity represents the database model for enhanced wallet balances
type EnhancedWalletBalanceEntity struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	WalletID  string    `gorm:"index;type:varchar(50);not null"`
	Asset     string    `gorm:"type:varchar(20);not null"`
	Free      float64   `gorm:"type:decimal(18,8);not null;default:0"`
	Locked    float64   `gorm:"type:decimal(18,8);not null;default:0"`
	Total     float64   `gorm:"type:decimal(18,8);not null;default:0"`
	USDValue  float64   `gorm:"type:decimal(18,8);not null;default:0"`
	CreatedAt time.Time `gorm:"not null;autoCreateTime"`
	UpdatedAt time.Time `gorm:"not null;autoUpdateTime"`
}

// TableName returns the table name for EnhancedWalletBalanceEntity
func (EnhancedWalletBalanceEntity) TableName() string {
	return "enhanced_wallet_balances"
}

// EnhancedWalletBalanceHistoryEntity represents the database model for enhanced wallet balance history
type EnhancedWalletBalanceHistoryEntity struct {
	ID            string    `gorm:"primaryKey;type:varchar(50)"`
	UserID        string    `gorm:"index;type:varchar(50);not null"`
	WalletID      string    `gorm:"index;type:varchar(50);not null"`
	BalancesJSON  []byte    `gorm:"type:json"`
	TotalUSDValue float64   `gorm:"type:decimal(18,8);not null;default:0"`
	Timestamp     time.Time `gorm:"index;not null"`
	CreatedAt     time.Time `gorm:"not null;autoCreateTime"`
}

// TableName returns the table name for EnhancedWalletBalanceHistoryEntity
func (EnhancedWalletBalanceHistoryEntity) TableName() string {
	return "enhanced_wallet_balance_history"
}
