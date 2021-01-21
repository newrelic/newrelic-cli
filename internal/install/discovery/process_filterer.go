package discovery

import (
	"context"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type ProcessFilterer interface {
	filter(context.Context, []types.GenericProcess) ([]types.ProcessInfoWrap, error)
}
