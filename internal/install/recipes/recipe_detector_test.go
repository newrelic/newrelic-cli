package recipes

import (
	"context"
	"testing"

	"github.com/newrelic/newrelic-cli/internal/install/execution"
	"github.com/newrelic/newrelic-cli/internal/install/types"
	"github.com/stretchr/testify/require"
)

var (
	recipe_detector_ProcessEvaluator = &mockRecipeEvaluator{}
	recipe_detector_ScriptEvaluator  = &mockRecipeEvaluator{}
	recipe_detector                  = *newRecipeDetector(recipe_detector_ProcessEvaluator, recipe_detector_ScriptEvaluator)
)

func TestRecipeDetector_ShouldGetProcessEvalStatusNull(t *testing.T) {
	recipe := givenRecipeWithNoProcessMatching()
	withProcessEvaluatorReturnStatus(execution.RecipeStatusTypes.NULL)
	withScriptEvaluatorReturnStatus(execution.RecipeStatusTypes.AVAILABLE)

	statusDetectionResult := recipe_detector.detectRecipe(ctx, recipe)
	actual := statusDetectionResult
	require.Equal(t, execution.RecipeStatusTypes.NULL, actual)
}

func TestRecipeDetector_ShouldGetProcessEvalStatusAvaliable(t *testing.T) {
	recipe := givenRecipeWithNoProcessMatching()
	withEmptyPreInstallRequiredAtDiscoverSection()
	withProcessEvaluatorReturnStatus(execution.RecipeStatusTypes.AVAILABLE)
	withScriptEvaluatorReturnStatus(execution.RecipeStatusTypes.DETECTED)

	statusDetectionResult := recipe_detector.detectRecipe(ctx, recipe)
	actual := statusDetectionResult
	require.Equal(t, execution.RecipeStatusTypes.AVAILABLE, actual)
}

func TestRecipeDetector_ShouldGetScriptEvalStatusNull(t *testing.T) {
	recipe := withPreInstallRequiredAtDiscoverSection()
	withPreInstallRequiredAtDiscoverSection()
	withProcessEvaluatorReturnStatus(execution.RecipeStatusTypes.AVAILABLE)
	withScriptEvaluatorReturnStatus(execution.RecipeStatusTypes.NULL)

	statusDetectionResult := recipe_detector.detectRecipe(ctx, recipe)
	actual := statusDetectionResult
	require.Equal(t, execution.RecipeStatusTypes.NULL, actual)
}

func TestRecipeDetector_ShouldGetScriptEvalStatusDetected(t *testing.T) {
	recipe := withPreInstallRequiredAtDiscoverSection()
	withProcessEvaluatorReturnStatus(execution.RecipeStatusTypes.AVAILABLE)
	withScriptEvaluatorReturnStatus(execution.RecipeStatusTypes.DETECTED)

	statusDetectionResult := recipe_detector.detectRecipe(ctx, recipe)
	actual := statusDetectionResult
	require.Equal(t, execution.RecipeStatusTypes.DETECTED, actual)
}
func TestRecipeDetector_ShouldGetScriptEvalStatusAvailable(t *testing.T) {
	recipe := withPreInstallRequiredAtDiscoverSection()
	withProcessEvaluatorReturnStatus(execution.RecipeStatusTypes.AVAILABLE)
	withScriptEvaluatorReturnStatus(execution.RecipeStatusTypes.AVAILABLE)

	statusDetectionResult := recipe_detector.detectRecipe(ctx, recipe)
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
	recipe_detector_ProcessEvaluator.status = status
}

func withScriptEvaluatorReturnStatus(status execution.RecipeStatusType) {
	recipe_detector_ScriptEvaluator.status = status
}

type mockRecipeEvaluator struct {
	status execution.RecipeStatusType
}

func (mre *mockRecipeEvaluator) DetectionStatus(ctx context.Context, recipe *types.OpenInstallationRecipe) execution.RecipeStatusType {
	return mre.status
}
