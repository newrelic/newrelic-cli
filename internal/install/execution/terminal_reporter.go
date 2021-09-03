package execution

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
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

	hasInstalledRecipes := status.hasAnyRecipeStatus(RecipeStatusTypes.INSTALLED)
	if hasInstalledRecipes {
		fmt.Println("  New Relic installation complete")

		printInstallationSummary(status)
	}

	linkToData := ""
	if status.PlatformLinkGenerator != nil {
		linkToData = status.PlatformLinkGenerator.GenerateRedirectURL(*status)
	}

	if hasInstalledRecipes && linkToData != "" {
		fmt.Println("  Use the link below to view your data.")
		fmt.Printf("  %s  %s", color.GreenString("\u2B95"), linkToData)
	}

	fmt.Println()

	return nil
}

func printInstallationSummary(status *InstallStatus) {
	fmt.Println()
	fmt.Println("  --------------------")
	fmt.Println("  Installation Summary")
	fmt.Println()

	successIcon := color.GreenString("\u2705")
	failedIcon := color.RedString("\u274C")
	neurtalIcon := color.RedString("\u2796")

	for _, s := range status.Statuses {
		var icon string
		suffix := fmt.Sprintf("(%s)", strings.ToLower(string(s.Status)))

		if s.Status == RecipeStatusTypes.INSTALLED {
			icon = successIcon
			suffix = fmt.Sprintf("(%s)", color.GreenString(strings.ToLower(string(s.Status))))
		}

		if s.Status == RecipeStatusTypes.FAILED {
			icon = failedIcon
			suffix = fmt.Sprintf("(%s)", color.RedString(strings.ToLower(string(s.Status))))
		}

		if s.Status == RecipeStatusTypes.SKIPPED || s.Status == RecipeStatusTypes.CANCELED {
			icon = neurtalIcon
		}

		fmt.Printf("  %s  %s %s\n", icon, s.DisplayName, suffix)
	}

	fmt.Println("  --------------------")
	fmt.Println()
}

func (r TerminalStatusReporter) InstallCanceled(status *InstallStatus) error {
	fmt.Println()
	fmt.Println("  Installation canceled.")
	fmt.Println("  To finish your installation please use New Relic's installation wizard using the following link.")
	fmt.Printf("  %s  %s", color.GreenString("\u2B95"), status.PlatformLinkGenerator.GenerateRedirectURL(*status))
	fmt.Print("\n\n")

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
