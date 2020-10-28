package diagnose

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var cmdRun = &cobra.Command{
	Use:   "run",
	Short: "Troubleshoot your New Relic-instrumented application",
	Long: `Troubleshoot your New Relic-instrumented application

The diagnose command runs New Relic Diagnostics, our troubleshooting suite. The first time you run this command the nrdiag binary appropriate for your system will be downloaded to .newrelic/bin in your home directory.\n
`,
	Example: "\tnewrelic diagnose run --suites java,infra",
	Run: func(cmd *cobra.Command, args []string) {
		nrdiagArgs := make([]string, 0)
		if options.listSuites {
			nrdiagArgs = append(nrdiagArgs, "-help", "suites")
		} else if options.suites != "" {
			nrdiagArgs = append(nrdiagArgs, "-suites", options.suites)
		}
		err := runDiagnostics(nrdiagArgs...)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	Command.AddCommand(cmdRun)
	cmdRun.Flags().StringVar(&options.attachmentKey, "attachment-key", "", "Attachment key for automatic upload to a support ticket (get key from an existing ticket).")
	cmdRun.Flags().BoolVar(&options.verbose, "verbose", false, "Display verbose logging during task execution.")
	cmdRun.Flags().StringVar(&options.suites, "suites", "", "The task suite or comma-separated list of suites to run. Use --list-suites for a list of available suites.")
	cmdRun.Flags().BoolVar(&options.listSuites, "list-suites", false, "List the task suites available for the --suites argument.")
}
