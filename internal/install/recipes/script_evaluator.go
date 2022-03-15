package recipes

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/internal/install/execution"
	"github.com/newrelic/newrelic-cli/internal/install/types"
	"github.com/newrelic/newrelic-cli/internal/utils"
)

type ScriptEvaluator struct {
	executor execution.RecipeExecutor
}

func NewScriptEvaluator() *ScriptEvaluator {
	return newScriptEvaluator(execution.NewShRecipeExecutor())
}

func newScriptEvaluator(executor execution.RecipeExecutor) *ScriptEvaluator {
	return &ScriptEvaluator{
		executor: executor,
	}
}

func (se *ScriptEvaluator) DetectionStatus(ctx context.Context, r *types.OpenInstallationRecipe) execution.RecipeStatusType {
	if err := se.executor.ExecutePreInstall(ctx, *r, types.RecipeVars{}); err != nil {
		log.Tracef("recipe %s failed script evaluation %s", r.Name, err)

		if utils.IsExitStatusCode(132, err) {
			return execution.RecipeStatusTypes.DETECTED
		}

		if utils.IsExitStatusCode(131, err) {
			return execution.RecipeStatusTypes.UNSUPPORTED
		}

		return execution.RecipeStatusTypes.NULL
	}

	return execution.RecipeStatusTypes.AVAILABLE
}
