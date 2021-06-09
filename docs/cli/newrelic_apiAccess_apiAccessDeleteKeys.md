## newrelic apiAccess apiAccessDeleteKeys

A mutation to delete keys.

```
newrelic apiAccess apiAccessDeleteKeys [flags]
```

### Examples

```
newrelic apiAccess apiAccessDeleteKeys --keys
```

### Options

```
  -h, --help      help for apiAccessDeleteKeys
      --keys id   A list of each key id that you want to delete. You can read more about managing keys on [this documentation page](https://docs.newrelic.com/docs/apis/nerdgraph/examples/use-nerdgraph-manage-license-keys-personal-api-keys).
```

### Options inherited from parent commands

```
      --format string   output text format [JSON, Text, YAML] (default "JSON")
      --plain           output compact text
```

### SEE ALSO

* [newrelic apiAccess](newrelic_apiAccess.md)	 - Manage New Relic API access keys

