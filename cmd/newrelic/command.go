package main

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// Command represents the base command when called without any subcommands
var Command = &cobra.Command{
	Use:     "newrelic-dev",
	Short:   "The New Relic CLI",
	Long:    `The New Relic CLI enables users to perform tasks against the New Relic APIs`,
	Version: "dev",
}

var (
	completionShell string
)

var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: "Generates shell completion functions",
	Long: `Generate shell completion functions

To load completion run the following.

. <(newrelic completion --shell bash)

To configure your shell to load the completions on start, include the following in your shell's rc file.

Using bash, for example.

# ~/.bashrc or ~/.profile
. <(newrelic completion --shell bash)


Using zsh, for example.

# ~/.zshrc
. <(newrelic completion --shell zsh)
`,
	Example: "newrelic completion --shell zsh",
	Run: func(cmd *cobra.Command, args []string) {

		switch shell := completionShell; shell {
		case "bash":
			err := Command.GenBashCompletion(os.Stdout)
			if err != nil {
				log.Error(err)
			}
		case "powershell":
			err := Command.GenPowerShellCompletion(os.Stdout)
			if err != nil {
				log.Error(err)
			}
		case "zsh":
			err := Command.GenZshCompletion(os.Stdout)
			if err != nil {
				log.Error(err)
			}
		default:
			log.Error("--shell must be one of [bash, powershell, zsh]")
		}
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the version of this tool",
	Long: `Use the version command to print out the version of this command.
`,
	Example: "newrelic version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("newrelic version %s\n", Version)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute(appName, version string) error {
	if appName != "" {
		Command.Use = appName
	}
	if version != "" {
		Command.Version = version
	}

	return Command.Execute()
}
