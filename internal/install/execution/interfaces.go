package execution

import (
	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type StatusReporter interface {
	RecipeDetected(status *InstallStatus, event RecipeStatusEvent) error
	RecipeFailed(status *InstallStatus, event RecipeStatusEvent) error
	RecipeInstalling(status *InstallStatus, event RecipeStatusEvent) error
	RecipeInstalled(status *InstallStatus, event RecipeStatusEvent) error
	RecipeSkipped(status *InstallStatus, event RecipeStatusEvent) error
	RecipeRecommended(status *InstallStatus, event RecipeStatusEvent) error
	RecipesSelected(status *InstallStatus, recipes []types.OpenInstallationRecipe) error
	RecipeAvailable(status *InstallStatus, event RecipeStatusEvent) error
	RecipeCanceled(status *InstallStatus, event RecipeStatusEvent) error
	InstallStarted(status *InstallStatus) error
	InstallComplete(status *InstallStatus) error
	InstallCanceled(status *InstallStatus) error
	DiscoveryComplete(status *InstallStatus, dm types.DiscoveryManifest) error
	RecipeUnsupported(status *InstallStatus, event RecipeStatusEvent) error
	UpdateRequired(status *InstallStatus) error
}
