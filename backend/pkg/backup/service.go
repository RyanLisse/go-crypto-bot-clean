package backup

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// backupService implements the BackupService interface
type backupService struct {
	storage BackupStorage
}

// NewBackupService creates a new instance of BackupService
func NewBackupService(storage BackupStorage) BackupService {
	return &backupService{
		storage: storage,
	}
}

// CreateBackup performs a backup operation and returns metadata about the backup
func (s *backupService) CreateBackup(opts BackupOptions) (*BackupMetadata, error) {
	// Validate options
	if err := validateBackupOptions(opts); err != nil {
		return nil, fmt.Errorf("invalid backup options: %w", err)
	}

	// Create metadata
	metadata := &BackupMetadata{
		ID:            generateBackupID(),
		Type:          opts.Type,
		SourceDir:     opts.SourceDir,
		DestDir:       opts.DestDir,
		Timestamp:     time.Now(),
		ExcludeRules:  opts.ExcludeRules,
		RetentionDays: opts.RetentionDays,
	}

	// Create backup writer
	writer, err := s.storage.CreateWriter(metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to create backup writer: %w", err)
	}
	defer writer.Close()

	// Walk source directory and add files
	err = filepath.Walk(opts.SourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip excluded files
		if isExcluded(path, opts.ExcludeRules) {
			return nil
		}

		// Create FileInfo
		relPath, err := filepath.Rel(opts.SourceDir, path)
		if err != nil {
			return fmt.Errorf("failed to get relative path: %w", err)
		}

		fileInfo := FileInfo{
			Path:         path,
			RelativePath: relPath,
			Size:         info.Size(),
			ModTime:      info.ModTime(),
			IsDir:        info.IsDir(),
		}

		// Calculate checksum for files
		if !info.IsDir() {
			checksum, err := calculateChecksum(path)
			if err != nil {
				return fmt.Errorf("failed to calculate checksum: %w", err)
			}
			fileInfo.Checksum = checksum

			// Add file to backup
			file, err := os.Open(path)
			if err != nil {
				return fmt.Errorf("failed to open file: %w", err)
			}
			defer file.Close()

			if err := writer.AddFile(fileInfo, file); err != nil {
				return fmt.Errorf("failed to add file to backup: %w", err)
			}

			// Update metadata
			metadata.Files = append(metadata.Files, fileInfo)
			metadata.TotalSize += info.Size()
			metadata.FileCount++
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create backup: %w", err)
	}

	// Write metadata
	if err := writer.WriteMetadata(metadata); err != nil {
		return nil, fmt.Errorf("failed to write metadata: %w", err)
	}

	return metadata, nil
}

// RestoreBackup restores files from a backup
func (s *backupService) RestoreBackup(backupID string, targetDir string) error {
	reader, err := s.storage.OpenReader(backupID)
	if err != nil {
		return fmt.Errorf("failed to open backup: %w", err)
	}
	defer reader.Close()

	files, err := reader.ListFiles()
	if err != nil {
		return fmt.Errorf("failed to list files: %w", err)
	}

	for _, file := range files {
		targetPath := filepath.Join(targetDir, file.RelativePath)

		// Create parent directories if needed
		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}

		if !file.IsDir {
			// Restore file
			fileReader, _, err := reader.GetFile(file.Path)
			if err != nil {
				return fmt.Errorf("failed to get file from backup: %w", err)
			}
			defer fileReader.Close()

			targetFile, err := os.Create(targetPath)
			if err != nil {
				return fmt.Errorf("failed to create target file: %w", err)
			}
			defer targetFile.Close()

			if _, err := io.Copy(targetFile, fileReader); err != nil {
				return fmt.Errorf("failed to restore file: %w", err)
			}
		}
	}

	return nil
}

// ListBackups returns a list of all available backups
func (s *backupService) ListBackups() ([]BackupMetadata, error) {
	return s.storage.List()
}

// GetBackupMetadata retrieves metadata for a specific backup
func (s *backupService) GetBackupMetadata(backupID string) (*BackupMetadata, error) {
	reader, err := s.storage.OpenReader(backupID)
	if err != nil {
		return nil, fmt.Errorf("failed to open backup: %w", err)
	}
	defer reader.Close()

	return reader.GetMetadata()
}

// DeleteBackup removes a backup and its associated files
func (s *backupService) DeleteBackup(backupID string) error {
	return s.storage.Delete(backupID)
}

// CleanupOldBackups removes backups that are older than their retention period
func (s *backupService) CleanupOldBackups() error {
	backups, err := s.storage.List()
	if err != nil {
		return fmt.Errorf("failed to list backups: %w", err)
	}

	now := time.Now()
	for _, backup := range backups {
		if backup.RetentionDays > 0 {
			retentionPeriod := time.Duration(backup.RetentionDays) * 24 * time.Hour
			if now.Sub(backup.Timestamp) > retentionPeriod {
				if err := s.storage.Delete(backup.ID); err != nil {
					return fmt.Errorf("failed to delete old backup: %w", err)
				}
			}
		}
	}

	return nil
}

// VerifyBackup checks the integrity of a backup
func (s *backupService) VerifyBackup(backupID string) error {
	reader, err := s.storage.OpenReader(backupID)
	if err != nil {
		return fmt.Errorf("failed to open backup: %w", err)
	}
	defer reader.Close()

	files, err := reader.ListFiles()
	if err != nil {
		return fmt.Errorf("failed to list files: %w", err)
	}

	for _, file := range files {
		if !file.IsDir {
			fileReader, info, err := reader.GetFile(file.Path)
			if err != nil {
				return fmt.Errorf("failed to get file from backup: %w", err)
			}
			defer fileReader.Close()

			checksum, err := calculateChecksumFromReader(fileReader)
			if err != nil {
				return fmt.Errorf("failed to calculate checksum: %w", err)
			}

			if checksum != info.Checksum {
				return fmt.Errorf("checksum mismatch for file %s", file.Path)
			}
		}
	}

	return nil
}

// Helper functions

func validateBackupOptions(opts BackupOptions) error {
	if opts.SourceDir == "" {
		return fmt.Errorf("source directory is required")
	}
	if opts.DestDir == "" {
		return fmt.Errorf("destination directory is required")
	}
	if opts.Type != FullBackup && opts.Type != IncrementalBackup {
		return fmt.Errorf("invalid backup type: %s", opts.Type)
	}
	return nil
}

func generateBackupID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func calculateChecksum(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	return calculateChecksumFromReader(file)
}

func calculateChecksumFromReader(reader io.Reader) (string, error) {
	hash := sha256.New()
	if _, err := io.Copy(hash, reader); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func isExcluded(path string, excludeRules []string) bool {
	for _, rule := range excludeRules {
		matched, err := filepath.Match(rule, filepath.Base(path))
		if err == nil && matched {
			return true
		}
	}
	return false
}
