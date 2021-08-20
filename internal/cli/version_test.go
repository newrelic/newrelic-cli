// +build unit

package cli

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVersion(t *testing.T) {
	t.Parallel()

	version = "1.0.0"
	require.Equal(t, "1.0.0", Version())
	require.NotEqual(t, "0.0.1", Version())
}

func TestIsLatestVersion(t *testing.T) {
	t.Parallel()

	version = "1.0.0"
	latest := "1.0.0"

	result, err := IsLatestVersion(context.Background(), latest)
	require.NoError(t, err)
	require.True(t, result)
}

func TestIsLatestVersion_False(t *testing.T) {
	t.Parallel()

	// Set installed version as an older version
	version = "0.30.0"
	latest := "0.31.0"

	result, err := IsLatestVersion(context.Background(), latest)
	require.NoError(t, err)
	require.False(t, result)
}

func TestIsDevEnvironment(t *testing.T) {
	t.Parallel()

	version = "0.32.1-10-gbe63a24-dirty"
	result := IsDevEnvironment()
	require.True(t, result)
}

func TestIsDevEnvironment_False(t *testing.T) {
	t.Parallel()

	version = "0.32.1"
	result := IsDevEnvironment()
	require.False(t, result)
}

func TestGetLatestReleaseVersion_Cached(t *testing.T) {
	latestVersion = "0.31.0"

	result, err := GetLatestReleaseVersion(context.Background())
	require.NoError(t, err)
	require.Equal(t, "0.31.0", result)
}
