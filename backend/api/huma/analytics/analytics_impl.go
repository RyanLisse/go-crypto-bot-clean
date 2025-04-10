package analytics

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"

	"go-crypto-bot-clean/backend/internal/api/dto/response"
	"go-crypto-bot-clean/backend/internal/core/analytics"
	domainmodels "go-crypto-bot-clean/backend/internal/domain/models" // Added import
)

// RegisterAnalyticsEndpoints registers analytics endpoints with the Huma API
func RegisterAnalyticsEndpoints(api huma.API, basePath string, analyticsService analytics.TradeAnalyticsService) {
	// GET /analytics/trade-analytics
	huma.Register(api, huma.Operation{
		OperationID: "get-trade-analytics",
		Method:      http.MethodGet,
		Path:        basePath + "/analytics/trade-analytics",
		Summary:     "Get trade analytics",
		Description: "Returns trade analytics for a specific time range and timeframe",
		Tags:        []string{"Analytics"},
	}, func(ctx context.Context, input *struct {
		TimeFrame string    `query:"timeFrame"`
		StartTime time.Time `query:"startTime"`
		EndTime   time.Time `query:"endTime"`
	}) (*struct {
		Body response.TradeAnalyticsResponse `json:"body"`
	}, error) {
		if input.StartTime.IsZero() {
			input.StartTime = time.Now().AddDate(0, 0, -30)
		}
		if input.EndTime.IsZero() {
			input.EndTime = time.Now()
		}

		timeFrame, err := parseTimeFrame(input.TimeFrame)
		if err != nil {
			return nil, fmt.Errorf("invalid timeframe: %w", err)
		}

		// Convert core TimeFrame to domainmodels.TimeFrame
		domainTimeFrame := domainmodels.TimeFrame(timeFrame)
		data, err := analyticsService.GetTradeAnalytics(ctx, domainTimeFrame, input.StartTime, input.EndTime)
		if err != nil {
			return nil, err
		}

		return &struct {
			Body response.TradeAnalyticsResponse `json:"body"`
		}{
			Body: response.TradeAnalyticsFromModel(data),
		}, nil
	})

	// GET /analytics/trade-performance
	huma.Register(api, huma.Operation{
		OperationID: "get-trade-performance",
		Method:      http.MethodGet,
		Path:        basePath + "/analytics/trade-performance",
		Summary:     "Get trade performance",
		Description: "Returns performance details for a specific trade",
		Tags:        []string{"Analytics"},
	}, func(ctx context.Context, input *struct {
		TradeID string `query:"tradeID"`
	}) (*struct {
		Body response.TradePerformanceResponse `json:"body"`
	}, error) {
		if input.TradeID == "" {
			return nil, fmt.Errorf("tradeID is required")
		}

		perf, err := analyticsService.GetTradePerformance(ctx, input.TradeID)
		if err != nil {
			return nil, err
		}

		return &struct {
			Body response.TradePerformanceResponse `json:"body"`
		}{
			Body: response.TradePerformanceFromModel(perf),
		}, nil
	})
}

// parseTimeFrame converts a string to the corresponding TimeFrame enum
func parseTimeFrame(tf string) (analytics.TimeFrame, error) {
	switch tf {
	case "day":
		return analytics.TimeFrameDay, nil
	case "week":
		return analytics.TimeFrameWeek, nil
	case "month":
		return analytics.TimeFrameMonth, nil
	case "quarter":
		return analytics.TimeFrameQuarter, nil
	case "year":
		return analytics.TimeFrameYear, nil
	case "":
		return analytics.TimeFrameAll, nil
	default:
		return analytics.TimeFrameAll, fmt.Errorf("invalid timeframe: %s", tf)
	}
}
