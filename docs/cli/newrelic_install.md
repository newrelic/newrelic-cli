## newrelic install

Install New Relic.

```
newrelic install [flags]
```

### Options

```
  -y, --assumeYes             use "yes" for all questions during install
  -h, --help                  help for install
      --localRecipes string   a path to local recipes to load instead of service other fetching
  -n, --recipe strings        the name of a recipe to install
  -c, --recipePath strings    the path to a recipe file to install
      --tag string            comma-separated list of tags ("key:value,key:value")
  -t, --testMode              fakes operations for UX testing
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

### Installing Specific Agent Versions

Install a specific version of the infrastructure agent by appending the version number to the recipe name with `@`. Currently, specific version installation is supported by the infrastructure agent only.

```
newrelic install -n infrastructure-agent-installer@1.65.0
```

Without a specified version, the command installs the latest available version. Supported on Linux and Windows hosts only; not available on macOS.

For a list of available versions, see the [Infrastructure agent release notes](https://docs.newrelic.com/docs/release-notes/infrastructure-release-notes/infrastructure-agent-release-notes/).

### SEE ALSO

- [newrelic](newrelic.md) - The New Relic CLI
