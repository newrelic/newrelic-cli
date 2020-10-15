package diagnose

import (
	"github.com/spf13/cobra"
)

// Command represents the diagnose command.
var Command = &cobra.Command{
	Use:   "diagnose",
	Short: "Troubleshoot your New Relic installation",
	Run: func(cmd *cobra.Command, args []string) {
		cmdDiag.Run(cmd, args)
	},
}
