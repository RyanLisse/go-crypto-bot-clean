package models

import (
	"time"
)

// DecisionType represents the type of trading decision
type DecisionType string

const (
	DecisionTypeBuy   DecisionType = "BUY"
	DecisionTypeSell  DecisionType = "SELL"
	DecisionTypeHold  DecisionType = "HOLD"
	DecisionTypeClose DecisionType = "CLOSE"
)

// DecisionStatus represents the status of a trading decision
type DecisionStatus string

const (
	DecisionStatusPending   DecisionStatus = "PENDING"
	DecisionStatusExecuted  DecisionStatus = "EXECUTED"
	DecisionStatusRejected  DecisionStatus = "REJECTED"
	DecisionStatusCancelled DecisionStatus = "CANCELLED"
)

// DecisionReason categorizes the reason for a trading decision
type DecisionReason string

const (
	// Buy reasons
	ReasonStrategySignal     DecisionReason = "STRATEGY_SIGNAL"
	ReasonManualEntry        DecisionReason = "MANUAL_ENTRY"
	ReasonRebalancing        DecisionReason = "REBALANCING"
	ReasonNewCoinListing     DecisionReason = "NEW_COIN_LISTING"
	ReasonTechnicalPattern   DecisionReason = "TECHNICAL_PATTERN"
	ReasonVolumeSpike        DecisionReason = "VOLUME_SPIKE"
	ReasonPriceBreakout      DecisionReason = "PRICE_BREAKOUT"
	
	// Sell reasons
	ReasonTakeProfit         DecisionReason = "TAKE_PROFIT"
	ReasonStopLoss           DecisionReason = "STOP_LOSS"
	ReasonTrailingStop       DecisionReason = "TRAILING_STOP"
	ReasonManualExit         DecisionReason = "MANUAL_EXIT"
	ReasonTimeBasedExit      DecisionReason = "TIME_BASED_EXIT"
	ReasonRiskLimitReached   DecisionReason = "RISK_LIMIT_REACHED"
	ReasonStrategyExit       DecisionReason = "STRATEGY_EXIT"
	
	// Rejection reasons
	ReasonInsufficientFunds  DecisionReason = "INSUFFICIENT_FUNDS"
	ReasonMaxPositionsReached DecisionReason = "MAX_POSITIONS_REACHED"
	ReasonRiskControlRejection DecisionReason = "RISK_CONTROL_REJECTION"
	ReasonLowConfidence      DecisionReason = "LOW_CONFIDENCE"
	ReasonMarketClosed       DecisionReason = "MARKET_CLOSED"
	ReasonInvalidParameters  DecisionReason = "INVALID_PARAMETERS"
)

// TradeDecision represents a decision to enter or exit a trade
type TradeDecision struct {
	ID              string          `json:"id" db:"id"`
	Symbol          string          `json:"symbol" db:"symbol"`
	Type            DecisionType    `json:"type" db:"type"`
	Status          DecisionStatus  `json:"status" db:"status"`
	Reason          DecisionReason  `json:"reason" db:"reason"`
	DetailedReason  string          `json:"detailed_reason" db:"detailed_reason"`
	Price           float64         `json:"price" db:"price"`
	Quantity        float64         `json:"quantity" db:"quantity"`
	TotalValue      float64         `json:"total_value" db:"total_value"`
	Confidence      float64         `json:"confidence" db:"confidence"`
	Strategy        string          `json:"strategy" db:"strategy"`
	StrategyParams  string          `json:"strategy_params" db:"strategy_params"`
	CreatedAt       time.Time       `json:"created_at" db:"created_at"`
	ExecutedAt      *time.Time      `json:"executed_at,omitempty" db:"executed_at"`
	PositionID      *string         `json:"position_id,omitempty" db:"position_id"`
	OrderID         *string         `json:"order_id,omitempty" db:"order_id"`
	StopLoss        *float64        `json:"stop_loss,omitempty" db:"stop_loss"`
	TakeProfit      *float64        `json:"take_profit,omitempty" db:"take_profit"`
	TrailingStop    *float64        `json:"trailing_stop,omitempty" db:"trailing_stop"`
	RiskRewardRatio *float64        `json:"risk_reward_ratio,omitempty" db:"risk_reward_ratio"`
	ExpectedProfit  *float64        `json:"expected_profit,omitempty" db:"expected_profit"`
	MaxRisk         *float64        `json:"max_risk,omitempty" db:"max_risk"`
	Tags            []string        `json:"-" db:"-"`
	TagsString      string          `json:"tags" db:"tags"`
	Metadata        map[string]interface{} `json:"-" db:"-"`
	MetadataJSON    string          `json:"metadata" db:"metadata_json"`
}

// TradeDecisionSummary provides a summary of trading decisions for a time period
type TradeDecisionSummary struct {
	Period          string    `json:"period"`
	StartTime       time.Time `json:"start_time"`
	EndTime         time.Time `json:"end_time"`
	TotalDecisions  int       `json:"total_decisions"`
	BuyDecisions    int       `json:"buy_decisions"`
	SellDecisions   int       `json:"sell_decisions"`
	ExecutedCount   int       `json:"executed_count"`
	RejectedCount   int       `json:"rejected_count"`
	SuccessRate     float64   `json:"success_rate"`
	AverageProfit   float64   `json:"average_profit"`
	TotalProfit     float64   `json:"total_profit"`
	ProfitableCount int       `json:"profitable_count"`
	LossCount       int       `json:"loss_count"`
	WinRate         float64   `json:"win_rate"`
	TopSymbols      []string  `json:"top_symbols"`
	TopStrategies   []string  `json:"top_strategies"`
	TopReasons      []string  `json:"top_reasons"`
}
