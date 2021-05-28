package execution

import (
	"fmt"

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
	return nil
}

func (r InstallEventsReporter) RecipesSelected(status *InstallStatus, recipes []types.OpenInstallationRecipe) error {
	return nil
}

func (r InstallEventsReporter) RecipeAvailable(status *InstallStatus, recipe types.OpenInstallationRecipe) error {
	return nil
}

func (r InstallEventsReporter) InstallComplete(status *InstallStatus) error {
	if status.hasAnyRecipeStatus(RecipeStatusTypes.CANCELED) {
		return nil
	}

	if status.hasAnyRecipeStatus(RecipeStatusTypes.FAILED) {
		fmt.Printf("  One or more installations failed.  Check the install log for more details: %s\n", status.LogFilePath)
	}

	recs := status.recommendations()

	if len(recs) > 0 {
		fmt.Println("  ---")
		fmt.Println("  Instrumentation recommendations")
		fmt.Println("  We discovered some additional instrumentation opportunities:")

		for _, recommendation := range recs {
			fmt.Printf("  - %s\n", recommendation.DisplayName)
		}

		fmt.Println("Please refer to the \"Data gaps\" section in the link to your data.")
		fmt.Println("  ---")
	}

	if status.hasAnyRecipeStatus(RecipeStatusTypes.INSTALLED) {
		fmt.Println("  New Relic installation complete!")
	}

	linkToData := ""
	if status.successLinkGenerator != nil {
		linkToData = status.successLinkGenerator.GenerateRedirectURL(*status)
	}

	if linkToData != "" {
		fmt.Printf("  Your data is available at %s", linkToData)
	}

	fmt.Println()

	return nil
}

func (r InstallEventsReporter) InstallCanceled(status *InstallStatus) error {
	return nil
}

func (r InstallEventsReporter) DiscoveryComplete(status *InstallStatus, dm types.DiscoveryManifest) error {
	return nil
}

func (r InstallEventsReporter) writeStatus(status *InstallStatus) error {

	return nil
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
