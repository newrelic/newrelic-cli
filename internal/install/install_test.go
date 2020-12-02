// +build unit

package install

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInstall(t *testing.T) {
	assert.True(t, true)
}

func TestMatchLogFilesFromRecipe(t *testing.T) {
	// Create a temp directory to work with
	tmpDir, err := ioutil.TempDir("/tmp", "logfiles")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Create a few "log" files within the directory
	f1, err := ioutil.TempFile(tmpDir, "*.log")
	defer f1.Close()
	require.NoError(t, err)
	f2, err := ioutil.TempFile(tmpDir, "*.log")
	require.NoError(t, err)
	defer f2.Close()
	f3, err := ioutil.TempFile(tmpDir, "*.nopelog")
	require.NoError(t, err)
	defer f3.Close()

	r := recipeFile{}
	r.MELTMatch.Logging = []logMatcher{
		{
			Name: "nginx",
			File: filepath.Join(tmpDir, "*.log"),
		},
	}

	for _, x := range r.MELTMatch.Logging {
		match, files := matchLogFilesFromRecipe(x)

		if x.Name == "nginx" {
			t.Logf("files: %+v", files)

			expectedMatch := []string{
				f1.Name(),
				f2.Name(),
			}

			sort.Strings(files)
			sort.Strings(expectedMatch)

			assert.True(t, match)
			assert.Equal(t, files, expectedMatch)
		}

	}

}
