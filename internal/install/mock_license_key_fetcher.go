package install

import (
	"context"
)

type MockLicenseKeyFetcher struct {
	FetchLicenseKeyFunc func(ctx context.Context) (string, error)
}

func NewMockLicenseKeyFetcher() *MockLicenseKeyFetcher {
	f := MockLicenseKeyFetcher{}
	f.FetchLicenseKeyFunc = defaultFetchLicenseKeyFunc
	return &f
}

func (f *MockLicenseKeyFetcher) FetchLicenseKey(ctx context.Context) (string, error) {
	return f.FetchLicenseKeyFunc(ctx)
}

func defaultFetchLicenseKeyFunc(ctx context.Context) (string, error) {
	return "mockLicenseKey", nil
}
