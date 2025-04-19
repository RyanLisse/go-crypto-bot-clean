package gateway

import (
	"context"
	"testing"

	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/config"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestAIGateway_GenerateText(t *testing.T) {
	// Create a logger
	logger := zerolog.Nop()

	// Create a config
	aiConfig := config.AIConfig{
		Provider:    "mock",
		APIKey:      "test-key",
		ModelName:   "test-model",
		Temperature: 0.7,
		MaxTokens:   100,
	}

	// Create a gateway
	gateway, err := NewAIGateway(aiConfig, &logger)
	assert.NoError(t, err)

	// Test with empty prompt
	text, err := gateway.GenerateText(context.Background(), "")
	assert.Error(t, err)
	assert.Empty(t, text)

	// Test with valid prompt
	text, err = gateway.GenerateText(context.Background(), "Hello, AI!")
	assert.NoError(t, err)
	assert.Contains(t, text, "Hello, AI!")
}

func TestAIGateway_AnalyzeSentiment(t *testing.T) {
	// Create a logger
	logger := zerolog.Nop()

	// Create a config
	aiConfig := config.AIConfig{
		Provider:    "mock",
		APIKey:      "test-key",
		ModelName:   "test-model",
		Temperature: 0.7,
		MaxTokens:   100,
	}

	// Create a gateway
	gateway, err := NewAIGateway(aiConfig, &logger)
	assert.NoError(t, err)

	// Test with empty text
	sentiment, err := gateway.AnalyzeSentiment(context.Background(), "")
	assert.Error(t, err)
	assert.Empty(t, sentiment)

	// Test with short text (should return "Neutral")
	sentiment, err = gateway.AnalyzeSentiment(context.Background(), "Hi")
	assert.NoError(t, err)
	assert.Equal(t, "Neutral", sentiment)

	// Test with even-length text (should return "Positive")
	sentiment, err = gateway.AnalyzeSentiment(context.Background(), "Hello there!")
	assert.NoError(t, err)
	assert.Equal(t, "Positive", sentiment)

	// Test with odd-length text (should return "Negative")
	// The string "Hello world!" has length 12, which is even
	// Let's use "Hello world" instead, which has length 11 (odd)
	sentiment, err = gateway.AnalyzeSentiment(context.Background(), "Hello world")
	assert.NoError(t, err)
	assert.Equal(t, "Negative", sentiment)
}
