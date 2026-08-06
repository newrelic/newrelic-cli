package fleetcontrol

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-client-go/v2/pkg/fleetcontrol"
)

// handleFleetCreateConfiguration implements the 'create-configuration' command to create a fleet configuration.
//
// This command creates a new fleet configuration with custom attributes.
// Configurations allow you to define settings and attributes for fleet management.
//
// The configuration body must be provided via one of two mutually exclusive flags:
//   - --configuration-file-path: Path to a file (recommended for production)
//   - --configuration-content: Inline content (for testing/development/emergency only)
//
// The command:
// 1. Validates flag values (done automatically by framework via YAML rules)
// 2. Validates mutual exclusivity of configuration flags
// 3. Derives organization ID if not provided
// 4. Builds custom headers with entity metadata
// 5. Calls the New Relic API to create the configuration
//
// Parameters:
//   - cmd: The cobra command being executed
//   - args: Command arguments (not used)
//   - flags: Validated flag values from YAML configuration
//
// Returns:
//   - Error if configuration creation fails, nil on success
func handleFleetCreateConfiguration(_ *cobra.Command, _ []string, flags *FlagValues) error {
	// Get typed flag values - no hardcoded strings!
	// Note: CreateConfiguration() returns error because it handles file reading
	f, err := flags.CreateConfiguration()
	if err != nil {
		return fmt.Errorf("failed to read flags: %w", err)
	}

	// Map validated managed entity type to client library type
	// YAML validation has already confirmed this value is in allowed_values
	entityType, err := MapManagedEntityType(f.ManagedEntityType)
	if err != nil {
		return PrintError(err)
	}

	// Validate operating system requirements based on entity type
	// For HOST configurations, operating system must be specified
	// For KUBERNETESCLUSTER configurations, operating system should not be specified
	if entityType == fleetcontrol.FleetControlManagedEntityTypeTypes.HOST {
		if f.OperatingSystem == "" {
			return PrintError(fmt.Errorf("--operating-system is required when --managed-entity-type is HOST (allowed values: LINUX, WINDOWS)"))
		}
	} else if entityType == fleetcontrol.FleetControlManagedEntityTypeTypes.KUBERNETESCLUSTER {
		if f.OperatingSystem != "" {
			return PrintError(fmt.Errorf("--operating-system should not be specified for KUBERNETESCLUSTER configurations"))
		}
	}

	// Validate that exactly one of --configuration-file-path or --configuration-content is provided
	hasFilePath := f.ConfigurationFilePath != ""
	hasContent := f.ConfigurationContent != ""

	if !hasFilePath && !hasContent {
		return fmt.Errorf("one of --configuration-file-path or --configuration-content must be provided")
	}

	if hasFilePath && hasContent {
		return fmt.Errorf("--configuration-file-path and --configuration-content are mutually exclusive, use only one")
	}

	// Determine which configuration content to use
	configBody := f.ConfigurationFilePath
	if hasContent {
		configBody = f.ConfigurationContent
	}

	// Convert the configuration body to []byte to prevent JSON marshaling
	// This preserves newlines and formatting in the configuration file
	configBodyBytes := []byte(configBody)

	// Get organization ID (provided or fetched from API)
	orgID := GetOrganizationID(f.OrganizationID)

	// Build custom headers required by the API
	// These headers specify the entity name, agent type, managed entity type, and operating system (for HOST)
	var entityHeader string
	if f.OperatingSystem != "" {
		entityHeader = fmt.Sprintf(
			`{"name": "%s", "agentType": "%s", "managedEntityType": "%s", "operatingSystem": {"type": "%s"}}`,
			f.Name,
			f.AgentType,
			f.ManagedEntityType,
			f.OperatingSystem,
		)
	} else {
		entityHeader = fmt.Sprintf(
			`{"name": "%s", "agentType": "%s", "managedEntityType": "%s"}`,
			f.Name,
			f.AgentType,
			f.ManagedEntityType,
		)
	}
	customHeaders := map[string]interface{}{
		"x-newrelic-client-go-custom-headers": map[string]string{
			"Newrelic-Entity": entityHeader,
		},
	}

	// Call New Relic API to create the configuration
	result, err := client.NRClient.FleetControl.FleetControlCreateConfiguration(
		configBodyBytes,
		customHeaders,
		orgID,
	)
	if err != nil {
		return PrintError(fmt.Errorf("failed to create fleet configuration: %w", err))
	}

	// Print the created configuration to stdout with status wrapper
	return PrintConfigurationSuccess(result)
}
