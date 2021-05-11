package discovery

import (
	"context"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type NoOpProcessFilterer struct{}

func NewNoOpProcessFilterer() *NoOpProcessFilterer {
	return &NoOpProcessFilterer{}
}

func (f *NoOpProcessFilterer) filter(ctx context.Context, processes []types.GenericProcess, manifest types.DiscoveryManifest) ([]types.MatchedProcess, error) {
	return []types.MatchedProcess{}, nil
}
