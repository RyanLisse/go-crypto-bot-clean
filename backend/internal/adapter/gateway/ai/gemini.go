package ai

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/config"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/rs/zerolog"
)

// GeminiAIService implements the AIService interface using Google's Gemini API
type GeminiAIService struct {
	config      *config.Config
	logger      zerolog.Logger
	functionReg map[string]FunctionHandler
}

// FunctionHandler is a function that handles a function call
type FunctionHandler func(ctx context.Context, params map[string]interface{}) (interface{}, error)

// NewGeminiAIService creates a new GeminiAIService
func NewGeminiAIService(cfg *config.Config, logger zerolog.Logger) (*GeminiAIService, error) {
	if cfg.AI.GeminiAPIKey == "" {
		return nil, errors.New("gemini API key is required")
	}

	// This is a mock implementation for development purposes
	service := &GeminiAIService{
		config:      cfg,
		logger:      logger.With().Str("component", "gemini_ai_service").Logger(),
		functionReg: make(map[string]FunctionHandler),
	}

	// Register default functions
	service.registerDefaultFunctions()

	return service, nil
}

// Chat sends a message to the AI and returns a response
func (s *GeminiAIService) Chat(ctx context.Context, message string, conversationID string) (*model.AIMessage, error) {
	// This is a mock implementation for development purposes
	s.logger.Info().Str("message", message).Msg("Received chat message")

	// Create a mock response
	responseText := "This is a mock response from the Gemini AI service. The actual implementation would connect to the Gemini API."

	aiMessage := &model.AIMessage{
		ID:             uuid.New().String(),
		ConversationID: conversationID,
		Role:           "assistant",
		Content:        responseText,
		Timestamp:      time.Now(),
		Metadata: map[string]interface{}{
			"model": s.config.AI.GeminiModel,
		},
	}

	return aiMessage, nil
}

// ChatWithHistory sends a message with conversation history to the AI
func (s *GeminiAIService) ChatWithHistory(ctx context.Context, messages []model.AIMessage) (*model.AIMessage, error) {
	// This is a mock implementation for development purposes
	s.logger.Info().Int("messageCount", len(messages)).Msg("Received chat with history")

	// Get the last user message for logging
	var lastUserMessage string
	for i := len(messages) - 1; i >= 0; i-- {
		if messages[i].Role == "user" {
			lastUserMessage = messages[i].Content
			break
		}
	}

	s.logger.Info().Str("lastMessage", lastUserMessage).Msg("Processing last message")

	// Create a mock response
	responseText := "This is a mock response from the Gemini AI service with conversation history. The actual implementation would connect to the Gemini API."

	conversationID := ""
	if len(messages) > 0 {
		conversationID = messages[0].ConversationID
	}

	aiMessage := &model.AIMessage{
		ID:             uuid.New().String(),
		ConversationID: conversationID,
		Role:           "assistant",
		Content:        responseText,
		Timestamp:      time.Now(),
		Metadata: map[string]interface{}{
			"model": s.config.AI.GeminiModel,
		},
	}

	return aiMessage, nil
}

// GenerateInsight generates an insight based on provided data
func (s *GeminiAIService) GenerateInsight(ctx context.Context, insightType string, data map[string]interface{}) (*model.AIInsight, error) {
	// This is a mock implementation for development purposes
	// Convert data to JSON string for logging
	dataJSON, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %w", err)
	}

	s.logger.Info().Str("insightType", insightType).RawJSON("data", dataJSON).Msg("Generating insight")

	// Create mock title and description
	title := fmt.Sprintf("Mock %s Insight", insightType)
	description := "This is a mock insight generated for development purposes. The actual implementation would use the Gemini API to generate a real insight based on the provided data."

	// Create insight
	insight := &model.AIInsight{
		ID:          uuid.New().String(),
		UserID:      data["user_id"].(string),
		Type:        insightType,
		Title:       title,
		Description: description,
		CreatedAt:   time.Now(),
		Confidence:  0.85, // Default confidence
		Metadata: map[string]interface{}{
			"model": s.config.AI.GeminiModel,
			"data":  data,
		},
	}

	return insight, nil
}

// GenerateTradeRecommendation generates a trade recommendation
func (s *GeminiAIService) GenerateTradeRecommendation(ctx context.Context, data map[string]interface{}) (*model.AITradeRecommendation, error) {
	// This is a mock implementation for development purposes
	// Convert data to JSON string for logging
	dataJSON, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %w", err)
	}

	s.logger.Info().RawJSON("data", dataJSON).Msg("Generating trade recommendation")

	// Create mock recommendation
	symbol := "BTC/USDT"
	action := "buy"
	quantity := 0.01
	reasoning := "This is a mock trade recommendation generated for development purposes. The actual implementation would use the Gemini API to generate a real recommendation based on market analysis."

	// Create recommendation
	recommendation := &model.AITradeRecommendation{
		ID:         uuid.New().String(),
		UserID:     data["user_id"].(string),
		Symbol:     symbol,
		Action:     action,
		Quantity:   quantity,
		Reasoning:  reasoning,
		CreatedAt:  time.Now(),
		ExpiresAt:  time.Now().Add(24 * time.Hour),
		Confidence: 0.75, // Default confidence
		Status:     "pending",
	}

	return recommendation, nil
}

// ExecuteFunction executes a function call from the AI
func (s *GeminiAIService) ExecuteFunction(ctx context.Context, functionCall model.AIFunctionCall) (*model.AIFunctionResponse, error) {
	// Check if function exists
	handler, exists := s.functionReg[functionCall.Name]
	if !exists {
		return nil, fmt.Errorf("function %s not found", functionCall.Name)
	}

	// Execute function
	result, err := handler(ctx, functionCall.Parameters)
	if err != nil {
		s.logger.Error().Err(err).Str("function", functionCall.Name).Msg("Failed to execute function")
		return nil, fmt.Errorf("failed to execute function %s: %w", functionCall.Name, err)
	}

	// Create response
	response := &model.AIFunctionResponse{
		Name:   functionCall.Name,
		Result: result,
	}

	return response, nil
}

// GenerateEmbedding generates a vector embedding for a text
func (s *GeminiAIService) GenerateEmbedding(ctx context.Context, text string) (*model.AIEmbedding, error) {
	// This is a mock implementation for development purposes
	s.logger.Info().Str("text", text).Msg("Generating embedding")

	// Create mock embedding
	embedding := &model.AIEmbedding{
		ID:         uuid.New().String(),
		SourceID:   uuid.New().String(),
		SourceType: "message",
		Vector:     make([]float64, 128), // Mock 128-dimensional vector
		CreatedAt:  time.Now(),
	}

	s.logger.Info().Msg("Generated mock embedding")

	return embedding, nil
}

// RegisterFunction registers a function handler
func (s *GeminiAIService) RegisterFunction(name string, handler FunctionHandler) {
	s.functionReg[name] = handler
}

// registerDefaultFunctions registers default function handlers
func (s *GeminiAIService) registerDefaultFunctions() {
	// Register get_market_data function
	s.RegisterFunction("get_market_data", func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
		symbol, ok := params["symbol"].(string)
		if !ok {
			return nil, errors.New("symbol parameter is required")
		}

		// Mock implementation
		return map[string]interface{}{
			"symbol":     symbol,
			"price":      42000.0,
			"change_24h": 2.5,
			"volume_24h": 1000000.0,
		}, nil
	})

	// Register get_portfolio function
	s.RegisterFunction("get_portfolio", func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
		// Mock implementation
		return map[string]interface{}{
			"total_value": 12345.67,
			"assets": []map[string]interface{}{
				{
					"symbol":   "BTC",
					"quantity": 0.5,
					"value":    21000.0,
				},
				{
					"symbol":   "ETH",
					"quantity": 5.0,
					"value":    10000.0,
				},
			},
		}, nil
	})
}
