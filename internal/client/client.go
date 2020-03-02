package client

import (
	"errors"
	"fmt"
	"os"

	"github.com/newrelic/newrelic-cli/internal/config"
	"github.com/newrelic/newrelic-cli/internal/credentials"
	"github.com/newrelic/newrelic-cli/internal/version"
	"github.com/newrelic/newrelic-client-go/newrelic"
)

// CreateNRClient initializes the New Relic client.
func CreateNRClient(cfg *config.Config, creds *credentials.Credentials) (*newrelic.NewRelic, error) {
	var (
		err                error
		personalAPIKey     string
		region             string
		defaultServiceName = "newrelic-cli"
	)

	// Create the New Relic Client
	defProfile := creds.Default()

	defProfile = applyOverrides(defProfile)

	if defProfile != nil {
		apiKey = defProfile.APIKey
		region = defProfile.Region
	}

	if apiKey == "" {
		return nil, errors.New("an API key is required, set a default profile or use the NEW_RELIC_API_KEY environment variable")
	}

	userAgent := fmt.Sprintf("newrelic/%s/%s (https://github.com/newrelic/%s)", defaultServiceName, version.Version, defaultServiceName)

	nrClient, err := newrelic.New(
		newrelic.ConfigPersonalAPIKey(apiKey),
		newrelic.ConfigLogLevel(cfg.LogLevel),
		newrelic.ConfigRegion(region),
		newrelic.ConfigUserAgent(userAgent),
		newrelic.ConfigServiceName(defaultServiceName),
	)

	if err != nil {
		return nil, fmt.Errorf("unable to create New Relic client with error: %s", err)
	}

	return nrClient, nil
}

// applyOverrides reads Profile info out of the Environment to override config
func applyOverrides(p *credentials.Profile) *credentials.Profile {
	envAPIKey := os.Getenv("NEW_RELIC_API_KEY")
	envRegion := os.Getenv("NEW_RELIC_REGION")

	if envAPIKey == "" && envRegion == "" {
		return p
	}

	var out credentials.Profile

	if p == nil {
		out = credentials.Profile{}
	} else {
		out = *p
	}

	if envAPIKey != "" {
		out.APIKey = envAPIKey
	}

	if envRegion != "" {
		out.Region = envRegion
	}

	return &out
}
