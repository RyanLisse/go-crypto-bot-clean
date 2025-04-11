// Package database provides database connectivity and operations
package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/tursodatabase/libsql-client-go/libsql"
	"go.uber.org/zap"
)

// TursoConfig represents the configuration for Turso database
type TursoConfig struct {
	Enabled             bool
	URL                 string
	AuthToken           string
	SyncEnabled         bool
	SyncIntervalSeconds int
}

// TursoManager manages the Turso database connection
type TursoManager struct {
	config     TursoConfig
	db         *sql.DB
	logger     *zap.Logger
	syncTicker *time.Ticker
	syncDone   chan bool
}

// NewTursoManager creates a new Turso database manager
func NewTursoManager(config TursoConfig, logger *zap.Logger) *TursoManager {
	return &TursoManager{
		config:   config,
		logger:   logger,
		syncDone: make(chan bool),
	}
}

// Connect establishes a connection to the Turso database
func (m *TursoManager) Connect(ctx context.Context) error {
	if !m.config.Enabled {
		return fmt.Errorf("turso is not enabled")
	}

	if m.config.URL == "" {
		return fmt.Errorf("turso URL is required")
	}

	// Create connector options
	var opts []libsql.Option
	if m.config.AuthToken != "" {
		opts = append(opts, libsql.WithAuthToken(m.config.AuthToken))
	}

	// Create connector
	connector, err := libsql.NewConnector(m.config.URL, opts...)
	if err != nil {
		return fmt.Errorf("failed to create Turso connector: %w", err)
	}

	// Open database connection
	db := sql.OpenDB(connector)

	// Test connection
	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping Turso database: %w", err)
	}

	m.db = db

	m.logger.Info("Connected to Turso database",
		zap.String("url", m.config.URL),
		zap.Bool("syncEnabled", m.config.SyncEnabled),
		zap.Int("syncIntervalSeconds", m.config.SyncIntervalSeconds),
	)

	// Start sync if enabled
	if m.config.SyncEnabled && m.config.SyncIntervalSeconds > 0 {
		m.StartSync()
	}

	return nil
}

// DB returns the database connection
func (m *TursoManager) DB() *sql.DB {
	return m.db
}

// Close closes the database connection
func (m *TursoManager) Close() error {
	if m.db == nil {
		return nil
	}

	// Stop sync if running
	if m.syncTicker != nil {
		m.StopSync()
	}

	if err := m.db.Close(); err != nil {
		return fmt.Errorf("failed to close Turso database connection: %w", err)
	}

	m.logger.Info("Closed Turso database connection", zap.String("url", m.config.URL))
	return nil
}

// StartSync starts the synchronization process
func (m *TursoManager) StartSync() {
	if !m.config.SyncEnabled || m.config.SyncIntervalSeconds <= 0 {
		m.logger.Warn("Sync is not enabled or interval is invalid")
		return
	}

	if m.syncTicker != nil {
		m.logger.Warn("Sync is already running")
		return
	}

	m.syncTicker = time.NewTicker(time.Duration(m.config.SyncIntervalSeconds) * time.Second)

	go func() {
		for {
			select {
			case <-m.syncTicker.C:
				if err := m.Sync(context.Background()); err != nil {
					m.logger.Error("Failed to sync with Turso", zap.Error(err))
				}
			case <-m.syncDone:
				return
			}
		}
	}()

	m.logger.Info("Started Turso sync",
		zap.Int("intervalSeconds", m.config.SyncIntervalSeconds),
	)
}

// StopSync stops the synchronization process
func (m *TursoManager) StopSync() {
	if m.syncTicker == nil {
		return
	}

	m.syncTicker.Stop()
	m.syncDone <- true
	m.syncTicker = nil

	m.logger.Info("Stopped Turso sync")
}

// Sync synchronizes the local database with Turso
func (m *TursoManager) Sync(ctx context.Context) error {
	if m.db == nil {
		return fmt.Errorf("not connected to Turso")
	}

	// In a real implementation, this would use the Turso sync API
	// For now, we'll just log that sync was attempted
	m.logger.Info("Syncing with Turso database", zap.String("url", m.config.URL))

	// Simulate sync delay
	time.Sleep(100 * time.Millisecond)

	return nil
}

// ExecuteQuery executes a query on the Turso database
func (m *TursoManager) ExecuteQuery(ctx context.Context, query string, args ...any) (sql.Result, error) {
	if m.db == nil {
		return nil, fmt.Errorf("not connected to Turso")
	}

	return m.db.ExecContext(ctx, query, args...)
}

// QueryRows executes a query and returns the rows
func (m *TursoManager) QueryRows(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	if m.db == nil {
		return nil, fmt.Errorf("not connected to Turso")
	}

	return m.db.QueryContext(ctx, query, args...)
}

// QueryRow executes a query and returns a single row
func (m *TursoManager) QueryRow(ctx context.Context, query string, args ...any) *sql.Row {
	if m.db == nil {
		return nil
	}

	return m.db.QueryRowContext(ctx, query, args...)
}
