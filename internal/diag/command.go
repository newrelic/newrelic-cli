package diag

import (
	"github.com/spf13/cobra"
)

// Command represents the diag command.
var Command = &cobra.Command{
	Use:   "diag",
	Short: "Troubleshoot your New Relic installation", // FIXME: something better
}
