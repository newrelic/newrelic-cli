package install

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-cli/internal/config"
	configAPI "github.com/newrelic/newrelic-cli/internal/config/api"
	"github.com/newrelic/newrelic-cli/internal/install/types"
	"github.com/newrelic/newrelic-cli/internal/utils"
)

var (
	assumeYes    bool
	localRecipes string
	recipeNames  []string
	recipePaths  []string
	testMode     bool
)

// Command represents the install command.
var Command = &cobra.Command{
	Use:    "install",
	Short:  "Install New Relic.",
	PreRun: client.RequireClient,
	Run: func(cmd *cobra.Command, args []string) {
		ic := types.InstallerContext{
			AssumeYes:    assumeYes,
			LocalRecipes: localRecipes,
			RecipeNames:  recipeNames,
			RecipePaths:  recipePaths,
		}

		logLevel := configAPI.GetLogLevel()
		config.InitFileLogger(logLevel)

		err := assertProfileIsValid()
		if err != nil {
			log.Fatal(err)
		}

		i := NewRecipeInstaller(ic, client.NRClient)

		// Run the install.
		if err := i.Install(); err != nil {
			if err == types.ErrInterrupt {
				return
			}

			if _, ok := err.(*types.UpdateRequiredError); ok {
				return
			}

			fallbackErrorMsg := fmt.Sprintf("We encountered an issue during the installation: %s.", err)
			fallbackHelpMsg := "If this problem persists, visit the documentation and support page for additional help here at https://one.newrelic.com/-/06vjAeZLKjP."

			// In the extremely rare case we run into an uncaught error (e.g. no recipes found),
			// we need to output something to user to sinc we probably haven't displayed anything yet.
			if _, ok := err.(*types.UncaughtError); ok {
				fmt.Print(fallbackErrorMsg)
				fmt.Printf("\nFor more information about about what happened, rerun the command using the --debug flag.\n")
				fmt.Print("\n  newrelic install --debug \n\n")
				fmt.Print(fallbackHelpMsg)
			} else {
				log.Debug(fallbackErrorMsg)
				log.Debug(fallbackHelpMsg)
			}

			fmt.Print("\n\n")
		}
	},
}

func assertProfileIsValid() error {
	accountID := configAPI.GetActiveProfileAccountID()
	if accountID == 0 {
		return fmt.Errorf("accountID is required")
	}

	if configAPI.GetActiveProfileString(config.APIKey) == "" {
		return fmt.Errorf("API key is required")
	}

	if configAPI.GetActiveProfileString(config.Region) == "" {
		return fmt.Errorf("region is required")
	}

	licenseKey, err := client.FetchLicenseKey(accountID, config.FlagProfileName)
	if err != nil {
		return fmt.Errorf("could not fetch license key for account %d: %s", accountID, err)
	}
	if licenseKey != configAPI.GetActiveProfileString(config.LicenseKey) {
		os.Setenv("NEW_RELIC_LICENSE_KEY", licenseKey)
		log.Debugf("using license key %s", utils.Obfuscate(licenseKey))
	}

	// Reinitialize client, overriding fetched values
	c, err := client.NewClient(configAPI.GetActiveProfileName())
	if err != nil {
		// An error was encountered initializing the client.  This may not be a
		// problem since many commands don't require the use of an initialized client
		log.Debugf("error initializing client: %s", err)
	}

	client.NRClient = c

	return nil
}

func init() {
	Command.Flags().StringSliceVarP(&recipePaths, "recipePath", "c", []string{}, "the path to a recipe file to install")
	Command.Flags().StringSliceVarP(&recipeNames, "recipe", "n", []string{}, "the name of a recipe to install")
	Command.Flags().BoolVarP(&testMode, "testMode", "t", false, "fakes operations for UX testing")
	Command.Flags().BoolVarP(&assumeYes, "assumeYes", "y", false, "use \"yes\" for all questions during install")
	Command.Flags().StringVarP(&localRecipes, "localRecipes", "", "", "a path to local recipes to load instead of service other fetching")
}
