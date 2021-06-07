package recipes

import (
	"context"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type ProcessMatchFinder interface {
	FindMatchesMultiple(context.Context, []types.GenericProcess, []types.OpenInstallationRecipe) ([]types.MatchedProcess, error)
	FindMatches(context.Context, []types.GenericProcess, types.OpenInstallationRecipe) ([]types.MatchedProcess, error)
}
