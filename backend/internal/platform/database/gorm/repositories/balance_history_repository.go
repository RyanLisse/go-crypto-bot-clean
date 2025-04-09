package repositories

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"go-crypto-bot-clean/backend/internal/domain/models"
	"go-crypto-bot-clean/backend/internal/domain/repository"
)

// BalanceHistory represents a point in the balance history with GORM tags
type BalanceHistory struct {
	ID            int64     `gorm:"primaryKey;autoIncrement"`
	Timestamp     time.Time `gorm:"index;not null"`
	Balance       float64   `gorm:"not null"`
	Equity        float64   `gorm:"not null"`
	FreeBalance   float64   `gorm:"not null"`
	LockedBalance float64   `gorm:"not null"`
	UnrealizedPnL float64   `gorm:"not null"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`
}

// TableName specifies the table name for the BalanceHistory model
func (BalanceHistory) TableName() string {
	return "balance_history"
}

// GORMBalanceHistoryRepository implements the BalanceHistoryRepository interface using GORM
type GORMBalanceHistoryRepository struct {
	db *gorm.DB
}

// NewGORMBalanceHistoryRepository creates a new GORM-based balance history repository
func NewGORMBalanceHistoryRepository(db *gorm.DB) repository.BalanceHistoryRepository {
	return &GORMBalanceHistoryRepository{
		db: db,
	}
}

// Create adds a new balance history point
func (r *GORMBalanceHistoryRepository) Create(ctx context.Context, history *repository.BalanceHistory) (int64, error) {
	balanceHistory := BalanceHistory{
		Timestamp:     history.Timestamp,
		Balance:       history.Balance,
		Equity:        history.Equity,
		FreeBalance:   history.FreeBalance,
		LockedBalance: history.LockedBalance,
		UnrealizedPnL: history.UnrealizedPnL,
	}

	result := r.db.WithContext(ctx).Create(&balanceHistory)
	if result.Error != nil {
		return 0, fmt.Errorf("failed to insert balance history: %w", result.Error)
	}

	history.ID = balanceHistory.ID
	return balanceHistory.ID, nil
}

// GetBalanceHistory retrieves balance history within a time range
func (r *GORMBalanceHistoryRepository) GetBalanceHistory(ctx context.Context, startTime, endTime time.Time) ([]*repository.BalanceHistory, error) {
	var balanceHistories []BalanceHistory

	result := r.db.WithContext(ctx).
		Where("timestamp BETWEEN ? AND ?", startTime, endTime).
		Order("timestamp ASC").
		Find(&balanceHistories)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to get balance history: %w", result.Error)
	}

	// Convert to domain model
	domainHistories := make([]*repository.BalanceHistory, len(balanceHistories))
	for i, history := range balanceHistories {
		domainHistories[i] = &repository.BalanceHistory{
			ID:            history.ID,
			Timestamp:     history.Timestamp,
			Balance:       history.Balance,
			Equity:        history.Equity,
			FreeBalance:   history.FreeBalance,
			LockedBalance: history.LockedBalance,
			UnrealizedPnL: history.UnrealizedPnL,
		}
	}

	return domainHistories, nil
}

// GetLatestBalance retrieves the latest balance history point
func (r *GORMBalanceHistoryRepository) GetLatestBalance(ctx context.Context) (*repository.BalanceHistory, error) {
	var balanceHistory BalanceHistory

	result := r.db.WithContext(ctx).
		Order("timestamp DESC").
		First(&balanceHistory)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get latest balance: %w", result.Error)
	}

	return &repository.BalanceHistory{
		ID:            balanceHistory.ID,
		Timestamp:     balanceHistory.Timestamp,
		Balance:       balanceHistory.Balance,
		Equity:        balanceHistory.Equity,
		FreeBalance:   balanceHistory.FreeBalance,
		LockedBalance: balanceHistory.LockedBalance,
		UnrealizedPnL: balanceHistory.UnrealizedPnL,
	}, nil
}

// GetBalancePoints retrieves balance points for equity curve
func (r *GORMBalanceHistoryRepository) GetBalancePoints(ctx context.Context, startTime, endTime time.Time) ([]models.BalancePoint, error) {
	var balancePoints []models.BalancePoint

	result := r.db.WithContext(ctx).
		Model(&models.BalanceHistory{}).
		Select("timestamp, balance").
		Where("timestamp BETWEEN ? AND ?", startTime, endTime).
		Order("timestamp ASC").
		Scan(&balancePoints)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to get balance points: %w", result.Error)
	}

	return balancePoints, nil
}
