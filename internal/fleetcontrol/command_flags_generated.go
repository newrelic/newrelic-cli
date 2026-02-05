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
	ID  string
	Ids []string
}

// Delete returns typed flags for the 'delete' command
func (fv *FlagValues) Delete() DeleteFlags {
	return DeleteFlags{
		ID:  fv.GetString("id"),
		Ids: fv.GetStringSlice("ids"),
	}
}

// GetFlags provides typed access to 'get' command flags
type GetFlags struct {
	ID string
}

// Get returns typed flags for the 'get' command
func (fv *FlagValues) Get() GetFlags {
	return GetFlags{
		ID: fv.GetString("id"),
	}
}

// SearchFlags provides typed access to 'search' command flags
type SearchFlags struct {
	NameEquals   string
	NameContains string
}

// Search returns typed flags for the 'search' command
func (fv *FlagValues) Search() SearchFlags {
	return SearchFlags{
		NameEquals:   fv.GetString("name-equals"),
		NameContains: fv.GetString("name-contains"),
	}
}

// AddMembersFlags provides typed access to 'add-members' command flags
type AddMembersFlags struct {
	FleetID   string
	Ring      string
	EntityIDs []string
}

// AddMembers returns typed flags for the 'add-members' command
func (fv *FlagValues) AddMembers() AddMembersFlags {
	return AddMembersFlags{
		FleetID:   fv.GetString("fleet-id"),
		Ring:      fv.GetString("ring"),
		EntityIDs: fv.GetStringSlice("entity-ids"),
	}
}

// RemoveMembersFlags provides typed access to 'remove-members' command flags
type RemoveMembersFlags struct {
	FleetID   string
	Ring      string
	EntityIDs []string
}

// RemoveMembers returns typed flags for the 'remove-members' command
func (fv *FlagValues) RemoveMembers() RemoveMembersFlags {
	return RemoveMembersFlags{
		FleetID:   fv.GetString("fleet-id"),
		Ring:      fv.GetString("ring"),
		EntityIDs: fv.GetStringSlice("entity-ids"),
	}
}

// ListMembersFlags provides typed access to 'list-members' command flags
type ListMembersFlags struct {
	FleetID     string
	Ring        string
	IncludeTags bool
}

// ListMembers returns typed flags for the 'list-members' command
func (fv *FlagValues) ListMembers() ListMembersFlags {
	return ListMembersFlags{
		FleetID:     fv.GetString("fleet-id"),
		Ring:        fv.GetString("ring"),
		IncludeTags: fv.GetBool("include-tags"),
	}
}

// CreateConfigurationFlags provides typed access to 'create-configuration' command flags
type CreateConfigurationFlags struct {
	EntityName             string
	AgentType              string
	ManagedEntityType      string
	OrganizationID         string
	ConfigurationFilePath  string // File content from path
	ConfigurationContent   string // Inline content
}

// CreateConfiguration returns typed flags for the 'create-configuration' command
func (fv *FlagValues) CreateConfiguration() (CreateConfigurationFlags, error) {
	configFilePath, err := fv.GetFile("configuration-file-path")
	if err != nil {
		return CreateConfigurationFlags{}, err
	}
	return CreateConfigurationFlags{
		EntityName:            fv.GetString("entity-name"),
		AgentType:             fv.GetString("agent-type"),
		ManagedEntityType:     fv.GetString("managed-entity-type"),
		OrganizationID:        fv.GetString("organization-id"),
		ConfigurationFilePath: configFilePath,
		ConfigurationContent:  fv.GetString("configuration-content"),
	}, nil
}

// GetConfigurationFlags provides typed access to 'get-configuration' command flags
type GetConfigurationFlags struct {
	EntityGUID     string
	OrganizationID string
	Mode           string
	Version        int
}

// GetConfiguration returns typed flags for the 'get-configuration' command
func (fv *FlagValues) GetConfiguration() GetConfigurationFlags {
	return GetConfigurationFlags{
		EntityGUID:     fv.GetString("entity-guid"),
		OrganizationID: fv.GetString("organization-id"),
		Mode:           fv.GetString("mode"),
		Version:        fv.GetInt("version"),
	}
}

// GetVersionsFlags provides typed access to 'get-versions' command flags
type GetVersionsFlags struct {
	ConfigurationGUID string
	OrganizationID    string
}

// GetVersions returns typed flags for the 'get-versions' command
func (fv *FlagValues) GetVersions() GetVersionsFlags {
	return GetVersionsFlags{
		ConfigurationGUID: fv.GetString("configuration-guid"),
		OrganizationID:    fv.GetString("organization-id"),
	}
}

// AddVersionFlags provides typed access to 'add-version' command flags
type AddVersionFlags struct {
	ConfigurationGUID     string
	OrganizationID        string
	ConfigurationFilePath string // File content from path
	ConfigurationContent  string // Inline content
}

// AddVersion returns typed flags for the 'add-version' command
func (fv *FlagValues) AddVersion() (AddVersionFlags, error) {
	configFilePath, err := fv.GetFile("configuration-file-path")
	if err != nil {
		return AddVersionFlags{}, err
	}
	return AddVersionFlags{
		ConfigurationGUID:     fv.GetString("configuration-guid"),
		OrganizationID:        fv.GetString("organization-id"),
		ConfigurationFilePath: configFilePath,
		ConfigurationContent:  fv.GetString("configuration-content"),
	}, nil
}

// DeleteConfigurationFlags provides typed access to 'delete-configuration' command flags
type DeleteConfigurationFlags struct {
	ConfigurationGUID string
	OrganizationID    string
}

// DeleteConfiguration returns typed flags for the 'delete-configuration' command
func (fv *FlagValues) DeleteConfiguration() DeleteConfigurationFlags {
	return DeleteConfigurationFlags{
		ConfigurationGUID: fv.GetString("configuration-guid"),
		OrganizationID:    fv.GetString("organization-id"),
	}
}

// DeleteVersionFlags provides typed access to 'delete-version' command flags
type DeleteVersionFlags struct {
	VersionGUID    string
	OrganizationID string
}

// DeleteVersion returns typed flags for the 'delete-version' command
func (fv *FlagValues) DeleteVersion() DeleteVersionFlags {
	return DeleteVersionFlags{
		VersionGUID:    fv.GetString("version-guid"),
		OrganizationID: fv.GetString("organization-id"),
	}
}

// CreateDeploymentFlags provides typed access to 'create-deployment' command flags
type CreateDeploymentFlags struct {
	FleetID                 string
	Name                    string
	ConfigurationVersionIDs []string
	Description             string
	Tags                    []string
}

// CreateDeployment returns typed flags for the 'create-deployment' command
func (fv *FlagValues) CreateDeployment() CreateDeploymentFlags {
	return CreateDeploymentFlags{
		FleetID:                 fv.GetString("fleet-id"),
		Name:                    fv.GetString("name"),
		ConfigurationVersionIDs: fv.GetStringSlice("configuration-version-ids"),
		Description:             fv.GetString("description"),
		Tags:                    fv.GetStringSlice("tags"),
	}
}

// UpdateDeploymentFlags provides typed access to 'update-deployment' command flags
type UpdateDeploymentFlags struct {
	ID                      string
	Name                    string
	Description             string
	ConfigurationVersionIDs []string
	Tags                    []string
}

// UpdateDeployment returns typed flags for the 'update-deployment' command
func (fv *FlagValues) UpdateDeployment() UpdateDeploymentFlags {
	return UpdateDeploymentFlags{
		ID:                      fv.GetString("id"),
		Name:                    fv.GetString("name"),
		Description:             fv.GetString("description"),
		ConfigurationVersionIDs: fv.GetStringSlice("configuration-version-ids"),
		Tags:                    fv.GetStringSlice("tags"),
	}
}
