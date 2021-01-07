package install

import (
	"errors"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-cli/internal/config"
	"github.com/newrelic/newrelic-cli/internal/credentials"
	"github.com/newrelic/newrelic-client-go/newrelic"
)

var (
	assumeYes          bool
	recipeNames        []string
	recipePaths        []string
	skipDiscovery      bool
	skipInfraInstall   bool
	skipIntegrations   bool
	skipLoggingInstall bool
	testMode           bool
	debug              bool
	trace              bool
)

// Command represents the install command.
var Command = &cobra.Command{
	Use:    "install",
	Short:  "Install New Relic.",
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		ic := InstallerContext{
			AssumeYes:          assumeYes,
			RecipeNames:        recipeNames,
			RecipePaths:        recipePaths,
			SkipDiscovery:      skipDiscovery,
			SkipInfraInstall:   skipInfraInstall,
			SkipIntegrations:   skipIntegrations,
			SkipLoggingInstall: skipLoggingInstall,
		}

		client.WithClientAndProfile(func(nrClient *newrelic.NewRelic, profile *credentials.Profile) {
			if trace {
				log.SetLevel(log.TraceLevel)
			} else if debug {
				log.SetLevel(log.DebugLevel)
			}

			err := assertProfileIsValid(profile)
			if err != nil {
				log.Fatal(err)
			}

			i := NewRecipeInstaller(ic, nrClient)

			// Run the install.
			if err := i.Install(); err != nil {
				log.Fatalf("Could not install New Relic: %s, check the install log for details: %s", err, config.DefaultLogFile)
			}
		})
	},
}

func assertProfileIsValid(profile *credentials.Profile) error {
	if profile == nil {
		return errors.New("default profile has not been set")
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
	Command.Flags().BoolVar(&debug, "debug", false, "debug level logging")
	Command.Flags().BoolVar(&trace, "trace", false, "trace level logging")
	Command.Flags().BoolVarP(&assumeYes, "assumeYes", "y", false, "use \"yes\" for all questions during install")
}
