//go:build integration
// +build integration

package discovery

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

func TestLogMatchFinder(t *testing.T) {
	// Create a temp directory to work with
	tmpDir, err := ioutil.TempDir("/tmp", "logfiles")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Create a few "log" files within the directory
	f1, err := ioutil.TempFile(tmpDir, "*.log")
	require.NoError(t, err)
	defer f1.Close()

	f2, err := ioutil.TempFile(tmpDir, "*.log")
	require.NoError(t, err)
	defer f2.Close()

	f3, err := ioutil.TempFile(tmpDir, "*.nopelog")
	require.NoError(t, err)
	defer f3.Close()

	recipes := []types.OpenInstallationRecipe{
		{
			ID: "test",
			LogMatch: []types.OpenInstallationLogMatch{
				{
					File: filepath.Join(tmpDir, "*.log"),
				},
			},
		},
		{
			ID: "nginx",
			LogMatch: []types.OpenInstallationLogMatch{
				{
					File: "/var/log/nope/*.log",
				},
			},
		},
	}

	f := NewLogMatchFinderer()
	filtered, err := f.Filter(context.Background(), recipes)

	require.NoError(t, err)
	require.NotNil(t, filtered)
	require.NotEmpty(t, filtered)
	require.Equal(t, 1, len(filtered))
}

func TestLogMatchFinderFromRecipe(t *testing.T) {

	// Create a temp directory to work with
	tmpDir, err := ioutil.TempDir("/tmp", "logfiles")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Create a few "log" files within the directory
	f1, err := ioutil.TempFile(tmpDir, "*.log")
	require.NoError(t, err)
	defer f1.Close()

	f2, err := ioutil.TempFile(tmpDir, "*.log")
	require.NoError(t, err)
	defer f2.Close()

	f3, err := ioutil.TempFile(tmpDir, "*.nopelog")
	require.NoError(t, err)
	defer f3.Close()

	recipes := []types.OpenInstallationRecipe{
		{
			ID: "test",
			LogMatch: []types.OpenInstallationLogMatch{
				{
					File: filepath.Join(tmpDir, "*.log"),
				},
			},
		},
	}

	matched, files := matchLogFilesFromRecipe(recipes[0].LogMatch[0])
	require.True(t, matched)
	require.Equal(t, 2, len(files))
}
