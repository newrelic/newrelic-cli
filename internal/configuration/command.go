package configuration

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/output"
	"github.com/newrelic/newrelic-cli/internal/utils"
	log "github.com/sirupsen/logrus"
)

var (
	// Display keys when printing output
	key   string
	value string
)

// Command is the base command for managing profiles
var Command = &cobra.Command{
	Use:   "config",
	Short: "Manage the configuration of the New Relic CLI",
}

var cmdSet = &cobra.Command{
	Use:   "set",
	Short: "Set a configuration value",
	Long: `Set a configuration value

The set command sets a persistent configuration value for the New Relic CLI.
`,
	Example: "newrelic config set --key <key> --value <value>",
	Run: func(cmd *cobra.Command, args []string) {
		if !isValidConfigKey() {
			log.Fatalf("%s is not a valid config field. valid values are %s", key, GetValidConfigKeys())
		}

		if err := SetConfigString(ConfigKey(key), value); err != nil {
			log.Fatal(err)
		}

		log.Info("success")
	},
}

func isValidConfigKey() (valid bool) {
	VisitAllConfigFields(func(fd FieldDefinition) {
		if strings.EqualFold(string(fd.Key), key) {
			valid = true
		}
	})

	return valid
}

var cmdGet = &cobra.Command{
	Use:   "get",
	Short: "Get a configuration value",
	Long: `Get a configuration value

The get command gets a persistent configuration value for the New Relic CLI.
`,
	Example: "newrelic config get --key <key>",
	Run: func(cmd *cobra.Command, args []string) {
		if !isValidConfigKey() {
			log.Fatalf("%s is not a valid config field. valid values are %s", key, GetValidConfigKeys())
		}

		output.Text(GetConfigString(ConfigKey(key)))
	},
}

var cmdList = &cobra.Command{
	Use:   "list",
	Short: "List the current configuration values",
	Long: `List the current configuration values

The list command lists all persistent configuration values for the New Relic CLI.
`,
	Example: "newrelic config list",
	Run: func(cmd *cobra.Command, args []string) {
		m := map[string]interface{}{}

		VisitAllConfigFields(func(fd FieldDefinition) {
			m[string(fd.Key)] = GetConfigString(fd.Key)
		})

		output.Text(m)
	},
	Aliases: []string{
		"ls",
	},
}

var cmdReset = &cobra.Command{
	Use:   "reset",
	Short: "Reset a configuration value to its default",
	Long: `Reset a configuration value

The reset command resets a configuration value to its default.
`,
	Example: "newrelic config reset --key <key>",
	Run: func(cmd *cobra.Command, args []string) {
		if !isValidConfigKey() {
			log.Fatalf("%s is not a valid config field. valid values are %s", key, GetValidConfigKeys())
		}

		fd := GetConfigFieldDefinition(ConfigKey(key))

		if fd.Default == nil {
			log.Fatalf("key %s cannot be reset to a default value since no default exists", fd.Key)
		}

		if err := SetConfigValue(ConfigKey(key), fd.Default); err != nil {
			log.Fatal(err)
		}

		log.Info("success")
	},
	Aliases: []string{
		"rm",
		"delete",
	},
}

func init() {
	Command.AddCommand(cmdList)

	Command.AddCommand(cmdSet)
	cmdSet.Flags().StringVarP(&key, "key", "k", "", "the key to set")
	cmdSet.Flags().StringVarP(&value, "value", "v", "", "the value to set")
	utils.LogIfError(cmdSet.MarkFlagRequired("key"))
	utils.LogIfError(cmdSet.MarkFlagRequired("value"))

	Command.AddCommand(cmdGet)
	cmdGet.Flags().StringVarP(&key, "key", "k", "", "the key to get")
	utils.LogIfError(cmdGet.MarkFlagRequired("key"))

	Command.AddCommand(cmdReset)
	cmdReset.Flags().StringVarP(&key, "key", "k", "", "the key to delete")
	utils.LogIfError(cmdReset.MarkFlagRequired("key"))
}
