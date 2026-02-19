package backup

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

// Manager handles backup lifecycle and rotation
type Manager struct {
	baseBackupDir string
	maxBackups    int
}

// NewManager creates a new backup manager
func NewManager(baseBackupDir string, maxBackups int) *Manager {
	return &Manager{
		baseBackupDir: baseBackupDir,
		maxBackups:    maxBackups,
	}
}

// RotateBackups keeps only the most recent N backups and deletes older ones
func (m *Manager) RotateBackups(ctx context.Context) error {
	log.Debugf("Rotating backups (keeping last %d)", m.maxBackups)

	// Check if backup directory exists
	if _, err := os.Stat(m.baseBackupDir); os.IsNotExist(err) {
		log.Debug("Backup directory does not exist, nothing to rotate")
		return nil
	}

	// List all backup directories
	entries, err := os.ReadDir(m.baseBackupDir)
	if err != nil {
		return fmt.Errorf("failed to read backup directory: %w", err)
	}

	// Parse timestamps and sort
	type backupEntry struct {
		name      string
		timestamp time.Time
	}

	var backups []backupEntry
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		name := entry.Name()
		if !strings.HasPrefix(name, "backup-") {
			continue
		}

		timestamp, err := m.parseBackupTimestamp(name)
		if err != nil {
			log.WithError(err).Warnf("Failed to parse timestamp from backup: %s", name)
			continue
		}

		backups = append(backups, backupEntry{
			name:      name,
			timestamp: timestamp,
		})
	}

	// Sort by timestamp (newest first)
	sort.Slice(backups, func(i, j int) bool {
		return backups[i].timestamp.After(backups[j].timestamp)
	})

	// Delete old backups
	if len(backups) > m.maxBackups {
		backupsToDelete := backups[m.maxBackups:]
		log.Infof("Rotating out %d old backups", len(backupsToDelete))

		for _, backup := range backupsToDelete {
			backupPath := filepath.Join(m.baseBackupDir, backup.name)
			log.Debugf("Deleting old backup: %s", backupPath)

			if err := os.RemoveAll(backupPath); err != nil {
				log.WithError(err).Warnf("Failed to delete backup: %s", backup.name)
				// Continue with other deletions
			}
		}
	} else {
		log.Debugf("Only %d backups exist, no rotation needed", len(backups))
	}

	return nil
}

// ListBackups returns all available backups sorted by timestamp (newest first)
func (m *Manager) ListBackups() ([]*Manifest, error) {
	// Check if backup directory exists
	if _, err := os.Stat(m.baseBackupDir); os.IsNotExist(err) {
		return []*Manifest{}, nil
	}

	entries, err := os.ReadDir(m.baseBackupDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read backup directory: %w", err)
	}

	var manifests []*Manifest
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		name := entry.Name()
		if !strings.HasPrefix(name, "backup-") {
			continue
		}

		manifest, err := m.GetManifest(name)
		if err != nil {
			log.WithError(err).Warnf("Failed to read manifest for backup: %s", name)
			continue
		}

		manifests = append(manifests, manifest)
	}

	// Sort by timestamp (newest first)
	sort.Slice(manifests, func(i, j int) bool {
		return manifests[i].Timestamp.After(manifests[j].Timestamp)
	})

	return manifests, nil
}

// GetManifest reads the manifest for a specific backup
func (m *Manager) GetManifest(backupID string) (*Manifest, error) {
	manifestPath := filepath.Join(m.baseBackupDir, backupID, "manifest.json")

	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read manifest: %w", err)
	}

	var manifest Manifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, fmt.Errorf("failed to parse manifest: %w", err)
	}

	return &manifest, nil
}

// DeleteBackup removes a specific backup directory
func (m *Manager) DeleteBackup(backupID string) error {
	backupPath := filepath.Join(m.baseBackupDir, backupID)

	// Verify it's a backup directory
	if !strings.HasPrefix(backupID, "backup-") {
		return fmt.Errorf("invalid backup ID: %s", backupID)
	}

	log.Infof("Deleting backup: %s", backupID)

	if err := os.RemoveAll(backupPath); err != nil {
		return fmt.Errorf("failed to delete backup: %w", err)
	}

	return nil
}

// parseBackupTimestamp parses the timestamp from a backup directory name
// Format: backup-2026-02-18-150405
func (m *Manager) parseBackupTimestamp(backupID string) (time.Time, error) {
	// Remove "backup-" prefix
	if !strings.HasPrefix(backupID, "backup-") {
		return time.Time{}, fmt.Errorf("invalid backup ID format: %s", backupID)
	}

	timestampStr := strings.TrimPrefix(backupID, "backup-")

	// Parse timestamp: 2006-01-02-150405
	timestamp, err := time.Parse("2006-01-02-150405", timestampStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse timestamp: %w", err)
	}

	return timestamp, nil
}

// GetBackupDir returns the full path to a backup directory
func (m *Manager) GetBackupDir(backupID string) string {
	return filepath.Join(m.baseBackupDir, backupID)
}
