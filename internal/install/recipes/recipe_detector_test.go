package recipes

import (
	"context"
	"testing"

	"github.com/newrelic/newrelic-cli/internal/install/execution"
	"github.com/newrelic/newrelic-cli/internal/install/types"
	"github.com/stretchr/testify/require"
)

var (
	mockProcessEvaluator = &mockRecipeEvaluator{}
	mockScriptEvaluator  = &mockRecipeEvaluator{}
	recipeDetectorSUT    = *newRecipeDetector(mockProcessEvaluator, mockScriptEvaluator)
)

func TestRecipeDetector_ShouldGetProcessEvalStatusNull(t *testing.T) {
	recipe := givenRecipeWithNoProcessMatching()
	withProcessEvaluatorReturnStatus(execution.RecipeStatusTypes.NULL)
	withScriptEvaluatorReturnStatus(execution.RecipeStatusTypes.AVAILABLE)

	statusDetectionResult := recipeDetectorSUT.DetectRecipes(ctx, []*types.OpenInstallationRecipe{recipe})
	actual := statusDetectionResult[recipe]
	require.Equal(t, 1, len(statusDetectionResult))
	require.Equal(t, execution.RecipeStatusTypes.NULL, actual)
}

func TestRecipeDetector_ShouldGetProcessEvalStatusAvaliable(t *testing.T) {
	recipe := givenRecipeWithNoProcessMatching()
	withEmptyPreInstallRequiredAtDiscoverSection()
	withProcessEvaluatorReturnStatus(execution.RecipeStatusTypes.AVAILABLE)
	withScriptEvaluatorReturnStatus(execution.RecipeStatusTypes.DETECTED)

	statusDetectionResult := recipeDetectorSUT.DetectRecipes(ctx, []*types.OpenInstallationRecipe{recipe})
	actual := statusDetectionResult[recipe]
	require.Equal(t, 1, len(statusDetectionResult))
	require.Equal(t, execution.RecipeStatusTypes.AVAILABLE, actual)
}

func TestRecipeDetector_ShouldGetScriptEvalStatusNull(t *testing.T) {
	givenRecipeWithNoProcessMatching()
	withPreInstallRequiredAtDiscoverSection()
	withProcessEvaluatorReturnStatus(execution.RecipeStatusTypes.AVAILABLE)
	withScriptEvaluatorReturnStatus(execution.RecipeStatusTypes.NULL)

	statusDetectionResult := recipeDetectorSUT.DetectRecipes(ctx, []*types.OpenInstallationRecipe{recipe})
	actual := statusDetectionResult[recipe]
	require.Equal(t, 1, len(statusDetectionResult))
	require.Equal(t, execution.RecipeStatusTypes.NULL, actual)
}

func TestRecipeDetector_ShouldGetScriptEvalStatusDetected(t *testing.T) {
	givenRecipeWithNoProcessMatching()
	withPreInstallRequiredAtDiscoverSection()
	withProcessEvaluatorReturnStatus(execution.RecipeStatusTypes.AVAILABLE)
	withScriptEvaluatorReturnStatus(execution.RecipeStatusTypes.DETECTED)

	statusDetectionResult := recipeDetectorSUT.DetectRecipes(ctx, []*types.OpenInstallationRecipe{recipe})
	actual := statusDetectionResult[recipe]
	require.Equal(t, 1, len(statusDetectionResult))
	require.Equal(t, execution.RecipeStatusTypes.DETECTED, actual)
}
func TestRecipeDetector_ShouldGetScriptEvalStatusAvailable(t *testing.T) {
	givenRecipeWithNoProcessMatching()
	withPreInstallRequiredAtDiscoverSection()
	withProcessEvaluatorReturnStatus(execution.RecipeStatusTypes.AVAILABLE)
	withScriptEvaluatorReturnStatus(execution.RecipeStatusTypes.AVAILABLE)

	statusDetectionResult := recipeDetectorSUT.DetectRecipes(ctx, []*types.OpenInstallationRecipe{recipe})
	actual := statusDetectionResult[recipe]
	require.Equal(t, 1, len(statusDetectionResult))
	require.Equal(t, execution.RecipeStatusTypes.AVAILABLE, actual)
}

func givenRecipeWithNoProcessMatching() *types.OpenInstallationRecipe {
	return createRecipe("0", "recipe1")
}

func withEmptyPreInstallRequiredAtDiscoverSection() {
}

func withPreInstallRequiredAtDiscoverSection() *types.OpenInstallationRecipe {
	recipe = createRecipe("0", "recipe1")
	recipe.PreInstall = types.OpenInstallationPreInstallConfiguration{
		RequireAtDiscovery: "pre-install script mock"}

	return recipe
}

func withProcessEvaluatorReturnStatus(status execution.RecipeStatusType) {
	mockProcessEvaluator.status = status
}

func withScriptEvaluatorReturnStatus(status execution.RecipeStatusType) {
	mockScriptEvaluator.status = status
}

type mockRecipeEvaluator struct {
	status execution.RecipeStatusType
}

func (mre *mockRecipeEvaluator) DetectionStatus(ctx context.Context, recipe *types.OpenInstallationRecipe) execution.RecipeStatusType {
	return mre.status
}
