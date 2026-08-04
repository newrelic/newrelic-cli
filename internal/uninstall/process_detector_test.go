package uninstall

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-cli/internal/install/execution"
	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type fakeDetectionStatusEvaluator struct {
	status execution.RecipeStatusType
}

func (f *fakeDetectionStatusEvaluator) DetectionStatus(ctx context.Context, r *types.OpenInstallationRecipe, recipeNames []string) execution.RecipeStatusType {
	return f.status
}

func TestEvaluatorProcessDetector_AvailableMeansRunning(t *testing.T) {
	d := NewEvaluatorProcessDetector(&fakeDetectionStatusEvaluator{status: execution.RecipeStatusTypes.AVAILABLE})

	require.True(t, d.IsRunning(context.Background(), types.OpenInstallationRecipe{}))
}

func TestEvaluatorProcessDetector_NullMeansNotRunning(t *testing.T) {
	d := NewEvaluatorProcessDetector(&fakeDetectionStatusEvaluator{status: execution.RecipeStatusTypes.NULL})

	require.False(t, d.IsRunning(context.Background(), types.OpenInstallationRecipe{}))
}
