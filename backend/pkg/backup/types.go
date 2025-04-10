package backup

import (
	"io"
	"time"
)

// BackupType represents the type of backup (full or incremental)
type BackupType string

const (
	// FullBackup represents a complete backup of all files
	FullBackup BackupType = "full"
	// IncrementalBackup represents a backup of only changed files since last backup
	IncrementalBackup BackupType = "incremental"
)

// FileInfo represents information about a file in the backup
type FileInfo struct {
	Path         string    `json:"path"`
	Size         int64     `json:"size"`
	ModTime      time.Time `json:"modTime"`
	IsDir        bool      `json:"isDir"`
	Checksum     string    `json:"checksum,omitempty"`
	RelativePath string    `json:"relativePath"`
}

// BackupMetadata contains information about a backup
type BackupMetadata struct {
	ID            string     `json:"id"`
	Type          BackupType `json:"type"`
	SourceDir     string     `json:"sourceDir"`
	DestDir       string     `json:"destDir"`
	Timestamp     time.Time  `json:"timestamp"`
	Files         []FileInfo `json:"files"`
	TotalSize     int64      `json:"totalSize"`
	FileCount     int        `json:"fileCount"`
	ExcludeRules  []string   `json:"excludeRules,omitempty"`
	Checksum      string     `json:"checksum,omitempty"`
	ParentBackup  string     `json:"parentBackup,omitempty"`
	RetentionDays int        `json:"retentionDays"`
}

// BackupService defines the interface for backup operations
type BackupService interface {
	// CreateBackup performs a backup operation and returns metadata about the backup
	CreateBackup(opts BackupOptions) (*BackupMetadata, error)

	// RestoreBackup restores files from a backup
	RestoreBackup(backupID string, targetDir string) error

	// ListBackups returns a list of all available backups
	ListBackups() ([]BackupMetadata, error)

	// GetBackupMetadata retrieves metadata for a specific backup
	GetBackupMetadata(backupID string) (*BackupMetadata, error)

	// DeleteBackup removes a backup and its associated files
	DeleteBackup(backupID string) error

	// CleanupOldBackups removes backups that are older than their retention period
	CleanupOldBackups() error

	// VerifyBackup checks the integrity of a backup
	VerifyBackup(backupID string) error
}

// BackupOptions contains options for creating a backup
type BackupOptions struct {
	Type          BackupType
	SourceDir     string
	DestDir       string
	ExcludeRules  []string
	RetentionDays int
}

// BackupWriter defines the interface for writing backup data
type BackupWriter interface {
	io.WriteCloser
	AddFile(info FileInfo, reader io.Reader) error
	WriteMetadata(metadata *BackupMetadata) error
}

// BackupReader defines the interface for reading backup data
type BackupReader interface {
	io.ReadCloser
	GetFile(path string) (io.ReadCloser, *FileInfo, error)
	GetMetadata() (*BackupMetadata, error)
	ListFiles() ([]FileInfo, error)
}

// BackupStorage defines the interface for backup storage operations
type BackupStorage interface {
	// CreateWriter creates a new backup writer
	CreateWriter(metadata *BackupMetadata) (BackupWriter, error)

	// OpenReader opens a backup for reading
	OpenReader(backupID string) (BackupReader, error)

	// Delete removes a backup and its associated files
	Delete(backupID string) error

	// List returns a list of all available backups
	List() ([]BackupMetadata, error)
}
