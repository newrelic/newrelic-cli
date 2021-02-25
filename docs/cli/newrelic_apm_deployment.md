## newrelic apm deployment

Manage New Relic APM deployment markers

### Synopsis

Manage New Relic APM deployment markers

A deployment marker is an event indicating that a deployment happened, and
it's paired with metadata available from your SCM system (for example,
the user, revision, or change-log). APM displays a vertical line, or
“marker,” on charts and graphs at the deployment event's timestamp.


### Examples

```
newrelic apm deployment list --applicationId <appID>
```

### Options

```
  -h, --help   help for deployment
```

### Options inherited from parent commands

```
  -a, --accountId string    A New Relic account ID
      --applicationId int   A New Relic APM application ID
      --format string       output text format [YAML, JSON, Text] (default "JSON")
      --plain               output compact text
```

### SEE ALSO

* [newrelic apm](newrelic_apm.md)	 - Interact with New Relic APM
* [newrelic apm deployment create](newrelic_apm_deployment_create.md)	 - Create a New Relic APM deployment
* [newrelic apm deployment delete](newrelic_apm_deployment_delete.md)	 - Delete a New Relic APM deployment
* [newrelic apm deployment list](newrelic_apm_deployment_list.md)	 - List New Relic APM deployments for an application

