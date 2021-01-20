## newrelic apm application search

Search for a New Relic application

### Synopsis

Search for a New Relic application

The search command performs a query for an APM application name and/or account ID.


```
newrelic apm application search [flags]
```

### Examples

```
newrelic apm application search --name <appName>
```

### Options

```
  -h, --help          help for search
  -n, --name string   search for results matching the given APM application name
```

### Options inherited from parent commands

```
  -a, --accountId string    A New Relic account ID
      --applicationId int   A New Relic APM application ID
      --format string       output text format [JSON, Text, YAML] (default "JSON")
  -g, --guid string         search for results matching the given APM application GUID
      --plain               output compact text
```

### SEE ALSO

* [newrelic apm application](newrelic_apm_application.md)	 - Interact with New Relic APM applications

