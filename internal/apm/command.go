package apm

import (
	"github.com/spf13/cobra"
)

var (
	apmAccountID string
	apmAppID     int
)

// Command represents the apm command
var Command = &cobra.Command{
	Use:   "apm",
	Short: "Interact with New Relic APM",
}

func init() {
	// Flags for all things APM
	Command.PersistentFlags().StringVarP(&apmAccountID, "accountId", "a", "", "A New Relic account ID")
	Command.PersistentFlags().IntVarP(&apmAppID, "applicationId", "", 0, "A New Relic APM application ID")
}
