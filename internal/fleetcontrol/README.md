# Fleet Control Command Framework

A YAML-driven, modular command framework that eliminates boilerplate and enables rapid development with type safety and declarative validation.

## ⚠️ Important: Rebuild After YAML Changes

**YAML files are embedded at compile time using `//go:embed configs/*.yaml`.**

**After changing any YAML file, you MUST rebuild:**
```bash
go build -o ./bin/darwin/newrelic ./cmd/newrelic
```

Changes to YAML files will NOT take effect until you rebuild the binary!

## 📁 Directory Structure

```
internal/fleetcontrol/
├── README.md                                      # This file
├── TEST_VALIDATION.md                            # Validation testing guide
├── command.go                                    # Main entry point (fleetcontrol command)
├── command_framework.go                          # Core framework (YAML loading, validation)
├── command_flags_generated.go                    # Generated typed flag accessors
├── command_fleet.go                              # Command registration and wiring
├── helpers.go                                    # Shared utility functions
│
├── configs/                                      # YAML configuration files (matching handler names)
│   ├── fleet_management_create.yaml              # Matches fleet_management_create.go
│   ├── fleet_management_update.yaml              # Matches fleet_management_update.go
│   ├── fleet_management_delete.yaml              # Matches fleet_management_delete.go
│   ├── fleet_members_add.yaml                    # Matches fleet_members_add.go
│   ├── fleet_members_remove.yaml                 # Matches fleet_members_remove.go
│   ├── fleet_configuration_create.yaml           # Matches fleet_configuration_create.go
│   ├── fleet_configuration_get.yaml              # Matches fleet_configuration_get.go
│   ├── fleet_configuration_version_list.yaml     # Matches fleet_configuration_version_list.go
│   ├── fleet_configuration_version_add.yaml      # Matches fleet_configuration_version_add.go
│   ├── fleet_configuration_delete.yaml           # Matches fleet_configuration_delete.go
│   ├── fleet_configuration_version_delete.yaml   # Matches fleet_configuration_version_delete.go
│   ├── fleet_deployment_create.yaml              # Matches fleet_deployment_create.go
│   └── fleet_deployment_update.yaml              # Matches fleet_deployment_update.go
│
└── Handler implementation files (one per command)
    ├── fleet_management_create.go                # 'create' command handler
    ├── fleet_management_update.go                # 'update' command handler
    ├── fleet_management_delete.go                # 'delete' command handler
    ├── fleet_members_add.go                      # 'add-members' command handler
    ├── fleet_members_remove.go                   # 'remove-members' command handler
    ├── fleet_configuration_create.go             # 'create-configuration' command handler
    ├── fleet_configuration_get.go                # 'get-configuration' command handler
    ├── fleet_configuration_version_list.go       # 'get-versions' command handler
    ├── fleet_configuration_version_add.go        # 'add-version' command handler
    ├── fleet_configuration_delete.go             # 'delete-configuration' command handler
    ├── fleet_configuration_version_delete.go     # 'delete-version' command handler
    ├── fleet_deployment_create.go                # 'create-deployment' command handler
    └── fleet_deployment_update.go                # 'update-deployment' command handler
```

## 🎯 Key Benefits

- **✅ Zero Hardcoded Strings**: All flag names defined once in YAML
- **✅ Type-Safe Access**: Generated structs with compile-time checking
- **✅ Declarative Validation**: YAML-driven validation (no redundant code)
- **✅ Single Source of Truth**: YAML defines everything
- **✅ Modular Organization**: Each command in its own file
- **✅ Strict File Reading**: `file` type flags read from file paths (with separate flags for inline content)
- **✅ Consistent Output Format**: All commands return JSON with status/error fields
- **✅ SOLID Principles**: Well-organized, readable, maintainable

## 🏗️ Architecture

### 1. **YAML Configuration** (`configs/*.yaml`)

Each command has its own YAML file defining:
- Command metadata (name, description, examples)
- Flags (name, type, required, validation rules)

```yaml
name: create
short: Create a new fleet
flags:
  - name: managed-entity-type
    type: string
    required: true
    validation:
      allowed_values: ["HOST", "KUBERNETESCLUSTER"]
      case_insensitive: true
```

### 2. **Framework** (`command_framework.go`)

The framework:
- Loads all YAML files from `configs/` directory
- Registers flags with cobra automatically
- Validates flag values against YAML rules **before** handlers run
- Provides type-safe flag accessors

### 3. **Generated Types** (`command_flags_generated.go`)

Generated typed structs for each command:

```go
type CreateFlags struct {
    Name              string
    ManagedEntityType string
    // ... all flags as typed fields
}

func (fv *FlagValues) Create() CreateFlags {
    return CreateFlags{
        Name:              fv.GetString("name"),
        ManagedEntityType: fv.GetString("managed-entity-type"),
        // ...
    }
}
```

### 4. **Command Handlers** (`command_<name>.go`)

Each handler:
- Gets typed flags: `f := flags.Create()`
- Uses helper functions for common operations
- Contains pure business logic
- No flag registration boilerplate
- No validation code (handled by framework)

### 5. **Shared Helpers** (`helpers.go`)

Common functions used across commands:

**Organization & Parsing:**
- `GetOrganizationID()`: Get or fetch organization ID
- `ParseTags()`: Parse tag strings in "key:value1,value2" format

**Type Mappers:**
- `MapManagedEntityType()`: Map validated strings to client types
- `MapScopeType()`: Map scope types
- `MapConfigurationMode()`: Map configuration modes

**Output Formatting:**
- `PrintError()`: Wrap errors with status/error fields
- `PrintFleetSuccess()`: Wrap fleet entity with status/error fields
- `PrintFleetListSuccess()`: Wrap fleet list with status/error fields
- `PrintDeleteSuccess()`: Wrap delete result with status/error fields
- `PrintDeleteBulkSuccess()`: Wrap bulk delete results (array)
- `PrintConfigurationSuccess()`: Wrap configuration data with status/error fields
- `PrintConfigurationDeleteSuccess()`: Wrap configuration delete with status/error fields
- `PrintDeploymentSuccess()`: Wrap deployment data with status/error fields
- `PrintMemberSuccess()`: Wrap member operation data with status/error fields

**Output Utilities:**
- `FilterFleetEntity()`: Convert API entity to clean output format
- `FilterFleetEntityFromEntityManagement()`: Convert EntityManagement entity to clean format
- `printJSON()`: Marshal and print JSON with preserved field order (bypasses prettyjson sorting)

## 📋 Available Commands

### Fleet Management Commands

**Create Fleet:**
```bash
newrelic fleetcontrol fleet create \
  --name "Production Fleet" \
  --managed-entity-type "HOST" \
  --description "Fleet for production hosts" \
  --product "Infrastructure" \
  --tags "env:prod,region:us-east-1"
```

**Get Fleet:**
```bash
newrelic fleetcontrol fleet get --id "fleet-123"
```

**Search Fleets:**
```bash
# Get all fleets
newrelic fleetcontrol fleet search

# Search by exact name
newrelic fleetcontrol fleet search --name-equals "Production Fleet"

# Search by name contains
newrelic fleetcontrol fleet search --name-contains "prod"
```

**Update Fleet:**
```bash
newrelic fleetcontrol fleet update \
  --id "fleet-123" \
  --name "Updated Fleet Name" \
  --description "New description"
```

**Delete Fleet:**
```bash
# Delete single fleet
newrelic fleetcontrol fleet delete --id "fleet-123"

# Bulk delete (requires 2+ IDs)
newrelic fleetcontrol fleet delete --ids "fleet-1,fleet-2,fleet-3"
```

### Fleet Member Commands

**Add Members:**
```bash
newrelic fleetcontrol fleet add-members \
  --fleet-id "fleet-123" \
  --ring "production" \
  --entity-ids "entity-1,entity-2,entity-3"
```

**Remove Members:**
```bash
newrelic fleetcontrol fleet remove-members \
  --fleet-id "fleet-123" \
  --ring "production" \
  --entity-ids "entity-1,entity-2"
```

### Configuration Commands

**Create Configuration:**
```bash
# From file (recommended)
newrelic fleetcontrol fleet create-configuration \
  --entity-name "My Config" \
  --agent-type "NRInfra" \
  --managed-entity-type "HOST" \
  --configuration-file-path ./config.json

# Inline content (testing only)
newrelic fleetcontrol fleet create-configuration \
  --entity-name "Test Config" \
  --agent-type "NRInfra" \
  --managed-entity-type "HOST" \
  --configuration-content '{"key": "value"}'
```

**Get Configuration:**
```bash
# Get latest version
newrelic fleetcontrol fleet get-configuration \
  --entity-guid "CONFIG_GUID"

# Get specific version
newrelic fleetcontrol fleet get-configuration \
  --entity-guid "CONFIG_GUID" \
  --version 2

# Get by version entity GUID
newrelic fleetcontrol fleet get-configuration \
  --entity-guid "VERSION_GUID" \
  --mode "ConfigVersionEntity"
```

**List Configuration Versions:**
```bash
newrelic fleetcontrol fleet get-versions \
  --configuration-guid "CONFIG_GUID"
```

**Add Version to Configuration:**
```bash
# From file (recommended)
newrelic fleetcontrol fleet add-version \
  --configuration-guid "CONFIG_GUID" \
  --configuration-file-path ./new-version.json

# Inline content (testing only)
newrelic fleetcontrol fleet add-version \
  --configuration-guid "CONFIG_GUID" \
  --configuration-content '{"updated": "config"}'
```

**Delete Configuration:**
```bash
newrelic fleetcontrol fleet delete-configuration \
  --configuration-guid "CONFIG_GUID"
```

**Delete Configuration Version:**
```bash
newrelic fleetcontrol fleet delete-version \
  --version-guid "VERSION_GUID"
```

### Deployment Commands

**Create Deployment:**
```bash
newrelic fleetcontrol fleet create-deployment \
  --fleet-id "fleet-123" \
  --name "Production Rollout" \
  --configuration-version-ids "version-1,version-2" \
  --description "Rolling out new configuration" \
  --tags "env:prod,release:v1.2.3"
```

**Update Deployment:**
```bash
newrelic fleetcontrol fleet update-deployment \
  --id "deployment-123" \
  --name "Updated Rollout" \
  --configuration-version-ids "version-3,version-4"
```

## 🚀 Adding a New Flag

### Step 1: Update YAML

Edit `configs/<command>.yaml`:

```yaml
flags:
  # ... existing flags ...
  - name: priority
    type: string
    required: false
    description: the priority level
    validation:
      allowed_values: ["LOW", "MEDIUM", "HIGH"]
      case_insensitive: true
```

### Step 2: Update Generated Struct

Edit `command_flags_generated.go`:

```go
type CreateFlags struct {
    // ... existing fields ...
    Priority string  // Add this line
}

func (fv *FlagValues) Create() CreateFlags {
    return CreateFlags{
        // ... existing fields ...
        Priority: fv.GetString("priority"),  // Add this line
    }
}
```

### Step 3: Use in Handler (Optional)

Only if you need the flag in your logic:

```go
func handleFleetCreate(cmd *cobra.Command, args []string, flags *FlagValues) error {
    f := flags.Create()

    // New flag is immediately available!
    if f.Priority != "" {
        // Use it...
    }
}
```

**That's it!** Framework automatically:
- Registers the flag
- Validates against allowed values
- Makes it available through typed struct

## 🔧 Adding a New Command

This guide shows how to add a new command using a real example: the `get-versions` command that was recently added.

### Step 1: Create YAML Config

Create `configs/fleet_configuration_version_list.yaml` (name matches the handler file):

```yaml
name: get-versions
short: Get all versions of a fleet configuration
long: |
  Retrieve all versions of a fleet configuration by configuration GUID.

  This command fetches the version history of a configuration, including
  version numbers, blob IDs, and timestamps.
example: |
  # Get all versions of a configuration
  newrelic fleetcontrol fleet get-versions --configuration-guid "ABC123DEF456"

  # Get versions with explicit organization ID
  newrelic fleetcontrol fleet get-versions --configuration-guid "ABC123DEF456" --organization-id "ORG_ID"
flags:
  - name: configuration-guid
    type: string
    required: true
    description: the configuration entity GUID
  - name: organization-id
    type: string
    description: the organization ID; if not provided, it will be fetched from the API
```

### Step 2: Generate Typed Flags

Edit `command_flags_generated.go`:

```go
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
```

### Step 3: Create Handler File

Create `fleet_configuration_version_list.go`:

```go
package fleetcontrol

import (
    "fmt"
    "github.com/spf13/cobra"
    "github.com/newrelic/newrelic-cli/internal/client"
)

// handleFleetGetConfigurationVersions implements the 'get-versions' command.
func handleFleetGetConfigurationVersions(cmd *cobra.Command, args []string, flags *FlagValues) error {
    // Get typed flag values - no hardcoded strings!
    f := flags.GetVersions()

    // Get organization ID (provided or fetched from API)
    orgID := GetOrganizationID(f.OrganizationID)

    // Call New Relic API
    result, err := client.NRClient.FleetControl.FleetControlGetConfigurationVersions(
        f.ConfigurationGUID,
        orgID,
    )
    if err != nil {
        return PrintError(fmt.Errorf("failed to get configuration versions: %w", err))
    }

    // Validate that versions were returned
    if result == nil || len(result.Versions) == 0 {
        return PrintError(fmt.Errorf("no version details found, please check the GUID of the configuration entity provided"))
    }

    // Print results with status wrapper
    return PrintConfigurationSuccess(result)
}
```

**Key changes from basic pattern:**
- Uses `PrintError()` for consistent error formatting with status/error fields
- Uses `PrintConfigurationSuccess()` to wrap response with status/error fields
- Added validation for empty results
- Removed unused imports (no log, no output package)

### Step 4: Register Handler

Edit `command_fleet.go`:

**Add variable declaration:**
```go
var (
    // ... existing vars ...
    cmdFleetGetConfigurationVersions *cobra.Command  // Add this line
)
```

**Add to handlers map:**
```go
handlers := map[string]CommandHandler{
    // ... existing handlers ...
    "get-versions": handleFleetGetConfigurationVersions,  // Add this line
}
```

**Add to switch statement:**
```go
switch cmdDef.Name {
    // ... existing cases ...
    case "get-versions":
        cmdFleetGetConfigurationVersions = cmd  // Add this case
}
```

**Add to registration:**
```go
cmdFleet.AddCommand(cmdFleetGetConfigurationVersions)  // Add this line
```

### Step 5: Rebuild and Test

```bash
# Rebuild the binary
go build -o ./bin/darwin/newrelic ./cmd/newrelic

# Test the command
./bin/darwin/newrelic fleetcontrol fleet get-versions --help
```

## 📝 Flag Types

### Supported Types

| Type | Description | Example |
|------|-------------|---------|
| `string` | Simple string value | `--name "My Fleet"` |
| `stringSlice` | Multiple values | `--tags "env:prod" --tags "team:platform"` |
| `int` | Integer value | `--version 1` |
| `bool` | Boolean flag | `--force` |
| `file` | File path (strict) | `--configuration-file-path ./config.json` |

### File Flags: Separated for Clarity

Configuration commands use **mutually exclusive** flags for file content:

- **`--configuration-file-path`** (type: `file`): Read from file path (recommended for production)
- **`--configuration-content`** (type: `string`): Inline content (for testing/development/emergency)

**Example (create-configuration):**
```bash
# Recommended: Read from file
newrelic fleetcontrol fleet create-configuration \
  --entity-name "My Config" \
  --agent-type "NRInfra" \
  --managed-entity-type "HOST" \
  --configuration-file-path ./config.json

# Alternative: Inline content (testing only)
newrelic fleetcontrol fleet create-configuration \
  --entity-name "Test Config" \
  --agent-type "NRInfra" \
  --managed-entity-type "HOST" \
  --configuration-content '{"key": "value"}'
```

**Example (add-version):**
```bash
# Recommended: Read from file
newrelic fleetcontrol fleet add-version \
  --configuration-guid "ABC123" \
  --configuration-file-path ./new-version.json

# Alternative: Inline content (testing only)
newrelic fleetcontrol fleet add-version \
  --configuration-guid "ABC123" \
  --configuration-content '{"updated": "config"}'
```

### Validation Options

```yaml
validation:
  allowed_values: ["VALUE1", "VALUE2"]  # Only these values accepted
  case_insensitive: true                # Ignore case when validating
```

**Common validation rules:**

**Agent Types** (case-insensitive):
```yaml
allowed_values: ["NRInfra", "NRDOT", "FluentBit", "NRPrometheusAgent"]
```

**Managed Entity Types** (case-insensitive):
```yaml
allowed_values: ["HOST", "KUBERNETESCLUSTER", "APM"]
```

**Configuration Modes** (case-insensitive):
```yaml
allowed_values: ["ConfigEntity", "ConfigVersionEntity"]
```

### Mutually Exclusive Flags

Some commands have mutually exclusive flags (only one can be used):

**Delete fleet:**
- `--id` (delete single fleet) OR `--ids` (bulk delete multiple fleets)
- Validation: Requires at least one; errors if both provided
- Bulk delete requires 2+ IDs (suggests using `--id` if only one)

**Configuration content:**
- `--configuration-file-path` OR `--configuration-content`
- Validation: Requires exactly one; errors if both or neither provided
- File path is recommended; inline content is for testing/development only

**Search fleets:**
- `--name-equals` OR `--name-contains` (both optional, but mutually exclusive)
- If neither provided, returns all fleets

Validation enforces mutual exclusivity at runtime.

### Special Validations

**Empty Results Validation:**
```go
// get-versions command validates that results are not empty
if result == nil || len(result.Versions) == 0 {
    return PrintError(fmt.Errorf("no version details found, please check the GUID of the configuration entity provided"))
}
```

**Tags Format:**
Tags must be in format `"key:value1,value2"`:
```bash
--tags "env:prod,staging" --tags "team:platform"
# Parsed as: [{Key: "env", Values: ["prod", "staging"]}, {Key: "team", Values: ["platform"]}]
```

## 📤 Response Format

All commands return consistent JSON output with status and error fields.

### Success Response Format

**Most commands** (create, update, delete, add-version, etc.):
```json
{
  "status": "success",
  "error": "",
  "result": {
    "entityGuid": "...",
    "blobId": "...",
    ...
  }
}
```

**Delete operations:**
```json
{
  "status": "success",
  "error": "",
  "id": "deleted-entity-id"
}
```

**Bulk delete operations** (returns array):
```json
[
  {
    "status": "success",
    "error": "",
    "id": "fleet-id-1"
  },
  {
    "status": "failed",
    "error": "failed to delete fleet: not found",
    "id": "fleet-id-2"
  }
]
```

### Failure Response Format

```json
{
  "status": "failed",
  "error": "failed to create fleet: organization ID required"
}
```

### Exceptions (Raw Output)

**Search and Get Configuration** commands return raw data without wrapper (for table formatting):
```json
[
  {
    "id": "fleet-123",
    "name": "Production Fleet",
    "managedEntityType": "HOST",
    ...
  }
]
```

Errors from these commands still use the status/error wrapper.

### Using jq with Responses

**Extract data from success response:**
```bash
# Get entityGuid from create-configuration
newrelic fleetcontrol fleet create-configuration ... | jq -r '.result.entityGuid'

# Get version number
newrelic fleetcontrol fleet create-configuration ... | jq -r '.result.blobVersionEntity.version'

# Get ID from delete
newrelic fleetcontrol fleet delete --id abc123 | jq -r '.id'
```

**Check status before extracting data:**
```bash
# Extract entityGuid only if successful
newrelic fleetcontrol fleet create-configuration ... | \
  jq -r 'select(.status == "success") | .result.entityGuid'

# Show error if failed, otherwise show entityGuid
newrelic fleetcontrol fleet create-configuration ... | \
  jq -r 'if .status == "success" then .result.entityGuid else .error end'
```

**Store result in variable with error handling:**
```bash
OUTPUT=$(newrelic fleetcontrol fleet create-configuration ...)
STATUS=$(echo "$OUTPUT" | jq -r '.status')

if [ "$STATUS" = "success" ]; then
  GUID=$(echo "$OUTPUT" | jq -r '.result.entityGuid')
  echo "Created configuration: $GUID"
else
  ERROR=$(echo "$OUTPUT" | jq -r '.error')
  echo "Failed: $ERROR"
  exit 1
fi
```

## ✅ Validation Flow

**CRITICAL:** Validation happens in the framework **before** handlers run.

1. **YAML defines allowed values**: `allowed_values: ["HOST", "KUBERNETESCLUSTER"]`
2. **Framework validates**: Checks flag value against YAML rules
3. **Handler receives validated value**: No need to validate again!

### Mapper Functions

Handlers use mapper functions to convert validated strings to client library types:

```go
// YAML has already validated that f.ManagedEntityType is "HOST" or "KUBERNETESCLUSTER"
entityType, err := mapManagedEntityType(f.ManagedEntityType)
// Now we have the typed value for the API client
```

**Why mappers?**
- Centralized type conversion
- Error handling (though validation should prevent errors)
- Clear separation: YAML validates, code maps

### ⚠️ Common Mistake

**DON'T** add switch statements in handlers:

```go
// ❌ BAD - Bypasses YAML validation
switch strings.ToUpper(f.ManagedEntityType) {
case "HOST":
    entityType = ...
case "KUBERNETESCLUSTER":
    entityType = ...
default:
    return fmt.Errorf("invalid type")  // This hides YAML validation failures!
}
```

**DO** use mapper functions:

```go
// ✅ GOOD - Trusts YAML validation
entityType, err := mapManagedEntityType(f.ManagedEntityType)
if err != nil {
    return err  // Should never happen after YAML validation
}
```

## 🧪 Testing Validation

**⚠️ CRITICAL: YAML files are embedded at compile time!**

The YAML files are embedded into the binary using `//go:embed configs/*.yaml`.
**You MUST rebuild the binary after changing any YAML file!**

To test that YAML validation is working:

1. Edit a YAML file with invalid allowed_values
2. **REBUILD THE BINARY** (critical step!)
3. Try to use a valid value (not in the broken YAML)
4. Should get validation error

Example:
```yaml
# Change in configs/create.yaml
validation:
  allowed_values: ["TYPO", "WRONG"]  # Intentionally wrong
```

```bash
# MUST REBUILD after changing YAML!
go build -o ./bin/darwin/newrelic ./cmd/newrelic

# This should now FAIL with validation error
./bin/darwin/newrelic fleetcontrol fleet create --managed-entity-type KUBERNETESCLUSTER
# Error: invalid value 'KUBERNETESCLUSTER' for flag --managed-entity-type: must be one of [TYPO, WRONG]
```

**See TEST_VALIDATION.md for detailed testing instructions.**

**Note**: Validation runs silently. You'll only see output if validation fails.

## 📚 Code Organization Principles

### SOLID Principles Applied

1. **Single Responsibility**
   - Each command file handles one command
   - Each YAML file defines one command
   - Helpers file contains shared utilities
   - Framework file handles YAML loading and validation

2. **Open/Closed**
   - Add new commands without modifying existing ones
   - Extend through YAML, not code changes

3. **Dependency Inversion**
   - Commands depend on abstractions (FlagValues interface)
   - Not on concrete flag implementations

### Code Readability

- **Comprehensive comments**: Every function documents purpose, parameters, returns
- **Clear naming**: `handleFleetCreate`, `CreateFlags`, `mapManagedEntityType`
- **Logical grouping**: Commands grouped by functionality
- **Minimal complexity**: Handlers focus on business logic only

## 🔍 For New Engineers

### Understanding the Flow

1. **User runs**: `newrelic fleetcontrol fleet create --name "Test" --managed-entity-type HOST`

2. **Framework loads**: Reads `configs/create.yaml`

3. **Framework registers**: Creates cobra flags from YAML

4. **User executes**: Command is invoked

5. **Framework validates**: Checks `managed-entity-type` is in `allowed_values`

6. **Handler runs**: `handleFleetCreate` receives validated flags

7. **Handler logic**:
   - Gets typed flags: `f := flags.Create()`
   - Maps to client types: `entityType, err := mapManagedEntityType(f.ManagedEntityType)`
   - Calls API: `client.NRClient.FleetControl.FleetControlCreateFleet(...)`

### Finding Your Way Around

- **Need to add a flag?** → Update YAML + regenerate `command_flags_generated.go`
- **Need to change validation?** → Update YAML validation section + rebuild binary
- **Need to add business logic?** → Edit handler in `fleet_<name>.go`
- **Need a shared function?** → Add to `helpers.go`
- **Need to change output format?** → Update wrapper functions in `helpers.go`
- **Need to understand framework?** → Read `command_framework.go`
- **Need to see all commands?** → Check "Available Commands" section above
- **Need response format examples?** → Check "Response Format" section above
- **Having issues?** → Check "Troubleshooting" section above

## 🔧 Go Client Integration

### Local Client Development

The CLI uses a local Go Client for development via go.mod replace directive:

```go
// go.mod
replace github.com/newrelic/newrelic-client-go/v2 => ../newrelic-client-go
```

### Recent Go Client Fixes

**Delete Operations Fixed:**
Both delete operations now correctly handle empty API responses (204 No Content):

1. **`FleetControlDeleteConfiguration`** - Pass `nil` for response body (not `&resp`)
2. **`FleetControlDeleteConfigurationVersion`** - Pass `nil` for response body (not `&resp`)

**Before (caused "unexpected end of JSON input"):**
```go
resp := DeleteBlobResponse{}
_, err := a.client.DeleteWithContext(ctx, url, nil, &resp)  // ❌ Tries to parse empty response
```

**After (correctly handles empty response):**
```go
_, err := a.client.DeleteWithContext(ctx, url, nil, nil)  // ✅ No JSON parsing attempted
```

These fixes are in `../newrelic-client-go/pkg/fleetcontrol/fleetcontrol_configurations.go`.

## 🐛 Troubleshooting

| Problem | Solution |
|---------|----------|
| **Validation not working** | **REBUILD BINARY!** YAML files are embedded at compile time. Run: `go build -o ./bin/darwin/newrelic ./cmd/newrelic` |
| YAML changes not taking effect | YAML is embedded at compile time - rebuild required |
| Flag not available in struct | Regenerate `command_flags_generated.go` |
| Validation passes but mapper fails | YAML validation passed (value in allowed_values), but mapper doesn't recognize it. Sync YAML with mapper. |
| File not reading | File flags now strictly read from file path. Use `--configuration-content` for inline content. |
| "unexpected end of JSON input" on delete | Fixed in Go Client - ensure you're using local client with `replace` directive in go.mod |
| Mutually exclusive flag error | Check that only one of the mutually exclusive flags is provided (e.g., `--id` OR `--ids`, not both) |
| "no version details found" error | Configuration GUID is invalid or has no versions. Verify GUID is correct. |
| Command not found | Check handler is registered in `command_fleet.go` |
| Type mismatch error | Ensure YAML type matches accessor method |
| JSON output not sorted | Status/error fields intentionally appear first (field order preserved via custom JSON marshaling) |
| Flag syntax error | Flags require `--` prefix: `--flag-name value` or `--flag-name=value` (not `flag-name=value`) |

## 🎓 Best Practices

1. **Always use typed accessors** - `f := flags.Create()` not `flags.GetString("name")`
2. **Trust YAML validation** - Don't re-validate in handlers
3. **Use mapper functions** - Centralize type conversions
4. **Use helper functions for output** - `PrintError()`, `PrintFleetSuccess()`, etc. for consistent formatting
5. **No log statements in handlers** - Use response wrappers (status/error fields) instead of log.Info/log.Debug
6. **Recommend file paths over inline content** - Document that `--configuration-file-path` is for production, `--configuration-content` is for testing
7. **Validate empty results** - Check for empty arrays/nil results and return meaningful errors
8. **Comment your code** - Explain the "why", not just the "what"
9. **Keep handlers focused** - Pure business logic, no boilerplate
10. **Update generated code** - When YAML changes, update the generated struct
11. **Test with jq** - Verify response structure works with common jq patterns
12. **Rebuild after YAML changes** - YAML files are embedded at compile time

## 📖 Additional Resources

- **Cobra Documentation**: https://github.com/spf13/cobra
- **YAML Specification**: https://yaml.org/spec/1.2/spec.html
- **Go Embedding**: https://pkg.go.dev/embed

---

## 📝 Recent Updates

**January 2026:**
- Split hybrid file flags into mutually exclusive `--configuration-file-path` and `--configuration-content`
- Renamed `--entity-guid` to `--configuration-guid` in get-versions for clarity
- Added consistent status/error response wrappers for all commands
- Added empty results validation for get-versions command
- Fixed Go Client delete operations to handle empty responses
- Removed log.Info/log.Debug statements from handlers
- Added comprehensive jq usage examples
- Updated all documentation with current examples

**Last Updated**: January 27, 2026
**Maintainer**: Virtuoso / Observability as Code
