package discovery

import (
	"context"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

// FileFilterer determines the existence of files on the underlying filesystem.
type FileFilterer interface {
	Filter(context.Context, []types.OpenInstallationRecipe) ([]types.OpenInstallationLogMatch, error)
}
