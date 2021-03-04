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
	if status.hasFailed() {
		return fmt.Errorf("one or more integrations failed to install, check the install log for more details: %s", status.LogFilePath)
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

	if len(status.EntityGUIDs) > 0 {
		infraLink := fmt.Sprintf("https://one.newrelic.com/redirect/entity/%s", status.EntityGUIDs[0])

		fmt.Printf("  Your data is available at %s", infraLink)
	}

	fmt.Println()

	return nil
}

func (r TerminalStatusReporter) DiscoveryComplete(status *InstallStatus, dm types.DiscoveryManifest) error {
	return nil
}
