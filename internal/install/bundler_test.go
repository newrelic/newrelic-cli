//go:build unit
// +build unit

package install

import (
	"testing"

	"github.com/newrelic/newrelic-cli/internal/install/types"

	"github.com/stretchr/testify/require"
)

func TestCreateBundles_ShouldCreateTwoBundles(t *testing.T) {
	recipes := []types.OpenInstallationRecipe{
		{
			Name: types.InfraAgentRecipeName,
		},
		{
			Name: types.LoggingRecipeName,
		},
		{
			Name: types.GoldenRecipeName,
		},
		{
			Name: "mysql-open-source-integration",
		},
	}
	coreBundle, extrasBundle := createBundles(coreBundleRecipeNames, recipes)

	require.Equal(t, 3, len(coreBundle))
	require.Greater(t, 1, len(extrasBundle))
}

func TestCreateBundles_ShouldCreateEmptyExtrasBundle(t *testing.T) {
	recipes := []types.OpenInstallationRecipe{
		{
			Name: types.InfraAgentRecipeName,
		},
		{
			Name: types.LoggingRecipeName,
		},
		{
			Name: types.GoldenRecipeName,
		},
	}
	coreBundle, extrasBundle := createBundles(coreBundleRecipeNames, recipes)

	require.Equal(t, 3, len(coreBundle))
	require.Equal(t, 0, len(extrasBundle))
}
