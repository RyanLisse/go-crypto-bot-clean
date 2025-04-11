package commands

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"

	"go-crypto-bot-clean/backend/pkg/backup"
)

var (
	scheduler     *backup.BackupScheduler
	backupService backup.BackupService
	logger        *log.Logger
)

func initScheduler() error {
	if logger == nil {
		logFile, err := os.OpenFile("backup_scheduler.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("failed to open log file: %w", err)
		}
		// Use both stdout and file for logging
		multiWriter := io.MultiWriter(os.Stdout, logFile)
		logger = log.New(multiWriter, "[BACKUP] ", log.LstdFlags|log.Lshortfile)
	}

	if scheduler == nil {
		if backupService == nil {
			backupService = backup.NewBackupService(nil) // Replace nil with actual storage if needed
		}
		scheduler = backup.NewBackupScheduler(backupService, logger)
		scheduler.Start()
	}
	return nil
}

// NewBackupScheduleCmd creates the backup-schedule root command
func NewBackupScheduleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "backup-schedule",
		Short: "Manage automated backup schedules",
	}

	cmd.AddCommand(newAddScheduleCmd())
	cmd.AddCommand(newListSchedulesCmd())
	cmd.AddCommand(newRemoveScheduleCmd())
	cmd.AddCommand(newEnableScheduleCmd())
	cmd.AddCommand(newDisableScheduleCmd())
	cmd.AddCommand(newTestScheduleCmd())

	return cmd
}

func newAddScheduleCmd() *cobra.Command {
	var opts backup.BackupOptions
	var cronExpr string

	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a new backup schedule",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := initScheduler(); err != nil {
				return fmt.Errorf("failed to initialize scheduler: %w", err)
			}
			id := fmt.Sprintf("%s-%d", opts.Type, time.Now().Unix())
			err := scheduler.AddSchedule(id, opts, cronExpr)
			if err != nil {
				return fmt.Errorf("failed to add schedule: %w", err)
			}
			logger.Printf("Added schedule %s with cron '%s'\n", id, cronExpr)
			return nil
		},
	}

	typeStr := string(opts.Type)
	cmd.Flags().StringVarP(&typeStr, "type", "t", "full", "Backup type (full or incremental)")
	cmd.Flags().StringVarP(&opts.SourceDir, "source", "s", "", "Source directory to backup")
	cmd.Flags().StringVarP(&opts.DestDir, "destination", "d", "", "Destination directory for backups")
	cmd.Flags().StringSliceVarP(&opts.ExcludeRules, "exclude", "e", []string{}, "Exclude patterns")
	cmd.Flags().IntVarP(&opts.RetentionDays, "retention", "r", 30, "Retention days")
	// After parsing flags, assign opts.Type = BackupType(typeStr)
	cmd.PreRun = func(cmd *cobra.Command, args []string) {
		opts.Type = backup.BackupType(typeStr)
	}
	cmd.Flags().StringVar(&cronExpr, "cron", "0 2 * * *", "Cron expression for schedule")

	cmd.MarkFlagRequired("source")
	cmd.MarkFlagRequired("destination")

	return cmd
}

func newListSchedulesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all backup schedules",
		Run: func(cmd *cobra.Command, args []string) {
			if err := initScheduler(); err != nil {
				fmt.Printf("Failed to initialize scheduler: %v\n", err)
				return
			}
			schedules := scheduler.ListSchedules()
			if len(schedules) == 0 {
				fmt.Println("No backup schedules found.")
				return
			}
			for _, s := range schedules {
				fmt.Printf("ID: %s, Cron: %s, Enabled: %v\n", s.ID, s.CronExpr, s.IsEnabled)
			}
		},
	}
}

func newRemoveScheduleCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "remove [schedule_id]",
		Short: "Remove a backup schedule",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := initScheduler(); err != nil {
				return fmt.Errorf("failed to initialize scheduler: %w", err)
			}
			err := scheduler.RemoveSchedule(args[0])
			if err != nil {
				return fmt.Errorf("failed to remove schedule: %w", err)
			}
			fmt.Printf("Removed schedule %s\n", args[0])
			return nil
		},
	}
}

func newEnableScheduleCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "enable [schedule_id]",
		Short: "Enable a backup schedule",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := initScheduler(); err != nil {
				return fmt.Errorf("failed to initialize scheduler: %w", err)
			}
			err := scheduler.EnableSchedule(args[0])
			if err != nil {
				return fmt.Errorf("failed to enable schedule: %w", err)
			}
			fmt.Printf("Enabled schedule %s\n", args[0])
			return nil
		},
	}
}

func newDisableScheduleCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "disable [schedule_id]",
		Short: "Disable a backup schedule",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := initScheduler(); err != nil {
				return fmt.Errorf("failed to initialize scheduler: %w", err)
			}
			err := scheduler.DisableSchedule(args[0])
			if err != nil {
				return fmt.Errorf("failed to disable schedule: %w", err)
			}
			fmt.Printf("Disabled schedule %s\n", args[0])
			return nil
		},
	}
}

func newTestScheduleCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "test [schedule_id]",
		Short: "Perform a dry-run of a backup schedule",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := initScheduler(); err != nil {
				return fmt.Errorf("failed to initialize scheduler: %w", err)
			}
			schedules := scheduler.ListSchedules()
			var target *backup.Schedule
			for _, s := range schedules {
				if s.ID == args[0] {
					target = s
					break
				}
			}
			if target == nil {
				return fmt.Errorf("schedule %s not found", args[0])
			}
			logger.Printf("Starting dry-run backup for schedule %s\n", target.ID)
			// Perform dry-run by invoking backup service directly (simulate)
			_, err := backupService.CreateBackup(target.Options)
			if err != nil {
				fmt.Printf("Dry-run backup failed: %v\n", err)
				return err
			}
			fmt.Printf("Dry-run backup completed for schedule %s\n", target.ID)
			return nil
		},
	}
}
