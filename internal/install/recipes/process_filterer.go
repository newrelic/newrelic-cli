package recipes

import (
	"context"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type ProcessFilterer interface {
	Filter(context.Context, []types.GenericProcess, []types.OpenInstallationRecipe) ([]types.MatchedProcess, error)
}
