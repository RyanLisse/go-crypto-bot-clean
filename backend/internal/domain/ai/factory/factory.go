package factory

import (
	"context"
	"fmt"
	"os"

	"go-crypto-bot-clean/backend/internal/domain/ai/repository"
	"go-crypto-bot-clean/backend/internal/domain/ai/service"
	"go-crypto-bot-clean/backend/internal/domain/ai/service/gemini"
	"go-crypto-bot-clean/backend/internal/domain/portfolio"
	"go-crypto-bot-clean/backend/internal/domain/risk"
	"go-crypto-bot-clean/backend/internal/domain/trade"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
	"gorm.io/gorm"
)

// CreateAIService creates an AI service based on configuration
func CreateAIService(
	db *gorm.DB,
	portfolioSvc portfolio.Service,
	tradeSvc trade.Service,
	riskSvc risk.Service,
) (service.AIService, error) {
	// Get AI provider from environment variable
	aiProvider := os.Getenv("AI_PROVIDER")
	if aiProvider == "" {
		aiProvider = "gemini" // Default to Gemini
	}

	switch aiProvider {
	case "gemini":
		return createGeminiAIService(db, portfolioSvc, tradeSvc, riskSvc)
	// Add other providers as needed
	default:
		return nil, fmt.Errorf("unsupported AI provider: %s", aiProvider)
	}
}

// createGeminiAIService creates a Gemini AI service
func createGeminiAIService(
	db *gorm.DB,
	portfolioSvc portfolio.Service,
	tradeSvc trade.Service,
	riskSvc risk.Service,
) (service.AIService, error) {
	// Get API key from environment variable
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY environment variable not set")
	}

	// Create Gemini client
	client, err := genai.NewClient(context.Background(), option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	// Create conversation memory repository
	memoryRepo, err := repository.NewGormConversationMemoryRepository(db)
	if err != nil {
		return nil, fmt.Errorf("failed to create conversation memory repository: %w", err)
	}

	// Create Gemini AI service
	return gemini.NewGeminiAIService(client, memoryRepo, portfolioSvc, tradeSvc, riskSvc)
}
