package port

import (
	"context"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
)

// TradeHistoryRepository defines the interface for trade history persistence
type TradeHistoryRepository interface {
	// SaveTradeRecord saves a trade record
	SaveTradeRecord(ctx context.Context, record *model.TradeRecord) error
	
	// GetTradeRecords retrieves trade records with filtering
	GetTradeRecords(ctx context.Context, filter TradeHistoryFilter) ([]*model.TradeRecord, error)
	
	// GetTradeRecordByID retrieves a trade record by ID
	GetTradeRecordByID(ctx context.Context, id string) (*model.TradeRecord, error)
	
	// GetTradeRecordsByOrderID retrieves trade records by order ID
	GetTradeRecordsByOrderID(ctx context.Context, orderID string) ([]*model.TradeRecord, error)
	
	// SaveDetectionLog saves a detection log
	SaveDetectionLog(ctx context.Context, log *model.DetectionLog) error
	
	// GetDetectionLogs retrieves detection logs with filtering
	GetDetectionLogs(ctx context.Context, filter DetectionLogFilter) ([]*model.DetectionLog, error)
	
	// MarkDetectionLogProcessed marks a detection log as processed
	MarkDetectionLogProcessed(ctx context.Context, id string, result string) error
	
	// GetUnprocessedDetectionLogs retrieves unprocessed detection logs
	GetUnprocessedDetectionLogs(ctx context.Context, limit int) ([]*model.DetectionLog, error)
}

// TradeHistoryFilter defines filters for retrieving trade records
type TradeHistoryFilter struct {
	UserID    string
	Symbol    string
	Side      model.OrderSide
	Strategy  string
	Tags      []string
	StartTime time.Time
	EndTime   time.Time
	Limit     int
	Offset    int
}

// DetectionLogFilter defines filters for retrieving detection logs
type DetectionLogFilter struct {
	Type       string
	Symbol     string
	Processed  *bool
	StartTime  time.Time
	EndTime    time.Time
	Limit      int
	Offset     int
}
