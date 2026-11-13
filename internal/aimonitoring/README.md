# AI Monitoring CLI - Command Reference

A command-line interface for querying New Relic AI Monitoring telemetry — the NRDB events reported by APM agents instrumenting LLM libraries (OpenAI, Bedrock, LangChain, etc.) and the applications that report them.

## Prerequisites

```bash
export NEW_RELIC_API_KEY="NRAK-YOUR-API-KEY-HERE"
export NEW_RELIC_ACCOUNT_ID="your-account-id"
```

## Commands

- **`newrelic aimonitoring events`** — look up your AI Monitoring data, like token usage, errors, or the actual prompts/responses sent to the LLM
  - `--type` — what kind of data you want: `summary` (default, token counts/cost), `message` (prompt/response text), or `embedding`, `feedback`, `tool`, `agent`, `vectorsearch`, `vectorsearchresult`
  - `--select`, `--where`, `--since`, `--limit` — pick which fields to show, filter the results, and set the time range
- **`newrelic aimonitoring application search`** — find out which of your applications are using AI/LLMs
  - `--tags` — only show apps with a matching tag (e.g. `aiEnabledApp:true`)
  - `--name` — look up one specific app by name
  - `--since` — how far back to look (default `7 days ago`)
- **`newrelic aimonitoring dataset search`** — list the AI Monitoring datasets you have access to
  - `--name` — only show datasets whose name contains this value (case-insensitive)
- **`newrelic aimonitoring dataset delete`** — permanently delete a dataset and all of its versions
  - `--id` — the dataset to delete (required)
  - `--yes` — skip the confirmation prompt
- **`newrelic aimonitoring dataset delete-version`** — permanently delete one version of a dataset
  - `--id`, `--version-id` — the dataset and the version to delete (both required)
  - `--yes` — skip the confirmation prompt
- **`newrelic aimonitoring dataset-version list`** — list every version of a dataset
  - `--dataset-id` — the dataset to list versions for (required, get it from `dataset search`)

## Command Reference

### `aimonitoring application search` - Find applications reporting AI Monitoring telemetry

Searches for entities that reported AI Monitoring events within a time window and resolves them to New Relic entities. Start broad with `--tags` to see all AI-related entities, then narrow further to a single application with `--name`. Output is trimmed to just the application `name` and `entityGuid` for each match.

**Flags:** `--since` (default `7 days ago`), `--tags` (optional, format `tagKey1:tagValue1,tagKey2:tagValue2` — entities must match every tag given), `--name` (optional, narrows to a single application name).

```bash
# see all AI-related entities reporting telemetry in the window
newrelic aimonitoring application search --since "7 days ago" --tags aiEnabledApp:true
```

**Response:**

```json
[
  {
    "name": "checkout-service",
    "entityGuid": "MTIzfEFQTXxBUFBMSUNBVElPTnwxMjM"
  }
]
```

```bash
# narrow further to a specific application by name
newrelic aimonitoring application search --tags aiEnabledApp:true --name checkout-service
```

> `--tags` filters the entities _resolved from_ the AI Monitoring NRDB query (tags live on the entity, not on the LLM events themselves). `--name` is applied as part of the underlying NRQL query, since `appName` is reported directly on the LLM events.

### `aimonitoring dataset search` - Find AI Monitoring datasets

Lists every `AIM_DATASET` entity the user has access to, following `nextCursor` so every page is returned. The entity search API only filters by type, so `--name` is applied client-side as a case-insensitive partial match: any dataset whose name contains the given value is returned.

**Flags:** `--name` (optional, omit to list every dataset).

```bash
newrelic aimonitoring dataset search --name "My Dataset"
```

**Response:**

```json
[
  {
    "id": "MTIzfEFJTXxEQVRBU0VUfDE",
    "name": "My Dataset",
    "description": "Customer support prompts",
    "variables": ["input"],
    "latestVersion": 2,
    "latestVersionId": "MTIzfEFJTXxEQVRBU0VUX1ZFUlNJT058Mg",
    "mode": "SINGLE_TURN",
    "tags": [{ "key": "env", "values": ["prod"] }]
  }
]
```

### `aimonitoring dataset-version list` - List the versions of a dataset

Lists every `AIM_DATASET_VERSION` entity whose parent `aimDataset` is the given dataset ID.

**Flags:** `--dataset-id` (required).

```bash
newrelic aimonitoring dataset-version list --dataset-id "MTIzfEFJTXxVRBU0VUfDE"
```

**Response:**

```json
[
  {
    "id": "MTIzfEFJTXxEQVRBU0VUX1NJT058MQ",
    "name": "My Dataset v1",
    "aimDataset": "MTIzfEFJTXxEQVRBU0VUfDE",
    "version": 1,
    "recordCount": 69,
    "variables": ["input"],
    "createdAt": 1793815928536,
    "updatedAt": 1793815928713
  }
]
```

### `aimonitoring dataset delete` / `delete-version` - Delete datasets and dataset versions

`delete` permanently deletes a dataset together with every one of its versions. `delete-version` permanently deletes a single version; **if it is the only version left, the dataset itself is deleted too**. Neither can be undone.

Both commands ask for confirmation first, defaulting to No. When no terminal is available (scripts, CI, AI assistants) the prompt cannot be answered, so the command exits with an error and nothing is deleted — pass `--yes` to delete without prompting.

```bash
newrelic aimonitoring dataset delete --id <datasetId>
newrelic aimonitoring dataset delete-version --id <datasetId> --version-id <versionId>

# non-interactive
newrelic aimonitoring dataset delete --id <datasetId> --yes
```

**Response:**

```json
{ "datasetId": "MTIzfEFJTXxEQVRBU0VUfDE" }
```

```json
{ "versionId": "MTIzfEFJTXxEQVRBU0VUX1ZFUlNJT058Mg" }
```

### `aimonitoring events` - Query AI Monitoring telemetry

Runs a NRQL `SELECT` against one of the NRDB event types reported for AI Monitoring, selected via `--type`:

| `--type`                       | NRDB event                 |
| ------------------------------ | -------------------------- |
| `summary` (default, **major**) | `LlmChatCompletionSummary` |
| `message` (**major**)          | `LlmChatCompletionMessage` |
| `embedding`                    | `LlmEmbedding`             |
| `feedback`                     | `LlmFeedbackMessage`       |
| `tool`                         | `LlmTool`                  |
| `agent`                        | `LlmAgent`                 |
| `vectorsearch`                 | `LlmVectorSearch`          |
| `vectorsearchresult`           | `LlmVectorSearchResult`    |

**Other flags:** `--select` (NRQL `SELECT` clause, default `*`), `--where` (optional NRQL `WHERE` clause), `--since` (default `1 day ago`), `--limit`.

```bash
newrelic aimonitoring events --type summary --select "count(*), average(duration)" --where "error is true" --since "1 day ago"
```

## More Query Examples

```bash
# token usage from LlmChatCompletionSummary
newrelic aimonitoring events --type summary --select "response.usage.completion_tokens, response.usage.prompt_tokens, response.usage.total_tokens"
```

```bash
# message content from LlmChatCompletionMessage
newrelic aimonitoring events --type message --select "content"
```

> If `content` comes back empty, it's likely because the APM agent has content recording turned off — many customers disable it by default to avoid capturing sensitive prompt/response text.

## Directory Structure

```
internal/aimonitoring/
├── README.md                        # This file
├── command.go                       # Root command
├── command_events.go                # events command
├── command_application.go           # application search command
├── command_dataset_search.go        # dataset search command
├── command_dataset_delete.go        # dataset delete command
├── command_dataset_delete_version.go # dataset delete-version command
├── helpers.go                       # shared helpers (confirmation prompt)
├── command_dataset_version_list.go  # dataset-version list command
├── entity_search.go                 # paginated EntityManagement entitySearch helper
├── command_test.go
├── command_events_test.go
├── command_application_test.go
├── command_dataset_search_test.go
├── command_dataset_delete_test.go
├── command_dataset_delete_version_test.go
└── command_dataset_version_list_test.go
```

## Additional Resources

- [Introduction to AI Monitoring](https://docs.newrelic.com/docs/ai-monitoring/intro-to-ai-monitoring/)
- [New Relic CLI Overview](https://github.com/newrelic/newrelic-cli)
