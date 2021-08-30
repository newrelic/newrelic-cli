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
  -a, --accountId int    the account ID to use. Can be overridden by setting NEW_RELIC_ACCOUNT_ID
      --debug            debug level logging
      --format string    output text format [JSON, Text, YAML] (default "JSON")
      --plain            output compact text
      --profile string   the authentication profile to use
      --trace            trace level logging
```

### SEE ALSO

* [newrelic apiAccess](newrelic_apiAccess.md)	 - Manage New Relic API access keys

