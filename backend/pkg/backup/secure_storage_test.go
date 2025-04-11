package backup

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"testing"
	"time"
)

func TestSecureBackupStorage(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "secure-backup-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Generate a random encryption key
	key := make([]byte, 32) // 256-bit key
	if _, err := rand.Read(key); err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	// Create base storage
	baseStorage, err := NewLocalBackupStorage(tempDir)
	if err != nil {
		t.Fatalf("Failed to create base storage: %v", err)
	}

	// Create secure storage
	config := SecurityConfig{
		EncryptionKey: key,
		AccessToken:   "test-token",
		AllowedPaths: []string{
			tempDir,
		},
		RetentionPolicy: RetentionPolicy{
			Days:           7,
			MaxBackups:     5,
			Strategy:       OldestFirst,
			MinimumBackups: 2,
		},
	}

	secureStorage, err := NewSecureBackupStorage(baseStorage, config)
	if err != nil {
		t.Fatalf("Failed to create secure storage: %v", err)
	}

	t.Run("TestBackupEncryption", func(t *testing.T) {
		testData := []byte("test data for encryption")
		metadata := &BackupMetadata{
			ID:        "test-backup-1",
			Timestamp: time.Now(),
			SourceDir: tempDir,
		}

		// Create a backup
		writer, err := secureStorage.CreateWriter(metadata)
		if err != nil {
			t.Fatalf("Failed to create writer: %v", err)
		}

		// Add a test file
		fileInfo := FileInfo{
			Path: "test.txt",
			Size: int64(len(testData)),
		}
		err = writer.AddFile(fileInfo, bytes.NewReader(testData))
		if err != nil {
			t.Fatalf("Failed to add file: %v", err)
		}

		err = writer.Close()
		if err != nil {
			t.Fatalf("Failed to close writer: %v", err)
		}

		// Read the backup
		reader, err := secureStorage.OpenReader(metadata.ID)
		if err != nil {
			t.Fatalf("Failed to open reader: %v", err)
		}
		defer reader.Close()

		// Get the file
		fileReader, _, err := reader.GetFile("test.txt")
		if err != nil {
			t.Fatalf("Failed to get file: %v", err)
		}
		defer fileReader.Close()

		// Read and verify the data
		data, err := io.ReadAll(fileReader)
		if err != nil {
			t.Fatalf("Failed to read file: %v", err)
		}

		if !bytes.Equal(data, testData) {
			t.Errorf("Data mismatch: got %q, want %q", data, testData)
		}
	})

	t.Run("TestAccessControl", func(t *testing.T) {
		// Create a directory outside of the allowed paths
		unauthorizedDir, err := os.MkdirTemp("", "unauthorized-backup-test")
		if err != nil {
			t.Fatalf("Failed to create unauthorized dir: %v", err)
		}
		defer os.RemoveAll(unauthorizedDir)

		metadata := &BackupMetadata{
			ID:        "test-backup-2",
			Timestamp: time.Now(),
			SourceDir: unauthorizedDir,
		}

		// Try to create a backup from unauthorized directory
		_, err = secureStorage.CreateWriter(metadata)
		if err == nil {
			t.Error("Expected access denied error, got nil")
		}
	})

	t.Run("TestRetentionPolicy", func(t *testing.T) {
		// Create multiple backups with different timestamps
		for i := 0; i < 10; i++ {
			metadata := &BackupMetadata{
				ID:        fmt.Sprintf("retention-test-%d", i),
				Timestamp: time.Now().AddDate(0, 0, -i), // Each backup is one day older
				SourceDir: tempDir,
			}

			writer, err := secureStorage.CreateWriter(metadata)
			if err != nil {
				t.Fatalf("Failed to create writer: %v", err)
			}

			err = writer.WriteMetadata(metadata)
			if err != nil {
				t.Fatalf("Failed to write metadata: %v", err)
			}

			err = writer.Close()
			if err != nil {
				t.Fatalf("Failed to close writer: %v", err)
			}
		}

		// Apply retention policy
		secureStore, ok := secureStorage.(*secureBackupStorage)
		if !ok {
			t.Fatal("Failed to cast to secureBackupStorage")
		}
		err := secureStore.ApplyRetentionPolicy()
		if err != nil {
			t.Fatalf("Failed to apply retention policy: %v", err)
		}

		// Check remaining backups
		backups, err := secureStorage.List()
		if err != nil {
			t.Fatalf("Failed to list backups: %v", err)
		}

		// Should have MaxBackups (5) or more if within retention period
		if len(backups) < config.RetentionPolicy.MinimumBackups {
			t.Errorf("Expected at least %d backups, got %d", config.RetentionPolicy.MinimumBackups, len(backups))
		}

		if len(backups) > config.RetentionPolicy.MaxBackups {
			t.Errorf("Expected at most %d backups, got %d", config.RetentionPolicy.MaxBackups, len(backups))
		}
	})

	t.Run("TestSelectiveRetention", func(t *testing.T) {
		// Create a new secure storage with selective retention
		selectiveConfig := config
		selectiveConfig.RetentionPolicy.Strategy = SelectiveRetention
		selectiveStorage, err := NewSecureBackupStorage(baseStorage, selectiveConfig)
		if err != nil {
			t.Fatalf("Failed to create selective storage: %v", err)
		}

		// Create backups with various ages
		timestamps := []time.Time{
			time.Now(),                    // Today
			time.Now().AddDate(0, 0, -2),  // 2 days ago
			time.Now().AddDate(0, 0, -8),  // 8 days ago
			time.Now().AddDate(0, 0, -15), // 15 days ago
			time.Now().AddDate(0, 0, -40), // 40 days ago
			time.Now().AddDate(0, -2, 0),  // 2 months ago
			time.Now().AddDate(-2, 0, 0),  // 2 years ago
		}

		for i, ts := range timestamps {
			metadata := &BackupMetadata{
				ID:        fmt.Sprintf("selective-test-%d", i),
				Timestamp: ts,
				SourceDir: tempDir,
			}

			writer, err := selectiveStorage.CreateWriter(metadata)
			if err != nil {
				t.Fatalf("Failed to create writer: %v", err)
			}

			err = writer.WriteMetadata(metadata)
			if err != nil {
				t.Fatalf("Failed to write metadata: %v", err)
			}

			err = writer.Close()
			if err != nil {
				t.Fatalf("Failed to close writer: %v", err)
			}
		}

		// Apply retention policy
		secureStore, ok := selectiveStorage.(*secureBackupStorage)
		if !ok {
			t.Fatal("Failed to cast to secureBackupStorage")
		}
		err = secureStore.ApplyRetentionPolicy()
		if err != nil {
			t.Fatalf("Failed to apply selective retention policy: %v", err)
		}

		// Check remaining backups
		backups, err := selectiveStorage.List()
		if err != nil {
			t.Fatalf("Failed to list backups: %v", err)
		}

		// Verify retention patterns
		now := time.Now()
		var (
			hasRecent   bool
			hasWeekly   bool
			hasMonthly  bool
			totalDaily  int
			totalWeekly int
		)

		for _, backup := range backups {
			age := now.Sub(backup.Timestamp)

			if age <= 24*time.Hour {
				hasRecent = true
			}
			if age > 7*24*time.Hour && age <= 30*24*time.Hour {
				hasWeekly = true
				totalWeekly++
			}
			if age > 30*24*time.Hour && age <= 365*24*time.Hour {
				hasMonthly = true
			}
			if age <= 7*24*time.Hour {
				totalDaily++
			}
		}

		if !hasRecent {
			t.Error("Expected to have recent backup")
		}
		if totalDaily > 7 {
			t.Errorf("Expected at most 7 daily backups, got %d", totalDaily)
		}
		if !hasWeekly {
			t.Error("Expected to have weekly backup")
		}
		if !hasMonthly {
			t.Error("Expected to have monthly backup")
		}
	})
}
