## newrelic entity tags replace

Replace tag:value pairs for the given entity

### Synopsis

Replace tag:value pairs for the given entity

The replace command replaces any existing tag:value pairs with those
provided for the given entity.


```
newrelic entity tags replace [flags]
```

### Examples

```
newrelic entity tags replace --guid <entityGUID> --tag tag1:value1
```

### Options

```
  -g, --guid string   the entity GUID to replace tag values on
  -h, --help          help for replace
  -t, --tag strings   the tag names to replace on the entity
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

