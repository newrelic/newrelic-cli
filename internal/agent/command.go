package agent

import (
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/agent/migrate"
	"github.com/newrelic/newrelic-cli/internal/agent/obfuscate"
	"github.com/newrelic/newrelic-cli/internal/output"
	"github.com/newrelic/newrelic-cli/internal/utils"
)

var (
	// cmdConfigObfuscate
	encodeKey    string
	textToEncode string

	// cmdMigrateV3toV4
	pathConfiguration string
	pathDefinition    string
	pathOutput        string
	overwrite         bool
)

// Command represents the agent command
var Command = &cobra.Command{
	Use:   "agent",
	Short: "Utilities for New Relic Agents",
	Long:  `Utilities for New Relic Agents`,
}

var cmdConfig = &cobra.Command{
	Use:     "config",
	Short:   "Configuration utilities/helpers for New Relic agents",
	Long:    "Configuration utilities/helpers for New Relic agents",
	Example: "newrelic agent config obfuscate --value <config_value> --key <obfuscation_key>",
}

var cmdConfigObfuscate = &cobra.Command{
	Use:   "obfuscate",
	Short: "Obfuscate a configuration value using a key",
	Long: `Obfuscate a configuration value using a key.  The obfuscated value
should be placed in the Agent configuration or in an environment variable." 
`,
	Example: "newrelic agent config obfuscate --value <config_value> --key <obfuscation_key>",
	Run: func(cmd *cobra.Command, args []string) {

		result := obfuscate.Result{
			Value: obfuscate.StringWithKey(textToEncode, encodeKey),
		}

		utils.LogIfFatal(output.Print(result))
	},
}

var cmdMigrateV3toV4 = &cobra.Command{
	Use:     "migrateV3toV4",
	Short:   "migrate V3 configuration to V4 configuration format",
	Long:    `migrate V3 configuration to V4 configuration format`,
	Example: "newrelic integrations config migrateV3toV4 --pathDefinition /file/path --pathConfiguration /file/path --pathOutput /file/path",
	Run: func(cmd *cobra.Command, args []string) {

		result := migrate.V3toV4Result{
			V3toV4Result: migrate.V3toV4(pathConfiguration, pathDefinition, pathOutput, overwrite),
		}

		utils.LogIfFatal(output.Print(result))
	},
}

func init() {

	Command.AddCommand(cmdConfig)

	cmdConfig.AddCommand(cmdConfigObfuscate)

	cmdConfigObfuscate.Flags().StringVarP(&encodeKey, "key", "k", "", "the key to use when obfuscating the clear-text value")
	cmdConfigObfuscate.Flags().StringVarP(&textToEncode, "value", "v", "", "the value, in clear text, to be obfuscated")

	utils.LogIfError(cmdConfigObfuscate.MarkFlagRequired("key"))
	utils.LogIfError(cmdConfigObfuscate.MarkFlagRequired("value"))

	cmdConfig.AddCommand(cmdMigrateV3toV4)

	cmdMigrateV3toV4.Flags().StringVarP(&pathConfiguration, "pathConfiguration", "c", "", "path configuration")
	cmdMigrateV3toV4.Flags().StringVarP(&pathDefinition, "pathDefinition", "d", "", "path definition")
	cmdMigrateV3toV4.Flags().StringVarP(&pathOutput, "pathOutput", "o", "", "path output")
	cmdMigrateV3toV4.Flags().BoolVar(&overwrite, "overwrite", false, "if set ti true and pathOutput file exists already the old file is removed ")

	utils.LogIfError(cmdMigrateV3toV4.MarkFlagRequired("pathConfiguration"))
	utils.LogIfError(cmdMigrateV3toV4.MarkFlagRequired("pathDefinition"))
	utils.LogIfError(cmdMigrateV3toV4.MarkFlagRequired("pathOutput"))
}
