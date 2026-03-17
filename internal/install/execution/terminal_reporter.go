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
	if len(status.Statuses) == 0 {
		return nil
	}

	hasInstalled := status.hasAnyRecipeStatus(RecipeStatusTypes.INSTALLED)
	hasFailures := status.hasAnyRecipeStatus(RecipeStatusTypes.FAILED) || status.hasAnyRecipeStatus(RecipeStatusTypes.UNSUPPORTED)
	isAgentControl := r.isAgentControlInstalled(status)
	isAgentControlAttempt := r.isOnlyAgentControlAttempt(status)

	if hasInstalled {
		fmt.Print("\n  New Relic installation complete \n\n")
	}

	fmt.Println("  --------------------")
	fmt.Println("  What's next?")
	r.printInstallationSummary(os.Stdout, status)

	var msg, link string
	dataLink := ""
	if status.PlatformLinkGenerator != nil && !isAgentControl && !isAgentControlAttempt {
		dataLink = status.PlatformLinkGenerator.GenerateRedirectURL(*status)
	}

	switch {
	case isAgentControl:
		msg = "Learn about configuring your agent and fleet:"
		if status.PlatformLinkGenerator != nil {
			link = status.PlatformLinkGenerator.GenerateFleetConfigurationDocLink()
		}

	case !hasInstalled:
		msg = "The installation is incomplete. Follow the instructions at the URL below to complete the installation process."
		if isAgentControlAttempt && status.PlatformLinkGenerator != nil {
			link = status.PlatformLinkGenerator.GenerateGuidedInstallDocLink()
		} else {
			link = dataLink
		}

	case hasFailures:
		msg = "Installation was successful overall, however, one or more installations could not be completed.\n  Follow the instructions at the URL below to complete the installation process."
		link = dataLink

	default:
		msg = "View your data at the link below:"
		link = dataLink
	}

	if link != "" {
		fmt.Printf("\n  %s\n", msg)
		fmt.Printf("  %s  %s\n", color.GreenString(ux.IconArrowRight), link)
	}

	r.printLoggingLink(status)
	r.printFleetLink(status)

	fmt.Printf("\n  --------------------\n\n")

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

func (r TerminalStatusReporter) printFleetLink(status *InstallStatus) {
	linkToFleet := ""
	fleetMsg := "Check out your agent in our Fleet Control UI:\n"
	statuses := r.getRecipesStatusesForInstallationSummary(status)

	for _, s := range statuses {
		isAgentControlRecipe := s.Name == types.AgentControlRecipeName

		if s.Status == RecipeStatusTypes.INSTALLED && isAgentControlRecipe {
			// Use NR_CLI_FLEET_ID environment variable if available, otherwise fall back to entity GUID
			fleetGUID := os.Getenv("NR_CLI_FLEET_ID")
			if fleetGUID == "" {
				fleetGUID = status.HostEntityGUID()
			}

			log.WithFields(log.Fields{
				"fleetGUID":      fleetGUID,
				"fromEnvVar":     os.Getenv("NR_CLI_FLEET_ID") != "",
				"hostEntityGUID": status.HostEntityGUID(),
			}).Debug("Generating fleet link")

			linkToFleet = status.PlatformLinkGenerator.GenerateFleetLink(fleetGUID)
		}
	}

	if linkToFleet != "" {
		fmt.Println("")
		fmt.Printf("\n  %s", fleetMsg)
		fmt.Printf("  %s  %s", color.GreenString(ux.IconArrowRight), linkToFleet)
	}
}

func (r TerminalStatusReporter) isAgentControlInstalled(status *InstallStatus) bool {
	statuses := r.getRecipesStatusesForInstallationSummary(status)
	isAgentControlRecipeInstalled := false

	for _, s := range statuses {
		if s.Status == RecipeStatusTypes.INSTALLED {
			isAgentControlRecipe := s.Name == types.AgentControlRecipeName
			if isAgentControlRecipe {
				isAgentControlRecipeInstalled = true
				break
			}
		}
	}
	return isAgentControlRecipeInstalled
}

func (r TerminalStatusReporter) isOnlyAgentControlAttempt(status *InstallStatus) bool {
	statuses := r.getRecipesStatusesForInstallationSummary(status)
	if len(statuses) == 0 {
		return false
	}

	for _, s := range statuses {
		isAgentControlRecipe := s.Name == types.AgentControlRecipeName
		if !isAgentControlRecipe {
			return false
		}
	}

	return true
}

func (r TerminalStatusReporter) printLoggingLink(status *InstallStatus) {
	linkToLogging := ""
	loggingMsg := "View your logs at the link below:\n"
	statusesToDisplay := r.getRecipesStatusesForInstallationSummary(status)

	for _, s := range statusesToDisplay {
		if s.Status == RecipeStatusTypes.INSTALLED && s.Name == types.LoggingRecipeName {
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
