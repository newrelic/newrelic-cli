package install

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/internal/config"
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
	var resp licenseKeyDataQueryResult

	vars := map[string]interface{}{}

	accountID := configAPI.GetActiveProfileString(config.AccountID)

	query := `
	query{
		actor {
			account(id: ` + accountID + `) {
				licenseKey
			}
		}
	}`

	if err := f.client.QueryWithResponseAndContext(ctx, query, vars, &resp); err != nil {
		return "", err
	}

	licenseKey := resp.Actor.Account.LicenseKey
	log.Debugf("Found license key %s", licenseKey)
	return licenseKey, nil
}

type licenseKeyDataQueryResult struct {
	Actor licenseKeyActorQueryResult `json:"actor"`
}

type licenseKeyActorQueryResult struct {
	Account licenseKeyAccountQueryResult `json:"account"`
}

type licenseKeyAccountQueryResult struct {
	LicenseKey string `json:"licenseKey"`
}
