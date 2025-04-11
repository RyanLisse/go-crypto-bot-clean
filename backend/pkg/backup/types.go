package backup

import (
	"io"
	"time"
)

// BackupType represents the type of backup to perform
type BackupType string

const (
	// FullBackup represents a complete backup of all files
	FullBackup BackupType = "full"
	// IncrementalBackup represents a backup of only changed files since last backup
	IncrementalBackup BackupType = "incremental"
)

// BackupOptions contains configuration for a backup operation
type BackupOptions struct {
	SourceDir     string
	DestDir       string
	Type          BackupType
	ExcludeRules  []string
	RetentionDays int
}

// BackupMetadata contains information about a backup
type BackupMetadata struct {
	ID            string
	Type          BackupType
	SourceDir     string
	DestDir       string
	Timestamp     time.Time
	Files         []FileInfo
	TotalSize     int64
	FileCount     int
	ExcludeRules  []string
	RetentionDays int
}

// FileInfo contains information about a file in a backup
type FileInfo struct {
	Path         string
	RelativePath string
	Size         int64
	ModTime      time.Time
	IsDir        bool
	Checksum     string
}

// BackupService defines the interface for backup operations
type BackupService interface {
	CreateBackup(opts BackupOptions) (*BackupMetadata, error)
	RestoreBackup(backupID string, targetDir string) error
	ListBackups() ([]BackupMetadata, error)
	GetBackupMetadata(backupID string) (*BackupMetadata, error)
	DeleteBackup(backupID string) error
	CleanupOldBackups() error
}

// BackupStorage defines the interface for backup storage operations
type BackupStorage interface {
	CreateWriter(metadata *BackupMetadata) (BackupWriter, error)
	OpenReader(backupID string) (BackupReader, error)
	Delete(backupID string) error
	List() ([]BackupMetadata, error)
}

// BackupWriter defines the interface for writing backup data
type BackupWriter interface {
	AddFile(info FileInfo, content io.Reader) error
	WriteMetadata(metadata *BackupMetadata) error
	Close() error
}

// BackupReader defines the interface for reading backup data
type BackupReader interface {
	GetFile(path string) (io.ReadCloser, FileInfo, error)
	ListFiles() ([]FileInfo, error)
	GetMetadata() (*BackupMetadata, error)
	Close() error
}
