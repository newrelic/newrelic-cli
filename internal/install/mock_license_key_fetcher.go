package install

import (
	"context"

	"github.com/newrelic/newrelic-cli/internal/config"
)

type MockLicenseKeyFetcher struct {
	FetchLicenseKeyFunc func(ctx context.Context) (string, error)
	maxTimeoutSeconds   int
}

func NewMockLicenseKeyFetcher() *MockLicenseKeyFetcher {
	f := MockLicenseKeyFetcher{
		maxTimeoutSeconds: config.DefaultPostMaxTimeoutSecs,
	}
	f.FetchLicenseKeyFunc = defaultFetchLicenseKeyFunc
	return &f
}

func (f *MockLicenseKeyFetcher) FetchLicenseKey(ctx context.Context) (string, error) {
	return f.FetchLicenseKeyFunc(ctx)
}

func defaultFetchLicenseKeyFunc(ctx context.Context) (string, error) {
	return "mockLicenseKey", nil
}
