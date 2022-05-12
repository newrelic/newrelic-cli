package recipes

import (
	"context"
	"fmt"
	"testing"

	"github.com/newrelic/newrelic-cli/internal/install/execution"
	"github.com/newrelic/newrelic-cli/internal/install/types"

	"strings"

	"github.com/stretchr/testify/require"
)

var (
	bundlerTestImpl = struct {
		discoveryManifest types.DiscoveryManifest
		recipeCache       []*types.OpenInstallationRecipe
		recipeRepository  *RecipeRepository
		ctx               context.Context
	}{}
)

func TestCreateAdditionalTargetedBundleShouldNotSkipCoreRecipes(t *testing.T) {
	infraRecipe := NewRecipeBuilder().Name(types.InfraAgentRecipeName).Build()
	loggingRecipe := NewRecipeBuilder().Name(types.LoggingRecipeName).Build()
	goldenRecipe := NewRecipeBuilder().Name(types.GoldenRecipeName).Build()
	mysqlRecipe := NewRecipeBuilder().Name("mysql").Build()
	bundler := createTestBundler()
	withAvailableRecipe(bundler, types.InfraAgentRecipeName, execution.RecipeStatusTypes.AVAILABLE, infraRecipe)
	withAvailableRecipe(bundler, types.LoggingRecipeName, execution.RecipeStatusTypes.AVAILABLE, loggingRecipe)
	withAvailableRecipe(bundler, types.GoldenRecipeName, execution.RecipeStatusTypes.AVAILABLE, goldenRecipe)
	withAvailableRecipe(bundler, "mysql", execution.RecipeStatusTypes.AVAILABLE, mysqlRecipe)

	recipeNames := []string{
		"mysql",
		types.InfraAgentRecipeName,
		types.LoggingRecipeName,
		types.GoldenRecipeName,
	}
	addBundle := bundler.CreateAdditionalTargetedBundle(recipeNames)

	require.Equal(t, 4, len(addBundle.BundleRecipes))
	require.NotNil(t, findRecipeByName(addBundle, types.InfraAgentRecipeName))
	require.NotNil(t, findRecipeByName(addBundle, types.LoggingRecipeName))
	require.NotNil(t, findRecipeByName(addBundle, types.GoldenRecipeName))
	require.NotNil(t, findRecipeByName(addBundle, "mysql"))
}

func TestCreateAdditionalTargetedBundleShouldNotDetectOtherRecipes(t *testing.T) {
	anotherRecipe := "x"
	xRecipe := NewRecipeBuilder().Name(anotherRecipe).Build()
	mysqlRecipe := NewRecipeBuilder().Name("mysql").Build()
	bundler := createTestBundler()
	withAvailableRecipe(bundler, anotherRecipe, execution.RecipeStatusTypes.AVAILABLE, xRecipe)
	withAvailableRecipe(bundler, "mysql", execution.RecipeStatusTypes.AVAILABLE, mysqlRecipe)

	recipeNames := []string{
		"mysql",
	}
	addBundle := bundler.CreateAdditionalTargetedBundle(recipeNames)

	r := findRecipeByName(addBundle, "mysql")
	require.Equal(t, 1, len(addBundle.BundleRecipes))
	require.NotNil(t, r)
	require.True(t, r.HasStatus(execution.RecipeStatusTypes.AVAILABLE))
}

func TestCreateCoreBundleShouldContainOnlyCoreBundleRecipes(t *testing.T) {
	infraRecipe := NewRecipeBuilder().Name(types.InfraAgentRecipeName).Build()
	loggingRecipe := NewRecipeBuilder().Name(types.LoggingRecipeName).Build()
	goldenRecipe := NewRecipeBuilder().Name(types.GoldenRecipeName).Build()
	mysqlRecipe := NewRecipeBuilder().Name("mysql").Build()
	bundler := createTestBundler()
	withAvailableRecipe(bundler, types.InfraAgentRecipeName, execution.RecipeStatusTypes.AVAILABLE, infraRecipe)
	withAvailableRecipe(bundler, types.LoggingRecipeName, execution.RecipeStatusTypes.AVAILABLE, loggingRecipe)
	withAvailableRecipe(bundler, types.GoldenRecipeName, execution.RecipeStatusTypes.AVAILABLE, goldenRecipe)
	withAvailableRecipe(bundler, "mysql", execution.RecipeStatusTypes.AVAILABLE, mysqlRecipe)

	coreBundle := bundler.CreateCoreBundle()

	require.Equal(t, 2, len(coreBundle.BundleRecipes))
	require.NotNil(t, findRecipeByName(coreBundle, types.InfraAgentRecipeName))
	require.NotNil(t, findRecipeByName(coreBundle, types.LoggingRecipeName))
	require.Nil(t, findRecipeByName(coreBundle, types.GoldenRecipeName))
	require.Nil(t, findRecipeByName(coreBundle, "mysql"))
}

func TestCreateAdditionalGuidedBundleShouldSkipCoreRecipes(t *testing.T) {
	infraRecipe := NewRecipeBuilder().Name(types.InfraAgentRecipeName).Build()
	loggingRecipe := NewRecipeBuilder().Name(types.LoggingRecipeName).Build()
	goldenRecipe := NewRecipeBuilder().Name(types.GoldenRecipeName).Build()
	mysqlRecipe := NewRecipeBuilder().Name("mysql").Build()
	bundler := createTestBundler()
	withAvailableRecipe(bundler, types.InfraAgentRecipeName, execution.RecipeStatusTypes.AVAILABLE, infraRecipe)
	withAvailableRecipe(bundler, types.LoggingRecipeName, execution.RecipeStatusTypes.AVAILABLE, loggingRecipe)
	withAvailableRecipe(bundler, types.GoldenRecipeName, execution.RecipeStatusTypes.AVAILABLE, goldenRecipe)
	withAvailableRecipe(bundler, "mysql", execution.RecipeStatusTypes.AVAILABLE, mysqlRecipe)

	addBundle := bundler.CreateAdditionalGuidedBundle()

	require.Equal(t, 2, len(addBundle.BundleRecipes))
	require.NotNil(t, findRecipeByName(addBundle, "mysql"))
}

func TestCreateCoreBundleShouldDetectAvailableStatus(t *testing.T) {
	infraRecipe := NewRecipeBuilder().Name(types.InfraAgentRecipeName).Build()
	loggingRecipe := NewRecipeBuilder().Name(types.LoggingRecipeName).Build()
	goldenRecipe := NewRecipeBuilder().Name(types.GoldenRecipeName).Build()
	mysqlRecipe := NewRecipeBuilder().Name("mysql").Build()
	bundler := createTestBundler()
	withAvailableRecipe(bundler, types.InfraAgentRecipeName, execution.RecipeStatusTypes.AVAILABLE, infraRecipe)
	withAvailableRecipe(bundler, types.LoggingRecipeName, execution.RecipeStatusTypes.AVAILABLE, loggingRecipe)
	withAvailableRecipe(bundler, types.GoldenRecipeName, execution.RecipeStatusTypes.AVAILABLE, goldenRecipe)
	withAvailableRecipe(bundler, "mysql", execution.RecipeStatusTypes.AVAILABLE, mysqlRecipe)

	coreBundle := bundler.CreateCoreBundle()

	require.Equal(t, 2, len(coreBundle.BundleRecipes))
	for _, r := range coreBundle.BundleRecipes {
		require.Equal(t, 1, len(r.DetectedStatuses))
		require.Equal(t, execution.RecipeStatusTypes.AVAILABLE, r.DetectedStatuses[0].Status)
	}
}

func TestCreateCoreBundleShouldIncludeDependencies(t *testing.T) {

	infraRecipe := NewRecipeBuilder().Name(types.InfraAgentRecipeName).Build()
	infraRecipe.Dependencies = []string{"dep1", "dep2"}
	loggingRecipe := NewRecipeBuilder().Name(types.LoggingRecipeName).Build()
	loggingRecipe.Dependencies = []string{"dep2"}
	dep1Recipe := NewRecipeBuilder().Name("dep1").Build()
	dep2Recipe := NewRecipeBuilder().Name("dep2").Build()

	bundler := createTestBundler()
	withAvailableRecipe(bundler, types.InfraAgentRecipeName, execution.RecipeStatusTypes.AVAILABLE, infraRecipe)
	withAvailableRecipe(bundler, types.LoggingRecipeName, execution.RecipeStatusTypes.AVAILABLE, loggingRecipe)
	withAvailableRecipe(bundler, "dep1", execution.RecipeStatusTypes.AVAILABLE, dep1Recipe)
	withAvailableRecipe(bundler, "dep2", execution.RecipeStatusTypes.AVAILABLE, dep2Recipe)

	coreBundle := bundler.CreateCoreBundle()

	require.Equal(t, 2, len(coreBundle.BundleRecipes))
	infraBundleRecipe := findRecipeByName(coreBundle, types.InfraAgentRecipeName)
	loggingBundleRecipe := findRecipeByName(coreBundle, types.LoggingRecipeName)
	require.NotNil(t, infraBundleRecipe)
	require.NotNil(t, loggingBundleRecipe)
	require.Nil(t, findRecipeByName(coreBundle, "dep1"))
	require.Nil(t, findRecipeByName(coreBundle, "dep2"))
	require.NotNil(t, findDependencyByName(infraBundleRecipe, "dep1"))
	require.NotNil(t, findDependencyByName(infraBundleRecipe, "dep2"))
	require.Nil(t, findDependencyByName(loggingBundleRecipe, "dep1"))
	require.NotNil(t, findDependencyByName(loggingBundleRecipe, "dep2"))

}

func TestCreateCoreBundleShouldNotIncludeRecipeWithInvalidDependencies(t *testing.T) {

	infraRecipe := NewRecipeBuilder().Name(types.InfraAgentRecipeName).Build()
	infraRecipe.Dependencies = []string{"dep1", "dep2", "dep3"}
	loggingRecipe := NewRecipeBuilder().Name(types.LoggingRecipeName).Build()
	loggingRecipe.Dependencies = []string{"dep1", "dep2"}
	dep1Recipe := NewRecipeBuilder().Name("dep1").Build()
	dep2Recipe := NewRecipeBuilder().Name("dep2").Build()

	bundler := createTestBundler()
	withAvailableRecipe(bundler, types.InfraAgentRecipeName, execution.RecipeStatusTypes.AVAILABLE, infraRecipe)
	withAvailableRecipe(bundler, types.LoggingRecipeName, execution.RecipeStatusTypes.AVAILABLE, loggingRecipe)
	withAvailableRecipe(bundler, "dep1", execution.RecipeStatusTypes.AVAILABLE, dep1Recipe)
	withAvailableRecipe(bundler, "dep2", execution.RecipeStatusTypes.AVAILABLE, dep2Recipe)

	coreBundle := bundler.CreateCoreBundle()

	require.Equal(t, 1, len(coreBundle.BundleRecipes))
	require.Nil(t, findRecipeByName(coreBundle, types.InfraAgentRecipeName))
	require.NotNil(t, findRecipeByName(coreBundle, types.LoggingRecipeName))
	require.Nil(t, findRecipeByName(coreBundle, "dep1"))
	require.Nil(t, findRecipeByName(coreBundle, "dep2"))
	require.NotNil(t, findDependencyByName(coreBundle.BundleRecipes[0], "dep1"))
	require.NotNil(t, findDependencyByName(coreBundle.BundleRecipes[0], "dep2"))
	require.Nil(t, findDependencyByName(coreBundle.BundleRecipes[0], "dep3"))
}

func TestCreateCoreBundleShouldCreateEmptyCore(t *testing.T) {
	bundler := createTestBundler()
	coreBundle := bundler.CreateCoreBundle()
	require.Equal(t, 0, len(coreBundle.BundleRecipes))
}

func TestCreateAdditionalGuidedBundleShouldBundleRecipesThatHaveNullStatus(t *testing.T) {
	bundler := createTestBundler()
	bundle := bundler.CreateAdditionalGuidedBundle()
	require.Equal(t, 0, len(bundle.BundleRecipes))
}

func TestCreateCoreBundleShouldBundleDependencyOrRecipeWhenNotDetected(t *testing.T) {
	infraRecipe := NewRecipeBuilder().Name(types.InfraAgentRecipeName).Build()
	infraRecipe.Dependencies = []string{"dep1"}
	dep1Recipe := NewRecipeBuilder().Name("dep1").Build()

	bundler := createTestBundler()
	coreBundle := bundler.CreateCoreBundle()
	withAvailableRecipe(bundler, types.InfraAgentRecipeName, execution.RecipeStatusTypes.AVAILABLE, dep1Recipe)
	withAvailableRecipe(bundler, "dep1", execution.RecipeStatusTypes.NULL, dep1Recipe)

	require.Equal(t, 0, len(coreBundle.BundleRecipes))
	require.Nil(t, findRecipeByName(coreBundle, types.InfraAgentRecipeName))
	require.Nil(t, findRecipeByName(coreBundle, "dep1"))
}

func TestBundleRecipeShouldNotBeAvailableWhenDependencyMissingInRepo(t *testing.T) {
	infraRecipe := NewRecipeBuilder().Name(types.InfraAgentRecipeName).Build()
	infraRecipe.Dependencies = []string{"dep1", "dep2"}
	bundler := createTestBundler()
	withAvailableRecipe(bundler, types.InfraAgentRecipeName, execution.RecipeStatusTypes.AVAILABLE, infraRecipe)

	coreBundle := bundler.CreateCoreBundle()
	fmt.Printf("coreBundle.BundleRecipes = %+v\n", coreBundle.BundleRecipes)
	require.Equal(t, 0, len(coreBundle.BundleRecipes))
}

func TestBundleRecipeShouldNotBeAvailableWhenRecipeNotAvailable(t *testing.T) {
	infraRecipe := NewRecipeBuilder().Name(types.InfraAgentRecipeName).Build()
	infraRecipe.Dependencies = []string{"dep1"}
	dep1Recipe := NewRecipeBuilder().Name("dep1").Build()

	bundler := createTestBundler()
	withAvailableRecipe(bundler, "dep1", execution.RecipeStatusTypes.AVAILABLE, dep1Recipe)
	coreBundle := bundler.CreateCoreBundle()

	require.Equal(t, 0, len(coreBundle.BundleRecipes))
}

func findRecipeByName(bundle *Bundle, name string) *BundleRecipe {
	for _, r := range bundle.BundleRecipes {
		if strings.EqualFold(r.Recipe.Name, name) {
			return r
		}
	}
	return nil
}

func findDependencyByName(recipe *BundleRecipe, name string) *types.OpenInstallationRecipe {
	for _, dep := range recipe.Dependencies {
		if strings.EqualFold(dep.Recipe.Name, name) {
			return dep.Recipe
		}
		found := findDependencyByName(dep, name)
		if found != nil {
			return found
		}
	}
	return nil
}

func newBundler(context context.Context, availableRecipes RecipeDetectionResults) *Bundler {
	return &Bundler{
		Context:             context,
		AvailableRecipes:    availableRecipes,
		cachedBundleRecipes: make(map[string]*BundleRecipe),
	}
}

func createTestBundler() *Bundler {

	d := RecipeDetectionResults{}
	bundler := newBundler(bundlerTestImpl.ctx, d)

	return bundler
}

func withAvailableRecipe(bundler *Bundler, recipeName string, status execution.RecipeStatusType, recipe *types.OpenInstallationRecipe) {

	bundler.AvailableRecipes = append(bundler.AvailableRecipes, &RecipeDetectionResult{
		Status: status,
		Recipe: recipe,
	})
}
