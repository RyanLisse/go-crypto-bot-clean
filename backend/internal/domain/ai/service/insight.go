package service

import (
	"time"
)

// Insight represents an AI-generated insight
type Insight struct {
	ID            string    `json:"id"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	Type          string    `json:"type"`
	Importance    string    `json:"importance"`
	Timestamp     time.Time `json:"timestamp"`
	Metrics       []Metric  `json:"metrics,omitempty"`
	Recommendation string    `json:"recommendation,omitempty"`
}

// Metric represents a metric associated with an insight
type Metric struct {
	Name   string  `json:"name"`
	Value  string  `json:"value"`
	Change float64 `json:"change,omitempty"`
}
