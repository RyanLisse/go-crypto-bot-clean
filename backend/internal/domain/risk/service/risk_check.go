package service

// RiskCheck represents the result of a risk check
type RiskCheck struct {
	Allowed   bool    `json:"allowed"`
	Threshold float64 `json:"threshold"`
}
