package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-crypto-bot-clean/backend/internal/notification"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

func main() {
	// Create logger
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	// Load configuration
	config, err := loadConfig("configs/notification.yaml")
	if err != nil {
		logger.Fatal("Failed to load configuration", zap.Error(err))
	}

	// Create notification service
	service := notification.NewNotificationService(logger)

	// Register providers
	telegramProvider := notification.NewTelegramProvider(logger)
	slackProvider := notification.NewSlackProvider(logger)

	// Initialize providers
	if telegramConfig, ok := config.Providers["telegram"].(map[string]interface{}); ok {
		if err := telegramProvider.Initialize(telegramConfig); err != nil {
			logger.Warn("Failed to initialize Telegram provider", zap.Error(err))
		} else {
			service.RegisterProvider(telegramProvider)
		}
	}

	if slackConfig, ok := config.Providers["slack"].(map[string]interface{}); ok {
		if err := slackProvider.Initialize(slackConfig); err != nil {
			logger.Warn("Failed to initialize Slack provider", zap.Error(err))
		} else {
			service.RegisterProvider(slackProvider)
		}
	}

	// Initialize service
	if err := service.Initialize(config); err != nil {
		logger.Fatal("Failed to initialize notification service", zap.Error(err))
	}

	// Start result handler
	go handleResults(service, logger)

	// Set up signal handling for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigCh
		logger.Info("Received signal", zap.String("signal", sig.String()))
		cancel()
	}()

	// Send example notifications
	sendExampleNotifications(ctx, service, logger)

	// Wait for shutdown signal
	<-ctx.Done()

	// Shutdown service
	service.Shutdown()
	logger.Info("Service shut down")
}

// loadConfig loads the notification service configuration from a YAML file
func loadConfig(filename string) (*notification.NotificationServiceConfig, error) {
	// Read file
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse YAML
	var rawConfig map[string]interface{}
	if err := yaml.Unmarshal(data, &rawConfig); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Extract notification section
	notificationConfig, ok := rawConfig["notification"].(map[interface{}]interface{})
	if !ok {
		return nil, fmt.Errorf("notification section not found in config")
	}

	// Convert to proper types
	config := &notification.NotificationServiceConfig{
		DefaultProviders: make([]string, 0),
		QueueCapacity:    1000,
		WorkerCount:      5,
		Templates:        make(map[string]interface{}),
		Providers:        make(map[string]interface{}),
	}

	// Default providers
	if providers, ok := notificationConfig["default_providers"].([]interface{}); ok {
		for _, provider := range providers {
			if providerStr, ok := provider.(string); ok {
				config.DefaultProviders = append(config.DefaultProviders, providerStr)
			}
		}
	}

	// Queue config
	if queueConfig, ok := notificationConfig["queue"].(map[interface{}]interface{}); ok {
		if capacity, ok := queueConfig["capacity"].(int); ok {
			config.QueueCapacity = capacity
		}
		if workers, ok := queueConfig["workers"].(int); ok {
			config.WorkerCount = workers
		}
	}

	// Templates
	if templates, ok := notificationConfig["templates"].(map[interface{}]interface{}); ok {
		for id, templateConfig := range templates {
			if idStr, ok := id.(string); ok {
				if templateMap, ok := templateConfig.(map[interface{}]interface{}); ok {
					// Convert map[interface{}]interface{} to map[string]interface{}
					convertedMap := make(map[string]interface{})
					for k, v := range templateMap {
						if kStr, ok := k.(string); ok {
							convertedMap[kStr] = v
						}
					}
					config.Templates[idStr] = convertedMap
				}
			}
		}
	}

	// Providers
	if providers, ok := notificationConfig["providers"].(map[interface{}]interface{}); ok {
		for name, providerConfig := range providers {
			if nameStr, ok := name.(string); ok {
				if providerMap, ok := providerConfig.(map[interface{}]interface{}); ok {
					// Convert map[interface{}]interface{} to map[string]interface{}
					convertedMap := make(map[string]interface{})
					for k, v := range providerMap {
						if kStr, ok := k.(string); ok {
							// Special handling for chat_ids and channels
							if kStr == "chat_ids" || kStr == "channels" {
								if items, ok := v.([]interface{}); ok {
									strItems := make([]string, 0, len(items))
									for _, item := range items {
										if itemStr, ok := item.(string); ok {
											strItems = append(strItems, itemStr)
										}
									}
									convertedMap[kStr] = strItems
								}
							} else {
								convertedMap[kStr] = v
							}
						}
					}
					config.Providers[nameStr] = convertedMap
				}
			}
		}
	}

	return config, nil
}

// handleResults handles notification results
func handleResults(service *notification.NotificationService, logger *zap.Logger) {
	for result := range service.GetResults() {
		if result.Success {
			logger.Info("Notification sent successfully",
				zap.String("notification_id", result.NotificationID),
				zap.String("provider", result.ProviderName),
				zap.Time("timestamp", result.Timestamp),
			)
		} else {
			logger.Error("Failed to send notification",
				zap.String("notification_id", result.NotificationID),
				zap.String("provider", result.ProviderName),
				zap.String("error", result.Error),
				zap.Time("timestamp", result.Timestamp),
			)
		}
	}
}

// sendExampleNotifications sends example notifications
func sendExampleNotifications(ctx context.Context, service *notification.NotificationService, logger *zap.Logger) {
	// Example 1: Simple notification
	notification1 := notification.NewNotification(
		"System Started",
		"Trading bot has started successfully",
		notification.LevelInfo,
	)
	notification1.Source = "SystemService"

	if err := service.Send(notification1); err != nil {
		logger.Error("Failed to send notification", zap.Error(err))
	}

	// Example 2: Using a template
	tradeData := map[string]interface{}{
		"Symbol":   "BTCUSDT",
		"Side":     "BUY",
		"Quantity": 0.1,
		"Price":    50000.0,
	}

	if err := service.SendWithTemplate("trade_executed", tradeData); err != nil {
		logger.Error("Failed to send notification with template", zap.Error(err))
	}

	// Example 3: Risk alert
	riskData := map[string]interface{}{
		"AlertType": "Drawdown",
		"Message":   "Portfolio drawdown has exceeded 5%",
	}

	if err := service.SendWithTemplate("risk_alert", riskData); err != nil {
		logger.Error("Failed to send notification with template", zap.Error(err))
	}

	// Example 4: System status
	statusData := map[string]interface{}{
		"Status":  "Healthy",
		"Message": "All systems operational",
	}

	if err := service.SendWithTemplate("system_status", statusData); err != nil {
		logger.Error("Failed to send notification with template", zap.Error(err))
	}

	// Example 5: Position opened
	positionData := map[string]interface{}{
		"Symbol":     "ETHUSDT",
		"Side":       "BUY",
		"Quantity":   2.5,
		"EntryPrice": 3000.0,
	}

	if err := service.SendWithTemplate("position_opened", positionData); err != nil {
		logger.Error("Failed to send notification with template", zap.Error(err))
	}

	// Wait a bit between notifications to avoid rate limiting
	time.Sleep(1 * time.Second)

	// Example 6: Position closed
	closedPositionData := map[string]interface{}{
		"Symbol":              "ETHUSDT",
		"Side":                "BUY",
		"Quantity":            2.5,
		"EntryPrice":          3000.0,
		"ExitPrice":           3150.0,
		"ProfitLoss":          375.0,
		"ProfitLossPercentage": 5.0,
	}

	if err := service.SendWithTemplate("position_closed", closedPositionData); err != nil {
		logger.Error("Failed to send notification with template", zap.Error(err))
	}

	// Example 7: Error alert
	errorData := map[string]interface{}{
		"Service":    "ExchangeService",
		"Message":    "Failed to connect to exchange API",
		"StackTrace": "goroutine 1 [running]:\nmain.main()\n\t/app/main.go:42 +0x7f",
	}

	if err := service.SendWithTemplate("error_alert", errorData); err != nil {
		logger.Error("Failed to send notification with template", zap.Error(err))
	}
}
