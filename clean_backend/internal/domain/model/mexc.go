package model

import "time"

// MEXCExchangeInfo extends the base ExchangeInfo with MEXC-specific fields
type MEXCExchangeInfo struct {
	ExchangeInfo
	Timezone   string    `json:"timezone"`
	ServerTime time.Time `json:"serverTime"`
}

// MEXCOrder extends the base Order with MEXC-specific fields
type MEXCOrder struct {
	Symbol           string      `json:"symbol"`
	OrderID          string      `json:"orderId"`
	ClientOrderID    string      `json:"clientOrderId"`
	Price            float64     `json:"price"`
	OriginalQuantity float64     `json:"origQty"`
	ExecutedQuantity float64     `json:"executedQty"`
	Status           OrderStatus `json:"status"`
	TimeInForce      TimeInForce `json:"timeInForce"`
	Type             OrderType   `json:"type"`
	Side             OrderSide   `json:"side"`
	Time             time.Time   `json:"time"`
	UpdateTime       time.Time   `json:"updateTime"`
}

// ToOrder converts a MEXCOrder to the canonical Order model
func (mo *MEXCOrder) ToOrder() *Order {
	return &Order{
		OrderID:       mo.OrderID,
		ClientOrderID: mo.ClientOrderID,
		Symbol:        mo.Symbol,
		Side:          mo.Side,
		Type:          mo.Type,
		Status:        mo.Status,
		Price:         mo.Price,
		Quantity:      mo.OriginalQuantity,
		ExecutedQty:   mo.ExecutedQuantity,
		TimeInForce:   mo.TimeInForce,
		CreatedAt:     mo.Time,
		UpdatedAt:     mo.UpdateTime,
		Exchange:      "MEXC",
	}
}
