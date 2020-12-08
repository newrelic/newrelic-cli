package install

import "context"

type mockDiscoverer struct {
	discoveryManifest *discoveryManifest
}

func newMockDiscoverer() *mockDiscoverer {
	m := &discoveryManifest{}

	return &mockDiscoverer{
		discoveryManifest: m,
	}
}

func (d *mockDiscoverer) discover(context.Context) (*discoveryManifest, error) {
	return d.discoveryManifest, nil
}
