package install

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/internal/client"
	configAPI "github.com/newrelic/newrelic-cli/internal/config/api"
	"github.com/newrelic/newrelic-cli/internal/install/recipes"
	"github.com/newrelic/newrelic-cli/internal/utils"
)

// relies on the Nerdgraph service
type ServiceLicenseKeyFetcher struct {
	client     recipes.NerdGraphClient
	LicenseKey string
}

type LicenseKeyFetcher interface {
	FetchLicenseKey(context.Context) (string, error)
}

func NewServiceLicenseKeyFetcher(client recipes.NerdGraphClient) LicenseKeyFetcher {
	f := ServiceLicenseKeyFetcher{
		client: client,
	}

	return &f
}

func (f *ServiceLicenseKeyFetcher) FetchLicenseKey(ctx context.Context) (string, error) {

	if f.LicenseKey != "" {
		return f.LicenseKey, nil
	}

	accountID := configAPI.GetActiveProfileAccountID()
	licenseKey, err := client.FetchLicenseKey(accountID, configAPI.GetActiveProfileName())
	if err != nil {
		return "", err
	}

	log.Debugf("Found license key %s", utils.Obfuscate(licenseKey))
	f.LicenseKey = licenseKey
	return licenseKey, nil
}
