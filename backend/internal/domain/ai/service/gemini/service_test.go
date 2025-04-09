package gemini

import (
	"context"
	"testing"

	"github.com/google/generative-ai-go/genai"
	"github.com/ryanlisse/go-crypto-bot/internal/domain/ai/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock dependencies
type MockMemoryRepository struct {
	mock.Mock
}

func (m *MockMemoryRepository) StoreConversation(ctx context.Context, userID int, sessionID string, messages []service.Message) error {
	args := m.Called(ctx, userID, sessionID, messages)
	return args.Error(0)
}

func (m *MockMemoryRepository) RetrieveConversation(ctx context.Context, userID int, sessionID string) (*service.ConversationMemory, error) {
	args := m.Called(ctx, userID, sessionID)
	return args.Get(0).(*service.ConversationMemory), args.Error(1)
}

// Add the DeleteSession method to the MockMemoryRepository
func (m *MockMemoryRepository) DeleteSession(ctx context.Context, userID int, sessionID string) error {
	args := m.Called(ctx, userID, sessionID)
	return args.Error(0)
}

// Add the ListUserSessions method to the MockMemoryRepository
func (m *MockMemoryRepository) ListUserSessions(ctx context.Context, userID int, limit int) ([]string, error) {
	args := m.Called(ctx, userID, limit)
	return args.Get(0).([]string), args.Error(1)
}

// TestGeminiAIService_ExecuteFunction tests the ExecuteFunction method
func TestGeminiAIService_ExecuteFunction(t *testing.T) {
	// Create mock memory repository
	mockMemoryRepo := new(MockMemoryRepository)

	// Create service
	service := &GeminiAIService{
		Client:     &genai.Client{},
		MemoryRepo: mockMemoryRepo,
	}

	// Test get_market_data function
	t.Run("get_market_data", func(t *testing.T) {
		// Set up parameters
		params := map[string]interface{}{
			"symbol":    "BTC",
			"timeframe": "1h",
		}

		// Call function
		result, err := service.ExecuteFunction(context.Background(), 1, "get_market_data", params)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)

		// Check result
		resultMap, ok := result.(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, "BTC", resultMap["symbol"])
		assert.Equal(t, "1h", resultMap["timeframe"])
	})

	// Test analyze_technical_indicators function
	t.Run("analyze_technical_indicators", func(t *testing.T) {
		// Set up parameters
		params := map[string]interface{}{
			"symbol":     "ETH",
			"indicators": []interface{}{"rsi", "macd"},
		}

		// Call function
		result, err := service.ExecuteFunction(context.Background(), 1, "analyze_technical_indicators", params)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)

		// Check result
		resultMap, ok := result.(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, "ETH", resultMap["symbol"])

		// Check indicators
		indicators, ok := resultMap["indicators"].(map[string]interface{})
		assert.True(t, ok)
		assert.Contains(t, indicators, "rsi")
		assert.Contains(t, indicators, "macd")
	})

	// Note: execute_trade function test removed as it's not implemented in the current version

	// Test unknown function
	t.Run("unknown_function", func(t *testing.T) {
		// Call function
		result, err := service.ExecuteFunction(context.Background(), 1, "unknown_function", nil)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "unknown function")
	})
}
