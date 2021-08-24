package execution

import (
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/internal/cli"
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

func (r InstallEventsReporter) RecipeInstalling(status *InstallStatus, event RecipeStatusEvent) error {
	err := r.createRecipeInstallEvent(status, installevents.InstallationRecipeStatusTypeTypes.INSTALLING, event)
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

func (r InstallEventsReporter) RecipeInstalled(status *InstallStatus, event RecipeStatusEvent) error {
	err := r.createRecipeInstallEvent(status, installevents.InstallationRecipeStatusTypeTypes.INSTALLED, event)
	return err
}

func (r InstallEventsReporter) InstallCanceled(status *InstallStatus) error {
	err := r.createMultipleRecipeInstallEvents(status, RecipeStatusEvent{})
	return err
}

func (r InstallEventsReporter) RecipeFailed(status *InstallStatus, event RecipeStatusEvent) error {
	err := r.createRecipeInstallEvent(status, installevents.InstallationRecipeStatusTypeTypes.FAILED, event)
	return err
}

func (r InstallEventsReporter) RecipesSelected(status *InstallStatus, recipes []types.OpenInstallationRecipe) error {
	return nil
}

func (r InstallEventsReporter) InstallComplete(status *InstallStatus) error {
	err := r.createInstallStatusEvent(installevents.InstallationInstallStateTypeTypes.COMPLETED, status, RecipeStatusEvent{})
	return err
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

func (r InstallEventsReporter) createMultipleRecipeInstallEvents(status *InstallStatus, event RecipeStatusEvent) error {
	for _, ss := range status.Statuses {
		i := installevents.InstallationRecipeStatus{
			CliVersion: status.CLIVersion,
			Complete:   status.Complete,
			Error: installevents.InstallationStatusErrorInput{
				Details: ss.Error.Details,
				Message: ss.Error.Message,
			},
			Status:                         installevents.InstallationRecipeStatusType(ss.Status),
			Name:                           ss.Name,
			DisplayName:                    ss.DisplayName,
			EntityGUID:                     entities.EntityGUID(ss.EntityGUID),
			ValidationDurationMilliseconds: ss.ValidationDurationMs,
			HostName:                       status.DiscoveryManifest.Hostname,
			KernelArch:                     status.DiscoveryManifest.KernelArch,
			KernelVersion:                  status.DiscoveryManifest.KernelVersion,
			LogFilePath:                    status.LogFilePath,
			Os:                             status.DiscoveryManifest.OS,
			Platform:                       status.DiscoveryManifest.Platform,
			PlatformFamily:                 status.DiscoveryManifest.PlatformFamily,
			PlatformVersion:                status.DiscoveryManifest.PlatformVersion,
			RedirectURL:                    status.RedirectURL,
			TargetedInstall:                status.targetedInstall,
			InstallId:                      status.InstallID,
			InstallLibraryVersion:          status.InstallLibraryVersion,
		}

		_, err := r.client.InstallationCreateRecipeEvent(r.accountID, i)
		if err != nil {
			log.Debugf("could not create multiple recipe install events: %s", err)
		}
	}
	return nil
}

func (r InstallEventsReporter) createRecipeInstallEvent(status *InstallStatus, statusType installevents.InstallationRecipeStatusType, event RecipeStatusEvent) error {
	s := buildRecipeStatus(status, &event, &statusType)
	_, err := r.client.InstallationCreateRecipeEvent(r.accountID, s)

	return err
}

func (r InstallEventsReporter) createInstallStatusEvent(state installevents.InstallationInstallStateType, status *InstallStatus, event RecipeStatusEvent) error {
	s := buildInstallStatus(state, status, &event)
	_, err := r.client.InstallationCreateInstallStatus(r.accountID, s)

	return err
}

func (r InstallEventsReporter) UpdateRequired(status *InstallStatus) error {
	err := r.createInstallStatusEvent(
		installevents.InstallationInstallStateType(installevents.InstallationRecipeStatusTypeTypes.UNSUPPORTED),
		status,
		RecipeStatusEvent{},
	)

	return err
}

func buildInstallStatus(state installevents.InstallationInstallStateType, status *InstallStatus, event *RecipeStatusEvent) installevents.InstallationInstallStatusInput {
	i := installevents.InstallationInstallStatusInput{
		CliVersion: getCLIVersion(status),
		Error: installevents.InstallationStatusErrorInput{
			Details: status.Error.Details,
			Message: status.Error.Message,
		},
		HostName:              status.DiscoveryManifest.Hostname,
		KernelArch:            status.DiscoveryManifest.KernelArch,
		KernelVersion:         status.DiscoveryManifest.KernelVersion,
		LogFilePath:           status.LogFilePath,
		Os:                    status.DiscoveryManifest.OS,
		Platform:              status.DiscoveryManifest.Platform,
		PlatformFamily:        status.DiscoveryManifest.PlatformFamily,
		PlatformVersion:       status.DiscoveryManifest.PlatformVersion,
		RedirectURL:           status.RedirectURL,
		TargetedInstall:       status.targetedInstall,
		IsUnsupported:         status.DiscoveryManifest.IsUnsupported || status.UpdateRequired,
		State:                 state,
		InstallId:             status.InstallID,
		InstallLibraryVersion: status.InstallLibraryVersion,
	}

	if status.HTTPSProxy != "" {
		i.EnabledProxy = true
	}

	return i
}

func buildRecipeStatus(status *InstallStatus, event *RecipeStatusEvent, statusType *installevents.InstallationRecipeStatusType) installevents.InstallationRecipeStatus {
	i := installevents.InstallationRecipeStatus{
		CliVersion: getCLIVersion(status),
		Complete:   status.Complete,
		Error: installevents.InstallationStatusErrorInput{
			Details: status.Error.Details,
			Message: status.Error.Message,
		},
		HostName:              status.DiscoveryManifest.Hostname,
		KernelArch:            status.DiscoveryManifest.KernelArch,
		KernelVersion:         status.DiscoveryManifest.KernelVersion,
		LogFilePath:           status.LogFilePath,
		Os:                    status.DiscoveryManifest.OS,
		Platform:              status.DiscoveryManifest.Platform,
		PlatformFamily:        status.DiscoveryManifest.PlatformFamily,
		PlatformVersion:       status.DiscoveryManifest.PlatformVersion,
		RedirectURL:           status.RedirectURL,
		TargetedInstall:       status.targetedInstall,
		InstallId:             status.InstallID,
		InstallLibraryVersion: status.InstallLibraryVersion,
	}

	if event != nil {
		i.Name = event.Recipe.Name
		i.DisplayName = event.Recipe.DisplayName
		i.EntityGUID = entities.EntityGUID(event.EntityGUID)
		i.ValidationDurationMilliseconds = event.ValidationDurationMs
		i.TaskPath = strings.Join(event.TaskPath, ",")
	}

	if statusType != nil {
		i.Status = *statusType
	}

	return i
}

func getCLIVersion(status *InstallStatus) string {
	cliVersion := status.CLIVersion
	if cliVersion == "" {
		return cli.Version()
	}

	return cliVersion
}
