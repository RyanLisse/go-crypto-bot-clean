package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/neo/crypto-bot/internal/domain/model"
	"github.com/neo/crypto-bot/internal/usecase"
	"github.com/rs/zerolog"
)

// AIHandler handles AI-related HTTP requests
type AIHandler struct {
	aiUsecase *usecase.AIUsecase
	logger    zerolog.Logger
}

// NewAIHandler creates a new AIHandler
func NewAIHandler(aiUsecase *usecase.AIUsecase, logger zerolog.Logger) *AIHandler {
	return &AIHandler{
		aiUsecase: aiUsecase,
		logger:    logger.With().Str("component", "ai_handler").Logger(),
	}
}

// ChatRequest represents a request to the chat endpoint
type ChatRequest struct {
	Message        string `json:"message"`
	ConversationID string `json:"conversation_id,omitempty"`
}

// ChatResponse represents a response from the chat endpoint
type ChatResponse struct {
	Message        model.AIMessage `json:"message"`
	ConversationID string          `json:"conversation_id"`
}

// HandleChat handles chat requests
func (h *AIHandler) HandleChat(c *gin.Context) {
	// Parse request
	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// Validate request
	if req.Message == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Message is required"})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID := getUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Send message to AI
	response, err := h.aiUsecase.Chat(c.Request.Context(), userID, req.Message, req.ConversationID)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get AI response")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get AI response"})
		return
	}

	// Return response
	c.JSON(http.StatusOK, ChatResponse{
		Message:        *response,
		ConversationID: response.ConversationID,
	})
}

// ListConversationsResponse represents a response from the list conversations endpoint
type ListConversationsResponse struct {
	Conversations []*model.AIConversation `json:"conversations"`
	Total         int                     `json:"total"`
}

// HandleListConversations handles list conversations requests
func (h *AIHandler) HandleListConversations(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID := getUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Parse pagination parameters
	limit, offset := getPaginationParams(c)

	// Get conversations
	conversations, err := h.aiUsecase.ListConversations(c.Request.Context(), userID, limit, offset)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to list conversations")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list conversations"})
		return
	}

	// Return response
	c.JSON(http.StatusOK, ListConversationsResponse{
		Conversations: conversations,
		Total:         len(conversations), // This is not accurate for pagination, but we don't have a count method
	})
}

// HandleGetConversation handles get conversation requests
func (h *AIHandler) HandleGetConversation(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID := getUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get conversation ID from path
	conversationID := c.Param("id")
	if conversationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Conversation ID is required"})
		return
	}

	// Get conversation
	conversation, err := h.aiUsecase.GetConversation(c.Request.Context(), userID, conversationID)
	if err != nil {
		h.logger.Error().Err(err).Str("conversation_id", conversationID).Msg("Failed to get conversation")
		c.JSON(http.StatusNotFound, gin.H{"error": "Conversation not found"})
		return
	}

	// Return response
	c.JSON(http.StatusOK, conversation)
}

// HandleDeleteConversation handles delete conversation requests
func (h *AIHandler) HandleDeleteConversation(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID := getUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get conversation ID from path
	conversationID := c.Param("id")
	if conversationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Conversation ID is required"})
		return
	}

	// Delete conversation
	err := h.aiUsecase.DeleteConversation(c.Request.Context(), userID, conversationID)
	if err != nil {
		h.logger.Error().Err(err).Str("conversation_id", conversationID).Msg("Failed to delete conversation")
		c.JSON(http.StatusNotFound, gin.H{"error": "Conversation not found"})
		return
	}

	// Return response
	c.JSON(http.StatusOK, gin.H{"message": "Conversation deleted"})
}

// InsightRequest represents a request to the insight endpoint
type InsightRequest struct {
	InsightType string                 `json:"insight_type"`
	Data        map[string]interface{} `json:"data"`
}

// HandleGenerateInsight handles generate insight requests
func (h *AIHandler) HandleGenerateInsight(c *gin.Context) {
	// Parse request
	var req InsightRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// Validate request
	if req.InsightType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insight type is required"})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID := getUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Generate insight
	insight, err := h.aiUsecase.GenerateInsight(c.Request.Context(), userID, req.InsightType, req.Data)
	if err != nil {
		h.logger.Error().Err(err).Str("insight_type", req.InsightType).Msg("Failed to generate insight")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate insight"})
		return
	}

	// Return response
	c.JSON(http.StatusOK, insight)
}

// TradeRecommendationRequest represents a request to the trade recommendation endpoint
type TradeRecommendationRequest struct {
	Data map[string]interface{} `json:"data"`
}

// HandleGenerateTradeRecommendation handles generate trade recommendation requests
func (h *AIHandler) HandleGenerateTradeRecommendation(c *gin.Context) {
	// Parse request
	var req TradeRecommendationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID := getUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Generate trade recommendation
	recommendation, err := h.aiUsecase.GenerateTradeRecommendation(c.Request.Context(), userID, req.Data)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to generate trade recommendation")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate trade recommendation"})
		return
	}

	// Return response
	c.JSON(http.StatusOK, recommendation)
}

// FunctionCallRequest represents a request to the function call endpoint
type FunctionCallRequest struct {
	Name       string                 `json:"name"`
	Parameters map[string]interface{} `json:"parameters"`
}

// HandleExecuteFunction handles execute function requests
func (h *AIHandler) HandleExecuteFunction(c *gin.Context) {
	// Parse request
	var req FunctionCallRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// Validate request
	if req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Function name is required"})
		return
	}

	// Create function call
	functionCall := model.AIFunctionCall{
		Name:       req.Name,
		Parameters: req.Parameters,
	}

	// Execute function
	response, err := h.aiUsecase.ExecuteFunction(c.Request.Context(), functionCall)
	if err != nil {
		h.logger.Error().Err(err).Str("function", req.Name).Msg("Failed to execute function")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute function"})
		return
	}

	// Return response
	c.JSON(http.StatusOK, response)
}

// RegisterRoutes registers the AI handler routes
func (h *AIHandler) RegisterRoutes(router *gin.Engine, authMiddleware gin.HandlerFunc) {
	aiGroup := router.Group("/api/v1/ai")
	aiGroup.Use(authMiddleware)
	{
		aiGroup.POST("/chat", h.HandleChat)
		aiGroup.GET("/conversations", h.HandleListConversations)
		aiGroup.GET("/conversations/:id", h.HandleGetConversation)
		aiGroup.DELETE("/conversations/:id", h.HandleDeleteConversation)
		aiGroup.POST("/insights", h.HandleGenerateInsight)
		aiGroup.POST("/recommendations", h.HandleGenerateTradeRecommendation)
		aiGroup.POST("/functions", h.HandleExecuteFunction)
	}
}

// Helper functions

// getUserID gets the user ID from the context
func getUserID(_ *gin.Context) string {
	// In a real implementation, this would get the user ID from the JWT token
	// For now, we'll use a dummy user ID
	return "user123"
}

// getPaginationParams gets pagination parameters from the request
func getPaginationParams(c *gin.Context) (int, int) {
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	return limit, offset
}
