# Fleet Control CLI - Command Reference

A command-line interface for managing New Relic Fleet Control entities including fleets, configurations, deployments, and members.

Fleet Control enables centralized management of New Relic agents across your infrastructure, allowing you to:
- Group entities into fleets for organized management
- Create and version agent configurations
- Deploy configurations to fleet members
- Manage fleet membership with ring-based deployments

## 📁 Directory Structure

```
internal/fleetcontrol/
├── README.md                                      # This file - Command reference
├── CONTRIBUTING.md                                # Technical guide for developers
├── TEST_VALIDATION.md                            # Validation testing guide
├── command.go                                    # Main entry point
├── command_framework.go                          # Core framework
├── command_flags_generated.go                    # Generated typed flag accessors
├── command_fleet.go                              # Command registration
├── helpers.go                                    # Shared utility functions
│
├── configs/                                      # YAML configuration files
│   ├── fleet_management_create.yaml
│   ├── fleet_management_update.yaml
│   ├── fleet_management_delete.yaml
│   ├── fleet_management_get.yaml
│   ├── fleet_management_search.yaml
│   ├── fleet_configuration_create.yaml
│   ├── fleet_configuration_get.yaml
│   ├── fleet_configuration_version_list.yaml
│   ├── fleet_configuration_version_add.yaml
│   ├── fleet_configuration_delete.yaml
│   ├── fleet_configuration_version_delete.yaml
│   ├── fleet_deployment_create.yaml
│   ├── fleet_deployment_update.yaml
│   ├── fleet_members_add.yaml
│   └── fleet_members_remove.yaml
│
└── Handler implementation files (one per command)
```

## 📋 Available Commands

### Fleet Management
- **`create`** - Create a new fleet to group and manage entities
- **`get`** - Retrieve details of a specific fleet by ID
- **`search`** - Search for fleets by name or list all fleets
- **`update`** - Update fleet properties (name, description, tags)
- **`delete`** - Delete one or more fleets (single or bulk)

### Configuration Management
- **`create-configuration`** - Create a new versioned configuration for fleet agents
- **`get-configuration`** - Retrieve configuration content by GUID and version
- **`get-versions`** - List all versions of a configuration
- **`add-version`** - Add a new version to an existing configuration
- **`delete-configuration`** - Delete an entire configuration and all versions
- **`delete-version`** - Delete a specific configuration version

### Deployment Management (⚠️ Experimental)
- **`create-deployment`** - Create a deployment to roll out configurations to fleet members
- **`update-deployment`** - Update deployment properties and configuration versions

### Member Management (⚠️ Experimental)
- **`add-members`** - Add entities to a fleet ring
- **`remove-members`** - Remove entities from a fleet ring

**Note:** Deployment and member management commands are functional but currently being tested with the right managed entities and may be unstable.

## 🔧 Prerequisites

Before using Fleet Control commands, ensure you have the following configured:

### Required Environment Variables

Set these environment variables for authentication and authorization:

```bash
# Required: Your New Relic User API Key
export NEW_RELIC_API_KEY="NRAK-YOUR-API-KEY-HERE"

# Required: Your New Relic Account ID
export NEW_RELIC_ACCOUNT_ID="your-account-id"

# Optional: Specify region (defaults to US)
export NEW_RELIC_REGION="US"  # or "EU" for European accounts
```

### Getting Your Credentials

1. **API Key**: Generate a User API Key from New Relic:
   - Go to [New Relic One](https://one.newrelic.com)
   - Click on your name in the bottom-left corner
   - Select "API Keys"
   - Create a "User" key (not "Browser" or "License")

2. **Account ID**: Find your account ID:
   - Go to [New Relic One](https://one.newrelic.com)
   - Look in the URL after `/accounts/` (e.g., `https://one.newrelic.com/accounts/1234567/...`)
   - Or find it in Account settings

### Building the CLI

```bash
# From the repository root
go build -o ./bin/darwin/newrelic ./cmd/newrelic

# Verify the build
./bin/darwin/newrelic fleetcontrol fleet --help
```

### Verifying Setup

Test your configuration with a simple command:

```bash
# List all fleets (should return empty array or existing fleets)
./bin/darwin/newrelic fleetcontrol fleet search
```

If you see authentication errors, verify your `NEW_RELIC_API_KEY` is set correctly.

### Organization ID

Most commands accept an optional `--organization-id` flag. If not provided, the CLI will automatically fetch your organization ID using your API credentials. You can find your organization ID in the New Relic UI under Account Settings.

---

## 📋 Command Reference

### Fleet Management Commands

#### create - Create a New Fleet

Create a fleet to group and manage entities of the same type.

**Required Flags:**
- `--name` - Fleet name
- `--managed-entity-type` - Type of entities this fleet will manage
  - Allowed values: `HOST`, `KUBERNETESCLUSTER`, `APM` (case-insensitive)

**Optional Flags:**
- `--description` - Fleet description
- `--product` - New Relic product associated with this fleet
- `--tags` - Tags in format `"key:value1,value2"` (can specify multiple times)
- `--organization-id` - Organization ID (auto-fetched if not provided)

**Examples:**

```bash
# Create a fleet for hosts
newrelic fleetcontrol fleet create \
  --name "Production Hosts" \
  --managed-entity-type "HOST" \
  --description "Production environment host fleet" \
  --product "Infrastructure" \
  --tags "env:prod" \
  --tags "region:us-east-1"

# Create a Kubernetes fleet
newrelic fleetcontrol fleet create \
  --name "K8s Prod Clusters" \
  --managed-entity-type "KUBERNETESCLUSTER" \
  --tags "env:prod,team:platform"
```

**Response:**
```json
{
  "status": "success",
  "error": "",
  "result": {
    "id": "fleet-abc-123",
    "name": "Production Hosts",
    "managedEntityType": "HOST",
    "description": "Production environment host fleet",
    ...
  }
}
```

---

#### get - Get Fleet by ID

Retrieve details of a specific fleet by its ID.

**Required Flags:**
- `--id` - Fleet ID

**Optional Flags:**
- `--organization-id` - Organization ID (auto-fetched if not provided)

**Examples:**

```bash
# Get fleet details
newrelic fleetcontrol fleet get --id "fleet-abc-123"

# With explicit organization ID
newrelic fleetcontrol fleet get \
  --id "fleet-abc-123" \
  --organization-id "ORG_ID"
```

**Response:**
```json
{
  "status": "success",
  "error": "",
  "result": {
    "id": "fleet-abc-123",
    "name": "Production Hosts",
    "managedEntityType": "HOST",
    "description": "Production environment host fleet",
    "createdAt": "2026-01-15T10:30:00Z",
    ...
  }
}
```

---

#### search - Search for Fleets

Search for fleets using name filters or retrieve all fleets.

**Optional Flags:**
- `--name-equals` - Exact name match (mutually exclusive with `--name-contains`)
- `--name-contains` - Partial name match (mutually exclusive with `--name-equals`)
- `--organization-id` - Organization ID (auto-fetched if not provided)

**Note:** If neither search flag is provided, returns all fleets.

**Examples:**

```bash
# Get all fleets
newrelic fleetcontrol fleet search

# Search by exact name
newrelic fleetcontrol fleet search --name-equals "Production Hosts"

# Search by name contains
newrelic fleetcontrol fleet search --name-contains "prod"

# Table format
newrelic fleetcontrol fleet search --format text
```

**Response (raw output, no wrapper):**
```json
[
  {
    "id": "fleet-abc-123",
    "name": "Production Hosts",
    "managedEntityType": "HOST",
    ...
  },
  {
    "id": "fleet-def-456",
    "name": "Production K8s",
    "managedEntityType": "KUBERNETESCLUSTER",
    ...
  }
]
```

---

#### update - Update an Existing Fleet

Update fleet properties such as name, description, or tags.

**Required Flags:**
- `--id` - Fleet ID to update

**Optional Flags:**
- `--name` - New fleet name
- `--description` - New description
- `--tags` - New tags (replaces existing tags)
- `--organization-id` - Organization ID (auto-fetched if not provided)

**Note:** Only provided fields will be updated. Omitted fields remain unchanged.

**Examples:**

```bash
# Update fleet name
newrelic fleetcontrol fleet update \
  --id "fleet-abc-123" \
  --name "Production Hosts - Updated"

# Update description and tags
newrelic fleetcontrol fleet update \
  --id "fleet-abc-123" \
  --description "New description" \
  --tags "env:prod" \
  --tags "updated:yes"

# Update only tags
newrelic fleetcontrol fleet update \
  --id "fleet-abc-123" \
  --tags "env:prod,region:us-west-2"
```

**Response:**
```json
{
  "status": "success",
  "error": "",
  "result": {
    "id": "fleet-abc-123",
    "name": "Production Hosts - Updated",
    "description": "New description",
    ...
  }
}
```

---

#### delete - Delete One or More Fleets

Delete a single fleet or multiple fleets in bulk.

**Required Flags (mutually exclusive):**
- `--id` - Delete a single fleet
- `--ids` - Delete multiple fleets (comma-separated, requires 2+ IDs)

**Optional Flags:**
- `--organization-id` - Organization ID (auto-fetched if not provided)

**Validation:**
- Must provide either `--id` or `--ids`, not both
- Bulk delete (`--ids`) requires at least 2 IDs (use `--id` for single deletion)

**Examples:**

```bash
# Delete single fleet
newrelic fleetcontrol fleet delete --id "fleet-abc-123"

# Bulk delete multiple fleets
newrelic fleetcontrol fleet delete --ids "fleet-1,fleet-2,fleet-3"
```

**Response (single delete):**
```json
{
  "status": "success",
  "error": "",
  "id": "fleet-abc-123"
}
```

**Response (bulk delete):**
```json
[
  {
    "status": "success",
    "error": "",
    "id": "fleet-1"
  },
  {
    "status": "success",
    "error": "",
    "id": "fleet-2"
  },
  {
    "status": "failed",
    "error": "failed to delete fleet: fleet not found",
    "id": "fleet-3"
  }
]
```

---

### Configuration Management Commands

#### create-configuration - Create a New Configuration

Create a versioned configuration for fleet agents.

**Required Flags:**
- `--entity-name` - Configuration name
- `--agent-type` - Type of agent this configuration targets
  - Allowed values: `NRInfra`, `NRDOT`, `FluentBit`, `NRPrometheusAgent` (case-insensitive)
- `--managed-entity-type` - Type of entities this configuration applies to
  - Allowed values: `HOST`, `KUBERNETESCLUSTER`, `APM` (case-insensitive)
- **Exactly one of:**
  - `--configuration-file-path` - Path to configuration file (JSON/YAML) - **recommended for production**
  - `--configuration-content` - Inline configuration content (JSON/YAML) - **for testing/development only**

**Optional Flags:**
- `--organization-id` - Organization ID (auto-fetched if not provided)

**Validation:**
- Mutually exclusive: Must provide either `--configuration-file-path` OR `--configuration-content`, not both
- File path is the recommended approach for production use
- Inline content should only be used for testing, development, or emergency purposes

**Examples:**

```bash
# Recommended: Read from file
newrelic fleetcontrol fleet create-configuration \
  --entity-name "Production Host Config" \
  --agent-type "NRInfra" \
  --managed-entity-type "HOST" \
  --configuration-file-path ./configs/prod-host.json

# Alternative: Inline content (testing only)
newrelic fleetcontrol fleet create-configuration \
  --entity-name "Test Config" \
  --agent-type "NRInfra" \
  --managed-entity-type "HOST" \
  --configuration-content '{"metrics_interval": 15}'
```

**Response:**
```json
{
  "status": "success",
  "error": "",
  "result": {
    "entityGuid": "CONFIG-ABC-123",
    "entityName": "Production Host Config",
    "blobVersionEntity": {
      "version": 1,
      "guid": "VERSION-XYZ-789",
      "blobId": "blob-456"
    },
    ...
  }
}
```

---

#### get-configuration - Get Configuration Content

Retrieve the configuration content for a specific configuration or version.

**Required Flags:**
- `--entity-guid` - Configuration entity GUID or version entity GUID

**Optional Flags:**
- `--version` - Specific version number to retrieve (defaults to latest if not provided)
- `--mode` - Entity mode:
  - `ConfigEntity` (default) - Use when `--entity-guid` is a configuration entity GUID
  - `ConfigVersionEntity` - Use when `--entity-guid` is a version entity GUID
- `--organization-id` - Organization ID (auto-fetched if not provided)

**Note:** This command returns raw configuration content on success (no status wrapper) but uses error wrapper on failure.

**Examples:**

```bash
# Get latest version
newrelic fleetcontrol fleet get-configuration \
  --entity-guid "CONFIG-ABC-123"

# Get specific version by number
newrelic fleetcontrol fleet get-configuration \
  --entity-guid "CONFIG-ABC-123" \
  --version 2

# Get by version entity GUID
newrelic fleetcontrol fleet get-configuration \
  --entity-guid "VERSION-XYZ-789" \
  --mode "ConfigVersionEntity"

# Table format
newrelic fleetcontrol fleet get-configuration \
  --entity-guid "CONFIG-ABC-123" \
  --format text
```

**Response (raw output on success):**
```json
{
  "metrics_interval": 15,
  "log_level": "info",
  "custom_attributes": {
    "environment": "production"
  }
}
```

**Response (on failure):**
```json
{
  "status": "failed",
  "error": "failed to get configuration: configuration not found"
}
```

---

#### get-versions - List All Configuration Versions

Retrieve version history for a configuration.

**Required Flags:**
- `--configuration-guid` - Configuration entity GUID (not version GUID)

**Optional Flags:**
- `--organization-id` - Organization ID (auto-fetched if not provided)

**Validation:**
- Returns error if no versions are found (invalid GUID or configuration with no versions)

**Examples:**

```bash
# List all versions
newrelic fleetcontrol fleet get-versions \
  --configuration-guid "CONFIG-ABC-123"

# With explicit organization ID
newrelic fleetcontrol fleet get-versions \
  --configuration-guid "CONFIG-ABC-123" \
  --organization-id "ORG_ID"
```

**Response:**
```json
{
  "status": "success",
  "error": "",
  "result": {
    "versions": [
      {
        "version": 3,
        "guid": "VERSION-XYZ-789",
        "blobId": "blob-789",
        "createdAt": "2026-01-20T14:30:00Z"
      },
      {
        "version": 2,
        "guid": "VERSION-XYZ-456",
        "blobId": "blob-456",
        "createdAt": "2026-01-15T10:00:00Z"
      },
      {
        "version": 1,
        "guid": "VERSION-XYZ-123",
        "blobId": "blob-123",
        "createdAt": "2026-01-10T09:00:00Z"
      }
    ]
  }
}
```

**Error (no versions found):**
```json
{
  "status": "failed",
  "error": "no version details found, please check the GUID of the configuration entity provided"
}
```

---

#### add-version - Add a New Version to Configuration

Add a new version to an existing configuration.

**Required Flags:**
- `--configuration-guid` - Configuration entity GUID
- **Exactly one of:**
  - `--configuration-file-path` - Path to configuration file (JSON/YAML) - **recommended for production**
  - `--configuration-content` - Inline configuration content (JSON/YAML) - **for testing/development only**

**Optional Flags:**
- `--organization-id` - Organization ID (auto-fetched if not provided)

**Validation:**
- Mutually exclusive: Must provide either `--configuration-file-path` OR `--configuration-content`, not both
- File path is the recommended approach for production use
- Inline content should only be used for testing, development, or emergency purposes

**Examples:**

```bash
# Recommended: Read from file
newrelic fleetcontrol fleet add-version \
  --configuration-guid "CONFIG-ABC-123" \
  --configuration-file-path ./configs/v2-config.json

# Alternative: Inline content (testing only)
newrelic fleetcontrol fleet add-version \
  --configuration-guid "CONFIG-ABC-123" \
  --configuration-content '{"metrics_interval": 30, "updated": true}'
```

**Response:**
```json
{
  "status": "success",
  "error": "",
  "result": {
    "blobVersionEntity": {
      "version": 2,
      "guid": "VERSION-XYZ-456",
      "blobId": "blob-456"
    },
    "entityGuid": "CONFIG-ABC-123",
    ...
  }
}
```

---

#### delete-configuration - Delete a Configuration

Delete an entire configuration and all its versions.

**Required Flags:**
- `--configuration-guid` - Configuration entity GUID to delete

**Optional Flags:**
- `--organization-id` - Organization ID (auto-fetched if not provided)

**Warning:** This deletes the configuration and all associated versions permanently.

**Examples:**

```bash
# Delete configuration
newrelic fleetcontrol fleet delete-configuration \
  --configuration-guid "CONFIG-ABC-123"
```

**Response:**
```json
{
  "status": "success",
  "error": "",
  "id": "CONFIG-ABC-123"
}
```

---

#### delete-version - Delete a Configuration Version

Delete a specific version of a configuration.

**Required Flags:**
- `--version-guid` - Version entity GUID to delete

**Optional Flags:**
- `--organization-id` - Organization ID (auto-fetched if not provided)

**Note:** The version GUID is the entity GUID of the specific version (found in get-versions output), not the configuration GUID.

**Examples:**

```bash
# Delete specific version
newrelic fleetcontrol fleet delete-version \
  --version-guid "VERSION-XYZ-456"
```

**Response:**
```json
{
  "status": "success",
  "error": "",
  "id": "VERSION-XYZ-456"
}
```

---

### Deployment Management Commands

**⚠️ Experimental:** These commands are functional but currently being tested with the right managed entities. They may be unstable in their current state.

#### create-deployment - Create a Deployment

Create a deployment to roll out configurations to fleet members.

**Required Flags:**
- `--fleet-id` - Fleet ID to deploy to
- `--name` - Deployment name
- `--configuration-version-ids` - Comma-separated list of configuration version IDs to deploy

**Optional Flags:**
- `--description` - Deployment description
- `--tags` - Tags in format `"key:value1,value2"` (can specify multiple times)
- `--organization-id` - Organization ID (auto-fetched if not provided)

**Examples:**

```bash
# Create deployment
newrelic fleetcontrol fleet create-deployment \
  --fleet-id "fleet-abc-123" \
  --name "Production Rollout v2" \
  --configuration-version-ids "version-1,version-2" \
  --description "Rolling out updated monitoring configuration" \
  --tags "env:prod" \
  --tags "release:v1.2.3"
```

**Response:**
```json
{
  "status": "success",
  "error": "",
  "result": {
    "id": "deployment-xyz-789",
    "name": "Production Rollout v2",
    "fleetId": "fleet-abc-123",
    "configurationVersionIds": ["version-1", "version-2"],
    ...
  }
}
```

---

#### update-deployment - Update a Deployment

Update an existing deployment's properties.

**Required Flags:**
- `--id` - Deployment ID to update

**Optional Flags:**
- `--name` - New deployment name
- `--configuration-version-ids` - New comma-separated list of configuration version IDs
- `--description` - New description
- `--tags` - New tags (replaces existing tags)
- `--organization-id` - Organization ID (auto-fetched if not provided)

**Note:** Only provided fields will be updated. Omitted fields remain unchanged.

**Examples:**

```bash
# Update deployment name
newrelic fleetcontrol fleet update-deployment \
  --id "deployment-xyz-789" \
  --name "Production Rollout v3"

# Update configuration versions
newrelic fleetcontrol fleet update-deployment \
  --id "deployment-xyz-789" \
  --configuration-version-ids "version-3,version-4"

# Update multiple fields
newrelic fleetcontrol fleet update-deployment \
  --id "deployment-xyz-789" \
  --name "Updated Deployment" \
  --description "New description" \
  --tags "env:prod,updated:yes"
```

**Response:**
```json
{
  "status": "success",
  "error": "",
  "result": {
    "id": "deployment-xyz-789",
    "name": "Production Rollout v3",
    "configurationVersionIds": ["version-3", "version-4"],
    ...
  }
}
```

---

### Member Management Commands

**⚠️ Experimental:** These commands are functional but currently being tested with the right managed entities. They may be unstable in their current state.

#### add-members - Add Members to Fleet

Add entities to a fleet ring.

**Required Flags:**
- `--fleet-id` - Fleet ID to add members to
- `--ring` - Ring name within the fleet
- `--entity-ids` - Comma-separated list of entity IDs to add

**Optional Flags:**
- `--organization-id` - Organization ID (auto-fetched if not provided)

**Examples:**

```bash
# Add members to fleet ring
newrelic fleetcontrol fleet add-members \
  --fleet-id "fleet-abc-123" \
  --ring "production" \
  --entity-ids "entity-1,entity-2,entity-3"
```

**Response:**
```json
{
  "status": "success",
  "error": "",
  "result": {
    "fleetId": "fleet-abc-123",
    "ring": "production",
    "addedEntityIds": ["entity-1", "entity-2", "entity-3"],
    ...
  }
}
```

---

#### remove-members - Remove Members from Fleet

Remove entities from a fleet ring.

**Required Flags:**
- `--fleet-id` - Fleet ID to remove members from
- `--ring` - Ring name within the fleet
- `--entity-ids` - Comma-separated list of entity IDs to remove

**Optional Flags:**
- `--organization-id` - Organization ID (auto-fetched if not provided)

**Examples:**

```bash
# Remove members from fleet ring
newrelic fleetcontrol fleet remove-members \
  --fleet-id "fleet-abc-123" \
  --ring "production" \
  --entity-ids "entity-1,entity-2"
```

**Response:**
```json
{
  "status": "success",
  "error": "",
  "result": {
    "fleetId": "fleet-abc-123",
    "ring": "production",
    "removedEntityIds": ["entity-1", "entity-2"],
    ...
  }
}
```

---

## 📤 Understanding Response Formats

All commands return consistent JSON output for easy parsing and automation.

### Success Response

Most commands wrap results with status and error fields:

```json
{
  "status": "success",
  "error": "",
  "result": {
    "entityGuid": "ABC123DEF456",
    "name": "Production Fleet",
    ...
  }
}
```

Delete operations:

```json
{
  "status": "success",
  "error": "",
  "id": "deleted-entity-id"
}
```

Bulk operations return an array:

```json
[
  {
    "status": "success",
    "error": "",
    "id": "fleet-1"
  },
  {
    "status": "failed",
    "error": "failed to delete fleet: not found",
    "id": "fleet-2"
  }
]
```

### Failure Response

```json
{
  "status": "failed",
  "error": "failed to create fleet: organization ID required"
}
```

### Commands with Raw Output

Search and get-configuration return raw data (no wrapper) for table formatting:

```json
[
  {
    "id": "fleet-123",
    "name": "Production Fleet",
    ...
  }
]
```

Errors from these commands still use the status/error wrapper.

---

## 🔍 Working with JSON Responses

### Using jq for Response Parsing

**Extract data from success:**
```bash
# Get entityGuid from create
newrelic fleetcontrol fleet create-configuration ... | jq -r '.result.entityGuid'

# Get version number
newrelic fleetcontrol fleet add-version ... | jq -r '.result.blobVersionEntity.version'

# Get ID from delete
newrelic fleetcontrol fleet delete --id abc123 | jq -r '.id'
```

**Check status before extracting:**
```bash
# Extract only if successful
newrelic fleetcontrol fleet create ... | \
  jq -r 'select(.status == "success") | .result.entityGuid'

# Show error on failure, result on success
newrelic fleetcontrol fleet create ... | \
  jq -r 'if .status == "success" then .result.entityGuid else .error end'
```

**Store result with error handling:**
```bash
OUTPUT=$(newrelic fleetcontrol fleet create ...)
STATUS=$(echo "$OUTPUT" | jq -r '.status')

if [ "$STATUS" = "success" ]; then
  GUID=$(echo "$OUTPUT" | jq -r '.result.entityGuid')
  echo "Created: $GUID"
else
  ERROR=$(echo "$OUTPUT" | jq -r '.error')
  echo "Failed: $ERROR"
  exit 1
fi
```

### Practical Workflow Examples

**Create fleet and store ID:**
```bash
# Create fleet and extract ID
FLEET_ID=$(newrelic fleetcontrol fleet create \
  --name "My Fleet" \
  --managed-entity-type "HOST" | jq -r '.result.id')

echo "Created fleet: $FLEET_ID"
```

**Create configuration and add version:**
```bash
# Create configuration
CONFIG_GUID=$(newrelic fleetcontrol fleet create-configuration \
  --entity-name "My Config" \
  --agent-type "NRInfra" \
  --managed-entity-type "HOST" \
  --configuration-file-path ./config.json | jq -r '.result.entityGuid')

# Add a new version
newrelic fleetcontrol fleet add-version \
  --configuration-guid "$CONFIG_GUID" \
  --configuration-file-path ./config-v2.json
```

**List and filter fleets:**
```bash
# Get all production fleets
newrelic fleetcontrol fleet search | jq '.[] | select(.name | contains("prod"))'

# Count total fleets
newrelic fleetcontrol fleet search | jq 'length'

# Get fleet names only
newrelic fleetcontrol fleet search | jq -r '.[].name'
```

---

## 🎯 Validation Rules Reference

### Agent Types

Used in configuration commands. Values are case-insensitive.

**Allowed values:**
- `NRInfra` - New Relic Infrastructure agent
- `NRDOT` - New Relic .NET agent
- `FluentBit` - Fluent Bit log forwarder
- `NRPrometheusAgent` - New Relic Prometheus agent

**Example:**
```bash
--agent-type "nrinfra"  # Case-insensitive, works fine
```

---

### Managed Entity Types

Used in fleet and configuration commands. Values are case-insensitive.

**Allowed values:**
- `HOST` - Physical or virtual hosts
- `KUBERNETESCLUSTER` - Kubernetes clusters
- `APM` - APM applications

**Example:**
```bash
--managed-entity-type "host"  # Case-insensitive, works fine
```

---

### Configuration Modes

Used in get-configuration command. Values are case-insensitive.

**Allowed values:**
- `ConfigEntity` (default) - Query by configuration entity GUID
- `ConfigVersionEntity` - Query by version entity GUID

**Example:**
```bash
--mode "configversionentity"  # Case-insensitive, works fine
```

---

### Tags Format

Tags must be in format `"key:value1,value2"`.

**Examples:**
```bash
# Single tag with single value
--tags "env:prod"

# Single tag with multiple values
--tags "env:prod,staging"

# Multiple tags
--tags "env:prod" --tags "team:platform" --tags "region:us-east-1"

# Result parsed as:
# [
#   {Key: "env", Values: ["prod"]},
#   {Key: "team", Values: ["platform"]},
#   {Key: "region", Values: ["us-east-1"]}
# ]
```

---

### Mutually Exclusive Flags

Some commands enforce mutual exclusivity between certain flags.

**Delete fleet:**
- `--id` (single delete) OR `--ids` (bulk delete)
- Cannot use both
- Bulk delete requires 2+ IDs

**Search fleet:**
- `--name-equals` OR `--name-contains`
- Both optional, but mutually exclusive
- If neither provided, returns all fleets

**Configuration content:**
- `--configuration-file-path` OR `--configuration-content`
- Must provide exactly one, not both or neither
- File path is recommended for production
- Inline content is for testing/development only

---

## 🐛 Troubleshooting

### Common Issues

| Problem | Solution |
|---------|----------|
| **"Authentication failed"** | Verify `NEW_RELIC_API_KEY` is set correctly. Ensure it's a User API key, not Browser or License key. |
| **"Account not found"** | Check `NEW_RELIC_ACCOUNT_ID` is correct. Find it in New Relic UI or URL. |
| **"required flag not set"** | Ensure flag syntax is correct: `--flag-name value` or `--flag-name=value` (not `flag-name=value`) |
| **"invalid value for flag"** | Check validation rules above. Values may need to match allowed values (case-insensitive) |
| **"mutually exclusive flags"** | Only one of the mutually exclusive flags should be provided (e.g., `--id` OR `--ids`, not both) |
| **"no version details found"** | Configuration GUID is invalid or has no versions. Verify GUID is correct using `search` command |
| **File not found error** | When using `--configuration-file-path`, ensure the file path is correct and file exists |
| **"organization ID required"** | Provide `--organization-id` explicitly if auto-fetch fails |
| **Empty response** | Check that entity exists using `search` or `get` commands |
| **JSON parse errors** | Ensure JSON/YAML configuration content is valid. Test with `jq` or YAML linter |

### Flag Syntax Examples

**Correct:**
```bash
--flag-name value
--flag-name=value
--flag-name "value with spaces"
```

**Incorrect:**
```bash
flag-name=value        # Missing -- prefix
-flag-name value       # Single dash instead of double
```

### Validation Errors

When you see "invalid value for flag", check:
1. Value is in the allowed values list (see Validation Rules Reference)
2. Spelling is correct (validation is case-insensitive but value must be in the list)
3. No extra spaces or quotes

**Example validation error:**
```json
{
  "status": "failed",
  "error": "invalid value 'KUBERNETES' for flag --managed-entity-type: must be one of [HOST, KUBERNETESCLUSTER, APM]"
}
```

**Solution:** Use `KUBERNETESCLUSTER` instead of `KUBERNETES`.

### Debug Mode

Enable verbose logging for troubleshooting:

```bash
# Set log level to debug
export NEW_RELIC_CLI_LOG_LEVEL=debug

# Run command
newrelic fleetcontrol fleet create --name "Test"
```

---

## 📚 Additional Resources

### For Developers

If you're looking to contribute or understand the technical architecture:

- See **[CONTRIBUTING.md](./CONTRIBUTING.md)** for technical details on:
  - YAML-driven framework architecture
  - Adding new commands and flags
  - Code organization principles
  - Testing and validation

- See **[TEST_VALIDATION.md](./TEST_VALIDATION.md)** for:
  - Testing YAML validation
  - Understanding validation flow
  - Common development issues

### New Relic Documentation

- [Fleet Control Documentation](https://docs.newrelic.com/docs/infrastructure/fleet-management/)
- [New Relic API Keys](https://docs.newrelic.com/docs/apis/intro-apis/new-relic-api-keys/)
- [New Relic CLI Overview](https://github.com/newrelic/newrelic-cli)

---

## 📝 Recent Updates

**January 2026:**
- Split hybrid file flags into mutually exclusive `--configuration-file-path` and `--configuration-content`
- Renamed `--entity-guid` to `--configuration-guid` in get-versions for clarity
- Added consistent status/error response wrappers for all commands
- Added empty results validation for get-versions command
- Fixed Go Client delete operations to handle empty responses
- Improved documentation structure for first-time users
- Added comprehensive prerequisites and setup guide

**Last Updated**: January 27, 2026
