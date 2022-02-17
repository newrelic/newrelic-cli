package recipes

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-cli/internal/install/execution"
	"github.com/newrelic/newrelic-cli/internal/install/types"
)

var (
	recipeDetectorProcessEvaluator = &mockRecipeEvaluator{}
	recipeDetectorScriptEvaluator  = &mockRecipeEvaluator{}
	recipeDetector                 = *newRecipeDetector(recipeDetectorProcessEvaluator, recipeDetectorScriptEvaluator)
)

func TestRecipeDetectorShouldGetProcessEvalStatusNull(t *testing.T) {
	recipe := givenRecipeWithNoProcessMatching()
	withProcessEvaluatorReturnStatus(execution.RecipeStatusTypes.NULL)
	withScriptEvaluatorReturnStatus(execution.RecipeStatusTypes.AVAILABLE)

	statusDetectionResult := recipeDetector.detectRecipe(context.Background(), recipe)
	actual := statusDetectionResult
	require.Equal(t, execution.RecipeStatusTypes.NULL, actual)
}

func TestRecipeDetectorShouldGetProcessEvalStatusAvaliable(t *testing.T) {
	recipe := givenRecipeWithNoProcessMatching()
	withEmptyPreInstallRequiredAtDiscoverSection()
	withProcessEvaluatorReturnStatus(execution.RecipeStatusTypes.AVAILABLE)
	withScriptEvaluatorReturnStatus(execution.RecipeStatusTypes.DETECTED)

	statusDetectionResult := recipeDetector.detectRecipe(context.Background(), recipe)
	actual := statusDetectionResult
	require.Equal(t, execution.RecipeStatusTypes.AVAILABLE, actual)
}

func TestRecipeDetectorShouldGetScriptEvalStatusNull(t *testing.T) {
	recipe := withPreInstallRequiredAtDiscoverSection()
	withPreInstallRequiredAtDiscoverSection()
	withProcessEvaluatorReturnStatus(execution.RecipeStatusTypes.AVAILABLE)
	withScriptEvaluatorReturnStatus(execution.RecipeStatusTypes.NULL)

	statusDetectionResult := recipeDetector.detectRecipe(context.Background(), recipe)
	actual := statusDetectionResult
	require.Equal(t, execution.RecipeStatusTypes.NULL, actual)
}

func TestRecipeDetectorShouldGetScriptEvalStatusDetected(t *testing.T) {
	recipe := withPreInstallRequiredAtDiscoverSection()
	withProcessEvaluatorReturnStatus(execution.RecipeStatusTypes.AVAILABLE)
	withScriptEvaluatorReturnStatus(execution.RecipeStatusTypes.DETECTED)

	statusDetectionResult := recipeDetector.detectRecipe(context.Background(), recipe)
	actual := statusDetectionResult
	require.Equal(t, execution.RecipeStatusTypes.DETECTED, actual)
}
func TestRecipeDetectorShouldGetScriptEvalStatusAvailable(t *testing.T) {
	recipe := withPreInstallRequiredAtDiscoverSection()
	withProcessEvaluatorReturnStatus(execution.RecipeStatusTypes.AVAILABLE)
	withScriptEvaluatorReturnStatus(execution.RecipeStatusTypes.AVAILABLE)

	statusDetectionResult := recipeDetector.detectRecipe(context.Background(), recipe)
	actual := statusDetectionResult
	require.Equal(t, execution.RecipeStatusTypes.AVAILABLE, actual)
}

func givenRecipeWithNoProcessMatching() *types.OpenInstallationRecipe {
	return createRecipe("0", "recipe1")
}

func withEmptyPreInstallRequiredAtDiscoverSection() {
}

func withPreInstallRequiredAtDiscoverSection() *types.OpenInstallationRecipe {
	recipe := createRecipe("0", "recipe1")
	recipe.PreInstall = types.OpenInstallationPreInstallConfiguration{
		RequireAtDiscovery: "pre-install script mock"}

	return recipe
}

func withProcessEvaluatorReturnStatus(status execution.RecipeStatusType) {
	recipeDetectorProcessEvaluator.status = status
}

func withScriptEvaluatorReturnStatus(status execution.RecipeStatusType) {
	recipeDetectorScriptEvaluator.status = status
}

type mockRecipeEvaluator struct {
	status execution.RecipeStatusType
}

func (mre *mockRecipeEvaluator) DetectionStatus(ctx context.Context, recipe *types.OpenInstallationRecipe) execution.RecipeStatusType {
	return mre.status
}
