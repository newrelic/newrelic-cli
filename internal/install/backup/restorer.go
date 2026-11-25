package backup

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"

	log "github.com/sirupsen/logrus"
)

// Restorer handles restoring configurations from backups
type Restorer struct {
	manager *Manager
}

// NewRestorer creates a new backup restorer
func NewRestorer(baseBackupDir string) *Restorer {
	return &Restorer{
		manager: NewManager(baseBackupDir, 5),
	}
}

// RestoreBackup restores files from a backup
func (r *Restorer) RestoreBackup(backupID string, verifyChecksums bool) error {
	log.WithFields(log.Fields{
		"backupID":        backupID,
		"verifyChecksums": verifyChecksums,
	}).Info("Starting backup restoration")

	// Read manifest
	manifest, err := r.manager.GetManifest(backupID)
	if err != nil {
		return fmt.Errorf("failed to read backup manifest: %w", err)
	}

	// Verify backup integrity if requested
	if verifyChecksums {
		valid, warnings := r.verifyBackupIntegrity(manifest)
		if !valid {
			return fmt.Errorf("backup integrity check failed: %v", warnings)
		}
		if len(warnings) > 0 {
			for _, warning := range warnings {
				log.Warn(warning)
			}
		}
	}

	// Restore each file
	restored := 0
	failed := 0

	for _, file := range manifest.Files {
		if err := r.restoreFile(file); err != nil {
			log.WithError(err).Errorf("Failed to restore file: %s", file.OriginalPath)
			failed++
		} else {
			restored++
		}
	}

	if failed > 0 {
		return fmt.Errorf("restoration completed with errors: %d restored, %d failed", restored, failed)
	}

	log.WithFields(log.Fields{
		"backupID": backupID,
		"restored": restored,
	}).Info("Backup restored successfully")

	return nil
}

// verifyBackupIntegrity verifies checksums of backed up files
func (r *Restorer) verifyBackupIntegrity(manifest *Manifest) (bool, []string) {
	log.Info("Verifying backup integrity")

	var warnings []string
	allValid := true

	for _, file := range manifest.Files {
		// Check if backup file exists
		if _, err := os.Stat(file.BackupPath); os.IsNotExist(err) {
			warning := fmt.Sprintf("Backup file missing: %s", file.BackupPath)
			warnings = append(warnings, warning)
			allValid = false
			continue
		}

		// Verify checksum
		checksum, err := r.calculateChecksum(file.BackupPath)
		if err != nil {
			warning := fmt.Sprintf("Failed to calculate checksum for %s: %v", file.BackupPath, err)
			warnings = append(warnings, warning)
			allValid = false
			continue
		}

		if checksum != file.SHA256Checksum {
			warning := fmt.Sprintf("Checksum mismatch for %s", file.BackupPath)
			warnings = append(warnings, warning)
			allValid = false
		}
	}

	if allValid {
		log.Info("Backup integrity verified successfully")
	} else {
		log.Warn("Backup integrity check failed")
	}

	return allValid, warnings
}

// restoreFile restores a single file from backup
func (r *Restorer) restoreFile(file BackedUpFile) error {
	log.Debugf("Restoring file: %s", file.OriginalPath)

	// Open backup file
	srcFile, err := os.Open(file.BackupPath)
	if err != nil {
		return fmt.Errorf("failed to open backup file: %w", err)
	}
	defer func() { _ = srcFile.Close() }()

	// Get backup file info
	srcInfo, err := srcFile.Stat()
	if err != nil {
		return fmt.Errorf("failed to stat backup file: %w", err)
	}

	// Create destination directory if it doesn't exist
	dstDir := filepath.Dir(file.OriginalPath)
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Create destination file
	dstFile, err := os.OpenFile(file.OriginalPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, srcInfo.Mode())
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}

	// Copy file
	_, copyErr := io.Copy(dstFile, srcFile)
	if closeErr := dstFile.Close(); closeErr != nil && copyErr == nil {
		copyErr = fmt.Errorf("failed to close destination file: %w", closeErr)
	}
	if copyErr != nil {
		return fmt.Errorf("failed to copy file: %w", copyErr)
	}

	// Restore permissions
	if err := r.restorePermissions(file); err != nil {
		log.WithError(err).Warnf("Failed to restore permissions for %s", file.OriginalPath)
		// Don't fail restoration on permission errors
	}

	log.Debugf("File restored successfully: %s", file.OriginalPath)
	return nil
}

// restorePermissions restores file permissions
func (r *Restorer) restorePermissions(file BackedUpFile) error {
	// Parse permission string
	// Unix format: "0640" (octal)
	// Windows format: "RW" (ignored on Windows as chmod doesn't work the same way)

	if len(file.Permissions) == 0 {
		return nil
	}

	// Try to parse as octal (Unix)
	if perm, err := strconv.ParseUint(file.Permissions, 8, 32); err == nil {
		if err := os.Chmod(file.OriginalPath, os.FileMode(perm)); err != nil {
			return fmt.Errorf("failed to set permissions: %w", err)
		}
	}
	// Windows permissions are descriptive only, chmod may not work as expected

	return nil
}

// calculateChecksum calculates SHA256 checksum of a file
func (r *Restorer) calculateChecksum(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer func() { _ = file.Close() }()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}
