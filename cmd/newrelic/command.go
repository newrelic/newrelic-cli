package main

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/text"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/credentials"
	"github.com/newrelic/newrelic-cli/internal/output"
	"github.com/newrelic/newrelic-cli/internal/utils"
)

var outputFormat string
var outputPlain bool

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
	credentials.WithCredentials(func (c *credentials.Credentials) {
		if c.DefaultProfile != "" {
			log.Infof("default profile already exists")
			return
		}

		envAPIKey := os.Getenv("NEW_RELIC_API_KEY")
		envRegion := os.Getenv("NEW_RELIC_REGION")

		if envAPIKey != "" && envRegion != "" {
			profileName := "default"

			// TODO: DRY this up (exists as well in credentials/command.go)
			err := c.AddProfile(profileName, envRegion, envAPIKey, "")
			if err != nil {
				log.Fatal(err)
			}

			log.Infof("profile %s added", text.FgCyan.Sprint(profileName))

			if len(c.Profiles) == 1 {
				err := c.SetDefaultProfile(profileName)
				if err != nil {
					log.Fatal(err)
				}

				log.Infof("setting %s as default profile", text.FgCyan.Sprint(profileName))
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
