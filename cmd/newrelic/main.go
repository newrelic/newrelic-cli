package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	// Commands
	"github.com/newrelic/newrelic-cli/internal/apm"
	"github.com/newrelic/newrelic-cli/internal/config"
	"github.com/newrelic/newrelic-cli/internal/credentials"
	"github.com/newrelic/newrelic-cli/internal/entities"
	"github.com/newrelic/newrelic-cli/internal/nerdgraph"

	"github.com/newrelic/newrelic-cli/internal/utils"
)

var (
	appName = "newrelic-dev"
	version = "dev"
)

func init() {
	Command.AddCommand(versionCmd)

	Command.AddCommand(completionCmd)
	completionCmd.Flags().StringVar(&completionShell, "shell", "", "Output completion for the specified shell.  (bash, powershell, zsh)")
	utils.LogIfError(completionCmd.MarkFlagRequired("shell"))

	// Bind imported sub-commands
	Command.AddCommand(entities.Command)
	Command.AddCommand(credentials.Command)
	Command.AddCommand(apm.Command)
	Command.AddCommand(config.Command)
	Command.AddCommand(nerdgraph.Command)
}

func main() {
	if err := Execute(); err != nil {
		if err != cobra.ErrSubCommandRequired {
			log.Fatal(err)
		}
	}
}
