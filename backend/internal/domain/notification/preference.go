package notification

// Preference represents a user's notification preference for a specific channel
type Preference struct {
	UserID    string // User ID
	Channel   string // Channel type (e.g., "email", "slack", "telegram")
	Recipient string // Channel-specific recipient (e.g., email address, chat ID)
	Enabled   bool   // Whether notifications are enabled for this channel
}
