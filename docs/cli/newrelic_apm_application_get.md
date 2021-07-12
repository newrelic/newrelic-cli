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
  -a, --accountId int       trace level logging
      --applicationId int   A New Relic APM application ID
      --debug               debug level logging
      --format string       output text format [JSON, Text, YAML] (default "JSON")
  -g, --guid string         search for results matching the given APM application GUID
      --plain               output compact text
      --profile string      the authentication profile to use
      --trace               trace level logging
```

### SEE ALSO

* [newrelic apm application](newrelic_apm_application.md)	 - Interact with New Relic APM applications

