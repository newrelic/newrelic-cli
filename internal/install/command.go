package install

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-cli/internal/config"
)

var (
	assumeYes          bool
	advancedMode       bool
	recipeNames        []string
	recipePaths        []string
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
	PreRun: func(cmd *cobra.Command, args []string) {
		config.FatalIfAccountIDNotPresent()
		config.FatalIfActiveProfileFieldStringNotPresent(config.APIKey)
	},
	Run: func(cmd *cobra.Command, args []string) {
		ic := InstallerContext{
			AdvancedMode:       advancedMode,
			AssumeYes:          assumeYes,
			RecipeNames:        recipeNames,
			RecipePaths:        recipePaths,
			SkipDiscovery:      skipDiscovery,
			SkipInfraInstall:   skipInfraInstall,
			SkipIntegrations:   skipIntegrations,
			SkipLoggingInstall: skipLoggingInstall,
		}

		activeProfile := config.GetActiveProfileName()
		if activeProfile == "" {
			log.Fatal("no active profile has been set")
		}

		i := NewRecipeInstaller(ic, client.Client)

		// Run the install.
		if err := i.Install(); err != nil {
			log.Fatalf("Could not install New Relic: %s, check the install log for details: %s", err, config.DefaultLogFile)
		}
	},
}

func init() {
	Command.Flags().StringSliceVarP(&recipePaths, "recipePath", "c", []string{}, "the path to a recipe file to install")
	Command.Flags().StringSliceVarP(&recipeNames, "recipe", "n", []string{}, "the name of a recipe to install")
	Command.Flags().BoolVarP(&skipDiscovery, "skipDiscovery", "d", false, "skips discovery of recommended New Relic integrations")
	Command.Flags().BoolVarP(&skipInfraInstall, "skipInfraInstall", "i", false, "skips installation of New Relic Infrastructure Agent")
	Command.Flags().BoolVarP(&skipIntegrations, "skipIntegrations", "r", false, "skips installation of recommended New Relic integrations")
	Command.Flags().BoolVarP(&skipLoggingInstall, "skipLoggingInstall", "l", false, "skips installation of New Relic Logging")
	Command.Flags().BoolVarP(&assumeYes, "assumeYes", "y", false, "use \"yes\" for all questions during install")
	Command.Flags().BoolVarP(&advancedMode, "advanced", "", false, "use \"advanced\" mode")
}
