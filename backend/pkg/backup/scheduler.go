package backup

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

// Schedule represents a backup schedule configuration
type Schedule struct {
	ID          string
	Options     BackupOptions
	CronExpr    string        // Cron expression for flexible scheduling
	Interval    time.Duration // Optional: for backward compatibility
	LastRunTime time.Time
	NextRunTime time.Time
	IsEnabled   bool
	cronID      cron.EntryID // Internal: ID of the cron job
}

// BackupScheduler manages automated backup schedules
type BackupScheduler struct {
	service   BackupService
	logger    *log.Logger
	schedules map[string]*Schedule
	cron      *cron.Cron
	stopChan  chan struct{}
	wg        sync.WaitGroup
	mu        sync.RWMutex
}

// NewBackupScheduler creates a new backup scheduler
func NewBackupScheduler(service BackupService, logger *log.Logger) *BackupScheduler {
	if logger == nil {
		logger = log.New(os.Stdout, "[BACKUP] ", log.LstdFlags)
	}

	return &BackupScheduler{
		service:   service,
		logger:    logger,
		schedules: make(map[string]*Schedule),
		cron:      cron.New(cron.WithSeconds()),
		stopChan:  make(chan struct{}),
	}
}

// AddSchedule adds a new backup schedule using a cron expression
func (s *BackupScheduler) AddSchedule(id string, opts BackupOptions, cronExpr string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.schedules[id]; exists {
		return fmt.Errorf("schedule with ID %s already exists", id)
	}

	_, err := cron.ParseStandard(cronExpr)
	if err != nil {
		return fmt.Errorf("invalid cron expression: %v", err)
	}

	schedule := &Schedule{
		ID:        id,
		Options:   opts,
		CronExpr:  cronExpr,
		IsEnabled: true,
	}

	entryID, err := s.cron.AddFunc(cronExpr, func() {
		if err := s.executeBackup(schedule); err != nil {
			s.logger.Printf("Failed to execute backup for schedule %s: %v", id, err)
		}
	})
	if err != nil {
		return fmt.Errorf("failed to add cron job: %v", err)
	}

	schedule.cronID = entryID
	s.schedules[id] = schedule
	s.logger.Printf("Added schedule %s with cron expression %s", id, cronExpr)
	return nil
}

// AddScheduleWithInterval adds a new backup schedule using a time interval (for backward compatibility)
func (s *BackupScheduler) AddScheduleWithInterval(id string, opts BackupOptions, interval time.Duration) error {
	// Convert interval to cron expression (run every X duration)
	cronExpr := fmt.Sprintf("@every %s", interval.String())
	return s.AddSchedule(id, opts, cronExpr)
}

// RemoveSchedule removes a backup schedule
func (s *BackupScheduler) RemoveSchedule(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	schedule, exists := s.schedules[id]
	if !exists {
		return fmt.Errorf("schedule with ID %s not found", id)
	}

	s.cron.Remove(schedule.cronID)
	delete(s.schedules, id)
	s.logger.Printf("Removed schedule %s", id)
	return nil
}

// Start starts the scheduler
func (s *BackupScheduler) Start() {
	s.logger.Println("Starting backup scheduler")
	s.cron.Start()
}

// Stop stops the scheduler
func (s *BackupScheduler) Stop() {
	s.logger.Println("Stopping backup scheduler")
	s.cron.Stop()
}

// executeBackup performs a scheduled backup
func (s *BackupScheduler) executeBackup(schedule *Schedule) error {
	s.logger.Printf("Starting backup for schedule %s", schedule.ID)

	opts := BackupOptions{
		Type:          schedule.Options.Type,
		SourceDir:     schedule.Options.SourceDir,
		DestDir:       schedule.Options.DestDir,
		ExcludeRules:  schedule.Options.ExcludeRules,
		RetentionDays: schedule.Options.RetentionDays,
	}

	metadata, err := s.service.CreateBackup(opts)
	if err != nil {
		s.logger.Printf("Backup failed for schedule %s: %v", schedule.ID, err)
		return err
	}

	s.logger.Printf("Successfully completed backup %s for schedule %s", metadata.ID, schedule.ID)
	return nil
}

// ListSchedules returns a list of all backup schedules
func (s *BackupScheduler) ListSchedules() []*Schedule {
	s.mu.RLock()
	defer s.mu.RUnlock()

	schedules := make([]*Schedule, 0, len(s.schedules))
	for _, schedule := range s.schedules {
		schedules = append(schedules, schedule)
	}
	return schedules
}

// EnableSchedule enables a backup schedule
func (s *BackupScheduler) EnableSchedule(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	schedule, exists := s.schedules[id]
	if !exists {
		return fmt.Errorf("schedule with ID %s not found", id)
	}

	if schedule.IsEnabled {
		return nil
	}

	entryID, err := s.cron.AddFunc(schedule.CronExpr, func() {
		if err := s.executeBackup(schedule); err != nil {
			s.logger.Printf("Failed to execute backup for schedule %s: %v", id, err)
		}
	})
	if err != nil {
		return fmt.Errorf("failed to add cron job: %v", err)
	}

	schedule.cronID = entryID
	schedule.IsEnabled = true
	s.logger.Printf("Enabled schedule %s", id)
	return nil
}

// DisableSchedule disables a backup schedule
func (s *BackupScheduler) DisableSchedule(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	schedule, exists := s.schedules[id]
	if !exists {
		return fmt.Errorf("schedule with ID %s not found", id)
	}

	if !schedule.IsEnabled {
		return nil
	}

	s.cron.Remove(schedule.cronID)
	schedule.IsEnabled = false
	s.logger.Printf("Disabled schedule %s", id)
	return nil
}
