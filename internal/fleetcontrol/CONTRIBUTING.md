# Contributing to Fleet Control

This guide covers the technical details of the YAML-driven command framework and how to contribute new commands.

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
├── README.md                                      # User guide (command reference)
├── CONTRIBUTING.md                                # This file (technical guide)
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
    ├── fleet_management_get.go                   # 'get' command handler
    ├── fleet_management_search.go                # 'search' command handler
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

### 4. **Command Handlers** (`fleet_*.go`)

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

This guide shows how to add a new command using a real example: the `get-versions` command.

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

**Key patterns:**
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

## ✅ Validation Flow

**CRITICAL:** Validation happens in the framework **before** handlers run.

1. **YAML defines allowed values**: `allowed_values: ["HOST", "KUBERNETESCLUSTER"]`
2. **Framework validates**: Checks flag value against YAML rules
3. **Handler receives validated value**: No need to validate again!

### Mapper Functions

Handlers use mapper functions to convert validated strings to client library types:

```go
// YAML has already validated that f.ManagedEntityType is "HOST" or "KUBERNETESCLUSTER"
entityType, err := MapManagedEntityType(f.ManagedEntityType)
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
entityType, err := MapManagedEntityType(f.ManagedEntityType)
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
# Change in configs/fleet_management_create.yaml
validation:
  allowed_values: ["TYPO", "WRONG"]  # Intentionally wrong
```

```bash
# MUST REBUILD after changing YAML!
go build -o ./bin/darwin/newrelic ./cmd/newrelic

# This should now FAIL with validation error
./bin/darwin/newrelic fleetcontrol fleet create --managed-entity-type KUBERNETESCLUSTER
# Error response will have: "status": "failed", "error": "invalid value 'KUBERNETESCLUSTER'..."
```

**See TEST_VALIDATION.md for detailed testing instructions.**

**Note**: Validation runs silently. You'll only see output if validation fails.

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

### Failure Response Format

```json
{
  "status": "failed",
  "error": "failed to create fleet: organization ID required"
}
```

### Field Ordering

Status and error fields are **always first** in the output. This is achieved through custom JSON marshaling that preserves struct field order and bypasses the go-prettyjson alphabetical sorting.

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
- **Clear naming**: `handleFleetCreate`, `CreateFlags`, `MapManagedEntityType`
- **Logical grouping**: Commands grouped by functionality
- **Minimal complexity**: Handlers focus on business logic only

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
| Command not found | Check handler is registered in `command_fleet.go` |
| Type mismatch error | Ensure YAML type matches accessor method |
| JSON output not sorted | Status/error fields intentionally appear first (field order preserved via custom JSON marshaling) |

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

## 🔍 For New Engineers

### Understanding the Flow

1. **User runs**: `newrelic fleetcontrol fleet create --name "Test" --managed-entity-type HOST`

2. **Framework loads**: Reads `configs/fleet_management_create.yaml`

3. **Framework registers**: Creates cobra flags from YAML

4. **User executes**: Command is invoked

5. **Framework validates**: Checks `managed-entity-type` is in `allowed_values`

6. **Handler runs**: `handleFleetCreate` receives validated flags

7. **Handler logic**:
   - Gets typed flags: `f := flags.Create()`
   - Maps to client types: `entityType, err := MapManagedEntityType(f.ManagedEntityType)`
   - Calls API: `client.NRClient.FleetControl.FleetControlCreateFleet(...)`

### Finding Your Way Around

- **Need to add a flag?** → Update YAML + regenerate `command_flags_generated.go`
- **Need to change validation?** → Update YAML validation section + rebuild binary
- **Need to add business logic?** → Edit handler in `fleet_*.go`
- **Need a shared function?** → Add to `helpers.go`
- **Need to change output format?** → Update wrapper functions in `helpers.go`
- **Need to understand framework?** → Read `command_framework.go`
- **Need to see all commands?** → Check README.md
- **Having issues?** → Check "Troubleshooting" section above

## 📖 Additional Resources

- **Cobra Documentation**: https://github.com/spf13/cobra
- **YAML Specification**: https://yaml.org/spec/1.2/spec.html
- **Go Embedding**: https://pkg.go.dev/embed

---

**Last Updated**: January 27, 2026
**Maintainer**: Virtuoso / Observability as Code
