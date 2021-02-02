package workload

import (
	"github.com/spf13/cobra"
)

// Command represents the workloads command.
var Command = &cobra.Command{
	Use:   "workload",
	Short: "Interact with New Relic One workloads",
}
