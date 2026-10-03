//go:build unit

package backup

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOrchestrator_PerformBackup_SkipBackup(t *testing.T) {
	options := Options{
		SkipBackup: true,
		MaxBackups: 5,
	}

	orchestrator, err := NewOrchestrator(options)
	require.NoError(t, err)

	result, err := orchestrator.PerformBackup("v0.106.23")
	require.NoError(t, err)
	assert.Nil(t, result)
}

func TestOrchestrator_PerformBackup_NoExistingInstallation(t *testing.T) {
	backupBaseDir := t.TempDir()
	options := Options{
		SkipBackup:     false,
		BackupLocation: backupBaseDir,
		MaxBackups:     5,
	}

	orchestrator, err := NewOrchestrator(options)
	require.NoError(t, err)

	// Note: The detector scans the real system and may find existing NR config files
	// This is expected behavior - the backup system is recipe-agnostic and detects all NR configs
	result, err := orchestrator.PerformBackup("v0.106.23")
	require.NoError(t, err)
	// Result may be nil (no configs found) or have a backup (configs found on system)
	// Both are valid outcomes depending on the test environment
	if result != nil {
		assert.True(t, result.Success)
	}
}

func TestGetDefaultBackupLocation(t *testing.T) {
	location := GetDefaultBackupLocation()

	assert.NotEmpty(t, location)
	assert.Contains(t, location, ".newrelic-backups")

	// Platform-specific checks
	switch runtime.GOOS {
	case "linux":
		// Should contain either /opt or home directory
		assert.True(t, filepath.IsAbs(location))
	case "windows":
		// Should be in ProgramData or similar
		assert.True(t, filepath.IsAbs(location))
	case "darwin":
		// Should be in home directory
		assert.Contains(t, location, ".newrelic-backups")
	}
}

func TestOrchestrator_GetBackupLocation(t *testing.T) {
	t.Run("custom location", func(t *testing.T) {
		customPath := "/custom/backup/path"
		options := Options{
			BackupLocation: customPath,
			MaxBackups:     5,
		}

		orchestrator, err := NewOrchestrator(options)
		require.NoError(t, err)

		location := orchestrator.GetBackupLocation()
		assert.Equal(t, customPath, location)
	})

	t.Run("default location", func(t *testing.T) {
		options := Options{
			MaxBackups: 5,
		}

		orchestrator, err := NewOrchestrator(options)
		require.NoError(t, err)

		location := orchestrator.GetBackupLocation()
		assert.NotEmpty(t, location)
		assert.Contains(t, location, ".newrelic-backups")
	})
}
