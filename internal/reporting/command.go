package reporting

import (
	"github.com/spf13/cobra"
)

// Command represents the reporting command.
var Command = &cobra.Command{
	Use:   "reporting",
	Short: "Commands for reporting data into New Relic",
}
