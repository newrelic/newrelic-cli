## newrelic documentation

Generate CLI documentation

### Synopsis

Generate CLI documentation

newrelic documentation --outputDir <my directory> --type (markdown|manpage)



```
newrelic documentation [flags]
```

### Examples

```
newrelic documentation --outputDir /tmp
```

### Options

```
  -h, --help               help for documentation
  -o, --outputDir string   Output directory for generated documentation
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

* [newrelic](newrelic.md)	 - The New Relic CLI

