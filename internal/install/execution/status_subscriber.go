package execution

import "github.com/newrelic/newrelic-cli/internal/install/types"

// StatusSubscriber is notified during the lifecycle of the recipe execution status.
type StatusSubscriber interface {
	InstallCanceled(status *InstallStatus) error
	InstallComplete(status *InstallStatus) error
	DiscoveryComplete(status *InstallStatus, dm types.DiscoveryManifest) error
	RecipeAvailable(status *InstallStatus, recipe types.Recipe) error
	RecipeFailed(status *InstallStatus, event RecipeStatusEvent) error
	RecipeInstalled(status *InstallStatus, event RecipeStatusEvent) error
	RecipeInstalling(status *InstallStatus, event RecipeStatusEvent) error
	RecipeRecommended(status *InstallStatus, event RecipeStatusEvent) error
	RecipeSkipped(status *InstallStatus, event RecipeStatusEvent) error
	RecipesAvailable(status *InstallStatus, recipes []types.Recipe) error
	RecipesSelected(status *InstallStatus, recipes []types.Recipe) error
}

// RecipeStatusEvent represents an event in a recipe's execution.
type RecipeStatusEvent struct {
	Recipe     types.Recipe
	Msg        string
	EntityGUID string
}
