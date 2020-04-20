## newrelic nerdgraph query

Execute a raw GraphQL query request to the NerdGraph API

### Synopsis

Execute a raw GraphQL query request to the NerdGraph API

The query command accepts a single argument in the form of a GraphQL query as a string.
This command accepts an optional flag, --variables, which should be a JSON string where the
keys are the variables to be referenced in the GraphQL query.


```
newrelic nerdgraph query [flags]
```

### Examples

```
newrelic nerdgraph query 'query($guid: EntityGuid!) { actor { entity(guid: $guid) { guid name domain entityType } } }' --variables '{"guid": "<GUID>"}'
```

### Options

```
  -h, --help               help for query
      --variables string   the variables to pass to the GraphQL query, represented as a JSON string (default "{}")
```

### Options inherited from parent commands

```
      --format string   output text format [JSON, YAML] (default "JSON")
      --plain           output compact text
```

### SEE ALSO

* [newrelic nerdgraph](newrelic_nerdgraph.md)	 - Execute GraphQL requests to the NerdGraph API

