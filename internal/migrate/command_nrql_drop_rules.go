package migrate

import (
	"github.com/spf13/cobra"
)

var cmdNRQLDropRules = &cobra.Command{
	Use:   "nrqldroprules",
	Short: "NRQL drop rules migration utilities",
	Long: `Commands for migrating and managing NRQL drop rules during platform transitions.
These utilities help with updating Terraform configurations and validating migrations.`,
}

func init() {
	Command.AddCommand(cmdNRQLDropRules)
}
