package config

// NotificationProvider represents a notification provider configuration
type NotificationProvider struct {
	Enabled   bool     `mapstructure:"enabled" validate:"required"`
	Token     string   `mapstructure:"token" validate:"required_if=Enabled true"`
	ChatIDs   []string `mapstructure:"chat_ids,omitempty"`
	Channels  []string `mapstructure:"channels,omitempty"`
	RateLimit int      `mapstructure:"rate_limit" validate:"required,min=1"`
}

// NotificationTemplate represents a notification template
type NotificationTemplate struct {
	Title     string   `mapstructure:"title" validate:"required"`
	Message   string   `mapstructure:"message" validate:"required"`
	Level     string   `mapstructure:"level" validate:"required,oneof=INFO WARNING ERROR TRADE"`
	Providers []string `mapstructure:"providers" validate:"required,min=1"`
	Priority  int      `mapstructure:"priority" validate:"required,min=1,max=100"`
}

// NotificationQueue represents the notification queue configuration
type NotificationQueue struct {
	Capacity int `mapstructure:"capacity" validate:"required,min=10"`
	Workers  int `mapstructure:"workers" validate:"required,min=1"`
}

// NotificationConfig represents the notification configuration
type NotificationConfig struct {
	Notification struct {
		Enabled          bool                               `mapstructure:"enabled" validate:"required"`
		DefaultProviders []string                           `mapstructure:"default_providers" validate:"required,min=1"`
		Queue            NotificationQueue                  `mapstructure:"queue" validate:"required"`
		Providers        map[string]NotificationProvider    `mapstructure:"providers" validate:"required,min=1"`
		Templates        map[string]NotificationTemplate    `mapstructure:"templates" validate:"required,min=1"`
	} `mapstructure:"notification" validate:"required"`
}
