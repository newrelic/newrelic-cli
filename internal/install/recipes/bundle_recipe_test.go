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

func TestAllDependenciesAvailable_Happy(t *testing.T) {
	recipe := NewRecipeBuilder().DependencyBuilder(NewRecipeBuilder().Name("dep1")).DependencyBuilder(NewRecipeBuilder().Name("dep2")).BuildBundleRecipe()
	recipe.Dependencies[0].AddDetectionStatus(execution.RecipeStatusTypes.AVAILABLE, 0)
	recipe.Dependencies[1].AddDetectionStatus(execution.RecipeStatusTypes.AVAILABLE, 0)

	require.True(t, recipe.AreAllDependenciesAvailable())
}

func TestAllDependenciesAvailable_OneDepNotAvailable(t *testing.T) {
	recipe := NewRecipeBuilder().DependencyBuilder(NewRecipeBuilder().Name("dep1")).DependencyBuilder(NewRecipeBuilder().Name("dep2")).BuildBundleRecipe()
	recipe.Dependencies[0].AddDetectionStatus(execution.RecipeStatusTypes.AVAILABLE, 0)

	require.False(t, recipe.AreAllDependenciesAvailable())
}

func TestAllDependenciesAvailable_AllDepNotAvailable(t *testing.T) {
	recipe := NewRecipeBuilder().DependencyBuilder(NewRecipeBuilder().Name("dep1")).DependencyBuilder(NewRecipeBuilder().Name("dep2")).BuildBundleRecipe()
	recipe.Dependencies[0].AddDetectionStatus(execution.RecipeStatusTypes.UNSUPPORTED, 0)
	recipe.Dependencies[1].AddDetectionStatus(execution.RecipeStatusTypes.DETECTED, 0)

	require.False(t, recipe.AreAllDependenciesAvailable())
}

func TestAllDependenciesAvailable_NoDep(t *testing.T) {
	recipe := testBundleRecipe()

	require.True(t, recipe.AreAllDependenciesAvailable())
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
