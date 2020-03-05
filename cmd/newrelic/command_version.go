package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the version of this tool",
	Long: `Use the version command to print out the version of this command.
`,
	Example: "newrelic version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("newrelic version %s\n", version)
	},
}
