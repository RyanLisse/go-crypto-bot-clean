package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go-crypto-bot-clean/backend/internal/domain/models"
	"go.uber.org/zap"
)

// ReportGenerator defines the interface for report generation
type ReportGenerator interface {
	GetLatestReport(ctx context.Context) (*models.PerformanceReport, error)
	GetReportByID(ctx context.Context, id string) (*models.PerformanceReport, error)
	GetReportsByPeriod(ctx context.Context, period models.ReportPeriod, limit int) ([]*models.PerformanceReport, error)
}

// ReportHandler handles requests for performance reports
type ReportHandler struct {
	reportGenerator ReportGenerator
	logger          *zap.Logger
}

// NewReportHandler creates a new ReportHandler
func NewReportHandler(reportGenerator ReportGenerator, logger *zap.Logger) *ReportHandler {
	return &ReportHandler{
		reportGenerator: reportGenerator,
		logger:          logger,
	}
}

// GetLatestReport godoc
// @Summary Get the latest performance report
// @Description Returns the latest generated performance report
// @Tags reports
// @Accept json
// @Produce json
// @Success 200 {object} models.PerformanceReport
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/reports/latest [get]
func (h *ReportHandler) GetLatestReport(c *gin.Context) {
	ctx := c.Request.Context()

	// Get the latest report
	report, err := h.reportGenerator.GetLatestReport(ctx)
	if err != nil {
		h.logger.Error("Failed to get latest report", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get latest report"})
		return
	}

	if report == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No reports found"})
		return
	}

	c.JSON(http.StatusOK, report)
}

// GetReportByID godoc
// @Summary Get a performance report by ID
// @Description Returns a specific performance report by ID
// @Tags reports
// @Accept json
// @Produce json
// @Param id path string true "Report ID"
// @Success 200 {object} models.PerformanceReport
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/reports/{id} [get]
func (h *ReportHandler) GetReportByID(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	// Check if ID is empty
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Report ID is required"})
		return
	}

	// Get the report by ID
	report, err := h.reportGenerator.GetReportByID(ctx, id)
	if err != nil {
		h.logger.Error("Failed to get report", zap.Error(err), zap.String("id", id))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get report"})
		return
	}

	if report == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Report not found"})
		return
	}

	c.JSON(http.StatusOK, report)
}

// GetReportsByPeriod godoc
// @Summary Get performance reports by period
// @Description Returns performance reports for a specific period
// @Tags reports
// @Accept json
// @Produce json
// @Param period query string true "Report period (hourly, daily, weekly, monthly)"
// @Param limit query int false "Maximum number of reports to return" default(10)
// @Success 200 {array} models.PerformanceReport
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/reports [get]
func (h *ReportHandler) GetReportsByPeriod(c *gin.Context) {
	ctx := c.Request.Context()
	periodStr := c.Query("period")
	limitStr := c.DefaultQuery("limit", "10")

	// Validate period
	if periodStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Period is required"})
		return
	}

	period := models.ReportPeriod(periodStr)
	if period != models.ReportPeriodHourly &&
		period != models.ReportPeriodDaily &&
		period != models.ReportPeriodWeekly &&
		period != models.ReportPeriodMonthly {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid period"})
		return
	}

	// Parse limit
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit"})
		return
	}

	// Get reports by period
	reports, err := h.reportGenerator.GetReportsByPeriod(ctx, period, limit)
	if err != nil {
		h.logger.Error("Failed to get reports", zap.Error(err), zap.String("period", string(period)))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get reports"})
		return
	}

	if len(reports) == 0 {
		c.JSON(http.StatusOK, []models.PerformanceReport{})
		return
	}

	c.JSON(http.StatusOK, reports)
}

// GetReportInsights godoc
// @Summary Get insights from a performance report
// @Description Returns the insights extracted from a specific performance report
// @Tags reports
// @Accept json
// @Produce json
// @Param id path string true "Report ID"
// @Success 200 {object} map[string][]string
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/reports/{id}/insights [get]
func (h *ReportHandler) GetReportInsights(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Report ID is required"})
		return
	}

	// Get the report by ID
	report, err := h.reportGenerator.GetReportByID(ctx, id)
	if err != nil {
		h.logger.Error("Failed to get report", zap.Error(err), zap.String("id", id))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get report"})
		return
	}

	if report == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Report not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"insights": report.Insights})
}

// GetReportMetrics godoc
// @Summary Get metrics from a performance report
// @Description Returns the metrics collected for a specific performance report
// @Tags reports
// @Accept json
// @Produce json
// @Param id path string true "Report ID"
// @Success 200 {object} models.PerformanceReportMetrics
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/reports/{id}/metrics [get]
func (h *ReportHandler) GetReportMetrics(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Report ID is required"})
		return
	}

	// Get the report by ID
	report, err := h.reportGenerator.GetReportByID(ctx, id)
	if err != nil {
		h.logger.Error("Failed to get report", zap.Error(err), zap.String("id", id))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get report"})
		return
	}

	if report == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Report not found"})
		return
	}

	c.JSON(http.StatusOK, report.Metrics)
}

// GetReportAnalysis godoc
// @Summary Get analysis from a performance report
// @Description Returns the AI-generated analysis for a specific performance report
// @Tags reports
// @Accept json
// @Produce json
// @Param id path string true "Report ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/reports/{id}/analysis [get]
func (h *ReportHandler) GetReportAnalysis(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Report ID is required"})
		return
	}

	// Get the report by ID
	report, err := h.reportGenerator.GetReportByID(ctx, id)
	if err != nil {
		h.logger.Error("Failed to get report", zap.Error(err), zap.String("id", id))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get report"})
		return
	}

	if report == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Report not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"analysis": report.Analysis})
}

// RegisterRoutes registers the report handler routes
func (h *ReportHandler) RegisterRoutes(router *gin.RouterGroup) {
	reports := router.Group("/reports")
	{
		reports.GET("", h.GetReportsByPeriod)
		reports.GET("/latest", h.GetLatestReport)
		reports.GET("/:id", h.GetReportByID)
		reports.GET("/:id/insights", h.GetReportInsights)
		reports.GET("/:id/metrics", h.GetReportMetrics)
		reports.GET("/:id/analysis", h.GetReportAnalysis)
	}
}
