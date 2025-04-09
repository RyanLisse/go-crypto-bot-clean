package websocket

import (
	"sync"
	"time"

	"go-crypto-bot-clean/backend/internal/domain/models"
	"go.uber.org/zap"
)

// ReportMessageType is the type for report messages
const (
	PerformanceReportType MessageType = "performance_report"
)

// PerformanceReportPayload represents a performance report summary
type PerformanceReportPayload struct {
	ID       string    `json:"id"`
	Time     time.Time `json:"time"`
	Period   string    `json:"period"`
	Summary  string    `json:"summary"`
	Insights []string  `json:"insights"`
}

// ReportService handles performance report broadcasting
type ReportService struct {
	hub    *Hub
	logger *zap.Logger
	mu     sync.Mutex
}

// NewReportService creates a new ReportService
func NewReportService(hub *Hub, logger *zap.Logger) *ReportService {
	return &ReportService{
		hub:    hub,
		logger: logger,
	}
}

// BroadcastReport broadcasts a report to all connected clients
func (s *ReportService) BroadcastReport(report *models.PerformanceReport) {
	// Create message payload
	payload := PerformanceReportPayload{
		ID:       report.ID,
		Time:     report.GeneratedAt,
		Period:   report.Period,
		Summary:  getSummary(report.Analysis),
		Insights: report.Insights,
	}

	// Create WebSocket message
	message := WSMessage{
		Type:      PerformanceReportType,
		Timestamp: time.Now().Unix(),
		Payload:   payload,
	}

	// Broadcast to all clients
	s.hub.Broadcast(message)

	s.logger.Debug("Broadcast report",
		zap.String("id", report.ID),
		zap.String("period", report.Period),
	)
}

// getSummary extracts a summary from the analysis
func getSummary(analysis string) string {
	// In a real implementation, you would extract a summary from the analysis
	// For now, just return the first 100 characters
	if len(analysis) > 100 {
		return analysis[:100] + "..."
	}
	return analysis
}
