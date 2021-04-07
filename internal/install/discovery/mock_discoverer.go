package discovery

import (
	"context"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type MockDiscoverer struct {
	DiscoveryManifest *types.DiscoveryManifest
}

func NewMockDiscoverer() *MockDiscoverer {
	m := &types.DiscoveryManifest{
		OS: "linux",
	}

	return &MockDiscoverer{
		DiscoveryManifest: m,
	}
}

func (d *MockDiscoverer) Os(os string) {
	d.DiscoveryManifest.OS = os
}

func (d *MockDiscoverer) Discover(context.Context) (*types.DiscoveryManifest, error) {
	return d.DiscoveryManifest, nil
}
