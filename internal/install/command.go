package install

import (
	"errors"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-cli/internal/credentials"
	"github.com/newrelic/newrelic-client-go/newrelic"
)

var (
	recipePaths        []string
	recipeNames        []string
	skipDiscovery      bool
	skipInfraInstall   bool
	skipIntegrations   bool
	skipLoggingInstall bool
)

// Command represents the install command.
var Command = &cobra.Command{
	Use:    "install",
	Short:  "Install New Relic.",
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		ic := installContext{
			recipePaths:        recipePaths,
			recipeNames:        recipeNames,
			skipDiscovery:      skipDiscovery,
			skipInfraInstall:   skipInfraInstall,
			skipIntegrations:   skipIntegrations,
			skipLoggingInstall: skipLoggingInstall,
		}

		client.WithClientAndProfile(func(nrClient *newrelic.NewRelic, profile *credentials.Profile) {
			err := assertProfileIsValid(profile)
			if err != nil {
				log.Fatal(err)
			}

			// Wire up the recipe installer with dependency injection.
			rf := newServiceRecipeFetcher(&nrClient.NerdGraph)
			pf := newRegexProcessFilterer(rf)
			ff := newRecipeFileFetcher()

			i := newRecipeInstaller(ic,
				newPSUtilDiscoverer(pf),
				newGlobFileFilterer(),
				rf,
				newGoTaskRecipeExecutor(),
				newPollingRecipeValidator(&nrClient.Nrdb),
				ff,
			)

			// Run the install.
			i.install()
		})
	},
}

func assertProfileIsValid(profile *credentials.Profile) error {
	if profile == nil {
		return errors.New("default profile has not been set")
	}
	if !strings.EqualFold(profile.Region, "US") {
		return fmt.Errorf("only the US region is supported at this time. region %s is not supported", profile.Region)
	}
	return nil
}

func init() {
	Command.Flags().StringSliceVarP(&recipePaths, "recipePath", "c", []string{}, "the path to a recipe file to install")
	Command.Flags().StringSliceVarP(&recipeNames, "recipe", "n", []string{}, "the name of a recipe to install")
	Command.Flags().BoolVarP(&skipDiscovery, "skipDiscovery", "d", false, "skips discovery of recommended New Relic integrations")
	Command.Flags().BoolVarP(&skipInfraInstall, "skipInfraInstall", "i", false, "skips installation of New Relic Infrastructure Agent")
	Command.Flags().BoolVarP(&skipIntegrations, "skipIntegrations", "r", false, "skips installation of recommended New Relic integrations")
	Command.Flags().BoolVarP(&skipLoggingInstall, "skipLoggingInstall", "l", false, "skips installation of New Relic Logging")
}
