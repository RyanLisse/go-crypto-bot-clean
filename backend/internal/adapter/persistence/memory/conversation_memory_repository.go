package memory

import (
	"context"
	"errors"
	"sort"
	"sync"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/rs/zerolog"
)

// ConversationMemoryRepository is an in-memory implementation of the ConversationMemoryRepository interface
type ConversationMemoryRepository struct {
	conversations map[string]*model.AIConversation
	messages      map[string][]*model.AIMessage
	mu            sync.RWMutex
	logger        zerolog.Logger
}

// NewConversationMemoryRepository creates a new ConversationMemoryRepository
func NewConversationMemoryRepository(logger zerolog.Logger) *ConversationMemoryRepository {
	return &ConversationMemoryRepository{
		conversations: make(map[string]*model.AIConversation),
		messages:      make(map[string][]*model.AIMessage),
		logger:        logger.With().Str("component", "conversation_memory_repository").Logger(),
	}
}

// SaveConversation saves a conversation
func (r *ConversationMemoryRepository) SaveConversation(ctx context.Context, conversation *model.AIConversation) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Update timestamp
	conversation.UpdatedAt = time.Now()

	// Save conversation
	r.conversations[conversation.ID] = conversation

	r.logger.Debug().Str("conversation_id", conversation.ID).Msg("Saved conversation")
	return nil
}

// GetConversation retrieves a conversation by ID
func (r *ConversationMemoryRepository) GetConversation(ctx context.Context, conversationID string) (*model.AIConversation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	conversation, exists := r.conversations[conversationID]
	if !exists {
		return nil, errors.New("conversation not found")
	}

	return conversation, nil
}

// ListConversations lists conversations for a user
func (r *ConversationMemoryRepository) ListConversations(ctx context.Context, userID string, limit, offset int) ([]*model.AIConversation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Filter conversations by user ID
	var userConversations []*model.AIConversation
	for _, conversation := range r.conversations {
		if conversation.UserID == userID {
			userConversations = append(userConversations, conversation)
		}
	}

	// Sort conversations by updated_at (newest first)
	sort.Slice(userConversations, func(i, j int) bool {
		return userConversations[i].UpdatedAt.After(userConversations[j].UpdatedAt)
	})

	// Apply pagination
	if offset >= len(userConversations) {
		return []*model.AIConversation{}, nil
	}

	end := offset + limit
	if end > len(userConversations) {
		end = len(userConversations)
	}

	return userConversations[offset:end], nil
}

// DeleteConversation deletes a conversation
func (r *ConversationMemoryRepository) DeleteConversation(ctx context.Context, conversationID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if conversation exists
	if _, exists := r.conversations[conversationID]; !exists {
		return errors.New("conversation not found")
	}

	// Delete conversation
	delete(r.conversations, conversationID)

	// Delete associated messages
	delete(r.messages, conversationID)

	r.logger.Debug().Str("conversation_id", conversationID).Msg("Deleted conversation")
	return nil
}

// SaveMessage saves a message
func (r *ConversationMemoryRepository) SaveMessage(ctx context.Context, message *model.AIMessage) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if conversation exists
	if _, exists := r.conversations[message.ConversationID]; !exists {
		return errors.New("conversation not found")
	}

	// Initialize messages slice if it doesn't exist
	if _, exists := r.messages[message.ConversationID]; !exists {
		r.messages[message.ConversationID] = []*model.AIMessage{}
	}

	// Save message
	r.messages[message.ConversationID] = append(r.messages[message.ConversationID], message)

	// Update conversation timestamp
	r.conversations[message.ConversationID].UpdatedAt = time.Now()

	r.logger.Debug().Str("conversation_id", message.ConversationID).Str("message_id", message.ID).Msg("Saved message")
	return nil
}

// GetMessages retrieves messages for a conversation
func (r *ConversationMemoryRepository) GetMessages(ctx context.Context, conversationID string, limit, offset int) ([]*model.AIMessage, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Check if conversation exists
	if _, exists := r.conversations[conversationID]; !exists {
		return nil, errors.New("conversation not found")
	}

	// Get messages
	messages, exists := r.messages[conversationID]
	if !exists {
		return []*model.AIMessage{}, nil
	}

	// Sort messages by timestamp (newest first)
	sort.Slice(messages, func(i, j int) bool {
		return messages[i].Timestamp.After(messages[j].Timestamp)
	})

	// Apply pagination
	if offset >= len(messages) {
		return []*model.AIMessage{}, nil
	}

	end := offset + limit
	if end > len(messages) {
		end = len(messages)
	}

	return messages[offset:end], nil
}

// GetMessage gets a message by ID
func (r *ConversationMemoryRepository) GetMessage(ctx context.Context, messageID string) (*model.AIMessage, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Search for message in all conversations
	for _, messages := range r.messages {
		for _, message := range messages {
			if message.ID == messageID {
				return message, nil
			}
		}
	}

	return nil, errors.New("message not found")
}

// DeleteMessage deletes a message
func (r *ConversationMemoryRepository) DeleteMessage(ctx context.Context, messageID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Search for message in all conversations
	for conversationID, messages := range r.messages {
		for i, message := range messages {
			if message.ID == messageID {
				// Remove message from slice
				r.messages[conversationID] = append(messages[:i], messages[i+1:]...)
				r.logger.Debug().Str("message_id", messageID).Msg("Deleted message")
				return nil
			}
		}
	}

	return errors.New("message not found")
}
