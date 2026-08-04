//go:build unit

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
	d := NewEvaluatorProcessDetector(
		&fakeDetectionStatusEvaluator{status: execution.RecipeStatusTypes.AVAILABLE},
		&fakeDetectionStatusEvaluator{status: execution.RecipeStatusTypes.AVAILABLE},
	)

	r := types.OpenInstallationRecipe{ProcessMatch: []string{"newrelic-infra"}}

	require.True(t, d.IsRunning(context.Background(), r))
}

func TestEvaluatorProcessDetector_NullMeansNotRunning(t *testing.T) {
	d := NewEvaluatorProcessDetector(
		&fakeDetectionStatusEvaluator{status: execution.RecipeStatusTypes.NULL},
		&fakeDetectionStatusEvaluator{status: execution.RecipeStatusTypes.NULL},
	)

	r := types.OpenInstallationRecipe{ProcessMatch: []string{"newrelic-infra"}}

	require.False(t, d.IsRunning(context.Background(), r))
}

func TestEvaluatorProcessDetector_NoProcessMatchMeansNotRunning(t *testing.T) {
	d := NewEvaluatorProcessDetector(
		&fakeDetectionStatusEvaluator{status: execution.RecipeStatusTypes.AVAILABLE},
		&fakeDetectionStatusEvaluator{status: execution.RecipeStatusTypes.AVAILABLE},
	)

	r := types.OpenInstallationRecipe{}

	require.False(t, d.IsRunning(context.Background(), r), "a recipe with no processMatch has no signal to check, so it should not be reported as running")
}

func TestEvaluatorProcessDetector_RequireAtDiscoveryCanDowngradeAvailable(t *testing.T) {
	d := NewEvaluatorProcessDetector(
		&fakeDetectionStatusEvaluator{status: execution.RecipeStatusTypes.AVAILABLE},
		&fakeDetectionStatusEvaluator{status: execution.RecipeStatusTypes.NULL},
	)

	r := types.OpenInstallationRecipe{
		ProcessMatch: []string{"newrelic-infra"},
		PreInstall:   types.OpenInstallationPreInstallConfiguration{RequireAtDiscovery: "some-script"},
	}

	require.False(t, d.IsRunning(context.Background(), r), "when RequireAtDiscovery is set, the script evaluator's result should win")
}

func TestEvaluatorProcessDetector_RequireAtDiscoveryNotConsultedWhenProcessNotAvailable(t *testing.T) {
	scriptEvaluator := &fakeDetectionStatusEvaluator{status: execution.RecipeStatusTypes.AVAILABLE}
	d := NewEvaluatorProcessDetector(
		&fakeDetectionStatusEvaluator{status: execution.RecipeStatusTypes.NULL},
		scriptEvaluator,
	)

	r := types.OpenInstallationRecipe{
		ProcessMatch: []string{"newrelic-infra"},
		PreInstall:   types.OpenInstallationPreInstallConfiguration{RequireAtDiscovery: "some-script"},
	}

	require.False(t, d.IsRunning(context.Background(), r))
}
