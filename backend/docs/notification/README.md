# Notification Service

The notification service provides a unified interface for sending notifications to various channels like Telegram and Slack. It allows the trading bot to send alerts, trade notifications, and system status updates to users.

## Table of Contents

1. [Overview](#overview)
2. [Features](#features)
3. [Supported Providers](#supported-providers)
4. [Configuration](#configuration)
5. [Usage](#usage)
6. [Templates](#templates)
7. [Extending](#extending)

## Overview

The notification service is designed to be a flexible and extensible system for sending notifications from the trading bot to various channels. It supports multiple notification providers, templating, attachments, and prioritization.

## Features

- **Multiple Providers**: Send notifications to Telegram, Slack, and more
- **Templating**: Use templates to format notifications
- **Attachments**: Send images, charts, and other files
- **Prioritization**: Prioritize important notifications
- **Rate Limiting**: Respect API rate limits
- **Retry Logic**: Automatically retry failed notifications
- **Asynchronous Processing**: Non-blocking notification sending
- **Extensible**: Easy to add new notification providers

## Supported Providers

### Telegram

The Telegram provider sends notifications to Telegram channels or direct messages using the Telegram Bot API. It supports:

- Text messages with Markdown formatting
- Image attachments
- Document attachments
- Multiple chat IDs

### Slack

The Slack provider sends notifications to Slack channels or direct messages using the Slack API. It supports:

- Text messages with formatting
- Image attachments
- Document attachments
- Multiple channels
- Color-coded messages based on notification level

## Configuration

The notification service is configured using a YAML file. Here's an example configuration:

```yaml
notification:
  enabled: true
  default_providers: ["telegram", "slack"]
  queue:
    capacity: 1000
    workers: 5
  providers:
    telegram:
      enabled: true
      token: "your-telegram-bot-token"
      chat_ids:
        - "-1001234567890" # Group chat ID
        - "123456789" # User chat ID
      rate_limit: 20 # messages per minute
    slack:
      enabled: true
      token: "your-slack-token"
      channels:
        - "#trading-alerts"
        - "#system-status"
      rate_limit: 30 # messages per minute
  templates:
    trade_executed:
      title: "Trade Executed: {{ .Symbol }}"
      message: "{{ .Side }} {{ .Quantity }} {{ .Symbol }} at {{ .Price }}"
      level: "TRADE"
      providers: ["telegram", "slack"]
      priority: 10
    risk_alert:
      title: "Risk Alert: {{ .AlertType }}"
      message: "{{ .Message }}"
      level: "WARNING"
      providers: ["telegram", "slack"]
      priority: 20
    system_status:
      title: "System Status: {{ .Status }}"
      message: "{{ .Message }}"
      level: "INFO"
      providers: ["slack"]
      priority: 5
```

### Configuration Options

#### Service Configuration

- `enabled`: Whether the notification service is enabled
- `default_providers`: Default providers to use if not specified in the notification
- `queue.capacity`: Maximum number of notifications in the queue
- `queue.workers`: Number of worker goroutines processing notifications

#### Provider Configuration

##### Telegram

- `enabled`: Whether the Telegram provider is enabled
- `token`: Telegram Bot API token
- `chat_ids`: List of chat IDs to send notifications to
- `rate_limit`: Maximum number of messages per minute

##### Slack

- `enabled`: Whether the Slack provider is enabled
- `token`: Slack API token
- `channels`: List of channels to send notifications to
- `rate_limit`: Maximum number of messages per minute

#### Template Configuration

- `title`: Template for the notification title
- `message`: Template for the notification message
- `level`: Notification level (INFO, WARNING, ERROR, CRITICAL, TRADE)
- `providers`: Providers to use for this template
- `priority`: Priority of the notification (higher means more important)

## Usage

### Initializing the Service

```go
// Create logger
logger, _ := zap.NewDevelopment()

// Create notification service
service := notification.NewNotificationService(logger)

// Register providers
telegramProvider := notification.NewTelegramProvider(logger)
slackProvider := notification.NewSlackProvider(logger)

// Initialize providers
telegramProvider.Initialize(telegramConfig)
slackProvider.Initialize(slackConfig)

// Register providers with the service
service.RegisterProvider(telegramProvider)
service.RegisterProvider(slackProvider)

// Initialize service
service.Initialize(config)
```

### Sending a Simple Notification

```go
notification := notification.NewNotification(
    "System Started",
    "Trading bot has started successfully",
    notification.LevelInfo,
)
notification.Source = "SystemService"

err := service.Send(notification)
```

### Using a Template

```go
data := map[string]interface{}{
    "Symbol":   "BTCUSDT",
    "Side":     "BUY",
    "Quantity": 0.1,
    "Price":    50000.0,
}

err := service.SendWithTemplate("trade_executed", data)
```

### Sending with Attachments

```go
// Create attachment from file
fileData, _ := ioutil.ReadFile("chart.png")
attachment := notification.Attachment{
    Type:        "image",
    ContentType: "image/png",
    Data:        fileData,
    Filename:    "chart.png",
}

// Create notification
notification := notification.NewNotification(
    "Price Chart: BTCUSDT",
    "4-hour chart with moving averages",
    notification.LevelInfo,
)
notification.Source = "ChartService"
notification.Attachments = []notification.Attachment{attachment}

// Send notification
err := service.Send(notification)
```

### Handling Results

```go
// Start a goroutine to handle results
go func() {
    for result := range service.GetResults() {
        if result.Success {
            logger.Info("Notification sent successfully",
                zap.String("notification_id", result.NotificationID),
                zap.String("provider", result.ProviderName),
            )
        } else {
            logger.Error("Failed to send notification",
                zap.String("notification_id", result.NotificationID),
                zap.String("provider", result.ProviderName),
                zap.String("error", result.Error),
            )
        }
    }
}()
```

### Shutting Down

```go
// Shutdown the service
service.Shutdown()
```

## Templates

Templates use Go's `html/template` package for formatting. You can use any valid Go template syntax in your templates.

### Available Template Variables

The template variables depend on the data you pass to the `SendWithTemplate` method. Here are some common variables:

#### Trade Executed

- `Symbol`: Trading pair symbol
- `Side`: Trade side (BUY or SELL)
- `Quantity`: Trade quantity
- `Price`: Trade price

#### Risk Alert

- `AlertType`: Type of risk alert
- `Message`: Alert message

#### System Status

- `Status`: System status
- `Message`: Status message

#### Position Opened

- `Symbol`: Trading pair symbol
- `Side`: Position side (BUY or SELL)
- `Quantity`: Position quantity
- `EntryPrice`: Entry price

#### Position Closed

- `Symbol`: Trading pair symbol
- `Side`: Position side (BUY or SELL)
- `Quantity`: Position quantity
- `EntryPrice`: Entry price
- `ExitPrice`: Exit price
- `ProfitLoss`: Profit or loss amount
- `ProfitLossPercentage`: Profit or loss percentage

#### Error Alert

- `Service`: Service that generated the error
- `Message`: Error message
- `StackTrace`: Stack trace

## Extending

### Adding a New Provider

To add a new notification provider, implement the `NotificationProvider` interface:

```go
type NotificationProvider interface {
    // Initialize sets up the provider with configuration
    Initialize(config map[string]interface{}) error
    
    // Send sends a notification
    Send(ctx context.Context, notification *Notification) (*NotificationResult, error)
    
    // GetName returns the provider name
    GetName() string
    
    // IsAvailable checks if the provider is available
    IsAvailable() bool
}
```

You can use the `BaseProvider` struct to implement common functionality:

```go
type MyProvider struct {
    *notification.BaseProvider
    // Provider-specific fields
}

func NewMyProvider(logger *zap.Logger) *MyProvider {
    return &MyProvider{
        BaseProvider: notification.NewBaseProvider("my_provider"),
        // Initialize provider-specific fields
    }
}

func (p *MyProvider) Initialize(config map[string]interface{}) error {
    p.SetConfig(config)
    
    // Initialize provider-specific fields
    
    p.SetAvailable(true)
    return nil
}

func (p *MyProvider) Send(ctx context.Context, notification *Notification) (*NotificationResult, error) {
    // Send notification using provider-specific logic
    
    return &NotificationResult{
        NotificationID: notification.ID,
        ProviderName:   p.GetName(),
        Success:        true,
        Timestamp:      time.Now(),
    }, nil
}
```

Then register the provider with the notification service:

```go
myProvider := NewMyProvider(logger)
myProvider.Initialize(myConfig)
service.RegisterProvider(myProvider)
```
