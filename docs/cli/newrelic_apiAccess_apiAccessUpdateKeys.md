## newrelic apiAccess apiAccessUpdateKeys

Update keys. You can update keys for multiple accounts at once. You can read more about managing keys on [this documentation page](https://docs.newrelic.com/docs/apis/nerdgraph/examples/use-nerdgraph-manage-license-keys-personal-api-keys).

```
newrelic apiAccess apiAccessUpdateKeys [flags]
```

### Examples

```
newrelic apiAccess apiAccessUpdateKeys --keys
```

### Options

```
  -h, --help          help for apiAccessUpdateKeys
      --keys string   The configurations of each key you want to update.
```

### Options inherited from parent commands

```
      --format string   output text format [JSON, Text, YAML] (default "JSON")
      --plain           output compact text
```

### SEE ALSO

* [newrelic apiAccess](newrelic_apiAccess.md)	 - Manage New Relic API access keys

