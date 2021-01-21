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
	userKey := config.GetProfileValueString(profileName, config.UserKey)
	insightsInsertKey := config.GetProfileValueString(profileName, config.InsightsInsertKey)
	if userKey == "" && insightsInsertKey == "" {
		return nil, errors.New("a User API key or Ingest API key is required, set a default profile or use the NEW_RELIC_API_KEY or NEW_RELIC_INSIGHTS_INSERT_KEY environment variables")
	}

	region := config.GetProfileValueString(profileName, config.Region)
	logLevel := config.GetLogLevel()
	userAgent := fmt.Sprintf("newrelic-cli/%s (https://github.com/newrelic/newrelic-cli)", version)

	cfgOpts := []newrelic.ConfigOption{
		newrelic.ConfigPersonalAPIKey(userKey),
		newrelic.ConfigInsightsInsertKey(insightsInsertKey),
		newrelic.ConfigLogLevel(logLevel),
		newrelic.ConfigRegion(region),
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
