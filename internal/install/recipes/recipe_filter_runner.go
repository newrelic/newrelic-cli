package recipes

import (
	"context"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/internal/install/execution"
	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type RecipeFilterRunner struct {
	availablilityFilters []RecipeFilterer
	userSkippedFilters   []RecipeFilterer
	installStatus        *execution.InstallStatus
}

func NewRecipeFilterRunner(ic types.InstallerContext, s *execution.InstallStatus) *RecipeFilterRunner {
	skipFilter := NewSkipFilterer(s)
	skipFilter.SkipNames(ic.SkipNames()...)
	skipFilter.SkipTypes(ic.SkipTypes()...)
	skipFilter.SkipKeywords(ic.SkipKeywords()...)
	skipFilter.OnlyNames(ic.RecipeNames...)

	return &RecipeFilterRunner{
		installStatus: s,
		availablilityFilters: []RecipeFilterer{
			NewProcessMatchRecipeFilterer(),
			NewScriptEvaluationRecipeFilterer(),
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
			log.Debugf("Filtering out unavailable recipe %s", r.Name)
			return true
		}
	}

	if r.HasApplicationTargetType() {
		if !r.HasKeyword(types.ApmKeyword) {
			rf.installStatus.RecipeRecommended(execution.RecipeStatusEvent{Recipe: *r})
		}
	} else {
		rf.installStatus.RecipeAvailable(*r)
	}

	for _, f := range rf.userSkippedFilters {
		filtered := f.Filter(ctx, r, m)

		if filtered {
			log.Debugf("Filtering out skipped recipe %s", r.Name)
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
			return fmt.Errorf("we couldn’t install the %s. Make sure %s is installed and running on this host and rerun the newrelic-cli command", recipe.DisplayName, recipeFirstName)
		}
	}

	return nil
}

type RecipeFilterer interface {
	Filter(context.Context, *types.OpenInstallationRecipe, *types.DiscoveryManifest) bool
}

type ProcessMatchRecipeFilterer struct {
	processMatchFinder ProcessMatchFinder
}

func NewProcessMatchRecipeFilterer() *ProcessMatchRecipeFilterer {
	return &ProcessMatchRecipeFilterer{
		processMatchFinder: NewRegexProcessMatchFinder(),
	}
}

func (f *ProcessMatchRecipeFilterer) Filter(ctx context.Context, r *types.OpenInstallationRecipe, m *types.DiscoveryManifest) bool {
	matches := f.processMatchFinder.FindMatches(ctx, m.DiscoveredProcesses, *r)

	filtered := len(r.ProcessMatch) > 0 && len(matches) == 0

	if filtered {
		log.Debugf("recipe %s not matching any process", r.Name)
	}

	return filtered
}

type ScriptEvaluationRecipeFilterer struct {
	recipeExecutor execution.RecipeExecutor
}

func NewScriptEvaluationRecipeFilterer() *ScriptEvaluationRecipeFilterer {
	recipeExecutor := execution.NewShRecipeExecutor()

	return &ScriptEvaluationRecipeFilterer{
		recipeExecutor: recipeExecutor,
	}
}

func (f *ScriptEvaluationRecipeFilterer) Filter(ctx context.Context, r *types.OpenInstallationRecipe, m *types.DiscoveryManifest) bool {
	if err := f.recipeExecutor.ExecutePreInstall(ctx, *r, types.RecipeVars{}); err != nil {
		log.Debugf("recipe %s failed script evaluation %s", r.Name, err)
		return true
	}

	return false
}

type SkipFilterer struct {
	*execution.InstallStatus
	skipNames    []string
	skipTypes    []string
	skipKeywords []string
	onlyNames    []string
}

func NewSkipFilterer(s *execution.InstallStatus) *SkipFilterer {
	return &SkipFilterer{
		InstallStatus: s,
		skipNames:     []string{},
		skipTypes:     []string{},
		skipKeywords:  []string{},
		onlyNames:     []string{},
	}
}

func (f *SkipFilterer) SkipNames(names ...string) {
	f.skipNames = append(f.skipNames, names...)
}

func (f *SkipFilterer) OnlyNames(names ...string) {
	f.onlyNames = append(f.onlyNames, names...)
}

func (f *SkipFilterer) SkipTypes(types ...string) {
	f.skipTypes = append(f.skipNames, types...)
}

func (f *SkipFilterer) SkipKeywords(keywords ...string) {
	f.skipKeywords = append(f.skipNames, keywords...)
}

func (f *SkipFilterer) Filter(ctx context.Context, r *types.OpenInstallationRecipe, m *types.DiscoveryManifest) bool {
	if len(f.onlyNames) > 0 {
		filtered := true
		for _, n := range f.onlyNames {
			if strings.EqualFold(strings.TrimSpace(n), strings.TrimSpace(r.Name)) {
				filtered = false
			}
		}

		if filtered {
			log.Debugf("recipe %s does not match provided names %s", r.Name, f.onlyNames)
			return true
		}
	}

	for _, n := range f.skipNames {
		if strings.EqualFold(strings.TrimSpace(n), strings.TrimSpace(r.Name)) {
			log.Debugf("recipe %s found in skip list %s", r.Name, f.skipNames)
			return true
		}
	}

	for _, k := range f.skipKeywords {
		if r.HasKeyword(k) {
			log.Debugf("recipe %s has keyword %s found in skip list %s", r.Name, k, f.skipKeywords)
			return true
		}
	}

	// Infra should never be skipped based on type
	if r.Name == types.InfraAgentRecipeName {
		return false
	}

	for _, t := range f.skipTypes {
		if r.HasTargetType(types.OpenInstallationTargetType(t)) {
			log.Debugf("recipe %s has type %s found in skip list %s", r.Name, t, f.skipTypes)
			return true
		}
	}

	return false
}
