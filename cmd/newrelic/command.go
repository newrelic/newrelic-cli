package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/jedib0t/go-pretty/v6/text"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-cli/internal/credentials"
	"github.com/newrelic/newrelic-cli/internal/install/types"
	"github.com/newrelic/newrelic-cli/internal/output"
	"github.com/newrelic/newrelic-cli/internal/utils"
	"github.com/newrelic/newrelic-client-go/newrelic"
	"github.com/newrelic/newrelic-client-go/pkg/accounts"
	"github.com/newrelic/newrelic-client-go/pkg/nerdgraph"
)

var outputFormat string
var outputPlain bool

const defaultProfileName string = "default"

// Command represents the base command when called without any subcommands
var Command = &cobra.Command{
	PersistentPreRun:  initializeCLI,
	Use:               appName,
	Short:             "The New Relic CLI",
	Long:              `The New Relic CLI enables users to perform tasks against the New Relic APIs`,
	Version:           version,
	DisableAutoGenTag: true, // Do not print generation date on documentation
}

func initializeCLI(cmd *cobra.Command, args []string) {
	initializeProfile()
}

func initializeProfile() {
	var accountID int
	var region string
	var licenseKey string
	var insightsInsertKey string
	var err error

	credentials.WithCredentials(func(c *credentials.Credentials) {
		if c.DefaultProfile != "" {
			err = errors.New("default profile already exists, not attempting to initialize")
			return
		}

		apiKey := os.Getenv("NEW_RELIC_API_KEY")
		envAccountID := os.Getenv("NEW_RELIC_ACCOUNT_ID")
		region = os.Getenv("NEW_RELIC_REGION")
		licenseKey = os.Getenv("NEW_RELIC_LICENSE_KEY")
		insightsInsertKey = os.Getenv("NEW_RELIC_INSIGHTS_INSERT_KEY")

		// If we don't have a personal API key we can't initialize a profile.
		if apiKey == "" {
			err = errors.New("api key not provided, not attempting to initialize default profile")
			return
		}

		// Default the region to US if it's not in the environment
		if region == "" {
			region = "US"
		}

		// Use the accountID from the environment if we have it.
		if envAccountID != "" {
			accountID, err = strconv.Atoi(envAccountID)
			if err != nil {
				err = fmt.Errorf("couldn't parse account ID: %s", err)
				return
			}
		}

		// We should have an API key by this point, initialize the client.
		client.WithClient(func(nrClient *newrelic.NewRelic) {
			// If we still don't have an account ID try to look one up from the API.
			if accountID == 0 {
				accountID, err = fetchAccountID(nrClient)
				if err != nil {
					return
				}
			}

			if licenseKey == "" {
				// We should have an account ID by now, so fetch the license key for it.
				licenseKey, err = fetchLicenseKey(nrClient, accountID)
				if err != nil {
					log.Error(err)
					return
				}
			}

			if insightsInsertKey == "" {
				// We should have an API key by now, so fetch the insights insert key for it.
				insightsInsertKey, err = fetchInsightsInsertKey(nrClient, accountID)
				if err != nil {
					log.Error(err)
				}
			}

			if !hasProfileWithDefaultName(c.Profiles) {
				p := credentials.Profile{
					Region:            region,
					APIKey:            apiKey,
					AccountID:         accountID,
					LicenseKey:        licenseKey,
					InsightsInsertKey: insightsInsertKey,
				}

				err = c.AddProfile(defaultProfileName, p)
				if err != nil {
					return
				}

				log.Infof("profile %s added", text.FgCyan.Sprint(defaultProfileName))
			}

			if len(c.Profiles) == 1 {
				err = c.SetDefaultProfile(defaultProfileName)
				if err != nil {
					err = fmt.Errorf("error setting %s as the default profile: %s", text.FgCyan.Sprint(defaultProfileName), err)
					return
				}

				log.Infof("setting %s as default profile", text.FgCyan.Sprint(defaultProfileName))
			}
		})
	})

	if err != nil {
		log.Debugf("couldn't initialize default profile: %s", err)
	}
}

func hasProfileWithDefaultName(profiles map[string]credentials.Profile) bool {
	for profileName := range profiles {
		if profileName == defaultProfileName {
			return true
		}
	}

	return false
}

func fetchLicenseKey(client *newrelic.NewRelic, accountID int) (string, error) {
	query := `query($accountId: Int!) { actor { account(id: $accountId) { licenseKey } } }`

	variables := map[string]interface{}{
		"accountId": accountID,
	}

	for i := 0; i < 3; i++ {
		resp, err := client.NerdGraph.Query(query, variables)
		if err != nil {
			return "", err
		}

		queryResp := resp.(nerdgraph.QueryResponse)
		actor := queryResp.Actor.(map[string]interface{})
		account := actor["account"].(map[string]interface{})

		if licenseKey, ok := account["licenseKey"]; ok {
			return licenseKey.(string), nil
		}

		time.Sleep(1 * time.Second)
	}

	return "", types.ErrorFetchingLicenseKey
}

func fetchInsightsInsertKey(client *newrelic.NewRelic, accountID int) (string, error) {
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

// fetchAccountID will try and retrieve an account ID for the given user.  If it
// finds more than one account it will returrn an error.
func fetchAccountID(client *newrelic.NewRelic) (int, error) {
	params := accounts.ListAccountsParams{
		Scope: &accounts.RegionScopeTypes.IN_REGION,
	}

	accounts, err := client.Accounts.ListAccounts(params)
	if err != nil {
		return 0, err
	}

	if len(accounts) == 1 {
		return accounts[0].ID, nil
	}

	return 0, errors.New("multiple accounts found, please set NEW_RELIC_ACCOUNT_ID")
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() error {
	Command.Use = appName
	Command.Version = version
	Command.SilenceUsage = os.Getenv("CI") != ""

	// Silence Cobra's internal handling of error messaging
	// since we have a custom error handler in main.go
	Command.SilenceErrors = true

	return Command.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	Command.PersistentFlags().StringVar(&outputFormat, "format", output.DefaultFormat.String(), "output text format ["+output.FormatOptions()+"]")
	Command.PersistentFlags().BoolVar(&outputPlain, "plain", false, "output compact text")
}

func initConfig() {
	utils.LogIfError(output.SetFormat(output.ParseFormat(outputFormat)))
	utils.LogIfError(output.SetPrettyPrint(!outputPlain))
}
