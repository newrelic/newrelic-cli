package discovery

import (
	"context"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type MockDiscoverer struct {
	DiscoveryManifest *types.DiscoveryManifest
}

func NewMockDiscoverer() *MockDiscoverer {
	m := &types.DiscoveryManifest{}

	return &MockDiscoverer{
		DiscoveryManifest: m,
	}
}

func (d *MockDiscoverer) Discover(context.Context) (*types.DiscoveryManifest, error) {
	return d.DiscoveryManifest, nil
}
