package fleetcontrol

import (
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
	// Command variables for all fleet subcommands
	// These are populated during initialization from YAML configs
	cmdFleetCreate                   *cobra.Command
	cmdFleetUpdate                   *cobra.Command
	cmdFleetDelete                   *cobra.Command
	cmdFleetGet                      *cobra.Command
	cmdFleetSearch                   *cobra.Command
	cmdFleetAddMembers               *cobra.Command
	cmdFleetRemoveMembers            *cobra.Command
	cmdFleetListMembers              *cobra.Command
	cmdFleetCreateConfiguration      *cobra.Command
	cmdFleetGetConfiguration         *cobra.Command
	cmdFleetGetConfigurationVersions *cobra.Command
	cmdFleetAddVersion               *cobra.Command
	cmdFleetDeleteConfiguration      *cobra.Command
	cmdFleetDeleteVersion            *cobra.Command
	cmdFleetCreateDeployment         *cobra.Command
	cmdFleetUpdateDeployment         *cobra.Command
)

// init initializes all fleet subcommands by loading their YAML configurations
// and wiring them to their handler functions.
//
// This function:
// 1. Loads all command definitions from configs/*.yaml files
// 2. Maps each command name to its handler function
// 3. Builds cobra commands with automatic flag registration and validation
// 4. Registers all commands with the parent cmdFleet command
//
// The initialization is driven entirely by YAML configuration.
// Adding a new command requires:
//   - Creating a new YAML file in configs/
//   - Creating a new handler function in command_<name>.go
//   - Adding the handler to the handlers map below
func init() {
	// Load all command configurations from the configs directory
	// This reads and parses all *.yaml files in configs/
	config, err := LoadCommandConfig()
	if err != nil {
		log.Fatalf("Failed to load command config: %v", err)
	}

	// Map command names to their handler functions
	// Each handler implements the business logic for one command
	// Handlers are defined in individual command files
	// Config files are in configs/ directory matching handler file names
	handlers := map[string]CommandHandler{
		"create":               handleFleetCreate,
		"update":               handleFleetUpdate,
		"delete":               handleFleetDelete,
		"get":                  handleFleetGet,
		"search":               handleFleetSearch,
		"add-members":          handleFleetAddMembers,
		"remove-members":       handleFleetRemoveMembers,
		"list-members":         handleFleetListMembers,
		"create-configuration": handleFleetCreateConfiguration,
		"get-configuration":    handleFleetGetConfiguration,
		"get-versions":         handleFleetGetConfigurationVersions,
		"add-version":          handleFleetAddVersion,
		"delete-configuration": handleFleetDeleteConfiguration,
		"delete-version":       handleFleetDeleteVersion,
		"create-deployment":    handleFleetCreateDeployment,
		"update-deployment":    handleFleetUpdateDeployment,
	}

	// Build cobra commands from YAML configurations
	// For each command definition, we:
	//   1. Find the corresponding handler function
	//   2. Build a cobra command with automatic flag registration
	//   3. Assign the command to the appropriate variable
	for _, cmdDef := range config.Commands {
		handler, ok := handlers[cmdDef.Name]
		if !ok {
			log.Fatalf("No handler found for command: %s", cmdDef.Name)
		}

		// BuildCommand creates a cobra command from the YAML definition
		// It automatically:
		//   - Registers all flags from the YAML
		//   - Sets up validation rules
		//   - Wires the handler function
		cmd := BuildCommand(cmdDef, handler)

		// All fleet commands require client initialization
		cmd.PreRun = client.RequireClient

		// Assign the built command to the appropriate variable
		// This makes commands accessible to tests and other code
		switch cmdDef.Name {
		case "create":
			cmdFleetCreate = cmd
		case "update":
			cmdFleetUpdate = cmd
		case "delete":
			cmdFleetDelete = cmd
		case "get":
			cmdFleetGet = cmd
		case "search":
			cmdFleetSearch = cmd
		case "add-members":
			cmdFleetAddMembers = cmd
		case "remove-members":
			cmdFleetRemoveMembers = cmd
		case "list-members":
			cmdFleetListMembers = cmd
		case "create-configuration":
			cmdFleetCreateConfiguration = cmd
		case "get-configuration":
			cmdFleetGetConfiguration = cmd
		case "get-versions":
			cmdFleetGetConfigurationVersions = cmd
		case "add-version":
			cmdFleetAddVersion = cmd
		case "delete-configuration":
			cmdFleetDeleteConfiguration = cmd
		case "delete-version":
			cmdFleetDeleteVersion = cmd
		case "create-deployment":
			cmdFleetCreateDeployment = cmd
		case "update-deployment":
			cmdFleetUpdateDeployment = cmd
		}
	}

	// Register all subcommands with the parent cmdFleet command
	// This makes them available as: newrelic fleetcontrol fleet <command>
	cmdFleet.AddCommand(cmdFleetCreate)
	cmdFleet.AddCommand(cmdFleetUpdate)
	cmdFleet.AddCommand(cmdFleetDelete)
	cmdFleet.AddCommand(cmdFleetGet)
	cmdFleet.AddCommand(cmdFleetSearch)
	cmdFleet.AddCommand(cmdFleetAddMembers)
	cmdFleet.AddCommand(cmdFleetRemoveMembers)
	cmdFleet.AddCommand(cmdFleetListMembers)
	cmdFleet.AddCommand(cmdFleetCreateConfiguration)
	cmdFleet.AddCommand(cmdFleetGetConfiguration)
	cmdFleet.AddCommand(cmdFleetGetConfigurationVersions)
	cmdFleet.AddCommand(cmdFleetAddVersion)
	cmdFleet.AddCommand(cmdFleetDeleteConfiguration)
	cmdFleet.AddCommand(cmdFleetDeleteVersion)
	cmdFleet.AddCommand(cmdFleetCreateDeployment)
	cmdFleet.AddCommand(cmdFleetUpdateDeployment)
}
