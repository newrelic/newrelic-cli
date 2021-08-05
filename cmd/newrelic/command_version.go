package main

import (
	"fmt"

	"github.com/newrelic/newrelic-cli/internal/cli"
	"github.com/spf13/cobra"
)

var cmdVersion = &cobra.Command{
	Use:   "version",
	Short: "Show the version of the New Relic CLI",
	Long: `Use the version command to print out the version of this command.
`,
	Example: "newrelic version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("newrelic version %s\n", cli.Version)
	},
}

func init() {
	Command.AddCommand(cmdVersion)
}
