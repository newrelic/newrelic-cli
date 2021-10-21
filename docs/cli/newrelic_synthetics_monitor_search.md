## newrelic synthetics monitor search

Search for a New Relic Synthetics Monitor

### Synopsis

Search for a New Relic Synthetics Monitor

The search command performs a query for a Synthetics Monitor by name.


```
newrelic synthetics monitor search [flags]
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

* [newrelic synthetics monitor](newrelic_synthetics_monitor.md)	 - Interact with New Relic Synthetics monitors

