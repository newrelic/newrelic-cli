package discovery

import (
	"context"
	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type Discoverer interface {
	Discover(context.Context) (*types.DiscoveryManifest, error)
}

type Validator interface {
	Validate(m *types.DiscoveryManifest) error
}
