// Code generated from command_config.yaml - DO NOT EDIT manually
// To regenerate: run `go generate` in this package
package fleetcontrol

// Typed flag accessors for each command - provides type-safe, zero-hardcoded-string access to flags

// CreateFlags provides typed access to 'create' command flags
type CreateFlags struct {
	Name              string
	ManagedEntityType string
	Description       string
	Product           string
	OrganizationID    string
	OperatingSystem   string
	Tags              []string
}

// Create returns typed flags for the 'create' command
func (fv *FlagValues) Create() CreateFlags {
	return CreateFlags{
		Name:              fv.GetString("name"),
		ManagedEntityType: fv.GetString("managed-entity-type"),
		Description:       fv.GetString("description"),
		Product:           fv.GetString("product"),
		OrganizationID:    fv.GetString("organization-id"),
		OperatingSystem:   fv.GetString("operating-system"),
		Tags:              fv.GetStringSlice("tags"),
	}
}

// UpdateFlags provides typed access to 'update' command flags
type UpdateFlags struct {
	ID          string
	Name        string
	Description string
	Tags        []string
}

// Update returns typed flags for the 'update' command
func (fv *FlagValues) Update() UpdateFlags {
	return UpdateFlags{
		ID:          fv.GetString("id"),
		Name:        fv.GetString("name"),
		Description: fv.GetString("description"),
		Tags:        fv.GetStringSlice("tags"),
	}
}

// DeleteFlags provides typed access to 'delete' command flags
type DeleteFlags struct {
	FleetID  string
	FleetIDs []string
}

// Delete returns typed flags for the 'delete' command
func (fv *FlagValues) Delete() DeleteFlags {
	return DeleteFlags{
		FleetID:  fv.GetString("fleet-id"),
		FleetIDs: fv.GetStringSlice("fleet-ids"),
	}
}

// GetFlags provides typed access to 'get' command flags
type GetFlags struct {
	FleetID string
}

// Get returns typed flags for the 'get' command
func (fv *FlagValues) Get() GetFlags {
	return GetFlags{
		FleetID: fv.GetString("fleet-id"),
	}
}

// SearchFlags provides typed access to 'search' command flags
type SearchFlags struct {
	NameEquals   string
	NameContains string
	ShowTags     bool
}

// Search returns typed flags for the 'search' command
func (fv *FlagValues) Search() SearchFlags {
	return SearchFlags{
		NameEquals:   fv.GetString("name-equals"),
		NameContains: fv.GetString("name-contains"),
		ShowTags:     fv.GetBool("show-tags"),
	}
}

// AddMembersFlags provides typed access to 'add' (members) command flags
type AddMembersFlags struct {
	FleetID   string
	Ring      string
	EntityIDs []string
}

// AddMembers returns typed flags for the 'add' (members) command
func (fv *FlagValues) AddMembers() AddMembersFlags {
	return AddMembersFlags{
		FleetID:   fv.GetString("fleet-id"),
		Ring:      fv.GetString("ring"),
		EntityIDs: fv.GetStringSlice("entity-ids"),
	}
}

// RemoveMembersFlags provides typed access to 'remove' (members) command flags
type RemoveMembersFlags struct {
	FleetID   string
	Ring      string
	EntityIDs []string
}

// RemoveMembers returns typed flags for the 'remove' (members) command
func (fv *FlagValues) RemoveMembers() RemoveMembersFlags {
	return RemoveMembersFlags{
		FleetID:   fv.GetString("fleet-id"),
		Ring:      fv.GetString("ring"),
		EntityIDs: fv.GetStringSlice("entity-ids"),
	}
}

// ListMembersFlags provides typed access to 'list' (members) command flags
type ListMembersFlags struct {
	FleetID  string
	Ring     string
	ShowTags bool
}

// ListMembers returns typed flags for the 'list' (members) command
func (fv *FlagValues) ListMembers() ListMembersFlags {
	return ListMembersFlags{
		FleetID:  fv.GetString("fleet-id"),
		Ring:     fv.GetString("ring"),
		ShowTags: fv.GetBool("show-tags"),
	}
}

// CreateConfigurationFlags provides typed access to 'create' (configuration) command flags
type CreateConfigurationFlags struct {
	Name                  string
	AgentType             string
	ManagedEntityType     string
	OrganizationID        string
	ConfigurationFilePath string // File content from path
	ConfigurationContent  string // Inline content
}

// CreateConfiguration returns typed flags for the 'create' (configuration) command
func (fv *FlagValues) CreateConfiguration() (CreateConfigurationFlags, error) {
	configFilePath, err := fv.GetFile("configuration-file-path")
	if err != nil {
		return CreateConfigurationFlags{}, err
	}
	return CreateConfigurationFlags{
		Name:                  fv.GetString("name"),
		AgentType:             fv.GetString("agent-type"),
		ManagedEntityType:     fv.GetString("managed-entity-type"),
		OrganizationID:        fv.GetString("organization-id"),
		ConfigurationFilePath: configFilePath,
		ConfigurationContent:  fv.GetString("configuration-content"),
	}, nil
}

// GetConfigurationFlags provides typed access to 'get' (configuration) command flags
type GetConfigurationFlags struct {
	ConfigurationID string
	OrganizationID  string
	Mode            string
	Version         int
}

// GetConfiguration returns typed flags for the 'get' (configuration) command
func (fv *FlagValues) GetConfiguration() GetConfigurationFlags {
	return GetConfigurationFlags{
		ConfigurationID: fv.GetString("configuration-id"),
		OrganizationID:  fv.GetString("organization-id"),
		Mode:            fv.GetString("mode"),
		Version:         fv.GetInt("version"),
	}
}

// GetVersionsFlags provides typed access to 'list' (configuration versions) command flags
type GetVersionsFlags struct {
	ConfigurationID string
	OrganizationID  string
}

// GetVersions returns typed flags for the 'list' (configuration versions) command
func (fv *FlagValues) GetVersions() GetVersionsFlags {
	return GetVersionsFlags{
		ConfigurationID: fv.GetString("configuration-id"),
		OrganizationID:  fv.GetString("organization-id"),
	}
}

// AddVersionFlags provides typed access to 'add' (configuration versions) command flags
type AddVersionFlags struct {
	ConfigurationID       string
	OrganizationID        string
	ConfigurationFilePath string // File content from path
	ConfigurationContent  string // Inline content
}

// AddVersion returns typed flags for the 'add' (configuration versions) command
func (fv *FlagValues) AddVersion() (AddVersionFlags, error) {
	configFilePath, err := fv.GetFile("configuration-file-path")
	if err != nil {
		return AddVersionFlags{}, err
	}
	return AddVersionFlags{
		ConfigurationID:       fv.GetString("configuration-id"),
		OrganizationID:        fv.GetString("organization-id"),
		ConfigurationFilePath: configFilePath,
		ConfigurationContent:  fv.GetString("configuration-content"),
	}, nil
}

// DeleteConfigurationFlags provides typed access to 'delete' (configuration) command flags
type DeleteConfigurationFlags struct {
	ConfigurationID string
	OrganizationID  string
}

// DeleteConfiguration returns typed flags for the 'delete' (configuration) command
func (fv *FlagValues) DeleteConfiguration() DeleteConfigurationFlags {
	return DeleteConfigurationFlags{
		ConfigurationID: fv.GetString("configuration-id"),
		OrganizationID:  fv.GetString("organization-id"),
	}
}

// DeleteVersionFlags provides typed access to 'delete' (configuration versions) command flags
type DeleteVersionFlags struct {
	VersionID      string
	OrganizationID string
}

// DeleteVersion returns typed flags for the 'delete' (configuration versions) command
func (fv *FlagValues) DeleteVersion() DeleteVersionFlags {
	return DeleteVersionFlags{
		VersionID:      fv.GetString("version-id"),
		OrganizationID: fv.GetString("organization-id"),
	}
}

// CreateDeploymentFlags provides typed access to 'create' (deployment) command flags
type CreateDeploymentFlags struct {
	FleetID                 string
	Name                    string
	Agent                   []string
	AgentType               string
	AgentVersion            string
	ConfigurationVersionIDs []string
	Description             string
	Tags                    []string
}

// CreateDeployment returns typed flags for the 'create' (deployment) command
func (fv *FlagValues) CreateDeployment() CreateDeploymentFlags {
	return CreateDeploymentFlags{
		FleetID:                 fv.GetString("fleet-id"),
		Name:                    fv.GetString("name"),
		Agent:                   fv.GetStringArray("agent"),
		AgentType:               fv.GetString("agent-type"),
		AgentVersion:            fv.GetString("agent-version"),
		ConfigurationVersionIDs: fv.GetStringSlice("configuration-version-ids"),
		Description:             fv.GetString("description"),
		Tags:                    fv.GetStringSlice("tags"),
	}
}

// UpdateDeploymentFlags provides typed access to 'update' (deployment) command flags
type UpdateDeploymentFlags struct {
	DeploymentID            string
	Name                    string
	Description             string
	Agent                   []string
	ConfigurationVersionIDs []string
	Tags                    []string
}

// UpdateDeployment returns typed flags for the 'update' (deployment) command
func (fv *FlagValues) UpdateDeployment() UpdateDeploymentFlags {
	return UpdateDeploymentFlags{
		DeploymentID:            fv.GetString("deployment-id"),
		Name:                    fv.GetString("name"),
		Description:             fv.GetString("description"),
		Agent:                   fv.GetStringArray("agent"),
		ConfigurationVersionIDs: fv.GetStringSlice("configuration-version-ids"),
		Tags:                    fv.GetStringSlice("tags"),
	}
}

// DeployFlags provides typed access to 'deploy' (deployment) command flags
type DeployFlags struct {
	DeploymentID  string
	RingsToDeploy []string
}

// Deploy returns typed flags for the 'deploy' (deployment) command
func (fv *FlagValues) Deploy() DeployFlags {
	return DeployFlags{
		DeploymentID:  fv.GetString("deployment-id"),
		RingsToDeploy: fv.GetStringSlice("rings-to-deploy"),
	}
}

// DeleteDeploymentFlags provides typed access to 'delete' (deployment) command flags
type DeleteDeploymentFlags struct {
	DeploymentID string
}

// DeleteDeployment returns typed flags for the 'delete' (deployment) command
func (fv *FlagValues) DeleteDeployment() DeleteDeploymentFlags {
	return DeleteDeploymentFlags{
		DeploymentID: fv.GetString("deployment-id"),
	}
}

// GetManagedFlags provides typed access to 'get-managed' (entities) command flags
type GetManagedFlags struct {
	EntityType  string
	Limit       int
	IncludeTags bool
}

// GetManaged returns typed flags for the 'get-managed' (entities) command
func (fv *FlagValues) GetManaged() GetManagedFlags {
	return GetManagedFlags{
		EntityType:  fv.GetString("entity-type"),
		Limit:       fv.GetInt("limit"),
		IncludeTags: fv.GetBool("include-tags"),
	}
}

// GetUnassignedFlags provides typed access to 'get-unassigned' (entities) command flags
type GetUnassignedFlags struct {
	EntityType  string
	Limit       int
	IncludeTags bool
}

// GetUnassigned returns typed flags for the 'get-unassigned' (entities) command
func (fv *FlagValues) GetUnassigned() GetUnassignedFlags {
	return GetUnassignedFlags{
		EntityType:  fv.GetString("entity-type"),
		Limit:       fv.GetInt("limit"),
		IncludeTags: fv.GetBool("include-tags"),
	}
}
