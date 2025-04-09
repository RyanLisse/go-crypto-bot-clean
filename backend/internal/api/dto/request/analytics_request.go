package request

import "time"

// TradeAnalyticsRequest represents a request for trade analytics
type TradeAnalyticsRequest struct {
	StartTime time.Time `form:"start_time" time_format:"2006-01-02T15:04:05Z07:00"`
	EndTime   time.Time `form:"end_time" time_format:"2006-01-02T15:04:05Z07:00"`
	TimeFrame string    `form:"time_frame" binding:"omitempty,oneof=day week month quarter year all"`
	Symbol    string    `form:"symbol"`
	Strategy  string    `form:"strategy"`
}

// TradePerformanceRequest represents a request for trade performance
type TradePerformanceRequest struct {
	StartTime time.Time `form:"start_time" time_format:"2006-01-02T15:04:05Z07:00"`
	EndTime   time.Time `form:"end_time" time_format:"2006-01-02T15:04:05Z07:00"`
	Limit     int       `form:"limit"`
}

// BalanceHistoryRequest represents a request for balance history
type BalanceHistoryRequest struct {
	StartTime time.Time `form:"start_time" time_format:"2006-01-02T15:04:05Z07:00"`
	EndTime   time.Time `form:"end_time" time_format:"2006-01-02T15:04:05Z07:00"`
	Interval  string    `form:"interval"`
}
