package execution

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/fatih/color"
	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/internal/install/types"
	"github.com/newrelic/newrelic-cli/internal/install/ux"
)

type TerminalStatusReporter struct{}

// NewTerminalStatusReporter is an implementation of the ExecutionStatusReporter interface that reports execution status to STDOUT.
func NewTerminalStatusReporter() *TerminalStatusReporter {
	r := TerminalStatusReporter{}

	return &r
}

func (r TerminalStatusReporter) RecipeDetected(status *InstallStatus, event RecipeStatusEvent) error {
	return nil
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

func (r TerminalStatusReporter) RecipeAvailable(status *InstallStatus, event RecipeStatusEvent) error {
	return nil
}

func (r TerminalStatusReporter) RecipeCanceled(status *InstallStatus, event RecipeStatusEvent) error {
	return nil
}

func (r TerminalStatusReporter) InstallStarted(status *InstallStatus) error {
	return nil
}

func (r TerminalStatusReporter) InstallComplete(status *InstallStatus) error {
	linkToData := ""
	if status.PlatformLinkGenerator != nil {
		linkToData = status.PlatformLinkGenerator.GenerateRedirectURL(*status)
	}

	hasStatuses := len(status.Statuses) > 0
	if hasStatuses {
		hasInstalledRecipes := status.hasAnyRecipeStatus(RecipeStatusTypes.INSTALLED)

		if hasInstalledRecipes {
			fmt.Print("\n  New Relic installation complete \n\n")
		}

		fmt.Println("  --------------------")
		fmt.Println("  Installation Summary")
		fmt.Println("")
		r.printInstallationSummary(os.Stdout, status)

		msg := "View your data at the link below:\n"
		followInstructionsMsg := "Follow the instructions at the URL below to complete the installation process."
		if hasInstalledRecipes && (status.hasAnyRecipeStatus(RecipeStatusTypes.FAILED) || status.hasAnyRecipeStatus(RecipeStatusTypes.UNSUPPORTED)) {
			msg = fmt.Sprintf("Installation was successful overall, however, one or more installations could not be completed.\n  %s \n\n", followInstructionsMsg)
		} else if !hasInstalledRecipes {
			msg = fmt.Sprintf("Installation incomplete. %s \n\n", followInstructionsMsg)
		}

		if linkToData != "" {
			fmt.Printf("\n  %s", msg)
			fmt.Printf("  %s  %s", color.GreenString(ux.IconArrowRight), linkToData)
		}

		r.printLoggingLink(status)

		fmt.Println()
		fmt.Println("\n  --------------------")
		fmt.Println()
	}

	return nil
}

func (r TerminalStatusReporter) InstallCanceled(status *InstallStatus) error {
	fmt.Print("\n\n")
	fmt.Println("  Installation canceled.")
	fmt.Println("  To finish your installation please use New Relic's installation wizard using the following link.")
	fmt.Printf("  %s  %s", color.GreenString(ux.IconArrowRight), status.PlatformLinkGenerator.GenerateRedirectURL(*status))
	fmt.Print("\n\n")

	return nil
}

func (r TerminalStatusReporter) DiscoveryComplete(status *InstallStatus, dm types.DiscoveryManifest) error {
	return nil
}

func (r TerminalStatusReporter) RecipeUnsupported(status *InstallStatus, event RecipeStatusEvent) error {
	return nil
}

func (r TerminalStatusReporter) UpdateRequired(status *InstallStatus) error {
	return nil
}

func (r TerminalStatusReporter) printLoggingLink(status *InstallStatus) {
	linkToLogging := ""
	loggingMsg := "View your logs at the link below:\n"
	statusesToDisplay := r.getRecipesStatusesForInstallationSummary(status)

	for _, s := range statusesToDisplay {
		if s.Status == RecipeStatusTypes.INSTALLED && (s.Name == types.LoggingRecipeName || s.Name == types.LoggingSuperAgentRecipeName) {
			linkToLogging = status.PlatformLinkGenerator.GenerateLoggingLink(status.HostEntityGUID())
		}
	}

	if linkToLogging != "" {
		fmt.Println("")
		fmt.Printf("\n  %s", loggingMsg)
		fmt.Printf("  %s  %s", color.GreenString(ux.IconArrowRight), linkToLogging)
	}
}

func (r TerminalStatusReporter) printInstallationSummary(w io.Writer, status *InstallStatus) {
	statusesToDisplay := r.getRecipesStatusesForInstallationSummary(status)

	for _, s := range statusesToDisplay {
		statusSuffix := strings.ToLower(string(s.Status))

		if s.Status == RecipeStatusTypes.INSTALLED {
			statusSuffix = color.GreenString(statusSuffix)
		}

		if s.Status == RecipeStatusTypes.FAILED {
			statusSuffix = color.YellowString("incomplete")
		}

		if s.Status == RecipeStatusTypes.CANCELED {
			statusSuffix = color.YellowString(statusSuffix)
		}

		if s.Status == RecipeStatusTypes.UNSUPPORTED {
			statusSuffix = color.RedString(statusSuffix)
		}

		fmt.Fprintf(w, "  %s  %s  (%s)  \n", StatusIconMap[s.Status], s.DisplayName, statusSuffix)
	}
}

// getRecipesStatusesForInstallationSummary returns the recipe installation results
// to show the user. Recipes with a DETECTED status are not displayed to the user
// because a DETECTED status at this point means the instrumentation was not installed.
func (r TerminalStatusReporter) getRecipesStatusesForInstallationSummary(status *InstallStatus) []*RecipeStatus {
	statusesToDisplay := []*RecipeStatus{}
	for _, s := range status.Statuses {
		if s.Status != RecipeStatusTypes.DETECTED && s.Status != RecipeStatusTypes.RECOMMENDED {
			statusesToDisplay = append(statusesToDisplay, s)
		}
	}

	return statusesToDisplay
}
