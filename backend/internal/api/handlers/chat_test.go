package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go-crypto-bot-clean/backend/internal/domain/ai/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAIService mocks the AIService interface
type MockAIService struct {
	mock.Mock
}

func (m *MockAIService) GenerateResponse(ctx context.Context, userID int, message string) (string, error) {
	args := m.Called(ctx, userID, message)
	return args.String(0), args.Error(1)
}

func (m *MockAIService) ExecuteFunction(ctx context.Context, userID int, functionName string, parameters map[string]interface{}) (interface{}, error) {
	args := m.Called(ctx, userID, functionName, parameters)
	return args.Get(0), args.Error(1)
}

func (m *MockAIService) StoreConversation(ctx context.Context, userID int, sessionID string, messages []service.Message) error {
	args := m.Called(ctx, userID, sessionID, messages)
	return args.Error(0)
}

func (m *MockAIService) RetrieveConversation(ctx context.Context, userID int, sessionID string) (*service.ConversationMemory, error) {
	args := m.Called(ctx, userID, sessionID)
	return args.Get(0).(*service.ConversationMemory), args.Error(1)
}

// TestChatHandler tests the ChatHandler
func TestChatHandler(t *testing.T) {
	// Create mock AI service
	mockAIService := new(MockAIService)

	// Set up expectations
	mockAIService.On("StoreConversation", mock.Anything, 1, mock.Anything, mock.Anything).Return(nil)
	mockAIService.On("GenerateResponse", mock.Anything, 1, "Hello, AI!").Return("Hello! How can I help you with your trading today?", nil)

	// Create handler
	handler := ChatHandler(mockAIService)

	// Create request
	requestBody := map[string]interface{}{
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": "Hello, AI!",
			},
		},
		"session_id": "test-session",
	}
	requestJSON, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/api/v1/ai/chat", bytes.NewBuffer(requestJSON))

	// Add user ID to context
	ctx := context.WithValue(req.Context(), "userID", 1)
	req = req.WithContext(ctx)

	// Create response recorder
	rr := httptest.NewRecorder()

	// Call handler
	handler(rr, req)

	// Check response
	assert.Equal(t, http.StatusOK, rr.Code)

	// Parse response
	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Check response fields
	assert.Equal(t, "Hello! How can I help you with your trading today?", response["output"])
	assert.Equal(t, "test-session", response["session_id"])

	// Verify expectations
	mockAIService.AssertExpectations(t)
}

// TestFunctionHandler tests the FunctionHandler
func TestFunctionHandler(t *testing.T) {
	// Create mock AI service
	mockAIService := new(MockAIService)

	// Set up expectations
	mockAIService.On("ExecuteFunction", mock.Anything, 1, "get_market_data", mock.Anything).Return(
		map[string]interface{}{
			"symbol":    "BTC",
			"price":     50000.0,
			"change":    2.5,
			"timestamp": time.Now().Format(time.RFC3339),
		}, nil)

	// Create handler
	handler := FunctionHandler(mockAIService)

	// Create request
	requestBody := map[string]interface{}{
		"function_name": "get_market_data",
		"parameters": map[string]interface{}{
			"symbol":    "BTC",
			"timeframe": "1h",
		},
		"session_id": "test-session",
	}
	requestJSON, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/api/v1/ai/function", bytes.NewBuffer(requestJSON))

	// Add user ID to context
	ctx := context.WithValue(req.Context(), "userID", 1)
	req = req.WithContext(ctx)

	// Create response recorder
	rr := httptest.NewRecorder()

	// Call handler
	handler(rr, req)

	// Check response
	assert.Equal(t, http.StatusOK, rr.Code)

	// Parse response
	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Check response fields
	result, ok := response["result"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "BTC", result["symbol"])
	assert.Equal(t, 50000.0, result["price"])
	assert.Equal(t, 2.5, result["change"])
	assert.Equal(t, "test-session", response["session_id"])

	// Verify expectations
	mockAIService.AssertExpectations(t)
}
