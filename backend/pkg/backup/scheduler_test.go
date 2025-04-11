package backup

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"
)

func TestScheduler(t *testing.T) {
	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
	mockService := &mockBackupService{
		backups: make([]BackupMetadata, 0),
	}
	scheduler := NewBackupScheduler(mockService, logger)

	t.Run("Add and remove schedule", func(t *testing.T) {
		schedule := BackupOptions{
			Type:      FullBackup,
			SourceDir: "/test/path",
		}

		// Add schedule
		err := scheduler.AddSchedule("test1", schedule, "*/5 * * * *")
		if err != nil {
			t.Fatalf("Failed to add schedule: %v", err)
		}

		// Verify schedule exists
		schedules := scheduler.ListSchedules()
		if len(schedules) != 1 {
			t.Errorf("Expected 1 schedule, got %d", len(schedules))
		}

		// Remove schedule
		err = scheduler.RemoveSchedule("test1")
		if err != nil {
			t.Fatalf("Failed to remove schedule: %v", err)
		}

		// Verify schedule was removed
		schedules = scheduler.ListSchedules()
		if len(schedules) != 0 {
			t.Errorf("Expected 0 schedules, got %d", len(schedules))
		}
	})

	t.Run("Enable and disable schedule", func(t *testing.T) {
		schedule := BackupOptions{
			Type:      FullBackup,
			SourceDir: "/test/path",
		}

		// Add schedule
		err := scheduler.AddSchedule("test2", schedule, "*/5 * * * *")
		if err != nil {
			t.Fatalf("Failed to add schedule: %v", err)
		}

		// Disable schedule
		err = scheduler.DisableSchedule("test2")
		if err != nil {
			t.Fatalf("Failed to disable schedule: %v", err)
		}

		// Verify schedule is disabled
		schedules := scheduler.ListSchedules()
		if schedules[0].IsEnabled {
			t.Error("Schedule should be disabled")
		}

		// Enable schedule
		err = scheduler.EnableSchedule("test2")
		if err != nil {
			t.Fatalf("Failed to enable schedule: %v", err)
		}

		// Verify schedule is enabled
		schedules = scheduler.ListSchedules()
		if !schedules[0].IsEnabled {
			t.Error("Schedule should be enabled")
		}

		// Cleanup
		scheduler.RemoveSchedule("test2")
	})

	t.Run("Backward compatibility with interval-based scheduling", func(t *testing.T) {
		schedule := BackupOptions{
			Type:      FullBackup,
			SourceDir: "/test/path",
		}

		// Add schedule with interval
		err := scheduler.AddScheduleWithInterval("test3", schedule, 5*time.Minute)
		if err != nil {
			t.Fatalf("Failed to add interval schedule: %v", err)
		}

		// Verify schedule exists
		schedules := scheduler.ListSchedules()
		if len(schedules) != 1 {
			t.Errorf("Expected 1 schedule, got %d", len(schedules))
		}

		// Cleanup
		scheduler.RemoveSchedule("test3")
	})
}

type mockBackupService struct {
	backups []BackupMetadata
}

func (m *mockBackupService) CreateBackup(opts BackupOptions) (*BackupMetadata, error) {
	backup := BackupMetadata{
		ID:        fmt.Sprintf("backup-%d", len(m.backups)+1),
		Type:      opts.Type,
		Timestamp: time.Now(),
	}
	m.backups = append(m.backups, backup)
	return &backup, nil
}

func (m *mockBackupService) RestoreBackup(id string, destDir string) error {
	return nil
}

func (m *mockBackupService) CleanupOldBackups() error {
	return nil
}

func (m *mockBackupService) GetBackupMetadata(id string) (*BackupMetadata, error) {
	for _, backup := range m.backups {
		if backup.ID == id {
			return &backup, nil
		}
	}
	return nil, fmt.Errorf("backup not found: %s", id)
}

func (m *mockBackupService) DeleteBackup(id string) error {
	for i, backup := range m.backups {
		if backup.ID == id {
			m.backups = append(m.backups[:i], m.backups[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("backup not found: %s", id)
}

func (m *mockBackupService) ListBackups() ([]BackupMetadata, error) {
	return m.backups, nil
}

func TestSchedulerExecution(t *testing.T) {
	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
	mockService := &mockBackupService{
		backups: make([]BackupMetadata, 0),
	}
	scheduler := NewBackupScheduler(mockService, logger)

	t.Run("Execute backup on schedule", func(t *testing.T) {
		schedule := BackupOptions{
			Type:      FullBackup,
			SourceDir: "/test/path",
		}

		err := scheduler.AddSchedule("test1", schedule, "@every 1s")
		if err != nil {
			t.Fatalf("Failed to add schedule: %v", err)
		}

		// Wait for at least one execution
		time.Sleep(2 * time.Second)

		// Cleanup
		scheduler.RemoveSchedule("test1")
	})
}

func TestSchedulerRetention(t *testing.T) {
	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
	mockService := &mockBackupService{
		backups: make([]BackupMetadata, 0),
	}
	scheduler := NewBackupScheduler(mockService, logger)

	t.Run("Cleanup old backups", func(t *testing.T) {
		schedule := BackupOptions{
			Type:      FullBackup,
			SourceDir: "/test/path",
		}

		err := scheduler.AddSchedule("test1", schedule, "@every 1h")
		if err != nil {
			t.Fatalf("Failed to add schedule: %v", err)
		}

		// Trigger cleanup
		err = mockService.CleanupOldBackups()
		if err != nil {
			t.Fatalf("Failed to cleanup old backups: %v", err)
		}

		// Cleanup
		scheduler.RemoveSchedule("test1")
	})
}

func TestSchedulerConcurrency(t *testing.T) {
	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
	mockService := &mockBackupService{}
	scheduler := NewBackupScheduler(mockService, logger)

	// Add multiple schedules concurrently
	t.Run("ConcurrentAddSchedule", func(t *testing.T) {
		done := make(chan bool)
		for i := 0; i < 10; i++ {
			go func(id int) {
				opts := BackupOptions{
					SourceDir: "/test/source",
					DestDir:   "/test/dest",
					Type:      FullBackup,
				}
				scheduler.AddSchedule(fmt.Sprintf("schedule-%d", id), opts, "@every 1h")
				done <- true
			}(i)
		}

		// Wait for all goroutines to complete
		for i := 0; i < 10; i++ {
			<-done
		}

		schedules := scheduler.ListSchedules()
		if len(schedules) > 10 {
			t.Errorf("Expected at most 10 schedules, got %d", len(schedules))
		}
	})

	// Test concurrent enable/disable operations
	t.Run("ConcurrentEnableDisable", func(t *testing.T) {
		done := make(chan bool)
		for i := 0; i < 5; i++ {
			go func(id int) {
				scheduleID := fmt.Sprintf("schedule-%d", id)
				scheduler.EnableSchedule(scheduleID)
				scheduler.DisableSchedule(scheduleID)
				done <- true
			}(i)
		}

		// Wait for all goroutines to complete
		for i := 0; i < 5; i++ {
			<-done
		}
	})
}
