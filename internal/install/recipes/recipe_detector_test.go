package recipes

import (
	"context"
	"testing"

	"github.com/newrelic/newrelic-cli/internal/install/execution"
	"github.com/newrelic/newrelic-cli/internal/install/types"
	"github.com/stretchr/testify/require"
)

var (
	recipe               = createRecipe("0", "recipe1")
	mockProcessEvaluator = &mockRecipeEvaluator{}
	mockScriptEvaluator  = &mockRecipeEvaluator{}
	sut                  = *newRecipeDetector(mockProcessEvaluator, mockScriptEvaluator)
)

func TestRecipeDetector_ShouldGetProcessEvalStatusNull(t *testing.T) {
	givenRecipeWithNoProcessMatching()
	withProcessEvaluatorReturnStatus(execution.RecipeStatusTypes.NULL)
	withScriptEvaluatorReturnStatus(execution.RecipeStatusTypes.AVAILABLE)

	statusDetectionResult := sut.DetectRecipes(ctx, []*types.OpenInstallationRecipe{recipe})
	actual := statusDetectionResult[recipe]
	require.Equal(t, 1, len(statusDetectionResult))
	require.Equal(t, execution.RecipeStatusTypes.NULL, actual)
}

func TestRecipeDetector_ShouldGetProcessEvalStatusAvaliable(t *testing.T) {
	givenRecipeWithNoProcessMatching()
	withEmptyPreInstallRequiredAtDiscoverSection()
	withProcessEvaluatorReturnStatus(execution.RecipeStatusTypes.AVAILABLE)
	withScriptEvaluatorReturnStatus(execution.RecipeStatusTypes.DETECTED)

	statusDetectionResult := sut.DetectRecipes(ctx, []*types.OpenInstallationRecipe{recipe})
	actual := statusDetectionResult[recipe]
	require.Equal(t, 1, len(statusDetectionResult))
	require.Equal(t, execution.RecipeStatusTypes.AVAILABLE, actual)
}

func TestRecipeDetector_ShouldGetScriptEvalStatusNull(t *testing.T) {
	givenRecipeWithNoProcessMatching()
	withPreInstallRequiredAtDiscoverSection()
	withProcessEvaluatorReturnStatus(execution.RecipeStatusTypes.AVAILABLE)
	withScriptEvaluatorReturnStatus(execution.RecipeStatusTypes.NULL)

	statusDetectionResult := sut.DetectRecipes(ctx, []*types.OpenInstallationRecipe{recipe})
	actual := statusDetectionResult[recipe]
	require.Equal(t, 1, len(statusDetectionResult))
	require.Equal(t, execution.RecipeStatusTypes.NULL, actual)
}

func TestRecipeDetector_ShouldGetScriptEvalStatusDetected(t *testing.T) {
	givenRecipeWithNoProcessMatching()
	withPreInstallRequiredAtDiscoverSection()
	withProcessEvaluatorReturnStatus(execution.RecipeStatusTypes.AVAILABLE)
	withScriptEvaluatorReturnStatus(execution.RecipeStatusTypes.DETECTED)

	statusDetectionResult := sut.DetectRecipes(ctx, []*types.OpenInstallationRecipe{recipe})
	actual := statusDetectionResult[recipe]
	require.Equal(t, 1, len(statusDetectionResult))
	require.Equal(t, execution.RecipeStatusTypes.DETECTED, actual)
}
func TestRecipeDetector_ShouldGetScriptEvalStatusAvailable(t *testing.T) {
	givenRecipeWithNoProcessMatching()
	withPreInstallRequiredAtDiscoverSection()
	withProcessEvaluatorReturnStatus(execution.RecipeStatusTypes.AVAILABLE)
	withScriptEvaluatorReturnStatus(execution.RecipeStatusTypes.AVAILABLE)

	statusDetectionResult := sut.DetectRecipes(ctx, []*types.OpenInstallationRecipe{recipe})
	actual := statusDetectionResult[recipe]
	require.Equal(t, 1, len(statusDetectionResult))
	require.Equal(t, execution.RecipeStatusTypes.AVAILABLE, actual)
}

func givenRecipeWithNoProcessMatching() {
	recipe = createRecipe("0", "recipe1")
}

func withEmptyPreInstallRequiredAtDiscoverSection() {
}

func withPreInstallRequiredAtDiscoverSection() {
	recipe.PreInstall = types.OpenInstallationPreInstallConfiguration{
		RequireAtDiscovery: "pre-install script mock"}
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
