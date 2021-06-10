package packs

import (
	"context"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type MockPacksInstaller struct{}

func NewMockPacksInstaller() *MockPacksInstaller {
	return &MockPacksInstaller{}
}

func (f *MockPacksInstaller) Install(ctx context.Context, packs []types.OpenInstallationObservabilityPack) error {
	return nil
}
