// +build integration

package install

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-cli/internal/install/recipes"
)

func TestFetchRecipeFile_FailedStatusCode(t *testing.T) {
	ff := recipes.NewRecipeFileFetcher()
	u, err := url.Parse("https://httpstat.us/404")
	require.NoError(t, err)

	f, err := ff.FetchRecipeFile(u)
	require.Error(t, err)
	require.Nil(t, f)
}
