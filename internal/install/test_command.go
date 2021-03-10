package install

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

var (
	testScenario string
)

// TestCommand represents the test command for the install command.
var TestCommand = &cobra.Command{
	Use:    "installTest",
	Short:  "Run a UX test of the install command.",
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		ic := InstallerContext{
			RecipePaths:        recipePaths,
			RecipeNames:        recipeNames,
			SkipDiscovery:      skipDiscovery,
			SkipIntegrations:   skipIntegrations,
			SkipLoggingInstall: skipLoggingInstall,
			AssumeYes:          assumeYes,
		}

		b := NewScenarioBuilder(ic)
		i := b.BuildScenario(TestScenario(testScenario))

		if i == nil {
			log.Fatalf("Scenario %s is not valid.  Valid values are %s", testScenario, strings.Join(TestScenarioValues(), ","))
		}

		if trace {
			log.SetLevel(log.TraceLevel)
		} else if debug {
			log.SetLevel(log.DebugLevel)
		}

		if err := i.Install(); err != nil {
			if err == types.ErrInterrupt {
				return
			}

			log.Fatalf("test failed: %s", err)
		}
	},
}

func init() {
	TestCommand.Flags().StringSliceVarP(&recipePaths, "recipePath", "c", []string{}, "the path to a recipe file to install")
	TestCommand.Flags().StringSliceVarP(&recipeNames, "recipe", "n", []string{}, "the name of a recipe to install")
	TestCommand.Flags().BoolVarP(&skipDiscovery, "skipDiscovery", "d", false, "skips discovery of recommended New Relic integrations")
	TestCommand.Flags().BoolVarP(&skipIntegrations, "skipIntegrations", "r", false, "skips installation of recommended New Relic integrations")
	TestCommand.Flags().BoolVarP(&skipLoggingInstall, "skipLoggingInstall", "l", false, "skips installation of New Relic Logging")
	TestCommand.Flags().StringVarP(&testScenario, "testScenario", "s", string(Basic), fmt.Sprintf("test scenario to run, defaults to BASIC.  Valid values are %s", strings.Join(TestScenarioValues(), ",")))
	TestCommand.Flags().BoolVar(&debug, "debug", false, "debug level logging")
	TestCommand.Flags().BoolVar(&trace, "trace", false, "trace level logging")
	TestCommand.Flags().BoolVarP(&assumeYes, "assumeYes", "y", false, "use \"yes\" for all questions during install")
}
