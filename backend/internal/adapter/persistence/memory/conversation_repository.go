package memory

import (
	"context"
	"errors"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/rs/zerolog"
)

// ConversationMemoryRepository implements the ConversationMemoryRepository interface using in-memory storage
type ConversationMemoryRepository struct {
	conversations map[string]*model.AIConversation
	messages      map[string][]*model.AIMessage
	mu            sync.RWMutex
	logger        zerolog.Logger
}

// NewConversationMemoryRepository creates a new ConversationMemoryRepository
func NewConversationMemoryRepository(logger zerolog.Logger) port.ConversationMemoryRepository {
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

	// If the conversation doesn't have an ID, generate one
	if conversation.ID == "" {
		conversation.ID = uuid.New().String()
	}

	// Update timestamps
	if conversation.CreatedAt.IsZero() {
		conversation.CreatedAt = time.Now()
	}
	conversation.UpdatedAt = time.Now()

	// Save the conversation
	r.conversations[conversation.ID] = conversation

	// Initialize messages array if it doesn't exist
	if _, exists := r.messages[conversation.ID]; !exists {
		r.messages[conversation.ID] = make([]*model.AIMessage, 0)
	}

	// Save any messages in the conversation
	for i := range conversation.Messages {
		msg := conversation.Messages[i]
		if msg.ID == "" {
			msg.ID = uuid.New().String()
		}
		msg.ConversationID = conversation.ID
		if msg.Timestamp.IsZero() {
			msg.Timestamp = time.Now()
		}

		// Check if message already exists
		exists := false
		for _, existingMsg := range r.messages[conversation.ID] {
			if existingMsg.ID == msg.ID {
				exists = true
				break
			}
		}

		// Add message if it doesn't exist
		if !exists {
			r.messages[conversation.ID] = append(r.messages[conversation.ID], &msg)
		}
	}

	r.logger.Debug().Str("conversation_id", conversation.ID).Msg("Saved conversation")
	return nil
}

// GetConversation retrieves a conversation by ID
func (r *ConversationMemoryRepository) GetConversation(ctx context.Context, id string) (*model.AIConversation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	conversation, exists := r.conversations[id]
	if !exists {
		return nil, errors.New("conversation not found")
	}

	// Create a copy of the conversation
	conversationCopy := *conversation

	// Get messages for the conversation
	messages, exists := r.messages[id]
	if exists {
		// Create copies of messages
		conversationCopy.Messages = make([]model.AIMessage, len(messages))
		for i, msg := range messages {
			conversationCopy.Messages[i] = *msg
		}

		// Sort messages by timestamp
		sort.Slice(conversationCopy.Messages, func(i, j int) bool {
			return conversationCopy.Messages[i].Timestamp.Before(conversationCopy.Messages[j].Timestamp)
		})
	} else {
		conversationCopy.Messages = make([]model.AIMessage, 0)
	}

	return &conversationCopy, nil
}

// ListConversations lists conversations for a user
func (r *ConversationMemoryRepository) ListConversations(ctx context.Context, userID string, limit, offset int) ([]*model.AIConversation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Find conversations for the user
	var conversations []*model.AIConversation
	for _, conversation := range r.conversations {
		if conversation.UserID == userID {
			// Create a copy of the conversation
			conversationCopy := *conversation
			conversationCopy.Messages = make([]model.AIMessage, 0) // Don't include messages in the list
			conversations = append(conversations, &conversationCopy)
		}
	}

	// Sort conversations by updated_at (newest first)
	sort.Slice(conversations, func(i, j int) bool {
		return conversations[i].UpdatedAt.After(conversations[j].UpdatedAt)
	})

	// Apply pagination
	if offset >= len(conversations) {
		return []*model.AIConversation{}, nil
	}

	end := offset + limit
	if end > len(conversations) {
		end = len(conversations)
	}

	return conversations[offset:end], nil
}

// SaveMessage saves a message to a conversation
func (r *ConversationMemoryRepository) SaveMessage(ctx context.Context, message *model.AIMessage) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if the conversation exists
	if _, exists := r.conversations[message.ConversationID]; !exists {
		return errors.New("conversation not found")
	}

	// If the message doesn't have an ID, generate one
	if message.ID == "" {
		message.ID = uuid.New().String()
	}

	// Set timestamp if not set
	if message.Timestamp.IsZero() {
		message.Timestamp = time.Now()
	}

	// Initialize messages array if it doesn't exist
	if _, exists := r.messages[message.ConversationID]; !exists {
		r.messages[message.ConversationID] = make([]*model.AIMessage, 0)
	}

	// Check if message already exists
	for i, existingMsg := range r.messages[message.ConversationID] {
		if existingMsg.ID == message.ID {
			// Update existing message
			r.messages[message.ConversationID][i] = message
			r.logger.Debug().Str("message_id", message.ID).Msg("Updated message")
			return nil
		}
	}

	// Add new message
	r.messages[message.ConversationID] = append(r.messages[message.ConversationID], message)

	// Update conversation's updated_at timestamp
	if conversation, exists := r.conversations[message.ConversationID]; exists {
		conversation.UpdatedAt = time.Now()
	}

	r.logger.Debug().Str("message_id", message.ID).Str("conversation_id", message.ConversationID).Msg("Saved message")
	return nil
}

// GetMessages retrieves messages for a conversation
func (r *ConversationMemoryRepository) GetMessages(ctx context.Context, conversationID string, limit, offset int) ([]*model.AIMessage, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Check if the conversation exists
	if _, exists := r.conversations[conversationID]; !exists {
		return nil, errors.New("conversation not found")
	}

	// Get messages for the conversation
	messages, exists := r.messages[conversationID]
	if !exists {
		return []*model.AIMessage{}, nil
	}

	// Create copies of messages
	messagesCopy := make([]*model.AIMessage, len(messages))
	for i, msg := range messages {
		msgCopy := *msg
		messagesCopy[i] = &msgCopy
	}

	// Sort messages by timestamp
	sort.Slice(messagesCopy, func(i, j int) bool {
		return messagesCopy[i].Timestamp.Before(messagesCopy[j].Timestamp)
	})

	// Apply pagination
	if offset >= len(messagesCopy) {
		return []*model.AIMessage{}, nil
	}

	end := offset + limit
	if end > len(messagesCopy) {
		end = len(messagesCopy)
	}

	return messagesCopy[offset:end], nil
}

// DeleteConversation deletes a conversation
func (r *ConversationMemoryRepository) DeleteConversation(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if the conversation exists
	if _, exists := r.conversations[id]; !exists {
		return errors.New("conversation not found")
	}

	// Delete the conversation
	delete(r.conversations, id)

	// Delete messages for the conversation
	delete(r.messages, id)

	r.logger.Debug().Str("conversation_id", id).Msg("Deleted conversation")
	return nil
}
