package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// BackupConfig represents the configuration for database backups
type BackupConfig struct {
	Enabled          bool
	BackupDir        string
	BackupInterval   time.Duration
	MaxBackups       int
	RetentionDays    int
	CompressBackups  bool
	IncludeTimestamp bool
}

// BackupManager manages database backups
type BackupManager struct {
	config       BackupConfig
	sqliteDB     *gorm.DB
	tursoDB      *sql.DB
	logger       *zap.Logger
	backupTicker *time.Ticker
	backupDone   chan bool
}

// NewBackupManager creates a new backup manager
func NewBackupManager(config BackupConfig, sqliteDB *gorm.DB, tursoDB *sql.DB, logger *zap.Logger) (*BackupManager, error) {
	// Create backup directory if it doesn't exist
	if config.Enabled && config.BackupDir != "" {
		if err := os.MkdirAll(config.BackupDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create backup directory: %w", err)
		}
	}

	return &BackupManager{
		config:     config,
		sqliteDB:   sqliteDB,
		tursoDB:    tursoDB,
		logger:     logger,
		backupDone: make(chan bool),
	}, nil
}

// StartScheduledBackups starts scheduled backups
func (m *BackupManager) StartScheduledBackups() error {
	if !m.config.Enabled {
		m.logger.Info("Backups are disabled")
		return nil
	}

	if m.config.BackupInterval <= 0 {
		return fmt.Errorf("invalid backup interval: %v", m.config.BackupInterval)
	}

	if m.backupTicker != nil {
		m.logger.Warn("Backup scheduler is already running")
		return fmt.Errorf("backup scheduler is already running")
	}

	m.backupTicker = time.NewTicker(m.config.BackupInterval)

	go func() {
		for {
			select {
			case <-m.backupTicker.C:
				if err := m.BackupDatabases(context.Background()); err != nil {
					m.logger.Error("Scheduled backup failed", zap.Error(err))
				}
			case <-m.backupDone:
				return
			}
		}
	}()

	m.logger.Info("Started scheduled database backups",
		zap.Duration("interval", m.config.BackupInterval),
		zap.String("backupDir", m.config.BackupDir),
	)

	return nil
}

// StopScheduledBackups stops scheduled backups
func (m *BackupManager) StopScheduledBackups() {
	if m.backupTicker == nil {
		return
	}

	m.backupTicker.Stop()
	m.backupDone <- true
	m.backupTicker = nil

	m.logger.Info("Stopped scheduled database backups")
}

// BackupDatabases performs a backup of both SQLite and Turso databases
func (m *BackupManager) BackupDatabases(ctx context.Context) error {
	if !m.config.Enabled {
		return fmt.Errorf("backups are not enabled")
	}

	m.logger.Info("Starting database backup")

	// Create backup timestamp
	timestamp := time.Now().Format("20060102-150405")

	// Backup SQLite database
	if err := m.backupSQLite(ctx, timestamp); err != nil {
		return fmt.Errorf("failed to backup SQLite database: %w", err)
	}

	// Backup Turso database (if connected)
	if m.tursoDB != nil {
		if err := m.backupTurso(ctx, timestamp); err != nil {
			return fmt.Errorf("failed to backup Turso database: %w", err)
		}
	}

	// Clean up old backups
	if err := m.cleanupOldBackups(); err != nil {
		m.logger.Error("Failed to clean up old backups", zap.Error(err))
		// Continue even if cleanup fails
	}

	m.logger.Info("Database backup completed successfully")
	return nil
}

// backupSQLite backs up the SQLite database
func (m *BackupManager) backupSQLite(ctx context.Context, timestamp string) error {
	// Get the SQLite database path
	var dbPath string
	err := m.sqliteDB.Raw("PRAGMA database_list").Row().Scan(nil, nil, &dbPath)
	if err != nil {
		return fmt.Errorf("failed to get SQLite database path: %w", err)
	}

	// Create backup filename
	backupFilename := "sqlite-backup"
	if m.config.IncludeTimestamp {
		backupFilename = fmt.Sprintf("%s-%s", backupFilename, timestamp)
	}
	backupFilename = fmt.Sprintf("%s.db", backupFilename)
	backupPath := filepath.Join(m.config.BackupDir, backupFilename)

	// Create a backup using the SQLite backup API
	// In a real implementation, we would use the SQLite backup API
	// For now, we'll just copy the file
	srcFile, err := os.Open(dbPath)
	if err != nil {
		return fmt.Errorf("failed to open source database file: %w", err)
	}
	defer srcFile.Close()

	dstFile, err := os.Create(backupPath)
	if err != nil {
		return fmt.Errorf("failed to create backup file: %w", err)
	}
	defer dstFile.Close()

	// Copy the database file
	if _, err := dstFile.ReadFrom(srcFile); err != nil {
		return fmt.Errorf("failed to copy database file: %w", err)
	}

	// Compress the backup if enabled
	if m.config.CompressBackups {
		// In a real implementation, we would compress the backup file
		// For now, we'll just log that compression would happen
		m.logger.Info("Backup compression would be applied here")
	}

	m.logger.Info("SQLite database backup completed",
		zap.String("source", dbPath),
		zap.String("backup", backupPath),
	)

	return nil
}

// backupTurso backs up the Turso database
func (m *BackupManager) backupTurso(ctx context.Context, timestamp string) error {
	// Create backup filename
	backupFilename := "turso-backup"
	if m.config.IncludeTimestamp {
		backupFilename = fmt.Sprintf("%s-%s", backupFilename, timestamp)
	}
	backupFilename = fmt.Sprintf("%s.sql", backupFilename)
	backupPath := filepath.Join(m.config.BackupDir, backupFilename)

	// Create backup file
	backupFile, err := os.Create(backupPath)
	if err != nil {
		return fmt.Errorf("failed to create Turso backup file: %w", err)
	}
	defer backupFile.Close()

	// In a real implementation, we would use the Turso API to create a backup
	// For now, we'll just export the schema and data as SQL statements

	// Export schema
	tables, err := m.getTursoTables(ctx)
	if err != nil {
		return fmt.Errorf("failed to get Turso tables: %w", err)
	}

	// Write schema creation statements
	for _, table := range tables {
		schema, err := m.getTursoTableSchema(ctx, table)
		if err != nil {
			return fmt.Errorf("failed to get schema for table %s: %w", table, err)
		}

		if _, err := fmt.Fprintf(backupFile, "%s;\n\n", schema); err != nil {
			return fmt.Errorf("failed to write schema to backup file: %w", err)
		}
	}

	// Export data
	for _, table := range tables {
		if err := m.exportTursoTableData(ctx, backupFile, table); err != nil {
			return fmt.Errorf("failed to export data for table %s: %w", table, err)
		}
	}

	// Compress the backup if enabled
	if m.config.CompressBackups {
		// In a real implementation, we would compress the backup file
		// For now, we'll just log that compression would happen
		m.logger.Info("Backup compression would be applied here")
	}

	m.logger.Info("Turso database backup completed",
		zap.String("backup", backupPath),
	)

	return nil
}

// getTursoTables gets a list of tables in the Turso database
func (m *BackupManager) getTursoTables(ctx context.Context) ([]string, error) {
	rows, err := m.tursoDB.QueryContext(ctx, "SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%'")
	if err != nil {
		return nil, fmt.Errorf("failed to query tables: %w", err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, fmt.Errorf("failed to scan table name: %w", err)
		}
		tables = append(tables, tableName)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating table rows: %w", err)
	}

	return tables, nil
}

// getTursoTableSchema gets the schema for a table in the Turso database
func (m *BackupManager) getTursoTableSchema(ctx context.Context, table string) (string, error) {
	var schema string
	err := m.tursoDB.QueryRowContext(ctx, fmt.Sprintf("SELECT sql FROM sqlite_master WHERE type='table' AND name='%s'", table)).Scan(&schema)
	if err != nil {
		return "", fmt.Errorf("failed to get table schema: %w", err)
	}
	return schema, nil
}

// exportTursoTableData exports the data from a table in the Turso database
func (m *BackupManager) exportTursoTableData(ctx context.Context, file *os.File, table string) error {
	// Get column names
	rows, err := m.tursoDB.QueryContext(ctx, fmt.Sprintf("PRAGMA table_info(%s)", table))
	if err != nil {
		return fmt.Errorf("failed to get table info: %w", err)
	}

	var columns []string
	for rows.Next() {
		var cid, notnull, pk int
		var name, dataType, dfltValue string
		if err := rows.Scan(&cid, &name, &dataType, &notnull, &dfltValue, &pk); err != nil {
			rows.Close()
			return fmt.Errorf("failed to scan column info: %w", err)
		}
		columns = append(columns, name)
	}
	rows.Close()

	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating column rows: %w", err)
	}

	// Get data
	rows, err = m.tursoDB.QueryContext(ctx, fmt.Sprintf("SELECT * FROM %s", table))
	if err != nil {
		return fmt.Errorf("failed to query table data: %w", err)
	}
	defer rows.Close()

	// Write INSERT statements
	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return fmt.Errorf("failed to scan row: %w", err)
		}

		// Build INSERT statement
		insertStmt := fmt.Sprintf("INSERT INTO %s (", table)
		for i, col := range columns {
			if i > 0 {
				insertStmt += ", "
			}
			insertStmt += col
		}
		insertStmt += ") VALUES ("

		for i, val := range values {
			if i > 0 {
				insertStmt += ", "
			}

			// Format value based on type
			switch v := val.(type) {
			case nil:
				insertStmt += "NULL"
			case int64, int32, int16, int8, int, uint64, uint32, uint16, uint8, uint:
				insertStmt += fmt.Sprintf("%d", v)
			case float64, float32:
				insertStmt += fmt.Sprintf("%f", v)
			case bool:
				if v {
					insertStmt += "1"
				} else {
					insertStmt += "0"
				}
			case []byte:
				insertStmt += fmt.Sprintf("X'%x'", v)
			case string:
				insertStmt += fmt.Sprintf("'%s'", escapeSQL(v))
			default:
				insertStmt += fmt.Sprintf("'%s'", escapeSQL(fmt.Sprintf("%v", v)))
			}
		}
		insertStmt += ");\n"

		if _, err := file.WriteString(insertStmt); err != nil {
			return fmt.Errorf("failed to write INSERT statement: %w", err)
		}
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating data rows: %w", err)
	}

	// Add a separator between tables
	if _, err := file.WriteString("\n"); err != nil {
		return fmt.Errorf("failed to write separator: %w", err)
	}

	return nil
}

// backupFile represents a backup file with metadata
type backupFile struct {
	path    string
	modTime time.Time
	size    int64
}

// cleanupOldBackups removes old backups based on retention policy
func (m *BackupManager) cleanupOldBackups() error {
	if m.config.RetentionDays <= 0 && m.config.MaxBackups <= 0 {
		// No retention policy, keep all backups
		return nil
	}

	// List backup files
	entries, err := os.ReadDir(m.config.BackupDir)
	if err != nil {
		return fmt.Errorf("failed to read backup directory: %w", err)
	}

	// Filter and sort backup files
	var backups []backupFile
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		// Check if it's a backup file
		name := entry.Name()
		if !isBackupFile(name) {
			continue
		}

		// Get file info
		info, err := entry.Info()
		if err != nil {
			m.logger.Warn("Failed to get file info", zap.String("file", name), zap.Error(err))
			continue
		}

		backups = append(backups, backupFile{
			path:    filepath.Join(m.config.BackupDir, name),
			modTime: info.ModTime(),
			size:    info.Size(),
		})
	}

	// Sort by modification time (oldest first)
	sortBackupsByTime(backups)

	// Remove old backups based on retention days
	if m.config.RetentionDays > 0 {
		cutoffTime := time.Now().AddDate(0, 0, -m.config.RetentionDays)
		for _, backup := range backups {
			if backup.modTime.Before(cutoffTime) {
				if err := os.Remove(backup.path); err != nil {
					m.logger.Warn("Failed to remove old backup",
						zap.String("file", backup.path),
						zap.Error(err),
					)
				} else {
					m.logger.Info("Removed old backup",
						zap.String("file", backup.path),
						zap.Time("modTime", backup.modTime),
					)
				}
			}
		}
	}

	// Remove excess backups based on max backups
	if m.config.MaxBackups > 0 && len(backups) > m.config.MaxBackups {
		// Re-sort by modification time (oldest first)
		sortBackupsByTime(backups)

		// Remove oldest backups
		for i := 0; i < len(backups)-m.config.MaxBackups; i++ {
			if err := os.Remove(backups[i].path); err != nil {
				m.logger.Warn("Failed to remove excess backup",
					zap.String("file", backups[i].path),
					zap.Error(err),
				)
			} else {
				m.logger.Info("Removed excess backup",
					zap.String("file", backups[i].path),
					zap.Time("modTime", backups[i].modTime),
				)
			}
		}
	}

	return nil
}

// isBackupFile checks if a file is a backup file
func isBackupFile(name string) bool {
	ext := filepath.Ext(name)
	return ext == ".db" || ext == ".sql" || ext == ".gz" || ext == ".zip"
}

// sortBackupsByTime sorts backup files by modification time (oldest first)
func sortBackupsByTime(backups []backupFile) {
	for i := 0; i < len(backups); i++ {
		for j := i + 1; j < len(backups); j++ {
			if backups[i].modTime.After(backups[j].modTime) {
				backups[i], backups[j] = backups[j], backups[i]
			}
		}
	}
}

// escapeSQL escapes a string for use in SQL
func escapeSQL(s string) string {
	return strings.Replace(s, "'", "''", -1)
}
