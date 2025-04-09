package models

import (
	"time"
)

// PerformanceMetrics represents the metrics to be analyzed
type PerformanceReportMetrics struct {
	Timestamp   time.Time              `json:"timestamp"`
	Metrics     map[string]interface{} `json:"metrics"`
	SystemState SystemState            `json:"system_state"`
}

// SystemState represents the system state at the time of metrics collection
type SystemState struct {
	CPUUsage    float64 `json:"cpu_usage"`
	MemoryUsage float64 `json:"memory_usage"`
	Latency     int64   `json:"latency_ms"`
	Goroutines  int     `json:"goroutines"`
	Uptime      string  `json:"uptime"`
}

// Report represents a generated performance report
type PerformanceReport struct {
	ID          string                  `json:"id"`
	GeneratedAt time.Time               `json:"generated_at"`
	Period      string                  `json:"period"`
	Analysis    string                  `json:"analysis"`
	Metrics     PerformanceReportMetrics `json:"metrics"`
	Insights    []string                `json:"insights"`
}

// ReportPeriod represents the period of a report
type ReportPeriod string

// Report periods
const (
	ReportPeriodHourly  ReportPeriod = "hourly"
	ReportPeriodDaily   ReportPeriod = "daily"
	ReportPeriodWeekly  ReportPeriod = "weekly"
	ReportPeriodMonthly ReportPeriod = "monthly"
)
