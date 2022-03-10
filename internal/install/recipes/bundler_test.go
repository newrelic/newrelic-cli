package recipes

import (
	"context"
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
	addRecipeToCache("id1", types.InfraAgentRecipeName)
	addRecipeToCache("id2", types.LoggingRecipeName)
	addRecipeToCache("id3", types.GoldenRecipeName)
	addRecipeToCache("id4", "mysql")
	bundler := createTestBundler()
	withRecipeStatusDetector(bundler, types.InfraAgentRecipeName, execution.RecipeStatusTypes.AVAILABLE)
	withRecipeStatusDetector(bundler, types.LoggingRecipeName, execution.RecipeStatusTypes.AVAILABLE)
	withRecipeStatusDetector(bundler, types.GoldenRecipeName, execution.RecipeStatusTypes.AVAILABLE)
	withRecipeStatusDetector(bundler, "mysql", execution.RecipeStatusTypes.AVAILABLE)

	recipeNames := []string{
		"mysql",
		types.LoggingRecipeName,
		types.InfraAgentRecipeName,
		types.GoldenRecipeName,
	}
	addBundle := bundler.CreateAdditionalTargetedBundle(recipeNames)

	require.Equal(t, 4, len(addBundle.BundleRecipes))
	require.NotNil(t, findRecipeByName(addBundle, types.InfraAgentRecipeName))
	require.NotNil(t, findRecipeByName(addBundle, types.LoggingRecipeName))
	require.NotNil(t, findRecipeByName(addBundle, types.GoldenRecipeName))
	require.NotNil(t, findRecipeByName(addBundle, "mysql"))
}

func TestCreateCoreBundleShouldContainOnlyCoreBundleRecipes(t *testing.T) {
	setup()
	addRecipeToCache("id1", types.InfraAgentRecipeName)
	addRecipeToCache("id2", types.LoggingRecipeName)
	addRecipeToCache("id3", types.GoldenRecipeName)
	addRecipeToCache("id4", "mysql")
	bundler := createTestBundler()

	coreBundle := bundler.CreateCoreBundle()

	require.Equal(t, 2, len(coreBundle.BundleRecipes))
	require.NotNil(t, findRecipeByName(coreBundle, types.InfraAgentRecipeName))
	require.NotNil(t, findRecipeByName(coreBundle, types.LoggingRecipeName))
	require.Nil(t, findRecipeByName(coreBundle, types.GoldenRecipeName))
	require.Nil(t, findRecipeByName(coreBundle, "mysql"))
}

func TestCreateAdditionalGuidedBundleShouldSkipCoreRecipes(t *testing.T) {
	setup()
	addRecipeToCache("id1", types.InfraAgentRecipeName)
	addRecipeToCache("id2", types.LoggingRecipeName)
	addRecipeToCache("id3", types.GoldenRecipeName)
	addRecipeToCache("id4", "mysql")
	bundler := createTestBundler()
	withRecipeStatusDetector(bundler, types.InfraAgentRecipeName, execution.RecipeStatusTypes.AVAILABLE)
	withRecipeStatusDetector(bundler, types.LoggingRecipeName, execution.RecipeStatusTypes.AVAILABLE)
	withRecipeStatusDetector(bundler, types.GoldenRecipeName, execution.RecipeStatusTypes.AVAILABLE)
	withRecipeStatusDetector(bundler, "mysql", execution.RecipeStatusTypes.AVAILABLE)

	addBundle := bundler.CreateAdditionalGuidedBundle()

	require.Equal(t, 2, len(addBundle.BundleRecipes))
	require.NotNil(t, findRecipeByName(addBundle, "mysql"))
}

func TestCreateCoreBundleShouldDetectAvailableStatus(t *testing.T) {
	setup()
	addRecipeToCache("id2", types.InfraAgentRecipeName)
	addRecipeToCache("id2", types.LoggingRecipeName)
	addRecipeToCache("id3", types.GoldenRecipeName)
	addRecipeToCache("id4", "mysql")
	bundler := createTestBundler()

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
	addRecipeWithDependenciesToCache("id1", types.InfraAgentRecipeName, []string{"dep1", "dep2"})
	addRecipeWithDependenciesToCache("id2", types.LoggingRecipeName, []string{"dep2"})
	addRecipeToCache("id3", "dep1")
	addRecipeToCache("id4", "dep2")
	bundler := createTestBundler()
	withRecipeStatusDetector(bundler, "dep1", execution.RecipeStatusTypes.AVAILABLE)
	withRecipeStatusDetector(bundler, "dep2", execution.RecipeStatusTypes.AVAILABLE)

	coreBundle := bundler.CreateCoreBundle()

	require.Equal(t, 2, len(coreBundle.BundleRecipes))
	infraRecipe := findRecipeByName(coreBundle, types.InfraAgentRecipeName)
	loggingRecipe := findRecipeByName(coreBundle, types.LoggingRecipeName)
	require.NotNil(t, infraRecipe)
	require.NotNil(t, loggingRecipe)
	require.Nil(t, findRecipeByName(coreBundle, "dep1"))
	require.Nil(t, findRecipeByName(coreBundle, "dep2"))
	require.NotNil(t, findDependencyByName(infraRecipe, "dep1"))
	require.NotNil(t, findDependencyByName(infraRecipe, "dep2"))
	require.Nil(t, findDependencyByName(loggingRecipe, "dep1"))
	require.NotNil(t, findDependencyByName(loggingRecipe, "dep2"))

}

func TestCreateCoreBundleShouldNotIncludeRecipeWithInvalidDependencies(t *testing.T) {
	setup()
	addRecipeWithDependenciesToCache("id1", types.InfraAgentRecipeName, []string{"dep1", "dep2", "dep3"})
	addRecipeWithDependenciesToCache("id2", types.LoggingRecipeName, []string{"dep1", "dep2"})
	addRecipeToCache("id3", "dep1")
	addRecipeToCache("id4", "dep2")
	bundler := createTestBundler()

	withRecipeStatusDetector(bundler, "dep1", execution.RecipeStatusTypes.AVAILABLE)
	withRecipeStatusDetector(bundler, "dep2", execution.RecipeStatusTypes.AVAILABLE)

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

func TestCreateAdditionalGuidedBundleShouldNotBundleRecipesThatHaveNullStatus(t *testing.T) {
	setup()
	addRecipeToCache("id1", "mysql")
	bundler := createTestBundler()
	withRecipeStatusDetector(bundler, "mysql", execution.RecipeStatusTypes.NULL)

	coreBundle := bundler.CreateAdditionalGuidedBundle()

	require.Equal(t, 0, len(coreBundle.BundleRecipes))
	require.Nil(t, findRecipeByName(coreBundle, "mysql"))
}

func TestCreateCoreBundleShouldNotBundleDependencyOrRecipeWhenNotDetected(t *testing.T) {
	setup()
	addRecipeWithDependenciesToCache("id1", types.InfraAgentRecipeName, []string{"dep1"})
	addRecipeToCache("id3", "dep1")
	bundler := createTestBundler()
	withRecipeStatusDetector(bundler, "dep1", execution.RecipeStatusTypes.NULL)

	coreBundle := bundler.CreateCoreBundle()

	require.Equal(t, 0, len(coreBundle.BundleRecipes))
	require.Nil(t, findRecipeByName(coreBundle, types.InfraAgentRecipeName))
	require.Nil(t, findRecipeByName(coreBundle, "dep1"))
}
func TestNewBundlerShouldCreate(t *testing.T) {
	setup()
	bundler := NewBundler(bundlerTestImpl.ctx, bundlerTestImpl.recipeRepository)

	require.NotNil(t, bundler)
	require.Equal(t, bundlerTestImpl.ctx, bundler.Context)
	require.Equal(t, bundlerTestImpl.recipeRepository, bundler.RecipeRepository)
	require.NotNil(t, bundler.RecipeDetector)
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

func newBundler(context context.Context, rr *RecipeRepository, rd *RecipeDetector) *Bundler {
	return &Bundler{
		Context:          context,
		RecipeRepository: rr,
		RecipeDetector:   rd,
	}
}

func createTestBundler() *Bundler {
	bundler := newBundler(bundlerTestImpl.ctx, bundlerTestImpl.recipeRepository, bundlerTestImpl.recipeDetector)

	withRecipeStatusDetector(bundler, types.InfraAgentRecipeName, execution.RecipeStatusTypes.AVAILABLE)
	withRecipeStatusDetector(bundler, types.LoggingRecipeName, execution.RecipeStatusTypes.AVAILABLE)
	withRecipeStatusDetector(bundler, types.GoldenRecipeName, execution.RecipeStatusTypes.AVAILABLE)

	return bundler
}

func withRecipeStatusDetector(bundler *Bundler, recipeName string, status execution.RecipeStatusType) {
	detectionStatus[recipeName] = status

	bundlerTestImpl.processEvaluator = &mockDetector{detectionStatus: detectionStatus}
	bundlerTestImpl.scriptedEvaluator = &mockDetector{detectionStatus: detectionStatus}
	bundlerTestImpl.recipeDetector = newRecipeDetector(bundlerTestImpl.processEvaluator, bundlerTestImpl.scriptedEvaluator)
	bundler.RecipeDetector = bundlerTestImpl.recipeDetector
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
