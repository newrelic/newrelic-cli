## newrelic apm deployment delete

Delete a New Relic APM deployment

### Synopsis

Delete a New Relic APM deployment

The delete command performs a delete operation for an APM deployment.


```
newrelic apm deployment delete [flags]
```

### Examples

```
newrelic apm deployment delete --applicationId <appID> --deploymentID <deploymentID>
```

### Options

```
  -d, --deploymentID int   the ID of the deployment to be deleted
  -h, --help               help for delete
```

### Options inherited from parent commands

```
  -a, --accountId string    A New Relic account ID
      --applicationId int   A New Relic APM application ID
      --format string       output text format [YAML, JSON, Text] (default "JSON")
      --plain               output compact text
```

### SEE ALSO

* [newrelic apm deployment](newrelic_apm_deployment.md)	 - Manage New Relic APM deployment markers

