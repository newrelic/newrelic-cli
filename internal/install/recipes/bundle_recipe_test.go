package recipes

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-cli/internal/install/execution"
)

func TestBundleRecipeAddsStatusWithTime(t *testing.T) {
	expectedStatus := execution.RecipeStatusTypes.INSTALLING
	br := testBundleRecipe()

	br.AddDetectionStatus(expectedStatus)

	require.Equal(t, len(br.DetectedStatuses), 1)
	require.Equal(t, expectedStatus, br.DetectedStatuses[0])
}

func TestBundleRecipeShouldAddStatusOnceAtFirstOccurrence(t *testing.T) {
	br := testBundleRecipe()
	expectedStatus := execution.RecipeStatusTypes.INSTALLING

	br.AddDetectionStatus(expectedStatus)
	br.AddDetectionStatus(expectedStatus)
	br.AddDetectionStatus(expectedStatus)

	require.Equal(t, len(br.DetectedStatuses), 1)
	require.Equal(t, expectedStatus, br.DetectedStatuses[0])
}

func TestBundleRecipeShouldAddStatusDetectedWhenAvailable(t *testing.T) {
	br := testBundleRecipe()

	br.AddDetectionStatus(execution.RecipeStatusTypes.AVAILABLE)

	require.Equal(t, len(br.DetectedStatuses), 2)
	require.Equal(t, execution.RecipeStatusTypes.DETECTED, br.DetectedStatuses[0])
	require.Equal(t, execution.RecipeStatusTypes.AVAILABLE, br.DetectedStatuses[1])

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
	bundleRecipe.DetectedStatuses = []execution.RecipeStatusType{
		status,
	}
	return bundleRecipe
}
