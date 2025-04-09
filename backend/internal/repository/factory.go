package repository

import (
	"github.com/ryanlisse/go-crypto-bot/internal/domain/repository"
	"github.com/ryanlisse/go-crypto-bot/internal/repository/balance"
	"github.com/ryanlisse/go-crypto-bot/internal/repository/database"
	"github.com/ryanlisse/go-crypto-bot/internal/repository/report"
	"go.uber.org/zap"
)

// Factory creates repository instances
type Factory struct {
	db     database.Repository
	logger *zap.Logger
}

// NewFactory creates a new repository factory
func NewFactory(db database.Repository, logger *zap.Logger) *Factory {
	return &Factory{
		db:     db,
		logger: logger,
	}
}

// NewBalanceHistoryRepository creates a new balance history repository
func (f *Factory) NewBalanceHistoryRepository() repository.BalanceHistoryRepository {
	return balance.NewBalanceHistoryRepository(f.db)
}

// NewReportRepository creates a new report repository
func (f *Factory) NewReportRepository() repository.ReportRepository {
	return report.NewReportRepository(f.db, f.logger)
}

// Additional repository factory methods will be added here as we migrate more repositories
