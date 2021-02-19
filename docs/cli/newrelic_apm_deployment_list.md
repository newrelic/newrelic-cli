## newrelic apm deployment list

List New Relic APM deployments for an application

### Synopsis

List New Relic APM deployments for an application

The list command returns deployments for a New Relic APM application.


```
newrelic apm deployment list [flags]
```

### Examples

```
newrelic apm deployment list --applicationId <appID>
```

### Options

```
  -h, --help   help for list
```

### Options inherited from parent commands

```
  -a, --accountId string    A New Relic account ID
      --applicationId int   A New Relic APM application ID
      --format string       output text format [Text, YAML, JSON] (default "JSON")
      --plain               output compact text
```

### SEE ALSO

* [newrelic apm deployment](newrelic_apm_deployment.md)	 - Manage New Relic APM deployment markers

