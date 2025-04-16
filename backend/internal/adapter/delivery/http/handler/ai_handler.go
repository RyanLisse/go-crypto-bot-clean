package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/delivery/http/response"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/usecase"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
)

type AIHandler struct {
	useCase *usecase.AIUsecase
	logger  *zerolog.Logger
}

func NewAIHandler(useCase *usecase.AIUsecase, logger *zerolog.Logger) *AIHandler {
	return &AIHandler{
		useCase: useCase,
		logger:  logger,
	}
}

// ChatRequest represents a request to the chat endpoint
type ChatRequest struct {
	UserID         string                 `json:"user_id"`
	Message        string                 `json:"message"`
	SessionID      string                 `json:"session_id,omitempty"`
	TradingContext map[string]interface{} `json:"trading_context,omitempty"`
}

// ChatResponse represents a response from the chat endpoint
type ChatResponse struct {
	Response      string                 `json:"response"`
	FunctionCalls map[string]interface{} `json:"function_calls,omitempty"`
}

// GetHistory returns the authenticated user's conversation history
func (h *AIHandler) GetHistory(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(string)
	if !ok || userID == "" {
		h.logger.Error().Msg("User ID not found in context")
		response.WriteJSON(w, http.StatusUnauthorized, response.Error("User not authenticated"))
		return
	}
	// Optionally support pagination
	limit := 10
	offset := 0
	convs, err := h.useCase.ListConversations(r.Context(), userID, limit, offset)
	if err != nil {
		h.logger.Error().Err(err).Str("user_id", userID).Msg("Failed to fetch conversation history")
		response.WriteJSON(w, http.StatusInternalServerError, response.Error("Failed to fetch conversation history"))
		return
	}
	response.WriteJSON(w, http.StatusOK, response.Success(convs))
}

// Chat handles chat requests
func (h *AIHandler) Chat(w http.ResponseWriter, r *http.Request) {
	var req ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error().Err(err).Msg("Failed to decode chat request")
		response.WriteJSON(w, http.StatusBadRequest, response.Error("Invalid request format"))
		return
	}

	if req.Message == "" {
		response.WriteJSON(w, http.StatusBadRequest, response.Error("Message cannot be empty"))
		return
	}

	// Log the incoming request
	h.logger.Info().Str("user_id", req.UserID).Str("session_id", req.SessionID).Msg("Received chat request")

	// Call the AI usecase to get a response
	aiMessage, err := h.useCase.Chat(r.Context(), req.UserID, req.Message, req.SessionID, req.TradingContext)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get AI response")
		response.WriteJSON(w, http.StatusInternalServerError, response.Error("Failed to process chat request"))
		return
	}

	// Extract function calls from metadata if present
	functionCalls := make(map[string]interface{})
	if aiMessage.Metadata != nil {
		if fc, ok := aiMessage.Metadata["function_calls"]; ok {
			functionCalls = fc.(map[string]interface{})
		}
	}

	// Create response
	resp := ChatResponse{
		Response:      aiMessage.Content,
		FunctionCalls: functionCalls,
	}

	response.WriteJSON(w, http.StatusOK, response.Success(resp))
}

func (h *AIHandler) RegisterRoutes(r chi.Router, authMiddleware func(http.Handler) http.Handler) {
	r.Route("/ai", func(r chi.Router) {
		// Chat endpoint
		r.With(authMiddleware).Post("/chat", h.Chat)
		// Conversation history endpoints
		r.With(authMiddleware).Get("/history", h.GetHistory)
		r.With(authMiddleware).Get("/conversations/{conversationID}", h.GetConversation)
		r.With(authMiddleware).Get("/conversations/{conversationID}/messages", h.GetConversationMessages)
		r.With(authMiddleware).Delete("/conversations/{conversationID}", h.DeleteConversation)
	})
}

// GetConversation returns details for a specific conversation
func (h *AIHandler) GetConversation(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(string)
	if !ok || userID == "" {
		h.logger.Error().Msg("User ID not found in context")
		response.WriteJSON(w, http.StatusUnauthorized, response.Error("User not authenticated"))
		return
	}

	conversationID := chi.URLParam(r, "conversationID")
	if conversationID == "" {
		response.WriteJSON(w, http.StatusBadRequest, response.Error("Conversation ID is required"))
		return
	}

	conversation, err := h.useCase.GetConversation(r.Context(), userID, conversationID)
	if err != nil {
		h.logger.Error().Err(err).Str("conversation_id", conversationID).Msg("Failed to fetch conversation")
		response.WriteJSON(w, http.StatusInternalServerError, response.Error("Failed to fetch conversation"))
		return
	}

	response.WriteJSON(w, http.StatusOK, response.Success(conversation))
}

// GetConversationMessages returns messages for a specific conversation
func (h *AIHandler) GetConversationMessages(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(string)
	if !ok || userID == "" {
		h.logger.Error().Msg("User ID not found in context")
		response.WriteJSON(w, http.StatusUnauthorized, response.Error("User not authenticated"))
		return
	}

	conversationID := chi.URLParam(r, "conversationID")
	if conversationID == "" {
		response.WriteJSON(w, http.StatusBadRequest, response.Error("Conversation ID is required"))
		return
	}

	// Verify the conversation belongs to the user
	_, err := h.useCase.GetConversation(r.Context(), userID, conversationID)
	if err != nil {
		h.logger.Error().Err(err).Str("conversation_id", conversationID).Msg("Failed to fetch conversation")
		response.WriteJSON(w, http.StatusNotFound, response.Error("Conversation not found or access denied"))
		return
	}

	// Get messages for the conversation
	// Support optional pagination
	limit := 50
	offset := 0

	// Parse query parameters for pagination
	if limitParam := r.URL.Query().Get("limit"); limitParam != "" {
		if parsedLimit, err := parseInt(limitParam, 1, 100); err == nil {
			limit = parsedLimit
		}
	}

	if offsetParam := r.URL.Query().Get("offset"); offsetParam != "" {
		if parsedOffset, err := parseInt(offsetParam, 0, 1000); err == nil {
			offset = parsedOffset
		}
	}

	messages, err := h.useCase.GetMessages(r.Context(), conversationID, limit, offset)
	if err != nil {
		h.logger.Error().Err(err).Str("conversation_id", conversationID).Msg("Failed to fetch messages")
		response.WriteJSON(w, http.StatusInternalServerError, response.Error("Failed to fetch messages"))
		return
	}

	response.WriteJSON(w, http.StatusOK, response.Success(messages))
}

// DeleteConversation deletes a specific conversation
func (h *AIHandler) DeleteConversation(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(string)
	if !ok || userID == "" {
		h.logger.Error().Msg("User ID not found in context")
		response.WriteJSON(w, http.StatusUnauthorized, response.Error("User not authenticated"))
		return
	}

	conversationID := chi.URLParam(r, "conversationID")
	if conversationID == "" {
		response.WriteJSON(w, http.StatusBadRequest, response.Error("Conversation ID is required"))
		return
	}

	if err := h.useCase.DeleteConversation(r.Context(), userID, conversationID); err != nil {
		h.logger.Error().Err(err).Str("conversation_id", conversationID).Msg("Failed to delete conversation")
		response.WriteJSON(w, http.StatusInternalServerError, response.Error("Failed to delete conversation"))
		return
	}

	response.WriteJSON(w, http.StatusOK, response.Success(map[string]bool{"deleted": true}))
}

// Helper function to parse integer parameters with bounds
func parseInt(s string, min, max int) (int, error) {
	var value int
	if _, err := fmt.Sscanf(s, "%d", &value); err != nil {
		return 0, err
	}
	if value < min {
		value = min
	} else if value > max {
		value = max
	}
	return value, nil
}
