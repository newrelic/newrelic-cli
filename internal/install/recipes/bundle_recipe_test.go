package recipes

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-cli/internal/install/execution"
)

func TestBundleRecipeAddsStatusWithTime(t *testing.T) {
	var durationMs int64 = 48
	expectedStatus := execution.RecipeStatusTypes.DETECTED
	br := testBundleRecipe()

	br.AddDetectionStatus(expectedStatus, durationMs)

	require.Equal(t, len(br.DetectedStatuses), 1)
	require.Equal(t, expectedStatus, br.DetectedStatuses[0].Status)
	require.Equal(t, durationMs, br.DetectedStatuses[0].DurationMs)
}

func TestBundleRecipeShouldAddStatusOnceAtFirstOccurrence(t *testing.T) {
	br := testBundleRecipe()
	expectedStatus := execution.RecipeStatusTypes.INSTALLING

	br.AddDetectionStatus(expectedStatus, 0)
	br.AddDetectionStatus(expectedStatus, 0)
	br.AddDetectionStatus(expectedStatus, 0)

	require.Equal(t, len(br.DetectedStatuses), 1)
	require.Equal(t, expectedStatus, br.DetectedStatuses[0].Status)
}

func TestBundleRecipeShouldAddStatusDetectedWhenAvailable(t *testing.T) {
	var durationMs int64 = 67
	br := testBundleRecipe()

	br.AddDetectionStatus(execution.RecipeStatusTypes.AVAILABLE, durationMs)

	require.Equal(t, len(br.DetectedStatuses), 2)
	require.Equal(t, execution.RecipeStatusTypes.DETECTED, br.DetectedStatuses[0].Status)
	require.Equal(t, durationMs, br.DetectedStatuses[1].DurationMs)
	require.Equal(t, execution.RecipeStatusTypes.AVAILABLE, br.DetectedStatuses[1].Status)
	require.Equal(t, durationMs, br.DetectedStatuses[1].DurationMs)

}

func TestBundleRecipeHasStatusReturnsTrue(t *testing.T) {
	br := testBundleRecipeWithStatus(execution.RecipeStatusTypes.DETECTED, time.Now())

	require.True(t, br.HasStatus(execution.RecipeStatusTypes.DETECTED))
}

func TestBundleRecipeHasStatusReturnsFalse(t *testing.T) {
	br := testBundleRecipeWithStatus(execution.RecipeStatusTypes.AVAILABLE, time.Now())

	require.False(t, br.HasStatus(execution.RecipeStatusTypes.DETECTED))
}

func TestBundleRecipeFlattenShouldReturnOneDistinct(t *testing.T) {
	br := testBundleRecipe()
	dependentRecipeWithSameName := testBundleRecipe()
	br.Dependencies = append(br.Dependencies, dependentRecipeWithSameName)

	actual := br.Flatten()
	require.Equal(t, 1, len(actual))
	require.Equal(t, actual[br.Recipe.Name], true)
}

func TestBundleRecipeFlattenShouldReturnTwoNonDistinct(t *testing.T) {
	br := testBundleRecipe()
	dependentRecipeWithDifferentName := testBundleRecipe()
	dependentRecipeWithDifferentName.Recipe.Name = "Fake Recipe"
	br.Dependencies = append(br.Dependencies, dependentRecipeWithDifferentName)

	actual := br.Flatten()
	require.Equal(t, 2, len(actual))
	require.Equal(t, actual[br.Recipe.Name], true)
	require.Equal(t, actual["Fake Recipe"], true)
}
func TestBundleRecipeFlattenMultiLevelShouldReturnTwoNonDistinct(t *testing.T) {
	br := testBundleRecipe()
	dependentRecipeWithDifferentName := testBundleRecipe()
	dependentRecipeWithDifferentName.Recipe.Name = "Fake Recipe"
	br.Dependencies = append(br.Dependencies, dependentRecipeWithDifferentName)
	layeredDependentRecipeWithDifferentName := testBundleRecipe()
	layeredDependentRecipeWithDifferentName.Recipe.Name = "Layered Fake Recipe"
	br.Dependencies[0].Dependencies = append(br.Dependencies[0].Dependencies, layeredDependentRecipeWithDifferentName)

	actual := br.Flatten()
	require.Equal(t, 3, len(actual))
	require.Equal(t, actual[br.Recipe.Name], true)
	require.Equal(t, actual["Fake Recipe"], true)
	require.Equal(t, actual["Layered Fake Recipe"], true)
}

func testBundleRecipe() *BundleRecipe {
	return &BundleRecipe{
		Recipe: NewRecipeBuilder().Build(),
	}
}

func testBundleRecipeWithStatus(status execution.RecipeStatusType, statusTime time.Time) *BundleRecipe {
	bundleRecipe := testBundleRecipe()
	bundleRecipe.DetectedStatuses = []*DetectedStatusType{{Status: status}}
	return bundleRecipe
}
