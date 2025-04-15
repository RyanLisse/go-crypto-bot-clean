package entity

import (
	"time"

	"gorm.io/gorm"
)

// AutoBuyRuleEntity represents an auto-buy rule in the database
type AutoBuyRuleEntity struct {
	ID                  string  `gorm:"primaryKey"`
	UserID              string  `gorm:"not null;index"`
	Name                string  `gorm:"not null"`
	Symbol              string  `gorm:"not null;index"`
	IsEnabled           bool    `gorm:"not null;default:true"`
	TriggerType         string  `gorm:"not null"`
	TriggerValue        float64 `gorm:"not null"`
	QuoteAsset          string  `gorm:"not null"`
	BuyAmountQuote      float64 `gorm:"not null"`
	MaxBuyPrice         *float64
	MinBaseAssetVolume  *float64
	MinQuoteAssetVolume *float64
	AllowPreTrading     bool   `gorm:"not null;default:false"`
	CooldownMinutes     int    `gorm:"not null;default:0"`
	OrderType           string `gorm:"not null;default:'MARKET'"`
	EnableRiskCheck     bool   `gorm:"not null;default:true"`
	ExecutionCount      int    `gorm:"not null;default:0"`
	LastTriggered       *time.Time
	LastPrice           float64    `gorm:"not null;default:0"`
	CreatedAt           time.Time  `gorm:"not null;autoCreateTime"`
	UpdatedAt           time.Time  `gorm:"not null;autoUpdateTime"`
	DeletedAt           *time.Time `gorm:"index"`

	// Deprecated fields (for migration purposes)
	UsePercentage    *bool
	PercentageAmount *float64
	FixedAmount      *float64
	MinOrderAmount   *float64
}

// TableName specifies the table name for the AutoBuyRuleEntity
func (AutoBuyRuleEntity) TableName() string {
	return "auto_buy_rules"
}

// BeforeCreate handles pre-creation hooks
func (e *AutoBuyRuleEntity) BeforeCreate(tx *gorm.DB) error {
	// Set default CreatedAt and UpdatedAt if not already set
	if e.CreatedAt.IsZero() {
		e.CreatedAt = time.Now()
	}
	if e.UpdatedAt.IsZero() {
		e.UpdatedAt = time.Now()
	}
	return nil
}

// AutoBuyExecutionEntity represents an auto-buy execution record in the database
type AutoBuyExecutionEntity struct {
	ID        string    `gorm:"primaryKey"`
	RuleID    string    `gorm:"not null;index"`
	UserID    string    `gorm:"not null;index"`
	Symbol    string    `gorm:"not null;index"`
	OrderID   string    `gorm:"not null"`
	Price     float64   `gorm:"not null"`
	Quantity  float64   `gorm:"not null"`
	Amount    float64   `gorm:"not null"`
	Timestamp time.Time `gorm:"not null;index"`
	CreatedAt time.Time `gorm:"not null;autoCreateTime"`
}

// TableName specifies the table name for the AutoBuyExecutionEntity
func (AutoBuyExecutionEntity) TableName() string {
	return "auto_buy_executions"
}

// BeforeCreate handles pre-creation hooks
func (e *AutoBuyExecutionEntity) BeforeCreate(tx *gorm.DB) error {
	if e.CreatedAt.IsZero() {
		e.CreatedAt = time.Now()
	}
	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now()
	}
	return nil
}
