package uninstall

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/install/discovery"
	"github.com/newrelic/newrelic-cli/internal/install/execution"
	"github.com/newrelic/newrelic-cli/internal/install/recipes"
	"github.com/newrelic/newrelic-cli/internal/install/types"
)

var (
	recipeName   string
	assumeYes    bool
	force        bool
	localRecipes string
	recipePaths  []string
)

// newProcessEvaluator is overridden in tests to exercise RunE end-to-end with a fake process list.
var newProcessEvaluator = func() recipes.ProcessEvaluatorInterface { return recipes.NewProcessEvaluator() }

// Command represents the uninstall command.
var Command = &cobra.Command{
	Use:     "uninstall",
	Short:   "Remove a New Relic integration previously installed by `newrelic install`.",
	Long:    "Remove a single named integration/agent from this host, following the same recipe model as `newrelic install`.",
	Example: "newrelic uninstall -n infrastructure-agent-installer",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		processEvaluator := newProcessEvaluator()
		if err := ensureSingleConcurrentUninstall(ctx, processEvaluator); err != nil {
			return err
		}

		manifest, err := discovery.NewPSUtilDiscoverer().Discover(ctx)
		if err != nil {
			return fmt.Errorf("could not discover host information: %w", err)
		}

		fetcher := newRecipeFetcher()
		loaderFunc := func() ([]*types.OpenInstallationRecipe, error) {
			return fetcher.FetchRecipes(ctx)
		}

		u := &RecipeUninstaller{
			Finder:       NewRepositoryRecipeFinder(loaderFunc, manifest),
			Detector:     NewEvaluatorProcessDetector(processEvaluator, recipes.NewScriptEvaluator()),
			Confirmer:    NewSurveyConfirmer(),
			TaskfileExec: execution.NewGoTaskRecipeExecutor(),
		}

		res := u.Run(ctx, Options{
			RecipeName: recipeName,
			AssumeYes:  assumeYes,
			Force:      force,
		})

		return reportResult(res)
	},
}

// ensureSingleConcurrentUninstall refuses to proceed if another install/uninstall is already running.
func ensureSingleConcurrentUninstall(ctx context.Context, evaluator recipes.ProcessEvaluatorInterface) error {
	processes := evaluator.GetOrLoadProcesses(ctx)
	if count := recipes.CountConcurrentNewRelicSubcommandProcesses(processes); count > 1 {
		return fmt.Errorf("only 1 newrelic install/uninstall command can run at one time, found %d currently executing. Please retry later, or terminate the other newrelic installations", count)
	}
	return nil
}

func newRecipeFetcher() recipes.RecipeFetcher {
	if localRecipes != "" {
		return &recipes.LocalRecipeFetcher{Path: localRecipes}
	}
	if len(recipePaths) > 0 {
		return recipes.NewRecipeFileFetcher(recipePaths)
	}
	return recipes.NewEmbeddedRecipeFetcher()
}

func reportResult(res Result) error {
	switch res.Status {
	case StatusUninstalled:
		fmt.Printf("Successfully removed %s.\n", recipeName)
		return nil

	case StatusNoAutomatedUninstall:
		fmt.Printf("No automated uninstall is available yet for %s.\n", recipeName)
		return nil

	case StatusAborted:
		fmt.Println("Uninstall cancelled.")
		return nil

	case StatusUnsupported:
		return fmt.Errorf("%s is not supported on this host", recipeName)

	default: // StatusFailed
		for _, w := range res.Warnings {
			log.Warn(w)
		}
		return fmt.Errorf("failed to remove %s: %w", recipeName, res.Err)
	}
}

func init() {
	Command.Flags().StringVarP(&recipeName, "recipe", "n", "", "the name of the recipe to uninstall")
	Command.Flags().BoolVarP(&assumeYes, "assumeYes", "y", false, "use \"yes\" for all confirmation prompts")
	Command.Flags().BoolVarP(&force, "force", "", false, "skip the not-detected safety check and proceed anyway")
	Command.Flags().StringVarP(&localRecipes, "localRecipes", "", "", "a path to local recipes to load instead of the default fetching")
	Command.Flags().StringSliceVarP(&recipePaths, "recipePath", "c", []string{}, "the path to a recipe file to uninstall")

	if err := Command.MarkFlagRequired("recipe"); err != nil {
		log.Error(err)
	}
}
