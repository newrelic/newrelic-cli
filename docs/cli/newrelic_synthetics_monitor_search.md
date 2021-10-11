## newrelic apm application search

### Synopsis

Search for a New Relic application

The search command performs a query for an APM application name and/or account ID.

```
newrelic apm application search [flags]
```

### Examples

```
newrelic synthetics monitor search --name <monitorName>
```

### Options

```
  -h, --help          help for search
  -n, --name string   search for results matching the given Synthetics monitor name
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

* [newrelic apm application](newrelic_apm_application.md)	 - Interact with New Relic APM applications

