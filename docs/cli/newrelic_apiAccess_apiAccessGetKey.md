## newrelic apiAccess apiAccessGetKey

Fetch a single key by ID and type.

---
**NR Internal** | [#help-unified-api](https://newrelic.slack.com/archives/CBHJRSPSA) | visibility(customer)



```
newrelic apiAccess apiAccessGetKey [flags]
```

### Examples

```
newrelic apiAccess apiAccessGetKey --id --keyType
```

### Options

```
  -h, --help             help for apiAccessGetKey
      --id id            The id of the key. This can be used to identify a key without revealing the key itself (used to update and delete).
      --keyType string   The type of key.
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

