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

func (se *ScriptEvaluator) DetectionStatus(ctx context.Context, r *types.OpenInstallationRecipe) (statusResult execution.RecipeStatusType) {

	defer func() {
		if err := recover(); err != nil {
			log.Debugf("recipe %s failed script evaluation with panic %s", r.Name, err)
			statusResult = execution.RecipeStatusTypes.NULL
		}
	}()

	if err := se.executor.ExecutePreInstall(ctx, *r, types.RecipeVars{}); err != nil {
		log.Debugf("recipe %s failed script evaluation %s", r.Name, err)

		if utils.IsExitStatusCode(132, err) {
			statusResult = execution.RecipeStatusTypes.DETECTED
		} else if utils.IsExitStatusCode(131, err) {
			statusResult = execution.RecipeStatusTypes.UNSUPPORTED
		} else {
			statusResult = execution.RecipeStatusTypes.NULL
		}
		return
	}

	statusResult = execution.RecipeStatusTypes.AVAILABLE
	return
}
