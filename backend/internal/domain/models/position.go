package models

import "time"

// Position structure represents a trading position
type Position struct {
	ID               string            `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	Symbol           string            `gorm:"index;not null;size:20" json:"symbol"`
	Side             OrderSide         `gorm:"type:varchar(4);not null" json:"side"` // BUY or SELL
	Quantity         float64           `gorm:"not null" json:"quantity"`
	Amount           float64           `gorm:"-" json:"amount"` // Alias for Quantity, not stored
	EntryPrice       float64           `gorm:"not null" json:"entry_price"`
	CurrentPrice     float64           `gorm:"not null" json:"current_price"`
	OpenTime         time.Time         `gorm:"index;not null" json:"open_time"`
	OpenedAt         time.Time         `gorm:"-" json:"opened_at"` // Alias for OpenTime
	CloseTime        *time.Time        `gorm:"index" json:"close_time,omitempty"`
	StopLoss         float64           `gorm:"not null" json:"stop_loss"`
	TakeProfit       float64           `gorm:"not null" json:"take_profit"`
	TrailingStop     *float64          `gorm:"default:null" json:"trailing_stop,omitempty"`
	CreatedAt        time.Time         `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time         `gorm:"autoUpdateTime" json:"updated_at"`
	PnL              float64           `gorm:"not null;default:0" json:"pnl"`
	PnLPercentage    float64           `gorm:"not null;default:0" json:"pnl_percentage"`
	Status           PositionStatus    `gorm:"index;type:varchar(10);not null;default:'open'" json:"status"` // open, closed
	Orders           []Order           `gorm:"foreignKey:PositionID" json:"orders"`                          // Entry and scaling orders
	EntryReason      string            `gorm:"type:varchar(100)" json:"entry_reason,omitempty"`
	ExitReason       string            `gorm:"type:varchar(100)" json:"exit_reason,omitempty"`
	Strategy         string            `gorm:"type:varchar(50)" json:"strategy,omitempty"`
	RiskRewardRatio  float64           `gorm:"default:0" json:"risk_reward_ratio,omitempty"`
	ExpectedProfit   float64           `gorm:"default:0" json:"expected_profit,omitempty"`
	MaxRisk          float64           `gorm:"default:0" json:"max_risk,omitempty"`
	TakeProfitLevels []TakeProfitLevel `gorm:"foreignKey:PositionID" json:"take_profit_levels,omitempty"`
	Tags             []string          `gorm:"-" json:"tags,omitempty"`
	Notes            string            `gorm:"-" json:"notes,omitempty"`
}

// CalculateValue returns the current value of the position
func (p *Position) CalculateValue() float64 {
	return p.CurrentPrice * p.Quantity
}

// CalculateUnrealizedPnL calculates the unrealized profit/loss
func (p *Position) CalculateUnrealizedPnL() float64 {
	return (p.CurrentPrice - p.EntryPrice) * p.Quantity
}

// IsOpen checks if the position is open
func (p *Position) IsOpen() bool {
	return p.Status == PositionStatusOpen
}

// Close closes the position with the given price
func (p *Position) Close(closePrice float64, closeTime time.Time) {
	p.Status = PositionStatusClosed
	p.CurrentPrice = closePrice
	p.CloseTime = &closeTime
	p.PnL = (closePrice - p.EntryPrice) * p.Quantity
	p.PnLPercentage = ((closePrice - p.EntryPrice) / p.EntryPrice) * 100
}
