package fleetcontrol

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-cli/internal/output"
)

// handleFleetGetConfiguration implements the 'get-configuration' command to retrieve a fleet configuration.
//
// This command retrieves a fleet configuration or a specific version of a configuration.
// You can:
//   - Get the latest version of a configuration
//   - Get a specific version by number
//   - Get a configuration version by its entity GUID
//
// The command:
// 1. Validates flag values (done automatically by framework via YAML rules)
// 2. Maps the validated mode string to the client library type
// 3. Derives organization ID if not provided
// 4. Calls the New Relic API to retrieve the configuration
//
// Parameters:
//   - cmd: The cobra command being executed
//   - args: Command arguments (not used)
//   - flags: Validated flag values from YAML configuration
//
// Returns:
//   - Error if configuration retrieval fails, nil on success
func handleFleetGetConfiguration(cmd *cobra.Command, args []string, flags *FlagValues) error {
	// Get typed flag values - no hardcoded strings!
	f := flags.GetConfiguration()

	// Get organization ID (provided or fetched from API)
	orgID := GetOrganizationID(f.OrganizationID)

	// Map validated mode to client library type
	// YAML validation has already confirmed this value is in allowed_values
	mode, err := MapConfigurationMode(f.Mode)
	if err != nil {
		return PrintError(err)
	}

	// Call New Relic API to get the configuration
	result, err := client.NRClient.FleetControl.FleetControlGetConfiguration(
		f.EntityGUID,
		orgID,
		mode,
		f.Version,
	)
	if err != nil {
		return PrintError(fmt.Errorf("failed to get configuration: %w", err))
	}

	// Print the configuration directly to stdout (no wrapper for success)
	// This allows the raw configuration to be used directly or formatted as table
	return output.Print(result)
}
