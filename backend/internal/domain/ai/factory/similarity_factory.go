package factory

import (
	"fmt"

	"go-crypto-bot-clean/backend/internal/domain/ai/repository"
	"go-crypto-bot-clean/backend/internal/domain/ai/service/embeddings"
	"go-crypto-bot-clean/backend/internal/domain/ai/similarity"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// CreateSimilarityService creates a new similarity service
func CreateSimilarityService(db *gorm.DB, logger *zap.Logger) (*similarity.Service, error) {
	// Create embeddings repository
	embeddingsRepo, err := repository.NewEmbeddingsRepository(db)
	if err != nil {
		return nil, fmt.Errorf("failed to create embeddings repository: %w", err)
	}

	// Create primary embeddings service (Gemini)
	geminiEmbeddingsService, err := embeddings.NewGeminiEmbeddingsService(logger)
	if err != nil {
		logger.Warn("Failed to create Gemini embeddings service, will use OpenAI only", zap.Error(err))
	}

	// Create fallback embeddings service (OpenAI)
	openaiEmbeddingsService, err := embeddings.NewOpenAIEmbeddingsService(logger)
	if err != nil {
		if geminiEmbeddingsService == nil {
			return nil, fmt.Errorf("failed to create any embeddings service: %w", err)
		}
		logger.Warn("Failed to create OpenAI embeddings service, will use Gemini only", zap.Error(err))
	}

	// Create fallback embeddings service
	var embeddingsService embeddings.EmbeddingsService
	if geminiEmbeddingsService != nil && openaiEmbeddingsService != nil {
		// Use both with Gemini as primary
		embeddingsService = embeddings.NewFallbackEmbeddingsService(geminiEmbeddingsService, openaiEmbeddingsService, logger)
		logger.Info("Using Gemini embeddings with OpenAI fallback")
	} else if geminiEmbeddingsService != nil {
		// Use Gemini only
		embeddingsService = geminiEmbeddingsService
		logger.Info("Using Gemini embeddings only")
	} else {
		// Use OpenAI only
		embeddingsService = openaiEmbeddingsService
		logger.Info("Using OpenAI embeddings only")
	}

	// Create similarity service
	similarityService := similarity.NewService(embeddingsRepo, embeddingsService, logger)

	return similarityService, nil
}
