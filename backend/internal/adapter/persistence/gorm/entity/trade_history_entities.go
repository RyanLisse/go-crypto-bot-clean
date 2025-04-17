package entity

import (
	"time"
)

// TradeRecordEntity represents a trade record in the database
type TradeRecordEntity struct {
	ID            string    `gorm:"primaryKey"`
	UserID        string    `gorm:"index"`
	Symbol        string    `gorm:"index"`
	Side          string    `gorm:"index"`
	Type          string
	Quantity      float64
	Price         float64
	Amount        float64
	Fee           float64
	FeeCurrency   string
	OrderID       string    `gorm:"index"`
	TradeID       string    `gorm:"index"`
	ExecutionTime time.Time `gorm:"index"`
	Strategy      string    `gorm:"index"`
	Notes         string
	Tags          string // JSON array
	Metadata      string // JSON object
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// TableName returns the table name for the trade record entity
func (TradeRecordEntity) TableName() string {
	return "trade_records"
}

// DetectionLogEntity represents a detection log in the database
type DetectionLogEntity struct {
	ID          string    `gorm:"primaryKey"`
	Type        string    `gorm:"index"`
	Symbol      string    `gorm:"index"`
	Value       float64
	Threshold   float64
	Description string
	Metadata    string // JSON object
	DetectedAt  time.Time `gorm:"index"`
	ProcessedAt *time.Time
	Processed   bool      `gorm:"index"`
	Result      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// TableName returns the table name for the detection log entity
func (DetectionLogEntity) TableName() string {
	return "detection_logs"
}
