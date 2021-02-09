package execution

import "github.com/newrelic/newrelic-cli/internal/install/types"

// StatusReporter is responsible for reporting the status of recipe execution.
type StatusReporter interface {
	ReportComplete(status *StatusRollup) error
	ReportRecipeAvailable(status *StatusRollup, recipe types.Recipe) error
	ReportRecipeFailed(status *StatusRollup, event RecipeStatusEvent) error
	ReportRecipeInstalled(status *StatusRollup, event RecipeStatusEvent) error
	ReportRecipeInstalling(status *StatusRollup, event RecipeStatusEvent) error
	ReportRecipeRecommended(status *StatusRollup, event RecipeStatusEvent) error
	ReportRecipeSkipped(status *StatusRollup, event RecipeStatusEvent) error
	ReportRecipesAvailable(status *StatusRollup, recipes []types.Recipe) error
}

// RecipeStatusEvent represents an event in a recipe's execution.
type RecipeStatusEvent struct {
	Recipe     types.Recipe
	Msg        string
	EntityGUID string
}
