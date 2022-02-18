package recipes

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-cli/internal/install/execution"
)

func TestBundleRecipeAddsStatusWithTime(t *testing.T) {
	expectedStatusTime := time.Now()
	expectedStatus := execution.RecipeStatusTypes.INSTALLING
	br := testBundleRecipe()

	br.AddStatus(expectedStatus, expectedStatusTime)

	require.Equal(t, len(br.RecipeStatuses), 1)
	require.Equal(t, expectedStatus, br.RecipeStatuses[0].Status)
	require.Equal(t, expectedStatusTime, br.RecipeStatuses[0].StatusTime)

}

func TestBundleRecipeShouldAddStatusOnceAtFirstOccurrence(t *testing.T) {
	br := testBundleRecipe()
	earliestTime := time.Now()
	moreRecentTime := earliestTime.Add(5 * time.Minute)
	mostRecentTime := earliestTime.Add(10 * time.Minute)
	expectedStatus := execution.RecipeStatusTypes.INSTALLING

	br.AddStatus(expectedStatus, earliestTime)
	br.AddStatus(expectedStatus, moreRecentTime)
	br.AddStatus(expectedStatus, mostRecentTime)

	require.Equal(t, len(br.RecipeStatuses), 1)
	require.Equal(t, expectedStatus, br.RecipeStatuses[0].Status)
	require.Equal(t, earliestTime, br.RecipeStatuses[0].StatusTime)
}

func TestBundleRecipeShouldAddStatusDetectedWhenAvailable(t *testing.T) {
	br := testBundleRecipe()
	expectedStatusTime := time.Now()

	br.AddStatus(execution.RecipeStatusTypes.AVAILABLE, expectedStatusTime)

	require.Equal(t, len(br.RecipeStatuses), 2)
	require.Equal(t, execution.RecipeStatusTypes.DETECTED, br.RecipeStatuses[0].Status)
	require.Equal(t, expectedStatusTime, br.RecipeStatuses[0].StatusTime)
	require.Equal(t, execution.RecipeStatusTypes.AVAILABLE, br.RecipeStatuses[1].Status)
	require.Equal(t, expectedStatusTime, br.RecipeStatuses[1].StatusTime)
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
	bundleRecipe.RecipeStatuses = []RecipeStatus{{
		Status:     status,
		StatusTime: statusTime,
	},
	}
	return bundleRecipe
}
