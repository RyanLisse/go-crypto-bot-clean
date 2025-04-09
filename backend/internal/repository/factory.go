package repository

import (
	"go-crypto-bot-clean/backend/internal/domain/repository"
	"go-crypto-bot-clean/backend/internal/repository/balance"
	"go-crypto-bot-clean/backend/internal/repository/database"
	"go-crypto-bot-clean/backend/internal/repository/report"
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
