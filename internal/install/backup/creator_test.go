//go:build unit

package backup

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreator_CreateBackup(t *testing.T) {
	// Create temp directories
	srcDir := t.TempDir()
	backupBaseDir := t.TempDir()

	// Create mock config files
	configFile1 := filepath.Join(srcDir, "newrelic-infra.yml")
	configFile2 := filepath.Join(srcDir, "integrations.d", "mysql.yml")

	require.NoError(t, os.WriteFile(configFile1, []byte("license_key: test1"), 0644))
	require.NoError(t, os.MkdirAll(filepath.Dir(configFile2), 0755))
	require.NoError(t, os.WriteFile(configFile2, []byte("interval: 30s"), 0644))

	// Create backup
	creator := NewCreator(backupBaseDir)
	result, err := creator.CreateBackup([]string{configFile1, configFile2}, "linux", "v0.106.23")

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.True(t, result.Success)
	assert.Equal(t, 2, result.FilesBackedUp)
	assert.NotEmpty(t, result.BackupID)
	assert.Contains(t, result.BackupID, "backup-")
	assert.NotEmpty(t, result.BackupDir)
	assert.NotEmpty(t, result.ManifestPath)

	// Verify backup directory was created
	_, err = os.Stat(result.BackupDir)
	assert.NoError(t, err)

	// Verify manifest exists
	manifestData, err := os.ReadFile(result.ManifestPath)
	require.NoError(t, err)

	var manifest Manifest
	require.NoError(t, json.Unmarshal(manifestData, &manifest))

	assert.Equal(t, 2, len(manifest.Files))
	assert.Equal(t, "linux", manifest.Platform)
	assert.Equal(t, "guided-install", manifest.Reason)
	assert.Equal(t, "v0.106.23", manifest.CLIVersion)

	// Verify checksums are not empty
	for _, file := range manifest.Files {
		assert.NotEmpty(t, file.SHA256Checksum)
		assert.Greater(t, file.Size, int64(0))
		assert.NotEmpty(t, file.Permissions)
	}

	// Verify the backed-up file exists at the correct nested path inside backupDir.
	vol := filepath.VolumeName(configFile1)
	afterVol := configFile1[len(vol):]
	if len(afterVol) > 0 && afterVol[0] == filepath.Separator {
		afterVol = afterVol[1:]
	}
	expectedPath := filepath.Join(result.BackupDir, afterVol)
	_, statErr := os.Stat(expectedPath)
	assert.NoError(t, statErr, "backed up file should exist at correct nested path inside backupDir")
}

func TestCreator_CreateBackup_NoFiles(t *testing.T) {
	backupBaseDir := t.TempDir()

	creator := NewCreator(backupBaseDir)
	result, err := creator.CreateBackup([]string{}, "linux", "v0.106.23")

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.True(t, result.Success)
	assert.Equal(t, 0, result.FilesBackedUp)
	assert.Len(t, result.Warnings, 1)
}

func TestCreator_CreateBackup_PartialFailure(t *testing.T) {
	srcDir := t.TempDir()
	backupBaseDir := t.TempDir()

	// Create one valid file
	validFile := filepath.Join(srcDir, "valid.yml")
	require.NoError(t, os.WriteFile(validFile, []byte("test"), 0644))

	// Reference one nonexistent file
	invalidFile := filepath.Join(srcDir, "nonexistent.yml")

	creator := NewCreator(backupBaseDir)
	result, err := creator.CreateBackup([]string{validFile, invalidFile}, "linux", "v0.106.23")

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.True(t, result.Success)
	assert.Equal(t, 1, result.FilesBackedUp)
	assert.NotEmpty(t, result.Warnings)
}

func TestCreator_generateBackupID(t *testing.T) {
	creator := &Creator{}

	backupID := creator.generateBackupID()

	assert.NotEmpty(t, backupID)
	assert.Contains(t, backupID, "backup-")
}

func TestCreator_copyFileWithChecksum(t *testing.T) {
	srcDir := t.TempDir()
	dstDir := t.TempDir()

	// Create source file
	srcFile := filepath.Join(srcDir, "test.yml")
	content := []byte("test content for checksum")
	require.NoError(t, os.WriteFile(srcFile, content, 0644))

	// Open files for copying
	src, err := os.Open(srcFile)
	require.NoError(t, err)
	defer src.Close()

	dstFile := filepath.Join(dstDir, "test.yml")
	dst, err := os.Create(dstFile)
	require.NoError(t, err)
	defer dst.Close()

	creator := &Creator{}
	checksum, size, err := creator.copyFileWithChecksum(src, dst)

	require.NoError(t, err)
	assert.NotEmpty(t, checksum)
	assert.Equal(t, int64(len(content)), size)
	assert.Len(t, checksum, 64) // SHA256 hex is 64 chars
}
