package entity

import (
	"time"
)

type RiskAssessmentEntity struct {
	ID        string    `gorm:"primaryKey"`
	AccountID string    `gorm:"not null;index"`
	Type      string    `gorm:"not null"`
	Level     string    `gorm:"not null"`
	Message   string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

func (RiskAssessmentEntity) TableName() string { return "risk_assessments" }
