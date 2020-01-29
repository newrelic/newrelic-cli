package cmd

import (
	"github.com/spf13/cobra"

	nr "github.com/newrelic/newrelic-client-go/newrelic"
	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/pkg/config"
	"github.com/newrelic/newrelic-cli/pkg/profile"
)

// LogLevel passsed into the CLI
var configFile string
var logLevel string

var globalConfig *config.Config
var profiles *map[string]profile.Profile

// Client is an instance of the New Relic client.
var Client *nr.NewRelic

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:     "newrelic-dev",
	Short:   "The New Relic CLI",
	Long:    `The New Relic CLI enables users to perform tasks against the New Relic APIs`,
	Version: "dev",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute(appName, version string) error {
	if appName != "" {
		RootCmd.Use = appName
	}
	if version != "" {
		RootCmd.Version = version
	}

	return RootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file (default is "+config.DefaultConfigFile()+")")
	RootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "", "log level [Panic,Fatal,Error,Warn,Info,Debug,Trace]")
}

// initConfig will be run ONLY after a valid command is called
func initConfig() {
	var (
		err    error
		apiKey string
		region string
	)

	if logLevel != "" {
		lvl, err := log.ParseLevel(logLevel)
		if err == nil {
			log.SetLevel(lvl)
		}
	}

	globalConfig, err = config.Load(configFile, logLevel)
	if err != nil {
		log.Fatalf("unable to load config with error: %s\n", err)
	}

	// Load profiles
	profiles, err = profile.Load(globalConfig)
	if err != nil {
		// TODO: Don't die here, we need to be able to run the profiles command to add one!
		log.Fatalf("unable to load profiles with error: %s\n", err)
	}

	log.Tracef("config: %+v\n", globalConfig)
	log.Tracef("profiles: %+v\n", profiles)

	if globalConfig == nil {
		log.Fatal("configuration required")
	}
	if profiles == nil {
		log.Fatal("at least one profile is required")
	}

	// Create the New Relic Client
	if val, ok := (*profiles)[globalConfig.ProfileName]; ok {
		apiKey = val.PersonalAPIKey
		region = val.Region
	} else {
		log.Fatalf("invalid profile name: '%s'", globalConfig.ProfileName)
	}

	Client, err = nr.New(nr.ConfigPersonalAPIKey(apiKey), nr.ConfigLogLevel(globalConfig.LogLevel), nr.ConfigRegion(region))
	if err != nil {
		log.Fatalf("unable to create New Relic client with error: %s", err)
	}
}
