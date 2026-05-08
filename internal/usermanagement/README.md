# User Management CLI - Command Reference

A command-line interface for managing New Relic users, groups, and authentication domains. These commands correspond to the **User Management** section of the New Relic Administration UI.

User management covers the creation and lifecycle of users and groups within your organization's authentication domains. To manage what those groups can access — roles, grants, and permissions — see the [Access Management CLI](../accessmanagement/README.md).

## Table of Contents

- [Command Hierarchy](#command-hierarchy)
- [Prerequisites](#prerequisites)
- [Command Reference](#command-reference)
  - [User Commands](#user-commands)
  - [Group Commands](#group-commands)
  - [Group Membership Commands](#group-membership-commands)
  - [Authentication Domain Commands](#authentication-domain-commands)
- [Working with JSON Responses](#working-with-json-responses)
- [End-to-End Workflow](#end-to-end-workflow)
- [Directory Structure](#directory-structure)
- [Additional Resources](#additional-resources)

---

## Command Hierarchy

```
newrelic usermanagement
├── users
│   ├── get        # Retrieve users from an authentication domain
│   ├── create     # Create a new user
│   ├── update     # Update an existing user
│   └── delete     # Delete a user
│
├── groups
│   ├── get        # Retrieve groups and their members
│   ├── create     # Create a new group
│   ├── update     # Update a group's display name
│   ├── delete     # Delete a group
│   └── members
│       ├── add    # Add a user to a group
│       └── remove # Remove a user from a group
│
└── auth-domains
    └── get        # Retrieve authentication domains
```

---

## Prerequisites

### Required Environment Variables

```bash
export NEW_RELIC_API_KEY="NRAK-YOUR-API-KEY-HERE"
export NEW_RELIC_ACCOUNT_ID="your-account-id"
```

### Required Permissions

Your API key must belong to a user with the **Authentication domain manager** or **Organization manager** role. Read-only operations require at minimum the **Authentication domain read-only** capability.

### Finding Your Authentication Domain ID

Most commands that create or query users and groups require an authentication domain ID. Retrieve yours with:

```bash
newrelic usermanagement auth-domains get
```

The `id` field in the response is what you pass as `--authDomainId` to other commands.

---

## Command Reference

### User Commands

#### `usermanagement users get` - Retrieve users

Returns users from the specified authentication domain, with optional filters by user ID, email, or name.

**Required Flags:**
- `--authDomainId` - Authentication domain ID

**Optional Flags:**
- `--id` - Filter by user ID
- `--email` - Filter by email address (exact match)
- `--name` - Filter by name (exact match)

**Examples:**

```bash
# Get all users in a domain
newrelic usermanagement users get --authDomainId <authDomainId>

# Filter by email
newrelic usermanagement users get \
  --authDomainId <authDomainId> \
  --email user@example.com

# Filter by user ID
newrelic usermanagement users get \
  --authDomainId <authDomainId> \
  --id <userId>
```

**Response:**
```json
{
  "authenticationDomains": [
    {
      "id": "<authDomainId>",
      "name": "Default",
      "users": {
        "users": [
          {
            "id": "1234567",
            "name": "Jane Smith",
            "email": "jane@example.com",
            "type": { "displayName": "Full platform", "id": "2" },
            "lastActive": "2026-05-01T10:00:00Z",
            "emailVerificationState": "Verified"
          }
        ]
      }
    }
  ]
}
```

---

#### `usermanagement users create` - Create a user

Creates a new user in the specified authentication domain. The user receives an email invitation to set their password.

**Required Flags:**
- `--authDomainId` - Authentication domain ID
- `--email` - User's email address
- `--name` - User's full name

**Optional Flags:**
- `--userType` - User type: `BASIC_USER_TIER`, `CORE_USER_TIER`, or `FULL_USER_TIER`

**Examples:**

```bash
# Create a full platform user
newrelic usermanagement users create \
  --authDomainId <authDomainId> \
  --email jane@example.com \
  --name "Jane Smith" \
  --userType FULL_USER_TIER

# Create a basic user
newrelic usermanagement users create \
  --authDomainId <authDomainId> \
  --email dev@example.com \
  --name "Dev User" \
  --userType BASIC_USER_TIER
```

**Response:**
```json
{
  "createdUser": {
    "authenticationDomainId": "<authDomainId>",
    "email": "jane@example.com",
    "id": "7654321",
    "name": "Jane Smith",
    "type": { "displayName": "Full platform", "id": "2" }
  }
}
```

---

#### `usermanagement users update` - Update a user

Updates one or more fields on an existing user. Only the fields you provide are changed.

**Required Flags:**
- `--id` - User ID

**Optional Flags:**
- `--name` - New full name
- `--email` - New email address
- `--userType` - New user type: `BASIC_USER_TIER`, `CORE_USER_TIER`, or `FULL_USER_TIER`
- `--timeZone` - Timezone string (e.g., `America/Chicago`)

**Examples:**

```bash
# Upgrade user type
newrelic usermanagement users update \
  --id <userId> \
  --userType FULL_USER_TIER

# Update name and timezone
newrelic usermanagement users update \
  --id <userId> \
  --name "Jane Smith-Jones" \
  --timeZone America/Los_Angeles
```

**Response:**
```json
{
  "user": {
    "id": "7654321",
    "name": "Jane Smith-Jones",
    "email": "jane@example.com",
    "type": { "displayName": "Full platform", "id": "2" },
    "timeZone": "America/Los_Angeles"
  }
}
```

---

#### `usermanagement users delete` - Delete a user

Permanently removes a user from your organization. This action cannot be undone.

**Required Flags:**
- `--id` - User ID

**Examples:**

```bash
newrelic usermanagement users delete --id <userId>
```

**Response:** Prints `success` on completion.

---

### Group Commands

#### `usermanagement groups get` - Retrieve groups

Returns groups and their members from the specified authentication domain.

**Required Flags:**
- `--authDomainId` - Authentication domain ID

**Optional Flags:**
- `--id` - Filter by group ID
- `--name` - Filter by group display name (exact match)

**Examples:**

```bash
# Get all groups in a domain
newrelic usermanagement groups get --authDomainId <authDomainId>

# Filter by name
newrelic usermanagement groups get \
  --authDomainId <authDomainId> \
  --name "Developers"
```

**Response:**
```json
{
  "authenticationDomains": [
    {
      "id": "<authDomainId>",
      "name": "Default",
      "groups": {
        "groups": [
          {
            "id": "group-123",
            "displayName": "Developers",
            "users": {
              "users": [
                { "id": "7654321", "email": "jane@example.com", "name": "Jane Smith" }
              ]
            }
          }
        ]
      }
    }
  ]
}
```

---

#### `usermanagement groups create` - Create a group

Creates a new group in the specified authentication domain.

**Required Flags:**
- `--authDomainId` - Authentication domain ID
- `--name` - Display name for the group

**Examples:**

```bash
newrelic usermanagement groups create \
  --authDomainId <authDomainId> \
  --name "Developers"
```

**Response:**
```json
{
  "group": {
    "id": "group-123",
    "displayName": "Developers"
  }
}
```

---

#### `usermanagement groups update` - Update a group

Updates the display name of an existing group.

**Required Flags:**
- `--id` - Group ID
- `--name` - New display name

**Examples:**

```bash
newrelic usermanagement groups update \
  --id <groupId> \
  --name "Senior Developers"
```

**Response:**
```json
{
  "group": {
    "id": "group-123",
    "displayName": "Senior Developers"
  }
}
```

---

#### `usermanagement groups delete` - Delete a group

Permanently removes a group. Users in the group are not deleted, but they lose any access grants associated with the group. See [Access Management](../accessmanagement/README.md) for grant details.

**Required Flags:**
- `--id` - Group ID

**Examples:**

```bash
newrelic usermanagement groups delete --id <groupId>
```

**Response:** Prints `success` on completion.

---

### Group Membership Commands

#### `usermanagement groups members add` - Add a user to a group

Adds a user to a group. Users inherit any access grants assigned to the group — see [accessmanagement grants get](../accessmanagement/README.md#accessmanagement-grants-get---retrieve-access-grants) to view what a group has access to.

**Required Flags:**
- `--groupId` - Group ID
- `--userId` - User ID

**Examples:**

```bash
newrelic usermanagement groups members add \
  --groupId <groupId> \
  --userId <userId>
```

**Response:**
```json
{
  "groups": [
    {
      "id": "group-123",
      "displayName": "Platform Engineers",
      "users": {
        "users": [
          { "id": "7654321", "email": "jane@example.com", "name": "Jane Smith" }
        ]
      }
    }
  ]
}
```

---

#### `usermanagement groups members remove` - Remove a user from a group

Removes a user from a group. The user loses any access inherited from that group's grants.

**Required Flags:**
- `--groupId` - Group ID
- `--userId` - User ID

**Examples:**

```bash
newrelic usermanagement groups members remove \
  --groupId <groupId> \
  --userId <userId>
```

**Response:** Prints `success` on completion.

---

### Authentication Domain Commands

#### `usermanagement auth-domains get` - Retrieve authentication domains

Returns the authentication domains in your organization. Authentication domains define how users are provisioned and authenticated (SAML, SCIM, username/password, etc.).

**Optional Flags:**
- `--id` - Filter by authentication domain ID

**Examples:**

```bash
# Get all authentication domains
newrelic usermanagement auth-domains get

# Get a specific domain
newrelic usermanagement auth-domains get --id <authDomainId>
```

**Response:**
```json
{
  "authenticationDomains": [
    {
      "id": "<authDomainId>",
      "name": "Default",
      "provisioningType": "MANUAL"
    }
  ],
  "nextCursor": null,
  "totalCount": 1
}
```

---

## Working with JSON Responses

```bash
# Get all user IDs in a domain
newrelic usermanagement users get --authDomainId <id> \
  | jq -r '.authenticationDomains[].users.users[].id'

# Get a user ID by email
newrelic usermanagement users get \
  --authDomainId <id> \
  --email jane@example.com \
  | jq -r '.authenticationDomains[].users.users[0].id'

# Get all group IDs
newrelic usermanagement groups get --authDomainId <id> \
  | jq -r '.authenticationDomains[].groups.groups[] | "\(.id) \(.displayName)"'

# Get your authentication domain ID
newrelic usermanagement auth-domains get \
  | jq -r '.authenticationDomains[0].id'
```

---

## End-to-End Workflow

The typical workflow is: create users and groups here, then grant those groups access to accounts and roles using [accessmanagement](../accessmanagement/README.md).

```bash
# 1. Find your authentication domain
AUTH_DOMAIN=$(newrelic usermanagement auth-domains get \
  | jq -r '.authenticationDomains[0].id')

# 2. Create a group
GROUP_ID=$(newrelic usermanagement groups create \
  --authDomainId "$AUTH_DOMAIN" \
  --name "Platform Engineers" \
  | jq -r '.group.id')

# 3. Create a user
USER_ID=$(newrelic usermanagement users create \
  --authDomainId "$AUTH_DOMAIN" \
  --email engineer@example.com \
  --name "Alex Engineer" \
  --userType FULL_USER_TIER \
  | jq -r '.createdUser.id')

# 4. Add the user to the group
newrelic usermanagement groups members add \
  --groupId "$GROUP_ID" \
  --userId "$USER_ID"

# 5. Grant the group access to an account with a role
#    (continues in accessmanagement — see link below)
```

For step 5, see [accessmanagement grants create](../accessmanagement/README.md#accessmanagement-grants-create---create-an-access-grant).

---

## Directory Structure

```
internal/usermanagement/
├── README.md                    # This file
├── command.go                   # Root command and shared flag variables
├── command_users.go             # User CRUD commands
├── command_groups.go            # Group CRUD and membership commands
├── command_auth_domains.go      # Authentication domain commands
├── command_users_test.go
├── command_groups_test.go
└── command_auth_domains_test.go
```

---

## Additional Resources

- [New Relic User Management Documentation](https://docs.newrelic.com/docs/accounts/accounts-billing/new-relic-one-user-management/user-management-ui-and-tasks/)
- [Authentication Domains](https://docs.newrelic.com/docs/accounts/accounts-billing/new-relic-one-user-management/authentication-domains-saml-sso-scim-more/)
- [Access Management CLI](../accessmanagement/README.md) — assign roles and grants to the groups you create here
- [New Relic CLI Overview](https://github.com/newrelic/newrelic-cli)
