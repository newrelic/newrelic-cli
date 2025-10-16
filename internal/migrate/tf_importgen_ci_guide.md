# New Relic CLI: NRQL Drop Rules to Pipeline Cloud Rules Migration Guide (CI/CD)

## Overview

The `tf-importgen-ci` command is a specialized tool designed to facilitate the migration from NRQL Drop Rules to Pipeline Cloud Rules in CI/CD environments. This command generates complete Terraform/OpenTofu configurations, including import blocks and provider configurations, to enable seamless migration with minimal manual intervention.

## Prerequisites

### Technical Requirements

- **Terraform/OpenTofu**: Version ≥ 1.5 (required for import blocks)
- **New Relic CLI**: Latest version
- **Environment Variables**: 
  - `NEW_RELIC_API_KEY` (required)
  - `NEW_RELIC_ACCOUNT_ID` (required)
  - `NEW_RELIC_REGION` (optional, defaults to 'US')

### Data Requirements

The command requires input data containing drop rule resource information in JSON format. This data typically comes from the output of the `tf-update` command or manual extraction from existing Terraform state.

**Required JSON Structure:**
```json
{
  "drop_rule_resource_ids": [
    {
      "name": "resource_name",
      "id": "account_id:rule_id",
      "pipeline_cloud_rule_entity_id": "account_id:pipeline_rule_id"
    }
  ]
}
```

**Note**: If duplicate resource names are found in the input data, the command will automatically resolve them by adding random alphabetic suffixes to ensure unique resource definitions.

## Command Syntax

```bash
newrelic migrate nrqldroprules tf-importgen-ci [flags]
```

### Flags

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--file` | string | conditional | Path to JSON file containing drop rule resource IDs |
| `--json` | string | conditional | JSON string containing drop rule resource IDs |
| `--workspacePath` | string | optional | Path to Terraform workspace (defaults to current directory) |
| `--tofu` | boolean | optional | Use OpenTofu instead of Terraform |

**Note**: Either `--file` or `--json` must be provided, but not both.

## Usage Examples

### Example 1: Basic Usage with File Input

```bash
# Generate import configuration from JSON file
newrelic migrate nrqldroprules tf-importgen-ci --file drop_rules.json
```

### Example 2: OpenTofu Usage

```bash
# Use OpenTofu instead of Terraform
newrelic migrate nrqldroprules tf-importgen-ci \
  --file drop_rules.json \
  --tofu
```

### Example 3: Custom Workspace Path

```bash
# Specify custom workspace directory
newrelic migrate nrqldroprules tf-importgen-ci \
  --file drop_rules.json \
  --workspacePath /path/to/terraform/workspace
```

### Example 4: Inline JSON Input

```bash
# Provide JSON data directly
newrelic migrate nrqldroprules tf-importgen-ci \
  --json '{"drop_rule_resource_ids":[{"name":"example_rule","id":"123456:rule123","pipeline_cloud_rule_entity_id":"123456:pcr456"}]}'
```

## Migration Workflow

### Phase 1: Pre-Migration Preparation

1. **Identify NRQL Drop Rules**: Catalog all existing NRQL drop rules in your infrastructure
2. **Run tf-update Command**: Execute the `tf-update` command to populate `pipeline_cloud_rule_entity_id` values
3. **Extract Resource Data**: Collect the updated resource information in the required JSON format
4. **Environment Setup**: Ensure all required environment variables are configured

### Phase 2: Local Workspace Generation

1. **Execute tf-importgen-ci Command**:
   ```bash
   newrelic migrate nrqldroprules tf-importgen-ci --file drop_rules.json
   ```

2. **Command Execution Flow**:
   - Validates input parameters and environment variables
   - Checks account ID consistency between environment and input data
   - Validates Terraform/OpenTofu installation and version
   - Creates or validates workspace directory
   - Generates provider configuration (`provider.tf`)
   - Generates import blocks (`imports.tf`)
   - Initializes Terraform/OpenTofu workspace
   - Runs plan to generate Pipeline Cloud Rules configuration (`pcrs.tf`)
   - Formats all configuration files

3. **Generated Files**:
   - `provider.tf`: New Relic provider configuration
   - `imports.tf`: Import blocks for Pipeline Cloud Rules
   - `pcrs.tf`: Generated Pipeline Cloud Rules resource configurations

### Phase 3: CI/CD Integration

The migration process differs depending on your Terraform/OpenTofu setup:

#### Standard Migration Process (Recommended)

1. **Copy Generated Files**: Transfer `provider.tf`, `imports.tf`, and `pcrs.tf` to your CI workspace
2. **Execute Import**: Run `terraform plan` and `terraform apply` to import Pipeline Cloud Rules
3. **Clean Up State**: After successful import, remove old NRQL drop rules from state:
   ```bash
   terraform state rm newrelic_nrql_drop_rule.rule1 newrelic_nrql_drop_rule.rule2
   ```

## Generated File Details

### provider.tf
Contains the New Relic provider configuration with version constraints:
```hcl
terraform {
  required_providers {
    newrelic = {
      source  = "newrelic/newrelic"
      version = "~> 3.0"
    }
  }
}

provider "newrelic" {
  # Configuration from environment variables
}
```

### imports.tf
Contains import blocks for each Pipeline Cloud Rule:
```hcl
import {
  to = newrelic_pipeline_cloud_rule.rule_name
  id = "account_id:pipeline_rule_id"
}
```

### pcrs.tf
Auto-generated Pipeline Cloud Rules resource configurations (generated by Terraform plan).

## Account ID Validation

The command performs automatic account ID consistency validation:

- Extracts account IDs from Pipeline Cloud Rule entity IDs
- Compares with `NEW_RELIC_ACCOUNT_ID` environment variable
- Displays warnings for any mismatches
- Highlights potential import failures due to account inconsistencies

## CI/CD Deployment Steps

**Requirements**: Terraform/OpenTofu ≥ 1.5

### Standard Process:

1. **Copy Generated Files**: Transfer `provider.tf`, `imports.tf`, and `pcrs.tf` to your CI workspace
2. **Execute Import**: Run `terraform plan` and `terraform apply` to import Pipeline Cloud Rules
3. **Remove Old Resources**: After successful import, clean up old NRQL drop rules from state:
   ```bash
   # Identify resources in state
   terraform state list | grep nrql_drop_rule
   
   # Remove them from state (use actual resource names from your state)
   terraform state rm newrelic_nrql_drop_rule.rule1 newrelic_nrql_drop_rule.rule2
   ```
4. **Verify Migration**: Confirm Pipeline Cloud Rules are working and old rules are removed from state

## Common Issues and Troubleshooting

### Duplicate Resource Names

**Issue**: Duplicate resource names in input data
```
⚠️ WARNING: DUPLICATE RESOURCE NAMES DETECTED
The following resource names appear multiple times in the input data:
  - log_filter_rule (appears 3 times at positions: [0, 2, 5])
```

**Solution**: The command automatically resolves this by renaming duplicates. No manual intervention required.

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

### Version Compatibility Issues

**Issue**: Terraform/OpenTofu version too old
```
Error: this command requires Terraform version >= 1.5 to generate import configuration
```

**Solution**: Upgrade Terraform/OpenTofu:
```bash
# For Terraform
terraform version
# Upgrade using your package manager or download from terraform.io

# For OpenTofu
tofu version
# Upgrade using your package manager or download from opentofu.org
```

### Account ID Mismatch

**Issue**: Account ID inconsistency warning
```
⚠️ WARNING: ACCOUNT ID MISMATCH DETECTED
Environment NEW_RELIC_ACCOUNT_ID: 123456
The following resources have different account IDs:
  - rule_name (rule account: 789012)
```

**Solution**: 
1. Verify correct account ID in environment variable
2. Ensure input data contains resources from the correct account
3. Update data source or environment variable as needed

### Workspace Conflicts

**Issue**: Conflicting files in workspace
```
Error: workspace directory contains conflicting files that may interfere with the import process
```

**Solution**: 
1. Use an empty directory for the workspace
2. Clean up existing Terraform files from the directory
3. Specify a different workspace path using `--workspacePath`

### Input Data Format Issues

**Issue**: JSON parsing errors
```
Error: failed to parse JSON: invalid character '}' looking for beginning of object key string
```

**Solution**:
1. Validate JSON format using tools like `jq`:
   ```bash
   cat drop_rules.json | jq .
   ```
2. Ensure proper JSON structure with required fields
3. Check for trailing commas or syntax errors

### Import Failures in CI

**Issue**: Import blocks fail during CI execution
```
Error: resource not found during import
```

**Solution**:
1. Verify Pipeline Cloud Rules exist in target account
2. Check account ID consistency
3. Ensure API key has sufficient permissions
4. Validate entity IDs are correct and accessible

