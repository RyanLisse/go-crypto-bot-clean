package models

import (
	"time"
)

// Transaction represents a financial transaction
type Transaction struct {
	ID          string    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	Amount      float64   `gorm:"not null" json:"amount"`
	Balance     float64   `gorm:"not null" json:"balance"`
	Reason      string    `gorm:"not null;size:255" json:"reason"`
	Timestamp   time.Time `gorm:"index;not null" json:"timestamp"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	
	// Foreign key relationships
	PositionID  string    `gorm:"index" json:"position_id,omitempty"`
	OrderID     string    `gorm:"index" json:"order_id,omitempty"`
	WalletID    string    `gorm:"index;not null" json:"wallet_id"`
}

// BalanceSummary provides an overview of wallet activity
type BalanceSummary struct {
	ID               string    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	CurrentBalance   float64   `gorm:"not null" json:"current_balance"`
	Deposits         float64   `gorm:"not null" json:"deposits"`
	Withdrawals      float64   `gorm:"not null" json:"withdrawals"`
	NetChange        float64   `gorm:"not null" json:"net_change"`
	TransactionCount int       `gorm:"not null" json:"transaction_count"`
	Period           int       `gorm:"not null" json:"period_days"`
	GeneratedAt      time.Time `gorm:"not null" json:"generated_at"`
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	
	// Foreign key relationship
	WalletID         string    `gorm:"index;not null" json:"wallet_id"
}

// TransactionAnalysis provides analysis of transaction history
type TransactionAnalysis struct {
	ID          string    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	StartTime   time.Time `gorm:"not null" json:"start_time"`
	EndTime     time.Time `gorm:"not null" json:"end_time"`
	TotalCount  int       `gorm:"not null" json:"total_count"`
	BuyCount    int       `gorm:"not null" json:"buy_count"`
	SellCount   int       `gorm:"not null" json:"sell_count"`
	TotalVolume float64   `gorm:"not null" json:"total_volume"`
	BuyVolume   float64   `gorm:"not null" json:"buy_volume"`
	SellVolume  float64   `gorm:"not null" json:"sell_volume"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	
	// Foreign key relationship
	WalletID    string    `gorm:"index;not null" json:"wallet_id"
}
