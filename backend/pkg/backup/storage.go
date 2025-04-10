package backup

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

// localBackupStorage implements BackupStorage using the local file system
type localBackupStorage struct {
	baseDir string
	mu      sync.RWMutex
}

// NewLocalBackupStorage creates a new instance of localBackupStorage
func NewLocalBackupStorage(baseDir string) (BackupStorage, error) {
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create backup directory: %w", err)
	}

	return &localBackupStorage{
		baseDir: baseDir,
	}, nil
}

// CreateWriter creates a new backup writer for the local file system
func (s *localBackupStorage) CreateWriter(metadata *BackupMetadata) (BackupWriter, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	backupDir := filepath.Join(s.baseDir, metadata.ID)
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create backup directory: %w", err)
	}

	archivePath := filepath.Join(backupDir, "backup.tar.gz")
	metadataPath := filepath.Join(backupDir, "metadata.json")

	archive, err := os.Create(archivePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create archive file: %w", err)
	}

	gzipWriter := gzip.NewWriter(archive)
	tarWriter := tar.NewWriter(gzipWriter)

	return &localBackupWriter{
		archive:      archive,
		gzipWriter:   gzipWriter,
		tarWriter:    tarWriter,
		metadataPath: metadataPath,
	}, nil
}

// OpenReader opens a backup for reading from the local file system
func (s *localBackupStorage) OpenReader(backupID string) (BackupReader, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	backupDir := filepath.Join(s.baseDir, backupID)
	archivePath := filepath.Join(backupDir, "backup.tar.gz")
	metadataPath := filepath.Join(backupDir, "metadata.json")

	archive, err := os.Open(archivePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open archive file: %w", err)
	}

	gzipReader, err := gzip.NewReader(archive)
	if err != nil {
		archive.Close()
		return nil, fmt.Errorf("failed to create gzip reader: %w", err)
	}

	tarReader := tar.NewReader(gzipReader)

	return &localBackupReader{
		archive:      archive,
		gzipReader:   gzipReader,
		tarReader:    tarReader,
		metadataPath: metadataPath,
	}, nil
}

// Delete removes a backup and its associated files from the local file system
func (s *localBackupStorage) Delete(backupID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	backupDir := filepath.Join(s.baseDir, backupID)
	return os.RemoveAll(backupDir)
}

// List returns a list of all available backups from the local file system
func (s *localBackupStorage) List() ([]BackupMetadata, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entries, err := os.ReadDir(s.baseDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read backup directory: %w", err)
	}

	var backups []BackupMetadata
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		metadataPath := filepath.Join(s.baseDir, entry.Name(), "metadata.json")
		metadata, err := readMetadataFile(metadataPath)
		if err != nil {
			continue // Skip invalid backups
		}

		backups = append(backups, *metadata)
	}

	return backups, nil
}

// localBackupWriter implements BackupWriter for the local file system
type localBackupWriter struct {
	archive      *os.File
	gzipWriter   *gzip.Writer
	tarWriter    *tar.Writer
	metadataPath string
}

func (w *localBackupWriter) AddFile(info FileInfo, reader io.Reader) error {
	header := &tar.Header{
		Name:    info.RelativePath,
		Size:    info.Size,
		Mode:    0644,
		ModTime: info.ModTime,
	}

	if err := w.tarWriter.WriteHeader(header); err != nil {
		return fmt.Errorf("failed to write tar header: %w", err)
	}

	if _, err := io.Copy(w.tarWriter, reader); err != nil {
		return fmt.Errorf("failed to write file content: %w", err)
	}

	return nil
}

func (w *localBackupWriter) WriteMetadata(metadata *BackupMetadata) error {
	data, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	if err := os.WriteFile(w.metadataPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write metadata file: %w", err)
	}

	return nil
}

func (w *localBackupWriter) Write(p []byte) (n int, err error) {
	return w.tarWriter.Write(p)
}

func (w *localBackupWriter) Close() error {
	if err := w.tarWriter.Close(); err != nil {
		return fmt.Errorf("failed to close tar writer: %w", err)
	}
	if err := w.gzipWriter.Close(); err != nil {
		return fmt.Errorf("failed to close gzip writer: %w", err)
	}
	return w.archive.Close()
}

// localBackupReader implements BackupReader for the local file system
type localBackupReader struct {
	archive      *os.File
	gzipReader   *gzip.Reader
	tarReader    *tar.Reader
	metadataPath string
}

// GetFile returns a reader for a specific file in the backup
func (r *localBackupReader) GetFile(path string) (io.ReadCloser, *FileInfo, error) {
	// Reset to beginning of archive
	if _, err := r.archive.Seek(0, io.SeekStart); err != nil {
		return nil, nil, fmt.Errorf("failed to seek archive: %w", err)
	}

	// Reset the gzip reader
	if err := r.gzipReader.Reset(r.archive); err != nil {
		return nil, nil, fmt.Errorf("failed to reset gzip reader: %w", err)
	}

	tr := tar.NewReader(r.gzipReader)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, nil, fmt.Errorf("failed to read tar header: %w", err)
		}

		if header.Name == path {
			info := &FileInfo{
				Path:         header.Name,
				Size:         header.Size,
				ModTime:      header.ModTime,
				IsDir:        header.Typeflag == tar.TypeDir,
				RelativePath: header.Name,
			}
			return ioutil.NopCloser(tr), info, nil
		}
	}

	return nil, nil, fmt.Errorf("file not found in backup: %s", path)
}

func (r *localBackupReader) GetMetadata() (*BackupMetadata, error) {
	return readMetadataFile(r.metadataPath)
}

func (r *localBackupReader) ListFiles() ([]FileInfo, error) {
	// Reset the tar reader to the beginning
	if _, err := r.archive.Seek(0, io.SeekStart); err != nil {
		return nil, fmt.Errorf("failed to seek archive: %w", err)
	}
	if err := r.gzipReader.Reset(r.archive); err != nil {
		return nil, fmt.Errorf("failed to reset gzip reader: %w", err)
	}

	var files []FileInfo
	for {
		header, err := r.tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read tar header: %w", err)
		}

		files = append(files, FileInfo{
			Path:         header.Name,
			RelativePath: header.Name,
			Size:         header.Size,
			ModTime:      header.ModTime,
			IsDir:        header.Typeflag == tar.TypeDir,
		})
	}

	return files, nil
}

func (r *localBackupReader) Read(p []byte) (n int, err error) {
	return r.tarReader.Read(p)
}

func (r *localBackupReader) Close() error {
	if err := r.gzipReader.Close(); err != nil {
		return fmt.Errorf("failed to close gzip reader: %w", err)
	}
	return r.archive.Close()
}

// Helper functions

func readMetadataFile(path string) (*BackupMetadata, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read metadata file: %w", err)
	}

	var metadata BackupMetadata
	if err := json.Unmarshal(data, &metadata); err != nil {
		return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
	}

	return &metadata, nil
}
