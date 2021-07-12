## newrelic entity tags delete-values

Delete the given tag/value pairs from the given entity

### Synopsis

Delete the given tag/value pairs from the given entity

The delete-values command deletes the specified tag:value pairs on a given entity.


```
newrelic entity tags delete-values [flags]
```

### Examples

```
newrelic entity tags delete-values --guid <guid> --tag tag1:value1
```

### Options

```
  -g, --guid string     the entity GUID to delete tag values on
  -h, --help            help for delete-values
  -v, --value strings   the tag key:value pairs to delete from the entity
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

