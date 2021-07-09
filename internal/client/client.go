package client

import (
	"context"
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/config"
	"github.com/newrelic/newrelic-cli/internal/install/types"
	"github.com/newrelic/newrelic-cli/internal/utils"
	"github.com/newrelic/newrelic-client-go/newrelic"
	"github.com/newrelic/newrelic-client-go/pkg/nerdgraph"

	log "github.com/sirupsen/logrus"
)

var (
	NRClient    *newrelic.NewRelic
	serviceName = "newrelic-cli"
	version     = "dev"
)

// NewClient initializes the New Relic client.
func NewClient(profileName string) (*newrelic.NewRelic, error) {
	apiKey := config.GetProfileString(profileName, config.APIKey)
	insightsInsertKey := config.GetProfileString(profileName, config.InsightsInsertKey)

	if apiKey == "" && insightsInsertKey == "" {
		return nil, errors.New("a User API key or Ingest API key is required, set a default profile or use the NEW_RELIC_API_KEY or NEW_RELIC_INSIGHTS_INSERT_KEY environment variables")
	}

	region := config.GetProfileString(profileName, config.Region)
	logLevel := config.GetLogLevelWithFlagOverride()
	userAgent := fmt.Sprintf("newrelic-cli/%s (https://github.com/newrelic/newrelic-cli)", version)

	cfgOpts := []newrelic.ConfigOption{
		newrelic.ConfigPersonalAPIKey(apiKey),
		newrelic.ConfigInsightsInsertKey(insightsInsertKey),
		newrelic.ConfigLogLevel(logLevel),
		newrelic.ConfigRegion(region),
		newrelic.ConfigUserAgent(userAgent),
		newrelic.ConfigServiceName(serviceName),
	}

	nrClient, err := newrelic.New(cfgOpts...)
	if err != nil {
		return nil, fmt.Errorf("unable to create New Relic client with error: %s", err)
	}

	return nrClient, nil
}

func RequireClient(cmd *cobra.Command, args []string) {
	if NRClient == nil {
		log.Fatalf("could not initialize New Relic client, make sure your profile is configured with `newrelic profile configure`")
	}
}

func FetchLicenseKey(accountID int, profileName string) (string, error) {
	client, err := NewClient(profileName)
	if err != nil {
		return "", err
	}

	var key string
	retryFunc := func() error {
		key, err = execLicenseKeyRequest(utils.SignalCtx, client, accountID)
		if err != nil {
			return err
		}

		return nil
	}

	r := utils.NewRetry(3, 1, retryFunc)
	if err := r.ExecWithRetries(utils.SignalCtx); err != nil {
		return "", err
	}

	return key, nil
}

func execLicenseKeyRequest(ctx context.Context, client *newrelic.NewRelic, accountID int) (string, error) {
	query := `query($accountId: Int!) { actor { account(id: $accountId) { licenseKey } } }`

	variables := map[string]interface{}{
		"accountId": accountID,
	}

	resp, err := client.NerdGraph.QueryWithContext(ctx, query, variables)
	if err != nil {
		return "", err
	}

	queryResp := resp.(nerdgraph.QueryResponse)
	actor := queryResp.Actor.(map[string]interface{})
	account := actor["account"].(map[string]interface{})

	if l, ok := account["licenseKey"]; ok {
		if l != nil {
			return l.(string), nil
		}
	}

	return "", types.ErrorFetchingLicenseKey
}

func FetchInsightsInsertKey(accountID int, profileName string) (string, error) {
	client, err := NewClient(profileName)
	if err != nil {
		return "", err
	}

	// Check for an existing key first
	keys, err := client.APIAccess.ListInsightsInsertKeys(accountID)
	if err != nil {
		return "", types.ErrorFetchingInsightsInsertKey
	}

	// We already have a key, return it
	if len(keys) > 0 {
		return keys[0].Key, nil
	}

	// Create a new key if one doesn't exist
	key, err := client.APIAccess.CreateInsightsInsertKey(accountID)
	if err != nil {
		return "", types.ErrorFetchingInsightsInsertKey
	}

	return key.Key, nil
}
