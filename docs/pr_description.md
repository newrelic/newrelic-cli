# Add Local Terraform Workspace Management Commands for NRQL Drop Rules Migration

## Problem Statement

With the upcoming end-of-life (EOL) of NRQL Drop Rules scheduled for January 7, 2026, customers managing `newrelic_nrql_drop_rule` resources via Terraform need a streamlined way to migrate to `newrelic_pipeline_cloud_rule` resources. While the existing `tf-importgen-ci` command addresses CI/CD environments, there was no solution for developers working directly with local Terraform workspaces who need:

- **Granular control** over each step of the migration process
- **Real-time feedback** and validation at each stage  
- **Interactive execution** with the ability to review changes before applying
- **Safe state management** that doesn't accidentally destroy active drop rules
- **Support for both Terraform and OpenTofu** workflows

This PR introduces three new commands that provide a comprehensive local workspace migration solution, allowing developers to migrate from NRQL Drop Rules to Pipeline Cloud Rules with precise control and safety guarantees.

## Solution Overview

This PR adds three specialized commands under `newrelic migrate nrqldroprules` that work together to provide a complete local migration workflow:

1. **`tf-update`** - Updates existing drop rule resources to include Pipeline Cloud Rule entity IDs
2. **`tf-importgen`** - Generates and executes import configurations for Pipeline Cloud Rules  
3. **`tf-delist`** - Safely removes drop rule resources from state without destroying actual resources

### Key Features

- ✅ **Interactive UX** with clear prompts and real-time command output
- ✅ **Auto-approve handling** to avoid Terraform interactive prompt issues
- ✅ **Comprehensive validation** at each step with detailed error reporting
- ✅ **OpenTofu support** via `--tofu` flag across all commands
- ✅ **Safe state management** with prominent warnings and confirmation prompts
- ✅ **Automatic cleanup** of temporary files and configurations
- ✅ **Detailed post-migration instructions** for completing the workflow

## Commands Added

### 1. `newrelic migrate nrqldroprules tf-update`

**Purpose**: Refreshes NRQL drop rule resources in Terraform state to populate `pipeline_cloud_rule_entity_id` values required for migration.

**Key Features**:
- Automatically discovers drop rule resources in Terraform state
- Validates New Relic provider version (≥3.68.0 required)
- Executes targeted `terraform apply -refresh-only` with auto-approve
- Provides real-time command output and validation results
- Supports resource filtering via `--resourceIdentifiers`

**Usage Examples**:
```bash
# Update all drop rules in current directory
newrelic migrate nrqldroprules tf-update

# Update with OpenTofu
newrelic migrate nrqldroprules tf-update --tofu

# Update specific resources
newrelic migrate nrqldroprules tf-update --resourceIdentifiers resource1,resource2
```

**Flags**:
- `--workspacePath` - Path to Terraform workspace (default: current directory)
- `--resourceIdentifiers` - Specific resource identifiers to update
- `--skipResponseToPrompt` - Skip user prompts (answers 'N' to all)
- `--tofu` - Use OpenTofu instead of Terraform

### 2. `newrelic migrate nrqldroprules tf-importgen`

**Purpose**: Generates Terraform import blocks for Pipeline Cloud Rules and optionally executes the import process.

**Key Features**:
- Validates that drop rules contain required `pipeline_cloud_rule_entity_id` values
- Generates import configuration blocks automatically
- Executes `terraform plan -generate-config-out` and `terraform apply` with auto-approve
- Supports both state-based discovery and manual ID specification
- Automatically saves import configs to files or displays in terminal

**Usage Examples**:
```bash
# Generate and execute import for discovered resources
newrelic migrate nrqldroprules tf-importgen

# Save import config to specific file
newrelic migrate nrqldroprules tf-importgen --fileName imports.tf

# Use specific Pipeline Cloud Rule IDs
newrelic migrate nrqldroprules tf-importgen --pipelineCloudRuleIDs id1,id2
```

**Flags**:
- `--workspacePath` - Path to Terraform workspace (default: current directory)
- `--pipelineCloudRuleIDs` - Specific Pipeline Cloud Rule IDs to import
- `--fileName` - File name for import configuration (optional)
- `--skipResponseToPrompt` - Skip user prompts (answers 'N' to all)
- `--tofu` - Use OpenTofu instead of Terraform

### 3. `newrelic migrate nrqldroprules tf-delist`

**Purpose**: Safely removes NRQL drop rule resources from Terraform state without destroying the actual resources in New Relic.

**Key Features**:
- **Safe-only operation** - uses `terraform state rm`, never destroys actual resources
- Automatic cleanup of temporary import configuration files
- Comprehensive post-execution instructions for configuration file cleanup
- Prominent safety warnings and confirmations throughout the process
- Validates successful delisting with before/after state comparison

**Usage Examples**:
```bash
# Delist all drop rule resources
newrelic migrate nrqldroprules tf-delist

# Delist specific resources
newrelic migrate nrqldroprules tf-delist --resourceIdentifiers resource1,resource2

# Use with OpenTofu
newrelic migrate nrqldroprules tf-delist --tofu
```

**Flags**:
- `--workspacePath` - Path to Terraform workspace (default: current directory)  
- `--resourceIdentifiers` - Specific resource identifiers to delist
- `--skipResponseToPrompt` - Skip user prompts (answers 'N' to all)
- `--tofu` - Use OpenTofu instead of Terraform

## Technical Implementation

### Shared Infrastructure

- **CommandContext**: Centralized context management for tool configuration, workspace paths, and command types
- **Auto-approve handling**: Resolves Terraform interactive prompt issues by automatically adding `-auto-approve` flags
- **Real-time output**: Commands stream Terraform/OpenTofu output directly to user for immediate feedback
- **Comprehensive validation**: Provider version checks, tool version validation, and state consistency verification

### Safety Features

- **Multiple confirmation prompts** before executing state-changing operations
- **Prominent warning messages** about the nature of each operation (especially for delisting)
- **Automatic file cleanup** to prevent workspace pollution
- **Detailed post-execution guidance** for completing migration steps

### Error Handling

- **Graceful state failure handling** with fallback to manual resource specification
- **Clear error messages** with specific resolution instructions
- **Comprehensive troubleshooting guidance** in generated documentation

## Migration Workflow

The typical user workflow follows this pattern:

```bash
# Step 1: Update drop rules to get Pipeline Cloud Rule entity IDs
newrelic migrate nrqldroprules tf-update

# Step 2: Generate and execute import configuration  
newrelic migrate nrqldroprules tf-importgen

# Step 3: Remove drop rules from Terraform state management
newrelic migrate nrqldroprules tf-delist

# Step 4: Comment out drop rule configurations (manual step with detailed instructions)
```

## Documentation

This PR includes comprehensive documentation:

- **`tf_importgen_guide.md`**: Complete user guide with detailed examples, troubleshooting, and best practices
- **Inline help**: Detailed command descriptions, examples, and flag documentation
- **Error messages**: Specific guidance for common issues and their resolutions

## Testing Considerations

These commands are designed to be tested in isolated Terraform workspaces with:
- Valid New Relic provider configuration
- Existing `newrelic_nrql_drop_rule` resources in state
- Proper environment variable setup (`NEW_RELIC_API_KEY`, `NEW_RELIC_ACCOUNT_ID`)

## Breaking Changes

None. This PR only adds new commands and doesn't modify existing functionality.

## Related Issues

This PR addresses the need for local workspace migration tools as a complement to the existing CI/CD automation (`tf-importgen-ci`) for the NRQL Drop Rules EOL scheduled for January 7, 2026.
