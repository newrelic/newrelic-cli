package entities

import (
	"github.com/spf13/cobra"
)

// Should these be moved out or made into higher-level flags?
var (
	entityName          string
	entityGUID          string
	entityValues        []string
	entityType          string
	entityAlertSeverity string
	entityDomain        string
	entityReporting     string
	entityFields        []string
)

// Command represents the entities command
var Command = &cobra.Command{
	Use:   "entities",
	Short: "Subcommands to interact with New Relic entities",
}
