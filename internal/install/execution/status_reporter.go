package execution

import "github.com/newrelic/newrelic-cli/internal/install/types"

// StatusReporter is responsible for reporting the status of recipe execution.
type StatusReporter interface {
	ReportRecipeFailed(event RecipeStatusEvent) error
	ReportRecipeInstalled(event RecipeStatusEvent) error
	ReportRecipeSkipped(event RecipeStatusEvent) error
	ReportRecipesAvailable(recipes []types.Recipe) error
	ReportComplete() error
}

// RecipeStatusEvent represents an event in a recipe's execution.
type RecipeStatusEvent struct {
	Recipe     types.Recipe
	Msg        string
	EntityGUID string
}
