/*
Package fleetcontrol provides commands for managing New Relic Fleet Control resources.

The fleetcontrol package allows users to manage fleets and their associated entities
through the New Relic CLI. It supports operations such as creating, updating, and
deleting fleets, as well as managing fleet members.
*/
package fleetcontrol

import (
	"github.com/spf13/cobra"
)

var (
	// Command represents the fleetcontrol command
	Command = &cobra.Command{
		Use:   "fleetcontrol",
		Short: "Manage New Relic Fleet Control resources",
		Long: `The fleetcontrol command provides subcommands for managing Fleet Control resources.
Fleet Control allows you to organize and manage collections of entities such as hosts
and Kubernetes clusters.`,
	}

	// Fleet subcommand
	cmdFleet = &cobra.Command{
		Use:   "fleet",
		Short: "Manage fleet entities",
		Long:  "Manage fleet entities including creation, updates, deletion, and member management.",
	}
)

func init() {
	Command.AddCommand(cmdFleet)
	// Note: Fleet subcommands are registered in command_fleet.go's init() after they're built from YAML
}
