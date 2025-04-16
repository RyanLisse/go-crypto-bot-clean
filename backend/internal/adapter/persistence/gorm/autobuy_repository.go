package gorm

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/persistence/gorm/entity"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// AutoBuyRuleRepository implements the port.AutoBuyRuleRepository interface using GORM
type AutoBuyRuleRepository struct {
	db     *gorm.DB
	logger zerolog.Logger
}

// NewAutoBuyRuleRepository creates a new AutoBuyRuleRepository
func NewAutoBuyRuleRepository(db *gorm.DB, logger zerolog.Logger) *AutoBuyRuleRepository {
	return &AutoBuyRuleRepository{
		db:     db,
		logger: logger.With().Str("component", "autobuy_rule_repository").Logger(),
	}
}

// Create adds a new auto-buy rule
func (r *AutoBuyRuleRepository) Create(ctx context.Context, rule *model.AutoBuyRule) error {
	entity := autoBuyRuleToEntity(rule)
	result := r.db.WithContext(ctx).Create(&entity)
	if result.Error != nil {
		r.logger.Error().Err(result.Error).Str("id", rule.ID).Msg("Failed to create auto-buy rule")
		return result.Error
	}

	// Update the domain model with the generated ID if not already set
	if rule.ID == "" {
		rule.ID = entity.ID
	}
	rule.CreatedAt = entity.CreatedAt
	rule.UpdatedAt = entity.UpdatedAt

	r.logger.Info().Str("id", rule.ID).Str("userId", rule.UserID).Str("symbol", rule.Symbol).Msg("Created auto-buy rule")
	return nil
}

// Update updates an existing auto-buy rule
func (r *AutoBuyRuleRepository) Update(ctx context.Context, rule *model.AutoBuyRule) error {
	entity := autoBuyRuleToEntity(rule)
	result := r.db.WithContext(ctx).Save(&entity)
	if result.Error != nil {
		r.logger.Error().Err(result.Error).Str("id", rule.ID).Msg("Failed to update auto-buy rule")
		return result.Error
	}

	rule.UpdatedAt = entity.UpdatedAt
	r.logger.Info().Str("id", rule.ID).Msg("Updated auto-buy rule")
	return nil
}

// GetByID retrieves an auto-buy rule by its ID
func (r *AutoBuyRuleRepository) GetByID(ctx context.Context, id string) (*model.AutoBuyRule, error) {
	var entity entity.AutoBuyRuleEntity
	result := r.db.WithContext(ctx).First(&entity, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		r.logger.Error().Err(result.Error).Str("id", id).Msg("Failed to get auto-buy rule by ID")
		return nil, result.Error
	}

	rule := autoBuyRuleFromEntity(&entity)
	return rule, nil
}

// GetByUserID retrieves auto-buy rules for a specific user
func (r *AutoBuyRuleRepository) GetByUserID(ctx context.Context, userID string) ([]*model.AutoBuyRule, error) {
	var entities []entity.AutoBuyRuleEntity
	result := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&entities)
	if result.Error != nil {
		r.logger.Error().Err(result.Error).Str("userId", userID).Msg("Failed to get auto-buy rules by user ID")
		return nil, result.Error
	}

	rules := make([]*model.AutoBuyRule, len(entities))
	for i, entity := range entities {
		rules[i] = autoBuyRuleFromEntity(&entity)
	}

	r.logger.Debug().Str("userId", userID).Int("count", len(rules)).Msg("Retrieved auto-buy rules by user ID")
	return rules, nil
}

// GetBySymbol retrieves auto-buy rules for a specific symbol
func (r *AutoBuyRuleRepository) GetBySymbol(ctx context.Context, symbol string) ([]*model.AutoBuyRule, error) {
	var entities []entity.AutoBuyRuleEntity
	result := r.db.WithContext(ctx).Where("symbol = ?", symbol).Find(&entities)
	if result.Error != nil {
		r.logger.Error().Err(result.Error).Str("symbol", symbol).Msg("Failed to get auto-buy rules by symbol")
		return nil, result.Error
	}

	rules := make([]*model.AutoBuyRule, len(entities))
	for i, entity := range entities {
		rules[i] = autoBuyRuleFromEntity(&entity)
	}

	r.logger.Debug().Str("symbol", symbol).Int("count", len(rules)).Msg("Retrieved auto-buy rules by symbol")
	return rules, nil
}

// GetActive retrieves all active auto-buy rules
func (r *AutoBuyRuleRepository) GetActive(ctx context.Context) ([]*model.AutoBuyRule, error) {
	var entities []entity.AutoBuyRuleEntity
	result := r.db.WithContext(ctx).Where("is_enabled = ?", true).Find(&entities)
	if result.Error != nil {
		r.logger.Error().Err(result.Error).Msg("Failed to get active auto-buy rules")
		return nil, result.Error
	}

	rules := make([]*model.AutoBuyRule, len(entities))
	for i, entity := range entities {
		rules[i] = autoBuyRuleFromEntity(&entity)
	}

	r.logger.Debug().Int("count", len(rules)).Msg("Retrieved active auto-buy rules")
	return rules, nil
}

// GetActiveByUserID retrieves active auto-buy rules for a specific user
func (r *AutoBuyRuleRepository) GetActiveByUserID(ctx context.Context, userID string) ([]*model.AutoBuyRule, error) {
	var entities []entity.AutoBuyRuleEntity
	result := r.db.WithContext(ctx).Where("user_id = ? AND is_enabled = ?", userID, true).Find(&entities)
	if result.Error != nil {
		r.logger.Error().Err(result.Error).Str("userId", userID).Msg("Failed to get active auto-buy rules by user ID")
		return nil, result.Error
	}

	rules := make([]*model.AutoBuyRule, len(entities))
	for i, entity := range entities {
		rules[i] = autoBuyRuleFromEntity(&entity)
	}

	r.logger.Debug().Str("userId", userID).Int("count", len(rules)).Msg("Retrieved active auto-buy rules by user ID")
	return rules, nil
}

// GetActiveBySymbol retrieves active auto-buy rules for a specific symbol
func (r *AutoBuyRuleRepository) GetActiveBySymbol(ctx context.Context, symbol string) ([]*model.AutoBuyRule, error) {
	var entities []entity.AutoBuyRuleEntity
	result := r.db.WithContext(ctx).Where("symbol = ? AND is_enabled = ?", symbol, true).Find(&entities)
	if result.Error != nil {
		r.logger.Error().Err(result.Error).Str("symbol", symbol).Msg("Failed to get active auto-buy rules by symbol")
		return nil, result.Error
	}

	rules := make([]*model.AutoBuyRule, len(entities))
	for i, entity := range entities {
		rules[i] = autoBuyRuleFromEntity(&entity)
	}

	r.logger.Debug().Str("symbol", symbol).Int("count", len(rules)).Msg("Retrieved active auto-buy rules by symbol")
	return rules, nil
}

// GetByTriggerType retrieves auto-buy rules with a specific trigger type
func (r *AutoBuyRuleRepository) GetByTriggerType(ctx context.Context, triggerType model.TriggerType) ([]*model.AutoBuyRule, error) {
	var entities []entity.AutoBuyRuleEntity
	result := r.db.WithContext(ctx).Where("trigger_type = ?", string(triggerType)).Find(&entities)
	if result.Error != nil {
		r.logger.Error().Err(result.Error).Str("triggerType", string(triggerType)).Msg("Failed to get auto-buy rules by trigger type")
		return nil, result.Error
	}

	rules := make([]*model.AutoBuyRule, len(entities))
	for i, entity := range entities {
		rules[i] = autoBuyRuleFromEntity(&entity)
	}

	r.logger.Debug().Str("triggerType", string(triggerType)).Int("count", len(rules)).Msg("Retrieved auto-buy rules by trigger type")
	return rules, nil
}

// Delete removes an auto-buy rule
func (r *AutoBuyRuleRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&entity.AutoBuyRuleEntity{}, "id = ?", id)
	if result.Error != nil {
		r.logger.Error().Err(result.Error).Str("id", id).Msg("Failed to delete auto-buy rule")
		return result.Error
	}

	if result.RowsAffected == 0 {
		r.logger.Warn().Str("id", id).Msg("Auto-buy rule not found for deletion")
		return fmt.Errorf("auto-buy rule not found for deletion: %s", id)
	}

	r.logger.Info().Str("id", id).Msg("Deleted auto-buy rule")
	return nil
}

// Count returns the total number of auto-buy rules matching the specified filters
func (r *AutoBuyRuleRepository) Count(ctx context.Context, filters map[string]interface{}) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&entity.AutoBuyRuleEntity{})

	// Apply filters
	for key, value := range filters {
		query = query.Where(fmt.Sprintf("%s = ?", key), value)
	}

	result := query.Count(&count)
	if result.Error != nil {
		r.logger.Error().Err(result.Error).Msg("Failed to count auto-buy rules")
		return 0, result.Error
	}

	return count, nil
}

// Helper functions to convert between domain model and entity

func autoBuyRuleToEntity(rule *model.AutoBuyRule) *entity.AutoBuyRuleEntity {
	return &entity.AutoBuyRuleEntity{
		ID:                  rule.ID,
		UserID:              rule.UserID,
		Name:                rule.Name,
		Symbol:              rule.Symbol,
		IsEnabled:           rule.IsEnabled,
		TriggerType:         string(rule.TriggerType),
		TriggerValue:        rule.TriggerValue,
		QuoteAsset:          rule.QuoteAsset,
		BuyAmountQuote:      rule.BuyAmountQuote,
		MaxBuyPrice:         rule.MaxBuyPrice,
		MinBaseAssetVolume:  rule.MinBaseAssetVolume,
		MinQuoteAssetVolume: rule.MinQuoteAssetVolume,
		AllowPreTrading:     rule.AllowPreTrading,
		CooldownMinutes:     rule.CooldownMinutes,
		OrderType:           string(rule.OrderType),
		EnableRiskCheck:     rule.EnableRiskCheck,
		ExecutionCount:      rule.ExecutionCount,
		LastTriggered:       rule.LastTriggered,
		LastPrice:           rule.LastPrice,
		CreatedAt:           rule.CreatedAt,
		UpdatedAt:           rule.UpdatedAt,
	}
}

func autoBuyRuleFromEntity(entity *entity.AutoBuyRuleEntity) *model.AutoBuyRule {
	return &model.AutoBuyRule{
		ID:                  entity.ID,
		UserID:              entity.UserID,
		Name:                entity.Name,
		Symbol:              entity.Symbol,
		IsEnabled:           entity.IsEnabled,
		TriggerType:         model.TriggerType(entity.TriggerType),
		TriggerValue:        entity.TriggerValue,
		QuoteAsset:          entity.QuoteAsset,
		BuyAmountQuote:      entity.BuyAmountQuote,
		MaxBuyPrice:         entity.MaxBuyPrice,
		MinBaseAssetVolume:  entity.MinBaseAssetVolume,
		MinQuoteAssetVolume: entity.MinQuoteAssetVolume,
		AllowPreTrading:     entity.AllowPreTrading,
		CooldownMinutes:     entity.CooldownMinutes,
		OrderType:           model.OrderType(entity.OrderType),
		EnableRiskCheck:     entity.EnableRiskCheck,
		ExecutionCount:      entity.ExecutionCount,
		LastTriggered:       entity.LastTriggered,
		LastPrice:           entity.LastPrice,
		CreatedAt:           entity.CreatedAt,
		UpdatedAt:           entity.UpdatedAt,
	}
}

// AutoBuyExecutionRepository implements the port.AutoBuyExecutionRepository interface using GORM
type AutoBuyExecutionRepository struct {
	db     *gorm.DB
	logger zerolog.Logger
}

// NewAutoBuyExecutionRepository creates a new AutoBuyExecutionRepository
func NewAutoBuyExecutionRepository(db *gorm.DB, logger zerolog.Logger) *AutoBuyExecutionRepository {
	return &AutoBuyExecutionRepository{
		db:     db,
		logger: logger.With().Str("component", "autobuy_execution_repository").Logger(),
	}
}

// Create adds a new auto-buy execution record
func (r *AutoBuyExecutionRepository) Create(ctx context.Context, execution *model.AutoBuyExecution) error {
	entity := autoBuyExecutionToEntity(execution)
	result := r.db.WithContext(ctx).Create(&entity)
	if result.Error != nil {
		r.logger.Error().Err(result.Error).Str("ruleId", execution.RuleID).Msg("Failed to create auto-buy execution record")
		return result.Error
	}

	// Update the domain model with the generated ID if not already set
	if execution.ID == "" {
		execution.ID = entity.ID
	}

	r.logger.Info().Str("id", execution.ID).Str("ruleId", execution.RuleID).Str("orderId", execution.OrderID).Msg("Created auto-buy execution record")
	return nil
}

// GetByID retrieves an auto-buy execution by its ID
func (r *AutoBuyExecutionRepository) GetByID(ctx context.Context, id string) (*model.AutoBuyExecution, error) {
	var entity entity.AutoBuyExecutionEntity
	result := r.db.WithContext(ctx).First(&entity, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		r.logger.Error().Err(result.Error).Str("id", id).Msg("Failed to get auto-buy execution by ID")
		return nil, result.Error
	}

	execution := autoBuyExecutionFromEntity(&entity)
	return execution, nil
}

// GetByRuleID retrieves execution records for a specific rule
func (r *AutoBuyExecutionRepository) GetByRuleID(ctx context.Context, ruleID string, limit, offset int) ([]*model.AutoBuyExecution, error) {
	var entities []entity.AutoBuyExecutionEntity
	query := r.db.WithContext(ctx).Where("rule_id = ?", ruleID).Order("timestamp DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	result := query.Find(&entities)
	if result.Error != nil {
		r.logger.Error().Err(result.Error).Str("ruleId", ruleID).Msg("Failed to get auto-buy executions by rule ID")
		return nil, result.Error
	}

	executions := make([]*model.AutoBuyExecution, len(entities))
	for i, entity := range entities {
		executions[i] = autoBuyExecutionFromEntity(&entity)
	}

	r.logger.Debug().Str("ruleId", ruleID).Int("count", len(executions)).Msg("Retrieved auto-buy executions by rule ID")
	return executions, nil
}

// GetByUserID retrieves execution records for a specific user
func (r *AutoBuyExecutionRepository) GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*model.AutoBuyExecution, error) {
	var entities []entity.AutoBuyExecutionEntity
	query := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("timestamp DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	result := query.Find(&entities)
	if result.Error != nil {
		r.logger.Error().Err(result.Error).Str("userId", userID).Msg("Failed to get auto-buy executions by user ID")
		return nil, result.Error
	}

	executions := make([]*model.AutoBuyExecution, len(entities))
	for i, entity := range entities {
		executions[i] = autoBuyExecutionFromEntity(&entity)
	}

	r.logger.Debug().Str("userId", userID).Int("count", len(executions)).Msg("Retrieved auto-buy executions by user ID")
	return executions, nil
}

// GetBySymbol retrieves execution records for a specific symbol
func (r *AutoBuyExecutionRepository) GetBySymbol(ctx context.Context, symbol string, limit, offset int) ([]*model.AutoBuyExecution, error) {
	var entities []entity.AutoBuyExecutionEntity
	query := r.db.WithContext(ctx).Where("symbol = ?", symbol).Order("timestamp DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	result := query.Find(&entities)
	if result.Error != nil {
		r.logger.Error().Err(result.Error).Str("symbol", symbol).Msg("Failed to get auto-buy executions by symbol")
		return nil, result.Error
	}

	executions := make([]*model.AutoBuyExecution, len(entities))
	for i, entity := range entities {
		executions[i] = autoBuyExecutionFromEntity(&entity)
	}

	r.logger.Debug().Str("symbol", symbol).Int("count", len(executions)).Msg("Retrieved auto-buy executions by symbol")
	return executions, nil
}

// GetByTimeRange retrieves execution records within a time range
func (r *AutoBuyExecutionRepository) GetByTimeRange(ctx context.Context, from, to time.Time, limit, offset int) ([]*model.AutoBuyExecution, error) {
	var entities []entity.AutoBuyExecutionEntity
	query := r.db.WithContext(ctx).Where("timestamp BETWEEN ? AND ?", from, to).Order("timestamp DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	result := query.Find(&entities)
	if result.Error != nil {
		r.logger.Error().Err(result.Error).Time("from", from).Time("to", to).Msg("Failed to get auto-buy executions by time range")
		return nil, result.Error
	}

	executions := make([]*model.AutoBuyExecution, len(entities))
	for i, entity := range entities {
		executions[i] = autoBuyExecutionFromEntity(&entity)
	}

	r.logger.Debug().Time("from", from).Time("to", to).Int("count", len(executions)).Msg("Retrieved auto-buy executions by time range")
	return executions, nil
}

// Count returns the total number of execution records matching the specified filters
func (r *AutoBuyExecutionRepository) Count(ctx context.Context, filters map[string]interface{}) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&entity.AutoBuyExecutionEntity{})

	// Apply filters
	for key, value := range filters {
		query = query.Where(fmt.Sprintf("%s = ?", key), value)
	}

	result := query.Count(&count)
	if result.Error != nil {
		r.logger.Error().Err(result.Error).Msg("Failed to count auto-buy executions")
		return 0, result.Error
	}

	return count, nil
}

// Helper functions to convert between domain model and entity

func autoBuyExecutionToEntity(execution *model.AutoBuyExecution) *entity.AutoBuyExecutionEntity {
	return &entity.AutoBuyExecutionEntity{
		ID:        execution.ID,
		RuleID:    execution.RuleID,
		UserID:    execution.UserID,
		Symbol:    execution.Symbol,
		OrderID:   execution.OrderID,
		Price:     execution.Price,
		Quantity:  execution.Quantity,
		Amount:    execution.Amount,
		Timestamp: execution.Timestamp,
	}
}

func autoBuyExecutionFromEntity(entity *entity.AutoBuyExecutionEntity) *model.AutoBuyExecution {
	return &model.AutoBuyExecution{
		ID:        entity.ID,
		RuleID:    entity.RuleID,
		UserID:    entity.UserID,
		Symbol:    entity.Symbol,
		OrderID:   entity.OrderID,
		Price:     entity.Price,
		Quantity:  entity.Quantity,
		Amount:    entity.Amount,
		Timestamp: entity.Timestamp,
	}
}
