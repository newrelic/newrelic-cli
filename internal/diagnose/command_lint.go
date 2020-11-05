package diagnose

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var cmdLint = &cobra.Command{
	Use:   "lint",
	Short: "Validate your agent config file",
	Long: `Validate your agent config file settings. Currently only available for the Java agent.

Checks the settings in the specified Java agent config file, making sure they have the correct type and structure.`,
	Example: "\tnewrelic diagnose lint --config-file ./newrelic.yml",
	Run: func(cmd *cobra.Command, args []string) {
		err := runDiagnostics("-t", "Java/Config/ValidateSettings", "-c", options.configFile)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	Command.AddCommand(cmdLint)
	cmdLint.Flags().StringVar(&options.configFile, "config-file", "", "Path to the config file to be validated.")
}
