## newrelic profile delete

Delete a profile

### Synopsis

Delete a profile

The delete command removes the profile specified by name.


```
newrelic profile delete [flags]
```

### Examples

```
newrelic profile delete --profile <profile>
```

### Options

```
  -h, --help   help for delete
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

* [newrelic profile](newrelic_profile.md)	 - Manage the authentication profiles for this tool

