package integrations

import (
	"github.com/spf13/cobra"
)

// Command represents the agent command
var Command = &cobra.Command{
	Use:   "integrations",
	Short: "Utilities for New Relic Agent's onHost integrations",
	Long:  `Utilities for New Relic Agent's onHost integrations`,
}
