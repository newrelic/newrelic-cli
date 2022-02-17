package recipes

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBundle_ShouldAddRecipes(t *testing.T) {
	bundle := Bundle{
		BundleRecipes: []*BundleRecipe{
			{
				Recipe: NewRecipeBuilder().ID("0").Name("recipe1").Build(),
			},
		},
	}
	newRecipe := &BundleRecipe{
		Recipe: NewRecipeBuilder().ID("2").Name("recipe3").Build(),
	}

	bundle.AddRecipe(newRecipe)
	require.Equal(t, len(bundle.BundleRecipes), 2)
	require.True(t, true, bundle.ContainsName(newRecipe.Recipe.Name))
}

func TestBundle_ShouldNotUpdateRecipe(t *testing.T) {
	bundleRecipe := &BundleRecipe{
		Recipe: NewRecipeBuilder().ID("0").Name("recipe1").Build(),
	}
	bundle := Bundle{
		BundleRecipes: []*BundleRecipe{bundleRecipe},
	}

	bundle.AddRecipe(bundleRecipe)
	bundle.AddRecipe(bundleRecipe)
	bundle.AddRecipe(bundleRecipe)

	require.Equal(t, len(bundle.BundleRecipes), 1)
	require.Equal(t, true, bundle.ContainsName(bundleRecipe.Recipe.Name))
}

func TestBundle_ShouldContainRecipeName(t *testing.T) {
	bundle := Bundle{
		BundleRecipes: []*BundleRecipe{
			{
				Recipe: NewRecipeBuilder().ID("0").Name("recipe1").Build(),
			},
		},
	}

	require.Equal(t, true, bundle.ContainsName("recipe1"))
}

func TestBundle_ShouldNotContainRecipeName(t *testing.T) {
	bundle := Bundle{
		BundleRecipes: []*BundleRecipe{
			{
				Recipe: NewRecipeBuilder().ID("0").Name("recipe1").Build(),
			},
		},
	}

	require.Equal(t, false, bundle.ContainsName("some other name"))
}
