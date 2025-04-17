package service

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/notification"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/persistence/csv"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model/status"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/rs/zerolog"
)

// TradingService is the main orchestration service for the trading bot
type TradingService struct {
	// Core components
	logger           *zerolog.Logger
	tradeExecutor    port.TradeExecutor
	tradeHistory     port.TradeHistoryRepository
	csvWriter        *csv.TradeHistoryWriter
	telegramNotifier *notification.TelegramNotifier

	// Service state
	ctx         context.Context
	cancel      context.CancelFunc
	wg          sync.WaitGroup
	status      status.Status
	statusMutex sync.RWMutex

	// Configuration
	shutdownTimeout time.Duration
	recoveryEnabled bool
}

// TradingServiceConfig contains configuration for the trading service
type TradingServiceConfig struct {
	ShutdownTimeout time.Duration
	RecoveryEnabled bool
}

// NewTradingService creates a new trading service
func NewTradingService(
	logger *zerolog.Logger,
	tradeExecutor port.TradeExecutor,
	tradeHistory port.TradeHistoryRepository,
	csvWriter *csv.TradeHistoryWriter,
	telegramNotifier *notification.TelegramNotifier,
	config TradingServiceConfig,
) *TradingService {
	ctx, cancel := context.WithCancel(context.Background())

	return &TradingService{
		logger:           logger,
		tradeExecutor:    tradeExecutor,
		tradeHistory:     tradeHistory,
		csvWriter:        csvWriter,
		telegramNotifier: telegramNotifier,
		ctx:              ctx,
		cancel:           cancel,
		status:           status.StatusStopped,
		shutdownTimeout:  config.ShutdownTimeout,
		recoveryEnabled:  config.RecoveryEnabled,
	}
}

// Start starts the trading service
func (s *TradingService) Start() error {
	s.statusMutex.Lock()
	if s.status == status.StatusRunning {
		s.statusMutex.Unlock()
		return fmt.Errorf("service is already running")
	}

	oldStatus := s.status
	s.status = status.StatusRunning
	s.statusMutex.Unlock()

	s.logger.Info().Msg("Starting trading service")

	// Notify status change
	if s.telegramNotifier != nil {
		if err := s.telegramNotifier.NotifySystemStatusChange(s.ctx, oldStatus, status.StatusRunning, "Trading service started"); err != nil {
			s.logger.Warn().Err(err).Msg("Failed to send start notification")
		}
	}

	// Start signal handler for graceful shutdown
	s.startSignalHandler()

	// Start worker goroutines
	s.startWorkers()

	return nil
}

// Stop stops the trading service
func (s *TradingService) Stop() error {
	s.statusMutex.Lock()
	if s.status == status.StatusStopped {
		s.statusMutex.Unlock()
		return fmt.Errorf("service is already stopped")
	}

	oldStatus := s.status
	s.status = status.StatusStopped
	s.statusMutex.Unlock()

	s.logger.Info().Msg("Stopping trading service")

	// Cancel context to signal all workers to stop
	s.cancel()

	// Wait for all workers to finish with timeout
	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		s.logger.Info().Msg("All workers stopped gracefully")
	case <-time.After(s.shutdownTimeout):
		s.logger.Warn().Msg("Shutdown timed out, some workers may still be running")
	}

	// Close CSV writer
	if s.csvWriter != nil {
		if err := s.csvWriter.Close(); err != nil {
			s.logger.Error().Err(err).Msg("Error closing CSV writer")
		}
	}

	// Notify status change
	if s.telegramNotifier != nil {
		if err := s.telegramNotifier.NotifySystemStatusChange(s.ctx, oldStatus, status.StatusStopped, "Trading service stopped"); err != nil {
			s.logger.Warn().Err(err).Msg("Failed to send stop notification")
		}
	}

	return nil
}

// GetStatus returns the current status of the trading service
func (s *TradingService) GetStatus() status.Status {
	s.statusMutex.RLock()
	defer s.statusMutex.RUnlock()
	return s.status
}

// ExecuteOrder executes a trade order and records it
func (s *TradingService) ExecuteOrder(ctx context.Context, request *model.OrderRequest) (*model.OrderResponse, error) {
	// Check if service is running
	if s.GetStatus() != status.StatusRunning {
		return nil, fmt.Errorf("trading service is not running")
	}

	// Execute the order with rate limiting and error handling
	response, err := s.tradeExecutor.ExecuteOrder(ctx, request)
	if err != nil {
		s.logger.Error().
			Err(err).
			Str("symbol", request.Symbol).
			Str("side", string(request.Side)).
			Str("type", string(request.Type)).
			Float64("quantity", request.Quantity).
			Float64("price", request.Price).
			Msg("Failed to execute order")

		// Notify about the error
		if s.telegramNotifier != nil {
			errMsg := fmt.Sprintf("Failed to execute %s %s order for %s: %v",
				request.Side, request.Type, request.Symbol, err)
			_ = s.telegramNotifier.NotifyAlert(ctx, "error", "Order Execution Failed", errMsg, "trade_executor")
		}

		return nil, err
	}

	// Record the trade
	tradeRecord := &model.TradeRecord{
		UserID:        request.UserID,
		Symbol:        request.Symbol,
		Side:          request.Side,
		Type:          request.Type,
		Quantity:      request.Quantity,
		Price:         request.Price,
		Amount:        request.Quantity * request.Price,
		OrderID:       response.Order.OrderID,
		ExecutionTime: time.Now(),
		// Strategy:      request.Metadata["strategy"].(string), // Removed: OrderRequest has no Metadata field
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Save to database
	if err := s.tradeHistory.SaveTradeRecord(ctx, tradeRecord); err != nil {
		s.logger.Error().
			Err(err).
			Str("orderID", response.Order.OrderID).
			Msg("Failed to save trade record to database")
	}

	// Save to CSV
	if s.csvWriter != nil {
		if err := s.csvWriter.WriteTradeRecord(ctx, tradeRecord); err != nil {
			s.logger.Error().
				Err(err).
				Str("orderID", response.Order.OrderID).
				Msg("Failed to save trade record to CSV")
		}
	}

	// Send notification
	if s.telegramNotifier != nil {
		_ = s.telegramNotifier.NotifyTrade(ctx,
			request.Symbol,
			string(request.Side),
			string(request.Type),
			request.Quantity,
			request.Price,
			response.Order.OrderID)
	}

	s.logger.Info().
		Str("orderID", response.Order.OrderID).
		Str("symbol", request.Symbol).
		Str("side", string(request.Side)).
		Str("type", string(request.Type)).
		Float64("quantity", request.Quantity).
		Float64("price", request.Price).
		Msg("Order executed and recorded successfully")

	return response, nil
}

// LogDetection logs a market event detection
func (s *TradingService) LogDetection(ctx context.Context, detectionType, symbol string, value, threshold float64, description string, metadata map[string]interface{}) (*model.DetectionLog, error) {
	// Create detection log
	log := &model.DetectionLog{
		Type:        detectionType,
		Symbol:      symbol,
		Value:       value,
		Threshold:   threshold,
		Description: description,
		Metadata:    metadata,
		DetectedAt:  time.Now(),
		Processed:   false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Save to database
	if err := s.tradeHistory.SaveDetectionLog(ctx, log); err != nil {
		s.logger.Error().
			Err(err).
			Str("type", detectionType).
			Str("symbol", symbol).
			Msg("Failed to save detection log to database")
		return nil, err
	}

	// Save to CSV
	if s.csvWriter != nil {
		if err := s.csvWriter.WriteDetectionLog(ctx, log); err != nil {
			s.logger.Error().
				Err(err).
				Str("id", log.ID).
				Msg("Failed to save detection log to CSV")
		}
	}

	// Send notification for significant detections
	if s.telegramNotifier != nil && isSignificantDetection(detectionType, value, threshold) {
		message := fmt.Sprintf("*%s Detection*\nSymbol: `%s`\nValue: `%.8f`\nThreshold: `%.8f`\nDescription: %s",
			detectionType, symbol, value, threshold, description)
		_ = s.telegramNotifier.NotifyAlert(ctx, "info", fmt.Sprintf("%s Detection", detectionType), message, "market_monitor")
	}

	s.logger.Info().
		Str("id", log.ID).
		Str("type", detectionType).
		Str("symbol", symbol).
		Float64("value", value).
		Float64("threshold", threshold).
		Msg("Detection logged successfully")

	return log, nil
}

// MarkDetectionProcessed marks a detection log as processed
func (s *TradingService) MarkDetectionProcessed(ctx context.Context, id, result string) error {
	if err := s.tradeHistory.MarkDetectionLogProcessed(ctx, id, result); err != nil {
		s.logger.Error().
			Err(err).
			Str("id", id).
			Msg("Failed to mark detection log as processed")
		return err
	}

	s.logger.Info().
		Str("id", id).
		Str("result", result).
		Msg("Detection log marked as processed")

	return nil
}

// startSignalHandler starts a goroutine to handle OS signals
func (s *TradingService) startSignalHandler() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		defer signal.Stop(signals)

		select {
		case sig := <-signals:
			s.logger.Info().
				Str("signal", sig.String()).
				Msg("Received signal, initiating shutdown")
			_ = s.Stop()
		case <-s.ctx.Done():
			// Context was cancelled, just exit
			return
		}
	}()
}

// startWorkers starts the worker goroutines
func (s *TradingService) startWorkers() {
	// Start detection processor
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		defer s.recoverPanic("detection_processor")

		s.runDetectionProcessor()
	}()

	// Add more workers as needed
}

// runDetectionProcessor processes unprocessed detection logs
func (s *TradingService) runDetectionProcessor() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Process unprocessed detection logs
			ctx, cancel := context.WithTimeout(s.ctx, 10*time.Second)
			logs, err := s.tradeHistory.GetUnprocessedDetectionLogs(ctx, 10)
			if err != nil {
				s.logger.Error().Err(err).Msg("Failed to get unprocessed detection logs")
				cancel()
				continue
			}

			for _, log := range logs {
				// Process the detection log
				// This would typically involve evaluating trading rules
				// and potentially executing trades

				// For now, just mark it as processed
				result := "processed without action"
				if err := s.tradeHistory.MarkDetectionLogProcessed(ctx, log.ID, result); err != nil {
					s.logger.Error().
						Err(err).
						Str("id", log.ID).
						Msg("Failed to mark detection log as processed")
				}
			}

			cancel()
		case <-s.ctx.Done():
			return
		}
	}
}

// recoverPanic recovers from panics in worker goroutines
func (s *TradingService) recoverPanic(workerName string) {
	if !s.recoveryEnabled {
		return
	}

	if r := recover(); r != nil {
		s.logger.Error().
			Interface("panic", r).
			Str("worker", workerName).
			Msg("Recovered from panic in worker")

		// Notify about the panic
		if s.telegramNotifier != nil {
			ctx := context.Background()
			message := fmt.Sprintf("Panic in %s worker: %v", workerName, r)
			_ = s.telegramNotifier.NotifyAlert(ctx, "critical", "Worker Panic", message, "recovery")
		}
	}
}

// isSignificantDetection determines if a detection is significant enough to notify
func isSignificantDetection(detectionType string, value, threshold float64) bool {
	// Implement logic to determine if a detection is significant
	// For example, if the value exceeds the threshold by a certain percentage

	switch detectionType {
	case "price_spike":
		return value > threshold*1.5
	case "volume_spike":
		return value > threshold*2.0
	case "breakout":
		return true // Always notify for breakouts
	default:
		return false
	}
}
