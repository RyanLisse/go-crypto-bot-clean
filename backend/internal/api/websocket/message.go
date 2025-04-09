package websocket

import "time"

// MessageType defines the type of message sent over WebSocket
type MessageType string

const (
	// Data message types
	MarketDataType        MessageType = "market_data"
	TradeNotificationType MessageType = "trade_notification"
	NewCoinAlertType      MessageType = "new_coin_alert"
	PortfolioUpdateType   MessageType = "portfolio_update"
	TradeUpdateType       MessageType = "trade_update"
	AccountUpdateType     MessageType = "account_update"

	// System message types
	ErrorType               MessageType = "error"
	SubscriptionSuccessType MessageType = "subscription_success"
	PingType                MessageType = "ping"
	PongType                MessageType = "pong"
	AuthSuccessType         MessageType = "auth_success"
	AuthFailureType         MessageType = "auth_failure"
)

// WSMessage is the envelope for all WebSocket messages
type WSMessage struct {
	Type      MessageType `json:"type"`
	Timestamp int64       `json:"timestamp"`
	Payload   any         `json:"payload"`
}

// MarketDataPayload represents real-time market data
type MarketDataPayload struct {
	Symbol    string  `json:"symbol"`
	Price     float64 `json:"price"`
	Volume    float64 `json:"volume"`
	Timestamp int64   `json:"timestamp"`
}

// TradeNotificationPayload represents executed trade info
type TradeNotificationPayload struct {
	Symbol      string    `json:"symbol"`
	Price       float64   `json:"price"`
	Quantity    float64   `json:"quantity"`
	PurchasedAt time.Time `json:"purchased_at"`
	TradeType   string    `json:"trade_type"` // buy/sell
}

// NewCoinAlertPayload represents a new coin listing alert
type NewCoinAlertPayload struct {
	Symbol      string `json:"symbol"`
	ListedAt    int64  `json:"listed_at"`
	Description string `json:"description"`
}

// ErrorPayload represents an error message
type ErrorPayload struct {
	Message string `json:"message"`
}

// SubscriptionSuccessPayload confirms a subscription
type SubscriptionSuccessPayload struct {
	Message string `json:"message"`
	Channel string `json:"channel"`
}

// PortfolioUpdatePayload represents portfolio update information
type PortfolioUpdatePayload struct {
	TotalValue float64        `json:"total_value"`
	Assets     []AssetPayload `json:"assets"`
	Timestamp  int64          `json:"timestamp"`
}

// AssetPayload represents a single asset in the portfolio
type AssetPayload struct {
	Symbol     string  `json:"symbol"`
	Amount     float64 `json:"amount"`
	ValueUSD   float64 `json:"value_usd"`
	Allocation float64 `json:"allocation_percentage"`
}

// TradeUpdatePayload represents a trade update
type TradeUpdatePayload struct {
	ID        string  `json:"id"`
	Symbol    string  `json:"symbol"`
	Side      string  `json:"side"`
	Price     float64 `json:"price"`
	Quantity  float64 `json:"quantity"`
	Total     float64 `json:"total"`
	Status    string  `json:"status"`
	Timestamp int64   `json:"timestamp"`
}

// PingPayload represents a ping message
type PingPayload struct {
	Timestamp int64 `json:"timestamp"`
}

// PongPayload represents a pong response
type PongPayload struct {
	Timestamp int64 `json:"timestamp"`
}

// AuthPayload represents authentication information
type AuthPayload struct {
	Token string `json:"token"`
}

// AccountUpdatePayload represents account update information
type AccountUpdatePayload struct {
	Balances  map[string]AssetBalancePayload `json:"balances"`
	UpdatedAt int64                          `json:"updatedAt"`
}

// AssetBalancePayload represents a single asset balance
type AssetBalancePayload struct {
	Asset  string  `json:"asset"`
	Free   float64 `json:"free"`
	Locked float64 `json:"locked"`
	Total  float64 `json:"total"`
}
