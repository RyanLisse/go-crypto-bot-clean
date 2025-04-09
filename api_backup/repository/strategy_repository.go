package repository

import (
	"context"
	"errors"
	"fmt"

	"go-crypto-bot-clean/api/models"
	"gorm.io/gorm"
)

// Common errors
var (
	ErrStrategyNotFound = errors.New("strategy not found")
)

// StrategyRepository defines the interface for strategy operations
type StrategyRepository interface {
	// Create creates a new strategy
	Create(ctx context.Context, strategy *models.Strategy) error
	
	// GetByID gets a strategy by ID
	GetByID(ctx context.Context, id string) (*models.Strategy, error)
	
	// GetByUserID gets all strategies for a user
	GetByUserID(ctx context.Context, userID string) ([]*models.Strategy, error)
	
	// Update updates a strategy
	Update(ctx context.Context, strategy *models.Strategy) error
	
	// Delete deletes a strategy
	Delete(ctx context.Context, id string) error
	
	// GetParameters gets all parameters for a strategy
	GetParameters(ctx context.Context, strategyID string) ([]*models.StrategyParameter, error)
	
	// AddParameter adds a parameter to a strategy
	AddParameter(ctx context.Context, parameter *models.StrategyParameter) error
	
	// UpdateParameter updates a parameter
	UpdateParameter(ctx context.Context, parameter *models.StrategyParameter) error
	
	// DeleteParameter deletes a parameter
	DeleteParameter(ctx context.Context, id uint) error
	
	// GetPerformance gets the performance metrics for a strategy
	GetPerformance(ctx context.Context, strategyID string) ([]*models.StrategyPerformance, error)
	
	// AddPerformance adds performance metrics for a strategy
	AddPerformance(ctx context.Context, performance *models.StrategyPerformance) error
}

// GormStrategyRepository implements StrategyRepository using GORM
type GormStrategyRepository struct {
	db *gorm.DB
}

// NewGormStrategyRepository creates a new GORM strategy repository
func NewGormStrategyRepository(db *gorm.DB) StrategyRepository {
	return &GormStrategyRepository{
		db: db,
	}
}

// Create creates a new strategy
func (r *GormStrategyRepository) Create(ctx context.Context, strategy *models.Strategy) error {
	if err := r.db.WithContext(ctx).Create(strategy).Error; err != nil {
		return fmt.Errorf("failed to create strategy: %w", err)
	}
	return nil
}

// GetByID gets a strategy by ID
func (r *GormStrategyRepository) GetByID(ctx context.Context, id string) (*models.Strategy, error) {
	var strategy models.Strategy
	if err := r.db.WithContext(ctx).First(&strategy, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrStrategyNotFound
		}
		return nil, fmt.Errorf("failed to get strategy by ID: %w", err)
	}
	return &strategy, nil
}

// GetByUserID gets all strategies for a user
func (r *GormStrategyRepository) GetByUserID(ctx context.Context, userID string) ([]*models.Strategy, error) {
	var strategies []*models.Strategy
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&strategies).Error; err != nil {
		return nil, fmt.Errorf("failed to get strategies by user ID: %w", err)
	}
	return strategies, nil
}

// Update updates a strategy
func (r *GormStrategyRepository) Update(ctx context.Context, strategy *models.Strategy) error {
	// Check if strategy exists
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.Strategy{}).Where("id = ?", strategy.ID).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to check strategy existence: %w", err)
	}
	if count == 0 {
		return ErrStrategyNotFound
	}

	// Update strategy
	if err := r.db.WithContext(ctx).Save(strategy).Error; err != nil {
		return fmt.Errorf("failed to update strategy: %w", err)
	}
	return nil
}

// Delete deletes a strategy
func (r *GormStrategyRepository) Delete(ctx context.Context, id string) error {
	// Check if strategy exists
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.Strategy{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to check strategy existence: %w", err)
	}
	if count == 0 {
		return ErrStrategyNotFound
	}

	// Delete strategy (soft delete)
	if err := r.db.WithContext(ctx).Delete(&models.Strategy{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("failed to delete strategy: %w", err)
	}
	return nil
}

// GetParameters gets all parameters for a strategy
func (r *GormStrategyRepository) GetParameters(ctx context.Context, strategyID string) ([]*models.StrategyParameter, error) {
	// Check if strategy exists
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.Strategy{}).Where("id = ?", strategyID).Count(&count).Error; err != nil {
		return nil, fmt.Errorf("failed to check strategy existence: %w", err)
	}
	if count == 0 {
		return nil, ErrStrategyNotFound
	}

	// Get parameters
	var parameters []*models.StrategyParameter
	if err := r.db.WithContext(ctx).Where("strategy_id = ?", strategyID).Find(&parameters).Error; err != nil {
		return nil, fmt.Errorf("failed to get strategy parameters: %w", err)
	}
	return parameters, nil
}

// AddParameter adds a parameter to a strategy
func (r *GormStrategyRepository) AddParameter(ctx context.Context, parameter *models.StrategyParameter) error {
	// Check if strategy exists
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.Strategy{}).Where("id = ?", parameter.StrategyID).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to check strategy existence: %w", err)
	}
	if count == 0 {
		return ErrStrategyNotFound
	}

	// Add parameter
	if err := r.db.WithContext(ctx).Create(parameter).Error; err != nil {
		return fmt.Errorf("failed to add strategy parameter: %w", err)
	}
	return nil
}

// UpdateParameter updates a parameter
func (r *GormStrategyRepository) UpdateParameter(ctx context.Context, parameter *models.StrategyParameter) error {
	// Check if parameter exists
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.StrategyParameter{}).Where("id = ?", parameter.ID).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to check parameter existence: %w", err)
	}
	if count == 0 {
		return fmt.Errorf("parameter not found")
	}

	// Update parameter
	if err := r.db.WithContext(ctx).Save(parameter).Error; err != nil {
		return fmt.Errorf("failed to update strategy parameter: %w", err)
	}
	return nil
}

// DeleteParameter deletes a parameter
func (r *GormStrategyRepository) DeleteParameter(ctx context.Context, id uint) error {
	// Check if parameter exists
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.StrategyParameter{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to check parameter existence: %w", err)
	}
	if count == 0 {
		return fmt.Errorf("parameter not found")
	}

	// Delete parameter
	if err := r.db.WithContext(ctx).Delete(&models.StrategyParameter{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("failed to delete strategy parameter: %w", err)
	}
	return nil
}

// GetPerformance gets the performance metrics for a strategy
func (r *GormStrategyRepository) GetPerformance(ctx context.Context, strategyID string) ([]*models.StrategyPerformance, error) {
	// Check if strategy exists
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.Strategy{}).Where("id = ?", strategyID).Count(&count).Error; err != nil {
		return nil, fmt.Errorf("failed to check strategy existence: %w", err)
	}
	if count == 0 {
		return nil, ErrStrategyNotFound
	}

	// Get performance metrics
	var performance []*models.StrategyPerformance
	if err := r.db.WithContext(ctx).Where("strategy_id = ?", strategyID).Order("period_end DESC").Find(&performance).Error; err != nil {
		return nil, fmt.Errorf("failed to get strategy performance: %w", err)
	}
	return performance, nil
}

// AddPerformance adds performance metrics for a strategy
func (r *GormStrategyRepository) AddPerformance(ctx context.Context, performance *models.StrategyPerformance) error {
	// Check if strategy exists
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.Strategy{}).Where("id = ?", performance.StrategyID).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to check strategy existence: %w", err)
	}
	if count == 0 {
		return ErrStrategyNotFound
	}

	// Add performance metrics
	if err := r.db.WithContext(ctx).Create(performance).Error; err != nil {
		return fmt.Errorf("failed to add strategy performance: %w", err)
	}
	return nil
}
