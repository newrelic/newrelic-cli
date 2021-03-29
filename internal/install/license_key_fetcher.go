package install

import (
	"context"
	"fmt"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/internal/install/recipes"
	"github.com/newrelic/newrelic-cli/internal/credentials"
)

// relies on the Nerdgraph service
type LicenseKeyFetcher struct {
	client recipes.NerdGraphClient
}

// NewLicenseKeyFetcher returns a new instance of LicenseKeyFetcher.
func NewLicenseKeyFetcher(client recipes.NerdGraphClient) *LicenseKeyFetcher {
	f := LicenseKeyFetcher{
		client: client,
	}

	return &f
}

func (f *LicenseKeyFetcher) FetchLicenseKey(ctx context.Context) (string, error) {
	var resp licenseKeyDataQueryResult

	vars := map[string]interface{}{
	}

	defaultProfile := credentials.DefaultProfile()

	query := fmt.Sprintf(licenseKeyQuery, strconv.Itoa(defaultProfile.AccountID))

	if err := f.client.QueryWithResponseAndContext(ctx, query, vars, &resp); err != nil {
		log.Fatal(err)
		return "", err
	}

	licenseKey := resp.Actor.Account.LicenseKey
	log.Debugf("Found license key %s", licenseKey)
	return licenseKey, nil
}

type licenseKeyDataQueryResult struct {
	Actor licenseKeyActorQueryResult `json:"docs"`
}

type licenseKeyActorQueryResult struct {
	Account licenseKeyAccountQueryResult `json:"docs"`
}

type licenseKeyAccountQueryResult struct {
	LicenseKey string `json:"licenseKey"`
}

const (
	licenseKeyQuery = `
	query{
		actor {
			account(id: %i) {
				licenseKey
			}
		}
	}`
)
