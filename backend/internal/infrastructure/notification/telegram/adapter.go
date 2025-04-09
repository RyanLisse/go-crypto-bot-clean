package telegram

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"go-crypto-bot-clean/backend/internal/domain/notification/ports"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const ChannelTelegram = "telegram"

// Config holds configuration for the Telegram adapter.
type Config struct {
	BotToken string `mapstructure:"bot_token"`
	// TODO: Consider how to manage chat IDs - per recipient or globally?
	// DefaultChatID string `mapstructure:"default_chat_id"`
}

// adapter implements the ports.Notifier interface for Telegram.
type adapter struct {
	config Config
	client *tgbotapi.BotAPI
}

// NewAdapter creates a new Telegram notification adapter.
// It will initialize the Telegram bot client.
func NewAdapter(cfg Config) (ports.Notifier, error) {
	// Initialize the actual Telegram client
	bot, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		log.Printf("Error initializing Telegram bot API: %v", err)
		return nil, fmt.Errorf("failed to initialize telegram bot api: %w", err)
	}
	// TODO: Make debug mode configurable via Config
	// bot.Debug = true

	log.Printf("Telegram Notifier: Authorized on account %s", bot.Self.UserName)

	return &adapter{
		config: cfg,
		client: bot,
	}, nil
}

// Send sends a notification via Telegram.
// The recipient is expected to be the Chat ID as a string.
func (a *adapter) Send(ctx context.Context, recipient string, subject string, message string) error {
	chatID, err := strconv.ParseInt(recipient, 10, 64)
	if err != nil {
		log.Printf("Error parsing Telegram recipient Chat ID '%s': %v", recipient, err)
		return fmt.Errorf("invalid telegram recipient format (expected chat ID): %w", err)
	}

	fullMessage := fmt.Sprintf("*%s*\n\n%s", escapeMarkdown(subject), escapeMarkdown(message))
	msg := tgbotapi.NewMessage(chatID, fullMessage)
	msg.ParseMode = tgbotapi.ModeMarkdown

	select {
	case <-ctx.Done():
		log.Printf("Context cancelled before sending Telegram notification to %d", chatID)
		return ctx.Err()
	default:
	}

	_, err = a.client.Send(msg)
	if err != nil {
		log.Printf("Error sending Telegram message to %d: %v", chatID, err)
		return fmt.Errorf("failed to send telegram message: %w", err)
	}

	log.Printf("Successfully sent Telegram notification to Chat ID: %d", chatID)
	return nil
}

// Supports checks if the channel is "telegram".
func (a *adapter) Supports(channel string) bool {
	return channel == ChannelTelegram
}

// escapeMarkdown escapes characters that have special meaning in Telegram's MarkdownV1.
// Note: Telegram's Markdown support can be tricky. This is a basic escaping function.
// For complex formatting, consider HTML or MarkdownV2 (which requires different escaping).
func escapeMarkdown(text string) string {
	var result string
	for _, r := range text {
		switch r {
		case '_', '*', '`', '[':
			result += "\\" + string(r)
		default:
			result += string(r)
		}
	}
	return result
}
