package execution

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type TerminalStatusReporter struct {
	successLinkGenerator SuccessLinkGenerator
}

// NewTerminalStatusReporter is an implementation of the ExecutionStatusReporter interface that reports execution status to STDOUT.
func NewTerminalStatusReporter() *TerminalStatusReporter {
	r := TerminalStatusReporter{
		successLinkGenerator: NewConcreteSuccessLinkGenerator(),
	}

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

func (r TerminalStatusReporter) RecipesAvailable(status *InstallStatus, recipes []types.Recipe) error {
	return nil
}

func (r TerminalStatusReporter) RecipesSelected(status *InstallStatus, recipes []types.Recipe) error {
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

func (r TerminalStatusReporter) RecipeAvailable(status *InstallStatus, recipe types.Recipe) error {
	return nil
}

func (r TerminalStatusReporter) InstallComplete(status *InstallStatus) error {
	if status.isCanceled() {
		return nil
	}

	if status.hasFailed() {
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

	fmt.Println("  New Relic installation complete!")

	linkToData := r.getSuccessLink(status)

	if linkToData != "" {
		fmt.Printf("  Your data is available at %s", linkToData)
	}

	fmt.Println()

	return nil
}

func (r *TerminalStatusReporter) getSuccessLink(status *InstallStatus) string {
	if status.hasAnyRecipeStatus(RecipeStatusTypes.INSTALLED) {
		switch t := status.successLinkConfig.Type; {
		case strings.EqualFold(t, "explorer"):
			return r.successLinkGenerator.GenerateExplorerLink(status.successLinkConfig.Filter)
		default:
			return r.successLinkGenerator.GenerateEntityLink(status.HostEntityGUID())
		}
	}
	return ""
}

func (r TerminalStatusReporter) InstallCanceled(status *InstallStatus) error {
	return nil
}

func (r TerminalStatusReporter) DiscoveryComplete(status *InstallStatus, dm types.DiscoveryManifest) error {
	return nil
}
