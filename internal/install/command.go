package install

import (
	"fmt"
	"errors"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-cli/internal/credentials"
	"github.com/newrelic/newrelic-client-go/newrelic"
)

var (
	specifyActions  bool
	interactiveMode bool
	installLogging  bool
	recipeNames     []string
	recipePaths     []string
)

// Command represents the install command.
var Command = &cobra.Command{
	Use:    "install",
	Short:  "Install New Relic.",
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		ic := installContext{
			interactiveMode: interactiveMode,
			installLogging:  installLogging,
			recipeNames:     recipeNames,
			recipePaths:     recipePaths,
			specifyActions:  specifyActions,
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
	if profile.Region != "US" {
		return fmt.Errorf("Only the US region is supported at this time. region %s is not supported.", profile.Region)
	}
	return nil
}

func init() {
	Command.Flags().BoolVarP(&interactiveMode, "interactive", "i", false, "enables interactive mode if specifyActions has been used")
	Command.Flags().BoolVarP(&installLogging, "installLogging", "l", false, "installs New Relic logging if specifyActions has been used")
	Command.Flags().BoolVarP(&specifyActions, "specifyActions", "s", false, "specify the actions to be run during install")
	Command.Flags().StringSliceVarP(&recipeNames, "recipe", "r", []string{}, "the name of a recipe to install")
	Command.Flags().StringSliceVarP(&recipePaths, "recipeFile", "c", []string{}, "a recipe file to install")
}
