package finstall

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/newrelic/foil/install"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-cli/internal/config"
	configAPI "github.com/newrelic/newrelic-cli/internal/config/api"
	"github.com/newrelic/newrelic-cli/internal/install/types"
)

var (
	assumeYes    bool     // non-interactive/robot-install
	localRecipes string   // redirect import() for this local repository via JIT, maybe, somehow... like go.mod `replace` directive
	recipeNames  []string // a list of recipes names; TODO: filter down discoverableRecipes to introduce targeted install
	// recipePaths  []string // we should cut down on the permutations of potential recipe location, and expect a defined, conventional directory structure
	// testMode     bool     // I am not sure if this actually does anything...
	tags    []string // tags!
	verbose bool     // much more output
)

// Command represents the install command.
var Command = &cobra.Command{
	Use:    "finstall",
	Short:  "Install New Relic.",
	PreRun: client.RequireClient,
	// we should handle errors here. returning an error implies handling in Execute()
	Run: func(cmd *cobra.Command, args []string) {

		ic := types.InstallerContext{
			AssumeYes:    assumeYes,
			LocalRecipes: localRecipes,
			RecipeNames:  recipeNames,
		}

		ic.SetTags(tags)

		logLevel := configAPI.GetLogLevel()

		config.InitFileLogger(logLevel)

		// TODO: set variables
		// i := NewRecipeInstaller(ic, client.NRClient)

		// TODO: context within Run()
		if err := install.Run(); err != nil {

			// must handle interrupts
			if err == types.ErrInterrupt {
				return
			}

			// TODO: handle
			// if _, ok := err.(*types.UpdateRequiredError); ok {
			// 	return
			// }

			// TODO: handle
			// if e, ok := err.(*nrErrors.PaymentRequiredError); ok {
			// 	return e // TODO
			// }

			fallbackErrorMsg := fmt.Sprintf("\nWe encountered an issue during the installation: %s.", err)

			fallbackHelpMsg := "If this problem persists, visit the documentation and support page for additional help here at https://docs.newrelic.com/docs/infrastructure/install-infrastructure-agent/get-started/requirements-infrastructure-agent/"

			// In the extremely rare case we run into an uncaught error (e.g. no recipes found),
			// we need to output something to user to sinc we probably haven't displayed anything yet.

			fmt.Println(fallbackErrorMsg)

			fmt.Println(fallbackHelpMsg)
		}
	},
}

func init() {
	Command.Flags().StringSliceVarP(&recipeNames, "recipe", "n", []string{}, "provide a comma-separated list of recipes to install, i.e.: --recipe nginx,redis")
	Command.Flags().BoolVarP(&assumeYes, "assumeYes", "y", false, "answer \"yes\" to all prompts; performs a non-interactive installation of New Relic")
	Command.Flags().BoolVarP(&verbose, "verbose", "v", false, "give detailed output, even on success")
	Command.Flags().StringVarP(&localRecipes, "localRecipes", "l", "", "provide a relative or absolute path to local open-install-libary directory")
	Command.Flags().StringSliceVarP(&tags, "tag", "t", []string{}, "tags the entity during installation with a comma-separated list of kev-value pairs, i.e.: --tag tag1:test,tag2:test")
}
