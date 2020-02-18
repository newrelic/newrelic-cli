package main

import (
	"github.com/spf13/cobra"
)

// Command represents the base command when called without any subcommands
var Command = &cobra.Command{
	Use:     "newrelic-dev",
	Short:   "The New Relic CLI",
	Long:    `The New Relic CLI enables users to perform tasks against the New Relic APIs`,
	Version: "dev",
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
