package diagnose

import (
	"os/exec"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var cmdUpdate = &cobra.Command{
	Use:   "update",
	Short: "Update the New Relic Diagnostics binary if necessary",
	Long: `Update the New Relic Diagnostics binary for your system, if it is out of date.

Checks the currently-installed version against the latest version, and if they are different, fetches and installs the latest New Relic Diagnostics build from https://download.newrelic.com/nrdiag.`,
	Example: "newrelic diagnose update",
	Run: func(cmd *cobra.Command, args []string) {
		err := runDiagnostics("-q", "-version")
		if err == nil {
			return
		}
		exitError, ok := err.(*exec.ExitError)
		if !ok || ok && exitError.ProcessState.ExitCode() != 1 {
			// Unexpected error
			logrus.Fatal(err)
		}
		err = downloadBinary()
		if err != nil {
			logrus.Fatal(err)
		}
	},
}

func init() {
	Command.AddCommand(cmdUpdate)
}
