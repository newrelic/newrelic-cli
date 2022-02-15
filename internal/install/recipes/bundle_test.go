package recipes

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	bundleRecipe = &BundleRecipe{
		Recipe: createRecipe("0", "recipe1"),
	}
	bundleRecipes = []*BundleRecipe{
		{
			Recipe: createRecipe("0", "recipe1"),
		},
		{
			Recipe: createRecipe("1", "recipe2"),
		},
	}
)

func TestBundle_ShouldAddRecipes(t *testing.T) {

	bundle := Bundle{
		BundleRecipes: bundleRecipes,
	}
	newRecipe := &BundleRecipe{
		Recipe: createRecipe("2", "recipe3"),
	}

	bundle.AddRecipe(newRecipe)
	require.Equal(t, len(bundle.BundleRecipes), len(bundleRecipes)+1)
	require.True(t, true, bundle.ContainsName(newRecipe.Recipe.Name))
}

func TestBundle_ShouldNotUpdateRecipe(t *testing.T) {

	bundle := Bundle{
		BundleRecipes: bundleRecipes,
	}

	bundle.AddRecipe(bundleRecipe)
	bundle.AddRecipe(bundleRecipe)
	bundle.AddRecipe(bundleRecipe)

	require.Equal(t, len(bundle.BundleRecipes), len(bundleRecipes))
	require.Equal(t, true, bundle.ContainsName(bundleRecipe.Recipe.Name))
}

func TestBundle_ShouldContainRecipeName(t *testing.T) {
	bundle := Bundle{
		BundleRecipes: bundleRecipes,
	}

	require.Equal(t, true, bundle.ContainsName(bundleRecipe.Recipe.Name))
}

func TestBundle_ShouldNotContainRecipeName(t *testing.T) {
	bundle := Bundle{
		BundleRecipes: bundleRecipes,
	}

	require.Equal(t, false, bundle.ContainsName("some name"))
}
