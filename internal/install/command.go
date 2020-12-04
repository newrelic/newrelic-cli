package install

import (
	"errors"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-cli/internal/credentials"
	"github.com/newrelic/newrelic-client-go/newrelic"
)

var (
	interactiveMode     bool
	autoDiscoveryMode   bool
	recipeFriendlyNames []string
)

// Command represents the install command.
var Command = &cobra.Command{
	Use:    "install",
	Short:  "Install New Relic.",
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		ic := installContext{
			interactiveMode:     interactiveMode,
			autoDiscoveryMode:   autoDiscoveryMode,
			recipeFriendlyNames: recipeFriendlyNames,
		}

		client.WithClientAndProfile(func(nrClient *newrelic.NewRelic, profile *credentials.Profile) {
			if profile == nil {
				log.Fatal(errors.New("default profile has not been set"))
			}

			// Wire up the recipe installer with dependency injection.
			rf := newServiceRecipeFetcher(&nrClient.NerdGraph)
			pf := newRegexProcessFilterer(rf)

			i := newRecipeInstaller(ic,
				newPSUtilDiscoverer(pf),
				rf,
				newGoTaskRecipeExecutor(),
				newPollingRecipeValidator(&nrClient.Nrdb),
			)

			// Run the install.
			i.install()
		})
	},
}

func init() {
	Command.Flags().BoolVarP(&interactiveMode, "interactive", "i", true, "enables interactive mode")
	Command.Flags().BoolVarP(&autoDiscoveryMode, "autoDiscovery", "d", true, "enables auto-discovery mode")
	Command.Flags().StringSliceVarP(&recipeFriendlyNames, "recipe", "r", []string{}, "the name of a recipe to install")
}
