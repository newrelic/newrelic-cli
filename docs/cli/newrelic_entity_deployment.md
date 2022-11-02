## newrelic entity deployment

Manage New Relic Entity changes

### Synopsis

Manage New Relic Entity changes

A deployment marker is an event indicating that a deployment happened, and
it's paired with metadata available from your SCM system (for example,
the user, revision, or change-log). New Relic displays a vertical line, or
“marker,” on charts and graphs at the deployment event's timestamp.


### Examples

```
newrelic entity deployment create --guid <guid> --version <1.0.0>
```

### Options

```
  -h, --help   help for deployment
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

* [newrelic entity](newrelic_entity.md)	 - Interact with New Relic entities
* [newrelic entity deployment create](newrelic_entity_deployment_create.md)	 - Create a deployment marker for a New Relic entity

