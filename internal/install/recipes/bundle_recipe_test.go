package recipes

import (
	"testing"
	"time"

	"github.com/newrelic/newrelic-cli/internal/install/execution"
	"github.com/stretchr/testify/require"
)

func TestBundleRecipeAddsStatusWithTime(t *testing.T) {
	expectedStatusTime := time.Now()
	expectedStatus := execution.RecipeStatusTypes.INSTALLING
	br := testBundleRecipe()

	br.AddStatus(expectedStatus, expectedStatusTime)

	require.Equal(t, len(br.RecipeStatuses), 1)
	require.Equal(t, expectedStatus, br.RecipeStatuses[0].status)
	require.Equal(t, expectedStatusTime, br.RecipeStatuses[0].statusTime)

}

func TestBundleRecipeShouldAddStatusOnceAtFirstOccurance(t *testing.T) {
	br := testBundleRecipe()
	earliestTime := time.Now()
	moreRecentTime := earliestTime.Add(5 * time.Minute)
	mostRecentTime := earliestTime.Add(10 * time.Minute)
	expectedStatus := execution.RecipeStatusTypes.INSTALLING

	br.AddStatus(expectedStatus, earliestTime)
	br.AddStatus(expectedStatus, moreRecentTime)
	br.AddStatus(expectedStatus, mostRecentTime)

	require.Equal(t, len(br.RecipeStatuses), 1)
	require.Equal(t, expectedStatus, br.RecipeStatuses[0].status)
	require.Equal(t, earliestTime, br.RecipeStatuses[0].statusTime)
}

func TestBundleRecipeShouldAddStatusDetectedWhenAvailable(t *testing.T) {
	br := testBundleRecipe()
	expectedStatusTime := time.Now()

	br.AddStatus(execution.RecipeStatusTypes.AVAILABLE, expectedStatusTime)

	require.Equal(t, len(br.RecipeStatuses), 2)
	require.Equal(t, execution.RecipeStatusTypes.DETECTED, br.RecipeStatuses[0].status)
	require.Equal(t, expectedStatusTime, br.RecipeStatuses[0].statusTime)
	require.Equal(t, execution.RecipeStatusTypes.AVAILABLE, br.RecipeStatuses[1].status)
	require.Equal(t, expectedStatusTime, br.RecipeStatuses[1].statusTime)
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
		Recipe: createRecipe("id1", "recipe2"),
	}
}

func testBundleRecipeWithStatus(status execution.RecipeStatusType, statusTime time.Time) *BundleRecipe {
	bundleRecipe := testBundleRecipe()
	bundleRecipe.RecipeStatuses = []recipeStatus{{
		status:     status,
		statusTime: statusTime,
	},
	}
	return bundleRecipe
}
