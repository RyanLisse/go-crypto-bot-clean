package models

// TakeProfitLevel represents a level at which to take profit
type TakeProfitLevel struct {
	Level       int     `json:"level"`
	Price       float64 `json:"price"`
	Percentage  float64 `json:"percentage"`
	Quantity    float64 `json:"quantity"`
	QuantityPct float64 `json:"quantity_pct"`
	Triggered   bool    `json:"triggered"`
	Executed    bool    `json:"executed"`
}
