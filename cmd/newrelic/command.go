package main

import (
	"os"
	"strconv"

	"github.com/jedib0t/go-pretty/v6/text"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-cli/internal/credentials"
	"github.com/newrelic/newrelic-cli/internal/output"
	"github.com/newrelic/newrelic-cli/internal/utils"
	"github.com/newrelic/newrelic-client-go/newrelic"
	"github.com/newrelic/newrelic-client-go/pkg/accounts"
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
	// TODO what do we want?

	// Create a license key from nerdgraph using the creds we have.  This may
	// require an accountID.  Set this in the profile.

	// Trim down the required number of elements in the environment as much as
	// possible.

	// When an accountID is required, we should try to retrieve it from
	// nerdgraph, in place of requiring the user to set it.

	// Determine the accounts this user has access to.  In the case we have only
	// one, then we have our answer about which account ID to use.
	// { actor { accounts(scope: IN_REGION) { id name } } }

	// Default to US region?
	var accountID int
	var err error

	credentials.WithCredentials(func(c *credentials.Credentials) {
		if c.DefaultProfile != "" {
			log.Infof("default profile already exists")
			return
		}

		envAPIKey := os.Getenv("NEW_RELIC_API_KEY")
		envRegion := os.Getenv("NEW_RELIC_REGION")
		envAccountID := os.Getenv("NEW_RELIC_ACCOUNT_ID")

		hasDefault := func() bool {
			for profileName := range c.Profiles {
				if profileName == defaultProfileName {
					return true
				}
			}

			return false
		}

		if envAPIKey != "" && envRegion != "" {

			// TODO: DRY this up (exists as well in credentials/command.go)

			// Use the accountID from the environment if we have it.
			if envAccountID != "" {
				accountID, err = strconv.Atoi(envAccountID)
				if err != nil {
					log.Errorf("error reading accountID from environment: %s", err)
				}
			}

			// If we still don't have an account ID, try to look one up from the API.
			if accountID == 0 {
				client.WithClient(func(nrClient *newrelic.NewRelic) {

					params := accounts.ListAccountsParams{
						Scope: &accounts.RegionScopeTypes.IN_REGION,
					}

					resp, err := nrClient.Accounts.ListAccounts(params)
					if err != nil {
						log.Fatal("error retrieving accounts:", err)
					}

					if len(resp) == 1 {
						// Set the accountID for the profile
						accountID = resp[0].ID
					} else {
						log.Warnf("more than one account found")
					}
				})
			}

			if !hasDefault() {
				err := c.AddProfile(defaultProfileName, envRegion, envAPIKey, "", accountID)
				if err != nil {
					log.Fatal(err)
				}

				log.Infof("profile %s added", text.FgCyan.Sprint(defaultProfileName))
			}

			if len(c.Profiles) == 1 {
				err := c.SetDefaultProfile(defaultProfileName)
				if err != nil {
					log.Fatal(err)
				}

				log.Infof("setting %s as default profile", text.FgCyan.Sprint(defaultProfileName))
			}
		}
	})
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
