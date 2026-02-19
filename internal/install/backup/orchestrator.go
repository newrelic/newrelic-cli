package backup

import (
	"context"
	"os"
	"path/filepath"
	"runtime"

	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

// Orchestrator coordinates all backup operations
type Orchestrator struct {
	manifest *types.DiscoveryManifest
	options  Options
	detector *InstallationDetector
	creator  *Creator
	manager  *Manager
}

// NewOrchestrator creates a new backup orchestrator
func NewOrchestrator(manifest *types.DiscoveryManifest, options Options, cliVersion string) (*Orchestrator, error) {
	// Set defaults
	if options.MaxBackups == 0 {
		options.MaxBackups = 5
	}

	// Determine backup location
	backupLocation := options.BackupLocation
	if backupLocation == "" {
		backupLocation = GetDefaultBackupLocation()
	}

	log.WithFields(log.Fields{
		"backupLocation": backupLocation,
		"skipBackup":     options.SkipBackup,
		"maxBackups":     options.MaxBackups,
	}).Debug("Initializing backup orchestrator")

	detector := NewInstallationDetector(manifest)
	creator := NewCreator(backupLocation)
	manager := NewManager(backupLocation, options.MaxBackups)

	return &Orchestrator{
		manifest: manifest,
		options:  options,
		detector: detector,
		creator:  creator,
		manager:  manager,
	}, nil
}

// PerformBackup executes the backup workflow (recipe-agnostic)
// This method detects ALL existing New Relic configurations and backs them up
// before ANY recipe installation, regardless of which recipe is being installed
func (o *Orchestrator) PerformBackup(ctx context.Context, cliVersion string) (*Result, error) {
	// 1. Check if backup is skipped
	if o.options.SkipBackup {
		log.Info("Backup skipped by user option")
		return nil, nil
	}

	// 2. Detect ALL existing New Relic installations (recipe-agnostic)
	// This scans for config files across ALL agent types (Infrastructure, APM, Logging, Integrations)
	installInfo, err := o.detector.DetectExistingInstallation(ctx)
	if err != nil {
		log.WithError(err).Warn("Failed to detect existing installations")
		return nil, nil // Don't fail installation on detection error
	}

	if !installInfo.IsInstalled || len(installInfo.ConfigFiles) == 0 {
		log.Info("No existing New Relic installations detected, skipping backup")
		return nil, nil
	}

	log.WithFields(log.Fields{
		"configFiles": len(installInfo.ConfigFiles),
		"platform":    runtime.GOOS,
	}).Infof("Detected %d New Relic config files across all agent types", len(installInfo.ConfigFiles))

	// 3. Create backup of ALL detected configs
	// Works the same regardless of which recipe is being installed
	result, err := o.creator.CreateBackup(ctx, installInfo.ConfigFiles, runtime.GOOS, cliVersion)
	if err != nil {
		log.WithError(err).Warn("Backup failed. Installation will continue.")
		return result, nil // Don't fail installation on backup error
	}

	if !result.Success {
		log.Warn("Backup completed with errors. Installation will continue.")
		if len(result.Warnings) > 0 {
			for _, warning := range result.Warnings {
				log.Warn(warning)
			}
		}
		return result, nil
	}

	// 4. Rotate old backups
	if err := o.manager.RotateBackups(ctx); err != nil {
		log.WithError(err).Warn("Failed to rotate old backups")
		// Don't fail on rotation error
	}

	return result, nil
}

// GetDefaultBackupLocation returns the platform-specific default backup location
func GetDefaultBackupLocation() string {
	var baseDir string

	switch runtime.GOOS {
	case "linux":
		// Check if running as root
		if os.Geteuid() == 0 {
			baseDir = "/opt/.newrelic-backups"
		} else {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				baseDir = "/tmp/.newrelic-backups"
			} else {
				baseDir = filepath.Join(homeDir, ".newrelic-backups")
			}
		}

	case "windows":
		programData := os.Getenv("ProgramData")
		if programData == "" {
			programData = "C:\\ProgramData"
		}
		baseDir = filepath.Join(programData, ".newrelic-backups")

	case "darwin":
		homeDir, err := os.UserHomeDir()
		if err != nil {
			baseDir = "/tmp/.newrelic-backups"
		} else {
			baseDir = filepath.Join(homeDir, ".newrelic-backups")
		}

	default:
		// Fallback for unknown platforms
		homeDir, err := os.UserHomeDir()
		if err != nil {
			baseDir = "/tmp/.newrelic-backups"
		} else {
			baseDir = filepath.Join(homeDir, ".newrelic-backups")
		}
	}

	return baseDir
}

// ListBackups returns all available backups
func (o *Orchestrator) ListBackups() ([]*Manifest, error) {
	return o.manager.ListBackups()
}

// RestoreBackup restores a specific backup
func (o *Orchestrator) RestoreBackup(ctx context.Context, backupID string, verifyChecksums bool) error {
	backupLocation := o.options.BackupLocation
	if backupLocation == "" {
		backupLocation = GetDefaultBackupLocation()
	}

	restorer := NewRestorer(backupLocation)
	return restorer.RestoreBackup(ctx, backupID, verifyChecksums)
}

// GetBackupLocation returns the configured backup location
func (o *Orchestrator) GetBackupLocation() string {
	if o.options.BackupLocation != "" {
		return o.options.BackupLocation
	}
	return GetDefaultBackupLocation()
}
