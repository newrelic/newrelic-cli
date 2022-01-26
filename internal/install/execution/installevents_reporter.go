package execution

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-client-go/pkg/common"
	"github.com/newrelic/newrelic-client-go/pkg/installevents"

	configAPI "github.com/newrelic/newrelic-cli/internal/config/api"
	"github.com/newrelic/newrelic-cli/internal/install/types"
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

func (r InstallEventsReporter) RecipeDetected(status *InstallStatus, recipe types.OpenInstallationRecipe, event RecipeStatusEvent) error {
	return r.createRecipeInstallEvent(status, installevents.InstallationRecipeStatusTypeTypes.DETECTED, event)
}

func (r InstallEventsReporter) RecipeAvailable(status *InstallStatus, recipe types.OpenInstallationRecipe) error {
	event := RecipeStatusEvent{Recipe: recipe}
	return r.createRecipeInstallEvent(status, installevents.InstallationRecipeStatusTypeTypes.AVAILABLE, event)
}

func (r InstallEventsReporter) RecipeInstalling(status *InstallStatus, event RecipeStatusEvent) error {
	return r.createRecipeInstallEvent(status, installevents.InstallationRecipeStatusTypeTypes.INSTALLING, event)
}

func (r InstallEventsReporter) RecipeSkipped(status *InstallStatus, event RecipeStatusEvent) error {
	return r.createRecipeInstallEvent(status, installevents.InstallationRecipeStatusTypeTypes.SKIPPED, event)
}

func (r InstallEventsReporter) RecipeRecommended(status *InstallStatus, event RecipeStatusEvent) error {
	return r.createRecipeInstallEvent(status, installevents.InstallationRecipeStatusTypeTypes.RECOMMENDED, event)
}

func (r InstallEventsReporter) RecipeUnsupported(status *InstallStatus, event RecipeStatusEvent) error {
	return r.createRecipeInstallEvent(status, installevents.InstallationRecipeStatusTypeTypes.UNSUPPORTED, event)
}

func (r InstallEventsReporter) RecipeInstalled(status *InstallStatus, event RecipeStatusEvent) error {
	return r.createRecipeInstallEvent(status, installevents.InstallationRecipeStatusTypeTypes.INSTALLED, event)
}

func (r InstallEventsReporter) InstallCanceled(status *InstallStatus) error {
	return r.createMultipleRecipeInstallEvents(status, RecipeStatusEvent{})
}

func (r InstallEventsReporter) RecipeFailed(status *InstallStatus, event RecipeStatusEvent) error {
	return r.createRecipeInstallEvent(status, installevents.InstallationRecipeStatusTypeTypes.FAILED, event)
}

func (r InstallEventsReporter) RecipesSelected(status *InstallStatus, recipes []types.OpenInstallationRecipe) error {
	return nil
}

func (r InstallEventsReporter) InstallStarted(status *InstallStatus) error {
	return r.createInstallStatusEvent(installevents.InstallationInstallStateTypeTypes.STARTED, status, RecipeStatusEvent{})
}

func (r InstallEventsReporter) InstallComplete(status *InstallStatus) error {
	return r.createInstallStatusEvent(installevents.InstallationInstallStateTypeTypes.COMPLETED, status, RecipeStatusEvent{})
}

func (r InstallEventsReporter) DiscoveryComplete(status *InstallStatus, dm types.DiscoveryManifest) error {
	return nil
}

func (r InstallEventsReporter) canSendRecipeStatusEvent(status *InstallStatus, recipeStatus installevents.InstallationRecipeStatusType) bool {
	// Prevent sending duplicate DETECTED events when installation is canceled because
	// the DETECTED event was already sent at the beginning of the installation process.
	if recipeStatus == installevents.InstallationRecipeStatusTypeTypes.DETECTED && status.HasCanceledRecipes {
		return false
	}

	return true
}

func (r InstallEventsReporter) createMultipleRecipeInstallEvents(status *InstallStatus, event RecipeStatusEvent) error {
	for _, ss := range status.Statuses {
		if !r.canSendRecipeStatusEvent(status, installevents.InstallationRecipeStatusType(ss.Status)) {
			continue
		}

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
			EntityGUID:                     common.EntityGUID(ss.EntityGUID),
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
			Metadata:                       event.Metadata,
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
	return r.createInstallStatusEvent(
		installevents.InstallationInstallStateType(installevents.InstallationRecipeStatusTypeTypes.UNSUPPORTED),
		status,
		RecipeStatusEvent{},
	)
}

func buildInstallStatus(state installevents.InstallationInstallStateType, status *InstallStatus, event *RecipeStatusEvent) installevents.InstallationInstallStatusInput {
	i := installevents.InstallationInstallStatusInput{
		CliVersion: status.CLIVersion,
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

	fmt.Print("\n\n **************************** \n")
	fmt.Printf("\n Event Reporter RECIPE event:  %+v \n", event)

	i := installevents.InstallationRecipeStatus{
		CliVersion: status.CLIVersion,
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
		i.EntityGUID = common.EntityGUID(event.EntityGUID)
		i.ValidationDurationMilliseconds = event.ValidationDurationMs
		i.TaskPath = strings.Join(event.TaskPath, ",")
		i.Metadata = event.Metadata
	}

	if statusType != nil {
		i.Status = *statusType
	}

	// // REMOVE BEFORE ASKIN FOR A REVIEW ;)
	fmt.Print("\n\n **************************** \n")
	fmt.Printf("\n Event Reporter RECIPE Status:  %+v \n", i.Metadata)
	fmt.Print("\n **************************** \n\n")

	return i
}
