package api

import (
	"database/sql"

	"go-crypto-bot-clean/backend/internal/domain/ai/service"
)

// AIServiceDependency adds AI service to the Dependencies struct
func (d *Dependencies) InitializeAIService(
	db *sql.DB,
) {
	// Create mock AI service
	d.logger.Info("Using mock AI service")
	d.AIService = &MockAIService{}
}

// GetAIService returns the AI service
func (d *Dependencies) GetAIService() service.AIService {
	return d.AIService
}
