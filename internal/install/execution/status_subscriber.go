package execution

import "github.com/newrelic/newrelic-cli/internal/install/types"

// StatusSubscriber is notified during the lifecycle of the recipe execution status.
type StatusSubscriber interface {
	UpdateRequired(status *InstallStatus) error
	InstallStarted(status *InstallStatus) error
	InstallCanceled(status *InstallStatus) error
	InstallComplete(status *InstallStatus) error
	DiscoveryComplete(status *InstallStatus, dm types.DiscoveryManifest) error
	RecipeDetected(status *InstallStatus, event RecipeStatusEvent) error
	RecipeCanceled(status *InstallStatus, event RecipeStatusEvent) error
	RecipeAvailable(status *InstallStatus, event RecipeStatusEvent) error
	RecipeFailed(status *InstallStatus, event RecipeStatusEvent) error
	RecipeInstalled(status *InstallStatus, event RecipeStatusEvent) error
	RecipeInstalling(status *InstallStatus, event RecipeStatusEvent) error
	RecipeRecommended(status *InstallStatus, event RecipeStatusEvent) error
	RecipeSkipped(status *InstallStatus, event RecipeStatusEvent) error
	RecipeUnsupported(status *InstallStatus, event RecipeStatusEvent) error
	RecipesSelected(status *InstallStatus, recipes []types.OpenInstallationRecipe) error
}

// RecipeStatusEvent represents an event in a recipe's execution.
type RecipeStatusEvent struct {
	Recipe               types.OpenInstallationRecipe
	Msg                  string
	TaskPath             []string
	EntityGUID           string
	ValidationDurationMs int64
}

func NewRecipeStatusEvent(recipe *types.OpenInstallationRecipe) RecipeStatusEvent {
	return RecipeStatusEvent{Recipe: *recipe}
}
