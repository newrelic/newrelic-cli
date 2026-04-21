# Fleet Control CLI - Command Reference

A command-line interface for managing New Relic Fleet Control entities including fleets, configurations, deployments, members, and entity queries.

Fleet Control enables centralized management of New Relic agents across your infrastructure, allowing you to:
- Group entities into fleets for organized management
- Create and version agent configurations
- Deploy configurations to fleet members
- Manage fleet membership with ring-based deployments
- Query managed and unassigned entities

## 📑 Table of Contents

- [📋 Command Hierarchy](#-command-hierarchy)
- [🔧 Prerequisites](#-prerequisites)
  - [Required Environment Variables](#required-environment-variables)
  - [Getting Your Credentials](#getting-your-credentials)
  - [Verifying Setup](#verifying-setup)
  - [Organization ID](#organization-id)
- [📋 Command Reference](#-command-reference)
  - [Fleet Management Commands](#fleet-management-commands)
    - [`fleetcontrol fleet create`](#fleetcontrol-fleet-create---create-a-new-fleet)
    - [`fleetcontrol fleet get`](#fleetcontrol-fleet-get---get-fleet-by-id)
    - [`fleetcontrol fleet search`](#fleetcontrol-fleet-search---search-for-fleets)
    - [`fleetcontrol fleet update`](#fleetcontrol-fleet-update---update-an-existing-fleet)
    - [`fleetcontrol fleet delete`](#fleetcontrol-fleet-delete---delete-one-or-more-fleets)
  - [Fleet Member Management Commands](#fleet-member-management-commands)
    - [`fleetcontrol fleet members add`](#fleetcontrol-fleet-members-add---add-members-to-fleet)
    - [`fleetcontrol fleet members remove`](#fleetcontrol-fleet-members-remove---remove-members-from-fleet)
    - [`fleetcontrol fleet members list`](#fleetcontrol-fleet-members-list---list-fleet-members)
  - [Configuration Management Commands](#configuration-management-commands)
    - [`fleetcontrol configuration create`](#fleetcontrol-configuration-create---create-a-new-configuration)
    - [`fleetcontrol configuration get`](#fleetcontrol-configuration-get---get-configuration-content)
    - [`fleetcontrol configuration delete`](#fleetcontrol-configuration-delete---delete-a-configuration)
  - [Configuration Version Commands](#configuration-version-commands)
    - [`fleetcontrol configuration versions list`](#fleetcontrol-configuration-versions-list---list-all-configuration-versions)
    - [`fleetcontrol configuration versions add`](#fleetcontrol-configuration-versions-add---add-a-new-version-to-configuration)
    - [`fleetcontrol configuration versions delete`](#fleetcontrol-configuration-versions-delete---delete-a-configuration-version)
  - [Deployment Management Commands](#deployment-management-commands)
    - [`fleetcontrol deployment create`](#fleetcontrol-deployment-create---create-a-deployment)
    - [`fleetcontrol deployment update`](#fleetcontrol-deployment-update---update-a-deployment)
    - [`fleetcontrol deployment deploy`](#fleetcontrol-deployment-deploy---trigger-deployment)
    - [`fleetcontrol deployment delete`](#fleetcontrol-deployment-delete---delete-a-deployment)
  - [Entity Query Commands](#entity-query-commands)
    - [`fleetcontrol entities get-managed`](#fleetcontrol-entities-get-managed---list-managed-entities)
    - [`fleetcontrol entities get-unassigned`](#fleetcontrol-entities-get-unassigned---list-unassigned-entities)
- [📤 Understanding Response Formats](#-understanding-response-formats)
  - [Success Response](#success-response)
  - [Failure Response](#failure-response)
  - [Commands with Raw Output](#commands-with-raw-output)
- [🔍 Working with JSON Responses](#-working-with-json-responses)
  - [Using jq for Response Parsing](#using-jq-for-response-parsing)
  - [Practical Workflow Examples](#practical-workflow-examples)
- [🎯 Validation Rules Reference](#-validation-rules-reference)
  - [Agent Types](#agent-types)
  - [Managed Entity Types](#managed-entity-types)
  - [Configuration Modes](#configuration-modes)
  - [Tags Format](#tags-format)
  - [Agent Specification Format](#agent-specification-format)
- [🐛 Troubleshooting](#-troubleshooting)
  - [Common Issues](#common-issues)
  - [Flag Syntax Examples](#flag-syntax-examples)
  - [Validation Errors](#validation-errors)
  - [Debug Mode](#debug-mode)
- [📁 Directory Structure](#-directory-structure)
- [📚 Additional Resources](#-additional-resources)
  - [New Relic Documentation](#new-relic-documentation)



## 📋 Command Hierarchy

Fleet Control commands are organized by resource type for intuitive navigation:

```
newrelic fleetcontrol
├── fleet                    # Fleet management
│   ├── create              # Create a new fleet
│   ├── get                 # Get fleet details
│   ├── search              # Search fleets
│   ├── update              # Update fleet
│   ├── delete              # Delete fleet(s)
│   └── members             # Manage fleet members
│       ├── add             # Add entities to ring
│       ├── remove          # Remove entities from ring
│       └── list            # List fleet members
│
├── configuration            # Configuration management
│   ├── create              # Create configuration
│   ├── get                 # Get configuration content
│   ├── delete              # Delete configuration
│   └── versions            # Manage configuration versions
│       ├── list            # List all versions
│       ├── add             # Add new version
│       └── delete          # Delete specific version
│
├── deployment               # Deployment management
│   ├── create              # Create deployment
│   ├── update              # Update deployment
│   ├── deploy              # Trigger deployment
│   └── delete              # Delete deployment
│
└── entities                 # Entity queries
    ├── get-managed         # List managed entities
    └── get-unassigned      # List available entities
```

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

1. **API Key**: Generate a User API Key from [New Relic One](https://one.newrelic.com) → Click your name (bottom-left) → "API Keys" → Create "User" key (not "Browser" or "License")

2. **Account ID**: Find it in the [New Relic One](https://one.newrelic.com) URL after `/accounts/` (e.g., `https://one.newrelic.com/accounts/1234567/...`) or in Account settings.

### Verifying Setup

Test your configuration with a simple command:

```bash
# List all fleets (should return empty array or existing fleets)
./bin/darwin/newrelic fleetcontrol fleet search
```

If you see authentication errors, verify your `NEW_RELIC_API_KEY` is set correctly. If the issue persists, ensure the user associated with the API key has the necessary capabilities to perform Fleet Control operations.

### Organization ID

Most commands accept an optional `--organization-id` flag. If not provided, the CLI will automatically fetch your organization ID using your API credentials. You can find your organization ID in the New Relic UI under Account Settings.

---

## 📋 Command Reference

### Fleet Management Commands

#### `fleetcontrol fleet create` - Create a New Fleet

Create a fleet to group and manage entities of the same type.

**Required Flags:**
- `--name` - Fleet name
- `--managed-entity-type` - Type of entities this fleet will manage
  - Allowed values: `HOST`, `KUBERNETESCLUSTER` (case-insensitive)
- `--operating-system` - Operating system type (**required for HOST fleets only**)
  - Allowed values: `LINUX`, `WINDOWS` (case-insensitive)
  - **Must be specified** when creating `HOST` fleets to ensure proper agent configuration
  - **Must NOT be specified** for `KUBERNETESCLUSTER` fleets (Kubernetes manages its own OS)

**Optional Flags:**
- `--description` - Fleet description
- `--product` - New Relic product associated with this fleet
- `--tags` - Tags in format `"key:value1,value2"` (can specify multiple times)
- `--organization-id` - Organization ID (auto-fetched if not provided)

**Examples:**

```bash
# Create a fleet for Linux hosts
newrelic fleetcontrol fleet create \
  --name "Production Linux Hosts" \
  --managed-entity-type "HOST" \
  --operating-system "LINUX" \
  --description "Production environment Linux host fleet" \
  --product "Infrastructure" \
  --tags "env:prod" \
  --tags "region:us-east-1"

# Create a fleet for Windows hosts
newrelic fleetcontrol fleet create \
  --name "Production Windows Hosts" \
  --managed-entity-type "HOST" \
  --operating-system "WINDOWS" \
  --description "Production environment Windows host fleet" \
  --product "Infrastructure" \
  --tags "env:prod"

# Create a Kubernetes fleet (no --operating-system flag)
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
  "id": "fleet-abc-123",
  "name": "Production Linux Hosts",
  "type": "FLEET",
  "managedEntityType": "HOST",
  "operatingSystem": {
    "type": "LINUX"
  },
  "description": "Production environment Linux host fleet",
  "createdAt": 1770977865600,
  "updatedAt": 1770977865882
}
```

---

#### `fleetcontrol fleet get` - Get Fleet by ID

Retrieve details of a specific fleet by its ID.

**Required Flags:**
- `--fleet-id` - Fleet ID

**Examples:**

```bash
# Get fleet details
newrelic fleetcontrol fleet get --fleet-id "fleet-abc-123"

# With explicit organization ID
newrelic fleetcontrol fleet get \
  --fleet-id "fleet-abc-123" \
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

#### `fleetcontrol fleet search` - Search for Fleets

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

#### `fleetcontrol fleet update` - Update an Existing Fleet

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

#### `fleetcontrol fleet delete` - Delete One or More Fleets

Delete a single fleet or multiple fleets in bulk.

**Required Flags (mutually exclusive):**
- `--fleet-id` - Delete a single fleet
- `--fleet-ids` - Delete multiple fleets (comma-separated, requires 2+ IDs)

**Validation:**
- Must provide either `--fleet-id` or `--fleet-ids`, not both
- Bulk delete (`--fleet-ids`) requires at least 2 IDs (use `--fleet-id` for single deletion)

**Examples:**

```bash
# Delete single fleet
newrelic fleetcontrol fleet delete --fleet-id "fleet-abc-123"

# Bulk delete multiple fleets
newrelic fleetcontrol fleet delete --fleet-ids "fleet-1,fleet-2,fleet-3"
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

### Fleet Member Management Commands

#### `fleetcontrol fleet members add` - Add Members to Fleet

Add entities to a fleet ring for controlled deployment rollouts.

**Required Flags:**
- `--fleet-id` - Fleet ID to add members to
- `--ring` - Ring name within the fleet
- `--entity-ids` - Comma-separated list of entity IDs to add

**Examples:**

```bash
# Add members to fleet ring
newrelic fleetcontrol fleet members add \
  --fleet-id "fleet-abc-123" \
  --ring "canary" \
  --entity-ids "entity-1,entity-2,entity-3"
```

**Response:**
```json
{
  "status": "success",
  "error": "",
  "result": {
    "fleetId": "fleet-abc-123",
    "ring": "canary",
    "addedEntityIds": ["entity-1", "entity-2", "entity-3"],
    ...
  }
}
```

---

#### `fleetcontrol fleet members remove` - Remove Members from Fleet

Remove entities from a fleet ring.

**Required Flags:**
- `--fleet-id` - Fleet ID to remove members from
- `--ring` - Ring name within the fleet
- `--entity-ids` - Comma-separated list of entity IDs to remove

**Examples:**

```bash
# Remove members from fleet ring
newrelic fleetcontrol fleet members remove \
  --fleet-id "fleet-abc-123" \
  --ring "default" \
  --entity-ids "entity-1,entity-2"
```

**Response:**
```json
{
  "status": "success",
  "error": "",
  "result": {
    "fleetId": "fleet-abc-123",
    "ring": "default",
    "removedEntityIds": ["entity-1", "entity-2"],
    ...
  }
}
```

---

#### `fleetcontrol fleet members list` - List Fleet Members

List all entities in a fleet, optionally filtered by ring. Results are paginated and include a cursor for retrieving additional pages.

**Required Flags:**
- `--fleet-id` - Fleet ID to list members from

**Optional Flags:**
- `--ring` - Filter by specific ring name
- `--next-cursor` - Cursor to continue from previous page (for pagination)
- `--show-tags` - Include entity tags in output (default: false)

**Pagination:**
The response includes a `nextCursor` field. If present, use it with `--next-cursor` to retrieve the next page of results.

**Examples:**

```bash
# List all members in a fleet (first page)
newrelic fleetcontrol fleet members list --fleet-id "fleet-abc-123"

# List members in a specific ring
newrelic fleetcontrol fleet members list \
  --fleet-id "fleet-abc-123" \
  --ring "canary"

# Include tags in the output
newrelic fleetcontrol fleet members list \
  --fleet-id "fleet-abc-123" \
  --show-tags

# Get next page using cursor from previous response
newrelic fleetcontrol fleet members list \
  --fleet-id "fleet-abc-123" \
  --next-cursor "eyJsYXN0SWQiOiIxMjM0NSJ9"
```

**Response:**
```json
{
  "items": [
    {
      "id": "entity-123",
      "name": "host-01",
      "type": "HOST",
      "scope": {...},
      "metadata": {...}
    }
  ],
  "nextCursor": "eyJsYXN0SWQiOiIxMjM0NSJ9"
}
```

**Note:** If `nextCursor` is present in the response, there are more results available. Use the cursor value with `--next-cursor` to retrieve the next page.

---

### Configuration Management Commands

#### `fleetcontrol configuration create` - Create a New Configuration

Create a versioned configuration for fleet agents.

**Required Flags:**
- `--name` - Configuration name
- `--agent-type` - Type of agent this configuration targets
  - Allowed values: `NRInfra`, `NRDOT`, `FluentBit`, `NRPrometheusAgent` (case-insensitive)
- `--managed-entity-type` - Type of entities this configuration applies to
  - Allowed values: `HOST`, `KUBERNETESCLUSTER` (case-insensitive)
- `--operating-system` - Operating system type (**required for HOST configurations only**)
  - Allowed values: `LINUX`, `WINDOWS` (case-insensitive)
- **Exactly one of:**
  - `--configuration-file-path` - Path to configuration file (JSON/YAML) - **recommended for production**
  - `--configuration-content` - Inline configuration content (JSON/YAML) - **for testing/development only**

**Optional Flags:**
- `--organization-id` - Organization ID (auto-fetched if not provided)

**Validation:**
- Mutually exclusive: Must provide either `--configuration-file-path` OR `--configuration-content`, not both
- `--operating-system` is required when `--managed-entity-type` is `HOST`
- `--operating-system` must not be specified when `--managed-entity-type` is `KUBERNETESCLUSTER`
- File path is the recommended approach for production use
- Inline content should only be used for testing, development, or emergency purposes.

**Examples:**

```bash
# Recommended: Read from file (Linux host)
newrelic fleetcontrol configuration create \
  --name "Production Host Config" \
  --agent-type "NRInfra" \
  --managed-entity-type "HOST" \
  --operating-system "LINUX" \
  --configuration-file-path ./configs/prod-host.json

# Recommended: Read from file (Kubernetes cluster)
newrelic fleetcontrol configuration create \
  --name "Production K8s Config" \
  --agent-type "NRInfra" \
  --managed-entity-type "KUBERNETESCLUSTER" \
  --configuration-file-path ./configs/prod-k8s.json

# Alternative: Inline content (testing only)
newrelic fleetcontrol configuration create \
  --name "Test Config" \
  --agent-type "NRInfra" \
  --managed-entity-type "HOST" \
  --operating-system "LINUX" \
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

#### `fleetcontrol configuration get` - Get Configuration Content

Retrieve the configuration content for a specific configuration or version.

**Required Flags:**
- `--configuration-id` - Configuration entity ID or version entity ID

**Optional Flags:**
- `--version` - Specific version number to retrieve (defaults to latest if not provided).
- `--mode` - Entity mode:
  - `ConfigEntity` (default) - Use when `--configuration-id` is a configuration entity ID
  - `ConfigVersionEntity` - Use when `--configuration-id` is a version entity ID
- `--organization-id` - Organization ID (auto-fetched if not provided)

**Note:** This command returns raw configuration content on success (no status wrapper) but uses error wrapper on failure.

**Examples:**

```bash
# Get latest version
newrelic fleetcontrol configuration get \
  --configuration-id "CONFIG-ABC-123"

# Get specific version by number
newrelic fleetcontrol configuration get \
  --configuration-id "CONFIG-ABC-123" \
  --version 2

# Get by version entity ID
newrelic fleetcontrol configuration get \
  --configuration-id "VERSION-XYZ-789" \
  --mode "ConfigVersionEntity"

# Table format
newrelic fleetcontrol configuration get \
  --configuration-id "CONFIG-ABC-123" \
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

#### `fleetcontrol configuration delete` - Delete a Configuration

Delete an entire configuration and all its versions.

**Required Flags:**
- `--configuration-id` - Configuration entity ID to delete

**Optional Flags:**
- `--organization-id` - Organization ID (auto-fetched if not provided)

**Warning:** This deletes the configuration and all associated versions permanently.

**Examples:**

```bash
# Delete configuration
newrelic fleetcontrol configuration delete \
  --configuration-id "CONFIG-ABC-123"
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

### Configuration Version Commands

#### `fleetcontrol configuration versions list` - List All Configuration Versions

Retrieve version history for a configuration.

**Required Flags:**
- `--configuration-id` - Configuration entity ID (not version ID)

**Optional Flags:**
- `--organization-id` - Organization ID (auto-fetched if not provided)

**Validation:**
- Returns error if no versions are found (invalid ID or configuration with no versions).

**Examples:**

```bash
# List all versions
newrelic fleetcontrol configuration versions list \
  --configuration-id "CONFIG-ABC-123"

# With explicit organization ID
newrelic fleetcontrol configuration versions list \
  --configuration-id "CONFIG-ABC-123" \
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
  "error": "no version details found, please check the ID of the configuration entity provided"
}
```

---

#### `fleetcontrol configuration versions add` - Add a New Version to Configuration

Add a new version to an existing configuration.

**Required Flags:**
- `--configuration-id` - Configuration entity ID
- **Exactly one of:**
  - `--configuration-file-path` - Path to configuration file (JSON/YAML) - **recommended for production**
  - `--configuration-content` - Inline configuration content (JSON/YAML) - **for testing/development only**

**Optional Flags:**
- `--organization-id` - Organization ID (auto-fetched if not provided)

**Validation:**
- Mutually exclusive: Must provide either `--configuration-file-path` OR `--configuration-content`, not both
- File path is the recommended approach for production use
- Inline content should only be used for testing, development, or emergency purposes.

**Examples:**

```bash
# Recommended: Read from file
newrelic fleetcontrol configuration versions add \
  --configuration-id "CONFIG-ABC-123" \
  --configuration-file-path ./configs/v2-config.json

# Alternative: Inline content (testing only)
newrelic fleetcontrol configuration versions add \
  --configuration-id "CONFIG-ABC-123" \
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

#### `fleetcontrol configuration versions delete` - Delete a Configuration Version

Delete a specific version of a configuration.

**Required Flags:**
- `--version-id` - Version entity ID to delete

**Optional Flags:**
- `--organization-id` - Organization ID (auto-fetched if not provided)

**Note:** The version ID is the entity ID of the specific version (found in versions list output), not the configuration ID.

**Examples:**

```bash
# Delete specific version
newrelic fleetcontrol configuration versions delete \
  --version-id "VERSION-XYZ-456"
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

#### `fleetcontrol deployment create` - Create a Deployment

Create a deployment to roll out configurations to fleet members. Supports single or multiple agents.

**Required Flags:**
- `--fleet-id` - Fleet ID to deploy to
- `--name` - Deployment name
- **Either** (new syntax - supports multiple agents):
  - `--agent` - Agent specification in format `"AgentType:Version:ConfigVersionID1,ConfigVersionID2"` (can specify multiple times for multiple agents)
- **Or** (legacy syntax - **SINGLE agent only**):
  - `--agent-type` - Agent type (e.g., NRInfra, NRDOT) - **creates ONE agent**
  - `--agent-version` - Agent version (e.g., 1.70.0, 2.0.0, or `*` for KUBERNETESCLUSTER fleets only)
  - `--configuration-version-ids` - Configuration version IDs to deploy (comma-separated values for multiple configs on the **same** agent)

**Optional Flags:**
- `--description` - Deployment description
- `--tags` - Tags in format `"key:value1,value2"` (can specify multiple times)

**Validation:**
- Must use either `--agent` OR all three legacy flags (`--agent-type`, `--agent-version`, `--configuration-version-ids`)
- Cannot mix syntaxes - using `--agent` with any legacy flag will error
- Agent version `"*"` (wildcard) is **only allowed for KUBERNETESCLUSTER fleets**
  - HOST fleets must specify an explicit version (e.g., `"1.70.0"`)
  - The CLI validates fleet type and rejects wildcards for HOST fleets with a clear error message.

**Examples:**

```bash
# New syntax: Single agent
newrelic fleetcontrol deployment create \
  --fleet-id "fleet-abc-123" \
  --name "Production Rollout v2" \
  --agent "NRInfra:1.70.0:version-1,version-2" \
  --description "Rolling out updated monitoring configuration" \
  --tags "env:prod" \
  --tags "release:v1.2.3"

# New syntax: Multiple agents (Infrastructure + .NET)
newrelic fleetcontrol deployment create \
  --fleet-id "fleet-abc-123" \
  --name "Multi-Agent Deployment" \
  --agent "NRInfra:1.70.0:version-1,version-2" \
  --agent "NRDOT:2.0.0:version-3" \
  --description "Deploying Infrastructure and .NET agent configs"

# Legacy syntax: SINGLE agent with one config (still supported for backward compatibility)
newrelic fleetcontrol deployment create \
  --fleet-id "fleet-abc-123" \
  --name "Single Config Deployment" \
  --agent-type "NRInfra" \
  --agent-version "1.70.0" \
  --configuration-version-ids "version-1" \
  --description "One agent, one configuration"

# Legacy syntax: SINGLE agent with multiple configs (comma-separated)
# This creates ONE Infrastructure agent with TWO configuration versions
newrelic fleetcontrol deployment create \
  --fleet-id "fleet-abc-123" \
  --name "Multi-Config Single Agent" \
  --agent-type "NRInfra" \
  --agent-version "1.70.0" \
  --configuration-version-ids "version-1,version-2" \
  --description "One agent, multiple configurations"

# Kubernetes fleet with wildcard version (only works for KUBERNETESCLUSTER fleets)
newrelic fleetcontrol deployment create \
  --fleet-id "k8s-fleet-456" \
  --name "K8s Wildcard Deployment" \
  --agent "NRInfra:*:version-1,version-2" \
  --description "Using wildcard version for Kubernetes"

# IMPORTANT DIFFERENCE: New vs Legacy Syntax
# Legacy: Can only create ONE agent type per deployment
# New: Can create MULTIPLE agent types per deployment

# ❌ Cannot do this with legacy syntax (would need two separate deployments):
# --agent-type "NRInfra" ... AND --agent-type "NRDOT" ...

# ✅ Can do this with new syntax (one deployment with multiple agents):
newrelic fleetcontrol deployment create \
  --fleet-id "fleet-abc-123" \
  --name "Multi-Agent Deployment" \
  --agent "NRInfra:1.70.0:version-1,version-2" \
  --agent "NRDOT:2.0.0:version-3" \
  --description "Two different agent types in one deployment"
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

#### `fleetcontrol deployment update` - Update a Deployment

Update an existing deployment's properties, including agent configurations. Note that a deployment can only be updated before it is triggered; a deployment that has already been completed cannot be updated.

**Required Flags:**
- `--deployment-id` - Deployment ID to update

**Optional Flags:**
- `--name` - New deployment name
- **Either** (new syntax - supports multiple agents):
  - `--agent` - Agent specification in format `"AgentType:Version:ConfigVersionID1,ConfigVersionID2"` (can specify multiple times for multiple agents)
- **Or** (legacy syntax - **SINGLE agent only**):
  - `--configuration-version-ids` - Configuration version IDs to update (comma-separated values for multiple configs on the **same** agent)
- `--description` - New description
- `--tags` - New tags (replaces existing tags)

**Important Notes:**
- Only provided fields will be updated. Omitted fields remain unchanged.
- Must use either `--agent` OR `--configuration-version-ids`, not both
- Using `--agent` allows you to update agent types, versions, and configuration versions
- Using `--configuration-version-ids` (legacy) only updates configuration versions
- Agent version `"*"` (wildcard) validation is applied during update but may be limited by API constraints.

**Examples:**

```bash
# New syntax: Update agents with new versions
newrelic fleetcontrol deployment update \
  --deployment-id "deployment-xyz-789" \
  --agent "NRInfra:1.71.0:version-3,version-4" \
  --agent "NRDOT:2.1.0:version-5"

# New syntax: Update single agent
newrelic fleetcontrol deployment update \
  --deployment-id "deployment-xyz-789" \
  --agent "NRInfra:1.71.0:version-3,version-4"

# Legacy syntax: Update configuration versions only
newrelic fleetcontrol deployment update \
  --deployment-id "deployment-xyz-789" \
  --configuration-version-ids "version-3,version-4"

# Update deployment name
newrelic fleetcontrol deployment update \
  --deployment-id "deployment-xyz-789" \
  --name "Production Rollout v3"

# Update multiple fields
newrelic fleetcontrol deployment update \
  --deployment-id "deployment-xyz-789" \
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

#### `fleetcontrol deployment deploy` - Trigger Deployment

Trigger deployment to roll out configurations across fleet rings.

**Required Flags:**
- `--deployment-id` - Deployment ID to trigger
- `--rings-to-deploy` - Comma-separated list of ring names to deploy to

**Examples:**

```bash
# Deploy to multiple rings
newrelic fleetcontrol deployment deploy \
  --deployment-id "deployment-xyz-789" \
  --rings-to-deploy "canary,default"

# Deploy to a single ring
newrelic fleetcontrol deployment deploy \
  --deployment-id "deployment-xyz-789" \
  --rings-to-deploy "default"
```

---

#### `fleetcontrol deployment delete` - Delete a Deployment

Delete a fleet deployment. This operation cannot be undone and will remove the deployment. The deployment must not be actively in progress.

**Required Flags:**
- `--deployment-id` - Deployment ID to delete

**Examples:**

```bash
# Delete a deployment
newrelic fleetcontrol deployment delete \
  --deployment-id "deployment-xyz-789"
```

**Response:**
```json
{
  "status": "success",
  "error": "",
  "id": "deployment-xyz-789"
}
```

---

### Entity Query Commands

These commands help you identify which entities are managed by fleets and which are available for fleet management.

#### `fleetcontrol entities get-managed` - List Managed Entities

Retrieve all entities that are currently managed by any fleet in the account.

Managed entities are identified by having both:
- `tags.nr.fleet IS NOT NULL`
- `tags.nr.supervisor IS NOT NULL`

**Optional Flags:**
- `--entity-type` - Filter by entity type (e.g., HOST, KUBERNETESCLUSTER)
- `--limit` - Maximum number of entities to return (default: 100)
- `--include-tags` - Include entity tags in output (default: false)

**Examples:**

```bash
# Get all managed entities
newrelic fleetcontrol entities get-managed

# Limit results to 50 entities
newrelic fleetcontrol entities get-managed --limit 50

# Filter by entity type
newrelic fleetcontrol entities get-managed --entity-type HOST

# Include tags in output
newrelic fleetcontrol entities get-managed --include-tags
```

---

#### `fleetcontrol entities get-unassigned` - List Unassigned Entities

Retrieve all entities that are NOT currently managed by any fleet but are available for fleet management.

Unassigned entities are identified by having:
- `tags.nr.fleet IS NULL`
- `tags.nr.supervisor IS NOT NULL`

This command helps you identify which entities can be added to a fleet.

**Optional Flags:**
- `--entity-type` - Filter by entity type (e.g., HOST, KUBERNETESCLUSTER)
- `--limit` - Maximum number of entities to return (default: 100)
- `--include-tags` - Include entity tags in output (default: false)

**Examples:**

```bash
# Get all unassigned entities
newrelic fleetcontrol entities get-unassigned

# Limit results to 50 entities
newrelic fleetcontrol entities get-unassigned --limit 50

# Filter by entity type
newrelic fleetcontrol entities get-unassigned --entity-type HOST

# Include tags in output
newrelic fleetcontrol entities get-unassigned --include-tags
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
newrelic fleetcontrol configuration create ... | jq -r '.result.entityGuid'

# Get version number
newrelic fleetcontrol configuration versions add ... | jq -r '.result.blobVersionEntity.version'

# Get ID from delete
newrelic fleetcontrol fleet delete --fleet-id abc123 | jq -r '.id'
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
CONFIG_GUID=$(newrelic fleetcontrol configuration create \
  --name "My Config" \
  --agent-type "NRInfra" \
  --managed-entity-type "HOST" \
  --configuration-file-path ./config.json | jq -r '.result.entityGuid')

# Add a new version
newrelic fleetcontrol configuration versions add \
  --configuration-id "$CONFIG_GUID" \
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

**Find entities to add to fleet:**
```bash
# Get unassigned hosts
ENTITIES=$(newrelic fleetcontrol entities get-unassigned \
  --entity-type HOST | jq -r '.[].id' | head -3 | paste -sd "," -)

# Add them to fleet
newrelic fleetcontrol fleet members add \
  --fleet-id "fleet-abc-123" \
  --ring "production" \
  --entity-ids "$ENTITIES"
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

**Example:**
```bash
--managed-entity-type "host"  # Case-insensitive, works fine
```

---

### Configuration Modes

Used in configuration get command. Values are case-insensitive.

**Allowed values:**
- `ConfigEntity` (default) - Query by configuration entity ID
- `ConfigVersionEntity` - Query by version entity ID

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

### Agent Specification Format

Used in deployment create command with the `--agent` flag. Format: `"AgentType:Version:ConfigVersionID1,ConfigVersionID2,..."`

**Format Components:**
- **AgentType** - The type of agent (e.g., NRInfra, NRDOT, FluentBit, NRPrometheusAgent)
- **Version** - The agent version to deploy (e.g., 1.70.0, 2.0.0, or `*` for KUBERNETESCLUSTER fleets only)
  - Use explicit versions like `"1.70.0"` for HOST fleets
  - Use `"*"` (wildcard) only for KUBERNETESCLUSTER fleets
  - The CLI validates fleet type before allowing wildcard versions
- **ConfigVersionIDs** - Comma-separated list of configuration version IDs (no spaces)

**Examples:**

```bash
# Single agent with one configuration version
--agent "NRInfra:1.70.0:version-abc-123"

# Single agent with multiple configuration versions
--agent "NRInfra:1.70.0:version-1,version-2,version-3"

# Multiple agents (Infrastructure and .NET)
--agent "NRInfra:1.70.0:version-1,version-2" \
--agent "NRDOT:2.0.0:version-3"

# Multiple agents (Infrastructure, .NET, and Fluent Bit)
--agent "NRInfra:1.70.0:config-infra-v1" \
--agent "NRDOT:2.0.0:config-dotnet-v1" \
--agent "FluentBit:1.9.0:config-logs-v1"

# Wildcard version for Kubernetes fleet (only valid for KUBERNETESCLUSTER type)
--agent "NRInfra:*:config-k8s-v1"
```

**Common Errors:**
```bash
# ❌ Incorrect: Spaces in format
--agent "NRInfra : 1.70.0 : version-1, version-2"

# ✅ Correct: No spaces
--agent "NRInfra:1.70.0:version-1,version-2"

# ❌ Incorrect: Missing version
--agent "NRInfra:version-1,version-2"

# ✅ Correct: All three parts present
--agent "NRInfra:1.70.0:version-1,version-2"

# ❌ Incorrect: Using wildcard "*" with HOST fleet
--agent "NRInfra:*:version-1"  # on a HOST fleet
# Error: agent version '*' (wildcard) is not supported for HOST fleets.
#        Please specify an explicit version (e.g., '1.70.0').

# ✅ Correct: Explicit version for HOST fleet
--agent "NRInfra:1.70.0:version-1"

# ✅ Correct: Wildcard for KUBERNETESCLUSTER fleet
--agent "NRInfra:*:version-1"  # on a KUBERNETESCLUSTER fleet
```
The syntax using separate flags can be preferred in the case of single-agent deployments:

```bash
# Legacy syntax (creates ONE agent with multiple configs)
--agent-type "NRInfra" \
--agent-version "1.70.0" \
--configuration-version-ids "version-1,version-2"
# Result: 1 Infrastructure agent with 2 configuration versions

# New syntax (can create MULTIPLE agents)
--agent "NRInfra:1.70.0:version-1,version-2"
# Result: 1 Infrastructure agent with 2 configuration versions

# New syntax (multiple agents - NOT possible with legacy syntax)
--agent "NRInfra:1.70.0:version-1,version-2" \
--agent "NRDOT:2.0.0:version-3"
# Result: 2 agents (Infrastructure + .NET)
```

**Important Notes:**
- The `--agent` flag and legacy flags are mutually exclusive - choose one syntax or the other, not both
- **Legacy syntax limitation**: Can only create ONE agent type per deployment. To deploy multiple agent types (e.g., Infrastructure + .NET), use the new `--agent` syntax
- `--configuration-version-ids` in legacy syntax: Comma-separated IDs all belong to the **same single agent**, not multiple agents

---

## 🐛 Troubleshooting

### Common Issues

| Problem | Solution |
|---------|----------|
| **"Authentication failed"** | Verify `NEW_RELIC_API_KEY` is set correctly. Ensure it's a User API key, not Browser or License key. |
| **"Account not found"** | Check `NEW_RELIC_ACCOUNT_ID` is correct. Find it in New Relic UI or URL. |
| **"required flag not set"** | Ensure flag syntax is correct: `--flag-name value` or `--flag-name=value` (not `flag-name=value`) |
| **"invalid value for flag"** | Check validation rules above. Values may need to match allowed values (case-insensitive) |
| **"mutually exclusive flags"** | Only one of the mutually exclusive flags should be provided (e.g., `--fleet-id` OR `--fleet-ids`, not both). For deployments, use either `--agent` or all three legacy flags, not a mix. |
| **"agent version '*' not supported for HOST fleets"** | Wildcard version (`"*"`) is only allowed for KUBERNETESCLUSTER fleets. Use an explicit version (e.g., `"1.70.0"`) for HOST fleets. |
| **"--configuration-version-ids is required"** | When using legacy deployment syntax, you must provide all three flags: `--agent-type`, `--agent-version`, AND `--configuration-version-ids`. |
| **"no version details found"** | Configuration ID is invalid or has no versions. Verify ID is correct using `search` command |
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
  "error": "invalid value 'KUBERNETES' for flag --managed-entity-type: must be one of [HOST, KUBERNETESCLUSTER]"
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
## 📁 Directory Structure

```
internal/fleetcontrol/
├── README.md                                     # This file - Command reference
├── command.go                                    # Main entry point with command hierarchy
├── command_framework.go                          # Core framework
├── command_flags_generated.go                    # Manually maintained typed flag accessors
├── command_fleet.go                              # Command registration
├── helpers.go                                    # Shared utility functions
│
├── configs/                                      # YAML configuration files
│   ├── fleet_management_create.yaml
│   ├── fleet_management_update.yaml
│   ├── fleet_management_delete.yaml
│   ├── fleet_management_get.yaml
│   ├── fleet_management_search.yaml
│   ├── fleet_members_add.yaml
│   ├── fleet_members_remove.yaml
│   ├── fleet_members_list.yaml
│   ├── fleet_configuration_create.yaml
│   ├── fleet_configuration_get.yaml
│   ├── fleet_configuration_delete.yaml
│   ├── fleet_configuration_version_list.yaml
│   ├── fleet_configuration_version_add.yaml
│   ├── fleet_configuration_version_delete.yaml
│   ├── fleet_deployment_create.yaml
│   ├── fleet_deployment_update.yaml
│   ├── fleet_deployment_deploy.yaml
│   ├── fleet_deployment_delete.yaml
│   ├── fleet_entities_get_managed.yaml
│   └── fleet_entities_get_unassigned.yaml
│
└── Handler implementation files (one per command)
```
---

## 📚 Additional Resources

### New Relic Documentation

- [Fleet Control Documentation](https://docs.newrelic.com/docs/infrastructure/fleet-management/)
- [New Relic CLI Overview](https://github.com/newrelic/newrelic-cli)

---
