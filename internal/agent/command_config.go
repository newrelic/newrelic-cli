package agent

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/output"
)

var (
	encodeKey    string
	textToEncode string
)

var cmdConfig = &cobra.Command{
	Use:   "config",
	Short: "Configuration utilities/helpers for New Relic agents",
}

var cmdConfigObfuscate = &cobra.Command{
	Use:   "obfuscate",
	Short: "Obfuscate a configuration value using a key",
	Long: `Obfuscate a configuration value using a key.  The obfuscated value
should be placed in the Agent configuration or in an environment variable." 
`,
	Example: "newrelic agent config obfuscate --value <config_value> --key <obfuscation_key>",
	Run: func(cmd *cobra.Command, args []string) {

		result := ObfuscationResult{
			ObfuscatedValue: obfuscateStringWithKey(textToEncode, encodeKey),
		}

		if err := output.Print(result); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {

	Command.AddCommand(cmdConfig)

	cmdConfig.AddCommand(cmdConfigObfuscate)

	cmdConfigObfuscate.Flags().StringVarP(&encodeKey, "key", "k", "", "the key to use when obfuscating the clear-text value")
	cmdConfigObfuscate.Flags().StringVarP(&textToEncode, "value", "v", "", "the value, in clear text, to be obfuscated")

	if err := cmdConfigObfuscate.MarkFlagRequired("key"); err != nil {
		log.Error(err)
	}
	if err := cmdConfigObfuscate.MarkFlagRequired("value"); err != nil {
		log.Error(err)
	}
}
