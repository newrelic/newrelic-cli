package install

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/internal/client"
	configAPI "github.com/newrelic/newrelic-cli/internal/config/api"
	"github.com/newrelic/newrelic-cli/internal/install/recipes"
)

// relies on the Nerdgraph service
type ServiceLicenseKeyFetcher struct {
	client recipes.NerdGraphClient
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
	accountID := configAPI.GetActiveProfileAccountID()
	licenseKey, err := client.FetchLicenseKey(accountID, configAPI.GetActiveProfileName())
	if err != nil {
		return "", err
	}

	log.Debugf("Found license key %s", licenseKey)
	return licenseKey, nil
}
