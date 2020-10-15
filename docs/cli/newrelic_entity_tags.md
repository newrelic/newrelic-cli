## newrelic entity tags

Manage tags on New Relic entities

### Synopsis

Manage entity tags

The tag command allows users to manage the tags applied on the requested
entity. Use --help for more information.


### Examples

```
newrelic entity tags get --guid <guid>
```

### Options

```
  -h, --help   help for tags
```

### Options inherited from parent commands

```
      --format string   output text format [JSON, Text, YAML] (default "JSON")
      --plain           output compact text
```

### SEE ALSO

* [newrelic entity](newrelic_entity.md)	 - Interact with New Relic entities
* [newrelic entity tags create](newrelic_entity_tags_create.md)	 - Create tag:value pairs for the given entity
* [newrelic entity tags delete](newrelic_entity_tags_delete.md)	 - Delete the given tag:value pairs from the given entity
* [newrelic entity tags delete-values](newrelic_entity_tags_delete-values.md)	 - Delete the given tag/value pairs from the given entity
* [newrelic entity tags get](newrelic_entity_tags_get.md)	 - Get the tags for a given entity
* [newrelic entity tags replace](newrelic_entity_tags_replace.md)	 - Replace tag:value pairs for the given entity

