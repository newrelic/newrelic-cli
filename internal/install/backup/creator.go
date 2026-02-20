package backup

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"time"

	log "github.com/sirupsen/logrus"
)

// Creator handles creating timestamped backups with checksums
type Creator struct {
	baseBackupDir string
}

// NewCreator creates a new backup creator
func NewCreator(baseBackupDir string) *Creator {
	return &Creator{
		baseBackupDir: baseBackupDir,
	}
}

// CreateBackup creates a timestamped backup of config files with checksums
func (c *Creator) CreateBackup(files []string, platform string, cliVersion string) (*Result, error) {
	if len(files) == 0 {
		log.Debug("No files to backup")
		return &Result{
			Success:  true,
			Warnings: []string{"No config files found to backup"},
		}, nil
	}

	backupID := c.generateBackupID()
	backupDir := filepath.Join(c.baseBackupDir, backupID)

	log.WithFields(log.Fields{
		"backupID":  backupID,
		"backupDir": backupDir,
		"files":     len(files),
	}).Debug("Creating configuration backup")

	// Create backup directory
	if err := os.MkdirAll(backupDir, 0750); err != nil {
		return &Result{
			Success: false,
			Error:   fmt.Errorf("failed to create backup directory: %w", err),
		}, err
	}

	result := &Result{
		BackupID:  backupID,
		BackupDir: backupDir,
		Success:   true,
		Warnings:  []string{},
	}

	manifest := &Manifest{
		Timestamp:  time.Now(),
		BackupID:   backupID,
		Platform:   platform,
		Reason:     "guided-install",
		CLIVersion: cliVersion,
		Files:      []BackedUpFile{},
	}

	// Backup each file
	for _, srcPath := range files {
		backedUpFile, err := c.backupSingleFile(srcPath, backupDir)
		if err != nil {
			warning := fmt.Sprintf("Failed to backup %s: %v", srcPath, err)
			result.Warnings = append(result.Warnings, warning)
			log.Warn(warning)
			continue
		}

		manifest.Files = append(manifest.Files, *backedUpFile)
		result.FilesBackedUp++
	}

	// Write manifest
	manifestPath, err := c.writeManifest(manifest, backupDir)
	if err != nil {
		result.Warnings = append(result.Warnings, fmt.Sprintf("Failed to write manifest: %v", err))
		log.WithError(err).Warn("Failed to write backup manifest")
	} else {
		result.ManifestPath = manifestPath
	}

	if result.FilesBackedUp == 0 {
		result.Success = false
		result.Error = fmt.Errorf("no files were successfully backed up")
		return result, result.Error
	}

	log.WithFields(log.Fields{
		"backupID":      backupID,
		"filesBackedUp": result.FilesBackedUp,
		"warnings":      len(result.Warnings),
	}).Debug("Backup completed")

	return result, nil
}

// backupSingleFile backs up a single file with checksum
func (c *Creator) backupSingleFile(srcPath string, backupDir string) (*BackedUpFile, error) {
	// Open source file
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open source file: %w", err)
	}
	defer func() { _ = srcFile.Close() }()

	// Get file info
	srcInfo, err := srcFile.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to stat source file: %w", err)
	}

	// Create destination path preserving relative structure.
	// Strip the volume name (e.g. "C:" on Windows, "" on Unix) and the
	// leading path separator so the result is always a relative path,
	// otherwise filepath.Join discards the backupDir prefix entirely.
	// e.g. Linux:   /etc/newrelic-infra.yml               → etc/newrelic-infra.yml
	//      Windows: C:\Program Files\New Relic\newrelic-infra.yml → Program Files\New Relic\newrelic-infra.yml
	relPath := srcPath
	if filepath.IsAbs(srcPath) {
		vol := filepath.VolumeName(srcPath) // "C:" on Windows, "" on Unix
		afterVol := srcPath[len(vol):]
		if len(afterVol) > 0 && afterVol[0] == filepath.Separator {
			relPath = afterVol[1:]
		} else {
			relPath = afterVol
		}
	}
	dstPath := filepath.Join(backupDir, relPath)

	// Create destination directory
	dstDir := filepath.Dir(dstPath)
	if err := os.MkdirAll(dstDir, 0750); err != nil {
		return nil, fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Create destination file
	dstFile, err := os.OpenFile(dstPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, srcInfo.Mode())
	if err != nil {
		return nil, fmt.Errorf("failed to create destination file: %w", err)
	}

	// Copy file and calculate checksum
	checksum, size, copyErr := c.copyFileWithChecksum(srcFile, dstFile)
	if closeErr := dstFile.Close(); closeErr != nil && copyErr == nil {
		copyErr = fmt.Errorf("failed to close destination file: %w", closeErr)
	}
	if copyErr != nil {
		return nil, fmt.Errorf("failed to copy file: %w", copyErr)
	}

	permissions := c.getFilePermissions(srcInfo)

	return &BackedUpFile{
		OriginalPath:   srcPath,
		BackupPath:     dstPath,
		SHA256Checksum: checksum,
		Size:           size,
		Permissions:    permissions,
	}, nil
}

// copyFileWithChecksum copies a file and calculates SHA256 checksum
func (c *Creator) copyFileWithChecksum(src io.Reader, dst io.Writer) (checksum string, size int64, err error) {
	hash := sha256.New()
	writer := io.MultiWriter(dst, hash)

	size, err = io.Copy(writer, src)
	if err != nil {
		return "", 0, err
	}

	checksum = hex.EncodeToString(hash.Sum(nil))
	return checksum, size, nil
}

// writeManifest writes the backup manifest to JSON
func (c *Creator) writeManifest(manifest *Manifest, backupDir string) (string, error) {
	manifestPath := filepath.Join(backupDir, "manifest.json")

	jsonData, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal manifest: %w", err)
	}

	if err := os.WriteFile(manifestPath, jsonData, 0640); err != nil {
		return "", fmt.Errorf("failed to write manifest file: %w", err)
	}

	log.Debugf("Manifest written to: %s", manifestPath)
	return manifestPath, nil
}

// getFilePermissions returns a string representation of file permissions
func (c *Creator) getFilePermissions(info os.FileInfo) string {
	if runtime.GOOS == "windows" {
		// Windows: simplified permission string
		mode := info.Mode()
		perms := ""
		if mode&0400 != 0 {
			perms += "R"
		}
		if mode&0200 != 0 {
			perms += "W"
		}
		if perms == "" {
			perms = "RO"
		}
		return perms
	}

	// Unix: octal permission string
	return fmt.Sprintf("%04o", info.Mode().Perm())
}

// generateBackupID generates a timestamp-based backup ID
func (c *Creator) generateBackupID() string {
	return fmt.Sprintf("backup-%s", time.Now().Format("2006-01-02-150405"))
}
