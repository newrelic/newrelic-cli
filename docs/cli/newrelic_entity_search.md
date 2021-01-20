## newrelic entity search

Search for New Relic entities

### Synopsis

Search for New Relic entities

The search command performs a search for New Relic entities.


```
newrelic entity search [flags]
```

### Examples

```
newrelic entity search --name <applicationName>
```

### Options

```
  -a, --alert-severity string   search for entities matching the given alert severity type
  -d, --domain string           search for entities matching the given entity domain
  -f, --fields-filter strings   filter search results to only return certain fields for each search result
  -h, --help                    help for search
  -n, --name string             search for entities matching the given name
  -r, --reporting string        search for entities based on whether or not an entity is reporting (true or false)
      --tag string              search for entities matching the given entity tag
  -t, --type string             search for entities matching the given type
```

### Options inherited from parent commands

```
      --format string   output text format [Text, YAML, JSON] (default "JSON")
      --plain           output compact text
```

### SEE ALSO

* [newrelic entity](newrelic_entity.md)	 - Interact with New Relic entities

