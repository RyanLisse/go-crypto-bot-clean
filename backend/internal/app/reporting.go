package app

import (
	"context"
	"time"

	"go-crypto-bot-clean/backend/internal/api/handlers"
	"go-crypto-bot-clean/backend/internal/api/websocket"
	"go-crypto-bot-clean/backend/internal/repository/sqlite"
	"go-crypto-bot-clean/backend/internal/services/gemini"
	"go-crypto-bot-clean/backend/internal/services/reporting"
	"go.uber.org/zap"
)

// SetupReportingSystem sets up the performance reporting system
func (a *BotApp) SetupReportingSystem(ctx context.Context) error {
	// Create Gemini client
	geminiClient := gemini.NewGeminiClient(a.config.Gemini.APIKey)

	// Create report repository
	reportRepo := sqlite.NewReportRepository(a.db, a.logger)
	if err := reportRepo.Initialize(ctx); err != nil {
		return err
	}

	// Create metrics collector
	metricsCollector := reporting.NewMetricsCollector(
		a.tradeAnalyticsRepo,
		a.balanceHistoryRepo,
		a.logger,
	)

	// Create report generator
	reportGenerator := reporting.NewReportGenerator(
		geminiClient,
		reportRepo,
		time.Duration(a.config.Reporting.Interval)*time.Minute,
		a.logger,
	)

	// Create report handler
	reportHandler := handlers.NewReportHandler(reportGenerator, a.logger)
	a.reportHandler = reportHandler

	// Create report service for WebSocket
	reportService := websocket.NewReportService(a.wsHub, a.logger)

	// Start metrics collection and processing
	go reportGenerator.StartCollection(ctx, metricsCollector)
	go reportGenerator.ProcessMetrics(ctx)

	// Set up a goroutine to broadcast new reports via WebSocket
	go func() {
		ticker := time.NewTicker(time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				// Get the latest report
				report, err := reportGenerator.GetLatestReport(ctx)
				if err != nil {
					a.logger.Error("Failed to get latest report for WebSocket broadcast", zap.Error(err))
					continue
				}

				if report != nil {
					// Broadcast the report
					reportService.BroadcastReport(report)
				}
			}
		}
	}()

	return nil
}
