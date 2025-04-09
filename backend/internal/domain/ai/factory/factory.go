package factory

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/google/generative-ai-go/genai"
	"github.com/ryanlisse/go-crypto-bot/internal/domain/ai/service"
	"github.com/ryanlisse/go-crypto-bot/internal/domain/ai/service/gemini"
	"google.golang.org/api/option"
)

// CreateAIService creates an AI service based on configuration
func CreateAIService(
	db *sql.DB,
) (service.AIService, error) {
	// Get AI provider from environment variable
	aiProvider := os.Getenv("AI_PROVIDER")
	if aiProvider == "" {
		aiProvider = "gemini" // Default to Gemini
	}

	switch aiProvider {
	case "gemini":
		return createGeminiAIService(db)
	// Add other providers as needed
	default:
		return nil, fmt.Errorf("unsupported AI provider: %s", aiProvider)
	}
}

// createGeminiAIService creates a Gemini AI service
func createGeminiAIService(
	db *sql.DB,
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

	// Create Gemini AI service
	return gemini.NewGeminiAIService(client, db)
}
