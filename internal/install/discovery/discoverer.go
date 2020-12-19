package discovery

import (
	"context"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

// Discoverer is reesponsible for discovering informataion about the host system.
type Discoverer interface {
	Discover(context.Context) (*types.DiscoveryManifest, error)
}
