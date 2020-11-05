package diagnose

import (
	"github.com/spf13/cobra"
)

var options struct {
	suites        string
	listSuites    bool
	verbose       bool
	attachmentKey string
	configFile    string
}

// Command represents the diagnose command.
var Command = &cobra.Command{
	Use:   "diagnose",
	Short: "Troubleshoot your New Relic installation",
}
