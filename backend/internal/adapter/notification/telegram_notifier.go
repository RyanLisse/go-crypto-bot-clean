package notification

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/config"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model/status"
	"github.com/rs/zerolog"
)

// TelegramNotifier implements notification via Telegram Bot API
type TelegramNotifier struct {
	logger           *zerolog.Logger
	config           config.TelegramConfig
	httpClient       *http.Client
	notificationChan chan TelegramMessage
	enabled          bool
	mutex            sync.RWMutex
	lastSentTime     time.Time
	minInterval      time.Duration // Minimum time between messages to avoid flooding
}

// TelegramMessage represents a message to be sent via Telegram
type TelegramMessage struct {
	ChatID                string `json:"chat_id"`
	Text                  string `json:"text"`
	ParseMode             string `json:"parse_mode,omitempty"`
	DisableWebPagePreview bool   `json:"disable_web_page_preview,omitempty"`
	DisableNotification   bool   `json:"disable_notification,omitempty"`
	ReplyToMessageID      int    `json:"reply_to_message_id,omitempty"`
	Type                  string `json:"-"` // Internal use only, not sent to Telegram
}

// NewTelegramNotifier creates a new Telegram notifier
func NewTelegramNotifier(config config.TelegramConfig, logger *zerolog.Logger) *TelegramNotifier {
	notifier := &TelegramNotifier{
		logger:           logger,
		config:           config,
		httpClient:       &http.Client{Timeout: 10 * time.Second},
		notificationChan: make(chan TelegramMessage, 100),
		enabled:          config.Enabled,
		minInterval:      500 * time.Millisecond, // Minimum 500ms between messages
	}

	// Start the notification processor if enabled
	if notifier.enabled {
		go notifier.processNotifications()
	} else {
		logger.Info().Msg("Telegram notifications are disabled")
	}

	return notifier
}

// NotifyStatusChange sends a notification about a status change
func (n *TelegramNotifier) NotifyStatusChange(ctx context.Context, component string, oldStatus, newStatus status.Status, message string) error {
	if !n.enabled {
		return nil
	}

	// Format the message
	text := fmt.Sprintf("*Status Change*\nComponent: `%s`\nStatus: `%s` → `%s`", component, oldStatus, newStatus)
	if message != "" {
		text += fmt.Sprintf("\nDetails: %s", message)
	}

	// Determine which chat ID to use
	chatID := n.config.ChatID
	if n.config.AlertChatID != "" {
		chatID = n.config.AlertChatID
	}

	// Send the message
	return n.sendMessage(ctx, TelegramMessage{
		ChatID:                chatID,
		Text:                  text,
		ParseMode:             n.config.ParseMode,
		DisableWebPagePreview: n.config.DisableWebPagePreview,
		Type:                  "status",
	})
}

// NotifySystemStatusChange sends a notification about a system status change
func (n *TelegramNotifier) NotifySystemStatusChange(ctx context.Context, oldStatus, newStatus status.Status, message string) error {
	if !n.enabled {
		return nil
	}

	// Format the message
	text := fmt.Sprintf("*System Status Change*\nStatus: `%s` → `%s`", oldStatus, newStatus)
	if message != "" {
		text += fmt.Sprintf("\nDetails: %s", message)
	}

	// Determine which chat ID to use
	chatID := n.config.ChatID
	if n.config.AlertChatID != "" {
		chatID = n.config.AlertChatID
	}

	// Send the message
	return n.sendMessage(ctx, TelegramMessage{
		ChatID:                chatID,
		Text:                  text,
		ParseMode:             n.config.ParseMode,
		DisableWebPagePreview: n.config.DisableWebPagePreview,
		Type:                  "system",
	})
}

// NotifyTrade sends a notification about a trade
func (n *TelegramNotifier) NotifyTrade(ctx context.Context, symbol, side, orderType string, quantity, price float64, orderID string) error {
	if !n.enabled {
		return nil
	}

	// Format the message
	text := fmt.Sprintf("*Trade Executed*\nSymbol: `%s`\nSide: `%s`\nType: `%s`\nQuantity: `%.8f`\nPrice: `%.8f`\nOrder ID: `%s`",
		symbol, side, orderType, quantity, price, orderID)

	// Determine which chat ID to use
	chatID := n.config.ChatID
	if n.config.TradeChatID != "" {
		chatID = n.config.TradeChatID
	}

	// Send the message
	return n.sendMessage(ctx, TelegramMessage{
		ChatID:                chatID,
		Text:                  text,
		ParseMode:             n.config.ParseMode,
		DisableWebPagePreview: n.config.DisableWebPagePreview,
		Type:                  "trade",
	})
}

// NotifyAlert sends an alert notification
func (n *TelegramNotifier) NotifyAlert(ctx context.Context, level, title, message, source string) error {
	if !n.enabled {
		return nil
	}

	// Format the message
	text := fmt.Sprintf("*Alert: %s*\nLevel: `%s`\nSource: `%s`\nDetails: %s",
		title, level, source, message)

	// Determine which chat ID to use
	chatID := n.config.ChatID
	if n.config.AlertChatID != "" {
		chatID = n.config.AlertChatID
	}

	// Send the message
	return n.sendMessage(ctx, TelegramMessage{
		ChatID:                chatID,
		Text:                  text,
		ParseMode:             n.config.ParseMode,
		DisableWebPagePreview: n.config.DisableWebPagePreview,
		DisableNotification:   level == "info", // Only send notification for non-info alerts
		Type:                  "alert",
	})
}

// NotifyDebug sends a debug message
func (n *TelegramNotifier) NotifyDebug(ctx context.Context, message string) error {
	if !n.enabled {
		return nil
	}

	// Only send debug messages if debug chat ID is configured
	if n.config.DebugChatID == "" {
		return nil
	}

	// Send the message
	return n.sendMessage(ctx, TelegramMessage{
		ChatID:                n.config.DebugChatID,
		Text:                  message,
		ParseMode:             n.config.ParseMode,
		DisableWebPagePreview: n.config.DisableWebPagePreview,
		DisableNotification:   true, // Don't send notification for debug messages
		Type:                  "debug",
	})
}

// sendMessage queues a message to be sent to Telegram
func (n *TelegramNotifier) sendMessage(ctx context.Context, message TelegramMessage) error {
	if !n.enabled {
		return nil
	}

	// Validate required fields
	if message.ChatID == "" {
		return fmt.Errorf("chat ID is required")
	}
	if message.Text == "" {
		return fmt.Errorf("message text is required")
	}

	// Try to send to channel with timeout
	select {
	case n.notificationChan <- message:
		n.logger.Debug().
			Str("type", message.Type).
			Str("chatID", message.ChatID).
			Msg("Telegram message queued")
	case <-time.After(time.Second):
		n.logger.Warn().
			Str("type", message.Type).
			Str("chatID", message.ChatID).
			Msg("Failed to queue Telegram message: channel full")
		return fmt.Errorf("notification channel is full")
	case <-ctx.Done():
		return ctx.Err()
	}

	return nil
}

// processNotifications processes queued notifications
func (n *TelegramNotifier) processNotifications() {
	for message := range n.notificationChan {
		// Rate limit to avoid flooding Telegram
		n.mutex.RLock()
		timeSinceLastSent := time.Since(n.lastSentTime)
		n.mutex.RUnlock()

		if timeSinceLastSent < n.minInterval {
			sleepTime := n.minInterval - timeSinceLastSent
			time.Sleep(sleepTime)
		}

		// Send the message
		err := n.doSendMessage(message)
		if err != nil {
			n.logger.Error().
				Err(err).
				Str("type", message.Type).
				Str("chatID", message.ChatID).
				Msg("Failed to send Telegram message")
		}

		// Update last sent time
		n.mutex.Lock()
		n.lastSentTime = time.Now()
		n.mutex.Unlock()
	}
}

// doSendMessage actually sends the message to Telegram
func (n *TelegramNotifier) doSendMessage(message TelegramMessage) error {
	// Construct the API URL
	url := fmt.Sprintf("%s/bot%s/sendMessage", n.config.APIBaseURL, n.config.BotToken)

	// Marshal the message to JSON
	jsonData, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	// Create the request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := n.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		var errorResponse struct {
			OK          bool   `json:"ok"`
			ErrorCode   int    `json:"error_code"`
			Description string `json:"description"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return fmt.Errorf("failed to decode error response: %w", err)
		}
		return fmt.Errorf("telegram API error: %d - %s", errorResponse.ErrorCode, errorResponse.Description)
	}

	n.logger.Debug().
		Str("type", message.Type).
		Str("chatID", message.ChatID).
		Msg("Telegram message sent successfully")

	return nil
}

// Enable enables the notifier
func (n *TelegramNotifier) Enable() {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	n.enabled = true
}

// Disable disables the notifier
func (n *TelegramNotifier) Disable() {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	n.enabled = false
}

// IsEnabled returns whether the notifier is enabled
func (n *TelegramNotifier) IsEnabled() bool {
	n.mutex.RLock()
	defer n.mutex.RUnlock()
	return n.enabled
}
