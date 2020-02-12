package config

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	cfg *Config
	// Display keys when printing output
	showKeys bool
)

// SetConfig takes a pointer to the loaded config for later reference
func SetConfig(c *Config) {
	cfg = c
}

// Command is the base command for managing profiles
var Command = &cobra.Command{
	Use:   "config",
	Short: "configuration management",
}

var cmdSet = &cobra.Command{
	Use:   "set",
	Short: "set a new configuration value",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("config set has not been implemented")
	},
}

var cmdGet = &cobra.Command{
	Use:   "get",
	Short: "get a configuration value",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("config get has not been implemented")
	},
}

var cmdList = &cobra.Command{
	Use:   "list",
	Short: "list configuration values",
	Run: func(cmd *cobra.Command, args []string) {
		if cfg != nil {
			cfg.List()
		} else {
			fmt.Println("no configuration values found")
		}
	},
	Aliases: []string{
		"ls",
	},
}

var cmdDelete = &cobra.Command{
	Use:   "delete",
	Short: "delete a configuration value",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("config delete has not been implemented")
	},
	Aliases: []string{
		"rm",
	},
}

func init() {
	Command.AddCommand(cmdSet)
	Command.AddCommand(cmdGet)
	Command.AddCommand(cmdList)
	Command.AddCommand(cmdDelete)
}
