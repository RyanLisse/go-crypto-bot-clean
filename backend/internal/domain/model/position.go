package model

import (
	"time"
)

// PositionSide represents the side of a position (long or short)
type PositionSide string

// PositionStatus represents the status of a position
type PositionStatus string

// PositionType represents the type of a position
type PositionType string

// Position side constants
const (
	PositionSideLong  PositionSide = "LONG"
	PositionSideShort PositionSide = "SHORT"
)

// Position status constants
const (
	PositionStatusOpen   PositionStatus = "OPEN"
	PositionStatusClosed PositionStatus = "CLOSED"
)

// Position type constants
const (
	PositionTypeManual    PositionType = "MANUAL"
	PositionTypeAutomatic PositionType = "AUTOMATIC"
	PositionTypeNewCoin   PositionType = "NEWCOIN"
)

// Position represents a trading position
type Position struct {
	ID              string         `json:"id"`
	Symbol          string         `json:"symbol"`
	Side            PositionSide   `json:"side"`
	Status          PositionStatus `json:"status"`
	Type            PositionType   `json:"type"`
	EntryPrice      float64        `json:"entryPrice"`
	Quantity        float64        `json:"quantity"`
	CurrentPrice    float64        `json:"currentPrice"`
	PnL             float64        `json:"pnl"`
	PnLPercent      float64        `json:"pnlPercent"`
	StopLoss        *float64       `json:"stopLoss,omitempty"`
	TakeProfit      *float64       `json:"takeProfit,omitempty"`
	StrategyID      *string        `json:"strategyId,omitempty"`
	OpenOrderIDs    []string       `json:"openOrderIds,omitempty"`
	EntryOrderIDs   []string       `json:"entryOrderIds"`
	ExitOrderIDs    []string       `json:"exitOrderIds,omitempty"`
	Notes           string         `json:"notes,omitempty"`
	OpenedAt        time.Time      `json:"openedAt"`
	ClosedAt        *time.Time     `json:"closedAt,omitempty"`
	LastUpdatedAt   time.Time      `json:"lastUpdatedAt"`
	MaxDrawdown     float64        `json:"maxDrawdown"`
	MaxProfit       float64        `json:"maxProfit"`
	RiskRewardRatio float64        `json:"riskRewardRatio,omitempty"`
}

// PositionCreateRequest represents data needed to create a position
type PositionCreateRequest struct {
	Symbol     string       `json:"symbol" binding:"required"`
	Side       PositionSide `json:"side" binding:"required,oneof=LONG SHORT"`
	Type       PositionType `json:"type" binding:"required"`
	EntryPrice float64      `json:"entryPrice" binding:"required,gt=0"`
	Quantity   float64      `json:"quantity" binding:"required,gt=0"`
	StopLoss   *float64     `json:"stopLoss"`
	TakeProfit *float64     `json:"takeProfit"`
	StrategyID *string      `json:"strategyId"`
	OrderIDs   []string     `json:"orderIds" binding:"required,min=1"`
	Notes      string       `json:"notes"`
}

// PositionUpdateRequest represents data for updating a position
type PositionUpdateRequest struct {
	CurrentPrice *float64  `json:"currentPrice"`
	StopLoss     *float64  `json:"stopLoss"`
	TakeProfit   *float64  `json:"takeProfit"`
	Notes        *string   `json:"notes"`
	Status       *string   `json:"status"`
	ClosedAt     *string   `json:"closedAt"`
	ExitOrderIDs *[]string `json:"exitOrderIds"`
}

// UpdateCurrentPrice updates the current price and recalculates PnL
func (p *Position) UpdateCurrentPrice(currentPrice float64) {
	p.CurrentPrice = currentPrice

	// Calculate PnL
	if p.Side == PositionSideLong {
		p.PnL = (currentPrice - p.EntryPrice) * p.Quantity
		p.PnLPercent = (currentPrice - p.EntryPrice) / p.EntryPrice * 100
	} else {
		p.PnL = (p.EntryPrice - currentPrice) * p.Quantity
		p.PnLPercent = (p.EntryPrice - currentPrice) / p.EntryPrice * 100
	}

	p.LastUpdatedAt = time.Now()

	// Update max profit/drawdown
	if p.PnL > p.MaxProfit {
		p.MaxProfit = p.PnL
	}

	if p.PnL < 0 && p.PnL < p.MaxDrawdown {
		p.MaxDrawdown = p.PnL
	}
}

// Close closes the position
func (p *Position) Close(exitPrice float64, exitOrderIDs []string) {
	p.Status = PositionStatusClosed
	p.UpdateCurrentPrice(exitPrice)
	now := time.Now()
	p.ClosedAt = &now
	p.ExitOrderIDs = exitOrderIDs
}

// CalculateRiskRewardRatio calculates the risk/reward ratio if stop-loss and take-profit are set
func (p *Position) CalculateRiskRewardRatio() {
	if p.StopLoss == nil || p.TakeProfit == nil {
		p.RiskRewardRatio = 0
		return
	}

	var reward, risk float64

	if p.Side == PositionSideLong {
		reward = (*p.TakeProfit - p.EntryPrice) * p.Quantity
		risk = (p.EntryPrice - *p.StopLoss) * p.Quantity
	} else {
		reward = (p.EntryPrice - *p.TakeProfit) * p.Quantity
		risk = (*p.StopLoss - p.EntryPrice) * p.Quantity
	}

	if risk <= 0 {
		p.RiskRewardRatio = 0
		return
	}

	p.RiskRewardRatio = reward / risk
}
