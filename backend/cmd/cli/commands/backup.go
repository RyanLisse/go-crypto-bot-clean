package commands

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// BackupOptions contains all the options for the backup command
type BackupOptions struct {
	Type        string
	Source      string
	Destination string
	Exclude     []string
	Retention   int
}

// NewBackupCmd creates a new backup command
func NewBackupCmd() *cobra.Command {
	opts := &BackupOptions{}

	cmd := &cobra.Command{
		Use:   "backup",
		Short: "Backup trading data",
		Long: `Backup trading data and configuration files.
Supports full and incremental backups with configurable retention period.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runBackup(opts)
		},
	}

	// Add flags
	cmd.Flags().StringVarP(&opts.Type, "type", "t", "full", "Backup type (full or incremental)")
	cmd.Flags().StringVarP(&opts.Source, "source", "s", "", "Source directory to backup")
	cmd.Flags().StringVarP(&opts.Destination, "destination", "d", "", "Destination directory for backups")
	cmd.MarkFlagRequired("source")
	cmd.MarkFlagRequired("destination")
	cmd.Flags().StringSliceVarP(&opts.Exclude, "exclude", "e", []string{}, "Patterns to exclude (can be specified multiple times)")
	cmd.Flags().IntVarP(&opts.Retention, "retention", "r", 30, "Number of days to retain backups")

	return cmd
}

func runBackup(opts *BackupOptions) error {
	// Setup logger
	logger, _ := zap.NewDevelopment()
	if !verbose {
		logger, _ = zap.NewProduction()
	}
	defer logger.Sync()

	logger.Info("Starting backup",
		zap.String("type", opts.Type),
		zap.String("source", opts.Source),
		zap.String("destination", opts.Destination),
		zap.Strings("exclude", opts.Exclude),
		zap.Int("retention", opts.Retention),
	)

	// Validate options
	if err := validateBackupOptions(opts); err != nil {
		return err
	}

	// Initialize backup service
	service := NewBackupService(logger)

	// Create backup context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Hour)
	defer cancel()

	// Start backup
	result, err := service.Backup(ctx, opts)
	if err != nil {
		return fmt.Errorf("backup failed: %w", err)
	}

	// Print backup summary
	printBackupSummary(result)

	return nil
}

func validateBackupOptions(opts *BackupOptions) error {
	// Validate backup type
	if opts.Type != "full" && opts.Type != "incremental" {
		return fmt.Errorf("invalid backup type: %s (must be 'full' or 'incremental')", opts.Type)
	}

	// Validate source directory
	srcInfo, err := os.Stat(opts.Source)
	if err != nil {
		return fmt.Errorf("invalid source directory: %w", err)
	}
	if !srcInfo.IsDir() {
		return fmt.Errorf("source path is not a directory: %s", opts.Source)
	}

	// Create destination directory if it doesn't exist
	err = os.MkdirAll(opts.Destination, 0755)
	if err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Validate retention period
	if opts.Retention < 1 {
		return fmt.Errorf("retention period must be at least 1 day")
	}

	return nil
}

// BackupResult contains information about the completed backup
type BackupResult struct {
	Type           string
	StartTime      time.Time
	EndTime        time.Time
	BackupFile     string
	FileCount      int
	TotalSize      int64
	CompressedSize int64
}

func printBackupSummary(result *BackupResult) {
	duration := result.EndTime.Sub(result.StartTime).Round(time.Second)
	fmt.Printf("\nBackup Summary:\n")
	fmt.Printf("Type: %s\n", result.Type)
	fmt.Printf("Start Time: %s\n", result.StartTime.Format(time.RFC1123))
	fmt.Printf("Duration: %s\n", duration)
	fmt.Printf("Backup File: %s\n", result.BackupFile)
	fmt.Printf("Files Processed: %d\n", result.FileCount)
	fmt.Printf("Total Size: %s\n", formatBytes(result.TotalSize))
	fmt.Printf("Compressed Size: %s\n", formatBytes(result.CompressedSize))
	fmt.Printf("Compression Ratio: %.1f%%\n", (1-float64(result.CompressedSize)/float64(result.TotalSize))*100)
}

// formatBytes converts bytes to a human-readable string
func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
