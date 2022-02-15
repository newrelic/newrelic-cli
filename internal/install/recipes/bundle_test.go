package recipes

import (
	"testing"

	"github.com/newrelic/newrelic-cli/internal/install/execution"
	"github.com/stretchr/testify/require"
)

var (
	bundle_recipe = &BundleRecipe{
		Recipe: createRecipe("0", "recipe1"),
	}
	bundle_recipes = []*BundleRecipe{
		&BundleRecipe{
			Recipe: createRecipe("0", "recipe1"),
		},
		&BundleRecipe{
			Recipe: createRecipe("1", "recipe2"),
		},
	}
)

func TestBundle_ShouldAddRecipes(t *testing.T) {

	bundle := Bundle{
		BundleRecipes: bundle_recipes,
	}
	newRecipe := &BundleRecipe{
		Recipe: createRecipe("2", "recipe3"),
	}

	bundle.AddRecipe(newRecipe)
	require.Equal(t, len(bundle.BundleRecipes), len(bundle_recipes)+1)
	require.True(t, true, bundle.ContainsName(newRecipe.Recipe.Name))
}

func TestBundle_ShouldUpdateRecipe(t *testing.T) {

	bundle := Bundle{
		BundleRecipes: bundle_recipes,
	}

	bundle.AddRecipe(bundle_recipe)

	require.Equal(t, len(bundle.BundleRecipes), len(bundle_recipes))
	require.Equal(t, true, bundle.ContainsName(bundle_recipe.Recipe.Name))
}

func TestBundle_ShouldContainRecipeName(t *testing.T) {
	bundle := Bundle{
		BundleRecipes: bundle_recipes,
	}

	require.Equal(t, true, bundle.ContainsName(bundle_recipe.Recipe.Name))
}

func TestBundle_ShouldNotContainRecipeName(t *testing.T) {
	bundle := Bundle{
		BundleRecipes: bundle_recipes,
	}

	require.Equal(t, false, bundle.ContainsName("some name"))
}

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
