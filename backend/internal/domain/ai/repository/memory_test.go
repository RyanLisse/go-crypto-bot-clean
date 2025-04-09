package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/ryanlisse/go-crypto-bot/internal/domain/ai/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)
	require.NotNil(t, db)
	return db
}

func TestSQLiteConversationMemoryRepository(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo, err := NewSQLiteConversationMemoryRepository(db)
	require.NoError(t, err)
	require.NotNil(t, repo)

	// Test storing and retrieving a conversation
	t.Run("StoreAndRetrieveConversation", func(t *testing.T) {
		ctx := context.Background()
		userID := 1
		sessionID := "test-session-1"
		messages := []service.Message{
			{
				Role:      "user",
				Content:   "Hello, AI!",
				Timestamp: time.Now().Add(-1 * time.Minute),
			},
			{
				Role:      "assistant",
				Content:   "Hello! How can I help you with your trading today?",
				Timestamp: time.Now(),
			},
		}

		// Store conversation
		err := repo.StoreConversation(ctx, userID, sessionID, messages)
		assert.NoError(t, err)

		// Retrieve conversation
		conversation, err := repo.RetrieveConversation(ctx, userID, sessionID)
		assert.NoError(t, err)
		assert.NotNil(t, conversation)
		assert.Equal(t, userID, conversation.UserID)
		assert.Equal(t, sessionID, conversation.SessionID)
		assert.Len(t, conversation.Messages, 2)
		assert.Equal(t, "user", conversation.Messages[0].Role)
		assert.Equal(t, "Hello, AI!", conversation.Messages[0].Content)
		assert.Equal(t, "assistant", conversation.Messages[1].Role)
		assert.Equal(t, "Hello! How can I help you with your trading today?", conversation.Messages[1].Content)
	})

	// Test updating an existing conversation
	t.Run("UpdateConversation", func(t *testing.T) {
		ctx := context.Background()
		userID := 1
		sessionID := "test-session-2"
		
		// Initial messages
		initialMessages := []service.Message{
			{
				Role:      "user",
				Content:   "What's the price of Bitcoin?",
				Timestamp: time.Now().Add(-2 * time.Minute),
			},
		}

		// Store initial conversation
		err := repo.StoreConversation(ctx, userID, sessionID, initialMessages)
		assert.NoError(t, err)

		// Updated messages
		updatedMessages := append(initialMessages, service.Message{
			Role:      "assistant",
			Content:   "The current price of Bitcoin is $50,000.",
			Timestamp: time.Now().Add(-1 * time.Minute),
		})

		// Update conversation
		err = repo.StoreConversation(ctx, userID, sessionID, updatedMessages)
		assert.NoError(t, err)

		// Retrieve updated conversation
		conversation, err := repo.RetrieveConversation(ctx, userID, sessionID)
		assert.NoError(t, err)
		assert.NotNil(t, conversation)
		assert.Len(t, conversation.Messages, 2)
		assert.Equal(t, "user", conversation.Messages[0].Role)
		assert.Equal(t, "What's the price of Bitcoin?", conversation.Messages[0].Content)
		assert.Equal(t, "assistant", conversation.Messages[1].Role)
		assert.Equal(t, "The current price of Bitcoin is $50,000.", conversation.Messages[1].Content)
	})

	// Test listing user sessions
	t.Run("ListUserSessions", func(t *testing.T) {
		ctx := context.Background()
		userID := 2
		
		// Create multiple sessions
		sessions := []string{"session-1", "session-2", "session-3"}
		for _, sessionID := range sessions {
			messages := []service.Message{
				{
					Role:      "user",
					Content:   "Test message for " + sessionID,
					Timestamp: time.Now(),
				},
			}
			err := repo.StoreConversation(ctx, userID, sessionID, messages)
			assert.NoError(t, err)
		}

		// List sessions
		retrievedSessions, err := repo.ListUserSessions(ctx, userID, 10)
		assert.NoError(t, err)
		assert.Len(t, retrievedSessions, 3)
		
		// Check that all sessions are in the list
		for _, sessionID := range sessions {
			assert.Contains(t, retrievedSessions, sessionID)
		}
	})

	// Test deleting a session
	t.Run("DeleteSession", func(t *testing.T) {
		ctx := context.Background()
		userID := 3
		sessionID := "session-to-delete"
		
		// Create session
		messages := []service.Message{
			{
				Role:      "user",
				Content:   "This session will be deleted",
				Timestamp: time.Now(),
			},
		}
		err := repo.StoreConversation(ctx, userID, sessionID, messages)
		assert.NoError(t, err)

		// Verify session exists
		conversation, err := repo.RetrieveConversation(ctx, userID, sessionID)
		assert.NoError(t, err)
		assert.NotNil(t, conversation)

		// Delete session
		err = repo.DeleteSession(ctx, userID, sessionID)
		assert.NoError(t, err)

		// Verify session no longer exists
		conversation, err = repo.RetrieveConversation(ctx, userID, sessionID)
		assert.NoError(t, err)
		assert.Nil(t, conversation)
	})
}
