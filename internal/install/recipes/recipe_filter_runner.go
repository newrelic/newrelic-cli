package recipes

import (
	"context"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/internal/install/execution"
	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type RecipeFilterer interface {
	Filter(context.Context, *types.OpenInstallationRecipe, *types.DiscoveryManifest) bool
	CheckCompatibility(context.Context, *types.OpenInstallationRecipe, *types.DiscoveryManifest) error
}

type RecipeFilterRunner struct {
	availablilityFilters []RecipeFilterer
	userSkippedFilters   []RecipeFilterer
	installStatus        *execution.InstallStatus
}

func NewRecipeFilterRunner(ic types.InstallerContext, s *execution.InstallStatus) *RecipeFilterRunner {
	skipFilter := NewSkipFilterer(s)
	skipFilter.OnlyNames(ic.RecipeNames...)

	return &RecipeFilterRunner{
		installStatus: s,
		availablilityFilters: []RecipeFilterer{
			NewProcessMatchRecipeFilterer(),
			NewScriptEvaluationRecipeFilterer(s),
		},
		userSkippedFilters: []RecipeFilterer{
			skipFilter,
		},
	}
}

func (rf *RecipeFilterRunner) RunFilter(ctx context.Context, r *types.OpenInstallationRecipe, m *types.DiscoveryManifest) bool {
	for _, f := range rf.availablilityFilters {
		filtered := f.Filter(ctx, r, m)
		if filtered {
			log.Tracef("Filtering out unavailable recipe %s", r.Name)
			return true
		}
	}

	// The DETECTED event must happen before AVAILABLE event
	event := execution.RecipeStatusEvent{
		Recipe: *r,
	}
	rf.installStatus.RecipeDetected(*r, event)

	if r.HasApplicationTargetType() {
		if !r.HasKeyword(types.ApmKeyword) {
			rf.installStatus.RecipeRecommended(execution.RecipeStatusEvent{Recipe: *r})
		}
	}

	rf.installStatus.RecipeAvailable(*r)

	for _, f := range rf.userSkippedFilters {
		filtered := f.Filter(ctx, r, m)

		if filtered {
			log.Tracef("Filtering out skipped recipe %s", r.Name)
			rf.installStatus.RecipeSkipped(execution.RecipeStatusEvent{Recipe: *r})
			return true
		}
	}

	return false
}

func (rf *RecipeFilterRunner) RunFilterAll(ctx context.Context, r []types.OpenInstallationRecipe, m *types.DiscoveryManifest) []types.OpenInstallationRecipe {
	results := []types.OpenInstallationRecipe{}

	for _, recipe := range r {
		filtered := rf.RunFilter(ctx, &recipe, m)

		if !filtered {
			results = append(results, recipe)
		}
	}

	return results
}

func getRecipeFirstName(r types.OpenInstallationRecipe) string {
	if len(r.DisplayName) > 0 {
		parts := strings.Split(r.DisplayName, " ")
		return parts[0]
	}
	return r.DisplayName
}

func (rf *RecipeFilterRunner) EnsureDoesNotFilter(ctx context.Context, r []types.OpenInstallationRecipe, m *types.DiscoveryManifest) error {
	for _, recipe := range r {
		filtered := rf.RunFilter(ctx, &recipe, m)

		if filtered {
			rf.installStatus.RecipeUnsupported(execution.RecipeStatusEvent{Recipe: recipe})
			recipeFirstName := getRecipeFirstName(recipe)
			return fmt.Errorf("we couldn't install the %s. Make sure %s is installed and running on this host and rerun the newrelic-cli command", recipe.DisplayName, recipeFirstName)
		}
	}

	return nil
}

func (rf *RecipeFilterRunner) ConfirmCompatibleRecipes(ctx context.Context, r []types.OpenInstallationRecipe, m *types.DiscoveryManifest) error {
	for _, recipe := range r {
		err := rf.runCompatibilityCheck(ctx, &recipe, m)

		if err != nil {
			recipeStatusEvent := execution.RecipeStatusEvent{
				Recipe: recipe,
				Msg:    err.Error(), // Set the message to the raw error message
			}

			if e, ok := err.(*types.IncomingMessage); ok {
				recipeStatusEvent.Metadata = e.ParseMetadata()
			}

			rf.installStatus.RecipeUnsupported(recipeStatusEvent)
			recipeFirstName := getRecipeFirstName(recipe)

			return fmt.Errorf("we couldn't install the %s. Make sure %s is installed and running on this host and rerun the newrelic-cli command", recipe.DisplayName, recipeFirstName)
		}
	}

	return nil
}

func (rf *RecipeFilterRunner) runCompatibilityCheck(ctx context.Context, r *types.OpenInstallationRecipe, m *types.DiscoveryManifest) error {
	for _, f := range rf.availablilityFilters {
		err := f.CheckCompatibility(ctx, r, m)
		if err != nil {
			log.Tracef("Filtering out unavailable recipe %s", r.Name)
			return err
		}
	}

	// The DETECTED event must happen before AVAILABLE event
	event := execution.RecipeStatusEvent{
		Recipe: *r,
	}
	rf.installStatus.RecipeDetected(*r, event)

	if r.HasApplicationTargetType() {
		if !r.HasKeyword(types.ApmKeyword) {
			rf.installStatus.RecipeRecommended(execution.RecipeStatusEvent{Recipe: *r})
		}
	}

	rf.installStatus.RecipeAvailable(*r)

	for _, f := range rf.userSkippedFilters {
		filtered := f.Filter(ctx, r, m)

		if filtered {
			log.Tracef("Filtering out skipped recipe %s", r.Name)
			rf.installStatus.RecipeSkipped(execution.RecipeStatusEvent{Recipe: *r})
			return fmt.Errorf("recipe %s skipped", r.Name)
		}
	}

	return nil
}
