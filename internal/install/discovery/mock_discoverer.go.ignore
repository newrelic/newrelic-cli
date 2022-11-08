package discovery

import (
	"context"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type MockDiscoverer struct {
	DiscoveryManifest *types.DiscoveryManifest
	Error             error
}

func NewMockDiscoverer() *MockDiscoverer {
	m := &types.DiscoveryManifest{
		OS:              "linux",
		Platform:        "",
		PlatformVersion: "",
		KernelArch:      "",
	}

	return &MockDiscoverer{
		DiscoveryManifest: m,
	}
}

func (d *MockDiscoverer) SetOs(os string) {
	d.DiscoveryManifest.OS = os
}

func (d *MockDiscoverer) SetPlatform(p string) {
	d.DiscoveryManifest.Platform = p
}

func (d *MockDiscoverer) SetPlatformVersion(pf string) {
	d.DiscoveryManifest.PlatformVersion = pf
}

func (d *MockDiscoverer) SetKernelArch(arch string) {
	d.DiscoveryManifest.KernelArch = arch
}

func (d *MockDiscoverer) GetManifest() *types.DiscoveryManifest {
	return d.DiscoveryManifest
}

func (d *MockDiscoverer) Discover(context.Context) (*types.DiscoveryManifest, error) {
	return d.DiscoveryManifest, d.Error
}
