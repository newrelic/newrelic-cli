package client

import (
	"context"
	"errors"
	"fmt"
	"sort"

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

const (
	PreferredIngestKeyName = "Installer Ingest License Key"
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

// FetchLicenseKey attempts to fetch and return a customer's license key.
// If the initial request to fetch a license key fails, we retry the request
// for a maximum time duration per config.DefaultMaxTimeoutSeconds.
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

	maxTimeoutSecs := config.DefaultMaxTimeoutSeconds
	if maxTimeoutSeconds != nil {
		maxTimeoutSecs = *maxTimeoutSeconds
	}

	retries := maxTimeoutSecs / config.DefaultPostRetryDelaySec
	retryDelay := config.DefaultPostRetryDelaySec * 1000

	r := utils.NewRetry(retries, retryDelay, retryFunc)
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

	key := getPreferredLicenseKey(keys)
	if key != "" {
		return key, nil
	}

	return "", types.ErrorFetchingLicenseKey
}

// Prefer using the earliest created APIKS license key named "Installer Ingest License Key" if exists.
// Otherwise, fallback to the Account Provisioning "Original account license key"
func getPreferredLicenseKey(keys []apiaccess.APIKey) string {
	key := ""
	if len(keys) > 0 {
		sort.Slice(keys, func(i, j int) bool {
			return keys[i].CreatedAt < keys[j].CreatedAt
		})

		key = keys[0].Key
		for _, k := range keys {
			if k.Name == PreferredIngestKeyName {
				key = k.Key
				break
			}
		}
	}

	return key
}
