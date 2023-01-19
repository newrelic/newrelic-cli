package install

import (
	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/internal/client"
	configAPI "github.com/newrelic/newrelic-cli/internal/config/api"
	"github.com/newrelic/newrelic-cli/internal/utils"
)

// relies on the Nerdgraph service
// FIXME not sure if this really relies on NerdGraph service

var (
	getActiveProfileAccountID = configAPI.GetActiveProfileAccountID
	getActiveProfileName      = configAPI.GetActiveProfileName
	fetchLicenseKey           = client.FetchLicenseKey
)

type LicenseKeyFetcher interface {
	FetchLicenseKey() (string, error)
}

type ServiceLicenseKeyFetcher struct {
	LicenseKey string
}

func NewServiceLicenseKeyFetcher() *ServiceLicenseKeyFetcher {
	return &ServiceLicenseKeyFetcher{}
}

func (f *ServiceLicenseKeyFetcher) FetchLicenseKey() (string, error) {

	if f.LicenseKey != "" {
		return f.LicenseKey, nil
	}

	accountID := getActiveProfileAccountID()
	licenseKey, err := fetchLicenseKey(accountID, getActiveProfileName())
	if err != nil {
		return "", err
	}

	log.Debugf("Found license key %s", utils.Obfuscate(licenseKey))
	f.LicenseKey = licenseKey
	return licenseKey, nil
}
