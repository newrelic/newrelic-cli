package client

import (
	"fmt"

	"github.com/newrelic/newrelic-cli/internal/config"
	"github.com/newrelic/newrelic-cli/internal/credentials"
	"github.com/newrelic/newrelic-client-go/newrelic"
)

// CreateNRClient initializes the New Relic client.
func CreateNRClient(cfg *config.Config, creds *credentials.Credentials) (*newrelic.NewRelic, error) {
	var (
		err            error
		personalAPIKey string
		adminAPIKey    string
		region         string
	)

	// Create the New Relic Client
	defProfile := creds.Default()
	if defProfile != nil {
		adminAPIKey = defProfile.AdminAPIKey
		personalAPIKey = defProfile.PersonalAPIKey
		region = defProfile.Region
	} else {
		return nil, fmt.Errorf("invalid profile name: '%s'", creds.DefaultProfile)
	}

	nrClient, err := newrelic.New(
		newrelic.ConfigAPIKey(adminAPIKey),
		newrelic.ConfigPersonalAPIKey(personalAPIKey),
		newrelic.ConfigLogLevel(cfg.LogLevel),
		newrelic.ConfigRegion(region),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to create New Relic client with error: %s", err)
	}

	return nrClient, nil
}
