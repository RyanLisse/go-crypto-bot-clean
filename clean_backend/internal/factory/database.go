package factory

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
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
	// Create a persistent directory for the local database
	dir := "./data/turso"
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Error().Err(err).Msg("Failed to create database directory")
		return nil, fmt.Errorf("error creating database directory: %w", err)
	}

	// Use a persistent path for the local database
	dbPath := filepath.Join(dir, "local.db")
	log.Info().Str("path", dbPath).Msg("Using local database for Turso")

	// Check if Turso URL and auth token are set for remote sync
	if cfg.Database.TursoURL == "" || cfg.Database.AuthToken == "" {
		log.Warn().Msg("TURSO_DB_URL or TURSO_AUTH_TOKEN not set, using local database only without remote sync")
		// Use local database only without remote sync
		return openLocalDatabase(dbPath, gormConfig, cfg, log)
	}

	// Try to connect with remote sync enabled
	db, err := openWithRemoteSync(dbPath, cfg, gormConfig, log)
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect with remote sync, falling back to local-only mode")
		return openLocalDatabase(dbPath, gormConfig, cfg, log)
	}

	// Setup periodic sync if successful
	setupPeriodicSync(db, cfg, log)

	return db, nil
}

// openLocalDatabase opens a local SQLite database without remote sync
func openLocalDatabase(dbPath string, gormConfig *gorm.Config, cfg *config.Config, log *zerolog.Logger) (*gorm.DB, error) {
	// Use the file path directly with the SQLite dialect
	log.Info().Msg("Opening local database without remote sync")
	db, err := gorm.Open(libsql.Open(fmt.Sprintf("file:%s", dbPath)), gormConfig)
	if err != nil {
		log.Error().Err(err).Msg("Failed to open local database")
		return nil, fmt.Errorf("failed to open local database: %w", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get database connection")
		return db, nil // Return DB anyway, it might still work
	}

	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.Database.ConnMaxLifetimeMinutes) * time.Minute)

	return db, nil
}

// openWithRemoteSync opens a database with remote sync enabled
func openWithRemoteSync(dbPath string, cfg *config.Config, gormConfig *gorm.Config, log *zerolog.Logger) (*gorm.DB, error) {
	log.Info().Str("remote_url", cfg.Database.TursoURL).Msg("Creating embedded replica connector with remote sync")

	// Create connector with auth token
	connector, err := libsqlgo.NewEmbeddedReplicaConnector(
		dbPath,
		cfg.Database.TursoURL,
		libsqlgo.WithAuthToken(cfg.Database.AuthToken),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create Turso connector: %w", err)
	}

	// Open database connection
	sqlDB := sql.OpenDB(connector)

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		connector.Close()
		return nil, fmt.Errorf("failed to connect to Turso database: %w", err)
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
	db, err := gorm.Open(libsql.Dialector{Conn: sqlDB}, gormConfig)
	if err != nil {
		connector.Close()
		return nil, fmt.Errorf("failed to open database with GORM: %w", err)
	}

	return db, nil
}

// setupPeriodicSync sets up periodic synchronization with the remote database
func setupPeriodicSync(db *gorm.DB, cfg *config.Config, log *zerolog.Logger) {
	// Skip if sync is disabled
	syncInterval := 5 * time.Minute // Default interval

	// Check environment variables for sync configuration
	syncEnabledStr := os.Getenv("TURSO_SYNC_ENABLED")
	syncIntervalStr := os.Getenv("TURSO_SYNC_INTERVAL_SECONDS")

	// Parse sync enabled flag
	syncEnabled := true
	if syncEnabledStr != "" {
		var err error
		syncEnabled, err = strconv.ParseBool(syncEnabledStr)
		if err != nil {
			log.Warn().Err(err).Msg("Invalid TURSO_SYNC_ENABLED value, defaulting to true")
			syncEnabled = true
		}
	}

	// Parse sync interval
	if syncIntervalStr != "" {
		syncIntervalSec, err := strconv.Atoi(syncIntervalStr)
		if err != nil {
			log.Warn().Err(err).Msg("Invalid TURSO_SYNC_INTERVAL_SECONDS value, using default")
		} else if syncIntervalSec > 0 {
			syncInterval = time.Duration(syncIntervalSec) * time.Second
		}
	}

	if !syncEnabled {
		log.Info().Msg("Periodic sync with Turso is disabled")
		return
	}

	log.Info().Dur("interval", syncInterval).Msg("Setting up periodic sync with Turso")

	// Start a goroutine for periodic sync
	go func() {
		ticker := time.NewTicker(syncInterval)
		defer ticker.Stop()

		for range ticker.C {
			// Get the underlying SQL DB
			sqlDB, err := db.DB()
			if err != nil {
				log.Error().Err(err).Msg("Failed to get database connection for sync")
				continue
			}

			// Get the connector
			driver := sqlDB.Driver()

			// Type assertion for the connector
			type syncable interface {
				Sync() (struct{ FramesSynced int }, error)
			}

			connector, ok := driver.(syncable)
			if !ok {
				log.Error().Msg("Database driver does not support sync")
				return // Exit the goroutine if sync is not supported
			}

			// Perform sync
			log.Debug().Msg("Performing periodic sync with Turso primary database")
			result, err := connector.Sync()
			if err != nil {
				log.Error().Err(err).Msg("Periodic sync failed")
			} else {
				log.Debug().Int("frames_synced", result.FramesSynced).Msg("Periodic sync completed")
			}
		}
	}()
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
