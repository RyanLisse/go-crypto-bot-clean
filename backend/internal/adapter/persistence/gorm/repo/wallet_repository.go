package repo

import (
	"context"
	"encoding/json"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// WalletEntity represents a wallet in the database
type WalletEntity struct {
	ID         string    `gorm:"primaryKey;type:varchar(50)"`
	UserID     string    `gorm:"uniqueIndex;type:varchar(50)"`
	Exchange   string    `gorm:"index;type:varchar(20)"`
	Balances   []byte    `gorm:"type:json"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`
	LastSyncAt time.Time
}

// BalanceHistoryEntity represents a balance history record in the database
type BalanceHistoryEntity struct {
	ID            string    `gorm:"primaryKey;type:varchar(50)"`
	UserID        string    `gorm:"index;type:varchar(50)"`
	WalletID      string    `gorm:"index;type:varchar(50)"`
	BalancesJSON  []byte    `gorm:"type:json"`
	TotalUSDValue float64   `gorm:"type:decimal(18,8);not null;default:0"`
	Timestamp     time.Time `gorm:"index"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
}

// GormWalletRepository implements port.WalletRepository using GORM
type GormWalletRepository struct {
	BaseRepository
}

// These are placeholder implementations to satisfy the interface
// In practice, we should use the ConsolidatedWalletRepository instead

// NewGormWalletRepository creates a new GormWalletRepository
func NewGormWalletRepository(db *gorm.DB, logger *zerolog.Logger) *GormWalletRepository {
	return &GormWalletRepository{
		BaseRepository: NewBaseRepository(db, logger),
	}
}

// Save saves a wallet to the database
func (r *GormWalletRepository) Save(ctx context.Context, wallet *model.Wallet) error {
	// Convert balances to JSON
	balancesJSON, err := json.Marshal(wallet.Balances)
	if err != nil {
		r.logger.Error().Err(err).Msg("Failed to marshal wallet balances")
		return err
	}

	// Create entity
	entity := &WalletEntity{
		ID:         wallet.UserID + "_" + wallet.Exchange,
		UserID:     wallet.UserID,
		Exchange:   wallet.Exchange,
		Balances:   balancesJSON,
		LastSyncAt: time.Now(),
	}

	// Save entity
	return r.Upsert(ctx, entity, []string{"user_id", "exchange"}, []string{
		"balances", "updated_at", "last_sync_at",
	})
}

// GetByUserID retrieves a wallet by user ID
func (r *GormWalletRepository) GetByUserID(ctx context.Context, userID string) (*model.Wallet, error) {
	var entity WalletEntity
	err := r.FindOne(ctx, &entity, "user_id = ?", userID)
	if err != nil {
		return nil, err
	}

	if entity.ID == "" {
		return nil, nil // Not found
	}

	return r.toDomain(&entity), nil
}

// SaveBalanceHistory saves a balance history record
func (r *GormWalletRepository) SaveBalanceHistory(ctx context.Context, history *model.BalanceHistory) error {
	// Create balances JSON
	balancesJSON, err := json.Marshal(history.Balances)
	if err != nil {
		r.logger.Error().Err(err).Str("userID", history.UserID).Msg("Failed to marshal balances")
		return err
	}

	// Create entity
	entity := &BalanceHistoryEntity{
		ID:            uuid.New().String(),
		UserID:        history.UserID,
		WalletID:      history.WalletID,
		BalancesJSON:  balancesJSON,
		TotalUSDValue: history.TotalUSDValue,
		Timestamp:     history.Timestamp,
	}

	// Save entity
	return r.Create(ctx, entity)
}

// GetBalanceHistory retrieves balance history for a user and asset within a time range
func (r *GormWalletRepository) GetBalanceHistory(ctx context.Context, userID string, asset model.Asset, from, to time.Time) ([]*model.BalanceHistory, error) {
	var entities []BalanceHistoryEntity

	query := r.GetDB(ctx).
		Where("user_id = ?", userID)

	// Add asset filter if provided
	if asset != "" {
		query = query.Where("asset = ?", asset)
	}

	// Add time range conditions
	if !from.IsZero() {
		query = query.Where("timestamp >= ?", from)
	}
	if !to.IsZero() {
		query = query.Where("timestamp <= ?", to)
	}

	// Execute query
	err := query.
		Order("timestamp ASC").
		Find(&entities).Error
	if err != nil {
		return nil, err
	}

	return r.historyToDomainSlice(entities), nil
}

// Helper methods for entity conversion

// toDomain converts a database entity to a domain wallet
func (r *GormWalletRepository) toDomain(entity *WalletEntity) *model.Wallet {
	if entity == nil {
		return nil
	}

	// Parse balances
	balanceMap := make(map[model.Asset]*model.Balance)
	if len(entity.Balances) > 0 {
		var balances []model.Balance
		if err := json.Unmarshal(entity.Balances, &balances); err != nil {
			r.logger.Error().Err(err).Msg("Failed to unmarshal wallet balances")
			balances = []model.Balance{}
		}

		// Convert slice to map
		for i := range balances {
			balanceMap[balances[i].Asset] = &balances[i]
		}
	}

	return &model.Wallet{
		UserID:        entity.UserID,
		Exchange:      entity.Exchange,
		Balances:      balanceMap,
		LastSyncAt:    entity.LastSyncAt,
		LastUpdated:   entity.UpdatedAt,
		TotalUSDValue: calculateTotalUSDValue(balanceMap),
	}
}

// calculateTotalUSDValue calculates the total USD value of all balances
func calculateTotalUSDValue(balances map[model.Asset]*model.Balance) float64 {
	total := 0.0
	for _, balance := range balances {
		total += balance.USDValue
	}
	return total
}

// historyToDomain converts a database entity to a domain balance history
func (r *GormWalletRepository) historyToDomain(entity *BalanceHistoryEntity) *model.BalanceHistory {
	if entity == nil {
		return nil
	}

	// Unmarshal balances JSON
	var balancesMap map[string]float64
	if err := json.Unmarshal(entity.BalancesJSON, &balancesMap); err != nil {
		r.logger.Error().Err(err).Str("id", entity.ID).Msg("Failed to unmarshal balances JSON")
		balancesMap = make(map[string]float64)
	}

	// Convert to map[model.Asset]*model.Balance
	balances := make(map[model.Asset]*model.Balance)
	for assetStr, value := range balancesMap {
		assetObj := model.Asset(assetStr)
		balances[assetObj] = &model.Balance{
			Asset: assetObj,
			Free:  value,
			Total: value,
		}
	}

	return &model.BalanceHistory{
		ID:            entity.ID,
		UserID:        entity.UserID,
		WalletID:      entity.WalletID,
		Balances:      balances,
		TotalUSDValue: entity.TotalUSDValue,
		Timestamp:     entity.Timestamp,
	}
}

// historyToDomainSlice converts a slice of database entities to domain balance histories
func (r *GormWalletRepository) historyToDomainSlice(entities []BalanceHistoryEntity) []*model.BalanceHistory {
	histories := make([]*model.BalanceHistory, len(entities))
	for i, entity := range entities {
		histories[i] = r.historyToDomain(&entity)
	}
	return histories
}

// GetByID retrieves a wallet by ID
func (r *GormWalletRepository) GetByID(ctx context.Context, id string) (*model.Wallet, error) {
	r.logger.Warn().Str("id", id).Msg("GetByID called on legacy GormWalletRepository, consider using ConsolidatedWalletRepository")
	// This is a placeholder implementation
	return nil, nil
}

// GetWalletsByUserID retrieves all wallets for a user
func (r *GormWalletRepository) GetWalletsByUserID(ctx context.Context, userID string) ([]*model.Wallet, error) {
	r.logger.Warn().Str("userID", userID).Msg("GetWalletsByUserID called on legacy GormWalletRepository, consider using ConsolidatedWalletRepository")
	// This is a placeholder implementation
	return nil, nil
}

// DeleteWallet deletes a wallet
func (r *GormWalletRepository) DeleteWallet(ctx context.Context, id string) error {
	r.logger.Warn().Str("id", id).Msg("DeleteWallet called on legacy GormWalletRepository, consider using ConsolidatedWalletRepository")
	// This is a placeholder implementation
	return nil
}
