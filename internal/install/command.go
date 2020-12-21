package install

import (
	"errors"
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-cli/internal/credentials"
	"github.com/newrelic/newrelic-cli/internal/utils"
	"github.com/newrelic/newrelic-client-go/newrelic"
)

var (
	recipePaths        []string
	recipeNames        []string
	skipDiscovery      bool
	skipInfraInstall   bool
	skipIntegrations   bool
	skipLoggingInstall bool
	testMode           bool
)

// Command represents the install command.
var Command = &cobra.Command{
	Use:    "install",
	Short:  "Install New Relic.",
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		go func() {
			<-utils.SignalCtx.Done()
			os.Exit(1)
		}()

		ic := InstallerContext{
			RecipePaths:        recipePaths,
			RecipeNames:        recipeNames,
			SkipDiscovery:      skipDiscovery,
			SkipInfraInstall:   skipInfraInstall,
			SkipIntegrations:   skipIntegrations,
			SkipLoggingInstall: skipLoggingInstall,
		}

		client.WithClientAndProfile(func(nrClient *newrelic.NewRelic, profile *credentials.Profile) {
			err := assertProfileIsValid(profile)
			if err != nil {
				log.Fatal(err)
			}

			i := NewRecipeInstaller(ic, nrClient)

			// Run the install.
			if err := i.Install(); err != nil {
				log.Fatalf("Could not install new Relic: %s", err)
			}
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
	Command.Flags().BoolVarP(&testMode, "testMode", "t", false, "fakes operations for UX testing")
}
