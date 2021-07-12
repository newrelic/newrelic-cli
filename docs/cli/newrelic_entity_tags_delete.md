## newrelic entity tags delete

Delete the given tag:value pairs from the given entity

### Synopsis

Delete the given tag:value pairs from the given entity

The delete command deletes all tags on the given entity 
that match the specified keys.


```
newrelic entity tags delete [flags]
```

### Examples

```
newrelic entity tags delete --guid <entityGUID> --tag tag1 --tag tag2 --tag tag3,tag4
```

### Options

```
  -g, --guid string   the entity GUID to delete tags on
  -h, --help          help for delete
  -t, --tag strings   the tag keys to delete from the entity
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

