## newrelic entity tags get

Get the tags for a given entity

### Synopsis

Get the tags for a given entity

The get command returns JSON output of the tags for the requested entity.


```
newrelic entity tags get [flags]
```

### Examples

```
newrelic entity tags get --guid <entityGUID>
```

### Options

```
  -g, --guid string   the entity GUID to retrieve tags for
  -h, --help          help for get
```

### Options inherited from parent commands

```
  -a, --accountId int    trace level logging
      --debug            debug level logging
      --format string    output text format [JSON, Text, YAML] (default "JSON")
      --plain            output compact text
      --profile string   the authentication profile to use
      --trace            trace level logging
```

### SEE ALSO

* [newrelic entity tags](newrelic_entity_tags.md)	 - Manage tags on New Relic entities

