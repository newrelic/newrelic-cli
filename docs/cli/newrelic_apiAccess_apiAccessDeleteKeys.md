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
  -a, --accountId int    the account ID to use. Can be overridden by setting NEW_RELIC_ACCOUNT_ID
      --debug            debug level logging
      --format string    output text format [JSON, Text, YAML] (default "JSON")
      --plain            output compact text
      --profile string   the authentication profile to use
      --trace            trace level logging
```

### SEE ALSO

* [newrelic apiAccess](newrelic_apiAccess.md)	 - Manage New Relic API access keys

