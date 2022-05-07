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

type mockDetector struct {
	detectionStatus map[string]execution.RecipeStatusType
}

var (
	detectionStatus = make(map[string]execution.RecipeStatusType)
	bundlerTestImpl = struct {
		discoveryManifest types.DiscoveryManifest
		recipeCache       []*types.OpenInstallationRecipe
		recipeRepository  *RecipeRepository
		ctx               context.Context
		processEvaluator  *mockDetector
		scriptedEvaluator *mockDetector
		recipeDetector    *RecipeDetector
	}{
		processEvaluator:  &mockDetector{detectionStatus: detectionStatus},
		scriptedEvaluator: &mockDetector{detectionStatus: detectionStatus},
	}
)

func TestCreateAdditionalTargetedBundleShouldNotSkipCoreRecipes(t *testing.T) {
	setup()
	infraRecipe := NewRecipeBuilder().Name(types.InfraAgentRecipeName).Build()
	loggingRecipe := NewRecipeBuilder().Name(types.LoggingRecipeName).Build()
	goldenRecipe := NewRecipeBuilder().Name(types.GoldenRecipeName).Build()
	mysqlRecipe := NewRecipeBuilder().Name("mysql").Build()
	bundler := createTestBundler()
	withRecipeStatusDetector(bundler, types.InfraAgentRecipeName, execution.RecipeStatusTypes.AVAILABLE, infraRecipe)
	withRecipeStatusDetector(bundler, types.LoggingRecipeName, execution.RecipeStatusTypes.AVAILABLE, loggingRecipe)
	withRecipeStatusDetector(bundler, types.GoldenRecipeName, execution.RecipeStatusTypes.AVAILABLE, goldenRecipe)
	withRecipeStatusDetector(bundler, "mysql", execution.RecipeStatusTypes.AVAILABLE, mysqlRecipe)

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
	setup()
	anotherRecipe := "x"
	xRecipe := NewRecipeBuilder().Name(anotherRecipe).Build()
	mysqlRecipe := NewRecipeBuilder().Name("mysql").Build()
	bundler := createTestBundler()
	withRecipeStatusDetector(bundler, anotherRecipe, execution.RecipeStatusTypes.AVAILABLE, xRecipe)
	withRecipeStatusDetector(bundler, "mysql", execution.RecipeStatusTypes.AVAILABLE, mysqlRecipe)

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
	setup()
	infraRecipe := NewRecipeBuilder().Name(types.InfraAgentRecipeName).Build()
	loggingRecipe := NewRecipeBuilder().Name(types.LoggingRecipeName).Build()
	goldenRecipe := NewRecipeBuilder().Name(types.GoldenRecipeName).Build()
	mysqlRecipe := NewRecipeBuilder().Name("mysql").Build()
	bundler := createTestBundler()
	withRecipeStatusDetector(bundler, types.InfraAgentRecipeName, execution.RecipeStatusTypes.AVAILABLE, infraRecipe)
	withRecipeStatusDetector(bundler, types.LoggingRecipeName, execution.RecipeStatusTypes.AVAILABLE, loggingRecipe)
	withRecipeStatusDetector(bundler, types.GoldenRecipeName, execution.RecipeStatusTypes.AVAILABLE, goldenRecipe)
	withRecipeStatusDetector(bundler, "mysql", execution.RecipeStatusTypes.AVAILABLE, mysqlRecipe)

	coreBundle := bundler.CreateCoreBundle()

	require.Equal(t, 2, len(coreBundle.BundleRecipes))
	require.NotNil(t, findRecipeByName(coreBundle, types.InfraAgentRecipeName))
	require.NotNil(t, findRecipeByName(coreBundle, types.LoggingRecipeName))
	require.Nil(t, findRecipeByName(coreBundle, types.GoldenRecipeName))
	require.Nil(t, findRecipeByName(coreBundle, "mysql"))
}

func TestCreateAdditionalGuidedBundleShouldSkipCoreRecipes(t *testing.T) {
	setup()
	infraRecipe := NewRecipeBuilder().Name(types.InfraAgentRecipeName).Build()
	loggingRecipe := NewRecipeBuilder().Name(types.LoggingRecipeName).Build()
	goldenRecipe := NewRecipeBuilder().Name(types.GoldenRecipeName).Build()
	mysqlRecipe := NewRecipeBuilder().Name("mysql").Build()
	bundler := createTestBundler()
	withRecipeStatusDetector(bundler, types.InfraAgentRecipeName, execution.RecipeStatusTypes.AVAILABLE, infraRecipe)
	withRecipeStatusDetector(bundler, types.LoggingRecipeName, execution.RecipeStatusTypes.AVAILABLE, loggingRecipe)
	withRecipeStatusDetector(bundler, types.GoldenRecipeName, execution.RecipeStatusTypes.AVAILABLE, goldenRecipe)
	withRecipeStatusDetector(bundler, "mysql", execution.RecipeStatusTypes.AVAILABLE, mysqlRecipe)

	addBundle := bundler.CreateAdditionalGuidedBundle()

	require.Equal(t, 2, len(addBundle.BundleRecipes))
	require.NotNil(t, findRecipeByName(addBundle, "mysql"))
}

func TestCreateCoreBundleShouldDetectAvailableStatus(t *testing.T) {
	setup()
	infraRecipe := NewRecipeBuilder().Name(types.InfraAgentRecipeName).Build()
	loggingRecipe := NewRecipeBuilder().Name(types.LoggingRecipeName).Build()
	goldenRecipe := NewRecipeBuilder().Name(types.GoldenRecipeName).Build()
	mysqlRecipe := NewRecipeBuilder().Name("mysql").Build()
	bundler := createTestBundler()
	withRecipeStatusDetector(bundler, types.InfraAgentRecipeName, execution.RecipeStatusTypes.AVAILABLE, infraRecipe)
	withRecipeStatusDetector(bundler, types.LoggingRecipeName, execution.RecipeStatusTypes.AVAILABLE, loggingRecipe)
	withRecipeStatusDetector(bundler, types.GoldenRecipeName, execution.RecipeStatusTypes.AVAILABLE, goldenRecipe)
	withRecipeStatusDetector(bundler, "mysql", execution.RecipeStatusTypes.AVAILABLE, mysqlRecipe)

	coreBundle := bundler.CreateCoreBundle()

	require.Equal(t, 2, len(coreBundle.BundleRecipes))
	for _, r := range coreBundle.BundleRecipes {
		lastStatusIndex := len(r.DetectedStatuses) - 1
		require.Equal(t, 2, len(r.DetectedStatuses))
		require.Equal(t, execution.RecipeStatusTypes.AVAILABLE, r.DetectedStatuses[lastStatusIndex].Status)
	}
}

func TestCreateCoreBundleShouldIncludeDependencies(t *testing.T) {

	setup()
	infraRecipe := NewRecipeBuilder().Name(types.InfraAgentRecipeName).Build()
	infraRecipe.Dependencies = []string{"dep1", "dep2"}
	loggingRecipe := NewRecipeBuilder().Name(types.LoggingRecipeName).Build()
	loggingRecipe.Dependencies = []string{"dep2"}
	dep1Recipe := NewRecipeBuilder().Name("dep1").Build()
	dep2Recipe := NewRecipeBuilder().Name("dep2").Build()

	bundler := createTestBundler()
	withRecipeStatusDetector(bundler, types.InfraAgentRecipeName, execution.RecipeStatusTypes.AVAILABLE, infraRecipe)
	withRecipeStatusDetector(bundler, types.LoggingRecipeName, execution.RecipeStatusTypes.AVAILABLE, loggingRecipe)
	withRecipeStatusDetector(bundler, "dep1", execution.RecipeStatusTypes.AVAILABLE, dep1Recipe)
	withRecipeStatusDetector(bundler, "dep2", execution.RecipeStatusTypes.AVAILABLE, dep2Recipe)

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

	setup()
	infraRecipe := NewRecipeBuilder().Name(types.InfraAgentRecipeName).Build()
	infraRecipe.Dependencies = []string{"dep1", "dep2", "dep3"}
	loggingRecipe := NewRecipeBuilder().Name(types.LoggingRecipeName).Build()
	loggingRecipe.Dependencies = []string{"dep1", "dep2"}
	dep1Recipe := NewRecipeBuilder().Name("dep1").Build()
	dep2Recipe := NewRecipeBuilder().Name("dep2").Build()

	bundler := createTestBundler()
	withRecipeStatusDetector(bundler, types.InfraAgentRecipeName, execution.RecipeStatusTypes.AVAILABLE, infraRecipe)
	withRecipeStatusDetector(bundler, types.LoggingRecipeName, execution.RecipeStatusTypes.AVAILABLE, loggingRecipe)
	withRecipeStatusDetector(bundler, "dep1", execution.RecipeStatusTypes.AVAILABLE, dep1Recipe)
	withRecipeStatusDetector(bundler, "dep2", execution.RecipeStatusTypes.AVAILABLE, dep2Recipe)

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
	setup()
	addRecipeToCache("id1", "mysql")
	bundler := createTestBundler()

	coreBundle := bundler.CreateCoreBundle()

	require.Equal(t, 0, len(coreBundle.BundleRecipes))
}

func TestCreateAdditionalGuidedBundleShouldBundleRecipesThatHaveNullStatus(t *testing.T) {

	setup()
	mysqlRecipe := NewRecipeBuilder().Name("mysql").Build()

	bundler := createTestBundler()
	withRecipeStatusDetector(bundler, types.InfraAgentRecipeName, execution.RecipeStatusTypes.NULL, mysqlRecipe)

	bundle := bundler.CreateAdditionalGuidedBundle()

	require.Equal(t, 0, len(bundle.BundleRecipes))
	require.Nil(t, findRecipeByName(bundle, "mysql"))
}

func TestCreateCoreBundleShouldBundleDependencyOrRecipeWhenNotDetected(t *testing.T) {

	setup()
	infraRecipe := NewRecipeBuilder().Name(types.InfraAgentRecipeName).Build()
	infraRecipe.Dependencies = []string{"dep1"}
	dep1Recipe := NewRecipeBuilder().Name("dep1").Build()

	bundler := createTestBundler()
	coreBundle := bundler.CreateCoreBundle()
	withRecipeStatusDetector(bundler, types.InfraAgentRecipeName, execution.RecipeStatusTypes.AVAILABLE, dep1Recipe)
	withRecipeStatusDetector(bundler, "dep1", execution.RecipeStatusTypes.NULL, dep1Recipe)

	require.Equal(t, 0, len(coreBundle.BundleRecipes))
	require.Nil(t, findRecipeByName(coreBundle, types.InfraAgentRecipeName))
	require.Nil(t, findRecipeByName(coreBundle, "dep1"))
}

func TestNewBundlerShouldCreate(t *testing.T) {
	setup()
	d := make(map[string]*RecipeDetection)
	bundler := NewBundler(bundlerTestImpl.ctx, d)

	require.NotNil(t, bundler)
	require.Equal(t, bundlerTestImpl.ctx, bundler.Context)
	require.Equal(t, d, bundler.Detections)
}

func TestBundleRecipeShouldNotBeAvailableWhenDependencyMissingInRepo(t *testing.T) {
	setup()
	r := addRecipeWithDependenciesToCache("id1", types.InfraAgentRecipeName, []string{"dep1"})
	bundler := createTestBundler()
	withRecipeStatusDetector(bundler, types.InfraAgentRecipeName, execution.RecipeStatusTypes.AVAILABLE, r)

	coreBundle := bundler.CreateCoreBundle()
	fmt.Printf("coreBundle.BundleRecipes = %+v\n", coreBundle.BundleRecipes)
	require.Equal(t, 0, len(coreBundle.BundleRecipes))
}

func TestBundleRecipeShouldNotBeAvailableWhenRecipeNotAvailable(t *testing.T) {

	setup()
	infraRecipe := NewRecipeBuilder().Name(types.InfraAgentRecipeName).Build()
	infraRecipe.Dependencies = []string{"dep1"}
	dep1Recipe := NewRecipeBuilder().Name("dep1").Build()

	bundler := createTestBundler()
	withRecipeStatusDetector(bundler, types.InfraAgentRecipeName, execution.RecipeStatusTypes.NULL, infraRecipe)
	withRecipeStatusDetector(bundler, "dep1", execution.RecipeStatusTypes.AVAILABLE, dep1Recipe)
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

func setup() {
	bundlerTestImpl.recipeDetector = newRecipeDetector(bundlerTestImpl.processEvaluator, bundlerTestImpl.scriptedEvaluator)
	bundlerTestImpl.discoveryManifest = types.DiscoveryManifest{
		OS: "linux",
	}
	bundlerTestImpl.recipeCache = []*types.OpenInstallationRecipe{}
	bundlerTestImpl.recipeRepository = NewRecipeRepository(bundlerRecipeLoader, &bundlerTestImpl.discoveryManifest)
}

func newBundler(context context.Context, detections map[string]*RecipeDetection) *Bundler {
	return &Bundler{
		Context:             context,
		Detections:          detections,
		cachedBundleRecipes: make(map[string]*BundleRecipe),
	}
}

func createTestBundler() *Bundler {

	d := make(map[string]*RecipeDetection)
	bundler := newBundler(bundlerTestImpl.ctx, d)

	return bundler
}

func withRecipeStatusDetector(bundler *Bundler, recipeName string, status execution.RecipeStatusType, recipe *types.OpenInstallationRecipe) {

	bundler.Detections[recipeName] = &RecipeDetection{
		Status: status,
		Recipe: recipe,
	}
}

func bundlerRecipeLoader() ([]*types.OpenInstallationRecipe, error) {
	return bundlerTestImpl.recipeCache, nil
}

func addRecipeToCache(id string, name string) *types.OpenInstallationRecipe {
	r := NewRecipeBuilder().ID(id).Name(name).TargetOs(types.OpenInstallationOperatingSystemTypes.LINUX).Build()
	bundlerTestImpl.recipeCache = append(bundlerTestImpl.recipeCache, r)
	return r
}

func addRecipeWithDependenciesToCache(id string, name string, dependencies []string) *types.OpenInstallationRecipe {
	r := NewRecipeBuilder().ID(id).Name(name).TargetOs(types.OpenInstallationOperatingSystemTypes.LINUX).Build()

	if len(dependencies) > 0 {
		r.Dependencies = dependencies
	}

	bundlerTestImpl.recipeCache = append(bundlerTestImpl.recipeCache, r)
	return r
}

func (d mockDetector) DetectionStatus(ctx context.Context, recipe *types.OpenInstallationRecipe) execution.RecipeStatusType {
	if v, ok := d.detectionStatus[recipe.Name]; ok {
		return v
	}
	return execution.RecipeStatusTypes.NULL
}
