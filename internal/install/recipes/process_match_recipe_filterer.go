package recipes

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

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
		log.Tracef("recipe %s not matching any process", r.Name)
	}

	return filtered
}

func (f *ProcessMatchRecipeFilterer) CheckCompatibility(ctx context.Context, r *types.OpenInstallationRecipe, m *types.DiscoveryManifest) error {
	matches := f.processMatchFinder.FindMatches(ctx, m.DiscoveredProcesses, *r)

	// If no `processMatch` is configured in recipe, we can ignore.
	if len(r.ProcessMatch) == 0 {
		return nil
	}

	// Recipe contains a `processMatch` list but no matches were found,
	// the recipe is incompatible with the environment.
	if len(matches) == 0 {
		return fmt.Errorf("recipe %s not matching any process", r.Name)
	}

	return nil
}
