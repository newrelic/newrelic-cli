package recipes

import (
	"testing"

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

func TestBundle_ShouldNotUpdateRecipe(t *testing.T) {

	bundle := Bundle{
		BundleRecipes: bundle_recipes,
	}

	bundle.AddRecipe(bundle_recipe)
	bundle.AddRecipe(bundle_recipe)
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
