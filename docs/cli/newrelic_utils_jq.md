## newrelic utils jq

Parse json strings

### Synopsis

Parse json strings

The jq subcommand makes use of gojq (https://github.com/itchyny/gojq) to provide
json parsing capabilities.


```
newrelic utils jq [flags]
```

### Examples

```
echo '{"foo": 128}' | newrelic utils jq '.foo'
```

### Options

```
  -h, --help   help for jq
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

* [newrelic utils](newrelic_utils.md)	 - Various utility methods

