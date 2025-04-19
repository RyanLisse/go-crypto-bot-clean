package repo

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/adapter/infrastructure/persistence/gorm/entity"
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/port"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// SymbolRepository implements the port.SymbolRepository interface using GORM
type SymbolRepository struct {
	BaseRepository
}

// NewSymbolRepository creates a new SymbolRepository
func NewSymbolRepository(db *gorm.DB, logger *zerolog.Logger) port.SymbolRepository {
	return &SymbolRepository{
		BaseRepository: NewBaseRepository(db, logger),
	}
}

// GetBySymbol retrieves a symbol by its name
func (r *SymbolRepository) GetBySymbol(ctx context.Context, symbol string) (*model.Symbol, error) {
	var symbolEntity entity.MexcSymbolEntity
	result := r.GetDB(ctx).Where("symbol = ?", symbol).First(&symbolEntity)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			r.logger.Debug().Str("symbol", symbol).Msg("Symbol not found")
			return nil, nil
		}
		r.logger.Error().Err(result.Error).Str("symbol", symbol).Msg("Failed to get symbol")
		return nil, result.Error
	}
	
	return &model.Symbol{
		Symbol:            symbolEntity.Symbol,
		BaseAsset:         symbolEntity.BaseAsset,
		QuoteAsset:        symbolEntity.QuoteAsset,
		Exchange:          symbolEntity.Exchange,
		Status:            model.SymbolStatus(symbolEntity.Status),
		MinPrice:          symbolEntity.MinPrice,
		MaxPrice:          symbolEntity.MaxPrice,
		PricePrecision:    symbolEntity.PricePrecision,
		MinQuantity:       symbolEntity.MinQty,
		MaxQuantity:       symbolEntity.MaxQty,
		QuantityPrecision: symbolEntity.QtyPrecision,
		AllowedOrderTypes: strings.Split(symbolEntity.AllowedOrderTypes, ","),
		CreatedAt:         symbolEntity.CreatedAt,
		UpdatedAt:         symbolEntity.UpdatedAt,
	}, nil
}

// GetBySymbolAndExchange retrieves a symbol by its name and exchange
func (r *SymbolRepository) GetBySymbolAndExchange(ctx context.Context, symbol, exchange string) (*model.Symbol, error) {
	var symbolEntity entity.MexcSymbolEntity
	result := r.GetDB(ctx).Where("symbol = ? AND exchange = ?", symbol, exchange).First(&symbolEntity)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			r.logger.Debug().Str("symbol", symbol).Str("exchange", exchange).Msg("Symbol not found")
			return nil, nil
		}
		r.logger.Error().Err(result.Error).Str("symbol", symbol).Str("exchange", exchange).Msg("Failed to get symbol")
		return nil, result.Error
	}
	
	return &model.Symbol{
		Symbol:            symbolEntity.Symbol,
		BaseAsset:         symbolEntity.BaseAsset,
		QuoteAsset:        symbolEntity.QuoteAsset,
		Exchange:          symbolEntity.Exchange,
		Status:            model.SymbolStatus(symbolEntity.Status),
		MinPrice:          symbolEntity.MinPrice,
		MaxPrice:          symbolEntity.MaxPrice,
		PricePrecision:    symbolEntity.PricePrecision,
		MinQuantity:       symbolEntity.MinQty,
		MaxQuantity:       symbolEntity.MaxQty,
		QuantityPrecision: symbolEntity.QtyPrecision,
		AllowedOrderTypes: strings.Split(symbolEntity.AllowedOrderTypes, ","),
		CreatedAt:         symbolEntity.CreatedAt,
		UpdatedAt:         symbolEntity.UpdatedAt,
	}, nil
}

// ListAll retrieves all symbols with pagination
func (r *SymbolRepository) ListAll(ctx context.Context, limit, offset int) ([]*model.Symbol, error) {
	var symbolEntities []entity.MexcSymbolEntity
	query := r.GetDB(ctx)
	
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}
	
	result := query.Find(&symbolEntities)
	if result.Error != nil {
		r.logger.Error().Err(result.Error).Msg("Failed to list all symbols")
		return nil, result.Error
	}
	
	symbols := make([]*model.Symbol, len(symbolEntities))
	for i, entity := range symbolEntities {
		symbols[i] = &model.Symbol{
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
			AllowedOrderTypes: strings.Split(entity.AllowedOrderTypes, ","),
			CreatedAt:         entity.CreatedAt,
			UpdatedAt:         entity.UpdatedAt,
		}
	}
	
	return symbols, nil
}

// ListByExchange retrieves all symbols for a specific exchange
func (r *SymbolRepository) ListByExchange(ctx context.Context, exchange string, limit, offset int) ([]*model.Symbol, error) {
	var symbolEntities []entity.MexcSymbolEntity
	query := r.GetDB(ctx).Where("exchange = ?", exchange)
	
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}
	
	result := query.Find(&symbolEntities)
	if result.Error != nil {
		r.logger.Error().Err(result.Error).Str("exchange", exchange).Msg("Failed to list symbols by exchange")
		return nil, result.Error
	}
	
	symbols := make([]*model.Symbol, len(symbolEntities))
	for i, entity := range symbolEntities {
		symbols[i] = &model.Symbol{
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
			AllowedOrderTypes: strings.Split(entity.AllowedOrderTypes, ","),
			CreatedAt:         entity.CreatedAt,
			UpdatedAt:         entity.UpdatedAt,
		}
	}
	
	return symbols, nil
}

// ListByBaseAsset retrieves all symbols with a specific base asset
func (r *SymbolRepository) ListByBaseAsset(ctx context.Context, baseAsset string, limit, offset int) ([]*model.Symbol, error) {
	var symbolEntities []entity.MexcSymbolEntity
	query := r.GetDB(ctx).Where("base_asset = ?", baseAsset)
	
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}
	
	result := query.Find(&symbolEntities)
	if result.Error != nil {
		r.logger.Error().Err(result.Error).Str("baseAsset", baseAsset).Msg("Failed to list symbols by base asset")
		return nil, result.Error
	}
	
	symbols := make([]*model.Symbol, len(symbolEntities))
	for i, entity := range symbolEntities {
		symbols[i] = &model.Symbol{
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
			AllowedOrderTypes: strings.Split(entity.AllowedOrderTypes, ","),
			CreatedAt:         entity.CreatedAt,
			UpdatedAt:         entity.UpdatedAt,
		}
	}
	
	return symbols, nil
}

// ListByQuoteAsset retrieves all symbols with a specific quote asset
func (r *SymbolRepository) ListByQuoteAsset(ctx context.Context, quoteAsset string, limit, offset int) ([]*model.Symbol, error) {
	var symbolEntities []entity.MexcSymbolEntity
	query := r.GetDB(ctx).Where("quote_asset = ?", quoteAsset)
	
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}
	
	result := query.Find(&symbolEntities)
	if result.Error != nil {
		r.logger.Error().Err(result.Error).Str("quoteAsset", quoteAsset).Msg("Failed to list symbols by quote asset")
		return nil, result.Error
	}
	
	symbols := make([]*model.Symbol, len(symbolEntities))
	for i, entity := range symbolEntities {
		symbols[i] = &model.Symbol{
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
			AllowedOrderTypes: strings.Split(entity.AllowedOrderTypes, ","),
			CreatedAt:         entity.CreatedAt,
			UpdatedAt:         entity.UpdatedAt,
		}
	}
	
	return symbols, nil
}

// Save saves a symbol
func (r *SymbolRepository) Save(ctx context.Context, symbol *model.Symbol) error {
	symbolEntity := &entity.MexcSymbolEntity{
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
	}
	
	result := r.GetDB(ctx).Create(symbolEntity)
	if result.Error != nil {
		r.logger.Error().Err(result.Error).Str("symbol", symbol.Symbol).Str("exchange", symbol.Exchange).Msg("Failed to save symbol")
		return result.Error
	}
	
	r.logger.Debug().Str("symbol", symbol.Symbol).Str("exchange", symbol.Exchange).Msg("Symbol saved successfully")
	return nil
}

// SaveBatch saves multiple symbols in a batch operation
func (r *SymbolRepository) SaveBatch(ctx context.Context, symbols []*model.Symbol) error {
	if len(symbols) == 0 {
		return nil
	}
	
	symbolEntities := make([]entity.MexcSymbolEntity, len(symbols))
	for i, symbol := range symbols {
		symbolEntities[i] = entity.MexcSymbolEntity{
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
		}
	}
	
	result := r.GetDB(ctx).Create(&symbolEntities)
	if result.Error != nil {
		r.logger.Error().Err(result.Error).Int("count", len(symbols)).Msg("Failed to save symbols in batch")
		return result.Error
	}
	
	r.logger.Debug().Int("count", len(symbols)).Msg("Symbols saved successfully in batch")
	return nil
}

// Update updates a symbol
func (r *SymbolRepository) Update(ctx context.Context, symbol *model.Symbol) error {
	symbolEntity := &entity.MexcSymbolEntity{
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
	}
	
	result := r.GetDB(ctx).Where("symbol = ? AND exchange = ?", symbol.Symbol, symbol.Exchange).Updates(symbolEntity)
	if result.Error != nil {
		r.logger.Error().Err(result.Error).Str("symbol", symbol.Symbol).Str("exchange", symbol.Exchange).Msg("Failed to update symbol")
		return result.Error
	}
	
	if result.RowsAffected == 0 {
		r.logger.Warn().Str("symbol", symbol.Symbol).Str("exchange", symbol.Exchange).Msg("No symbol found to update")
		return fmt.Errorf("symbol not found: %s on %s", symbol.Symbol, symbol.Exchange)
	}
	
	r.logger.Debug().Str("symbol", symbol.Symbol).Str("exchange", symbol.Exchange).Msg("Symbol updated successfully")
	return nil
}

// UpdateStatus updates the status of a symbol
func (r *SymbolRepository) UpdateStatus(ctx context.Context, symbol, exchange, status string) error {
	result := r.GetDB(ctx).
		Model(&entity.MexcSymbolEntity{}).
		Where("symbol = ? AND exchange = ?", symbol, exchange).
		Update("status", status)
	
	if result.Error != nil {
		r.logger.Error().Err(result.Error).Str("symbol", symbol).Str("exchange", exchange).Str("status", status).Msg("Failed to update symbol status")
		return result.Error
	}
	
	if result.RowsAffected == 0 {
		r.logger.Warn().Str("symbol", symbol).Str("exchange", exchange).Msg("No symbol found to update status")
		return fmt.Errorf("symbol not found: %s on %s", symbol, exchange)
	}
	
	r.logger.Debug().Str("symbol", symbol).Str("exchange", exchange).Str("status", status).Msg("Symbol status updated successfully")
	return nil
}

// Delete deletes a symbol
func (r *SymbolRepository) Delete(ctx context.Context, symbol, exchange string) error {
	result := r.GetDB(ctx).
		Where("symbol = ? AND exchange = ?", symbol, exchange).
		Delete(&entity.MexcSymbolEntity{})
	
	if result.Error != nil {
		r.logger.Error().Err(result.Error).Str("symbol", symbol).Str("exchange", exchange).Msg("Failed to delete symbol")
		return result.Error
	}
	
	if result.RowsAffected == 0 {
		r.logger.Warn().Str("symbol", symbol).Str("exchange", exchange).Msg("No symbol found to delete")
		return fmt.Errorf("symbol not found: %s on %s", symbol, exchange)
	}
	
	r.logger.Debug().Str("symbol", symbol).Str("exchange", exchange).Msg("Symbol deleted successfully")
	return nil
}

// GetLastUpdated retrieves the timestamp of the last symbol update for an exchange
func (r *SymbolRepository) GetLastUpdated(ctx context.Context, exchange string) (time.Time, error) {
	var syncState entity.MexcSyncStateEntity
	result := r.GetDB(ctx).Where("exchange = ? AND entity_type = ?", exchange, "symbols").First(&syncState)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return time.Time{}, nil
		}
		r.logger.Error().Err(result.Error).Str("exchange", exchange).Msg("Failed to get last updated timestamp")
		return time.Time{}, result.Error
	}
	
	return syncState.LastSyncTime, nil
}

// SetLastUpdated sets the timestamp of the last symbol update for an exchange
func (r *SymbolRepository) SetLastUpdated(ctx context.Context, exchange string, timestamp time.Time) error {
	syncState := &entity.MexcSyncStateEntity{
		Exchange:     exchange,
		EntityType:   "symbols",
		LastSyncTime: timestamp,
	}
	
	result := r.GetDB(ctx).Where("exchange = ? AND entity_type = ?", exchange, "symbols").
		Assign(syncState).
		FirstOrCreate(syncState)
	
	if result.Error != nil {
		r.logger.Error().Err(result.Error).Str("exchange", exchange).Time("timestamp", timestamp).Msg("Failed to set last updated timestamp")
		return result.Error
	}
	
	r.logger.Debug().Str("exchange", exchange).Time("timestamp", timestamp).Msg("Last updated timestamp set successfully")
	return nil
}

// Search searches for symbols by name, base asset, or quote asset
func (r *SymbolRepository) Search(ctx context.Context, query string, limit, offset int) ([]*model.Symbol, error) {
	var symbolEntities []entity.MexcSymbolEntity
	
	// Create a query with LIKE conditions for symbol, base_asset, and quote_asset
	dbQuery := r.GetDB(ctx).
		Where("symbol LIKE ? OR base_asset LIKE ? OR quote_asset LIKE ?",
			"%"+query+"%", "%"+query+"%", "%"+query+"%")
	
	if limit > 0 {
		dbQuery = dbQuery.Limit(limit)
	}
	if offset > 0 {
		dbQuery = dbQuery.Offset(offset)
	}
	
	result := dbQuery.Find(&symbolEntities)
	if result.Error != nil {
		r.logger.Error().Err(result.Error).Str("query", query).Msg("Failed to search symbols")
		return nil, result.Error
	}
	
	symbols := make([]*model.Symbol, len(symbolEntities))
	for i, entity := range symbolEntities {
		symbols[i] = &model.Symbol{
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
			AllowedOrderTypes: strings.Split(entity.AllowedOrderTypes, ","),
			CreatedAt:         entity.CreatedAt,
			UpdatedAt:         entity.UpdatedAt,
		}
	}
	
	return symbols, nil
}

// Count returns the total number of symbols
func (r *SymbolRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	result := r.GetDB(ctx).Model(&entity.MexcSymbolEntity{}).Count(&count)
	if result.Error != nil {
		r.logger.Error().Err(result.Error).Msg("Failed to count symbols")
		return 0, result.Error
	}
	
	return count, nil
}

// CountByExchange returns the total number of symbols for a specific exchange
func (r *SymbolRepository) CountByExchange(ctx context.Context, exchange string) (int64, error) {
	var count int64
	result := r.GetDB(ctx).Model(&entity.MexcSymbolEntity{}).Where("exchange = ?", exchange).Count(&count)
	if result.Error != nil {
		r.logger.Error().Err(result.Error).Str("exchange", exchange).Msg("Failed to count symbols by exchange")
		return 0, result.Error
	}
	
	return count, nil
}

// CountByBaseAsset returns the total number of symbols with a specific base asset
func (r *SymbolRepository) CountByBaseAsset(ctx context.Context, baseAsset string) (int64, error) {
	var count int64
	result := r.GetDB(ctx).Model(&entity.MexcSymbolEntity{}).Where("base_asset = ?", baseAsset).Count(&count)
	if result.Error != nil {
		r.logger.Error().Err(result.Error).Str("baseAsset", baseAsset).Msg("Failed to count symbols by base asset")
		return 0, result.Error
	}
	
	return count, nil
}

// CountByQuoteAsset returns the total number of symbols with a specific quote asset
func (r *SymbolRepository) CountByQuoteAsset(ctx context.Context, quoteAsset string) (int64, error) {
	var count int64
	result := r.GetDB(ctx).Model(&entity.MexcSymbolEntity{}).Where("quote_asset = ?", quoteAsset).Count(&count)
	if result.Error != nil {
		r.logger.Error().Err(result.Error).Str("quoteAsset", quoteAsset).Msg("Failed to count symbols by quote asset")
		return 0, result.Error
	}
	
	return count, nil
}

// Ensure SymbolRepository implements port.SymbolRepository
var _ port.SymbolRepository = (*SymbolRepository)(nil)
