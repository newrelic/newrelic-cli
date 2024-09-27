package agent

import (
	"errors"
	"fmt"
	"strings"

	"github.com/newrelic/newrelic-cli/internal/agent/migrate"
	"github.com/newrelic/newrelic-cli/internal/agent/obfuscate"
	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-cli/internal/output"
	"github.com/newrelic/newrelic-cli/internal/utils"
	ng "github.com/newrelic/newrelic-client-go/v2/pkg/nerdgraph"
	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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

const (
	ANDROID        = "ANDROID"
	BROWSER        = "BROWSER"
	DOTNET         = "DOTNET"
	ELIXIR         = "ELIXIR"
	GO             = "GO"
	INFRASTRUCTURE = "INFRASTRUCTURE"
	IOS            = "IOS"
	JAVA           = "JAVA"
	NODEJS         = "NODEJS"
	PHP            = "PHP"
	PYTHON         = "PYTHON"
	RUBY           = "RUBY"
	SDK            = "SDK"
)

func newAgentNameList() []string {
	return []string{
		ANDROID,
		BROWSER,
		DOTNET,
		ELIXIR,
		GO,
		INFRASTRUCTURE,
		IOS,
		JAVA,
		NODEJS,
		PHP,
		PYTHON,
		RUBY,
		SDK,
	}
}

func isValidAgentName(agentName string) bool {
	for _, a := range newAgentNameList() {
		if a == agentName {
			return true
		}
	}

	return false
}

func agentNameTitleCase(agentName string) string {
	caser := cases.Title(language.AmericanEnglish)

	switch agentName {
	case DOTNET:
		return ".NET"
	case IOS:
		return "iOS"
	case NODEJS:
		return "Node.js"
	case SDK:
		return "SDK"
	default:
		return caser.String(agentName)
	}
}

var cmdAgentVersion = &cobra.Command{
	Use:   "version",
	Short: "Show latest agent versions.",
	Long: `Show latest agent versions. Valid agent names include:
android, browser, dotnet, elixir, go, infrastructure, ios, java, nodejs, php, python, ruby, sdk"
`,
	Example: "newrelic agent version <agent_name>",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("invalid number of arguments")
		}

		agentName := strings.ToUpper(args[0])

		if !isValidAgentName(agentName) {
			return fmt.Errorf("invalid agent name: %s, use --help for a list of valid agent names", args[0])
		}

		return nil
	},
	PreRun: client.RequireClient,
	RunE: func(cmd *cobra.Command, args []string) error {
		agentName := strings.ToUpper(args[0])

		query := `
		query CurrentAgentRelease ($agentName: AgentReleasesFilter!) {
			docs {
				currentAgentRelease(agentName: $agentName) {
				    version
				}
			}
		}`

		variables := map[string]interface{}{
			"agentName": agentName,
		}

		result, err := client.NRClient.NerdGraph.QueryWithContext(utils.SignalCtx, query, variables)

		if err != nil {
			return err
		}

		queryResp := result.(ng.QueryResponse)
		release := queryResp.Docs.(map[string]interface{})["currentAgentRelease"]
		version := release.(map[string]interface{})["version"].(string)

		agentNameTitleCase := agentNameTitleCase(agentName)

		fmt.Printf("%s: %s\n", agentNameTitleCase, version)

		return nil
	},
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

	Command.AddCommand(cmdAgentVersion)

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
