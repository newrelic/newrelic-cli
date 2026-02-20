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

	// Fleet subcommand - manages fleet entities
	cmdFleet = &cobra.Command{
		Use:   "fleet",
		Short: "Manage fleet entities",
		Long:  "Manage fleet entities including creation, updates, deletion, and member management.",
	}

	// Fleet members subcommand - nested under fleet
	cmdFleetMembers = &cobra.Command{
		Use:   "members",
		Short: "Manage fleet members",
		Long:  "Manage fleet members including adding, removing, and listing entities in fleet rings.",
	}

	// Configuration subcommand - top-level command for configurations
	cmdConfiguration = &cobra.Command{
		Use:   "configuration",
		Short: "Manage fleet configurations",
		Long:  "Manage fleet configurations including creating, retrieving, and deleting configurations.",
	}

	// Configuration versions subcommand - nested under configuration
	cmdConfigurationVersions = &cobra.Command{
		Use:   "versions",
		Short: "Manage configuration versions",
		Long:  "Manage configuration versions including listing, adding, and deleting versions.",
	}

	// Deployment subcommand - top-level command for deployments
	cmdDeployment = &cobra.Command{
		Use:   "deployment",
		Short: "Manage fleet deployments",
		Long:  "Manage fleet deployments including creating, updating, and triggering deployments across fleet rings.",
	}

	// Entities subcommand - top-level command for entity queries
	cmdEntities = &cobra.Command{
		Use:   "entities",
		Short: "Query entities for fleet management",
		Long:  "Query entities to identify which are managed by fleets and which are available for assignment.",
	}
)

func init() {
	// Register fleet members as nested under fleet
	cmdFleet.AddCommand(cmdFleetMembers)

	// Register configuration versions as nested under configuration
	cmdConfiguration.AddCommand(cmdConfigurationVersions)

	// Register top-level commands
	Command.AddCommand(cmdFleet)
	Command.AddCommand(cmdConfiguration)
	Command.AddCommand(cmdDeployment)
	Command.AddCommand(cmdEntities)

	// Note: Subcommands are registered in command_fleet.go's init() after they're built from YAML
}
