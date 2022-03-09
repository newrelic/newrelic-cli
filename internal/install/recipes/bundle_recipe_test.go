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
