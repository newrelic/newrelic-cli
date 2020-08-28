## newrelic apm application get

Get a New Relic application

### Synopsis

Get a New Relic application

The get command performs a query for an APM application by GUID.


```
newrelic apm application get [flags]
```

### Examples

```
newrelic apm application get --guid <entityGUID>
```

### Options

```
  -h, --help   help for get
```

### Options inherited from parent commands

```
  -a, --accountId string    A New Relic account ID
      --applicationId int   A New Relic APM application ID
      --format string       output text format [Text, YAML, JSON] (default "JSON")
  -g, --guid string         search for results matching the given APM application GUID
      --plain               output compact text
```

### SEE ALSO

* [newrelic apm application](newrelic_apm_application.md)	 - Interact with New Relic APM applications

