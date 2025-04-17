package repo

import (
	"context"
	"encoding/json"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"gorm.io/gorm"
)

// RiskAssessmentEntity is the GORM entity for risk assessment
type RiskAssessmentEntity struct {
	ID             string     `gorm:"column:id;primaryKey"`
	UserID         string     `gorm:"column:user_id;index"`
	Type           string     `gorm:"column:type;index"`   // RiskType as string
	Level          string     `gorm:"column:level;index"`  // RiskLevel as string
	Status         string     `gorm:"column:status;index"` // RiskStatus as string
	Symbol         string     `gorm:"column:symbol;index"`
	PositionID     string     `gorm:"column:position_id;index"`
	OrderID        string     `gorm:"column:order_id;index"`
	Score          float64    `gorm:"column:score"`
	Message        string     `gorm:"column:message;type:text"`
	Recommendation string     `gorm:"column:recommendation;type:text"`
	MetadataJSON   string     `gorm:"column:metadata_json;type:text"` // JSON string of metadata
	CreatedAt      time.Time  `gorm:"column:created_at;index"`
	UpdatedAt      time.Time  `gorm:"column:updated_at"`
	ResolvedAt     *time.Time `gorm:"column:resolved_at"`
}

// TableName overrides the table name
func (RiskAssessmentEntity) TableName() string {
	return "risk_assessments"
}

// toEntity converts a risk assessment model to entity
func toRiskAssessmentEntity(model *model.RiskAssessment) (*RiskAssessmentEntity, error) {
	var metadataJSON string
	if model.Metadata != nil {
		metadata, err := json.Marshal(model.Metadata)
		if err != nil {
			return nil, err
		}
		metadataJSON = string(metadata)
	}

	return &RiskAssessmentEntity{
		ID:             model.ID,
		UserID:         model.UserID,
		Type:           string(model.Type),
		Level:          string(model.Level),
		Status:         string(model.Status),
		Symbol:         model.Symbol,
		PositionID:     model.PositionID,
		OrderID:        model.OrderID,
		Score:          model.Score,
		Message:        model.Message,
		Recommendation: model.Recommendation,
		MetadataJSON:   metadataJSON,
		CreatedAt:      model.CreatedAt,
		UpdatedAt:      time.Now(),
		ResolvedAt:     model.ResolvedAt,
	}, nil
}

// toDomain converts a risk assessment entity to domain model
func (e *RiskAssessmentEntity) toDomain() (*model.RiskAssessment, error) {
	var metadata interface{}
	if e.MetadataJSON != "" {
		if err := json.Unmarshal([]byte(e.MetadataJSON), &metadata); err != nil {
			return nil, err
		}
	}

	return &model.RiskAssessment{
		ID:             e.ID,
		UserID:         e.UserID,
		Type:           model.RiskType(e.Type),
		Level:          model.RiskLevel(e.Level),
		Status:         model.RiskStatus(e.Status),
		Symbol:         e.Symbol,
		PositionID:     e.PositionID,
		OrderID:        e.OrderID,
		Score:          e.Score,
		Message:        e.Message,
		Recommendation: e.Recommendation,
		Metadata:       metadata,
		CreatedAt:      e.CreatedAt,
		UpdatedAt:      e.UpdatedAt,
		ResolvedAt:     e.ResolvedAt,
	}, nil
}

// GormRiskAssessmentRepository implements the RiskAssessmentRepository using GORM
type GormRiskAssessmentRepository struct {
	db *gorm.DB
}

// NewGormRiskAssessmentRepository creates a new instance of GormRiskAssessmentRepository
func NewGormRiskAssessmentRepository(db *gorm.DB) *GormRiskAssessmentRepository {
	return &GormRiskAssessmentRepository{
		db: db,
	}
}

// Create adds a new risk assessment
func (r *GormRiskAssessmentRepository) Create(ctx context.Context, assessment *model.RiskAssessment) error {
	entity, err := toRiskAssessmentEntity(assessment)
	if err != nil {
		return err
	}

	return r.db.WithContext(ctx).Create(entity).Error
}

// Update updates an existing risk assessment
func (r *GormRiskAssessmentRepository) Update(ctx context.Context, assessment *model.RiskAssessment) error {
	entity, err := toRiskAssessmentEntity(assessment)
	if err != nil {
		return err
	}

	return r.db.WithContext(ctx).Save(entity).Error
}

// GetByID retrieves a risk assessment by its ID
func (r *GormRiskAssessmentRepository) GetByID(ctx context.Context, id string) (*model.RiskAssessment, error) {
	var entity RiskAssessmentEntity
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&entity).Error; err != nil {
		return nil, err
	}
	return entity.toDomain()
}

// GetByUserID retrieves risk assessments for a specific user with pagination
func (r *GormRiskAssessmentRepository) GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*model.RiskAssessment, error) {
	var entities []RiskAssessmentEntity

	query := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&entities).Error; err != nil {
		return nil, err
	}

	assessments := make([]*model.RiskAssessment, 0, len(entities))
	for _, entity := range entities {
		assessment, err := entity.toDomain()
		if err != nil {
			return nil, err
		}
		assessments = append(assessments, assessment)
	}

	return assessments, nil
}

// GetActiveByUserID retrieves active risk assessments for a user
func (r *GormRiskAssessmentRepository) GetActiveByUserID(ctx context.Context, userID string) ([]*model.RiskAssessment, error) {
	var entities []RiskAssessmentEntity

	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND status = ?", userID, string(model.RiskStatusActive)).
		Order("created_at DESC").
		Find(&entities).Error; err != nil {
		return nil, err
	}

	assessments := make([]*model.RiskAssessment, 0, len(entities))
	for _, entity := range entities {
		assessment, err := entity.toDomain()
		if err != nil {
			return nil, err
		}
		assessments = append(assessments, assessment)
	}

	return assessments, nil
}

// GetBySymbol retrieves risk assessments for a specific symbol
func (r *GormRiskAssessmentRepository) GetBySymbol(ctx context.Context, symbol string, limit, offset int) ([]*model.RiskAssessment, error) {
	var entities []RiskAssessmentEntity

	query := r.db.WithContext(ctx).
		Where("symbol = ?", symbol).
		Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&entities).Error; err != nil {
		return nil, err
	}

	assessments := make([]*model.RiskAssessment, 0, len(entities))
	for _, entity := range entities {
		assessment, err := entity.toDomain()
		if err != nil {
			return nil, err
		}
		assessments = append(assessments, assessment)
	}

	return assessments, nil
}

// GetByType retrieves risk assessments of a specific type
func (r *GormRiskAssessmentRepository) GetByType(ctx context.Context, riskType model.RiskType, limit, offset int) ([]*model.RiskAssessment, error) {
	var entities []RiskAssessmentEntity

	query := r.db.WithContext(ctx).
		Where("type = ?", string(riskType)).
		Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&entities).Error; err != nil {
		return nil, err
	}

	assessments := make([]*model.RiskAssessment, 0, len(entities))
	for _, entity := range entities {
		assessment, err := entity.toDomain()
		if err != nil {
			return nil, err
		}
		assessments = append(assessments, assessment)
	}

	return assessments, nil
}

// GetByLevel retrieves risk assessments of a specific level
func (r *GormRiskAssessmentRepository) GetByLevel(ctx context.Context, level model.RiskLevel, limit, offset int) ([]*model.RiskAssessment, error) {
	var entities []RiskAssessmentEntity

	query := r.db.WithContext(ctx).
		Where("level = ?", string(level)).
		Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&entities).Error; err != nil {
		return nil, err
	}

	assessments := make([]*model.RiskAssessment, 0, len(entities))
	for _, entity := range entities {
		assessment, err := entity.toDomain()
		if err != nil {
			return nil, err
		}
		assessments = append(assessments, assessment)
	}

	return assessments, nil
}

// GetByTimeRange retrieves risk assessments within a time range
func (r *GormRiskAssessmentRepository) GetByTimeRange(ctx context.Context, from, to time.Time, limit, offset int) ([]*model.RiskAssessment, error) {
	var entities []RiskAssessmentEntity

	query := r.db.WithContext(ctx).
		Where("created_at BETWEEN ? AND ?", from, to).
		Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&entities).Error; err != nil {
		return nil, err
	}

	assessments := make([]*model.RiskAssessment, 0, len(entities))
	for _, entity := range entities {
		assessment, err := entity.toDomain()
		if err != nil {
			return nil, err
		}
		assessments = append(assessments, assessment)
	}

	return assessments, nil
}

// Count returns the total number of risk assessments matching the specified filters
func (r *GormRiskAssessmentRepository) Count(ctx context.Context, filters map[string]interface{}) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&RiskAssessmentEntity{})

	for key, value := range filters {
		query = query.Where(key+" = ?", value)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

// Delete removes a risk assessment
func (r *GormRiskAssessmentRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&RiskAssessmentEntity{}, "id = ?", id).Error
}

// Ensure GormRiskAssessmentRepository implements port.RiskAssessmentRepository
var _ port.RiskAssessmentRepository = (*GormRiskAssessmentRepository)(nil)
