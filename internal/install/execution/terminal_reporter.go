package execution

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type TerminalStatusReporter struct {
}

// NewTerminalStatusReporter is an implementation of the ExecutionStatusReporter interface that reports execution status to STDOUT.
func NewTerminalStatusReporter() *TerminalStatusReporter {
	r := TerminalStatusReporter{}

	return &r
}

func (r TerminalStatusReporter) ReportRecipeFailed(status *StatusRollup, event RecipeStatusEvent) error {
	return nil
}

func (r TerminalStatusReporter) ReportRecipeInstalling(status *StatusRollup, event RecipeStatusEvent) error {
	return nil
}

func (r TerminalStatusReporter) ReportRecipeInstalled(status *StatusRollup, event RecipeStatusEvent) error {
	return nil
}

func (r TerminalStatusReporter) ReportRecipeSkipped(status *StatusRollup, event RecipeStatusEvent) error {
	return nil
}

func (r TerminalStatusReporter) ReportRecipeRecommended(status *StatusRollup, event RecipeStatusEvent) error {
	return nil
}

func (r TerminalStatusReporter) ReportRecipesAvailable(status *StatusRollup, recipes []types.Recipe) error {
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

func (r TerminalStatusReporter) ReportRecipeAvailable(status *StatusRollup, recipe types.Recipe) error {
	return nil
}

func (r TerminalStatusReporter) ReportComplete(status *StatusRollup) error {
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

		fmt.Println("Please refer to the \"Detected observability gaps\" section in the link to your data.")
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
