package model

import (
	"time"

	"github.com/google/uuid"
)

// RiskConstraint represents a specific risk rule or constraint
type RiskConstraint struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	UserID      string    `json:"user_id" gorm:"index"`
	Type        RiskType  `json:"type" gorm:"index"`
	Parameter   string    `json:"parameter"`           // e.g., "max_position_size", "max_drawdown"
	Operator    string    `json:"operator"`            // e.g., "GT" (greater than), "LT" (less than)
	Value       float64   `json:"value"`               // The threshold value
	Action      string    `json:"action"`              // e.g., "BLOCK", "WARN"
	Description string    `json:"description"`         // Human-readable description
	Active      bool      `json:"active" gorm:"index"` // Whether this constraint is active
	Symbol      string    `json:"symbol,omitempty"`    // Optional: specific symbol this applies to
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// NewRiskConstraint creates a new risk constraint
func NewRiskConstraint(userID string, riskType RiskType, parameter string, operator string, value float64, action string) *RiskConstraint {
	return &RiskConstraint{
		ID:        uuid.New().String(),
		UserID:    userID,
		Type:      riskType,
		Parameter: parameter,
		Operator:  operator,
		Value:     value,
		Action:    action,
		Active:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// Deactivate deactivates a risk constraint
func (r *RiskConstraint) Deactivate() {
	r.Active = false
	r.UpdatedAt = time.Now()
}

// Activate activates a risk constraint
func (r *RiskConstraint) Activate() {
	r.Active = true
	r.UpdatedAt = time.Now()
}
