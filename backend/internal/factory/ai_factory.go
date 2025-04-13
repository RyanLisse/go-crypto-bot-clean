package factory

import (
	"github.com/neo/crypto-bot/internal/adapter/gateway/ai"
	"github.com/neo/crypto-bot/internal/adapter/http/handler"
	"github.com/neo/crypto-bot/internal/adapter/persistence/memory"
	"github.com/neo/crypto-bot/internal/config"
	"github.com/neo/crypto-bot/internal/domain/port"
	"github.com/neo/crypto-bot/internal/usecase"
	"github.com/rs/zerolog"
)

// AIFactory creates AI-related components
type AIFactory struct {
	config *config.Config
	logger zerolog.Logger
}

// NewAIFactory creates a new AIFactory
func NewAIFactory(config *config.Config, logger zerolog.Logger) *AIFactory {
	return &AIFactory{
		config: config,
		logger: logger.With().Str("component", "ai_factory").Logger(),
	}
}

// CreateAIService creates an AIService based on the configuration
func (f *AIFactory) CreateAIService() (port.AIService, error) {
	// Create AI service based on provider
	switch f.config.AI.Provider {
	case "gemini":
		return ai.NewGeminiAIService(f.config, f.logger)
	default:
		f.logger.Warn().Str("provider", f.config.AI.Provider).Msg("Unknown AI provider, using Gemini")
		return ai.NewGeminiAIService(f.config, f.logger)
	}
}

// CreateConversationMemoryRepository creates a ConversationMemoryRepository
func (f *AIFactory) CreateConversationMemoryRepository() port.ConversationMemoryRepository {
	// For now, we'll use an in-memory repository
	// In a real implementation, this would use a database
	return memory.NewConversationMemoryRepository(f.logger)
}

// CreateEmbeddingRepository creates an EmbeddingRepository
func (f *AIFactory) CreateEmbeddingRepository() port.EmbeddingRepository {
	// For now, we'll return nil
	// In a real implementation, this would use a vector database
	return nil
}

// CreateAIUsecase creates an AIUsecase
func (f *AIFactory) CreateAIUsecase() (*usecase.AIUsecase, error) {
	// Create dependencies
	aiService, err := f.CreateAIService()
	if err != nil {
		return nil, err
	}
	
	conversationMemoryRepo := f.CreateConversationMemoryRepository()
	embeddingRepo := f.CreateEmbeddingRepository()
	
	// Create usecase
	return usecase.NewAIUsecase(aiService, conversationMemoryRepo, embeddingRepo, f.logger), nil
}

// CreateAIHandler creates an AIHandler
func (f *AIFactory) CreateAIHandler() (*handler.AIHandler, error) {
	// Create usecase
	aiUsecase, err := f.CreateAIUsecase()
	if err != nil {
		return nil, err
	}
	
	// Create handler
	return handler.NewAIHandler(aiUsecase, f.logger), nil
}
