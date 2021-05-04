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
      --format string   output text format [JSON, Text, YAML] (default "JSON")
      --plain           output compact text
```

### SEE ALSO

* [newrelic utils](newrelic_utils.md)	 - Various utility methods

