package extensions

import (
	"github.com/spf13/cobra"
)

// Command represents the apm command
var Command = &cobra.Command{
	Use:   "extensions",
	Short: "Run an extensions RPC server",
	Run:   runExecutionServer,
}

func init() {
}

func runExecutionServer(cmd *cobra.Command, args []string) {
	NewServer()
}
