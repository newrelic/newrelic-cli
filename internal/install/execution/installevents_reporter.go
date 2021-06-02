package execution

import (
	"github.com/newrelic/newrelic-cli/internal/install/types"
	"github.com/newrelic/newrelic-client-go/pkg/installevents"
)

type InstallEventsReporter struct {
	client InstalleventsClient
}

func NewInstallEventsReporter(client InstalleventsClient) *InstallEventsReporter {
	r := InstallEventsReporter{
		client: client,
	}

	return &r
}

func (r InstallEventsReporter) RecipeFailed(status *InstallStatus, event RecipeStatusEvent) error {
	_, err := r.client.CreateInstallEvent(buildInstallStatus(event, installevents.RecipeStatusTypeTypes.FAILED))
	return err
}

func (r InstallEventsReporter) RecipeInstalling(status *InstallStatus, event RecipeStatusEvent) error {
	_, err := r.client.CreateInstallEvent(buildInstallStatus(event, installevents.RecipeStatusTypeTypes.INSTALLING))
	return err
}

func (r InstallEventsReporter) RecipeInstalled(status *InstallStatus, event RecipeStatusEvent) error {
	_, err := r.client.CreateInstallEvent(buildInstallStatus(event, installevents.RecipeStatusTypeTypes.INSTALLED))
	return err
}

func (r InstallEventsReporter) RecipeSkipped(status *InstallStatus, event RecipeStatusEvent) error {
	_, err := r.client.CreateInstallEvent(buildInstallStatus(event, installevents.RecipeStatusTypeTypes.SKIPPED))
	return err
}

func (r InstallEventsReporter) RecipeRecommended(status *InstallStatus, event RecipeStatusEvent) error {
	_, err := r.client.CreateInstallEvent(buildInstallStatus(event, installevents.RecipeStatusTypeTypes.SKIPPED))
	return err
}

func (r InstallEventsReporter) RecipesAvailable(status *InstallStatus, recipes []types.OpenInstallationRecipe) error {
	return r.writeStatus(status)
}

func (r InstallEventsReporter) RecipesSelected(status *InstallStatus, recipes []types.OpenInstallationRecipe) error {
	return r.writeStatus(status)
}

func (r InstallEventsReporter) RecipeAvailable(status *InstallStatus, recipe types.OpenInstallationRecipe) error {
	return r.writeStatus(status)
}

func (r InstallEventsReporter) InstallComplete(status *InstallStatus) error {
	return r.writeStatus(status)
}

func (r InstallEventsReporter) InstallCanceled(status *InstallStatus) error {
	return r.writeStatus(status)
}

func (r InstallEventsReporter) DiscoveryComplete(status *InstallStatus, dm types.DiscoveryManifest) error {
	return r.writeStatus(status)
}

func (r InstallEventsReporter) writeStatus(status *InstallStatus) error {
	_, err := r.client.CreateInstallMetadata(buildInstallMetadata(status))
	return err
}

func buildInstallStatus(event RecipeStatusEvent, status installevents.RecipeStatusType) installevents.InstallStatus {
	i := installevents.InstallStatus{
		// Error:                          "",
		DisplayName: event.Recipe.DisplayName,
		EntityGUID:  event.EntityGUID,
		Name:        event.Recipe.Name,
		Status:      status,
		// ValidationDurationMilliseconds: event.ValidationDurationMilliseconds,
	}

	return i
}

func buildInstallMetadata(status *InstallStatus) installevents.InputInstallMetadata {
	i := installevents.InputInstallMetadata{}

	return i
}
