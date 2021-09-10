package execution

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type TerminalStatusReporter struct{}

// NewTerminalStatusReporter is an implementation of the ExecutionStatusReporter interface that reports execution status to STDOUT.
func NewTerminalStatusReporter() *TerminalStatusReporter {
	r := TerminalStatusReporter{}

	return &r
}

func (r TerminalStatusReporter) RecipeFailed(status *InstallStatus, event RecipeStatusEvent) error {
	return nil
}

func (r TerminalStatusReporter) RecipeInstalling(status *InstallStatus, event RecipeStatusEvent) error {
	return nil
}

func (r TerminalStatusReporter) RecipeInstalled(status *InstallStatus, event RecipeStatusEvent) error {
	return nil
}

func (r TerminalStatusReporter) RecipeSkipped(status *InstallStatus, event RecipeStatusEvent) error {
	return nil
}

func (r TerminalStatusReporter) RecipeRecommended(status *InstallStatus, event RecipeStatusEvent) error {
	return nil
}

func (r TerminalStatusReporter) RecipesSelected(status *InstallStatus, recipes []types.OpenInstallationRecipe) error {
	if len(recipes) > 0 {
		fmt.Println("The following will be installed:")
	}

	for _, r := range recipes {
		log.WithFields(log.Fields{
			"name": r.Name,
		}).Debug("found available integration")

		if r.DisplayName != "" {
			fmt.Printf("  %s\n", r.DisplayName)
		} else {
			fmt.Printf("  %s\n", r.Name)
		}
	}

	fmt.Println()

	return nil
}

func (r TerminalStatusReporter) RecipeAvailable(status *InstallStatus, recipe types.OpenInstallationRecipe) error {
	return nil
}

func (r TerminalStatusReporter) InstallStarted(status *InstallStatus) error {
	return nil
}

func (r TerminalStatusReporter) InstallComplete(status *InstallStatus) error {
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
	if status.PlatformLinkGenerator != nil {
		linkToData = status.PlatformLinkGenerator.GenerateRedirectURL(*status)
	}

	if linkToData != "" {
		fmt.Printf("  Your data is available at %s", linkToData)
	}

	fmt.Println()

	return nil
}

func (r TerminalStatusReporter) InstallCanceled(status *InstallStatus) error {
	return nil
}

func (r TerminalStatusReporter) DiscoveryComplete(status *InstallStatus, dm types.DiscoveryManifest) error {
	return nil
}

func (r TerminalStatusReporter) RecipeUnsupported(status *InstallStatus, event RecipeStatusEvent) error {
	return nil
}

func (r TerminalStatusReporter) ObservabilityPackFetchPending(status *InstallStatus) error {
	return nil
}

func (r TerminalStatusReporter) ObservabilityPackFetchSuccess(status *InstallStatus) error {
	return nil
}

func (r TerminalStatusReporter) ObservabilityPackFetchFailed(status *InstallStatus) error {
	return nil
}

func (r TerminalStatusReporter) ObservabilityPackInstallPending(status *InstallStatus) error {
	return nil
}

func (r TerminalStatusReporter) ObservabilityPackInstallSuccess(status *InstallStatus) error {
	return nil
}

func (r TerminalStatusReporter) ObservabilityPackInstallFailed(status *InstallStatus) error {
	return nil
}

func (r TerminalStatusReporter) UpdateRequired(status *InstallStatus) error {
	return nil
}
