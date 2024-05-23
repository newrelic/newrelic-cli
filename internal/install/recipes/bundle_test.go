package recipes

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
	index, found := bundle.ContainsName(newRecipe.Recipe.Name)
	require.Equal(t, len(bundle.BundleRecipes), 2)
	require.Equal(t, true, found)
	require.Equal(t, 1, index)
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
	_, found := bundle.ContainsName(bundleRecipe.Recipe.Name)

	require.Equal(t, true, found)
}

func TestBundle_ShouldContainRecipeName(t *testing.T) {
	bundle := Bundle{
		BundleRecipes: []*BundleRecipe{
			{
				Recipe: NewRecipeBuilder().ID("0").Name("recipe1").Build(),
			},
		},
	}
	_, found := bundle.ContainsName("recipe1")

	require.Equal(t, true, found)
}

func TestBundle_ShouldNotContainRecipeName(t *testing.T) {
	bundle := Bundle{
		BundleRecipes: []*BundleRecipe{
			{
				Recipe: NewRecipeBuilder().ID("0").Name("recipe1").Build(),
			},
		},
	}
	_, found := bundle.ContainsName("some other name")

	require.Equal(t, false, found)
}

func TestBundle_ShouldBeGuided(t *testing.T) {
	bundle := Bundle{
		Type: BundleTypes.ADDITIONALGUIDED,
	}
	assert.True(t, bundle.IsAdditionalGuided())
}

func TestBundle_ShouldNotBeGuidedWhenCore(t *testing.T) {
	bundle := Bundle{
		Type: BundleTypes.CORE,
	}
	assert.False(t, bundle.IsAdditionalGuided())
}

func TestBundle_ShouldNotBeGuidedWhenTargeted(t *testing.T) {
	bundle := Bundle{
		Type: BundleTypes.ADDITIONALTARGETED,
	}
	assert.False(t, bundle.IsAdditionalGuided())
}
