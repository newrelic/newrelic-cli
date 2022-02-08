package recipes

import (
	"context"

	"github.com/newrelic/newrelic-cli/internal/install/execution"
	"github.com/newrelic/newrelic-cli/internal/install/types"
	"github.com/newrelic/newrelic-cli/internal/utils"
	log "github.com/sirupsen/logrus"
)

type ScriptEvaluator struct {
	executor execution.RecipeExecutor
}

func NewScriptEvaluator() *ScriptEvaluator {
	executor := execution.NewShRecipeExecutor()

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

		return ""
	}

	return execution.RecipeStatusTypes.AVAILABLE
}
