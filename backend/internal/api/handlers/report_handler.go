package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"go-crypto-bot-clean/backend/internal/domain/models"

	"github.com/go-chi/chi/v5"
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
func (h *ReportHandler) GetLatestReport(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	report, err := h.reportGenerator.GetLatestReport(ctx)
	if err != nil {
		h.logger.Error("Failed to get latest report", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": "Failed to get latest report"})
		return
	}

	if report == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": "No reports found"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(report)
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
func (h *ReportHandler) GetReportByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": "Report ID is required"})
		return
	}

	report, err := h.reportGenerator.GetReportByID(ctx, id)
	if err != nil {
		h.logger.Error("Failed to get report", zap.Error(err), zap.String("id", id))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": "Failed to get report"})
		return
	}

	if report == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": "Report not found"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(report)
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
func (h *ReportHandler) GetReportsByPeriod(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	periodStr := r.URL.Query().Get("period")
	limitStr := r.URL.Query().Get("limit")
	if limitStr == "" {
		limitStr = "10"
	}

	if periodStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": "Period is required"})
		return
	}

	period := models.ReportPeriod(periodStr)
	if period != models.ReportPeriodHourly &&
		period != models.ReportPeriodDaily &&
		period != models.ReportPeriodWeekly &&
		period != models.ReportPeriodMonthly {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": "Invalid period"})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": "Invalid limit"})
		return
	}

	reports, err := h.reportGenerator.GetReportsByPeriod(ctx, period, limit)
	if err != nil {
		h.logger.Error("Failed to get reports", zap.Error(err), zap.String("period", string(period)))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": "Failed to get reports"})
		return
	}

	if len(reports) == 0 {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]models.PerformanceReport{})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(reports)
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
func (h *ReportHandler) GetReportInsights(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": "Report ID is required"})
		return
	}

	report, err := h.reportGenerator.GetReportByID(ctx, id)
	if err != nil {
		h.logger.Error("Failed to get report", zap.Error(err), zap.String("id", id))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": "Failed to get report"})
		return
	}

	if report == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": "Report not found"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"insights": report.Insights})
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
func (h *ReportHandler) GetReportMetrics(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": "Report ID is required"})
		return
	}

	report, err := h.reportGenerator.GetReportByID(ctx, id)
	if err != nil {
		h.logger.Error("Failed to get report", zap.Error(err), zap.String("id", id))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": "Failed to get report"})
		return
	}

	if report == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": "Report not found"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(report.Metrics)
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
func (h *ReportHandler) GetReportAnalysis(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": "Report ID is required"})
		return
	}

	report, err := h.reportGenerator.GetReportByID(ctx, id)
	if err != nil {
		h.logger.Error("Failed to get report", zap.Error(err), zap.String("id", id))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": "Failed to get report"})
		return
	}

	if report == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": "Report not found"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"analysis": report.Analysis})
}

// RegisterRoutes registers the report handler routes
