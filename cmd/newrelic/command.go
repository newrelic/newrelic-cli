package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/cli"
	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-cli/internal/config"
	configAPI "github.com/newrelic/newrelic-cli/internal/config/api"
	diagnose "github.com/newrelic/newrelic-cli/internal/diagnose"
	"github.com/newrelic/newrelic-cli/internal/output"
	"github.com/newrelic/newrelic-cli/internal/split"
	"github.com/newrelic/newrelic-cli/internal/utils"
	"github.com/newrelic/newrelic-client-go/newrelic"
	nrErrors "github.com/newrelic/newrelic-client-go/pkg/errors"
)

var (
	outputFormat string
	outputPlain  bool
)

// Command represents the base command when called without any subcommands
var Command = &cobra.Command{
	PersistentPreRun:  initializeCLI,
	Use:               appName,
	Short:             "The New Relic CLI",
	Long:              `The New Relic CLI enables users to perform tasks against the New Relic APIs`,
	Version:           cli.Version(),
	DisableAutoGenTag: true, // Do not print generation date on documentation
}

func initializeCLI(cmd *cobra.Command, args []string) {
	// Initialize logger
	logLevel := configAPI.GetLogLevel()
	config.InitLogger(log.StandardLogger(), logLevel)

	// Initialize feature flag service
	split.Init()

	// Initialize client
	if client.NRClient == nil {
		client.NRClient = createClient()
	}
}

func createClient() *newrelic.NewRelic {
	c, err := client.NewClient(configAPI.GetActiveProfileName())
	if err != nil {
		// An error was encountered initializing the client.  This may not be a
		// problem since many commands don't require the use of an initialized client
		log.Debugf("error initializing client: %s", err)
	}

	return c
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() error {
	Command.Use = appName
	Command.Version = cli.Version()

	// Silence Cobra's internal handling of command usage help text.
	// Note, the help text is still displayed if any command arg or
	// flag validation fails.
	Command.SilenceUsage = true

	// Silence Cobra's internal handling of error messaging
	// since we have a custom error handler in main.go
	Command.SilenceErrors = true

	err := Command.Execute()
	if _, ok := err.(*nrErrors.PaymentRequiredError); ok {
		diagnose.PrintPaymentRequiredErrorMessage()
		log.Debug(err)
		return nil
	}

	return err
}

func init() {
	cobra.OnInitialize(initConfig)

	Command.PersistentFlags().StringVar(&outputFormat, "format", output.DefaultFormat.String(), "output text format ["+output.FormatOptions()+"]")
	Command.PersistentFlags().StringVar(&config.FlagProfileName, "profile", "", "the authentication profile to use")
	Command.PersistentFlags().BoolVar(&outputPlain, "plain", false, "output compact text")
	Command.PersistentFlags().BoolVar(&config.FlagDebug, "debug", false, "debug level logging")
	Command.PersistentFlags().BoolVar(&config.FlagTrace, "trace", false, "trace level logging")
	Command.PersistentFlags().IntVarP(&config.FlagAccountID, "accountId", "a", 0, "the account ID to use. Can be overridden by setting NEW_RELIC_ACCOUNT_ID")
}

func initConfig() {
	utils.LogIfError(output.SetFormat(output.ParseFormat(outputFormat)))
	utils.LogIfError(output.SetPrettyPrint(!outputPlain))
}
