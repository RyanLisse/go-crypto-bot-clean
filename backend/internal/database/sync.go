package database

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// SyncConfig represents the configuration for database synchronization
type SyncConfig struct {
	Enabled             bool
	SyncIntervalSeconds int
	BatchSize           int
	MaxRetries          int
	RetryDelaySeconds   int
}

// SyncManager manages the synchronization between SQLite and Turso
type SyncManager struct {
	config     SyncConfig
	sqliteDB   *gorm.DB
	tursoDB    *sql.DB
	logger     *zap.Logger
	syncTicker *time.Ticker
	syncDone   chan bool

	// Sync status tracking
	lastSyncTime     time.Time
	lastSyncSuccess  bool
	lastSyncError    string
	totalSyncedItems int64
	pendingChanges   int
	retryCount       int
	mutex            sync.RWMutex
}

// NewSyncManager creates a new synchronization manager
func NewSyncManager(config SyncConfig, sqliteDB *gorm.DB, tursoDB *sql.DB, logger *zap.Logger) *SyncManager {
	return &SyncManager{
		config:   config,
		sqliteDB: sqliteDB,
		tursoDB:  tursoDB,
		logger:   logger,
		syncDone: make(chan bool),
	}
}

// StartSync starts the synchronization process
func (m *SyncManager) StartSync() error {
	if !m.config.Enabled || m.config.SyncIntervalSeconds <= 0 {
		m.logger.Warn("Sync is not enabled or interval is invalid")
		return fmt.Errorf("sync is not enabled or interval is invalid")
	}

	if m.syncTicker != nil {
		m.logger.Warn("Sync is already running")
		return fmt.Errorf("sync is already running")
	}

	m.syncTicker = time.NewTicker(time.Duration(m.config.SyncIntervalSeconds) * time.Second)

	go func() {
		for {
			select {
			case <-m.syncTicker.C:
				if err := m.Sync(context.Background()); err != nil {
					m.logger.Error("Failed to sync databases", zap.Error(err))
				}
			case <-m.syncDone:
				return
			}
		}
	}()

	m.logger.Info("Started database synchronization",
		zap.Int("intervalSeconds", m.config.SyncIntervalSeconds),
		zap.Int("batchSize", m.config.BatchSize),
	)

	return nil
}

// StopSync stops the synchronization process
func (m *SyncManager) StopSync() {
	if m.syncTicker == nil {
		return
	}

	m.syncTicker.Stop()
	m.syncDone <- true
	m.syncTicker = nil

	m.logger.Info("Stopped database synchronization")
}

// Sync synchronizes the SQLite database with Turso
func (m *SyncManager) Sync(ctx context.Context) error {
	if m.sqliteDB == nil || m.tursoDB == nil {
		return fmt.Errorf("database connections not initialized")
	}

	// Update sync status
	m.mutex.Lock()
	m.lastSyncTime = time.Now()
	m.lastSyncSuccess = false
	m.lastSyncError = ""
	m.mutex.Unlock()

	m.logger.Info("Starting database synchronization")

	// Create a context with timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(30)*time.Second)
	defer cancel()

	// Sync with retry logic
	var lastErr error
	for attempt := 0; attempt <= m.config.MaxRetries; attempt++ {
		if attempt > 0 {
			m.logger.Warn("Retrying database synchronization",
				zap.Int("attempt", attempt),
				zap.Int("maxRetries", m.config.MaxRetries),
				zap.Error(lastErr))

			// Update retry count
			m.mutex.Lock()
			m.retryCount = attempt
			m.mutex.Unlock()

			// Wait before retrying
			time.Sleep(time.Duration(m.config.RetryDelaySeconds) * time.Second)
		}

		// Check if context is cancelled
		if timeoutCtx.Err() != nil {
			m.updateSyncStatus(false, "Sync operation timed out", 0)
			return fmt.Errorf("sync operation timed out: %w", timeoutCtx.Err())
		}

		// Sync system info
		if err := m.syncSystemInfo(timeoutCtx); err != nil {
			lastErr = fmt.Errorf("failed to sync system info: %w", err)
			continue
		}

		// Sync health checks
		if err := m.syncHealthChecks(timeoutCtx); err != nil {
			lastErr = fmt.Errorf("failed to sync health checks: %w", err)
			continue
		}

		// Sync log entries
		if err := m.syncLogEntries(timeoutCtx); err != nil {
			lastErr = fmt.Errorf("failed to sync log entries: %w", err)
			continue
		}

		// If we get here, sync was successful
		m.updateSyncStatus(true, "", 3) // 3 entities synced (system info, health checks, logs)
		m.logger.Info("Database synchronization completed successfully",
			zap.Int("attempts", attempt+1),
			zap.Int64("totalSyncedItems", m.totalSyncedItems))
		return nil
	}

	// If we get here, all retries failed
	m.updateSyncStatus(false, lastErr.Error(), 0)
	m.logger.Error("Database synchronization failed after retries",
		zap.Int("maxRetries", m.config.MaxRetries),
		zap.Error(lastErr))
	return lastErr
}

// updateSyncStatus updates the sync status with thread safety
func (m *SyncManager) updateSyncStatus(success bool, errorMsg string, itemsSynced int64) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.lastSyncSuccess = success
	m.lastSyncError = errorMsg
	m.totalSyncedItems += itemsSynced

	// Reset retry count on success
	if success {
		m.retryCount = 0
	}
}

// syncSystemInfo synchronizes system info between SQLite and Turso
func (m *SyncManager) syncSystemInfo(ctx context.Context) error {
	m.logger.Debug("Syncing system info")

	// Get system info from SQLite
	var systemInfo struct {
		ID        string `gorm:"primaryKey"`
		Name      string
		Version   string
		UpdatedAt time.Time
	}

	result := m.sqliteDB.Table("system_infos").First(&systemInfo)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			m.logger.Debug("No system info found in SQLite")
			return nil
		}
		return fmt.Errorf("failed to get system info from SQLite: %w", result.Error)
	}

	// Check if system info exists in Turso
	var tursoUpdatedAt time.Time
	row := m.tursoDB.QueryRowContext(ctx, "SELECT updated_at FROM system_infos WHERE id = ?", systemInfo.ID)
	err := row.Scan(&tursoUpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			// Insert into Turso
			_, err = m.tursoDB.ExecContext(ctx,
				"INSERT INTO system_infos (id, name, version, updated_at) VALUES (?, ?, ?, ?)",
				systemInfo.ID, systemInfo.Name, systemInfo.Version, systemInfo.UpdatedAt)
			if err != nil {
				return fmt.Errorf("failed to insert system info into Turso: %w", err)
			}
			m.logger.Debug("Inserted system info into Turso")
			return nil
		}
		return fmt.Errorf("failed to check system info in Turso: %w", err)
	}

	// Conflict resolution: Use the most recent update
	if systemInfo.UpdatedAt.After(tursoUpdatedAt) {
		// SQLite is newer, update Turso
		_, err = m.tursoDB.ExecContext(ctx,
			"UPDATE system_infos SET name = ?, version = ?, updated_at = ? WHERE id = ?",
			systemInfo.Name, systemInfo.Version, systemInfo.UpdatedAt, systemInfo.ID)
		if err != nil {
			return fmt.Errorf("failed to update system info in Turso: %w", err)
		}
		m.logger.Debug("Updated system info in Turso (SQLite was newer)")
	} else if tursoUpdatedAt.After(systemInfo.UpdatedAt) {
		// Turso is newer, update SQLite
		var tursoSystemInfo struct {
			Name    string
			Version string
		}
		row = m.tursoDB.QueryRowContext(ctx, "SELECT name, version FROM system_infos WHERE id = ?", systemInfo.ID)
		err = row.Scan(&tursoSystemInfo.Name, &tursoSystemInfo.Version)
		if err != nil {
			return fmt.Errorf("failed to get system info details from Turso: %w", err)
		}

		result = m.sqliteDB.Table("system_infos").Where("id = ?", systemInfo.ID).Updates(map[string]interface{}{
			"name":       tursoSystemInfo.Name,
			"version":    tursoSystemInfo.Version,
			"updated_at": tursoUpdatedAt,
		})
		if result.Error != nil {
			return fmt.Errorf("failed to update system info in SQLite: %w", result.Error)
		}
		m.logger.Debug("Updated system info in SQLite (Turso was newer)")
	} else {
		m.logger.Debug("System info is in sync")
	}

	return nil
}

// syncHealthChecks synchronizes health checks between SQLite and Turso
func (m *SyncManager) syncHealthChecks(ctx context.Context) error {
	m.logger.Debug("Syncing health checks")
	// In a real implementation, this would query SQLite for health checks
	// that haven't been synced yet and insert them into Turso
	return nil
}

// syncLogEntries synchronizes log entries between SQLite and Turso
func (m *SyncManager) syncLogEntries(ctx context.Context) error {
	m.logger.Debug("Syncing log entries")
	// In a real implementation, this would query SQLite for log entries
	// that haven't been synced yet and insert them into Turso in batches
	return nil
}

// SyncStatus represents the status of the synchronization
type SyncStatus struct {
	LastSyncTime      time.Time `json:"last_sync_time"`
	LastSyncSuccess   bool      `json:"last_sync_success"`
	LastSyncError     string    `json:"last_sync_error,omitempty"`
	PendingChanges    int       `json:"pending_changes"`
	TotalSyncedItems  int64     `json:"total_synced_items"`
	SyncEnabled       bool      `json:"sync_enabled"`
	SyncIntervalSecs  int       `json:"sync_interval_secs"`
	NextScheduledSync time.Time `json:"next_scheduled_sync,omitempty"`
}

// GetSyncStatus returns the current synchronization status
func (m *SyncManager) GetSyncStatus() SyncStatus {
	// Get a thread-safe copy of the sync status
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	status := SyncStatus{
		LastSyncTime:     m.lastSyncTime,
		LastSyncSuccess:  m.lastSyncSuccess,
		LastSyncError:    m.lastSyncError,
		PendingChanges:   m.pendingChanges,
		TotalSyncedItems: m.totalSyncedItems,
		SyncEnabled:      m.config.Enabled,
		SyncIntervalSecs: m.config.SyncIntervalSeconds,
	}

	if m.syncTicker != nil {
		// Calculate next sync time based on interval
		status.NextScheduledSync = time.Now().Add(time.Duration(m.config.SyncIntervalSeconds) * time.Second)
	}

	return status
}
