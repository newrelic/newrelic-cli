package fleetcontrol

import (
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
)

// Command registration and initialization
//
// This file wires up all the fleet subcommands by:
// 1. Loading YAML configurations from the configs directory
// 2. Mapping each command name to its handler function
// 3. Building cobra commands from the YAML definitions
// 4. Registering all commands with the parent cmdFleet command
//
// Each command's business logic lives in its own file (command_<name>.go)
// Each command's configuration lives in its own YAML file (configs/<name>.yaml)

var (
	// Fleet management command variables
	cmdFleetCreate *cobra.Command
	cmdFleetUpdate *cobra.Command
	cmdFleetDelete *cobra.Command
	cmdFleetGet    *cobra.Command
	cmdFleetSearch *cobra.Command

	// Fleet members command variables (nested under fleet)
	cmdFleetMembersAdd    *cobra.Command
	cmdFleetMembersRemove *cobra.Command
	cmdFleetMembersList   *cobra.Command

	// Configuration command variables
	cmdConfigurationCreate *cobra.Command
	cmdConfigurationGet    *cobra.Command
	cmdConfigurationDelete *cobra.Command

	// Configuration versions command variables (nested under configuration)
	cmdConfigurationVersionsList   *cobra.Command
	cmdConfigurationVersionsAdd    *cobra.Command
	cmdConfigurationVersionsDelete *cobra.Command

	// Deployment command variables
	cmdDeploymentCreate *cobra.Command
	cmdDeploymentUpdate *cobra.Command
	cmdDeploymentDeploy *cobra.Command
	cmdDeploymentDelete *cobra.Command

	// Entities command variables
	cmdEntitiesGetManaged    *cobra.Command
	cmdEntitiesGetUnassigned *cobra.Command
)

// init initializes all fleet control commands by loading their YAML configurations
// and registering them to their appropriate parent commands.
func init() {
	// Load all command configurations from the configs directory
	config, err := LoadCommandConfig()
	if err != nil {
		log.Fatalf("Failed to load command config: %v", err)
	}

	// Initialize each command group
	initFleetCommands(config)
	initConfigurationCommands(config)
	initDeploymentCommands(config)
	initEntitiesCommands(config)
}

// initFleetCommands initializes fleet management and fleet members commands
func initFleetCommands(config *CommandConfig) {
	// Map of filename patterns to handlers for fleet commands
	handlers := map[string]CommandHandler{
		"fleet_management_create": handleFleetCreate,
		"fleet_management_update": handleFleetUpdate,
		"fleet_management_delete": handleFleetDelete,
		"fleet_management_get":    handleFleetGet,
		"fleet_management_search": handleFleetSearch,
		"fleet_members_add":       handleFleetAddMembers,
		"fleet_members_remove":    handleFleetRemoveMembers,
		"fleet_members_list":      handleFleetListMembers,
	}

	// Build commands from YAML and assign to variables
	for _, cmdDef := range config.Commands {
		// Determine handler based on the YAML filename pattern
		var handler CommandHandler
		var cmdVar **cobra.Command

		// Match by checking what the command is about
		if cmdDef.Name == "create" && contains(cmdDef.Short, "fleet") && !contains(cmdDef.Short, "configuration") && !contains(cmdDef.Short, "deployment") {
			handler = handlers["fleet_management_create"]
			cmdVar = &cmdFleetCreate
		} else if cmdDef.Name == "update" && contains(cmdDef.Short, "fleet") && !contains(cmdDef.Short, "deployment") {
			handler = handlers["fleet_management_update"]
			cmdVar = &cmdFleetUpdate
		} else if cmdDef.Name == "delete" && contains(cmdDef.Short, "fleet") && !contains(cmdDef.Short, "configuration") && !contains(cmdDef.Short, "version") {
			handler = handlers["fleet_management_delete"]
			cmdVar = &cmdFleetDelete
		} else if cmdDef.Name == "get" && contains(cmdDef.Short, "fleet") && !contains(cmdDef.Short, "configuration") {
			handler = handlers["fleet_management_get"]
			cmdVar = &cmdFleetGet
		} else if cmdDef.Name == "search" {
			handler = handlers["fleet_management_search"]
			cmdVar = &cmdFleetSearch
		} else if cmdDef.Name == "add" && contains(cmdDef.Short, "entities") {
			handler = handlers["fleet_members_add"]
			cmdVar = &cmdFleetMembersAdd
		} else if cmdDef.Name == "remove" && contains(cmdDef.Short, "entities") {
			handler = handlers["fleet_members_remove"]
			cmdVar = &cmdFleetMembersRemove
		} else if cmdDef.Name == "list" && contains(cmdDef.Short, "entities") {
			handler = handlers["fleet_members_list"]
			cmdVar = &cmdFleetMembersList
		} else {
			continue // Not a fleet command
		}

		cmd := BuildCommand(cmdDef, handler)
		cmd.PreRun = client.RequireClient
		*cmdVar = cmd
	}

	// Register fleet management commands
	cmdFleet.AddCommand(cmdFleetCreate)
	cmdFleet.AddCommand(cmdFleetUpdate)
	cmdFleet.AddCommand(cmdFleetDelete)
	cmdFleet.AddCommand(cmdFleetGet)
	cmdFleet.AddCommand(cmdFleetSearch)

	// Register fleet members commands (nested under fleet members)
	cmdFleetMembers.AddCommand(cmdFleetMembersAdd)
	cmdFleetMembers.AddCommand(cmdFleetMembersRemove)
	cmdFleetMembers.AddCommand(cmdFleetMembersList)
}

// initConfigurationCommands initializes configuration and configuration versions commands
func initConfigurationCommands(config *CommandConfig) {
	handlers := map[string]CommandHandler{
		"configuration_create":        handleFleetCreateConfiguration,
		"configuration_get":           handleFleetGetConfiguration,
		"configuration_delete":        handleFleetDeleteConfiguration,
		"configuration_versions_list": handleFleetGetConfigurationVersions,
		"configuration_versions_add":  handleFleetAddVersion,
		"configuration_versions_delete": handleFleetDeleteVersion,
	}

	for _, cmdDef := range config.Commands {
		var handler CommandHandler
		var cmdVar **cobra.Command

		if cmdDef.Name == "create" && contains(cmdDef.Short, "configuration") {
			handler = handlers["configuration_create"]
			cmdVar = &cmdConfigurationCreate
		} else if cmdDef.Name == "get" && contains(cmdDef.Short, "configuration") {
			handler = handlers["configuration_get"]
			cmdVar = &cmdConfigurationGet
		} else if cmdDef.Name == "delete" && contains(cmdDef.Short, "configuration") && !contains(cmdDef.Short, "version") {
			handler = handlers["configuration_delete"]
			cmdVar = &cmdConfigurationDelete
		} else if cmdDef.Name == "list" && contains(cmdDef.Short, "versions") {
			handler = handlers["configuration_versions_list"]
			cmdVar = &cmdConfigurationVersionsList
		} else if cmdDef.Name == "add" && contains(cmdDef.Short, "version") {
			handler = handlers["configuration_versions_add"]
			cmdVar = &cmdConfigurationVersionsAdd
		} else if cmdDef.Name == "delete" && contains(cmdDef.Short, "version") {
			handler = handlers["configuration_versions_delete"]
			cmdVar = &cmdConfigurationVersionsDelete
		} else {
			continue // Not a configuration command
		}

		cmd := BuildCommand(cmdDef, handler)
		cmd.PreRun = client.RequireClient
		*cmdVar = cmd
	}

	// Register configuration commands
	cmdConfiguration.AddCommand(cmdConfigurationCreate)
	cmdConfiguration.AddCommand(cmdConfigurationGet)
	cmdConfiguration.AddCommand(cmdConfigurationDelete)

	// Register configuration versions commands (nested under configuration versions)
	cmdConfigurationVersions.AddCommand(cmdConfigurationVersionsList)
	cmdConfigurationVersions.AddCommand(cmdConfigurationVersionsAdd)
	cmdConfigurationVersions.AddCommand(cmdConfigurationVersionsDelete)
}

// initDeploymentCommands initializes deployment commands
func initDeploymentCommands(config *CommandConfig) {
	handlers := map[string]CommandHandler{
		"deployment_create": handleFleetCreateDeployment,
		"deployment_update": handleFleetUpdateDeployment,
		"deployment_deploy": handleFleetDeploy,
		"deployment_delete": handleFleetDeleteDeployment,
	}

	for _, cmdDef := range config.Commands {
		var handler CommandHandler
		var cmdVar **cobra.Command

		if cmdDef.Name == "create" && contains(cmdDef.Short, "deployment") {
			handler = handlers["deployment_create"]
			cmdVar = &cmdDeploymentCreate
		} else if cmdDef.Name == "update" && contains(cmdDef.Short, "deployment") {
			handler = handlers["deployment_update"]
			cmdVar = &cmdDeploymentUpdate
		} else if cmdDef.Name == "deploy" {
			handler = handlers["deployment_deploy"]
			cmdVar = &cmdDeploymentDeploy
		} else if cmdDef.Name == "delete" && contains(cmdDef.Short, "deployment") {
			handler = handlers["deployment_delete"]
			cmdVar = &cmdDeploymentDelete
		} else {
			continue // Not a deployment command
		}

		cmd := BuildCommand(cmdDef, handler)
		cmd.PreRun = client.RequireClient
		*cmdVar = cmd
	}

	// Register deployment commands
	cmdDeployment.AddCommand(cmdDeploymentCreate)
	cmdDeployment.AddCommand(cmdDeploymentUpdate)
	cmdDeployment.AddCommand(cmdDeploymentDeploy)
	cmdDeployment.AddCommand(cmdDeploymentDelete)
}

// initEntitiesCommands initializes entities commands
func initEntitiesCommands(config *CommandConfig) {
	handlers := map[string]CommandHandler{
		"entities_get_managed":    handleFleetGetManagedEntities,
		"entities_get_unassigned": handleFleetGetUnassignedEntities,
	}

	for _, cmdDef := range config.Commands {
		var handler CommandHandler
		var cmdVar **cobra.Command

		if cmdDef.Name == "get-managed" && contains(cmdDef.Short, "managed entities") {
			handler = handlers["entities_get_managed"]
			cmdVar = &cmdEntitiesGetManaged
		} else if cmdDef.Name == "get-unassigned" && contains(cmdDef.Short, "unassigned entities") {
			handler = handlers["entities_get_unassigned"]
			cmdVar = &cmdEntitiesGetUnassigned
		} else {
			continue // Not an entities command
		}

		cmd := BuildCommand(cmdDef, handler)
		cmd.PreRun = client.RequireClient
		*cmdVar = cmd
	}

	// Register entities commands
	cmdEntities.AddCommand(cmdEntitiesGetManaged)
	cmdEntities.AddCommand(cmdEntitiesGetUnassigned)
}

// contains checks if a string contains a substring (case-insensitive)
func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr ||
		 len(s) > len(substr) &&
		 (strings.Contains(strings.ToLower(s), strings.ToLower(substr))))
}
