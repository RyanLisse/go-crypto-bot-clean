package gorm

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// TradeRecordEntity represents a trade record in the database
type TradeRecordEntity struct {
	ID            string    `gorm:"primaryKey"`
	UserID        string    `gorm:"index"`
	Symbol        string    `gorm:"index"`
	Side          string    `gorm:"index"`
	Type          string
	Quantity      float64
	Price         float64
	Amount        float64
	Fee           float64
	FeeCurrency   string
	OrderID       string    `gorm:"index"`
	TradeID       string    `gorm:"index"`
	ExecutionTime time.Time `gorm:"index"`
	Strategy      string    `gorm:"index"`
	Notes         string
	Tags          string // JSON array
	Metadata      string // JSON object
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// TableName returns the table name for the trade record entity
func (TradeRecordEntity) TableName() string {
	return "trade_records"
}

// DetectionLogEntity represents a detection log in the database
type DetectionLogEntity struct {
	ID          string    `gorm:"primaryKey"`
	Type        string    `gorm:"index"`
	Symbol      string    `gorm:"index"`
	Value       float64
	Threshold   float64
	Description string
	Metadata    string // JSON object
	DetectedAt  time.Time `gorm:"index"`
	ProcessedAt *time.Time
	Processed   bool      `gorm:"index"`
	Result      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// TableName returns the table name for the detection log entity
func (DetectionLogEntity) TableName() string {
	return "detection_logs"
}

// TradeHistoryRepository implements the TradeHistoryRepository interface using GORM
type TradeHistoryRepository struct {
	db     *gorm.DB
	logger *zerolog.Logger
}

// NewTradeHistoryRepository creates a new trade history repository
func NewTradeHistoryRepository(db *gorm.DB, logger *zerolog.Logger) *TradeHistoryRepository {
	return &TradeHistoryRepository{
		db:     db,
		logger: logger,
	}
}

// SaveTradeRecord saves a trade record
func (r *TradeHistoryRepository) SaveTradeRecord(ctx context.Context, record *model.TradeRecord) error {
	// Generate ID if not provided
	if record.ID == "" {
		record.ID = uuid.New().String()
	}

	// Set timestamps
	now := time.Now()
	if record.CreatedAt.IsZero() {
		record.CreatedAt = now
	}
	record.UpdatedAt = now

	// Convert to entity
	entity, err := r.tradeRecordToEntity(record)
	if err != nil {
		return fmt.Errorf("failed to convert trade record to entity: %w", err)
	}

	// Save to database
	if err := r.db.WithContext(ctx).Create(entity).Error; err != nil {
		r.logger.Error().Err(err).Str("id", record.ID).Msg("Failed to save trade record")
		return fmt.Errorf("failed to save trade record: %w", err)
	}

	r.logger.Debug().Str("id", record.ID).Msg("Trade record saved successfully")
	return nil
}

// GetTradeRecords retrieves trade records with filtering
func (r *TradeHistoryRepository) GetTradeRecords(ctx context.Context, filter port.TradeHistoryFilter) ([]*model.TradeRecord, error) {
	var entities []TradeRecordEntity
	query := r.db.WithContext(ctx)

	// Apply filters
	if filter.UserID != "" {
		query = query.Where("user_id = ?", filter.UserID)
	}
	if filter.Symbol != "" {
		query = query.Where("symbol = ?", filter.Symbol)
	}
	if filter.Side != "" {
		query = query.Where("side = ?", string(filter.Side))
	}
	if filter.Strategy != "" {
		query = query.Where("strategy = ?", filter.Strategy)
	}
	if !filter.StartTime.IsZero() {
		query = query.Where("execution_time >= ?", filter.StartTime)
	}
	if !filter.EndTime.IsZero() {
		query = query.Where("execution_time <= ?", filter.EndTime)
	}
	if len(filter.Tags) > 0 {
		// For each tag, check if it's in the JSON array
		for _, tag := range filter.Tags {
			query = query.Where("tags LIKE ?", "%"+tag+"%")
		}
	}

	// Apply pagination
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}
	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}

	// Order by execution time descending
	query = query.Order("execution_time DESC")

	// Execute query
	if err := query.Find(&entities).Error; err != nil {
		r.logger.Error().Err(err).Msg("Failed to retrieve trade records")
		return nil, fmt.Errorf("failed to retrieve trade records: %w", err)
	}

	// Convert to domain models
	records := make([]*model.TradeRecord, len(entities))
	for i, entity := range entities {
		record, err := r.entityToTradeRecord(&entity)
		if err != nil {
			r.logger.Error().Err(err).Str("id", entity.ID).Msg("Failed to convert entity to trade record")
			return nil, fmt.Errorf("failed to convert entity to trade record: %w", err)
		}
		records[i] = record
	}

	return records, nil
}

// GetTradeRecordByID retrieves a trade record by ID
func (r *TradeHistoryRepository) GetTradeRecordByID(ctx context.Context, id string) (*model.TradeRecord, error) {
	var entity TradeRecordEntity
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&entity).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.logger.Error().Err(err).Str("id", id).Msg("Failed to retrieve trade record")
		return nil, fmt.Errorf("failed to retrieve trade record: %w", err)
	}

	record, err := r.entityToTradeRecord(&entity)
	if err != nil {
		r.logger.Error().Err(err).Str("id", id).Msg("Failed to convert entity to trade record")
		return nil, fmt.Errorf("failed to convert entity to trade record: %w", err)
	}

	return record, nil
}

// GetTradeRecordsByOrderID retrieves trade records by order ID
func (r *TradeHistoryRepository) GetTradeRecordsByOrderID(ctx context.Context, orderID string) ([]*model.TradeRecord, error) {
	var entities []TradeRecordEntity
	if err := r.db.WithContext(ctx).Where("order_id = ?", orderID).Find(&entities).Error; err != nil {
		r.logger.Error().Err(err).Str("orderID", orderID).Msg("Failed to retrieve trade records by order ID")
		return nil, fmt.Errorf("failed to retrieve trade records by order ID: %w", err)
	}

	records := make([]*model.TradeRecord, len(entities))
	for i, entity := range entities {
		record, err := r.entityToTradeRecord(&entity)
		if err != nil {
			r.logger.Error().Err(err).Str("id", entity.ID).Msg("Failed to convert entity to trade record")
			return nil, fmt.Errorf("failed to convert entity to trade record: %w", err)
		}
		records[i] = record
	}

	return records, nil
}

// SaveDetectionLog saves a detection log
func (r *TradeHistoryRepository) SaveDetectionLog(ctx context.Context, log *model.DetectionLog) error {
	// Generate ID if not provided
	if log.ID == "" {
		log.ID = uuid.New().String()
	}

	// Set timestamps
	now := time.Now()
	if log.CreatedAt.IsZero() {
		log.CreatedAt = now
	}
	log.UpdatedAt = now

	// Convert to entity
	entity, err := r.detectionLogToEntity(log)
	if err != nil {
		return fmt.Errorf("failed to convert detection log to entity: %w", err)
	}

	// Save to database
	if err := r.db.WithContext(ctx).Create(entity).Error; err != nil {
		r.logger.Error().Err(err).Str("id", log.ID).Msg("Failed to save detection log")
		return fmt.Errorf("failed to save detection log: %w", err)
	}

	r.logger.Debug().Str("id", log.ID).Msg("Detection log saved successfully")
	return nil
}

// GetDetectionLogs retrieves detection logs with filtering
func (r *TradeHistoryRepository) GetDetectionLogs(ctx context.Context, filter port.DetectionLogFilter) ([]*model.DetectionLog, error) {
	var entities []DetectionLogEntity
	query := r.db.WithContext(ctx)

	// Apply filters
	if filter.Type != "" {
		query = query.Where("type = ?", filter.Type)
	}
	if filter.Symbol != "" {
		query = query.Where("symbol = ?", filter.Symbol)
	}
	if filter.Processed != nil {
		query = query.Where("processed = ?", *filter.Processed)
	}
	if !filter.StartTime.IsZero() {
		query = query.Where("detected_at >= ?", filter.StartTime)
	}
	if !filter.EndTime.IsZero() {
		query = query.Where("detected_at <= ?", filter.EndTime)
	}

	// Apply pagination
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}
	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}

	// Order by detected time descending
	query = query.Order("detected_at DESC")

	// Execute query
	if err := query.Find(&entities).Error; err != nil {
		r.logger.Error().Err(err).Msg("Failed to retrieve detection logs")
		return nil, fmt.Errorf("failed to retrieve detection logs: %w", err)
	}

	// Convert to domain models
	logs := make([]*model.DetectionLog, len(entities))
	for i, entity := range entities {
		log, err := r.entityToDetectionLog(&entity)
		if err != nil {
			r.logger.Error().Err(err).Str("id", entity.ID).Msg("Failed to convert entity to detection log")
			return nil, fmt.Errorf("failed to convert entity to detection log: %w", err)
		}
		logs[i] = log
	}

	return logs, nil
}

// MarkDetectionLogProcessed marks a detection log as processed
func (r *TradeHistoryRepository) MarkDetectionLogProcessed(ctx context.Context, id string, result string) error {
	now := time.Now()
	if err := r.db.WithContext(ctx).Model(&DetectionLogEntity{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"processed":    true,
			"processed_at": now,
			"result":       result,
			"updated_at":   now,
		}).Error; err != nil {
		r.logger.Error().Err(err).Str("id", id).Msg("Failed to mark detection log as processed")
		return fmt.Errorf("failed to mark detection log as processed: %w", err)
	}

	r.logger.Debug().Str("id", id).Msg("Detection log marked as processed")
	return nil
}

// GetUnprocessedDetectionLogs retrieves unprocessed detection logs
func (r *TradeHistoryRepository) GetUnprocessedDetectionLogs(ctx context.Context, limit int) ([]*model.DetectionLog, error) {
	var entities []DetectionLogEntity
	query := r.db.WithContext(ctx).
		Where("processed = ?", false).
		Order("detected_at ASC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&entities).Error; err != nil {
		r.logger.Error().Err(err).Msg("Failed to retrieve unprocessed detection logs")
		return nil, fmt.Errorf("failed to retrieve unprocessed detection logs: %w", err)
	}

	logs := make([]*model.DetectionLog, len(entities))
	for i, entity := range entities {
		log, err := r.entityToDetectionLog(&entity)
		if err != nil {
			r.logger.Error().Err(err).Str("id", entity.ID).Msg("Failed to convert entity to detection log")
			return nil, fmt.Errorf("failed to convert entity to detection log: %w", err)
		}
		logs[i] = log
	}

	return logs, nil
}

// Helper methods for entity conversion

func (r *TradeHistoryRepository) tradeRecordToEntity(record *model.TradeRecord) (*TradeRecordEntity, error) {
	// Convert tags to JSON
	tagsJSON, err := json.Marshal(record.Tags)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal tags: %w", err)
	}

	// Convert metadata to JSON
	var metadataJSON string
	if record.Metadata != nil {
		metadataBytes, err := json.Marshal(record.Metadata)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal metadata: %w", err)
		}
		metadataJSON = string(metadataBytes)
	}

	return &TradeRecordEntity{
		ID:            record.ID,
		UserID:        record.UserID,
		Symbol:        record.Symbol,
		Side:          string(record.Side),
		Type:          string(record.Type),
		Quantity:      record.Quantity,
		Price:         record.Price,
		Amount:        record.Amount,
		Fee:           record.Fee,
		FeeCurrency:   record.FeeCurrency,
		OrderID:       record.OrderID,
		TradeID:       record.TradeID,
		ExecutionTime: record.ExecutionTime,
		Strategy:      record.Strategy,
		Notes:         record.Notes,
		Tags:          string(tagsJSON),
		Metadata:      metadataJSON,
		CreatedAt:     record.CreatedAt,
		UpdatedAt:     record.UpdatedAt,
	}, nil
}

func (r *TradeHistoryRepository) entityToTradeRecord(entity *TradeRecordEntity) (*model.TradeRecord, error) {
	// Parse tags from JSON
	var tags []string
	if entity.Tags != "" {
		if err := json.Unmarshal([]byte(entity.Tags), &tags); err != nil {
			return nil, fmt.Errorf("failed to unmarshal tags: %w", err)
		}
	}

	// Parse metadata from JSON
	var metadata map[string]interface{}
	if entity.Metadata != "" {
		if err := json.Unmarshal([]byte(entity.Metadata), &metadata); err != nil {
			return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
		}
	}

	return &model.TradeRecord{
		ID:            entity.ID,
		UserID:        entity.UserID,
		Symbol:        entity.Symbol,
		Side:          model.OrderSide(entity.Side),
		Type:          model.OrderType(entity.Type),
		Quantity:      entity.Quantity,
		Price:         entity.Price,
		Amount:        entity.Amount,
		Fee:           entity.Fee,
		FeeCurrency:   entity.FeeCurrency,
		OrderID:       entity.OrderID,
		TradeID:       entity.TradeID,
		ExecutionTime: entity.ExecutionTime,
		Strategy:      entity.Strategy,
		Notes:         entity.Notes,
		Tags:          tags,
		Metadata:      metadata,
		CreatedAt:     entity.CreatedAt,
		UpdatedAt:     entity.UpdatedAt,
	}, nil
}

func (r *TradeHistoryRepository) detectionLogToEntity(log *model.DetectionLog) (*DetectionLogEntity, error) {
	// Convert metadata to JSON
	var metadataJSON string
	if log.Metadata != nil {
		metadataBytes, err := json.Marshal(log.Metadata)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal metadata: %w", err)
		}
		metadataJSON = string(metadataBytes)
	}

	return &DetectionLogEntity{
		ID:          log.ID,
		Type:        log.Type,
		Symbol:      log.Symbol,
		Value:       log.Value,
		Threshold:   log.Threshold,
		Description: log.Description,
		Metadata:    metadataJSON,
		DetectedAt:  log.DetectedAt,
		ProcessedAt: log.ProcessedAt,
		Processed:   log.Processed,
		Result:      log.Result,
		CreatedAt:   log.CreatedAt,
		UpdatedAt:   log.UpdatedAt,
	}, nil
}

func (r *TradeHistoryRepository) entityToDetectionLog(entity *DetectionLogEntity) (*model.DetectionLog, error) {
	// Parse metadata from JSON
	var metadata map[string]interface{}
	if entity.Metadata != "" {
		if err := json.Unmarshal([]byte(entity.Metadata), &metadata); err != nil {
			return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
		}
	}

	return &model.DetectionLog{
		ID:          entity.ID,
		Type:        entity.Type,
		Symbol:      entity.Symbol,
		Value:       entity.Value,
		Threshold:   entity.Threshold,
		Description: entity.Description,
		Metadata:    metadata,
		DetectedAt:  entity.DetectedAt,
		ProcessedAt: entity.ProcessedAt,
		Processed:   entity.Processed,
		Result:      entity.Result,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
	}, nil
}

// Ensure TradeHistoryRepository implements port.TradeHistoryRepository
var _ port.TradeHistoryRepository = (*TradeHistoryRepository)(nil)
