package execution

import (
	"fmt"

	"github.com/newrelic/newrelic-cli/internal/install/types"
	log "github.com/sirupsen/logrus"
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

func (r TerminalStatusReporter) ReportRecipesAvailable(status *StatusRollup, recipes []types.Recipe) error {
	if len(recipes) > 0 {
		fmt.Println("The following will be installed, based on what has been discovered on your system.")
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

	msg := `
  Success! Your data is available in New Relic.

  Go to New Relic to confirm and start exploring your data.`

	fmt.Println(msg)

	for _, entityGUID := range status.EntityGUIDs {
		fmt.Printf("\n  https://one.newrelic.com/redirect/entity/%s\n", entityGUID)
	}

	fmt.Println()

	return nil
}
