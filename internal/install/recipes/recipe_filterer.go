package recipes

import (
	"context"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/internal/install/execution"
	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type RecipeFilterer struct {
	availablilityFilters []RecipeFilter
	userSkippedFilters   []RecipeFilter
	installStatus        *execution.InstallStatus
}

func NewRecipeFilterer(ic types.InstallerContext, s *execution.InstallStatus) *RecipeFilterer {
	skipFilter := NewSkipFilter(s)
	skipFilter.SkipNames(ic.SkipNames()...)
	skipFilter.SkipKeywords(ic.SkipKeywords()...)

	return &RecipeFilterer{
		installStatus: s,
		availablilityFilters: []RecipeFilter{
			NewProcessMatchRecipeFilter(),
			NewScriptEvaluationRecipeFilter(),
		},
		userSkippedFilters: []RecipeFilter{
			skipFilter,
		},
	}
}

func (rf *RecipeFilterer) Filter(ctx context.Context, r *types.OpenInstallationRecipe, m *types.DiscoveryManifest) (bool, error) {
	for _, f := range rf.availablilityFilters {
		filtered, err := f.Execute(ctx, r, m)
		if err != nil {
			return false, err
		}

		if filtered {
			log.Debugf("Filtering out unavailable recipe %s", r.Name)
			return true, nil
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
		filtered, err := f.Execute(ctx, r, m)
		if err != nil {
			return false, err
		}

		if filtered {
			log.Debugf("Filtering out skipped recipe %s", r.Name)
			rf.installStatus.RecipeSkipped(execution.RecipeStatusEvent{Recipe: *r})
			return true, nil
		}
	}

	return false, nil
}

func (rf *RecipeFilterer) FilterMultiple(ctx context.Context, r []types.OpenInstallationRecipe, m *types.DiscoveryManifest) ([]types.OpenInstallationRecipe, error) {
	results := []types.OpenInstallationRecipe{}
	for _, recipe := range r {
		filtered, err := rf.Filter(ctx, &recipe, m)
		if err != nil {
			return nil, err
		}

		if !filtered {
			results = append(results, recipe)
		}
	}

	return results, nil
}

type RecipeFilter interface {
	Execute(context.Context, *types.OpenInstallationRecipe, *types.DiscoveryManifest) (bool, error)
}

type ProcessMatchRecipeFilter struct {
	processMatchFinder ProcessMatchFinder
}

func NewProcessMatchRecipeFilter() *ProcessMatchRecipeFilter {
	return &ProcessMatchRecipeFilter{
		processMatchFinder: NewRegexProcessMatchFinder(),
	}
}

func (f *ProcessMatchRecipeFilter) Execute(ctx context.Context, r *types.OpenInstallationRecipe, m *types.DiscoveryManifest) (bool, error) {
	matches, err := f.processMatchFinder.FindMatches(ctx, m.DiscoveredProcesses, *r)
	if err != nil {
		return false, err
	}

	filtered := len(r.ProcessMatch) > 0 && len(matches) == 0

	if filtered {
		log.Debugf("recipe %s failed process match", r.Name)
	}

	return filtered, nil
}

type ScriptEvaluationRecipeFilter struct {
	recipeExecutor execution.RecipeExecutor
}

func NewScriptEvaluationRecipeFilter() *ScriptEvaluationRecipeFilter {
	return &ScriptEvaluationRecipeFilter{
		recipeExecutor: execution.NewShRecipeExecutor(),
	}
}

func (f *ScriptEvaluationRecipeFilter) Execute(ctx context.Context, r *types.OpenInstallationRecipe, m *types.DiscoveryManifest) (bool, error) {
	if err := f.recipeExecutor.ExecuteDiscovery(ctx, *r, types.RecipeVars{}); err != nil {
		log.Debugf("recipe %s failed script evaluation", r.Name)
		return true, nil
	}

	return false, nil
}

type SkipFilter struct {
	*execution.InstallStatus
	skipNames    []string
	skipKeywords []string
}

func NewSkipFilter(s *execution.InstallStatus) *SkipFilter {
	return &SkipFilter{
		InstallStatus: s,
		skipNames:     []string{},
		skipKeywords:  []string{},
	}
}

func (f *SkipFilter) SkipNames(names ...string) {
	f.skipNames = append(f.skipNames, names...)
}

func (f *SkipFilter) SkipKeywords(keywords ...string) {
	f.skipKeywords = append(f.skipNames, keywords...)
}

func (f *SkipFilter) Execute(ctx context.Context, r *types.OpenInstallationRecipe, m *types.DiscoveryManifest) (bool, error) {
	for _, n := range f.skipNames {
		if strings.EqualFold(strings.TrimSpace(n), strings.TrimSpace(r.Name)) {
			return true, nil
		}
	}

	for _, n := range f.skipKeywords {
		for _, k := range r.Keywords {
			if strings.EqualFold(strings.TrimSpace(n), strings.TrimSpace(k)) {
				return true, nil
			}
		}
	}

	return false, nil
}
