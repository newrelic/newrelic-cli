package main

import (
	log "github.com/sirupsen/logrus"

	// Commands
	"github.com/newrelic/newrelic-cli/internal/apm"
	"github.com/newrelic/newrelic-cli/internal/config"
	"github.com/newrelic/newrelic-cli/internal/credentials"
	"github.com/newrelic/newrelic-cli/internal/entities"
)

var (
	// AppName for this CMD
	AppName = "newrelic-dev"
	// Version of the CLI
	Version = "dev"
)

func init() {
	// Bind imported sub-commands
	Command.AddCommand(entities.Command)
	Command.AddCommand(credentials.Command)
	Command.AddCommand(apm.Command)
	Command.AddCommand(config.Command)

	Command.AddCommand(completionCmd)
	completionCmd.Flags().StringVar(&completionShell, "shell", "", "Output completion for shell.  (bash, powershell, zsh)")
	completionCmd.MarkFlagRequired("shell")
}

func main() {
	if err := Execute(AppName, Version); err != nil {
		log.Fatal(err)
	}
}
