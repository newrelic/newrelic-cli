package execution

import (
	configAPI "github.com/newrelic/newrelic-cli/internal/config/api"
	"github.com/newrelic/newrelic-cli/internal/install/types"
	"github.com/newrelic/newrelic-client-go/pkg/entities"
	"github.com/newrelic/newrelic-client-go/pkg/installevents"
)

type InstallEventsReporter struct {
	client    InstallEventsClient
	accountID int
}

func NewInstallEventsReporter(client InstallEventsClient) *InstallEventsReporter {
	r := InstallEventsReporter{
		client: client,
	}

	r.accountID = configAPI.GetActiveProfileAccountID()

	return &r
}

func (r InstallEventsReporter) RecipeAvailable(status *InstallStatus, recipe types.OpenInstallationRecipe) error {
	event := RecipeStatusEvent{Recipe: recipe}
	err := r.createRecipeInstallEvent(status, installevents.InstallationRecipeStatusTypeTypes.AVAILABLE, event)
	return err
}

func (r InstallEventsReporter) RecipeCanceled(status *InstallStatus, event RecipeStatusEvent) error {
	err := r.createRecipeInstallEvent(status, installevents.InstallationRecipeStatusTypeTypes.CANCELED, event)
	return err
}

func (r InstallEventsReporter) RecipeFailed(status *InstallStatus, event RecipeStatusEvent) error {
	err := r.createRecipeInstallEvent(status, installevents.InstallationRecipeStatusTypeTypes.FAILED, event)
	return err
}

func (r InstallEventsReporter) RecipeInstalling(status *InstallStatus, event RecipeStatusEvent) error {
	err := r.createRecipeInstallEvent(status, installevents.InstallationRecipeStatusTypeTypes.INSTALLING, event)
	return err
}

func (r InstallEventsReporter) RecipeInstalled(status *InstallStatus, event RecipeStatusEvent) error {
	err := r.createRecipeInstallEvent(status, installevents.InstallationRecipeStatusTypeTypes.INSTALLED, event)
	return err
}

func (r InstallEventsReporter) RecipeSkipped(status *InstallStatus, event RecipeStatusEvent) error {
	err := r.createRecipeInstallEvent(status, installevents.InstallationRecipeStatusTypeTypes.SKIPPED, event)
	return err
}

func (r InstallEventsReporter) RecipeRecommended(status *InstallStatus, event RecipeStatusEvent) error {
	err := r.createRecipeInstallEvent(status, installevents.InstallationRecipeStatusTypeTypes.RECOMMENDED, event)
	return err
}

func (r InstallEventsReporter) RecipeUnsupported(status *InstallStatus, event RecipeStatusEvent) error {
	err := r.createRecipeInstallEvent(status, installevents.InstallationRecipeStatusTypeTypes.UNSUPPORTED, event)
	return err
}

func (r InstallEventsReporter) RecipesSelected(status *InstallStatus, recipes []types.OpenInstallationRecipe) error {
	return nil
}

func (r InstallEventsReporter) InstallComplete(status *InstallStatus) error {
	return nil
}

func (r InstallEventsReporter) InstallCanceled(status *InstallStatus) error {
	return nil
}

func (r InstallEventsReporter) ObservabilityPackFetchPending(status *InstallStatus) error {
	return nil
}

func (r InstallEventsReporter) ObservabilityPackFetchSuccess(status *InstallStatus) error {
	return nil
}

func (r InstallEventsReporter) ObservabilityPackFetchFailed(status *InstallStatus) error {
	return nil
}

func (r InstallEventsReporter) ObservabilityPackInstallPending(status *InstallStatus) error {
	return nil
}

func (r InstallEventsReporter) ObservabilityPackInstallSuccess(status *InstallStatus) error {
	return nil
}

func (r InstallEventsReporter) ObservabilityPackInstallFailed(status *InstallStatus) error {
	return nil
}

func (r InstallEventsReporter) DiscoveryComplete(status *InstallStatus, dm types.DiscoveryManifest) error {
	return nil
}

func (r InstallEventsReporter) createRecipeInstallEvent(status *InstallStatus, statusType installevents.InstallationRecipeStatusType, event RecipeStatusEvent) error {
	s := buildInstallStatus(status, &event, &statusType)
	_, err := r.client.InstallationCreateRecipeEvent(r.accountID, s)

	return err
}

func buildInstallStatus(status *InstallStatus, event *RecipeStatusEvent, statusType *installevents.InstallationRecipeStatusType) installevents.InstallationRecipeStatus {
	i := installevents.InstallationRecipeStatus{
		CliVersion: status.CLIVersion,
		Complete:   status.Complete,
		Error: installevents.InstallationStatusErrorInput{
			Details: status.Error.Details,
			Message: status.Error.Message,
		},
		HostName:        status.DiscoveryManifest.Hostname,
		KernelArch:      status.DiscoveryManifest.KernelArch,
		KernelVersion:   status.DiscoveryManifest.KernelVersion,
		LogFilePath:     status.LogFilePath,
		Os:              status.DiscoveryManifest.OS,
		Platform:        status.DiscoveryManifest.Platform,
		PlatformFamily:  status.DiscoveryManifest.PlatformFamily,
		PlatformVersion: status.DiscoveryManifest.PlatformVersion,
		RedirectURL:     status.RedirectURL,
		TargetedInstall: status.targetedInstall,
	}

	if event != nil {
		i.Name = event.Recipe.Name
		i.DisplayName = event.Recipe.DisplayName
		i.EntityGUID = entities.EntityGUID(event.EntityGUID)
		i.ValidationDurationMilliseconds = event.ValidationDurationMilliseconds
	}

	if statusType != nil {
		i.Status = *statusType
	}

	return i
}
