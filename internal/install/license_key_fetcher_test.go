package install

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFetchLicenseKeyReturnsSavedKey(t *testing.T) {
	expectedLicenseKey := "12345somelicensekey"
	licenseKeyFetcher := NewServiceLicenseKeyFetcher()
	licenseKeyFetcher.LicenseKey = expectedLicenseKey

	actualLicenseKey, err := licenseKeyFetcher.FetchLicenseKey()

	assert.NoError(t, err)
	assert.Equal(t, expectedLicenseKey, actualLicenseKey)
}

func TestFetchLicenseKeyReturnsErrorIfFetchErrors(t *testing.T) {
	realFetchLicenseKey := fetchLicenseKey
	defer func() {
		fetchLicenseKey = realFetchLicenseKey
	}()
	fetchLicenseKey = func(int, string) (string, error) {
		return "", errors.New("couldn't fetch key")
	}

	licenseKey, err := NewServiceLicenseKeyFetcher().FetchLicenseKey()

	assert.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "couldn't fetch key"))
	assert.True(t, 0 == len(licenseKey))
}

func TestFetchLicenseKeyFetchesAndStoresKey(t *testing.T) {
	realFetchLicenseKey := fetchLicenseKey
	defer func() {
		fetchLicenseKey = realFetchLicenseKey
	}()
	fetchLicenseKey = func(int, string) (string, error) {
		return "hey-ima-key", nil
	}

	licenseKeyFetcher := NewServiceLicenseKeyFetcher()
	actualLicenseKey, err := licenseKeyFetcher.FetchLicenseKey()

	assert.NoError(t, err)
	assert.Equal(t, "hey-ima-key", actualLicenseKey)
	assert.Equal(t, "hey-ima-key", licenseKeyFetcher.LicenseKey)
}
