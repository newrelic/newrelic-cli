package client

import (
	"errors"
	"fmt"

	"github.com/newrelic/newrelic-client-go/newrelic"

	"github.com/newrelic/newrelic-cli/internal/config"
	"github.com/newrelic/newrelic-cli/internal/credentials"
)

var (
	serviceName = "newrelic-cli"
	version     = "dev"
)

// CreateNRClient initializes the New Relic client.
func CreateNRClient(cfg *config.Config, creds *credentials.Credentials) (*newrelic.NewRelic, *credentials.Profile, error) {
	var (
		err               error
		apiKey            string
		insightsInsertKey string
		regionValue       string
	)

	// Create the New Relic Client
	defProfile := creds.Default()

	if defProfile != nil {
		apiKey = defProfile.APIKey
		insightsInsertKey = defProfile.InsightsInsertKey
		regionValue = defProfile.Region
	}

	if apiKey == "" {
		return nil, nil, errors.New("an API key is required, set a default profile or use the NEW_RELIC_API_KEY environment variable")
	}

	userAgent := fmt.Sprintf("newrelic-cli/%s (https://github.com/newrelic/newrelic-cli)", version)

	nrClient, err := newrelic.New(
		newrelic.ConfigPersonalAPIKey(apiKey),
		newrelic.ConfigInsightsInsertKey(insightsInsertKey),
		newrelic.ConfigLogLevel(cfg.LogLevel),
		newrelic.ConfigRegion(regionValue),
		newrelic.ConfigUserAgent(userAgent),
		newrelic.ConfigServiceName(serviceName),
	)

	if err != nil {
		return nil, nil, fmt.Errorf("unable to create New Relic client with error: %s", err)
	}

	return nrClient, defProfile, nil
}
