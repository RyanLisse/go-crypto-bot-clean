package api

import (
	"database/sql"
	"log"

	"github.com/ryanlisse/go-crypto-bot/internal/domain/ai/factory"
	"github.com/ryanlisse/go-crypto-bot/internal/domain/ai/service"
)

// AIServiceDependency adds AI service to the Dependencies struct
func (d *Dependencies) InitializeAIService(
	db *sql.DB,
) {
	// Create AI service
	aiSvc, err := factory.CreateAIService(db)
	if err != nil {
		log.Printf("Failed to create AI service: %v", err)
		return
	}

	// Store AI service in dependencies
	d.AIService = aiSvc
}

// GetAIService returns the AI service
func (d *Dependencies) GetAIService() service.AIService {
	return d.AIService
}
