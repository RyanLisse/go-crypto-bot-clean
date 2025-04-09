package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go-crypto-bot-clean/backend/internal/domain/ai/service"
	"gorm.io/gorm"
)

// ConfirmationModel is the GORM model for trade confirmations
type ConfirmationModel struct {
	ID                 string    `gorm:"primaryKey;column:id"`
	UserID             int       `gorm:"column:user_id;index"`
	TradeRequestJSON   string    `gorm:"column:trade_request_json;type:text"`
	RecommendationJSON string    `gorm:"column:recommendation_json;type:text"`
	RiskAssessmentJSON string    `gorm:"column:risk_assessment_json;type:text"`
	Status             string    `gorm:"column:status;index"`
	ConfirmationReason string    `gorm:"column:confirmation_reason"`
	CreatedAt          time.Time `gorm:"column:created_at;autoCreateTime"`
	ExpiresAt          time.Time `gorm:"column:expires_at;index"`
	ConfirmedAt        *time.Time `gorm:"column:confirmed_at"`
}

// TableName specifies the table name for GORM
func (ConfirmationModel) TableName() string {
	return "trade_confirmations"
}

// GormConfirmationRepository implements ConfirmationRepository using GORM
type GormConfirmationRepository struct {
	db *gorm.DB
}

// NewGormConfirmationRepository creates a new GormConfirmationRepository
func NewGormConfirmationRepository(db *gorm.DB) (*GormConfirmationRepository, error) {
	// Auto-migrate the schema
	err := db.AutoMigrate(&ConfirmationModel{})
	if err != nil {
		return nil, fmt.Errorf("failed to migrate trade_confirmations table: %w", err)
	}

	return &GormConfirmationRepository{db: db}, nil
}

// StoreConfirmation saves a trade confirmation to the database
func (r *GormConfirmationRepository) StoreConfirmation(
	ctx context.Context,
	confirmation *service.TradeConfirmation,
) error {
	// Convert trade request to JSON
	tradeRequestJSON, err := json.Marshal(confirmation.TradeRequest)
	if err != nil {
		return fmt.Errorf("failed to marshal trade request: %w", err)
	}

	// Convert recommendation to JSON
	recommendationJSON, err := json.Marshal(confirmation.Recommendation)
	if err != nil {
		return fmt.Errorf("failed to marshal recommendation: %w", err)
	}

	// Convert risk assessment to JSON
	riskAssessmentJSON, err := json.Marshal(confirmation.RiskAssessment)
	if err != nil {
		return fmt.Errorf("failed to marshal risk assessment: %w", err)
	}

	// Create model
	model := ConfirmationModel{
		ID:                 confirmation.ID,
		UserID:             confirmation.UserID,
		TradeRequestJSON:   string(tradeRequestJSON),
		RecommendationJSON: string(recommendationJSON),
		RiskAssessmentJSON: string(riskAssessmentJSON),
		Status:             string(confirmation.Status),
		ConfirmationReason: confirmation.ConfirmationReason,
		CreatedAt:          confirmation.CreatedAt,
		ExpiresAt:          confirmation.ExpiresAt,
		ConfirmedAt:        confirmation.ConfirmedAt,
	}

	// Insert confirmation
	result := r.db.WithContext(ctx).Create(&model)
	if result.Error != nil {
		return fmt.Errorf("failed to store confirmation: %w", result.Error)
	}

	return nil
}

// GetConfirmation gets a trade confirmation from the database
func (r *GormConfirmationRepository) GetConfirmation(
	ctx context.Context,
	id string,
) (*service.TradeConfirmation, error) {
	var model ConfirmationModel

	result := r.db.WithContext(ctx).Where("id = ?", id).First(&model)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("confirmation not found: %s", id)
		}
		return nil, fmt.Errorf("failed to get confirmation: %w", result.Error)
	}

	// Parse trade request from JSON
	var tradeRequest service.TradeRequest
	err := json.Unmarshal([]byte(model.TradeRequestJSON), &tradeRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal trade request: %w", err)
	}

	// Parse recommendation from JSON
	var recommendation service.TradeRecommendation
	err = json.Unmarshal([]byte(model.RecommendationJSON), &recommendation)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal recommendation: %w", err)
	}

	// Parse risk assessment from JSON
	var riskAssessment service.RiskAssessment
	err = json.Unmarshal([]byte(model.RiskAssessmentJSON), &riskAssessment)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal risk assessment: %w", err)
	}

	// Create confirmation
	confirmation := &service.TradeConfirmation{
		ID:                 model.ID,
		UserID:             model.UserID,
		TradeRequest:       &tradeRequest,
		Recommendation:     &recommendation,
		RiskAssessment:     &riskAssessment,
		Status:             service.ConfirmationStatus(model.Status),
		ConfirmationReason: model.ConfirmationReason,
		CreatedAt:          model.CreatedAt,
		ExpiresAt:          model.ExpiresAt,
		ConfirmedAt:        model.ConfirmedAt,
	}

	return confirmation, nil
}

// UpdateConfirmationStatus updates the status of a trade confirmation
func (r *GormConfirmationRepository) UpdateConfirmationStatus(
	ctx context.Context,
	id string,
	status service.ConfirmationStatus,
) error {
	now := time.Now()
	result := r.db.WithContext(ctx).Model(&ConfirmationModel{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":       string(status),
			"confirmed_at": now,
		})

	if result.Error != nil {
		return fmt.Errorf("failed to update confirmation status: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("confirmation not found: %s", id)
	}

	return nil
}

// ListPendingConfirmations lists all pending confirmations for a user
func (r *GormConfirmationRepository) ListPendingConfirmations(
	ctx context.Context,
	userID int,
) ([]*service.TradeConfirmation, error) {
	var models []ConfirmationModel

	result := r.db.WithContext(ctx).
		Where("user_id = ? AND status = ? AND expires_at > ?", userID, string(service.ConfirmationPending), time.Now()).
		Order("created_at DESC").
		Find(&models)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to list pending confirmations: %w", result.Error)
	}

	confirmations := make([]*service.TradeConfirmation, 0, len(models))
	for _, model := range models {
		// Parse trade request from JSON
		var tradeRequest service.TradeRequest
		err := json.Unmarshal([]byte(model.TradeRequestJSON), &tradeRequest)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal trade request: %w", err)
		}

		// Parse recommendation from JSON
		var recommendation service.TradeRecommendation
		err = json.Unmarshal([]byte(model.RecommendationJSON), &recommendation)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal recommendation: %w", err)
		}

		// Parse risk assessment from JSON
		var riskAssessment service.RiskAssessment
		err = json.Unmarshal([]byte(model.RiskAssessmentJSON), &riskAssessment)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal risk assessment: %w", err)
		}

		// Create confirmation
		confirmation := &service.TradeConfirmation{
			ID:                 model.ID,
			UserID:             model.UserID,
			TradeRequest:       &tradeRequest,
			Recommendation:     &recommendation,
			RiskAssessment:     &riskAssessment,
			Status:             service.ConfirmationStatus(model.Status),
			ConfirmationReason: model.ConfirmationReason,
			CreatedAt:          model.CreatedAt,
			ExpiresAt:          model.ExpiresAt,
			ConfirmedAt:        model.ConfirmedAt,
		}

		confirmations = append(confirmations, confirmation)
	}

	return confirmations, nil
}

// CleanupExpiredConfirmations cleans up expired confirmations
func (r *GormConfirmationRepository) CleanupExpiredConfirmations(
	ctx context.Context,
) error {
	result := r.db.WithContext(ctx).
		Model(&ConfirmationModel{}).
		Where("status = ? AND expires_at <= ?", string(service.ConfirmationPending), time.Now()).
		Updates(map[string]interface{}{
			"status": string(service.ConfirmationExpired),
		})

	if result.Error != nil {
		return fmt.Errorf("failed to cleanup expired confirmations: %w", result.Error)
	}

	return nil
}
