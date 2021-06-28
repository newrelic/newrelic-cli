package packs

import (
	"context"

	"github.com/newrelic/newrelic-cli/internal/install/execution"
	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type MockPacksInstaller struct {
	InstallCallCount int
	InstallErr       error
	installStatus    *execution.InstallStatus
}

func NewMockPacksInstaller(s *execution.InstallStatus) *MockPacksInstaller {
	return &MockPacksInstaller{
		installStatus: s,
	}
}

func (f *MockPacksInstaller) Install(ctx context.Context, packs []types.OpenInstallationObservabilityPack) error {
	for _, pack := range packs {
		f.installStatus.ObservabilityPackInstallPending(execution.ObservabilityPackStatusEvent{
			ObservabilityPack: pack,
		})
		f.installStatus.ObservabilityPackInstallSuccess(execution.ObservabilityPackStatusEvent{
			ObservabilityPack: pack,
		})
	}

	f.InstallCallCount++
	return f.InstallErr
}
