## newrelic install

Install New Relic.

```
newrelic install [flags]
```

### Options

```
  -y, --assumeYes                  use "yes" for all questions during install
      --backup-location string     custom location for backup files (default: platform-specific)
  -h, --help                       help for install
      --list-backups               list all available configuration backups and exit
      --localRecipes string        a path to local recipes to load instead of service other fetching
  -n, --recipe strings             the name of a recipe to install
  -c, --recipePath strings         the path to a recipe file to install
      --restore-backup string      restore configuration from a specific backup ID (e.g. backup-2026-02-19-143022)
      --skip-backup                skip backing up existing configuration files before install
      --tag string                 comma-separated list of tags ("key:value,key:value")
  -t, --testMode                   fakes operations for UX testing
```

### Options inherited from parent commands

```
  -a, --accountId int    the account ID to use. Can be overridden by setting NEW_RELIC_ACCOUNT_ID
      --debug            debug level logging
      --format string    output text format [JSON, Text, YAML] (default "JSON")
      --plain            output compact text
      --profile string   the authentication profile to use
      --trace            trace level logging
```

### Installing Specific Versions of the Infrastructure Agent

To install a specific version of the infrastructure agent, append the version number to the recipe name using `@`. This functionality is currently supported only for the infrastructure agent.
```
newrelic install -n infrastructure-agent-installer@1.65.0
```

If no version is specified, the latest available version will be installed. This feature is supported on Linux and Windows hosts but is not available on macOS.

For a list of available versions, please refer to the [Infrastructure Agent Release Notes](https://docs.newrelic.com/docs/release-notes/infrastructure-release-notes/infrastructure-agent-release-notes/).

### Configuration Backup

Before any recipe is executed, the CLI automatically detects and backs up existing New Relic configuration files (Infrastructure agent, APM agents, Logging, Integrations). Backups are timestamped and stored in a platform-specific location:

| Platform         | Default backup location            |
|------------------|------------------------------------|
| Linux (root)     | `/opt/.newrelic-backups/`          |
| Linux (non-root) | `~/.newrelic-backups/`             |
| Windows          | `%ProgramData%\.newrelic-backups\` |
| macOS            | `~/.newrelic-backups/`             |

The last 5 backups are retained automatically; older ones are removed.

Each backup contains a `manifest.json` with the timestamp, list of files, SHA256 checksums, and CLI version used.

**Skip backup (CI/CD environments):**
```
newrelic install --skip-backup
```

**Use a custom backup directory:**
```
newrelic install --backup-location /var/backups/newrelic
```

**List all available backups:**
```
newrelic install --list-backups
```

**Restore configuration from a specific backup:**
```
newrelic install --restore-backup backup-2026-02-19-143022
```

Use `-y` / `--assumeYes` with `--restore-backup` to skip the confirmation prompt:
```
newrelic install --restore-backup backup-2026-02-19-143022 --assumeYes
```

### SEE ALSO

- [newrelic](newrelic.md) - The New Relic CLI
