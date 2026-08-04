package uninstall

import "github.com/newrelic/newrelic-cli/internal/install/types"

// Tier identifies which removal mechanism applies to a recipe.
type Tier int

const (
	// TierNone means the recipe has no automated uninstall path defined.
	TierNone Tier = iota
	// TierTaskfile means the recipe defines its own go-task Uninstall taskfile.
	TierTaskfile
	// TierGeneric means the recipe declares UninstallMeta for the generic remover.
	TierGeneric
)

// SelectTier decides how a recipe should be uninstalled, preferring the
// recipe-authored Taskfile over the generic packages/services/paths metadata.
func SelectTier(r types.OpenInstallationRecipe) Tier {
	if r.Uninstall != "" {
		return TierTaskfile
	}

	m := r.UninstallMeta
	if len(m.Packages) > 0 || len(m.Services) > 0 || len(m.Paths) > 0 {
		return TierGeneric
	}

	return TierNone
}
