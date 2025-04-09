package notification

import (
	"context"
	"fmt"
	"time"

	"go-crypto-bot-clean/backend/internal/domain/models"
	"go.uber.org/zap"
)

// NotificationHandler handles notifications for various trading bot events
type NotificationHandler struct {
	service *NotificationService
	logger  *zap.Logger
}

// NewNotificationHandler creates a new notification handler
func NewNotificationHandler(service *NotificationService, logger *zap.Logger) *NotificationHandler {
	if logger == nil {
		logger, _ = zap.NewProduction()
	}

	return &NotificationHandler{
		service: service,
		logger:  logger,
	}
}

// NotifyTradeExecuted sends a notification when a trade is executed
func (h *NotificationHandler) NotifyTradeExecuted(ctx context.Context, order *models.Order) error {
	data := map[string]interface{}{
		"Symbol":    order.Symbol,
		"Side":      string(order.Side),
		"Quantity":  order.Quantity,
		"Price":     order.Price,
		"OrderID":   order.ID,
		"OrderType": string(order.Type),
		"Status":    string(order.Status),
		"Time":      order.Time.Format(time.RFC3339),
	}

	return h.service.SendWithTemplate("trade_executed", data)
}

// NotifyPositionOpened sends a notification when a position is opened
func (h *NotificationHandler) NotifyPositionOpened(ctx context.Context, position *models.Position) error {
	data := map[string]interface{}{
		"Symbol":     position.Symbol,
		"Side":       string(position.Side),
		"Quantity":   position.Quantity,
		"EntryPrice": position.EntryPrice,
		"OpenTime":   position.OpenTime.Format(time.RFC3339),
		"StopLoss":   position.StopLoss,
		"TakeProfit": position.TakeProfit,
	}

	return h.service.SendWithTemplate("position_opened", data)
}

// NotifyPositionClosed sends a notification when a position is closed
func (h *NotificationHandler) NotifyPositionClosed(ctx context.Context, position *models.ClosedPosition) error {
	data := map[string]interface{}{
		"Symbol":               position.Symbol,
		"Quantity":             position.Amount,
		"EntryPrice":           position.EntryPrice,
		"ExitPrice":            position.ExitPrice,
		"ProfitLoss":           position.ProfitLoss,
		"ProfitLossPercentage": position.ProfitLossPercentage,
		"OpenTime":             position.OpenTime.Format(time.RFC3339),
		"CloseTime":            position.CloseTime.Format(time.RFC3339),
		"Reason":               position.ExitReason,
	}

	return h.service.SendWithTemplate("position_closed", data)
}

// NotifyNewCoinDetected sends a notification when a new coin is detected
func (h *NotificationHandler) NotifyNewCoinDetected(ctx context.Context, symbol string, exchange string, price float64) error {
	data := map[string]interface{}{
		"Symbol":   symbol,
		"Exchange": exchange,
		"Price":    price,
		"Time":     time.Now().Format(time.RFC3339),
	}

	return h.service.SendWithTemplate("new_coin_detected", data)
}

// NotifyRiskAlert sends a notification for a risk alert
func (h *NotificationHandler) NotifyRiskAlert(ctx context.Context, alertType string, message string) error {
	data := map[string]interface{}{
		"AlertType": alertType,
		"Message":   message,
		"Time":      time.Now().Format(time.RFC3339),
	}

	return h.service.SendWithTemplate("risk_alert", data)
}

// NotifySystemStatus sends a notification for system status
func (h *NotificationHandler) NotifySystemStatus(ctx context.Context, status string, message string) error {
	data := map[string]interface{}{
		"Status":  status,
		"Message": message,
		"Time":    time.Now().Format(time.RFC3339),
	}

	return h.service.SendWithTemplate("system_status", data)
}

// NotifyError sends a notification for an error
func (h *NotificationHandler) NotifyError(ctx context.Context, service string, err error, stackTrace string) error {
	data := map[string]interface{}{
		"Service":    service,
		"Message":    err.Error(),
		"StackTrace": stackTrace,
		"Time":       time.Now().Format(time.RFC3339),
	}

	return h.service.SendWithTemplate("error_alert", data)
}

// NotifyBacktestResult sends a notification for a backtest result
func (h *NotificationHandler) NotifyBacktestResult(ctx context.Context, strategyName string, symbols []string, startTime, endTime time.Time, initialCapital, finalCapital float64, totalReturn, sharpeRatio, maxDrawdown, winRate float64, totalTrades int) error {
	// Create a custom notification since we don't have a template for this
	title := fmt.Sprintf("Backtest Result: %s", strategyName)
	message := fmt.Sprintf(
		"Strategy: %s\nSymbols: %v\nPeriod: %s to %s\nInitial Capital: $%.2f\nFinal Capital: $%.2f\nTotal Return: %.2f%%\nSharpe Ratio: %.2f\nMax Drawdown: %.2f%%\nWin Rate: %.2f%%\nTotal Trades: %d",
		strategyName,
		symbols,
		startTime.Format("2006-01-02"),
		endTime.Format("2006-01-02"),
		initialCapital,
		finalCapital,
		totalReturn,
		sharpeRatio,
		maxDrawdown,
		winRate,
		totalTrades,
	)

	notification := NewNotification(title, message, LevelInfo)
	notification.Source = "BacktestService"

	return h.service.Send(notification)
}
