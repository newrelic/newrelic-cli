package packs

import (
	"context"

	"github.com/newrelic/newrelic-cli/internal/install/execution"
	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type MockPacksFetcher struct {
	FetchPacksCallCount int
	FetchPacksErr       error
	FetchPacksVal       []types.OpenInstallationObservabilityPack
	installStatus       *execution.InstallStatus
}

func NewMockPacksFetcher(s *execution.InstallStatus) *MockPacksFetcher {
	return &MockPacksFetcher{
		installStatus: s,
	}
}

func (f *MockPacksFetcher) FetchPacks(ctx context.Context, recipes []types.OpenInstallationRecipe) ([]types.OpenInstallationObservabilityPack, error) {
	for _, r := range recipes {
		for _, pack := range r.ObservabilityPacks {
			f.installStatus.ObservabilityPackFetchPending(execution.ObservabilityPackStatusEvent{Name: pack.Name})
			f.installStatus.ObservabilityPackFetchSuccess(execution.ObservabilityPackStatusEvent{
				ObservabilityPack: types.OpenInstallationObservabilityPack{
					Name: pack.Name,
				},
			})
		}
	}

	f.FetchPacksCallCount++
	f.FetchPacksVal = []types.OpenInstallationObservabilityPack{
		{
			Name: "test-pack",
		},
	}
	return f.FetchPacksVal, f.FetchPacksErr
}
