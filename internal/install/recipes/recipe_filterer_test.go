package recipes

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-cli/internal/install/execution"
	"github.com/newrelic/newrelic-cli/internal/install/types"
)

func TestRecommend_CustomScript_Success(t *testing.T) {
	recipe := types.OpenInstallationRecipe{
		Name: "test-recipe",
		PreInstall: types.OpenInstallationPreInstallConfiguration{
			RequireAtDiscovery: "echo 1234",
		},
	}

	m := &types.DiscoveryManifest{}

	r := NewRecipeFilterer(types.InstallerContext{}, &execution.InstallStatus{})

	filtered, err := r.Filter(context.Background(), &recipe, m)
	require.NoError(t, err)
	require.False(t, filtered)
}

func TestRecommend_CustomScript_Failure(t *testing.T) {
	recipe := types.OpenInstallationRecipe{
		Name: "test-recipe",
		PreInstall: types.OpenInstallationPreInstallConfiguration{
			RequireAtDiscovery: "bogus command",
		},
	}

	m := &types.DiscoveryManifest{}

	r := NewRecipeFilterer(types.InstallerContext{}, &execution.InstallStatus{})

	filtered, err := r.Filter(context.Background(), &recipe, m)
	require.NoError(t, err)
	require.True(t, filtered)
}
