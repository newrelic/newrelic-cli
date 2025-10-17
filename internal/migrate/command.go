package migrate

import (
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   "migrate",
	Short: "Commands to support migration of New Relic Resources for EOLs and more.",
}
