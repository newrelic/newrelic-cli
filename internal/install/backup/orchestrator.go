package backup

import (
	"os"
	"path/filepath"
	"runtime"

	log "github.com/sirupsen/logrus"
)

// Orchestrator coordinates all backup operations
type Orchestrator struct {
	options  Options
	detector *InstallationDetector
	creator  *Creator
	manager  *Manager
}

// NewOrchestrator creates a new backup orchestrator
func NewOrchestrator(options Options) (*Orchestrator, error) {
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

	detector := NewInstallationDetector()
	creator := NewCreator(backupLocation)
	manager := NewManager(backupLocation, options.MaxBackups)

	return &Orchestrator{
		options:  options,
		detector: detector,
		creator:  creator,
		manager:  manager,
	}, nil
}

// PerformBackup executes the backup workflow (recipe-agnostic)
// This method detects ALL existing New Relic configurations and backs them up
// before ANY recipe installation, regardless of which recipe is being installed
func (o *Orchestrator) PerformBackup(cliVersion string) (*Result, error) {
	// 1. Check if backup is skipped
	if o.options.SkipBackup {
		return nil, nil
	}

	// 2. Detect ALL existing New Relic installations (recipe-agnostic)
	// This scans for config files across ALL agent types (Infrastructure, APM, Logging, Integrations)
	installInfo, err := o.detector.DetectExistingInstallation()
	if err != nil {
		log.WithError(err).Warn("Failed to detect existing installations")
		return nil, nil // Don't fail installation on detection error
	}

	if !installInfo.IsInstalled || len(installInfo.ConfigFiles) == 0 {
		return nil, nil
	}

	// 3. Create backup of ALL detected configs
	// Works the same regardless of which recipe is being installed
	result, err := o.creator.CreateBackup(installInfo.ConfigFiles, runtime.GOOS, cliVersion)
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
	if err := o.manager.RotateBackups(); err != nil {
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

	default:
		// darwin and other platforms: use home directory
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
func (o *Orchestrator) RestoreBackup(backupID string, verifyChecksums bool) error {
	restorer := NewRestorer(o.GetBackupLocation())
	return restorer.RestoreBackup(backupID, verifyChecksums)
}

// GetBackupLocation returns the configured backup location
func (o *Orchestrator) GetBackupLocation() string {
	if o.options.BackupLocation != "" {
		return o.options.BackupLocation
	}
	return GetDefaultBackupLocation()
}
