package install

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-cli/internal/config"
	"github.com/newrelic/newrelic-cli/internal/configuration"
	"github.com/newrelic/newrelic-cli/internal/install/types"
)

var (
	assumeYes          bool
	localRecipes       string
	recipeNames        []string
	recipePaths        []string
	skipIntegrations   bool
	skipLoggingInstall bool
	skipApm            bool
	skipInfra          bool
	testMode           bool
	debug              bool
	trace              bool
)

// Command represents the install command.
var Command = &cobra.Command{
	Use:   "install",
	Short: "Install New Relic.",
	Run: func(cmd *cobra.Command, args []string) {
		ic := types.InstallerContext{
			AssumeYes:          assumeYes,
			LocalRecipes:       localRecipes,
			RecipeNames:        recipeNames,
			RecipePaths:        recipePaths,
			SkipIntegrations:   skipIntegrations,
			SkipLoggingInstall: skipLoggingInstall,
			SkipApm:            skipApm,
			SkipInfraInstall:   skipInfra,
		}

		config.InitFileLogger()

		if trace {
			log.SetLevel(log.TraceLevel)
			client.Client.SetLogLevel("trace")
		} else if debug {
			log.SetLevel(log.DebugLevel)
			client.Client.SetLogLevel("debug")
		}

		err := assertProfileIsValid()
		if err != nil {
			log.Fatal(err)
		}

		i := NewRecipeInstaller(ic, client.Client)

		// Run the install.
		if err := i.Install(); err != nil {
			if err == types.ErrInterrupt {
				return
			}

			log.Fatalf("We encountered an error during the installation: %s. If this problem persists please visit the documentation and support page for additional help here: https://one.newrelic.com/-/06vjAeZLKjP", err)
		}
	},
}

func assertProfileIsValid() error {
	if configuration.GetActiveProfileInt(configuration.AccountID) == 0 {
		return fmt.Errorf("accountID is required")
	}

	if configuration.GetActiveProfileString(configuration.APIKey) == "" {
		return fmt.Errorf("API key is required")
	}

	if configuration.GetActiveProfileString(configuration.Region) == "" {
		return fmt.Errorf("region is required")
	}

	return nil
}

func init() {
	Command.Flags().StringSliceVarP(&recipePaths, "recipePath", "c", []string{}, "the path to a recipe file to install")
	Command.Flags().StringSliceVarP(&recipeNames, "recipe", "n", []string{}, "the name of a recipe to install")
	Command.Flags().BoolVarP(&skipIntegrations, "skipIntegrations", "r", false, "skips installation of recommended New Relic integrations")
	Command.Flags().BoolVarP(&skipLoggingInstall, "skipLoggingInstall", "l", false, "skips installation of New Relic Logging")
	Command.Flags().BoolVarP(&skipApm, "skipApm", "a", false, "skips installation for APM")
	Command.Flags().BoolVarP(&skipInfra, "skipInfra", "i", false, "skips installation for infrastructure agent (only for targeted install)")
	Command.Flags().BoolVarP(&testMode, "testMode", "t", false, "fakes operations for UX testing")
	Command.Flags().BoolVar(&debug, "debug", false, "debug level logging")
	Command.Flags().BoolVar(&trace, "trace", false, "trace level logging")
	Command.Flags().BoolVarP(&assumeYes, "assumeYes", "y", false, "use \"yes\" for all questions during install")
	Command.Flags().StringVarP(&localRecipes, "localRecipes", "", "", "a path to local recipes to load instead of service other fetching")
}
