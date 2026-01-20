# New Relic CLI: `newrelic_nrql_drop_rule` -> `newrelic_pipeline_cloud_rule` Migration Guide for Local Terraform Workspaces

This guide describes a comprehensive set of local Terraform workspace management commands designed to assist in migrating `newrelic_nrql_drop_rule` resources (Drop Rules, managed via Terraform) to `newrelic_pipeline_cloud_rule` resources in local development environments. This migration is necessary due to the upcoming end-of-life (EOL) of NRQL Drop Rules, scheduled for June 30, 2026. After this date, any Drop Rules managed via the `newrelic_nrql_drop_rule` Terraform resource will no longer function. For a general overview and more details on the EOL, its implications, and the required actions to replace `newrelic_nrql_drop_rule` resources, refer to [this detailed article](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/guides/drop_rules_eol_guide).

Unlike the CI/CD automation workflow (covered in the `tf-importgen-ci` guide), this guide focuses on **interactive local workspace management** using three specialized commands that provide granular control over the migration process:

- **`tf-update`** - Updates existing NRQL drop rules in Terraform state to include `pipeline_cloud_rule_entity_id` values
- **`tf-importgen`** - Generates import configuration blocks for Pipeline Cloud Rules based on updated drop rule resources  
- **`tf-delist`** - Safely removes NRQL drop rules from Terraform state without destroying the actual resources in New Relic

These commands are designed for developers working directly with Terraform workspaces who need precise control over each step of the migration process, real-time feedback, and the ability to validate changes at each stage. The commands support both Terraform and OpenTofu, with automatic approval handling for seamless execution.

## Overview

The local workspace migration commands provide a step-by-step approach to migrating from NRQL Drop Rules to Pipeline Cloud Rules. Each command serves a specific purpose in the migration workflow and can be used independently or as part of a complete migration sequence. All commands are designed to be executed within your existing Terraform workspace directory and will interact directly with your Terraform state and configuration files.

The typical migration workflow follows this pattern:
1. **Update**: Refresh existing drop rule resources to obtain Pipeline Cloud Rule entity IDs
2. **Import**: Generate and execute import configurations for Pipeline Cloud Rules
3. **Delist**: Remove legacy drop rule resources from Terraform state management

**Important Note**: It is strongly recommended to run these commands directly in your Terraform workspace directory where your `newrelic_nrql_drop_rule` resources are defined. This ensures proper state management and configuration file handling.

## Prerequisites

### Technical Requirements

Before using any of the migration commands, ensure that the following technical requirements are met in your local environment:

- **New Relic CLI**: Latest version installed and accessible in your PATH
- **Terraform/OpenTofu**: Version 1.5 or higher must be installed for import generation support
- **New Relic Terraform Provider**: Version 3.68.0 or higher is required for `pipeline_cloud_rule_entity_id` support
- **Environment Variables**: The following must be set:
  - `NEW_RELIC_API_KEY` (required) - Your New Relic User API key with appropriate permissions
  - `NEW_RELIC_ACCOUNT_ID` (required) - The New Relic account ID where your Drop Rules are located
  - `NEW_RELIC_REGION` (optional) - Set to 'US' or 'EU' based on your account region (defaults to 'US')
- **Terraform Workspace**: A valid Terraform workspace with existing `newrelic_nrql_drop_rule` resources

### Workspace Preparation

Ensure your Terraform workspace is properly initialized and contains valid `newrelic_nrql_drop_rule` resources. The commands will automatically detect these resources in your Terraform state.

## Quick Start: Three-Command Workflow

Here's a brief overview of what each command accomplishes and their typical usage:

### 1. Update Command (`tf-update`)
**Purpose**: Refreshes NRQL drop rule resources in Terraform state to populate `pipeline_cloud_rule_entity_id` values needed for Pipeline Cloud Rule migration.

```bash
# Basic usage - run in your Terraform workspace directory
newrelic migrate nrqldroprules tf-update

# With OpenTofu
newrelic migrate nrqldroprules tf-update --tofu
```

**What it does**: Executes `terraform apply -refresh-only` on drop rule resources to update state with the latest Pipeline Cloud Rule entity IDs from New Relic.

### 2. Import Generation Command (`tf-importgen`)
**Purpose**: Generates Terraform import blocks for Pipeline Cloud Rules and optionally executes the import process to bring existing Pipeline Cloud Rules under Terraform management.

```bash
# Basic usage - generates import configuration
newrelic migrate nrqldroprules tf-importgen

# With OpenTofu
newrelic migrate nrqldroprules tf-importgen --tofu
```

**What it does**: Creates import blocks based on `pipeline_cloud_rule_entity_id` values from updated drop rules, then runs `terraform plan -generate-config-out` and `terraform apply` to import Pipeline Cloud Rules.

### 3. Delist Command (`tf-delist`)
**Purpose**: Safely removes NRQL drop rule resources from Terraform state without destroying the actual drop rules in New Relic, allowing you to stop managing them via Terraform.

```bash
# Basic usage - removes resources from state
newrelic migrate nrqldroprules tf-delist

# With OpenTofu
newrelic migrate nrqldroprules tf-delist --tofu
```

**What it does**: Executes `terraform state rm` commands to remove drop rule resources from state while keeping the actual resources active in New Relic.

## OpenTofu Support

All three commands support OpenTofu as an alternative to Terraform by using the `--tofu` flag. When this flag is specified:

- Commands will use `tofu` instead of `terraform` for all operations
- Version requirements and compatibility checks will be performed for OpenTofu
- All functionality remains identical, with tool-specific messaging and validation

**Recommendation**: Run these commands directly in your Terraform/OpenTofu workspace directory to ensure proper state file access and configuration management.

## Detailed Command Reference

### tf-update Command

The `tf-update` command refreshes existing NRQL drop rule resources in your Terraform state to include the `pipeline_cloud_rule_entity_id` attribute, which is essential for the migration process.

#### Command Syntax
```bash
newrelic migrate nrqldroprules tf-update [flags]
```

#### Flags

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--workspacePath` | string | optional | Path to Terraform workspace (defaults to current directory) |
| `--resourceIdentifiers` | string slice | optional | List of specific resource identifiers to update |
| `--skipResponseToPrompt` | boolean | optional | Skip all user prompts (answers 'N' to all prompts) |
| `--tofu` | boolean | optional | Use OpenTofu instead of Terraform |

#### Usage Examples

```bash
# Update all drop rules in current directory
newrelic migrate nrqldroprules tf-update

# Update specific resources
newrelic migrate nrqldroprules tf-update \
  --resourceIdentifiers newrelic_nrql_drop_rule.example1,newrelic_nrql_drop_rule.example2

# Use custom workspace path
newrelic migrate nrqldroprules tf-update \
  --workspacePath /path/to/terraform/workspace

# Automated execution without prompts
newrelic migrate nrqldroprules tf-update --skipResponseToPrompt

# Use with OpenTofu
newrelic migrate nrqldroprules tf-update --tofu
```

#### What tf-update Does

1. **State Discovery**: Scans Terraform state for `newrelic_nrql_drop_rule` resources
2. **Provider Validation**: Checks that New Relic provider version ≥ 3.68.0 for `pipeline_cloud_rule_entity_id` support
3. **Command Generation**: Creates targeted `terraform plan -refresh-only` and `terraform apply -refresh-only` commands
4. **User Interaction**: Displays commands and prompts for execution confirmation
5. **Execution**: Runs plan to show changes, then applies refresh with auto-approve
6. **Validation**: Verifies that resources now contain `pipeline_cloud_rule_entity_id` values

#### Output Example
```
✅ Found 3 NRQL drop rule resources in Terraform state

📋 Resources to be updated:
  1. newrelic_nrql_drop_rule.log_filter
  2. newrelic_nrql_drop_rule.error_filter  
  3. newrelic_nrql_drop_rule.debug_filter

📝 Generated Terraform commands for resource updates:
1. terraform plan -refresh-only -target=newrelic_nrql_drop_rule.log_filter -target=newrelic_nrql_drop_rule.error_filter -target=newrelic_nrql_drop_rule.debug_filter
2. terraform apply -refresh-only -target=newrelic_nrql_drop_rule.log_filter -target=newrelic_nrql_drop_rule.error_filter -target=newrelic_nrql_drop_rule.debug_filter

✅ All 3 NRQL drop rule resources have been successfully updated with pipeline_cloud_rule_entity_id
```

### tf-importgen Command

The `tf-importgen` command generates Terraform import blocks for Pipeline Cloud Rules based on the `pipeline_cloud_rule_entity_id` values from updated drop rule resources.

#### Command Syntax
```bash
newrelic migrate nrqldroprules tf-importgen [flags]
```

#### Flags

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--workspacePath` | string | optional | Path to Terraform workspace (defaults to current directory) |
| `--pipelineCloudRuleIDs` | string slice | optional | List of Pipeline Cloud Rule IDs to generate import configuration with |
| `--fileName` | string | optional | File name to write import blocks to (prints to terminal if not specified) |
| `--skipResponseToPrompt` | boolean | optional | Skip all user prompts (answers 'N' to all prompts) |
| `--tofu` | boolean | optional | Use OpenTofu instead of Terraform |

#### Usage Examples

```bash
# Generate import config for all updated drop rules
newrelic migrate nrqldroprules tf-importgen

# Save import config to specific file
newrelic migrate nrqldroprules tf-importgen --fileName pipeline_imports.tf

# Use specific Pipeline Cloud Rule IDs
newrelic migrate nrqldroprules tf-importgen \
  --pipelineCloudRuleIDs "123456:pcr_789,123456:pcr_abc"

# Automated execution without prompts
newrelic migrate nrqldroprules tf-importgen --skipResponseToPrompt

# Use with OpenTofu
newrelic migrate nrqldroprules tf-importgen --tofu
```

#### What tf-importgen Does

1. **State Validation**: Verifies drop rule resources contain `pipeline_cloud_rule_entity_id` values
2. **Version Check**: Ensures Terraform/OpenTofu version ≥ 1.5 for import block support
3. **Import Generation**: Creates import blocks mapping Pipeline Cloud Rules to Terraform resources
4. **File Management**: Optionally writes import configuration to specified file
5. **Plan Execution**: Runs `terraform plan -generate-config-out` to create Pipeline Cloud Rule configurations
6. **Import Execution**: Applies changes to import Pipeline Cloud Rules into Terraform state

#### Generated Import Configuration Example
```hcl
import {
  to = newrelic_pipeline_cloud_rule.log_filter
  id = newrelic_nrql_drop_rule.log_filter.pipeline_cloud_rule_entity_id
}

import {
  to = newrelic_pipeline_cloud_rule.error_filter
  id = newrelic_nrql_drop_rule.error_filter.pipeline_cloud_rule_entity_id
}
```

#### Output Example
```
📝 Generated Terraform commands for Pipeline Cloud Rule import:
1. terraform plan -generate-config-out=generated_pipeline_rules.tf
2. terraform apply

🎉 Pipeline Cloud Rule import completed successfully!
Your Pipeline Cloud Rules have been imported into Terraform state.
```

### tf-delist Command

The `tf-delist` command safely removes NRQL drop rule resources from Terraform state without destroying the actual drop rules in New Relic.

#### Command Syntax
```bash
newrelic migrate nrqldroprules tf-delist [flags]
```

#### Flags

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--workspacePath` | string | optional | Path to Terraform workspace (defaults to current directory) |
| `--resourceIdentifiers` | string slice | optional | List of specific resource identifiers to delist |
| `--skipResponseToPrompt` | boolean | optional | Skip all user prompts (answers 'N' to all prompts) |
| `--tofu` | boolean | optional | Use OpenTofu instead of Terraform |

#### Usage Examples

```bash
# Delist all drop rule resources from state
newrelic migrate nrqldroprules tf-delist

# Delist specific resources
newrelic migrate nrqldroprules tf-delist \
  --resourceIdentifiers newrelic_nrql_drop_rule.example1,newrelic_nrql_drop_rule.example2

# Automated execution without prompts
newrelic migrate nrqldroprules tf-delist --skipResponseToPrompt

# Use with OpenTofu
newrelic migrate nrqldroprules tf-delist --tofu
```

#### What tf-delist Does

1. **Safety Warning**: Displays prominent warnings about the safe nature of the operation
2. **File Cleanup**: Checks for and attempts to remove `import_config_pipeline_rules.tf` if present
3. **State Discovery**: Identifies drop rule resources in Terraform state
4. **Command Generation**: Creates `terraform state rm` commands for each resource
5. **State Removal**: Executes state removal commands to delist resources
6. **Validation**: Confirms resources have been removed from state
7. **Instructions**: Provides detailed post-migration cleanup instructions

#### Safety Features
- **No Resource Destruction**: Only removes from state, never destroys actual New Relic resources
- **Cleanup Automation**: Automatically removes temporary import configuration files
- **Clear Messaging**: Extensive warnings and confirmations about the operation's safety
- **Post-Action Guidance**: Detailed instructions for completing the migration

#### Output Example
```
⚠️⚠️⚠️⚠️⚠️⚠️⚠️⚠️⚠️⚠️⚠️⚠️⚠️⚠️⚠️⚠️⚠️⚠️⚠️⚠️
🛡️  SAFE DELISTING MODE: Resources will be REMOVED FROM STATE ONLY
📋 The actual drop rules in New Relic will remain ACTIVE and UNCHANGED
🔄 This allows you to stop managing them via Terraform safely
⚠️⚠️⚠️⚠️⚠️⚠️⚠️⚠️⚠️⚠️⚠️⚠️⚠️⚠️⚠️⚠️⚠️⚠️⚠️⚠️

✅ All 3 NRQL drop rule resources successfully delisted from state

🔔 IMPORTANT POST-DELISTING INSTRUCTIONS
====================================================================

📝 NEXT REQUIRED STEP - Comment out configurations:
   To prevent these resources from being recreated, you MUST comment out
   or remove the following resource configurations from your .tf files:

   📄 Resource: newrelic_nrql_drop_rule.log_filter
   📄 Resource: newrelic_nrql_drop_rule.error_filter
   📄 Resource: newrelic_nrql_drop_rule.debug_filter

💡 Example of what to do in your .tf files:
   # Comment out the entire resource block like this:
   # resource "newrelic_nrql_drop_rule" "my_rule" {
   #   account_id  = var.account_id
   #   name        = "Drop rule name"
   #   action      = "drop_data"
   #   nrql        = "SELECT * FROM Log"
   # }

🎉 Delisting process completed! Don't forget to comment out configurations!
```

## Complete Migration Workflow

Here's a step-by-step example of a complete migration using all three commands:

### Step 1: Update Drop Rules
```bash
# Navigate to your Terraform workspace
cd /path/to/terraform/workspace

# Update drop rules to get Pipeline Cloud Rule entity IDs
newrelic migrate nrqldroprules tf-update
```

### Step 2: Import Pipeline Cloud Rules
```bash
# Generate and execute import configuration
newrelic migrate nrqldroprules tf-importgen
```

### Step 3: Delist Legacy Drop Rules
```bash
# Remove drop rules from Terraform state
newrelic migrate nrqldroprules tf-delist
```

### Step 4: Clean Up Configuration Files
```bash
# Comment out or remove newrelic_nrql_drop_rule blocks from .tf files
# This step must be done manually as shown in tf-delist output
```

### Step 5: Verify Migration
```bash
# Verify Pipeline Cloud Rules are managed by Terraform
terraform state list | grep pipeline_cloud_rule

# Confirm no drop rules remain in state
terraform state list | grep nrql_drop_rule

# Test that a plan shows no changes
terraform plan
```

## Common Issues and Troubleshooting

### Missing pipeline_cloud_rule_entity_id

**Issue**: tf-importgen fails because drop rules don't have `pipeline_cloud_rule_entity_id`
```
Error: resource newrelic_nrql_drop_rule.example is missing pipeline_cloud_rule_entity_id. Please run tf-update first
```

**Solution**: Run `tf-update` command first to refresh the state:
```bash
newrelic migrate nrqldroprules tf-update
```

### Provider Version Incompatibility

**Issue**: tf-update fails due to old provider version
```
Error: changes to add pipeline_cloud_rule_entity_id corresponding to drop rules would not be added to the state with New Relic Terraform Provider version 3.67.0. Provider version >= 3.68.0 is required
```

**Solution**: Update your New Relic provider version in your Terraform configuration:
```hcl
terraform {
  required_providers {
    newrelic = {
      source  = "newrelic/newrelic"
      version = "~> 3.68.0"
    }
  }
}
```

### Terraform Version Issues

**Issue**: Import generation fails due to old Terraform version
```
Error: This command requires Terraform version >= 1.5 to generate import configuration
```

**Solution**: Upgrade Terraform or OpenTofu:
```bash
# Check current version
terraform version

# Upgrade using your package manager or download from terraform.io
```

### State File Issues

**Issue**: Commands can't access Terraform state
```
Error: Terraform state list failed: No such file or directory
```

**Solution**: 
1. Ensure you're in the correct Terraform workspace directory
2. Initialize Terraform if needed: `terraform init`
3. Verify state file exists and is accessible

### No Drop Rules Found

**Issue**: Commands report no drop rule resources found
```
INFO[0001] Unable to list Terraform state and no --resourceIdentifiers provided
```

**Solution**:
1. Verify drop rule resources exist in your Terraform configuration
2. Ensure Terraform state is up to date: `terraform refresh`
3. Use `--resourceIdentifiers` flag to specify resources manually:
   ```bash
   newrelic migrate nrqldroprules tf-update \
     --resourceIdentifiers newrelic_nrql_drop_rule.example1
   ```

### Environment Variable Issues

**Issue**: Missing required environment variables
```
Error: missing required environment variables: [NEW_RELIC_API_KEY NEW_RELIC_ACCOUNT_ID]
```

**Solution**: Set required environment variables:
```bash
export NEW_RELIC_API_KEY="your-api-key"
export NEW_RELIC_ACCOUNT_ID="your-account-id"
export NEW_RELIC_REGION="US"  # Optional
```

### Import Configuration File Issues

**Issue**: tf-delist warns about import configuration file deletion
```
⚠️ WARNING: Could not delete import configuration file
File: /path/to/workspace/import_config_pipeline_rules.tf
Error: permission denied
```

**Solution**: Manually delete the file after the command completes:
```bash
rm import_config_pipeline_rules.tf
```

### Resource Still in State After Delist

**Issue**: tf-delist reports some resources weren't removed
```
⚠️ 2 out of 3 resources were delisted. Some resources may still be in state.
```

**Solution**: 
1. Check the specific error messages for failed resources
2. Manually remove stubborn resources:
   ```bash
   terraform state rm newrelic_nrql_drop_rule.problematic_resource
   ```
3. Verify resource names are correct in state:
   ```bash
   terraform state list | grep nrql_drop_rule
   ```

### Post-Migration Validation Issues

**Issue**: Terraform plan shows unwanted changes after migration
```
Plan: 3 to add, 0 to change, 0 to destroy.
```

**Solution**: Ensure you've commented out all `newrelic_nrql_drop_rule` resource blocks in your `.tf` files as instructed by tf-delist output.