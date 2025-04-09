package models

import "time"

// Account represents a user's trading account
type Account struct {
	ID        string    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	UserID    string    `gorm:"index;not null" json:"user_id"`
	Name      string    `gorm:"size:100;not null" json:"name"`
	Type      string    `gorm:"type:varchar(20);not null" json:"type"` // spot, margin, futures
	Status    string    `gorm:"type:varchar(20);not null;default:'active'" json:"status"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	Wallets []Wallet `gorm:"foreignKey:AccountID" json:"wallets,omitempty"`
}

// Wallet represents a user's cryptocurrency wallet
type Wallet struct {
	ID        string                   `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	AccountID string                   `gorm:"index;not null" json:"account_id"`
	Type      string                   `gorm:"type:varchar(20);not null" json:"type"` // spot, margin, futures
	Balances  map[string]*AssetBalance `gorm:"-" json:"balances"`
	CreatedAt time.Time               `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time               `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	Transactions []Transaction `gorm:"foreignKey:WalletID" json:"transactions,omitempty"`
}

// AssetBalance represents the balance of a specific cryptocurrency asset
type AssetBalance struct {
	ID        string    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	WalletID  string    `gorm:"index;not null" json:"wallet_id"`
	Asset     string    `gorm:"type:varchar(20);not null;index" json:"asset"`
	Free      float64   `gorm:"not null;default:0" json:"free"`
	Locked    float64   `gorm:"not null;default:0" json:"locked"`
	Total     float64   `gorm:"not null;default:0" json:"total"`
	Price     float64   `gorm:"not null;default:0" json:"price"` // Current price in USDT
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
