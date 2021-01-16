package client

import (
	"errors"
	"fmt"
	"os"

	"github.com/newrelic/newrelic-cli/internal/config"
	"github.com/newrelic/newrelic-client-go/newrelic"
)

var (
	Client      *newrelic.NewRelic
	serviceName = "newrelic-cli"
	version     = "dev"
)

func NewClient(profileName string) (*newrelic.NewRelic, error) {
	apiKey := config.GetProfileValueString(profileName, config.APIKey)
	if apiKey == "" {
		return nil, errors.New("an API key is required, set a default profile or use the NEW_RELIC_API_KEY environment variable")
	}

	region := config.GetProfileValueString(profileName, config.Region)
	insightsInsertKey := config.GetProfileValueString(profileName, config.InsightsInsertKey)

	logLevel := config.GetConfigValueString(config.LogLevel)
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
