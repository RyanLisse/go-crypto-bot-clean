package backup

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

// SecurityConfig holds the configuration for secure backup storage
type SecurityConfig struct {
	// EncryptionKey is the key used for encrypting backup data
	EncryptionKey []byte
	// AccessToken is used for authenticating backup operations
	AccessToken string
	// AllowedPaths contains the paths that are allowed for backup operations
	AllowedPaths []string
	// RetentionPolicy defines how long backups should be kept
	RetentionPolicy RetentionPolicy
}

// RetentionPolicy defines the rules for backup retention
type RetentionPolicy struct {
	// Days to keep backups (0 means keep forever)
	Days int
	// MaxBackups is the maximum number of backups to keep (0 means unlimited)
	MaxBackups int
	// Strategy determines how to select which backups to remove
	Strategy RetentionStrategy
	// MinimumBackups is the minimum number of backups to keep regardless of age
	MinimumBackups int
}

// RetentionStrategy determines how to select backups for removal
type RetentionStrategy string

const (
	// OldestFirst removes the oldest backups first
	OldestFirst RetentionStrategy = "oldest_first"
	// LargestFirst removes the largest backups first
	LargestFirst RetentionStrategy = "largest_first"
	// SelectiveRetention keeps important backups longer
	SelectiveRetention RetentionStrategy = "selective"
)

// secureBackupStorage implements BackupStorage with encryption and access control
type secureBackupStorage struct {
	baseStorage BackupStorage
	config      SecurityConfig
	cipher      cipher.AEAD
	mu          sync.RWMutex
}

// NewSecureBackupStorage creates a new secure backup storage wrapper
func NewSecureBackupStorage(baseStorage BackupStorage, config SecurityConfig) (BackupStorage, error) {
	if len(config.EncryptionKey) == 0 {
		return nil, fmt.Errorf("encryption key is required")
	}

	// Create AES cipher
	block, err := aes.NewCipher(config.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	return &secureBackupStorage{
		baseStorage: baseStorage,
		config:      config,
		cipher:      gcm,
	}, nil
}

// encryptData encrypts data using AES-GCM
func (s *secureBackupStorage) encryptData(data []byte) ([]byte, error) {
	nonce := make([]byte, s.cipher.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt data
	encrypted := s.cipher.Seal(nonce, nonce, data, nil)
	return encrypted, nil
}

// decryptData decrypts data using AES-GCM
func (s *secureBackupStorage) decryptData(data []byte) ([]byte, error) {
	if len(data) < s.cipher.NonceSize() {
		return nil, fmt.Errorf("encrypted data too short")
	}

	nonce := data[:s.cipher.NonceSize()]
	ciphertext := data[s.cipher.NonceSize():]

	// Decrypt data
	decrypted, err := s.cipher.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt data: %w", err)
	}

	return decrypted, nil
}

// validateAccess checks if the operation is allowed based on access token and paths
func (s *secureBackupStorage) validateAccess(path string) error {
	// Check if path is in allowed paths
	if len(s.config.AllowedPaths) > 0 {
		allowed := false
		for _, allowedPath := range s.config.AllowedPaths {
			if isSubPath(path, allowedPath) {
				allowed = true
				break
			}
		}
		if !allowed {
			return fmt.Errorf("access denied: path not allowed")
		}
	}

	// Check access token if set
	if s.config.AccessToken != "" {
		// Additional token-based validation could be added here
	}

	return nil
}

// isSubPath checks if child path is a subdirectory of parent path
func isSubPath(child, parent string) bool {
	childAbs, err := filepath.Abs(child)
	if err != nil {
		return false
	}
	parentAbs, err := filepath.Abs(parent)
	if err != nil {
		return false
	}
	rel, err := filepath.Rel(parentAbs, childAbs)
	if err != nil {
		return false
	}
	return !filepath.IsAbs(rel) && !strings.HasPrefix(rel, "..")
}

// CreateWriter implements BackupStorage interface
func (s *secureBackupStorage) CreateWriter(metadata *BackupMetadata) (BackupWriter, error) {
	// Check if the source directory is allowed
	if err := s.validateAccess(metadata.SourceDir); err != nil {
		return nil, err
	}

	// Create writer from base storage
	writer, err := s.baseStorage.CreateWriter(metadata)
	if err != nil {
		return nil, err
	}

	return &secureBackupWriter{
		baseWriter: writer,
		storage:    s,
	}, nil
}

// OpenReader implements BackupStorage interface with decryption
func (s *secureBackupStorage) OpenReader(backupID string) (BackupReader, error) {
	reader, err := s.baseStorage.OpenReader(backupID)
	if err != nil {
		return nil, err
	}

	return &secureBackupReader{
		baseReader: reader,
		storage:    s,
	}, nil
}

// Delete implements BackupStorage interface with access control
func (s *secureBackupStorage) Delete(backupID string) error {
	metadata, err := s.baseStorage.OpenReader(backupID)
	if err != nil {
		return err
	}
	defer metadata.Close()

	meta, err := metadata.GetMetadata()
	if err != nil {
		return err
	}

	if err := s.validateAccess(meta.SourceDir); err != nil {
		return err
	}

	return s.baseStorage.Delete(backupID)
}

// List implements BackupStorage interface
func (s *secureBackupStorage) List() ([]BackupMetadata, error) {
	return s.baseStorage.List()
}

// ApplyRetentionPolicy implements BackupStorage interface
func (s *secureBackupStorage) ApplyRetentionPolicy() error {
	backups, err := s.List()
	if err != nil {
		return fmt.Errorf("failed to list backups: %w", err)
	}

	// Sort backups by timestamp (newest first)
	sort.Slice(backups, func(i, j int) bool {
		return backups[i].Timestamp.After(backups[j].Timestamp)
	})

	// Keep track of what we've retained
	retained := make(map[string]bool)
	toDelete := make([]string, 0)

	// Always keep the minimum number of backups
	for i := 0; i < len(backups) && i < s.config.RetentionPolicy.MinimumBackups; i++ {
		retained[backups[i].ID] = true
	}

	// Apply retention strategy
	if s.config.RetentionPolicy.Strategy == SelectiveRetention {
		now := time.Now()

		// Keep daily backups for the last week
		dailyCount := 0
		for _, backup := range backups {
			age := now.Sub(backup.Timestamp)
			if age <= 7*24*time.Hour && dailyCount < 7 {
				retained[backup.ID] = true
				dailyCount++
			}
		}

		// Keep weekly backups for the last month
		weeklyCount := 0
		for _, backup := range backups {
			age := now.Sub(backup.Timestamp)
			if age <= 30*24*time.Hour && age > 7*24*time.Hour && weeklyCount < 4 {
				retained[backup.ID] = true
				weeklyCount++
			}
		}

		// Keep monthly backups for the last year
		monthlyCount := 0
		for _, backup := range backups {
			age := now.Sub(backup.Timestamp)
			if age <= 365*24*time.Hour && age > 30*24*time.Hour && monthlyCount < 12 {
				retained[backup.ID] = true
				monthlyCount++
			}
		}
	} else {
		// Keep backups within retention period
		now := time.Now()
		for _, backup := range backups {
			age := now.Sub(backup.Timestamp)
			if age <= time.Duration(s.config.RetentionPolicy.Days)*24*time.Hour && len(retained) < s.config.RetentionPolicy.MaxBackups {
				retained[backup.ID] = true
			}
		}
	}

	// Mark backups for deletion
	for _, backup := range backups {
		if !retained[backup.ID] {
			toDelete = append(toDelete, backup.ID)
		}
	}

	// Delete backups
	for _, id := range toDelete {
		if err := s.Delete(id); err != nil {
			return fmt.Errorf("failed to delete backup %s: %w", id, err)
		}
	}

	return nil
}

// secureBackupWriter implements BackupWriter with encryption
type secureBackupWriter struct {
	baseWriter BackupWriter
	storage    *secureBackupStorage
}

// AddFile implements BackupWriter interface with encryption
func (w *secureBackupWriter) AddFile(info FileInfo, reader io.Reader) error {
	// Read all data
	data, err := io.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("failed to read file data: %w", err)
	}

	// Encrypt data
	encrypted, err := w.storage.encryptData(data)
	if err != nil {
		return fmt.Errorf("failed to encrypt file data: %w", err)
	}

	// Write encrypted data
	return w.baseWriter.AddFile(info, bytes.NewReader(encrypted))
}

func (w *secureBackupWriter) WriteMetadata(metadata *BackupMetadata) error {
	return w.baseWriter.WriteMetadata(metadata)
}

func (w *secureBackupWriter) Close() error {
	return w.baseWriter.Close()
}

// secureBackupReader implements BackupReader with decryption
type secureBackupReader struct {
	baseReader BackupReader
	storage    *secureBackupStorage
}

// GetFile implements BackupReader interface with decryption
func (r *secureBackupReader) GetFile(path string) (io.ReadCloser, FileInfo, error) {
	reader, fileInfo, err := r.baseReader.GetFile(path)
	if err != nil {
		return nil, FileInfo{}, err
	}

	// Read encrypted data
	encrypted, err := io.ReadAll(reader)
	if err != nil {
		return nil, FileInfo{}, fmt.Errorf("failed to read encrypted data: %w", err)
	}

	// Decrypt data
	decrypted, err := r.storage.decryptData(encrypted)
	if err != nil {
		return nil, FileInfo{}, fmt.Errorf("failed to decrypt file data: %w", err)
	}

	// Update file info with original size
	fileInfo.Size = int64(len(decrypted))

	return io.NopCloser(bytes.NewReader(decrypted)), fileInfo, nil
}

func (r *secureBackupReader) ListFiles() ([]FileInfo, error) {
	return r.baseReader.ListFiles()
}

func (r *secureBackupReader) GetMetadata() (*BackupMetadata, error) {
	return r.baseReader.GetMetadata()
}

func (r *secureBackupReader) Close() error {
	return r.baseReader.Close()
}
