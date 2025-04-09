package models

import (
	"time"

	"gorm.io/gorm"
)

// Backtest represents a backtest run
type Backtest struct {
	ID             string         `gorm:"primaryKey;type:varchar(36)" json:"id"`
	UserID         string         `gorm:"index;not null;type:varchar(36)" json:"userId"`
	StrategyID     string         `gorm:"index;not null;type:varchar(36)" json:"strategyId"`
	Name           string         `gorm:"not null;type:varchar(100)" json:"name"`
	Description    string         `gorm:"type:text" json:"description"`
	StartDate      time.Time      `gorm:"not null" json:"startDate"`
	EndDate        time.Time      `gorm:"not null" json:"endDate"`
	InitialBalance float64        `gorm:"not null" json:"initialBalance"`
	FinalBalance   float64        `gorm:"not null" json:"finalBalance"`
	TotalTrades    int            `gorm:"not null" json:"totalTrades"`
	WinningTrades  int            `gorm:"not null" json:"winningTrades"`
	LosingTrades   int            `gorm:"not null" json:"losingTrades"`
	WinRate        float64        `gorm:"not null" json:"winRate"`
	ProfitFactor   float64        `gorm:"not null" json:"profitFactor"`
	SharpeRatio    float64        `gorm:"not null" json:"sharpeRatio"`
	MaxDrawdown    float64        `gorm:"not null" json:"maxDrawdown"`
	Parameters     Parameters     `gorm:"type:json" json:"parameters"`
	Status         string         `gorm:"not null;type:varchar(20)" json:"status"`
	CreatedAt      time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for the Backtest model
func (Backtest) TableName() string {
	return "backtests"
}

// BacktestTrade represents a trade executed during a backtest
type BacktestTrade struct {
	ID           uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	BacktestID   string         `gorm:"index;not null;type:varchar(36)" json:"backtestId"`
	Symbol       string         `gorm:"not null;type:varchar(20)" json:"symbol"`
	EntryTime    time.Time      `gorm:"not null" json:"entryTime"`
	EntryPrice   float64        `gorm:"not null" json:"entryPrice"`
	ExitTime     *time.Time     `json:"exitTime,omitempty"`
	ExitPrice    *float64       `json:"exitPrice,omitempty"`
	Quantity     float64        `gorm:"not null" json:"quantity"`
	Direction    string         `gorm:"not null;type:varchar(10)" json:"direction"`
	ProfitLoss   *float64       `json:"profitLoss,omitempty"`
	ProfitLossPct *float64      `json:"profitLossPct,omitempty"`
	ExitReason   *string        `gorm:"type:varchar(50)" json:"exitReason,omitempty"`
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for the BacktestTrade model
func (BacktestTrade) TableName() string {
	return "backtest_trades"
}

// BacktestEquity represents the equity curve of a backtest
type BacktestEquity struct {
	ID         uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	BacktestID string         `gorm:"index;not null;type:varchar(36)" json:"backtestId"`
	Timestamp  time.Time      `gorm:"not null" json:"timestamp"`
	Equity     float64        `gorm:"not null" json:"equity"`
	Balance    float64        `gorm:"not null" json:"balance"`
	Drawdown   float64        `gorm:"not null" json:"drawdown"`
	CreatedAt  time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for the BacktestEquity model
func (BacktestEquity) TableName() string {
	return "backtest_equity"
}
