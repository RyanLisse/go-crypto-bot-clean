package commands

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go.uber.org/zap"
)

// BackupService handles the backup operations
type BackupService struct {
	logger *zap.Logger
}

// NewBackupService creates a new backup service
func NewBackupService(logger *zap.Logger) *BackupService {
	return &BackupService{
		logger: logger,
	}
}

// BackupMetadata contains information about the backup
type BackupMetadata struct {
	Type           string    `json:"type"`
	Timestamp      time.Time `json:"timestamp"`
	SourceDir      string    `json:"source_dir"`
	Files          []string  `json:"files"`
	ExcludePattern []string  `json:"exclude_pattern"`
	Checksums      []string  `json:"checksums"`
	TotalSize      int64     `json:"total_size"`
	CompressedSize int64     `json:"compressed_size"`
}

// Backup performs the backup operation
func (s *BackupService) Backup(ctx context.Context, opts *BackupOptions) (*BackupResult, error) {
	startTime := time.Now()

	// Create backup filename with timestamp
	timestamp := time.Now().Format("20060102-150405")
	backupFile := filepath.Join(opts.Destination, fmt.Sprintf("backup-%s-%s.tar.gz", opts.Type, timestamp))
	metadataFile := backupFile + ".json"

	// Create the backup file
	file, err := os.Create(backupFile)
	if err != nil {
		return nil, fmt.Errorf("failed to create backup file: %w", err)
	}
	defer file.Close()

	// Create gzip writer
	gw := gzip.NewWriter(file)
	defer gw.Close()

	// Create tar writer
	tw := tar.NewWriter(gw)
	defer tw.Close()

	// Initialize metadata
	metadata := &BackupMetadata{
		Type:           opts.Type,
		Timestamp:      startTime,
		SourceDir:      opts.Source,
		ExcludePattern: opts.Exclude,
	}

	// Walk through source directory
	fileCount := 0
	err = filepath.Walk(opts.Source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check context cancellation
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Skip if matches exclude patterns
		for _, pattern := range opts.Exclude {
			if matched, _ := filepath.Match(pattern, filepath.Base(path)); matched {
				s.logger.Debug("Skipping excluded file", zap.String("path", path))
				return nil
			}
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Get relative path
		relPath, err := filepath.Rel(opts.Source, path)
		if err != nil {
			return fmt.Errorf("failed to get relative path: %w", err)
		}

		// Open the file
		file, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("failed to open file %s: %w", path, err)
		}
		defer file.Close()

		// Calculate checksum
		hash := sha256.New()
		if _, err := io.Copy(hash, file); err != nil {
			return fmt.Errorf("failed to calculate checksum for %s: %w", path, err)
		}
		checksum := hex.EncodeToString(hash.Sum(nil))

		// Reset file pointer
		if _, err := file.Seek(0, 0); err != nil {
			return fmt.Errorf("failed to reset file pointer for %s: %w", path, err)
		}

		// Create tar header
		header := &tar.Header{
			Name:    relPath,
			Size:    info.Size(),
			Mode:    int64(info.Mode()),
			ModTime: info.ModTime(),
		}

		// Write header
		if err := tw.WriteHeader(header); err != nil {
			return fmt.Errorf("failed to write tar header for %s: %w", path, err)
		}

		// Copy file content
		if _, err := io.Copy(tw, file); err != nil {
			return fmt.Errorf("failed to write file content for %s: %w", path, err)
		}

		// Update metadata
		metadata.Files = append(metadata.Files, relPath)
		metadata.Checksums = append(metadata.Checksums, checksum)
		metadata.TotalSize += info.Size()
		fileCount++

		s.logger.Debug("Added file to backup",
			zap.String("path", relPath),
			zap.Int64("size", info.Size()),
			zap.String("checksum", checksum),
		)

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk source directory: %w", err)
	}

	// Close writers to ensure all data is written
	tw.Close()
	gw.Close()

	// Get compressed size
	stat, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get backup file stats: %w", err)
	}
	metadata.CompressedSize = stat.Size()

	// Write metadata file
	metadataJSON, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal metadata: %w", err)
	}
	if err := os.WriteFile(metadataFile, metadataJSON, 0644); err != nil {
		return nil, fmt.Errorf("failed to write metadata file: %w", err)
	}

	// Clean up old backups if retention is set
	if opts.Retention > 0 {
		if err := s.cleanOldBackups(opts.Destination, opts.Retention); err != nil {
			s.logger.Error("Failed to clean old backups", zap.Error(err))
		}
	}

	// Return backup result
	return &BackupResult{
		Type:           opts.Type,
		StartTime:      startTime,
		EndTime:        time.Now(),
		BackupFile:     backupFile,
		FileCount:      fileCount,
		TotalSize:      metadata.TotalSize,
		CompressedSize: metadata.CompressedSize,
	}, nil
}

// cleanOldBackups removes backup files older than the retention period
func (s *BackupService) cleanOldBackups(backupDir string, retentionDays int) error {
	cutoff := time.Now().AddDate(0, 0, -retentionDays)

	entries, err := os.ReadDir(backupDir)
	if err != nil {
		return fmt.Errorf("failed to read backup directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() && (strings.HasSuffix(entry.Name(), ".tar.gz") || strings.HasSuffix(entry.Name(), ".tar.gz.json")) {
			info, err := entry.Info()
			if err != nil {
				s.logger.Error("Failed to get file info", zap.String("file", entry.Name()), zap.Error(err))
				continue
			}

			if info.ModTime().Before(cutoff) {
				path := filepath.Join(backupDir, entry.Name())
				if err := os.Remove(path); err != nil {
					s.logger.Error("Failed to remove old backup file",
						zap.String("file", path),
						zap.Error(err),
					)
				} else {
					s.logger.Info("Removed old backup file",
						zap.String("file", path),
						zap.Time("modified", info.ModTime()),
					)
				}
			}
		}
	}

	return nil
}
