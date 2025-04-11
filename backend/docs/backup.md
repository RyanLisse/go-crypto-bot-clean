# Crypto-Brute-Dash Backup System

## Overview

This document explains how to perform manual and automated backups of the backend system, how to restore from backups, configure backup schedules, and test the backup process.

---

## Manual Backups

Run a manual backup using the CLI:

```bash
./backend_app backup --type full --source /path/to/data --destination /path/to/backups --retention 30
```

### Options:

- `--type`: `full` or `incremental`
- `--source`: Source directory to back up
- `--destination`: Directory where backups are stored
- `--exclude`: Patterns to exclude (can be repeated)
- `--retention`: Days to retain backups (default 30)

Example:

```bash
./backend_app backup --type incremental --source ./data --destination ./backups --exclude "*.tmp" --retention 14
```

---

## Automated Backups

### Adding a Schedule

Create an automated backup schedule with a cron expression:

```bash
./backend_app backup-schedule add --type full --source ./data --destination ./backups --cron "0 2 * * *" --retention 30
```

- `--cron`: Cron expression (default: daily at 2 AM)
- Other options same as manual backup

### Listing Schedules

```bash
./backend_app backup-schedule list
```

### Enabling/Disabling Schedules

```bash
./backend_app backup-schedule enable <schedule_id>
./backend_app backup-schedule disable <schedule_id>
```

### Removing a Schedule

```bash
./backend_app backup-schedule remove <schedule_id>
```

---

## Dry-Run / Test Mode

To simulate a backup without affecting data:

```bash
./backend_app backup-schedule test <schedule_id>
```

This performs a dry-run backup using the schedule's configuration.

---

## Restore from Backup

To restore data from a backup archive:

1. Locate the backup archive `.tar.gz` file in your backup directory.
2. Extract it:

```bash
tar -xzvf backup-full-YYYYMMDD-HHMMSS.tar.gz -C /restore/target/directory
```

3. Alternatively, use the CLI restore command (if implemented):

```bash
./backend_app restore --backup-id <backup_id> --target /restore/target/directory
```

---

## Configuring Backup Frequency

- Use cron expressions when adding schedules to control frequency.
- Examples:
  - Daily at 2 AM: `"0 2 * * *"`
  - Weekly on Sunday at 3 AM: `"0 3 * * 0"`
  - Every 6 hours: `"0 */6 * * *"`

---

## Retention Policy

- Set `--retention` days to control how long backups are kept.
- Old backups beyond this period are automatically deleted.

---

## Notes

- Backups are compressed `.tar.gz` archives.
- Both full and incremental backups are supported.
- Exclusion patterns help skip temp files, logs, etc.
- Logs of backup operations are stored in the console output and can be redirected to files.
- The backup scheduler runs in the background once started.

---

## Disaster Recovery

- Regularly test backups using dry-run and restore procedures.
- Store backups in secure, redundant locations.
- Document recovery steps and keep this guide updated.