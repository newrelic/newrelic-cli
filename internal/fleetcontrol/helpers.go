package fleetcontrol

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-client-go/v2/pkg/fleetcontrol"
)

// Shared helper functions used across multiple commands
// These provide common functionality without code duplication

// GetOrganizationID retrieves the organization ID from the provided value or fetches it from the API.
// This is used by commands that need an organization ID but make it optional in their flags.
//
// Parameters:
//   - providedOrgID: The organization ID provided by the user (may be empty)
//
// Returns:
//   - The organization ID (either provided or fetched)
func GetOrganizationID(providedOrgID string) string {
	if providedOrgID != "" {
		return providedOrgID
	}

	org, err := client.NRClient.Organization.GetOrganization()
	if err != nil {
		log.Warnf("Failed to get organization: %v", err)
		return ""
	}
	return org.ID
}

// ParseTags converts tag strings in the format "key:value1,value2" into FleetControlTagInput structs.
// Tags are used to organize and categorize fleet resources.
//
// Format: Each tag string should be "key:value1,value2" where:
//   - key: The tag key (required, cannot be empty)
//   - value1,value2: Comma-separated values for this key (required, cannot be empty)
//
// Parameters:
//   - tagStrings: Array of tag strings to parse
//
// Returns:
//   - Array of FleetControlTagInput structs
//   - Error if any tag has invalid format
//
// Example:
//   tags, err := ParseTags([]string{"env:prod,staging", "team:platform"})
//   // Returns: [{Key: "env", Values: ["prod", "staging"]}, {Key: "team", Values: ["platform"]}]
func ParseTags(tagStrings []string) ([]fleetcontrol.FleetControlTagInput, error) {
	if len(tagStrings) == 0 {
		return nil, nil
	}

	tags := make([]fleetcontrol.FleetControlTagInput, 0, len(tagStrings))

	for _, tagStr := range tagStrings {
		// Split on first colon to separate key from values
		parts := strings.SplitN(tagStr, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid tag format '%s': expected 'key:value1,value2'", tagStr)
		}

		key := strings.TrimSpace(parts[0])
		if key == "" {
			return nil, fmt.Errorf("tag key cannot be empty in '%s'", tagStr)
		}

		valueStr := strings.TrimSpace(parts[1])
		if valueStr == "" {
			return nil, fmt.Errorf("tag values cannot be empty in '%s'", tagStr)
		}

		// Split values on comma and trim whitespace
		values := strings.Split(valueStr, ",")
		for i := range values {
			values[i] = strings.TrimSpace(values[i])
		}

		tags = append(tags, fleetcontrol.FleetControlTagInput{
			Key:    key,
			Values: values,
		})
	}

	return tags, nil
}

// ParseAgentSpec parses an agent specification string into a FleetControlAgentInput struct.
// Agent specs are used to define agent configurations for deployments.
//
// Format: "AgentType:Version:ConfigVersionID1,ConfigVersionID2,..."
//   - AgentType: The type of agent (e.g., NRInfra, NRDOT)
//   - Version: The agent version (e.g., 1.70.0, 2.0.0)
//   - ConfigVersionIDs: Comma-separated list of configuration version IDs
//
// Parameters:
//   - agentSpec: Agent specification string to parse
//
// Returns:
//   - FleetControlAgentInput struct
//   - Error if format is invalid
//
// Example:
//   agent, err := ParseAgentSpec("NRInfra:1.70.0:version1,version2")
//   // Returns: {AgentType: "NRInfra", Version: "1.70.0", ConfigurationVersionList: [{ID: "version1"}, {ID: "version2"}]}
func ParseAgentSpec(agentSpec string) (fleetcontrol.FleetControlAgentInput, error) {
	// Split on colon to separate AgentType:Version:ConfigVersionIDs
	parts := strings.SplitN(agentSpec, ":", 3)
	if len(parts) != 3 {
		return fleetcontrol.FleetControlAgentInput{}, fmt.Errorf("expected format 'AgentType:Version:ConfigVersionID1,ConfigVersionID2' but got '%s'", agentSpec)
	}

	agentType := strings.TrimSpace(parts[0])
	version := strings.TrimSpace(parts[1])
	configVersionsStr := strings.TrimSpace(parts[2])

	// Validate required fields
	if agentType == "" {
		return fleetcontrol.FleetControlAgentInput{}, fmt.Errorf("agent type cannot be empty")
	}
	if version == "" {
		return fleetcontrol.FleetControlAgentInput{}, fmt.Errorf("agent version cannot be empty")
	}
	if configVersionsStr == "" {
		return fleetcontrol.FleetControlAgentInput{}, fmt.Errorf("configuration version IDs cannot be empty")
	}

	// Split configuration version IDs on comma
	configVersionIDs := strings.Split(configVersionsStr, ",")
	configVersionList := make([]fleetcontrol.FleetControlConfigurationVersionListInput, 0, len(configVersionIDs))

	for _, versionID := range configVersionIDs {
		versionID = strings.TrimSpace(versionID)
		if versionID == "" {
			continue // Skip empty entries
		}
		configVersionList = append(configVersionList, fleetcontrol.FleetControlConfigurationVersionListInput{
			ID: versionID,
		})
	}

	if len(configVersionList) == 0 {
		return fleetcontrol.FleetControlAgentInput{}, fmt.Errorf("at least one configuration version ID is required")
	}

	return fleetcontrol.FleetControlAgentInput{
		AgentType:                agentType,
		Version:                  version,
		ConfigurationVersionList: configVersionList,
	}, nil
}

// ValidateAgentVersionsForFleet validates that agent versions are compatible with the fleet type.
// The "*" wildcard version is only allowed for KUBERNETESCLUSTER fleets, not HOST fleets.
//
// Parameters:
//   - fleetID: The ID of the fleet to validate against
//   - agents: The agent configurations to validate
//
// Returns:
//   - Error if validation fails (e.g., "*" used with HOST fleet)
func ValidateAgentVersionsForFleet(fleetID string, agents []fleetcontrol.FleetControlAgentInput) error {
	// Fetch the fleet entity to check its managed entity type
	entityInterface, err := client.NRClient.FleetControl.GetEntity(fleetID)
	if err != nil {
		return fmt.Errorf("failed to fetch fleet details for validation: %w", err)
	}

	if entityInterface == nil {
		return fmt.Errorf("fleet with ID '%s' not found", fleetID)
	}

	// Type assert to fleet entity
	fleetEntity, ok := (*entityInterface).(*fleetcontrol.EntityManagementFleetEntity)
	if !ok {
		return fmt.Errorf("entity '%s' is not a fleet", fleetID)
	}

	// Check if this is a HOST fleet
	isHostFleet := string(fleetEntity.ManagedEntityType) == "HOST"

	// If it's a HOST fleet, validate that no agent uses "*" as version
	if isHostFleet {
		for _, agent := range agents {
			if agent.Version == "*" {
				return fmt.Errorf(
					"agent version '*' (wildcard) is not supported for HOST fleets. "+
						"Please specify an explicit version (e.g., '1.70.0'). "+
						"Wildcard versions are only supported for KUBERNETESCLUSTER fleets. "+
						"Fleet '%s' is of type: %s",
					fleetID, string(fleetEntity.ManagedEntityType))
			}
		}
	}

	// KUBERNETESCLUSTER fleets allow "*", so no validation needed
	return nil
}

// Type mappers - Convert YAML-validated string values to client library types
// These mappers ensure we use the YAML-validated values without bypassing validation

// MapManagedEntityType converts a validated managed entity type string to the client library type.
// This function should ONLY be called after YAML validation has confirmed the value is allowed.
//
// Parameters:
//   - typeStr: The managed entity type string (already validated by framework)
//
// Returns:
//   - The corresponding FleetControlManagedEntityType
//   - Error if the type is not recognized (should never happen after YAML validation)
func MapManagedEntityType(typeStr string) (fleetcontrol.FleetControlManagedEntityType, error) {
	// Note: YAML validation has already confirmed this value is in allowed_values
	// This mapping must match the YAML allowed_values exactly
	switch strings.ToUpper(typeStr) {
	case "HOST":
		return fleetcontrol.FleetControlManagedEntityTypeTypes.HOST, nil
	case "KUBERNETESCLUSTER":
		return fleetcontrol.FleetControlManagedEntityTypeTypes.KUBERNETESCLUSTER, nil
	default:
		// This should never happen if YAML validation is working correctly
		return fleetcontrol.FleetControlManagedEntityType(""), fmt.Errorf(
			"unrecognized managed entity type '%s' - YAML validation may have failed", typeStr)
	}
}

// MapScopeType converts a validated scope type string to the client library type.
// This function should ONLY be called after YAML validation has confirmed the value is allowed.
//
// Parameters:
//   - typeStr: The scope type string (already validated by framework)
//
// Returns:
//   - The corresponding FleetControlEntityScope
//   - Error if the type is not recognized (should never happen after YAML validation)
func MapScopeType(typeStr string) (fleetcontrol.FleetControlEntityScope, error) {
	// Note: YAML validation has already confirmed this value is in allowed_values
	// This mapping must match the YAML allowed_values exactly
	switch strings.ToUpper(typeStr) {
	case "ACCOUNT":
		return fleetcontrol.FleetControlEntityScopeTypes.ACCOUNT, nil
	case "ORGANIZATION":
		return fleetcontrol.FleetControlEntityScopeTypes.ORGANIZATION, nil
	default:
		// This should never happen if YAML validation is working correctly
		return fleetcontrol.FleetControlEntityScope(""), fmt.Errorf(
			"unrecognized scope type '%s' - YAML validation may have failed", typeStr)
	}
}

// MapConfigurationMode converts a validated mode string to the client library type.
// This function should ONLY be called after YAML validation has confirmed the value is allowed.
//
// Parameters:
//   - modeStr: The mode string (already validated by framework)
//
// Returns:
//   - The corresponding GetConfigurationMode
//   - Error if the mode is not recognized (should never happen after YAML validation)
func MapConfigurationMode(modeStr string) (fleetcontrol.GetConfigurationMode, error) {
	// Note: YAML validation has already confirmed this value is in allowed_values
	// This mapping must match the YAML allowed_values exactly
	switch strings.ToLower(modeStr) {
	case "configentity", "":
		return fleetcontrol.GetConfigurationModeTypes.ConfigEntity, nil
	case "configversionentity":
		return fleetcontrol.GetConfigurationModeTypes.ConfigVersionEntity, nil
	default:
		// This should never happen if YAML validation is working correctly
		return fleetcontrol.GetConfigurationMode(""), fmt.Errorf(
			"unrecognized configuration mode '%s' - YAML validation may have failed", modeStr)
	}
}

// MapOperatingSystemType converts a validated operating system type string to the client library type.
// This function should ONLY be called after YAML validation has confirmed the value is allowed.
//
// Parameters:
//   - osStr: The operating system type string (already validated by framework)
//
// Returns:
//   - The corresponding FleetControlOperatingSystemType
//   - Error if the type is not recognized (should never happen after YAML validation)
func MapOperatingSystemType(osStr string) (fleetcontrol.FleetControlOperatingSystemType, error) {
	// Note: YAML validation has already confirmed this value is in allowed_values
	// This mapping must match the YAML allowed_values exactly
	switch strings.ToUpper(osStr) {
	case "LINUX":
		return fleetcontrol.FleetControlOperatingSystemTypeTypes.LINUX, nil
	case "WINDOWS":
		return fleetcontrol.FleetControlOperatingSystemTypeTypes.WINDOWS, nil
	default:
		// This should never happen if YAML validation is working correctly
		return fleetcontrol.FleetControlOperatingSystemType(""), fmt.Errorf(
			"unrecognized operating system type '%s' - YAML validation may have failed", osStr)
	}
}

// ErrorResponse represents an error response with consistent field ordering.
// Status field comes first, followed by error message.
type ErrorResponse struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

// PrintError outputs an error as JSON with "status" and "error" keys.
// This allows errors to be parsed by consumers expecting JSON output.
// Uses custom JSON marshaling to ensure field order: status first, then error.
//
// Parameters:
//   - err: The error to print
//
// Returns:
//   - Error from output.Print (typically nil)
func PrintError(err error) error {
	response := ErrorResponse{
		Status: "failed",
		Error:  err.Error(),
	}
	return printJSON(response)
}

// printJSON marshals and prints data as JSON with preserved field order.
// This bypasses the output.Print prettyjson formatter which sorts keys alphabetically.
//
// Parameters:
//   - data: The data to marshal and print
//
// Returns:
//   - Error if marshaling fails
func printJSON(data interface{}) error {
	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	fmt.Fprintln(os.Stdout, string(jsonBytes))
	return nil
}

// FleetResponseWrapper wraps fleet entity output with consistent status and error fields.
// Used to provide a uniform response structure across all fleet management commands.
type FleetResponseWrapper struct {
	Status            string                                    `json:"status"`
	Error             string                                    `json:"error"`
	ID                string                                    `json:"id,omitempty"`
	Name              string                                    `json:"name,omitempty"`
	Type              string                                    `json:"type,omitempty"`
	ManagedEntityType string                                    `json:"managedEntityType,omitempty"`
	OperatingSystem   *fleetcontrol.FleetControlOperatingSystem `json:"operatingSystem,omitempty"`
	Scope             *fleetcontrol.FleetControlScopedReference `json:"scope,omitempty"`
	Tags              []fleetcontrol.FleetControlTag            `json:"tags,omitempty"`
	Product           []string                                  `json:"product,omitempty"`
	Description       string                                    `json:"description,omitempty"`
	CreatedAt         int64                                     `json:"createdAt,omitempty"`
	UpdatedAt         int64                                     `json:"updatedAt,omitempty"`
}

// FleetListResponseWrapper wraps a list of fleet entities with consistent status and error fields.
type FleetListResponseWrapper struct {
	Status string              `json:"status"`
	Error  string              `json:"error"`
	Fleets []FleetEntityOutput `json:"fleets,omitempty"`
}

// PrintFleetSuccess outputs successful fleet data wrapped with status and error fields.
// This provides a consistent response structure: status="success", error="", <fleet fields>
// Uses custom JSON marshaling to preserve field order.
//
// Parameters:
//   - fleet: The fleet entity output to wrap and print
//
// Returns:
//   - Error from printJSON (typically nil)
func PrintFleetSuccess(fleet *FleetEntityOutput) error {
	response := FleetResponseWrapper{
		Status:            "success",
		Error:             "",
		ID:                fleet.ID,
		Name:              fleet.Name,
		Type:              fleet.Type,
		ManagedEntityType: fleet.ManagedEntityType,
		Description:       fleet.Description,
		CreatedAt:         fleet.CreatedAt,
		UpdatedAt:         fleet.UpdatedAt,
	}

	if fleet.Scope.ID != "" {
		response.Scope = &fleet.Scope
	}

	if fleet.OperatingSystem != nil {
		response.OperatingSystem = fleet.OperatingSystem
	}

	if len(fleet.Tags) > 0 {
		response.Tags = fleet.Tags
	}

	if len(fleet.Product) > 0 {
		response.Product = fleet.Product
	}

	return printJSON(response)
}

// PrintFleetListSuccess outputs a successful list of fleets wrapped with status and error fields.
// Uses custom JSON marshaling to preserve field order.
//
// Parameters:
//   - fleets: The list of fleet entities to wrap and print
//
// Returns:
//   - Error from printJSON (typically nil)
func PrintFleetListSuccess(fleets []FleetEntityOutput) error {
	response := FleetListResponseWrapper{
		Status: "success",
		Error:  "",
		Fleets: fleets,
	}
	return printJSON(response)
}

// FleetDeleteResponseWrapper wraps delete result with consistent status and error fields.
type FleetDeleteResponseWrapper struct {
	Status string `json:"status"`
	Error  string `json:"error"`
	ID     string `json:"id,omitempty"`
}

// PrintDeleteSuccess outputs successful delete operation with status and error fields.
// Uses custom JSON marshaling to preserve field order.
//
// Parameters:
//   - id: The ID of the deleted fleet
//
// Returns:
//   - Error from printJSON (typically nil)
func PrintDeleteSuccess(id string) error {
	response := FleetDeleteResponseWrapper{
		Status: "success",
		Error:  "",
		ID:     id,
	}
	return printJSON(response)
}

// PrintDeleteBulkSuccess outputs successful bulk delete operations as a list.
// Each element in the list has status and error fields.
// Uses custom JSON marshaling to preserve field order.
//
// Parameters:
//   - results: The list of delete results to print
//
// Returns:
//   - Error from printJSON (typically nil)
func PrintDeleteBulkSuccess(results []FleetDeleteResponseWrapper) error {
	return printJSON(results)
}

// FleetEntityOutput is a filtered representation of a fleet entity containing only user-relevant fields.
// This removes verbose metadata and internal structures that aren't useful for command output.
type FleetEntityOutput struct {
	ID                string                                   `json:"id"`
	Name              string                                   `json:"name"`
	Type              string                                   `json:"type,omitempty"`
	ManagedEntityType string                                   `json:"managedEntityType,omitempty"`
	OperatingSystem   *fleetcontrol.FleetControlOperatingSystem `json:"operatingSystem,omitempty"`
	Scope             fleetcontrol.FleetControlScopedReference `json:"scope,omitempty"`
	Tags              []fleetcontrol.FleetControlTag           `json:"tags,omitempty"`
	Product           []string                                 `json:"product,omitempty"`
	Description       string                                   `json:"description,omitempty"`
	CreatedAt         int64                                    `json:"createdAt,omitempty"`
	UpdatedAt         int64                                    `json:"updatedAt,omitempty"`
}

// FilterFleetEntity creates a filtered output from a fleet entity, removing verbose fields.
// This keeps only the fields that correspond to command flags or essential metadata.
//
// Parameters:
//   - entity: The fleet entity from the API response
//
// Returns:
//   - Filtered fleet entity output
func FilterFleetEntity(entity fleetcontrol.FleetControlFleetEntityResult) *FleetEntityOutput {
	output := &FleetEntityOutput{
		ID:   entity.ID,
		Name: entity.Name,
	}

	// Add optional string fields
	if entity.Type != "" {
		output.Type = entity.Type
	}

	if entity.ManagedEntityType != "" {
		output.ManagedEntityType = string(entity.ManagedEntityType)
	}

	// Add scope (always present)
	output.Scope = entity.Scope

	// Add operating system if present (only applicable for HOST fleets)
	if entity.OperatingSystem.Type != "" {
		output.OperatingSystem = &entity.OperatingSystem
	}

	// Add tags if present
	if len(entity.Tags) > 0 {
		output.Tags = entity.Tags
	}

	// Add product if present and not empty
	if len(entity.Product) > 0 {
		output.Product = entity.Product
	}

	// Add description if present
	if entity.Description != "" {
		output.Description = entity.Description
	}

	// Add timestamps from metadata
	// EpochMilliseconds is a time.Time under the hood, convert to Unix milliseconds
	createdTime := time.Time(entity.Metadata.CreatedAt)
	if !createdTime.IsZero() {
		output.CreatedAt = createdTime.UnixMilli()
	}

	updatedTime := time.Time(entity.Metadata.UpdatedAt)
	if !updatedTime.IsZero() {
		output.UpdatedAt = updatedTime.UnixMilli()
	}

	return output
}

// FilterFleetEntityFromEntityManagement creates a filtered output from an EntityManagementFleetEntity.
// This is used for the get and search commands which use the EntityManagement API.
//
// Parameters:
//   - entity: The fleet entity from EntityManagement API
//   - showTags: Whether to include tags in the output
//
// Returns:
//   - Filtered fleet entity output
func FilterFleetEntityFromEntityManagement(entity fleetcontrol.EntityManagementFleetEntity, showTags bool) *FleetEntityOutput {
	output := &FleetEntityOutput{
		ID:   entity.ID,
		Name: entity.Name,
	}

	// Add optional string fields
	if entity.Type != "" {
		output.Type = entity.Type
	}

	if entity.ManagedEntityType != "" {
		output.ManagedEntityType = string(entity.ManagedEntityType)
	}

	// Convert EntityManagementScopedReference to FleetControlScopedReference
	output.Scope = fleetcontrol.FleetControlScopedReference{
		ID:   entity.Scope.ID,
		Type: fleetcontrol.FleetControlEntityScope(entity.Scope.Type),
	}

	// Add operating system if present (only applicable for HOST fleets)
	// Convert EntityManagementOperatingSystem to FleetControlOperatingSystem
	if entity.OperatingSystem.Type != "" {
		output.OperatingSystem = &fleetcontrol.FleetControlOperatingSystem{
			Type: fleetcontrol.FleetControlOperatingSystemType(entity.OperatingSystem.Type),
		}
	}

	// Convert tags from EntityManagementTag to FleetControlTag
	// Only include tags if showTags flag is true
	if showTags && len(entity.Tags) > 0 {
		tags := make([]fleetcontrol.FleetControlTag, len(entity.Tags))
		for i, tag := range entity.Tags {
			tags[i] = fleetcontrol.FleetControlTag(tag)
		}
		output.Tags = tags
	}

	// Add product if present and not empty
	if len(entity.Product) > 0 {
		output.Product = entity.Product
	}

	// Add description if present
	if entity.Description != "" {
		output.Description = entity.Description
	}

	// Add timestamps from metadata
	// EpochMilliseconds is a time.Time under the hood, convert to Unix milliseconds
	createdTime := time.Time(entity.Metadata.CreatedAt)
	if !createdTime.IsZero() {
		output.CreatedAt = createdTime.UnixMilli()
	}

	updatedTime := time.Time(entity.Metadata.UpdatedAt)
	if !updatedTime.IsZero() {
		output.UpdatedAt = updatedTime.UnixMilli()
	}

	return output
}

// ConfigurationResponseWrapper wraps configuration output with consistent status and error fields.
// Used for create, add-version, and get operations on configurations.
type ConfigurationResponseWrapper struct {
	Status string      `json:"status"`
	Error  string      `json:"error"`
	Result interface{} `json:"result,omitempty"`
}

// PrintConfigurationSuccess outputs successful configuration data wrapped with status and error fields.
// Uses custom JSON marshaling to preserve field order.
//
// Parameters:
//   - result: The configuration result to wrap and print
//
// Returns:
//   - Error from printJSON (typically nil)
func PrintConfigurationSuccess(result interface{}) error {
	response := ConfigurationResponseWrapper{
		Status: "success",
		Error:  "",
		Result: result,
	}
	return printJSON(response)
}

// ConfigurationDeleteResponseWrapper wraps configuration delete result with consistent status and error fields.
type ConfigurationDeleteResponseWrapper struct {
	Status string `json:"status"`
	Error  string `json:"error"`
	ID     string `json:"id,omitempty"`
}

// PrintConfigurationDeleteSuccess outputs successful configuration delete operation with status and error fields.
// Uses custom JSON marshaling to preserve field order.
//
// Parameters:
//   - id: The ID of the deleted configuration
//
// Returns:
//   - Error from printJSON (typically nil)
func PrintConfigurationDeleteSuccess(id string) error {
	response := ConfigurationDeleteResponseWrapper{
		Status: "success",
		Error:  "",
		ID:     id,
	}
	return printJSON(response)
}

// DeploymentResponseWrapper wraps deployment output with consistent status and error fields.
// Used for create and update operations on deployments.
type DeploymentResponseWrapper struct {
	Status string      `json:"status"`
	Error  string      `json:"error"`
	Result interface{} `json:"result,omitempty"`
}

// PrintDeploymentSuccess outputs successful deployment data wrapped with status and error fields.
// Uses custom JSON marshaling to preserve field order.
//
// Parameters:
//   - result: The deployment result to wrap and print
//
// Returns:
//   - Error from printJSON (typically nil)
func PrintDeploymentSuccess(result interface{}) error {
	response := DeploymentResponseWrapper{
		Status: "success",
		Error:  "",
		Result: result,
	}
	return printJSON(response)
}

// DeploymentDeleteResponseWrapper wraps deployment delete result with consistent status and error fields.
type DeploymentDeleteResponseWrapper struct {
	Status string `json:"status"`
	Error  string `json:"error"`
	ID     string `json:"id,omitempty"`
}

// PrintDeploymentDeleteSuccess outputs successful deployment delete operation with status and error fields.
// Uses custom JSON marshaling to preserve field order.
//
// Parameters:
//   - id: The ID of the deleted deployment
//
// Returns:
//   - Error from printJSON (typically nil)
func PrintDeploymentDeleteSuccess(id string) error {
	response := DeploymentDeleteResponseWrapper{
		Status: "success",
		Error:  "",
		ID:     id,
	}
	return printJSON(response)
}

// MemberResponseWrapper wraps member operation output with consistent status and error fields.
// Used for add-members and remove-members operations.
type MemberResponseWrapper struct {
	Status string      `json:"status"`
	Error  string      `json:"error"`
	Result interface{} `json:"result,omitempty"`
}

// PrintMemberSuccess outputs successful member operation data wrapped with status and error fields.
// Uses custom JSON marshaling to preserve field order.
//
// Parameters:
//   - result: The member operation result to wrap and print
//
// Returns:
//   - Error from printJSON (typically nil)
func PrintMemberSuccess(result interface{}) error {
	response := MemberResponseWrapper{
		Status: "success",
		Error:  "",
		Result: result,
	}
	return printJSON(response)
}
