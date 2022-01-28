package recipes

import (
	"context"

	"github.com/newrelic/newrelic-cli/internal/install/execution"
	"github.com/newrelic/newrelic-cli/internal/install/types"
	"github.com/newrelic/newrelic-cli/internal/utils"
)

type ScriptEvaluationRecipeFilterer struct {
	recipeExecutor execution.RecipeExecutor
	installStatus  *execution.InstallStatus
}

func NewScriptEvaluationRecipeFilterer(installStatus *execution.InstallStatus) *ScriptEvaluationRecipeFilterer {
	recipeExecutor := execution.NewShRecipeExecutor()

	return &ScriptEvaluationRecipeFilterer{
		recipeExecutor: recipeExecutor,
		installStatus:  installStatus,
	}
}

func (f *ScriptEvaluationRecipeFilterer) CheckCompatibility(ctx context.Context, r *types.OpenInstallationRecipe, m *types.DiscoveryManifest) error {
	err := f.recipeExecutor.ExecutePreInstall(ctx, *r, types.RecipeVars{})

	if err != nil {
		var metadata map[string]interface{}
		var message string
		if e, ok := err.(*types.CustomStdError); ok {
			metadata = e.Metadata
		} else {
			message = err.Error()
		}

		event := execution.RecipeStatusEvent{
			Recipe:   *r,
			Msg:      message,
			Metadata: metadata,
		}

		if utils.IsExitStatusCode(132, err) {
			f.installStatus.RecipeDetected(*r, event)
		} else {
			f.installStatus.RecipeUnsupported(event)
		}
	}

	return err
}
