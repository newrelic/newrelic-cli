package decode

import (
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:     "decode",
	Short:   "Decodes NR1 URL Strings ",
	Example: `newrelic decode <subcommand>`,
}
