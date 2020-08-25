package apm

import (
	"encoding/base64"

	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/output"
	"github.com/newrelic/newrelic-cli/internal/utils"
)

var (
	encodeKey    string
	textToEncode string
)

var cmdObfuscate = &cobra.Command{
	Use:   "obfuscate",
	Short: "Obfuscate a configuration value using a key",
	Long: `Obfuscate a configuration value using a key.  Obfuscated value
is placed in the Agent configuration or environment variable." 
`,
	Example: "newrelic agent config obfuscate --value <config_value> --key <obfuscation_key>",
	Run: func(cmd *cobra.Command, args []string) {

		encoding := base64.NewEncoding(encodeKey)

		textToEncodeBytes := []byte(textToEncode)

		encodedString := encoding.EncodeToString(textToEncodeBytes)

		utils.LogIfFatal(output.Print(encodedString))
	},
}

func init() {

	cmdConfig.AddCommand(cmdObfuscate)

	cmdObfuscate.Flags().StringVarP(&encodeKey, "kye", "", "", "the key to use when obfuscating the clear-text value")
	cmdObfuscate.Flags().StringVarP(&textToEncode, "value", "", "", "the value, in clear text, to be obfuscated")

	utils.LogIfError(cmdObfuscate.MarkFlagRequired("key"))
	utils.LogIfError(cmdObfuscate.MarkFlagRequired("value"))
}
