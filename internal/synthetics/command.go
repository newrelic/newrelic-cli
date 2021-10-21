package synthetics

import (
	"github.com/spf13/cobra"
)

// Command represents the synthetics command
var Command = &cobra.Command{
	Use:   "synthetics",
	Short: "Interact with New Relic Synthetics",
}
