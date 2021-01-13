package execution

import (
	"fmt"

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

func (r TerminalStatusReporter) ReportRecipesAvailable(status *StatusRollup, recipes []types.Recipe) error {
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
		fmt.Printf("\n\thttps://one.newrelic.com/redirect/entity/%s\n", entityGUID)
	}

	return nil
}
