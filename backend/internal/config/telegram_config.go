package config

// TelegramConfig contains configuration for Telegram notifications
type TelegramConfig struct {
	Enabled  bool   `mapstructure:"enabled"`
	BotToken string `mapstructure:"bot_token"`
	ChatID   string `mapstructure:"chat_id"`
	// Optional additional chat IDs for different notification types
	AlertChatID    string `mapstructure:"alert_chat_id"`
	TradeChatID    string `mapstructure:"trade_chat_id"`
	DebugChatID    string `mapstructure:"debug_chat_id"`
	APIBaseURL     string `mapstructure:"api_base_url"`
	DisableWebPagePreview bool `mapstructure:"disable_web_page_preview"`
	ParseMode      string `mapstructure:"parse_mode"` // "Markdown" or "HTML"
}

// GetDefaultTelegramConfig returns the default Telegram configuration
func GetDefaultTelegramConfig() TelegramConfig {
	return TelegramConfig{
		Enabled:             false,
		BotToken:            "",
		ChatID:              "",
		AlertChatID:         "",
		TradeChatID:         "",
		DebugChatID:         "",
		APIBaseURL:          "https://api.telegram.org",
		DisableWebPagePreview: true,
		ParseMode:           "Markdown",
	}
}
