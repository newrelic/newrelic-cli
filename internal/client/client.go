package client

import (
	"errors"
	"fmt"
	"os"

	"github.com/newrelic/newrelic-client-go/newrelic"

	"github.com/newrelic/newrelic-cli/internal/configuration"
)

var (
	serviceName = "newrelic-cli"
	version     = "dev"
)

func NewClient(profileName string) (*newrelic.NewRelic, error) {
	apiKey, err := configuration.GetProfileValue(profileName, configuration.APIKey)
	if err != nil {
		return nil, err
	}

	if apiKey == "" {
		return nil, errors.New("an API key is required, set a default profile or use the NEW_RELIC_API_KEY environment variable")
	}

	region, err := configuration.GetProfileValue(profileName, configuration.Region)
	if err != nil {
		return nil, err
	}

	insightsInsertKey, err := configuration.GetProfileValue(profileName, configuration.InsightsInsertKey)
	if err != nil {
		return nil, err
	}

	logLevel, err := configuration.GetConfigValue(configuration.LogLevel)
	if err != nil {
		return nil, err
	}

	userAgent := fmt.Sprintf("newrelic-cli/%s (https://github.com/newrelic/newrelic-cli)", version)

	cfgOpts := []newrelic.ConfigOption{
		newrelic.ConfigPersonalAPIKey(apiKey.(string)),
		newrelic.ConfigInsightsInsertKey(insightsInsertKey.(string)),
		newrelic.ConfigLogLevel(logLevel.(string)),
		newrelic.ConfigRegion(region.(string)),
		newrelic.ConfigUserAgent(userAgent),
		newrelic.ConfigServiceName(serviceName),
	}

	nerdGraphURLOverride := os.Getenv("NEW_RELIC_NERDGRAPH_URL")
	if nerdGraphURLOverride != "" {
		cfgOpts = append(cfgOpts, newrelic.ConfigNerdGraphBaseURL(nerdGraphURLOverride))
	}

	nrClient, err := newrelic.New(cfgOpts...)
	if err != nil {
		return nil, fmt.Errorf("unable to create New Relic client with error: %s", err)
	}

	return nrClient, nil
}
