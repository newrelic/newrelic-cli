## newrelic entity changetracker

Manage New Relic Entity changes

### Synopsis

Manage New Relic Entity changes

A change tracker is an event indicating that a change happened, and
it's paired with metadata available from your SCM system (for example,
the user, revision, or change-log). New Relic displays a vertical line, or
“marker,” on charts and graphs at the change event's timestamp.


### Examples

```
newrelic entity changetracker create --guid <guid> --version <1.0.0>
```

### Options

```
  -h, --help   help for changetracker
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
* [newrelic entity changetracker create](newrelic_entity_changetracker_create.md)	 - Create a change for a New Relic entity

