package gateway

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/config"
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/port/gateway"
	"github.com/rs/zerolog"
)

// Ensure AIGateway implements the gateway.AIGateway interface.
var _ gateway.AIGateway = (*AIGateway)(nil)

// AIGateway implements the gateway.AIGateway interface
type AIGateway struct {
	config     *config.Config
	logger     *zerolog.Logger
	httpClient *http.Client
	apiKey     string
}

// NewAIGateway creates a new adapter for interacting with AI services.
func NewAIGateway(config *config.Config, logger *zerolog.Logger) gateway.AIGateway {
	// Check if AI API key is set
	if config.AI.APIKey == "" {
		logger.Warn().Msg("AI_API_KEY not set, AI features will not work properly")
	}

	return &AIGateway{
		config:     config,
		logger:     logger,
		httpClient: &http.Client{Timeout: 30 * time.Second},
		apiKey:     config.AI.APIKey,
	}
}

// GenerateText generates text based on a given prompt
func (g *AIGateway) GenerateText(ctx context.Context, prompt string) (string, error) {
	g.logger.Debug().Str("prompt", prompt).Msg("Generating text with AI")

	// This is a placeholder implementation
	// In a real implementation, you would call an AI service API
	if prompt == "" {
		return "", fmt.Errorf("empty prompt")
	}

	// For now, return a mock response
	return fmt.Sprintf("AI generated response for: %s", prompt), nil
}

// AnalyzeSentiment analyzes the sentiment of a given text
func (g *AIGateway) AnalyzeSentiment(ctx context.Context, text string) (string, error) {
	g.logger.Debug().Str("text", text).Msg("Analyzing sentiment with AI")

	// This is a placeholder implementation
	// In a real implementation, you would call an AI service API
	if text == "" {
		return "", fmt.Errorf("empty text")
	}

	// For now, return a mock response
	// In a real implementation, this would be based on the AI service's analysis
	if len(text) < 10 {
		return "Neutral", nil
	} else if len(text)%2 == 0 {
		return "Positive", nil
	} else {
		return "Negative", nil
	}
}
