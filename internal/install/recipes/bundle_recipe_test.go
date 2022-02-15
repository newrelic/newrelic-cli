package recipes

import (
	"testing"

	"github.com/newrelic/newrelic-cli/internal/install/execution"
	"github.com/stretchr/testify/require"
)

func TestBundleRecipe_ShouldAddStatus(t *testing.T) {
	br := givenBundleRecipe()

	br.AddStatus(execution.RecipeStatusTypes.INSTALLING)

	require.Equal(t, len(br.Statuses), 1)
	require.Equal(t, br.Statuses[0], execution.RecipeStatusTypes.INSTALLING)
}

func TestBundleRecipe_ShouldAddStatusOnce(t *testing.T) {
	br := givenBundleRecipe()

	br.AddStatus(execution.RecipeStatusTypes.INSTALLING)
	br.AddStatus(execution.RecipeStatusTypes.INSTALLING)
	br.AddStatus(execution.RecipeStatusTypes.INSTALLING)

	require.Equal(t, len(br.Statuses), 1)
	require.Equal(t, br.Statuses[0], execution.RecipeStatusTypes.INSTALLING)
}

func TestBundleRecipe_ShouldAddStatusDetectedWhenAvailable(t *testing.T) {
	br := givenBundleRecipe()

	br.AddStatus(execution.RecipeStatusTypes.AVAILABLE)

	require.Equal(t, len(br.Statuses), 2)
	require.Equal(t, br.Statuses[0], execution.RecipeStatusTypes.DETECTED)
	require.Equal(t, br.Statuses[1], execution.RecipeStatusTypes.AVAILABLE)
}

func givenBundleRecipe() *BundleRecipe {
	return &BundleRecipe{
		Recipe: createRecipe("id1", "recipe2"),
	}
}
