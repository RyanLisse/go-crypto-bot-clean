package backup

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// localBackupStorage implements BackupStorage using the local filesystem
type localBackupStorage struct {
	baseDir string
}

// NewLocalBackupStorage creates a new local backup storage
func NewLocalBackupStorage(baseDir string) (BackupStorage, error) {
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create base directory: %w", err)
	}

	return &localBackupStorage{
		baseDir: baseDir,
	}, nil
}

// CreateWriter implements BackupStorage interface
func (s *localBackupStorage) CreateWriter(metadata *BackupMetadata) (BackupWriter, error) {
	if metadata.ID == "" {
		metadata.ID = fmt.Sprintf("%d", time.Now().UnixNano())
	}
	if metadata.Timestamp.IsZero() {
		metadata.Timestamp = time.Now()
	}

	backupDir := filepath.Join(s.baseDir, metadata.ID)
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create backup directory: %w", err)
	}

	return &localBackupWriter{
		storage:   s,
		metadata:  metadata,
		backupDir: backupDir,
	}, nil
}

// OpenReader implements BackupStorage interface
func (s *localBackupStorage) OpenReader(backupID string) (BackupReader, error) {
	backupDir := filepath.Join(s.baseDir, backupID)
	if _, err := os.Stat(backupDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("backup not found: %s", backupID)
	}

	return &localBackupReader{
		storage:   s,
		backupDir: backupDir,
	}, nil
}

// Delete implements BackupStorage interface
func (s *localBackupStorage) Delete(backupID string) error {
	backupDir := filepath.Join(s.baseDir, backupID)
	return os.RemoveAll(backupDir)
}

// List implements BackupStorage interface
func (s *localBackupStorage) List() ([]BackupMetadata, error) {
	entries, err := os.ReadDir(s.baseDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	var backups []BackupMetadata
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		reader, err := s.OpenReader(entry.Name())
		if err != nil {
			continue
		}
		defer reader.Close()

		metadata, err := reader.GetMetadata()
		if err != nil {
			continue
		}

		backups = append(backups, *metadata)
	}

	return backups, nil
}

// ApplyRetentionPolicy implements BackupStorage interface
func (s *localBackupStorage) ApplyRetentionPolicy() error {
	// For local storage, we don't implement retention policy
	// This should be handled by the secure storage wrapper
	return nil
}

// localBackupWriter implements BackupWriter for local storage
type localBackupWriter struct {
	storage   *localBackupStorage
	metadata  *BackupMetadata
	backupDir string
}

// AddFile implements BackupWriter interface
func (w *localBackupWriter) AddFile(info FileInfo, reader io.Reader) error {
	filePath := filepath.Join(w.backupDir, info.Path)
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	if _, err := io.Copy(file, reader); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// WriteMetadata implements BackupWriter interface
func (w *localBackupWriter) WriteMetadata(metadata *BackupMetadata) error {
	metadataPath := filepath.Join(w.backupDir, "metadata.json")
	file, err := os.Create(metadataPath)
	if err != nil {
		return fmt.Errorf("failed to create metadata file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(metadata); err != nil {
		return fmt.Errorf("failed to encode metadata: %w", err)
	}

	return nil
}

func (w *localBackupWriter) Write(p []byte) (n int, err error) {
	// Not implemented for local storage
	return 0, fmt.Errorf("direct writing not supported for local storage")
}

func (w *localBackupWriter) Close() error {
	return w.WriteMetadata(w.metadata)
}

// localBackupReader implements BackupReader for local storage
type localBackupReader struct {
	storage   *localBackupStorage
	backupDir string
}

// GetFile implements BackupReader interface
func (r *localBackupReader) GetFile(path string) (io.ReadCloser, FileInfo, error) {
	filePath := filepath.Join(r.backupDir, path)
	file, err := os.Open(filePath)
	if err != nil {
		return nil, FileInfo{}, fmt.Errorf("failed to open file: %w", err)
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, FileInfo{}, fmt.Errorf("failed to get file info: %w", err)
	}

	fileInfo := FileInfo{
		Path:    path,
		Size:    info.Size(),
		ModTime: info.ModTime(),
	}

	return file, fileInfo, nil
}

// GetMetadata implements BackupReader interface
func (r *localBackupReader) GetMetadata() (*BackupMetadata, error) {
	metadataPath := filepath.Join(r.backupDir, "metadata.json")
	file, err := os.Open(metadataPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open metadata file: %w", err)
	}
	defer file.Close()

	var metadata BackupMetadata
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&metadata); err != nil {
		return nil, fmt.Errorf("failed to decode metadata: %w", err)
	}

	return &metadata, nil
}

// ListFiles implements BackupReader interface
func (r *localBackupReader) ListFiles() ([]FileInfo, error) {
	var files []FileInfo
	err := filepath.Walk(r.backupDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || info.Name() == "metadata.json" {
			return nil
		}

		relPath, err := filepath.Rel(r.backupDir, path)
		if err != nil {
			return err
		}

		files = append(files, FileInfo{
			Path:    relPath,
			Size:    info.Size(),
			ModTime: info.ModTime(),
		})
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to list files: %w", err)
	}

	return files, nil
}

func (r *localBackupReader) Read(p []byte) (n int, err error) {
	// Not implemented for local storage
	return 0, fmt.Errorf("direct reading not supported for local storage")
}

func (r *localBackupReader) Close() error {
	return nil
}
