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

### SEE ALSO

- [newrelic](newrelic.md) - The New Relic CLI
