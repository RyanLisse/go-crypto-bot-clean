package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ryanlisse/go-crypto-bot/internal/domain/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap/zaptest"
)

// MockReportGenerator is a mock implementation of the report generator
type MockReportGenerator struct {
	mock.Mock
}

func (m *MockReportGenerator) GetLatestReport(ctx context.Context) (*models.PerformanceReport, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PerformanceReport), args.Error(1)
}

func (m *MockReportGenerator) GetReportByID(ctx context.Context, id string) (*models.PerformanceReport, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PerformanceReport), args.Error(1)
}

func (m *MockReportGenerator) GetReportsByPeriod(ctx context.Context, period models.ReportPeriod, limit int) ([]*models.PerformanceReport, error) {
	args := m.Called(ctx, period, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.PerformanceReport), args.Error(1)
}

func TestNewReportHandler(t *testing.T) {
	// Create dependencies
	reportGenerator := &MockReportGenerator{}
	logger := zaptest.NewLogger(t)

	// Create handler
	handler := NewReportHandler(reportGenerator, logger)

	// Assert
	assert.NotNil(t, handler)
	assert.Equal(t, reportGenerator, handler.reportGenerator)
	assert.Equal(t, logger, handler.logger)
}

func TestReportHandler_GetLatestReport(t *testing.T) {
	// Set up Gin
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Create dependencies
	reportGenerator := &MockReportGenerator{}
	logger := zaptest.NewLogger(t)

	// Create handler
	handler := NewReportHandler(reportGenerator, logger)

	// Register routes
	router.GET("/api/v1/reports/latest", handler.GetLatestReport)

	// Create test data
	report := &models.PerformanceReport{
		ID:          "test-id",
		GeneratedAt: time.Now(),
		Period:      "hourly",
		Analysis:    "Test analysis",
		Insights:    []string{"Insight 1", "Insight 2"},
	}

	// Test cases
	tests := []struct {
		name           string
		setupMock      func()
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name: "Success",
			setupMock: func() {
				reportGenerator.On("GetLatestReport", mock.Anything).Return(report, nil).Once()
			},
			expectedStatus: http.StatusOK,
			expectedBody:   report,
		},
		{
			name: "No reports found",
			setupMock: func() {
				reportGenerator.On("GetLatestReport", mock.Anything).Return(nil, nil).Once()
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   gin.H{"error": "No reports found"},
		},
		{
			name: "Error",
			setupMock: func() {
				reportGenerator.On("GetLatestReport", mock.Anything).Return(nil, errors.New("test error")).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   gin.H{"error": "Failed to get latest report"},
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up mock
			tt.setupMock()

			// Create request
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/reports/latest", nil)
			w := httptest.NewRecorder()

			// Serve request
			router.ServeHTTP(w, req)

			// Assert status code
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Assert body
			if tt.expectedStatus == http.StatusOK {
				var response models.PerformanceReport
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, report.ID, response.ID)
				assert.Equal(t, report.Period, response.Period)
				assert.Equal(t, report.Analysis, response.Analysis)
				assert.Equal(t, report.Insights, response.Insights)
			} else {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBody.(gin.H)["error"], response["error"])
			}

			// Verify all expectations were met
			reportGenerator.AssertExpectations(t)
		})
	}
}

func TestReportHandler_GetReportByID(t *testing.T) {
	// Set up Gin
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Create dependencies
	reportGenerator := &MockReportGenerator{}
	logger := zaptest.NewLogger(t)

	// Create handler
	handler := NewReportHandler(reportGenerator, logger)

	// Register routes
	router.GET("/api/v1/reports/:id", handler.GetReportByID)
	// Add a route for the root path to test missing ID
	router.GET("/api/v1/reports", func(c *gin.Context) {
		handler.GetReportByID(c)
	})

	// Create test data
	id := "test-id"
	report := &models.PerformanceReport{
		ID:          id,
		GeneratedAt: time.Now(),
		Period:      "hourly",
		Analysis:    "Test analysis",
		Insights:    []string{"Insight 1", "Insight 2"},
	}

	// Test cases
	tests := []struct {
		name           string
		id             string
		setupMock      func()
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name: "Success",
			id:   id,
			setupMock: func() {
				reportGenerator.On("GetReportByID", mock.Anything, id).Return(report, nil).Once()
			},
			expectedStatus: http.StatusOK,
			expectedBody:   report,
		},
		{
			name: "Report not found",
			id:   id,
			setupMock: func() {
				reportGenerator.On("GetReportByID", mock.Anything, id).Return(nil, nil).Once()
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   gin.H{"error": "Report not found"},
		},
		{
			name: "Error",
			id:   id,
			setupMock: func() {
				reportGenerator.On("GetReportByID", mock.Anything, id).Return(nil, errors.New("test error")).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   gin.H{"error": "Failed to get report"},
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up mock
			tt.setupMock()

			// Create request
			url := "/api/v1/reports/"
			if tt.id != "" {
				url += tt.id
			}
			req, _ := http.NewRequest(http.MethodGet, url, nil)
			w := httptest.NewRecorder()

			// Serve request
			router.ServeHTTP(w, req)

			// Assert status code
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Assert body
			if tt.expectedStatus == http.StatusOK {
				var response models.PerformanceReport
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, report.ID, response.ID)
				assert.Equal(t, report.Period, response.Period)
				assert.Equal(t, report.Analysis, response.Analysis)
				assert.Equal(t, report.Insights, response.Insights)
			} else {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBody.(gin.H)["error"], response["error"])
			}

			// Verify all expectations were met
			reportGenerator.AssertExpectations(t)
		})
	}

	// Test missing ID separately
	t.Run("Missing ID", func(t *testing.T) {
		// Create a new Gin context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// Create a request with no ID parameter
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/reports", nil)
		c.Request = req

		// Call the handler directly
		handler.GetReportByID(c)

		// Assert status code
		assert.Equal(t, http.StatusBadRequest, w.Code)

		// Assert body
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Report ID is required", response["error"])
	})
}

func TestReportHandler_GetReportsByPeriod(t *testing.T) {
	// Set up Gin
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Create dependencies
	reportGenerator := &MockReportGenerator{}
	logger := zaptest.NewLogger(t)

	// Create handler
	handler := NewReportHandler(reportGenerator, logger)

	// Register routes
	router.GET("/api/v1/reports", handler.GetReportsByPeriod)

	// Create test data
	period := models.ReportPeriodHourly
	limit := 10
	reports := []*models.PerformanceReport{
		{
			ID:          "test-id-1",
			GeneratedAt: time.Now(),
			Period:      string(period),
			Analysis:    "Test analysis 1",
			Insights:    []string{"Insight 1", "Insight 2"},
		},
		{
			ID:          "test-id-2",
			GeneratedAt: time.Now().Add(-time.Hour),
			Period:      string(period),
			Analysis:    "Test analysis 2",
			Insights:    []string{"Insight 3", "Insight 4"},
		},
	}

	// Test cases
	tests := []struct {
		name           string
		period         string
		limit          string
		setupMock      func()
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:   "Success",
			period: string(period),
			limit:  "10",
			setupMock: func() {
				reportGenerator.On("GetReportsByPeriod", mock.Anything, period, limit).Return(reports, nil).Once()
			},
			expectedStatus: http.StatusOK,
			expectedBody:   reports,
		},
		{
			name:   "No reports found",
			period: string(period),
			limit:  "10",
			setupMock: func() {
				reportGenerator.On("GetReportsByPeriod", mock.Anything, period, limit).Return([]*models.PerformanceReport{}, nil).Once()
			},
			expectedStatus: http.StatusOK,
			expectedBody:   []*models.PerformanceReport{},
		},
		{
			name:   "Error",
			period: string(period),
			limit:  "10",
			setupMock: func() {
				reportGenerator.On("GetReportsByPeriod", mock.Anything, period, limit).Return(nil, errors.New("test error")).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   gin.H{"error": "Failed to get reports"},
		},
		{
			name:           "Missing period",
			period:         "",
			limit:          "10",
			setupMock:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   gin.H{"error": "Period is required"},
		},
		{
			name:           "Invalid period",
			period:         "invalid",
			limit:          "10",
			setupMock:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   gin.H{"error": "Invalid period"},
		},
		{
			name:           "Invalid limit",
			period:         string(period),
			limit:          "invalid",
			setupMock:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   gin.H{"error": "Invalid limit"},
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up mock
			tt.setupMock()

			// Create request
			url := "/api/v1/reports?period=" + tt.period
			if tt.limit != "" {
				url += "&limit=" + tt.limit
			}
			req, _ := http.NewRequest(http.MethodGet, url, nil)
			w := httptest.NewRecorder()

			// Serve request
			router.ServeHTTP(w, req)

			// Assert status code
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Assert body
			if tt.expectedStatus == http.StatusOK {
				var response []*models.PerformanceReport
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, len(tt.expectedBody.([]*models.PerformanceReport)), len(response))
				if len(response) > 0 {
					assert.Equal(t, tt.expectedBody.([]*models.PerformanceReport)[0].ID, response[0].ID)
				}
			} else {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBody.(gin.H)["error"], response["error"])
			}

			// Verify all expectations were met
			reportGenerator.AssertExpectations(t)
		})
	}
}
