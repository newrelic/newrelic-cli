## newrelic apm deployment create

Create a New Relic APM deployment

### Synopsis

Create a New Relic APM deployment

The create command creates a new deployment marker for a New Relic APM
application.


```
newrelic apm deployment create [flags]
```

### Examples

```
newrelic apm deployment create --applicationId <appID> --revision <deploymentRevision>
```

### Options

```
      --change-log string    the change log stored with the deployment
      --description string   the description stored with the deployment
  -h, --help                 help for create
  -r, --revision string      a freeform string representing the revision of the deployment
      --user string          the user creating with the deployment
```

### Options inherited from parent commands

```
  -a, --accountId int       trace level logging
      --applicationId int   A New Relic APM application ID
      --debug               debug level logging
      --format string       output text format [JSON, Text, YAML] (default "JSON")
      --plain               output compact text
      --profile string      the authentication profile to use
      --trace               trace level logging
```

### SEE ALSO

* [newrelic apm deployment](newrelic_apm_deployment.md)	 - Manage New Relic APM deployment markers

