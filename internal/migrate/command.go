package migrate

import (
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   "migrate",
	Short: "Commands for interacting with the New Relic Database",
}
