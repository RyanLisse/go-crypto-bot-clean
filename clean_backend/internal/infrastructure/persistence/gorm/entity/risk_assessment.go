package entity

import (
	"time"
)

// RiskAssessmentEntity represents the GORM entity for risk assessment
type RiskAssessmentEntity struct {
	ID             string     `gorm:"column:id;primaryKey"`
	UserID         string     `gorm:"column:user_id;index"`
	Type           string     `gorm:"column:type;index"`   // RiskType as string
	Level          string     `gorm:"column:level;index"`  // RiskLevel as string
	Status         string     `gorm:"column:status;index"` // RiskStatus as string
	Symbol         string     `gorm:"column:symbol;index"`
	PositionID     string     `gorm:"column:position_id;index"`
	OrderID        string     `gorm:"column:order_id;index"`
	Score          float64    `gorm:"column:score"`
	Message        string     `gorm:"column:message;type:text"`
	Recommendation string     `gorm:"column:recommendation;type:text"`
	MetadataJSON   string     `gorm:"column:metadata_json;type:text"` // JSON string of metadata
	CreatedAt      time.Time  `gorm:"column:created_at;index"`
	UpdatedAt      time.Time  `gorm:"column:updated_at"`
	ResolvedAt     *time.Time `gorm:"column:resolved_at"`
}

// TableName overrides the table name
func (RiskAssessmentEntity) TableName() string {
	return "risk_assessments"
}
