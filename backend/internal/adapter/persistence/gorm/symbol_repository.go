package gorm

import (
	"context"
	"fmt"
	"strings"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// SymbolRepository implements the port.SymbolRepository interface using GORM
type SymbolRepository struct {
	db     *gorm.DB
	logger *zerolog.Logger
}

// NewSymbolRepository creates a new SymbolRepository
func NewSymbolRepository(db *gorm.DB, logger *zerolog.Logger) port.SymbolRepository {
	return &SymbolRepository{
		db:     db,
		logger: logger,
	}
}

// Create stores a new Symbol
func (r *SymbolRepository) Create(ctx context.Context, symbol *model.Symbol) error {
	entity := r.symbolToEntity(symbol)

	result := r.db.WithContext(ctx).Create(&entity)
	if result.Error != nil {
		r.logger.Error().Err(result.Error).Str("symbol", symbol.Symbol).Msg("Failed to create symbol")
		return fmt.Errorf("failed to create symbol: %w", result.Error)
	}

	r.logger.Info().Str("symbol", symbol.Symbol).Str("exchange", symbol.Exchange).Msg("Symbol created successfully")
	return nil
}

// GetSymbolsByStatus returns symbols by status with pagination
func (r *SymbolRepository) GetSymbolsByStatus(ctx context.Context, status string, limit int, offset int) ([]*model.Symbol, error) {
	var entities []SymbolEntity
	result := r.db.WithContext(ctx).Where("status = ?", status).Limit(limit).Offset(offset).Find(&entities)
	if result.Error != nil {
		r.logger.Error().Err(result.Error).Str("status", status).Msg("Failed to get symbols by status")
		return nil, fmt.Errorf("failed to get symbols by status: %w", result.Error)
	}
	symbols := make([]*model.Symbol, 0, len(entities))
	for _, entity := range entities {
		symbols = append(symbols, r.symbolToDomain(&entity))
	}
	return symbols, nil
}

// GetBySymbol returns a Symbol by its symbol string (e.g., "BTCUSDT")
func (r *SymbolRepository) GetBySymbol(ctx context.Context, symbol string) (*model.Symbol, error) {
	var entity SymbolEntity

	result := r.db.WithContext(ctx).Where("symbol = ?", symbol).First(&entity)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			r.logger.Info().Str("symbol", symbol).Msg("Symbol not found")
			return nil, fmt.Errorf("symbol not found: %s", symbol)
		}
		r.logger.Error().Err(result.Error).Str("symbol", symbol).Msg("Failed to get symbol")
		return nil, fmt.Errorf("failed to get symbol: %w", result.Error)
	}

	return r.symbolToDomain(&entity), nil
}

// GetByExchange returns all Symbols from a specific exchange
func (r *SymbolRepository) GetByExchange(ctx context.Context, exchange string) ([]*model.Symbol, error) {
	var entities []SymbolEntity

	result := r.db.WithContext(ctx).Where("exchange = ?", exchange).Find(&entities)
	if result.Error != nil {
		r.logger.Error().Err(result.Error).Str("exchange", exchange).Msg("Failed to get symbols by exchange")
		return nil, fmt.Errorf("failed to get symbols by exchange: %w", result.Error)
	}

	symbols := make([]*model.Symbol, len(entities))
	for i, entity := range entities {
		symbols[i] = r.symbolToDomain(&entity)
	}

	r.logger.Info().Str("exchange", exchange).Int("count", len(symbols)).Msg("Retrieved symbols by exchange")
	return symbols, nil
}

// GetAll returns all available Symbols
func (r *SymbolRepository) GetAll(ctx context.Context) ([]*model.Symbol, error) {
	var entities []SymbolEntity

	result := r.db.WithContext(ctx).Find(&entities)
	if result.Error != nil {
		r.logger.Error().Err(result.Error).Msg("Failed to get all symbols")
		return nil, fmt.Errorf("failed to get all symbols: %w", result.Error)
	}

	symbols := make([]*model.Symbol, len(entities))
	for i, entity := range entities {
		symbols[i] = r.symbolToDomain(&entity)
	}

	r.logger.Info().Int("count", len(symbols)).Msg("Retrieved all symbols")
	return symbols, nil
}

// Update updates an existing Symbol
func (r *SymbolRepository) Update(ctx context.Context, symbol *model.Symbol) error {
	entity := r.symbolToEntity(symbol)

	result := r.db.WithContext(ctx).Where("symbol = ?", symbol.Symbol).Updates(&entity)
	if result.Error != nil {
		r.logger.Error().Err(result.Error).Str("symbol", symbol.Symbol).Msg("Failed to update symbol")
		return fmt.Errorf("failed to update symbol: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		r.logger.Warn().Str("symbol", symbol.Symbol).Msg("No symbol found to update")
		return fmt.Errorf("symbol not found: %s", symbol.Symbol)
	}

	r.logger.Info().Str("symbol", symbol.Symbol).Msg("Symbol updated successfully")
	return nil
}

// Delete removes a Symbol
func (r *SymbolRepository) Delete(ctx context.Context, symbol string) error {
	result := r.db.WithContext(ctx).Where("symbol = ?", symbol).Delete(&SymbolEntity{})
	if result.Error != nil {
		r.logger.Error().Err(result.Error).Str("symbol", symbol).Msg("Failed to delete symbol")
		return fmt.Errorf("failed to delete symbol: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		r.logger.Warn().Str("symbol", symbol).Msg("No symbol found to delete")
		return fmt.Errorf("symbol not found: %s", symbol)
	}

	r.logger.Info().Str("symbol", symbol).Msg("Symbol deleted successfully")
	return nil
}

// Helper methods for entity conversion
func (r *SymbolRepository) symbolToEntity(symbol *model.Symbol) *SymbolEntity {
	return &SymbolEntity{
		Symbol:            symbol.Symbol,
		BaseAsset:         symbol.BaseAsset,
		QuoteAsset:        symbol.QuoteAsset,
		Exchange:          symbol.Exchange,
		Status:            string(symbol.Status),
		MinPrice:          symbol.MinPrice,
		MaxPrice:          symbol.MaxPrice,
		PricePrecision:    symbol.PricePrecision,
		MinQty:            symbol.MinQuantity,
		MaxQty:            symbol.MaxQuantity,
		QtyPrecision:      symbol.QuantityPrecision,
		AllowedOrderTypes: strings.Join(symbol.AllowedOrderTypes, ","),
		CreatedAt:         symbol.CreatedAt,
		UpdatedAt:         symbol.UpdatedAt,
	}
}

func (r *SymbolRepository) symbolToDomain(entity *SymbolEntity) *model.Symbol {
	var allowedOrderTypes []string
	if entity.AllowedOrderTypes != "" {
		allowedOrderTypes = strings.Split(entity.AllowedOrderTypes, ",")
	}

	return &model.Symbol{
		Symbol:            entity.Symbol,
		BaseAsset:         entity.BaseAsset,
		QuoteAsset:        entity.QuoteAsset,
		Exchange:          entity.Exchange,
		Status:            model.SymbolStatus(entity.Status),
		MinPrice:          entity.MinPrice,
		MaxPrice:          entity.MaxPrice,
		PricePrecision:    entity.PricePrecision,
		MinQuantity:       entity.MinQty,
		MaxQuantity:       entity.MaxQty,
		QuantityPrecision: entity.QtyPrecision,
		AllowedOrderTypes: allowedOrderTypes,
		CreatedAt:         entity.CreatedAt,
		UpdatedAt:         entity.UpdatedAt,
	}
}
