package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"go-crypto-bot-clean/backend/internal/domain/ai/service"
)

// ConversationMemoryRepository defines the interface for storing and retrieving conversation memories
type ConversationMemoryRepository interface {
	StoreConversation(ctx context.Context, userID int, sessionID string, messages []service.Message) error
	RetrieveConversation(ctx context.Context, userID int, sessionID string) (*service.ConversationMemory, error)
	ListUserSessions(ctx context.Context, userID int, limit int) ([]string, error)
	DeleteSession(ctx context.Context, userID int, sessionID string) error
}

// SQLiteConversationMemoryRepository implements ConversationMemoryRepository using SQLite
type SQLiteConversationMemoryRepository struct {
	db *sql.DB
}

// NewSQLiteConversationMemoryRepository creates a new SQLiteConversationMemoryRepository
func NewSQLiteConversationMemoryRepository(db *sql.DB) (*SQLiteConversationMemoryRepository, error) {
	// Create table if it doesn't exist
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS conversation_memories (
		user_id INTEGER NOT NULL,
		session_id TEXT NOT NULL,
		messages_json TEXT NOT NULL,
		summary TEXT,
		last_accessed TIMESTAMP NOT NULL,
		PRIMARY KEY (user_id, session_id)
	);
	CREATE INDEX IF NOT EXISTS idx_conversation_memories_user_id ON conversation_memories(user_id);
	CREATE INDEX IF NOT EXISTS idx_conversation_memories_last_accessed ON conversation_memories(last_accessed);
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to create conversation_memories table: %w", err)
	}

	return &SQLiteConversationMemoryRepository{db: db}, nil
}

// StoreConversation saves a conversation to the database
func (r *SQLiteConversationMemoryRepository) StoreConversation(
	ctx context.Context,
	userID int,
	sessionID string,
	messages []service.Message,
) error {
	// Convert messages to JSON
	messagesJSON, err := json.Marshal(messages)
	if err != nil {
		return fmt.Errorf("failed to marshal messages: %w", err)
	}

	// Insert or update conversation
	_, err = r.db.ExecContext(
		ctx,
		`INSERT INTO conversation_memories (user_id, session_id, messages_json, last_accessed)
		VALUES (?, ?, ?, ?)
		ON CONFLICT(user_id, session_id) DO UPDATE SET
		messages_json = ?,
		last_accessed = ?`,
		userID, sessionID, messagesJSON, time.Now().UTC(),
		messagesJSON, time.Now().UTC(),
	)
	if err != nil {
		return fmt.Errorf("failed to store conversation: %w", err)
	}

	return nil
}

// RetrieveConversation gets a conversation from the database
func (r *SQLiteConversationMemoryRepository) RetrieveConversation(
	ctx context.Context,
	userID int,
	sessionID string,
) (*service.ConversationMemory, error) {
	var messagesJSON string
	var summary sql.NullString
	var lastAccessed time.Time

	err := r.db.QueryRowContext(
		ctx,
		`SELECT messages_json, summary, last_accessed FROM conversation_memories
		WHERE user_id = ? AND session_id = ?`,
		userID, sessionID,
	).Scan(&messagesJSON, &summary, &lastAccessed)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No conversation found
		}
		return nil, fmt.Errorf("failed to retrieve conversation: %w", err)
	}

	// Parse messages JSON
	var messages []service.Message
	err = json.Unmarshal([]byte(messagesJSON), &messages)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal messages: %w", err)
	}

	// Update last accessed time
	_, err = r.db.ExecContext(
		ctx,
		`UPDATE conversation_memories SET last_accessed = ? WHERE user_id = ? AND session_id = ?`,
		time.Now().UTC(), userID, sessionID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update last accessed time: %w", err)
	}

	return &service.ConversationMemory{
		UserID:       userID,
		SessionID:    sessionID,
		Messages:     messages,
		Summary:      summary.String,
		LastAccessed: lastAccessed,
	}, nil
}

// ListUserSessions lists all sessions for a user
func (r *SQLiteConversationMemoryRepository) ListUserSessions(
	ctx context.Context,
	userID int,
	limit int,
) ([]string, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT session_id FROM conversation_memories
		WHERE user_id = ?
		ORDER BY last_accessed DESC
		LIMIT ?`,
		userID, limit,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list user sessions: %w", err)
	}
	defer rows.Close()

	var sessions []string
	for rows.Next() {
		var sessionID string
		err := rows.Scan(&sessionID)
		if err != nil {
			return nil, fmt.Errorf("failed to scan session ID: %w", err)
		}
		sessions = append(sessions, sessionID)
	}

	return sessions, nil
}

// DeleteSession deletes a session
func (r *SQLiteConversationMemoryRepository) DeleteSession(
	ctx context.Context,
	userID int,
	sessionID string,
) error {
	_, err := r.db.ExecContext(
		ctx,
		`DELETE FROM conversation_memories WHERE user_id = ? AND session_id = ?`,
		userID, sessionID,
	)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	return nil
}
