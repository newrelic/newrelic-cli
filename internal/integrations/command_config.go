package integrations

import (
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/output"
	"github.com/newrelic/newrelic-cli/internal/utils"
)

var (
	pathConfiguration string
	pathDefinition    string
	pathOutput        string
)

var cmdConfig = &cobra.Command{
	Use:     "config",
	Short:   "Configuration helpers for New Relic onHost integrations",
	Long:    `Configuration helpers for New Relic onHost integrations`,
	Example: "newrelic integrations config migrateV3toV4 --pathDefinition /file/path --pathConfiguration /file/path --pathOutput /file/path",
}

var cmdMigrateV3toV4 = &cobra.Command{
	Use:     "migrateV3toV4",
	Short:   "migrate V3 configuration to V4 configuration format",
	Long:    `migrate V3 configuration to V4 configuration format`,
	Example: "newrelic integrations config migrateV3toV4 --pathDefinition /file/path --pathConfiguration /file/path --pathOutput /file/path",
	Run: func(cmd *cobra.Command, args []string) {

		result := MigrateV3toV4Result{
			MigrateV3toV4Result: migrateV3toV4(pathConfiguration, pathDefinition, pathOutput),
		}

		utils.LogIfFatal(output.Print(result))
	},
}

func init() {

	Command.AddCommand(cmdConfig)

	cmdConfig.AddCommand(cmdMigrateV3toV4)

	cmdMigrateV3toV4.Flags().StringVarP(&pathConfiguration, "pathConfiguration", "c", "", "path configuration")
	cmdMigrateV3toV4.Flags().StringVarP(&pathDefinition, "pathDefinition", "d", "", "path definition")
	cmdMigrateV3toV4.Flags().StringVarP(&pathOutput, "pathOutput", "o", "", "path output")

	utils.LogIfError(cmdMigrateV3toV4.MarkFlagRequired("pathConfiguration"))
	utils.LogIfError(cmdMigrateV3toV4.MarkFlagRequired("pathDefinition"))
	utils.LogIfError(cmdMigrateV3toV4.MarkFlagRequired("pathOutput"))

}
