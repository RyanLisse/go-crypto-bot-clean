package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
)

// Strategy represents a trading strategy
type Strategy struct {
	ID          string         `gorm:"primaryKey;type:varchar(36)" json:"id"`
	Name        string         `gorm:"not null;type:varchar(100)" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	Parameters  Parameters     `gorm:"type:json" json:"parameters"`
	IsEnabled   bool           `gorm:"not null;default:false" json:"isEnabled"`
	UserID      string         `gorm:"index;type:varchar(36)" json:"userId"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for the Strategy model
func (Strategy) TableName() string {
	return "strategies"
}

// Parameters represents strategy parameters stored as JSON
type Parameters map[string]interface{}

// Scan implements the sql.Scanner interface for Parameters
func (p *Parameters) Scan(value interface{}) error {
	if value == nil {
		*p = make(Parameters)
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to unmarshal Parameters value")
	}

	result := make(Parameters)
	err := json.Unmarshal(bytes, &result)
	if err != nil {
		return err
	}

	*p = result
	return nil
}

// Value implements the driver.Valuer interface for Parameters
func (p Parameters) Value() (driver.Value, error) {
	if p == nil {
		return nil, nil
	}

	return json.Marshal(p)
}

// StrategyPerformance represents the performance metrics of a strategy
type StrategyPerformance struct {
	ID           uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	StrategyID   string         `gorm:"index;not null;type:varchar(36)" json:"strategyId"`
	WinRate      float64        `gorm:"not null" json:"winRate"`
	ProfitFactor float64        `gorm:"not null" json:"profitFactor"`
	SharpeRatio  float64        `gorm:"not null" json:"sharpeRatio"`
	MaxDrawdown  float64        `gorm:"not null" json:"maxDrawdown"`
	TotalTrades  int            `gorm:"not null" json:"totalTrades"`
	PeriodStart  time.Time      `gorm:"not null" json:"periodStart"`
	PeriodEnd    time.Time      `gorm:"not null" json:"periodEnd"`
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for the StrategyPerformance model
func (StrategyPerformance) TableName() string {
	return "strategy_performance"
}

// StrategyParameter represents a parameter for a strategy
type StrategyParameter struct {
	ID          uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	StrategyID  string         `gorm:"index;not null;type:varchar(36)" json:"strategyId"`
	Name        string         `gorm:"not null;type:varchar(100)" json:"name"`
	Type        string         `gorm:"not null;type:varchar(50)" json:"type"`
	Description string         `gorm:"type:text" json:"description"`
	Default     string         `gorm:"type:text" json:"default"`
	Min         *string        `gorm:"type:text" json:"min,omitempty"`
	Max         *string        `gorm:"type:text" json:"max,omitempty"`
	Options     Options        `gorm:"type:json" json:"options,omitempty"`
	Required    bool           `gorm:"not null;default:false" json:"required"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for the StrategyParameter model
func (StrategyParameter) TableName() string {
	return "strategy_parameters"
}

// Options represents parameter options stored as JSON
type Options []string

// Scan implements the sql.Scanner interface for Options
func (o *Options) Scan(value interface{}) error {
	if value == nil {
		*o = make(Options, 0)
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to unmarshal Options value")
	}

	result := make(Options, 0)
	err := json.Unmarshal(bytes, &result)
	if err != nil {
		return err
	}

	*o = result
	return nil
}

// Value implements the driver.Valuer interface for Options
func (o Options) Value() (driver.Value, error) {
	if o == nil {
		return nil, nil
	}

	return json.Marshal(o)
}
