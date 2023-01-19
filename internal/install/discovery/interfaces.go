package discovery

import (
	"context"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

// ManifestDiscoverer is responsible for discovering information about the host system.
type ManifestDiscoverer interface {
	DiscoverManifest(context.Context, Discoverer) (*types.DiscoveryManifest, error)
	ValidateManifest(*types.DiscoveryManifest, Validator) error
}

type Discoverer interface {
	Discover(context.Context) (*types.DiscoveryManifest, error)
}

type Validator interface {
	Validate(m *types.DiscoveryManifest) error
}
