package execution

import (
	"time"

	"github.com/newrelic/newrelic-cli/internal/credentials"
	"github.com/newrelic/newrelic-cli/internal/install/types"
	"github.com/newrelic/newrelic-client-go/pkg/entities"
	"github.com/newrelic/newrelic-client-go/pkg/installevents"
	"github.com/newrelic/newrelic-client-go/pkg/nrtime"
)

type InstallEventsReporter struct {
	client    InstalleventsClient
	accountID int
}

func NewInstallEventsReporter(client InstalleventsClient) *InstallEventsReporter {
	r := InstallEventsReporter{
		client: client,
	}

	r.accountID = credentials.DefaultProfile().AccountID

	return &r
}

func (r InstallEventsReporter) RecipeFailed(status *InstallStatus, event RecipeStatusEvent) error {
	s := buildInstallStatus(status, &event, &installevents.InstallationRecipeStatusTypeTypes.FAILED)

	_, err := r.client.InstallationCreateRecipeEvent(r.accountID, s)
	return err
}

func (r InstallEventsReporter) RecipeInstalling(status *InstallStatus, event RecipeStatusEvent) error {
	s := buildInstallStatus(status, &event, &installevents.InstallationRecipeStatusTypeTypes.INSTALLING)

	_, err := r.client.InstallationCreateRecipeEvent(r.accountID, s)
	return err
}

func (r InstallEventsReporter) RecipeInstalled(status *InstallStatus, event RecipeStatusEvent) error {
	s := buildInstallStatus(status, &event, &installevents.InstallationRecipeStatusTypeTypes.INSTALLED)

	_, err := r.client.InstallationCreateRecipeEvent(r.accountID, s)
	return err
}

func (r InstallEventsReporter) RecipeSkipped(status *InstallStatus, event RecipeStatusEvent) error {
	s := buildInstallStatus(status, &event, &installevents.InstallationRecipeStatusTypeTypes.SKIPPED)

	_, err := r.client.InstallationCreateRecipeEvent(r.accountID, s)
	return err
}

func (r InstallEventsReporter) RecipeRecommended(status *InstallStatus, event RecipeStatusEvent) error {
	s := buildInstallStatus(status, &event, &installevents.InstallationRecipeStatusTypeTypes.RECOMMENDED)

	_, err := r.client.InstallationCreateRecipeEvent(r.accountID, s)
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

func (r InstallEventsReporter) RecipeUnsupported(status *InstallStatus, event RecipeStatusEvent) error {
	return r.writeStatus(status)
}

func (r InstallEventsReporter) ObservabilityPackFetchPending(status *InstallStatus) error {
	return r.writeStatus(status)
}

func (r InstallEventsReporter) ObservabilityPackFetchSuccess(status *InstallStatus) error {
	return r.writeStatus(status)
}

func (r InstallEventsReporter) ObservabilityPackFetchFailed(status *InstallStatus) error {
	return r.writeStatus(status)
}

func (r InstallEventsReporter) ObservabilityPackInstallPending(status *InstallStatus) error {
	return r.writeStatus(status)
}

func (r InstallEventsReporter) ObservabilityPackInstallSuccess(status *InstallStatus) error {
	return r.writeStatus(status)
}

func (r InstallEventsReporter) ObservabilityPackInstallFailed(status *InstallStatus) error {
	return r.writeStatus(status)
}

func (r InstallEventsReporter) DiscoveryComplete(status *InstallStatus, dm types.DiscoveryManifest) error {
	return r.writeStatus(status)
}

func (r InstallEventsReporter) writeStatus(status *InstallStatus) error {
	_, err := r.client.InstallationCreateRecipeEvent(r.accountID, buildInstallStatus(status, nil, nil))
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
		i.ValidationDurationMilliseconds = nrtime.EpochMilliseconds(time.Unix(event.ValidationDurationMilliseconds, 0).UTC())
	}

	if statusType != nil {
		i.Status = *statusType
	}

	return i
}
