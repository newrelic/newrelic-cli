# NRQL Macros CLI - Command Reference

Commands for creating and managing NRQL Macros in your New Relic organization.

A NRQL Macro is a named, reusable expression fragment. Wherever you write `$macroName`
in a NRQL query, New Relic substitutes the macro's expression before executing the query.
Macros are available across every account in your organization.

## Table of Contents

- [Command Hierarchy](#command-hierarchy)
- [Prerequisites](#prerequisites)
- [Command Reference](#command-reference)
  - [`nrql macros create`](#nrql-macros-create)
  - [`nrql macros get`](#nrql-macros-get)
  - [`nrql macros list`](#nrql-macros-list)
  - [`nrql macros update`](#nrql-macros-update)
  - [`nrql macros delete`](#nrql-macros-delete)
- [Positional Parameters](#positional-parameters)
- [Output Format](#output-format)
- [Error Reference](#error-reference)

## Command Hierarchy

```
newrelic nrql macros
├── create          Create a new macro
├── get             Fetch a macro by entity GUID
├── list            List all macros in your organization
├── update          Update a macro's expression or parameter count
└── delete          Permanently delete a macro
```

## Prerequisites

**Required environment variables:**

```bash
export NEW_RELIC_API_KEY=<your-user-api-key>
export NEW_RELIC_ACCOUNT_ID=<your-account-id>
```

Macros are scoped to your **organization**, not a specific account. The organization is
resolved automatically from your API key — you do not need to pass an organization ID
for any of these commands.

## Command Reference

---

### `nrql macros create`

Creates a new NRQL Macro in your organization. The macro name is permanent and cannot
be changed after creation.

**Flags:**

| Flag | Short | Required | Description |
|------|-------|----------|-------------|
| `--name` | `-n` | Yes | Unique macro name; used as `$name` in queries. Immutable after creation. |
| `--expression` | `-e` | Yes | NRQL expression fragment the macro expands to. |
| `--parameter-count` | | No | Number of positional `$1`, `$2`, ... tokens in the expression. Required when the expression contains positional parameters. |
| `--org-id` | | No | Organization GUID. Resolved from your API key automatically if omitted. |

**Examples:**

```bash
# Simple filter macro — no parameters
newrelic nrql macros create \
  --name filterProduction \
  --expression "environment = 'production' AND appName NOT LIKE '%test%'"

# Using the short flags
newrelic nrql macros create \
  -n calculateApdex \
  -e "average(duration) / 1000"

# Parameterized macro — requires --parameter-count
newrelic nrql macros create \
  --name filterStatusRegion \
  --expression "status = $1 AND region = $2" \
  --parameter-count 2
```

**Example output:**

```json
{
  "expression": "environment = 'production' AND appName NOT LIKE '%test%'",
  "id": "NDgyOTY3M3xOR0...",
  "metadata": {
    "createdAt": 1793251957069,
    "createdBy": { "id": "1004672904", "type": "USER" },
    "updatedAt": 1793251957375,
    "updatedBy": { "id": "1004672904", "type": "USER" },
    "version": 1
  },
  "name": "filterProduction",
  "scope": { "id": "fb33fea3-...", "type": "ORGANIZATION" },
  "tags": null,
  "type": "NRQL_MACRO"
}
```

---

### `nrql macros get`

Fetches a single macro by its entity GUID.

**Flags:**

| Flag | Required | Description |
|------|----------|-------------|
| `--id` | Yes | Entity GUID of the macro to fetch. |

**Example:**

```bash
newrelic nrql macros get --id "NDgyOTY3M3xOR0..."
```

---

### `nrql macros list`

Lists all macros in your organization.

**Flags:** None (the list is automatically scoped to the organization associated with your API key).

**Example:**

```bash
newrelic nrql macros list
```

**Example output:**

```json
[
  {
    "expression": "environment = 'production' AND appName NOT LIKE '%test%'",
    "id": "NDgyOTY3M3xOR0...",
    "metadata": { "version": 1, ... },
    "name": "filterProduction",
    "scope": { "id": "fb33fea3-...", "type": "ORGANIZATION" },
    "type": "NRQL_MACRO"
  }
]
```

---

### `nrql macros update`

Updates the expression or parameter count of an existing macro. At least one of
`--expression` or `--parameter-count` must be provided.

> **Warning:** Expression changes take effect immediately. Every query, dashboard widget,
> and alert condition that references this macro will use the new expression on its next
> execution.

**Flags:**

| Flag | Short | Required | Description |
|------|-------|----------|-------------|
| `--id` | | Yes | Entity GUID of the macro to update. |
| `--expression` | `-e` | No | New NRQL expression fragment. Cannot be empty. |
| `--parameter-count` | | No | New positional parameter count. Pass `0` to remove all parameters. |

**Examples:**

```bash
# Update the expression
newrelic nrql macros update \
  --id "NDgyOTY3M3xOR0..." \
  --expression "environment = 'staging'"

# Remove positional parameters (parameter-count 2 → 0)
newrelic nrql macros update \
  --id "NDgyOTY3M3xOR0..." \
  --expression "environment = 'prod'" \
  --parameter-count 0

# Update both expression and parameter count together
newrelic nrql macros update \
  --id "NDgyOTY3M3xOR0..." \
  --expression "status = $1" \
  --parameter-count 1
```

---

### `nrql macros delete`

Permanently deletes a macro. After deletion, any query that uses `$macroName` will
no longer expand the macro. This cannot be undone.

**Flags:**

| Flag | Required | Description |
|------|----------|-------------|
| `--id` | Yes | Entity GUID of the macro to delete. |

**Example:**

```bash
newrelic nrql macros delete --id "NDgyOTY3M3xOR0..."
```

---

## Positional Parameters

NRQL Macros support positional parameters using bare `$1`, `$2`, ... tokens (without
quotes). When creating or updating a parameterized macro, `--parameter-count` is required
and must equal the exact number of distinct `$N` tokens in the expression.

```bash
# Expression uses $1 and $2 → parameter-count must be 2
newrelic nrql macros create \
  --name filterStatusRegion \
  --expression "status = $1 AND region = $2" \
  --parameter-count 2
```

Omitting `--parameter-count` when the expression contains `$N` tokens results in an
API error: `"Macros with parameters must be created/updated using a NerdGraph query
that provides the numParameters field."`

## Output Format

All commands (create, get, list, update) return a consistent flat entity structure.
The `delete` command prints an info log message rather than JSON.

Use global flags to control output format:

```bash
newrelic --format YAML nrql macros get --id "NDgyOTY3M3xOR0..."
newrelic --plain nrql macros list   # compact single-line JSON
```

## Error Reference

| Message | Cause |
|---------|-------|
| `required flag(s) "name" not set` | Missing `--name` on create |
| `required flag(s) "expression" not set` | Missing `--expression` on create |
| `required flag(s) "id" not set` | Missing `--id` on get/update/delete |
| `at least one of --expression or --parameter-count must be specified` | Called update with no changes |
| `--expression must not be empty` | Passed `--expression ""` |
| `<id>: Entity not found.` | The provided entity GUID does not exist |
| `Macros with parameters must be created/updated...` | Expression has `$N` tokens but `--parameter-count` was omitted |
| `...is not unique` | A macro with this name already exists in your organization |
