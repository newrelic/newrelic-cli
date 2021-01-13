package config

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/configuration"
	"github.com/newrelic/newrelic-cli/internal/output"
	"github.com/newrelic/newrelic-cli/internal/utils"
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
		err := configuration.SetConfigValue(configuration.ConfigFieldKey(key), value)
		if err != nil {
			log.Fatal(err)
		}
	},
}

var cmdGet = &cobra.Command{
	Use:   "get",
	Short: "Get a configuration value",
	Long: `Get a configuration value

The get command gets a persistent configuration value for the New Relic CLI.
`,
	Example: "newrelic config get --key <key>",
	Run: func(cmd *cobra.Command, args []string) {
		val, err := configuration.GetConfigValue(configuration.ConfigFieldKey(key))
		if err != nil {
			log.Fatal(err)
		}
		output.Text(val)
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
		vals := []Value{}
		for _, v := range configuration.ConfigFields {
			// Error can be ignored here since all keys will be valid
			val, _ := configuration.GetConfigValue(v.Key)
			vals = append(vals, Value{
				Name:    v.Name,
				Value:   val,
				Default: v.Default,
			})
		}

		output.Text(vals)
	},
	Aliases: []string{
		"ls",
	},
}

var cmdDelete = &cobra.Command{
	Use:   "delete",
	Short: "Delete a configuration value",
	Long: `Delete a configuration value

The delete command deletes a persistent configuration value for the New Relic CLI.
This will have the effect of resetting the value to its default.
`,
	Example: "newrelic config delete --key <key>",
	Run: func(cmd *cobra.Command, args []string) {
		err := configuration.SetConfigValue(configuration.ConfigFieldKey(key), "")
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
	},
	Aliases: []string{
		"rm",
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

	Command.AddCommand(cmdDelete)
	cmdDelete.Flags().StringVarP(&key, "key", "k", "", "the key to delete")
	utils.LogIfError(cmdDelete.MarkFlagRequired("key"))
}
