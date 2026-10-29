# Access Management CLI - Command Reference

A command-line interface for managing New Relic access grants, roles, and permissions. These commands correspond to the **Access Management** section of the New Relic Administration UI.

Access management covers what groups can access — the roles they hold and the accounts or organization they hold them in. To manage the users and groups that receive access, see the [User Management CLI](../usermanagement/README.md).

## Table of Contents

- [Command Hierarchy](#command-hierarchy)
- [Prerequisites](#prerequisites)
- [Command Reference](#command-reference)
  - [Grant Commands](#grant-commands)
  - [Role Commands](#role-commands)
  - [Permission Commands](#permission-commands)
- [Working with JSON Responses](#working-with-json-responses)
- [End-to-End Workflow](#end-to-end-workflow)
- [Directory Structure](#directory-structure)
- [Additional Resources](#additional-resources)

---

## Command Hierarchy

```
newrelic accessmanagement
├── grants
│   ├── get     # Retrieve access grants, optionally filtered by group
│   ├── create  # Grant a role to a group (account or organization scope)
│   └── revoke  # Revoke a role from a group
│
├── roles
│   └── get     # Retrieve available roles, optionally filtered by name or group
│
└── permissions
    └── get     # Retrieve permissions, optionally filtered by role or scope
```

---

## Prerequisites

### Required Environment Variables

```bash
export NEW_RELIC_API_KEY="NRAK-YOUR-API-KEY-HERE"
export NEW_RELIC_ACCOUNT_ID="your-account-id"
```

### Required Permissions

Your API key must belong to a user with the **Organization manager** role. Read operations (grants get, roles get, permissions get) require at minimum the **Authentication domain read-only** capability.

### Key IDs You Will Need

- **Group ID** — from `newrelic usermanagement groups get`
- **Role ID** — from `newrelic accessmanagement roles get`
- **Account ID** — your New Relic account number (required for account-scoped grants)

---

## Command Reference

### Grant Commands

#### `accessmanagement grants get` - Retrieve access grants

Returns all access grants in your organization, with an optional filter by group.

**Optional Flags:**
- `--groupId` - Filter grants by group ID

**Examples:**

```bash
# Get all grants
newrelic accessmanagement grants get

# Get grants for a specific group
newrelic accessmanagement grants get --groupId <groupId>
```

**Response:**
```json
{
  "grants": [
    {
      "groupId": "group-123",
      "roleId": "role-456",
      "roleType": "ORGANIZATION",
      "accountId": null,
      "displayName": "All product admin"
    },
    {
      "groupId": "group-123",
      "roleId": "role-789",
      "roleType": "ACCOUNT",
      "accountId": 12345678,
      "displayName": "Standard user"
    }
  ]
}
```

---

#### `accessmanagement grants create` - Create an access grant

Grants a role to a group. You must specify whether the grant applies to a specific account or to the entire organization.

**Required Flags:**
- `--groupId` - Group ID
- `--roleId` - Role ID
- `--scope` - Grant scope: `account` or `organization`

**Optional Flags:**
- `--accountId` - Account ID (required when `--scope` is `account`)

**Examples:**

```bash
# Grant account-scoped access
newrelic accessmanagement grants create \
  --groupId <groupId> \
  --roleId <roleId> \
  --scope account \
  --accountId 12345678

# Grant organization-scoped access
newrelic accessmanagement grants create \
  --groupId <groupId> \
  --roleId <roleId> \
  --scope organization
```

**Response:**
```json
{
  "roles": [
    {
      "displayName": "Standard user",
      "organizationId": "org-abc123",
      "roleId": "role-789",
      "roleType": "ACCOUNT"
    }
  ]
}
```

---

#### `accessmanagement grants revoke` - Revoke an access grant

Removes a role from a group. The group's users immediately lose any access that role provided.

**Required Flags:**
- `--groupId` - Group ID
- `--roleId` - Role ID
- `--scope` - Grant scope: `account` or `organization`

**Optional Flags:**
- `--accountId` - Account ID (required when `--scope` is `account`)

**Examples:**

```bash
# Revoke account-scoped access
newrelic accessmanagement grants revoke \
  --groupId <groupId> \
  --roleId <roleId> \
  --scope account \
  --accountId 12345678

# Revoke organization-scoped access
newrelic accessmanagement grants revoke \
  --groupId <groupId> \
  --roleId <roleId> \
  --scope organization
```

**Response:** Prints `success` on completion.

---

### Role Commands

#### `accessmanagement roles get` - Retrieve roles

Returns all roles available in your organization. Use this to find the role ID you need for grant operations.

**Optional Flags:**
- `--name` - Filter by role name (partial match, case-insensitive)
- `--groupId` - Filter to roles currently granted to a specific group

**Examples:**

```bash
# Get all roles
newrelic accessmanagement roles get

# Filter by name (partial match)
newrelic accessmanagement roles get --name "Admin"

# Find roles already granted to a group
newrelic accessmanagement roles get --groupId <groupId>
```

**Response:**
```json
{
  "roles": [
    {
      "id": "role-456",
      "displayName": "All product admin",
      "name": "all_product_admin",
      "type": "ORGANIZATION",
      "scope": "organization",
      "organizationId": "org-abc123"
    },
    {
      "id": "role-789",
      "displayName": "Standard user",
      "name": "standard_user",
      "type": "ACCOUNT",
      "scope": "account",
      "organizationId": "org-abc123"
    }
  ]
}
```

---

### Permission Commands

#### `accessmanagement permissions get` - Retrieve permissions

Returns all permissions in your organization, optionally filtered by role or scope. Permissions represent individual capabilities that are bundled into roles.

**Optional Flags:**
- `--roleId` - Filter by role ID
- `--scope` - Filter by scope: `account` or `organization`

**Examples:**

```bash
# Get all permissions
newrelic accessmanagement permissions get

# Get permissions for a specific role
newrelic accessmanagement permissions get --roleId <roleId>

# Get only account-scoped permissions
newrelic accessmanagement permissions get --scope account

# Combine filters
newrelic accessmanagement permissions get --roleId <roleId> --scope account
```

**Response:**
```json
{
  "permissions": [
    {
      "category": "ALERTS",
      "feature": "alerts_conditions",
      "id": "perm-001",
      "scope": "account"
    },
    {
      "category": "ENTITY_EXPLORER",
      "feature": "workloads",
      "id": "perm-002",
      "scope": "account"
    }
  ]
}
```

---

## Working with JSON Responses

```bash
# Get a role ID by name
newrelic accessmanagement roles get --name "Standard user" \
  | jq -r '.roles[] | select(.displayName == "Standard user") | .id'

# List all role IDs and display names
newrelic accessmanagement roles get \
  | jq -r '.roles[] | "\(.id)\t\(.displayName)"'

# Get all grants for a group (group ID from usermanagement)
GROUP_ID=$(newrelic usermanagement groups get \
  --authDomainId <authDomainId> --name "Platform Engineers" \
  | jq -r '.authenticationDomains[].groups.groups[0].id')
newrelic accessmanagement grants get --groupId "$GROUP_ID"

# List permission categories for a role
newrelic accessmanagement permissions get --roleId <roleId> \
  | jq -r '[.permissions[].category] | unique[]'
```

---

## End-to-End Workflow

The access management workflow begins after you have created users and groups in [usermanagement](../usermanagement/README.md). This section picks up at step 5 of that workflow.

```bash
# Prerequisites: AUTH_DOMAIN, GROUP_ID, and USER_ID set from usermanagement steps
# See: usermanagement/README.md#end-to-end-workflow

# 5. Find the role you want to grant
ROLE_ID=$(newrelic accessmanagement roles get --name "Standard user" \
  | jq -r '.roles[] | select(.displayName == "Standard user") | .id')

# 6. Grant the group account-scoped access
newrelic accessmanagement grants create \
  --groupId "$GROUP_ID" \
  --roleId "$ROLE_ID" \
  --scope account \
  --accountId "$YOUR_ACCOUNT_ID"

# 7. Verify the grant was created
newrelic accessmanagement grants get --groupId "$GROUP_ID"

# To revoke later:
newrelic accessmanagement grants revoke \
  --groupId "$GROUP_ID" \
  --roleId "$ROLE_ID" \
  --scope account \
  --accountId "$YOUR_ACCOUNT_ID"
```

For steps 1–4 (creating users and groups), see [usermanagement — End-to-End Workflow](../usermanagement/README.md#end-to-end-workflow).

---

## Directory Structure

```
internal/accessmanagement/
├── README.md                     # This file
├── command.go                    # Root command and shared flag variables
├── command_grants.go             # Grant get/create/revoke commands
├── command_roles.go              # Role get command
├── command_permissions.go        # Permission get command
├── command_grants_test.go
├── command_roles_test.go
└── command_permissions_test.go
```

---

## Additional Resources

- [New Relic Access Management Documentation](https://docs.newrelic.com/docs/accounts/accounts-billing/new-relic-one-user-management/user-management-ui-and-tasks/#where)
- [Roles and Permissions Reference](https://docs.newrelic.com/docs/accounts/accounts-billing/new-relic-one-user-management/user-management-concepts/#roles)
- [User Management CLI](../usermanagement/README.md) — create the users and groups that receive the access grants you configure here
- [New Relic CLI Overview](https://github.com/newrelic/newrelic-cli)
