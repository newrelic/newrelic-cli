package backup

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

// InstallationDetector detects existing New Relic installations
type InstallationDetector struct {
	manifest *types.DiscoveryManifest
	platform string
}

// NewInstallationDetector creates a new installation detector
func NewInstallationDetector(manifest *types.DiscoveryManifest) *InstallationDetector {
	return &InstallationDetector{
		manifest: manifest,
		platform: runtime.GOOS,
	}
}

// DetectExistingInstallation checks for existing New Relic installations
// and returns information about all detected config files across all agent types
func (d *InstallationDetector) DetectExistingInstallation(ctx context.Context) (*InstallationInfo, error) {
	log.Debug("Detecting existing New Relic installations (all agent types)")

	configPaths := d.GetAllConfigPaths()
	existingFiles, err := d.findExistingFiles(configPaths)
	if err != nil {
		return nil, err
	}

	if len(existingFiles) == 0 {
		log.Debug("No existing New Relic config files detected")
		return &InstallationInfo{
			IsInstalled: false,
			ConfigFiles: []string{},
		}, nil
	}

	log.WithFields(log.Fields{
		"platform": d.platform,
		"files":    len(existingFiles),
	}).Infof("Detected %d New Relic config files across all agent types", len(existingFiles))

	return &InstallationInfo{
		IsInstalled: true,
		ConfigFiles: existingFiles,
	}, nil
}

// GetAllConfigPaths returns all possible config paths for New Relic agents
// across Infrastructure, APM agents, Integrations, and Logging
func (d *InstallationDetector) GetAllConfigPaths() []string {
	switch d.platform {
	case "linux":
		return d.getLinuxPaths()
	case "windows":
		return d.getWindowsPaths()
	case "darwin":
		return d.getDarwinPaths()
	default:
		log.Warnf("Unsupported platform: %s", d.platform)
		return []string{}
	}
}

// getLinuxPaths returns all New Relic config paths for Linux
func (d *InstallationDetector) getLinuxPaths() []string {
	return []string{
		// Infrastructure Agent
		"/etc/newrelic-infra.yml",
		"/etc/newrelic-infra/",
		"/var/db/newrelic-infra/",

		// APM Agents (generic)
		"/etc/newrelic/newrelic.yml",
		"/etc/newrelic/",
		"/usr/local/etc/newrelic/",

		// APM Agents (language-specific)
		"/etc/newrelic-java/newrelic.yml",
		"/etc/newrelic-java/",
		"/etc/php.d/newrelic.ini",
		"/etc/php/*/conf.d/newrelic.ini",

		// User-level configs
		filepath.Join(os.Getenv("HOME"), ".newrelic/"),

		// Logging
		"/etc/newrelic-infra/logging.d/",
		"/var/log/newrelic/",
	}
}

// getWindowsPaths returns all New Relic config paths for Windows
func (d *InstallationDetector) getWindowsPaths() []string {
	programFiles := os.Getenv("ProgramFiles")
	if programFiles == "" {
		programFiles = "C:\\Program Files"
	}

	programData := os.Getenv("ProgramData")
	if programData == "" {
		programData = "C:\\ProgramData"
	}

	userProfile := os.Getenv("USERPROFILE")
	if userProfile == "" {
		userProfile = "C:\\Users\\*"
	}

	return []string{
		// Infrastructure Agent
		filepath.Join(programFiles, "New Relic", "newrelic-infra", "newrelic-infra.yml"),
		filepath.Join(programFiles, "New Relic", "newrelic-infra", "integrations.d"),
		filepath.Join(programFiles, "New Relic", "newrelic-infra"),

		// .NET Agent
		filepath.Join(programData, "New Relic", ".NET Agent"),
		filepath.Join(programFiles, "New Relic", ".NET Agent"),

		// Generic APM configs
		filepath.Join(programFiles, "New Relic", "newrelic.yml"),
		filepath.Join(programData, "New Relic"),

		// User-level configs
		filepath.Join(userProfile, "AppData", "Local", "New Relic"),
		filepath.Join(userProfile, ".newrelic"),
	}
}

// getDarwinPaths returns all New Relic config paths for macOS
func (d *InstallationDetector) getDarwinPaths() []string {
	homeDir := os.Getenv("HOME")
	if homeDir == "" {
		homeDir = "/Users/*"
	}

	return []string{
		// Infrastructure Agent
		"/usr/local/etc/newrelic-infra/newrelic-infra.yml",
		"/usr/local/etc/newrelic-infra/",
		"/opt/newrelic-infra/",

		// APM Agents
		"/usr/local/etc/newrelic/",
		"/Library/Application Support/New Relic/",

		// User-level configs
		filepath.Join(homeDir, ".newrelic/"),
		filepath.Join(homeDir, "Library", "Application Support", "New Relic"),
	}
}

// findExistingFiles recursively finds all New Relic config files
func (d *InstallationDetector) findExistingFiles(paths []string) ([]string, error) {
	var existingFiles []string
	seen := make(map[string]bool) // Deduplicate files

	// Expand wildcard paths into concrete paths before iterating
	var expandedPaths []string
	for _, path := range paths {
		if strings.Contains(path, "*") {
			matches, err := filepath.Glob(path)
			if err != nil {
				log.WithError(err).Debugf("Error globbing path: %s", path)
				continue
			}
			expandedPaths = append(expandedPaths, matches...)
		} else {
			expandedPaths = append(expandedPaths, path)
		}
	}

	for _, path := range expandedPaths {
		info, err := os.Stat(path)
		if err != nil {
			// Path doesn't exist, skip silently
			continue
		}

		if info.IsDir() {
			// Walk directory to find config files
			err = filepath.Walk(path, func(filePath string, fileInfo os.FileInfo, err error) error {
				if err != nil {
					return nil // Continue on errors
				}

				if fileInfo.IsDir() {
					return nil
				}

				// Check if this is a New Relic config file
				if d.isNewRelicConfigFile(filePath) && !seen[filePath] {
					existingFiles = append(existingFiles, filePath)
					seen[filePath] = true
				}

				return nil
			})
			if err != nil {
				log.WithError(err).Debugf("Error walking directory: %s", path)
			}
		} else {
			// Single file
			if d.isNewRelicConfigFile(path) && !seen[path] {
				existingFiles = append(existingFiles, path)
				seen[path] = true
			}
		}
	}

	log.Debugf("Found %d existing New Relic config files", len(existingFiles))
	return existingFiles, nil
}

// isNewRelicConfigFile checks if a file is a New Relic config file
func (d *InstallationDetector) isNewRelicConfigFile(filePath string) bool {
	fileName := strings.ToLower(filepath.Base(filePath))
	ext := strings.ToLower(filepath.Ext(filePath))
	dirName := strings.ToLower(filepath.Dir(filePath))

	// Check for known config file extensions
	validExtensions := []string{".yml", ".yaml", ".xml", ".ini", ".json"}
	hasValidExtension := false
	for _, validExt := range validExtensions {
		if ext == validExt {
			hasValidExtension = true
			break
		}
	}

	if !hasValidExtension {
		return false
	}

	// Check if filename or directory contains "newrelic"
	containsNewRelic := strings.Contains(fileName, "newrelic") ||
		strings.Contains(dirName, "newrelic") ||
		strings.Contains(dirName, "new relic")

	// Check for known config file patterns
	knownPatterns := []string{
		"newrelic.yml",
		"newrelic.yaml",
		"newrelic-infra.yml",
		"newrelic.ini",
		"newrelic.xml",
		"newrelic.config",
	}

	for _, pattern := range knownPatterns {
		if fileName == pattern {
			return true
		}
	}

	return containsNewRelic
}
