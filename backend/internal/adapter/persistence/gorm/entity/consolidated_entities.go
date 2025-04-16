package entity

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
)

// ===== API Credential Entities =====

// APICredentialEntity represents the database model for API credentials
type APICredentialEntity struct {
	ID           string     `gorm:"primaryKey;type:varchar(50)"`
	UserID       string     `gorm:"not null;index;type:varchar(50)"`
	Exchange     string     `gorm:"not null;index;type:varchar(20)"`
	APIKey       string     `gorm:"not null;type:varchar(100)"`
	APISecret    []byte     `gorm:"not null;type:blob"` // Encrypted
	Label        string     `gorm:"type:varchar(50)"`
	Status       string     `gorm:"type:varchar(20);not null;default:'active'"`
	LastUsed     *time.Time `gorm:"column:last_used"`
	LastVerified *time.Time `gorm:"column:last_verified"`
	ExpiresAt    *time.Time `gorm:"column:expires_at"`
	RotationDue  *time.Time `gorm:"column:rotation_due"`
	FailureCount int        `gorm:"not null;default:0"`
	Metadata     []byte     `gorm:"type:json"`
	CreatedAt    time.Time  `gorm:"autoCreateTime"`
	UpdatedAt    time.Time  `gorm:"autoUpdateTime"`
}

// TableName returns the table name for the APICredentialEntity
func (APICredentialEntity) TableName() string {
	return "api_credentials"
}

// ===== User Entities =====

// UserEntity represents the database model for users
type UserEntity struct {
	ID        string    `gorm:"primaryKey;type:varchar(50)"`
	Email     string    `gorm:"uniqueIndex;type:varchar(100);not null"`
	Name      string    `gorm:"type:varchar(100)"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// TableName returns the table name for the UserEntity
func (UserEntity) TableName() string {
	return "users"
}

// ===== Wallet Entities =====

// WalletEntity represents a wallet in the database (legacy)
type WalletEntity struct {
	ID         string    `gorm:"primaryKey"`
	AccountID  string    `gorm:"not null;index"`
	Exchange   string    `gorm:"not null"`
	TotalUSD   float64   `gorm:"not null"`
	LastUpdate time.Time `gorm:"autoUpdateTime"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`
}

func (WalletEntity) TableName() string { return "legacy_wallets" }

// WalletMetadata represents the metadata for a wallet
type WalletMetadata struct {
	Name        string            `json:"name,omitempty"`
	Description string            `json:"description,omitempty"`
	Tags        []string          `json:"tags,omitempty"`
	IsPrimary   bool              `json:"is_primary,omitempty"`
	Network     string            `json:"network,omitempty"`
	Address     string            `json:"address,omitempty"`
	Custom      map[string]string `json:"custom,omitempty"`
}

// Value implements the driver.Valuer interface for WalletMetadata
func (m WalletMetadata) Value() (driver.Value, error) {
	return json.Marshal(m)
}

// Scan implements the sql.Scanner interface for WalletMetadata
func (m *WalletMetadata) Scan(value interface{}) error {
	if value == nil {
		*m = WalletMetadata{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, m)
}

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

// EnhancedWalletEntity represents a wallet in the database
type EnhancedWalletEntity struct {
	ID            string               `gorm:"primaryKey;type:varchar(50)"`
	UserID        string               `gorm:"index;type:varchar(50);not null"`
	Exchange      string               `gorm:"index;type:varchar(50)"`
	Type          string               `gorm:"index;type:varchar(20);not null"`
	Status        string               `gorm:"index;type:varchar(20);not null"`
	TotalUSDValue float64              `gorm:"type:decimal(18,8);not null;default:0"`
	Metadata      WalletMetadataEntity `gorm:"type:json"`
	LastUpdated   time.Time            `gorm:"not null"`
	LastSyncAt    *time.Time
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`
}

// TableName returns the table name for EnhancedWalletEntity
func (EnhancedWalletEntity) TableName() string {
	return "enhanced_wallets"
}

// EnhancedWalletBalanceEntity represents a balance in the database
type EnhancedWalletBalanceEntity struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	WalletID  string    `gorm:"index;type:varchar(50);not null"`
	Asset     string    `gorm:"index;type:varchar(20);not null"`
	Free      float64   `gorm:"type:decimal(18,8);not null;default:0"`
	Locked    float64   `gorm:"type:decimal(18,8);not null;default:0"`
	Total     float64   `gorm:"type:decimal(18,8);not null;default:0"`
	USDValue  float64   `gorm:"type:decimal(18,8);not null;default:0"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// TableName returns the table name for EnhancedWalletBalanceEntity
func (EnhancedWalletBalanceEntity) TableName() string {
	return "enhanced_wallet_balances"
}

// EnhancedWalletBalanceHistoryEntity represents a balance history record in the database
type EnhancedWalletBalanceHistoryEntity struct {
	ID            string    `gorm:"primaryKey;type:varchar(50)"`
	UserID        string    `gorm:"index;type:varchar(50);not null"`
	WalletID      string    `gorm:"index;type:varchar(50);not null"`
	BalancesJSON  []byte    `gorm:"type:json"`
	TotalUSDValue float64   `gorm:"type:decimal(18,8);not null;default:0"`
	Timestamp     time.Time `gorm:"index;not null"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
}

// TableName returns the table name for EnhancedWalletBalanceHistoryEntity
func (EnhancedWalletBalanceHistoryEntity) TableName() string {
	return "enhanced_wallet_balance_history"
}

// Wallet represents a wallet in the database
type Wallet struct {
	ID        string         `gorm:"primaryKey"`
	UserID    string         `gorm:"not null;index"`
	Exchange  string         `gorm:"not null"`
	Name      string         `gorm:"not null"`
	IsActive  bool           `gorm:"not null;default:true"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TableName sets the table name for Wallet
func (Wallet) TableName() string {
	return "wallets"
}

// WalletBalance represents a balance in the database
type WalletBalance struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	WalletID  string    `gorm:"index;type:varchar(50);not null"`
	Asset     string    `gorm:"index;type:varchar(20);not null"`
	Free      float64   `gorm:"type:decimal(18,8);not null;default:0"`
	Locked    float64   `gorm:"type:decimal(18,8);not null;default:0"`
	Total     float64   `gorm:"type:decimal(18,8);not null;default:0"`
	USDValue  float64   `gorm:"type:decimal(18,8);not null;default:0"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// TableName sets the table name for WalletBalance
func (WalletBalance) TableName() string {
	return "wallet_balances"
}

// WalletBalanceHistory represents a balance history record in the database
type WalletBalanceHistory struct {
	ID            string    `gorm:"primaryKey;type:varchar(50)"`
	UserID        string    `gorm:"index;type:varchar(50);not null"`
	WalletID      string    `gorm:"index;type:varchar(50);not null"`
	BalancesJSON  []byte    `gorm:"type:json"`
	TotalUSDValue float64   `gorm:"type:decimal(18,8);not null;default:0"`
	Timestamp     time.Time `gorm:"index;not null"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
}

// TableName sets the table name for WalletBalanceHistory
func (WalletBalanceHistory) TableName() string {
	return "wallet_balance_history"
}

// Symbol represents a trading symbol in the database
type Symbol struct {
	ID             string    `gorm:"primaryKey;type:varchar(50)"`
	Symbol         string    `gorm:"uniqueIndex;type:varchar(20);not null"`
	Exchange       string    `gorm:"index;type:varchar(20);not null"`
	BaseAsset      string    `gorm:"index;type:varchar(20);not null"`
	QuoteAsset     string    `gorm:"index;type:varchar(20);not null"`
	Status         string    `gorm:"type:varchar(20);not null"`
	MinPrice       float64   `gorm:"type:decimal(18,8);not null"`
	MaxPrice       float64   `gorm:"type:decimal(18,8);not null"`
	PricePrecision int       `gorm:"not null"`
	MinQty         float64   `gorm:"type:decimal(18,8);not null"`
	MaxQty         float64   `gorm:"type:decimal(18,8);not null"`
	QtyPrecision   int       `gorm:"not null"`
	CreatedAt      time.Time `gorm:"autoCreateTime"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`
}

// TableName sets the table name for Symbol
func (Symbol) TableName() string {
	return "symbols"
}

// Ticker represents a market ticker in the database
type Ticker struct {
	ID            uint      `gorm:"primaryKey;autoIncrement"`
	Symbol        string    `gorm:"index;type:varchar(20);not null"`
	Exchange      string    `gorm:"index;type:varchar(20);not null"`
	Price         float64   `gorm:"type:decimal(18,8);not null"`
	PriceChange   float64   `gorm:"type:decimal(18,8);not null"`
	PercentChange float64   `gorm:"type:decimal(18,8);not null"`
	High24h       float64   `gorm:"type:decimal(18,8);not null"`
	Low24h        float64   `gorm:"type:decimal(18,8);not null"`
	Volume        float64   `gorm:"type:decimal(18,8);not null"`
	LastUpdated   time.Time `gorm:"index;not null"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`
}

// TableName sets the table name for Ticker
func (Ticker) TableName() string {
	return "tickers"
}

// OrderBook represents an order book in the database
type OrderBook struct {
	ID          uint      `gorm:"primaryKey;autoIncrement"`
	Symbol      string    `gorm:"index;type:varchar(20);not null"`
	Exchange    string    `gorm:"index;type:varchar(20);not null"`
	BidsJSON    []byte    `gorm:"type:json;not null"`
	AsksJSON    []byte    `gorm:"type:json;not null"`
	LastUpdated time.Time `gorm:"index;not null"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

// TableName sets the table name for OrderBook
func (OrderBook) TableName() string {
	return "order_books"
}

// Candle represents a candlestick in the database
type Candle struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	Symbol    string    `gorm:"index;type:varchar(20);not null"`
	Exchange  string    `gorm:"index;type:varchar(20);not null"`
	Interval  string    `gorm:"index;type:varchar(10);not null"`
	OpenTime  time.Time `gorm:"index;not null"`
	CloseTime time.Time `gorm:"not null"`
	Open      float64   `gorm:"type:decimal(18,8);not null"`
	High      float64   `gorm:"type:decimal(18,8);not null"`
	Low       float64   `gorm:"type:decimal(18,8);not null"`
	Close     float64   `gorm:"type:decimal(18,8);not null"`
	Volume    float64   `gorm:"type:decimal(18,8);not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// TableName sets the table name for Candle
func (Candle) TableName() string {
	return "candles"
}

// Position represents a trading position in the database
type Position struct {
	ID         string    `gorm:"primaryKey"`
	UserID     string    `gorm:"not null;index"`
	Symbol     string    `gorm:"not null;index"`
	Side       string    `gorm:"not null"` // "LONG" or "SHORT"
	Quantity   float64   `gorm:"type:decimal(18,8);not null"`
	EntryPrice float64   `gorm:"type:decimal(18,8);not null"`
	Status     string    `gorm:"not null"` // "OPEN", "CLOSED"
	OpenedAt   time.Time `gorm:"not null"`
	ClosedAt   *time.Time
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`
}

// TableName sets the table name for Position
func (Position) TableName() string {
	return "positions"
}
