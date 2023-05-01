package client

import (
	"context"
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/cli"
	"github.com/newrelic/newrelic-cli/internal/config"
	configAPI "github.com/newrelic/newrelic-cli/internal/config/api"
	"github.com/newrelic/newrelic-cli/internal/install/types"
	"github.com/newrelic/newrelic-cli/internal/utils"
	"github.com/newrelic/newrelic-client-go/v2/newrelic"
	"github.com/newrelic/newrelic-client-go/v2/pkg/apiaccess"
	clientLogging "github.com/newrelic/newrelic-client-go/v2/pkg/logging"

	log "github.com/sirupsen/logrus"
)

var (
	NRClient    *newrelic.NewRelic
	serviceName = "newrelic-cli"
)

// NewClient initializes the New Relic client.
func NewClient(profileName string) (*newrelic.NewRelic, error) {
	apiKey := configAPI.GetProfileString(profileName, config.APIKey)
	licenseKey := configAPI.GetProfileString(profileName, config.LicenseKey)

	if apiKey == "" && licenseKey == "" {
		return nil, errors.New("a User API key or License key is required, set a default profile or use the NEW_RELIC_API_KEY or NEW_RELIC_LICENSE_KEY environment variables")
	}

	region := configAPI.GetProfileString(profileName, config.Region)
	userAgent := fmt.Sprintf("newrelic-cli/%s (https://github.com/newrelic/newrelic-cli)", cli.Version())

	// Feed our logrus instance to the client's logrus adapter
	logger := clientLogging.NewLogrusLogger(clientLogging.ConfigLoggerInstance(config.Logger))

	cfgOpts := []newrelic.ConfigOption{
		newrelic.ConfigPersonalAPIKey(apiKey),
		newrelic.ConfigInsightsInsertKey(licenseKey),
		newrelic.ConfigLogger(logger),
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

func FetchLicenseKey(accountID int, profileName string, maxTimeoutSeconds *int) (string, error) {
	var client *newrelic.NewRelic
	var err error
	if profileName == "" {
		client = NRClient
	} else {
		client, err = NewClient(profileName)
		if err != nil {
			return "", err
		}
	}

	var key string
	retryFunc := func() error {
		key, err = execLicenseKeyRequest(utils.SignalCtx, client, accountID)
		if err != nil {
			return err
		}

		return nil
	}

	maxTimeoutSecs := config.DefaultPostMaxTimeoutSecs
	if maxTimeoutSeconds != nil {
		maxTimeoutSecs = *maxTimeoutSeconds
	}

	retries := maxTimeoutSecs / config.DefaultPostRetryDelaySec
	r := utils.NewRetry(retries, (config.DefaultPostRetryDelaySec * 1000), retryFunc)
	retryCtx := r.ExecWithRetries(utils.SignalCtx)

	if !retryCtx.Success {
		return "", retryCtx.MostRecentError()
	}

	return key, nil
}

func execLicenseKeyRequest(ctx context.Context, client *newrelic.NewRelic, accountID int) (string, error) {
	params := apiaccess.APIAccessKeySearchQuery{
		Scope: apiaccess.APIAccessKeySearchScope{
			AccountIDs:  []int{accountID},
			IngestTypes: []apiaccess.APIAccessIngestKeyType{apiaccess.APIAccessIngestKeyTypeTypes.LICENSE},
		},
		Types: []apiaccess.APIAccessKeyType{apiaccess.APIAccessKeyTypeTypes.INGEST},
	}

	keys, err := client.APIAccess.SearchAPIAccessKeysWithContext(ctx, params)
	if err != nil {
		return "", err
	}

	if len(keys) > 0 {
		return keys[0].Key, nil
	}

	return "", types.ErrorFetchingLicenseKey
}
