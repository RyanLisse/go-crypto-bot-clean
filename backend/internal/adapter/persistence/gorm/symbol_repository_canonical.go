package gorm

import (
	"context"
	"fmt"
	"strings"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/apperror"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model/market"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// SymbolRepositoryCanonical implements the port.SymbolRepository interface using GORM
// with the canonical model types
type SymbolRepositoryCanonical struct {
	db     *gorm.DB
	logger *zerolog.Logger
}

// NewSymbolRepositoryCanonical creates a new SymbolRepositoryCanonical
func NewSymbolRepositoryCanonical(db *gorm.DB, logger *zerolog.Logger) port.SymbolRepository {
	return &SymbolRepositoryCanonical{
		db:     db,
		logger: logger,
	}
}

// Create stores a new Symbol
func (r *SymbolRepositoryCanonical) Create(ctx context.Context, symbol *model.Symbol) error {
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
func (r *SymbolRepositoryCanonical) GetSymbolsByStatus(ctx context.Context, status string, limit int, offset int) ([]*model.Symbol, error) {
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
func (r *SymbolRepositoryCanonical) GetBySymbol(ctx context.Context, symbol string) (*model.Symbol, error) {
	var entity SymbolEntity

	result := r.db.WithContext(ctx).Where("symbol = ?", symbol).First(&entity)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			r.logger.Info().Str("symbol", symbol).Msg("Symbol not found")
			return nil, apperror.ErrNotFound
		}
		r.logger.Error().Err(result.Error).Str("symbol", symbol).Msg("Failed to get symbol")
		return nil, fmt.Errorf("failed to get symbol: %w", result.Error)
	}

	return r.symbolToDomain(&entity), nil
}

// GetByExchange returns all Symbols from a specific exchange
func (r *SymbolRepositoryCanonical) GetByExchange(ctx context.Context, exchange string) ([]*model.Symbol, error) {
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
func (r *SymbolRepositoryCanonical) GetAll(ctx context.Context) ([]*model.Symbol, error) {
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
func (r *SymbolRepositoryCanonical) Update(ctx context.Context, symbol *model.Symbol) error {
	entity := r.symbolToEntity(symbol)

	result := r.db.WithContext(ctx).Where("symbol = ?", symbol.Symbol).Updates(&entity)
	if result.Error != nil {
		r.logger.Error().Err(result.Error).Str("symbol", symbol.Symbol).Msg("Failed to update symbol")
		return fmt.Errorf("failed to update symbol: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		r.logger.Warn().Str("symbol", symbol.Symbol).Msg("No symbol found to update")
		return apperror.ErrNotFound
	}

	r.logger.Info().Str("symbol", symbol.Symbol).Msg("Symbol updated successfully")
	return nil
}

// Delete removes a Symbol
func (r *SymbolRepositoryCanonical) Delete(ctx context.Context, symbol string) error {
	result := r.db.WithContext(ctx).Where("symbol = ?", symbol).Delete(&SymbolEntity{})
	if result.Error != nil {
		r.logger.Error().Err(result.Error).Str("symbol", symbol).Msg("Failed to delete symbol")
		return fmt.Errorf("failed to delete symbol: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		r.logger.Warn().Str("symbol", symbol).Msg("No symbol found to delete")
		return apperror.ErrNotFound
	}

	r.logger.Info().Str("symbol", symbol).Msg("Symbol deleted successfully")
	return nil
}

// Helper methods for entity conversion
func (r *SymbolRepositoryCanonical) symbolToEntity(symbol *model.Symbol) *SymbolEntity {
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

func (r *SymbolRepositoryCanonical) symbolToDomain(entity *SymbolEntity) *model.Symbol {
	var allowedOrderTypes []string
	if entity.AllowedOrderTypes != "" {
		allowedOrderTypes = strings.Split(entity.AllowedOrderTypes, ",")
	}

	// Convert string status to model.SymbolStatus
	var status model.SymbolStatus = model.SymbolStatusHalt
	if entity.Status == string(model.SymbolStatusTrading) {
		status = model.SymbolStatusTrading
	} else if entity.Status == string(model.SymbolStatusBreak) {
		status = model.SymbolStatusBreak
	}

	return &model.Symbol{
		Symbol:            entity.Symbol,
		BaseAsset:         entity.BaseAsset,
		QuoteAsset:        entity.QuoteAsset,
		Exchange:          entity.Exchange,
		Status:            status,
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

// Legacy methods for backward compatibility

// CreateLegacy stores a new Symbol using the legacy model
func (r *SymbolRepositoryCanonical) CreateLegacy(ctx context.Context, symbol *market.Symbol) error {
	// Convert legacy model to canonical model
	var status model.SymbolStatus = model.SymbolStatusHalt
	if symbol.Status == string(model.SymbolStatusTrading) {
		status = model.SymbolStatusTrading
	} else if symbol.Status == string(model.SymbolStatusBreak) {
		status = model.SymbolStatusBreak
	}

	canonicalSymbol := &model.Symbol{
		Symbol:            symbol.Symbol,
		BaseAsset:         symbol.BaseAsset,
		QuoteAsset:        symbol.QuoteAsset,
		Exchange:          symbol.Exchange,
		Status:            status,
		MinPrice:          symbol.MinPrice,
		MaxPrice:          symbol.MaxPrice,
		PricePrecision:    symbol.PricePrecision,
		MinQuantity:       symbol.MinQty,
		MaxQuantity:       symbol.MaxQty,
		QuantityPrecision: symbol.QtyPrecision,
		AllowedOrderTypes: symbol.AllowedOrderTypes,
		CreatedAt:         symbol.CreatedAt,
		UpdatedAt:         symbol.UpdatedAt,
	}

	// Use the canonical implementation
	return r.Create(ctx, canonicalSymbol)
}

// GetBySymbolLegacy returns a Symbol by its symbol string using the legacy model
func (r *SymbolRepositoryCanonical) GetBySymbolLegacy(ctx context.Context, symbol string) (*market.Symbol, error) {
	// Use the canonical implementation
	canonicalSymbol, err := r.GetBySymbol(ctx, symbol)
	if err != nil {
		return nil, err
	}

	// Convert canonical model to legacy model
	return &market.Symbol{
		Symbol:            canonicalSymbol.Symbol,
		BaseAsset:         canonicalSymbol.BaseAsset,
		QuoteAsset:        canonicalSymbol.QuoteAsset,
		Exchange:          canonicalSymbol.Exchange,
		Status:            string(canonicalSymbol.Status),
		MinPrice:          canonicalSymbol.MinPrice,
		MaxPrice:          canonicalSymbol.MaxPrice,
		PricePrecision:    canonicalSymbol.PricePrecision,
		MinQty:            canonicalSymbol.MinQuantity,
		MaxQty:            canonicalSymbol.MaxQuantity,
		QtyPrecision:      canonicalSymbol.QuantityPrecision,
		AllowedOrderTypes: canonicalSymbol.AllowedOrderTypes,
		CreatedAt:         canonicalSymbol.CreatedAt,
		UpdatedAt:         canonicalSymbol.UpdatedAt,
	}, nil
}

// GetByExchangeLegacy returns all Symbols from a specific exchange using the legacy model
func (r *SymbolRepositoryCanonical) GetByExchangeLegacy(ctx context.Context, exchange string) ([]*market.Symbol, error) {
	// Use the canonical implementation
	canonicalSymbols, err := r.GetByExchange(ctx, exchange)
	if err != nil {
		return nil, err
	}

	// Convert canonical models to legacy models
	legacySymbols := make([]*market.Symbol, len(canonicalSymbols))
	for i, canonicalSymbol := range canonicalSymbols {
		legacySymbols[i] = &market.Symbol{
			Symbol:            canonicalSymbol.Symbol,
			BaseAsset:         canonicalSymbol.BaseAsset,
			QuoteAsset:        canonicalSymbol.QuoteAsset,
			Exchange:          canonicalSymbol.Exchange,
			Status:            string(canonicalSymbol.Status),
			MinPrice:          canonicalSymbol.MinPrice,
			MaxPrice:          canonicalSymbol.MaxPrice,
			PricePrecision:    canonicalSymbol.PricePrecision,
			MinQty:            canonicalSymbol.MinQuantity,
			MaxQty:            canonicalSymbol.MaxQuantity,
			QtyPrecision:      canonicalSymbol.QuantityPrecision,
			AllowedOrderTypes: canonicalSymbol.AllowedOrderTypes,
			CreatedAt:         canonicalSymbol.CreatedAt,
			UpdatedAt:         canonicalSymbol.UpdatedAt,
		}
	}

	return legacySymbols, nil
}

// GetAllLegacy returns all available Symbols using the legacy model
func (r *SymbolRepositoryCanonical) GetAllLegacy(ctx context.Context) ([]*market.Symbol, error) {
	// Use the canonical implementation
	canonicalSymbols, err := r.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	// Convert canonical models to legacy models
	legacySymbols := make([]*market.Symbol, len(canonicalSymbols))
	for i, canonicalSymbol := range canonicalSymbols {
		legacySymbols[i] = &market.Symbol{
			Symbol:            canonicalSymbol.Symbol,
			BaseAsset:         canonicalSymbol.BaseAsset,
			QuoteAsset:        canonicalSymbol.QuoteAsset,
			Exchange:          canonicalSymbol.Exchange,
			Status:            string(canonicalSymbol.Status),
			MinPrice:          canonicalSymbol.MinPrice,
			MaxPrice:          canonicalSymbol.MaxPrice,
			PricePrecision:    canonicalSymbol.PricePrecision,
			MinQty:            canonicalSymbol.MinQuantity,
			MaxQty:            canonicalSymbol.MaxQuantity,
			QtyPrecision:      canonicalSymbol.QuantityPrecision,
			AllowedOrderTypes: canonicalSymbol.AllowedOrderTypes,
			CreatedAt:         canonicalSymbol.CreatedAt,
			UpdatedAt:         canonicalSymbol.UpdatedAt,
		}
	}

	return legacySymbols, nil
}

// UpdateLegacy updates an existing Symbol using the legacy model
func (r *SymbolRepositoryCanonical) UpdateLegacy(ctx context.Context, symbol *market.Symbol) error {
	// Convert legacy model to canonical model
	var status model.SymbolStatus = model.SymbolStatusHalt
	if symbol.Status == string(model.SymbolStatusTrading) {
		status = model.SymbolStatusTrading
	} else if symbol.Status == string(model.SymbolStatusBreak) {
		status = model.SymbolStatusBreak
	}

	canonicalSymbol := &model.Symbol{
		Symbol:            symbol.Symbol,
		BaseAsset:         symbol.BaseAsset,
		QuoteAsset:        symbol.QuoteAsset,
		Exchange:          symbol.Exchange,
		Status:            status,
		MinPrice:          symbol.MinPrice,
		MaxPrice:          symbol.MaxPrice,
		PricePrecision:    symbol.PricePrecision,
		MinQuantity:       symbol.MinQty,
		MaxQuantity:       symbol.MaxQty,
		QuantityPrecision: symbol.QtyPrecision,
		AllowedOrderTypes: symbol.AllowedOrderTypes,
		CreatedAt:         symbol.CreatedAt,
		UpdatedAt:         symbol.UpdatedAt,
	}

	// Use the canonical implementation
	return r.Update(ctx, canonicalSymbol)
}

// GetSymbolsByStatusLegacy returns symbols by status with pagination using the legacy model
func (r *SymbolRepositoryCanonical) GetSymbolsByStatusLegacy(ctx context.Context, status string, limit int, offset int) ([]*market.Symbol, error) {
	// Use the canonical implementation
	canonicalSymbols, err := r.GetSymbolsByStatus(ctx, status, limit, offset)
	if err != nil {
		return nil, err
	}

	// Convert canonical models to legacy models
	legacySymbols := make([]*market.Symbol, len(canonicalSymbols))
	for i, canonicalSymbol := range canonicalSymbols {
		legacySymbols[i] = &market.Symbol{
			Symbol:            canonicalSymbol.Symbol,
			BaseAsset:         canonicalSymbol.BaseAsset,
			QuoteAsset:        canonicalSymbol.QuoteAsset,
			Exchange:          canonicalSymbol.Exchange,
			Status:            string(canonicalSymbol.Status),
			MinPrice:          canonicalSymbol.MinPrice,
			MaxPrice:          canonicalSymbol.MaxPrice,
			PricePrecision:    canonicalSymbol.PricePrecision,
			MinQty:            canonicalSymbol.MinQuantity,
			MaxQty:            canonicalSymbol.MaxQuantity,
			QtyPrecision:      canonicalSymbol.QuantityPrecision,
			AllowedOrderTypes: canonicalSymbol.AllowedOrderTypes,
			CreatedAt:         canonicalSymbol.CreatedAt,
			UpdatedAt:         canonicalSymbol.UpdatedAt,
		}
	}

	return legacySymbols, nil
}
