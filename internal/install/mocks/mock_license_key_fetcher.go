package mocks

type MockLicenseKeyFetcher struct {
	FetchLicenseKeyFunc func() (string, error)
}

func NewMockLicenseKeyFetcher() *MockLicenseKeyFetcher {
	f := MockLicenseKeyFetcher{}
	f.FetchLicenseKeyFunc = defaultFetchLicenseKeyFunc
	return &f
}

func (f *MockLicenseKeyFetcher) FetchLicenseKey() (string, error) {
	return f.FetchLicenseKeyFunc()
}

func defaultFetchLicenseKeyFunc() (string, error) {
	return "mockLicenseKey", nil
}
