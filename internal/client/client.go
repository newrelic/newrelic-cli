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

	// licenseKeyAutoFetched is set when FetchLicenseKey finds a key from the account
	// during an install run. Downstream validation can read this via LicenseKeyAutoFetched()
	// to skip the redundant key-list search.
	licenseKeyAutoFetched bool
)

// LicenseKeyAutoFetched reports whether the license key in use was discovered
// automatically by the CLI this run, rather than provided by the user.
func LicenseKeyAutoFetched() bool {
	return licenseKeyAutoFetched
}

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
	licenseKeyAutoFetched = false // reset per-call so multiple invocations behave correctly

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
		key, err = searchForPreferredLicenseKey(utils.SignalCtx, client, accountID)
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

	licenseKeyAutoFetched = true
	return key, nil
}

// FetchAllLicenseKeys returns every LICENSE ingest key for the given account,
// following pagination cursors until all pages are collected.
func FetchAllLicenseKeys(ctx context.Context, client *newrelic.NewRelic, accountID int) ([]apiaccess.APIKey, error) {
	return client.APIAccess.SearchAPIAccessKeysWithContext(ctx, apiaccess.APIAccessKeySearchQuery{
		Scope: apiaccess.APIAccessKeySearchScope{
			AccountIDs:  []int{accountID},
			IngestTypes: []apiaccess.APIAccessIngestKeyType{apiaccess.APIAccessIngestKeyTypeTypes.LICENSE},
		},
		Types: []apiaccess.APIAccessKeyType{apiaccess.APIAccessKeyTypeTypes.INGEST},
	})
}

func searchForPreferredLicenseKey(ctx context.Context, client *newrelic.NewRelic, accountID int) (string, error) {
	log.Debugf("NEW_RELIC_LICENSE_KEY not set; searching for a license key to use for account %d", accountID)

	// Try to find the preferred key by name directly — avoids fetching all keys on
	// accounts with many ingest keys.
	log.Debugf("querying for preferred key %q by name", PreferredIngestKeyName)
	nameKeys, nameErr := client.APIAccess.SearchAPIAccessKeysWithContext(ctx, apiaccess.APIAccessKeySearchQuery{
		Scope: apiaccess.APIAccessKeySearchScope{
			AccountIDs:  []int{accountID},
			IngestTypes: []apiaccess.APIAccessIngestKeyType{apiaccess.APIAccessIngestKeyTypeTypes.LICENSE},
			Name:        &apiaccess.APIAccessKeyNameQuery{Contains: PreferredIngestKeyName},
		},
		Types: []apiaccess.APIAccessKeyType{apiaccess.APIAccessKeyTypeTypes.INGEST},
	})
	if nameErr != nil {
		log.Debugf("name-based search failed (%s); falling back to scanning all LICENSE keys", nameErr)
	} else {
		for _, k := range nameKeys {
			if k.Name == PreferredIngestKeyName {
				log.Debugf("found preferred key %q; using it for the install", PreferredIngestKeyName)
				return k.Key, nil
			}
		}
		log.Debugf("preferred key %q not present; falling back to scanning all LICENSE keys", PreferredIngestKeyName)
	}

	// Preferred key not found or name search failed — fetch all LICENSE keys and
	// select the best available one.
	allKeys, err := FetchAllLicenseKeys(ctx, client, accountID)
	if err != nil {
		return "", err
	}

	key := getPreferredLicenseKey(allKeys)
	if key != "" {
		log.Debugf("selected license key from full scan (%d LICENSE keys returned)", len(allKeys))
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
