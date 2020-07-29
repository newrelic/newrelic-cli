package events

import (
	"github.com/spf13/cobra"
)

// Command represents the events command.
var Command = &cobra.Command{
	Use:   "events",
	Short: "Send custom events to New Relic",
}
