package repo

import (
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"gorm.io/gorm"
)

// NewRiskAssessmentRepository creates a new risk assessment repository
func NewRiskAssessmentRepository(db *gorm.DB) port.RiskAssessmentRepository {
	return NewGormRiskAssessmentRepository(db)
}

// NewRiskProfileRepository creates a new risk profile repository
func NewRiskProfileRepository(db *gorm.DB) port.RiskProfileRepository {
	return NewGormRiskProfileRepository(db)
}
