package factory

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/config"
	libsql "github.com/ekristen/gorm-libsql"
	"github.com/rs/zerolog"
	libsqlgo "github.com/tursodatabase/go-libsql" // Import the Turso libsql package
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

// NewDBConnection creates a new database connection based on the configuration
func NewDBConnection(cfg *config.Config, log *zerolog.Logger) (*gorm.DB, error) {
	// Create a GORM logger that uses zerolog
	gormLogger := gormLogger.New(
		&GormLogAdapter{Logger: log},
		gormLogger.Config{
			SlowThreshold:             time.Second,     // Log slow queries
			LogLevel:                  gormLogger.Warn, // Log level
			IgnoreRecordNotFoundError: true,            // Ignore not found errors
			Colorful:                  false,           // Disable color
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

	switch cfg.Database.Type {
	case "turso", "libsql":
		db, err = connectTurso(cfg, gormConfig, log)
	case "sqlite":
		log.Warn().Msg("SQLite driver is deprecated, please use Turso instead")
		db, err = connectTurso(cfg, gormConfig, log) // Fallback to Turso
	case "mysql":
		db, err = gorm.Open(mysql.Open(cfg.Database.DSN), gormConfig)
	case "postgres":
		// Implement PostgreSQL connection if needed
		return nil, fmt.Errorf("postgres database not implemented yet")
	default:
		return nil, fmt.Errorf("unsupported database type: %s", cfg.Database.Type)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection: %w", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.Database.ConnMaxLifetimeMinutes) * time.Minute)

	log.Info().
		Str("type", cfg.Database.Type).
		Str("dsn", maskDSN(cfg.Database.DSN)).
		Msg("Database connection established")

	return db, nil
}

// connectTurso creates a connection to a Turso database using the gorm-libsql driver
func connectTurso(cfg *config.Config, gormConfig *gorm.Config, log *zerolog.Logger) (*gorm.DB, error) {
	// Check if Turso URL and auth token are set
	if cfg.Database.TursoURL == "" || cfg.Database.AuthToken == "" {
		log.Warn().Msg("TURSO_DB_URL or TURSO_AUTH_TOKEN not set, using local database only")
		// Use local database only
		return gorm.Open(libsql.Open(cfg.Database.DSN), gormConfig)
	}

	// Create a persistent directory for the local database
	dir := "./data/turso"
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("error creating database directory: %w", err)
	}

	// Use a persistent path for the local database
	dbPath := filepath.Join(dir, "local.db")
	log.Info().Str("path", dbPath).Msg("Using local database for Turso")

	// Create connector with auth token
	connector, err := libsqlgo.NewEmbeddedReplicaConnector(
		dbPath,
		cfg.Database.TursoURL,
		libsqlgo.WithAuthToken(cfg.Database.AuthToken),
	)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create Turso connector, falling back to local database")
		return gorm.Open(libsql.Open(cfg.Database.DSN), gormConfig)
	}

	// Open database connection
	sqlDB := sql.OpenDB(connector)
	if err := sqlDB.Ping(); err != nil {
		connector.Close()
		log.Error().Err(err).Msg("Failed to connect to Turso database, falling back to local database")
		return gorm.Open(libsql.Open(cfg.Database.DSN), gormConfig)
	}

	// Perform initial sync
	log.Info().Msg("Performing initial sync with Turso primary database")
	result, syncErr := connector.Sync()
	if syncErr != nil {
		log.Warn().Err(syncErr).Msg("Initial sync failed, will retry later")
	} else {
		log.Info().Int("frames_synced", result.FramesSynced).Msg("Initial sync completed successfully")
	}

	// Configure connection pool
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.Database.ConnMaxLifetimeMinutes) * time.Minute)

	// Create GORM DB with the SQLite dialect using the Turso connector
	return gorm.Open(libsql.Dialector{Conn: sqlDB}, gormConfig)
}

// maskDSN masks sensitive information in the DSN for logging
func maskDSN(dsn string) string {
	// For SQLite, just return the path
	return dsn
}

// GormLogAdapter adapts zerolog to GORM's logger interface
type GormLogAdapter struct {
	Logger *zerolog.Logger
}

// Printf implements GORM's logger interface
func (l *GormLogAdapter) Printf(format string, args ...interface{}) {
	l.Logger.Debug().Msgf(format, args...)
}

// CloseDBConnection closes the database connection
func CloseDBConnection(db *gorm.DB) {
	// Implement the logic to close the database connection
}
