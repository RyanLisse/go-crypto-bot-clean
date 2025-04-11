// Package huma provides the status endpoint for the API.
package huma

import (
	"context"
	"time"
)

// StatusResponse matches the frontend StatusResponse interface.
type StatusResponse struct {
	Body struct {
		Status      string `json:"status"`
		Version     string `json:"version"`
		Uptime      string `json:"uptime"`
		StartTime   string `json:"start_time"`
		MemoryUsage struct {
			Allocated string `json:"allocated"`
			Total     string `json:"total"`
			System    string `json:"system"`
		} `json:"memory_usage"`
		Goroutines   int `json:"goroutines"`
		ProcessCount int `json:"process_count"`
	} `json:"body"`
}

// StatusHandler handles GET /api/v1/status
func StatusHandler(ctx context.Context, input *struct{}) (*StatusResponse, error) {
	resp := &StatusResponse{}
	resp.Body.Status = "ok"
	resp.Body.Version = "1.0.0"
	resp.Body.Uptime = "24h"
	resp.Body.StartTime = time.Now().Add(-24 * time.Hour).Format(time.RFC3339)
	resp.Body.MemoryUsage.Allocated = "100MB"
	resp.Body.MemoryUsage.Total = "200MB"
	resp.Body.MemoryUsage.System = "300MB"
	resp.Body.Goroutines = 12
	resp.Body.ProcessCount = 2
	return resp, nil
}
