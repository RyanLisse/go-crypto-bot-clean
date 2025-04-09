package notification

// Preference represents a user's notification preference for a specific channel.
type Preference struct {
	UserID    string // Or appropriate user identifier type
	Channel   string // e.g., "telegram", "slack", "email"
	Recipient string // Channel-specific recipient (e.g., Telegram Chat ID, Slack User ID, email address)
	Enabled   bool   // Whether notifications are enabled for this channel
}

// TODO: Potentially add other fields like notification types (e.g., "alerts", "updates") if granularity is needed.
