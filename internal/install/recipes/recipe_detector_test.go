package recipes

import (
	"context"
	"testing"

	"github.com/newrelic/newrelic-cli/internal/install/execution"
	"github.com/newrelic/newrelic-cli/internal/install/types"
	"github.com/stretchr/testify/require"
)

func TestDetectorShouldReturnProcessEvaluatorStatus(t *testing.T) {

	tests := []struct {
		name     string
		recipe   *types.OpenInstallationRecipe
		detector RecipeDetector
		expected execution.RecipeStatusType
	}{
		{
			"Null Status from process detector should given recipe null status",
			createRecipe("0", "recipe1"),
			*newRecipeDetector(&mockRecipeEvaluator{execution.RecipeStatusTypes.NULL}, &mockRecipeEvaluator{execution.RecipeStatusTypes.AVAILABLE}),
			execution.RecipeStatusTypes.NULL,
		},
		{
			"Avilabe Status from process detector and null require discovery script should given recipe available status",
			createRecipe("0", "recipe1"),
			*newRecipeDetector(&mockRecipeEvaluator{execution.RecipeStatusTypes.AVAILABLE}, &mockRecipeEvaluator{execution.RecipeStatusTypes.DETECTED}),
			execution.RecipeStatusTypes.AVAILABLE,
		},
	}

	for _, d := range tests {
		t.Run(d.name, func(t *testing.T) {
			sut := d.detector
			var ctx context.Context
			actual := sut.DetectRecipes(ctx, []*types.OpenInstallationRecipe{d.recipe})
			actualRecipeStatus := actual[d.recipe]
			require.Equal(t, len(actual), 1)
			require.Equal(t, d.expected, actualRecipeStatus)
		})
	}
}

func TestDetectorShouldReturnScriptEvaluatorStatus(t *testing.T) {

	tests := []struct {
		name     string
		recipe   *types.OpenInstallationRecipe
		detector RecipeDetector
		expected execution.RecipeStatusType
	}{
		{
			"Null Status from process detector should given recipe null status",
			createRecipeWithPreInstall("0", "recipe1"),
			*newRecipeDetector(&mockRecipeEvaluator{execution.RecipeStatusTypes.AVAILABLE}, &mockRecipeEvaluator{execution.RecipeStatusTypes.NULL}),
			execution.RecipeStatusTypes.NULL,
		},
		{
			"Avilabe Status from process detector and null require discovery script should given recipe available status",
			createRecipeWithPreInstall("0", "recipe1"),
			*newRecipeDetector(&mockRecipeEvaluator{execution.RecipeStatusTypes.AVAILABLE}, &mockRecipeEvaluator{execution.RecipeStatusTypes.AVAILABLE}),
			execution.RecipeStatusTypes.AVAILABLE,
		},
		{
			"Avilabe Status from process detector and detected from script detector should given recipe detected status",
			createRecipeWithPreInstall("0", "recipe1"),
			*newRecipeDetector(&mockRecipeEvaluator{execution.RecipeStatusTypes.AVAILABLE}, &mockRecipeEvaluator{execution.RecipeStatusTypes.DETECTED}),
			execution.RecipeStatusTypes.DETECTED,
		},
	}

	for _, d := range tests {
		t.Run(d.name, func(t *testing.T) {
			sut := d.detector
			var ctx context.Context
			actual := sut.DetectRecipes(ctx, []*types.OpenInstallationRecipe{d.recipe})
			actualRecipeStatus := actual[d.recipe]
			require.Equal(t, len(actual), 1)
			require.Equal(t, d.expected, actualRecipeStatus)
		})
	}
}

type mockRecipeEvaluator struct {
	status execution.RecipeStatusType
}

func (mre *mockRecipeEvaluator) DetectionStatus(ctx context.Context, recipe *types.OpenInstallationRecipe) execution.RecipeStatusType {
	return mre.status
}

func createRecipeWithPreInstall(id string, name string) *types.OpenInstallationRecipe {
	r := &types.OpenInstallationRecipe{
		ID:   id,
		Name: name,
	}
	r.PreInstall = types.OpenInstallationPreInstallConfiguration{
		RequireAtDiscovery: "pre-install script mock",
	}
	return r
}
