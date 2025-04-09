package slack

import (
	"context"
	"fmt"
	"log"

	"go-crypto-bot-clean/backend/internal/domain/notification/ports"

	"github.com/slack-go/slack"
)

const ChannelSlack = "slack"

// Config holds configuration for the Slack adapter.
type Config struct {
	BotToken string `mapstructure:"bot_token"`
	// Recipient could be a channel ID, user ID, or email
}

// adapter implements the ports.Notifier interface for Slack.
type adapter struct {
	config Config
	client *slack.Client
}

// NewAdapter creates a new Slack notification adapter.
func NewAdapter(cfg Config) (ports.Notifier, error) {
	// Initialize the actual Slack client
	client := slack.New(cfg.BotToken)
	_, err := client.AuthTest() // Verify authentication
	if err != nil {
		log.Printf("Error authenticating Slack client: %v", err)
		return nil, fmt.Errorf("failed to authenticate slack client: %w", err)
	}
	log.Printf("Slack Notifier: Authorized")

	return &adapter{
		config: cfg,
		client: client,
	}, nil
}

// Send sends a notification via Slack.
// Recipient can be a channel ID (e.g., C12345) or user ID (e.g., U12345).
func (a *adapter) Send(ctx context.Context, recipient string, subject string, message string) error {
	// Format the message using Slack's mrkdwn
	fullMessage := fmt.Sprintf("*%s*\n\n%s", subject, message)

	// Create message options
	options := []slack.MsgOption{
		slack.MsgOptionText(fullMessage, false),
		slack.MsgOptionAsUser(true),
	}

	// Add context to ensure we respect cancellation
	select {
	case <-ctx.Done():
		log.Printf("Context cancelled before sending Slack notification to %s", recipient)
		return ctx.Err()
	default:
	}

	// Send the message
	_, _, err := a.client.PostMessageContext(ctx, recipient, options...)
	if err != nil {
		log.Printf("Error sending Slack message to %s: %v", recipient, err)
		return fmt.Errorf("failed to send slack message: %w", err)
	}

	log.Printf("Successfully sent Slack notification to %s", recipient)
	return nil
}

// Supports checks if the channel is "slack".
func (a *adapter) Supports(channel string) bool {
	return channel == ChannelSlack
}
