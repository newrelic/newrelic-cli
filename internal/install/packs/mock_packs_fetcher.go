package packs

import (
	"context"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type MockPacksFetcher struct{}

func NewMockPacksFetcher() *MockPacksFetcher {
	return &MockPacksFetcher{}
}

func (f *MockPacksFetcher) FetchPacks(context.Context, []types.OpenInstallationRecipe) ([]types.OpenInstallationObservabilityPack, error) {
	return nil, nil
}
