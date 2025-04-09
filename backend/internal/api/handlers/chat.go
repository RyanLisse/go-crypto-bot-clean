package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/ryanlisse/go-crypto-bot/internal/domain/ai/service"
)

// ChatHandler handles chat requests
func ChatHandler(
	aiSvc service.AIService,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract user ID from context (set by auth middleware)
		userID, ok := r.Context().Value("userID").(int)
		if !ok {
			http.Error(w, "User ID not found in context", http.StatusUnauthorized)
			return
		}

		// Parse request body
		var requestBody struct {
			Messages []struct {
				Role    string `json:"role"`
				Content string `json:"content"`
			} `json:"messages"`
			SessionID string `json:"session_id,omitempty"`
		}
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Get the last user message
		var userMessage string
		for i := len(requestBody.Messages) - 1; i >= 0; i-- {
			if requestBody.Messages[i].Role == "user" {
				userMessage = requestBody.Messages[i].Content
				break
			}
		}
		if userMessage == "" {
			http.Error(w, "No user message provided", http.StatusBadRequest)
			return
		}

		// Generate session ID if not provided
		sessionID := requestBody.SessionID
		if sessionID == "" {
			sessionID = uuid.New().String()
		}

		// Convert messages to internal format
		var messages []service.Message
		for _, msg := range requestBody.Messages {
			messages = append(messages, service.Message{
				Role:      msg.Role,
				Content:   msg.Content,
				Timestamp: time.Now(),
			})
		}

		// Store conversation
		if err := aiSvc.StoreConversation(r.Context(), userID, sessionID, messages); err != nil {
			log.Printf("Error storing conversation: %v", err)
			// Continue without storing
		}

		// Generate response
		response, err := aiSvc.GenerateResponse(r.Context(), userID, userMessage)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to generate response: %v", err), http.StatusInternalServerError)
			return
		}

		// Add assistant message to conversation
		assistantMessage := service.Message{
			Role:      "assistant",
			Content:   response,
			Timestamp: time.Now(),
		}
		messages = append(messages, assistantMessage)

		// Update conversation in database
		if err := aiSvc.StoreConversation(r.Context(), userID, sessionID, messages); err != nil {
			log.Printf("Error updating conversation: %v", err)
			// Continue without updating
		}

		// Return response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"output":     response,
			"session_id": sessionID,
		})
	}
}

// FunctionHandler handles function execution requests
func FunctionHandler(
	aiSvc service.AIService,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract user ID from context (set by auth middleware)
		userID, ok := r.Context().Value("userID").(int)
		if !ok {
			http.Error(w, "User ID not found in context", http.StatusUnauthorized)
			return
		}

		// Parse request body
		var requestBody struct {
			FunctionName string                 `json:"function_name"`
			Parameters   map[string]interface{} `json:"parameters"`
			SessionID    string                 `json:"session_id,omitempty"`
		}
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Execute function
		result, err := aiSvc.ExecuteFunction(r.Context(), userID, requestBody.FunctionName, requestBody.Parameters)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to execute function: %v", err), http.StatusInternalServerError)
			return
		}

		// Return result
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"result":     result,
			"session_id": requestBody.SessionID,
		})
	}
}
