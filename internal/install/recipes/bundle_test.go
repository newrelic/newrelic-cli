package recipes

import (
	"testing"

	"github.com/newrelic/newrelic-cli/internal/install/types"
	"github.com/stretchr/testify/require"
)

var (
	bundle_recipe  = createRecipe("0", "recipe1")
	bundle_recipes = []*types.OpenInstallationRecipe{
		createRecipe("0", "recipe1"),
		createRecipe("1", "recipe2"),
	}
)

func TestBundleRecipe_ShouldReferenceSameRecipe(t *testing.T) {
	bundleRecipe := NewBundleRecipe(bundle_recipe)
	actual := bundleRecipe.recipe
	require.Equal(t, bundle_recipe, actual)
}

func TestBundle_ShouldCreateWithSameRecipeReference(t *testing.T) {
	bundle := NewBundle(bundle_recipes)
	actual := bundle.BundleRecipes

	require.Equal(t, len(bundle_recipes), len(actual))

	for i := 0; i < len(bundle_recipes); i++ {
		_, ok := bundle.Contains(bundle_recipes[i])
		require.Equal(t, true, ok)
	}
}

func TestBundle_ShouldAddRecipe(t *testing.T) {

	bundle := NewBundle(bundle_recipes)
	newRecipe := createRecipe("2", "recipe3")

	bundle.AddRecipe(newRecipe)
	require.Equal(t, len(bundle.BundleRecipes), len(bundle_recipes)+1)
	_, ok := bundle.Contains(newRecipe)
	require.Equal(t, true, ok)
}

func TestBundle_ShouldUpdateRecipe(t *testing.T) {

	bundle := NewBundle(bundle_recipes)
	bundle.AddRecipe(bundle_recipe)

	require.Equal(t, len(bundle.BundleRecipes), len(bundle_recipes))
	_, ok := bundle.Contains(bundle_recipe)
	require.Equal(t, true, ok)
}

func TestBundle_ShouldUpdateRecipes(t *testing.T) {

	bundle := NewBundle(bundle_recipes)
	newRecipe := createRecipe("2", "recipe3")
	newRecipes := []*types.OpenInstallationRecipe{
		newRecipe,
		bundle_recipe,
	}

	bundle.AddRecipes(newRecipes)
	require.Equal(t, len(bundle.BundleRecipes), len(bundle_recipes)+1)
	_, ok := bundle.Contains(newRecipe)
	require.Equal(t, true, ok)
	_, ok = bundle.Contains(bundle_recipe)
	require.Equal(t, true, ok)
}

func TestBundle_ShouldContainRecipeName(t *testing.T) {
	bundle := NewBundle(bundle_recipes)
	_, ok := bundle.ContainsName(bundle_recipe.Name)
	require.Equal(t, true, ok)
}

func TestBundle_ShouldNotContainRecipeName(t *testing.T) {
	bundle := NewBundle(bundle_recipes)
	_, ok := bundle.ContainsName("some name")
	require.Equal(t, false, ok)
}
func TestBundle_ShouldNotContainRecipe(t *testing.T) {
	bundle := NewBundle(bundle_recipes)
	newRecipe := createRecipe("1", "recipe2")
	_, ok := bundle.Contains(newRecipe)
	require.Equal(t, false, ok)
}
