package csv

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/rs/zerolog"
)

// TradeHistoryWriter implements CSV writing for trade history
type TradeHistoryWriter struct {
	logger          *zerolog.Logger
	baseDir         string
	tradeFile       *os.File
	detectionFile   *os.File
	tradeWriter     *csv.Writer
	detectionWriter *csv.Writer
	mutex           sync.Mutex
	enabled         bool
}

// TradeHistoryWriterConfig contains configuration for the CSV writer
type TradeHistoryWriterConfig struct {
	Enabled           bool
	BaseDirectory     string
	TradeFilename     string
	DetectionFilename string
	FlushInterval     time.Duration
}

// NewTradeHistoryWriter creates a new CSV trade history writer
func NewTradeHistoryWriter(config TradeHistoryWriterConfig, logger *zerolog.Logger) (*TradeHistoryWriter, error) {
	writer := &TradeHistoryWriter{
		logger:  logger,
		baseDir: config.BaseDirectory,
		enabled: config.Enabled,
	}

	if !config.Enabled {
		logger.Info().Msg("CSV trade history writer is disabled")
		return writer, nil
	}

	// Create base directory if it doesn't exist
	if err := os.MkdirAll(config.BaseDirectory, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	// Initialize trade file
	tradeFilePath := filepath.Join(config.BaseDirectory, config.TradeFilename)
	tradeFileExists := fileExists(tradeFilePath)

	tradeFile, err := os.OpenFile(tradeFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open trade file: %w", err)
	}
	writer.tradeFile = tradeFile
	writer.tradeWriter = csv.NewWriter(tradeFile)

	// Write header if file is new
	if !tradeFileExists {
		if err := writer.writeTradeHeader(); err != nil {
			return nil, fmt.Errorf("failed to write trade header: %w", err)
		}
	}

	// Initialize detection file
	detectionFilePath := filepath.Join(config.BaseDirectory, config.DetectionFilename)
	detectionFileExists := fileExists(detectionFilePath)

	detectionFile, err := os.OpenFile(detectionFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open detection file: %w", err)
	}
	writer.detectionFile = detectionFile
	writer.detectionWriter = csv.NewWriter(detectionFile)

	// Write header if file is new
	if !detectionFileExists {
		if err := writer.writeDetectionHeader(); err != nil {
			return nil, fmt.Errorf("failed to write detection header: %w", err)
		}
	}

	// Start periodic flush if interval is set
	if config.FlushInterval > 0 {
		go writer.periodicFlush(config.FlushInterval)
	}

	logger.Info().
		Str("tradeFile", tradeFilePath).
		Str("detectionFile", detectionFilePath).
		Msg("CSV trade history writer initialized")

	return writer, nil
}

// WriteTradeRecord writes a trade record to the CSV file
func (w *TradeHistoryWriter) WriteTradeRecord(ctx context.Context, record *model.TradeRecord) error {
	if !w.enabled {
		return nil
	}

	w.mutex.Lock()
	defer w.mutex.Unlock()

	// Convert tags to JSON string
	tagsJSON, err := json.Marshal(record.Tags)
	if err != nil {
		return fmt.Errorf("failed to marshal tags: %w", err)
	}

	// Convert metadata to JSON string
	metadataJSON, err := json.Marshal(record.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	// Prepare row
	row := []string{
		record.ID,
		record.UserID,
		record.Symbol,
		string(record.Side),
		string(record.Type),
		strconv.FormatFloat(record.Quantity, 'f', -1, 64),
		strconv.FormatFloat(record.Price, 'f', -1, 64),
		strconv.FormatFloat(record.Amount, 'f', -1, 64),
		strconv.FormatFloat(record.Fee, 'f', -1, 64),
		record.FeeCurrency,
		record.OrderID,
		record.TradeID,
		record.ExecutionTime.Format(time.RFC3339),
		record.Strategy,
		record.Notes,
		string(tagsJSON),
		string(metadataJSON),
		record.CreatedAt.Format(time.RFC3339),
		record.UpdatedAt.Format(time.RFC3339),
	}

	// Write row
	if err := w.tradeWriter.Write(row); err != nil {
		w.logger.Error().Err(err).Str("id", record.ID).Msg("Failed to write trade record to CSV")
		return fmt.Errorf("failed to write trade record to CSV: %w", err)
	}

	// Flush to ensure data is written
	w.tradeWriter.Flush()
	if err := w.tradeWriter.Error(); err != nil {
		w.logger.Error().Err(err).Msg("Error flushing trade writer")
		return fmt.Errorf("error flushing trade writer: %w", err)
	}

	w.logger.Debug().Str("id", record.ID).Msg("Trade record written to CSV")
	return nil
}

// WriteDetectionLog writes a detection log to the CSV file
func (w *TradeHistoryWriter) WriteDetectionLog(ctx context.Context, log *model.DetectionLog) error {
	if !w.enabled {
		return nil
	}

	w.mutex.Lock()
	defer w.mutex.Unlock()

	// Convert metadata to JSON string
	metadataJSON, err := json.Marshal(log.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	// Format processed_at
	processedAtStr := ""
	if log.ProcessedAt != nil {
		processedAtStr = log.ProcessedAt.Format(time.RFC3339)
	}

	// Prepare row
	row := []string{
		log.ID,
		log.Type,
		log.Symbol,
		strconv.FormatFloat(log.Value, 'f', -1, 64),
		strconv.FormatFloat(log.Threshold, 'f', -1, 64),
		log.Description,
		string(metadataJSON),
		log.DetectedAt.Format(time.RFC3339),
		processedAtStr,
		strconv.FormatBool(log.Processed),
		log.Result,
		log.CreatedAt.Format(time.RFC3339),
		log.UpdatedAt.Format(time.RFC3339),
	}

	// Write row
	if err := w.detectionWriter.Write(row); err != nil {
		w.logger.Error().Err(err).Str("id", log.ID).Msg("Failed to write detection log to CSV")
		return fmt.Errorf("failed to write detection log to CSV: %w", err)
	}

	// Flush to ensure data is written
	w.detectionWriter.Flush()
	if err := w.detectionWriter.Error(); err != nil {
		w.logger.Error().Err(err).Msg("Error flushing detection writer")
		return fmt.Errorf("error flushing detection writer: %w", err)
	}

	w.logger.Debug().Str("id", log.ID).Msg("Detection log written to CSV")
	return nil
}

// Close closes the CSV files
func (w *TradeHistoryWriter) Close() error {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if !w.enabled {
		return nil
	}

	// Flush writers
	if w.tradeWriter != nil {
		w.tradeWriter.Flush()
	}
	if w.detectionWriter != nil {
		w.detectionWriter.Flush()
	}

	// Close files
	var tradeErr, detectionErr error
	if w.tradeFile != nil {
		tradeErr = w.tradeFile.Close()
		w.tradeFile = nil
		w.tradeWriter = nil
	}
	if w.detectionFile != nil {
		detectionErr = w.detectionFile.Close()
		w.detectionFile = nil
		w.detectionWriter = nil
	}

	// Return first error
	if tradeErr != nil {
		return fmt.Errorf("failed to close trade file: %w", tradeErr)
	}
	if detectionErr != nil {
		return fmt.Errorf("failed to close detection file: %w", detectionErr)
	}

	w.logger.Info().Msg("CSV trade history writer closed")
	return nil
}

// Flush flushes the CSV writers
func (w *TradeHistoryWriter) Flush() error {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if !w.enabled {
		return nil
	}

	// Flush writers
	if w.tradeWriter != nil {
		w.tradeWriter.Flush()
		if err := w.tradeWriter.Error(); err != nil {
			return fmt.Errorf("error flushing trade writer: %w", err)
		}
	}
	if w.detectionWriter != nil {
		w.detectionWriter.Flush()
		if err := w.detectionWriter.Error(); err != nil {
			return fmt.Errorf("error flushing detection writer: %w", err)
		}
	}

	return nil
}

// writeTradeHeader writes the header row to the trade CSV file
func (w *TradeHistoryWriter) writeTradeHeader() error {
	header := []string{
		"id",
		"user_id",
		"symbol",
		"side",
		"type",
		"quantity",
		"price",
		"amount",
		"fee",
		"fee_currency",
		"order_id",
		"trade_id",
		"execution_time",
		"strategy",
		"notes",
		"tags",
		"metadata",
		"created_at",
		"updated_at",
	}

	if err := w.tradeWriter.Write(header); err != nil {
		return fmt.Errorf("failed to write trade header: %w", err)
	}
	w.tradeWriter.Flush()
	return w.tradeWriter.Error()
}

// writeDetectionHeader writes the header row to the detection CSV file
func (w *TradeHistoryWriter) writeDetectionHeader() error {
	header := []string{
		"id",
		"type",
		"symbol",
		"value",
		"threshold",
		"description",
		"metadata",
		"detected_at",
		"processed_at",
		"processed",
		"result",
		"created_at",
		"updated_at",
	}

	if err := w.detectionWriter.Write(header); err != nil {
		return fmt.Errorf("failed to write detection header: %w", err)
	}
	w.detectionWriter.Flush()
	return w.detectionWriter.Error()
}

// periodicFlush periodically flushes the CSV writers
func (w *TradeHistoryWriter) periodicFlush(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		if err := w.Flush(); err != nil {
			w.logger.Error().Err(err).Msg("Error during periodic flush")
		}
	}
}

// fileExists checks if a file exists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
