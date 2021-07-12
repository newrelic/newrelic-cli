## newrelic apiAccess apiAccessCreateKeys

Create keys. You can create keys for multiple accounts at once. You can read more about managing keys on [this documentation page](https://docs.newrelic.com/docs/apis/nerdgraph/examples/use-nerdgraph-manage-license-keys-personal-api-keys).

```
newrelic apiAccess apiAccessCreateKeys [flags]
```

### Examples

```
newrelic apiAccess apiAccessCreateKeys --keys
```

### Options

```
  -h, --help          help for apiAccessCreateKeys
      --keys string   A list of the configurations for each key you want to create.
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

* [newrelic apiAccess](newrelic_apiAccess.md)	 - Manage New Relic API access keys

