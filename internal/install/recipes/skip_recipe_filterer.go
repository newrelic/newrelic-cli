package recipes

import (
	"context"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/internal/install/execution"
	"github.com/newrelic/newrelic-cli/internal/install/types"
)

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
			log.Tracef("recipe %s does not match provided names %s", r.Name, f.onlyNames)
			return true
		}
	}

	for _, n := range f.skipNames {
		if strings.EqualFold(strings.TrimSpace(n), strings.TrimSpace(r.Name)) {
			log.Tracef("recipe %s found in skip list %s", r.Name, f.skipNames)
			return true
		}
	}

	for _, k := range f.skipKeywords {
		if r.HasKeyword(k) {
			log.Tracef("recipe %s has keyword %s found in skip list %s", r.Name, k, f.skipKeywords)
			return true
		}
	}

	// Infra should never be skipped based on type
	if r.Name == types.InfraAgentRecipeName {
		return false
	}

	for _, t := range f.skipTypes {
		if r.HasTargetType(types.OpenInstallationTargetType(t)) {
			log.Tracef("recipe %s has type %s found in skip list %s", r.Name, t, f.skipTypes)
			return true
		}
	}

	return false
}

// This is only here to satisfy the interface
func (f *SkipFilterer) CheckCompatibility(ctx context.Context, r *types.OpenInstallationRecipe, m *types.DiscoveryManifest) error {
	return nil
}
