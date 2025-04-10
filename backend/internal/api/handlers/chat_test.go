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
	"go-crypto-bot-clean/backend/internal/domain/ai/service/function"
	"go-crypto-bot-clean/backend/internal/domain/ai/service/templates"
	"go-crypto-bot-clean/backend/internal/domain/ai/types"

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

func (m *MockAIService) GenerateResponseWithTemplate(ctx context.Context, userID int, templateName string, templateData templates.TemplateData) (string, error) {
	args := m.Called(ctx, userID, templateName, templateData)
	return args.String(0), args.Error(1)
}

func (m *MockAIService) ExecuteFunction(ctx context.Context, userID int, functionName string, parameters map[string]interface{}) (interface{}, error) {
	args := m.Called(ctx, userID, functionName, parameters)
	return args.Get(0), args.Error(1)
}

func (m *MockAIService) GetAvailableFunctions(ctx context.Context) []function.FunctionDefinition {
	args := m.Called(ctx)
	return args.Get(0).([]function.FunctionDefinition)
}

func (m *MockAIService) GetAvailableTemplates(ctx context.Context) []string {
	args := m.Called(ctx)
	return args.Get(0).([]string)
}

func (m *MockAIService) StoreConversation(ctx context.Context, userID int, sessionID string, messages []service.Message) error {
	args := m.Called(ctx, userID, sessionID, messages)
	return args.Error(0)
}

func (m *MockAIService) RetrieveConversation(ctx context.Context, userID int, sessionID string) (*service.ConversationMemory, error) {
	args := m.Called(ctx, userID, sessionID)
	return args.Get(0).(*service.ConversationMemory), args.Error(1)
}

func (m *MockAIService) ListUserSessions(ctx context.Context, userID int, limit int) ([]string, error) {
	args := m.Called(ctx, userID, limit)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockAIService) DeleteSession(ctx context.Context, userID int, sessionID string) error {
	args := m.Called(ctx, userID, sessionID)
	return args.Error(0)
}

func (m *MockAIService) ApplyRiskGuardrails(ctx context.Context, userID int, recommendation *service.TradeRecommendation) (*service.GuardrailsResult, error) {
	args := m.Called(ctx, userID, recommendation)
	return args.Get(0).(*service.GuardrailsResult), args.Error(1)
}

func (m *MockAIService) CreateTradeConfirmation(ctx context.Context, userID int, trade *service.TradeRequest, recommendation *service.TradeRecommendation) (*service.TradeConfirmation, error) {
	args := m.Called(ctx, userID, trade, recommendation)
	return args.Get(0).(*service.TradeConfirmation), args.Error(1)
}

func (m *MockAIService) ConfirmTrade(ctx context.Context, confirmationID string, approve bool) (*service.TradeConfirmation, error) {
	args := m.Called(ctx, confirmationID, approve)
	return args.Get(0).(*service.TradeConfirmation), args.Error(1)
}

func (m *MockAIService) ListPendingTradeConfirmations(ctx context.Context, userID int) ([]*service.TradeConfirmation, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*service.TradeConfirmation), args.Error(1)
}

func (m *MockAIService) FindSimilarMessages(ctx context.Context, query string, limit int) ([]types.SimilarMessage, error) {
	args := m.Called(ctx, query, limit)
	return args.Get(0).([]types.SimilarMessage), args.Error(1)
}

func (m *MockAIService) IndexMessage(ctx context.Context, conversationID, messageID, content string, metadata map[string]interface{}) error {
	args := m.Called(ctx, conversationID, messageID, content, metadata)
	return args.Error(0)
}

func (m *MockAIService) GenerateInsights(ctx context.Context, userID int, portfolio map[string]interface{}, tradeHistory []map[string]interface{}, insightTypes []string) ([]service.Insight, error) {
	args := m.Called(ctx, userID, portfolio, tradeHistory, insightTypes)
	return args.Get(0).([]service.Insight), args.Error(1)
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
