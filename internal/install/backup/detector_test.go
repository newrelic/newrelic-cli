//go:build unit

package backup

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInstallationDetector_DetectExistingInstallation(t *testing.T) {
	// Create temp directory with mock config files
	tmpDir := t.TempDir()

	// Create mock New Relic config files
	configDir := filepath.Join(tmpDir, "etc", "newrelic-infra")
	require.NoError(t, os.MkdirAll(configDir, 0755))

	configFile := filepath.Join(configDir, "newrelic-infra.yml")
	require.NoError(t, os.WriteFile(configFile, []byte("license_key: test"), 0644))

	// Create detector with mock paths
	detector := NewInstallationDetector()

	// Override paths to use temp directory
	paths := []string{filepath.Join(tmpDir, "etc")}

	files := detector.findExistingFiles(paths)
	assert.NotEmpty(t, files)
}

func TestInstallationDetector_isNewRelicConfigFile(t *testing.T) {
	detector := &InstallationDetector{}

	tests := []struct {
		name     string
		filePath string
		want     bool
	}{
		{
			name:     "newrelic-infra.yml",
			filePath: "/etc/newrelic-infra.yml",
			want:     true,
		},
		{
			name:     "newrelic.yml",
			filePath: "/etc/newrelic/newrelic.yml",
			want:     true,
		},
		{
			name:     "file in newrelic directory",
			filePath: "/etc/newrelic-infra/config.yml",
			want:     true,
		},
		{
			name:     "non-newrelic file",
			filePath: "/etc/apache2/apache2.conf",
			want:     false,
		},
		{
			name:     "newrelic.ini",
			filePath: "/etc/php.d/newrelic.ini",
			want:     true,
		},
		{
			name:     "agent-control local_config.yaml",
			filePath: "/etc/newrelic-agent-control/local-data/agent-control/local_config.yaml",
			want:     true,
		},
		{
			name:     "agent-control environment_variables.yaml",
			filePath: "/etc/newrelic-agent-control/environment_variables.yaml",
			want:     true,
		},
		{
			name:     "ebpf agent .conf file",
			filePath: "/etc/newrelic-ebpf-agent/newrelic-ebpf-agent.conf",
			want:     true,
		},
		{
			name:     "dotnet agent .config file on Windows",
			filePath: "C:\\ProgramData\\New Relic\\.NET Agent\\newrelic.config",
			want:     true,
		},
		{
			name:     "wrong extension",
			filePath: "/etc/newrelic-infra.txt",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := detector.isNewRelicConfigFile(tt.filePath)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestInstallationDetector_DetectExistingInstallation_NoFiles(t *testing.T) {
	detector := NewInstallationDetector()

	// Use a directory that doesn't exist
	paths := []string{"/nonexistent/path"}
	files := detector.findExistingFiles(paths)
	assert.Empty(t, files)
}

func TestInstallationDetector_GetAllConfigPaths(t *testing.T) {
	detector := NewInstallationDetector()

	paths := detector.GetAllConfigPaths()

	// Should return paths for current platform
	assert.NotEmpty(t, paths)

	// Verify paths contain newrelic in some form
	foundNewRelic := false
	for _, path := range paths {
		if containsNewRelicPath(path) {
			foundNewRelic = true
			break
		}
	}
	assert.True(t, foundNewRelic, "Expected at least one path to contain 'newrelic'")
}

func containsNewRelicPath(path string) bool {
	lower := strings.ToLower(filepath.ToSlash(path))
	return strings.Contains(lower, "newrelic")
}
