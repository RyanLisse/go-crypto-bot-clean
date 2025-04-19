package factory

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/adapter/infrastructure/persistence/turso"
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/config"
	"github.com/rs/zerolog"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DatabaseFactory creates database connections
type DatabaseFactory struct {
	config *config.Config
	logger *zerolog.Logger
}

// NewDatabaseFactory creates a new DatabaseFactory
func NewDatabaseFactory(config *config.Config, logger *zerolog.Logger) *DatabaseFactory {
	return &DatabaseFactory{
		config: config,
		logger: logger,
	}
}

// CreateSQLDB creates a new SQL database connection
func (f *DatabaseFactory) CreateSQLDB() (*sql.DB, error) {
	// Check if Turso is configured
	if f.config.Database.Type == "turso" {
		return f.createTursoDB()
	}

	// Default to SQLite
	return f.createSQLiteDB()
}

// CreateGormDB creates a new GORM database connection
func (f *DatabaseFactory) CreateGormDB() (*gorm.DB, error) {
	// Create a GORM logger that uses zerolog
	gormLogger := logger.New(
		&GormLogAdapter{Logger: f.logger},
		logger.Config{
			SlowThreshold:             time.Second, // Log slow queries
			LogLevel:                  logger.Warn, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore not found errors
			Colorful:                  false,       // Disable color
		},
	)

	// Configure GORM
	gormConfig := &gorm.Config{
		Logger:                 gormLogger,
		SkipDefaultTransaction: true, // For better performance
	}

	// Connect to the database based on the configuration
	var db *gorm.DB
	var err error

	switch f.config.Database.Type {
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(f.config.Database.DSN), gormConfig)
	case "turso":
		// For Turso, we need to get the SQL DB first and then create a GORM connection
		sqlDB, err := f.createTursoDB()
		if err != nil {
			return nil, err
		}

		// Create a GORM connection using the SQL DB
		db, err = gorm.Open(sqlite.Dialector{Conn: sqlDB}, gormConfig)
	default:
		return nil, fmt.Errorf("unsupported database type: %s", f.config.Database.Type)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool for GORM
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection: %w", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(f.config.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(f.config.Database.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(f.config.Database.ConnMaxLifetimeMinutes) * time.Minute)

	f.logger.Info().
		Str("type", f.config.Database.Type).
		Str("dsn", maskDSN(f.config.Database.DSN)).
		Msg("Database connection established")

	return db, nil
}

// createTursoDB creates a connection to a Turso database
func (f *DatabaseFactory) createTursoDB() (*sql.DB, error) {
	// Check if Turso URL and auth token are set
	if f.config.Database.TursoURL == "" || f.config.Database.AuthToken == "" {
		return nil, fmt.Errorf("turso URL or auth token not set")
	}

	// Get sync interval from environment or use default
	syncIntervalStr := f.config.Database.ConnMaxLifetimeMinutes
	syncInterval := 5 * time.Minute // Default sync interval

	if syncIntervalStr > 0 {
		syncInterval = time.Duration(syncIntervalStr) * time.Minute
	}

	f.logger.Info().
		Str("url", f.config.Database.TursoURL).
		Dur("sync_interval", syncInterval).
		Msg("Initializing Turso database")

	// Create Turso database
	tursoDB, err := turso.NewTursoDB(
		f.config.Database.TursoURL,
		f.config.Database.AuthToken,
		syncInterval,
		f.logger,
	)
	if err != nil {
		f.logger.Error().Err(err).Msg("Failed to create Turso database")
		return nil, err
	}

	// Return the database connection
	return tursoDB.DB(), nil
}

// createSQLiteDB creates a connection to a SQLite database
func (f *DatabaseFactory) createSQLiteDB() (*sql.DB, error) {
	// Open SQLite database
	db, err := sql.Open("sqlite3", f.config.Database.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to open SQLite database: %w", err)
	}

	// Configure connection pool
	db.SetMaxIdleConns(f.config.Database.MaxIdleConns)
	db.SetMaxOpenConns(f.config.Database.MaxOpenConns)
	db.SetConnMaxLifetime(time.Duration(f.config.Database.ConnMaxLifetimeMinutes) * time.Minute)

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping SQLite database: %w", err)
	}

	return db, nil
}
