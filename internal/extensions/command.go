package extensions

import (
	"time"

	"github.com/spf13/cobra"
)

// Command represents the apm command
var Command = &cobra.Command{
	Use:   "extensions",
	Short: "Run an extensions RPC server",
	Run:   runExecutionServer,
}

func runExecutionServer(cmd *cobra.Command, args []string) {
	s := NewServer(cmd.CalledAs(), args)

	timeout := time.NewTimer(30 * time.Second)
	<-timeout.C
	s.Stop()
}
