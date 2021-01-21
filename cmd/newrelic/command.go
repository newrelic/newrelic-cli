package main

import (
	"errors"
	"os"
	"strconv"

	"github.com/jedib0t/go-pretty/v6/text"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-cli/internal/config"
	"github.com/newrelic/newrelic-cli/internal/output"
	"github.com/newrelic/newrelic-client-go/newrelic"
	"github.com/newrelic/newrelic-client-go/pkg/accounts"
	"github.com/newrelic/newrelic-client-go/pkg/nerdgraph"
)

var (
	outputFormat string
	outputPlain  bool
	debug        bool
	trace        bool
)

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
	initializeLogger()

	// If default profile has not been set, atteempt to initialize it
	if config.GetDefaultProfileName() == "" {
		initializeDefaultProfile()
	}

	// If profile has been overridden, verify it exists
	if config.ProfileOverride == "" {
		if !config.ProfileExists(config.ProfileOverride) {
			log.Fatalf("profile not found: %s", config.ProfileOverride)
		}
	}

	if client.Client == nil {
		client.Client = createClient()
	}
}

func initializeLogger() {
	var logLevel string
	if debug {
		logLevel = "DEBUG"
		config.LogLevelOverride = logLevel
	} else if trace {
		logLevel = "TRACE"
		config.LogLevelOverride = logLevel
	} else {
		logLevel = config.GetLogLevel()
	}

	config.InitLogger(logLevel)
}

func createClient() *newrelic.NewRelic {
	c, err := client.NewClient(config.GetActiveProfileName())
	if err != nil {
		// An error was encountered initializing the client.  This may not be a
		// problem since many commands don't require the use of an initialized client
		log.Debugf("error initializing client: %s", err)
	}

	return c
}

func initializeDefaultProfile() {
	var accountID int
	var region string
	var licenseKey string
	var err error

	userKey := os.Getenv("NEW_RELIC_API_KEY")
	envAccountID := os.Getenv("NEW_RELIC_ACCOUNT_ID")
	region = os.Getenv("NEW_RELIC_REGION")
	licenseKey = os.Getenv("NEW_RELIC_LICENSE_KEY")

	// If we don't have a personal API key we can't initialize a profile.
	if userKey == "" {
		log.Debugf("NEW_RELIC_API_KEY key not set, cannot initialize default profile")
		return
	}

	log.Infof("default profile does not exist and API key detected. attempting to initialize")

	if config.ProfileExists(config.DefaultDefaultProfileName) {
		log.Warnf("a profile named %s already exists, cannot initialize default profile", config.DefaultDefaultProfileName)
		return
	}

	// Saving an initial value will also set the default profile.
	if err = config.SaveValueToProfile(config.DefaultDefaultProfileName, config.UserKey, userKey); err != nil {
		log.Warnf("error saving API key to profile, cannot initialize default profile: %s", err)
		return
	}

	// Default the region to US if it's not in the environment
	if region == "" {
		region = "US"
	}

	if err = config.SaveValueToActiveProfile(config.Region, region); err != nil {
		log.Warnf("couldn't save region to default profile: %s", err)
	}

	// Initialize a client.
	client.Client = createClient()

	// Use the accountID from the environment if we have it.
	if envAccountID != "" {
		accountID, err = strconv.Atoi(envAccountID)
		if err != nil {
			log.Warnf("NEW_RELIC_ACCOUNT_ID has invalid value %s, attempting to fetch account ID", envAccountID)
		}
	}

	// If we still don't have an account ID try to look one up from the API.
	if accountID == 0 {
		accountID, err = fetchAccountID()
		if err != nil {
			log.Warnf("couldn't fetch account ID: %s", err)
		}
	}

	if accountID != 0 {
		if err = config.SaveValueToActiveProfile(config.AccountID, accountID); err != nil {
			log.Warnf("couldn't save account ID to default profile: %s", err)
		}

		if licenseKey == "" {
			log.Infof("attempting to resolve license key for account ID %d", accountID)

			licenseKey, err = fetchLicenseKey(accountID)
			if err != nil {
				log.Warnf("couldn't fetch license key for account ID %d: %s", accountID, err)
			}
		}
	}

	if licenseKey != "" {
		if err = config.SaveValueToActiveProfile(config.LicenseKey, licenseKey); err != nil {
			log.Warnf("couldn't save license key to default profile: %s", err)
		}
	}

	log.Infof("profile %s added", text.FgCyan.Sprint(config.DefaultDefaultProfileName))
}

func fetchLicenseKey(accountID int) (string, error) {
	query := ` query($accountId: Int!) { actor { account(id: $accountId) { licenseKey } } }`

	variables := map[string]interface{}{
		"accountId": accountID,
	}

	resp, err := client.Client.NerdGraph.Query(query, variables)
	if err != nil {
		return "", err
	}

	queryResp := resp.(nerdgraph.QueryResponse)
	actor := queryResp.Actor.(map[string]interface{})
	account := actor["account"].(map[string]interface{})
	licenseKey := account["licenseKey"].(string)

	return licenseKey, nil
}

// fetchAccountID will try and retrieve an account ID for the given user.  If it
// finds more than one account it will returrn an error.
func fetchAccountID() (int, error) {
	params := accounts.ListAccountsParams{
		Scope: &accounts.RegionScopeTypes.IN_REGION,
	}

	accounts, err := client.Client.Accounts.ListAccounts(params)
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
	Command.PersistentFlags().BoolVar(&debug, "debug", false, "debug level logging")
	Command.PersistentFlags().BoolVar(&trace, "trace", false, "trace level logging")
	Command.PersistentFlags().StringVar(&config.ProfileOverride, "profile", "", "the authentication profile to use")
	Command.PersistentFlags().IntVar(&config.AccountIDOverride, "accountId", 0, "the account ID to use for this command")
}

func initConfig() {
	if err := output.SetFormat(output.ParseFormat(outputFormat)); err != nil {
		log.Error(err)
	}

	if err := output.SetPrettyPrint(!outputPlain); err != nil {
		log.Error(err)
	}
}
