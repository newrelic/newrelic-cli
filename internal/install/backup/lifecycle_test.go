//go:build unit

package backup

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestManager_RotateBackups(t *testing.T) {
	backupBaseDir := t.TempDir()

	// Create 7 mock backups
	backupIDs := []string{
		"backup-2026-02-10-100000",
		"backup-2026-02-11-100000",
		"backup-2026-02-12-100000",
		"backup-2026-02-13-100000",
		"backup-2026-02-14-100000",
		"backup-2026-02-15-100000",
		"backup-2026-02-16-100000",
	}

	for _, id := range backupIDs {
		backupDir := filepath.Join(backupBaseDir, id)
		require.NoError(t, os.MkdirAll(backupDir, 0755))

		// Create a manifest file
		manifest := Manifest{
			BackupID:   id,
			Platform:   "linux",
			Timestamp:  time.Now(),
			CLIVersion: "v0.106.23",
			Reason:     "test",
		}
		manifestData, _ := json.Marshal(manifest)
		manifestPath := filepath.Join(backupDir, "manifest.json")
		require.NoError(t, os.WriteFile(manifestPath, manifestData, 0644))
	}

	// Keep only last 5
	manager := NewManager(backupBaseDir, 5)
	err := manager.RotateBackups(context.Background())
	require.NoError(t, err)

	// Verify only 5 backups remain
	entries, err := os.ReadDir(backupBaseDir)
	require.NoError(t, err)

	count := 0
	for _, entry := range entries {
		if entry.IsDir() {
			count++
		}
	}
	assert.Equal(t, 5, count)

	// Verify the oldest 2 backups were deleted
	_, err = os.Stat(filepath.Join(backupBaseDir, backupIDs[0]))
	assert.True(t, os.IsNotExist(err))

	_, err = os.Stat(filepath.Join(backupBaseDir, backupIDs[1]))
	assert.True(t, os.IsNotExist(err))

	// Verify the newest 5 still exist
	for i := 2; i < 7; i++ {
		_, err = os.Stat(filepath.Join(backupBaseDir, backupIDs[i]))
		assert.NoError(t, err)
	}
}

func TestManager_ListBackups(t *testing.T) {
	backupBaseDir := t.TempDir()

	// Create 3 mock backups with different timestamps
	backupData := []struct {
		id        string
		timestamp time.Time
	}{
		{"backup-2026-02-16-100000", time.Date(2026, 2, 16, 10, 0, 0, 0, time.UTC)},
		{"backup-2026-02-15-100000", time.Date(2026, 2, 15, 10, 0, 0, 0, time.UTC)},
		{"backup-2026-02-14-100000", time.Date(2026, 2, 14, 10, 0, 0, 0, time.UTC)},
	}

	for _, data := range backupData {
		backupDir := filepath.Join(backupBaseDir, data.id)
		require.NoError(t, os.MkdirAll(backupDir, 0755))

		manifest := Manifest{
			BackupID:   data.id,
			Platform:   "linux",
			Timestamp:  data.timestamp,
			CLIVersion: "v0.106.23",
			Reason:     "guided-install",
			Files:      []BackedUpFile{},
		}
		manifestData, _ := json.Marshal(manifest)
		require.NoError(t, os.WriteFile(filepath.Join(backupDir, "manifest.json"), manifestData, 0644))
	}

	manager := NewManager(backupBaseDir, 5)
	backups, err := manager.ListBackups()

	require.NoError(t, err)
	assert.Equal(t, 3, len(backups))

	// Verify sorted by timestamp (newest first)
	assert.Equal(t, "backup-2026-02-16-100000", backups[0].BackupID)
	assert.Equal(t, "backup-2026-02-15-100000", backups[1].BackupID)
	assert.Equal(t, "backup-2026-02-14-100000", backups[2].BackupID)
}

func TestManager_GetManifest(t *testing.T) {
	backupBaseDir := t.TempDir()

	backupID := "backup-2026-02-18-143022"
	backupDir := filepath.Join(backupBaseDir, backupID)
	require.NoError(t, os.MkdirAll(backupDir, 0755))

	expectedManifest := Manifest{
		BackupID:   backupID,
		Platform:   "linux",
		Timestamp:  time.Now(),
		CLIVersion: "v0.106.23",
		Reason:     "guided-install",
		Files: []BackedUpFile{
			{
				OriginalPath:   "/etc/newrelic-infra.yml",
				BackupPath:     "/backup/etc/newrelic-infra.yml",
				SHA256Checksum: "abc123",
				Size:           100,
				Permissions:    "0644",
			},
		},
	}

	manifestData, _ := json.Marshal(expectedManifest)
	require.NoError(t, os.WriteFile(filepath.Join(backupDir, "manifest.json"), manifestData, 0644))

	manager := NewManager(backupBaseDir, 5)
	manifest, err := manager.GetManifest(backupID)

	require.NoError(t, err)
	assert.Equal(t, backupID, manifest.BackupID)
	assert.Equal(t, "linux", manifest.Platform)
	assert.Equal(t, 1, len(manifest.Files))
	assert.Equal(t, "/etc/newrelic-infra.yml", manifest.Files[0].OriginalPath)
}

func TestManager_parseBackupTimestamp(t *testing.T) {
	manager := &Manager{}

	tests := []struct {
		name      string
		backupID  string
		wantError bool
	}{
		{
			name:      "valid timestamp",
			backupID:  "backup-2026-02-18-143022",
			wantError: false,
		},
		{
			name:      "invalid format",
			backupID:  "invalid-format",
			wantError: true,
		},
		{
			name:      "missing prefix",
			backupID:  "2026-02-18-143022",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			timestamp, err := manager.parseBackupTimestamp(tt.backupID)

			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotZero(t, timestamp)
			}
		})
	}
}
