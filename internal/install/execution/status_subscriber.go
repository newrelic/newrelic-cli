package execution

import "github.com/newrelic/newrelic-cli/internal/install/types"

// StatusSubscriber is notified during the lifecycle of the recipe execution status.
type StatusSubscriber interface {
	UpdateRequired(status *InstallStatus) error
	InstallStarted(status *InstallStatus) error
	InstallCanceled(status *InstallStatus) error
	InstallComplete(status *InstallStatus) error
	DiscoveryComplete(status *InstallStatus, dm types.DiscoveryManifest) error
	RecipeAvailable(status *InstallStatus, recipe types.OpenInstallationRecipe) error
	RecipeFailed(status *InstallStatus, event RecipeStatusEvent) error
	RecipeInstalled(status *InstallStatus, event RecipeStatusEvent) error
	RecipeInstalling(status *InstallStatus, event RecipeStatusEvent) error
	RecipeRecommended(status *InstallStatus, event RecipeStatusEvent) error
	RecipeSkipped(status *InstallStatus, event RecipeStatusEvent) error
	RecipeUnsupported(status *InstallStatus, event RecipeStatusEvent) error
	RecipesSelected(status *InstallStatus, recipes []types.OpenInstallationRecipe) error
	ObservabilityPackFetchPending(status *InstallStatus) error
	ObservabilityPackFetchSuccess(status *InstallStatus) error
	ObservabilityPackFetchFailed(status *InstallStatus) error
	ObservabilityPackInstallPending(status *InstallStatus) error
	ObservabilityPackInstallSuccess(status *InstallStatus) error
	ObservabilityPackInstallFailed(status *InstallStatus) error
}

// RecipeStatusEvent represents an event in a recipe's execution.
type RecipeStatusEvent struct {
	Recipe               types.OpenInstallationRecipe
	Msg                  string
	TaskPath             []string
	EntityGUID           string
	ValidationDurationMs int64
}

type ObservabilityPackStatusEvent struct {
	ObservabilityPack types.OpenInstallationObservabilityPack
	Name              string
	Msg               string
}
