package backup

import "time"

// Manifest contains metadata about a backup
type Manifest struct {
	Timestamp  time.Time      `json:"timestamp"`
	BackupID   string         `json:"backupId"` // backup-2026-02-18-143022
	Platform   string         `json:"platform"` // linux, windows, darwin
	Files      []BackedUpFile `json:"files"`
	Reason     string         `json:"reason"` // "guided-install"
	CLIVersion string         `json:"cliVersion"`
}

// BackedUpFile represents a single backed up file
type BackedUpFile struct {
	OriginalPath   string `json:"originalPath"`
	BackupPath     string `json:"backupPath"`
	SHA256Checksum string `json:"sha256Checksum"`
	Size           int64  `json:"size"`
	Permissions    string `json:"permissions"` // Unix: "0640", Windows: "RW"
}

// Options controls backup behavior
type Options struct {
	SkipBackup     bool
	BackupLocation string
	MaxBackups     int
}

// Result contains the outcome of a backup operation
type Result struct {
	BackupID      string
	BackupDir     string
	ManifestPath  string
	FilesBackedUp int
	Success       bool
	Error         error
	Warnings      []string
}

// InstallationInfo contains detected installation details
type InstallationInfo struct {
	IsInstalled bool
	ConfigFiles []string
}
